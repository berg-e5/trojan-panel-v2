#!/bin/bash
# ============================================================
# Trojan Panel 一键部署脚本
# 支持：后端 + 前端 + MySQL + Redis
# ============================================================

set -e

# -------------------- 配置区 --------------------
# 修改以下配置为你的实际值
MYSQL_ROOT_PASSWORD="TrojanPanel@2024"
MYSQL_HOST="127.0.0.1"
MYSQL_PORT="3306"
MYSQL_DATABASE="trojan_panel"
MYSQL_USER="root"

REDIS_PASSWORD="TrojanPanel@2024"
REDIS_HOST="127.0.0.1"
REDIS_PORT="6379"

SERVER_PORT="8081"
TZ="Asia/Shanghai"

BACKEND_IMAGE="ghcr.io/berg-e5/trojan-panel"
FRONTEND_IMAGE="ghcr.io/berg-e5/trojan-panel-ui"
# ------------------------------------------------

# 颜色
RED='\033[31m'
GREEN='\033[32m'
YELLOW='\033[33m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_err()  { echo -e "${RED}[ERROR]${NC} $1"; }

# 检测 Docker
check_docker() {
    if ! command -v docker &> /dev/null; then
        log_warn "Docker 未安装，正在安装..."
        curl -fsSL https://get.docker.com | sh
        systemctl enable docker
        systemctl start docker
        log_info "Docker 安装完成"
    else
        log_info "Docker 已安装: $(docker --version)"
    fi
}

# 拉取镜像
pull_images() {
    log_info "正在拉取镜像..."

    docker pull ghcr.io/berg-e5/trojan-panel:latest
    docker pull ghcr.io/berg-e5/trojan-panel-ui:latest
    docker pull mysql:8.0
    docker pull redis:7

    log_info "镜像拉取完成"
}

# 启动 MySQL
start_mysql() {
    log_info "启动 MySQL..."

    if docker ps -a | grep -q "^mysql "; then
        docker rm -f mysql 2>/dev/null || true
    fi

    docker run -d \
        --name mysql \
        --network host \
        -e MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD} \
        -e MYSQL_DATABASE=${MYSQL_DATABASE} \
        -v /opt/mysql/data:/var/lib/mysql \
        --restart unless-stopped \
        mysql:8.0 \
        --character-set-server=utf8mb4 \
        --collation-server=utf8mb4_unicode_ci

    log_info "MySQL 启动中，等待初始化..."
    sleep 15
    log_info "MySQL 启动完成"
}

# 启动 Redis
start_redis() {
    log_info "启动 Redis..."

    if docker ps -a | grep -q "^redis "; then
        docker rm -f redis 2>/dev/null || true
    fi

    docker run -d \
        --name redis \
        --network host \
        -e REDIS_PASSWORD=${REDIS_PASSWORD} \
        -v /opt/redis/data:/data \
        --restart unless-stopped \
        redis:7 --requirepass ${REDIS_PASSWORD}

    log_info "Redis 启动完成"
}

# 启动后端
start_backend() {
    log_info "启动后端服务..."

    if docker ps | grep -q "^trojan-panel "; then
        docker rm -f trojan-panel 2>/dev/null || true
    fi

    docker run -d \
        --name trojan-panel \
        --network host \
        --restart unless-stopped \
        -e mariadb_ip=${MYSQL_HOST} \
        -e mariadb_port=${MYSQL_PORT} \
        -e mariadb_user=${MYSQL_USER} \
        -e mariadb_pas=${MYSQL_ROOT_PASSWORD} \
        -e redis_host=${REDIS_HOST} \
        -e redis_port=${REDIS_PORT} \
        -e redis_pass=${REDIS_PASSWORD} \
        -e server_port=${SERVER_PORT} \
        -e TZ=${TZ} \
        ${BACKEND_IMAGE}:latest

    log_info "后端启动完成，端口: ${SERVER_PORT}"
}

# 启动前端
start_frontend() {
    log_info "启动前端服务..."

    if docker ps | grep -q "^trojan-panel-ui "; then
        docker rm -f trojan-panel-ui 2>/dev/null || true
    fi

    docker run -d \
        --name trojan-panel-ui \
        --network host \
        --restart unless-stopped \
        -e TZ=${TZ} \
        ${FRONTEND_IMAGE}:latest

    log_info "前端启动完成，端口: 80"
}

# 查看状态
show_status() {
    echo ""
    echo "=========================================="
    echo "              服务状态"
    echo "=========================================="
    docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep -E "mysql|redis|trojan"
    echo ""
    echo "=========================================="
    echo "              访问地址"
    echo "=========================================="
    echo "  前端面板: http://localhost:80"
    echo "  后端 API: http://localhost:${SERVER_PORT}"
    echo "=========================================="
}

# 主流程
main() {
    echo ""
    echo "=========================================="
    echo "       Trojan Panel 一键部署"
    echo "=========================================="
    echo ""

    check_docker
    pull_images
    start_mysql
    start_redis
    start_backend
    start_frontend

    sleep 3
    show_status

    log_info "部署完成!"
}

main "$@"
