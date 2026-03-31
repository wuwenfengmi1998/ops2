import { watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

/**
 * 自动设置页面标题，响应语言变化
 *
 * @param {string} i18nKey - i18n 翻译键，如 'appname.home'
 * @param {string} prefix  - 标题前缀，默认 'Operations.'
 *
 * @example
 *   usePageTitle('appname.home')
 *   // → document.title = "Operations.Home"（英文）
 *   // → document.title = "Operations.主页"（中文）
 */
export function usePageTitle(i18nKey, prefix = 'Operations.') {
  const { t, locale } = useI18n()

  function update() {
    document.title = prefix + t(i18nKey)
  }

  onMounted(update)
  watch(locale, update)
}
