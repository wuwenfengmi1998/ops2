<script setup>
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { IconCloud, IconDeviceFloppy, IconLoader2, IconPhoto, IconPlus, IconRobot, IconSend, IconTrash, IconUser, IconX } from '@tabler/icons-vue'
import {
  deleteAIChatConversation,
  fetchAIChatConversation,
  fetchAIChatConversations,
  fetchOpenAIProfiles,
  streamChat,
} from '@/api/aichat'
import { usePageTitle } from '@/composables/usePageTitle'
import { useToastStore } from '@/stores/toast'
import { useUserStore } from '@/stores/user'

const props = defineProps({
  embedded: {
    type: Boolean,
    default: false,
  },
})

const { t } = useI18n()
const toast = useToastStore()
const userStore = useUserStore()

if (!props.embedded) {
  usePageTitle('appname.aichat')
}

const messages = ref([])
const inputText = ref('')
const selectedImage = ref(null)
const pending = ref(false)
const traces = ref([])
const tracesCollapsed = ref(true)
const reasoning = ref('')
const reasoningCollapsed = ref(true)
const stats = ref(null)
const profiles = ref([])
const activeProfile = ref('')
const toolRouter = ref(null)
const messageListRef = ref(null)
const fileInputRef = ref(null)

const localConversations = ref([])
const serverConversations = ref([])
const activeSource = ref('local')
const activeLocalId = ref('')
const activeServerId = ref(0)
const loadingConversations = ref(false)

const MAX_IMAGE_SIZE = 4 * 1024 * 1024
const ALLOWED_IMAGE_TYPES = ['image/jpeg', 'image/png', 'image/webp', 'image/gif']
const LOCAL_STORAGE_KEY = 'ops:aichat:local:v1'

const sortedLocalConversations = computed(() => sortConversations(localConversations.value))
const sortedServerConversations = computed(() => sortConversations(serverConversations.value))

onMounted(async () => {
  loadLocalConversations()
  if (userStore.isLoggedIn) {
    activeSource.value = 'server'
    activeLocalId.value = createLocalId()
  } else if (localConversations.value.length === 0) {
    createLocalConversation(false)
  } else {
    selectLocalConversation(localConversations.value[0].localId)
  }
  await Promise.all([loadProfiles(), loadServerConversations()])
})

watch(() => userStore.isLoggedIn, async (loggedIn) => {
  serverConversations.value = []
  activeServerId.value = 0
  if (loggedIn) {
    activeSource.value = 'server'
    activeLocalId.value = createLocalId()
    await loadServerConversations()
  } else if (!activeLocalId.value && localConversations.value.length) {
    selectLocalConversation(localConversations.value[0].localId)
  }
})

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

async function loadServerConversations() {
  if (!userStore.isLoggedIn) return
  loadingConversations.value = true
  try {
    const res = await fetchAIChatConversations({ page: 1, page_size: 100 })
    if (res.errCode === 0 && res.data) {
      serverConversations.value = res.data.conversations || []
    } else {
      toast.error(t('aichat.load_conversations_failed'))
    }
  } catch (error) {
    toast.error(error instanceof Error ? error.message : t('aichat.load_conversations_failed'))
  } finally {
    loadingConversations.value = false
  }
}

function loadLocalConversations() {
  try {
    const raw = localStorage.getItem(LOCAL_STORAGE_KEY)
    const parsed = raw ? JSON.parse(raw) : []
    localConversations.value = Array.isArray(parsed) ? parsed.map(normalizeLocalConversation) : []
  } catch {
    localConversations.value = []
  }
}

function saveLocalConversations() {
  try {
    localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(localConversations.value))
    return true
  } catch {
    toast.error(t('aichat.storage_full'))
    return false
  }
}

