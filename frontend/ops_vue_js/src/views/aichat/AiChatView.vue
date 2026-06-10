<script setup>
import { nextTick, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { IconLoader2, IconPhoto, IconRobot, IconSend, IconTrash, IconUser, IconX } from '@tabler/icons-vue'
import { fetchOpenAIProfiles, streamChat } from '@/api/aichat'
import { usePageTitle } from '@/composables/usePageTitle'
import { useToastStore } from '@/stores/toast'

const { t } = useI18n()
const toast = useToastStore()

usePageTitle('appname.aichat')

const messages = ref([])
const inputText = ref('')
const selectedImage = ref(null)
const pending = ref(false)
const traces = ref([])
const reasoning = ref('')
const stats = ref(null)
const profiles = ref([])
const activeProfile = ref('')
const toolRouter = ref(null)
const messageListRef = ref(null)
const fileInputRef = ref(null)

const MAX_IMAGE_SIZE = 4 * 1024 * 1024
const ALLOWED_IMAGE_TYPES = ['image/jpeg', 'image/png', 'image/webp', 'image/gif']

onMounted(loadProfiles)

async function loadProfiles() {
  try {
    const res = await fetchOpenAIProfiles()
    if (res.errCode === 0 && res.data) {
      profiles.value = res.data.profiles || []
      activeProfile.value = res.data.active || profiles.value[0]?.name || ''
      toolRouter.value = res.data.toolRouter || null
    }
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error)
    toast.error(message)
  }
}

function scrollToBottom() {
  nextTick(() => {
    const el = messageListRef.value
    if (el) {
      el.scrollTop = el.scrollHeight
    }
  })
}

function onKeydown(event) {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    sendMessage()
  }
}

function clearChat() {
  if (pending.value) return
  messages.value = []
  traces.value = []
  reasoning.value = ''
  stats.value = null
  clearSelectedImage()
}

function triggerImagePicker() {
  if (pending.value) return
  fileInputRef.value?.click()
}

function onImageSelected(event) {
  const file = event.target.files?.[0]
  event.target.value = ''
  if (!file) return

  if (!ALLOWED_IMAGE_TYPES.includes(file.type)) {
    toast.error(t('aichat.image_type_error'))
    return
  }
  if (file.size > MAX_IMAGE_SIZE) {
    toast.error(t('aichat.image_size_error'))
    return
  }

  const reader = new FileReader()
  reader.onload = (loadEvent) => {
    selectedImage.value = {
      dataUrl: loadEvent.target?.result || '',
      name: file.name,
      size: file.size,
      type: file.type,
    }
  }
  reader.onerror = () => {
    toast.error(t('aichat.image_read_error'))
  }
  reader.readAsDataURL(file)
}

function clearSelectedImage() {
  selectedImage.value = null
}

function formatFileSize(size) {
  if (size >= 1024 * 1024) {
    return `${(size / 1024 / 1024).toFixed(1)} MB`
  }
  return `${Math.max(1, Math.round(size / 1024))} KB`
}

function messageImage(message) {
  return message.image_url || message.imageURL || ''
}

function formatTraceData(data) {
  if (!data) return []
  const parts = []
  if (data.database) parts.push(`${t('aichat.trace_database')}: ${data.database}`)
  if (data.sql) parts.push(data.sql)
  if (typeof data.rows === 'number') parts.push(`${t('aichat.trace_rows')}: ${data.rows}`)
  if (typeof data.columns === 'number') parts.push(`${t('aichat.trace_columns')}: ${data.columns}`)
  if (typeof data.count === 'number') parts.push(`${t('aichat.trace_count')}: ${data.count}`)
  if (Array.isArray(data.tools)) parts.push(`${t('aichat.trace_tools')}: ${data.tools.join(', ') || '-'}`)
  if (Array.isArray(data.selections) && data.selections.length) {
    parts.push(data.selections.map((item) => `${item.name}: ${item.reason || '-'}`).join('\n'))
  }
  if (data.reason) parts.push(`${t('aichat.trace_reason')}: ${data.reason}`)
  if (data.error) parts.push(`${t('aichat.trace_error')}: ${data.error}`)
  if (data.truncated) parts.push(t('aichat.trace_truncated'))
  if (data.max_rows) parts.push(`max_rows: ${data.max_rows}`)
  return parts
}

function formatFixed(value) {
  return typeof value === 'number' ? value.toFixed(1) : '0.0'
}

function formatTokenStats(value) {
  if (!value) return ''
  const toolTokens = (value.tool_prompt_tokens || 0) + (value.tool_completion_tokens || 0)
  const parts = [
    `${t('aichat.tokens_avg_speed')}: ${formatFixed(value.completion_tokens_per_sec)} tokens/sec`,
    `${t('aichat.tokens_peak_speed')}: ${formatFixed(value.peak_completion_tokens_per_sec)} tokens/sec`,
    `${t('aichat.tokens_total')}: ${value.total_tokens || 0}`,
    `${t('aichat.tokens_prompt')}: ${value.prompt_tokens || 0}`,
    `${t('aichat.tokens_completion')}: ${value.completion_tokens || 0}`,
  ]
  if (toolTokens) parts.push(`${t('aichat.tokens_tool')}: ${toolTokens}`)
  if (value.estimated) parts.push(t('aichat.tokens_estimated'))
  return parts.join(' ｜ ')
}

