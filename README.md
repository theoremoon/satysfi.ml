# SATySFi-Online
(Work in Progress)

SATySFiがOnlineで使えたら→→うれしい！！

## Requirements

- go 1.13
- Elm 0.19.1

## Dependencies

- GNU Make
- Docker >= 19.03
- yarn
- [cespare/reflex](https://github.com/cespare/reflex)

## How to use

(with reflex)

1. make build-satysfi
2. make dev
3. access `http://loacalhost:8888`

(withou reflex)
1. make build-satysfi
2. make build
3. PORT=8888 make run
4. access `http://loacalhost:8888`


## Directory Structure

```
.
├── docker-compose.yml        調整中
├── Dockerfile              　調整中
├── go.mod
├── go.sum
├── main.go                   サーバ。labstack/echo製。htmlの配信とSATySFiのコンパイルを担ってる
├── Makefile                  便利
├── README.md
├── reflex.conf
├── SATySFi                   SATySFi用のDocker Imageを作るやつ。docker-composeはおまけ
│   ├── docker-compose.yml
│   └── Dockerfile
└── web                       フロント。Elm製
    ├── elm.json
    ├── index.html
    ├── index.js
    ├── main.css
    ├── package.json          Elm をビルドしてCSSをbundleするためにparcelを使ってる
    ├── src
    │   └── Main.elm          フロント本体。Elm楽しい
    └── yarn.lock
```


## Author

theoremoon

## License

Apache 2.0