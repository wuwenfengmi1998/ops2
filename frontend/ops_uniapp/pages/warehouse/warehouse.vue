<template>
  <view class="container">
    <view class="header">
      <text class="title">仓库管理</text>
    </view>
    
    <!-- Tab切换 -->
    <view class="tabs">
      <view 
        class="tab" 
        :class="{ active: currentTab === 'containers' }"
        @click="switchTab('containers')"
      >
        全部
      </view>
      <view 
        class="tab" 
        :class="{ active: currentTab === 'items' }"
        @click="switchTab('items')"
      >
        仅物品
      </view>
      <view class="tab-search" @click="showSearchModal = true">
        <text class="search-icon">🔍</text>
      </view>
    </view>
    
    <!-- 列表（容器+物品） -->
    <scroll-view scroll-y class="list-container" refresher-enabled @refresherrefresh="onRefresh" :refresher-triggered="refreshing" v-if="currentTab === 'containers'">
      <!-- 面包屑导航和返回按钮 -->
      <view class="nav-bar" v-if="currentContainerId > 0">
        <text class="back-btn" @click="goBack">‹ 返回上级</text>
        <view class="breadcrumb" v-if="breadcrumbs.length > 0">
          <text 
            class="breadcrumb-item" 
            v-for="(crumb, index) in breadcrumbs" 
            :key="crumb.id"
          >
            <text @click="goToContainer(crumb.id)">{{ crumb.title }}</text>
            <text class="breadcrumb-sep" v-if="index < breadcrumbs.length - 1"> / </text>
          </text>
        </view>
        <text class="print-btn" @click="printContainer">打印</text>
        <text class="edit-btn" @click="editCurrentContainer">✎ 编辑</text>
      </view>
      
      <view v-if="loading" class="loading">
        <text>加载中...</text>
      </view>
      <view v-else-if="containers.length === 0 && items.length === 0" class="empty">
        <text>暂无内容</text>
      </view>
      <view v-else>
        <!-- 容器卡片 -->
        <view
          v-for="item in containers"
          :key="'c_' + item.ID"
          class="item-card container-card"
          @click="onContainerClick(item)"
        >
          <view class="item-icon">
            <text class="iconfont">📁</text>
          </view>
          <view class="item-info">
            <text class="item-name">{{ item.Title }}</text>
            <text class="item-meta">
              子容器: {{ item.ChildCount }} | 物品: {{ item.ItemCount }}
            </text>
          </view>
          <view class="item-arrow">›</view>
        </view>
        
        <!-- 物品卡片 -->
        <view 
          v-for="item in items" 
          :key="'i_' + item.ID" 
          class="item-card"
          @click="goItemDetail(item.ID)"
        >
          <view class="item-icon">
            <text class="iconfont">📦</text>
          </view>
          <view class="item-info">
            <text class="item-name">{{ item.Name }}</text>
            <text class="item-meta" v-if="item.SerialNumber">编号: {{ item.SerialNumber }}</text>
            <text class="item-meta">数量: {{ item.Quantity ?? 1 }}</text>
          </view>
          <view class="item-arrow">›</view>
        </view>
      </view>
    </scroll-view>
    
    <!-- 物品列表 -->
    <scroll-view scroll-y class="list-container" refresher-enabled @refresherrefresh="onRefresh" :refresher-triggered="refreshing" v-else>
      <view v-if="loading" class="loading">
        <text>加载中...</text>
      </view>
      <view v-else-if="items.length === 0" class="empty">
        <text>暂无物品</text>
      </view>
      <view v-else>
        <!-- 物品卡片 -->
        <view 
          v-for="item in items" 
          :key="item.ID" 
          class="item-card"
          @click="goItemDetail(item.ID)"
        >
          <view class="item-icon">
            <text class="iconfont">📦</text>
          </view>
          <view class="item-info">
            <text class="item-name">{{ item.Name }}</text>
            <text class="item-meta" v-if="item.SerialNumber">编号: {{ item.SerialNumber }}</text>
            <text class="item-meta">数量: {{ item.Quantity ?? 1 }}</text>
            <text class="item-location" v-if="item.ContainerBreadcrumb">
              位置: {{ item.ContainerBreadcrumb }}
            </text>
          </view>
          <view class="item-arrow">›</view>
        </view>
        
        <view v-if="hasMore" class="load-more-btn" @click="loadMore">
          <text>加载更多</text>
        </view>
      </view>
    </scroll-view>
    
    <!-- 底部操作栏 -->
    <view :class="['action-bar', (currentTab === 'containers' && currentContainerId > 0) ? 'two-btn' : 'single-btn']">
      <!-- 容器内部：显示新增容器、新增物品两个按钮 -->
      <template v-if="currentTab === 'containers' && currentContainerId > 0">
        <view class="action-btn half" @click="showAddModal">
          <text class="action-icon">+</text>
          <text class="action-text">新增容器</text>
        </view>
        <view class="action-btn half" @click="goAddItem">
          <text class="action-icon">+</text>
          <text class="action-text">新增物品</text>
        </view>
      </template>
      <!-- 根目录：显示新增容器 -->
      <view class="action-btn" v-else-if="currentTab === 'containers'" @click="showAddModal">
        <text class="action-icon">+</text>
        <text class="action-text">新增容器</text>
      </view>
      <!-- 物品模式：显示新增物品 -->
      <view class="action-btn" v-else @click="showAddModal">
        <text class="action-icon">+</text>
        <text class="action-text">新增物品</text>
      </view>
    </view>
    
    <!-- 搜索弹窗 -->
    <view v-if="showSearchModal" class="modal-overlay" @click="showSearchModal = false">
      <view class="modal-container search-modal" @click.stop>
        <view class="modal-header">
          <text class="modal-title">搜索</text>
          <text class="modal-close" @click="showSearchModal = false">×</text>
        </view>
        <view class="modal-content">
          <input 
            class="form-input" 
            :placeholder="currentTab === 'containers' ? '搜索容器...' : '搜索物品...'" 
            v-model="searchText"
            @confirm="onSearch"
            focus
          />
          <view class="btn-confirm mt-20" @click="onSearch">搜索</view>
        </view>
      </view>
    </view>
    
    <!-- 新增/编辑容器弹窗 -->
    <view v-if="showModal" class="modal-overlay">
      <view class="modal-container">
        <view class="modal-header">
          <text class="modal-title">{{ modalMode === 'add' ? '新增容器' : '编辑容器' }}</text>
          <text class="modal-close" @click="closeModal">×</text>
        </view>
        
        <view class="modal-content">
          <view class="form-row">
            <text class="form-label">名称 *</text>
            <input class="form-input" v-model="containerForm.title" placeholder="请输入容器名称" />
          </view>
          
          <view class="form-row">
            <text class="form-label">所属目录</text>
            <view class="picker-box" @click="openParentPicker">
              <text class="picker-value">{{ selectedParentName }}</text>
              <text class="picker-arrow">›</text>
            </view>
          </view>
          
          <view class="form-row">
            <text class="form-label">备注</text>
            <textarea class="form-textarea" v-model="containerForm.remark" placeholder="请输入备注" />
          </view>
        </view>
        
        <view class="modal-footer">
          <view class="btn-cancel" @click="closeModal">取消</view>
          <view class="btn-confirm" @click="submitContainer">确定</view>
        </view>
      </view>
    </view>
    
    <!-- 父容器选择弹窗 -->
    <view v-if="showParentPicker" class="modal-overlay">
      <view class="modal-container">
        <view class="modal-header">
          <text class="modal-title">选择所属目录</text>
          <text class="modal-close" @click="showParentPicker = false">×</text>
        </view>
        <view class="modal-content" style="max-height: 400rpx;">
          <view
            v-for="container in parentContainerList"
            :key="container.ID"
            class="picker-item"
            :class="{ active: (container.ID === 0 && containerForm.parent_id === null) || container.ID === containerForm.parent_id }"
            @click="selectParent(container)"
          >
            <text class="picker-item-icon">📁</text>
            <text class="picker-item-text">{{ container.Title }}</text>
            <text v-if="(container.ID === 0 && containerForm.parent_id === null) || container.ID === containerForm.parent_id" class="picker-item-check">✓</text>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
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

