FROM golang:1.24-alpine AS base

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
COPY external/ external/

FROM base AS test
ENV JWT_SECRET=test-secret-for-docker-build
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server
EXPOSE 8080
CMD ["./main"]

FROM base AS builder
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

FROM alpine:latest AS production
RUN apk --no-cache add ca-certificates curl
WORKDIR /app
COPY --from=builder /build/main .
EXPOSE 8080
CMD ["./main"]
