FROM golang:1.13-alpine as builder

RUN apk update
RUN apk add --no-cache curl gzip git make

# elm
RUN curl -sfSL --retry 3 https://github.com/elm/compiler/releases/download/0.19.1/binary-for-linux-64-bit.gz -o elm.gz && gunzip elm.gz && mv elm /usr/bin/elm && chmod a+x /usr/bin/elm

WORKDIR /go/app
COPY . .
RUN make build

FROM alpine

WORKDIR /app
COPY --from=builder /go/app/app /app/app
RUN addgroup go \
        && adduser -D -G go go \
        && chown -R go:go /app
CMD ./app


