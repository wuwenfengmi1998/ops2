# MEMORY.md - 长效记忆

## 项目架构

**项目名称**: OPS 运营管理系统
**技术栈**: Vue 3 + TypeScript (前端) / Go + Gin + GORM (后端)

### 目录结构
```
ops2/
├── backend/my_work/        # Go 后端（端口 8080）
├── frontend/
│   ├── ops_vue_js/         # Vue 3 Web 前端
│   └── ops2_uniapp/        # uni-app 移动端（待开发）
└── DOC/
```

## 后端架构

### 入口: `backend/my_work/main.go`
- 配置读取: `data/config.yaml`，无则复制 `defConfig/configTemp.yaml`
- 支持 SQLite/MySQL/PostgreSQL
- 按顺序初始化路由：User → Files → Schedule → Purchase → WorkOrder → Warehouse

### 核心模块 (`backend/my_work/routers/`)
| 文件 | 用途 |
|------|------|
| `apiUsers.go` | 用户认证（登录/注册/Cookie） |
| `apiFiles.go` | 文件上传/管理 |
| `apiSchedule.go` | 日程排班 |
| `apiPurchase.go` | 采购订单 |
| `apiWorkOrder.go` | 工单管理 |
| `apiWarehouse.go` | 仓库管理（容器+物品） |
| `apiStatic.go` | 静态资源 |

## 用户认证模块 (`apiUsers.go`)

### 核心函数
- `AuthenticationAuthorityFromCookie(c string)` - 验证 Cookie 并返回用户
- `AuthenticationAuthority(ctx)` - 通用认证函数，分离 Cookie 和 data
- `GetUserInfoFromUserID(userID uint)` - 通过 ID 获取用户详情

### 用户组
- 自动创建 `admins` 组和 `admin` 用户（默认密码：adminpassword）
- 各功能模块独立创建管理员组：`purchase_admin`、`work_order_admin`、`schedule_admin`、`warehouse_admin`

### API 路由 (`/api/users/*`)
| 路由 | 用途 |
|------|------|
| `POST /login` | 用户登录（返回 Cookie） |
| `POST /register` | 用户注册 |
| `POST /getinfo` | 获取当前用户信息 |
| `POST /changePassword` | 修改密码 |
| `POST /changeEmail` | 修改邮箱 |
| `POST /updateAvatar` | 更新头像（FormData 上传） |
| `POST /updateInfo` | 更新用户详情 |
| `GET /getuserinfo/:id` | 获取指定用户信息 |
| `GET/POST /test` | 测试接口 |

### 密码机制
- 密码加盐哈希（Salt + Hash）
- 支持 `text` / `md5` / `md5salt` 三种哈希类型（配置指定）

## 文件管理模块 (`apiFiles.go`)

### 数据表
- `TabFileInfo_` - 文件元数据（SHA256 哈希为唯一标识）

| 字段 | 说明 |
|------|------|
| `Sha256` | 文件哈希（主键/索引） |
| `Name` | 原始文件名 |
| `Path` | 存储路径 |
| `Mime` | MIME 类型 |
| `Type` | 文件类型（image/video/pdf 等） |
| `Const` | 引用计数（同文件多次上传只存一份） |

### API 路由 (`/api/files/*`)
| 路由 | 用途 |
|------|------|
| `POST /upload/image` | 上传图片（FormData，含 SHA256 去重） |
| `GET /:mode/:hash` | 获取文件（mode=get 下载，mode=download 预览） |

### 存储结构
```
data/
├── static/avatar/           # 用户头像
└── upload/
    ├── image/              # 图片（以 SHA256 命名）
    ├── video/
    ├── music/
    └── pdf/
```

## 日程排班模块 (`apiSchedule.go`)

### 数据表
| 表 | 用途 |
|---|---|
| `TabSchedule` | 日程（软删除） |
| `TabScheduleLog` | 操作日志 |

### 日程结构
| 字段 | 说明 |
|------|------|
| `Title` | 日程标题 |
| `StartDate` | 开始日期（YYYY-MM-DD） |
| `EndDate` | 结束日期（YYYY-MM-DD） |
| `BgColor` | 背景颜色（默认 #3788d9） |
| `Remark` | 备注 |

