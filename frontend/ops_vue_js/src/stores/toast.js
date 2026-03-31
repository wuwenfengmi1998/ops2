import { defineStore } from 'pinia'
import { ref } from 'vue'

/**
 * 全局 Toast 通知 store
 * 用法：
 *   const toast = useToastStore()
 *   toast.success('保存成功')
 *   toast.error('网络错误')
 *   toast.warning('请注意')
 *   toast.info('提示信息')
 */
export const useToastStore = defineStore('toast', () => {
  const visible = ref(false)
  const type = ref('info')   // success | warning | danger | info
  const message = ref('')
  const dismissTimer = ref(null)

  function show(newType, newMessage, duration = 5000) {
    if (dismissTimer.value) {
      clearTimeout(dismissTimer.value)
    }

    type.value = newType
    message.value = newMessage
    visible.value = true

    if (duration > 0) {
      dismissTimer.value = setTimeout(() => {
        visible.value = false
      }, duration)
    }
  }

  function success(msg, duration) { show('success', msg, duration) }
  function warning(msg, duration) { show('warning', msg, duration) }
  function danger(msg, duration) { show('danger', msg, duration) }
  function error(msg, duration) { show('danger', msg, duration) }
  function info(msg, duration) { show('info', msg, duration) }

  function hide() {
    if (dismissTimer.value) {
      clearTimeout(dismissTimer.value)
    }
    visible.value = false
  }

  return { visible, type, message, show, success, warning, danger, error, info, hide }
})
