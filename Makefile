
PROJECT=chetan/go-go-github-badge

.DEFAULT_GOAL := run

build-docker:
	docker buildx build --platform linux/amd64,linux/arm64 --pull --push -t ${PROJECT}:latest .

build:
	go build ./bin/go-go-github-badge

run:
	go run ./bin/go-go-github-badge

serve: run
