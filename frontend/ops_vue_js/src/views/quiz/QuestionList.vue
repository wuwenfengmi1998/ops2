<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToastStore } from '@/stores/toast'
import { useUserStore } from '@/stores/user'
import { usePageTitle } from '@/composables/usePageTitle'
import { quizApi } from '@/api/quiz'
import {
  IconPlayerPlay,
  IconHistory,
  IconTrophy,
  IconBook2,
  IconClock,
  IconListNumbers,
  IconMedal,
  IconPencil,
  IconChevronLeftPipe,
  IconChevronRightPipe,
  IconChevronsLeft,
  IconChevronsRight,
  IconReload,
} from '@tabler/icons-vue'

usePageTitle('appname.quiz_bank')
const { t, locale } = useI18n()
const router = useRouter()
const toast = useToastStore()
const userStore = useUserStore()
const isQuizAdmin = computed(() => userStore.isQuizAdmin)

const tab = ref('banks')

// ── 题库卡片 ──
const banks = ref([])
const loadingBanks = ref(false)

// ── 我的成绩 ──
const sessions = ref([])
const historyTotal = ref(0)
const historyPageSize = ref(10)
const historyPage = ref(1)
const loadingHistory = ref(false)

const historyTotalPages = computed(() => Math.ceil(historyTotal.value / historyPageSize.value) || 1)

const historyPageRange = computed(() => {
  const total = historyTotalPages.value
  const cur = historyPage.value
  let start = Math.max(1, cur - 2)
  let end = Math.min(cur + 4, total)
  if (end - start < 4) start = Math.max(1, end - 4)
  return Array.from({ length: end - start + 1 }, (_, i) => start + i)
})

// ── 排行榜 ──
const board = ref([])
const loadingBoard = ref(false)
const boardBankId = ref(0)
const boardOptions = ref([])

function durationText(sec) {
  if (!sec || sec <= 0) return t('quiz.unlimited')
  const m = Math.floor(sec / 60)
  const s = sec % 60
  return m > 0 ? `${m}${t('quiz.minutes')} ${s}${t('quiz.seconds')}` : `${s}${t('quiz.seconds')}`
}

function drawCountText(bank) {
  const single = bank.countSingle ?? 20
  const multiple = bank.countMultiple ?? 10
  const judge = bank.countJudge ?? 20
  const blank = bank.countBlank ?? 10
  const parts = [
    `${t('banks.single_short')}${single}`,
    `${t('banks.multiple_short')}${multiple}`,
    `${t('banks.judge_short')}${judge}`,
    `${t('banks.blank_short')}${blank}`,
  ]
  return parts.join(' · ')
}

