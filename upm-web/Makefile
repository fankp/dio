.PHONY: build
build:
	go build -o upm-web main.go plugin.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t user-web:latest