function normalizeLocalConversation(conversation) {
  const now = new Date().toISOString()
  return {
    localId: conversation.localId || createLocalId(),
    serverId: Number(conversation.serverId || 0),
    title: conversation.title || t('aichat.untitled_chat'),
    openaiName: conversation.openaiName || '',
    messages: Array.isArray(conversation.messages) ? conversation.messages : [],
    createdAt: conversation.createdAt || now,
    updatedAt: conversation.updatedAt || now,
  }
}

function createLocalId() {
  const suffix = globalThis.crypto?.randomUUID?.() || `${Date.now()}-${Math.random().toString(36).slice(2)}`
  return `local-${suffix}`
}

function sortConversations(items) {
  return [...items].sort((a, b) => new Date(b.updatedAt || b.lastMessageAt || 0) - new Date(a.updatedAt || a.lastMessageAt || 0))
}

function createLocalConversation(select = true) {
  const now = new Date().toISOString()
  const conversation = {
    localId: createLocalId(),
    serverId: 0,
    title: t('aichat.untitled_chat'),
    openaiName: activeProfile.value || '',
    messages: [],
    createdAt: now,
    updatedAt: now,
  }
  localConversations.value.unshift(conversation)
  saveLocalConversations()
  if (select) {
    selectLocalConversation(conversation.localId)
  }
  return conversation
}

function newChat() {
  if (pending.value) return
  if (userStore.isLoggedIn) {
    activeSource.value = 'server'
    activeServerId.value = 0
    activeLocalId.value = createLocalId()
    messages.value = []
    resetStreamDetails()
    scrollToBottom()
    return
  }

  const emptyConversation = sortedLocalConversations.value.find((conversation) => isEmptyLocalConversation(conversation))
  if (emptyConversation) {
    selectLocalConversation(emptyConversation.localId)
    return
  }
  createLocalConversation(true)
}

function isEmptyLocalConversation(conversation) {
  return !Array.isArray(conversation.messages) || conversation.messages.length === 0
}

function selectLocalConversation(localId) {
  if (pending.value) return
  const conversation = localConversations.value.find((item) => item.localId === localId)
  if (!conversation) return
  activeSource.value = 'local'
  activeLocalId.value = conversation.localId
  activeServerId.value = Number(conversation.serverId || 0)
  messages.value = cloneMessages(conversation.messages)
  activeProfile.value = conversation.openaiName || activeProfile.value
  resetStreamDetails()
  scrollToBottom()
}

async function selectServerConversation(id) {
  if (pending.value) return
  const serverId = Number(id)
  activeSource.value = 'server'
  activeServerId.value = serverId
  resetStreamDetails()
  try {
    const res = await fetchAIChatConversation(serverId)
    if (res.errCode === 0 && res.data) {
      const conversation = res.data.conversation || {}
      messages.value = cloneMessages(res.data.messages || [])
      activeLocalId.value = findLocalByServer(conversation.id, conversation.clientLocalId)?.localId || ''
      activeProfile.value = conversation.openaiName || activeProfile.value
      scrollToBottom()
    } else {
      toast.error(t('aichat.load_chat_failed'))
    }
  } catch (error) {
    toast.error(error instanceof Error ? error.message : t('aichat.load_chat_failed'))
  }
}

function cloneMessages(items) {
  return items.map((item) => ({ ...item }))
}

function findLocalByServer(serverId, clientLocalId = '') {
  return localConversations.value.find((item) => Number(item.serverId || 0) === Number(serverId))
    || localConversations.value.find((item) => clientLocalId && item.localId === clientLocalId)
}

function ensureActiveLocalConversation() {
  let conversation = localConversations.value.find((item) => item.localId === activeLocalId.value)
  if (conversation) return conversation

  conversation = createLocalConversation(false)
  conversation.messages = cloneMessages(messages.value)
  activeLocalId.value = conversation.localId
  saveLocalConversations()
  return conversation
}

function updateActiveLocalConversation() {
  if (userStore.isLoggedIn) return
  const conversation = ensureActiveLocalConversation()
  const now = new Date().toISOString()
  conversation.messages = cloneMessages(messages.value)
  conversation.title = makeTitle(conversation.messages)
  conversation.openaiName = activeProfile.value || conversation.openaiName || ''
  conversation.updatedAt = now
  saveLocalConversations()
}

