# SATySFi-Online
(Work in Progress)

[SATySFi](https://github.com/gfngfn/SATySFi)がOnlineで使えたら→→うれしい！！

## Screenshot

![](docs/screenshot.png)

## Requirements

- go 1.13
- Elm 0.19.1
- Docker >= 19.03

## DevDependencies

- GNU Make
- yarn
- [cespare/reflex](https://github.com/cespare/reflex)

## How to use

(with reflex)

1. make satysfi
2. make dev
3. access `http://loacalhost:8888`

(without reflex)
1. make satysfi
2. make build
3. PORT=8888 make run
4. access `http://loacalhost:8888`


## Directory Structure

```
.
├── config.json             # 設定ファイル。いる
├── docs
│   └── screenshot.png
├── go.mod
├── go.sum
├── LICENSE
├── main.go                 # 本体。ファイルの配信とかコンパイルとかやる
├── Makefile
├── README.md
├── reflex.conf
├── SATySFi                 # SATySFiが使えるDockerfileとおまけのdocker-compose
│   ├── docker-compose.yml
│   └── Dockerfile
├── template                # プロジェクトのデフォルトの構成。名前はconfig.jsonに依る
│   ├── assets
│   ├── demo.saty
│   └── local.satyh
├── ui                      # フロント。Elm
│   ├── elm.json
│   ├── elm-stuff
│   ├── index.html
│   ├── index.js
│   ├── main.css
│   ├── package.json
│   ├── src
│   └── yarn.lock
└── work                    # Project保管ばしょ。config.jsonに書いてある

```


## Author

theoremoon

## License

Apache 2.0 (SATySFiはLGPLだけど大丈夫かな……？)
