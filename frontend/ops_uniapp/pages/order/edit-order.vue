<template>
  <view class="container">
    <!-- 页面标题 -->
    <view class="page-header">
      <text class="page-title">编辑订单</text>
    </view>

    <scroll-view scroll-y class="form-content" v-if="!loading">
      <!-- 基本信息 -->
      <view class="card">
        <view class="card-title">基本信息</view>
        
        <!-- 标题（必填） -->
        <view class="form-item">
          <text class="form-label required">配件名称</text>
          <input 
            class="form-input" 
            v-model="form.title" 
            placeholder="请输入配件名称"
            maxlength="50"
          />
        </view>

        <!-- 备注 -->
        <view class="form-item">
          <text class="form-label">备注</text>
          <textarea 
            class="form-textarea" 
            v-model="form.remark" 
            placeholder="请输入备注信息"
            maxlength="256"
          />
          <text class="char-count">{{ form.remark.length }}/256</text>
        </view>

        <!-- 采购链接 -->
        <view class="form-item">
          <text class="form-label">采购链接</text>
          <textarea 
            class="form-textarea" 
            v-model="form.link" 
            placeholder="请输入采购链接"
            rows="2"
          />
        </view>

        <!-- 款式标签 -->
        <view class="form-item">
          <text class="form-label">款式标签</text>
          <view class="tags-input">
            <view class="tag-list">
              <view 
                class="tag" 
                v-for="(tag, index) in tags" 
                :key="index"
              >
                {{ tag }}
                <text class="tag-remove" @click="removeTag(index)">×</text>
              </view>
            </view>
            <input 
              class="tag-input" 
              v-model="newTag" 
              placeholder="输入标签后回车添加"
              @confirm="addTag"
            />
          </view>
        </view>
      </view>

      <!-- 费用明细 -->
      <view class="card">
        <view class="card-title">费用明细</view>
        
        <!-- 已添加的费用列表 -->
        <view class="cost-table" v-if="costEntries.length > 0">
          <view class="cost-header">
            <text class="cost-col">类型</text>
            <text class="cost-col">数量</text>
            <text class="cost-col">单价</text>
            <text class="cost-col">总计</text>
            <text class="cost-col">货币</text>
            <text class="cost-col">操作</text>
          </view>
          <view 
            class="cost-row" 
            v-for="(item, index) in costEntries" 
            :key="index"
          >
            <text class="cost-col">{{ costType[item.type] || item.type }}</text>
            <text class="cost-col">{{ item.int }}</text>
            <text class="cost-col">{{ item.cost }}</text>
            <text class="cost-col">{{ item.costt }}</text>
            <text class="cost-col">{{ currencyOptions[item.currencytype] || item.currencytype }}</text>
            <text class="cost-col delete" @click="removeCostEntry(index)">删除</text>
          </view>
        </view>

        <!-- 添加费用表单 -->
        <view class="cost-form">
          <view class="cost-form-row">
            <view class="cost-field">
              <text class="cost-label">费用类型</text>
              <picker mode="selector" :range="costTypeOptions" range-key="label" @change="onCostTypeChange">
                <view class="cost-picker">
                  {{ costTypeOptions[costTypeIndex].label }}
                </view>
              </picker>
            </view>
            <view class="cost-field">
              <text class="cost-label">数量</text>
              <input class="cost-input" type="number" v-model="newCost.int" min="1" />
            </view>
          </view>
          <view class="cost-form-row">
            <view class="cost-field">
              <text class="cost-label">单价</text>
              <input class="cost-input" type="digit" v-model="newCost.cost" step="0.01" />
            </view>
            <view class="cost-field">
              <text class="cost-label">货币</text>
              <picker mode="selector" :range="currencyOptionsList" range-key="label" @change="onCurrencyChange">
                <view class="cost-picker">
                  {{ currencyOptionsList[currencyIndex].label }}
                </view>
              </picker>
            </view>
          </view>
          <view class="cost-total" v-if="newCostCost > 0">
            总计：{{ currencyOptionsList[currencyIndex].symbol }}{{ newCostTotal }}
          </view>
          <view class="cost-error" v-if="costError">单价必须大于0</view>
          <button class="add-cost-btn" @click="addCostEntry">添加费用</button>
        </view>
      </view>

      <!-- 图片上传 -->
      <view class="card">
        <view class="card-title">图片</view>
        <view class="photo-upload">
          <view
            v-for="(photo, index) in photos"
            :key="index"
            class="photo-item"
            @click="previewImage(index)"
          >
            <image
              class="photo-img"
              :src="photo.url"
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

      <!-- 提交按钮 -->
      <view class="submit-area">
        <button class="submit-btn" :disabled="submitting" @click="handleSubmit">
          {{ submitting ? '提交中...' : '保存修改' }}
        </button>
      </view>
    </scroll-view>
    
    <view v-else class="loading-view"><text>加载中...</text></view>
  </view>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { purchaseApi } from '@/api/purchase.js'
