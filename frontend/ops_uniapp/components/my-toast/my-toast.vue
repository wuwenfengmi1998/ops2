<template>
  <view v-if="visible" class="my-toast" :class="typeClass">
    <text>{{ message }}</text>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'

const visible = ref(false)
const message = ref('')
const type = ref('info') // info | success | error | warning

const typeClass = computed(() => `toast-${type.value}`)

// 暴露方法给外部调用
function show(msg, t = 'info', duration = 2000) {
  message.value = msg
  type.value = t
  visible.value = true

  setTimeout(() => {
    visible.value = false
  }, duration)
}

// 快捷方法
function success(msg, duration = 2000) { show(msg, 'success', duration) }
function error(msg, duration = 2000) { show(msg, 'error', duration) }
function warning(msg, duration = 2000) { show(msg, 'warning', duration) }
function info(msg, duration = 2000) { show(msg, 'info', duration) }

defineExpose({ show, success, error, warning, info })
</script>

<style scoped>
/* 固定顶部居中 */
.my-toast {
  position: fixed;
  top: 120rpx;
  left: 50%;
  transform: translateX(-50%);
  padding: 20rpx 40rpx;
  border-radius: 12rpx;
  font-size: 28rpx;
  z-index: 9999;
  text-align: center;
  max-width: 600rpx;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.15);
  word-break: break-all;
  line-height: 1.5;
}

/* 四种类型 */
.toast-info {
  background-color: #007AFF;
  color: #fff;
}

.toast-success {
  background-color: #34C759;
  color: #fff;
}

.toast-error {
  background-color: #FF3B30;
  color: #fff;
}

.toast-warning {
  background-color: #FF9500;
  color: #fff;
}
</style>
