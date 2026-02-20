# ginchat
基于 Gin+GORM 开发的高性能、易扩展的实时通讯 IM 项目，支持单聊、群聊、在线状态管理、WebSocket长连接等核心能力。

## 技术栈
- 后端框架：Gin
- ORM：GORM
- 配置管理：Viper
- 日志：Zap
- 数据库：MySQL/PostgreSQL
- 缓存：Redis
- 长连接：WebSocket

## 核心功能
- ✅ 用户注册/登录/鉴权
- ✅ 单聊实时消息收发
- ✅ 群聊/房间消息广播
- ✅ 用户在线状态管理
- ✅ 离线消息推送
- ✅ 心跳保活机制

## 快速启动
### 环境要求
- Go 1.21+
- MySQL 8.0+ / PostgreSQL 14+
- Redis 6.0+

### 部署步骤
1. 克隆项目
```bash
git clone https://github.com/USERON8/ginchat.git
cd ginchat

开发规范
项目开发遵循 RULE.md 开发规范，参与开发前请先阅读。
开源协议
本项目基于 MIT License 开源，详见 LICENSE 文件。
plaintext
