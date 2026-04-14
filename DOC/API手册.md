# OPS 后端 API 使用手册

## 概述

- **基础URL**: `/api`
- **请求格式**: JSON Body（通过 `userCookieValue` 传cookie，`data` 传业务数据）
- **认证方式**: Cookie-based authentication
- **响应格式**: JSON

## 通用请求格式

```json
{
  "userCookieValue": "cookie字符串",
  "data": {
    // 业务参数
  }
}
```

## 通用响应格式

```json
{
  "err_code": 0,
  "err_msg": "apiOK",
  "return": {
    // 响应数据
  }
}
```

## 错误码

| err_code | err_msg | 说明 |
|----------|---------|------|
| 0 | apiOK | 成功 |
| -1 | apiErr | 服务器错误 |
| -2 | postErr | POST请求错误 |
| -3 | jsonErr | JSON解析错误 |
| -31 | jsonErr_1 | 数据格式错误 |
| -4 | userNameDup | 用户名已存在 |
| -41 | userNameNoFund | 用户不存在 |
| -42 | userPassIncorrect | 密码错误 |
| -43 | userEmailFormatError | 邮箱格式错误 |
| -44 | userCookieError | Cookie错误/未登录 |
| -51 | file_mime_err | 文件类型不允许 |
| -52 | file_size_err | 文件大小超限 |
| -53 | file_name_err | 文件名错误 |
| -54 | file_get_err | 文件获取错误 |
| -55 | file_hash_err | 文件哈希计算错误 |
| -56 | file_save_err | 文件保存失败 |
| -57 | file_not_found | 文件不存在 |
| -58 | file_part_err | 文件参数错误 |
| -61 | schedule_event_not_find | 日程不存在 |
| -62 | schedule_permission_denied | 无权限操作 |
| -1001 | order_not_found | 订单不存在 |
| -1002 | invalid_status | 无效的订单状态 |
| -1003 | status_no_change | 状态未变化 |
| -1004 | no_permission | 无权限 |
| -1005 | photo_hash_invalid | 图片哈希无效 |

---

## 用户模块 `/api/users`

### 用户登录

```
POST /api/users/login
```

**请求参数:**
```json
{
  "userCookieValue": "",
  "data": {
    "username": "用户名",
    "password": "密码",
    "remember": false
  }
}
```

**响应数据:**
```json
{
  "cookie": {
    "ID": 1,
    "Value": "32位随机字符串",
    "ExpiresAt": "过期时间"
  }
}
```

---

### 用户注册

```
POST /api/users/register
```

**请求参数:**
```json
{
  "userCookieValue": "",
  "data": {
    "username": "用户名",
    "useremail": "邮箱",
    "userpass": "密码"
  }
}
```

---

### 获取当前用户信息

```
POST /api/users/getinfo
```

**需要认证**

**响应数据:**
```json
{
  "user": {
    "ID": 1,
    "Name": "用户名",
    "Email": "邮箱"
  },
  "userInfo": {
    "UserID": 1,
    "Username": "显示名",
    "FirstName": "名",
    "Birthdate": "生日",
    "AvatarPath": "头像路径",
    "Gender": "M/F/U",
    "Region": "地区",
    "Language": "语言"
  }
}
```

---

### 获取指定用户信息

```
GET /api/users/getuserinfo/:id
```

**需要认证**

**响应数据:**
```json
{
  "userinfo": {
    "UserID": 1,
    "Username": "显示名"
  }
}
```

---

### 修改密码

```
POST /api/users/changePassword
```

**需要认证**

**请求参数:**
```json
{
  "userCookieValue": "cookie",
  "data": {
    "oldpass": "旧密码",
    "newpass": "新密码"
  }
}
```

---

### 修改邮箱

```
POST /api/users/changeEmail
```

**需要认证**

**请求参数:**
```json
{
  "userCookieValue": "cookie",
  "data": {
    "newemail": "新邮箱"
  }
}
```

---

### 更新用户信息

```
POST /api/users/updateInfo
```

**需要认证**

**请求参数:**
```json
{
  "userCookieValue": "cookie",
  "data": {
    "username": "用户名",
    "remark": "备注(作为FirstName)",
    "birthday": "2006-01-02"
  }
}
```

---

### 上传头像

```
POST /api/users/updateAvatar
```

**需要认证**

**请求格式:** `multipart/form-data`

| 字段 | 类型 | 说明 |
|------|------|------|
| cookie | string | 认证cookie |
| file | file | 图片文件(512字节-1MB) |

