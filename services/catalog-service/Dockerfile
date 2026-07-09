FROM golang:1.24-alpine AS builder

ENV GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct
ENV GOSUMDB=sum.golang.google.cn

WORKDIR /app
COPY go.mod go.sum ./
COPY pkg/ pkg/
COPY common/ common/
COPY services/catalog-service/ services/catalog-service/

WORKDIR /app/services/catalog-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /catalog-service ./cmd

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /catalog-service /app/catalog-service
EXPOSE 8888
CMD ["/app/catalog-service"]
