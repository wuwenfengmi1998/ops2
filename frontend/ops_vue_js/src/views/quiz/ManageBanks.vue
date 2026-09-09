<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { quizApi } from '@/api/quiz'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import {
  IconPlus,
  IconPencil,
  IconTrash,
  IconX,
  IconChevronLeftPipe,
  IconChevronRightPipe,
  IconChevronsLeft,
  IconChevronsRight,
  IconArrowLeft,
  IconBook2,
  IconPlayerPlay,
  IconSettings2,
} from '@tabler/icons-vue'

usePageTitle('appname.quiz_questions')
const { t, locale } = useI18n()
const router = useRouter()
const toast = useToastStore()

const banks = ref([])
const totalCount = ref(0)
const pageSize = ref(10)
const currentPage = ref(1)
const keyword = ref('')
const loading = ref(false)

const totalPages = computed(() => Math.ceil(totalCount.value / pageSize.value) || 1)

const pageRange = computed(() => {
  const total = totalPages.value
  const cur = currentPage.value
  let start = Math.max(1, cur - 2)
  let end = Math.min(cur + 4, total)
  if (end - start < 4) start = Math.max(1, end - 4)
  return Array.from({ length: end - start + 1 }, (_, i) => start + i)
})

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
  loading.value = true
  try {
    const { errCode, data } = await quizApi.listBanks({
      keyword: keyword.value,
      page: currentPage.value,
      pageSize: pageSize.value,
    })
    if (errCode === 0) {
      banks.value = data.items ?? []
      totalCount.value = data.total ?? 0
    } else if (errCode === -74) {
      toast.error(t('quiz.permission_denied'))
      router.replace('/questions')
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  } finally {
    loading.value = false
  }
}

function goToPage(page) {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  fetchBanks()
}

function handlePageSizeInput(e) {
  let val = parseInt(e.target.value) || 10
  if (val > 100) val = 100
  if (val < 1) val = 1
  pageSize.value = val
  currentPage.value = 1
  fetchBanks()
}

function handleJumpPageInput(e) {
  const val = parseInt(e.target.value)
  if (val > 0 && val <= totalPages.value) {
    currentPage.value = val
    fetchBanks()
  }
}

function handleSearch() {
  currentPage.value = 1
  fetchBanks()
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  return new Intl.DateTimeFormat(locale.value, {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false,
  }).format(new Date(dateStr))
}

function handleStart(bank, event) {
  event.stopPropagation()
  if (!bank.questionCnt || bank.questionCnt <= 0) {
    toast.warning(t('quiz.bank_empty'))
    return
  }
  router.push({ path: '/quiz/play', query: { bank: bank.id } })
}

function openBankQuestions(row) {
  router.push(`/questions/bank/${row.id}`)
}

// ── 新增/编辑弹窗 ──
const showModal = ref(false)
const saving = ref(false)
const editingId = ref(0)
const form = reactive({
  name: '',
  description: '',
  countSingle: 20,
  countMultiple: 10,
  countJudge: 20,
  countBlank: 10,
  durationSec: 300,
})

function openAdd() {
  editingId.value = 0
  form.name = ''
  form.description = ''
  form.countSingle = 20
  form.countMultiple = 10
  form.countJudge = 20
  form.countBlank = 10
  form.durationSec = 300
  showModal.value = true
}

function openEdit(row, event) {
  event.stopPropagation()
  editingId.value = row.id
  form.name = row.name
  form.description = row.description || ''
  form.countSingle = row.countSingle || 20
  form.countMultiple = row.countMultiple || 10
  form.countJudge = row.countJudge || 20
  form.countBlank = row.countBlank || 10
  form.durationSec = row.durationSec
  showModal.value = true
}