async function sendMessage() {
  const text = inputText.value.trim()
  const image = selectedImage.value
  if ((!text && !image) || pending.value) return

  inputText.value = ''
  clearSelectedImage()
  traces.value = []
  reasoning.value = ''
  stats.value = null

  const userMessage = { role: 'user', content: text }
  if (image) {
    userMessage.image_url = image.dataUrl
  }
  messages.value.push(userMessage)
  const assistantMessage = { role: 'assistant', content: '' }
  messages.value.push(assistantMessage)
  pending.value = true
  scrollToBottom()

  const history = messages.value
    .filter((message) => message.role === 'user' || message.role === 'assistant')
    .map((message) => {
      const item = { role: message.role, content: message.content || '' }
      if (message.image_url) item.image_url = message.image_url
      return item
    })
    .slice(0, -1)

  try {
    await streamChat(history, { openaiName: activeProfile.value }, {
      onDelta(delta) {
        assistantMessage.content += delta
        scrollToBottom()
      },
      onTrace(frame) {
        traces.value.push(frame)
        scrollToBottom()
      },
      onReasoning(delta) {
        reasoning.value += delta
        scrollToBottom()
      },
      onStats(value) {
        stats.value = value
      },
      onError(message) {
        if (!assistantMessage.content) {
          assistantMessage.content = t('aichat.error_prefix') + message
        }
        toast.error(message)
        scrollToBottom()
      },
    })
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error)
    assistantMessage.content = t('aichat.error_prefix') + message
    toast.error(message)
  } finally {
    pending.value = false
    scrollToBottom()
  }
}
</script>

