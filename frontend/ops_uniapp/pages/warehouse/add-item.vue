<template>
  <view class="container">
    <view class="header">
      <text class="back-btn" @click="goBack">‹ 返回</text>
      <text class="title">新增物品</text>
    </view>
    
    <scroll-view scroll-y class="content">
      <view class="form-card">
        <view class="form-item">
          <text class="form-label required">名称 *</text>
          <input class="form-input" v-model="form.name" placeholder="请输入物品名称" />
        </view>
        
        <view class="form-item">
          <text class="form-label">编号</text>
          <input class="form-input" v-model="form.serialNumber" placeholder="请输入物品编号（选填）" />
        </view>
        
        <view class="form-item">
          <text class="form-label">数量</text>
          <view class="quantity-row">
            <view class="qty-btn" @click="form.quantity > 1 && form.quantity--">
              <text class="qty-btn-text">−</text>
            </view>
            <input class="qty-input" type="number" v-model.number="form.quantity" />
            <view class="qty-btn" @click="form.quantity++">
              <text class="qty-btn-text">+</text>
            </view>
          </view>
        </view>
        
        <view class="form-item">
          <text class="form-label">备注</text>
          <textarea class="form-textarea" v-model="form.remark" placeholder="请输入备注（选填）" />
        </view>
        
        <!-- 关联客户 -->
        <view class="form-item">
          <text class="form-label">关联客户</text>
          
          <!-- 已选中的客户 -->
          <view v-if="selectedCustomers.length > 0" class="selected-customers">
            <view 
              v-for="customer in selectedCustomers" 
              :key="customer.id" 
              class="customer-tag"
            >
              <text>{{ (customer.last_name || '') + (customer.first_name ? ' ' + customer.first_name : '') }}</text>
              <text class="customer-tag-remove" @click="removeCustomer(customer.id)">×</text>
            </view>
          </view>
          
          <!-- 搜索框 -->
          <view class="customer-search-box">
            <input 
              class="customer-search-input"
              v-model="customerSearchQuery" 
              placeholder="搜索客户..."
              @input="onCustomerSearchInput"
              @focus="onCustomerSearchFocus"
            />
          </view>
          
          <!-- 搜索结果 -->
          <view v-if="showCustomerDropdown && customerSearchResults.length > 0" class="customer-dropdown">
            <view 
              v-for="customer in customerSearchResults" 
              :key="customer.id" 
              class="customer-dropdown-item"
              @click="selectCustomer(customer)"
            >
              <text class="customer-dropdown-name">{{ (customer.last_name || '') + (customer.first_name ? ' ' + customer.first_name : '') }}</text>
              <text v-if="customer.primary_phone" class="customer-dropdown-phone">{{ customer.primary_phone }}</text>
            </view>
          </view>
          <view v-else-if="showCustomerDropdown && customerSearchQuery && customerSearchResults.length === 0 && !customerSearchLoading" class="customer-dropdown-empty">
            未找到匹配的客户
          </view>
        </view>
      </view>
      
      <view class="form-card">
        <view class="form-item">
          <text class="form-label">物品图片</text>
          <view class="image-list">
            <view
              v-for="(photo, index) in form.photos"
              :key="index"
              class="image-item"
            >
              <image class="preview-img" :src="getPhotoUrl(photo)" mode="aspectFill" />
              <view class="image-remove" @click="removePhoto(index)">×</view>
            </view>
            <view v-if="form.photos.length < 9" class="image-add" @click="chooseImage">
              <text class="add-icon">+</text>
              <text class="add-text">添加图片</text>
            </view>
          </view>
        </view>
      </view>
      
      <view class="submit-section">
        <view class="submit-btn" @click="submitForm">
          <text>保存</text>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { warehouseApi } from '@/api/warehouse.js'
import { customerApi } from '@/api/customer.js'
import { useConfigStore } from '@/stores/config.js'
import api from '@/api/index.js'

const configStore = useConfigStore()

// 容器ID（从URL参数获取）
const containerId = ref(0)

// 表单数据
const form = ref({
  name: '',
  serialNumber: '',
  remark: '',
  quantity: 1,
  photos: []
})

// 关联客户
const selectedCustomers = ref([])
const customerSearchQuery = ref('')
const customerSearchResults = ref([])
const showCustomerDropdown = ref(false)
const customerSearchLoading = ref(false)
let customerSearchTimer = null

function getPhotoUrl(sha256) {
  const baseUrl = configStore.getFileBaseUrl ? configStore.getFileBaseUrl() : ''
  return `${baseUrl}/api/files/get/${sha256}`
}

function goBack() {
  uni.navigateBack()
}