function getPhotoUrl(sha256) {
  const baseUrl = configStore.getFileBaseUrl ? configStore.getFileBaseUrl() : ''
  return `${baseUrl}/api/files/get/image/${sha256}`
}

const currentTab = ref('containers')
const loading = ref(false)
const searchText = ref('')
const containers = ref([])
const items = ref([])
const currentContainerId = ref(0)
const breadcrumbs = ref([])
const page = ref(1)
const pageSize = ref(20)
const hasMore = ref(false)

// 弹窗相关
const showModal = ref(false)
const modalMode = ref('add')
const containerForm = ref({
  id: 0,
  title: '',
  remark: '',
  parent_id: null
})





// 父容器选择器
const showParentPicker = ref(false)
const parentContainerList = ref([])
const selectedParentName = ref('当前目录')

// 所有容器列表（用于选择父容器）
const allContainers = ref([])

// 当前容器对象（用于编辑当前容器）
const currentContainer = ref(null)

// 搜索弹窗
const showSearchModal = ref(false)
const refreshing = ref(false)

function switchTab(tab) {
  if (currentTab.value === tab) return
  currentTab.value = tab
  searchText.value = ''
  page.value = 1
  if (tab === 'containers') {
    fetchContainers()
  } else {
    fetchItems(true)
  }
}