async function fetchBanks() {
  loadingBanks.value = true
  try {
    const { errCode, data } = await quizApi.availableBanks()
    if (errCode === 0) {
      banks.value = data.items ?? []
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  } finally {
    loadingBanks.value = false
  }
}

async function fetchBoardOptions() {
  try {
    const { errCode, data } = await quizApi.availableBanks()
    if (errCode === 0) {
      boardOptions.value = data.items ?? []
      if (boardBankId.value === 0 && boardOptions.value.length > 0) {
        boardBankId.value = boardOptions.value[0].id
        fetchBoard()
      }
    }
  } catch {
    // 拦截器已处理
  }
}

function handleStart(bank) {
  if (!bank.questionCnt || bank.questionCnt <= 0) {
    toast.warning(t('quiz.bank_empty'))
    return
  }
  router.push({ path: '/quiz/play', query: { bank: bank.id } })
}

function handleLearn(bank) {
  router.push({ path: '/quiz/learn', query: { bank: bank.id } })
}

function goManage() {
  router.push('/questions/manage')
}

// ── 我的成绩 ──
async function fetchSessions() {
  loadingHistory.value = true
  try {
    const { errCode, data } = await quizApi.mySessions({ page: historyPage.value, pageSize: historyPageSize.value })
    if (errCode === 0) {
      sessions.value = data.items ?? []
      historyTotal.value = data.total ?? 0
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  } finally {
    loadingHistory.value = false
  }
}

function goHistoryPage(page) {
  if (page < 1 || page > historyTotalPages.value) return
  historyPage.value = page
  fetchSessions()
}

function handleHistoryPageSizeInput(e) {
  let val = parseInt(e.target.value) || 10
  if (val > 100) val = 100
  if (val < 1) val = 1
  historyPageSize.value = val
  historyPage.value = 1
  fetchSessions()
}

function handleHistoryJumpInput(e) {
  const val = parseInt(e.target.value)
  if (val > 0 && val <= historyTotalPages.value) {
    historyPage.value = val
    fetchSessions()
  }
}

function viewSession(id) {
  router.push(`/quiz/result/${id}`)
}

function redoWrong(s) {
  router.push(`/quiz/redo/${s.id}`)
}

// ── 排行榜 ──
async function fetchBoard() {
  loadingBoard.value = true
  try {
    const { errCode, data } = await quizApi.leaderboard({ bankId: boardBankId.value, limit: 10 })
    if (errCode === 0) {
      board.value = data.items ?? []
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  } finally {
    loadingBoard.value = false
  }
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  return new Intl.DateTimeFormat(locale.value, {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false,
  }).format(new Date(dateStr))
}

onMounted(() => {
  fetchBanks()
  fetchBoardOptions()
  fetchSessions()
})
</script>

<template>
  <div class="mx-auto max-w-6xl px-6 py-6">
    <!-- Tabs -->
    <div class="mb-6 flex gap-1 rounded-xl border border-gray-200 bg-white p-1 text-sm font-medium shadow-sm dark:border-dk-muted dark:bg-dk-card">
      <button
        class="flex flex-1 items-center justify-center gap-1.5 rounded-lg px-4 py-2 transition-colors"
        :class="tab === 'banks' ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dk-base'"
        @click="tab = 'banks'"
      >
        <IconBook2 :size="16" />
        {{ t('quiz.tab_banks') }}
      </button>
      <button
        class="flex flex-1 items-center justify-center gap-1.5 rounded-lg px-4 py-2 transition-colors"
        :class="tab === 'history' ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dk-base'"
        @click="tab = 'history'"
      >
        <IconHistory :size="16" />
        {{ t('quiz.tab_history') }}
      </button>
      <button
        class="flex flex-1 items-center justify-center gap-1.5 rounded-lg px-4 py-2 transition-colors"
        :class="tab === 'board' ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dk-base'"
        @click="tab = 'board'"
      >
        <IconTrophy :size="16" />
        {{ t('quiz.tab_leaderboard') }}
      </button>
    </div>

    <!-- ═══ 题库列表(卡片) ═══ -->
    <div v-if="tab === 'banks'" class="rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card">
      <div class="flex flex-wrap items-center justify-between gap-3 border-b border-gray-100 px-6 py-4 dark:border-dk-muted">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('appname.quiz_bank') }}</h3>
        <button
          v-if="isQuizAdmin"
          class="inline-flex items-center gap-1.5 rounded-lg border border-gray-300 px-3 py-1.5 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-base dark:text-gray-300 dark:hover:bg-dk-muted"
          @click="goManage"
        >
          <IconPencil :size="16" />
          {{ t('banks.manage') }}
        </button>
      </div>

      <div class="grid grid-cols-1 gap-4 p-6 sm:grid-cols-2 lg:grid-cols-3">
        <div v-if="loadingBanks" class="col-span-full py-10 text-center text-gray-400">Loading...</div>
        <div
          v-for="bank in banks"
          :key="bank.id"
          class="flex flex-col gap-3 rounded-xl border p-5 transition-colors"
          :class="bank.questionCnt > 0
            ? 'border-gray-200 hover:border-blue-400 hover:shadow-md dark:border-dk-muted dark:hover:border-blue-500'
            : 'border-gray-200 opacity-60 dark:border-dk-muted'"
        >
          <div class="flex items-center gap-2">
            <IconBook2 :size="18" class="shrink-0 text-blue-500" />
            <h4 class="truncate font-semibold text-gray-900 dark:text-white">{{ bank.name }}</h4>
            <div class="ml-auto flex shrink-0 items-center gap-1.5 text-xs text-gray-500 dark:text-gray-400">
              <img
                :src="bank.creatorAvatar ? '/api/static/avatar/' + bank.creatorAvatar : '/ava.svg'"
                class="h-5 w-5 rounded-full object-cover"
                :alt="bank.creatorName || ''"
              />
              <span class="font-medium">{{ bank.creatorName || '-' }}</span>
            </div>
          </div>
          <p v-if="bank.description" class="line-clamp-2 text-xs text-gray-500 dark:text-gray-400">{{ bank.description }}</p>
          <div class="space-y-2 text-xs text-gray-500 dark:text-gray-400">
            <div class="flex items-center justify-between gap-2">
              <span class="inline-flex shrink-0 items-center gap-1">
                <IconListNumbers :size="13" />
                {{ t('quiz.draw_count') }}
              </span>
              <span class="break-words text-right font-medium text-gray-700 dark:text-gray-300">{{ drawCountText(bank) }}</span>
            </div>
            <div class="flex items-center justify-between gap-2">
              <span class="inline-flex shrink-0 items-center gap-1">
                <IconClock :size="13" />
                {{ t('quiz.time_limit') }}
              </span>
              <span class="font-medium text-gray-700 dark:text-gray-300">{{ durationText(bank.durationSec) }}</span>
            </div>
            <div class="flex items-center justify-between gap-2">
              <span class="inline-flex shrink-0 items-center gap-1">
                <IconBook2 :size="13" />
                {{ t('quiz.questions_available') }}
              </span>
              <span class="font-medium text-gray-700 dark:text-gray-300">{{ bank.questionCnt }}</span>
            </div>
            <div v-if="bank.myBest > 0" class="flex items-center justify-between gap-2">
              <span class="inline-flex shrink-0 items-center gap-1 text-amber-600 dark:text-amber-400">
                <IconMedal :size="13" />
                {{ t('quiz.my_best') }}
              </span>
              <span class="font-medium text-amber-600 dark:text-amber-400">{{ bank.myBest }}</span>
            </div>
          </div>
          <div class="mt-auto flex items-center gap-2">
            <button
              class="inline-flex flex-1 items-center justify-center gap-1.5 rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-base dark:text-gray-300 dark:hover:bg-dk-muted"
              @click="handleLearn(bank)"
            >
              <IconBook2 :size="16" />
              {{ t('quiz.learn') }}
            </button>
            <button
              class="inline-flex flex-1 items-center justify-center gap-1.5 rounded-lg px-4 py-2 text-sm font-medium transition-colors disabled:opacity-50"
              :class="bank.questionCnt > 0
                ? 'bg-blue-600 text-white hover:bg-blue-700'
                : 'bg-gray-200 text-gray-500 dark:bg-dk-base dark:text-gray-400'"
              :disabled="bank.questionCnt === 0"
              @click="handleStart(bank)"
            >
              <IconPlayerPlay :size="16" />
              {{ t('quiz.start_btn') }}
            </button>
          </div>
          <p v-if="bank.questionCnt === 0" class="text-center text-xs text-gray-400">{{ t('quiz.bank_empty') }}</p>
        </div>
        <div v-if="!loadingBanks && banks.length === 0" class="col-span-full py-10 text-center text-gray-400">{{ t('quiz.no_banks') }}</div>
      </div>
    </div>

    <!-- ═══ 我的成绩 ═══ -->
    <div v-else-if="tab === 'history'" class="rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card">
      <div class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dk-muted">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('quiz.tab_history') }}</h3>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm text-gray-900">
          <thead>
            <tr class="border-b border-gray-200 bg-gray-50 text-gray-500 dark:border-dk-muted dark:bg-dk-base">
              <th class="w-16 px-6 py-3 font-medium text-gray-500 dark:text-gray-400">No.</th>
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('quiz.bank') }}</th>
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('quiz.score') }}</th>
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('quiz.total_score') }}</th>
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('quiz.correct') }}</th>
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('quiz.wrong') }}</th>
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('quiz.duration') }}</th>
              <th class="whitespace-nowrap px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('quiz.date') }}</th>
              <th class="w-28 px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('quiz.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loadingHistory">
              <td colspan="9" class="px-6 py-8 text-center text-gray-400">Loading...</td>
            </tr>
            <tr v-for="s in sessions" :key="s.id"
              class="cursor-pointer border-b border-gray-100 transition-colors hover:bg-blue-50/50 dark:border-dk-muted/50 dark:hover:bg-dk-base/50"
              @click="viewSession(s.id)">
              <td class="px-6 py-3 text-gray-400">{{ s.id }}</td>
              <td class="max-w-[180px] truncate px-6 py-3 font-medium text-gray-900 dark:text-white">{{ s.bankName || '-' }}</td>
              <td class="px-6 py-3">
                <span class="inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-semibold"
                  :class="s.score >= s.totalScore * 0.6 ? 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400' : 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400'">
                  {{ s.score }}
                </span>
              </td>
              <td class="px-6 py-3 text-gray-600 dark:text-gray-300">{{ s.totalScore }}</td>
              <td class="px-6 py-3 text-green-600 dark:text-green-400">{{ s.correctCount }}</td>
              <td class="px-6 py-3 text-red-500 dark:text-red-400">{{ s.wrongCount }}</td>
              <td class="px-6 py-3 text-gray-500 dark:text-gray-400">{{ s.durationSec }} {{ t('quiz.seconds') }}</td>
              <td class="whitespace-nowrap px-6 py-3 text-gray-500 dark:text-gray-400">{{ formatDate(s.createdAt) }}</td>
              <td class="px-6 py-3">
                <button
                  v-if="s.wrongCount > 0"
                  class="inline-flex items-center gap-1 whitespace-nowrap rounded-lg border border-blue-300 px-2.5 py-1 text-xs font-medium text-blue-600 transition-colors hover:bg-blue-50 dark:border-blue-700 dark:text-blue-400 dark:hover:bg-blue-900/20"
                  @click.stop="redoWrong(s)"
                >
                  <IconReload :size="13" />
                  {{ t('quiz.redo') }}
                </button>
              </td>
            </tr>
            <tr v-if="!loadingHistory && sessions.length === 0">
              <td colspan="9" class="px-6 py-8 text-center text-gray-400">{{ t('quiz.history_empty') }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="flex flex-col items-center justify-between gap-3 border-t border-gray-200 px-6 py-3 sm:flex-row dark:border-dk-muted">
        <div class="flex items-center gap-1.5 text-sm text-gray-500">
          <label>{{ t('quiz.show') }}</label>
          <input type="text" class="w-14 rounded border border-gray-300 px-2 py-1 text-center text-sm text-gray-900 dark:border-dk-muted dark:bg-dk-base dark:text-white" :value="historyPageSize" @change="handleHistoryPageSizeInput" />
          <label>{{ t('quiz.entries') }}</label>
          <span class="ml-1">{{ t('quiz.total_items') }} {{ historyTotal }}</span>
        </div>
        <div class="flex items-center gap-1">
          <button class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 disabled:opacity-40 dark:hover:bg-dk-card" :disabled="historyPage <= 1" @click="goHistoryPage(1)"><IconChevronsLeft :size="16" /></button>
          <button class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 disabled:opacity-40 dark:hover:bg-dk-card" :disabled="historyPage <= 1" @click="goHistoryPage(historyPage - 1)"><IconChevronLeftPipe :size="16" /></button>
          <template v-for="a in historyPageRange" :key="a">
            <button
              class="min-w-[32px] rounded px-2 py-1 text-sm font-medium transition-colors"
              :class="a === historyPage ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-dk-card'"
              @click="goHistoryPage(a)"
            >{{ a }}</button>
          </template>
          <button class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 disabled:opacity-40 dark:hover:bg-dk-card" :disabled="historyPage >= historyTotalPages" @click="goHistoryPage(historyPage + 1)"><IconChevronRightPipe :size="16" /></button>
          <button class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 disabled:opacity-40 dark:hover:bg-dk-card" :disabled="historyPage >= historyTotalPages" @click="goHistoryPage(historyTotalPages)"><IconChevronsRight :size="16" /></button>
          <input type="text" class="ml-2 w-14 rounded border border-gray-300 px-2 py-1 text-center text-sm text-gray-900 dark:border-dk-muted dark:bg-dk-base dark:text-white" @change="handleHistoryJumpInput" />
        </div>
      </div>
    </div>

    <!-- ═══ 排行榜 ═══ -->
    <div v-else class="rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card">
      <div class="flex flex-wrap items-center gap-3 border-b border-gray-100 px-6 py-4 dark:border-dk-muted">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('quiz.tab_leaderboard') }}</h3>
        <select
          v-model="boardBankId"
          class="ml-auto rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
          @change="fetchBoard()"
        >
          <option v-for="b in boardOptions" :key="b.id" :value="b.id">{{ b.name }}</option>
        </select>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm text-gray-900">
          <thead>
            <tr class="border-b border-gray-200 bg-gray-50 text-gray-500 dark:border-dk-muted dark:bg-dk-base">
              <th class="w-16 px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('quiz.rank') }}</th>
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('quiz.player') }}</th>
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('quiz.best_score') }}</th>
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('quiz.attempts') }}</th>
              <th class="whitespace-nowrap px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('quiz.last_at') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loadingBoard">
              <td colspan="5" class="px-6 py-8 text-center text-gray-400">Loading...</td>
            </tr>
            <tr v-for="item in board" :key="item.userId" class="border-b border-gray-100 dark:border-dk-muted/50">
              <td class="px-6 py-3">
                <span class="inline-flex h-6 w-6 items-center justify-center rounded-full text-xs font-semibold"
                  :class="item.rank === 1 ? 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/40 dark:text-yellow-400'
                    : item.rank === 2 ? 'bg-gray-200 text-gray-700 dark:bg-gray-700 dark:text-gray-300'
                    : item.rank === 3 ? 'bg-orange-100 text-orange-700 dark:bg-orange-900/40 dark:text-orange-400'
                    : 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-400'">
                  {{ item.rank }}
                </span>
              </td>
              <td class="px-6 py-3 font-medium text-gray-900 dark:text-white">
                <span class="inline-flex items-center gap-2">
                  <img
                    :src="item.avatar ? '/api/static/avatar/' + item.avatar : '/ava.svg'"
                    class="h-6 w-6 rounded-full object-cover"
                    :alt="item.name || ''"
                  />
                  {{ item.name || '-' }}
                </span>
              </td>
              <td class="px-6 py-3 font-semibold text-green-600 dark:text-green-400">{{ item.best }}</td>
              <td class="px-6 py-3 text-gray-600 dark:text-gray-300">{{ item.sessions }}</td>
              <td class="whitespace-nowrap px-6 py-3 text-gray-500 dark:text-gray-400">{{ formatDate(item.lastAt) }}</td>
            </tr>
            <tr v-if="!loadingBoard && board.length === 0">
              <td colspan="5" class="px-6 py-8 text-center text-gray-400">{{ t('quiz.board_empty') }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
