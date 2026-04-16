<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { purchaseApi } from '../../@api/purchase.js'
import { userStore } from '../../store/user.js'

// ───────────── 状态定义 ─────────────
const loading = ref(false)
const orders = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const hasMore = ref(true)
const loadingMore = ref(false)

// 搜索和过滤
const searchText = ref('')
const activeStatus = ref('')  // '' = 全部

// 订单详情
const showDetail = ref(false)
const detailLoading = ref(false)
const detailOrder = ref(null)
const detailCosts = ref([])
const detailPhotos = ref([])
const detailCommits = ref([])
const detailCanModify = ref(false)

// 状态变更弹窗
const showStatusModal = ref(false)
const pendingStatus = ref('')
const pendingComment = ref('')
const updatingStatus = ref(false)

// 统计数字
const countMap = reactive({
  pending: 0, ordered: 0, arrived: 0, received: 0, lost: 0, returned: 0,
})

// ───────────── 常量 ─────────────
const statusOptions = [
  { value: '',         label: '全部',   color: '#6b7280' },
  { value: 'pending',  label: '待订购', color: '#d97706' },
  { value: 'ordered',  label: '已下单', color: '#2563eb' },
  { value: 'arrived',  label: '已到货', color: '#7c3aed' },
  { value: 'received', label: '已收货', color: '#059669' },
  { value: 'lost',     label: '已丢失', color: '#dc2626' },
  { value: 'returned', label: '已退回', color: '#9ca3af' },
]

const currencyMap = { 1: 'CNY', 2: 'MOP', 3: 'HKD', 4: 'USD' }
const costTypeMap = { 1: '单价', 2: '运费' }

// ───────────── 工具函数 ─────────────
function getStatusLabel(val) {
  return statusOptions.find(o => o.value === val)?.label || val
}

function getStatusColor(val) {
  return statusOptions.find(o => o.value === val)?.color || '#6b7280'
}