---

### 测试接口

```
GET /api/users/test
POST /api/users/test
```

---

## 订单模块 `/api/purchase`

### 获取订单列表

```
POST /api/purchase/getorders
```

**需要认证**

**请求参数:**
```json
{
  "userCookieValue": "cookie",
  "data": {
    "search": "搜索关键词(可选)",
    "status": "pending/ordered/arrived/received/lost/returned(可选)",
    "entries": 20,
    "page": 1
  }
}
```

**响应数据:**
```json
{
  "all_count": 100,
  "all_orders": [
    {
      "ID": 1,
      "UserID": 1,
      "Title": "订单标题",
      "Remark": "备注",
      "Link": "链接",
      "Styles": "样式",
      "OrderStatus": "pending",
      "CreatedAt": "创建时间",
      "UpdatedAt": "更新时间"
    }
  ]
}
```

---

### 获取单个订单详情

```
POST /api/purchase/getorder
```

**需要认证**

**请求参数:**
```json
{
  "userCookieValue": "cookie",
  "data": {
    "id": 1
  }
}
```

**响应数据:**
```json
{
  "order": {
    "ID": 1,
    "UserID": 1,
    "Title": "订单标题",
    "Remark": "备注",
    "Link": "链接",
    "Styles": "样式",
    "OrderStatus": "pending"
  },
  "canModify": true,
  "costs": [
    {
      "ID": 1,
      "OrderID": 1,
      "Price": 1000,
      "Quantity": 2,
      "CurrencyType": 1,
      "CostType": 1
    }
  ],
  "photos": [
    {
      "ID": 1,
      "Name": "图片名",
      "Sha256": "哈希值"
    }
  ],
  "commits": [
    {
      "ID": 1,
      "OrderID": 1,
      "UserID": 1,
      "Action": "create_status",
      "Status": "pending",
      "OldStatus": "",
      "Comment": "状态变更为: pending",
      "Photos": [],
      "CreatedAt": "时间"
    }
  ]
}
```

---

### 创建订单

```
POST /api/purchase/addorder
```

**需要认证**

**请求参数:**
```json
{
  "userCookieValue": "cookie",
  "data": {
    "title": "订单标题(必填)",
    "remark": "备注",
    "link": "链接",
    "styles": "样式",
    "costs": [
      {
        "cost": 1000,
        "costt": 2000,
        "currencytype": 1,
        "int": 2,
        "type": 1
      }
    ],
    "photos": ["图片sha256哈希"]
  }
}
```

**字段说明:**
- `cost`: 费用(分)
- `currencytype`: 货币类型 (1-CNY, 2-MOP, 3-HKD, 4-USD)
- `int`: 数量
- `type`: 费用类型 (1-单价, 2-运费)
- `photos`: 图片SHA256哈希数组

---

### 更新订单

```
POST /api/purchase/updateorder
```

**需要认证** (创建者或管理员)

**请求参数:**
```json
{
  "userCookieValue": "cookie",
  "data": {
    "id": 1,
    "title": "订单标题",
    "remark": "备注",
    "link": "链接",
    "styles": "样式",
    "costs": [
      {
        "cost": 1000,
        "currencytype": 1,
        "int": 2,
        "type": 1
      }
    ],
    "photos": ["哈希"]
  }
}
```

---

### 更新订单状态

```
POST /api/purchase/updatestatus
```

**需要认证** (创建者或管理员)

**请求参数:**
```json
{
  "userCookieValue": "cookie",
  "data": {
    "id": 1,
    "status": "ordered",
    "comment": "备注",
    "photos": ["变更图片哈希"]
  }
}
```

**状态值:**
| 值 | 说明 |
|----|------|
| pending | 待处理 |
| ordered | 已下单 |
| arrived | 已到达 |
| received | 已收件 |
| lost | 丢件 |
| returned | 退件 |

---

### 删除订单

```
POST /api/purchase/deleteorder
```

**需要认证** (创建者或管理员)

**请求参数:**
```json
{
  "userCookieValue": "cookie",
  "data": {
    "id": 1
  }
}
```

---

### 获取订单统计

```
POST /api/purchase/getordercount
```

**需要认证**

**响应数据:**
```json
{
  "pending": 10,
  "ordered": 5,
  "arrived": 3,
  "received": 20,
  "lost": 1,
  "returned": 2,
  "total": 41
}
```

---

## 日程模块 `/api/schedule`

