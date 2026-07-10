FROM golang:1.24-alpine AS builder

ENV GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct
ENV GOSUMDB=sum.golang.google.cn

WORKDIR /app
COPY go.mod go.sum ./
COPY pkg/ pkg/
COPY common/ common/
COPY api/ api/
COPY services/order-service/ services/order-service/

WORKDIR /app/services/order-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /order-service ./cmd

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /order-service /app/order-service
COPY services/order-service/config.yaml /app/config.yaml
EXPOSE 8888
CMD ["/app/order-service"]
