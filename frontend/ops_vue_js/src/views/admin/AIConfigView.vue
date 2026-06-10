<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { IconDeviceFloppy, IconPlus, IconRefresh, IconTrash } from '@tabler/icons-vue'
import { fetchAIChatAdminConfig, refreshAIChatAdminConfig, updateAIChatAdminConfig } from '@/api/aichat'
import { usePageTitle } from '@/composables/usePageTitle'
import { useToastStore } from '@/stores/toast'

const { t } = useI18n()
const toast = useToastStore()

usePageTitle('aiconfig.title')

const loading = ref(false)
const saving = ref(false)
const refreshing = ref(false)

const form = reactive({
  enabled: false,
  openai: [],
  toolRouter: {
    enabled: true,
    openaiName: '',
    timeout: 30,
    maxTokens: 512,
    tools: [],
  },
})

const profileNames = computed(() => form.openai.map((profile) => profile.name.trim()).filter(Boolean))
const activeProfileName = computed(() => form.openai.find((profile) => profile.active)?.name || form.openai[0]?.name || '')

onMounted(loadConfig)

function normalizeProfile(profile = {}) {
  return {
    name: profile.name || '',
    active: !!profile.active,
    apiKey: '',
    apiKeySet: !!profile.apiKeySet,
    baseUrl: profile.baseUrl || '',
    model: profile.model || '',
    timeout: Number(profile.timeout || 120),
    maxTokens: Number(profile.maxTokens || 4096),
    systemPrompt: profile.systemPrompt || '',
  }
}

function normalizeTool(tool = {}) {
  return {
    name: tool.name || '',
    enabled: tool.enabled !== false,
    description: tool.description || '',
  }
}

function applyConfig(data) {
  form.enabled = !!data?.enabled
  form.openai = Array.isArray(data?.openai) ? data.openai.map(normalizeProfile) : []
  if (form.openai.length === 0) {
    addProfile()
  }

  const router = data?.toolRouter || {}
  form.toolRouter = {
    enabled: router.enabled !== false,
    openaiName: router.openaiName || '',
    timeout: Number(router.timeout || 30),
    maxTokens: Number(router.maxTokens || 512),
    tools: Array.isArray(router.tools) ? router.tools.map(normalizeTool) : [],
  }
  if (form.toolRouter.tools.length === 0) {
    addTool()
  }
}

async function loadConfig() {
  loading.value = true
  try {
    const res = await fetchAIChatAdminConfig()
    if (res.errCode === 0) {
      applyConfig(res.data || {})
    } else {
      toast.error(t('aiconfig.load_failed'))
    }
  } catch (error) {
    toast.error(error instanceof Error ? error.message : String(error))
  } finally {
    loading.value = false
  }
}

async function refreshCache() {
  refreshing.value = true
  try {
    const res = await refreshAIChatAdminConfig()
    if (res.errCode === 0) {
      toast.success(t('aiconfig.refresh_success'))
      await loadConfig()
    } else {
      toast.error(t('aiconfig.refresh_failed'))
    }
  } catch (error) {
    toast.error(error instanceof Error ? error.message : String(error))
  } finally {
    refreshing.value = false
  }
}

function addProfile() {
  form.openai.push({
    name: '',
    active: form.openai.length === 0,
    apiKey: '',
    apiKeySet: false,
    baseUrl: 'https://ark.cn-beijing.volces.com/api/v3',
    model: '',
    timeout: 120,
    maxTokens: 4096,
    systemPrompt: '',
  })
}

function removeProfile(index) {
  const removed = form.openai[index]
  form.openai.splice(index, 1)
  if (form.openai.length === 0) {
    addProfile()
  }
  if (!form.openai.some((profile) => profile.active)) {
    form.openai[0].active = true
  }
  if (removed?.name && form.toolRouter.openaiName === removed.name) {
    form.toolRouter.openaiName = activeProfileName.value
  }
}

function setActiveProfile(index) {
  form.openai.forEach((profile, i) => {
    profile.active = i === index
  })
}

