FROM golang:1.24-alpine

WORKDIR /backend

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/server

EXPOSE 8080

CMD ["/main"]

