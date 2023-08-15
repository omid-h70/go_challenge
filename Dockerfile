FROM golang:1.20.alpine AS builder
WORKDIR /app
COPY . .

FROM alpine:3.13
WORKDIR /app
COPY --from=builder /app/main .

#its just for info, not actually exposing a port !
EXPOSE 8080
CMD ["app/main"]