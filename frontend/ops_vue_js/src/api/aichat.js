import { api } from './index'
import { useUserStore } from '@/stores/user'

export async function fetchOpenAIProfiles() {
  return api.get('/aichat/openai')
}

export async function fetchAIChatAdminConfig() {
  return api.post('/aichat/admin/config', {})
}

export async function updateAIChatAdminConfig(config) {
  return api.post('/aichat/admin/config/update', config)
}

export async function refreshAIChatAdminConfig() {
  return api.post('/aichat/admin/refresh', {})
}

function parseSSEBlock(block) {
  const lines = block.split('\n')
  const dataLines = []

  for (const line of lines) {
    if (line.startsWith('data:')) {
      dataLines.push(line.slice(5).trimStart())
    }
  }

  if (dataLines.length === 0) return null
  return dataLines.join('\n')
}

export async function streamChat(messages, options = {}, handlers = {}) {
  if (typeof options.onDelta === 'function' || typeof options.onError === 'function') {
    handlers = options
    options = {}
  }

  const userStore = useUserStore()
  const response = await fetch('/api/aichat/chat', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      userCookieValue: userStore.cookieValue || '',
      data: { messages, openaiName: options.openaiName || '' },
    }),
  })

  if (!response.ok) {
    handlers.onError?.(`HTTP ${response.status}`)
    return
  }
  if (!response.body) {
    handlers.onError?.('ReadableStream is not supported')
    return
  }

  const reader = response.body.getReader()
  const decoder = new TextDecoder()
  let buffer = ''

  while (true) {
    const { done, value } = await reader.read()
    if (done) break

    buffer += decoder.decode(value, { stream: true })
    const blocks = buffer.split('\n\n')
    buffer = blocks.pop() || ''

    for (const block of blocks) {
      const payload = parseSSEBlock(block)
      if (!payload) continue
      if (payload === '[DONE]') {
        handlers.onDone?.()
        return
      }

      try {
        const frame = JSON.parse(payload)
        switch (frame.type) {
          case 'delta':
            handlers.onDelta?.(frame.text || '')
            break
          case 'trace':
            handlers.onTrace?.(frame)
            break
          case 'stats':
            handlers.onStats?.(frame.stats || null)
            break
          case 'error':
            handlers.onError?.(frame.error || frame.message || 'AI request failed')
            break
          case '[DONE]':
            handlers.onDone?.()
            return
          default:
            handlers.onFrame?.(frame)
            break
        }
      } catch (error) {
        handlers.onError?.(error.message)
      }
    }
  }

  handlers.onDone?.()
}
