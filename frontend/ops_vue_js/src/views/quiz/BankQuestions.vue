<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
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
} from '@tabler/icons-vue'

usePageTitle('appname.quiz_questions')
const { t, locale } = useI18n()
const route = useRoute()
const router = useRouter()
const toast = useToastStore()

const bankId = Number(route.params.id)
const bankInfo = ref(null)

const items = ref([])
const totalCount = ref(0)
const pageSize = ref(10)
const currentPage = ref(1)
const typeFilter = ref('')
const keyword = ref('')
const loading = ref(false)

const typeOptions = [
  { value: '', labelKey: 'questions.filter_all' },
  { value: 'single', labelKey: 'questions.type_single' },
  { value: 'multiple', labelKey: 'questions.type_multiple' },
  { value: 'blank', labelKey: 'questions.type_blank' },
  { value: 'judge', labelKey: 'questions.type_judge' },
]

const typeLabels = {
  single: 'questions.type_single',
  multiple: 'questions.type_multiple',
  blank: 'questions.type_blank',
  judge: 'questions.type_judge',
}

const typeColors = {
  single: 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-400',
  multiple: 'bg-purple-100 text-purple-700 dark:bg-purple-900/40 dark:text-purple-400',
  blank: 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-400',
  judge: 'bg-teal-100 text-teal-700 dark:bg-teal-900/40 dark:text-teal-400',
}

const totalPages = computed(() => Math.ceil(totalCount.value / pageSize.value) || 1)

const pageRange = computed(() => {
  const total = totalPages.value
  const cur = currentPage.value
  let start = Math.max(1, cur - 2)
  let end = Math.min(cur + 4, total)
  if (end - start < 4) start = Math.max(1, end - 4)
  return Array.from({ length: end - start + 1 }, (_, i) => start + i)
})

async function fetchBankInfo() {
  try {
    const { errCode, data } = await quizApi.listBanks({ page: 1, pageSize: 100 })
    if (errCode === 0) {
      const found = (data.items ?? []).find(b => b.id === bankId)
      if (found) bankInfo.value = found
    }
  } catch {
    // 拦截器已处理
  }
}

