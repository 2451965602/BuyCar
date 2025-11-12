# BuyCar

基于 CloudWeGo Hertz 框架的 Go 语言后端服务项目，为汽车购买咨询平台提供 API 服务。

## 目录

- [项目简介](#项目简介)
- [技术栈](#技术栈)
- [环境要求](#环境要求)
- [快速启动](#快速启动)
- [配置文件说明](#配置文件说明)
- [项目结构](#项目结构)
- [核心功能](#核心功能)

## 项目简介

本项目采用 CloudWeGo 生态的 Hertz 框架开发，提供高性能的 HTTP 服务。项目集成了 MySQL、JWT 认证、AI 咨询等功能模块，为用户提供汽车购买咨询、评分、用户管理等一站式服务。

## 技术栈

- **CloudWeGo Hertz** - 高性能 HTTP 框架，用于构建 API 服务
- **MySQL** - 关系型数据库，用于存储结构化数据
- **GORM** - ORM 框架，用于简化数据库操作
- **JWT** - JSON Web Token，用于用户认证和授权
- **Viper** - 配置管理工具，用于读取和管理应用配置
- **OpenAI API** - AI 智能咨询服务

## 环境要求

- Go 1.25.4+
- MySQL 5.7+
- Docker & Docker Compose (可选)

## 快速启动

### 本地部署

1. **克隆项目**

```bash
git clone https://github.com/2451965602/BuyCar.git
cd BuyCar
```

2. **安装依赖**

```bash
go mod tidy
```

3. **配置数据库**

创建配置文件：

```bash
cp config/config.example.yaml config/config.yaml
```

修改 `config/config.yaml` 中的数据库连接信息：

```yaml
mysql:
  addr: "localhost:3306"
  database: "buycar"
  username: "root"
  password: "your_password"
  charset: "utf8mb4"
```

4. **初始化数据库**

创建数据库：

```sql
CREATE DATABASE buycar CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

5. **启动服务**

```bash
go run main.go
```

服务将在配置的地址和端口上启动（默认配置请查看 `config/config.yaml`）。

## 配置文件说明

项目配置文件使用 YAML 格式，位于 `config/config.yaml`

### 配置项说明

```yaml
# 数据库配置
mysql:
  addr: "localhost:3306"        # MySQL 服务器地址
  database: "buycar"             # 数据库名称
  username: "root"               # 数据库用户名
  password: "password"           # 数据库密码
  charset: "utf8mb4"             # 字符集

# 服务器配置
server:
  addr: "0.0.0.0"                # 监听地址
  port: 8080                     # 监听端口

# AI 服务配置
aiEndpoint:
  url: "https://api.openai.com"  # AI API 地址
  apiKey: "your-api-key"         # API 密钥
  model: "gpt-3.5-turbo"         # 使用的模型
```

**注意：** 请勿将包含敏感信息的 `config.yaml` 提交到版本控制系统。

## 项目结构

```
BuyCar/
├── biz/                        # 业务逻辑层
│   ├── dal/                    # 数据访问层
│   │   ├── db/                 # 数据库操作
│   │   │   ├── admin.go        # 管理员相关
│   │   │   ├── consult.go      # 咨询相关
│   │   │   ├── model.go        # 模型相关
│   │   │   ├── score.go        # 评分相关
│   │   │   └── user.go         # 用户相关
│   │   └── init.go             # 数据层初始化
│   ├── handler/                # HTTP 处理器
│   │   ├── admin/              # 管理员接口
│   │   ├── purchase/           # 购买咨询接口
│   │   ├── score/              # 评分接口
│   │   └── user/               # 用户接口
│   ├── middleware/             # 中间件
│   │   └── jwt.go              # JWT 认证中间件
│   ├── model/                  # 数据模型
│   │   ├── admin/              # 管理员模型
│   │   ├── module/             # 模块模型
│   │   ├── purchase/           # 购买咨询模型
│   │   ├── score/              # 评分模型
│   │   └── user/               # 用户模型
│   ├── pack/                   # 数据打包
│   │   ├── pack.go             # 响应打包
│   │   └── user.go             # 用户数据打包
│   ├── router/                 # 路由
│   │   ├── admin/              # 管理员路由
│   │   ├── auth/               # 认证路由
│   │   ├── module/             # 模块路由
│   │   ├── purchase/           # 购买咨询路由
│   │   ├── score/              # 评分路由
│   │   ├── user/               # 用户路由
│   │   └── register.go         # 路由注册
│   └── service/                # 业务服务层
│       ├── admin.go            # 管理员服务
│       ├── purchase.go         # 购买咨询服务
│       ├── score.go            # 评分服务
│       ├── service.go          # 服务基础
│       └── user.go             # 用户服务
├── config/                     # 配置文件
│   ├── config.go               # 配置加载逻辑
│   ├── type.go                 # 配置类型定义
│   ├── config.yaml             # 配置文件（不提交）
│   └── config.example.yaml     # 配置文件示例
├── pkg/                        # 公共包
│   ├── AIAgent/                # AI 服务集成
│   │   └── openAI.go           # OpenAI 集成
│   ├── constants/              # 常量定义
│   │   ├── db.go               # 数据库常量
│   │   ├── jwt.go              # JWT 常量
│   │   └── service.go          # 服务常量
│   ├── errno/                  # 错误码定义
│   │   ├── code.go             # 错误码
│   │   ├── default.go          # 默认错误
│   │   └── errno.go            # 错误处理
│   └── utils/                  # 工具函数
│       ├── encrypt.go          # 加密工具
│       └── utils.go            # 通用工具
├── idl/                        # 接口定义语言文件
├── main.go                     # 程序入口
├── router.go                   # 路由定义
├── router_gen.go               # 生成的路由代码
├── go.mod                      # Go 模块定义
└── go.sum                      # Go 依赖校验

```

## 核心功能

### 1. 用户管理
- 用户注册与登录
- JWT 认证与授权
- 用户信息管理
- 用户反馈功能

### 2. 购买咨询
- AI 智能购车咨询
- 咨询记录管理
- 个性化推荐

### 3. 评分系统
- 用户评分功能
- 评分统计与分析

### 4. 管理员功能
- 后台管理
- 数据审核
- 权限管理

## API 接口

### 用户相关
- `POST /api/user/register` - 用户注册
- `POST /api/user/login` - 用户登录
- `POST /api/user/feedback` - 用户反馈

### 认证相关
- JWT Token 刷新机制
- Access Token 和 Refresh Token 双令牌体系

## 中间件

项目使用以下中间件：
- **JWT 认证中间件** - 保护需要认证的接口
- **CORS 中间件** - 处理跨域请求
- **日志中间件** - 记录请求日志

## 数据库设计

项目使用 GORM 进行数据库操作，支持自动迁移。主要表结构包括：
- `users` - 用户表
- `admins` - 管理员表
- `consults` - 咨询记录表
- `scores` - 评分表
- 其他业务相关表

## 开发说明

### 添加新的 API

1. 在 `biz/model/` 下定义请求和响应结构
2. 在 `biz/dal/db/` 下实现数据库操作
3. 在 `biz/service/` 下实现业务逻辑
4. 在 `biz/handler/` 下实现 HTTP 处理器
5. 在 `biz/router/` 下注册路由

### 错误处理

项目使用统一的错误处理机制，错误码定义在 `pkg/errno/` 目录下。

## 安全性

- 密码使用加密存储
- JWT Token 认证机制
- 请求参数验证
- SQL 注入防护（通过 GORM）

## 贡献指南

欢迎提交 Issue 和 Pull Request！

## 许可证

[请添加适当的许可证信息]

## 联系方式

[请添加联系方式]
