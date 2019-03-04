APP_NAME := isayme/go-wstunnel
APP_VERSION := $(shell git describe --tags --always)
APP_PKG := $(shell echo ${PWD} | sed -e "s\#${GOPATH}/src/\#\#g")
BUILD_TIME := $(shell date -u +"%FT%TZ")
GIT_REVISION := $(shell git rev-parse HEAD)

.PHONY: build
build:
	go build -ldflags "-X ${APP_PKG}/wstunnel/Version=${APP_VERSION} \
	-X ${APP_PKG}/wstunnel/BuildTime=${BUILD_TIME} \
	-X ${APP_PKG}/wstunnel/GitRevision=${GIT_REVISION}" \
	-o ./dist/server cmd/server/main.go
	go build -ldflags "-X ${APP_PKG}/wstunnel/Version=${APP_VERSION} \
	-X ${APP_PKG}/wstunnel/BuildTime=${BUILD_TIME} \
	-X ${APP_PKG}/wstunnel/GitRevision=${GIT_REVISION}" \
	-o ./dist/local cmd/local/main.go

.PHONY: image
image:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build
	docker build --rm -t ${APP_NAME}:${APP_VERSION} .

.PHONY: publish
publish: image
	docker tag ${APP_NAME}:${APP_VERSION} isayme/${APP_NAME}:${APP_VERSION}
	docker push isayme/${APP_NAME}:${APP_VERSION}
	docker tag ${APP_NAME}:${APP_VERSION} isayme/${APP_NAME}:latest
	docker push isayme/${APP_NAME}:latest