function updateServerConversationList(conversation) {
  if (!conversation?.id) return
  const index = serverConversations.value.findIndex((item) => Number(item.id) === Number(conversation.id))
  const next = {
    ...(index >= 0 ? serverConversations.value[index] : {}),
    ...conversation,
    updatedAt: conversation.updatedAt || new Date().toISOString(),
  }
  if (index >= 0) {
    serverConversations.value.splice(index, 1, next)
  } else {
    serverConversations.value.unshift(next)
  }
}

function bindServerConversation(conversation) {
  if (!conversation?.id) return
  activeServerId.value = Number(conversation.id)
  updateServerConversationList(conversation)
}

function resetStreamDetails() {
  traces.value = []
  tracesCollapsed.value = true
  reasoning.value = ''
  reasoningCollapsed.value = true
  stats.value = null
  clearSelectedImage()
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
  resetStreamDetails()
  updateActiveLocalConversation()
}

async function deleteLocalConversation(localId) {
  if (pending.value) return
  localConversations.value = localConversations.value.filter((item) => item.localId !== localId)
  saveLocalConversations()
  if (activeLocalId.value === localId) {
    if (localConversations.value.length) {
      selectLocalConversation(localConversations.value[0].localId)
    } else {
      createLocalConversation(true)
    }
  }
}