import { api } from '@/api/index.js'
import { useConfigStore } from '@/stores/config.js'

const configStore = useConfigStore()

// 订单ID
const orderId = ref(null)

// 加载状态
const loading = ref(true)

// 费用类型选项
const costTypeOptions = [
  { value: 1, label: '单价' },
  { value: 2, label: '运费' }
]
const costTypeIndex = ref(0)
const costType = { 1: '单价', 2: '运费' }

// 货币选项
const currencyOptionsList = [
  { value: 1, label: 'CNY', symbol: '¥' },
  { value: 2, label: 'MOP', symbol: 'MOP' },
  { value: 3, label: 'HKD', symbol: 'HK$' },
  { value: 4, label: 'USD', symbol: '$' }
]
const currencyOptions = { 1: 'CNY', 2: 'MOP', 3: 'HKD', 4: 'USD' }
const currencyIndex = ref(0)

// 表单数据
const form = reactive({
  title: '',
  remark: '',
  link: '',
  styles: '',
  photos: []
})

// 标签
const tags = ref([])
const newTag = ref('')

// 费用明细
const costEntries = reactive([])
const newCost = reactive({
  type: 1,
  int: 1,
  cost: 0,
  currencyType: 1
})
const costError = ref(false)

// 图片
const photos = ref([])
const newCostCost = computed(() => parseFloat(newCost.cost) || 0)
const newCostTotal = computed(() => (newCost.int * newCostCost.value).toFixed(2))

// 提交状态
const submitting = ref(false)

// 初始化页面
onLoad((options) => {
  if (options?.id) {
    orderId.value = parseInt(options.id)
    fetchOrderDetail()
  } else {
    uni.showToast({ title: '缺少订单ID', icon: 'none' })
    setTimeout(() => uni.navigateBack(), 1500)
  }
})

