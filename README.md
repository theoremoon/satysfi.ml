# SATySFi-Online
(Work in Progress)

[SATySFi](https://github.com/gfngfn/SATySFi)がOnlineで使えたら→→うれしい！！

## Screenshot

![](docs/screenshot.png)

## Requirements

- go 1.13
- Vue 2.6.10
- Docker >= 19.03

## DevDependencies

- GNU Make
- GNU Parallel
- yarn
- [cespare/reflex](https://github.com/cespare/reflex)
- [statik](https://github.com/rakyll/statik)

## How to use

(with GNU parallel and reflex)

1. make dev
2. access `http://localhost:8888`

(without reflex)
1. make setup
2. make build
3. make run
4. access `http://localhost:8888`

**NOTE**: This tool uses a Docker image for PDF generation. When starting the application local environment, execute the following command before compiling the SATySFi document.

```
$ docker pull theoldmoon0602/satysfi:latest
```

## Directory Structure

```
.
├── app                           # 成果物。gitignoreされてる
├── dist                          # ビルドされたフロントが個々にできる。gitignoreされてる
├── docs
│   └── screenshot.png
├── go.mod
├── go.sum
├── LICENSE
├── main.go                       # ファイルの配信とAPIを捌くのをあｙる
├── Makefile
├── README.md
├── statik                        # dist/以下をstatikでまとめてる。gitignoreされてる
│   └── statik.go
├── template                      # プロジェクトのテンプレート
│   ├── assets
│   ├── demo.saty
│   └── local.satyh
├── ui                            # Vue。なんやかんやでElmやめちゃった
│   ├── favicon.png
│   ├── index.html
│   ├── index.js
│   ├── node_modules
│   ├── package.json
│   ├── src
│   └── yarn.lock
└── work                          # 作業ファイルがここに作られる。gitignoreされてる
```

## Want to DO

- [ ] 実ディレクトリをいじる代わりにDBをいじるようにしたい（なんかアプリケーションを分散させられる気がするので）
- [ ] エディタをまともにしたい（シンタックスハイライト・補完……）
- [ ] いい感じビルド＆デプロイ機構がほしい

## Author

theoremoon

## License

Apache 2.0 (SATySFiはLGPLだけど大丈夫かな……？)
