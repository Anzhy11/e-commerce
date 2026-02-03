# Builder
FROM golang:1.25.6-alpine AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .
RUN go run ./scripts/build.go

# Runner
FROM golang:1.25.6-alpine AS runner

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/bin .

EXPOSE 8080

CMD ["./api"]
