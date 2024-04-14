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
	SessionTokenName    string   `env:"SESSION_TOKEN_NAME"`
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

func (app *App) ErrorHandler(w http.ResponseWriter, req *http.Request, err error) {
	// TODO ErrorHandler
	slog.Error("ErrorHandler", "ctx", req.Context())
}

func (app *App) HttpHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func (app *App) HttpCallback(w http.ResponseWriter, r *http.Request) {
	oa := auth.NewGoogleOAuth(app.ResolveRedirectURL(r), app.OAuthClientId, app.OAuthClientSecret)

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

	userInfo, err := oa.GetUserInfo(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("InternalServerError"))
		return
	}
	slog.Info("HttpCallback: Authorized.", "user", userInfo)

	// TODO どっかに保存しておく
	http.SetCookie(w, &http.Cookie{
		Name:     app.SessionTokenName,
		Value:    token.AccessToken,
		Path:     "/",
		Expires:  token.Expiry,
		Secure:   false, // TODO
		HttpOnly: true,
	})

	// TODO 最初のURL覚えといてそっちに飛ばす
	http.Redirect(w, r, PATH_PREFIX+"/success", http.StatusFound)
}

func (app *App) Rewrite(r *httputil.ProxyRequest) {
	// slog.Info("xxxx",
	// 	"Host", r.In.Host,
	// 	"URL", r.In.URL,
	// 	"Proto", r.In.Proto,
	// 	"Header", r.In.Header,
	// 	"Host", r.In.Host,
	// 	"RemoteAddr", r.In.RemoteAddr,
	// 	"RequestURI", r.In.RequestURI,
	// )
	target, err := url.Parse("http://" + r.In.Host)
	if err != nil {
		panic(err)
	}
	r.SetURL(target)
	r.SetXForwarded()
}

func (app *App) ResolveRedirectURL(r *http.Request) string {
	// TODO https, localhost or
	// redirectURL := "http://" + r.Host + PATH_PREFIX + "/callback"
	redirectURL := "http://localhost:9999" + PATH_PREFIX + "/callback"
	return redirectURL
}

func (app *App) ReverseProxyFunc(rp *httputil.ReverseProxy) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie(app.SessionTokenName)
		if err != nil {
			oa := auth.NewGoogleOAuth(app.ResolveRedirectURL(r), app.OAuthClientId, app.OAuthClientSecret)
			http.Redirect(w, r, oa.GetAuthCodeURL(app.CsrfToken), http.StatusFound)
			return
		}
		// TODO どっかに保存してある
		slog.Info("ReverseProxyFunc: Token.", "value", token)
		rp.ServeHTTP(w, r)
	}
}

func main() {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	slog.SetDefault(slog.New(logHandler))

	app, err := NewApp()
	if err != nil {
		panic(err)
	}

	// TODO リクエストログどこで出す
	http.HandleFunc(PATH_PREFIX+"/ping", app.HttpHealthCheck)
	http.HandleFunc(PATH_PREFIX+"/callback", app.HttpCallback)
	http.HandleFunc("/", app.ReverseProxyFunc(&httputil.ReverseProxy{
		// TODO Transport
		Rewrite:       app.Rewrite,
		FlushInterval: app.TargetFlushInterval,
		ErrorHandler:  app.ErrorHandler,
	}))

	slog.Info("Server stated.", "listen", app.Listen)
	err = http.ListenAndServe(app.Listen, nil)
	if err != nil {
		panic(err)
	}
}
