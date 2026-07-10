FROM golang:1.24-alpine AS builder

ENV GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct
ENV GOSUMDB=sum.golang.google.cn

WORKDIR /app
COPY go.mod go.sum ./
COPY pkg/ pkg/
COPY common/ common/
COPY api/ api/
COPY services/user-service/ services/user-service/

WORKDIR /app/services/user-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /user-service ./cmd

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /user-service /app/user-service
COPY services/user-service/config.yaml /app/config.yaml
EXPOSE 8888 9090
CMD ["/app/user-service"]
