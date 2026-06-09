<template>
  <view class="container">
    <!-- 搜索栏 -->
    <view class="search-bar">
      <text class="search-back" @click="goBack">‹</text>
      <input
        class="search-input"
        v-model="keyword"
        placeholder="输入关键词搜索..."
        confirm-type="search"
        @input="onInput"
        @confirm="doSearch"
      />
      <text class="search-scan" @click="onScan">📷</text>
      <text class="search-clear" v-if="keyword" @click="clearKeyword">×</text>
    </view>

    <!-- 分类Tab -->
    <view class="tabs">
      <view
        class="tab"
        v-for="cat in categories"
        :key="cat.key"
        :class="{ active: currentCat === cat.key }"
        @click="switchCat(cat.key)"
      >
        {{ cat.label }}
      </view>
    </view>

    <!-- 加载状态 -->
    <view v-if="loading" class="loading">
      <text>搜索中...</text>
    </view>

    <!-- 空状态 -->
    <view v-else-if="results.length === 0 && searched" class="empty">
      <text>未找到相关内容</text>
    </view>

    <!-- 搜索结果列表 -->
    <scroll-view v-else scroll-y class="result-list">
      <view
        v-for="item in results"
        :key="item._type + ':' + item.ID"
        class="result-card"
        @click="goDetail(item)"
      >
        <view class="result-icon">
          <text>{{ getIcon(item._type) }}</text>
        </view>
        <view class="result-info">
          <text class="result-title">{{ getTitle(item) }}</text>
          <text class="result-desc">{{ getDesc(item) }}</text>
        </view>
        <view class="result-arrow">›</view>
      </view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { purchaseApi } from '@/api/purchase.js'
import { warehouseApi } from '@/api/warehouse.js'
import { workOrderApi } from '@/api/work_order.js'

const keyword = ref('')
const currentCat = ref('all')
const loading = ref(false)
const results = ref([])
const searched = ref(false)

let searchTimer = null

const categories = [
  { key: 'all', label: '全部' },
  { key: 'order', label: '订单' },
  { key: 'workorder', label: '工单' },
  { key: 'item', label: '物品' },
  { key: 'container', label: '容器' },
]

function onInput() {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    if (keyword.value.trim()) {
      doSearch()
    } else {
      results.value = []
      searched.value = false
    }
  }, 500)
}

function doSearch() {
  if (!keyword.value.trim()) {
    results.value = []
    searched.value = false
    return
  }

  // 条码识别：wo:ID / item:ID / warehouse:ID
  const val = keyword.value.trim().toLowerCase()
  let match

  // wo:ID → 工单
  match = val.match(/^wo:(\d+)$/)
  if (match) {
    uni.navigateTo({ url: '/pages/workorder/show-workorder?id=' + match[1] })
    return
  }

  // item:ID → 物品
  match = val.match(/^item:(\d+)$/)
  if (match) {
    uni.navigateTo({ url: '/pages/warehouse/item-detail?id=' + match[1] })
    return
  }

  // warehouse:ID → 容器（tabBar页面，用事件传参）
  match = val.match(/^warehouse:(\d+)$/)
  if (match) {
    uni.$emit('barcode-navigate-container', { id: match[1] })
    uni.switchTab({ url: '/pages/warehouse/warehouse' })
    return
  }

  // po:ID → 采购订单
  match = val.match(/^po:(\d+)$/)
  if (match) {
    uni.navigateTo({ url: '/pages/order/order-detail?id=' + match[1] })
    return
  }

  searched.value = true
  loading.value = true

  const cat = currentCat.value
  const tasks = []

  if (cat === 'all' || cat === 'order') {
    tasks.push(
      purchaseApi.getOrders({ search: keyword.value, entries: 50, page: 1 }).then(res => {
        if (res.errCode === 0 && res.data && res.data.all_orders) {
          return (res.data.all_orders || []).map(o => ({ ...o, _type: 'order' }))
        }
        return []
      }).catch(() => [])
    )
  }

  if (cat === 'all' || cat === 'workorder') {
    tasks.push(
      workOrderApi.list({ search: keyword.value, entries: 50, page: 1 }).then(res => {
        if (res.errCode === 0 && res.data && res.data.all_orders) {
          return (res.data.all_orders || []).map(o => ({ ...o, _type: 'workorder' }))
        }
        return []
      }).catch(() => [])
    )
  }

  if (cat === 'all' || cat === 'item') {
    tasks.push(
      warehouseApi.listItem({ search: keyword.value, entries: 50, page: 1 }).then(res => {
        if (res.errCode === 0 && res.data && res.data.items) {
          return (res.data.items || []).map(i => ({ ...i, _type: 'item' }))
        }
        return []
      }).catch(() => [])
    )
  }

  if (cat === 'all' || cat === 'container') {
    tasks.push(
      warehouseApi.listContainer({ search: keyword.value, entries: 50, page: 1 }).then(res => {
        if (res.errCode === 0 && res.data && res.data.containers) {
          return (res.data.containers || []).map(c => ({ ...c, _type: 'container' }))
        }
        return []
      }).catch(() => [])
    )
  }

  Promise.all(tasks).then(allResults => {
    const merged = allResults.flat()
    results.value = merged
  }).finally(() => {
    loading.value = false
  })
}

