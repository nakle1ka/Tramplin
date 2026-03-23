FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o ./bin/app ./cmd/app/main.go
RUN go build -o ./bin/seed ./cmd/admin-seed/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/bin/app .
COPY --from=builder /app/bin/seed .
COPY .env .

CMD ["./app"]
