.PHONY: all
all:
	@cat Makefile | grep -E "^[A-Za-z0-9-]+:"
# don't use in the docker because this command use docker-compose
.PHONY: dev
dev:
	cd dev; docker-compose up --build

.PHONY: release
release:
	docker-compose up --build

# used in docker
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
	elm make src/Main.elm --output=index.html


