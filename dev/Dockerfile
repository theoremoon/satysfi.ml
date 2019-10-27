# for developing
FROM golang:1.13-alpine as dev

RUN apk update
RUN apk add --no-cache curl gzip git make

# elm
RUN curl -sfSL --retry 3 https://github.com/elm/compiler/releases/download/0.19.1/binary-for-linux-64-bit.gz -o elm.gz && gunzip elm.gz && mv elm /usr/bin/elm && chmod a+x /usr/bin/elm

WORKDIR /develop

# reflex
RUN go get github.com/cespare/reflex

