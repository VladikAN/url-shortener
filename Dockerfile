# Build image definition
FROM golang:1.14 AS builder
WORKDIR /src/
COPY . .
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Runtime image definition
FROM alpine:latest
LABEL maintainer="https://github.com/vladikan/url-shortener"
WORKDIR /root/
COPY --from=builder src/app .
COPY --from=builder src/config.yaml .

ENV US_LOG_LEVEL="info" \
    US_HOST_ADDR=":80" \
    US_HOST_SSL="false"

EXPOSE 443
EXPOSE 80
VOLUME ["/autocert"]
ENTRYPOINT ["./app"]