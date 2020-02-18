.DEFAULT_GOAL := run

.PHONY: build
build: gen
	go build ./cmd/koverto

.PHONY: docker
docker: build
	docker build . -t koverto/koverto:latest

.PHONY: gen
gen:
	go generate ./api

.PHONY: run
run: gen
	go run ./cmd/koverto
