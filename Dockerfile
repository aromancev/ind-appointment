FROM golang:1.16-alpine as builder

WORKDIR /app
COPY . /app

RUN go build -o watch .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/watch /app/

CMD ["/app/watch"]
