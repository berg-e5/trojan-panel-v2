FROM golang:1.22-alpine AS builder

WORKDIR /app

# 安装编译依赖
RUN apk add --no-cache bash git ca-certificates

# 设置 Go proxy（避免网络问题）
ENV GOPROXY=https://proxy.golang.org,direct

# 安装 garble（固定版本）
RUN go install mvdan.cc/garble/cmd/garble@v0.15.0

# 复制源码
COPY . .

# 编译（支持多架构）
ARG TARGETARCH
ARG TARGETVARIANT
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH}${TARGETVARIANT} \
    garble -literals -tiny build \
    -o trojan-panel \
    -ldflags "-s -w"

# 运行镜像
FROM alpine:3.15

WORKDIR /tpdata/trojan-panel/

ENV TZ=Asia/Shanghai \
    GIN_MODE=release

RUN apk add --no-cache bash tzdata ca-certificates

COPY --from=builder /app/trojan-panel .

ENTRYPOINT chmod 777 ./trojan-panel && \
    ./trojan-panel \
    -host=${mariadb_ip:-127.0.0.1} \
    -port=${mariadb_port:-9507} \
    -user=${mariadb_user:-root} \
    -password=${mariadb_pas:-123456} \
    -redisHost=${redis_host:-127.0.0.1} \
    -redisPort=${redis_port:-6378} \
    -redisPassword=${redis_pass:-123456} \
    -serverPort=${server_port:-8081}
