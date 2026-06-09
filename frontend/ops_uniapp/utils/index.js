
/**
 * 检查字符串是否是合法的 URL
 * @param {string} str
 * @returns {boolean}
 */
export function isUrl(str) {
 if (!str || typeof str !== 'string') return false
   
   // 必须以 http:// 或 https:// 开头（接口地址必须带）
   const reg = /^https?:\/\/.+/i
   return reg.test(str.trim())
}