async function handleSave() {
  if (saving.value) return
  if (!form.name.trim()) {
    toast.error(t('banks.err_name_required'))
    return
  }
  saving.value = true
  try {
    const payload = {
      name: form.name.trim(),
      description: form.description.trim(),
      countSingle: form.countSingle,
      countMultiple: form.countMultiple,
      countJudge: form.countJudge,
      countBlank: form.countBlank,
      durationSec: form.durationSec,
    }
    const { errCode } = editingId.value
      ? await quizApi.updateBank({ id: editingId.value, ...payload })
      : await quizApi.addBank(payload)
    if (errCode === 0) {
      toast.success(t('message.save_success'))
      showModal.value = false
      fetchBanks()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  } finally {
    saving.value = false
  }
}

// ── 删除 ──
const deleteTarget = ref(null)
const confirmDelete = ref(false)

function askDelete(row, event) {
  event.stopPropagation()
  deleteTarget.value = row
  confirmDelete.value = true
}

async function doDelete() {
  if (!deleteTarget.value) return
  try {
    const { errCode } = await quizApi.deleteBank(deleteTarget.value.id)
    if (errCode === 0) {
      toast.success(t('message.delete_success'))
      confirmDelete.value = false
      deleteTarget.value = null
      fetchBanks()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  }
}

async function toggleActive(row, event) {
  event.stopPropagation()
  try {
    const { errCode } = await quizApi.updateBank({ id: row.id, active: !row.active })
    if (errCode === 0) {
      row.active = !row.active
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  }
}

onMounted(fetchBanks)
</script>

<template>
  <div class="mx-auto max-w-6xl px-6 py-6">
    <div class="flex flex-col gap-6 rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card">
      <!-- Header -->
      <div class="flex flex-wrap items-center gap-3 border-b border-gray-100 px-6 py-4 dark:border-dk-muted">
        <button
          class="inline-flex items-center gap-1 rounded-lg border border-gray-300 px-2.5 py-1.5 text-xs text-gray-600 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-base dark:text-gray-300 dark:hover:bg-dk-muted"
          @click="router.push('/questions')"
        >
          <IconArrowLeft :size="14" />
          {{ t('banks.back') }}
        </button>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('banks.title') }}</h3>
        <button
          class="ml-auto inline-flex items-center gap-1.5 rounded-lg bg-blue-600 px-3 py-1.5 text-sm font-medium text-white transition-colors hover:bg-blue-700"
          @click="openAdd"
        >
          <IconPlus :size="16" />
          {{ t('banks.add') }}
        </button>
      </div>

      <!-- 筛选栏 -->
      <div class="flex items-center gap-2 px-6 py-3">
        <input
          v-model="keyword"
          type="text"
          :placeholder="t('banks.search_placeholder')"
          class="w-48 rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm text-gray-900 placeholder-gray-400 outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white dark:placeholder-gray-500"
          @input="handleSearch"
          @keydown.enter="handleSearch"
        />
        <button
          class="inline-flex items-center gap-1 rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm text-gray-600 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-base dark:text-gray-300 dark:hover:bg-dk-muted"
          @click="handleSearch"
        >
          {{ t('banks.search') }}
        </button>
      </div>

      <!-- Table -->
      <div class="overflow-x-auto px-0">
        <table class="w-full text-left text-sm text-gray-900">
          <thead>
            <tr class="border-b border-gray-200 bg-gray-50 text-gray-500 dark:border-dk-muted dark:bg-dk-base">
              <th class="w-16 px-6 py-3 font-medium text-gray-500 dark:text-gray-400">No.</th>
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('banks.name') }}</th>
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400 whitespace-nowrap">{{ t('quiz.draw_count') }}</th>
              <th class="w-24 whitespace-nowrap px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('quiz.time_limit') }}</th>
              <th class="w-20 px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('banks.question_count') }}</th>
              <th class="w-16 px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('banks.active') }}</th>
              <th class="w-40 whitespace-nowrap px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('banks.created_at') }}</th>
              <th class="w-32 px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('banks.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="8" class="px-6 py-8 text-center text-gray-400">Loading...</td>
            </tr>
            <tr
              v-for="row in banks"
              :key="row.id"
              class="cursor-pointer border-b border-gray-100 transition-colors hover:bg-blue-50/50 dark:border-dk-muted/50 dark:hover:bg-dk-base/50"
              @click="openBankQuestions(row)"
            >
              <td class="px-6 py-3 text-gray-400">{{ row.id }}</td>
              <td class="px-6 py-3">
                <span class="inline-flex items-center gap-1.5 font-medium text-gray-900 dark:text-white">
                  <IconBook2 :size="16" class="text-blue-500" />
                  {{ row.name }}
                </span>
              </td>
              <td class="whitespace-nowrap px-6 py-3 text-xs text-gray-600 dark:text-gray-300">{{ drawCountText(row) }}</td>
              <td class="whitespace-nowrap px-6 py-3 text-gray-600 dark:text-gray-300">{{ durationText(row.durationSec) }}</td>
              <td class="px-6 py-3">
                <span class="rounded-full bg-blue-100 px-2.5 py-0.5 text-xs font-semibold text-blue-700 dark:bg-blue-900/40 dark:text-blue-400">{{ row.questionCnt }}</span>
              </td>
              <td class="px-6 py-3">
                <button
                  class="relative inline-flex h-5 w-9 items-center rounded-full transition-colors"
                  :class="row.active ? 'bg-blue-600' : 'bg-gray-300 dark:bg-gray-600'"
                  :title="t('banks.active')"
                  @click="toggleActive(row, $event)"
                >
                  <span
                    class="inline-block h-4 w-4 transform rounded-full bg-white transition-transform"
                    :class="row.active ? 'translate-x-4' : 'translate-x-0.5'"
                  ></span>
                </button>
              </td>
              <td class="whitespace-nowrap px-6 py-3 text-gray-500 dark:text-gray-400">{{ formatDate(row.createdAt) }}</td>
              <td class="px-6 py-3">
                <div class="flex items-center gap-1">
                  <button class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-green-600 dark:hover:bg-dk-card" :title="t('quiz.start_btn')" @click="handleStart(row, $event)">
                    <IconPlayerPlay :size="16" />
                  </button>
                  <button class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-blue-600 dark:hover:bg-dk-card" :title="t('banks.manage_questions')" @click="openBankQuestions(row)">
                    <IconSettings2 :size="16" />
                  </button>
                  <button class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-blue-600 dark:hover:bg-dk-card" :title="t('banks.edit')" @click="openEdit(row, $event)">
                    <IconPencil :size="16" />
                  </button>
                  <button class="rounded p-1.5 text-gray-500 transition-colors hover:bg-red-50 hover:text-red-500 dark:hover:bg-dk-card" :title="t('banks.delete')" @click="askDelete(row, $event)">
                    <IconTrash :size="16" />
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="!loading && banks.length === 0">
              <td colspan="8" class="px-6 py-8 text-center text-gray-400">{{ t('banks.empty') }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="flex flex-col items-center justify-between gap-3 border-t border-gray-200 px-6 py-3 sm:flex-row dark:border-dk-muted">
        <div class="flex items-center gap-1.5 text-sm text-gray-500">
          <label>{{ t('quiz.show') }}</label>
          <input type="text" class="w-14 rounded border border-gray-300 px-2 py-1 text-center text-sm text-gray-900 dark:border-dk-muted dark:bg-dk-base dark:text-white" :value="pageSize" @change="handlePageSizeInput" />
          <label>{{ t('quiz.entries') }}</label>
          <span class="ml-1">{{ t('quiz.total_items') }} {{ totalCount }}</span>
        </div>
        <div class="flex items-center gap-1">
          <button class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 disabled:opacity-40 dark:hover:bg-dk-card" :disabled="currentPage <= 1" @click="goToPage(1)"><IconChevronsLeft :size="16" /></button>
          <button class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 disabled:opacity-40 dark:hover:bg-dk-card" :disabled="currentPage <= 1" @click="goToPage(currentPage - 1)"><IconChevronLeftPipe :size="16" /></button>
          <template v-for="a in pageRange" :key="a">
            <button
              class="min-w-[32px] rounded px-2 py-1 text-sm font-medium transition-colors"
              :class="a === currentPage ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-dk-card'"
              @click="goToPage(a)"
            >{{ a }}</button>
          </template>
          <button class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 disabled:opacity-40 dark:hover:bg-dk-card" :disabled="currentPage >= totalPages" @click="goToPage(currentPage + 1)"><IconChevronRightPipe :size="16" /></button>
          <button class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 disabled:opacity-40 dark:hover:bg-dk-card" :disabled="currentPage >= totalPages" @click="goToPage(totalPages)"><IconChevronsRight :size="16" /></button>
          <input type="text" class="ml-2 w-14 rounded border border-gray-300 px-2 py-1 text-center text-sm text-gray-900 dark:border-dk-muted dark:bg-dk-base dark:text-white" @change="handleJumpPageInput" />
        </div>
      </div>
    </div>

    <!-- Add/Edit Modal -->
    <Teleport to="body">
      <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4" @click.self="showModal = false">
        <div class="flex max-h-[90vh] w-full max-w-lg flex-col rounded-xl bg-white shadow-xl dark:bg-dk-card">
          <div class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dk-muted">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ editingId > 0 ? t('banks.edit') : t('banks.add') }}
            </h3>
            <button class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 dark:hover:bg-dk-base" @click="showModal = false">
              <IconX :size="18" />
            </button>
          </div>

          <div class="flex flex-col gap-4 overflow-y-auto px-6 py-5">
            <div class="flex items-center gap-3">
              <label class="w-24 shrink-0 text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('banks.name') }}</label>
              <input
                v-model="form.name"
                type="text"
                :placeholder="t('banks.name_placeholder')"
                maxlength="100"
                class="flex-1 rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm text-gray-900 placeholder-gray-400 outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white dark:placeholder-gray-500"
              />
            </div>
            <div class="flex items-center gap-3">
              <label class="w-24 shrink-0 text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('banks.description') }}</label>
              <input
                v-model="form.description"
                type="text"
                :placeholder="t('banks.description_placeholder')"
                maxlength="200"
                class="flex-1 rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm text-gray-900 placeholder-gray-400 outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white dark:placeholder-gray-500"
              />
            </div>
            <div class="flex flex-col gap-2">
              <label class="w-24 shrink-0 text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('quiz.draw_count') }}</label>
              <div class="grid grid-cols-2 gap-3 pl-24">
                <div class="flex items-center gap-2">
                  <label class="text-xs font-medium text-gray-500 dark:text-gray-400 whitespace-nowrap">{{ t('questions.type_single') }}</label>
                  <input
                    v-model.number="form.countSingle"
                    type="number"
                    min="1"
                    max="100"
                    class="w-20 rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm text-gray-900 outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white"
                  />
                </div>
                <div class="flex items-center gap-2">
                  <label class="text-xs font-medium text-gray-500 dark:text-gray-400 whitespace-nowrap">{{ t('questions.type_multiple') }}</label>
                  <input
                    v-model.number="form.countMultiple"
                    type="number"
                    min="1"
                    max="100"
                    class="w-20 rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm text-gray-900 outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white"
                  />
                </div>
                <div class="flex items-center gap-2">
                  <label class="text-xs font-medium text-gray-500 dark:text-gray-400 whitespace-nowrap">{{ t('questions.type_judge') }}</label>
                  <input
                    v-model.number="form.countJudge"
                    type="number"
                    min="1"
                    max="100"
                    class="w-20 rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm text-gray-900 outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white"
                  />
                </div>
                <div class="flex items-center gap-2">
                  <label class="text-xs font-medium text-gray-500 dark:text-gray-400 whitespace-nowrap">{{ t('questions.type_blank') }}</label>
                  <input
                    v-model.number="form.countBlank"
                    type="number"
                    min="1"
                    max="100"
                    class="w-20 rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm text-gray-900 outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white"
                  />
                </div>
              </div>
              <p class="pl-24 text-xs text-gray-400 dark:text-gray-500">{{ t('banks.draw_count_hint') }}</p>
            </div>
            <div class="flex items-center gap-3">
              <label class="w-24 shrink-0 text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('quiz.time_limit') }}</label>
              <input
                v-model.number="form.durationSec"
                type="number"
                min="0"
                max="7200"
                class="w-24 rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm text-gray-900 outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white"
              />
              <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('banks.duration_seconds_hint') }}</span>
            </div>
          </div>

          <div class="flex justify-end gap-2 border-t border-gray-100 px-6 py-4 dark:border-dk-muted">
            <button
              class="rounded-lg border border-gray-300 px-4 py-1.5 text-sm text-gray-600 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-base dark:text-gray-300 dark:hover:bg-dk-muted"
              @click="showModal = false"
            >
              {{ t('message.cancel') }}
            </button>
            <button
              class="rounded-lg bg-blue-600 px-4 py-1.5 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:opacity-50"
              :disabled="saving"
              @click="handleSave"
            >
              {{ t('message.save') }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <ConfirmDialog
      v-model="confirmDelete"
      :title="t('banks.delete_title')"
      :message="t('banks.delete_msg', { name: deleteTarget?.name || '' })"
      danger
      @confirm="doDelete"
    />
  </div>
</template>
