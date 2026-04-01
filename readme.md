# OPS 运营管理系统

> Operations（运营）的缩写，一个前后端分离的工作流/运营管理系统。

## 项目简介

OPS 是一个功能完善的运营管理系统，支持用户认证、文件管理、采购订单管理、日程排班等功能。前端采用 Vue 3 + Tailwind CSS，后端采用 Go + Gin 框架。

## 技术栈

### 后端
- **语言**: Go
- **框架**: Gin (HTTP) + GORM (ORM)
- **数据库**: SQLite / MySQL / PostgreSQL (通过配置切换)
- **日志**: Uber Zap

### 前端
- **框架**: Vue 3 (JavaScript)
- **构建工具**: Vite 7
- **CSS**: Tailwind CSS v4
- **状态管理**: Pinia
- **国际化**: vue-i18n
- **图标**: @tabler/icons-vue
- **组件库**: FullCalendar, CropperJS, FilePond, flatpickr

## 项目结构

```
ops2/
├── backend/           # Go 后端
│   ├── api/           # API 路由定义
│   ├── internal/      # 内部包
│   │   ├── config/    # 配置管理
│   │   ├── database/ # 数据库连接
│   │   ├── handler/   # HTTP 处理器
│   │   ├── middleware/# 中间件
│   │   ├── repository/# 数据访问层
│   │   └── service/  # 业务逻辑层
│   ├── models/        # 数据模型
│   ├── routers/       # 兼容性路由
│   ├── pkg/           # 公共包
│   └── main.go        # 应用入口
│
├── frontend/
│   └── ops_vue_js/   # Vue 3 前端（主力开发目录）
│
└── DOC/              # 项目文档
```

## 功能模块

- **用户认证**: 登录、注册、密码找回、Session 管理
- **文件管理**: 图片/文件上传、下载、头像裁剪
- **采购订单**: 订单创建、列表查看、费用明细
- **日程排班**: FullCalendar 日历视图
- **仓库管理**: 库存管理功能

## 快速开始

### 前置要求

- Go 1.21+
- Node.js 18+
- GCC (TDM-GCC for Windows)

### 后端启动

```bash
fresh
```

服务默认运行在 http://localhost:8080

### 前端启动

```bash
cd frontend/ops_vue_js
npm install
npm run dev
```

## API 路由

### RESTful API v1
- `GET    /api/v1/`          - API 根路径
- `POST   /api/v1/users/login`      - 用户登录
- `POST   /api/v1/users/register`   - 用户注册
- `GET    /api/v1/users/profile`    - 获取用户信息
- `POST   /api/v1/files/upload`     - 文件上传
- `GET    /api/v1/files/list`       - 文件列表
- `POST   /api/v1/purchase/orders`  - 创建采购订单
- `GET    /api/v1/purchase/orders`  - 获取采购订单列表

### 兼容性 API
- `/api/users/*`    - 用户相关
- `/api/files/*`    - 文件相关
- `/api/purchase/*` - 采购相关

## 配置

后端配置文件位于 `./backend/data/config.yaml`:

```yaml
server:
  port: 8080
  mode: debug  # debug / release

database:
  type: sqlite
  path: ./data/ops.db
```

## 开发规范

### 后端架构 (分层)
```
Handler → Service → Repository → Database
```

### 认证方式
- 请求体中携带 `userCookieValue` 字段
- 或使用 Authorization Bearer Token

### 响应格式
```json
{
  "code": "0",        // 错误码，0 表示成功
  "message": "Success",
  "data": {}
}
```

### 错误码
- `0`: 成功
- `-1`: 内部错误
- `-2`: 参数错误
- `-3`: 未登录
- `-4`: 用户已存在
- `-5`: 用户不存在
- `-42`: 凭证错误

## 许可证

MIT License - 见 LICENSE 文件