async function deleteServerConversation(id) {
  if (pending.value) return
  if (!window.confirm(t('aichat.delete_chat'))) return
  try {
    const res = await deleteAIChatConversation(id)
    if (res.errCode === 0) {
      serverConversations.value = serverConversations.value.filter((item) => Number(item.id) !== Number(id))
      localConversations.value.forEach((item) => {
        if (Number(item.serverId || 0) === Number(id)) item.serverId = 0
      })
      saveLocalConversations()
      if (Number(activeServerId.value) === Number(id)) {
        activeServerId.value = 0
        activeSource.value = 'local'
      }
    } else {
      toast.error(t('aichat.delete_chat_failed'))
    }
  } catch (error) {
    toast.error(error instanceof Error ? error.message : t('aichat.delete_chat_failed'))
  }
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

function messageTraces(message, index) {
  if (Array.isArray(message.traces) && message.traces.length) return message.traces
  if (message.role !== 'user' && index === messages.value.length - 1) return traces.value
  return []
}

function messageReasoning(message, index) {
  if (message.reasoning) return message.reasoning
  if (message.role !== 'user' && index === messages.value.length - 1) return reasoning.value
  return ''
}

function messageStats(message, index) {
  if (message.stats) return message.stats
  if (message.role !== 'user' && index === messages.value.length - 1) return stats.value
  return null
}

function formatTraceData(data) {
  if (!data) return []
  const parts = []
  const stringify = (value) => {
    if (typeof value === 'string') return value
    try {
      return JSON.stringify(value, null, 2)
    } catch (e) {
      return String(value)
    }
  }
  if (data.database) parts.push(`${t('aichat.trace_database')}: ${data.database}`)
  if (data.sql) parts.push(data.sql)
  if (data.tool) parts.push(`tool: ${data.tool}`)
  if (typeof data.round === 'number') parts.push(`round: ${data.round}`)
  if (typeof data.rows === 'number') parts.push(`${t('aichat.trace_rows')}: ${data.rows}`)
  if (typeof data.columns === 'number') parts.push(`${t('aichat.trace_columns')}: ${data.columns}`)
  if (typeof data.count === 'number') parts.push(`${t('aichat.trace_count')}: ${data.count}`)
  if (Array.isArray(data.tools)) parts.push(`${t('aichat.trace_tools')}: ${data.tools.join(', ') || '-'}`)
  if (Array.isArray(data.selections) && data.selections.length) {
    parts.push(data.selections.map((item) => `${item.name}: ${item.reason || '-'}`).join('\n'))
  }
  if (data.arguments !== undefined) parts.push(`arguments:\n${stringify(data.arguments)}`)
  if (data.result !== undefined) parts.push(`result:\n${stringify(data.result)}`)
  if (data.start_date) parts.push(`start_date: ${data.start_date}`)
  if (data.end_date) parts.push(`end_date: ${data.end_date}`)
  if (data.limit) parts.push(`limit: ${data.limit}`)
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

function formatDate(value) {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return date.toLocaleString()
}

function makeTitle(items) {
  const firstUser = items.find((item) => item.role === 'user' && (item.content || messageImage(item)))
  if (!firstUser) return t('aichat.untitled_chat')
  const text = (firstUser.content || t('aichat.attach_image')).trim()
  return text.length > 40 ? `${text.slice(0, 40)}...` : text
}

const markdownCache = new Map()

function escapeHtml(value) {
  return String(value)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

function renderInlineMarkdown(value) {
  return escapeHtml(value)
    .replace(/`([^`]+)`/g, '<code>$1</code>')
    .replace(/!\[([^\]]*)\]\((https?:\/\/[^\s)]+)\)/g, '<a href="$2" target="_blank" rel="noopener noreferrer">$1</a>')
    .replace(/(?<!!)\[([^\]]+)\]\((https?:\/\/[^\s)]+)\)/g, '<a href="$2" target="_blank" rel="noopener noreferrer">$1</a>')
    .replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
    .replace(/__([^_]+)__/g, '<strong>$1</strong>')
    .replace(/~~([^~]+)~~/g, '<del>$1</del>')
    .replace(/(^|\s)\*([^*]+)\*(?=\s|$)/g, '$1<em>$2</em>')
    .replace(/(^|\s)_([^_]+)_(?=\s|$)/g, '$1<em>$2</em>')
}

function renderMarkdown(value) {
  const source = String(value || '')
  if (!source) return ''
  const cached = markdownCache.get(source)
  if (cached) return cached

  const lines = source.replace(/\r\n/g, '\n').split('\n')
  const html = []
  let paragraph = []

  const flushParagraph = () => {
    if (!paragraph.length) return
    html.push(`<p>${renderInlineMarkdown(paragraph.join('\n')).replace(/\n/g, '<br>')}</p>`)
    paragraph = []
  }

  for (let index = 0; index < lines.length; index += 1) {
    const line = lines[index]
    const fence = line.match(/^\s*(```|~~~)\s*(.*)$/)
    if (fence) {
      flushParagraph()
      const fenceMarker = fence[1]
      const language = fence[2].trim()
      const code = []
      index += 1
      while (index < lines.length && !lines[index].trim().startsWith(fenceMarker)) {
        code.push(lines[index])
        index += 1
      }
      const languageClass = language ? ` class="language-${escapeHtml(language)}"` : ''
      html.push(`<pre><code${languageClass}>${escapeHtml(code.join('\n'))}</code></pre>`)
      continue
    }

    if (!line.trim()) {
      flushParagraph()
      continue
    }

    const heading = line.match(/^(#{1,6})\s+(.+)$/)
    if (heading) {
      flushParagraph()
      const level = heading[1].length
      html.push(`<h${level}>${renderInlineMarkdown(heading[2])}</h${level}>`)
      continue
    }

    if (/^\s*[-*_]{3,}\s*$/.test(line)) {
      flushParagraph()
      html.push('<hr>')
      continue
    }

    const quote = line.match(/^>\s?(.*)$/)
    if (quote) {
      flushParagraph()
      const quotes = [quote[1]]
      while (index + 1 < lines.length && /^>\s?/.test(lines[index + 1])) {
        index += 1
        quotes.push(lines[index].replace(/^>\s?/, ''))
      }
      html.push(`<blockquote>${renderInlineMarkdown(quotes.join('\n')).replace(/\n/g, '<br>')}</blockquote>`)
      continue
    }

    const unordered = line.match(/^\s*[-*+]\s+(.+)$/)
    if (unordered) {
      flushParagraph()
      const items = [unordered[1]]
      while (index + 1 < lines.length && /^\s*[-*+]\s+/.test(lines[index + 1])) {
        index += 1
        items.push(lines[index].replace(/^\s*[-*+]\s+/, ''))
      }
      html.push(`<ul>${items.map((item) => `<li>${renderInlineMarkdown(item)}</li>`).join('')}</ul>`)
      continue
    }

    const ordered = line.match(/^\s*\d+[.)]\s+(.+)$/)
    if (ordered) {
      flushParagraph()
      const items = [ordered[1]]
      while (index + 1 < lines.length && /^\s*\d+[.)]\s+/.test(lines[index + 1])) {
        index += 1
        items.push(lines[index].replace(/^\s*\d+[.)]\s+/, ''))
      }
      html.push(`<ol>${items.map((item) => `<li>${renderInlineMarkdown(item)}</li>`).join('')}</ol>`)
      continue
    }

    paragraph.push(line)
  }

  flushParagraph()
  const rendered = html.join('')
  if (markdownCache.size > 200) markdownCache.clear()
  markdownCache.set(source, rendered)
  return rendered
}