// 获取订单详情
async function fetchOrderDetail() {
  loading.value = true
  try {
    const res = await purchaseApi.getOrder({ id: orderId.value })
    if (res.errCode === 0 && res.data?.order) {
      const order = res.data.order
      
      // 填充表单 (注意：Go 结构体字段首字母大写)
      form.title = order.Title || ''
      form.remark = order.Remark || ''
      form.link = order.Link || ''
      form.styles = order.Styles || ''
      
      // 解析标签
      if (order.Styles) {
        tags.value = order.Styles.split(',').filter(t => t.trim())
      }
      
      // 解析费用 (API返回的是数据库模型 Price/Quantity/CurrencyType/CostType)
      costEntries.length = 0
      if (res.data.costs && res.data.costs.length > 0) {
        res.data.costs.forEach(c => {
          costEntries.push({
            type: c.CostType,
            int: c.Quantity,
            cost: (c.Price / 100).toFixed(2),
            costt: ((c.Price * c.Quantity) / 100).toFixed(2),
            currencytype: c.CurrencyType
          })
        })
      }
      
      // 解析图片
      photos.value = []
      if (res.data.photos && res.data.photos.length > 0) {
        res.data.photos.forEach(p => {
          photos.value.push({
            hash: p.Sha256,
            url: `${configStore.getFileBaseUrl()}/api/files/get/${p.Sha256}`
          })
        })
      }
    } else {
      uni.showToast({ title: '获取订单信息失败', icon: 'none' })
    }
  } catch (e) {
    console.error('获取订单详情失败', e)
    uni.showToast({ title: '获取订单信息失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

// 添加标签
function addTag() {
  if (newTag.value.trim()) {
    tags.value.push(newTag.value.trim())
    newTag.value = ''
    form.styles = tags.value.join(',')
  }
}

// 删除标签
function removeTag(index) {
  tags.value.splice(index, 1)
  form.styles = tags.value.join(',')
}

// 费用类型选择
function onCostTypeChange(e) {
  costTypeIndex.value = e.detail.value
  newCost.type = costTypeOptions[e.detail.value].value
}

// 货币类型选择
function onCurrencyChange(e) {
  currencyIndex.value = e.detail.value
  newCost.currencyType = currencyOptionsList[e.detail.value].value
}

// 添加费用
function addCostEntry() {
  if (newCost.cost <= 0) {
    costError.value = true
    return
  }
  costError.value = false
  costEntries.push({
    type: newCost.type,
    int: newCost.int,
    cost: parseFloat(newCost.cost).toFixed(2),
    costt: newCostTotal.value,
    currencytype: newCost.currencyType
  })
  // 重置
  newCost.type = 1
  newCost.int = 1
  newCost.cost = 0
  newCost.currencyType = 1
  costTypeIndex.value = 0
  currencyIndex.value = 0
}

// 删除费用
function removeCostEntry(index) {
  costEntries.splice(index, 1)
}

// 选择图片
function chooseImage() {
  uni.chooseImage({
    count: 9 - photos.value.length,
    success: (res) => {
      res.tempFiles.forEach(file => {
        uploadImage(file.path)
      })
    }
  })
}

// 上传图片
async function uploadImage(filePath) {
  uni.showLoading({ title: '上传中...', mask: true })
  
  try {
    const res = await api.upload('/files/upload/image', { uri: filePath, name: 'file' })
    if (res.errCode === 0 && res.data?.hash) {
      photos.value.push({
        hash: res.data.hash,
        url: `${configStore.getFileBaseUrl()}/api/files/get/${res.data.hash}`
      })
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

// 预览图片
function previewImage(index) {
  const urls = photos.value.map(p => p.url)
  uni.previewImage({
    urls,
    current: index
  })
}

// 删除图片
function removePhoto(index) {
  photos.value.splice(index, 1)
}

// 提交表单
async function handleSubmit() {
  // 验证必填项
  if (!form.title.trim()) {
    uni.showToast({ title: '请输入配件名称', icon: 'none' })
    return
  }

  submitting.value = true
  
  try {
    // 转换费用数据（元转分）
    const costs = costEntries.map(h => ({
      type: h.type,
      int: h.int,
      cost: Math.round(parseFloat(h.cost) * 100),
      costt: Math.round(parseFloat(h.costt) * 100),
      currencytype: h.currencytype
    }))

    // 转换图片
    const photoHashes = photos.value.map(p => p.hash)

    const res = await purchaseApi.updateOrder({
      id: orderId.value,
      title: form.title,
      remark: form.remark,
      link: form.link,
      styles: form.styles,
      costs,
      photos: photoHashes
    })

    if (res.errCode === 0) {
      uni.showToast({ title: '保存成功', icon: 'success' })
      uni.$emit('purchase-refresh')
      setTimeout(() => {
        uni.navigateBack()
      }, 1500)
    } else {
      uni.showToast({ title: '保存失败', icon: 'none' })
    }
  } catch (e) {
    console.error('提交失败', e)
    uni.showToast({ title: '保存失败', icon: 'none' })
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

.page-header {
  background-color: #fff;
  padding: 30rpx;
  border-bottom: 1px solid #eee;
}

.page-title {
  font-size: 36rpx;
  font-weight: 600;
  color: #333;
}

.form-content {
  padding: 20rpx;
}

.card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
}

.card-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 30rpx;
  padding-bottom: 20rpx;
  border-bottom: 1px solid #eee;
}

.form-item {
  margin-bottom: 30rpx;
}

.form-label {
  display: block;
  font-size: 28rpx;
  color: #666;
  margin-bottom: 12rpx;
}

.form-label.required::after {
  content: ' *';
  color: #f56c6c;
}

.form-input {
  width: 100%;
  height: 80rpx;
  padding: 0 24rpx;
  font-size: 28rpx;
  background-color: #f8f8f8;
  border-radius: 12rpx;
  box-sizing: border-box;
}

.form-textarea {
  width: 100%;
  padding: 20rpx 24rpx;
  font-size: 28rpx;
  background-color: #f8f8f8;
  border-radius: 12rpx;
  box-sizing: border-box;
  min-height: 120rpx;
}

.char-count {
  display: block;
  text-align: right;
  font-size: 24rpx;
  color: #999;
  margin-top: 8rpx;
}

/* 标签输入 */
.tags-input {
  background-color: #f8f8f8;
  border-radius: 12rpx;
  padding: 16rpx;
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
  margin-bottom: 12rpx;
}

.tag {
  display: flex;
  align-items: center;
  background-color: #e6f0ff;
  color: #1890ff;
  padding: 8rpx 16rpx;
  border-radius: 8rpx;
  font-size: 24rpx;
}

.tag-remove {
  margin-left: 8rpx;
  font-size: 28rpx;
  color: #999;
}

.tag-input {
  width: 100%;
  height: 60rpx;
  font-size: 26rpx;
}

/* 费用明细 */
.cost-table {
  margin-bottom: 30rpx;
  border: 1px solid #eee;
  border-radius: 12rpx;
  overflow: hidden;
}

.cost-header,
.cost-row {
  display: flex;
  padding: 16rpx 12rpx;
  font-size: 24rpx;
}

.cost-header {
  background-color: #f8f8f8;
  color: #666;
}

.cost-row {
  border-top: 1px solid #eee;
  color: #333;
}

.cost-col {
  flex: 1;
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cost-col.delete {
  color: #f56c6c;
}

.cost-form {
  background-color: #f8f8f8;
  border-radius: 12rpx;
  padding: 20rpx;
}

.cost-form-row {
  display: flex;
  gap: 20rpx;
  margin-bottom: 20rpx;
}

.cost-field {
  flex: 1;
}

.cost-label {
  display: block;
  font-size: 24rpx;
  color: #666;
  margin-bottom: 8rpx;
}

.cost-input {
  width: 100%;
  height: 72rpx;
  padding: 0 16rpx;
  font-size: 28rpx;
  background-color: #fff;
  border-radius: 8rpx;
  box-sizing: border-box;
}

.cost-picker {
  height: 72rpx;
  padding: 0 16rpx;
  font-size: 28rpx;
  background-color: #fff;
  border-radius: 8rpx;
  line-height: 72rpx;
  box-sizing: border-box;
}

.cost-total {
  text-align: right;
  font-size: 28rpx;
  color: #1890ff;
  margin-bottom: 16rpx;
}

.cost-error {
  font-size: 24rpx;
  color: #f56c6c;
  margin-bottom: 16rpx;
}

.add-cost-btn {
  width: 100%;
  height: 80rpx;
  line-height: 80rpx;
  background-color: #1890ff;
  color: #fff;
  font-size: 28rpx;
  border-radius: 12rpx;
  border: none;
}

/* 图片上传 */
.photo-upload {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.photo-item {
  position: relative;
  width: 200rpx;
  height: 200rpx;
}

.photo-img {
  width: 100%;
  height: 100%;
  border-radius: 12rpx;
}

.photo-remove {
  position: absolute;
  top: -12rpx;
  right: -12rpx;
  width: 40rpx;
  height: 40rpx;
  background-color: #f56c6c;
  color: #fff;
  border-radius: 50%;
  font-size: 28rpx;
  text-align: center;
  line-height: 40rpx;
}

.photo-add {
  width: 200rpx;
  height: 200rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: #f5f5f5;
  border: 2rpx dashed #ddd;
  border-radius: 12rpx;
}

.photo-add-icon {
  font-size: 60rpx;
  color: #ccc;
  line-height: 1;
}

.photo-add-text {
  font-size: 24rpx;
  color: #999;
  margin-top: 8rpx;
}

/* 提交按钮 */
.submit-area {
  padding: 20rpx 0 40rpx;
}

.submit-btn {
  width: 100%;
  height: 88rpx;
  background-color: #1890ff;
  color: #fff;
  font-size: 32rpx;
  font-weight: 600;
  border-radius: 44rpx;
  border: none;
}

.submit-btn[disabled] {
  background-color: #a0cfff;
}

.loading-view {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 300rpx;
  color: #999;
}
</style>