function switchCat(cat) {
  currentCat.value = cat
  if (keyword.value.trim()) {
    doSearch()
  }
}

function clearKeyword() {
  keyword.value = ''
  results.value = []
  searched.value = false
}

function getIcon(type) {
  const icons = {
    order: '📦',
    workorder: '🔧',
    item: '📦',
    container: '📁',
  }
  return icons[type] || '📄'
}

function getTitle(item) {
  switch (item._type) {
    case 'order':
      return item.Title || '订单'
    case 'workorder':
      return item.Title || '工单'
    case 'item':
      return item.Name || '物品'
    case 'container':
      return item.Title || '容器'
    default:
      return ''
  }
}

function getDesc(item) {
  switch (item._type) {
    case 'order':
      return '状态: ' + (item.Status || '') + '  |  ' + (item.Vendor || '')
    case 'workorder':
      return '状态: ' + (item.Status || '')
    case 'item':
      return item.SerialNumber || (item.Remark ? item.Remark.substring(0, 30) : '')
    case 'container':
      return '子容器: ' + (item.ChildCount || 0) + ' | 物品: ' + (item.ItemCount || 0)
    default:
      return ''
  }
}

function goDetail(item) {
  switch (item._type) {
    case 'order':
      uni.navigateTo({ url: '/pages/order/order-detail?id=' + item.ID })
      break
    case 'workorder':
      uni.navigateTo({ url: '/pages/workorder/show-workorder?id=' + item.ID })
      break
    case 'item':
      uni.navigateTo({ url: '/pages/warehouse/item-detail?id=' + item.ID })
      break
    case 'container':
      uni.switchTab({ url: '/pages/warehouse/warehouse?container_id=' + item.ID })
      break
  }
}

function goBack() {
  uni.navigateBack({ delta: 1 })
}