async function sendMessage() {
  const text = inputText.value.trim()
  const image = selectedImage.value
  if ((!text && !image) || pending.value) return

  const clientLocalId = activeLocalId.value || createLocalId()
  activeLocalId.value = clientLocalId
  inputText.value = ''
  clearSelectedImage()
  traces.value = []
  tracesCollapsed.value = true
  reasoning.value = ''
  reasoningCollapsed.value = true
  stats.value = null

  const userMessage = { role: 'user', content: text }
  if (image) {
    userMessage.image_url = image.dataUrl
  }
  messages.value.push(userMessage)
  const assistantMessage = { role: 'assistant', content: '' }
  messages.value.push(assistantMessage)
  pending.value = true
  updateActiveLocalConversation()
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
    await streamChat(history, {
      openaiName: activeProfile.value,
      conversationId: activeServerId.value || 0,
      clientLocalId,
      saveToServer: userStore.isLoggedIn,
    }, {
      onConversation(conversation) {
        bindServerConversation(conversation)
      },
      onDelta(delta) {
        assistantMessage.content += delta
        updateActiveLocalConversation()
        scrollToBottom()
      },
      onTrace(frame) {
        traces.value.push(frame)
        assistantMessage.traces = [...traces.value]
        updateActiveLocalConversation()
        scrollToBottom()
      },
      onReasoning(delta) {
        reasoning.value += delta
        assistantMessage.reasoning = reasoning.value
        updateActiveLocalConversation()
        scrollToBottom()
      },
      onStats(value) {
        stats.value = value
        assistantMessage.stats = value
        updateActiveLocalConversation()
      },
      onError(message) {
        if (!assistantMessage.content) {
          assistantMessage.content = t('aichat.error_prefix') + message
        }
        updateActiveLocalConversation()
        toast.error(message)
        scrollToBottom()
      },
    })
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error)
    assistantMessage.content = t('aichat.error_prefix') + message
    updateActiveLocalConversation()
    toast.error(message)
  } finally {
    pending.value = false
    if (traces.value.length) {
      tracesCollapsed.value = true
    }
    if (reasoning.value) {
      reasoningCollapsed.value = true
    }
    updateActiveLocalConversation()
    if (userStore.isLoggedIn) {
      loadServerConversations()
    }
    scrollToBottom()
  }
}
</script>

