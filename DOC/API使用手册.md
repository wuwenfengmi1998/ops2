# OPS API v1 使用手册

## 概述

OPS (Operations 运营管理系统) 是一个前后端分离的工作流/运营管理系统。本文档详细描述了后端API的接口规范和使用方法。

**当前版本**: v1  
**基础URL**: `http://localhost:8080/api/v1/`  
**响应格式**: JSON

---

## 目录

1. [快速开始](#快速开始)
2. [认证系统](#认证系统)
3. [用户管理](#用户管理)
4. [文件管理](#文件管理)
5. [采购管理](#采购管理)
6. [系统管理](#系统管理)
7. [错误码说明](#错误码说明)
8. [请求示例](#请求示例)

---

## 快速开始

### 安装与运行

```bash
# 进入后端目录
cd backend

# 安装依赖
go mod download

# 运行服务器
go run cmd/ops-server/main.go

# 或使用编译版本
go build -o ops-server cmd/ops-server/main.go
./ops-server
```

### 环境配置

默认配置位于 `./data/config.yaml`，模板配置在 `./defConfig/configTemp.yaml`

```yaml
# 默认配置示例
server:
  host: "localhost"
  port: 8080
  tls_enable: false

database:
  driver: "sqlite"
  path: "./data/db.db"

file:
  paths:
    image: "./data/images"
  max_size: 5242880  # 5MB
```

### 首次使用

1. 启动后端服务器
2. 访问前端界面：`http://localhost:8080`
3. 注册新用户或使用默认账户
4. 开始使用API

---

## 认证系统

OPS系统支持多种认证方式，确保与现有前端的兼容性。

### 认证方式

**1. Authorization Header (推荐)**
```http
Authorization: Bearer <token>
```

**2. Cookie Header**
```http
Cookie: ops_session=<session_token>
```

**3. POST Body (兼容现有前端)**
```json
{
  "userCookieValue": "<session_token>",
  "otherField": "value"
}
```

### 响应格式

所有API响应都遵循统一格式：

```json
{
  "code": "0",           // 错误码，0表示成功
  "message": "Success",  // 人类可读的消息
  "data": {}            // 实际数据，根据接口不同而变化
}
```

---

## 用户管理

### 用户注册

**注册新用户**

```http
POST /api/v1/users/register
Content-Type: application/json
```

**请求体**:
```json
{
  "name": "username",
  "password": "password123",
  "password_confirm": "password123",
  "email": "user@example.com"
}
```

**响应**:
```json
{
  "code": "0",
  "message": "注册成功",
  "data": {
    "user_id": 1,
    "name": "username",
    "email": "user@example.com"
  }
}
```

### 用户登录

**用户登录获取令牌**

```http
POST /api/v1/users/login
Content-Type: application/json
```

**请求体**:
```json
{
  "name": "username",
  "password": "password123"
}
```

**响应**:
```json
{
  "code": "0",
  "message": "登录成功",
  "data": {
    "user_id": 1,
    "name": "username",
    "cookie_value": "session_token_here",
    "expires_at": "2026-04-01T10:00:00Z"
  }
}
```

### 忘记密码

**请求密码重置**

```http
POST /api/v1/users/forgot-password
Content-Type: application/json
```

**请求体**:
```json
{
  "email": "user@example.com"
}
```

**响应**:
```json
{
  "code": "0",
  "message": "重置邮件已发送",
  "data": null
}
```

### 重置密码

**使用重置令牌设置新密码**

```http
POST /api/v1/users/reset-password
Content-Type: application/json
```

**请求体**:
```json
{
  "token": "reset_token",
  "new_password": "newpassword123",
  "new_password_confirm": "newpassword123"
}
```

### 获取用户资料

**获取当前用户信息**

```http
GET /api/v1/users/profile
Authorization: Bearer <token>
```

**响应**:
```json
{
  "code": "0",
  "message": "Success",
  "data": {
    "user_id": 1,
    "name": "username",
    "email": "user@example.com",
    "avatar": "/files/get/avatar_hash",
    "gender": "male",
    "language": "zh-CN",
    "created_at": "2026-03-31T10:00:00Z",
    "last_login": "2026-03-31T18:00:00Z"
  }
}
```

### 更新用户资料

**更新用户个人信息**

```http
PUT /api/v1/users/profile
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**:
```json
{
  "email": "newemail@example.com",
  "avatar_hash": "new_avatar_hash",
  "gender": "female",
  "language": "en-US"
}
```

### 用户登出

**注销当前会话**

```http
POST /api/v1/users/logout
Authorization: Bearer <token>
```

---

## 文件管理

### 上传文件

**上传图片文件**

```http
POST /api/v1/files/upload
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**表单字段**:
- `file`: 文件内容 (图片文件，支持PNG、JPEG、GIF、WebP)
- `type`: 文件类型 (可选，默认: "image")
- `description`: 文件描述 (可选)

**响应**:
```json
{
  "code": "0",
  "message": "文件上传成功",
  "data": {
    "file_id": 1,
    "name": "example.png",
    "sha256": "abc123...",
    "mime": "image/png",
    "size": 123456,
    "download_url": "/api/v1/files/download/abc123...",
    "preview_url": "/api/v1/files/get/abc123...",
    "created_at": "2026-03-31T18:15:00Z"
  }
}
```

### 获取文件列表

**获取用户上传的文件列表**

```http
GET /api/v1/files/list
Authorization: Bearer <token>
```

**查询参数**:
- `type`: 文件类型过滤 (可选)
- `page`: 页码 (默认: 1)
- `entries`: 每页数量 (默认: 20, 最大: 100)

**响应**:
```json
{
  "code": "0",
  "message": "Success",
  "data": {
    "files": [
      {
        "file_id": 1,
        "name": "example.png",
        "sha256": "abc123...",
        "mime": "image/png",
        "size": 123456,
        "type": "image",
        "created_at": "2026-03-31T18:15:00Z"
      }
    ],
    "total": 10,
    "page": 1,
    "pages": 1
  }
}
```

### 获取文件信息

**获取单个文件详情**

```http
GET /api/v1/files/{file_id}
Authorization: Bearer <token>
```

**路径参数**:
- `file_id`: 文件ID

### 下载文件

**下载文件内容**

```http
GET /api/v1/files/download/{hash}
```

**路径参数**:
- `hash`: 文件SHA256哈希值

### 预览文件

**预览文件（浏览器直接显示）**

```http
GET /api/v1/files/get/{hash}
```

**路径参数**:
- `hash`: 文件SHA256哈希值

### 删除文件

**删除用户文件**

```http
DELETE /api/v1/files/{file_id}
Authorization: Bearer <token>
```

---

## 采购管理

### 获取采购订单列表

**方式1: POST方式（兼容现有前端）**

```http
POST /api/v1/purchase/getorders
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**:
```json
{
  "search": "keyword",
  "page": 1,
  "entries": 20
}
```

**方式2: GET方式（RESTful API）**

```http
GET /api/v1/purchase/orders
Authorization: Bearer <token>
```

**查询参数**:
- `search`: 搜索关键词 (可选)
- `page`: 页码 (默认: 1)
- `entries`: 每页数量 (默认: 20, 最大: 300)

**响应**:
```json
{
  "code": "0",
  "message": "Success",
  "data": {
    "all_count": 150,
    "all_orders": [
      {
        "id": 1,
        "user_id": 1,
        "title": "服务器硬件采购",
        "part_name": "服务器",
        "order_status": "pending",
        "tracking_number": "TN123456789",
        "photos": ["hash1", "hash2"],
        "created_at": "2026-03-30T10:00:00Z",
        "update_time": "2026-03-31T15:00:00Z"
      }
    ]
  }
}
```

### 创建采购订单

**方式1: POST方式（兼容现有前端）**

```http
POST /api/v1/purchase/addorder
Authorization: Bearer <token>
Content-Type: application/json
```

**方式2: POST方式（RESTful API）**

```http
POST /api/v1/purchase/orders
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**:
```json
{
  "costs": [
    {
      "cost": 1200,
      "costt": 1200,
      "currencytype": "CNY",
      "int": 2,
      "type": "服务器"
    },
    {
      "cost": 300,
      "costt": 300,
      "currencytype": "CNY",
      "int": 1,
      "type": "内存条"
    }
  ],
  "title": "服务器硬件采购",
  "order_status": "pending",
  "part_name": "服务器配件",
  "remark": "需要尽快发货",
  "link": "https://example.com/product/123",
  "photos": ["image_hash_1", "image_hash_2"],
  "styles": "{\"priority\": \"high\"}",
  "tracking_number": "TN123456789",
  "update_time": "2026-03-31T18:00:00"
}
```

**响应**:
```json
{
  "code": "0",
  "message": "订单创建成功",
  "data": {
    "order_id": 1,
    "total_cost": 2700,
    "created_at": "2026-03-31T18:17:00Z"
  }
}
```

### 获取订单详情

**获取单个订单的详细信息**

```http
GET /api/v1/purchase/orders/{order_id}
Authorization: Bearer <token>
```

**路径参数**:
- `order_id`: 订单ID

**响应**:
```json
{
  "code": "0",
  "message": "Success",
  "data": {
    "order": {
      "id": 1,
      "user_id": 1,
      "title": "服务器硬件采购",
      "remark": "需要尽快发货",
      "photos": ["hash1", "hash2"],
      "link": "https://example.com/product/123",
      "part_name": "服务器配件",
      "styles": "{\"priority\": \"high\"}",
      "update_time": "2026-03-31T18:00:00Z",
      "tracking_number": "TN123456789",
      "order_status": "pending",
      "created_at": "2026-03-31T18:17:00Z"
    },
    "costs": [
      {
        "id": 1,
        "order_id": 1,
        "user_id": 1,
        "price": 1200,
        "quantity": 2,
        "created_at": "2026-03-31T18:17:00Z"
      },
      {
        "id": 2,
        "order_id": 1,
        "user_id": 1,
        "price": 300,
        "quantity": 1,
        "created_at": "2026-03-31T18:17:00Z"
      }
    ]
  }
}
```

---

## 系统管理

### 系统状态

**获取系统运行状态**

```http
GET /api/v1/system/status
```

**响应**:
```json
{
  "code": "0",
  "message": "系统运行正常",
  "data": {
    "version": "1.0.0",
    "uptime": "10h30m",
    "database": "connected",
    "memory_usage": "45%",
    "active_sessions": 5,
    "total_users": 50,
    "server_time": "2026-03-31T18:17:30Z"
  }
}
```

### 获取系统配置

**获取系统配置（需要管理员权限）**

```http
GET /api/v1/system/config
Authorization: Bearer <token>
```

### 更新系统配置

**更新系统配置（需要管理员权限）**

```http
PUT /api/v1/system/config
Authorization: Bearer <token>
Content-Type: application/json
```

---

## 错误码说明

### 核心错误码

| 错误码 | 含义 | HTTP状态码 |
|--------|------|-----------|
| `0`    | 成功 | 200 |
| `-1`   | 内部服务器错误 | 500 |
| `-2`   | 参数错误 | 400 |
| `-3`   | 用户未登录 | 401 |
| `-4`   | 用户已存在 | 409 |
| `-5`   | 用户不存在 | 404 |
| `-42`  | 用户名或密码错误 | 401 |

### 文件相关错误码

| 错误码 | 含义 | HTTP状态码 |
|--------|------|-----------|
| `-100` | 文件不存在 | 404 |
| `-101` | 文件类型不支持 | 415 |
| `-102` | 文件大小超出限制 | 413 |
| `-103` | 文件上传失败 | 500 |
| `-104` | 文件哈希计算失败 | 500 |

### 采购订单相关错误码

| 错误码 | 含义 | HTTP状态码 |
|--------|------|-----------|
| `-200` | 订单不存在 | 404 |
| `-201` | 订单创建失败 | 500 |
| `-202` | 费用明细错误 | 400 |
| `-203` | 订单状态无效 | 400 |

---

## 请求示例

### 使用cURL

**用户登录**:
```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"name": "admin", "password": "password123"}'
```

**获取采购订单**:
```bash
curl -X GET "http://localhost:8080/api/v1/purchase/orders?page=1&entries=20" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**上传文件**:
```bash
curl -X POST http://localhost:8080/api/v1/files/upload \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -F "file=@/path/to/image.png" \
  -F "type=image" \
  -F "description=产品图片"
```

### 使用JavaScript (Fetch API)

**用户注册**:
```javascript
async function registerUser(username, password, email) {
  const response = await fetch('http://localhost:8080/api/v1/users/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      name: username,
      password: password,
      password_confirm: password,
      email: email
    })
  });
  
  const result = await response.json();
  
  if (result.code === '0') {
    console.log('注册成功:', result.data);
    return result.data.cookie_value;
  } else {
    console.error('注册失败:', result.message);
    throw new Error(result.message);
  }
}
```

**创建采购订单**:
```javascript
async function createPurchaseOrder(token, orderData) {
  const response = await fetch('http://localhost:8080/api/v1/purchase/orders', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify(orderData)
  });
  
  const result = await response.json();
  
  if (result.code === '0') {
    console.log('订单创建成功:', result.data);
    return result.data.order_id;
  } else {
    console.error('订单创建失败:', result.message);
    throw new Error(result.message);
  }
}
```

### 使用Python (Requests)

**获取用户资料**:
```python
import requests

def get_user_profile(token):
    url = "http://localhost:8080/api/v1/users/profile"
    headers = {
        "Authorization": f"Bearer {token}"
    }
    
    response = requests.get(url, headers=headers)
    result = response.json()
    
    if result["code"] == "0":
        return result["data"]
    else:
        raise Exception(result["message"])

# 使用示例
try:
    token = "your_token_here"
    profile = get_user_profile(token)
    print(f"用户ID: {profile['user_id']}")
    print(f"用户名: {profile['name']}")
except Exception as e:
    print(f"错误: {e}")
```

---

## 附录

### 数据模型说明

**用户表 (TabUser_)**
```go
type TabUser_ struct {
    ID           uint   `gorm:"primarykey"`      // 用户ID
    Name         string `gorm:"unique"`          // 用户名（唯一）
    PasswordHash string                          // 密码哈希
    Email        string                          // 邮箱
    CreatedAt    time.Time                       // 创建时间
}
```

**文件信息表 (TabFileInfo_)**
```go
type TabFileInfo_ struct {
    ID     uint      `gorm:"primaryKey"`       // 文件ID
    Name   string    `gorm:"not null"`         // 文件名
    Path   string    `gorm:"not null"`         // 文件路径
    Sha256 string    `gorm:"not null;index"`   // SHA256哈希
    Mime   string    `gorm:"index"`            // MIME类型
    Type   string    `gorm:"index"`            // 文件类型
    UserID uint      `gorm:"not null;index"`   // 用户ID
    Date   time.Time                           // 上传时间
}
```

**采购订单表 (TabPurchaseOrder)**
```go
type TabPurchaseOrder struct {
    ID             uint           `gorm:"primarykey"`          // 订单ID
    UserID         uint           `gorm:"not null"`            // 用户ID
    Title          string         `gorm:"size:200"`            // 标题
    Remark         string         `gorm:"type:text"`           // 备注
    Photos         datatypes.JSON `gorm:"type:json"`           // 照片哈希数组
    Link           string         `gorm:"size:1000"`           // 链接
    PartName       string         `gorm:"size:200;not null"`   // 物品名称
    Styles         string         `gorm:"type:text"`           // 样式数组
    TrackingNumber string         `gorm:"size:100;Index"`      // 快递单号
    OrderStatus    string         `gorm:"default:1"`           // 订单状态
    CreatedAt      *time.Time     `gorm:"type:datetime"`       // 创建时间
    UpdateTime     *time.Time     `gorm:"type:datetime"`       // 更新时间
}
```

### 版本历史

| 版本 | 日期 | 变更说明 |
|------|------|----------|
| v1.0 | 2026-03-31 | 初始版本发布，包含完整的API文档 |
| v1.1 | 2026-04-01 | 添加文件管理API，优化错误处理 |
| v1.2 | 2026-04-02 | 增加采购订单管理功能 |

### 支持与反馈

如需技术支持或有任何建议，请联系：
- **GitHub仓库**: https://github.com/yourusername/ops
- **邮箱**: support@example.com
- **文档更新**: 定期查看本文档获取最新API信息

---

*文档最后更新: 2026-03-31*
*OPS系统开发团队 版权所有 © 2026*