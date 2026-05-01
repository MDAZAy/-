FROM golang:1.23 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/agent ./cmd/agent

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /out/agent /app/agent
COPY config /app/config
COPY web /app/web

EXPOSE 8080

ENTRYPOINT ["/app/agent"]
