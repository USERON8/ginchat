# GinChat

基于 Go + Gin 构建的单体 IM 即时通讯服务，支持私聊、群聊、好友管理，使用 WebSocket 实现实时消息推送。

---

## 技术栈

| 分类     | 技术                                |
|--------|-----------------------------------|
| Web 框架 | Gin                               |
| 数据库    | PostgreSQL + GORM                 |
| 缓存     | Redis                             |
| 实时通讯   | WebSocket（gorilla/websocket）      |
| 认证     | JWT（access token + refresh token） |
| 日志     | Zap + Lumberjack                  |
| 配置     | Viper                             |
| 部署     | Docker Compose                    |

---

## 功能特性

### 用户模块

- 注册 / 登录 / 退出
- 查询 / 修改个人信息
- 搜索用户
- 修改密码
- 用户信息 Redis 缓存

### 好友模块

- 发送 / 同意 / 拒绝好友申请
- 好友列表
- 删除好友

### 群组模块

- 创建群组
- 邀请 / 踢出成员
- 成员列表
- 我的群列表

### 消息模块

- 私聊实时消息（WebSocket）
- 群聊实时消息（WebSocket）
- 历史消息（游标分页）
- 私聊已读回执
- 未读消息数
- 消息去重（Redis SETNX）

### 通讯优化

- 在线状态管理（Redis）
- 心跳检测（Ping/Pong）
- Worker Pool 异步写库
- IP 限流（全局 + 严格模式）

---

## 项目结构

```
ginchat/
├── config/
│   ├── config.yaml             # 本地配置（已 gitignore）
│   └── config.example.yaml     # 配置模板
├── internal/
│   ├── handler/                # 路由处理层
│   │   ├── user.go
│   │   ├── friend.go
│   │   ├── group.go
│   │   ├── message.go
│   │   └── ws.go
│   ├── model/                  # 数据模型
│   │   ├── user.go
│   │   ├── friendship.go
│   │   ├── message.go
│   │   └── group.go
│   └── ws/                     # WebSocket 核心
│       ├── client.go           # ReadPump / WritePump / handleMessage
│       ├── manager.go          # 在线用户管理
│       ├── message.go          # 消息结构体
│       └── worker.go           # 异步写库 Worker Pool
├── middleware/
│   ├── auth.go                 # JWT 鉴权
│   ├── logger.go               # 请求日志
│   └── ratelimit.go            # IP 限流
├── pkg/
│   ├── config/                 # 配置加载
│   ├── database/               # DB 初始化 + 连接池
│   ├── jwt/                    # Token 生成 / 解析
│   ├── logger/                 # Zap 日志封装
│   ├── redis/                  # Redis 操作封装
│   │   ├── redis.go
│   │   ├── user_cache.go       # 用户信息缓存
│   │   ├── online.go           # 在线状态
│   │   ├── dedup.go            # 消息去重
│   │   └── token.go            # Refresh Token
│   └── response/               # 统一响应 + 错误码
├── router/
│   └── router.go               # 路由注册
├── docker-compose.yml
├── init.sql                    # 建表 SQL
├── test.sql                    # 测试数据
├── go.mod
└── main.go
```

---

## 快速开始

### 环境要求

- Go 1.21+
- Docker + Docker Compose

### 1. 克隆项目

```bash
git clone https://github.com/yourname/ginchat.git
cd ginchat
```

### 2. 配置文件

```bash
cp config/config.example.yaml config/config.yaml
```

按需修改 `config/config.yaml`：

```yaml
server:
  port: 8080

database:
  host: 127.0.0.1
  port: 5432
  username: root
  password: root
  dbname: im_db
  sslmode: disable
  max_idle_conns: 20
  max_open_conns: 100
  conn_max_lifetime: 3600

redis:
  host: 127.0.0.1
  port: 6379
  password: "root"
  db: 0
  pool_size: 100

jwt:
  secret: "your-secret-key"
  expire: 2
```

### 3. 启动依赖服务

```bash
# 启动 PostgreSQL + Redis
docker-compose up -d

# 查看状态
docker-compose ps
```

### 4. 启动服务

```bash
go run main.go
```

服务启动后访问 `http://localhost:8080`

---

## 接口文档

### 认证说明

需要登录的接口在 Header 中携带：

```
Authorization: Bearer <access_token>
```

WebSocket 连接通过 URL 参数传 token：

```
ws://localhost:8080/api/ws?token=<access_token>
```

---

### 用户接口