// 扫码
function onScan() {
  // #ifdef H5
  uni.showModal({
    title: '提示',
    content: 'H5 端不支持摄像头扫码，请手动输入条码（如 wo:1、item:1、warehouse:1）',
    showCancel: false
  })
  // #endif

  // #ifdef APP-PLUS
  // App 端使用 5+ API 检查并申请相机权限
  const osName = plus.os.name
  if (osName === 'Android') {
    const main = plus.android.runtimeMainActivity()
    const ContextCompat = plus.android.importClass('androidx.core.content.ContextCompat')
    const Manifest = plus.android.importClass('android.Manifest$permission')
    const permission = Manifest.CAMERA
    const result = ContextCompat.checkSelfPermission(main, permission)
    const PackageManager = plus.android.importClass('android.content.pm.PackageManager')
    if (result !== PackageManager.PERMISSION_GRANTED) {
      // 没有权限，申请权限
      plus.android.requestPermissions(['android.permission.CAMERA'], (e) => {
        if (e.deniedAlways.length > 0) {
          // 永久拒绝，引导去设置
          uni.showModal({
            title: '需要相机权限',
            content: '请在系统设置中开启相机权限后再使用扫码功能',
            confirmText: '去设置',
            success: (res) => {
              if (res.confirm) {
                if (plus.os.name === 'Android') {
                  const Intent = plus.android.importClass('android.content.Intent')
                  const Settings = plus.android.importClass('android.provider.Settings')
                  const intent = new Intent(Settings.ACTION_APPLICATION_DETAILS_SETTINGS)
                  const Uri = plus.android.importClass('android.net.Uri')
                  intent.setData(Uri.fromParts('package', main.getPackageName(), null))
                  main.startActivity(intent)
                }
              }
            }
          })
        } else if (e.granted.length > 0) {
          // 授权成功，调用扫码
          callScan()
        }
      }, (e) => {
        console.error('申请权限失败', e)
        callScan()
      })
      return
    }
  }
  // iOS 或已有权限，直接调用扫码
  callScan()
  // #endif

  // #ifdef MP-WEIXIN
  uni.getSetting({
    success: (res) => {
      if (res.authSetting['scope.camera'] === false) {
        uni.showModal({
          title: '需要相机权限',
          content: '扫码功能需要相机权限，请在设置中开启',
          success: (modalRes) => {
            if (modalRes.confirm) {
              uni.openSetting()
            }
          }
        })
      } else {
        callScan()
      }
    },
    fail: () => {
      callScan()
    }
  })
  // #endif
}

// 调用扫码
function callScan() {
  uni.scanCode({
    onlyFromCamera: false,
    scanType: ['barCode', 'qrCode'],
    success: (res) => {
      console.log('扫码结果:', res)
      keyword.value = res.result
      doSearch()
    },
    fail: (err) => {
      console.error('扫码失败', err)
      if (err.errMsg && err.errMsg.includes('cancel')) return
      uni.showToast({ title: '扫码失败: ' + (err.errMsg || '未知错误'), icon: 'none' })
    }
  })
}
</script>

<style scoped>
.container {
  min-height: 100vh;
  background-color: #f5f5f5;
  display: flex;
  flex-direction: column;
}

/* 搜索栏 */
.search-bar {
  display: flex;
  align-items: center;
  padding: 20rpx 24rpx;
  background-color: #007AFF;
  gap: 16rpx;
}

.search-back {
  font-size: 48rpx;
  color: #fff;
  width: 60rpx;
  text-align: center;
}

.search-input {
  flex: 1;
  height: 72rpx;
  background-color: #fff;
  border-radius: 36rpx;
  padding: 0 30rpx;
  font-size: 28rpx;
}

.search-clear {
  font-size: 40rpx;
  color: #fff;
  width: 60rpx;
  text-align: center;
}

.search-scan {
  font-size: 40rpx;
  color: #fff;
  width: 60rpx;
  text-align: center;
}

/* 分类Tab */
.tabs {
  display: flex;
  background-color: #fff;
  padding: 0 10rpx;
  border-bottom: 1rpx solid #eee;
}

.tab {
  flex: 1;
  text-align: center;
  padding: 24rpx 0;
  font-size: 26rpx;
  color: #666;
  border-bottom: 4rpx solid transparent;
}

.tab.active {
  color: #007AFF;
  border-bottom-color: #007AFF;
  font-weight: bold;
}

/* 加载/空状态 */
.loading,
.empty {
  text-align: center;
  padding: 120rpx 0;
  color: #999;
  font-size: 28rpx;
}

/* 结果列表 */
.result-list {
  flex: 1;
  height: 0;
  padding: 20rpx;
}

.result-card {
  display: flex;
  align-items: center;
  background-color: #fff;
  padding: 28rpx 24rpx;
  border-radius: 12rpx;
  margin-bottom: 16rpx;
}

.result-icon {
  width: 80rpx;
  height: 80rpx;
  background-color: #f0f0f0;
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 36rpx;
  margin-right: 20rpx;
  flex-shrink: 0;
}

.result-info {
  flex: 1;
  overflow: hidden;
}

.result-title {
  display: block;
  font-size: 30rpx;
  color: #333;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.result-desc {
  display: block;
  font-size: 24rpx;
  color: #999;
  margin-top: 8rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.result-arrow {
  font-size: 36rpx;
  color: #ccc;
  margin-left: 16rpx;
}
</style>
