<template>
  <view class="container">
    <view class="header">
      <text class="back-btn" @click="goBack">‹ 返回</text>
      <text class="title">订单详情</text>
      <view class="header-actions">
        <text v-if="order" class="print-btn" @click="printOrder">🖨</text>
        <text v-if="canModify" class="edit-btn" @click="goEdit">编辑</text>
        <view v-else class="header-right"></view>
      </view>
    </view>
    
    <scroll-view scroll-y class="content" refresher-enabled @refresherrefresh="onRefresh" :refresher-triggered="refreshing">
      <view v-if="loading" class="loading"><text>加载中...</text></view>
      <view v-else-if="order">
        <!-- 基本信息 -->
        <view class="card">
          <view class="card-header">
            <text class="card-title">订单信息</text>
            <text class="order-status" :class="order.OrderStatus">{{ getStatusText(order.OrderStatus) }}</text>
          </view>
          <view class="info-row">
            <text class="info-label">标题</text>
            <text class="info-value">{{ order.Title || '-' }}</text>
          </view>
          <view class="info-row" v-if="order.Link">
            <text class="info-label">链接</text>
            <text class="info-value link" @click="openLink">{{ order.Link }}</text>
          </view>
          <view class="info-row" v-if="order.Styles">
            <text class="info-label">样式</text>
            <text class="info-value">{{ order.Styles }}</text>
          </view>
          <view class="info-row" v-if="order.Remark">
            <text class="info-label">备注</text>
            <text class="info-value">{{ order.Remark }}</text>
          </view>
          <view class="info-row">
            <text class="info-label">创建时间</text>
            <text class="info-value">{{ formatDate(order.CreatedAt) }}</text>
          </view>
        </view>
        
        <!-- 状态切换 -->
        <view class="card">
          <view class="card-header"><text class="card-title">变更状态</text></view>
          <view class="status-buttons">
            <view v-for="opt in statusOptions" :key="opt.value" class="status-btn" :class="order.OrderStatus === opt.value ? 'active ' + opt.value : ''" @click="openStatusDialog(opt.value)">
              <text>{{ opt.label }}</text>
            </view>
          </view>
        </view>
        
        <!-- 费用明细 -->
        <view class="card" v-if="costs.length > 0">
          <view class="card-header"><text class="card-title">费用明细</text></view>
          <view class="cost-table">
            <view class="cost-header">
              <text class="cost-th">类型</text>
              <text class="cost-th">数量</text>
              <text class="cost-th">单价</text>
              <text class="cost-th">小计</text>
            </view>
            <view class="cost-row" v-for="(cost, idx) in costs" :key="idx">
              <text class="cost-td">{{ getCostTypeText(cost.CostType) }}</text>
              <text class="cost-td">{{ cost.Quantity }}</text>
              <text class="cost-td">{{ formatPrice(cost.Price) }}</text>
              <text class="cost-td">{{ formatPrice(cost.Price * cost.Quantity) }}</text>
            </view>
          </view>
          <view class="cost-total">
            <text v-for="g in costsByCurrency" :key="g.currency" class="total-tag">{{ g.currency }} {{ g.total }}</text>
          </view>
        </view>
        
        <!-- 图片 -->
        <view class="card" v-if="photos.length > 0">
          <view class="card-header"><text class="card-title">图片 ({{ photos.length }})</text></view>
          <view class="photo-grid">
            <view v-for="photo in photos" :key="photo.ID" class="photo-item">
              <image :src="getImageUrl(photo.Sha256)" mode="aspectFill" @click="previewImages(photo.Sha256)" />
            </view>
          </view>
        </view>
        
        <!-- 关联工单 -->
        <view class="card" v-if="workOrders.length > 0">
          <view class="card-header"><text class="card-title">关联工单</text></view>
          <view v-for="wo in workOrders" :key="wo.id" class="workorder-item" @click="goToWorkOrder(wo.id)">
            <text class="wo-id">#{{ wo.id }}</text>
            <text class="wo-title">{{ wo.title }}</text>
            <text class="wo-status" :class="wo.status">{{ getWOStatusText(wo.status) }}</text>
          </view>
        </view>
        
        <!-- 状态记录 -->
        <view class="card" v-if="commits.length > 0">
          <view class="card-header"><text class="card-title">状态记录</text></view>
          <view class="timeline">
            <view v-for="commit in commits" :key="commit.id" class="timeline-item">
              <view class="timeline-dot" :class="commit.status"></view>
              <view class="timeline-content">
                <view class="timeline-header">
                  <text class="timeline-status" :class="commit.status">{{ getStatusText(commit.status) }}</text>
                  <text class="timeline-date">{{ formatDate(commit.createdAt) }}</text>
                </view>
                <text class="timeline-user">by {{ getUsernameById(commit.userId) }}</text>
                <text class="timeline-comment" v-if="commit.comment">{{ commit.comment }}</text>
                <view class="timeline-photos" v-if="commit.photos && commit.photos.length">
                  <image v-for="(hash, idx) in commit.photos" :key="hash" class="timeline-photo" :src="getImageUrl(hash)" mode="aspectFill" @click="previewCommitImages(commit.photos, idx)" />
                </view>
              </view>
            </view>
          </view>
        </view>
      </view>
      <view v-else class="empty"><text>未找到订单信息</text></view>
    </scroll-view>
    
    <!-- 状态变更弹窗 -->
    <view class="modal" v-if="statusDialogVisible" @click="closeStatusDialog">
      <view class="modal-content" @click.stop>
        <view class="modal-header">
          <text class="modal-title">变更状态</text>
          <text class="modal-close" @click="closeStatusDialog">×</text>
        </view>
        <view class="modal-body">
          <view class="form-item">
            <text class="form-label">新状态</text>
            <view class="status-display">{{ getStatusText(pendingStatus) }}</view>
          </view>
          <view class="form-item">
            <text class="form-label">备注</text>
            <textarea v-model="pendingComment" class="form-textarea" placeholder="可选" />
          </view>
          <view class="form-item">
            <text class="form-label">图片</text>
            <view class="photo-upload-area">
              <view v-for="(photo, idx) in pendingPhotos" :key="idx" class="uploaded-photo">
                <image :src="photo.url" mode="aspectFill" />
                <view class="photo-remove" @click="pendingPhotos.splice(idx, 1)">×</view>
              </view>
              <view class="add-photo-btn" @click="chooseImage" v-if="pendingPhotos.length < 9">
                <text>+</text>
              </view>
            </view>
          </view>
        </view>
        <view class="modal-footer">
          <text class="btn-cancel" @click="closeStatusDialog">取消</text>
          <text class="btn-confirm" @click="confirmStatusChange">确认</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { purchaseApi } from '@/api/purchase.js'