<template>
  <div :class="embedded ? 'flex h-[70vh] min-h-[38rem] flex-col' : 'mx-auto flex h-[calc(100vh-7rem)] max-w-7xl flex-col px-4 py-6'">
    <div :class="['mb-4 flex flex-wrap items-center justify-between gap-3', { 'px-1': embedded }]">
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
          class="inline-flex items-center gap-2 rounded-md bg-blue-600 px-3 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-60"
          :disabled="pending"
          @click="newChat"
        >
          <IconPlus :size="16" />
          {{ t('aichat.new_chat') }}
        </button>
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

    <div class="grid min-h-0 flex-1 gap-4 lg:grid-cols-[18rem_minmax(0,1fr)]">
      <aside class="flex min-h-0 flex-col overflow-hidden rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card">
        <div class="border-b border-gray-200 p-4 dark:border-dk-muted">
          <button
            type="button"
            class="flex w-full items-center justify-center gap-2 rounded-lg bg-blue-600 px-3 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="pending"
            @click="newChat"
          >
            <IconPlus :size="16" />
            {{ t('aichat.new_chat') }}
          </button>
          <p v-if="!userStore.isLoggedIn" class="mt-3 text-xs text-gray-500 dark:text-dk-subtle">
            {{ t('aichat.login_to_sync') }}
          </p>
        </div>

        <div class="min-h-0 flex-1 overflow-y-auto p-3">
          <section v-if="userStore.isLoggedIn" class="mb-5">
            <div class="mb-2 flex items-center justify-between text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dk-subtle">
              <span>{{ t('aichat.server_chats') }}</span>
              <IconLoader2 v-if="loadingConversations" :size="14" class="animate-spin" />
            </div>
            <div v-if="sortedServerConversations.length === 0" class="rounded-lg border border-dashed border-gray-200 p-3 text-xs text-gray-500 dark:border-dk-muted dark:text-dk-subtle">
              {{ t('aichat.no_server_chats') }}
            </div>
            <div v-else class="space-y-2">
              <button
                v-for="conversation in sortedServerConversations"
                :key="conversation.id"
                type="button"
                :class="[
                  'group flex w-full items-start gap-2 rounded-lg border px-3 py-2 text-left transition-colors',
                  activeSource === 'server' && Number(activeServerId) === Number(conversation.id)
                    ? 'border-blue-200 bg-blue-50 dark:border-blue-900/50 dark:bg-blue-900/20'
                    : 'border-transparent hover:bg-gray-50 dark:hover:bg-dk-base',
                ]"
                @click="selectServerConversation(conversation.id)"
              >
                <IconCloud :size="16" class="mt-0.5 shrink-0 text-blue-600 dark:text-blue-300" />
                <span class="min-w-0 flex-1">
                  <span class="block truncate text-sm font-medium text-gray-800 dark:text-dk-text">
                    {{ conversation.title || t('aichat.untitled_chat') }}
                  </span>
                  <span class="mt-1 block truncate text-xs text-gray-500 dark:text-dk-subtle">
                    {{ formatDate(conversation.lastMessageAt || conversation.updatedAt) }}
                  </span>
                </span>
                <span
                  class="rounded p-1 text-gray-400 opacity-0 transition hover:bg-red-50 hover:text-red-600 group-hover:opacity-100 dark:hover:bg-red-900/20"
                  :title="t('aichat.delete_chat')"
                  @click.stop="deleteServerConversation(conversation.id)"
                >
                  <IconTrash :size="14" />
                </span>
              </button>
            </div>
          </section>

          <section>
            <div class="mb-2 text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dk-subtle">
              {{ t('aichat.browser_chats') }}
            </div>
            <div v-if="sortedLocalConversations.length === 0" class="rounded-lg border border-dashed border-gray-200 p-3 text-xs text-gray-500 dark:border-dk-muted dark:text-dk-subtle">
              {{ t('aichat.no_browser_chats') }}
            </div>
            <div v-else class="space-y-2">
              <button
                v-for="conversation in sortedLocalConversations"
                :key="conversation.localId"
                type="button"
                :class="[
                  'group flex w-full items-start gap-2 rounded-lg border px-3 py-2 text-left transition-colors',
                  activeLocalId === conversation.localId && activeSource === 'local'
                    ? 'border-blue-200 bg-blue-50 dark:border-blue-900/50 dark:bg-blue-900/20'
                    : 'border-transparent hover:bg-gray-50 dark:hover:bg-dk-base',
                ]"
                @click="selectLocalConversation(conversation.localId)"
              >
                <IconDeviceFloppy :size="16" class="mt-0.5 shrink-0 text-gray-500 dark:text-dk-subtle" />
                <span class="min-w-0 flex-1">
                  <span class="block truncate text-sm font-medium text-gray-800 dark:text-dk-text">
                    {{ conversation.title || t('aichat.untitled_chat') }}
                  </span>
                  <span class="mt-1 block truncate text-xs text-gray-500 dark:text-dk-subtle">
                    {{ formatDate(conversation.updatedAt) }}
                  </span>
                </span>
                <span
                  class="rounded p-1 text-gray-400 opacity-0 transition hover:bg-red-50 hover:text-red-600 group-hover:opacity-100 dark:hover:bg-red-900/20"
                  :title="t('aichat.delete_chat')"
                  @click.stop="deleteLocalConversation(conversation.localId)"
                >
                  <IconTrash :size="14" />
                </span>
              </button>
            </div>
          </section>
        </div>
      </aside>

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
                <div v-if="message.role !== 'user' && messageTraces(message, index).length" class="mb-3 rounded-lg border border-blue-100 bg-blue-50 px-3 py-2 text-xs text-blue-700 dark:border-blue-900/40 dark:bg-blue-900/20 dark:text-blue-200">
                  <button
                    type="button"
                    class="flex w-full items-center justify-between gap-2 text-left font-medium"
                    @click="tracesCollapsed = !tracesCollapsed"
                  >
                    <span>{{ t('aichat.trace_details') }}</span>
                    <span class="text-[11px] opacity-75">
                      {{ tracesCollapsed ? t('aichat.expand') : t('aichat.collapse') }}
                    </span>
                  </button>
                  <div v-if="!tracesCollapsed" class="mt-2 space-y-2">
                    <div
                      v-for="(trace, traceIndex) in messageTraces(message, index)"
                      :key="traceIndex"
                      class="rounded-lg border border-blue-100 bg-white/70 px-3 py-2 dark:border-blue-900/40 dark:bg-dk-card/70"
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
                </div>

                <div v-if="message.role !== 'user' && messageReasoning(message, index)" class="mb-3 rounded-lg border border-purple-100 bg-purple-50 px-3 py-2 text-xs text-purple-800 dark:border-purple-900/40 dark:bg-purple-900/20 dark:text-purple-100">
                  <button
                    type="button"
                    class="flex w-full items-center justify-between gap-2 text-left font-medium"
                    @click="reasoningCollapsed = !reasoningCollapsed"
                  >
                    <span>{{ t('aichat.reasoning') }}</span>
                    <span class="text-[11px] opacity-75">
                      {{ reasoningCollapsed ? t('aichat.expand') : t('aichat.collapse') }}
                    </span>
                  </button>
                  <div v-if="!reasoningCollapsed" class="mt-1 whitespace-pre-wrap break-words">
                    {{ messageReasoning(message, index) }}
                  </div>
                </div>

                <img
                  v-if="messageImage(message)"
                  :src="messageImage(message)"
                  :alt="message.content || t('aichat.attach_image')"
                  class="mb-2 max-h-64 max-w-full rounded-lg object-contain"
                />

                <div
                  v-if="message.role !== 'user' && (message.content || (pending && index === messages.length - 1))"
                  class="markdown-body break-words"
                  v-html="renderMarkdown(message.content || (pending ? t('aichat.thinking') : ''))"
                />
                <p v-else-if="message.content" class="whitespace-pre-wrap break-words">
                  {{ message.content }}
                </p>

                <div v-if="message.role !== 'user' && index === messages.length - 1 && pending" class="mt-2 inline-flex items-center gap-1 text-xs text-gray-500 dark:text-dk-subtle">
                  <IconLoader2 :size="14" class="animate-spin" />
                  {{ t('aichat.streaming') }}
                </div>

                <div v-if="message.role !== 'user' && messageStats(message, index)" class="mt-3 text-xs text-gray-500 dark:text-dk-subtle">
                  {{ formatTokenStats(messageStats(message, index)) }}
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
  </div>
