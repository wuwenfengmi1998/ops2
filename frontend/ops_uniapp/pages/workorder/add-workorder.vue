<template>
  <view class="container">
    <view class="header">
      <text class="back-btn" @click="goBack">‹ 返回</text>
      <text class="title">新建工单</text>
      <view class="header-right"></view>
    </view>
    
    <scroll-view scroll-y class="content">
      <!-- 表单卡片 -->
      <view class="card">
        <!-- 工单标题 -->
        <view class="form-item">
          <text class="form-label">工单标题 <text class="required">*</text></text>
          <input
            v-model="form.title"
            class="form-input"
            type="text"
            maxlength="200"
            placeholder="请输入工单标题"
          />
        </view>
        
        <!-- 问题描述 -->
        <view class="form-item">
          <text class="form-label">问题描述</text>
          <textarea
            v-model="form.description"
            class="form-textarea"
            placeholder="请输入问题描述"
          />
        </view>
        
        <!-- 关联物品 -->
        <view class="form-item">
          <text class="form-label">关联物品</text>
          
          <!-- 已选择物品 -->
          <view v-if="selectedItems.length > 0" class="selected-items">
            <view
              v-for="item in selectedItems"
              :key="item.ID"
              class="selected-item-tag"
            >
              <view class="selected-item-tag-info">
                <text class="selected-item-tag-name">{{ item.Name }}</text>
                <text v-if="item.SerialNumber" class="selected-item-tag-serial">{{ item.SerialNumber }}</text>
              </view>
              <text class="remove-item" @click="removeSelectedItem(item.ID)">×</text>
            </view>
          </view>
          
          <!-- 搜索框 -->
          <view class="search-box">
            <input
              v-model="searchQuery"
              class="form-input"
              type="text"
              placeholder="搜索物品名称或编号"
              @input="onSearchInput"
              @focus="onSearchFocus"
            />
            
            <!-- 搜索结果下拉 -->
            <view v-if="showDropdown && searchResults.length > 0" class="dropdown">
              <view
                v-for="item in searchResults"
                :key="item.ID"
                class="dropdown-item"
                @click="selectItem(item)"
              >
                <view class="dropdown-item-info">
                  <text class="dropdown-name">{{ item.Name }}</text>
                  <text v-if="item.SerialNumber" class="dropdown-serial">{{ item.SerialNumber }}</text>
                </view>
              </view>
            </view>
            
            <view v-if="showDropdown && searchLoading" class="dropdown">
              <view class="dropdown-loading">搜索中...</view>
            </view>
            
            <view v-if="showDropdown && searchResults.length === 0 && !searchLoading && searchQuery.trim()" class="dropdown">
              <view class="dropdown-empty">未找到匹配的物品</view>
            </view>
          </view>
        </view>
        
        <!-- 关联客户 -->
        <view class="form-item">
          <text class="form-label">关联客户</text>
          
          <!-- 已选择客户 -->
          <view v-if="selectedCustomers.length > 0" class="selected-customers">
            <view
              v-for="customer in selectedCustomers"
              :key="customer.id"
              class="selected-customer-tag"
            >
              <text class="customer-name">{{ (customer.last_name || '') + (customer.first_name ? ' ' + customer.first_name : '') }}</text>
              <text class="remove-customer" @click="removeSelectedCustomer(customer.id)">×</text>
            </view>
          </view>
          
          <!-- 搜索框 -->
          <view class="search-box">
            <input
              v-model="customerSearchQuery"
              class="form-input"
              type="text"
              placeholder="搜索客户姓名"
              @input="onCustomerSearchInput"
              @focus="onCustomerSearchFocus"
            />
            
            <!-- 搜索结果下拉 -->
            <view v-if="showCustomerDropdown && customerSearchResults.length > 0" class="dropdown">
              <view
                v-for="customer in customerSearchResults"
                :key="customer.id"
                class="dropdown-item"
                @click="selectCustomer(customer)"
              >
                <text class="dropdown-name">{{ (customer.last_name || '') + (customer.first_name ? ' ' + customer.first_name : '') }}</text>
                <text v-if="customer.primary_phone" class="dropdown-phone">{{ customer.primary_phone }}</text>
              </view>
            </view>
            
            <view v-if="showCustomerDropdown && customerSearchLoading" class="dropdown">
              <view class="dropdown-loading">搜索中...</view>
            </view>
            
            <view v-if="showCustomerDropdown && customerSearchResults.length === 0 && !customerSearchLoading && customerSearchQuery.trim()" class="dropdown">
              <view class="dropdown-empty">未找到匹配的客户</view>
            </view>
          </view>
        </view>
        
        <!-- 图片上传 -->
        <view class="form-item">
          <text class="form-label">图片</text>
          <view class="photo-upload">
            <view
              v-for="(hash, index) in photos"
              :key="index"
              class="photo-item"
              @click="previewImage(index)"
            >
              <image
                class="photo-img"
                :src="getPhotoUrl(hash)"
                mode="aspectFill"
              />
              <view class="photo-remove" @click.stop="removePhoto(index)">×</view>
            </view>
            <view v-if="photos.length < 9" class="photo-add" @click="chooseImage">
              <text class="photo-add-icon">+</text>
              <text class="photo-add-text">添加图片</text>
            </view>
          </view>
        </view>
      </view>
      
      <!-- 提交按钮 -->
      <view class="submit-bar">
        <view class="submit-btn" :class="{ disabled: submitting }" @click="submitForm">
          <text v-if="submitting">提交中...</text>
          <text v-else>提交工单</text>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref, reactive } from 'vue'
