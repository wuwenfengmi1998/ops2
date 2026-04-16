<script>
/**
 * 用户设置页（个人信息）
 * 对标 PC 前端 src/views/settings/AccountView.vue
 */
import { userStore } from '../../store/user.js'
import { authApi } from '../../@api/auth.js'

export default {
  data() {
    return {
      isLoggedIn: false,

      // 表单
      form: {
        username: '',
        remark: '',
        birthday: '',
      },
      avatarUrl: '/static/ava.svg',

      // 状态
      loading: false,
      avatarUploading: false,
    }
  },

  methods: {
    // ── 初始化 ──

    initForm() {
      this.isLoggedIn   = userStore.isLoggedIn
      this.avatarUrl    = userStore.getAvatarUrl()
      const info        = userStore.userInfo
      if (info) {
        this.form.username = info.Username || ''
        this.form.remark   = info.FirstName || ''
        this.form.birthday = userStore.getBirthday()
      }
    },

    // ── 日期选择 ──

    onBirthdayChange(e) {
      this.form.birthday = e.detail.value
    },

    // ── 头像上传 ──

    chooseAvatar() {
      if (!this.isLoggedIn) {
        uni.showToast({ title: '请先登录', icon: 'none' })
        return
      }
      uni.chooseImage({
        count: 1,
        sourceType: ['album', 'camera'],
        success: async (res) => {
          const filePath = res.tempFilePaths[0]
          if (!filePath) return
          this.avatarUploading = true
          try {
            const { errCode } = await authApi.updateAvatar(filePath)
            if (errCode === 0) {
              uni.showToast({ title: '头像更新成功', icon: 'success' })
              // 刷新用户信息（含新头像路径）
              await userStore.fetchUserInfo()
              this.avatarUrl = userStore.getAvatarUrl()
            } else {
              uni.showToast({ title: '头像上传失败', icon: 'none' })
            }
          } catch {
            // request.js 已处理
          } finally {
            this.avatarUploading = false
          }
        },
      })
    },

    // ── 保存信息 ──

    async saveInfo() {
      if (!this.isLoggedIn) {
        uni.showToast({ title: '请先登录', icon: 'none' })
        return
      }
      if (!this.form.username.trim()) {
        uni.showToast({ title: '名字不能为空', icon: 'none' })
        return
      }
      this.loading = true
      try {
        const { errCode } = await authApi.updateInfo({
          username: this.form.username.trim(),
          remark:   this.form.remark.trim(),
          birthday: this.form.birthday,
        })
        if (errCode === 0) {
          uni.showToast({ title: '保存成功', icon: 'success' })
          await userStore.fetchUserInfo()
        } else {
          uni.showToast({ title: '保存失败', icon: 'none' })
        }
      } catch {
        // request.js 已处理
      } finally {
        this.loading = false
      }
    },

    // ── 登出 ──

    handleLogout() {
      uni.showModal({
        title: '确认登出',
        content: '确定要退出登录吗？',
        success: (res) => {
          if (res.confirm) {
            userStore.logout()
            uni.reLaunch({ url: '/pages/index/index' })
          }
        },
      })
    },

    goToLogin() {
      uni.navigateTo({ url: '/pages/signin' })
    },
  },

  onShow() {
    this.initForm()
  },

  onLoad() {
    this.initForm()
  },
}
</script>

<template>
  <view class="settings-page">

    <!-- 顶部导航 -->
    <view class="navbar">
      <view class="back-btn" @tap="() => uni.navigateBack()">
        <text class="back-icon">‹</text>
      </view>
      <text class="navbar-title">个人设置</text>
      <view style="width: 80rpx;" />
    </view>

    <scroll-view scroll-y class="content">

      <!-- ── 未登录提示 ── -->
      <view v-if="!isLoggedIn" class="not-login-card">
        <text class="icon-emoji" style="font-size: 80rpx;">👤</text>
        <text class="not-login-title">尚未登录</text>
        <text class="not-login-desc">登录后可查看和修改个人信息</text>
        <button class="primary-btn" @tap="goToLogin">立即登录</button>
      </view>

      <!-- ── 已登录内容 ── -->
      <view v-else>

        <!-- 头像区域 -->
        <view class="avatar-card">
          <view class="avatar-wrap" @tap="chooseAvatar">
            <image :src="avatarUrl" class="avatar" mode="aspectFill" />
            <view class="avatar-overlay">
              <text class="camera-icon">📷</text>
            </view>
            <view v-if="avatarUploading" class="uploading-mask">
              <text class="uploading-text">上传中…</text>
            </view>
          </view>
          <text class="avatar-hint">点击更换头像</text>
        </view>

        <!-- 信息表单 -->
        <view class="form-card">
          <text class="section-title">基本信息</text>

          <!-- 名字 -->
          <view class="field">
            <text class="label">名字</text>
            <input
              v-model="form.username"
              class="input"
              placeholder="输入你的名字"
              maxlength="30"
            />
          </view>

          <!-- 备注 -->
          <view class="field">
            <text class="label">备注</text>
            <input
              v-model="form.remark"
              class="input"
              placeholder="个人简介或备注"
              maxlength="50"
            />
          </view>

          <!-- 生日 -->
          <view class="field">
            <text class="label">生日</text>
            <picker
              mode="date"
              :value="form.birthday"
              @change="onBirthdayChange"
            >
              <view class="picker-display">
                <text :class="form.birthday ? 'picker-value' : 'picker-placeholder'">
                  {{ form.birthday || '选择生日' }}
                </text>
                <text class="picker-arrow">›</text>
              </view>
            </picker>
          </view>

          <!-- 保存按钮 -->
          <button
            class="primary-btn"
            :class="{ 'btn-loading': loading }"
            :disabled="loading"
            @tap="saveInfo"
          >
            {{ loading ? '保存中…' : '保存修改' }}
          </button>
        </view>

        <!-- 安全操作 -->
        <view class="danger-card">
          <button class="logout-btn" @tap="handleLogout">退出登录</button>
        </view>

      </view>
    </scroll-view>
  </view>
