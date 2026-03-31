<script setup>
import { ref, onMounted, watch } from 'vue'
import flatpickr from 'flatpickr'
import 'flatpickr/dist/flatpickr.css'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()

const props = defineProps({
  modelValue: { type: String, default: '' },
  placeholder: { type: String, default: '' },
})

const emit = defineEmits(['update:modelValue'])

const inputEl = ref(null)
let picker = null

onMounted(() => {
  picker = flatpickr(inputEl.value, {
    dateFormat: 'Y-m-d',
    defaultDate: props.modelValue || null,
    allowInput: true,
    disableMobile: true,
    parseDate(datestr) { return new Date(datestr) },
    onChange(selectedDates, dateStr) { emit('update:modelValue', dateStr) },
  })
})

watch(() => props.modelValue, (val) => {
  if (picker && val !== picker.input.value) { picker.setDate(val, false) }
})
</script>

<template>
  <div class="relative">
    <span class="pointer-events-none absolute inset-y-0 left-3 flex items-center text-gray-400 dark:text-dk-subtle">
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M4 7a2 2 0 0 1 2 -2h12a2 2 0 0 1 2 2v12a2 2 0 0 1 -2 2h-12a2 2 0 0 1 -2 -2v-12z"/><path d="M16 3v4"/><path d="M8 3v4"/><path d="M4 11l16 0"/><path d="M11 15h1"/><path d="M12 15v3"/>
      </svg>
    </span>
    <input ref="inputEl" type="text" :placeholder="placeholder || t('message.select_date')" class="w-full rounded-lg border border-gray-300 bg-white py-2 pr-3 pl-9 text-sm focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 focus:outline-none dark:border-dk-muted dark:bg-dk-card dark:text-dk-text" :value="modelValue" readonly />
  </div>
</template>
