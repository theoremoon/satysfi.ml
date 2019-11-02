.PHONY: help
help:
	@cat Makefile | grep -E "^[A-Za-z0-9-]+:"

.PHONY: dev
dev:
	PORT=8888 reflex -g reflex.conf -s -- reflex -c reflex.conf

.PHONY: build
build: ui go

.PHONY: run
run:
	./app

.PHONY: go
go: main.go
	go generate
	go build -o app

.PHONY: ui
ui:
	cd ui; \
	yarn;\
	yarn run parcel build index.html;
	rm -rf dist
	mv ui/dist ./dist

.PHONY: satysfi
satysfi:
	cd SATySFi; \
	docker build -t satysfi .
