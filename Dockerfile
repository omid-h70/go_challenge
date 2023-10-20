FROM golang:1.20.alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o -mode=vendor main.o
RUN apk --no-cache add curl
RUN $ curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-arm64.tar.gz | tar xvz

FROM alpine:3.13
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY /app.env .
COPY  start.sh .
COPY  wait-for-sh.sh .
COPY db/migration ./migaration

#its just for info, not actually exposing a port !
EXPOSE 8080
# when we use cmd instruction before entry point its like doing
# ENTRYPOINT["/app.start.sh", "app/main"]
CMD ["app/main"]
ENTRYPOINT["/app.start.sh"]