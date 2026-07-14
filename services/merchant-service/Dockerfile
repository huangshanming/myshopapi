ARG GOLANG_IMAGE=golang:1.24-alpine
ARG ALPINE_IMAGE=alpine:3.19
FROM ${GOLANG_IMAGE} AS builder

ENV GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct
ENV GOSUMDB=sum.golang.google.cn

WORKDIR /app
COPY go.mod go.sum ./
COPY pkg/ pkg/
COPY common/ common/
COPY api/ api/
COPY services/merchant-service/ services/merchant-service/

WORKDIR /app/services/merchant-service
RUN go mod tidy && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /merchant-service ./cmd

ARG ALPINE_IMAGE=alpine:3.19
FROM ${ALPINE_IMAGE}
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /merchant-service /app/merchant-service
COPY services/merchant-service/config.yaml /app/config.yaml
EXPOSE 8888
CMD ["/app/merchant-service"]
