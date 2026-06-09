<template>
  <view class="container">
    <view class="header">
      <text class="back-btn" @click="goBack">‹ 返回</text>
      <text class="title">工单详情</text>
      <view class="header-right">
        <text v-if="canModify" class="edit-btn" @click="goEdit">编辑</text>
        <text class="print-btn" @click="printWorkOrder">打印</text>
      </view>
    </view>
    
    <scroll-view scroll-y class="content" refresher-enabled @refresherrefresh="onRefresh" :refresher-triggered="refreshing">
      <!-- 加载中 -->
      <view v-if="loading" class="loading">
        <text>加载中...</text>
      </view>
      
      <!-- 未找到 -->
      <view v-else-if="notFound" class="not-found">
        <text>工单不存在</text>
      </view>
      
      <!-- 工单详情 -->
      <view v-else class="order-detail">
        <!-- 基本信息卡片 -->
        <view class="card">
          <view class="card-header">
            <view class="card-title-row">
              <text class="order-id">#{{ orderId }}</text>
              <text class="order-status" :class="order.CurrentStatus">{{ getStatusText(order.CurrentStatus) }}</text>
            </view>
            <text class="order-date">创建: {{ formatDate(order.CreatedAt) }}</text>
          </view>
          
          <view class="info-item">
            <text class="info-label">标题</text>
            <text class="info-value">{{ order.Title || '-' }}</text>
          </view>
          
          <view v-if="order.Description" class="info-item">
            <text class="info-label">描述</text>
            <text class="info-value desc">{{ order.Description }}</text>
          </view>
          
          <!-- 关联物品 -->
          <view v-if="linkedItems.length > 0" class="info-item">
            <text class="info-label">关联物品</text>
            <view class="linked-items">
              <view
                v-for="item in linkedItems"
                :key="item.ID"
                class="linked-item"
                @click="goToItem(item.ID)"
              >
                <text class="linked-item-name">{{ item.Name }}</text>
                <text v-if="item.SerialNumber" class="linked-item-serial">-{{ item.SerialNumber }}</text>
              </view>
            </view>
          </view>
          
          <!-- 关联客户 -->
          <view v-if="linkedCustomers.length > 0" class="info-item">
            <text class="info-label">关联客户</text>
            <view class="linked-customers">
              <view
                v-for="customer in linkedCustomers"
                :key="customer.id"
                class="linked-customer"
              >
                <text class="linked-customer-name">{{ (customer.last_name || '') + (customer.first_name ? ' ' + customer.first_name : '') }}</text>
                <text v-if="customer.primary_phone" class="linked-customer-phone">{{ customer.primary_phone }}</text>
              </view>
            </view>
          </view>
          
          <!-- 关联采购订单汇总（去重） -->
          <view v-if="allPurchaseOrders.length > 0" class="info-item">
            <text class="info-label">关联采购订单</text>
            <view class="linked-pos">
              <view
                v-for="po in allPurchaseOrders"
                :key="po.id"
                class="linked-po"
                @click="goToPurchaseOrder(po.id)"
              >
                <text class="linked-po-id">#{{ po.id }}</text>
                <text class="linked-po-title">{{ po.title || '' }}</text>
                <text class="linked-po-status" :class="po.status">{{ getPurchaseStatusText(po.status) }}</text>
              </view>
            </view>
          </view>
        </view>
        
        <!-- 图片卡片 -->
        <view v-if="photos.length > 0" class="card">
          <text class="card-title">图片</text>
          <view class="photo-grid">
            <view
              v-for="photo in photos"
              :key="photo.ID"
              class="photo-item"
              @click="previewPhoto(photo.Sha256)"
            >
              <image
                class="photo-img"
                :src="getPhotoUrl(photo.Sha256)"
                mode="aspectFill"
              />
            </view>
          </view>
        </view>
        
        <!-- 进度历史卡片 -->
        <view class="card">
          <text class="card-title">进度历史</text>
          
          <view v-if="commits.length === 0" class="empty-text">
            <text>暂无进度记录</text>
          </view>
          
          <view v-else class="commit-list">
            <view
              v-for="(commit, index) in commits"
              :key="commit.ID"
              class="commit-item"
            >
              <view class="commit-dot" :class="commit.Status"></view>
              <view class="commit-content">
                <view class="commit-header">
                  <text class="commit-status" :class="commit.Status">{{ getStatusText(commit.Status) }}</text>
                  <text class="commit-date">{{ formatDateTime(commit.CreatedAt) }}</text>
                </view>
                <text v-if="commit.Comment" class="commit-comment">{{ commit.Comment }}</text>
                <text v-if="commit.Action === 'create'" class="commit-action">创建工单</text>
                
                <!-- 进度图片 -->
                <view v-if="commit.photos && commit.photos.length > 0" class="commit-photos">
                  <view
                    v-for="(photo, idx) in commit.photos"
                    :key="photo.ID"
                    class="commit-photo"
                    @click="previewCommitPhotos(commit.photos, idx)"
                  >
                    <image
                      class="commit-photo-img"
                      :src="getPhotoUrl(photo.Sha256)"
                      mode="aspectFill"
                    />
                  </view>
                </view>
                
                <!-- 关联采购订单 -->
                <view v-if="commit.purchaseOrders && commit.purchaseOrders.length > 0" class="commit-po">
                  <text class="commit-po-label">关联采购订单:</text>
                  <view
                    v-for="po in commit.purchaseOrders"
                    :key="po.id"
                    class="commit-po-item"
                    @click="goToPurchaseOrder(po.id)"
                  >
                    <text class="commit-po-id">#{{ po.id }}</text>
                    <text class="commit-po-title">{{ po.title || '' }}</text>
                    <text class="commit-po-status" :class="po.status">{{ getPurchaseStatusText(po.status) }}</text>
                  </view>
                </view>
              </view>
            </view>
          </view>
        </view>
        
        <!-- 新增进度卡片 -->
        <view v-if="canCommit" class="card">
          <text class="card-title">新增进度</text>
          
          <view class="form-item">
            <text class="form-label">进度状态</text>
            <picker mode="selector" :range="statusOptions" range-key="label" @change="onStatusChange">
              <view class="picker">
                {{ selectedStatus.label || '请选择状态' }}
              </view>
            </picker>
          </view>
          
          <!-- 关联采购订单（仅已下单零件状态显示） -->
          <view v-if="selectedStatus.value === 'parts_ordered'" class="form-item">
            <text class="form-label">关联采购订单</text>
            
            <!-- 已选中的订单 -->
            <view v-if="selectedPurchaseOrders.length > 0" class="selected-pos">
              <view 
                v-for="po in selectedPurchaseOrders" 
                :key="po.id" 
                class="po-tag"
              >
                <text>#{{ po.id }} {{ po.title || '' }}</text>
                <text class="po-tag-remove" @click="removePurchaseOrder(po.id)">×</text>
              </view>
            </view>
            
            <!-- 搜索框 -->
            <view class="po-search-box">
              <input 
                class="po-search-input"
                v-model="purchaseSearchQuery" 
                placeholder="搜索采购订单..."
                @input="onPurchaseSearchInput"
                @focus="onPurchaseSearchFocus"
              />
            </view>
            
            <!-- 搜索结果 -->
            <view v-if="showPurchaseDropdown && purchaseSearchResults.length > 0" class="po-dropdown">
              <view 
                v-for="po in purchaseSearchResults" 
                :key="po.id" 
                class="po-dropdown-item"
                @click="selectPurchaseOrder(po)"
              >
                <text class="po-dropdown-id">#{{ po.id }}</text>
                <text class="po-dropdown-title">{{ po.title || '-' }}</text>
                <text class="po-dropdown-status" :class="po.status">{{ getPurchaseStatusText(po.status) }}</text>
              </view>
            </view>
            <view v-else-if="showPurchaseDropdown && purchaseSearchQuery && purchaseSearchResults.length === 0 && !purchaseSearchLoading" class="po-dropdown-empty">
              未找到匹配的订单
            </view>
          </view>
          
          <view class="form-item">
            <text class="form-label">备注</text>
            <textarea
              v-model="commitComment"
              class="textarea"
              placeholder="请输入备注（选填）"
              maxlength="500"
            />
          </view>
          
          <view class="form-item">
            <text class="form-label">图片</text>
            <view class="photo-upload">
              <view
                v-for="(hash, index) in commitPhotos"
                :key="index"
                class="photo-item"
                @click="previewCommitPhoto(index)"
              >
                <image
                  class="photo-img"
                  :src="getPhotoUrl(hash)"
                  mode="aspectFill"
                />
                <view class="photo-remove" @click.stop="removeCommitPhoto(index)">×</view>
              </view>
              <view v-if="commitPhotos.length < 10" class="photo-add" @click="chooseCommitImage">
                <text class="photo-add-icon">+</text>
                <text class="photo-add-text">添加图片</text>
              </view>
            </view>
          </view>
          
          <view class="submit-btn" :class="{ disabled: submitting || !canSubmit }" @click="submitCommit">
            <text v-if="submitting">提交中...</text>
            <text v-else>提交进度</text>
          </view>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { workOrderApi } from '@/api/work_order.js'
