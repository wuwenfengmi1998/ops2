# MEMORY.md - 长效记忆

> 最后更新: 2026-04-29

---

## 项目概览

**项目名称**: OPS 运营管理系统
**技术栈**: Vue 3 + Vite (前端 Web) / uni-app (移动端) / Go + Gin + GORM (后端)

```
ops2/
├── backend/my_work/          # Go 后端（端口由 config.yaml 决定，默认 8080）
├── frontend/
│   ├── ops_vue_js/           # Vue 3 Web 前端
│   └── ops2_uniapp/          # uni-app 移动端
└── DOC/
```

---

## 后端架构

### 启动流程 (`main.go`)
1. 检查 `./data/config.yaml`，不存在则复制 `./defConfig/configTemp.yaml`
2. `configed` 必须为 `true` 才能启动
3. 初始化顺序：`ReturnInit` → `ApiUserInit` → `ApiFilesInit` → `ApiScheduleInit` → `ApiPurchaseInit` → `ApiWorkOrderInit` → `ApiWarehouseInit` → `BindsInit`
4. 静态文件服务 `./dist`，所有非 `/api` 请求转发给前端 HTML
5. 支持 TLS（证书路径在 config 中配置）
6. 版本信息通过 `-ldflags -X` 编译注入（`GitVersion / GitCommit / BuildTime`）

### 路由根 (`api.go`)
```
/api
├── /static      → ApiStatic
├── /users       → ApiUser
├── /files       → ApiFiles
├── /purchase    → ApiPurchase
├── /schedule    → ApiSchedule
├── /work_order  → ApiWorkOrder
├── /warehouse   → ApiWarehouse
└── /admin       → ApiSysAdmin
```

### 请求/响应格式
```json
// 请求体
{ "userCookieValue": "xxx", "data": { ...业务参数 } }

// 响应体
{ "err_code": 0, "err_msg": "apiOK", "return": { ... } }
```
错误码从 `./defConfig/errorCodes.json` 读取，常用：`apiOK`=0，`userNoLogin`=-44，`permission_denied`，`parameErr`，`dbErr`

### 认证机制
- 登录成功返回 Cookie 对象（存 `TabUserCookie` 表）
- 后续请求通过 JSON body 中 `userCookieValue` 传递
- Cookie 过期后端返回 err_code=-44，前端拦截器自动处理
- `AuthenticationAuthority(ctx)` → 分离 cookie 和 data，验证返回 `(isAuth bool, user TabUser, data map)`

---

## 用户认证模块 (`apiUsers.go`)

### 数据表
| 表 | 说明 |
|---|---|
| `TabUser` | 用户（Name 唯一索引） |
| `TabUserGroups` | 用户组 |
| `TabUserGroupBinds` | 用户-组绑定 |
| `TabUserInfo` | 用户详情（头像/昵称/性别/生日/语言等） |
| `TabUserCookie` | 登录 Cookie（ExpiresAt，Remember 字段） |
| `TabUserLoginFailLog` | 登录失败日志（24小时内聚合，累计 Count） |

### 初始化
- 自动创建 `admins` 组 + `admin` 用户（默认密码 `adminpassword`）
- 密码：加盐哈希，支持 `text` / `md5` / `md5salt`（config 指定）

### 权限缓存（内存）
各模块维护各自的管理员 ID 列表：
| 变量 | 所属模块 | 刷新函数 |
|---|---|---|
| `sysAdmins []uint` | apiUsers.go | `updateSysAdminsCash()` |
| `scheduleAdmins []uint` | apiSchedule.go | `ScheduleUpdateAdminsCash()` |
| `purchaseAdmins []uint` | apiPurchase.go | `PurchaseUpdateAdminsCash()` |
| `workOrderAdmins []uint` | apiWorkOrder.go | `WorkOrderUpdateAdminsCash()` |
| `warehouseAdmins []uint` | apiWarehouse.go | `WarehouseUpdateAdminsCash()` |

