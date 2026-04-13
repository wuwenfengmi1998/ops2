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
