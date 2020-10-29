FROM golang:1.15.3-alpine as builder
WORKDIR /app

ARG APP_NAME
ENV APP_NAME ${APP_NAME}
ARG APP_VERSION
ENV APP_VERSION ${APP_VERSION}
ARG BUILD_TIME
ENV BUILD_TIME ${BUILD_TIME}
ARG GIT_REVISION
ENV GIT_REVISION ${GIT_REVISION}

COPY . .
RUN GO111MODULE=on GOPROXY=https://goproxy.io,direct go mod download \
  && go build -ldflags "-X github.com/isayme/go-wstunnel/wstunnel.Name=${APP_NAME} \
  -X github.com/isayme/go-wstunnel/wstunnel.Version=${APP_VERSION} \
  -X github.com/isayme/go-wstunnel/wstunnel.BuildTime=${BUILD_TIME} \
  -X github.com/isayme/go-wstunnel/wstunnel.GitRevision=${GIT_REVISION}" \
  -o ./dist/wstunnel main.go

FROM alpine
WORKDIR /app

ARG APP_NAME
ENV APP_NAME ${APP_NAME}
ARG APP_VERSION
ENV APP_VERSION ${APP_VERSION}

COPY --from=builder /app/dist/wstunnel ./

CMD ["/app/wstunnel", "server"]