import { useConfigStore } from '@/stores/config.js'
import { fetchUserInfo, getUsername } from '@/stores/users.js'
import { api } from '@/api/index.js'
import { onLoad, onShow } from '@dcloudio/uni-app'

const configStore = useConfigStore()

const orderId = ref(0)
const loading = ref(false)
const order = ref(null)
const costs = ref([])
const photos = ref([])
const commits = ref([])
const workOrders = ref([])
const canModify = ref(false)

const statusOptions = [
  { value: 'pending', label: '待处理' },
  { value: 'ordered', label: '已下单' },
  { value: 'arrived', label: '已到达' },
  { value: 'received', label: '已收件' },
  { value: 'lost', label: '丢件' },
  { value: 'returned', label: '退件' }
]

const currencyOptions = { 1: 'CNY', 2: 'MOP', 3: 'HKD', 4: 'USD' }
const costTypeOptions = { 1: '单价', 2: '运费' }

const costsByCurrency = computed(() => {
  const groups = {}
  costs.value.forEach(c => {
    const cur = currencyOptions[c.CurrencyType] || 'Unknown'
    const amount = c.Price && c.Quantity ? ((c.Price * c.Quantity) / 100).toFixed(2) : '0.00'
    if (!groups[cur]) groups[cur] = 0
    groups[cur] += parseFloat(amount)
  })
  return Object.entries(groups).map(([currency, total]) => ({ currency, total: total.toFixed(2) }))
})

const statusDialogVisible = ref(false)
const pendingStatus = ref('')
const pendingComment = ref('')
const pendingPhotos = ref([])
const submitting = ref(false)
const refreshing = ref(false)

function getStatusText(status) {
  return statusOptions.find(s => s.value === status)?.label || status
}

