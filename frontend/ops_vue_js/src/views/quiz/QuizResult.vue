<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { quizApi } from '@/api/quiz'
import { IconChevronLeft, IconCircleCheck, IconCircleX, IconTrophy, IconBook2 } from '@tabler/icons-vue'

usePageTitle('appname.quiz')
const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const toast = useToastStore()

const loading = ref(true)
const session = ref(null)
const review = ref([])

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

const scorePercent = computed(() => {
  if (!session.value || !session.value.totalScore) return 0
  return Math.round((session.value.score / session.value.totalScore) * 100)
})

const passText = computed(() => {
  const s = session.value
  if (!s || !s.totalScore) return t('quiz.result_pass')
  return scorePercent.value >= 60 ? t('quiz.result_pass') : t('quiz.result_fail')
})

function judgeLabel(val) {
  if (val === undefined || val === null || val === '') return t('quiz.no_answer')
  return String(val) === 'true' ? t('quiz.true') : t('quiz.false')
}

function answerText(item) {
  const vals = item.yourAnswerText || []
  if (vals.length === 0) return t('quiz.no_answer')
  return vals.join(', ')
}

onMounted(async () => {
  try {
    const sessionId = Number(route.params.id) || 0
    if (!sessionId) {
      loading.value = false
      toast.error(t('message.server_error'))
      router.replace('/quiz')
      return
    }
    const { errCode, data } = await quizApi.getSession(sessionId)
    if (errCode === 0) {
      session.value = data.session ?? null
      review.value = data.review ?? []
    } else {
      toast.error(t('message.server_error'))
      router.replace('/quiz')
    }
  } catch {
    // 拦截器已处理
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="mx-auto max-w-4xl px-6 py-6">
    <button
      class="mb-4 inline-flex items-center gap-1 text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300"
      @click="router.push('/quiz')"
    >
      <IconChevronLeft :size="16" />
      {{ t('quiz.back_home') }}
    </button>

    <div v-if="loading" class="py-20 text-center text-gray-400">
      <svg class="mx-auto mb-2 h-6 w-6 animate-spin text-gray-400" viewBox="0 0 24 24" fill="none">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
      </svg>
      Loading...
    </div>

    <template v-else-if="session">
      <!-- Score card -->
      <div class="mb-6 flex flex-col items-center rounded-xl border border-gray-200 bg-white p-8 shadow-sm dark:border-dk-muted dark:bg-dk-card">
        <span v-if="session.bankName" class="inline-flex items-center gap-1.5 text-sm font-medium text-gray-600 dark:text-gray-300">
          <IconBook2 :size="15" class="text-blue-500" />
          {{ session.bankName }}
        </span>
        <span
          class="inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-sm font-semibold"
          :class="scorePercent >= 60
            ? 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400'
            : 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400'"
        >
          <IconTrophy :size="16" />
          {{ passText }}
        </span>
        <div class="mt-4 flex items-end gap-2">
          <span class="text-5xl font-bold text-gray-900 dark:text-white">{{ session.score }}</span>
          <span class="mb-1 text-lg text-gray-400 dark:text-gray-500">/ {{ session.totalScore }}</span>
        </div>
        <div class="mt-5 grid grid-cols-2 gap-x-12 gap-y-2 text-center sm:grid-cols-4">
          <div>
            <div class="text-sm text-gray-400 dark:text-gray-500">{{ t('quiz.correct') }}</div>
            <div class="text-lg font-semibold text-green-600 dark:text-green-400">{{ session.correctCount }}</div>
          </div>
          <div>
            <div class="text-sm text-gray-400 dark:text-gray-500">{{ t('quiz.wrong') }}</div>
            <div class="text-lg font-semibold text-red-500 dark:text-red-400">{{ session.wrongCount }}</div>
          </div>
          <div>
            <div class="text-sm text-gray-400 dark:text-gray-500">{{ t('quiz.questions_title') }}</div>
            <div class="text-lg font-semibold text-gray-900 dark:text-white">{{ session.count }}</div>
          </div>
          <div>
            <div class="text-sm text-gray-400 dark:text-gray-500">{{ t('quiz.duration') }}</div>
            <div class="text-lg font-semibold text-gray-900 dark:text-white">{{ session.durationSec }}<span class="text-xs font-normal text-gray-400">s</span></div>
          </div>
        </div>
      </div>

      <!-- Review -->
      <div class="rounded-xl border border-gray-200 bg-white shadow-sm dark:border-dk-muted dark:bg-dk-card">
        <div class="border-b border-gray-100 px-6 py-4 dark:border-dk-muted">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('quiz.review') }}</h3>
        </div>
        <div class="flex flex-col gap-5 p-6">
          <div
            v-for="item in review"
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
                  {{ item.type === 'judge' ? judgeLabel(item.yourAnswerText?.[0]) : answerText(item) }}
                </span>
              </div>
              <div v-if="!item.isCorrect" class="text-sm">
                <span class="text-gray-400 dark:text-gray-500">{{ t('quiz.correct_answer') }}: </span>
                <span class="font-medium text-green-600 dark:text-green-400">
                  {{ item.type === 'judge' ? judgeLabel(item.correctAnswerText?.[0]) : (item.correctAnswerText || []).join(', ') }}
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

          <div v-if="review.length === 0" class="py-8 text-center text-gray-400">{{ t('quiz.review_empty') }}</div>
        </div>
      </div>
    </template>
  </div>
</template>
