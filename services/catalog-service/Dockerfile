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
COPY services/catalog-service/ services/catalog-service/
RUN echo "cachebust=${CACHEBUST}" >/tmp/cachebust

WORKDIR /app/services/catalog-service
RUN go mod tidy && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /catalog-service ./cmd

ARG ALPINE_IMAGE=alpine:3.19
FROM ${ALPINE_IMAGE}
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /catalog-service /app/catalog-service
COPY services/catalog-service/etc/ /app/etc/
ENV CONFIG_PATH=/app/etc/catalog-service.yaml
EXPOSE 8888 9090
CMD ["/app/catalog-service"]