<template>
  <div class="mx-auto flex h-[calc(100vh-7rem)] max-w-5xl flex-col px-4 py-6">
    <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-dk-text">
          {{ t('aichat.title') }}
        </h1>
        <p class="mt-1 text-sm text-gray-500 dark:text-dk-subtle">
          {{ t('aichat.subtitle') }}
        </p>
      </div>
      <div class="flex flex-wrap items-center gap-2">
        <button
          type="button"
          class="inline-flex items-center gap-2 rounded-md border border-gray-200 px-3 py-2 text-sm text-gray-600 transition-colors hover:bg-gray-100 disabled:cursor-not-allowed disabled:opacity-60 dark:border-dk-muted dark:text-dk-subtle dark:hover:bg-dk-card"
          :disabled="pending || messages.length === 0"
          @click="clearChat"
        >
          <IconTrash :size="16" />
          {{ t('aichat.clear') }}
        </button>
      </div>
    </div>

    <div class="flex min-h-0 flex-1 flex-col overflow-hidden rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card">
      <div ref="messageListRef" class="min-h-0 flex-1 overflow-y-auto px-4 py-5 sm:px-6">
        <div v-if="messages.length === 0" class="flex h-full items-center justify-center text-center">
          <div class="max-w-md">
            <div class="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-full bg-blue-50 text-blue-600 dark:bg-blue-900/20 dark:text-blue-300">
              <IconRobot :size="30" />
            </div>
            <h2 class="text-lg font-semibold text-gray-900 dark:text-dk-text">
              {{ t('aichat.empty_title') }}
            </h2>
            <p class="mt-2 text-sm text-gray-500 dark:text-dk-subtle">
              {{ t('aichat.empty_hint') }}
            </p>
          </div>
        </div>

        <div v-else class="space-y-5">
          <div
            v-for="(message, index) in messages"
            :key="index"
            :class="['flex gap-3', message.role === 'user' ? 'justify-end' : 'justify-start']"
          >
            <div v-if="message.role !== 'user'" class="mt-1 flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-blue-50 text-blue-600 dark:bg-blue-900/20 dark:text-blue-300">
              <IconRobot :size="18" />
            </div>

            <div :class="['max-w-[82%] rounded-2xl px-4 py-3 text-sm leading-6 shadow-sm', message.role === 'user' ? 'bg-blue-600 text-white' : 'border border-gray-200 bg-gray-50 text-gray-800 dark:border-dk-muted dark:bg-dk-base dark:text-dk-text']">
              <div v-if="message.role !== 'user' && index === messages.length - 1 && traces.length" class="mb-3 space-y-2">
                <div
                  v-for="(trace, traceIndex) in traces"
                  :key="traceIndex"
                  class="rounded-lg border border-blue-100 bg-blue-50 px-3 py-2 text-xs text-blue-700 dark:border-blue-900/40 dark:bg-blue-900/20 dark:text-blue-200"
                >
                  <div class="font-medium">
                    {{ trace.tool || 'tool' }} · {{ trace.status || trace.stage || 'trace' }}
                  </div>
                  <div v-if="trace.message" class="mt-1 opacity-90">
                    {{ trace.message }}
                  </div>
                  <div v-if="formatTraceData(trace.data).length" class="mt-2 space-y-1 rounded-md border border-blue-100 bg-white/70 px-2 py-1 font-mono text-[11px] leading-5 text-blue-800 dark:border-blue-900/40 dark:bg-dk-card/70 dark:text-blue-100">
                    <div v-for="(line, dataIndex) in formatTraceData(trace.data)" :key="dataIndex" class="whitespace-pre-wrap break-words">
                      {{ line }}
                    </div>
                  </div>
                </div>
              </div>

              <div v-if="message.role !== 'user' && index === messages.length - 1 && reasoning" class="mb-3 rounded-lg border border-purple-100 bg-purple-50 px-3 py-2 text-xs text-purple-800 dark:border-purple-900/40 dark:bg-purple-900/20 dark:text-purple-100">
                <div class="font-medium">
                  {{ t('aichat.reasoning') }}
                </div>
                <div class="mt-1 whitespace-pre-wrap break-words">
                  {{ reasoning }}
                </div>
              </div>

              <img
                v-if="messageImage(message)"
                :src="messageImage(message)"
                :alt="message.content || t('aichat.attach_image')"
                class="mb-2 max-h-64 max-w-full rounded-lg object-contain"
              />

              <p v-if="message.content || (message.role === 'assistant' && pending)" class="whitespace-pre-wrap break-words">
                {{ message.content || (message.role === 'assistant' && pending ? t('aichat.thinking') : '') }}
              </p>

              <div v-if="message.role !== 'user' && index === messages.length - 1 && pending" class="mt-2 inline-flex items-center gap-1 text-xs text-gray-500 dark:text-dk-subtle">
                <IconLoader2 :size="14" class="animate-spin" />
                {{ t('aichat.streaming') }}
              </div>

              <div v-if="message.role !== 'user' && index === messages.length - 1 && stats" class="mt-3 text-xs text-gray-500 dark:text-dk-subtle">
                {{ formatTokenStats(stats) }}
              </div>
            </div>

            <div v-if="message.role === 'user'" class="mt-1 flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-gray-100 text-gray-600 dark:bg-dk-muted dark:text-dk-text">
              <IconUser :size="18" />
            </div>
          </div>
        </div>
      </div>

      <div class="border-t border-gray-200 bg-gray-50 p-4 dark:border-dk-muted dark:bg-dk-base">
        <input
          ref="fileInputRef"
          type="file"
          accept="image/jpeg,image/png,image/webp,image/gif"
          class="hidden"
          :disabled="pending"
          @change="onImageSelected"
        />

        <div v-if="selectedImage" class="mb-3 flex items-center gap-3 rounded-lg border border-gray-200 bg-white p-2 dark:border-dk-muted dark:bg-dk-card">
          <img :src="selectedImage.dataUrl" :alt="selectedImage.name" class="h-14 w-14 rounded object-cover" />
          <div class="min-w-0 flex-1">
            <div class="truncate text-sm font-medium text-gray-800 dark:text-dk-text">
              {{ selectedImage.name }}
            </div>
            <div class="text-xs text-gray-500 dark:text-dk-subtle">
              {{ formatFileSize(selectedImage.size) }}
            </div>
          </div>
          <button
            type="button"
            class="inline-flex h-8 w-8 items-center justify-center rounded-md text-gray-500 transition-colors hover:bg-gray-100 disabled:cursor-not-allowed disabled:opacity-60 dark:text-dk-subtle dark:hover:bg-dk-muted"
            :title="t('aichat.remove_image')"
            :disabled="pending"
            @click="clearSelectedImage"
          >
            <IconX :size="16" />
          </button>
        </div>

        <div class="flex items-end gap-3">
          <button
            type="button"
            class="inline-flex h-[52px] w-[52px] shrink-0 items-center justify-center rounded-lg border border-gray-300 bg-white text-gray-600 transition-colors hover:bg-gray-100 disabled:cursor-not-allowed disabled:opacity-60 dark:border-dk-muted dark:bg-dk-card dark:text-dk-subtle dark:hover:bg-dk-muted"
            :title="t('aichat.attach_image')"
            :disabled="pending"
            @click="triggerImagePicker"
          >
            <IconPhoto :size="20" />
          </button>
          <textarea
            v-model="inputText"
            rows="2"
            class="min-h-[52px] flex-1 resize-none rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 disabled:cursor-not-allowed disabled:opacity-60 dark:border-dk-muted dark:bg-dk-card dark:text-dk-text"
            :placeholder="t('aichat.input_placeholder')"
            :disabled="pending"
            @keydown="onKeydown"
          />
          <button
            type="button"
            class="inline-flex h-[52px] items-center gap-2 rounded-lg bg-blue-600 px-4 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="pending || (!inputText.trim() && !selectedImage)"
            @click="sendMessage"
          >
            <IconLoader2 v-if="pending" :size="18" class="animate-spin" />
            <IconSend v-else :size="18" />
            {{ t('aichat.send') }}
          </button>
        </div>
        <p class="mt-2 text-xs text-gray-500 dark:text-dk-subtle">
          {{ t('aichat.enter_hint') }}
        </p>
      </div>
    </div>
  </div>
</template>