function onSearch() {
  page.value = 1
  if (currentTab.value === 'containers') {
    fetchContainers()
  } else {
    fetchItems(true)
  }
}

async function onRefresh() {
  refreshing.value = true
  if (currentTab.value === 'containers') {
    await fetchContainers()
  } else {
    await fetchItems(true)
  }
  refreshing.value = false
}

async function fetchContainers() {
  loading.value = true
  try {
    // 同时获取容器和物品
    const [containerRes, itemRes] = await Promise.all([
      warehouseApi.listContainer({
        parent_id: currentContainerId.value || null,
        search: searchText.value || '',
        all_levels: false,
        entries: pageSize.value,
        page: 1
      }),
      warehouseApi.listItem({
        search: searchText.value || '',
        container_id: currentContainerId.value || null,
        entries: pageSize.value,
        page: 1
      })
    ])
    
    if (containerRes.errCode === 0 && containerRes.data) {
      containers.value = containerRes.data.containers || []
    }
    if (itemRes.errCode === 0 && itemRes.data) {
      items.value = itemRes.data.items || []
      hasMore.value = (itemRes.data.items || []).length >= pageSize.value
    }
  } catch (e) {
    console.error('获取数据失败', e)
  } finally {
    loading.value = false
  }
}

async function fetchItems(reset = false) {
  if (reset) {
    page.value = 1
    items.value = []
  }
  loading.value = true
  
  try {
    const params = {
      search: searchText.value || '',
      container_id: currentContainerId.value || null,
      entries: pageSize.value,
      page: page.value
    }
    const res = await warehouseApi.listItem(params)
    if (res.errCode === 0 && res.data) {
      const list = res.data.items || []
      if (reset) {
        items.value = list
      } else {
        items.value = [...items.value, ...list]
      }
      hasMore.value = list.length >= pageSize.value
      page.value++
    }
  } catch (e) {
    console.error('获取物品列表失败', e)
  } finally {
    loading.value = false
  }
}

function loadMore() {
  if (hasMore.value) {
    fetchItems(false)
  }
}

function onContainerClick(container) {
  currentContainerId.value = container.ID
  currentContainer.value = container
  breadcrumbs.value.push({ id: container.ID, title: container.Title })
  fetchContainers()
}

async function goBack() {
  if (breadcrumbs.value.length > 1) {
    // 返回上一级（倒数第二个）
    const parentId = breadcrumbs.value[breadcrumbs.value.length - 2].id
    goToContainer(parentId)
  } else {
    // 返回根目录
    goToContainer(0)
  }
}

function goToContainer(id) {
  currentContainerId.value = id
  if (id === 0) {
    breadcrumbs.value = []
    currentContainer.value = null
  } else {
    // 找到对应索引，截断面包屑
    const index = breadcrumbs.value.findIndex(c => c.id === id)
    if (index >= 0) {
      breadcrumbs.value = breadcrumbs.value.slice(0, index + 1)
      // 更新当前容器为截断后的最后一个
      currentContainer.value = containers.value.find(c => c.ID === id) || null
    }
  }
  fetchContainers()
}

