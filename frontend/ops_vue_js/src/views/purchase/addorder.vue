<script setup>
import { reactive, ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { useValidation } from '@/composables'
import { purchaseApi } from '@/api/purchase'
import tagadder from '@/components/tagadder.vue'
import datePicker from '@/components/datePicker.vue'
import useDropzone from '@/components/useDropzone.vue'

usePageTitle('purchase.add_part')
const { t, locale } = useI18n()
const toast = useToastStore()
const { validate, errors, clearErrors } = useValidation()

const textMaxLen = 256
const photosRef = ref(null)

const currencyOptions = { 1: 'RMB', 2: 'MOA', 3: 'HKD', 4: 'USD' }

const costType = computed(() => ({
  1: t('cost_type.unit_price'),
  2: t('cost_type.freight'),
}))

const orderStatus = computed(() => ({
  1: t('order_status.pending_order'),
  2: t('order_status.order_placed'),
  3: t('order_status.in_transit'),
  4: t('order_status.completed'),
  5: t('order_status.refund_requested'),
  6: t('order_status.returning'),
  7: t('order_status.refunded'),
  8: t('order_status.lost_package'),
}))

const costEntries = reactive([])
const newCost = reactive({
  type: '1', int: 1, cost: 0, currencyType: '1',
})

const newCostTotal = computed(() =>
  parseFloat((newCost.int * newCost.cost).toFixed(2))
)

function addCostEntry() {
  if (newCost.cost <= 0) return
  costEntries.push({
    type: newCost.type,
    int: newCost.int,
    cost: newCost.cost,
    costt: newCostTotal.value,
    currencytype: newCost.currencyType,
  })
  newCost.type = '1'
  newCost.int = 1
  newCost.cost = 0
  newCost.currencyType = '1'
}

function removeCostEntry(index) {
  costEntries.splice(index, 1)
}

watch(() => newCost.cost, (val) => {
  const fixed = parseFloat(val).toFixed(2)
  if (parseFloat(fixed) !== val) newCost.cost = parseFloat(fixed)
})

const form = reactive({
  title: '',
  remark: '',
  photos: [],
  link: '',
  partname: '',
  styles: '',
  costs: [],
  tracking_number: '',
  updatetime: '',
  order_status: '1',
})

const loading = ref(false)

async function handleSubmit() {
  clearErrors()
  const err = validate('title', form.title, t('purchase_addorder.title'))
  if (!err) return

  form.photos = []
  if (photosRef.value?.has_some_files) {
    const result = photosRef.value.get_some_files()
    form.photos = result.map(f => f.name)
  }

  form.costs = costEntries.map(h => ({
    ...h,
    cost: Math.round(h.cost * 100),
    costt: Math.round(h.costt * 100),
  }))

  loading.value = true
  try {
    const { errCode } = await purchaseApi.addOrder(form)
    if (errCode === 0) {
      toast.success(t('message.save_ok'))
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // interceptor handled
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-6xl px-6 py-6">
    <div class="flex flex-col gap-6 rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card">
      <!-- Order Info -->
      <div class="border-b border-gray-200 px-6 py-4 dark:border-dk-muted">
        <h4 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('purchase_addorder.order_info') }}</h4>
      </div>
      <div class="space-y-4 px-6 py-5">
        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('purchase_addorder.title') }} <span class="text-red-500">*</span>
          </label>
          <input
            v-model="form.title"
            type="text"
            class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
            :class="errors.title ? 'border-red-500' : 'border-gray-300'"
            :placeholder="t('purchase_addorder.input_title')"
          />
          <span v-if="errors.title" class="mt-1 block text-xs text-red-500">{{ errors.title }}</span>
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('purchase_addorder.remarks') }}
            <span class="text-gray-400">{{ form.remark.length }}/{{ textMaxLen }}</span>
          </label>
          <textarea
            v-model="form.remark"
            class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
            rows="4"
            :placeholder="t('purchase_addorder.remarks_text')"
          />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('purchase_addorder.photo_remarks') }}</label>
          <useDropzone acceptFiles="image/*" uploadURL="/api/files/upload/image" maxFiles="10" ref="photosRef" />
        </div>
      </div>

      <!-- Purchase Channel -->
      <div class="border-t border-gray-200 px-6 py-4 dark:border-dk-muted">
        <h4 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('purchase_addorder.purchase_channel') }}</h4>
      </div>
      <div class="space-y-4 px-6 py-5">
        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('purchase_addorder.link') }}</label>
          <textarea
            v-model="form.link"
            class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
            rows="2"
            placeholder="url"
          />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('purchase_addorder.part_name') }}</label>
          <input
            v-model="form.partname"
            type="text"
            class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
            :placeholder="t('purchase_addorder.part_name')"
          />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('purchase_addorder.style_remarks') }}</label>
          <tagadder :placeholder="t('purchase_addorder.add_style')" v-model="form.styles" />
        </div>

        <!-- costs table -->
        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('purchase_addorder.cost') }}</label>
          <div v-if="costEntries.length" class="mb-4 overflow-x-auto">
            <table class="w-full text-left text-sm text-gray-900">
              <thead>
                <tr class="border-b border-gray-200 bg-gray-50 text-gray-500 dark:border-dk-muted dark:bg-dk-base">
                  <th class="px-3 py-2 font-medium">{{ t('purchase_addorder.type') }}</th>
                  <th class="px-3 py-2 font-medium">{{ t('purchase_addorder.quantity') }}</th>
                  <th class="px-3 py-2 font-medium">{{ t('purchase_addorder.fee') }}</th>
                  <th class="px-3 py-2 font-medium">{{ t('purchase_addorder.total_price') }}</th>
                  <th class="px-3 py-2 font-medium">{{ t('purchase_addorder.currency') }}</th>
                  <th class="px-3 py-2 font-medium">{{ t('purchase_addorder.operation') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(item, idx) in costEntries" :key="idx" class="border-b border-gray-100 dark:border-dk-muted">
                  <td class="px-3 py-2 font-medium text-gray-900 dark:text-white">{{ costType[item.type] }}</td>
                  <td class="px-3 py-2 text-gray-500">{{ item.int }}</td>
                  <td class="px-3 py-2 text-gray-500">{{ item.cost }}</td>
                  <td class="px-3 py-2 text-gray-500">{{ item.costt }}</td>
                  <td class="px-3 py-2 text-gray-500">{{ currencyOptions[item.currency_type] }}</td>
                  <td class="px-3 py-2">
                    <button class="rounded px-2 py-1 text-xs font-medium text-red-600 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/20" @click="removeCostEntry(idx)">{{ t('purchase_addorder.remove') }}</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <div class="grid grid-cols-2 gap-4 sm:grid-cols-5">
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-500">{{ t('purchase_addorder.fee_type') }}</label>
              <select v-model="newCost.type" class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white">
                <template v-for="(label, key) in costType" :key="key">
                  <option :value="key">{{ label }}</option>
                </template>
              </select>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-500">{{ t('purchase_addorder.input_quantity') }}</label>
              <input v-model.number="newCost.int" type="number" class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white" min="1" />
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-500">{{ t('purchase_addorder.input_fee') }}</label>
              <input v-model="newCost.cost" type="number" class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white" step="0.01" min="0" />
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-500">{{ t('purchase_addorder.select_currency') }}</label>
              <select v-model="newCost.currencyType" class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white">
                <template v-for="(label, key) in currencyOptions" :key="key">
                  <option :value="key">{{ label }}</option>
                </template>
              </select>
            </div>
            <div class="flex items-end">
              <button class="w-full rounded-lg border border-gray-300 bg-blue-600 px-3 py-2 text-sm font-semibold text-blue-100 transition-colors hover:bg-blue-700 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:text-white dark:bg-blue-600" @click="addCostEntry">{{ t('purchase_addorder.add') }}</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Order Status -->
      <div class="border-t border-gray-200 px-6 py-4 dark:border-dk-muted">
        <h4 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('purchase_addorder.order_status') }}</h4>
      </div>
      <div class="px-6 py-5">
        <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <div>
            <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
              {{ t('purchase_addorder.update_time') }} <span class="text-red-500">*</span>
            </label>
            <datePicker v-model="form.updatetime" />
          </div>
          <div>
            <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('purchase_addorder.tracking_number') }}</label>
            <input
              v-model="form.tracking_number"
              type="text"
              class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
              :placeholder="t('purchase_addorder.input_tracking_number')"
            />
          </div>
          <div>
            <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('purchase_addorder.order_status') }}</label>
            <select v-model="form.order_status" class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white">
              <template v-for="(label, key) in orderStatus" :key="key">
                <option :value="key">{{ label }}</option>
              </template>
            </select>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="flex justify-end border-t border-gray-200 px-6 py-4 dark:border-dk-muted">
        <button
          class="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-blue-700 focus:ring-2 focus:ring-blue-500/20 focus:outline-none disabled:active:scale-100"
          :disabled="loading"
          @click="handleSubmit"
        >
          <svg v-if="loading" class="h-4 w-4 animate-spin text-white" viewBox="0 0 24 24" fill="none">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
          {{ t('purchase_addorder.submit') }}
        </button>
      </div>
    </div>
  </div>
</template>
