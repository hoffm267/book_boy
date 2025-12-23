FROM golang:1.24-alpine AS base

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
COPY external/ external/

FROM base AS development
RUN apk --no-cache add curl
ENV JWT_SECRET=test-secret-for-docker-build
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server
EXPOSE 8080
CMD ["./main"]

FROM golang:1.24-alpine AS debug
RUN go install github.com/go-delve/delve/cmd/dlv@latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ cmd/
COPY internal/ internal/
COPY external/ external/
EXPOSE 8080
EXPOSE 2345
CMD go build -gcflags="all=-N -l" -o main ./cmd/server && \
    dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./main

FROM base AS builder
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

FROM alpine:latest AS production
RUN apk --no-cache --no-scripts add ca-certificates curl && \
    apk fix --no-scripts
WORKDIR /app
COPY --from=builder /build/main .
EXPOSE 8080
CMD ["./main"]
