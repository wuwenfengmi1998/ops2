<template>
	<view class="container">
		<!-- 欢迎区域 -->
		<view class="welcome-section">
			<view class="welcome-header">
				<text class="welcome-title">{{ welcomeText }}</text>
				<text class="search-icon" @click="goSearch">🔍</text>
			</view>
			<text class="welcome-date">{{ todayDisplay }}</text>
		</view>

		<!-- 今日日程 -->
		<view class="card">
			<view class="card-header">
				<text class="card-title">📅 今日日程</text>
				<text class="card-count" v-if="todaySchedules.length > 0">{{ todaySchedules.length }}个</text>
			</view>

			<!-- 加载状态 -->
			<view v-if="loadingSchedules" class="loading">
				<text>加载中...</text>
			</view>

			<!-- 日程列表 -->
			<view v-else-if="todaySchedules.length > 0" class="schedule-list">
				<view
					v-for="schedule in todaySchedules"
					:key="schedule.ID"
					class="schedule-item"
				>
					<view class="schedule-date" :style="{ backgroundColor: schedule.BgColor || '#007AFF' }">
						<text>{{ formatScheduleDate(schedule) }}</text>
					</view>
					<view class="schedule-content">
						<text class="schedule-title">{{ schedule.Title }}</text>
						<text class="schedule-user">创建人: {{ getCreatorName(schedule.UserID) }}</text>
					</view>
				</view>
			</view>

			<!-- 无日程 -->
			<view v-else class="empty">
				<text class="empty-text">今日暂无日程</text>
			</view>
		</view>

		<!-- 订单统计 -->
		<view class="card">
			<view class="card-header">
				<text class="card-title">📦 待处理订单</text>
			</view>

			<!-- 加载状态 -->
			<view v-if="loadingOrders" class="loading">
				<text>加载中...</text>
			</view>

			<!-- 订单数量 -->
			<view v-else class="order-stats">
				<view class="stat-item" @click="switchToTab('/pages/order/order')">
					<text class="stat-num" :class="{ 'stat-warning': pendingCount > 0 }">{{ pendingCount || '—' }}</text>
					<text class="stat-label">待处理</text>
				</view>
				<view class="stat-item" @click="switchToTab('/pages/order/order')">
					<text class="stat-num">{{ arrivedCount || '—' }}</text>
					<text class="stat-label">已到达</text>
				</view>
				<view class="stat-item" @click="switchToTab('/pages/order/order')">
					<text class="stat-num">{{ receivedCount || '—' }}</text>
					<text class="stat-label">已收件</text>
				</view>
			</view>
		</view>

		<!-- 功能入口 -->
		<view class="card">
			<view class="card-header">
				<text class="card-title">🚀 快捷入口</text>
			</view>
			<view class="quick-links">
				<view class="quick-item" @click="switchToTab('/pages/order/order')">
					<text class="quick-icon">📦</text>
					<text class="quick-text">订单管理</text>
				</view>
				<view class="quick-item" @click="switchToTab('/pages/warehouse/warehouse')">
					<text class="quick-icon">🏭</text>
					<text class="quick-text">仓库管理</text>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
import { onShow, onHide } from '@dcloudio/uni-app'
import { ref, computed } from 'vue'
import { useUserStore } from '../../stores/user'
import { scheduleApi } from '../../api/schedule'
import { purchaseApi } from '../../api/purchase'
import { fetchUserInfo, getUsername } from '../../stores/users'

const userStore = useUserStore()

// 今日日期
const today = new Date()
const todayStr = today.toISOString().split('T')[0]

// 欢迎语
const welcomeText = computed(() => {
	if (userStore.isLoggedIn) {
		return `${userStore.username || '用户'}，您好！`
	}
	return '欢迎使用 OPS 系统'
})

// 今日日期显示
const todayDisplay = computed(() => {
	const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
	return `${today.getMonth() + 1}月${today.getDate()}日 ${weekdays[today.getDay()]}`
})

// 日程数据
const todaySchedules = ref([])
const loadingSchedules = ref(false)

// 订单数据
const pendingCount = ref(0)
const arrivedCount = ref(0)
const receivedCount = ref(0)
const loadingOrders = ref(false)

// 获取今日日程
async function fetchTodaySchedules(silent = false) {
	if (!silent) loadingSchedules.value = true
	try {
		const res = await scheduleApi.getEvents({
			start: todayStr,
			end: todayStr
		})
		if (res.errCode === 0 && res.data?.list) {
			todaySchedules.value = res.data.list
			// 预加载所有创建人的用户名
			const userIDs = [...new Set(res.data.list.map(s => s.UserID).filter(Boolean))]
			userIDs.forEach(userID => fetchUserInfo(userID))
		}
	} catch (e) {
		console.error('获取今日日程失败', e)
	} finally {
		if (!silent) loadingSchedules.value = false
	}
}

// 获取订单统计
async function fetchOrderStats(silent = false) {
	if (!silent) loadingOrders.value = true
	try {
		const res = await purchaseApi.getOrderCount()
		if (res.errCode === 0 && res.data) {
			pendingCount.value = res.data.pending || 0
			arrivedCount.value = res.data.arrived || 0
			receivedCount.value = res.data.received || 0
		}
	} catch (e) {
		console.error('获取订单统计失败', e)
	} finally {
		if (!silent) loadingOrders.value = false
	}
}

