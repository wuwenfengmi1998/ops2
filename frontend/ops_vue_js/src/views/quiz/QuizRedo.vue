<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { quizApi } from '@/api/quiz'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import { IconChevronLeft, IconCircleCheck, IconCircleX, IconClock, IconTrophy } from '@tabler/icons-vue'

usePageTitle('appname.quiz')
const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const toast = useToastStore()

const sessionId = Number(route.params.id) || 0
const loading = ref(true)
const questions = ref([])
const answers = reactive({})
const seconds = ref(0)
const submitting = ref(false)
const confirmShow = ref(false)
const result = ref(null)
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

const answeredCount = computed(() => {
  let n = 0
  for (const q of questions.value) {
    const ans = answers[q.qid]
    if (q.type === 'single' || q.type === 'judge') {
      if (ans !== undefined && ans !== null && ans !== '') n++
    } else if (q.type === 'multiple') {
      if (Array.isArray(ans) && ans.length > 0) n++
    } else if (q.type === 'blank') {
      const arr = Array.isArray(ans) ? ans : []
      const has = arr.some(v => (v || '').trim() !== '')
      if (has) n++
    }
  }
  return n
})

const timeText = computed(() => {
  const m = String(Math.floor(seconds.value / 60)).padStart(2, '0')
  const s = String(seconds.value % 60).padStart(2, '0')
  return `${m}:${s}`
})

function initAnswers() {
  for (const q of questions.value) {
    if (q.type === 'single' || q.type === 'judge') {
      answers[q.qid] = ''
    } else if (q.type === 'multiple') {
      answers[q.qid] = []
    } else if (q.type === 'blank') {
      answers[q.qid] = Array.from({ length: q.blanks }, () => '')
    }
  }
}

function toggleOption(q, idx) {
  if (q.type === 'single') {
    answers[q.qid] = idx
  } else {
    const cur = answers[q.qid] || []
    const pos = cur.indexOf(idx)
    if (pos >= 0) cur.splice(pos, 1)
    else cur.push(idx)
  }
}

function chooseJudge(q, val) {
  answers[q.qid] = val
}

