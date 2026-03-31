# OPS 后端重构架构设计文档

## 当前状态
✅ 已完成基础架构重构
✅ 新目录结构已创建
✅ 配置管理模块完成
✅ 数据库连接层完成
✅ 响应统一处理完成

## 目录结构

```
backend/
├── cmd/                       # 入口点
│   └── ops-server/           # 主应用程序入口
│       └── main.go
├── internal/                 # 私有应用程序代码
│   ├── config/               # 配置管理
│   │   ├── config.go         # 配置结构定义
│   │   └── [其他配置组件]
│   ├── database/             # 数据库层
│   │   ├── connection.go     # 数据库连接
│   │   └── migration.go      # 数据库迁移和模型定义
│   ├── models/               # 数据模型（待重构）
│   ├── repository/           # 数据访问层（待创建）
│   │   ├── user_repository.go
│   │   ├── purchase_repository.go
│   │   └── file_repository.go
│   ├── service/              # 业务逻辑层（待创建）
│   │   ├── auth_service.go
│   │   ├── purchase_service.go
│   │   └── file_service.go
│   ├── handler/              # HTTP处理器（待创建）
│   │   ├── auth_handler.go
│   │   ├── purchase_handler.go
│   │   └── file_handler.go
│   ├── middleware/           # 中间件（待创建）
│   │   ├── auth.go
│   │   ├── logging.go
│   │   └── recovery.go
│   └── pkg/                  # 内部公共库（待创建）
│       ├── errors/
│       ├── validation/
│       └── utils/
├── api/                      # API定义（待创建）
│   └── v1/
│       └── routes.go
├── pkg/                      # 公共库
│   └── response/
│       └── response.go       # API响应统一处理
├── data/                     # 数据目录（配置文件、数据库）
├── defConfig/                # 默认配置模板
├── dist/                     # 前端构建产物（编译后）
└── tests/                    # 测试文件
```

## 主要改进

### 1. 配置管理重构
**旧方案问题：**
- 使用全局变量 `var Configs map[string]interface{}`
- 使用奇怪的命名如 `ConfigsWed`, `ConfigsFile`
- 拼写错误：`Pahts` 应该是 `Paths`
- 缺少类型安全和验证

**新方案：**
- 使用结构体定义配置，支持类型安全
- 支持默认配置自动生成
- 支持热重载（未来扩展）
- 统一配置路径管理

### 2. 数据库层重构
**旧方案问题：**
- 直接在路由层进行数据库操作
- 缺少连接池配置
- 错误处理不一致

**新方案：**
- 集中管理数据库连接
- 自动连接池配置
- 支持多种数据库（SQLite/MySQL/PostgreSQL）
- 统一错误处理

### 3. 错误处理统一
**旧方案问题：**
- 混合使用 panic、return error 和日志
- HTTP 响应格式不一致
- 错误码管理混乱

**新方案：**
- 统一 API 响应格式
- 标准错误码映射
- 结构化错误信息
- 详细的 HTTP 状态码

### 4. 中间件系统
**新功能：**
- CORS 跨域支持
- 请求日志记录
- 认证中间件
- 性能监控
- 限流保护

## 迁移步骤

### 第一阶段：基础架构 ✅
- [x] 创建新的目录结构
- [x] 重构配置管理模块
- [x] 重构数据库连接层
- [x] 创建统一响应处理

### 第二阶段：数据访问层
- [ ] 创建 Repository 层
- [ ] 迁移用户相关数据访问
- [ ] 迁移采购订单数据访问  
- [ ] 迁移文件管理数据访问

### 第三阶段：业务逻辑层
- [ ] 创建 Service 层
- [ ] 迁移用户认证逻辑
- [ ] 迁移采购订单管理逻辑
- [ ] 迁移文件上传逻辑

### 第四阶段：HTTP 层
- [ ] 创建 Handler 层
- [ ] 迁移用户认证 API
- [ ] 迁移采购订单 API
- [ ] 迁移文件管理 API

### 第五阶段：中间件和测试
- [ ] 创建中间件系统
- [ ] 添加单元测试
- [ ] 添加集成测试
- [ ] 性能测试

## 运行说明

### 准备步骤
1. 运行配置迁移脚本：
```bash
python migrate-config.py
```

2. 检查前端构建：
```bash
# 确保前端已构建，在frontend目录中运行
npm run build
# 或手动复制dist目录到backend/dist
```

3. 启动新版本服务器：
```bash
# Windows
.\start-dev.bat

# Linux/Mac
go run ./cmd/ops-server/main.go
```

### 迁移过程中的注意事项
1. **保持向后兼容**：逐步迁移，确保 API 接口不变
2. **数据库数据安全**：旧数据库文件会自动迁移
3. **配置文件备份**：迁移前自动备份旧配置
4. **增量测试**：每次迁移后测试新功能

## API 兼容性保证

### 保持不变的接口
- 用户登录：`POST /api/users/login`
- 用户注册：`POST /api/users/register`
- 获取采购订单：`GET /api/purchase/orders`
- 上传文件：`POST /api/files/upload`

### 响应的格式统一
```json
{
  "code": "0",        // 错误码，0表示成功
  "message": "Success", // 人类可读的消息
  "data": {}          // 实际数据
}
```

## 性能改进目标
1. **响应时间**：平均 < 200ms
2. **并发连接**：支持 1000+ 并发
3. **内存使用**：< 200MB
4. **启动时间**：< 5s

## 监控和日志
- 结构化日志输出
- 请求追踪ID
- 性能指标收集
- 错误聚合报告