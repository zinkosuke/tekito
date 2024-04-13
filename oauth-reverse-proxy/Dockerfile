FROM golang:1.22-bookworm AS dev

WORKDIR /opt

RUN set -eux \
 && apt-get update -yqq \
 && apt-get install -y --no-install-recommends \
        curl \
        jq \
        make \
 # go-ci \
 && go install github.com/kisielk/errcheck@latest \
 && go install golang.org/x/tools/cmd/goimports@latest \
 && go install honnef.co/go/tools/cmd/staticcheck@latest \
 # go-api \
 && go install github.com/cosmtrek/air@latest \
 # clean \
 && apt-get autoremove -yqq --purge \
 && apt-get -y clean \
 && rm -rf \
        /tmp/* \
        /usr/share/doc \
        /usr/share/doc-base \
        /usr/share/man \
        /var/lib/apt/lists/* \
        /var/tmp/* \
 && true

COPY . .
