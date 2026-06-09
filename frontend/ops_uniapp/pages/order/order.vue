<template>
  <view class="container">
    <view class="header">
      <text class="title">采购订单</text>
      <view class="add-btn" @click="goToAdd">
        <text>+ 新增</text>
      </view>
    </view>
    
    <view class="stats-row">
      <view class="stat-item" @click="filterStatus('pending')">
        <text class="stat-num">{{ stats.pending || 0 }}</text>
        <text class="stat-label">待处理</text>
      </view>
      <view class="stat-item" @click="filterStatus('ordered')">
        <text class="stat-num">{{ stats.ordered || 0 }}</text>
        <text class="stat-label">已下单</text>
      </view>
      <view class="stat-item" @click="filterStatus('arrived')">
        <text class="stat-num">{{ stats.arrived || 0 }}</text>
        <text class="stat-label">已到达</text>
      </view>
      <view class="stat-item" @click="filterStatus('received')">
        <text class="stat-num">{{ stats.received || 0 }}</text>
        <text class="stat-label">已收件</text>
      </view>
    </view>
    
    <view class="filter-bar">
      <view class="search-box">
        <input 
          class="search-input" 
          placeholder="搜索订单..." 
          v-model="searchText"
          @confirm="onSearch"
        />
        <text class="search-btn" @click="onSearch">搜索</text>
      </view>
      <picker mode="selector" :range="statusOptions" range-key="label" @change="onStatusChange">
        <view class="picker">
          {{ currentFilter.label || '全部状态' }}
        </view>
      </picker>
    </view>
    
    <scroll-view scroll-y class="order-list" refresher-enabled @refresherrefresh="onRefresh" :refresher-triggered="refreshing" @scrolltolower="loadMore">
      <view v-if="loading && orders.length === 0" class="loading">
        <text>加载中...</text>
      </view>
      <view v-else-if="orders.length === 0" class="empty">
        <text>暂无订单</text>
      </view>
      <view v-else>
        <view 
          v-for="item in orders" 
          :key="item.ID" 
          class="order-card"
          @click="goDetail(item.ID)"
        >
          <view class="order-header">
            <text class="order-id">#{{ item.ID }}</text>
            <text class="order-status" :class="item.OrderStatus">{{ getStatusText(item.OrderStatus) }}</text>
          </view>
          <view class="order-content">
            <text class="order-title">{{ item.Title || '无标题' }}</text>
            <text class="order-remark" v-if="item.Remark">{{ item.Remark }}</text>
          </view>
          <view class="order-footer">
            <text class="order-date">{{ formatDate(item.CreatedAt) }}</text>
            <text class="order-link" v-if="item.Link">有链接</text>
          </view>
        </view>
        
        <view v-if="loadingMore" class="loading-more">
          <text>加载更多...</text>
        </view>
        <view v-else-if="hasMore" class="load-more-btn" @click="loadMore">
          <text>加载更多</text>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { purchaseApi } from '@/api/purchase.js'

const loading = ref(false)
const loadingMore = ref(false)
const orders = ref([])
const stats = ref({})
const page = ref(1)
const pageSize = ref(10)
const hasMore = ref(false)
const currentFilter = ref({})
const searchText = ref('')
const refreshing = ref(false)
const statusOptions = [
  { value: '', label: '全部状态' },
  { value: 'pending', label: '待处理' },
  { value: 'ordered', label: '已下单' },
  { value: 'arrived', label: '已到达' },
  { value: 'received', label: '已收件' },
  { value: 'lost', label: '丢件' },
  { value: 'returned', label: '退件' }
]

