<template>
  <view class="container">
    <view class="header">
      <text class="title">工单管理</text>
      <text class="add-btn" @click="goAddWorkOrder">+ 新建</text>
    </view>
    
    <view class="stats-row">
      <view class="stat-item" @click="filterStatus('pending')">
        <text class="stat-num">{{ stats.pending || 0 }}</text>
        <text class="stat-label">待处理</text>
      </view>
      <view class="stat-item" @click="filterStatus('checked')">
        <text class="stat-num">{{ stats.checked || 0 }}</text>
        <text class="stat-label">已检查</text>
      </view>
      <view class="stat-item" @click="filterStatus('parts_ordered')">
        <text class="stat-num">{{ stats.parts_ordered || 0 }}</text>
        <text class="stat-label">已下单零件</text>
      </view>
      <view class="stat-item" @click="filterStatus('repaired')">
        <text class="stat-num">{{ stats.repaired || 0 }}</text>
        <text class="stat-label">已维修</text>
      </view>
    </view>
    
    <view class="filter-bar">
      <picker mode="selector" :range="statusOptions" range-key="label" @change="onStatusChange">
        <view class="picker">
          {{ currentFilter.label || '全部状态' }}
        </view>
      </picker>
    </view>
    
    <scroll-view scroll-y class="order-list" refresher-enabled @refresherrefresh="onRefresh" :refresher-triggered="refreshing" @scrolltolower="loadMore">
      <view v-if="loading && workOrders.length === 0" class="loading">
        <text>加载中...</text>
      </view>
      <view v-else-if="workOrders.length === 0" class="empty">
        <text>暂无工单</text>
      </view>
      <view v-else>
        <view 
          v-for="item in workOrders" 
          :key="item.ID" 
          class="order-card"
          @click="goDetail(item.ID)"
        >
          <view class="order-header">
            <text class="order-id">#{{ item.ID }}</text>
            <text class="order-status" :class="item.CurrentStatus">{{ getStatusText(item.CurrentStatus) }}</text>
          </view>
          <view class="order-content">
            <text class="order-title">{{ item.Title || '无标题' }}</text>
            <text class="order-desc">{{ item.Description || '无描述' }}</text>
          </view>
          <view class="order-footer">
            <text class="order-date">创建: {{ formatDate(item.CreatedAt) }}</text>
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
import { workOrderApi } from '@/api/work_order.js'

const loading = ref(false)
const loadingMore = ref(false)
const workOrders = ref([])
const stats = ref({})
const page = ref(1)
const pageSize = ref(10)
const hasMore = ref(false)
const currentFilter = ref({})
const refreshing = ref(false)
const statusOptions = [
  { value: '', label: '全部状态' },
  { value: 'pending', label: '待处理' },
  { value: 'checked', label: '已检查' },
  { value: 'parts_ordered', label: '已下单零件' },
  { value: 'repaired', label: '已维修' },
  { value: 'returned', label: '已送还' },
  { value: 'unrepairable', label: '无法维修' }
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
    const res = await workOrderApi.getCount()
    if (res.errCode === 0 && res.data) {
      stats.value = res.data
    }
  } catch (e) {
    console.error('获取工单统计失败', e)
  }
}

async function fetchWorkOrders(reset = false) {
  if (reset) {
    page.value = 1
    workOrders.value = []
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
      page: page.value,
      entries: pageSize.value
    }
    if (currentFilter.value.value) {
      params.status = currentFilter.value.value
    }
    
    const res = await workOrderApi.list(params)
    if (res.errCode === 0 && res.data) {
      const list = res.data.all_orders || []
      if (reset) {
        workOrders.value = list
      } else {
        workOrders.value = [...workOrders.value, ...list]
      }
      hasMore.value = list.length >= pageSize.value
      page.value++
    }
  } catch (e) {
    console.error('获取工单列表失败', e)
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

function onStatusChange(e) {
  currentFilter.value = statusOptions[e.detail.value]
  fetchWorkOrders(true)
}

function filterStatus(status) {
  currentFilter.value = statusOptions.find(s => s.value === status) || statusOptions[0]
  fetchWorkOrders(true)
}

function loadMore() {
  if (hasMore.value && !loadingMore.value) {
    fetchWorkOrders(false)
  }
}

function goDetail(id) {
  uni.navigateTo({
    url: `/pages/workorder/show-workorder?id=${id}`
  })
}

function goAddWorkOrder() {
  uni.navigateTo({
    url: '/pages/workorder/add-workorder'
  })
}

async function onRefresh() {
  refreshing.value = true
  await Promise.all([fetchStats(), fetchWorkOrders(true)])
  refreshing.value = false
}

onMounted(() => {
  fetchStats()
  fetchWorkOrders(true)
})

// 监听新增/编辑工单后刷新
uni.$on('workorder-refresh', () => {
  fetchStats()
  fetchWorkOrders(true)
})
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
  text-align: center;
}

.title {
  font-size: 36rpx;
  font-weight: bold;
  color: #333;
}

.add-btn {
  font-size: 28rpx;
  color: #007AFF;
  position: absolute;
  right: 30rpx;
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

.picker {
  background-color: #f0f0f0;
  padding: 20rpx 30rpx;
  border-radius: 10rpx;
  font-size: 28rpx;
}

.order-list {
  height: calc(100vh - 400rpx);
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
  background-color: #fef3c7;
  color: #b45309;
}

.order-status.checked {
  background-color: #f3e8ff;
  color: #7c3aed;
}

.order-status.parts_ordered {
  background-color: #dbeafe;
  color: #1d4ed8;
}

.order-status.repaired {
  background-color: #dcfce7;
  color: #15803d;
}

.order-status.returned {
  background-color: #e5e7eb;
  color: #4b5563;
}

.order-status.unrepairable {
  background-color: #fee2e2;
  color: #dc2626;
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

.order-desc {
  display: block;
  font-size: 26rpx;
  color: #666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.order-footer {
  border-top: 1rpx solid #f0f0f0;
  padding-top: 15rpx;
}

.order-date {
  font-size: 24rpx;
  color: #999;
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
