ARG GOLANG_IMAGE=golang:1.24-alpine
ARG ALPINE_IMAGE=alpine:3.19
FROM ${GOLANG_IMAGE} AS builder
ARG CACHEBUST=1

ENV GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct
ENV GOSUMDB=sum.golang.google.cn

WORKDIR /app
COPY go.mod go.sum ./
COPY pkg/ pkg/
COPY common/ common/
COPY api/ api/
COPY services/order-service/ services/order-service/
RUN echo "cachebust=${CACHEBUST}" >/tmp/cachebust

WORKDIR /app/services/order-service
RUN go mod tidy && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /order-service ./cmd

ARG ALPINE_IMAGE=alpine:3.19
FROM ${ALPINE_IMAGE}
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /order-service /app/order-service
COPY services/order-service/etc/ /app/etc/
ENV CONFIG_PATH=/app/etc/order-service.yaml
EXPOSE 8888
CMD ["/app/order-service"]
