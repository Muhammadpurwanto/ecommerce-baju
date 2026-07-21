FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
# 1. Copy file go.mod dan go.sum dari common dan cart service saja untuk caching dependencies
COPY services/common/go.mod services/common/go.sum ./services/common/
COPY services/cart/go.mod services/cart/go.sum ./services/cart/

WORKDIR /app/services/cart
RUN go mod download

# 2. Copy source code asli setelah dependencies di-cached
COPY services/common /app/services/common
COPY services/cart /app/services/cart

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/services/cart/main .
EXPOSE 8084
CMD ["./main"]