function chooseImage() {
  uni.chooseImage({
    count: 9 - form.value.photos.length,
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
      form.value.photos.push(res.data.hash)
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
  form.value.photos.splice(index, 1)
}

// 客户搜索
async function searchCustomers() {
  if (!customerSearchQuery.value.trim()) {
    customerSearchResults.value = []
    return
  }
  customerSearchLoading.value = true
  try {
    const res = await customerApi.list(1, 20, customerSearchQuery.value)
    if (res.errCode === 0 && res.data) {
      customerSearchResults.value = res.data.customers || res.data || []
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

function onCustomerSearchInput() {
  clearTimeout(customerSearchTimer)
  customerSearchTimer = setTimeout(() => {
    searchCustomers()
  }, 300)
}

function onCustomerSearchFocus() {
  showCustomerDropdown.value = true
  searchCustomers()
}

function selectCustomer(customer) {
  const exists = selectedCustomers.value.find(c => c.id === customer.id)
  if (!exists) {
    selectedCustomers.value.push(customer)
  }
  customerSearchQuery.value = ''
  showCustomerDropdown.value = false
  customerSearchResults.value = []
}

function removeCustomer(id) {
  const index = selectedCustomers.value.findIndex(c => c.id === id)
  if (index >= 0) {
    selectedCustomers.value.splice(index, 1)
  }
}

async function submitForm() {
  if (!form.value.name.trim()) {
    uni.showToast({ title: '请输入名称', icon: 'none' })
    return
  }
  
  try {
    uni.showLoading({ title: '保存中...' })
    const data = {
      name: form.value.name,
      serial_number: form.value.serialNumber,
      remark: form.value.remark,
      quantity: form.value.quantity > 0 ? form.value.quantity : 1,
      container_id: containerId.value || null,
      photos: form.value.photos,
      customer_ids: selectedCustomers.value.map(c => c.id)
    }
    
    const res = await warehouseApi.addItem(data)
    
    if (res.errCode === 0) {
      uni.showToast({ title: '保存成功', icon: 'success' })
      setTimeout(() => {
        goBack()
      }, 1500)
    } else {
      uni.showToast({ title: res.errMsg || '保存失败', icon: 'none' })
    }
  } catch (e) {
    console.error('保存失败', e)
    uni.showToast({ title: '保存失败', icon: 'none' })
  } finally {
    uni.hideLoading()
  }
}

// 获取页面参数
import { onLoad } from '@dcloudio/uni-app'
onLoad((options) => {
  if (options && options.container_id) {
    containerId.value = parseInt(options.container_id)
  }
})
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

.content {
  height: calc(100vh - 120rpx);
}

.form-card {
  background-color: #fff;
  margin: 20rpx;
  border-radius: 12rpx;
  padding: 20rpx;
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
  color: #666;
  margin-bottom: 15rpx;
}

.form-label.required::after {
  content: ' *';
  color: #e74c3c;
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
  min-height: 150rpx;
  box-sizing: border-box;
}

/* 图片上传 */
.image-list {
  display: flex;
  flex-wrap: wrap;
  gap: 20rpx;
}

.image-item {
  position: relative;
  width: 180rpx;
  height: 180rpx;
  border-radius: 12rpx;
  overflow: hidden;
}

.preview-img {
  width: 100%;
  height: 100%;
}

.image-remove {
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
  border-radius: 0 0 0 12rpx;
}

.image-add {
  width: 180rpx;
  height: 180rpx;
  background-color: #f5f5f5;
  border-radius: 12rpx;
  border: 2rpx dashed #ddd;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.add-icon {
  font-size: 60rpx;
  color: #999;
  line-height: 1;
}

.add-text {
  font-size: 22rpx;
  color: #999;
  margin-top: 10rpx;
}

.submit-section {
  padding: 40rpx 20rpx;
}

.submit-btn {
  background-color: #007AFF;
  color: #fff;
  text-align: center;
  padding: 30rpx;
  border-radius: 12rpx;
  font-size: 32rpx;
}

/* 数量选择器 */
.quantity-row {
  display: flex;
  align-items: center;
  gap: 0;
}

.qty-btn {
  width: 72rpx;
  height: 72rpx;
  background-color: #f0f0f0;
  border-radius: 10rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.qty-btn-text {
  font-size: 36rpx;
  color: #333;
  line-height: 1;
}

.qty-input {
  width: 100rpx;
  height: 72rpx;
  text-align: center;
  background-color: #f5f5f5;
  font-size: 30rpx;
  border-radius: 10rpx;
  margin: 0 16rpx;
  box-sizing: border-box;
}

/* 关联客户 */
.selected-customers {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
  margin-bottom: 16rpx;
}

.customer-tag {
  display: inline-flex;
  align-items: center;
  gap: 8rpx;
  padding: 8rpx 16rpx;
  background-color: #e6f7ff;
  border: 1rpx solid #91d5ff;
  border-radius: 20rpx;
  font-size: 24rpx;
  color: #1890ff;
}

.customer-tag-remove {
  font-size: 28rpx;
  color: #999;
  margin-left: 4rpx;
}

.customer-tag-remove:active {
  color: #f56c6c;
}

.customer-search-box {
  margin-top: 12rpx;
}

.customer-search-input {
  width: 100%;
  height: 72rpx;
  padding: 0 20rpx;
  background-color: #f5f5f5;
  border-radius: 10rpx;
  font-size: 28rpx;
  box-sizing: border-box;
}

.customer-dropdown {
  margin-top: 12rpx;
  background-color: #fff;
  border: 1rpx solid #eee;
  border-radius: 10rpx;
  max-height: 400rpx;
  overflow-y: auto;
}

.customer-dropdown-item {
  display: flex;
  align-items: center;
  padding: 20rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.customer-dropdown-item:last-child {
  border-bottom: none;
}

.customer-dropdown-item:active {
  background-color: #f5f5f5;
}

.customer-dropdown-name {
  flex: 1;
  font-size: 28rpx;
  color: #333;
}

.customer-dropdown-phone {
  font-size: 24rpx;
  color: #999;
  margin-left: 16rpx;
}

.customer-dropdown-empty {
  padding: 30rpx;
  text-align: center;
  color: #999;
  font-size: 26rpx;
}
</style>
