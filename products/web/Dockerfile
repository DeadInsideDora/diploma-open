FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

RUN apk add --no-cache ca-certificates

CMD ["./main", "/app/configs/categories_info.json"]