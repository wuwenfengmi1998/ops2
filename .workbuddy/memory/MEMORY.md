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
- **CSS 框架**：Tailwind CSS v4（@tailwindcss/vite 插件）
- **图标**：@tabler/icons-vue
- **状态管理**：Pinia
- **国际化**：vue-i18n
- **日期选择**：flatpickr / litepicker
- **文件上传**：FilePond
- **图片裁剪**：CropperJS（@cropper/elements）
- **日历**：FullCalendar（含 daygrid/timegrid/list/interaction）
- **其他组件**：imageCropper、tagadder、dateTimePicker、useDropzone 等
- ~~**UI 框架**：Tabler（Bootstrap 5 + @tabler/core）~~ 已弃用（2026-03-31 迁移至 Tailwind）

> 注意：`frontend/ops_vue/`（TypeScript 版）是旧目录，已弃用

### 前端页面路由（ops_vue_js）
- `/` — 首页（HomeView）
- `/login` — 登录（AuthLayout）
- `/register` — 注册（AuthLayout）
- `/forgot_password` — 找回密码（AuthLayout）
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

## 项目现状（2026-03-31 更新）
### 前端
- 前端 `ops_vue_js` 目录是主力开发目录（Vue 3 + Tailwind CSS v4）
- **已完成前端整体重构**：API 层 async/await、Router 导航守卫、composables、布局分离
- **已完成 Tabler → Tailwind CSS v4 迁移**
- **已修复所有字符损坏文件**（20 个 Vue 文件，因批量脚本偏移错误）
- **已修复国际化翻译缺失问题**：补充 `account_information` 翻译键
- **已修复头像裁剪功能**：
  - 修复事件通信问题（`crop_to_canvas` → `crop-data-url`）
  - 修正坐标计算逻辑，解决预览不正确问题
  - 添加多层容错机制和详细调试信息
- 所有页面构建通过，6170+ modules, 0 errors
- 前端构建产物放在 `backend/dist/` 供后端 serve
- `frontend/ops_vue/`（TypeScript 版）是旧目录，已弃用

### 后端（重构完成 ✅）
- **已完成基础架构重构**：cmd/internal/pkg 三层架构
- **用户认证模块重构完成**：Handler → Service → Repository 分层
- **采购订单模块重构完成**：新增分层架构，兼容现有前端API
- **新增中间件系统**：认证、日志、CORS、恢复中间件
- **统一API响应**：标准错误码映射和响应格式
- **模块化路由系统**：API v1 版本路由定义清晰分离
- **新目录结构**：
  - `main.go` - 应用入口（已从 cmd/ops-server/main.go 合并至根目录，2026-04-01）
  - `internal/config/` - 配置管理
  - `internal/database/` - 数据库连接和迁移
  - `internal/handler/` - HTTP处理器（auth_handler.go, purchase_handler.go）
  - `internal/service/` - 业务逻辑层（auth_service.go, purchase_service.go）
  - `internal/repository/` - 数据访问层（user_repository.go, purchase_repository.go）
  - `internal/middleware/` - 中间件系统（auth.go, logging.go, cors.go）
  - `api/v1/` - API定义（routes.go）
  - `pkg/response/` - 统一响应处理

### 重构进展总结
- ✅ **用户认证模块**：完整迁移到分层架构
- ✅ **采购订单模块**：完整迁移，同时支持原始POST路由和RESTful API
- ✅ **文件管理模块**：完整迁移，支持分层架构
- ✅ **基础架构**：所有中间件、配置、数据库连接已完成
- ✅ **路由和中间件系统**：已完成统一管理和配置（2026-03-31）
- ✅ **编译状态**：项目编译成功（需要CGO_ENABLED=1以支持SQLite）

### 新路由架构（2026-03-31）
- **主入口**：`main.go`（根目录）- 现代化主入口，支持优雅关机
- **路由配置**：`api/`包统一管理所有路由
- **兼容性**：完全兼容现有前端API `/api/*`
- **新增API**：RESTful API v1 `/api/v1/*`
- **中间件系统**：环境感知的日志、CORS、认证、恢复中间件
- **静态文件**：智能SPA支持，支持Vue Router history模式

### 中间件系统
- **CORS中间件**：完整跨域支持
- **日志中间件**：开发环境用简易日志，生产环境用详细日志
- **认证中间件**：支持多种认证方式（Bearer令牌、userCookieValue）
- **恢复中间件**：Panic恢复和错误处理

### 已完成模块
1. 文件管理模块的分层重构 ✅
2. 静态文件服务整合 ✅
3. API请求日志模块 ✅
4. 管理员权限控制 ⏳
5. 系统配置管理 ✅

### 技术架构升级
1. **分层架构完成**：Handler → Service → Repository
2. **统一错误处理**：标准错误码和响应格式
3. **路由系统整合**：兼容性路由 + RESTful API v1
4. **中间件规范化**：统一的中间件加载和配置
5. **开发工具完善**：run-dev.bat启动脚本，配置文档

### 技术规范
- **认证方式**：兼容前端 `userCookieValue` POST字段、Authorization头、Cookie头
- **响应格式**：统一使用 `pkg/response` 包的标准响应
- **错误码**："0"成功、"-1"内部错误、"-2"参数错误、"-3"未登录、"-4"用户存在、"-5"用户不存在、"-42"凭证错误
- **数据库**：支持SQLite/MySQL/PostgreSQL切换
- **API版本**：v1 API统一在 `/api/v1/` 路径下

## 经验教训
- **批量字符替换脚本危险**：需在源码上使用前先备份，并限定替换范围
- **`@tabler/icons-vue` 不包含所有图标**：如 `IconFileTypeText` 不存在，使用前需确认

## 前端重构后架构（2026-03-31）
- **API 层**：`src/api/` — axios 实例 + 拦截器，async/await 封装
- **Composables**：`src/composables/` — usePageTitle、useValidation、isValidEmail
- **Stores**：`src/stores/user.js`（精简）、`src/stores/toast.js`（全局通知）
- **布局**：`src/layouts/DefaultLayout.vue`（主站）、`AuthLayout.vue`（认证页）
- **公共组件**：AppHeader、AppFooter、AppToast、SettingNav
- **命名规范**：PascalCase 文件名，camelCase 函数名

## 开发规范

### 后端架构规范
- **分层架构**：Handler → Service → Repository → Database
- **认证方式**：
  - 兼容现有前端：POST JSON中的 `userCookieValue` 字段
  - 标准方式：Authorization: Bearer token 或 Cookie header
- **响应格式**：
  ```json
  {
    "code": "0",          // 错误码，0表示成功
    "message": "Success", // 人类可读的消息
    "data": {}            // 实际数据
  }
  ```
- **错误码系统**：
  - "0": 成功
  - "-1": 内部错误
  - "-2": 参数错误
  - "-3": 用户未登录
  - "-4": 用户已存在
  - "-5": 用户不存在
  - "-42": 用户名或密码错误
- **依赖注入**：Handler通过Service，Service通过Repository访问数据库

### 前端规范（保持不变）
- API 请求统一携带 `userCookieValue` 做身份验证
- 错误码定义在 `./defConfig/errorCodes.json`