function formatDate(str) {
  if (!str) return '-'
  const d = new Date(str)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

function formatPrice(cents) {
  if (!cents) return '0.00'
  return (cents / 100).toFixed(2)
}

function getPhotoUrl(hash) {
  return `/api/files/get/${hash}`
}

// 按货币分组合计
const costsByCurrency = computed(() => {
  const groups = {}
  detailCosts.value.forEach(c => {
    const cur = currencyMap[c.CurrencyType] || 'Unknown'
    const amount = c.Price && c.Quantity ? (c.Price * c.Quantity / 100) : 0
    groups[cur] = (groups[cur] || 0) + amount
  })
  return Object.entries(groups).map(([currency, total]) => ({
    currency,
    total: total.toFixed(2),
  }))
})

// ───────────── 数据加载 ─────────────
async function loadOrderCount() {
  try {
    const { errCode, data } = await purchaseApi.getOrderCount()
    if (errCode === 0 && data) {
      Object.assign(countMap, data)
    }
  } catch {}
}

async function loadOrders(reset = true) {
  if (reset) {
    page.value = 1
    hasMore.value = true
    orders.value = []
  }
  if (!hasMore.value) return
  loading.value = reset
  loadingMore.value = !reset
  try {
    const { errCode, data } = await purchaseApi.getOrders({
      page: page.value,
      pageSize,
      search: searchText.value.trim(),
      status: activeStatus.value,
    })
    if (errCode === 0 && data) {
      const list = data.list || []
      total.value = data.total || 0
      if (reset) {
        orders.value = list
      } else {
        orders.value.push(...list)
      }
      hasMore.value = orders.value.length < total.value
      page.value++
    }
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

function onSearch() {
  loadOrders(true)
}

function setStatus(val) {
  activeStatus.value = val
  loadOrders(true)
}

function loadMore() {
  if (!loadingMore.value && hasMore.value) {
    loadOrders(false)
  }
}

// ───────────── 订单详情 ─────────────
async function openDetail(order) {
  showDetail.value = true
  detailLoading.value = true
  detailOrder.value = null
  try {
    const { errCode, data } = await purchaseApi.getOrder(order.ID)
    if (errCode === 0 && data) {
      detailOrder.value = data.order ?? null
      detailCanModify.value = data.canModify ?? false
      detailCosts.value = data.costs ?? []
      detailPhotos.value = data.photos ?? []
      detailCommits.value = data.commits ?? []
    }
  } finally {
    detailLoading.value = false
  }
}

function closeDetail() {
  showDetail.value = false
}

// ───────────── 状态变更 ─────────────
function openStatusModal(status) {
  if (status === detailOrder.value?.OrderStatus) return
  pendingStatus.value = status
  pendingComment.value = ''
  showStatusModal.value = true
}

function closeStatusModal() {
  showStatusModal.value = false
  pendingStatus.value = ''
  pendingComment.value = ''
}

async function confirmStatusChange() {
  if (!pendingStatus.value) return
  updatingStatus.value = true
  try {
    const { errCode } = await purchaseApi.updateStatus({
      id: detailOrder.value.ID,
      status: pendingStatus.value,
      comment: pendingComment.value,
    })
    if (errCode === 0) {
      uni.showToast({ title: '状态已更新', icon: 'success' })
      closeStatusModal()
      // 刷新详情
      await openDetail(detailOrder.value)
      // 刷新列表
      loadOrders(true)
      loadOrderCount()
    } else {
      uni.showToast({ title: '操作失败', icon: 'none' })
    }
  } catch {
    uni.showToast({ title: '网络错误', icon: 'none' })
  } finally {
    updatingStatus.value = false
  }
}

// ───────────── 生命周期 ─────────────
onMounted(() => {
  if (!userStore.isLoggedIn) {
    uni.reLaunch({ url: '/pages/signin' })
    return
  }
  loadOrderCount()
  loadOrders()
})
</script>

<template>
  <view class="page-wrap">
    <!-- 顶部搜索栏 -->
    <view class="search-bar">
      <view class="search-input-wrap">
        <text class="search-icon">🔍</text>
        <input
          v-model="searchText"
          class="search-input"
          placeholder="搜索订单名称…"
          confirm-type="search"
          @confirm="onSearch"
        />
        <text v-if="searchText" class="search-clear" @tap="searchText = ''; onSearch()">✕</text>
      </view>
    </view>

    <!-- 状态过滤标签 -->
    <scroll-view scroll-x class="status-tabs">
      <view class="status-tabs-inner">
        <view
          v-for="opt in statusOptions"
          :key="opt.value"
          class="status-tab"
          :class="{ 'status-tab-active': activeStatus === opt.value }"
          :style="activeStatus === opt.value ? { borderColor: opt.color, color: opt.color } : {}"
          @tap="setStatus(opt.value)"
        >
          {{ opt.label }}
          <text v-if="opt.value && countMap[opt.value]" class="tab-badge">
            {{ countMap[opt.value] }}
          </text>
        </view>
      </view>
    </scroll-view>

    <!-- 加载中 -->
    <view v-if="loading" class="loading-box">
      <text class="loading-text">加载中…</text>
    </view>

    <!-- 空状态 -->
    <view v-else-if="!orders.length" class="empty-box">
      <text class="empty-icon">📦</text>
      <text class="empty-text">暂无采购订单</text>
    </view>

    <!-- 订单列表 -->
    <scroll-view v-else scroll-y class="list-scroll" @scrolltolower="loadMore">
      <view
        v-for="order in orders"
        :key="order.ID"
        class="order-card"
        @tap="openDetail(order)"
      >
        <view class="order-card-top">
          <text class="order-title">{{ order.Title }}</text>
          <view
            class="order-status-badge"
            :style="{ backgroundColor: getStatusColor(order.OrderStatus) + '20', color: getStatusColor(order.OrderStatus) }"
          >
            {{ getStatusLabel(order.OrderStatus) }}
          </view>
        </view>
        <view class="order-card-meta">
          <text class="meta-text">#{{ order.ID }}</text>
          <text class="meta-sep">·</text>
          <text class="meta-text">{{ formatDate(order.CreatedAt) }}</text>
        </view>
      </view>

      <view v-if="loadingMore" class="load-more-tip">
        <text class="loading-text">加载更多…</text>
      </view>
      <view v-else-if="!hasMore && orders.length" class="load-more-tip">
        <text style="color:#d1d5db;font-size:24rpx;">已加载全部 {{ total }} 条</text>
      </view>
      <view style="height:40rpx;" />
    </scroll-view>

    <!-- ─────────────────────────────────────────
         订单详情底部抽屉
    ───────────────────────────────────────── -->
    <view v-if="showDetail" class="drawer-mask" @tap.self="closeDetail">
      <view class="drawer-box">
        <!-- 抽屉头部 -->
        <view class="drawer-header">
          <text class="drawer-title">订单详情</text>
          <text class="drawer-close" @tap="closeDetail">✕</text>
        </view>

        <!-- 加载中 -->
        <view v-if="detailLoading" class="loading-box" style="height:300rpx;">
          <text class="loading-text">加载中…</text>
        </view>

        <scroll-view v-else-if="detailOrder" scroll-y class="drawer-scroll">
          <!-- 基本信息 -->
          <view class="detail-section">
            <view class="detail-row">
              <text class="detail-label">名称</text>
              <text class="detail-value">{{ detailOrder.Title }}</text>
            </view>
            <view class="detail-row">
              <text class="detail-label">当前状态</text>
              <view
                class="order-status-badge"
                :style="{ backgroundColor: getStatusColor(detailOrder.OrderStatus) + '20', color: getStatusColor(detailOrder.OrderStatus) }"
              >
                {{ getStatusLabel(detailOrder.OrderStatus) }}
              </view>
            </view>
            <view v-if="detailOrder.Styles" class="detail-row">
              <text class="detail-label">款式备注</text>
              <text class="detail-value">{{ detailOrder.Styles }}</text>
            </view>
            <view v-if="detailOrder.Remark" class="detail-row">
              <text class="detail-label">备注</text>
              <text class="detail-value">{{ detailOrder.Remark }}</text>
            </view>
            <view v-if="detailOrder.Link" class="detail-row">
              <text class="detail-label">链接</text>
              <text class="detail-value detail-link" @tap="uni.setClipboardData({ data: detailOrder.Link, success: () => uni.showToast({ title: '链接已复制', icon: 'none' }) })">
                {{ detailOrder.Link }}（点击复制）
              </text>
            </view>
          </view>

          <!-- 状态切换（仅可修改者显示） -->
          <view v-if="detailCanModify" class="detail-section">
            <text class="section-title">变更状态</text>
            <view class="status-change-grid">
              <view
                v-for="opt in statusOptions.slice(1)"
                :key="opt.value"
                class="status-change-btn"
                :class="{ 'status-change-active': detailOrder.OrderStatus === opt.value }"
                :style="detailOrder.OrderStatus === opt.value
                  ? { backgroundColor: opt.color, color: '#fff' }
                  : { borderColor: opt.color, color: opt.color }"
                @tap="openStatusModal(opt.value)"
              >
                {{ opt.label }}
              </view>
            </view>
          </view>

          <!-- 费用明细 -->
          <view v-if="detailCosts.length" class="detail-section">
            <text class="section-title">费用明细</text>
            <view v-for="c in detailCosts" :key="c.ID" class="cost-row">
              <text class="cost-type">{{ costTypeMap[c.CostType] || c.CostType }}</text>
              <text class="cost-num">x{{ c.Quantity }}</text>
              <text class="cost-price">{{ formatPrice(c.Price) }}</text>
              <text class="cost-cur">{{ currencyMap[c.CurrencyType] || '-' }}</text>
            </view>
            <!-- 合计 -->
            <view class="cost-total-row">
              <text class="cost-total-label">合计</text>
              <view class="cost-total-values">
                <text v-for="g in costsByCurrency" :key="g.currency" class="cost-total-tag">
                  {{ g.currency }} {{ g.total }}
                </text>
              </view>
            </view>
          </view>

          <!-- 图片备注 -->
          <view v-if="detailPhotos.length" class="detail-section">
            <text class="section-title">图片备注</text>
            <view class="photo-grid">
              <image
                v-for="p in detailPhotos"
                :key="p.ID"
                :src="getPhotoUrl(p.Sha256)"
                mode="aspectFill"
                class="photo-thumb"
                @tap="uni.previewImage({ urls: detailPhotos.map(ph => getPhotoUrl(ph.Sha256)), current: getPhotoUrl(p.Sha256) })"
              />
            </view>
          </view>

          <!-- 状态记录 -->
          <view v-if="detailCommits.length" class="detail-section">
            <text class="section-title">变更记录</text>
            <view v-for="c in detailCommits" :key="c.id" class="commit-row">
              <view
                class="commit-dot"
                :style="{ backgroundColor: getStatusColor(c.status) }"
              />
              <view class="commit-body">
                <view class="commit-top">
                  <view
                    class="commit-badge"
                    :style="{ backgroundColor: getStatusColor(c.status) + '20', color: getStatusColor(c.status) }"
                  >{{ getStatusLabel(c.status) }}</view>
                  <text class="commit-date">{{ formatDate(c.createdAt) }}</text>
                </view>
                <text v-if="c.comment" class="commit-comment">{{ c.comment }}</text>
                <view v-if="c.photos?.length" class="commit-photos">
                  <image
                    v-for="hash in c.photos"
                    :key="hash"
                    :src="getPhotoUrl(hash)"
                    mode="aspectFill"
                    class="commit-photo"
                    @tap="uni.previewImage({ urls: c.photos.map(h => getPhotoUrl(h)), current: getPhotoUrl(hash) })"
                  />
                </view>
              </view>
            </view>
          </view>

          <view style="height:60rpx;" />
        </scroll-view>
      </view>
    </view>

    <!-- ─────────────────────────────────────────
         状态变更弹窗
    ───────────────────────────────────────── -->
    <view v-if="showStatusModal" class="modal-mask" @tap.self="closeStatusModal">
      <view class="modal-box">
        <view class="modal-header">
          <text class="modal-title">变更状态</text>
          <text class="modal-close" @tap="closeStatusModal">✕</text>
        </view>
        <view class="modal-body">
          <view class="form-item">
            <text class="form-label">目标状态</text>
            <view
              class="status-preview"
              :style="{ backgroundColor: getStatusColor(pendingStatus) + '20', color: getStatusColor(pendingStatus) }"
            >
              {{ getStatusLabel(pendingStatus) }}
            </view>
          </view>
          <view class="form-item">
            <text class="form-label">变更备注（可选）</text>
            <textarea
              v-model="pendingComment"
              class="form-textarea"
              placeholder="填写变更原因或说明"
              maxlength="500"
              auto-height
            />
          </view>
        </view>
        <view class="modal-footer" style="justify-content: flex-end;">
          <view class="footer-right">
            <button class="btn-cancel" @tap="closeStatusModal">取消</button>
            <button class="btn-primary" :disabled="updatingStatus" @tap="confirmStatusChange">
              {{ updatingStatus ? '提交中…' : '确认变更' }}
            </button>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<style scoped>
.page-wrap {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #f3f4f6;
}

/* 搜索栏 */
.search-bar {
  background: #fff;
  padding: 48rpx 24rpx 16rpx;
  box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.06);
}
.search-input-wrap {
  display: flex;
  align-items: center;
  background: #f3f4f6;
  border-radius: 48rpx;
  padding: 12rpx 24rpx;
  gap: 12rpx;
}
.search-icon { font-size: 28rpx; color: #9ca3af; }
.search-input { flex: 1; font-size: 28rpx; color: #111827; background: transparent; }
.search-clear { font-size: 28rpx; color: #9ca3af; padding: 4rpx; }

/* 状态标签 */
.status-tabs {
  background: #fff;
  white-space: nowrap;
  border-bottom: 1rpx solid #f3f4f6;
}
.status-tabs-inner {
  display: inline-flex;
  padding: 0 16rpx;
  gap: 8rpx;
}
.status-tab {
  display: inline-flex;
  align-items: center;
  gap: 6rpx;
  padding: 16rpx 20rpx;
  font-size: 26rpx;
  color: #6b7280;
  border-bottom: 4rpx solid transparent;
  white-space: nowrap;
}
.status-tab-active {
  font-weight: 600;
  border-bottom-color: currentColor;
}
.tab-badge {
  background: #f3f4f6;
  color: inherit;
  font-size: 20rpx;
  border-radius: 20rpx;
  padding: 2rpx 10rpx;
}

/* 加载/空状态 */
.loading-box, .empty-box {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16rpx;
}
.loading-text { color: #9ca3af; font-size: 28rpx; }
.empty-icon { font-size: 80rpx; }
.empty-text { color: #9ca3af; font-size: 28rpx; }

/* 订单列表 */
.list-scroll { flex: 1; padding: 20rpx; }
.order-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
  box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.06);
}
.order-card-top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16rpx;
  margin-bottom: 12rpx;
}
.order-title {
  flex: 1;
  font-size: 30rpx;
  font-weight: 600;
  color: #111827;
}
.order-status-badge {
  font-size: 22rpx;
  font-weight: 500;
  padding: 4rpx 16rpx;
  border-radius: 20rpx;
  flex-shrink: 0;
}
.order-card-meta {
  display: flex;
  align-items: center;
  gap: 8rpx;
}
.meta-text { font-size: 24rpx; color: #9ca3af; }
.meta-sep { font-size: 24rpx; color: #d1d5db; }

.load-more-tip {
  display: flex;
  justify-content: center;
  padding: 16rpx 0;
}

/* ─── 抽屉 ─── */
.drawer-mask {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
  z-index: 100;
  display: flex;
  align-items: flex-end;
}
.drawer-box {
  width: 100%;
  background: #fff;
  border-radius: 32rpx 32rpx 0 0;
  height: 85vh;
  display: flex;
  flex-direction: column;
}
.drawer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 32rpx;
  border-bottom: 1rpx solid #f3f4f6;
  flex-shrink: 0;
}
.drawer-title { font-size: 32rpx; font-weight: 600; color: #111827; }
.drawer-close { font-size: 36rpx; color: #9ca3af; }
.drawer-scroll { flex: 1; }

/* 详情区块 */
.detail-section {
  padding: 24rpx 32rpx;
  border-bottom: 1rpx solid #f9fafb;
}
.section-title {
  display: block;
  font-size: 26rpx;
  color: #6b7280;
  font-weight: 600;
  margin-bottom: 16rpx;
}
.detail-row {
  display: flex;
  align-items: flex-start;
  gap: 16rpx;
  margin-bottom: 16rpx;
}
.detail-label {
  font-size: 26rpx;
  color: #9ca3af;
  width: 120rpx;
  flex-shrink: 0;
}
.detail-value {
  flex: 1;
  font-size: 28rpx;
  color: #111827;
  word-break: break-all;
}
.detail-link { color: #2563eb; }

/* 状态切换网格 */
.status-change-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16rpx;
}
.status-change-btn {
  text-align: center;
  padding: 16rpx 8rpx;
  border-radius: 12rpx;
  font-size: 26rpx;
  border: 2rpx solid;
}
.status-change-active { font-weight: 600; }

/* 费用明细 */
.cost-row {
  display: flex;
  gap: 16rpx;
  padding: 12rpx 0;
  border-bottom: 1rpx solid #f9fafb;
  font-size: 28rpx;
}
.cost-type { flex: 1; color: #374151; }
.cost-num { color: #9ca3af; }
.cost-price { font-weight: 500; color: #111827; }
.cost-cur { color: #6b7280; width: 80rpx; text-align: right; }
.cost-total-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 16rpx;
}
.cost-total-label { font-size: 26rpx; color: #6b7280; font-weight: 600; }
.cost-total-values { display: flex; gap: 8rpx; flex-wrap: wrap; }
.cost-total-tag {
  background: #dbeafe;
  color: #1d4ed8;
  font-size: 22rpx;
  font-weight: 600;
  padding: 4rpx 16rpx;
  border-radius: 20rpx;
}

/* 图片网格 */
.photo-grid { display: flex; flex-wrap: wrap; gap: 12rpx; }
.photo-thumb {
  width: 140rpx;
  height: 140rpx;
  border-radius: 12rpx;
  background: #f3f4f6;
}

/* 变更记录 */
.commit-row {
  display: flex;
  gap: 16rpx;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #f9fafb;
}
.commit-dot {
  width: 16rpx;
  height: 16rpx;
  border-radius: 50%;
  flex-shrink: 0;
  margin-top: 10rpx;
}
.commit-body { flex: 1; min-width: 0; }
.commit-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8rpx;
}
.commit-badge {
  font-size: 22rpx;
  padding: 4rpx 12rpx;
  border-radius: 20rpx;
}
.commit-date { font-size: 22rpx; color: #9ca3af; }
.commit-comment { font-size: 26rpx; color: #374151; display: block; margin-top: 4rpx; }
.commit-photos { display: flex; flex-wrap: wrap; gap: 8rpx; margin-top: 8rpx; }
.commit-photo {
  width: 80rpx;
  height: 80rpx;
  border-radius: 8rpx;
  background: #f3f4f6;
}

/* ─── 模态框 ─── */
.modal-mask {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
  z-index: 200;
  display: flex;
  align-items: flex-end;
}
.modal-box {
  width: 100%;
  background: #fff;
  border-radius: 32rpx 32rpx 0 0;
  display: flex;
  flex-direction: column;
}
.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 32rpx 32rpx 24rpx;
  border-bottom: 1rpx solid #f3f4f6;
}
.modal-title { font-size: 32rpx; font-weight: 600; color: #111827; }
.modal-close { font-size: 36rpx; color: #9ca3af; }
.modal-body { padding: 24rpx 32rpx; }
.form-item { margin-bottom: 28rpx; }
.form-label {
  display: block;
  font-size: 26rpx;
  color: #6b7280;
  margin-bottom: 12rpx;
}
.status-preview {
  display: inline-flex;
  padding: 8rpx 24rpx;
  border-radius: 20rpx;
  font-size: 28rpx;
  font-weight: 600;
}
.form-textarea {
  width: 100%;
  min-height: 120rpx;
  border: 2rpx solid #e5e7eb;
  border-radius: 12rpx;
  padding: 16rpx 20rpx;
  font-size: 28rpx;
  color: #111827;
  background: #fff;
  box-sizing: border-box;
}
.modal-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20rpx 32rpx 48rpx;
  border-top: 1rpx solid #f3f4f6;
}
.footer-right { display: flex; gap: 16rpx; margin-left: auto; }
.btn-cancel {
  background: #f3f4f6;
  color: #374151;
  border: none;
  border-radius: 12rpx;
  padding: 16rpx 32rpx;
  font-size: 28rpx;
}
.btn-primary {
  background: #2563eb;
  color: #fff;
  border: none;
  border-radius: 12rpx;
  padding: 16rpx 32rpx;
  font-size: 28rpx;
}
button[disabled] { opacity: 0.5; }
</style>
