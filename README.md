# Trojan Panel V2

基于 [trojanpanel/trojan-panel](https://github.com/trojanpanel/trojan-panel) 的二次开发版本，支持深色模式、批量管理、数据统计等新功能。

## 功能特性

- 🌙 **深色模式** - 一键切换亮/暗主题
- 📧 **通知系统** - 支持邮件 + Telegram 账号通知
- 📊 **数据统计** - 用户/节点/流量可视化图表（ECharts）
- 🔧 **批量管理** - 批量启用/禁用/续期/删除/重置流量
- 🚀 **CI/CD** - GitHub Actions 自动构建 Docker 镜像

## 架构

```
┌─────────────────┐      ┌─────────────────┐
│   Frontend      │      │    Backend      │
│  (Vue + Nginx)  │ ─── │   (Go + Gin)    │
│   Port: 80      │      │   Port: 8081   │
└─────────────────┘      └────────┬────────┘
                                  │
                    ┌─────────────┼─────────────┐
                    │             │             │
              ┌─────▼─────┐ ┌────▼─────┐ ┌────▼─────┐
              │   MySQL    │ │   Redis   │ │ Xray Core │
              │   Port 3306│ │ Port 6379 │ │  节点同步 │
              └───────────┘ └───────────┘ └───────────┘
```

## 快速部署

### 一键部署（推荐）

```bash
curl -fsSL https://raw.githubusercontent.com/Berg-E5/Trojan-Panel-V2/main/deploy.sh | bash
```

> 脚本会自动安装 Docker、启动 MySQL、Redis、前端、后端全部服务

### Docker 手动部署

```bash
# 1. 启动 MySQL
docker run -d --name mysql \
  --network host \
  -e MYSQL_ROOT_PASSWORD=你的密码 \
  -e MYSQL_DATABASE=trojan_panel \
  -v /opt/mysql/data:/var/lib/mysql \
  --restart unless-stopped \
  mysql:8.0

# 2. 启动 Redis
docker run -d --name redis \
  --network host \
  -e REDIS_PASSWORD=你的密码 \
  -v /opt/redis/data:/data \
  --restart unless-stopped \
  redis:7 --requirepass 你的密码

# 3. 启动后端
docker run -d --name trojan-panel \
  --network host \
  -e mariadb_ip=127.0.0.1 \
  -e mariadb_port=3306 \
  -e mariadb_user=root \
  -e mariadb_pas=你的数据库密码 \
  -e redis_host=127.0.0.1 \
  -e redis_port=6379 \
  -e redis_pass=你的Redis密码 \
  -e server_port=8081 \
  --restart unless-stopped \
  ghcr.io/berg-e5/trojan-panel:latest

# 4. 启动前端
docker run -d --name trojan-panel-ui \
  --network host \
  --restart unless-stopped \
  ghcr.io/berg-e5/trojan-panel-ui:latest
```

## 访问地址

| 服务 | 地址 |
|------|------|
| 前端面板 | http://你的服务器IP |
| 后端 API | http://你的服务器IP:8081 |

## 项目结构

```
Berg-E5/Trojan-Panel-V2      # 后端服务
Berg-E5/Trojan-Panel-UI      # 前端界面
Berg-E5/Trojan-Panel-Core    # 节点同步核心
```

## 开发构建

```bash
# 后端编译
cd trojan-panel
go mod download
go build -o trojan-panel .

# 前端构建
cd trojan-panel-ui
yarn install
yarn build
```

## 主要文件

| 文件 | 说明 |
|------|------|
| `Dockerfile` | 多架构 Docker 构建（amd64/arm64 等） |
| `deploy.sh` | 一键部署脚本 |
| `.github/workflows/` | GitHub Actions CI/CD |

## Docker 镜像

| 镜像 | 地址 |
|------|------|
| 后端 | `ghcr.io/berg-e5/trojan-panel` |
| 前端 | `ghcr.io/berg-e5/trojan-panel-ui` |

## License

MIT