import api from '@/api/index.js'
import { workOrderApi } from '@/api/work_order.js'
import { warehouseApi } from '@/api/warehouse.js'
import { customerApi } from '@/api/customer.js'
import { useConfigStore } from '@/stores/config.js'
import { onLoad } from '@dcloudio/uni-app'

const configStore = useConfigStore()
const goBack = () => uni.navigateBack()

// 预填物品ID（从物品详情跳转时）
const presetItemIds = ref([])

onLoad((options) => {
  // 优先检查预填数据（从物品详情跳转时）
  const prefillStr = uni.getStorageSync('prefill_work_order')
  if (prefillStr) {
    try {
      const prefill = JSON.parse(prefillStr)
      form.title = prefill.title || ''
      form.description = prefill.description || ''
      if (prefill.itemId) {
        presetItemIds.value = [prefill.itemId]
        fetchPresetItem(prefill.itemId)
      }
      // 如果有预填客户信息，自动填充
      if (prefill.customers && prefill.customers.length > 0) {
        selectedCustomers.value = prefill.customers.map(c => ({
          id: c.id,
          first_name: c.first_name || '',
          last_name: c.last_name || '',
          primary_phone: c.primary_phone || ''
        }))
      }
      uni.removeStorageSync('prefill_work_order')
    } catch (e) {
      console.error('读取预填数据失败', e)
    }
  } else if (options && options.item_id) {
    // 兼容旧方式：直接传 item_id 参数
    presetItemIds.value = [parseInt(options.item_id)]
    fetchPresetItem(parseInt(options.item_id))
  }
})

async function fetchPresetItem(itemId) {
  try {
    const res = await warehouseApi.getItem(itemId)
    if (res.errCode === 0 && res.data && res.data.item) {
      const item = res.data.item
      // 检查是否已选中
      if (!selectedItems.value.find(i => i.ID === item.ID)) {
        selectedItems.value.push(item)
        linkedItemIds.value.push(item.ID)
      }
    }
  } catch (e) {
    console.error('获取物品信息失败', e)
  }
}

// 表单数据
const form = reactive({
  title: '',
  description: ''
})

// 关联物品
const selectedItems = ref([])
const linkedItemIds = ref([])
const searchQuery = ref('')
const searchResults = ref([])
const showDropdown = ref(false)
const searchLoading = ref(false)
let searchTimer = null

// 关联客户
const selectedCustomers = ref([])
const customerSearchQuery = ref('')
const customerSearchResults = ref([])
const showCustomerDropdown = ref(false)
const customerSearchLoading = ref(false)
let customerSearchTimer = null

// 图片
const photos = ref([])

// 提交状态
const submitting = ref(false)

function onSearchInput() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    doSearch()
  }, 300)
}

function onSearchFocus() {
  if (searchQuery.value.trim()) {
    showDropdown.value = true
  }
}

