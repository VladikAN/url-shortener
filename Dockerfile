# Build image definition
FROM golang:1.14 AS builder
WORKDIR /src
COPY . .
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime image definition
FROM alpine:latest
LABEL maintainer="https://github.com/vladikan/url-shortener"
WORKDIR /app
COPY --from=builder main .
EXPOSE 80
CMD ["./main"]