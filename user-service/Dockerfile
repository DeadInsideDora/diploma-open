FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