function showAddModal() {
  modalMode.value = 'add'
  containerForm.value = { 
    id: 0, 
    title: '', 
    remark: '',
    parent_id: currentContainerId.value || null
  }
  selectedParentName.value = currentContainerId.value > 0 
    ? (breadcrumbs.value[breadcrumbs.value.length - 1]?.title || '当前目录')
    : '根目录'
  showModal.value = true
}



function showEditModal(container) {
  modalMode.value = 'edit'
  containerForm.value = {
    id: container.ID,
    title: container.Title,
    remark: container.Remark || '',
    parent_id: container.ParentID
  }
  selectedParentName.value = container.ParentID > 0
    ? (breadcrumbs.value.find(c => c.id === container.ParentID)?.title || '未知')
    : '根目录'
  showModal.value = true
}

function editCurrentContainer() {
  if (currentContainer.value) {
    showEditModal(currentContainer.value)
  } else {
    uni.showToast({ title: '无法编辑', icon: 'none' })
  }
}

function goAddItem() {
  uni.navigateTo({
    url: `/pages/warehouse/add-item?container_id=${currentContainerId.value}`
  })
}

function confirmDeleteContainer(container) {
  if (container.ChildCount > 0 || container.ItemCount > 0) {
    uni.showToast({ title: '请先删除子容器和物品', icon: 'none' })
    return
  }
  uni.showModal({
    title: '确认删除',
    content: `确定要删除容器"${container.Title}"吗？`,
    success: async (res) => {
      if (res.confirm) {
        await deleteContainer(container.ID)
      }
    }
  })
}

async function deleteContainer(id) {
  try {
    const res = await warehouseApi.deleteContainer(id)
    if (res.errCode === 0) {
      uni.showToast({ title: '删除成功', icon: 'success' })
      fetchContainers()
    } else {
      uni.showToast({ title: res.errCode || '删除失败', icon: 'none' })
    }
  } catch (e) {
    console.error('删除失败', e)
    uni.showToast({ title: '删除失败', icon: 'none' })
  }
}

async function loadAllContainers() {
  try {
    const res = await warehouseApi.listContainer({
      all_levels: true,
      entries: 1000,
      page: 1
    })
    if (res.errCode === 0 && res.data) {
      allContainers.value = res.data.containers || []
    }
  } catch (e) {
    console.error('获取容器列表失败', e)
  }
}

function openParentPicker() {
  loadAllContainers()
  parentContainerList.value = [
    { ID: 0, Title: '根目录', ParentID: null }
  ].concat(allContainers.value.filter(c => c.ID !== containerForm.value.id))
  showParentPicker.value = true
}

function selectParent(container) {
  containerForm.value.parent_id = container.ID === 0 ? null : container.ID
  selectedParentName.value = container.Title
  showParentPicker.value = false
}

function closeModal() {
  showModal.value = false
}

async function submitContainer() {
  if (!containerForm.value.title.trim()) {
    uni.showToast({ title: '请输入名称', icon: 'none' })
    return
  }
  
  try {
    const data = {
      title: containerForm.value.title,
      remark: containerForm.value.remark,
      parent_id: containerForm.value.parent_id
    }
    
    let res
    if (modalMode.value === 'add') {
      res = await warehouseApi.addContainer(data)
    } else {
      data.id = containerForm.value.id
      res = await warehouseApi.updateContainer(data)
    }
    
    if (res.errCode === 0) {
      uni.showToast({ title: '操作成功', icon: 'success' })
      closeModal()
      fetchContainers()
    } else {
      uni.showToast({ title: res.errCode || '操作失败', icon: 'none' })
    }
  } catch (e) {
    console.error('提交失败', e)
    uni.showToast({ title: '操作失败', icon: 'none' })
  }
}

function goItemDetail(id) {
  uni.navigateTo({
    url: `/pages/warehouse/item-detail?id=${id}`
  })
}

