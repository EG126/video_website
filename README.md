# Video Website Backend

视频网站后端服务，提供用户管理、视频上传、社交互动、实时聊天等核心功能。

## 📋 文档目录

| 文档 | 描述 | 路径 |
| --- | --- | --- |
| API 接口文档 | 完整的 API 接口定义 | [API接口文档](https://s.apifox.cn/2c42f18c-ed41-434f-9fde-3166b87f69b9)  |
| Redis 架构文档 | Redis 缓存设计与接口原理 | [REDIS架构文档](https://gcntmecan6g1.feishu.cn/wiki/JSRrw6TbyiTISWkOQzkcpRX1npb) |
| 消息队列文档 | Redis Pub/Sub 消息队列设计 | [消息队列架构文档](https://gcntmecan6g1.feishu.cn/wiki/OpVHwZSA8idrNkkKqdecaTqlnbf) |
| Docker 部署文档 | 容器化部署指南 | [Docker部署文档](https://gcntmecan6g1.feishu.cn/wiki/MDjCwpjW9iziipkVfrLcalwkn7e) |

## 🛠️ 技术栈

### 核心技术

- **语言**: Go
- **Web 框架**: Hertz
- **ORM**: GORM
- **数据库驱动**: MySQL Driver
- **缓存客户端**: Redis Client
- **实时通信**: WebSocket
- **认证**: JWT
- **二步验证**: OTP
- **配置管理**: Viper
- **ID 生成**: Snowflake
- **密码加密**: Crypto

### 基础设施

- **数据库**: MySQL
- **缓存**: Redis
- **容器化**: Docker
- **容器编排**: Docker Compose

## 🚀 快速开始

### 本地开发

```bash
# 克隆项目
git clone https://github.com/EG126/video_website
cd video_website

# 安装依赖
go mod tidy

# 配置数据库连接 (config/config.yaml)
# 修改 MySQL 和 Redis 连接信息

# 运行服务
go run main.go
```

### Docker 部署

详细部署指南请参考 [Docker部署文档](https://gcntmecan6g1.feishu.cn/wiki/MDjCwpjW9iziipkVfrLcalwkn7e)


## 📁 项目结构

```
video_website/
├── .github/
│   └── workflows/          # CI/CD 配置
├── biz/                    # 业务逻辑层
│   ├── dal/                # 数据访问层
│   │   ├── mysql/          # MySQL 操作
│   │   │   └── entity/     # 数据库实体
│   │   └── redis/          # Redis 操作
│   ├── handler/            # HTTP 处理器
│   ├── middleware/         # 中间件
│   ├── model/              # 数据模型
│   ├── router/             # 路由注册
│   └── service/            # 业务服务
├── config/                 # 配置文件
├── idl/                    # IDL 定义文件
├── pkg/                    # 公共包
│   ├── bcrypt/             # 密码加密
│   ├── constants/          # 常量定义
│   ├── errno/              # 错误定义
│   ├── jwt/                # JWT 工具
│   ├── response/           # 响应封装
│   └── utils/              # 工具函数
├── static/                 # 静态资源
│   ├── avatars/            # 用户头像
│   └── videos/             # 视频文件
├── .dockerignore
├── .gitignore
├── .golangci.yml
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
├── README.md
├── router.go
└── router_gen.go
```


## 🔌 API 接口

完整的 API 接口定义请参考 [API接口文档](https://s.apifox.cn/2c42f18c-ed41-434f-9fde-3166b87f69b9)

### 用户模块

| 方法 | 路径 | 描述 |
| --- | --- | --- |
| POST | /user/register | 用户注册 |
| POST | /user/login | 用户登录 |
| GET | /user/info | 获取用户信息 |
| PUT | /user/avatar/upload | 上传头像 |
| GET | /auth/mfa/qrcode | 获取 MFA 二维码 |
| POST | /auth/mfa/bind | 绑定 MFA |
| POST | /user/refresh | 刷新 Token |

### 视频模块

| 方法 | 路径 | 描述 |
| --- | --- | --- |
| POST | /video/publish | 发布视频 |
| GET | /video/list | 获取用户视频列表 |
| GET | /video/popular | 获取热门视频 |
| GET | /video/feed/ | 获取视频流 |
| POST | /video/search | 搜索视频 |

### 社交模块

| 方法 | 路径 | 描述 |
| --- | --- | --- |
| POST | /relation/action | 关注/取消关注 |
| GET | /follower/list | 获取粉丝列表 |
| GET | /following/list | 获取关注列表 |
| GET | /friends/list | 获取好友列表 |

### 互动模块

| 方法 | 路径 | 描述 |
| --- | --- | --- |
| POST | /like/action | 点赞/取消点赞 |
| GET | /like/list | 获取点赞列表 |
| POST | /comment/publish | 发布评论 |
| GET | /comment/list | 获取评论列表 |
| DELETE | /comment/delete | 删除评论 |


### 实时聊天（WebSocket）

| 连接方式 | 地址 | 描述 |
| --- | --- | --- |
| WebSocket | ws://localhost:6666/ws | 实时聊天服务 |

**消息类型**

| 类型 | 描述 |
| --- | --- |
| type1 | 私聊消息 |
| type2 | 获取私聊历史 |
| type3 | 获取未读消息 |
| type4 | 群聊消息 |
| type5 | 获取群聊历史 |
| check_online | 检查用户在线状态 |
| get_unread_count | 获取未读消息数量 |

### 健康检查

| 方法 | 路径 | 描述 |
| --- | --- | --- |
| GET | /ping | 服务健康检查 |


## 📖 架构文档

### Redis 缓存设计

详细的 Redis 架构设计请参考 [REDIS架构文档](https://gcntmecan6g1.feishu.cn/wiki/JSRrw6TbyiTISWkOQzkcpRX1npb)


### 消息队列架构

消息队列设计与实现请参考 [消息队列架构文档](https://gcntmecan6g1.feishu.cn/wiki/OpVHwZSA8idrNkkKqdecaTqlnbf)
