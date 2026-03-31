# OPS API 快速参考

## 基础信息

**Base URL**: `http://localhost:8080/api/v1/`  
**认证方式**: Bearer Token, Cookie, 或POST body `userCookieValue`  
**响应格式**: JSON  

## 认证端点

| 方法 | 端点 | 描述 | 需要认证 |
|------|------|------|----------|
| POST | `/users/register` | 用户注册 | ❌ |
| POST | `/users/login` | 用户登录 | ❌ |
| POST | `/users/forgot-password` | 忘记密码 | ❌ |
| POST | `/users/reset-password` | 重置密码 | ❌ |
| GET | `/users/profile` | 获取用户资料 | ✅ |
| PUT | `/users/profile` | 更新用户资料 | ✅ |
| POST | `/users/logout` | 用户登出 | ✅ |

## 文件管理端点

| 方法 | 端点 | 描述 | 需要认证 |
|------|------|------|----------|
| POST | `/files/upload` | 上传文件 | ✅ |
| GET | `/files/list` | 获取文件列表 | ✅ |
| GET | `/files/{id}` | 获取文件信息 | ✅ |
| DELETE | `/files/{id}` | 删除文件 | ✅ |
| GET | `/files/download/{hash}` | 下载文件 | ❌ |
| GET | `/files/get/{hash}` | 预览文件 | ❌ |

## 采购订单端点

### 兼容前端API
| 方法 | 端点 | 描述 | 需要认证 |
|------|------|------|----------|
| POST | `/purchase/getorders` | 获取订单列表 | ✅ |
| POST | `/purchase/addorder` | 创建订单 | ✅ |

### RESTful API
| 方法 | 端点 | 描述 | 需要认证 |
|------|------|------|----------|
| GET | `/purchase/orders` | 获取订单列表 | ✅ |
| POST | `/purchase/orders` | 创建订单 | ✅ |
| GET | `/purchase/orders/{id}` | 获取订单详情 | ✅ |

## 系统管理端点

| 方法 | 端点 | 描述 | 需要认证 | 需要管理员 |
|------|------|------|----------|------------|
| GET | `/system/status` | 系统状态 | ❌ | ❌ |
| GET | `/system/config` | 获取配置 | ✅ | ✅ |
| PUT | `/system/config` | 更新配置 | ✅ | ✅ |

## 请求头示例

### Bearer Token
```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Cookie
```http
Cookie: ops_session=abc123def456
```

### Content-Type
```http
Content-Type: application/json
```

## 响应格式

```json
{
  "code": "0",           // 错误码
  "message": "Success",  // 消息
  "data": {}            // 数据
}
```

## 核心错误码

| 错误码 | 含义 | 说明 |
|--------|------|------|
| `0` | 成功 | 请求成功完成 |
| `-1` | 内部错误 | 服务器内部错误 |
| `-2` | 参数错误 | 请求参数不正确 |
| `-3` | 未登录 | 用户未登录或令牌无效 |
| `-4` | 用户已存在 | 注册时用户名已存在 |
| `-5` | 用户不存在 | 用户记录不存在 |
| `-42` | 凭证错误 | 用户名或密码错误 |

## 分页参数

所有列表接口支持分页：
- `page`: 页码 (默认: 1)
- `entries`: 每页数量 (默认: 20)

```http
GET /api/v1/purchase/orders?page=2&entries=50
```

## 文件上传

支持的文件类型：
- 图片: PNG, JPEG, GIF, WebP
- 最大文件大小: 5MB

```http
POST /api/v1/files/upload
Content-Type: multipart/form-data
```

## 采购订单数据结构

### 创建订单请求体
```json
{
  "title": "订单标题",
  "order_status": "pending",
  "part_name": "物品名称",
  "remark": "备注信息",
  "link": "https://example.com",
  "photos": ["hash1", "hash2"],
  "tracking_number": "TN123456789",
  "update_time": "2026-03-31T18:00:00",
  "costs": [
    {
      "cost": 1000,
      "costt": 1000,
      "currencytype": "CNY",
      "int": 2,
      "type": "类型"
    }
  ]
}
```

### 订单状态值
- `pending`: 待处理
- `processing`: 处理中
- `shipped`: 已发货
- `completed`: 已完成
- `cancelled`: 已取消

## 快速示例

### 1. 用户登录
```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"name": "admin", "password": "password123"}'
```

### 2. 获取订单列表
```bash
curl -X GET "http://localhost:8080/api/v1/purchase/orders?page=1" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 3. 上传文件
```bash
curl -X POST http://localhost:8080/api/v1/files/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@image.png"
```

## 状态码映射

| OPS错误码 | HTTP状态码 | 含义 |
|-----------|------------|------|
| `0` | 200 OK | 成功 |
| `-2` | 400 Bad Request | 参数错误 |
| `-3`, `-42` | 401 Unauthorized | 认证失败 |
| `-4` | 409 Conflict | 资源冲突 |
| `-5`, `-100` | 404 Not Found | 资源不存在 |
| `-102` | 413 Payload Too Large | 文件太大 |
| `-101` | 415 Unsupported Media Type | 文件类型不支持 |
| `-1` | 500 Internal Server Error | 服务器错误 |

## 开发环境配置

### 数据库配置
```yaml
# backend/data/config.yaml
database:
  driver: "sqlite"      # sqlite, mysql, postgres
  path: "./data/db.db"  # SQLite文件路径
  # 或MySQL配置:
  # host: "localhost"
  # port: 3306
  # name: "ops"
  # user: "root"
  # password: "password"
```

### 启动命令
```bash
# 开发模式
go run cmd/ops-server/main.go

# 生产模式
go build -o ops-server cmd/ops-server/main.go
./ops-server
```

## 前端集成示例

### 存储令牌
```javascript
// 登录后存储令牌
localStorage.setItem('ops_token', response.data.cookie_value);

// 后续请求自动附加令牌
const token = localStorage.getItem('ops_token');
fetch('/api/v1/users/profile', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});
```

### 处理响应
```javascript
async function apiRequest(endpoint, options = {}) {
  const token = localStorage.getItem('ops_token');
  
  const response = await fetch(`/api/v1${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
      ...options.headers
    }
  });
  
  const result = await response.json();
  
  if (result.code === '0') {
    return result.data;
  } else {
    throw new Error(result.message);
  }
}
```

## 故障排除

### 常见问题

1. **401 Unauthorized**
   - 检查令牌是否过期
   - 确保请求头格式正确
   - 尝试重新登录获取新令牌

2. **400 Bad Request**
   - 检查请求体JSON格式
   - 验证必填字段是否提供
   - 检查字段类型和值范围

3. **500 Internal Server Error**
   - 查看服务器日志
   - 检查数据库连接
   - 验证文件权限

### 日志位置
- 后端日志: 控制台输出或日志文件
- 访问日志: `backend/logs/access.log`
- 错误日志: `backend/logs/error.log`

---

*文档版本: v1.0.0*  
*最后更新: 2026-03-31*  
*保持联系: support@example.com*