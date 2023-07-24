FROM golang:1.19-alpine3.17 AS builder

COPY . /app

WORKDIR /app

RUN go build -o srvd cmd/main.go

# stage 2
FROM alpine:3.17

ENV APP_HOME /app

WORKDIR $APP_HOME

STOPSIGNAL SIGTERM

COPY --from=builder /app/srvd /usr/local/bin/srvd

CMD ["srvd", "server", "-a", "0.0.0.0:5001"]
