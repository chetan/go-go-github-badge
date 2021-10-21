# BUILDER image
FROM golang:alpine AS builder

COPY . /app

RUN cd /app && go build ./bin/go-go-github-badge


# RUNTIME image
FROM alpine AS runtime

# personal access token for making github v4 api requests
ENV GITHUB_TOKEN=""

# restrict the users for which we can return badges
ENV ALLOWED_USERS=""

ENV GIN_MODE=release

RUN mkdir /app

COPY --from=builder /app/go-go-github-badge /app/
COPY ./static /app/static
COPY ./templates /app/templates

WORKDIR /app
ENTRYPOINT ["/app/go-go-github-badge"]
