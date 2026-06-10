<script setup>
import { nextTick, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { IconRobot, IconSend, IconTrash, IconUser, IconLoader2 } from '@tabler/icons-vue'
import { fetchOpenAIProfiles, streamChat } from '@/api/aichat'
import { usePageTitle } from '@/composables/usePageTitle'
import { useToastStore } from '@/stores/toast'

const { t } = useI18n()
const toast = useToastStore()

usePageTitle('appname.aichat')

const messages = ref([])
const inputText = ref('')
const pending = ref(false)
const traces = ref([])
const stats = ref(null)
const profiles = ref([])
const activeProfile = ref('')
const toolRouter = ref(null)
const messageListRef = ref(null)

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
  stats.value = null
}

async function sendMessage() {
  const text = inputText.value.trim()
  if (!text || pending.value) return

  inputText.value = ''
  traces.value = []
  stats.value = null

  messages.value.push({ role: 'user', content: text })
  const assistantMessage = { role: 'assistant', content: '' }
  messages.value.push(assistantMessage)
  pending.value = true
  scrollToBottom()

  const history = messages.value
    .filter((message) => message.role === 'user' || message.role === 'assistant')
    .map((message) => ({ role: message.role, content: message.content }))
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
                </div>
              </div>

              <p class="whitespace-pre-wrap break-words">
                {{ message.content || (message.role === 'assistant' && pending ? t('aichat.thinking') : '') }}
              </p>

              <div v-if="message.role !== 'user' && index === messages.length - 1 && pending" class="mt-2 inline-flex items-center gap-1 text-xs text-gray-500 dark:text-dk-subtle">
                <IconLoader2 :size="14" class="animate-spin" />
                {{ t('aichat.streaming') }}
              </div>

              <div v-if="message.role !== 'user' && index === messages.length - 1 && stats" class="mt-3 text-xs text-gray-500 dark:text-dk-subtle">
                {{ t('aichat.tokens') }}: {{ stats.total_tokens || 0 }}
              </div>
            </div>

            <div v-if="message.role === 'user'" class="mt-1 flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-gray-100 text-gray-600 dark:bg-dk-muted dark:text-dk-text">
              <IconUser :size="18" />
            </div>
          </div>
        </div>
      </div>

      <div class="border-t border-gray-200 bg-gray-50 p-4 dark:border-dk-muted dark:bg-dk-base">
        <div class="flex items-end gap-3">
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
            :disabled="pending || !inputText.trim()"
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