async function doSearch() {
  if (!searchQuery.value.trim()) {
    searchResults.value = []
    showDropdown.value = false
    return
  }
  
  searchLoading.value = true
  showDropdown.value = true
  
  try {
    const res = await warehouseApi.listItem({
      search: searchQuery.value.trim()
    })
    if (res.errCode === 0 && res.data) {
      searchResults.value = (res.data.items || []).slice(0, 10)
    } else {
      searchResults.value = []
    }
  } catch (e) {
    console.error('搜索物品失败', e)
    searchResults.value = []
  } finally {
    searchLoading.value = false
  }
}

function selectItem(item) {
  // 检查是否已选中
  if (!selectedItems.value.find(i => i.ID === item.ID)) {
    selectedItems.value.push(item)
    linkedItemIds.value.push(item.ID)
  }
  searchQuery.value = ''
  searchResults.value = []
  showDropdown.value = false
}

function removeSelectedItem(itemId) {
  selectedItems.value = selectedItems.value.filter(i => i.ID !== itemId)
  linkedItemIds.value = linkedItemIds.value.filter(id => id !== itemId)
}

function onCustomerSearchInput() {
  clearTimeout(customerSearchTimer)
  customerSearchTimer = setTimeout(() => {
    doCustomerSearch()
  }, 300)
}

function onCustomerSearchFocus() {
  if (customerSearchQuery.value.trim()) {
    showCustomerDropdown.value = true
  }
}

async function doCustomerSearch() {
  if (!customerSearchQuery.value.trim()) {
    customerSearchResults.value = []
    showCustomerDropdown.value = false
    return
  }
  
  customerSearchLoading.value = true
  showCustomerDropdown.value = true
  
  try {
    const res = await customerApi.list({
      search: customerSearchQuery.value.trim(),
      page: 1,
      page_size: 10
    })
    if (res.errCode === 0 && res.data) {
      customerSearchResults.value = (res.data.customers || []).slice(0, 10)
    } else {
      customerSearchResults.value = []
    }
  } catch (e) {
    console.error('搜索客户失败', e)
    customerSearchResults.value = []
  } finally {
    customerSearchLoading.value = false
  }
}

function selectCustomer(customer) {
  // 检查是否已选中
  if (!selectedCustomers.value.find(c => c.id === customer.id)) {
    selectedCustomers.value.push(customer)
  }
  customerSearchQuery.value = ''
  customerSearchResults.value = []
  showCustomerDropdown.value = false
}

function removeSelectedCustomer(customerId) {
  selectedCustomers.value = selectedCustomers.value.filter(c => c.id !== customerId)
}

function getPhotoUrl(hash) {
  return configStore.getFileBaseUrl() + '/api/files/get/' + hash
}

function chooseImage() {
  uni.chooseImage({
    count: 9 - photos.value.length,
    success: (res) => {
      res.tempFilePaths.forEach(path => {
        uploadImage(path)
      })
    }
  })
}

async function uploadImage(filePath) {
  try {
    uni.showLoading({ title: '上传中...' })
    const res = await api.upload('/files/upload/image', {
      uri: filePath,
      name: 'file'
    })
    
    if (res.errCode === 0 && res.data && res.data.hash) {
      photos.value.push(res.data.hash)
      uni.showToast({ title: '上传成功', icon: 'success' })
    } else {
      uni.showToast({ title: res.message || '上传失败', icon: 'none' })
    }
  } catch (e) {
    console.error('上传失败', e)
    uni.showToast({ title: '上传失败', icon: 'none' })
  } finally {
    uni.hideLoading()
  }
}

function removePhoto(index) {
  photos.value.splice(index, 1)
}

function previewImage(index) {
  const urls = photos.value.map(hash => getPhotoUrl(hash))
  uni.previewImage({
    current: index,
    urls: urls
  })
}

async function submitForm() {
  if (!form.title.trim()) {
    uni.showToast({ title: '请输入工单标题', icon: 'none' })
    return
  }
  
  if (submitting.value) return
  submitting.value = true
  
  try {
    const data = {
      title: form.title.trim(),
      description: form.description.trim(),
      photos: photos.value,
      customer_ids: selectedCustomers.value.map(c => c.id)
    }
    
    if (linkedItemIds.value.length > 0) {
      data.item_ids = linkedItemIds.value
    }
    
    const res = await workOrderApi.add(data)
    
    if (res.errCode === 0) {
      uni.showToast({ title: '提交成功', icon: 'success' })
      uni.$emit('workorder-refresh')
      setTimeout(() => {
        uni.navigateBack()
      }, 500)
    } else {
      uni.showToast({ title: res.errMsg || '提交失败', icon: 'none' })
    }
  } catch (e) {
    console.error('提交工单失败', e)
    uni.showToast({ title: '提交失败', icon: 'none' })
  } finally {
    submitting.value = false
  }
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
  margin-right: 20rpx;
}

