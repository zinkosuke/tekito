package main

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/caarlos0/env/v10"

	"proxy/pkg/auth"
)

const PATH_PREFIX = "/hec4XUHvwm"

type App struct {
	OAuthClientId       string   `env:"OAUTH_CLIENT_ID"`
	OAuthClientSecret   string   `env:"OAUTH_CLIENT_SECRET"`
	CsrfToken           string   `env:"CSRF_TOKEN"`
	TargetHosts         []string `env:"TARGET_HOSTS" envSeparator:":"`
	TargetFlushInterval time.Duration
	Listen              string
}

func NewApp() (*App, error) {
	flushInterval, _ := time.ParseDuration("3s")
	app := &App{
		TargetFlushInterval: flushInterval,
		Listen:              ":80",
	}
	if err := env.Parse(app); err != nil {
		slog.Error("NewApp: Application initialization failed.", err)
		return nil, err
	}
	return app, nil
}

func (app *App) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func (app *App) Redirect(w http.ResponseWriter, r *http.Request) {
	// TODO ALB配下でhttps判定
	redirectURL := "http://" + r.Host + PATH_PREFIX + "/callback"
	oa := auth.NewGoogleOAuth(redirectURL, app.OAuthClientId, app.OAuthClientSecret)
	http.Redirect(w, r, oa.GetAuthCodeURL(app.CsrfToken), http.StatusFound)
}

func (app *App) Callback(w http.ResponseWriter, r *http.Request) {
	redirectURL := "http://" + r.Host + PATH_PREFIX + "/callback"
	oa := auth.NewGoogleOAuth(redirectURL, app.OAuthClientId, app.OAuthClientSecret)

	state := r.URL.Query().Get("state")
	if state != app.CsrfToken {
		slog.Warn("Callback: Invalid state.", "raw_query", r.URL.RawQuery)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("BadRequest 'state'"))
		return
	}

	code := r.URL.Query().Get("code")
	token, err := oa.GetAccessToken(code)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("BadRequest 'code'"))
		return
	}
	// TODO token はどっかに保存 & cookieに直接セットしない

	userInfo, err := oa.GetUserInfo(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("InternalServerError"))
		return
	}
	slog.Info("TODO GetUserInfo: getuser", "userInfo", userInfo)

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token.AccessToken,
		Secure:   false, // TODO
		HttpOnly: true,
	})
	slog.Info("TODO expire", "x", token.Expiry)

	// TODO IndexのURL覚えといてそっちに飛ばす
	http.Redirect(w, r, PATH_PREFIX+"/success", http.StatusFound)
}

func (app *App) Success(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("access_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	slog.Info("TODO token", "xxxxxxxxxxx", token)
	w.Write([]byte("Success"))
}

func (app *App) Rewrite(r *httputil.ProxyRequest) {
    // ALBで設定できるもの
    // ユーザー認証を設定
    // セッション cookie 名
    // セッションタイムアウト
    // スコープ
    // 認証されていないリクエストのアクション
    // 追加のリクエストパラメータ
    //   ID プロバイダーのパラメータ
	slog.Info("TODO xxxx",
		"Host", r.In.Host,
		"URL", r.In.URL,
		"Proto", r.In.Proto,
		"Header", r.In.Header,
		"Host", r.In.Host,
		"RemoteAddr", r.In.RemoteAddr,
		"RequestURI", r.In.RequestURI,
	)
	target, err := url.Parse("http://" + r.In.Host)
	if err != nil {
		panic(err)
	}
	r.SetURL(target)
	r.Out.Host = r.In.Host // if desired
	r.SetXForwarded()
}

func (app *App) ErrorHandler(w http.ResponseWriter, req *http.Request, err error) {
	slog.Error("ErrorHandler", "ctx", req.Context())
}

func main() {
	app, err := NewApp()
	if err != nil {
		panic(err)
	}

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	slog.SetDefault(slog.New(logHandler))

	// ReverseProxy
	rp := &httputil.ReverseProxy{
		// TODO Transport
		Rewrite:       app.Rewrite,
		FlushInterval: app.TargetFlushInterval,
		ErrorHandler:  app.ErrorHandler,
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rp.ServeHTTP(w, r)
	})

	// OAuth
	http.HandleFunc(PATH_PREFIX+"/ping", app.HealthCheck)
	http.HandleFunc(PATH_PREFIX+"/redirect", app.Redirect)
	http.HandleFunc(PATH_PREFIX+"/success", app.Success)
	http.HandleFunc(PATH_PREFIX+"/callback", app.Callback)

	slog.Info("Server stated.", "listen", app.Listen)
	err = http.ListenAndServe(app.Listen, nil)
	if err != nil {
		panic(err)
	}
}
