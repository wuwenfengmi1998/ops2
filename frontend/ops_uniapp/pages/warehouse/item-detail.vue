<template>
  <view class="container">
    <view class="header">
      <text class="back-btn" @click="goBack">‹ 返回</text>
      <text class="title">物品详情</text>
      <view class="header-right">
        <text class="print-btn" @click="printItem">打印</text>
      </view>
    </view>
    
    <scroll-view scroll-y class="content" refresher-enabled @refresherrefresh="onRefresh" :refresher-triggered="refreshing">
      <!-- 加载状态 -->
      <view v-if="loading" class="loading">
        <text>加载中...</text>
      </view>
      
      <!-- 详情内容 -->
      <view v-else-if="item">
        <!-- 基本信息卡片 -->
        <view class="card">
          <view class="card-header">
            <text class="card-title">基本信息</text>
          </view>
          <view class="info-row">
            <text class="info-label">名称</text>
            <text class="info-value">{{ item.Name }}</text>
          </view>
          <view class="info-row" v-if="item.SerialNumber">
            <text class="info-label">编号</text>
            <text class="info-value">{{ item.SerialNumber }}</text>
          </view>
          <view class="info-row" v-if="item.Specification">
            <text class="info-label">规格</text>
            <text class="info-value">{{ item.Specification }}</text>
          </view>
          <view class="info-row" v-if="item.ContainerBreadcrumb">
            <text class="info-label">位置</text>
            <text class="info-value location">{{ item.ContainerBreadcrumb }}</text>
          </view>
          <view class="info-row" v-if="item.Remark">
            <text class="info-label">备注</text>
            <text class="info-value">{{ item.Remark }}</text>
          </view>
          <view class="info-row">
            <text class="info-label">创建时间</text>
            <text class="info-value">{{ formatDate(item.CreatedAt) }}</text>
          </view>
          <view class="info-row">
            <text class="info-label">更新时间</text>
            <text class="info-value">{{ formatDate(item.UpdatedAt) }}</text>
          </view>
        </view>
        
        <!-- 图片卡片 -->
        <view class="card" v-if="photos.length > 0">
          <view class="card-header">
            <text class="card-title">图片 ({{ photos.length }})</text>
          </view>
          <view class="photo-grid">
            <view 
              class="photo-item" 
              v-for="photo in photos" 
              :key="photo.ID"
            >
              <image 
                class="photo-img" 
                :src="getImageUrl(photo.Sha256)" 
                mode="aspectFill"
                @click="previewAllImages(photo.Sha256)"
              />
            </view>
          </view>
        </view>
        
        <!-- 移动历史 -->
        <view class="card" v-if="commits.length > 0">
          <view class="card-header">
            <text class="card-title">移动历史</text>
          </view>
          <view class="timeline">
            <view 
              class="timeline-item" 
              v-for="commit in commits" 
              :key="commit.ID"
            >
              <view class="timeline-dot"></view>
              <view class="timeline-content">
                <text class="timeline-text">
                  {{ commit.OldContainerBreadcrumb || '无' }} → {{ commit.NewContainerBreadcrumb || '已移除' }}
                </text>
                <text class="timeline-date">{{ formatDate(commit.CreatedAt) }}</text>
              </view>
            </view>
          </view>
        </view>
        
        <!-- 关联工单 -->
        <view class="card" v-if="workOrders.length > 0">
          <view class="card-header">
            <text class="card-title">关联工单</text>
          </view>
          <view 
            class="workorder-item" 
            v-for="wo in workOrders" 
            :key="wo.ID"
            @click="goToWorkOrder(wo.id)"
          >
            <text class="wo-id">#{{ wo.id }}</text>
            <text class="wo-title">{{ wo.title }}</text>
            <text class="wo-status" :class="wo.status">{{ getWorkOrderStatusText(wo.status) }}</text>
          </view>
        </view>
        
        <!-- 关联客户 -->
        <view class="card" v-if="linkedCustomers.length > 0">
          <view class="card-header">
            <text class="card-title">关联客户</text>
          </view>
          <view class="linked-customers">
            <view 
              class="linked-customer-item" 
              v-for="customer in linkedCustomers" 
              :key="customer.ID"
            >
              <text class="customer-name">{{ (customer.last_name || '') + (customer.first_name ? ' ' + customer.first_name : '') }}</text>
              <text v-if="customer.title" class="customer-title">{{ customer.title }}</text>
            </view>
          </view>
        </view>
      </view>
      
      <!-- 空状态 -->
      <view v-else class="empty">
        <text>未找到物品信息</text>
      </view>
    </scroll-view>
    
    <!-- 底部操作栏 -->
    <view class="action-bar" v-if="item">
      <view class="action-btn workorder" @click="addWorkOrder">
        <text>+ 工单</text>
      </view>
      <view class="action-btn edit" @click="editItem">
        <text>编辑</text>
      </view>
      <view class="action-btn move" @click="openMoveModal">
        <text>移动</text>
      </view>
    </view>
    
    <!-- 移动弹窗 -->
    <view class="modal" v-if="moveModalVisible" @click="closeMoveModal">
      <view class="modal-content" @click.stop>
        <view class="modal-header">
          <text class="modal-title">移动物品</text>
          <text class="modal-close" @click="closeMoveModal">×</text>
        </view>
        <view class="modal-body">
          <view class="form-item">
            <text class="form-label">选择容器</text>
            <picker mode="selector" :range="containerOptions" range-key="Title" @change="onContainerChange">
              <view class="picker">
                {{ selectedContainer.Title || '请选择目标容器' }}
              </view>
            </picker>
          </view>
        </view>
        <view class="modal-footer">
          <text class="btn-cancel" @click="closeMoveModal">取消</text>
          <text class="btn-confirm" @click="confirmMove">确定</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { warehouseApi } from '@/api/warehouse.js'