function addTool() {
  form.toolRouter.tools.push({ name: '', enabled: true, description: '' })
}

function removeTool(index) {
  form.toolRouter.tools.splice(index, 1)
}

function validate() {
  if (form.openai.length === 0) {
    toast.error(t('aiconfig.error_profile_required'))
    return false
  }

  const names = new Set()
  let hasActive = false
  for (const profile of form.openai) {
    profile.name = profile.name.trim()
    profile.baseUrl = profile.baseUrl.trim()
    profile.model = profile.model.trim()
    profile.apiKey = profile.apiKey.trim()

    if (!profile.name) {
      toast.error(t('aiconfig.error_profile_name_required'))
      return false
    }
    if (names.has(profile.name)) {
      toast.error(t('aiconfig.error_profile_name_duplicate'))
      return false
    }
    names.add(profile.name)
    if (profile.active) {
      hasActive = true
    }

    profile.timeout = Number(profile.timeout)
    profile.maxTokens = Number(profile.maxTokens)
    if (!Number.isInteger(profile.timeout) || profile.timeout <= 0) {
      toast.error(t('aiconfig.error_timeout'))
      return false
    }
    if (!Number.isInteger(profile.maxTokens) || profile.maxTokens <= 0) {
      toast.error(t('aiconfig.error_max_tokens'))
      return false
    }
  }

  if (!hasActive) {
    form.openai[0].active = true
  }

  const active = form.openai.find((profile) => profile.active)
  if (form.enabled && active) {
    if (!active.baseUrl || !active.model) {
      toast.error(t('aiconfig.error_active_profile_required'))
      return false
    }
    if (!active.apiKeySet && !active.apiKey) {
      toast.error(t('aiconfig.error_api_key_required'))
      return false
    }
  }

  if (form.toolRouter.openaiName && !names.has(form.toolRouter.openaiName)) {
    toast.error(t('aiconfig.error_router_profile'))
    return false
  }

  form.toolRouter.timeout = Number(form.toolRouter.timeout)
  form.toolRouter.maxTokens = Number(form.toolRouter.maxTokens)
  if (!Number.isInteger(form.toolRouter.timeout) || form.toolRouter.timeout <= 0) {
    toast.error(t('aiconfig.error_timeout'))
    return false
  }
  if (!Number.isInteger(form.toolRouter.maxTokens) || form.toolRouter.maxTokens <= 0) {
    toast.error(t('aiconfig.error_max_tokens'))
    return false
  }

  const toolNames = new Set()
  for (const tool of form.toolRouter.tools) {
    tool.name = tool.name.trim()
    if (!tool.name) {
      toast.error(t('aiconfig.error_tool_name_required'))
      return false
    }
    if (toolNames.has(tool.name)) {
      toast.error(t('aiconfig.error_tool_name_duplicate'))
      return false
    }
    toolNames.add(tool.name)
  }

  return true
}

function buildPayload() {
  return {
    enabled: form.enabled,
    openai: form.openai.map((profile) => ({
      name: profile.name.trim(),
      active: !!profile.active,
      apiKey: profile.apiKey.trim(),
      baseUrl: profile.baseUrl.trim(),
      model: profile.model.trim(),
      timeout: Number(profile.timeout),
      maxTokens: Number(profile.maxTokens),
      systemPrompt: profile.systemPrompt || '',
    })),
    toolRouter: {
      enabled: !!form.toolRouter.enabled,
      openaiName: form.toolRouter.openaiName || '',
      timeout: Number(form.toolRouter.timeout),
      maxTokens: Number(form.toolRouter.maxTokens),
      tools: form.toolRouter.tools.map((tool) => ({
        name: tool.name.trim(),
        enabled: !!tool.enabled,
        description: tool.description || '',
      })),
    },
  }
}

