FROM alpine
WORKDIR /app

COPY ./dist/server /app/server
COPY ./dist/local /app/local

CMD ["/app/server"]