### API 路由 (`/api/schedule/*`)
| 路由 | 用途 |
|------|------|
| `POST /getevents` | 获取日程列表（按日期范围） |
| `POST /addevent` | 新增日程 |
| `POST `/editevent`` | 编辑日程 |
| `POST /deleevent` | 删除日程（软删除） |

### 查询逻辑
```sql
WHERE start_date <= :end AND end_date >= :start
```

## 静态资源模块 (`apiStatic.go`)

| 路由 | 用途 |
|------|------|
| `GET /static/avatar/:filename` | 获取用户头像 |

### 数据库模型 (`backend/my_work/models/sql.go`)
- `TabUser_` - 用户
- `TabUserGroups_` - 用户组
- `TabUserInfo_` - 用户详情
- `TabCookie_` - 登录 Cookie（有效期 604800 秒）
- `APIRequestLog_` - API 日志

### 仓库模块核心表 (`apiWarehouse.go`)
- `TabWarehouseContainer` - 容器（树形，最多5层嵌套）
- `TabWarehouseItem` - 物品
- `TabWarehouseItemCommit` - 物品移动记录
- `TabWarehouseLog` - 操作日志
- `TabWarehouseItemWorkOrderBind` - 物品-工单关联

## 采购模块 (`backend/my_work/routers/apiPurchase.go`)

### 数据表
| 表 | 用途 |
|---|---|
| `TabPurchaseOrder` | 采购订单（软删除） |
| `TabPurchaseCosts` | 费用明细（单价/运费，支持多币种） |
| `TabPurchaseFileBind` | 图片关联 |
| `TabPurchaseCommit` | 状态变更记录 |
| `TabPurchaseLog` | 操作日志 |

### 订单状态流程
```
pending(待处理) → ordered(已下单) → arrived(已到达) → received(已收件)
                                    ↓
                              lost(丢件) / returned(退件)
```

### 货币类型
`1-CNY` / `2-MOP` / `3-HKD` / `4-USD`

### API 路由 (`/api/purchase/*`)
| 路由 | 用途 |
|------|------|
| `POST /getorder` | 获取订单详情（含费用、图片、状态记录、关联工单） |
| `POST /getorders` | 获取订单列表（支持搜索、分页、状态筛选） |
| `POST /addorder` | 新增订单 |
| `POST /updateorder` | 编辑订单（含费用、图片重建） |
| `POST /deleteorder` | 删除订单 |
| `POST /updatestatus` | 更新订单状态（可附评论/图片） |
| `POST /delete_commit` | 删除状态记录 |
| `POST /getordercount` | 统计各状态数量 |
| `POST /search_work_orders` | 搜索工单（用于关联） |

## 工单模块 (`backend/my_work/routers/apiWorkOrder.go`)

### 数据表
| 表 | 用途 |
|---|---|
| `TabWorkOrder` | 工单（软删除） |
| `TabWorkOrderFileBind` | 工单图片关联 |
| `TabWorkOrderCommit` | 进度记录 |
| `TabWorkOrderLog` | 操作日志 |
| `TabWorkOrderCommitFileBind` | 进度关联图片 |
| `TabWorkOrderPurchaseOrderBind` | 工单-采购订单关联 |

### 工单状态流程
```
pending(待处理) → checked(已检查) → parts_ordered(已下单零件) → repaired(已维修) → returned(已送还)
                                       ↓
                                  unrepairable(无法维修)
```

### 关联关系
- 工单 ↔ 仓库物品 (`TabWarehouseItemWorkOrderBind`)
- 工单 ↔ 采购订单 (`TabWorkOrderPurchaseOrderBind`)
- 特殊逻辑：状态变更为 `returned` 时，自动移除物品的容器绑定

### API 路由 (`/api/work_order/*`)
| 路由 | 用途 |
|------|------|
| `POST /add` | 新增工单（可关联物品） |
| `POST /update` | 编辑工单 |
| `POST /list` | 获取工单列表 |
| `POST /get` | 获取工单详情（含图片、进度、关联物品/采购订单） |
| `POST /commit` | 提交进度（更新状态，可关联采购订单） |
| `POST /delete` | 删除工单 |
| `POST /delete_commit` | 删除进度 |
| `POST /count` | 统计各状态数量 |
| `POST /search_purchase_orders` | 搜索采购订单（用于关联） |