</template>

<style scoped>
.markdown-body :deep(*) {
  overflow-wrap: anywhere;
}

.markdown-body :deep(p) {
  margin: 0.25rem 0;
  white-space: normal;
}

.markdown-body :deep(p:first-child),
.markdown-body :deep(ul:first-child),
.markdown-body :deep(ol:first-child),
.markdown-body :deep(pre:first-child),
.markdown-body :deep(blockquote:first-child),
.markdown-body :deep(h1:first-child),
.markdown-body :deep(h2:first-child),
.markdown-body :deep(h3:first-child),
.markdown-body :deep(h4:first-child),
.markdown-body :deep(h5:first-child),
.markdown-body :deep(h6:first-child) {
  margin-top: 0;
}

.markdown-body :deep(p:last-child),
.markdown-body :deep(ul:last-child),
.markdown-body :deep(ol:last-child),
.markdown-body :deep(pre:last-child),
.markdown-body :deep(blockquote:last-child) {
  margin-bottom: 0;
}

.markdown-body :deep(h1),
.markdown-body :deep(h2),
.markdown-body :deep(h3),
.markdown-body :deep(h4),
.markdown-body :deep(h5),
.markdown-body :deep(h6) {
  margin: 0.75rem 0 0.35rem;
  font-weight: 700;
  line-height: 1.35;
}