async function fetchQuestions() {
  loading.value = true
  try {
    const { errCode, data } = await quizApi.listQuestions({
      bankId,
      type: typeFilter.value,
      keyword: keyword.value,
      page: currentPage.value,
      pageSize: pageSize.value,
    })
    if (errCode === 0) {
      items.value = data.items ?? []
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
  fetchQuestions()
}

function handlePageSizeInput(e) {
  let val = parseInt(e.target.value) || 10
  if (val > 100) val = 100
  if (val < 1) val = 1
  pageSize.value = val
  currentPage.value = 1
  fetchQuestions()
}

function handleJumpPageInput(e) {
  const val = parseInt(e.target.value)
  if (val > 0 && val <= totalPages.value) {
    currentPage.value = val
    fetchQuestions()
  }
}

function handleSearch() {
  currentPage.value = 1
  fetchQuestions()
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  return new Intl.DateTimeFormat(locale.value, {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false,
  }).format(new Date(dateStr))
}

// ── 新增/编辑弹窗 ──
const showModal = ref(false)
const saving = ref(false)
const editingId = ref(0)
const form = reactive({
  type: 'single',
  title: '',
  options: ['', ''],
  singleAnswer: null,
  multipleAnswer: [],
  blankAnswers: [''],
  judgeAnswer: true,
  explain: '',
  score: 5,
  active: true,
})

const isOptionType = computed(() => form.type === 'single' || form.type === 'multiple')

function openAdd() {
  editingId.value = 0
  form.type = 'single'
  form.title = ''
  form.options = ['', '']
  form.singleAnswer = null
  form.multipleAnswer = []
  form.blankAnswers = ['']
  form.judgeAnswer = true
  form.explain = ''
  form.score = 5
  form.active = true
  showModal.value = true
}

function openEdit(row) {
  editingId.value = row.id
  form.type = row.type
  form.title = row.title
  form.score = row.score
  form.active = row.active
  form.explain = row.explain || ''
  form.options = row.options?.length ? [...row.options] : ['', '']
  form.singleAnswer = null
  form.multipleAnswer = []
  form.blankAnswers = ['']
  form.judgeAnswer = true
  const ans = row.answer
  if (row.type === 'single') {
    form.singleAnswer = typeof ans === 'number' ? ans : parseInt(ans) || null
  } else if (row.type === 'multiple') {
    form.multipleAnswer = Array.isArray(ans) ? ans.map(Number) : []
  } else if (row.type === 'blank') {
    form.blankAnswers = Array.isArray(ans) && ans.length ? ans.map(String) : ['']
  } else if (row.type === 'judge') {
    form.judgeAnswer = ans === true || ans === 'true'
  }
  showModal.value = true
}

function toggleMultiple(i) {
  const arr = form.multipleAnswer
  const idx = arr.indexOf(i)
  if (idx >= 0) arr.splice(idx, 1)
  else arr.push(i)
}

function addOption() {
  form.options.push('')
}

function removeOption(idx) {
  if (form.options.length <= 2) return
  form.options.splice(idx, 1)
  if (form.type === 'single' && form.singleAnswer === idx) form.singleAnswer = null
  if (form.singleAnswer !== null && form.singleAnswer >= idx) {
    form.singleAnswer -= 1
  }
  if (form.type === 'multiple') {
    form.multipleAnswer = form.multipleAnswer
      .filter(i => i !== idx)
      .map(i => (i > idx ? i - 1 : i))
  }
}

function addBlank() {
  form.blankAnswers.push('')
}

function removeBlank(idx) {
  if (form.blankAnswers.length <= 1) return
  form.blankAnswers.splice(idx, 1)
}

function validate() {
  if (!form.title.trim()) {
    toast.error(t('questions.err_title_required'))
    return false
  }
  if (isOptionType.value) {
    const opts = form.options.filter(o => o.trim() !== '')
    if (opts.length < 2) {
      toast.error(t('questions.err_options_required'))
      return false
    }
    form.options = form.options.map(o => (o.trim() === '' ? ' ' : o))
  }
  if (form.type === 'single' && (form.singleAnswer === null || form.singleAnswer === undefined)) {
    toast.error(t('questions.err_answer_required'))
    return false
  }
  if (form.type === 'multiple' && form.multipleAnswer.length === 0) {
    toast.error(t('questions.err_answer_required'))
    return false
  }
  if (form.type === 'blank' && !form.blankAnswers.some(b => b.trim() !== '')) {
    toast.error(t('questions.err_answer_required'))
    return false
  }
  return true
}

function buildAnswer() {
  if (form.type === 'single') return form.singleAnswer
  if (form.type === 'multiple') return form.multipleAnswer
  if (form.type === 'blank') return form.blankAnswers.filter(b => b.trim() !== '')
  return form.judgeAnswer
}

async function handleSave() {
  if (saving.value || !validate()) return
  saving.value = true
  try {
    const payload = {
      type: form.type,
      title: form.title.trim(),
      options: isOptionType.value ? form.options.filter(o => o.trim() !== '') : [],
      answer: buildAnswer(),
      explain: form.explain.trim(),
      score: form.score,
      active: form.active,
    }
    const { errCode } = editingId.value
      ? await quizApi.updateQuestion({ id: editingId.value, ...payload })
      : await quizApi.addQuestion({ bankId, ...payload })
    if (errCode === 0) {
      toast.success(t('message.save_success'))
      showModal.value = false
      fetchQuestions()
      fetchBankInfo()
    } else if (errCode === -76) {
      toast.error(t('questions.err_answer_format'))
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

function askDelete(row) {
  deleteTarget.value = row
  confirmDelete.value = true
}

async function doDelete() {
  if (!deleteTarget.value) return
  try {
    const { errCode } = await quizApi.deleteQuestion(deleteTarget.value.id)
    if (errCode === 0) {
      toast.success(t('message.delete_success'))
      confirmDelete.value = false
      deleteTarget.value = null
      fetchQuestions()
      fetchBankInfo()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  }
}

async function toggleActive(row) {
  try {
    const { errCode } = await quizApi.updateQuestion({ id: row.id, active: !row.active })
    if (errCode === 0) {
      row.active = !row.active
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  }
}

onMounted(() => {
  fetchBankInfo()
  fetchQuestions()
})
</script>

<template>
  <div class="mx-auto max-w-6xl px-6 py-6">
    <div class="flex flex-col gap-6 rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card">
      <!-- Header -->
      <div class="flex items-center gap-3 border-b border-gray-100 px-6 py-4 dark:border-dk-muted">
        <button
          class="inline-flex items-center gap-1 rounded-lg border border-gray-300 px-2.5 py-1.5 text-xs text-gray-600 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-base dark:text-gray-300 dark:hover:bg-dk-muted"
          @click="router.push('/questions')"
        >
          <IconArrowLeft :size="14" />
          {{ t('banks.back') }}
        </button>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ bankInfo?.name || '...' }}</h3>
        <span v-if="bankInfo" class="rounded-full bg-blue-100 px-2.5 py-0.5 text-xs font-semibold text-blue-700 dark:bg-blue-900/40 dark:text-blue-400">
          {{ t('banks.question_count') }} {{ bankInfo.questionCnt }}
        </span>
        <button
          class="ml-auto inline-flex items-center gap-1.5 rounded-lg bg-blue-600 px-3 py-1.5 text-sm font-medium text-white transition-colors hover:bg-blue-700"
          @click="openAdd"
        >
          <IconPlus :size="16" />
          {{ t('questions.add') }}
        </button>
      </div>

      <!-- Toolbar -->
      <div class="flex flex-col gap-3 px-6 py-3 sm:flex-row sm:items-center">
        <div class="flex flex-wrap items-center gap-2">
          <select
            v-model="typeFilter"
            class="rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
            @change="currentPage = 1; fetchQuestions()"
          >
            <option v-for="opt in typeOptions" :key="opt.value" :value="opt.value">
              {{ t(opt.labelKey) }}
            </option>
          </select>
        </div>
        <div class="flex items-center gap-2 sm:ml-auto">
          <input
            v-model="keyword"
            type="text"
            :placeholder="t('questions.search_placeholder')"
            class="w-48 rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm text-gray-900 placeholder-gray-400 outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white dark:placeholder-gray-500"
            @input="handleSearch"
            @keydown.enter="handleSearch"
          />
          <button
            class="inline-flex items-center gap-1 rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm text-gray-600 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-base dark:text-gray-300 dark:hover:bg-dk-muted"
            @click="handleSearch"
          >
            {{ t('questions.search') }}
          </button>
        </div>
      </div>

      <!-- Table -->
      <div class="overflow-x-auto px-0">
        <table class="w-full text-left text-sm text-gray-900">
          <thead>
            <tr class="border-b border-gray-200 bg-gray-50 text-gray-500 dark:border-dk-muted dark:bg-dk-base">
              <th class="w-16 px-6 py-3 font-medium text-gray-500 dark:text-gray-400">No.</th>
              <th class="w-32 px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('questions.type') }}</th>
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('questions.question') }}</th>
              <th class="w-16 px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('questions.score') }}</th>
              <th class="w-16 px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('questions.active') }}</th>
              <th class="w-32 whitespace-nowrap px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('questions.creator') }}</th>
              <th class="w-40 whitespace-nowrap px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('questions.created_at') }}</th>
              <th class="w-28 px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('questions.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="8" class="px-6 py-8 text-center text-gray-400">Loading...</td>
            </tr>
            <tr v-for="row in items" :key="row.id" class="border-b border-gray-100 transition-colors hover:bg-blue-50/50 dark:border-dk-muted/50 dark:hover:bg-dk-base/50">
              <td class="px-6 py-3 text-gray-400">{{ row.id }}</td>
              <td class="px-6 py-3">
                <span class="rounded-full px-2.5 py-0.5 text-xs font-semibold" :class="typeColors[row.type]">
                  {{ t(typeLabels[row.type]) }}
                </span>
              </td>
              <td class="max-w-[260px] truncate px-6 py-3 font-medium text-gray-900 dark:text-white">{{ row.title }}</td>
              <td class="px-6 py-3 text-gray-600 dark:text-gray-300">{{ row.score }}</td>
              <td class="px-6 py-3">
                <button
                  class="relative inline-flex h-5 w-9 items-center rounded-full transition-colors"
                  :class="row.active ? 'bg-blue-600' : 'bg-gray-300 dark:bg-gray-600'"
                  :title="t('questions.active')"
                  @click="toggleActive(row)"
                >
                  <span
                    class="inline-block h-4 w-4 transform rounded-full bg-white transition-transform"
                    :class="row.active ? 'translate-x-4' : 'translate-x-0.5'"
                  ></span>
                </button>
              </td>
              <td class="px-6 py-3 text-gray-600 dark:text-gray-300">{{ row.creatorName || '-' }}</td>
              <td class="whitespace-nowrap px-6 py-3 text-gray-500 dark:text-gray-400">{{ formatDate(row.createdAt) }}</td>
              <td class="px-6 py-3">
                <div class="flex items-center gap-1">
                  <button class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-blue-600 dark:hover:bg-dk-card" :title="t('questions.edit')" @click="openEdit(row)">
                    <IconPencil :size="16" />
                  </button>
                  <button class="rounded p-1.5 text-gray-500 transition-colors hover:bg-red-50 hover:text-red-500 dark:hover:bg-dk-card" :title="t('questions.delete')" @click="askDelete(row)">
                    <IconTrash :size="16" />
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="!loading && items.length === 0">
              <td colspan="8" class="px-6 py-8 text-center text-gray-400">{{ t('questions.empty') }}</td>
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
        <div class="flex max-h-[90vh] w-full max-w-2xl flex-col rounded-xl bg-white shadow-xl dark:bg-dk-card">
          <div class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dk-muted">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ editingId > 0 ? t('questions.edit') : t('questions.add') }}
            </h3>
            <button class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 dark:hover:bg-dk-base" @click="showModal = false">
              <IconX :size="18" />
            </button>
          </div>

          <div class="flex flex-col gap-4 overflow-y-auto px-6 py-5">
            <!-- Type -->
            <div class="flex items-center gap-3">
              <label class="w-20 text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('questions.type') }}</label>
              <select
                v-model="form.type"
                class="rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
              >
                <option value="single">{{ t('questions.type_single') }}</option>
                <option value="multiple">{{ t('questions.type_multiple') }}</option>
                <option value="blank">{{ t('questions.type_blank') }}</option>
                <option value="judge">{{ t('questions.type_judge') }}</option>
              </select>
            </div>

            <!-- Title -->
            <div class="flex items-start gap-3">
              <label class="w-20 pt-2 text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('questions.question') }}</label>
              <textarea
                v-model="form.title"
                rows="3"
                :placeholder="t('questions.title_placeholder')"
                class="flex-1 resize-y rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 placeholder-gray-400 outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white dark:placeholder-gray-500"
              ></textarea>
            </div>

            <!-- Options -->
            <div v-if="isOptionType" class="flex flex-col gap-2">
              <div class="flex items-center gap-3">
                <label class="w-20 text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('questions.options') }}</label>
                <button
                  class="inline-flex items-center gap-1 rounded-lg border border-gray-300 px-2.5 py-1 text-xs text-gray-600 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-base dark:text-gray-300 dark:hover:bg-dk-muted"
                  @click="addOption"
                >
                  <IconPlus :size="14" />
                  {{ t('questions.add_option') }}
                </button>
              </div>
              <div v-for="(opt, i) in form.options" :key="i" class="flex items-center gap-2 pl-20">
                <span class="text-sm font-semibold text-gray-500 dark:text-gray-400">{{ String.fromCharCode(65 + i) }}</span>
                <input
                  v-model="form.options[i]"
                  type="text"
                  :placeholder="t('questions.option_placeholder')"
                  class="flex-1 rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm text-gray-900 placeholder-gray-400 outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white dark:placeholder-gray-500"
                />
                <button
                  class="rounded p-1 text-gray-400 transition-colors hover:text-red-500 disabled:opacity-30"
                  :disabled="form.options.length <= 2"
                  @click="removeOption(i)"
                >
                  <IconX :size="15" />
                </button>
                <!-- 单选：选正确答案 -->
                <label
                  v-if="form.type === 'single'"
                  class="flex shrink-0 cursor-pointer items-center gap-1 text-xs text-gray-600 dark:text-gray-300"
                  :class="form.singleAnswer === i ? 'text-blue-600 dark:text-blue-400' : ''"
                >
                  <input
                    type="radio"
                    name="singleAnswer"
                    :checked="form.singleAnswer === i"
                    class="h-3.5 w-3.5"
                    @change="form.singleAnswer = i"
                  />
                  {{ t('questions.correct') }}
                </label>
                <!-- 多选：选正确答案 -->
                <label
                  v-if="form.type === 'multiple'"
                  class="flex shrink-0 cursor-pointer items-center gap-1 text-xs text-gray-600 dark:text-gray-300"
                  :class="form.multipleAnswer.includes(i) ? 'text-purple-600 dark:text-purple-400' : ''"
                >
                  <input
                    type="checkbox"
                    :checked="form.multipleAnswer.includes(i)"
                    class="h-3.5 w-3.5"
                    @change="toggleMultiple(i)"
                  />
                  {{ t('questions.correct') }}
                </label>
              </div>
            </div>

            <!-- Blank answers -->
            <div v-else-if="form.type === 'blank'" class="flex flex-col gap-2">
              <div class="flex items-center gap-3 pl-20">
                <p class="text-xs text-gray-400 dark:text-gray-500">{{ t('questions.blank_hint') }}</p>
                <button
                  class="inline-flex items-center gap-1 rounded-lg border border-gray-300 px-2.5 py-1 text-xs text-gray-600 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-base dark:text-gray-300 dark:hover:bg-dk-muted"
                  @click="addBlank"
                >
                  <IconPlus :size="14" />
                  {{ t('questions.add_blank') }}
                </button>
              </div>
              <div v-for="(_, bi) in form.blankAnswers" :key="bi" class="flex items-center gap-2 pl-20">
                <span class="text-xs font-semibold text-gray-500 dark:text-gray-400">[{{ bi + 1 }}]</span>
                <input
                  v-model="form.blankAnswers[bi]"
                  type="text"
                  :placeholder="t('questions.blank_placeholder')"
                  class="flex-1 rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm text-gray-900 placeholder-gray-400 outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white dark:placeholder-gray-500"
                />
                <button
                  class="rounded p-1 text-gray-400 transition-colors hover:text-red-500 disabled:opacity-30"
                  :disabled="form.blankAnswers.length <= 1"
                  @click="removeBlank(bi)"
                >
                  <IconX :size="15" />
                </button>
              </div>
            </div>

            <!-- Judge answer -->
            <div v-else class="flex items-center gap-3">
              <label class="w-20 text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('questions.answer') }}</label>
              <div class="flex gap-2">
                <button
                  class="rounded-lg border px-4 py-1.5 text-sm transition-colors"
                  :class="form.judgeAnswer
                    ? 'border-green-500 bg-green-50 text-green-700 dark:border-green-500 dark:bg-green-900/30 dark:text-green-300'
                    : 'border-gray-300 text-gray-600 hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-base dark:text-gray-300'"
                  @click="form.judgeAnswer = true"
                >
                  {{ t('quiz.true') }}
                </button>
                <button
                  class="rounded-lg border px-4 py-1.5 text-sm transition-colors"
                  :class="!form.judgeAnswer
                    ? 'border-red-500 bg-red-50 text-red-700 dark:border-red-500 dark:bg-red-900/30 dark:text-red-300'
                    : 'border-gray-300 text-gray-600 hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-base dark:text-gray-300'"
                  @click="form.judgeAnswer = false"
                >
                  {{ t('quiz.false') }}
                </button>
              </div>
            </div>

            <!-- 答案解析 -->
            <div class="flex items-start gap-3">
              <label class="w-20 pt-2 text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('questions.explain') }}</label>
              <textarea
                v-model="form.explain"
                rows="2"
                :placeholder="t('questions.explain_placeholder')"
                class="flex-1 resize-y rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 placeholder-gray-400 outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white dark:placeholder-gray-500"
              ></textarea>
            </div>

            <!-- Score + Active -->
            <div class="flex items-center gap-3">
              <label class="w-20 text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('questions.score') }}</label>
              <input
                v-model.number="form.score"
                type="number"
                min="1"
                max="100"
                class="w-24 rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm text-gray-900 outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white"
              />
              <label class="flex items-center gap-1.5 pl-4 text-sm text-gray-700 dark:text-gray-300">
                <input v-model="form.active" type="checkbox" class="h-3.5 w-3.5" />
                {{ t('questions.active') }}
              </label>
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
      :title="t('questions.delete_title')"
      :message="t('questions.delete_msg', { title: deleteTarget?.title || '' })"
      danger
      @confirm="doDelete"
    />
  </div>
</template>