// 打印容器标签
async function printContainer() {
  if (!printer) {
    uni.showToast({ title: '打印机插件未加载（仅在 App 环境可用）', icon: 'none' })
    return
  }

  if (!currentContainerId.value || currentContainerId.value === 0) {
    uni.showToast({ title: '容器信息未加载', icon: 'none' })
    return
  }

  uni.showLoading({ title: '获取容器信息...' })

  try {
    // 获取完整容器详情（包含 Remark, CreatedAt 等字段）
    const res = await warehouseApi.getContainer(currentContainerId.value)
    uni.hideLoading()

    if (res.errCode !== 0 || !res.data || !res.data.container) {
      uni.showToast({ title: '获取容器信息失败', icon: 'none' })
      return
    }

    const container = res.data.container

    // 初始化打印机
    printer.initPrinter({})
    printer.setConcentration({ level: 39 })
    printer.setLineSpacing({ spacing: 1 })

    // 标签打印，使用黑标
    printer.printEnableMark({
      enable: true
    })

    // 第一行：容器名（加粗）
    printer.setFontSize({ fontSize: 0 })
    printer.setTextBold({ bold: true })
    printer.printText({
      content: (container.Title || '容器') + '\n'
    })
    printer.setTextBold({ bold: false })

    // 第二行：完整面包屑
    const breadcrumbText = breadcrumbs.value.map(b => b.title).join(' / ')
    if (breadcrumbText) {
      printer.setFontSize({ fontSize: 1 })
      printer.printText({
        content: breadcrumbText + '\n'
      })
    }

    // 第三行：备注
    if (container.Remark) {
      printer.setFontSize({ fontSize: 1 })
      printer.printText({
        content: container.Remark + '\n'
      })
    }

    // 第四行：创建日期
    if (container.CreatedAt) {
      printer.setFontSize({ fontSize: 1 })
      printer.printText({
        content: formatDate(container.CreatedAt)
      })
    }

    // 条形码（高度40）
    printer.printBarcode({
      text: 'warehouse:' + container.ID,
      height: 40,
      barcodeType: 73
    })
    printer.printLine({ line_length: 2 })

    // 提交打印
    printer.start()

    uni.showToast({ title: '打印成功', icon: 'success' })
  } catch (e) {
    uni.hideLoading()
    console.error('打印失败', e)
    uni.showToast({ title: '打印失败', icon: 'none' })
  }
}