async function doSubmit() {
  if (submitting.value) return
  submitting.value = true
  try {
    const answerList = questions.value.map(q => ({
      qid: q.qid,
      answer: answers[q.qid] ?? '',
    }))
    const { errCode, data } = await quizApi.redoQuiz(sessionId, seconds.value, answerList)
    if (errCode === 0) {
      result.value = data
      if (timer) clearInterval(timer)
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  } finally {
    submitting.value = false
  }
}

function handleSubmit() {
  confirmShow.value = false
  doSubmit()
}

function restart() {
  result.value = null
  seconds.value = 0
  initAnswers()
  timer = setInterval(() => { seconds.value++ }, 1000)
}

function judgeLabel(val) {
  if (val === undefined || val === null || val === '') return t('quiz.no_answer')
  return String(val) === 'true' ? t('quiz.true') : t('quiz.false')
}

function answerText(vals) {
  const arr = vals || []
  if (arr.length === 0) return t('quiz.no_answer')
  return arr.join(', ')
}

onMounted(async () => {
  if (!sessionId) {
    router.replace('/questions')
    return
  }
  try {
    const { errCode, data } = await quizApi.getSession(sessionId)
    if (errCode === 0) {
      const wrong = (data.review ?? []).filter(r => !r.isCorrect)
      if (wrong.length === 0) {
        toast.info(t('quiz.redo_none'))
        router.replace('/questions')
        return
      }
      questions.value = wrong.map(r => ({
        qid: r.qid,
        type: r.type,
        title: r.title,
        options: r.options || [],
        score: r.score,
        blanks: r.type === 'blank' ? (r.correctAnswerText?.length || 1) : 0,
      }))
      initAnswers()
      timer = setInterval(() => { seconds.value++ }, 1000)
    } else if (errCode === -73) {
      toast.error(t('message.server_error'))
      router.replace('/questions')
    } else {
      toast.error(t('message.server_error'))
      router.replace('/questions')
    }
  } catch {
    router.replace('/questions')
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <div class="mx-auto max-w-4xl px-6 py-6">
    <button
      class="mb-4 inline-flex items-center gap-1 text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300"
      @click="router.push('/questions')"
    >
      <IconChevronLeft :size="16" />
      {{ t('quiz.back_home') }}
    </button>

    <!-- 成绩(本地) -->
    <div v-if="result" class="mb-6 flex flex-col items-center rounded-xl border border-gray-200 bg-white p-8 shadow-sm dark:border-dk-muted dark:bg-dk-card">
      <span class="inline-flex items-center gap-1.5 text-sm font-medium text-amber-600 dark:text-amber-400">
        {{ t('quiz.redo_local_hint') }}
      </span>
      <span
        class="mt-3 inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-sm font-semibold"
        :class="result.totalScore > 0 && result.score / result.totalScore >= 0.6
          ? 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400'
          : 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400'"
      >
        <IconTrophy :size="16" />
        {{ result.totalScore > 0 && result.score / result.totalScore >= 0.6 ? t('quiz.result_pass') : t('quiz.result_fail') }}
      </span>
      <div class="mt-4 flex items-end gap-2">
        <span class="text-5xl font-bold text-gray-900 dark:text-white">{{ result.score }}</span>
        <span class="mb-1 text-lg text-gray-400 dark:text-gray-500">/ {{ result.totalScore }}</span>
      </div>
      <div class="mt-5 grid grid-cols-2 gap-x-12 gap-y-2 text-center sm:grid-cols-4">
        <div>
          <div class="text-sm text-gray-400 dark:text-gray-500">{{ t('quiz.correct') }}</div>
          <div class="text-lg font-semibold text-green-600 dark:text-green-400">{{ result.correct }}</div>
        </div>
        <div>
          <div class="text-sm text-gray-400 dark:text-gray-500">{{ t('quiz.wrong') }}</div>
          <div class="text-lg font-semibold text-red-500 dark:text-red-400">{{ result.wrong }}</div>
        </div>
        <div>
          <div class="text-sm text-gray-400 dark:text-gray-500">{{ t('quiz.questions_title') }}</div>
          <div class="text-lg font-semibold text-gray-900 dark:text-white">{{ result.correct + result.wrong }}</div>
        </div>
        <div>
          <div class="text-sm text-gray-400 dark:text-gray-500">{{ t('quiz.duration') }}</div>
          <div class="text-lg font-semibold text-gray-900 dark:text-white">{{ result.durationSec }}<span class="text-xs font-normal text-gray-400">s</span></div>
        </div>
      </div>
      <div class="mt-6 flex gap-2">
        <button
          class="rounded-lg border border-gray-300 px-4 py-1.5 text-sm text-gray-600 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-base dark:text-gray-300 dark:hover:bg-dk-muted"
          @click="restart"
        >
          {{ t('quiz.redo_again') }}
        </button>
        <button
          class="rounded-lg bg-blue-600 px-4 py-1.5 text-sm font-medium text-white transition-colors hover:bg-blue-700"
          @click="router.push('/questions')"
        >
          {{ t('quiz.back_home') }}
        </button>
      </div>
    </div>

    <!-- Loading -->
    <div v-else-if="loading" class="py-20 text-center text-gray-400">
      <svg class="mx-auto mb-2 h-6 w-6 animate-spin text-gray-400" viewBox="0 0 24 24" fill="none">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
      </svg>
      Loading...
    </div>

    <template v-else>
      <!-- 顶部: 标题+计时 -->
      <div class="mb-6 flex items-center gap-3 rounded-xl border border-gray-200 bg-white px-6 py-3 shadow-sm dark:border-dk-muted dark:bg-dk-card">
        <span class="inline-flex items-center gap-1.5 text-sm font-medium text-gray-700 dark:text-gray-300">
          <IconClock :size="16" />
          {{ t('quiz.elapsed') }}: {{ timeText }}
        </span>
        <span class="ml-auto text-sm text-gray-500 dark:text-gray-400">{{ t('quiz.redo_title') }} {{ answeredCount }} / {{ questions.length }}</span>
      </div>

      <div
        v-for="(q, i) in questions"
        :key="q.qid"
        class="mb-5 rounded-xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dk-muted dark:bg-dk-card"
      >
        <div class="mb-4 flex items-center gap-3">
          <span class="flex h-7 min-w-7 items-center justify-center rounded-full bg-blue-100 px-1 text-sm font-bold text-blue-700 dark:bg-blue-900/40 dark:text-blue-400">{{ i + 1 }}</span>
          <span class="rounded-full px-2.5 py-0.5 text-xs font-semibold" :class="typeColors[q.type]">{{ t(typeLabels[q.type]) }}</span>
          <span class="ml-auto text-xs font-medium text-gray-400 dark:text-gray-500">{{ q.score }} {{ t('quiz.points') }}</span>
        </div>
        <p class="mb-4 whitespace-pre-wrap text-sm font-medium text-gray-900 dark:text-white">{{ q.title }}</p>

        <!-- 单选 -->
        <div v-if="q.type === 'single'" class="flex flex-col gap-2">
          <button
            v-for="(opt, oi) in q.options"
            :key="oi"
            class="flex items-center gap-3 rounded-lg border px-4 py-2.5 text-left text-sm transition-colors"
            :class="answers[q.qid] === oi
              ? 'border-blue-500 bg-blue-50 text-blue-700 dark:border-blue-500 dark:bg-blue-900/30 dark:text-blue-300'
              : 'border-gray-200 text-gray-700 hover:border-gray-300 hover:bg-gray-50 dark:border-dk-muted dark:text-gray-300 dark:hover:bg-dk-base'"
            @click="toggleOption(q, oi)"
          >
            <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded-full border"
              :class="answers[q.qid] === oi ? 'border-blue-500 bg-blue-500' : 'border-gray-300 dark:border-dk-muted'">
              <span v-if="answers[q.qid] === oi" class="h-2 w-2 rounded-full bg-white"></span>
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
            :class="(answers[q.qid] || []).includes(oi)
              ? 'border-purple-500 bg-purple-50 text-purple-700 dark:border-purple-500 dark:bg-purple-900/30 dark:text-purple-300'
              : 'border-gray-200 text-gray-700 hover:border-gray-300 hover:bg-gray-50 dark:border-dk-muted dark:text-gray-300 dark:hover:bg-dk-base'"
            @click="toggleOption(q, oi)"
          >
            <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded border"
              :class="(answers[q.qid] || []).includes(oi) ? 'border-purple-500 bg-purple-500' : 'border-gray-300 dark:border-dk-muted'">
              <span v-if="(answers[q.qid] || []).includes(oi)" class="h-2 w-2 rounded-sm bg-white"></span>
            </span>
            <span>{{ opt }}</span>
          </button>
        </div>

        <!-- 填空 -->
        <div v-else-if="q.type === 'blank'" class="flex flex-col gap-3">
          <div v-for="(_, bi) in (answers[q.qid] || [])" :key="bi" class="flex items-center gap-3">
            <span class="text-sm text-gray-400 dark:text-gray-500">{{ t('quiz.blank') }}{{ bi + 1 }}</span>
            <input
              v-model="answers[q.qid][bi]"
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
            :class="answers[q.qid] === 'true'
              ? 'border-green-500 bg-green-50 text-green-700 dark:border-green-500 dark:bg-green-900/30 dark:text-green-300'
              : 'border-gray-200 text-gray-700 hover:border-gray-300 hover:bg-gray-50 dark:border-dk-muted dark:text-gray-300 dark:hover:bg-dk-base'"
            @click="chooseJudge(q, 'true')"
          >
            {{ t('quiz.true') }}
          </button>
          <button
            class="flex-1 rounded-lg border px-4 py-2.5 text-sm font-medium transition-colors"
            :class="answers[q.qid] === 'false'
              ? 'border-red-500 bg-red-50 text-red-700 dark:border-red-500 dark:bg-red-900/30 dark:text-red-300'
              : 'border-gray-200 text-gray-700 hover:border-gray-300 hover:bg-gray-50 dark:border-dk-muted dark:text-gray-300 dark:hover:bg-dk-base'"
            @click="chooseJudge(q, 'false')"
          >
            {{ t('quiz.false') }}
          </button>
        </div>
      </div>

      <!-- 提交 -->
      <div class="flex justify-center pb-8 pt-2">
        <button
          class="inline-flex items-center gap-1.5 rounded-lg bg-blue-600 px-10 py-2.5 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:opacity-50"
          :disabled="submitting"
          @click="confirmShow = true"
        >
          {{ t('quiz.submit') }}
        </button>
      </div>
    </template>

    <!-- 本地结果: 逐题解析 -->
    <div v-if="result" class="rounded-xl border border-gray-200 bg-white shadow-sm dark:border-dk-muted dark:bg-dk-card">
      <div class="border-b border-gray-100 px-6 py-4 dark:border-dk-muted">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('quiz.review') }}</h3>
      </div>
      <div class="flex flex-col gap-5 p-6">
        <div
          v-for="item in result.review"
          :key="item.qid"
          class="rounded-lg border p-5"
          :class="item.isCorrect
            ? 'border-green-200 bg-green-50/50 dark:border-green-900/40 dark:bg-green-900/10'
            : 'border-red-200 bg-red-50/50 dark:border-red-900/40 dark:bg-red-900/10'"
        >
          <div class="mb-3 flex items-center gap-3">
            <span class="flex h-6 min-w-6 items-center justify-center rounded-full bg-gray-100 px-1 text-xs font-bold text-gray-600 dark:bg-dk-base dark:text-gray-400">{{ item.index }}</span>
            <span class="rounded-full px-2.5 py-0.5 text-xs font-semibold" :class="typeColors[item.type]">{{ t(typeLabels[item.type]) }}</span>
            <span
              class="ml-auto inline-flex items-center gap-1 text-xs font-semibold"
              :class="item.isCorrect ? 'text-green-600 dark:text-green-400' : 'text-red-500 dark:text-red-400'"
            >
              <IconCircleCheck v-if="item.isCorrect" :size="15" />
              <IconCircleX v-else :size="15" />
              {{ item.isCorrect ? '+' + item.gotScore : '0' }} / {{ item.score }}
            </span>
          </div>
          <p class="mb-3 whitespace-pre-wrap text-sm font-medium text-gray-900 dark:text-white">{{ item.title }}</p>
          <div v-if="item.options.length > 0" class="mb-3 flex flex-col gap-1.5">
            <div
              v-for="(opt, oi) in item.options"
              :key="oi"
              class="rounded-md px-3 py-1.5 text-sm"
              :class="item.correctAnswerText.includes(opt)
                ? 'bg-green-100 text-green-800 dark:bg-green-900/40 dark:text-green-300'
                : 'bg-gray-50 text-gray-600 dark:bg-dk-base dark:text-gray-400'"
            >
              {{ String.fromCharCode(65 + oi) }}. {{ opt }}
            </div>
          </div>
          <div class="flex flex-col gap-1.5 sm:flex-row sm:gap-8">
            <div class="text-sm">
              <span class="text-gray-400 dark:text-gray-500">{{ t('quiz.your_answer') }}: </span>
              <span class="font-medium" :class="item.isCorrect ? 'text-green-600 dark:text-green-400' : 'text-red-500 dark:text-red-400'">
                {{ item.type === 'judge' ? judgeLabel(item.yourAnswerText?.[0]) : answerText(item.yourAnswerText) }}
              </span>
            </div>
            <div v-if="!item.isCorrect" class="text-sm">
              <span class="text-gray-400 dark:text-gray-500">{{ t('quiz.correct_answer') }}: </span>
              <span class="font-medium text-green-600 dark:text-green-400">
                {{ item.type === 'judge' ? judgeLabel(item.correctAnswerText?.[0]) : answerText(item.correctAnswerText) }}
              </span>
            </div>
          </div>
          <div
            v-if="item.explain"
            class="mt-3 rounded-lg border border-blue-100 bg-blue-50/50 px-4 py-2.5 text-sm leading-relaxed text-gray-700 dark:border-blue-900/40 dark:bg-blue-900/10 dark:text-gray-300"
          >
            <span class="font-medium text-blue-600 dark:text-blue-400">{{ t('quiz.explain') }}: </span>
            <span class="whitespace-pre-wrap">{{ item.explain }}</span>
          </div>
        </div>
      </div>
    </div>

    <ConfirmDialog
      v-model="confirmShow"
      :title="t('quiz.submit_confirm_title')"
      :message="t('quiz.submit_confirm')"
      @confirm="handleSubmit"
    />
  </div>
</template>