.title {
  font-size: 36rpx;
  font-weight: bold;
  color: #333;
  flex: 1;
  text-align: center;
  padding-right: 60rpx;
}

.header-right {
  width: 60rpx;
}

.content {
  height: calc(100vh - 180rpx);
  padding: 20rpx;
}

.card {
  background-color: #fff;
  border-radius: 12rpx;
  padding: 30rpx;
}

.form-item {
  margin-bottom: 30rpx;
}

.form-item:last-child {
  margin-bottom: 0;
}

.form-label {
  display: block;
  font-size: 28rpx;
  color: #333;
  margin-bottom: 15rpx;
}

.required {
  color: #ff4d4f;
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
  min-height: 200rpx;
  padding: 20rpx;
  background-color: #f5f5f5;
  border-radius: 10rpx;
  font-size: 28rpx;
  box-sizing: border-box;
}

.selected-customers {
  display: flex;
  flex-wrap: wrap;
  gap: 15rpx;
  margin-bottom: 20rpx;
}

.selected-customer-tag {
  display: flex;
  align-items: center;
  gap: 10rpx;
  padding: 15rpx 20rpx;
  background-color: #e6f7ff;
  border: 1rpx solid #91d5ff;
  border-radius: 30rpx;
  font-size: 26rpx;
  color: #1890ff;
}

.customer-name {
  display: block;
  font-size: 28rpx;
  color: #333;
}

.dropdown-phone {
  display: block;
  font-size: 24rpx;
  color: #999;
  margin-top: 5rpx;
}

.remove-customer {
  font-size: 28rpx;
  color: #ff4d4f;
  padding: 5rpx;
}

.search-box {
  position: relative;
}

.dropdown {
  position: absolute;
  top: 90rpx;
  left: 0;
  right: 0;
  background-color: #fff;
  border: 1rpx solid #e5e5e5;
  border-radius: 10rpx;
  max-height: 400rpx;
  overflow-y: auto;
  z-index: 100;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.1);
}

.dropdown-item {
  padding: 20rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.dropdown-item:last-child {
  border-bottom: none;
}

.dropdown-name {
  display: block;
  font-size: 28rpx;
  color: #333;
}

.dropdown-serial {
  display: block;
  font-size: 24rpx;
  color: #999;
  margin-top: 5rpx;
}

.dropdown-loading,
.dropdown-empty {
  padding: 30rpx;
  text-align: center;
  font-size: 26rpx;
  color: #999;
}

.photo-upload {
  display: flex;
  flex-wrap: wrap;
  gap: 20rpx;
}

.photo-item {
  width: 200rpx;
  height: 200rpx;
  border-radius: 10rpx;
  overflow: hidden;
  position: relative;
}

.photo-img {
  width: 100%;
  height: 100%;
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

.selected-items {
  display: flex;
  flex-wrap: wrap;
  gap: 15rpx;
  margin-bottom: 20rpx;
}

.selected-item-tag {
  display: flex;
  align-items: flex-start;
  gap: 10rpx;
  padding: 15rpx 20rpx;
  background-color: #e6f7ff;
  border: 1rpx solid #91d5ff;
  border-radius: 10rpx;
  font-size: 26rpx;
  color: #1890ff;
  max-width: 100%;
}

.selected-item-tag-info {
  flex: 1;
  overflow: hidden;
}

.selected-item-tag-name {
  display: block;
  font-size: 26rpx;
  color: #1890ff;
  word-break: break-all;
}

.selected-item-tag-serial {
  display: block;
  font-size: 22rpx;
  color: #8c8c8c;
  margin-top: 5rpx;
  word-break: break-all;
}

.remove-item {
  font-size: 28rpx;
  color: #ff4d4f;
  padding: 5rpx;
  flex-shrink: 0;
}

.dropdown-item-info {
  display: flex;
  flex-direction: column;
}

.submit-bar {
  margin-top: 40rpx;
  padding: 0 20rpx;
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
</style>
