.PHONY: help
help:
	@cat Makefile | grep -E "^[A-Za-z0-9-]+:"

.PHONY: dev
dev:
	parallel --line-buffer bash -c ::: 'make watch-ui' 'reflex -g "dist/*" -s -- sh -c "sleep 5; make go && PORT=8888 make run"'

.PHONY: build
build: ui go

.PHONY: run
run:
	./app

.PHONY: go
go: main.go
	go generate
	go build -o app

.PHONY: watch-ui
watch-ui:
	cd ui;\
		yarn run parcel watch index.html -d ../dist;

.PHONY: ui
ui:
	cd ui; \
	yarn;\
	yarn run parcel build index.html -d ../dist;

.PHONY: setup
setup:
	go get github.com/rakyll/statik

