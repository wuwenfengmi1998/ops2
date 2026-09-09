<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { quizApi } from '@/api/quiz'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import { IconCircleCheck, IconClock, IconBook2 } from '@tabler/icons-vue'

usePageTitle('appname.quiz')
const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const toast = useToastStore()

const loading = ref(true)
const sessionId = ref(0)
const bankName = ref('')
const durationLimit = ref(0)
const questions = ref([])
const answers = reactive({})
const seconds = ref(0)
const remaining = ref(0)
const submitting = ref(false)
const confirmShow = ref(false)
const timedOut = ref(false)
const activeIndex = ref(0)
const questionEls = ref([])
let timer = null

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

function isAnswered(q) {
  const ans = answers[q.id]
  if (q.type === 'single' || q.type === 'judge') {
    return ans !== undefined && ans !== null && ans !== ''
  }
  if (q.type === 'multiple') {
    return Array.isArray(ans) && ans.length > 0
  }
  if (q.type === 'blank') {
    return Array.isArray(ans) && ans.some(v => (v || '').trim() !== '')
  }
  return false
}

const answeredCount = computed(() => {
  let n = 0
  for (const q of questions.value) {
    if (isAnswered(q)) n++
  }
  return n
})

const progressPercent = computed(() => {
  if (questions.value.length === 0) return 0
  return Math.round((answeredCount.value / questions.value.length) * 100)
})

const hasTimeLimit = computed(() => durationLimit.value > 0)
const remainingText = computed(() => {
  const r = Math.max(0, remaining.value)
  const m = String(Math.floor(r / 60)).padStart(2, '0')
  const s = String(r % 60).padStart(2, '0')
  return `${m}:${s}`
})

const timeText = computed(() => {
  const m = String(Math.floor(seconds.value / 60)).padStart(2, '0')
  const s = String(seconds.value % 60).padStart(2, '0')
  return `${m}:${s}`
})

function initAnswers() {
  for (const q of questions.value) {
    if (q.type === 'single' || q.type === 'judge') {
      answers[q.id] = ''
    } else if (q.type === 'multiple') {
      answers[q.id] = []
    } else if (q.type === 'blank') {
      const n = q.blanks > 0 ? q.blanks : 1
      answers[q.id] = Array.from({ length: n }, () => '')
    }
  }
}

function setQuestionRef(el, index) {
  if (el) questionEls.value[index] = el
}

function activate(index) {
  activeIndex.value = index
}

function jumpTo(index) {
  activeIndex.value = index
  const el = questionEls.value[index]
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
}

function toggleOption(q, idx) {
  if (q.type === 'single') {
    answers[q.id] = idx
  } else {
    const cur = answers[q.id] || []
    const pos = cur.indexOf(idx)
    if (pos >= 0) cur.splice(pos, 1)
    else cur.push(idx)
  }
}

function chooseJudge(q, val) {
  answers[q.id] = val
}

async function doSubmit(force = false) {
  if (submitting.value) return
  submitting.value = true
  try {
    const answerList = questions.value.map(q => ({
      qid: q.id,
      answer: answers[q.id] ?? '',
    }))
    const { errCode, data } = await quizApi.submitQuiz(sessionId.value, seconds.value, answerList)
    if (errCode === 0) {
      router.replace(`/quiz/result/${sessionId.value}`)
    } else if (errCode === -75) {
      toast.warning(t('quiz.session_finished'))
      router.replace('/quiz')
    } else {
      toast.error(t('message.server_error'))
      if (force) router.replace('/quiz')
    }
  } catch {
    // 拦截器已处理
    if (force) router.replace('/quiz')
  } finally {
    submitting.value = false
  }
}

function handleSubmit() {
  confirmShow.value = false
  doSubmit(false)
}

function onTimeout() {
  if (timedOut.value) return
  timedOut.value = true
  if (timer) {
    clearInterval(timer)
    timer = null
  }
  toast.warning(t('quiz.timeout'))
  doSubmit(true)
}

