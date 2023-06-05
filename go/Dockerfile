FROM golang:1.20-buster AS dev

WORKDIR /opt

RUN set -eux \
 && apt-get update -yqq \
 && apt-get install -y --no-install-recommends \
        curl \
        jq \
        make \
 # go \
 && go install github.com/kisielk/errcheck@latest \
 && go install github.com/kyleconroy/sqlc/cmd/sqlc@latest \
 && go install golang.org/x/tools/cmd/goimports@latest \
 && go install honnef.co/go/tools/cmd/staticcheck@latest \
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
