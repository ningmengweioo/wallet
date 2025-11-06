# Cex 钱包后端服务

## 项目介绍

这是一个基于 Go 语言开发的 Cex 钱包后端服务，提供用户管理、钱包管理、资金转账等功能。

## 技术栈

- 编程语言：Go
- Web 框架：Gin
- 数据库：MySQL
- 配置管理：YAML


## 文件结构

```
├── Read.md           # 项目说明文档
├── config/           # 配置相关
│   ├── config.go     # 配置结构和加载逻辑
│   ├── config.yaml   # 配置文件
│   ├── db.go         # 数据库连接配置
│   └── logger.go     # 日志配置
├── controller/       # 控制器层
│   └── WalletController.go # 钱包相关控制器
├── go.mod            # Go 模块文件
├── go.sum            # Go 依赖校验文件
├── main.go           # 应用入口
├── models/           # 数据模型
│   ├── transaction.go # 交易记录模型
│   ├── users.go      # 用户模型
│   └── wallets.go    # 钱包模型
├── question.md       # 问题记录
├── router/           # 路由配置
│   └── router.go     # 路由设置
├── service/          # 业务逻辑层
│   ├── transaction.go # 交易相关业务逻辑
│   ├── user.go       # 用户相关业务逻辑
│   └── wallet.go     # 钱包相关业务逻辑
├── test/             # 测试目录
│   └── api_test.go   # API 测试文件
└── utils/            # 工具函数
    └── response.go   # 响应处理工具
```

## 主要功能模块

### 1. 用户管理
- 用户注册
- 获取所有用户
- 获取单个用户详情

### 2. 钱包管理
- 查询余额
- 存款
- 取款
- 转账

### 3. 交易记录
- 查询用户交易历史

## 数据库设计

### 用户表 (users)
- id: 主键，自增长
- username: 用户名
- email: 邮箱，唯一索引
- created_at: 创建时间
- updated_at: 更新时间
- deleted_at: 软删除时间

### 钱包表 (wallets)
- id: 主键，自增长
- user_id: 用户ID，外键，唯一索引
- balance: 余额，默认0
- created_at: 创建时间
- updated_at: 更新时间
- deleted_at: 软删除时间

### 交易记录表 (transactions)
- 包含交易ID、用户ID、交易类型、金额、状态等字段

## API 接口

### 健康检查
- GET /health - 检查API服务是否正常运行

### 用户相关接口
- POST /api/v1/users - 注册用户
- GET /api/v1/users - 获取所有用户
- GET /api/v1/users/:id - 获取用户详情

### 钱包相关接口
- GET /api/v1/wallets/:user_id/balance - 查询余额
- POST /api/v1/wallets/:user_id/deposit - 存款
- POST /api/v1/wallets/:user_id/withdraw - 取款
- POST /api/v1/wallets/transfer - 转账

### 交易记录接口
- GET /api/v1/transactions/:user_id - 获取用户交易记录

## 部署说明



### 配置文件设置

修改 `config/config.yaml` 文件，设置数据库连接信息和服务端口：

```yaml
http:
  port: 8090

mysql:
  host: localhost
  port: 3306
  db_name: wallet
  user: root
  password: your_password
  charset: utf8mb4

log:
  level: info
```

### 环境变量

可以通过环境变量覆盖配置文件中的设置：
- `CONFIG_PATH`: 配置文件路径，默认为 `./config/config.yaml`

### 启动服务

1. 安装依赖：
```bash
go mod tidy
```

2. 启动服务：
```bash
go run main.go
```

服务将在配置的端口上启动，默认端口为 8090。

## 测试

项目包含 API 测试，可以通过以下命令运行：

```bash
go test ./test -v
```