// 格式化日期
function formatDate(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

onMounted(() => {
  fetchContainers()
})

// 每次页面显示时刷新数据（如从新增物品页面返回时）
onShow(() => {
  // 检查是否有条码导航事件
  uni.$once('barcode-navigate-container', (data) => {
    if (data && data.id) {
      currentContainerId.value = parseInt(data.id)
      // 清空搜索，确保显示容器内容
      searchText.value = ''
    }
  })
  fetchContainers()
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
  text-align: center;
}

.title {
  font-size: 36rpx;
  font-weight: bold;
  color: #333;
}

.tabs {
  display: flex;
  background-color: #fff;
  padding: 0 20rpx;
  margin-bottom: 20rpx;
}

.tab {
  flex: 1;
  text-align: center;
  padding: 25rpx 0;
  font-size: 28rpx;
  color: #666;
  border-bottom: 4rpx solid transparent;
}

.tab.active {
  color: #007AFF;
  border-bottom-color: #007AFF;
}

.filter-bar {
  background-color: #fff;
  padding: 20rpx 30rpx;
  margin-bottom: 20rpx;
}

.search-box {
  display: flex;
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

.list-container {
  height: calc(100vh - 400rpx);
}

.loading, .empty {
  text-align: center;
  padding: 100rpx 0;
  color: #999;
}

.nav-bar {
  background-color: #fff;
  padding: 20rpx 30rpx;
  margin-bottom: 20rpx;
  display: flex;
  align-items: center;
}

.back-btn {
  font-size: 28rpx;
  color: #007AFF;
  margin-right: 20rpx;
  white-space: nowrap;
}

.breadcrumb {
  flex: 1;
  font-size: 26rpx;
  color: #666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.breadcrumb-item {
  color: #007AFF;
}

.breadcrumb-item.active {
  color: #333;
}

.breadcrumb-sep {
  margin: 0 10rpx;
  color: #ccc;
}

.edit-btn {
  font-size: 28rpx;
  color: #007AFF;
  white-space: nowrap;
  margin-left: 20rpx;
}

.print-btn {
  font-size: 28rpx;
  color: #007AFF;
  white-space: nowrap;
  margin-left: 20rpx;
}

.item-card {
  display: flex;
  align-items: center;
  background-color: #fff;
  margin: 0 20rpx 20rpx;
  padding: 30rpx;
  border-radius: 12rpx;
  border-left: 6rpx solid #3788d9;
}

.container-card {
  border-left: 6rpx solid;
}

.item-icon {
  width: 80rpx;
  height: 80rpx;
  background-color: #f0f0f0;
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 25rpx;
}

.iconfont {
  font-size: 40rpx;
}

.item-info {
  flex: 1;
}

.item-name {
  display: block;
  font-size: 30rpx;
  color: #333;
  margin-bottom: 8rpx;
}

.item-meta {
  display: block;
  font-size: 24rpx;
  color: #999;
  margin-bottom: 4rpx;
}

.item-location {
  display: block;
  font-size: 22rpx;
  color: #007AFF;
}

.item-arrow {
  font-size: 40rpx;
  color: #ccc;
}

.action-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background-color: #fff;
  padding: 20rpx 30rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.1);
}

.action-bar.single-btn {
  display: flex;
  justify-content: center;
}

.action-bar.single-btn .action-btn {
  flex: 1;
}

.action-bar.multi-btn {
  display: flex;
  gap: 20rpx;
}

.action-bar.three-btn {
  display: flex;
  gap: 10rpx;
}

.action-bar.two-btn {
  display: flex;
  gap: 20rpx;
}

.action-btn.half {
  flex: 1;
}

.action-btn.third {
  flex: 1;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #007AFF;
  color: #fff;
  padding: 25rpx;
  border-radius: 12rpx;
}

.action-icon {
  font-size: 36rpx;
  margin-right: 10rpx;
}

.action-text {
  font-size: 30rpx;
}

.load-more-btn {
  text-align: center;
  padding: 30rpx;
  color: #007AFF;
  background-color: #fff;
  margin: 0 20rpx;
  border-radius: 10rpx;
}

/* 弹窗样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 999;
}

.modal-container {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 85%;
  max-height: 80vh;
  background-color: #fff;
  border-radius: 16rpx;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 30rpx;
  border-bottom: 1rpx solid #eee;
}

.modal-title {
  font-size: 32rpx;
  font-weight: bold;
  color: #333;
}

.modal-close {
  font-size: 50rpx;
  color: #999;
  line-height: 1;
}

.modal-content {
  padding: 30rpx;
  max-height: 55vh;
  overflow-y: auto;
}

.modal-footer {
  display: flex;
  border-top: 1rpx solid #eee;
}

.form-row {
  margin-bottom: 30rpx;
}

.form-label {
  display: block;
  font-size: 28rpx;
  color: #666;
  margin-bottom: 15rpx;
}

.form-input {
  width: 100%;
  height: 80rpx;
  padding: 0 20rpx;
  background-color: #f5f5f5;
  border-radius: 10rpx;
  font-size: 28rpx;
  box-sizing: border-box;
}

.form-textarea {
  width: 100%;
  padding: 20rpx;
  background-color: #f5f5f5;
  border-radius: 10rpx;
  font-size: 28rpx;
  min-height: 120rpx;
  box-sizing: border-box;
}

.btn-cancel, .btn-confirm {
  flex: 1;
  text-align: center;
  padding: 30rpx;
  font-size: 30rpx;
}

.btn-cancel {
  color: #666;
  border-right: 1rpx solid #eee;
}

.btn-confirm {
  color: #007AFF;
}

.required {
  color: #e74c3c;
}

.picker-box {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20rpx;
  background-color: #f5f5f5;
  border-radius: 10rpx;
}

.picker-value {
  font-size: 28rpx;
  color: #333;
}

.picker-arrow {
  font-size: 36rpx;
  color: #999;
}

/* 父容器选择列表 */
.picker-item {
  display: flex;
  align-items: center;
  padding: 25rpx;
  border-bottom: 1rpx solid #f5f5f5;
  background-color: #fff;
}

.picker-item:last-child {
  border-bottom: none;
}

.picker-item.active {
  background-color: #f0f7ff;
}

.picker-item-icon {
  font-size: 36rpx;
  margin-right: 20rpx;
}

.picker-item-text {
  flex: 1;
  font-size: 28rpx;
  color: #333;
}

.picker-item-check {
  color: #007AFF;
  font-size: 32rpx;
  font-weight: bold;
}

/* Tab搜索按钮 */
.tab-search {
  width: 80rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.search-icon {
  font-size: 36rpx;
}

/* 搜索弹窗 */
.search-modal {
  width: 70%;
}

.mt-20 {
  margin-top: 20rpx;
}
</style>
