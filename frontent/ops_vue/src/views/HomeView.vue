<script setup lang="ts">
import { useI18n } from 'vue-i18n'

// 使用 vue-i18n 的 Composition API
const { t, locale } = useI18n()

// 切换语言的函数
const changeLanguage = (lang: string) => {
  locale.value = lang
  // 可选：将用户选择的语言保存到 localStorage
  //localStorage.setItem('userLanguage', lang)
}
</script>

<template>
  <div>
    <h3>{{ t('message.hello') }}</h3>
    <button>{{ t('button.submit') }}</button>
    <button @click="changeLanguage('en')">English</button>
    <button @click="changeLanguage('zh-CN')">中文</button>
  </div>
</template>
