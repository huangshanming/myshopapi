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
COPY services/lottery-service/ services/lottery-service/
RUN echo "cachebust=${CACHEBUST}" >/tmp/cachebust

WORKDIR /app/services/lottery-service
RUN go mod tidy && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /lottery-service ./cmd

ARG ALPINE_IMAGE=alpine:3.19
FROM ${ALPINE_IMAGE}
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /lottery-service /app/lottery-service
COPY services/lottery-service/etc/ /app/etc/
ENV CONFIG_PATH=/app/etc/lottery-service.yaml
EXPOSE 8888
CMD ["/app/lottery-service"]