onMounted(async () => {
  const bankId = parseInt(route.query.bank) || 0
  if (!bankId) {
    router.replace('/quiz')
    return
  }
  try {
    const { errCode, data } = await quizApi.startQuiz(bankId)
    if (errCode === 0) {
      sessionId.value = data.sessionId ?? 0
      bankName.value = data.bankName ?? ''
      durationLimit.value = data.durationLimit ?? 0
      questions.value = data.questions ?? []
      initAnswers()
      timer = setInterval(() => {
        seconds.value++
        if (hasTimeLimit.value) {
          remaining.value = durationLimit.value - seconds.value
          if (remaining.value <= 0) onTimeout()
        }
      }, 1000)
    } else if (errCode === -71) {
      toast.error(t('quiz.no_questions'))
      router.replace('/quiz')
    } else if (errCode === -77) {
      toast.error(t('quiz.bank_unavailable'))
      router.replace('/quiz')
    } else {
      toast.error(t('message.server_error'))
      router.replace('/quiz')
    }
  } catch {
    router.replace('/quiz')
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <div class="mx-auto max-w-5xl px-6 py-6">
    <!-- Loading -->
    <div v-if="loading" class="py-20 text-center text-gray-400">
      <svg class="mx-auto mb-2 h-6 w-6 animate-spin text-gray-400" viewBox="0 0 24 24" fill="none">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
      </svg>
      Loading...
    </div>

    <template v-else>
      <!-- 题库名 -->
      <div class="mb-4 inline-flex items-center gap-1.5 text-sm font-medium text-gray-700 dark:text-gray-300">
        <IconBook2 :size="16" class="text-blue-500" />
        {{ bankName }}
      </div>

      <div class="flex flex-col gap-6 lg:flex-row">
        <!-- 左栏：题目 -->
        <main class="min-w-0 flex-1">
          <div
            v-for="(q, i) in questions"
            :key="q.id"
            :ref="el => setQuestionRef(el, i)"
            class="mb-5 cursor-pointer scroll-mt-4 rounded-xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dk-muted dark:bg-dk-card"
            :class="activeIndex === i ? 'border-blue-400 ring-1 ring-blue-400/40 dark:border-blue-500' : ''"
            @click="activate(i)"
          >
            <div class="mb-4 flex items-center gap-3">
              <span class="flex h-7 min-w-7 items-center justify-center rounded-full bg-blue-100 px-1 text-sm font-bold text-blue-700 dark:bg-blue-900/40 dark:text-blue-400">{{ i + 1 }}</span>
              <span
                class="rounded-full px-2.5 py-0.5 text-xs font-semibold"
                :class="typeColors[q.type]"
              >{{ t(typeLabels[q.type]) }}</span>
              <span class="ml-auto text-xs font-medium text-gray-400 dark:text-gray-500">{{ q.score }} {{ t('quiz.points') }}</span>
            </div>
            <p class="mb-4 whitespace-pre-wrap text-sm font-medium text-gray-900 dark:text-white">{{ q.title }}</p>

            <!-- 单选 -->
            <div v-if="q.type === 'single'" class="flex flex-col gap-2">
              <button
                v-for="(opt, oi) in q.options"
                :key="oi"
                class="flex items-center gap-3 rounded-lg border px-4 py-2.5 text-left text-sm transition-colors"
                :class="answers[q.id] === oi
                  ? 'border-blue-500 bg-blue-50 text-blue-700 dark:border-blue-500 dark:bg-blue-900/30 dark:text-blue-300'
                  : 'border-gray-200 text-gray-700 hover:border-gray-300 hover:bg-gray-50 dark:border-dk-muted dark:text-gray-300 dark:hover:bg-dk-base'"
                @click="toggleOption(q, oi)"
              >
                <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded-full border"
                  :class="answers[q.id] === oi ? 'border-blue-500 bg-blue-500' : 'border-gray-300 dark:border-dk-muted'">
                  <span v-if="answers[q.id] === oi" class="h-2 w-2 rounded-full bg-white"></span>
                </span>
                <span>{{ opt }}</span>
              </button>
            </div>

            <!-- 多选 -->
            <div v-else-if="q.type === 'multiple'" class="flex flex-col gap-2">
              <button
                v-for="(opt, oi) in q.options"
                :key="oi"
                class="flex items-center gap-3 rounded-lg border px-4 py-2.5 text-left text-sm transition-colors"
                :class="(answers[q.id] || []).includes(oi)
                  ? 'border-purple-500 bg-purple-50 text-purple-700 dark:border-purple-500 dark:bg-purple-900/30 dark:text-purple-300'
                  : 'border-gray-200 text-gray-700 hover:border-gray-300 hover:bg-gray-50 dark:border-dk-muted dark:text-gray-300 dark:hover:bg-dk-base'"
                @click="toggleOption(q, oi)"
              >
                <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded border"
                  :class="(answers[q.id] || []).includes(oi) ? 'border-purple-500 bg-purple-500' : 'border-gray-300 dark:border-dk-muted'">
                  <span v-if="(answers[q.id] || []).includes(oi)" class="h-2 w-2 rounded-sm bg-white"></span>
                </span>
                <span>{{ opt }}</span>
              </button>
            </div>

            <!-- 填空 -->
            <div v-else-if="q.type === 'blank'" class="flex flex-col gap-3">
              <div v-for="(_, bi) in (answers[q.id] || [])" :key="bi" class="flex items-center gap-3">
                <span class="text-sm text-gray-400 dark:text-gray-500">{{ t('quiz.blank') }}{{ bi + 1 }}</span>
                <input
                  v-model="answers[q.id][bi]"
                  type="text"
                  :placeholder="t('quiz.blank_placeholder')"
                  class="flex-1 rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 placeholder-gray-400 outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white dark:placeholder-gray-500"
                />
              </div>
            </div>

            <!-- 判断 -->
            <div v-else class="flex gap-3">
              <button
                class="flex-1 rounded-lg border px-4 py-2.5 text-sm font-medium transition-colors"
                :class="answers[q.id] === 'true'
                  ? 'border-green-500 bg-green-50 text-green-700 dark:border-green-500 dark:bg-green-900/30 dark:text-green-300'
                  : 'border-gray-200 text-gray-700 hover:border-gray-300 hover:bg-gray-50 dark:border-dk-muted dark:text-gray-300 dark:hover:bg-dk-base'"
                @click="chooseJudge(q, 'true')"
              >
                <span class="inline-flex items-center gap-1.5">
                  <IconCircleCheck :size="16" />
                  {{ t('quiz.true') }}
                </span>
              </button>
              <button
                class="flex-1 rounded-lg border px-4 py-2.5 text-sm font-medium transition-colors"
                :class="answers[q.id] === 'false'
                  ? 'border-red-500 bg-red-50 text-red-700 dark:border-red-500 dark:bg-red-900/30 dark:text-red-300'
                  : 'border-gray-200 text-gray-700 hover:border-gray-300 hover:bg-gray-50 dark:border-dk-muted dark:text-gray-300 dark:hover:bg-dk-base'"
                @click="chooseJudge(q, 'false')"
              >
                {{ t('quiz.false') }}
              </button>
            </div>
          </div>

          <div v-if="questions.length === 0" class="py-16 text-center text-gray-400">{{ t('quiz.no_questions') }}</div>
        </main>

        <!-- 右栏：计时/进度/答题卡/提交 -->
        <aside class="lg:w-72 lg:shrink-0">
          <div class="flex flex-col gap-4 lg:sticky lg:top-4">
            <!-- 计时 -->
            <div class="rounded-xl border border-gray-200 bg-white p-5 shadow-sm dark:border-dk-muted dark:bg-dk-card">
              <div class="flex items-center justify-between">
                <span class="inline-flex items-center gap-1.5 text-xs font-medium text-gray-500 dark:text-gray-400">
                  <IconClock :size="14" />
                  {{ hasTimeLimit ? t('quiz.remaining') : t('quiz.elapsed') }}
                </span>
                <span
                  class="font-mono text-2xl font-bold tabular-nums"
                  :class="hasTimeLimit && remaining <= 60 ? 'text-red-500 dark:text-red-400' : 'text-gray-900 dark:text-white'"
                  :title="hasTimeLimit ? t('quiz.elapsed') + ': ' + timeText : ''"
                >
                  {{ hasTimeLimit ? remainingText : timeText }}
                </span>
              </div>
              <div v-if="hasTimeLimit" class="mt-2 text-right text-xs text-gray-400 dark:text-gray-500">
                {{ t('quiz.elapsed') }}: {{ timeText }}
              </div>
            </div>

            <!-- 进度 -->
            <div class="rounded-xl border border-gray-200 bg-white p-5 shadow-sm dark:border-dk-muted dark:bg-dk-card">
              <div class="mb-2 flex items-center justify-between text-sm">
                <span class="text-gray-500 dark:text-gray-400">{{ t('quiz.answered_q') }}</span>
                <span class="font-semibold text-gray-900 dark:text-white">{{ answeredCount }} / {{ questions.length }}</span>
              </div>
              <div class="h-2 w-full overflow-hidden rounded-full bg-gray-100 dark:bg-dk-base">
                <div class="h-full rounded-full bg-blue-600 transition-all" :style="{ width: progressPercent + '%' }"></div>
              </div>
            </div>

            <!-- 答题卡 -->
            <div class="rounded-xl border border-gray-200 bg-white p-5 shadow-sm dark:border-dk-muted dark:bg-dk-card">
              <h4 class="mb-3 text-sm font-semibold text-gray-900 dark:text-white">{{ t('quiz.card_title') }}</h4>
              <div class="grid grid-cols-8 gap-1.5">
                <button
                  v-for="(q, i) in questions"
                  :key="q.id"
                  class="h-8 rounded-md text-xs font-medium transition-colors"
                  :class="[
                    isAnswered(q)
                      ? 'bg-green-100 text-green-700 dark:bg-green-900/50 dark:text-green-300'
                      : 'bg-gray-100 text-gray-500 hover:bg-gray-200 dark:bg-dk-base dark:text-gray-400 dark:hover:bg-dk-muted',
                    activeIndex === i ? 'ring-2 ring-blue-500 ring-offset-1 dark:ring-offset-dk-card' : '',
                  ]"
                  @click="jumpTo(i)"
                >{{ i + 1 }}</button>
              </div>
            </div>

            <!-- 提交 -->
            <button
              class="inline-flex items-center justify-center gap-1.5 rounded-lg bg-blue-600 px-6 py-2.5 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:opacity-50"
              :disabled="submitting"
              @click="confirmShow = true"
            >
              {{ t('quiz.submit') }}
            </button>
          </div>
        </aside>
      </div>
    </template>

    <ConfirmDialog
      v-model="confirmShow"
      :title="t('quiz.submit_confirm_title')"
      :message="t('quiz.submit_confirm')"
      @confirm="handleSubmit"
    />
  </div>
</template>