function getCurrencyText(type) { return currencyOptions[type] || '-' }
function getCostTypeText(type) { return costTypeOptions[type] || type }
function getWOStatusText(status) {
  return { pending: '待处理', checked: '已检查', parts_ordered: '已下单零件', repaired: '已维修', returned: '已送还', unrepairable: '无法维修' }[status] || status
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

function formatPrice(priceInCents) {
  return priceInCents ? (priceInCents / 100).toFixed(2) : '0.00'
}

function getImageUrl(sha256) {
  return configStore.getFileBaseUrl() + '/api/files/get/' + sha256
}

function previewImages(currentSha) {
  const allPhotos = photos.value || []
  const urls = allPhotos.map(p => getImageUrl(p.Sha256))
  const idx = allPhotos.findIndex(p => p.Sha256 === currentSha)
  uni.previewImage({ urls, current: idx >= 0 ? idx : 0 })
}

function previewCommitImages(photoHashes, currentIndex) {
  const urls = photoHashes.map(h => getImageUrl(h))
  uni.previewImage({ urls, current: currentIndex })
}

function getUsernameById(userId) {
  return getUsername(userId) || `用户${userId}`
}

async function fetchOrderDetail() {
  loading.value = true
  try {
    const res = await purchaseApi.getOrder({ id: orderId.value })
    if (res.errCode === 0 && res.data) {
      order.value = res.data.order || null
      canModify.value = res.data.canModify || false
      costs.value = res.data.costs || []
      photos.value = res.data.photos || []
      commits.value = res.data.commits || []
      workOrders.value = res.data.workOrders || []
      
      // 预取用户信息
      if (order.value?.UserID) fetchUserInfo(order.value.UserID)
      commits.value.forEach(c => {
        if (c.userId) fetchUserInfo(c.userId)
      })
    } else {
      order.value = null
    }
  } catch (e) {
    console.error('获取订单详情失败', e)
    order.value = null
  } finally {
    loading.value = false
  }
}

function openStatusDialog(status) {
  pendingStatus.value = status
  pendingComment.value = ''
  pendingPhotos.value = []
  statusDialogVisible.value = true
}

function closeStatusDialog() {
  statusDialogVisible.value = false
  pendingStatus.value = ''
  pendingComment.value = ''
  pendingPhotos.value = []
}

async function chooseImage() {
  uni.chooseImage({
    count: 9 - pendingPhotos.value.length,
    success: (res) => res.tempFilePaths.forEach(path => uploadImage(path))
  })
}

async function uploadImage(filePath) {
  try {
    const res = await api.upload('/files/upload/image', { uri: filePath, name: 'file' })
    if (res.errCode === 0 && res.data?.hash) {
      pendingPhotos.value.push({ hash: res.data.hash, url: getImageUrl(res.data.hash) })
    }
  } catch (e) {
    console.error('上传失败', e)
  }
}

async function confirmStatusChange() {
  if (submitting.value) return
  submitting.value = true
  try {
    const res = await api.post('/purchase/updatestatus', {
      id: orderId.value,
      status: pendingStatus.value,
      comment: pendingComment.value,
      photos: pendingPhotos.value.map(p => p.hash)
    })
    if (res.errCode === 0) {
      uni.showToast({ title: '状态更新成功', icon: 'success' })
      closeStatusDialog()
      fetchOrderDetail()
      uni.$emit('purchase-refresh')
    } else {
      uni.showToast({ title: res.errMsg || '更新失败', icon: 'none' })
    }
  } catch (e) {
    uni.showToast({ title: '更新失败', icon: 'none' })
  } finally {
    submitting.value = false
  }
}

function goBack() { uni.navigateBack() }

function printOrder() {
  if (!order.value) return

  // #ifndef APP-PLUS
  uni.showToast({ title: '打印功能仅在 App 端可用', icon: 'none' })
  return
  // #endif

  // #ifdef APP-PLUS
  const printer = uni.requireNativePlugin('LcPrinter')

  // 初始化打印机
  printer.initPrinter({})
  printer.setConcentration({ level: 39 })
  printer.setLineSpacing({ spacing: 1 })

  // 标签打印模式（使用黑标定位）
  printer.printEnableMark({ enable: true })

  // 第一行：标题（加粗大字）
  printer.setFontSize({ fontSize: 1 })
  printer.setTextBold({ bold: true })
  printer.printText({ content: (order.value.Title || '(无标题)')+'\n' })
  //printer.printLine({ line_length: 1 })

  // 第二行：备注
  printer.setFontSize({ fontSize: 0 })
  printer.setTextBold({ bold: false })
  printer.printText({ content: '备注: ' + (order.value.Remark || '(无备注)')+'\n' })
  //printer.printLine({ line_length: 1 })

  // 第三行：样式
  printer.printText({ content: '样式: ' + (order.value.Styles || '(无样式)')+'\n' })
  //printer.printLine({ line_length: 1 })

  // 第四行：创建日期
  printer.printText({ content: '日期: ' + formatDate(order.value.CreatedAt) })
  //printer.printLine({ line_length: 1 })

  // 条形码：内容 po:ID，高度 4
  printer.printBarcode({
    text: 'po:' + orderId.value,
    height: 40,
    barcodeType: 73
  })

  printer.printGoToNextMark()
  // #endif
}



function goEdit() {
  uni.navigateTo({ url: `/pages/order/edit-order?id=${orderId.value}` })
}

function goToWorkOrder(id) {
  uni.navigateTo({ url: `/pages/workorder/show-workorder?id=${id}` })
}

function openLink() {
  if (!order.value?.Link) return
  let url = order.value.Link.trim()
  if (!/^https?:\/\//i.test(url)) url = 'https://' + url
  uni.setClipboardData({
    data: url,
    success: () => uni.showToast({ title: '链接已复制', icon: 'success' })
  })
}

async function onRefresh() {
  refreshing.value = true
  await fetchOrderDetail()
  refreshing.value = false
}

onLoad((options) => {
  if (options?.id) {
    orderId.value = parseInt(options.id)
    fetchOrderDetail()
  }
})

// 每次页面显示时刷新数据
onShow(() => {
  if (orderId.value) {
    fetchOrderDetail()
  }
})
</script>

<style scoped>
.container { min-height: 100vh; background-color: #f5f5f5; padding-bottom: 20rpx; }
.header { background-color: #fff; padding: 30rpx; display: flex; align-items: center; }
.back-btn { font-size: 32rpx; color: #007AFF; margin-right: 20rpx; }
.title { font-size: 36rpx; font-weight: bold; color: #333; flex: 1; text-align: center; }
.header-actions { display: flex; align-items: center; gap: 20rpx; }
.print-btn { font-size: 36rpx; color: #007AFF; }
.edit-btn { font-size: 28rpx; color: #007AFF; }
.header-right { width: 60rpx; }
.content { padding: 20rpx; height: calc(100vh - 120rpx); }
.loading, .empty { display: flex; justify-content: center; align-items: center; height: 300rpx; color: #999; }
.card { background-color: #fff; border-radius: 16rpx; margin-bottom: 20rpx; overflow: hidden; }
.card-header { padding: 24rpx 30rpx; border-bottom: 1rpx solid #f0f0f0; display: flex; justify-content: space-between; align-items: center; }
.card-title { font-size: 30rpx; font-weight: bold; color: #333; }
.order-status { font-size: 24rpx; padding: 8rpx 20rpx; border-radius: 20rpx; }
.order-status.pending, .timeline-status.pending, .timeline-dot.pending { background-color: #fff7e6; color: #faad14; }
.order-status.ordered, .timeline-status.ordered, .timeline-dot.ordered { background-color: #e6f7ff; color: #1890ff; }
.order-status.arrived, .timeline-status.arrived, .timeline-dot.arrived { background-color: #f9f0ff; color: #722ed1; }
.order-status.received, .timeline-status.received, .timeline-dot.received { background-color: #f6ffed; color: #52c41a; }
.order-status.lost, .timeline-status.lost, .timeline-dot.lost { background-color: #fff1f0; color: #ff4d4f; }
.order-status.returned, .timeline-status.returned, .timeline-dot.returned { background-color: #f5f5f5; color: #999; }
.info-row { padding: 20rpx 30rpx; display: flex; border-bottom: 1rpx solid #f0f0f0; }
.info-label { width: 160rpx; color: #999; font-size: 28rpx; flex-shrink: 0; }
.info-value { flex: 1; color: #333; font-size: 28rpx; word-break: break-all; }
.info-value.link { color: #007AFF; }
.status-buttons { display: flex; flex-wrap: wrap; padding: 20rpx; gap: 20rpx; }
.status-btn { padding: 16rpx 30rpx; border-radius: 30rpx; background-color: #f5f5f5; color: #666; font-size: 26rpx; }
.status-btn.active { color: #fff; }
.status-btn.active.pending { background-color: #faad14; }
.status-btn.active.ordered { background-color: #1890ff; }
.status-btn.active.arrived { background-color: #722ed1; }
.status-btn.active.received { background-color: #52c41a; }
.status-btn.active.lost { background-color: #ff4d4f; }
.status-btn.active.returned { background-color: #999; }
.cost-table { padding: 0 20rpx 20rpx; }
.cost-header, .cost-row { display: flex; padding: 20rpx 10rpx; }
.cost-header { background-color: #fafafa; border-radius: 8rpx; }
.cost-th, .cost-td { flex: 1; text-align: center; font-size: 24rpx; }
.cost-th { color: #999; font-weight: 500; }
.cost-total { padding: 20rpx; border-top: 1rpx solid #f0f0f0; display: flex; gap: 20rpx; }
.total-tag { background-color: #e6f7ff; color: #1890ff; padding: 8rpx 20rpx; border-radius: 8rpx; font-size: 24rpx; }
.photo-grid { display: flex; flex-wrap: wrap; padding: 20rpx; gap: 20rpx; }
.photo-item { width: 200rpx; height: 200rpx; border-radius: 12rpx; overflow: hidden; }
.photo-item image { width: 100%; height: 100%; }
.workorder-item { padding: 24rpx 30rpx; border-bottom: 1rpx solid #f0f0f0; display: flex; align-items: center; }
.wo-id { color: #999; font-size: 26rpx; margin-right: 16rpx; }
.wo-title { flex: 1; color: #333; font-size: 28rpx; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.wo-status { font-size: 24rpx; padding: 8rpx 16rpx; border-radius: 16rpx; margin-left: 16rpx; }
.wo-status.pending { background-color: #fff7e6; color: #faad14; }
.wo-status.checked { background-color: #e6f7ff; color: #1890ff; }
.wo-status.parts_ordered { background-color: #f9f0ff; color: #722ed1; }
.wo-status.repaired { background-color: #f6ffed; color: #52c41a; }
.wo-status.returned { background-color: #f5f5f5; color: #999; }
.wo-status.unrepairable { background-color: #fff1f0; color: #ff4d4f; }
.timeline { padding: 20rpx 30rpx; }
.timeline-item { display: flex; padding-bottom: 30rpx; position: relative; }
.timeline-item:not(:last-child)::before { content: ''; position: absolute; left: 10rpx; top: 30rpx; bottom: 0; width: 2rpx; background-color: #e5e5e5; }
.timeline-dot { width: 20rpx; height: 20rpx; border-radius: 50%; margin-top: 8rpx; flex-shrink: 0; margin-right: 20rpx; }
.timeline-content { flex: 1; }
.timeline-header { display: flex; justify-content: space-between; align-items: center; }
.timeline-date { font-size: 24rpx; color: #999; }
.timeline-user { display: block; font-size: 24rpx; color: #999; margin-top: 8rpx; }
.timeline-comment { display: block; font-size: 28rpx; color: #333; margin-top: 12rpx; }
.timeline-photos { display: flex; flex-wrap: wrap; gap: 16rpx; margin-top: 16rpx; }
.timeline-photo { width: 100rpx; height: 100rpx; border-radius: 8rpx; }
.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background-color: rgba(0, 0, 0, 0.5); display: flex; align-items: flex-end; z-index: 999; }
.modal-content { width: 100%; background-color: #fff; border-radius: 24rpx 24rpx 0 0; max-height: 80vh; }
.modal-header { padding: 30rpx; display: flex; justify-content: space-between; align-items: center; border-bottom: 1rpx solid #f0f0f0; }
.modal-title { font-size: 32rpx; font-weight: bold; color: #333; }
.modal-close { font-size: 48rpx; color: #999; line-height: 1; }
.modal-body { padding: 30rpx; max-height: 60vh; overflow-y: auto; }
.form-item { margin-bottom: 30rpx; }
.form-label { display: block; font-size: 28rpx; color: #666; margin-bottom: 16rpx; }
.status-display { font-size: 32rpx; color: #333; font-weight: 500; }
.form-textarea { width: 100%; height: 160rpx; border: 1rpx solid #e5e5e5; border-radius: 12rpx; padding: 20rpx; font-size: 28rpx; box-sizing: border-box; }
.photo-upload-area { display: flex; flex-wrap: wrap; gap: 20rpx; }
.uploaded-photo { width: 150rpx; height: 150rpx; border-radius: 12rpx; overflow: hidden; position: relative; }
.uploaded-photo image { width: 100%; height: 100%; }
.photo-remove { position: absolute; top: -10rpx; right: -10rpx; width: 40rpx; height: 40rpx; background-color: rgba(0, 0, 0, 0.5); color: #fff; border-radius: 50%; display: flex; justify-content: center; align-items: center; font-size: 28rpx; }
.add-photo-btn { width: 150rpx; height: 150rpx; border-radius: 12rpx; border: 2rpx dashed #ddd; display: flex; justify-content: center; align-items: center; }
.add-photo-btn text { font-size: 60rpx; color: #999; }
.modal-footer { padding: 30rpx; display: flex; gap: 20rpx; border-top: 1rpx solid #f0f0f0; }
.btn-cancel, .btn-confirm { flex: 1; text-align: center; padding: 24rpx; border-radius: 12rpx; font-size: 30rpx; }
.btn-cancel { background-color: #f5f5f5; color: #666; }
.btn-confirm { background-color: #007AFF; color: #fff; }
</style>
