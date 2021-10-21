
PROJECT=chetan/go-go-github-badge

.DEFAULT_GOAL := run

all:
	docker buildx build --platform linux/amd64,linux/arm64 --pull --push -t ${PROJECT}:latest .

build: build-amd64

build-amd64:
	docker buildx build --platform linux/amd64 --pull --push -t ${PROJECT}:latest .

build-arm64:
	docker buildx build --platform linux/arm64/v8 --pull --push -t ${PROJECT}:latest .

push:
	docker push ${PROJECT}:latest

run:
	go run ./bin/go-go-github-badge

serve: run