import { useConfigStore } from '@/stores/config.js'

// 仅在 App 环境下加载原生插件
let printer = null
try {
  if (uni.requireNativePlugin) {
    printer = uni.requireNativePlugin('LcPrinter')
  }
} catch (e) {
  console.error('打印机插件加载失败：', e.message)
}

const configStore = useConfigStore()
const itemId = ref(0)
const loading = ref(false)
const item = ref(null)
const photos = ref([])
const commits = ref([])
const workOrders = ref([])
const linkedCustomers = ref([])
const canModify = ref(false)

// 移动相关
const moveModalVisible = ref(false)
const containerOptions = ref([])
const selectedContainer = ref({})
const refreshing = ref(false)

const workOrderStatusOptions = [
  { value: 'pending', label: '待处理' },
  { value: 'checked', label: '已检查' },
  { value: 'parts_ordered', label: '已下单零件' },
  { value: 'repaired', label: '已维修' },
  { value: 'returned', label: '已送还' },
  { value: 'unrepairable', label: '无法维修' }
]

function getWorkOrderStatusText(status) {
  const option = workOrderStatusOptions.find(s => s.value === status)
  return option ? option.label : status
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

function getImageUrl(sha256) {
  const baseUrl = configStore.getFileBaseUrl()
  return `${baseUrl}/api/files/get/${sha256}`
}

function previewAllImages(currentSha) {
  const urls = photos.value.map(p => getImageUrl(p.Sha256))
  const currentIndex = photos.value.findIndex(p => p.Sha256 === currentSha)
  uni.previewImage({
    urls: urls,
    current: currentIndex >= 0 ? currentIndex : 0
  })
}

async function fetchDetail() {
  loading.value = true
  try {
    const res = await warehouseApi.getItem(itemId.value)
    if (res.errCode === 0 && res.data) {
      item.value = res.data.item
      photos.value = res.data.photos || []
      commits.value = res.data.commits || []
      workOrders.value = res.data.work_orders || []
      linkedCustomers.value = res.data.customers || []
      canModify.value = res.data.canModifyItem
    } else {
      uni.showToast({ title: '获取详情失败', icon: 'none' })
    }
  } catch (e) {
    console.error('获取物品详情失败', e)
    uni.showToast({ title: '获取详情失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

async function fetchContainers() {
  try {
    const res = await warehouseApi.listContainer({
      all_levels: true,
      entries: 5000,
      page: 1
    })
    if (res.errCode === 0 && res.data) {
      // 添加根目录选项，并格式化容器数据
      const containers = res.data.containers || []
      containerOptions.value = [
        { ID: 0, Title: '根目录' }
      ].concat(containers.map(c => ({
        ID: c.ID,
        Title: c.Title
      })))
    }
  } catch (e) {
    console.error('获取容器列表失败', e)
  }
}

function openMoveModal() {
  moveModalVisible.value = true
  fetchContainers()
}

function closeMoveModal() {
  moveModalVisible.value = false
  selectedContainer.value = {}
}

function onContainerChange(e) {
  selectedContainer.value = containerOptions.value[e.detail.value]
}

async function confirmMove() {
  try {
    uni.showLoading({ title: '移动中...' })
    // ID 为 0 表示移动到根目录
    const targetId = selectedContainer.value.ID === 0 ? null : selectedContainer.value.ID
    const res = await warehouseApi.moveItem(itemId.value, targetId)
    if (res.errCode === 0) {
      uni.showToast({ title: '移动成功', icon: 'success' })
      moveModalVisible.value = false
      selectedContainer.value = {}
      fetchDetail()
    } else {
      uni.showToast({ title: res.errMsg || '移动失败', icon: 'none' })
    }
  } catch (e) {
    console.error('移动失败', e)
    uni.showToast({ title: '移动失败', icon: 'none' })
  } finally {
    uni.hideLoading()
  }
}

function editItem() {
  uni.navigateTo({
    url: `/pages/warehouse/item-edit?id=${itemId.value}`
  })
}

function addWorkOrder() {
  // 存储预填数据到 storage（参考 web 端实现）
  const prefillData = {
    itemId: itemId.value,
    title: item.value.SerialNumber
      ? `${item.value.Name}-${item.value.SerialNumber}`
      : item.value.Name,
    description: item.value.Remark || '',
  }
  
  // 如果有已关联的客户，自动填充到工单
  if (linkedCustomers.value && linkedCustomers.value.length > 0) {
    prefillData.customer_ids = linkedCustomers.value.map(c => c.ID)
    prefillData.customers = linkedCustomers.value.map(c => ({
      id: c.ID,
      first_name: c.first_name || '',
      last_name: c.last_name || '',
      primary_phone: c.primary_phone || ''
    }))
  }
  
  uni.setStorageSync('prefill_work_order', JSON.stringify(prefillData))
  uni.navigateTo({
    url: '/pages/workorder/add-workorder'
  })
}

function goToWorkOrder(id) {
  uni.navigateTo({
    url: `/pages/workorder/show-workorder?id=${id}`
  })
}

function goBack() {
  uni.navigateBack({ delta: 1 })
}

// 打印物品标签
function printItem() {
  if (!printer) {
    uni.showToast({ title: '打印机插件未加载（仅在 App 环境可用）', icon: 'none' })
    return
  }

  if (!item.value) {
    uni.showToast({ title: '物品信息未加载', icon: 'none' })
    return
  }

  // 初始化打印机
  printer.initPrinter({})
  printer.setConcentration({ level: 39 })
  printer.setLineSpacing({ spacing: 1 })

  // 标签打印，使用黑标
  printer.printEnableMark({
    enable: true
  })

  // 第一行：标题（名称）
  printer.setFontSize({ fontSize: 0 })
  printer.setTextBold({ bold: true })
  printer.printText({
    content: (item.value.Name || '物品')+'\n'
  })
  printer.setTextBold({ bold: false })
  //printer.printLine({ line_length: 1 })

  // 第二行：编号
  if (item.value.SerialNumber) {
    printer.setFontSize({ fontSize: 1 })
    printer.printText({
      content: '序列号：' + item.value.SerialNumber +'\n'
    })
    //printer.printLine({ line_length: 1 })
  }

  // 第三行：备注
  if (item.value.Remark) {
    printer.setFontSize({ fontSize: 1 })
    printer.printText({
      content: '备注：' + item.value.Remark+'\n'
    })
    //printer.printLine({ line_length: 1 })
  }

  // 第四行：创建日期
  if (item.value.CreatedAt) {
    printer.setFontSize({ fontSize: 1 })
    printer.printText({
      content: '创建日期：' + formatDate(item.value.CreatedAt)
    })
    //printer.printLine({ line_length: 1 })
  }

  // 条形码（高度4）
  printer.printBarcode({
    text: 'item:' + itemId.value,
    height: 40,
    barcodeType: 73
  })
  printer.printLine({ line_length: 2 })

  // 提交打印
  printer.start()

  uni.showToast({ title: '打印成功', icon: 'success' })
}

async function onRefresh() {
  refreshing.value = true
  await fetchDetail()
  refreshing.value = false
}

onMounted(() => {
  // 备用：在某些情况下 onLoad 可能未触发
})

// 使用 onLoad 获取页面参数（推荐方式）
import { onLoad, onShow } from '@dcloudio/uni-app'
onLoad((options) => {
  console.log('item-detail onLoad options:', options)
  if (options && options.id) {
    itemId.value = parseInt(options.id)
    fetchDetail()
  }
})

// 每次显示页面时刷新数据（从编辑页返回时自动更新）
onShow(() => {
  if (itemId.value > 0) {
    fetchDetail()
  }
})
</script>

<style scoped>
.container {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 120rpx;
}

.header {
  background-color: #fff;
  padding: 30rpx;
  display: flex;
  align-items: center;
}

.back-btn {
  font-size: 32rpx;
  color: #007AFF;
  margin-right: 20rpx;
}

.title {
  font-size: 36rpx;
  font-weight: bold;
  color: #333;
  flex: 1;
  text-align: center;
}

.header-right {
  display: flex;
  align-items: center;
}

.print-btn {
  font-size: 28rpx;
  color: #007AFF;
}

.content {
  height: calc(100vh - 180rpx);
}

.loading, .empty {
  text-align: center;
  padding: 100rpx 0;
  color: #999;
}

.card {
  background-color: #fff;
  margin: 20rpx;
  border-radius: 12rpx;
  overflow: hidden;
}

.card-header {
  padding: 25rpx 30rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.card-title {
  font-size: 30rpx;
  font-weight: bold;
  color: #333;
}

.info-row {
  display: flex;
  padding: 20rpx 30rpx;
  border-bottom: 1rpx solid #f5f5f5;
}

.info-row:last-child {
  border-bottom: none;
}

.info-label {
  width: 150rpx;
  font-size: 28rpx;
  color: #666;
}

.info-value {
  flex: 1;
  font-size: 28rpx;
  color: #333;
}

.info-value.location {
  color: #007AFF;
}

.photo-grid {
  padding: 20rpx 30rpx;
  display: flex;
  flex-wrap: wrap;
  gap: 15rpx;
}

.photo-item {
  width: calc((100% - 30rpx) / 3);
  height: 200rpx;
  border-radius: 10rpx;
  overflow: hidden;
}

.photo-img {
  width: 100%;
  height: 100%;
}

.timeline {
  padding: 20rpx 30rpx;
}

.timeline-item {
  display: flex;
  padding: 15rpx 0;
}

.timeline-dot {
  width: 16rpx;
  height: 16rpx;
  background-color: #007AFF;
  border-radius: 50%;
  margin-right: 20rpx;
  margin-top: 8rpx;
}

.timeline-content {
  flex: 1;
}

.timeline-text {
  display: block;
  font-size: 28rpx;
  color: #333;
  margin-bottom: 8rpx;
}

.timeline-date {
  display: block;
  font-size: 24rpx;
  color: #999;
}

.workorder-item {
  display: flex;
  align-items: center;
  padding: 20rpx 30rpx;
  border-bottom: 1rpx solid #f5f5f5;
}

.workorder-item:last-child {
  border-bottom: none;
}

.wo-id {
  font-size: 28rpx;
  color: #007AFF;
  margin-right: 15rpx;
}

.wo-title {
  flex: 1;
  font-size: 28rpx;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.wo-status {
  font-size: 22rpx;
  padding: 6rpx 16rpx;
  border-radius: 20rpx;
  background-color: #f0f0f0;
  color: #666;
}

.wo-status.pending {
  background-color: #fff3e0;
  color: #ff9800;
}

.wo-status.checked {
	background-color: #f3e5f5;
	color: #9c27b0;
  
}

.wo-status.parts_ordered {
	background-color: #e3f2fd;
	color: #2196f3;
}

.wo-status.repaired {
  background-color: #e8f5e9;
  color: #4caf50;
}

.wo-status.returned {
  background-color: #e5e5e5;
  color: #7a7a7a;
}

.wo-status.unrepairable {
  background-color: #f2bdbe;
  color: #ff575a;
}

.linked-customers {
  padding: 20rpx 30rpx;
}

.linked-customer-item {
  display: flex;
  align-items: center;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.linked-customer-item:last-child {
  border-bottom: none;
}

.customer-name {
  font-size: 28rpx;
  color: #333;
  flex: 1;
}

.customer-title {
  font-size: 24rpx;
  color: #999;
  margin-left: 16rpx;
}

.action-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  background-color: #fff;
  padding: 20rpx 30rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.1);
}

.action-btn {
  flex: 1;
  text-align: center;
  padding: 25rpx;
  border-radius: 10rpx;
  font-size: 30rpx;
}

.action-btn.workorder {
  background-color: #4CAF50;
  color: #fff;
  margin-right: 20rpx;
}

.action-btn.edit {
  background-color: #007AFF;
  color: #fff;
  margin-right: 20rpx;
}

.action-btn.move {
  background-color: #ff9800;
  color: #fff;
}

/* 弹窗样式 */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  width: 85%;
  background-color: #fff;
  border-radius: 16rpx;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 30rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.modal-title {
  font-size: 32rpx;
  font-weight: bold;
  color: #333;
}

.modal-close {
  font-size: 50rpx;
  color: #999;
}

.modal-body {
  padding: 30rpx;
}

.form-item {
  margin-bottom: 20rpx;
}

.form-label {
  display: block;
  font-size: 28rpx;
  color: #666;
  margin-bottom: 15rpx;
}

.picker {
  background-color: #f5f5f5;
  padding: 25rpx;
  border-radius: 10rpx;
  font-size: 28rpx;
}

.modal-footer {
  display: flex;
  border-top: 1rpx solid #f0f0f0;
}

.btn-cancel, .btn-confirm {
  flex: 1;
  text-align: center;
  padding: 30rpx;
  font-size: 30rpx;
}

.btn-cancel {
  color: #666;
  border-right: 1rpx solid #f0f0f0;
}

.btn-confirm {
  color: #007AFF;
}
</style>
