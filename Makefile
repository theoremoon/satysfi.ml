.PHONY: help
help:
	@cat Makefile | grep -E "^[A-Za-z0-9-]+:"

.PHONY: dev
dev:
	docker-compose down
	docker-compose up -d db
	DSN="postgres://127.0.0.1:54321/satysfi?user=satysfi&password=satysfi&sslmode=disable" PORT=8888 reflex -g reflex.conf -s -- reflex -c reflex.conf

.PHONY: run
run:
	./app

.PHONY: build
build: ui go

.PHONY: go
go:
	go build -o app

.PHONY: ui
ui:
	cd ui;\
		yarn;\
		yarn run parcel build index.html;

.PHONY: satysfi
satysfi:
	cd SATySFi; \
	docker build -t satysfi .
