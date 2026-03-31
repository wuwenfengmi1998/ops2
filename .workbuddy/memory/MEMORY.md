# 项目长期记忆 - OPS2

## 项目概况
- **项目名称**：Ops（Operations 运营管理系统）
- **类型**：前后端分离的工作流/运营管理系统
- **工作区路径**：`c:\Users\wuwen\Documents\prj\ops2`

## 技术栈

### 后端（backend/）
- **语言**：Go
- **框架**：Gin（HTTP 框架）+ GORM（ORM）
- **数据库**：支持 SQLite / MySQL / PostgreSQL（通过配置切换）
- **配置**：YAML 格式，路径 `./data/config.yaml`，模板在 `./defConfig/configTemp.yaml`
- **静态文件**：后端直接 serve `./dist` 目录下的前端构建产物
- **TLS**：支持 HTTPS（可配置）

### 前端（frontend/ops_vue_js/）⬅️ 主力开发目录
- **框架**：Vue 3 + JavaScript（非 TypeScript）
- **路由**：Vue Router（Hash 模式）
- **构建工具**：Vite 7
- **UI 框架**：Tabler（Bootstrap 5 + @tabler/core）
- **图标**：@tabler/icons-vue
- **状态管理**：Pinia
- **国际化**：vue-i18n
- **日期选择**：flatpickr / litepicker
- **文件上传**：FilePond
- **图片裁剪**：CropperJS / @cropper/*
- **日历**：FullCalendar（含 daygrid/timegrid/list/interaction）
- **其他组件**：MyOffcanvas、imageCropper、datePicker、dateTimePicker、tagadder、useDropzone 等

> 注意：`frontend/ops_vue/`（TypeScript 版）是旧目录，已弃用

### 前端页面路由（ops_vue_js）
- `/` — 首页（HomeView）
- `/login` — 登录
- `/register` — 注册
- `/forgot_password` — 找回密码
- `/admin` — 管理后台
- `/schedule` — 日程/排班（FullCalendar）
- `/purchase` — 采购订单列表
- `/purchase/addorder` — 新增采购订单
- `/purchase/showorder/:id` — 查看采购订单详情
- `/warehouse` — 仓库管理
- `/settings/account` — 账户设置
- `/settings/contact` — 联系信息设置
- `/settings/security` — 安全设置
- `/404` — 404 页面

## 数据模型（GORM 表结构）
- `TabUser_` - 用户表（name 唯一索引，支持 md5/md5salt 密码哈希）
- `TabUserGroups_` - 用户组表
- `TabUserGroupBinds_` - 用户-组绑定关系表
- `TabUserInfo_` - 用户详情表（头像、性别、语言等）
- `TabCookie_` - Session Cookie 表（含过期时间、记住我功能）
- `TabFileInfo_` - 文件信息表（支持图片/视频/音乐/PDF）
- `APIRequestLog_` - API 请求日志表
- `TabPurchaseOrder` - 采购订单表（含照片JSON、快递单号、订单状态）
- `TabPurchaseCosts` - 采购费用明细表

## API 路由结构
- `POST /api/users/...` - 用户相关（登录、注册、鉴权）
- `POST /api/files/...` - 文件上传管理
- `POST /api/purchase/getorders` - 获取采购订单列表（分页）
- `POST /api/purchase/addorder` - 新增采购订单
- `GET /api/static/...` - 静态资源访问
- 认证方式：请求体中携带 `userCookieValue` 字段

## 前端页面
- 见上方"前端页面路由"章节

## 项目现状（2026-03-31）
- 后端基础架构完整，采购模块已有基础实现
- 前端 `ops_vue_js` 目录是主力开发目录（JS 版 Vue3，Tabler UI）
- 已有页面：登录/注册/采购/日程/仓库/设置等
- 前端构建产物放在 `backend/dist/` 供后端 serve
- `frontend/ops_vue/`（TypeScript 版）是旧目录，已弃用
- 目录名原为 `frontent`（拼写错误），实际文件系统路径为 `frontend`

## 开发规范
- API 请求统一携带 `userCookieValue` 做身份验证
- 响应统一用 `ReturnJson(ctx, errorCode, data)` 格式
- 错误码定义在 `./defConfig/errorCodes.json`
