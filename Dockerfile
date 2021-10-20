# BUILDER image
FROM golang:alpine AS builder

COPY . /app

RUN cd /app && go get ./bin/...


# RUNTIME image
FROM alpine AS runtime

ENV GIN_MODE=release

RUN mkdir /app

COPY --from=builder /go/bin/go-go-github-badge /app/
COPY ./static /app/static
COPY ./templates /app/templates

WORKDIR /app
ENTRYPOINT ["/app/go-go-github-badge"]