import api from '@/api/index.js'
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
const goBack = () => uni.navigateBack()

// 路由参数
const orderId = ref(null)

// 数据
const loading = ref(true)
const notFound = ref(false)
const order = ref({})
const photos = ref([])
const commits = ref([])
const linkedItems = ref([])
const linkedCustomers = ref([])
const canCommit = ref(false)
const canModify = ref(false)

// 状态选项
const statusOptions = [
  { value: 'pending', label: '待处理' },
  { value: 'checked', label: '已检查' },
  { value: 'parts_ordered', label: '已下单零件' },
  { value: 'repaired', label: '已维修' },
  { value: 'returned', label: '已送还' },
  { value: 'unrepairable', label: '无法维修' }
]

const purchaseStatusOptions = [
  { value: 'pending', label: '待处理' },
  { value: 'ordered', label: '已下单' },
  { value: 'arrived', label: '已到达' },
  { value: 'received', label: '已收件' },
  { value: 'lost', label: '丢件' },
  { value: 'returned', label: '退件' }
]

function getStatusText(status) {
  const option = statusOptions.find(s => s.value === status)
  return option ? option.label : status || ''
}

function getPurchaseStatusText(status) {
  const option = purchaseStatusOptions.find(s => s.value === status)
  return option ? option.label : status || ''
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

function formatDateTime(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${formatDate(dateStr)} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

function getPhotoUrl(hash) {
  return configStore.getFileBaseUrl() + '/api/files/get/' + hash
}

// 获取工单详情
async function fetchOrderDetail() {
  loading.value = true
  try {
    const res = await workOrderApi.get(orderId.value)
    if (res.errCode === 0 && res.data) {
      order.value = res.data.order || {}
      photos.value = res.data.photos || []
      commits.value = res.data.commits || []
      linkedItems.value = res.data.linkedItems || []
      linkedCustomers.value = res.data.linkedCustomers || []
      canCommit.value = res.data.canCommit || false
      canModify.value = res.data.canModify || false
    } else {
      notFound.value = true
    }
  } catch (e) {
    console.error('获取工单详情失败', e)
    notFound.value = true
  } finally {
    loading.value = false
  }
}

// 打印工单标签
function printWorkOrder() {
  if (!printer) {
    uni.showToast({ title: '打印机插件未加载（仅在 App 环境可用）', icon: 'none' })
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
  
  // 第一行：标题
  printer.setFontSize({ fontSize: 0 })
  printer.setTextBold({ bold: true })
  printer.printText({
    content: (order.value.Title || '工单')+'\n'
  })
  printer.setTextBold({ bold: false })
  //printer.printLine({ line_length: 1 })
  
  // 第二行：描述
  if (order.value.Description) {
    //printer.setFontSize({ fontSize: 1 })
    printer.printText({
      content: order.value.Description+'\n'
    })
    //printer.printLine({ line_length: 1 })
  }
  
  // 第三行：创建时间
  if (order.value.CreatedAt) {
    printer.setFontSize({ fontSize: 1 })
    printer.printText({
      content: formatDate(order.value.CreatedAt)
    })
    //printer.printLine({ line_length: 1 })
  }
  
  // 条形码（高40）
  printer.printBarcode({
    text: 'wo:' + orderId.value,
    height: 40,
    barcodeType: 73
  })
  //printer.printLine({ line_length: 2 })
  
  // 提交打印
  printer.start()
  
  uni.showToast({ title: '打印成功', icon: 'success' })
}

// 新增进度
const selectedStatus = ref({})
const commitComment = ref('')
const commitPhotos = ref([])
const submitting = ref(false)

// 采购订单关联
const selectedPurchaseOrders = ref([])
const purchaseSearchQuery = ref('')
const purchaseSearchResults = ref([])
const showPurchaseDropdown = ref(false)
const purchaseSearchLoading = ref(false)
let purchaseSearchTimer = null
const refreshing = ref(false)

const canSubmit = computed(() => {
  return selectedStatus.value.value && (commitComment.value.trim() || commitPhotos.value.length > 0)
})

// 所有 commits 中关联的采购订单（去重）
const allPurchaseOrders = computed(() => {
  const map = new Map()
  for (const commit of commits.value) {
    if (commit.purchaseOrders) {
      for (const po of commit.purchaseOrders) {
        if (!map.has(po.id)) {
          map.set(po.id, po)
        }
      }
    }
  }
  return [...map.values()]
})

function onStatusChange(e) {
  selectedStatus.value = statusOptions[e.detail.value]
}

// 采购订单搜索
async function searchPurchaseOrders() {
  purchaseSearchLoading.value = true
  try {
    const res = await workOrderApi.searchPurchaseOrders(purchaseSearchQuery.value, 10)
    if (res.errCode === 0 && res.data) {
      purchaseSearchResults.value = res.data.orders || res.data || []
    } else {
      purchaseSearchResults.value = []
    }
  } catch (e) {
    console.error('搜索采购订单失败', e)
    purchaseSearchResults.value = []
  } finally {
    purchaseSearchLoading.value = false
  }
}

function onPurchaseSearchInput() {
  clearTimeout(purchaseSearchTimer)
  purchaseSearchTimer = setTimeout(() => {
    searchPurchaseOrders()
  }, 300)
}

function onPurchaseSearchFocus() {
  showPurchaseDropdown.value = true
  searchPurchaseOrders()
}

function selectPurchaseOrder(po) {
  // 检查是否已选中
  const exists = selectedPurchaseOrders.value.find(p => p.id === po.id)
  if (!exists) {
    selectedPurchaseOrders.value.push(po)
  }
  purchaseSearchQuery.value = ''
  showPurchaseDropdown.value = false
  purchaseSearchResults.value = []
}

function removePurchaseOrder(id) {
  const index = selectedPurchaseOrders.value.findIndex(p => p.id === id)
  if (index >= 0) {
    selectedPurchaseOrders.value.splice(index, 1)
  }
}

function previewPhoto(sha256) {
  const urls = photos.value.map(p => getPhotoUrl(p.Sha256))
  const current = urls.findIndex(u => u.includes(sha256))
  uni.previewImage({
    current: current >= 0 ? current : 0,
    urls: urls
  })
}

// 预览进度历史中的图片
function previewCommitPhotos(photos, currentIndex) {
  const urls = photos.map(p => getPhotoUrl(p.Sha256))
  uni.previewImage({ urls, current: currentIndex })
}

function previewCommitPhoto(index) {
  const urls = commitPhotos.value.map(h => getPhotoUrl(h))
  uni.previewImage({
    current: index,
    urls: urls
  })
}

function chooseCommitImage() {
  uni.chooseImage({
    count: 10 - commitPhotos.value.length,
    success: (res) => {
      res.tempFilePaths.forEach(path => {
        uploadCommitImage(path)
      })
    }
  })
}

async function uploadCommitImage(filePath) {
  try {
    uni.showLoading({ title: '上传中...' })
    const res = await api.upload('/files/upload/image', {
      uri: filePath,
      name: 'file'
    })
    
    if (res.errCode === 0 && res.data && res.data.hash) {
      commitPhotos.value.push(res.data.hash)
    } else {
      uni.showToast({ title: '上传失败', icon: 'none' })
    }
  } catch (e) {
    console.error('上传失败', e)
    uni.showToast({ title: '上传失败', icon: 'none' })
  } finally {
    uni.hideLoading()
  }
}

function removeCommitPhoto(index) {
  commitPhotos.value.splice(index, 1)
}

async function submitCommit() {
  if (!selectedStatus.value.value) {
    uni.showToast({ title: '请选择状态', icon: 'none' })
    return
  }
  
  if (!canSubmit.value) {
    uni.showToast({ title: '请填写备注或上传图片', icon: 'none' })
    return
  }
  
  if (submitting.value) return
  submitting.value = true
  
  try {
    const data = {
      id: orderId.value,
      status: selectedStatus.value.value,
      comment: commitComment.value.trim(),
      photos: commitPhotos.value
    }
    
    // 如果选择了采购订单，添加关联
    if (selectedPurchaseOrders.value.length > 0) {
      data.purchaseOrderIds = selectedPurchaseOrders.value.map(po => po.id)
    }
    
    const res = await workOrderApi.commit(data)
    
    if (res.errCode === 0) {
      uni.showToast({ title: '提交成功', icon: 'success' })
      // 重置表单
      commitComment.value = ''
      commitPhotos.value = []
      selectedStatus.value = {}
      selectedPurchaseOrders.value = []
      // 刷新详情
      await fetchOrderDetail()
    } else {
      uni.showToast({ title: res.errMsg || '提交失败', icon: 'none' })
    }
  } catch (e) {
    console.error('提交进度失败', e)
    uni.showToast({ title: '提交失败', icon: 'none' })
  } finally {
    submitting.value = false
  }
}

// 跳转到物品详情
function goToItem(id) {
  uni.navigateTo({
    url: `/pages/warehouse/item-detail?id=${id}`
  })
}

// 跳转到采购订单详情
function goToPurchaseOrder(id) {
  uni.navigateTo({
    url: `/pages/order/order-detail?id=${id}`
  })
}

// 跳转到编辑页面
function goEdit() {
  uni.navigateTo({
    url: `/pages/workorder/edit-workorder?id=${orderId.value}`
  })
}

onMounted(() => {
  initPage()
})

function initPage() {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  const options = currentPage.options || currentPage.$page?.options || {}
  orderId.value = options.id ? parseInt(options.id) : null
  
  if (orderId.value) {
    fetchOrderDetail()
  } else {
    notFound.value = true
    loading.value = false
  }
}

// uni-app 页面显示时刷新数据
uni.$on('page-refresh', () => {
  if (orderId.value) {
    fetchOrderDetail()
  }
})

async function onRefresh() {
  refreshing.value = true
  await fetchOrderDetail()
  refreshing.value = false
}
</script>

<style scoped>
.container {
  min-height: 100vh;
  background-color: #f5f5f5;
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
}

.title {
  font-size: 36rpx;
  font-weight: bold;
  color: #333;
  flex: 1;
  text-align: center;
}

.edit-btn {
  font-size: 28rpx;
  color: #007AFF;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 40rpx;
}

.print-btn {
  font-size: 28rpx;
  color: #007AFF;
}

.content {
  height: calc(100vh - 120rpx);
  padding: 20rpx;
}

.loading, .not-found {
  text-align: center;
  padding: 100rpx 0;
  color: #999;
}

.card {
  background-color: #fff;
  border-radius: 12rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
}

.card-header {
  margin-bottom: 30rpx;
}

.card-title-row {
  display: flex;
  align-items: center;
  gap: 20rpx;
  margin-bottom: 10rpx;
}

.card-title {
  display: block;
  font-size: 32rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 30rpx;
}

.order-id {
  font-size: 36rpx;
  font-weight: bold;
  color: #333;
}

.order-status {
  font-size: 22rpx;
  padding: 6rpx 16rpx;
  border-radius: 20rpx;
}

.order-date {
  font-size: 24rpx;
  color: #999;
}

.info-item {
  margin-bottom: 24rpx;
}

.info-label {
  display: block;
  font-size: 24rpx;
  color: #999;
  margin-bottom: 10rpx;
}

.info-value {
  font-size: 28rpx;
  color: #333;
}

.info-value.desc {
  white-space: pre-wrap;
}

.linked-items {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.linked-item {
  display: inline-flex;
  align-items: center;
  padding: 10rpx 20rpx;
  background-color: #e6f7ff;
  border: 1rpx solid #91d5ff;
  border-radius: 30rpx;
}

.linked-item-name {
  font-size: 24rpx;
  color: #1890ff;
}

.linked-item-serial {
  font-size: 24rpx;
  color: #8c8c8c;
}

.linked-customers {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.linked-customer {
  display: inline-flex;
  align-items: center;
  gap: 8rpx;
  padding: 10rpx 20rpx;
  background-color: #e6f7ff;
  border: 1rpx solid #91d5ff;
  border-radius: 30rpx;
}

.linked-customer-name {
  font-size: 24rpx;
  color: #1890ff;
}

.linked-customer-phone {
  font-size: 22rpx;
  color: #8c8c8c;
  margin-left: 8rpx;
}

/* 关联采购订单 */
.linked-pos {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.linked-po {
  display: inline-flex;
  align-items: center;
  gap: 8rpx;
  padding: 10rpx 20rpx;
  background-color: #e6f0ff;
  border: 1rpx solid #91d5ff;
  border-radius: 30rpx;
}

.linked-po-id {
  font-size: 24rpx;
  font-weight: bold;
  color: #1890ff;
}

.linked-po-title {
  font-size: 24rpx;
  color: #333;
  max-width: 200rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.linked-po-status {
  font-size: 20rpx;
  padding: 4rpx 10rpx;
  border-radius: 10rpx;
}

/* 采购订单状态颜色 */
.linked-po-status.pending {
  background-color: #fef3c7;
  color: #b45309;
}

.linked-po-status.ordered {
  background-color: #dbeafe;
  color: #1d4ed8;
}

.linked-po-status.arrived {
  background-color: #e0e7ff;
  color: #4338ca;
}

.linked-po-status.received {
  background-color: #dcfce7;
  color: #15803d;
}

.linked-po-status.lost,
.linked-po-status.returned {
  background-color: #fee2e2;
  color: #dc2626;
}

/* 状态颜色 */
.order-status.pending,
.commit-status.pending,
.commit-dot.pending {
  background-color: #fef3c7;
  color: #b45309;
}

.order-status.checked,
.commit-status.checked,
.commit-dot.checked {
  background-color: #f3e8ff;
  color: #7c3aed;
}

.order-status.parts_ordered,
.commit-status.parts_ordered,
.commit-dot.parts_ordered {
  background-color: #dbeafe;
  color: #1d4ed8;
}

.order-status.repaired,
.commit-status.repaired,
.commit-dot.repaired {
  background-color: #dcfce7;
  color: #15803d;
}

.order-status.returned,
.commit-status.returned,
.commit-dot.returned {
  background-color: #e5e7eb;
  color: #4b5563;
}

.order-status.unrepairable,
.commit-status.unrepairable,
.commit-dot.unrepairable {
  background-color: #fee2e2;
  color: #dc2626;
}

/* 图片 */
.photo-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.photo-item {
  width: 200rpx;
  height: 200rpx;
  border-radius: 10rpx;
  overflow: hidden;
}

.photo-img {
  width: 100%;
  height: 100%;
}

/* 进度历史 */
.empty-text {
  text-align: center;
  padding: 40rpx 0;
  color: #999;
  font-size: 28rpx;
}

.commit-list {
  position: relative;
}

.commit-item {
  display: flex;
  padding-left: 40rpx;
  position: relative;
  padding-bottom: 40rpx;
}

.commit-item:last-child {
  padding-bottom: 0;
}

.commit-dot {
  position: absolute;
  left: 0;
  top: 8rpx;
  width: 20rpx;
  height: 20rpx;
  border-radius: 50%;
}

.commit-content {
  flex: 1;
}

.commit-header {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-bottom: 10rpx;
}

.commit-status {
  font-size: 22rpx;
  padding: 4rpx 12rpx;
  border-radius: 16rpx;
}

.commit-date {
  font-size: 22rpx;
  color: #999;
}

.commit-comment {
  display: block;
  font-size: 28rpx;
  color: #333;
  line-height: 1.6;
  margin-bottom: 16rpx;
}

.commit-action {
  display: block;
  font-size: 26rpx;
  color: #666;
  font-style: italic;
}

.commit-photos {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
  margin-top: 16rpx;
}

.commit-photo {
  width: 140rpx;
  height: 140rpx;
  border-radius: 8rpx;
  overflow: hidden;
}

.commit-photo-img {
  width: 100%;
  height: 100%;
}

.commit-po {
  margin-top: 16rpx;
}

.commit-po-label {
  font-size: 24rpx;
  color: #999;
}

.commit-po-item {
  display: inline-flex;
  align-items: center;
  gap: 10rpx;
  padding: 10rpx 16rpx;
  background-color: #e6f7ff;
  border: 1rpx solid #91d5ff;
  border-radius: 8rpx;
  margin-top: 10rpx;
}

.commit-po-id {
  font-size: 24rpx;
  color: #1890ff;
  font-weight: bold;
}

.commit-po-title {
  font-size: 24rpx;
  color: #333;
}

.commit-po-status {
  font-size: 20rpx;
  padding: 2rpx 8rpx;
  border-radius: 8rpx;
}

.commit-po-status.pending {
  background-color: #fef3c7;
  color: #b45309;
}

.commit-po-status.ordered {
  background-color: #dbeafe;
  color: #1d4ed8;
}

.commit-po-status.arrived {
  background-color: #f3e8ff;
  color: #7c3aed;
}

.commit-po-status.received {
  background-color: #dcfce7;
  color: #15803d;
}

.commit-po-status.lost {
  background-color: #fee2e2;
  color: #dc2626;
}

.commit-po-status.returned {
  background-color: #e5e7eb;
  color: #4b5563;
}

/* 新增进度表单 */
.form-item {
  margin-bottom: 30rpx;
}

.form-label {
  display: block;
  font-size: 28rpx;
  color: #333;
  margin-bottom: 15rpx;
}

.picker {
  background-color: #f5f5f5;
  padding: 20rpx 30rpx;
  border-radius: 10rpx;
  font-size: 28rpx;
}

.textarea {
  width: 100%;
  min-height: 160rpx;
  padding: 20rpx;
  background-color: #f5f5f5;
  border-radius: 10rpx;
  font-size: 28rpx;
  box-sizing: border-box;
}

.photo-upload {
  display: flex;
  flex-wrap: wrap;
  gap: 20rpx;
}

.photo-add {
  width: 200rpx;
  height: 200rpx;
  background-color: #f5f5f5;
  border: 2rpx dashed #d9d9d9;
  border-radius: 10rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.photo-add-icon {
  font-size: 60rpx;
  color: #999;
  line-height: 1;
}

.photo-add-text {
  font-size: 24rpx;
  color: #999;
  margin-top: 10rpx;
}

.photo-remove {
  position: absolute;
  top: 0;
  right: 0;
  width: 40rpx;
  height: 40rpx;
  background-color: rgba(0, 0, 0, 0.5);
  color: #fff;
  font-size: 32rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.submit-btn {
  background-color: #007AFF;
  color: #fff;
  font-size: 32rpx;
  text-align: center;
  padding: 30rpx;
  border-radius: 12rpx;
}

.submit-btn.disabled {
  background-color: #ccc;
}

/* 关联采购订单 */
.selected-pos {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
  margin-bottom: 16rpx;
}

.po-tag {
  display: inline-flex;
  align-items: center;
  gap: 8rpx;
  padding: 8rpx 16rpx;
  background-color: #e6f0ff;
  border: 1rpx solid #91d5ff;
  border-radius: 20rpx;
  font-size: 24rpx;
  color: #1890ff;
}

.po-tag-remove {
  font-size: 28rpx;
  color: #999;
  margin-left: 4rpx;
}

.po-tag-remove:active {
  color: #f56c6c;
}

.po-search-box {
  margin-top: 12rpx;
}

.po-search-input {
  width: 100%;
  height: 72rpx;
  padding: 0 20rpx;
  background-color: #f5f5f5;
  border-radius: 10rpx;
  font-size: 28rpx;
  box-sizing: border-box;
}

.po-dropdown {
  margin-top: 12rpx;
  background-color: #fff;
  border: 1rpx solid #eee;
  border-radius: 10rpx;
  max-height: 400rpx;
  overflow-y: auto;
}

.po-dropdown-item {
  display: flex;
  align-items: center;
  padding: 20rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.po-dropdown-item:last-child {
  border-bottom: none;
}

.po-dropdown-id {
  font-size: 28rpx;
  font-weight: bold;
  color: #1890ff;
  margin-right: 16rpx;
}

.po-dropdown-title {
  flex: 1;
  font-size: 26rpx;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.po-dropdown-status {
  font-size: 22rpx;
  padding: 4rpx 12rpx;
  border-radius: 8rpx;
  margin-left: 12rpx;
}

.po-dropdown-empty {
  padding: 30rpx;
  text-align: center;
  color: #999;
  font-size: 26rpx;
}
</style>