</template>

<style scoped>
/* ── 页面 ── */
.settings-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #f3f4f6;
}

/* ── 顶部导航 ── */
.navbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 88rpx 24rpx 24rpx;
  background: #fff;
  border-bottom: 1rpx solid #e5e7eb;
}
.back-btn {
  width: 80rpx;
  height: 80rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}
.back-icon {
  font-size: 56rpx;
  color: #374151;
  line-height: 1;
}
.navbar-title {
  font-size: 36rpx;
  font-weight: 600;
  color: #111827;
}

/* ── 内容区 ── */
.content {
  flex: 1;
  padding: 32rpx;
}

/* ── 未登录 ── */
.not-login-card {
  background: #fff;
  border-radius: 24rpx;
  padding: 80rpx 40rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20rpx;
  margin-top: 40rpx;
}
.not-login-title {
  font-size: 36rpx;
  font-weight: 600;
  color: #374151;
}
.not-login-desc {
  font-size: 28rpx;
  color: #9ca3af;
  text-align: center;
}

/* ── 头像区 ── */
.avatar-card {
  background: #fff;
  border-radius: 24rpx;
  padding: 48rpx 32rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20rpx;
  margin-bottom: 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.06);
}
.avatar-wrap {
  position: relative;
  width: 180rpx;
  height: 180rpx;
}
.avatar {
  width: 180rpx;
  height: 180rpx;
  border-radius: 50%;
  background: #e5e7eb;
}
.avatar-overlay {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 56rpx;
  height: 56rpx;
  background: #2563eb;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 4rpx solid #fff;
}
.camera-icon { font-size: 28rpx; }
.uploading-mask {
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0.5);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.uploading-text {
  font-size: 24rpx;
  color: #fff;
}
.avatar-hint {
  font-size: 26rpx;
  color: #9ca3af;
}

/* ── 表单卡片 ── */
.form-card {
  background: #fff;
  border-radius: 24rpx;
  padding: 36rpx;
  margin-bottom: 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.06);
}
.section-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #374151;
  display: block;
  margin-bottom: 32rpx;
}
.field {
  margin-bottom: 36rpx;
}
.label {
  font-size: 28rpx;
  font-weight: 500;
  color: #374151;
  display: block;
  margin-bottom: 12rpx;
}
.input {
  width: 100%;
  height: 88rpx;
  border: 2rpx solid #d1d5db;
  border-radius: 16rpx;
  padding: 0 28rpx;
  font-size: 30rpx;
  color: #111827;
  background: #fff;
  box-sizing: border-box;
}

/* ── Picker ── */
.picker-display {
  height: 88rpx;
  border: 2rpx solid #d1d5db;
  border-radius: 16rpx;
  padding: 0 28rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #fff;
}
.picker-value {
  font-size: 30rpx;
  color: #111827;
}
.picker-placeholder {
  font-size: 30rpx;
  color: #9ca3af;
}
.picker-arrow {
  font-size: 36rpx;
  color: #9ca3af;
}

/* ── 按钮 ── */
.primary-btn {
  width: 100%;
  height: 96rpx;
  background: #2563eb;
  color: #fff;
  border-radius: 16rpx;
  font-size: 32rpx;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  margin-top: 8rpx;
}
.primary-btn[disabled] {
  background: #93c5fd;
}

/* ── 危险区 ── */
.danger-card {
  background: #fff;
  border-radius: 24rpx;
  padding: 36rpx;
  margin-bottom: 40rpx;
  box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.06);
}
.logout-btn {
  width: 100%;
  height: 96rpx;
  background: #fee2e2;
  color: #dc2626;
  border-radius: 16rpx;
  font-size: 32rpx;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
}
</style>