function getStatusText(status) {
  const option = statusOptions.find(s => s.value === status)
  return option ? option.label : status
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

async function fetchStats() {
  try {
    const res = await purchaseApi.getOrderCount()
    if (res.errCode === 0 && res.data) {
      stats.value = res.data
    }
  } catch (e) {
    console.error('获取订单统计失败', e)
  }
}

async function fetchOrders(reset = false) {
  if (reset) {
    page.value = 1
    orders.value = []
  }
  
  if (loading.value || loadingMore.value) return
  if (!reset && !hasMore.value) return
  
  if (reset) {
    loading.value = true
  } else {
    loadingMore.value = true
  }
  
  try {
    const params = {
      entries: pageSize.value,
      page: page.value
    }
    if (currentFilter.value.value) {
      params.status = currentFilter.value.value
    }
    if (searchText.value) {
      params.search = searchText.value
    }
    
    const res = await purchaseApi.getOrders(params)
    if (res.errCode === 0 && res.data) {
      const list = res.data.all_orders || []
      if (reset) {
        orders.value = list
      } else {
        orders.value = [...orders.value, ...list]
      }
      hasMore.value = list.length >= pageSize.value
      page.value++
    }
  } catch (e) {
    console.error('获取订单列表失败', e)
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

function onStatusChange(e) {
  currentFilter.value = statusOptions[e.detail.value]
  fetchOrders(true)
}

function filterStatus(status) {
  currentFilter.value = statusOptions.find(s => s.value === status) || statusOptions[0]
  fetchOrders(true)
}

function onSearch() {
  fetchOrders(true)
}

function loadMore() {
  if (hasMore.value && !loadingMore.value) {
    fetchOrders(false)
  }
}

function goDetail(id) {
  uni.navigateTo({
    url: `/pages/order/order-detail?id=${id}`
  })
}

onMounted(() => {
  fetchStats()
  fetchOrders(true)
})

// 监听状态变更后刷新
uni.$on('purchase-refresh', () => {
  fetchStats()
  fetchOrders(true)
})

function goToAdd() {
  uni.navigateTo({
    url: '/pages/order/order-add'
  })
}

async function onRefresh() {
  refreshing.value = true
  await Promise.all([fetchStats(), fetchOrders(true)])
  refreshing.value = false
}
</script>

<style scoped>
.container {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 20rpx;
}

.header {
  background-color: #fff;
  padding: 30rpx;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 36rpx;
  font-weight: bold;
  color: #333;
}

.add-btn {
  padding: 12rpx 24rpx;
  background-color: #1890ff;
  color: #fff;
  font-size: 26rpx;
  border-radius: 8rpx;
}

.stats-row {
  display: flex;
  background-color: #fff;
  padding: 20rpx 10rpx;
  margin-bottom: 20rpx;
}

.stat-item {
  flex: 1;
  text-align: center;
  padding: 15rpx 0;
}

.stat-num {
  display: block;
  font-size: 36rpx;
  font-weight: bold;
  color: #007AFF;
}

.stat-label {
  display: block;
  font-size: 22rpx;
  color: #666;
  margin-top: 8rpx;
}

.filter-bar {
  background-color: #fff;
  padding: 20rpx 30rpx;
  margin-bottom: 20rpx;
}

.search-box {
  display: flex;
  margin-bottom: 20rpx;
}

.search-input {
  flex: 1;
  background-color: #f0f0f0;
  padding: 15rpx 25rpx;
  border-radius: 10rpx 0 0 10rpx;
  font-size: 28rpx;
}

.search-btn {
  background-color: #007AFF;
  color: #fff;
  padding: 15rpx 35rpx;
  border-radius: 0 10rpx 10rpx 0;
  font-size: 28rpx;
}

.picker {
  background-color: #f0f0f0;
  padding: 20rpx 30rpx;
  border-radius: 10rpx;
  font-size: 28rpx;
}

.order-list {
  height: calc(100vh - 480rpx);
}

.loading, .empty {
  text-align: center;
  padding: 100rpx 0;
  color: #999;
}

.order-card {
  background-color: #fff;
  margin: 0 20rpx 20rpx;
  padding: 30rpx;
  border-radius: 12rpx;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.order-id {
  font-size: 28rpx;
  font-weight: bold;
  color: #333;
}

.order-status {
  font-size: 22rpx;
  padding: 6rpx 16rpx;
  border-radius: 20rpx;
  background-color: #f0f0f0;
  color: #666;
}

.order-status.pending {
  background-color: #fff3e0;
  color: #ff9800;
}

.order-status.ordered {
  background-color: #e3f2fd;
  color: #2196f3;
}

.order-status.arrived {
  background-color: #e8f5e9;
  color: #4caf50;
}

.order-status.received {
  background-color: #e0f2f1;
  color: #009688;
}

.order-status.lost {
  background-color: #ffebee;
  color: #f44336;
}

.order-status.returned {
  background-color: #f3e5f5;
  color: #9c27b0;
}

.order-content {
  margin-bottom: 20rpx;
}

.order-title {
  display: block;
  font-size: 30rpx;
  color: #333;
  margin-bottom: 10rpx;
}

.order-remark {
  display: block;
  font-size: 26rpx;
  color: #666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.order-footer {
  display: flex;
  justify-content: space-between;
  border-top: 1rpx solid #f0f0f0;
  padding-top: 15rpx;
}

.order-date {
  font-size: 24rpx;
  color: #999;
}

.order-link {
  font-size: 24rpx;
  color: #007AFF;
}

.loading-more, .load-more-btn {
  text-align: center;
  padding: 30rpx;
  color: #007AFF;
}

.load-more-btn {
  background-color: #fff;
  margin: 0 20rpx;
  border-radius: 10rpx;
}
</style>
