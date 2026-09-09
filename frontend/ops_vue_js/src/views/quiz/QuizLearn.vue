<script setup>
import { ref, computed, nextTick, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { quizApi } from '@/api/quiz'
import { IconChevronLeft, IconCircleCheck, IconBook2 } from '@tabler/icons-vue'

usePageTitle('appname.quiz_bank')
const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const toast = useToastStore()

const bankId = parseInt(route.query.bank) || 0
const bankName = ref('')
const items = ref([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)
const loadingMore = ref(false)
const cardEls = ref([])
let sentinel = null
let sentinelObserver = null

const posKey = 'quiz_learn_pos_' + bankId

const typeLabels = {
  single: 'quiz.type_single',
  multiple: 'quiz.type_multiple',
  blank: 'quiz.type_blank',
  judge: 'quiz.type_judge',
}

const typeColors = {
  single: 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-400',
  multiple: 'bg-purple-100 text-purple-700 dark:bg-purple-900/40 dark:text-purple-400',
  blank: 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-400',
  judge: 'bg-teal-100 text-teal-700 dark:bg-teal-900/40 dark:text-teal-400',
}

const hasMore = computed(() => items.value.length < total.value)

function avatarUrl(bankItem) {
  return bankItem.creatorAvatar ? '/api/static/avatar/' + bankItem.creatorAvatar : '/ava.svg'
}

// 正确答案的展示文本
function correctText(item) {
  if (item.type === 'blank') {
    const arr = Array.isArray(item.answer) ? item.answer : []
    return arr.map(a => String(a).split('|')[0]).join(' / ') || '-'
  }
  if (item.type === 'judge') {
    return String(item.answer) === 'true' ? t('quiz.true') : t('quiz.false')
  }
  if (item.type === 'single') {
    const idx = typeof item.answer === 'number' ? item.answer : parseInt(item.answer)
    if (Number.isFinite(idx) && item.options[idx]) return item.options[idx]
    return '-'
  }
  // multiple
  const idxs = Array.isArray(item.answer) ? item.answer : []
  const txt = idxs.map(i => item.options[i]).filter(Boolean)
  return txt.join(', ') || '-'
}

function isCorrectOption(item, idx) {
  if (item.type === 'multiple') {
    return Array.isArray(item.answer) && item.answer.map(Number).includes(idx)
  }
  if (item.type === 'single') {
    const a = typeof item.answer === 'number' ? item.answer : parseInt(item.answer)
    return a === idx
  }
  return false
}

async function loadFirst() {
  loading.value = true
  try {
    const { errCode, data } = await quizApi.learnBank(bankId, 1, 10)
    if (errCode === 0) {
      bankName.value = data.bankName ?? ''
      items.value = data.items ?? []
      total.value = data.total ?? 0
      page.value = 1
    } else if (errCode === -77) {
      toast.error(t('quiz.bank_unavailable'))
      router.replace('/questions')
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  } finally {
    loading.value = false
    await nextTick()
    sentinel = document.getElementById('learn-sentinel')
    setupSentinel()
  }
}

async function loadMore() {
  if (!hasMore.value || loadingMore.value || loading.value) return
  loadingMore.value = true
  try {
    const next = page.value + 1
    const { errCode, data } = await quizApi.learnBank(bankId, next, 10)
    if (errCode === 0) {
      items.value = (items.value ?? []).concat(data.items ?? [])
      total.value = data.total ?? total.value
      page.value = next
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  } finally {
    loadingMore.value = false
  }
}

function setupSentinel() {
  if (sentinelObserver) sentinelObserver.disconnect()
  if (!sentinel) return
  sentinelObserver = new IntersectionObserver((entries) => {
    if (entries.some(e => e.isIntersecting)) {
      loadMore()
    }
  }, { rootMargin: '200px 0px' })
  sentinelObserver.observe(sentinel)
}

function setCardRef(el, index) {
  if (el) cardEls.value[index] = el
}

// 视口上半个区域内最靠下的题卡序号
function topVisibleIndex() {
  let idx = 0
  cardEls.value.forEach((el, i) => {
    if (!el) return
    const r = el.getBoundingClientRect()
    if (r.top < window.innerHeight * 0.5) idx = i
  })
  return idx
}

// 保存阅读位置（题卡序号 + 滚动偏移）并返回
function saveAndBack() {
  const pos = { itemIndex: topVisibleIndex(), scrollY: window.scrollY }
  try {
    localStorage.setItem(posKey, JSON.stringify(pos))
  } catch {
    // 忽略存储失败，仍返回
  }
  router.push('/questions')
}

// 恢复阅读位置：按需加载至目标题所在批后定位滚动
async function restorePosition() {
  let pos = null
  try {
    pos = JSON.parse(localStorage.getItem(posKey))
  } catch {
    pos = null
  }
  if (!pos || typeof pos.itemIndex !== 'number' || typeof pos.scrollY !== 'number') return
  const targetIdx = Math.min(pos.itemIndex, total.value - 1)
  if (targetIdx < 0) return
  let guard = 0
  while (items.value.length <= targetIdx && hasMore.value && guard < 50) {
    await loadMore()
    guard++
  }
  await nextTick()
  const el = cardEls.value[targetIdx]
  if (el) {
    window.scrollTo(0, pos.scrollY)
  }
}

onMounted(async () => {
  if (!bankId) {
    router.replace('/questions')
    return
  }
  await loadFirst()
  restorePosition()
})

onUnmounted(() => {
  if (sentinelObserver) sentinelObserver.disconnect()
})
</script>

<template>
  <div class="mx-auto max-w-2xl px-6 py-6 lg:max-w-5xl">
    <div class="mb-6 flex items-center gap-2">
      <IconBook2 :size="18" class="text-blue-500" />
      <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ bankName || '...' }}</h3>
      <span class="rounded-full bg-blue-100 px-2.5 py-0.5 text-xs font-semibold text-blue-700 dark:bg-blue-900/40 dark:text-blue-400">
        {{ t('quiz.learn_title') }}
      </span>
      <span class="ml-auto text-sm text-gray-500 dark:text-gray-400">{{ items.length }} / {{ total }}</span>
      <button
        class="ml-2 inline-flex items-center gap-1 text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300 lg:hidden"
        @click="saveAndBack"
      >
        <IconChevronLeft :size="16" />
        {{ t('quiz.back_home') }}
      </button>
    </div>

    <div v-if="loading" class="py-20 text-center text-gray-400">
      <svg class="mx-auto mb-2 h-6 w-6 animate-spin text-gray-400" viewBox="0 0 24 24" fill="none">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
      </svg>
      Loading...
    </div>

    <template v-else>
      <div class="flex flex-col gap-6 lg:flex-row">
        <!-- 左栏：题目 -->
        <main class="min-w-0 flex-1">
          <div v-if="items.length === 0" class="py-20 text-center text-gray-400">{{ t('quiz.learn_empty') }}</div>

          <!-- 单列 -->
          <div class="columns-1">
            <div
              v-for="(item, i) in items"
              :key="item.id"
              :ref="el => setCardRef(el, i)"
              class="mb-4 break-inside-avoid rounded-xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dk-muted dark:bg-dk-card"
            >
              <div class="mb-3 flex items-center gap-3">
                <span class="flex h-7 min-w-7 items-center justify-center rounded-full bg-blue-100 px-1 text-sm font-bold text-blue-700 dark:bg-blue-900/40 dark:text-blue-400">{{ i + 1 }}</span>
                <span class="rounded-full px-2.5 py-0.5 text-xs font-semibold" :class="typeColors[item.type]">{{ t(typeLabels[item.type]) }}</span>
                <span class="ml-auto text-xs font-medium text-gray-400 dark:text-gray-500">{{ item.score }} {{ t('quiz.points') }}</span>
              </div>

              <p class="mb-4 whitespace-pre-wrap text-sm font-medium text-gray-900 dark:text-white">{{ item.title }}</p>

              <!-- 选项（正确高亮） -->
              <div v-if="item.options.length > 0" class="mb-4 flex flex-col gap-2">
                <div
                  v-for="(opt, oi) in item.options"
                  :key="oi"
                  class="flex items-center gap-2.5 rounded-lg border px-3.5 py-2 text-sm"
                  :class="isCorrectOption(item, oi)
                    ? 'border-green-300 bg-green-50 text-green-700 dark:border-green-700 dark:bg-green-900/30 dark:text-green-300'
                    : 'border-gray-200 text-gray-600 dark:border-dk-muted dark:text-gray-300'"
                >
                  <span class="shrink-0 font-semibold">{{ String.fromCharCode(65 + oi) }}.</span>
                  <span class="min-w-0 flex-1 break-words">{{ opt }}</span>
                  <IconCircleCheck v-if="isCorrectOption(item, oi)" :size="16" class="shrink-0 text-green-500" />
                </div>
              </div>

              <!-- 答案 -->
              <div class="mb-2 text-sm">
                <span class="text-gray-400 dark:text-gray-500">{{ t('quiz.correct_answer') }}: </span>
                <template v-if="item.type === 'judge'">
                  <span
                    class="rounded-full px-2.5 py-0.5 text-xs font-semibold"
                    :class="String(item.answer) === 'true'
                      ? 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400'
                      : 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400'"
                  >{{ correctText(item) }}</span>
                </template>
                <span v-else class="font-medium text-green-600 dark:text-green-400">{{ correctText(item) }}</span>
              </div>

              <!-- 解析 -->
              <div
                v-if="item.explain"
                class="mb-3 rounded-lg border border-blue-100 bg-blue-50/60 px-4 py-2.5 text-sm leading-relaxed text-gray-700 dark:border-blue-900/40 dark:bg-blue-900/10 dark:text-gray-300"
              >
                <span class="font-medium text-blue-600 dark:text-blue-400">{{ t('quiz.explain') }}: </span>
                <span class="whitespace-pre-wrap">{{ item.explain }}</span>
              </div>

              <!-- 作者(录入人) -->
              <div class="flex items-center justify-end gap-1.5 text-xs text-gray-500 dark:text-gray-400">
                <img :src="avatarUrl(item)" class="h-5 w-5 rounded-full object-cover" :alt="item.creatorName || ''" />
                <span class="font-medium">{{ item.creatorName || '-' }}</span>
              </div>
            </div>
          </div>

          <!-- 加载更多 -->
          <div v-if="hasMore" class="flex h-10 items-center justify-center text-sm text-gray-400">
            <span v-if="loadingMore">{{ t('quiz.loading_more') }}</span>
          </div>
          <div v-else class="flex justify-center py-6 text-sm text-gray-400">{{ t('quiz.all_loaded') }}</div>
          <div id="learn-sentinel" class="h-px"></div>
        </main>

        <!-- 右栏：返回 -->
        <aside class="hidden lg:block lg:w-24 lg:shrink-0">
          <div class="lg:sticky lg:top-6">
            <button
              class="inline-flex items-center gap-1 rounded-lg border border-gray-300 px-3 py-1.5 text-sm text-gray-600 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-base dark:text-gray-300 dark:hover:bg-dk-muted"
              @click="saveAndBack"
            >
              <IconChevronLeft :size="16" />
              {{ t('quiz.back_home') }}
            </button>
          </div>
        </aside>
      </div>
    </template>
  </div>
</template>
