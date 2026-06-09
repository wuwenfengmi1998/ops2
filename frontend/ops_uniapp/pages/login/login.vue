<template>
	<view class="container">
		<!-- Logo 区域 -->
		<view class="logo-section">
			<image class="logo-img" src="/static/logo.png" mode="aspectFit" />
			<text class="app-name">OPS 管理系统</text>
		</view>

		<!-- 登录表单 -->
		<view class="form">
			<input
				class="input"
				type="text"
				v-model="username"
				placeholder="请输入用户名"
				:disabled="loading"
			/>
			<input
				class="input"
				type="password"
				v-model="password"
				placeholder="请输入密码"
				:disabled="loading"
				@confirm="handleLogin"
			/>

			<!-- 记住登录 -->
			<view class="remember-row">
				<checkbox-group @change="onRememberChange">
					<label class="remember-label">
						<checkbox :checked="remember" value="1" color="#007AFF" />
						<text>记住登录</text>
					</label>
				</checkbox-group>
			</view>

			<button
				class="submit-btn"
				@click="handleLogin"
				:disabled="loading"
			>
				{{ loading ? '登录中...' : '登录' }}
			</button>
		</view>

		<!-- 错误提示 -->
		<view v-if="errorMsg" class="error-tip">
			{{ errorMsg }}
		</view>
	</view>
	<my-toast ref="toast" />
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '../../stores/user'
import { useConfigStore } from '../../stores/config'
import { userApi } from '../../api/user'

const userStore = useUserStore()
const configStore = useConfigStore()

const username = ref('')
const password = ref('')
const remember = ref(true)
const loading = ref(false)
const errorMsg = ref('')
const toast = ref(null)

// 页面加载时检查是否已配置 API 地址
onMounted(() => {
  const baseUrl = configStore.getApiBaseUrl()
  if (!baseUrl) {
    uni.showModal({
      title: '提示',
      content: '请先在设置页面配置 API 地址',
      showCancel: false,
      success: () => {
        uni.navigateTo({ url: '/pages/settings/settings' })
      }
    })
  }
})

// 记住登录切换
const onRememberChange = (e) => {
  remember.value = e.detail.value.includes('1')
}

// 处理登录
const handleLogin = async () => {
  // 验证输入
  if (!username.value.trim()) {
    errorMsg.value = '请输入用户名'
    return
  }
  if (!password.value) {
    errorMsg.value = '请输入密码'
    return
  }

  // 检查 API 地址
  const baseUrl = configStore.getApiBaseUrl()
  if (!baseUrl) {
    uni.showModal({
      title: '提示',
      content: '请先在设置页面配置 API 地址',
      showCancel: false,
      success: () => {
        uni.navigateTo({ url: '/pages/settings/settings' })
      }
    })
    return
  }

  errorMsg.value = ''
  loading.value = true

  try {
    const res = await userApi.login(username.value, password.value, remember.value)

    if (res.errCode === 0 && res.data?.cookie) {
      // 登录成功
      toast.value.success('登录成功')
      userStore.login(res.data.cookie)

      // 延迟跳转，等待 Toast 显示
      setTimeout(() => {
        // 跳转到用户中心（TabBar 页面）
        uni.switchTab({
          url: '/pages/user/user'
        })
      }, 500)
    } else {
      // 登录失败，显示统一错误信息（不区分具体错误，防止暴力破解）
      errorMsg.value = '用户名或密码错误'
      toast.value.error(errorMsg.value)
    }
  } catch (e) {
    console.error('登录异常', e)
    errorMsg.value = '网络错误，请检查网络连接'
    toast.value.error('网络错误')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.container {
	flex: 1;
	display: flex;
	flex-direction: column;
	align-items: center;
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	min-height: 100vh;
	padding: 100rpx 40rpx 40rpx;
}

/* Logo 区域 */
.logo-section {
	display: flex;
	flex-direction: column;
	align-items: center;
	margin-bottom: 80rpx;
}

.logo-img {
	width: 160rpx;
	height: 160rpx;
	margin-bottom: 20rpx;
}

.app-name {
	font-size: 40rpx;
	font-weight: bold;
	color: #FFFFFF;
	letter-spacing: 4rpx;
}

/* 表单区域 */
.form {
	width: 100%;
	background-color: #FFFFFF;
	border-radius: 20rpx;
	padding: 40rpx;
	box-shadow: 0 10rpx 40rpx rgba(0, 0, 0, 0.1);
}

.input {
	height: 90rpx;
	line-height: 90rpx;
	padding: 0 30rpx;
	margin-bottom: 30rpx;
	background-color: #F8F8F8;
	border-radius: 10rpx;
	font-size: 28rpx;
}

.input:last-of-type {
	margin-bottom: 20rpx;
}

/* 记住登录 */
.remember-row {
	margin-bottom: 30rpx;
}

.remember-label {
	display: flex;
	align-items: center;
	font-size: 26rpx;
	color: #666;
}

.remember-label checkbox {
	margin-right: 10rpx;
}

/* 登录按钮 */
.submit-btn {
	width: 100%;
	height: 90rpx;
	line-height: 90rpx;
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	color: #FFFFFF;
	font-size: 32rpx;
	border-radius: 45rpx;
	border: none;
}

.submit-btn[disabled] {
	background: #CCCCCC;
}

/* 错误提示 */
.error-tip {
	width: 100%;
	margin-top: 30rpx;
	padding: 20rpx;
	background-color: rgba(255, 255, 255, 0.9);
	border-radius: 10rpx;
	text-align: center;
	color: #FF3B30;
	font-size: 26rpx;
}
</style>