**修改用户组时自动刷新缓存**（`apiSysAdmin.go` 的 `add_group_member` 和 `remove_group_member`）：
```go
switch group.Name {
case "admins":         updateSysAdminsCash()
case "schedule_admin": ScheduleUpdateAdminsCash()
case "purchase_admin": PurchaseUpdateAdminsCash()
case "work_order_admin": WorkOrderUpdateAdminsCash()
case "warehouse_admin":  WarehouseUpdateAdminsCash()
}
```

### API 路由 (`/api/users/*`)
| 路由 | 用途 |
|------|------|
| `POST /login` | 登录（返回 Cookie 对象，含 Remember） |
| `POST /register` | 注册 |
| `POST /getinfo` | 获取当前用户信息（含 isSysAdmin 字段） |
| `POST /changePassword` | 修改密码（oldpass/newpass） |
| `POST /changeEmail` | 修改邮箱 |
| `POST /updateAvatar` | 更新头像（FormData） |
| `POST /updateInfo` | 更新详情（firstName/username/birthdate/gender/region/language） |
| `GET /getuserinfo/:id` | 获取指定用户信息 |

---

## 系统管理模块 (`apiSysAdmin.go`)

**路由前缀**: `/api/admin/*`，全部要求 `SysAdminCheck`

| 路由 | 用途 |
|------|------|
| `POST /sysadmins` | 获取系统管理员 ID 列表 |
| `POST /users` | 用户列表（分页+搜索，返回含头像路径） |
| `POST /groups` | 用户组列表（含成员数 + 前5个成员ID） |
| `POST /group_members` | 指定组的成员列表（分页） |
| `POST /user_detail` | 用户详情（基本信息 + userinfo） |
| `POST /reset_user_password` | 重置密码（同时注销该用户所有 cookie） |
| `POST /add_group_member` | 添加用户到组（含缓存刷新） |
| `POST /remove_group_member` | 从组移除用户（含缓存刷新） |
| `POST /login_fail_logs` | 登录失败日志（分页+搜索 username/IP） |

`SysAdminCheck(userID)` 直接查内存 `sysAdmins` 列表。

---

## 文件管理模块 (`api_Files.go`)

### 数据表 `TabFileInfo`
| 字段 | 说明 |
|---|---|
| `Sha256` | 唯一标识，去重键 |
| `Name` | 原始文件名 |
| `Path` | 存储路径 |
| `Mime` | MIME 类型 |
| `Type` | image/video/pdf 等 |
| `Const` | 引用计数 |

### 存储结构
```
data/
├── static/avatar/        # 用户头像（/api/static/avatar/:filename）
└── upload/
    ├── image/            # 以 SHA256 命名
    ├── video/
    ├── music/
    └── pdf/
```

### API (`/api/files/*`)
| 路由 | 用途 |
|---|---|
| `POST /upload/image` | 上传图片（FormData + SHA256 去重） |
| `GET /:mode/:hash` | `get`=下载（带文件名），`download`=预览 |

---

## 日程模块 (`apiSchedule.go`)

