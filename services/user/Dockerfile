FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app

# 1. Copy file go.mod dan go.sum dari common dan user service saja untuk caching dependencies
COPY services/common/go.mod services/common/go.sum ./services/common/
COPY services/user/go.mod services/user/go.sum ./services/user/

WORKDIR /app/services/user
RUN go mod download

# 2. Copy source code asli setelah dependencies di-cached
COPY services/common /app/services/common
COPY services/user /app/services/user

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/services/user/main .
EXPOSE 8082
CMD ["./main"]