.markdown-body :deep(h1) {
  font-size: 1.25rem;
}

.markdown-body :deep(h2) {
  font-size: 1.15rem;
}

.markdown-body :deep(h3),
.markdown-body :deep(h4),
.markdown-body :deep(h5),
.markdown-body :deep(h6) {
  font-size: 1rem;
}

.markdown-body :deep(ul),
.markdown-body :deep(ol) {
  margin: 0.4rem 0;
  padding-left: 1.25rem;
}

.markdown-body :deep(ul) {
  list-style: disc;
}

.markdown-body :deep(ol) {
  list-style: decimal;
}

.markdown-body :deep(li + li) {
  margin-top: 0.2rem;
}

.markdown-body :deep(blockquote) {
  margin: 0.5rem 0;
  border-left: 3px solid rgb(147 197 253);
  padding-left: 0.75rem;
  color: rgb(75 85 99);
}

.markdown-body :deep(pre) {
  margin: 0.5rem 0;
  overflow-x: auto;
  border-radius: 0.5rem;
  background: rgb(17 24 39);
  padding: 0.75rem;
  color: rgb(243 244 246);
}

.markdown-body :deep(code) {
  border-radius: 0.25rem;
  background: rgb(229 231 235);
  padding: 0.1rem 0.3rem;
  font-size: 0.85em;
}

.markdown-body :deep(pre code) {
  background: transparent;
  padding: 0;
  font-size: 0.85rem;
}

.markdown-body :deep(a) {
  color: rgb(37 99 235);
  text-decoration: underline;
}

.markdown-body :deep(hr) {
  margin: 0.75rem 0;
  border: 0;
  border-top: 1px solid rgb(229 231 235);
}

:global(.dark) .markdown-body :deep(blockquote) {
  color: rgb(209 213 219);
}

:global(.dark) .markdown-body :deep(code) {
  background: rgb(55 65 81);
}

:global(.dark) .markdown-body :deep(pre) {
  background: rgb(3 7 18);
}

:global(.dark) .markdown-body :deep(a) {
  color: rgb(147 197 253);
}

:global(.dark) .markdown-body :deep(hr) {
  border-top-color: rgb(55 65 81);
}
</style>