// 格式化日程日期
function formatScheduleDate(schedule) {
	if (schedule.StartDate !== schedule.EndDate) {
		return `${schedule.StartDate}\n至\n${schedule.EndDate}`
	}
	return schedule.StartDate
}

// 获取创建人名称（使用用户缓存）
function getCreatorName(userId) {
	if (!userId) return '未知'
	// 如果是当前用户，显示"我"
	if (userStore.userInfo && userStore.userInfo.ID === userId) {
		return userStore.username || '我'
	}
	// 从缓存获取用户名
	return getUsername(userId)
}

// 跳转到 TabBar 页面
function switchToTab(url) {
	uni.switchTab({ url })
}

function goSearch() {
	uni.navigateTo({ url: '/pages/search/search' })
}

// 定时刷新
let refreshTimer = null
const REFRESH_INTERVAL = 5000 // 5秒

// 开始定时刷新（静默刷新，不显示 loading）
function startRefreshTimer() {
	stopRefreshTimer()
	refreshTimer = setInterval(() => {
		fetchTodaySchedules(true) // true = 静默刷新
		if (userStore.isLoggedIn) {
			fetchOrderStats(true) // true = 静默刷新
		}
	}, REFRESH_INTERVAL)
}

// 停止定时刷新
function stopRefreshTimer() {
	if (refreshTimer) {
		clearInterval(refreshTimer)
		refreshTimer = null
	}
}

// 页面显示时加载数据
onShow(() => {
	// 恢复会话（如果需要）
	if (!userStore.isLoggedIn) {
		userStore.restoreSession()
	}

	// 加载数据
	fetchTodaySchedules()
	if (userStore.isLoggedIn) {
		fetchOrderStats()
	}

	// 开始定时刷新
	startRefreshTimer()
})

// 页面离开时停止定时器
onHide(() => {
	stopRefreshTimer()
})
</script>

<style scoped>
.container {
	flex: 1;
	background-color: #F5F5F5;
	min-height: 100vh;
	padding: 20rpx 30rpx;
}

/* 欢迎区域 */
.welcome-section {
	padding: 30rpx 0;
}

.welcome-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 10rpx;
}

.welcome-title {
	font-size: 44rpx;
	font-weight: bold;
	color: #333;
}

.search-icon {
	font-size: 40rpx;
	padding: 10rpx 20rpx;
	background-color: #f0f0f0;
	border-radius: 50%;
}

.welcome-date {
	font-size: 28rpx;
	color: #999;
}

/* 卡片 */
.card {
	background-color: #FFFFFF;
	border-radius: 20rpx;
	padding: 30rpx;
	margin-bottom: 24rpx;
	box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.05);
}

.card-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 24rpx;
}

.card-title {
	font-size: 32rpx;
	font-weight: bold;
	color: #333;
}

.card-count {
	font-size: 26rpx;
	color: #007AFF;
}

/* 加载状态 */
.loading {
	display: flex;
	justify-content: center;
	padding: 40rpx 0;
}

.loading text {
	font-size: 28rpx;
	color: #999;
}

/* 空状态 */
.empty {
	display: flex;
	justify-content: center;
	padding: 60rpx 0;
}

.empty-text {
	font-size: 28rpx;
	color: #999;
}

/* 日程列表 */
.schedule-list {
	display: flex;
	flex-direction: column;
	gap: 20rpx;
}

.schedule-item {
	display: flex;
	align-items: flex-start;
	background-color: #F8F9FA;
	border-radius: 12rpx;
	padding: 20rpx;
}

.schedule-date {
	display: flex;
	align-items: center;
	justify-content: center;
	padding: 10rpx 16rpx;
	border-radius: 8rpx;
	color: #FFFFFF;
	font-size: 22rpx;
	text-align: center;
	white-space: pre-line;
	line-height: 1.4;
	min-width: 120rpx;
}

.schedule-content {
	flex: 1;
	margin-left: 20rpx;
}

.schedule-title {
	display: block;
	font-size: 30rpx;
	color: #333;
	font-weight: 500;
	margin-bottom: 8rpx;
}

.schedule-user {
	font-size: 24rpx;
	color: #999;
}

/* 订单统计 */
.order-stats {
	display: flex;
	justify-content: space-around;
}

.stat-item {
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 20rpx;
}

.stat-num {
	font-size: 48rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 10rpx;
}

.stat-num.stat-warning {
	color: #FF9500;
}

.stat-label {
	font-size: 26rpx;
	color: #999;
}

/* 快捷入口 */
.quick-links {
	display: flex;
	justify-content: space-around;
}

.quick-item {
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 30rpx 50rpx;
	background-color: #F8F9FA;
	border-radius: 16rpx;
}

.quick-icon {
	font-size: 60rpx;
	margin-bottom: 16rpx;
}

.quick-text {
	font-size: 28rpx;
	color: #333;
}
</style>
