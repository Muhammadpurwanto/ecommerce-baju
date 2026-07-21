FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
# 1. Copy file go.mod dan go.sum dari common dan order service saja untuk caching dependencies
COPY services/common/go.mod services/common/go.sum ./services/common/
COPY services/order/go.mod services/order/go.sum ./services/order/

WORKDIR /app/services/order
RUN go mod download

# 2. Copy source code asli setelah dependencies di-cached
COPY services/common /app/services/common
COPY services/order /app/services/order

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/services/order/main .
EXPOSE 8085
CMD ["./main"]