## 前端架构

### Web 前端 (`frontend/ops_vue_js/`)

**API 封装** (`src/api/index.js`):
- 基础 URL: `/api`
- 请求自动注入 `userCookieValue`
- 响应统一处理：err_code=-44 表示 Cookie 过期，自动登出
- 返回格式: `{ errCode, data }`

**路由** (`src/router/index.js`):
- 使用 `createWebHashHistory`（hash 模式）
- 认证页面: `/login`, `/register`, `/forgot_password`
- 需要登录的页面在白名单外

**国际化**: `src/i18n/en.json`, `zh-CN.json`

### 移动端 (`frontend/ops2_uniapp/`)

**技术栈**: uni-app + Vue 3 + Pinia + HBuilderX

**项目结构**:
```
ops2_uniapp/
├── api/                    # API 接口封装
│   ├── index.js            # 基础请求工具（框架已有，方法待实现）
│   ├── request.js          # 请求配置（待完善）
│   └── user.js             # 用户接口
├── components/
│   └── my-toast/           # 自定义 Toast 组件
├── pages/                  # 页面
│   ├── index/index.vue     # 主页 TabBar（占位）
│   ├── order/order.vue     # 订单 TabBar（占位）
│   ├── message/message.vue # 消息 TabBar（占位）
│   ├── user/user.vue       # 用户 TabBar（基础框架）
│   ├── login/login.vue     # 登录页（已完成 85%）
│   └── settings/settings.vue # 设置页（已完成 90%）
├── stores/
│   ├── config.js           # 配置 Store（完整）
│   └── user.js            # 用户 Store（基础）
├── utils/
│   └── index.js            # 工具函数（isUrl）
├── pages.json              # 路由配置
├── manifest.json           # 应用配置
└── package.json           # 依赖（pinia）
```

**Stores 状态管理**:
- `useConfigStore`: apiBaseUrl / appName / version / theme
  - `setApiBaseUrl()` / `getApiBaseUrl()` - API 地址持久化
- `useUserStore`: username / token
  - `setUser()` / `logout()`

**API 封装状态**:
- `api/index.js`: 框架有，get/post/upload 方法空实现
- 登录页直接使用 `uni.request()` 调用，绕过了封装层
- 需要完善 Cookie 认证机制（参照 Web 前端）

**页面完成度**:
| 页面 | 完成度 | 说明 |
|------|--------|------|
| index | 5% | 占位文本 |
| order | 5% | 占位文本 |
| message | 5% | 占位文本 |
| user | 20% | 登录按钮 + 设置入口 |
| login | **85%** | 表单 + 验证 + 请求 + Toast |
| settings | **90%** | API 地址编辑 + 连接测试 |

**当前总完成度**: ~35-40%

**移动端待开发**:
1. 完善 API 封装层（Cookie 认证）
2. 实现各功能页面（主页仪表盘、订单列表、消息列表、用户中心）
3. 添加更多组件（Loading、确认对话框、空状态）
4. 对接后端各模块（仓库、工单、采购等）

## 前后端交互协议

### 请求格式 (POST JSON)
```json
{
  "data": {
    "userCookieValue": "xxx",
    ...业务参数
  }
}
```

### 响应格式
```json
{
  "err_code": 0,
  "return": { ... }
}
```

### 认证机制
- 登录成功后服务端返回 Cookie（存储在 `TabCookie_` 表）
- 后续请求通过 `userCookieValue` 字段传递
- Cookie 过期码: -44

## 开发注意事项

1. **移动端开发时**: 需要完善 `api/index.js` 的请求封装，参照 Web 前端实现 Cookie 认证
2. **仓库模块**: 是当前开发重点，支持树形容器、物品管理、工单关联
3. **同源部署**: 后端直接 serve `./dist` 静态文件，简化部署

## 更新记录
- 2026-04-24: 首次梳理项目运行逻辑，保存长效记忆
