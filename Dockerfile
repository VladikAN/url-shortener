# Build image definition
FROM golang:1.14 AS builder
WORKDIR /src/
COPY src/ .
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Runtime image definition
FROM alpine:latest
LABEL maintainer="https://github.com/vladikan/url-shortener"
WORKDIR /root/
COPY --from=builder src/app src/config.yml.

ENV US_log_level="info" \
    US_host_port=":80"

EXPOSE 80
ENTRYPOINT ["./app"]