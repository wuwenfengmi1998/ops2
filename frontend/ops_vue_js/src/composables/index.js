import { ref } from 'vue'

/**
 * 表单验证工具
 *
 * @example
 *   const { validate, errors, clearErrors } = useValidation()
 *
 *   const form = reactive({ name: '', email: '' })
 *
 *   function handleSubmit() {
 *     clearErrors()
 *     validate('name', form.name, '请输入名称')
 *     validate('email', form.email, '请输入邮箱', v => isValidEmail(v) || '邮箱格式不正确')
 *     if (!Object.keys(errors).length) { ... }
 *   }
 */
export function useValidation() {
  const errors = ref({})

  function clearErrors() {
    errors.value = {}
  }

  function clearError(field) {
    delete errors.value[field]
    errors.value = { ...errors.value }
  }

  /**
   * @param {string} field - 字段名
   * @param {*} value - 值
   * @param {string} emptyMsg - 空值提示
   * @param {function|undefined} extraCheck - 额外校验函数，返回 false 或错误消息
   */
  function validate(field, value, emptyMsg, extraCheck) {
    clearError(field)

    if (!value || (typeof value === 'string' && !value.trim())) {
      errors.value[field] = emptyMsg
      return false
    }

    if (extraCheck) {
      const result = extraCheck(value)
      if (result === false) {
        errors.value[field] = emptyMsg
        return false
      }
      if (typeof result === 'string') {
        errors.value[field] = result
        return false
      }
    }

    return true
  }

  function hasErrors() {
    return Object.keys(errors.value).length > 0
  }

  return { errors, clearErrors, clearError, validate, hasErrors }
}

/** 邮箱格式校验 */
export function isValidEmail(email) {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
}