| 方法   | 路径                 | 说明       | 是否需要登录 |
|------|--------------------|----------|--------|
| POST | /api/user/register | 注册       | ❌      |
| POST | /api/user/login    | 登录       | ❌      |
| POST | /api/user/refresh  | 刷新 Token | ❌      |
| GET  | /api/user/info     | 获取个人信息   | ✅      |
| PUT  | /api/user/info     | 修改个人信息   | ✅      |
| PUT  | /api/user/password | 修改密码     | ✅      |
| GET  | /api/user/search   | 搜索用户     | ✅      |
| POST | /api/user/logout   | 退出登录     | ✅      |

**注册**

```json
POST /api/user/register
{
  "username": "test1",
  "password": "123456",
  "phone": "13800000001",
  "email": "test1@example.com"
}
```

**登录**

```json
POST /api/user/login
{
  "username": "test1",
  "password": "123456"
}

// 返回
{
  "code": 0,
  "msg": "success",
  "data": {
    "accessToken": "eyJ...",
    "refreshToken": "eyJ...",
    "userID": 1
  }
}
```

**刷新 Token**

```json
POST /api/user/refresh
{
  "refreshToken": "eyJ..."
}
```

---

### 好友接口

| 方法     | 路径                    | 说明      |
|--------|-----------------------|---------|
| POST   | /api/friend/apply/:id | 发送好友申请  |
| PUT    | /api/friend/apply     | 同意/拒绝申请 |
| GET    | /api/friend/applies   | 收到的申请列表 |
| GET    | /api/friend/list      | 好友列表    |
| DELETE | /api/friend           | 删除好友    |

**同意/拒绝申请**

```json
PUT /api/friend/apply
{
  "applyId": 1,
  "action": "accepted"
  // accepted / rejected
}
```

---

### 群组接口

| 方法     | 路径                      | 说明     |
|--------|-------------------------|--------|
| POST   | /api/group              | 创建群    |
| GET    | /api/group/list         | 我的群列表  |
| GET    | /api/group/:id          | 群详情    |
| GET    | /api/group/:id/members  | 成员列表   |
| POST   | /api/group/:id/member   | 邀请成员   |
| DELETE | /api/group/:id/member   | 踢出成员   |
| GET    | /api/group/:id/messages | 群聊历史消息 |

---

### 消息接口

| 方法  | 路径                   | 说明     |
|-----|----------------------|--------|
| GET | /api/message/history | 私聊历史消息 |
| GET | /api/message/unread  | 未读消息数  |
| PUT | /api/message/read    | 标记已读   |

**历史消息（游标分页）**

```
GET /api/message/history?friendId=2&lastId=0&size=20

// lastId=0 → 获取最新 20 条
// lastId=100 → 获取 id < 100 的 20 条（加载更多）
```

---

### WebSocket 消息格式

**连接**

```
ws://localhost:8080/api/ws?token=<access_token>
```

**发送私聊消息**

```json
{
  "msgId": "唯一ID（用于去重）",
  "type": "private",
  "toId": 2,
  "content": "你好"
}
```

**发送群聊消息**

```json
{
  "msgId": "唯一ID",
  "type": "group",
  "toId": 1,
  "content": "大家好"
}
```

**心跳**

```json
// 发送
{
  "type": "ping"
}

// 服务端返回
{
  "type": "pong"
}
```

---

## 消息链路

```
发送方
  │
  │  WebSocket 发送 JSON
  ▼
ReadPump 接收
  │
  ├─► 消息去重（Redis SETNX，60s 内同一 msgId 只处理一次）
  │
  ├─► Worker Pool 异步写数据库（不阻塞转发）
  │
  └─► 转发
        │
        ├─ 私聊 → Manager 找到目标用户连接 → WritePump 推送
        │
        └─ 群聊 → 查群成员 → 逐一推送在线成员


目标用户不在线
  └─► 消息已落库，上线后调历史消息接口拉取
```

---

## 统一响应格式

```json
{
  "code": 0,
  // 0=成功，非0=失败
  "msg": "success",
  "data": {}
}
```

**错误码说明**

| 错误码  | 说明       |
|------|----------|
| 0    | 成功       |
| 1001 | 参数错误     |
| 1002 | 服务器错误    |
| 1003 | 请求过于频繁   |
| 2001 | 用户不存在    |
| 2002 | 用户名已存在   |
| 2003 | 密码错误     |
| 2004 | Token 无效 |
| 2005 | 未登录      |
| 2006 | 注册失败     |
| 2007 | 更新失败     |
| 3001 | 已申请或已是好友 |
| 3002 | 申请不存在    |
| 3003 | 无权操作     |
| 3004 | 不能添加自己   |
| 5001 | 群组不存在    |
| 5002 | 群组无权操作   |
| 5003 | 用户已在群中   |

---

## License

MIT