async function saveConfig() {
  if (!validate()) return
  saving.value = true
  try {
    const res = await updateAIChatAdminConfig(buildPayload())
    if (res.errCode === 0) {
      toast.success(t('aiconfig.save_success'))
      await loadConfig()
    } else {
      toast.error(t('aiconfig.save_failed'))
    }
  } catch (error) {
    toast.error(error instanceof Error ? error.message : String(error))
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 p-6 dark:bg-dk-base">
    <div class="mx-auto max-w-6xl space-y-6">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-dk-text">{{ t('aiconfig.title') }}</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-dk-subtle">{{ t('aiconfig.subtitle') }}</p>
        </div>
        <div class="flex flex-wrap gap-2">
          <button class="btn-secondary" :disabled="loading" @click="loadConfig">
            <IconRefresh :size="16" />
            {{ t('aiconfig.reload') }}
          </button>
          <button class="btn-secondary" :disabled="refreshing" @click="refreshCache">
            <IconRefresh :size="16" />
            {{ t('aiconfig.refresh_cache') }}
          </button>
          <button class="btn-primary" :disabled="saving" @click="saveConfig">
            <IconDeviceFloppy :size="16" />
            {{ t('aiconfig.save') }}
          </button>
        </div>
      </div>

      <section class="card">
        <label class="flex items-center gap-3 text-sm font-medium text-gray-700 dark:text-dk-text">
          <input v-model="form.enabled" type="checkbox" class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500" />
          {{ t('aiconfig.enabled') }}
        </label>
      </section>

      <section class="card space-y-4">
        <div class="flex items-center justify-between">
          <h2 class="section-title">{{ t('aiconfig.profiles') }}</h2>
          <button class="btn-secondary" @click="addProfile">
            <IconPlus :size="16" />
            {{ t('aiconfig.add_profile') }}
          </button>
        </div>

        <div v-for="(profile, index) in form.openai" :key="index" class="rounded-lg border border-gray-200 p-4 dark:border-dk-muted">
          <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
            <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-dk-text">
              <input :checked="profile.active" type="radio" name="active-profile" class="text-blue-600" @change="setActiveProfile(index)" />
              {{ t('aiconfig.active') }}
            </label>
            <button class="text-sm text-red-600 hover:text-red-700" @click="removeProfile(index)">
              <IconTrash :size="16" class="inline" />
              {{ t('aiconfig.remove') }}
            </button>
          </div>

          <div class="grid gap-4 md:grid-cols-2">
            <label class="field">
              <span>{{ t('aiconfig.name') }}</span>
              <input v-model="profile.name" class="input" />
            </label>
            <label class="field">
              <span>{{ t('aiconfig.api_key') }}</span>
              <input v-model="profile.apiKey" class="input" type="password" :placeholder="profile.apiKeySet ? t('aiconfig.api_key_keep') : ''" />
            </label>
            <label class="field md:col-span-2">
              <span>{{ t('aiconfig.base_url') }}</span>
              <input v-model="profile.baseUrl" class="input" />
            </label>
            <label class="field">
              <span>{{ t('aiconfig.model') }}</span>
              <input v-model="profile.model" class="input" />
            </label>
            <label class="field">
              <span>{{ t('aiconfig.timeout') }}</span>
              <input v-model.number="profile.timeout" class="input" type="number" min="1" />
            </label>
            <label class="field">
              <span>{{ t('aiconfig.max_tokens') }}</span>
              <input v-model.number="profile.maxTokens" class="input" type="number" min="1" />
            </label>
            <label class="field md:col-span-2">
              <span>{{ t('aiconfig.system_prompt') }}</span>
              <textarea v-model="profile.systemPrompt" class="input min-h-24 resize-y" />
            </label>
          </div>
        </div>
      </section>

      <section class="card space-y-4">
        <h2 class="section-title">{{ t('aiconfig.tool_router') }}</h2>
        <div class="grid gap-4 md:grid-cols-4">
          <label class="flex items-center gap-3 text-sm font-medium text-gray-700 dark:text-dk-text">
            <input v-model="form.toolRouter.enabled" type="checkbox" class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500" />
            {{ t('aiconfig.enabled') }}
          </label>
          <label class="field md:col-span-1">
            <span>{{ t('aiconfig.router_profile') }}</span>
            <select v-model="form.toolRouter.openaiName" class="input">
              <option value="">{{ t('aiconfig.none') }}</option>
              <option v-for="name in profileNames" :key="name" :value="name">{{ name }}</option>
            </select>
          </label>
          <label class="field">
            <span>{{ t('aiconfig.timeout') }}</span>
            <input v-model.number="form.toolRouter.timeout" class="input" type="number" min="1" />
          </label>
          <label class="field">
            <span>{{ t('aiconfig.max_tokens') }}</span>
            <input v-model.number="form.toolRouter.maxTokens" class="input" type="number" min="1" />
          </label>
        </div>
      </section>

      <section class="card space-y-4">
        <div class="flex items-center justify-between">
          <h2 class="section-title">{{ t('aiconfig.tools') }}</h2>
          <button class="btn-secondary" @click="addTool">
            <IconPlus :size="16" />
            {{ t('aiconfig.add_tool') }}
          </button>
        </div>
        <div v-for="(tool, index) in form.toolRouter.tools" :key="index" class="grid gap-4 rounded-lg border border-gray-200 p-4 dark:border-dk-muted md:grid-cols-12">
          <label class="field md:col-span-3">
            <span>{{ t('aiconfig.name') }}</span>
            <input v-model="tool.name" class="input" />
          </label>
          <label class="flex items-center gap-3 text-sm font-medium text-gray-700 dark:text-dk-text md:col-span-2 md:pt-6">
            <input v-model="tool.enabled" type="checkbox" class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500" />
            {{ t('aiconfig.enabled') }}
          </label>
          <label class="field md:col-span-6">
            <span>{{ t('aiconfig.description') }}</span>
            <input v-model="tool.description" class="input" />
          </label>
          <div class="flex items-end md:col-span-1">
            <button class="text-sm text-red-600 hover:text-red-700" @click="removeTool(index)">
              <IconTrash :size="18" />
            </button>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.card {
  border-radius: 0.75rem;
  border: 1px solid rgb(229 231 235);
  background: white;
  padding: 1.25rem;
  box-shadow: 0 1px 2px rgb(0 0 0 / 0.05);
}
.section-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: rgb(17 24 39);
}
.field {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: rgb(55 65 81);
}
.input {
  border-radius: 0.375rem;
  border: 1px solid rgb(209 213 219);
  background: white;
  padding: 0.5rem 0.75rem;
  font-size: 0.875rem;
  color: rgb(17 24 39);
  outline: none;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}
.input:focus {
  border-color: rgb(59 130 246);
  box-shadow: 0 0 0 2px rgb(59 130 246 / 0.2);
}
.btn-primary,
.btn-secondary {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  border-radius: 0.375rem;
  padding: 0.5rem 0.75rem;
  font-size: 0.875rem;
  font-weight: 500;
  transition: background-color 0.15s ease, color 0.15s ease;
}
.btn-primary {
  background: rgb(37 99 235);
  color: white;
}
.btn-primary:hover {
  background: rgb(29 78 216);
}
.btn-secondary {
  border: 1px solid rgb(229 231 235);
  color: rgb(75 85 99);
}
.btn-secondary:hover {
  background: rgb(243 244 246);
}
.btn-primary:disabled,
.btn-secondary:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}
:global(.dark) .card {
  border-color: var(--color-dk-muted);
  background: var(--color-dk-card);
}
:global(.dark) .section-title,
:global(.dark) .field {
  color: var(--color-dk-text);
}
:global(.dark) .input {
  border-color: var(--color-dk-muted);
  background: var(--color-dk-base);
  color: var(--color-dk-text);
}
:global(.dark) .btn-secondary {
  border-color: var(--color-dk-muted);
  color: var(--color-dk-subtle);
}
:global(.dark) .btn-secondary:hover {
  background: var(--color-dk-card);
}
</style>
