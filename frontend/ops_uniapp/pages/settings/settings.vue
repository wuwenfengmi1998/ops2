<template>
  <!-- 页面根容器，灰色背景 -->
  <view class="page">

    <!-- ===== 分组：账号设置 ===== -->
    <view class="group">
      <!-- 分组标题 -->
      <text class="group-title">base url</text>
      <!-- 列表容器（白色卡片） -->
      <view class="group-body">
        <!-- 列表项：点击跳转或弹窗编辑 -->
        <view class="cell" @tap="onEditApiUrl">
           <text class="cell-value">{{ useConfig.getApiBaseUrl() || '未设置' }}</text> 
        </view>
		<view class="btn-wrapper">
		  <button v-if="!isTesting" class="save-btn" @tap="onTest">测试连接</button>
		  <button v-if="isTesting" class="save-btn" disabled>正在测试...</button>
		</view>
		<view v-if="isTested&&!isTesting" :class="isTestok ? 'result success' : 'result error'">
		  {{ isTestok ? '✅ 连接成功' : '❌ 连接失败' }}
		</view>
		</view>
    </view>

    <!-- ===== 分组：打印机 ===== -->
    <view class="group">
      <text class="group-title">打印机</text>
      <view class="group-body">
        <view class="cell" @tap="goPrinterTest">
          <text class="cell-label">打印机测试</text>
          <text class="cell-arrow">›</text>
        </view>
      </view>
    </view>

    <!-- ===== 分组：关于 ===== -->
    <view class="group">
      <text class="group-title">关于</text>
      <view class="group-body">
        <view class="cell">
          <text class="cell-label">版本号</text>
          <text class="cell-value">{{ version }}</text>
        </view>
      </view>
    </view>

  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { isUrl } from '@/utils/index.js'
import { useConfigStore } from '../../stores/config'

const useConfig=useConfigStore()

// ---- 版本号 ----
const version = ref('')

// ---- 响应式状态 ----
const isTested=ref(false)//是否点了测试
const isTesting=ref(false)//是否正在测试
const isTestok=ref(false)//测试是否通过

// ---- 事件处理 ----

function onTest(){
	//console.log("test")
	isTested.value=true
	isTesting.value=true
	uni.request({
	  url: useConfig.getApiBaseUrl(), 
	  method: 'GET',
	timeout:1000,
	  success: (res) => {
	    //console.log('成功', res.data)
		if(res.data.return.isOpsApiRoot==true){
			isTestok.value=true
		}
		
	  },
	  fail: (err) => {
	    //console.log('失败', err)
		isTestok.value=false
	  },
	  complete() {
	      // 这个回调不管成功失败**一定**会执行
	      //console.log('请求完成')
		  isTesting.value=false
	    }
	})
	
}

// 点击 API URL → 弹出输入框（示例用 uni.showModal 简单实现）
function onEditApiUrl() {
  // 实际可跳转子页面或弹出自定义弹窗
  uni.showModal({
    title: 'API Base URL',
    editable: true,
    placeholderText: 'https://example.com',
    content: useConfig.getApiBaseUrl(),
    success(res) {
      if (res.confirm) {
		  console.log(res)
		  if (isUrl(res.content)){
			  useConfig.setApiBaseUrl(res.content)
			  console.log(res.content,":ok")
			
		  }else{
			  console.log(res.content,":not a url")
		  }
		
        
		
        
      }
    }
  })
}

// 跳转打印机测试页
function goPrinterTest() {
  uni.navigateTo({ url: '/pages/printer-test/printer-test' })
}


onMounted(() => {
  // 真机/打包环境：通过 plus.runtime 获取
  // #ifndef H5
  version.value = plus.runtime.version || ''
  // #endif

  // H5 端：从 manifest.json 读取
  // #ifdef H5
  fetch('/manifest.json')
    .then(r => r.text())
    .then(text => {
      // manifest.json 可能含 JS 注释，需先去除再解析
      const cleaned = text.replace(/\/\*[\s\S]*?\*\/|\/\/.*$/gm, '')
      const m = JSON.parse(cleaned)
      version.value = (m.versionName || '') + ' (' + (m.versionCode || '') + ')'
    })
    .catch(() => {
      version.value = ''
    })
  // #endif
})

</script>

<style scoped>
/* 页面底色 */
.page {
  min-height: 100vh;
  background-color: #f2f2f7; /* iOS 系统级灰色 */
  padding-top: 20rpx;
}

/* 分组间距 */
.group {
  margin-bottom: 32rpx;
}

/* 分组标题：小号灰色，左侧缩进 */
.group-title {
  font-size: 24rpx;
  color: #8e8e93;
  padding: 0 32rpx 12rpx;
}

/* 白色卡片容器 */
.group-body {
  background-color: #ffffff;
  border-radius: 12rpx;
  margin: 0 24rpx;
  overflow: hidden; /* 保证圆角裁剪子项 */
}

/* 每一行 cell */
.cell {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  padding: 0 32rpx;
  height: 96rpx;
  background-color: #ffffff;
}

/* 左侧标签 */
.cell-label {
  font-size: 30rpx;
  color: #1c1c1e;
}

/* 点击高亮（uni-app 用 hover-class 实现，这里仅演示结构） */
.cell:active {
  background-color: #f2f2f7;
}

/* 左侧标签 */
.cell-label {
  font-size: 30rpx;
  color: #1c1c1e;
}

/* 右侧容器（值 + 箭头） */
.cell-right {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 8rpx;
}

/* 右侧当前值 */
.cell-value {
  font-size: 28rpx;
  color: #8e8e93;
  max-width: 300rpx;          /* 防止超长文本撑破布局 */
  
  white-space: nowrap;
  text-overflow: ellipsis;
}

/* 右箭头 */
.cell-arrow {
  font-size: 36rpx;
  color: #c7c7cc;
  line-height: 1;
}

/* switch 开关右对齐 */
.cell-switch {
  transform: scale(0.85); /* 微调大小 */
}

/* 分割线：左侧缩进，不跨满 */
.divider {
  height: 1rpx;
  background-color: #e5e5ea;
  margin-left: 32rpx; /* 左缩进，右边贴边，iOS 风格 */
}

/* 危险操作行（退出登录） */
.cell-danger {
  justify-content: center; /* 文字居中 */
}

.cell-label-danger {
  font-size: 30rpx;
  color: #ff3b30; /* iOS 红色 */
}
</style>
