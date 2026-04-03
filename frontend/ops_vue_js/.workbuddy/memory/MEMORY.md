# 项目长期记忆

## 项目概况
- **项目路径**：`c:/Users/wuwen/Documents/prj/ops2/frontend/ops_vue_js`
- **技术栈**：Vue 3 + Vite + Tailwind CSS + FullCalendar + Vue I18n + Pinia
- **主要页面**：scheduleView.vue（日程日历视图）

## 关键实现决策

### FullCalendar 日程 hover 展示完整标题（2026-04-04，已验证可用方案）
- **问题**：月视图中日程标题过长被 ellipsis 截断，用户无法看到完整内容
- **方案**：
  1. `eventDidMount` 回调用 `setTimeout(0)` 等布局完成后，检测 `titleEl.scrollWidth > titleEl.clientWidth`，**只有实际被截断时**才把完整标题写到事件根元素 `info.el` 的 `data-full-title` 属性上
  2. CSS 用非 scoped 的 `<style>` 块（避免 Vue scoped + ::after 兼容问题），在 `.fc-daygrid-event[data-full-title]:hover::after` 上用 `content: attr(data-full-title)` 展示 tooltip
  3. 只把 `.fc-daygrid-event-harness`、`.fc-daygrid-event` 设为 `overflow: visible`，**不能**动 `.fc-daygrid-day`、`.fc-daygrid-day-frame`、`.fc-daygrid-day-events`，否则会破坏拖拽时的层叠上下文，导致拖动中日程消失
- **坑**：Vue `scoped` CSS 中 `:deep(selector):hover::after` 写法编译后伪元素不会生效，必须用独立的非 scoped `<style>` 块
- **坑**：伪元素必须挂到 `.fc-daygrid-event` 根元素而非 `.fc-event-title`，因为后者有 `overflow:hidden` 会裁切伪元素
