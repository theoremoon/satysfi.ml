.PHONY: all
all:
	@cat Makefile | grep -E "^[A-Za-z0-9-]+:"
.PHONY: dev
dev:
	PORT=8888 reflex -g reflex.conf -s -- reflex -c reflex.conf

.PHONY: release
release:
	docker-compose up --build

.PHONY: build
build: build-web build-go

.PHONY: run
run:
	./app

.PHONY: build-go
build-go: main.go
	go build -o app

.PHONY: build-web
build-web: web/src/*.elm
	cd web; \
	yarn;\
	yarn run parcel build index.html;

.PHONY: build-satysfi
build-satysfi:
	cd SATySFi; \
	docker build -t satysfi .
