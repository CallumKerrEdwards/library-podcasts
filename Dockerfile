FROM golang:1.19-alpine as dependencies
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

FROM dependencies as src
COPY cmd/ ./cmd
COPY internal/ ./internal

FROM src as test
RUN CGO_ENABLED=0 GOOS=linux go test -v ./...

FROM test as builder
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/podcasts cmd/podcasts/main.go

FROM alpine:latest AS production
WORKDIR /app
COPY --from=builder /bin/podcasts .

HEALTHCHECK  --interval=5s --timeout=2s --start-period=10s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthcheck || exit 1

CMD ["/app/podcasts"]