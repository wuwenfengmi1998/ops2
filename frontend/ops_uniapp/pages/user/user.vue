<template>
	<view class="container">
		<!-- 顶部设置按钮 -->
		<view class="header">
			<text class="settings-icon" @click="goToSettings">⚙️</text>
		</view>

		<!-- 未登录状态 -->
		<view v-if="!userStore.isLoggedIn" class="not-login">
			<view class="avatar-placeholder">
				<text class="avatar-icon">👤</text>
			</view>
			<text class="not-login-text">您还未登录</text>
			<button class="login-btn" @click="goToLogin">前往登录</button>
		</view>

		<!-- 已登录状态 -->
		<view v-else class="user-info">
			<!-- 用户信息卡片 -->
			<view class="user-card">
				<view class="avatar-section">
					<image
						v-if="userStore.avatarUrl"
						class="avatar-img"
						:src="userStore.avatarUrl"
						mode="aspectFill"
					/>
					<view v-else class="avatar-placeholder">
						<text class="avatar-icon">👤</text>
					</view>
				</view>
				<view class="user-detail">
					<text class="username">{{ userStore.username || '未设置昵称' }}</text>
					<text class="user-role">{{ getUserRole() }}</text>
				</view>
			</view>

<!-- 功能菜单 -->
		<view class="menu-section">
			<view class="menu-item" @click="switchToTab('/pages/order/order')">
				<text class="menu-icon">📦</text>
				<text class="menu-text">我的订单</text>
				<text class="menu-arrow">›</text>
			</view>
			<view class="menu-item" @click="switchToTab('/pages/warehouse/warehouse')">
				<text class="menu-icon">🏭</text>
				<text class="menu-text">仓库管理</text>
				<text class="menu-arrow">›</text>
			</view>
			<view class="menu-item" @click="switchToTab('/pages/schedule/list')">
				<text class="menu-icon">📅</text>
				<text class="menu-text">我的日程</text>
				<text class="menu-arrow">›</text>
			</view>
		</view>

			<!-- 退出登录 -->
			<view class="logout-section">
				<button class="logout-btn" @click="handleLogout">退出登录</button>
			</view>
		</view>
	</view>
</template>

<script setup>
import { onShow } from '@dcloudio/uni-app'
import { useUserStore } from '../../stores/user'

const userStore = useUserStore()

// 每次显示页面时检查登录状态
onShow(() => {
  // 如果未登录，尝试从存储恢复会话
  if (!userStore.isLoggedIn) {
    userStore.restoreSession()
  }
})

// 跳转到登录页
const goToLogin = () => {
	uni.navigateTo({
		url: '/pages/login/login'
	})
}

// 跳转到设置页
const goToSettings = () => {
	uni.navigateTo({
		url: '/pages/settings/settings'
	})
}

// 通用跳转（支持普通页面和 TabBar 页面）
const navigateTo = (url) => {
	// 检查是否是 TabBar 页面
	const tabBarPages = ['/pages/index/index', '/pages/order/order', '/pages/warehouse/warehouse', '/pages/user/user']
	if (tabBarPages.includes(url)) {
		uni.switchTab({ url })
	} else {
		uni.navigateTo({ url })
	}
}

// 跳转到 TabBar 页面
const switchToTab = (url) => {
	uni.switchTab({ url })
}

// 获取用户角色
const getUserRole = () => {
	if (!userStore.user) return ''
	// 根据用户组判断角色（这里简单处理）
	return '用户'
}

// 退出登录
const handleLogout = () => {
	uni.showModal({
		title: '确认退出',
		content: '确定要退出登录吗？',
		confirmColor: '#FF3B30',
		success: (res) => {
			if (res.confirm) {
				userStore.logout()
				uni.showToast({
					title: '已退出登录',
					icon: 'success'
				})
			}
		}
	})
}
</script>

<style scoped>
.container {
	flex: 1;
	display: flex;
	flex-direction: column;
	background-color: #F5F5F5;
	min-height: 100vh;
}

/* 顶部设置按钮 */
.header {
	width: 100%;
	display: flex;
	justify-content: flex-end;
	padding: 20rpx 30rpx;
}

.settings-icon {
	font-size: 48rpx;
}

/* 未登录状态 */
.not-login {
	flex: 1;
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	padding: 100rpx 40rpx;
}

.avatar-placeholder {
	width: 160rpx;
	height: 160rpx;
	border-radius: 80rpx;
	background-color: #E0E0E0;
	display: flex;
	align-items: center;
	justify-content: center;
	margin-bottom: 30rpx;
}

.avatar-icon {
	font-size: 80rpx;
}

.not-login-text {
	font-size: 32rpx;
	color: #666;
	margin-bottom: 40rpx;
}

.login-btn {
	width: 300rpx;
	height: 80rpx;
	line-height: 80rpx;
	background-color: #007AFF;
	color: #FFFFFF;
	font-size: 28rpx;
	border-radius: 40rpx;
}

/* 已登录状态 */
.user-info {
	flex: 1;
	padding: 0 30rpx;
}

/* 用户信息卡片 */
.user-card {
	display: flex;
	align-items: center;
	background-color: #FFFFFF;
	border-radius: 20rpx;
	padding: 40rpx;
	margin-bottom: 30rpx;
	box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.05);
}

.avatar-section {
	margin-right: 30rpx;
}

.avatar-img {
	width: 120rpx;
	height: 120rpx;
	border-radius: 60rpx;
}

.avatar-placeholder {
	width: 120rpx;
	height: 120rpx;
	border-radius: 60rpx;
	background-color: #E0E0E0;
	display: flex;
	align-items: center;
	justify-content: center;
}

.user-detail {
	display: flex;
	flex-direction: column;
}

.username {
	font-size: 36rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 10rpx;
}

.user-role {
	font-size: 26rpx;
	color: #999;
}

/* 功能菜单 */
.menu-section {
	background-color: #FFFFFF;
	border-radius: 20rpx;
	margin-bottom: 30rpx;
	overflow: hidden;
	box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.05);
}

.menu-item {
	display: flex;
	align-items: center;
	padding: 30rpx;
	border-bottom: 1rpx solid #F0F0F0;
}

.menu-item:last-child {
	border-bottom: none;
}

.menu-icon {
	font-size: 40rpx;
	margin-right: 20rpx;
}

.menu-text {
	flex: 1;
	font-size: 30rpx;
	color: #333;
}

.menu-arrow {
	font-size: 36rpx;
	color: #CCCCCC;
}

/* 退出登录 */
.logout-section {
	margin-top: 40rpx;
}

.logout-btn {
	width: 100%;
	height: 90rpx;
	line-height: 90rpx;
	background-color: #FFFFFF;
	color: #FF3B30;
	font-size: 30rpx;
	border-radius: 45rpx;
	border: none;
}
</style>