### 数据表
- `TabSchedule`（软删除）：Title / StartDate / EndDate / BgColor(默认#3788d9) / Remark
- `TabScheduleLog`：操作日志

### API (`/api/schedule/*`)
| 路由 | 用途 |
|---|---|
| `POST /getevents` | 按日期范围查询（`start_date <= end AND end_date >= start`） |
| `POST /addevent` | 新增 |
| `POST /editevent` | 编辑 |
| `POST /deleevent` | 软删除 |

---

## 采购模块 (`apiPurchase.go`)

### 数据表
| 表 | 用途 |
|---|---|
| `TabPurchaseOrder` | 采购订单（软删除） |
| `TabPurchaseCosts` | 费用明细（单价/运费，多币种） |
| `TabPurchaseFileBind` | 图片关联 |
| `TabPurchaseCommit` | 状态变更记录 |
| `TabPurchaseLog` | 操作日志 |

### 状态流
```
pending → ordered → arrived → received
                    ↓
               lost / returned
```

### 货币：`1-CNY / 2-MOP / 3-HKD / 4-USD`

### API (`/api/purchase/*`)
| 路由 | 用途 |
|---|---|
| `POST /getorder` | 订单详情（含费用/图片/状态记录/关联工单） |
| `POST /getorders` | 列表（搜索/分页/状态筛选） |
| `POST /addorder` | 新增 |
| `POST /updateorder` | 编辑（费用/图片重建） |
| `POST /deleteorder` | 删除 |
| `POST /updatestatus` | 更新状态（可附评论/图片） |
| `POST /delete_commit` | 删除状态记录 |
| `POST /getordercount` | 各状态数量统计 |
| `POST /search_work_orders` | 搜索工单（用于关联） |

---

## 工单模块 (`apiWorkOrder.go`)

### 数据表
| 表 | 用途 |
|---|---|
| `TabWorkOrder` | 工单（软删除） |
| `TabWorkOrderFileBind` | 工单图片关联 |
| `TabWorkOrderCommit` | 进度记录 |
| `TabWorkOrderCommitFileBind` | 进度关联图片 |
| `TabWorkOrderPurchaseOrderBind` | 工单-采购订单关联（含 CommitID） |
| `TabWorkOrderLog` | 操作日志 |

### 状态流
```
pending → checked → parts_ordered → repaired → returned
                          ↓
                    unrepairable
```

### 特殊逻辑
- 状态变为 `returned` 时，自动移除物品的容器绑定（`ContainerID = nil`）
- 工单 ↔ 物品：`TabWarehouseItemWorkOrderBind`
- 工单 ↔ 采购：`TabWorkOrderPurchaseOrderBind`

### API (`/api/work_order/*`)
| 路由 | 用途 |
|---|---|
| `POST /add` | 新增（可关联物品） |
| `POST /update` | 编辑 |
| `POST /list` | 列表 |
| `POST /get` | 详情（含图片/进度/关联物品/采购） |
| `POST /commit` | 提交进度（更新状态，可关联采购） |
| `POST /delete` | 删除 |
| `POST /delete_commit` | 删除进度 |
| `POST /count` | 各状态统计 |
| `POST /search_purchase_orders` | 搜索采购订单（用于关联） |

---

## 仓库模块 (`apiWarehouse.go`)

### 数据表
| 表 | 用途 |
|---|---|
| `TabWarehouseContainer` | 容器（树形，最多5层，ParentID=nil为顶级） |
| `TabWarehouseItem` | 物品（ContainerID=nil表示未入库） |
| `TabWarehouseItemCommit` | 物品移动记录 |
| `TabWarehouseLog` | 操作日志 |

### 跨模块绑定表 (`binds.go`)
| 表 | 用途 |
|---|---|
| `TabWarehouseItemWorkOrderBind` | 物品-工单 |
| `TabWarehouseItemFileBind` | 物品-图片 |
| `TabWarehouseContainerFileBind` | 容器-图片 |
| `TabPurchaseFileBind` | 采购-图片 |
| `TabWorkOrderFileBind` | 工单-图片 |
| `TabWorkOrderCommitFileBind` | 工单进度-图片 |
| `TabWorkOrderPurchaseOrderBind` | 工单-采购 |

---

## Web 前端架构 (`frontend/ops_vue_js/`)

**技术栈**: Vue 3 + Vite 7 + Pinia + Vue Router (hash 模式) + Vue I18n + Tailwind CSS v4 + Tabler Icons

### 目录结构
```
src/
├── api/
│   ├── index.js        # Axios 实例，基础 URL /api，请求拦截注入 cookie，响应拦截处理 -44
│   ├── auth.js         # 认证 + sysadmin 管理 API
│   ├── purchase.js     # 采购 API
│   ├── warehouse.js    # 仓库 API
│   ├── work_order.js   # 工单 API
│   ├── schedule.js     # 日程 API
│   └── users.js        # 其他用户信息 API（按需加载头像/用户名）
├── components/
│   ├── AppHeader.vue          # 导航栏（含系统管理入口，权限判断 isSysAdmin）
│   ├── AppFooter.vue
│   ├── AppToast.vue
│   ├── ConfirmDialog.vue
│   ├── PurchaseOrderForm.vue
│   ├── SettingNav.vue
│   ├── tagadder.vue
│   ├── useDropzone.vue        # 文件拖拽上传
│   ├── imageCropper.vue       # 图片裁剪
│   ├── datePicker.vue
│   ├── dateTimePicker.vue
│   └── datatimePickerForFullCalendar.vue
├── composables/
├── i18n/
│   ├── en.json                # 英文翻译
│   └── zh-CN.json             # 中文翻译
├── layouts/
│   ├── DefaultLayout.vue      # 需要登录的页面布局
│   └── AuthLayout.vue         # 认证页面全屏布局
├── router/index.js
├── stores/
│   ├── user.js        # 当前用户（isLoggedIn / isSysAdmin / cookie / avatarUrl）
│   ├── users.js       # 其他用户信息缓存（按需拉取，防重复请求）
│   └── toast.js       # 全局 Toast
└── views/
    ├── HomeView.vue
    ├── ScheduleView.vue       # FullCalendar
    ├── SysAdminView.vue       # 系统管理（仅 sysAdmin 可访问，meta.requireSysAdmin）
    ├── AdminView.vue
    ├── purchase/              # PurchaseList / addorder / ShowOrder / editorder
    ├── work_order/            # WorkOrderList / AddEditWorkOrder / ShowWorkOrder
    ├── warehouse/             # WarehouseOverview / ContainerList / ContainerDetail / ItemList / ItemDetail / AddItem / ItemEdit
    └── settings/              # AccountView / ContactView / SecurityView
```

### 路由说明
- 公开页：`/` `/login` `/register` `/forgot_password` `/schedule` `/404`
- `/sysadmin` 需要 `meta.requireSysAdmin`，不满足跳回首页
- 未登录跳转 `/login?redirect=原路径`

### 状态管理 (`stores/user.js`)
- Cookie 持久化：`Remember=true` 存 localStorage，否则只存 sessionStorage
- `isSysAdmin` 由 `/users/getinfo` 返回的 `isSysAdmin` 字段驱动
- `fetchUserInfo()` 在 `login()` 后自动调用

### i18n 翻译节点
`week / errorpage / appname / tagadder / dropzone / cropper / purchase / work_order / warehouse / purchase_addorder / schedule / home / message / settings / button / footer / cost_type / order_status / sysadmin`

### 构建配置
- 输出目录: `../../backend/my_work/dist`（后端直接 serve）
- 开发代理: `/api` → `http://127.0.0.1:8080`
- 路径别名: `@` → `./src`

---

## 移动端 (`frontend/ops2_uniapp/`)

**技术栈**: uni-app + Vue 3 + Pinia

**当前完成度**: ~40%

| 页面 | 完成度 |
|---|---|
| login | 85%（表单/验证/请求/Toast） |
| settings | 90%（API 地址配置/连接测试） |
| index/order/message | 5%（占位） |
| user | 20%（登录入口） |

**待完成**:
- 完善 `api/index.js`（Cookie 认证，参照 Web 前端）
- 实现各功能页面（仓库/工单/采购等）

---

## 开发注意事项

1. **后端 JSON 字段命名**: 结构体字段为 PascalCase（如 `UserID / AvatarPath`），前端需对应
2. **头像**: `usersStore.getAvatarUrlFromUserID(id)` 返回 `/api/static/avatar/{path}` 或默认 `/ava.svg`
3. **i18n 修改**: 同时修改 `en.json` 和 `zh-CN.json` 两个文件
4. **同源部署**: 后端 serve `./dist` 静态资源，前端 build 直接输出到后端目录
5. **API 封装**: `api/index.js` 统一返回 `{ errCode, data, raw }`

---

## 更新记录
- 2026-04-24: 首次梳理
- 2026-04-28: 修复 SysAdminView 中 `<img .../>` 误删 bug；添加 message.sysadmin 英文翻译
- 2026-04-29: 全面重新分析代码，更新并精简记忆文档
