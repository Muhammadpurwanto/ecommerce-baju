FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
# 1. Copy file go.mod dan go.sum dari common dan gateway service saja untuk caching dependencies
COPY services/common/go.mod services/common/go.sum ./services/common/
COPY services/gateway/go.mod services/gateway/go.sum ./services/gateway/

WORKDIR /app/services/gateway
RUN go mod download

# 2. Copy source code asli setelah dependencies di-cached
COPY services/common /app/services/common
COPY services/gateway /app/services/gateway

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/services/gateway/main .
EXPOSE 8080
CMD ["./main"]