### 获取日程列表

```
POST /api/schedule/getevents
```

**请求参数:**
```json
{
  "userCookieValue": "cookie",
  "data": {
    "start": "2026-01-01",
    "end": "2026-12-31"
  }
}
```

**响应数据:**
```json
{
  "list": [
    {
      "ID": 1,
      "UserID": 1,
      "Title": "日程标题",
      "StartDate": "2026-01-01",
      "EndDate": "2026-01-01",
      "BgColor": "#3788d9",
      "Remark": "备注",
      "edit": true
    }
  ]
}
```

---

### 创建日程

```
POST /api/schedule/addevent
```

**需要认证**

**请求参数:**
```json
{
  "userCookieValue": "cookie",
  "data": {
    "id": 0,
    "title": "日程标题",
    "start": "2026-01-01",
    "end": "2026-01-01",
    "color": "#3788d9"
  }
}
```

---

### 编辑日程

```
POST /api/schedule/editevent
```

**需要认证** (创建者或管理员)

**请求参数:**
```json
{
  "userCookieValue": "cookie",
  "data": {
    "id": 1,
    "title": "日程标题",
    "start": "2026-01-01",
    "end": "2026-01-01",
    "color": "#3788d9"
  }
}
```

---

### 删除日程

```
POST /api/schedule/deleevent
```

**需要认证** (创建者或管理员)

**请求参数:**
```json
{
  "userCookieValue": "cookie",
  "data": {
    "id": 1
  }
}
```

---

## 文件模块 `/api/files`

### 上传图片

```
POST /api/files/upload/image
```

**需要认证**

**请求格式:** `multipart/form-data`

| 字段 | 类型 | 说明 |
|------|------|------|
| cookie | string | 认证cookie |
| file | file | 图片文件 |

**响应数据:**
```json
{
  "download": "/api/files/download/文件哈希",
  "get": "/api/files/get/文件哈希",
  "hash": "sha256哈希值"
}
```

---

### 获取/下载文件

```
GET /api/files/:mode/:hash

mode: get - 获取(预览)
mode: download - 下载
```

**参数说明:**
- `get`: 返回文件预览
- `download`: 返回文件下载

---

## 静态资源 `/api/static`

### 获取头像

```
GET /api/static/avatar/:filename
```

---

## 数据模型

### TabUser_ (用户表)

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| Name | string | 用户名(唯一) |
| Email | string | 邮箱 |
| Pass | string | 密码(哈希) |
| Salt | string | 盐值 |
| Type | string | 用户类型 |
| Date | datetime | 创建时间 |

### TabUserInfo_ (用户信息表)

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| UserID | uint | 用户ID(唯一) |
| Username | string | 显示名 |
| FirstName | string | 名 |
| Birthdate | datetime | 生日 |
| Gender | char | 性别 M/F/U |
| AvatarPath | string | 头像路径 |
| Region | string | 地区 |
| Language | string | 语言 |

### TabPurchaseOrder (订单表)

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| UserID | uint | 创建者ID |
| Title | string | 标题 |
| Remark | text | 备注 |
| Link | string | 链接 |
| Styles | text | 样式 |
| OrderStatus | string | 状态 |
| CreatedAt | datetime | 创建时间 |
| UpdatedAt | datetime | 更新时间 |

### TabPurchaseCosts (订单费用表)

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| OrderID | uint | 订单ID |
| UserID | uint | 用户ID |
| Price | int | 单价(分) |
| Quantity | int | 数量 |
| CurrencyType | int | 货币类型 |
| CostType | int | 费用类型 |

### TabSchedule (日程表)

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| UserID | uint | 创建者ID |
| Title | string | 标题 |
| StartDate | string | 开始日期 |
| EndDate | string | 结束日期 |
| BgColor | string | 背景颜色 |
| Remark | text | 备注 |

### TabFileInfo_ (文件表)

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| Name | string | 文件名 |
| Size | int64 | 文件大小 |
| Path | string | 存储路径 |
| Sha256 | string | 文件哈希 |
| Mime | string | MIME类型 |
| Type | string | 文件类型 |
| UserID | uint | 上传用户ID |
| Date | datetime | 上传时间 |

---

## 权限说明

### 订单权限
- 创建者：可修改/删除自己的订单
- purchase_admin 用户组成员：可修改/删除所有订单

### 日程权限
- 创建者：可修改/删除自己的日程
- schedule_admin 用户组成员：可修改/删除所有日程
