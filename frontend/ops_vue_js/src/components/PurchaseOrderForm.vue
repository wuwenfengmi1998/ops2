<script setup>
/**
 * PurchaseOrderForm —— 采购订单表单公共组件
 *
 * 供 addorder.vue 和 editorder.vue 共用，包含：
 * - 标题、备注、链接、款式标签
 * - 费用明细（添加/删除）
 * - 图片上传（useDropzone）
 *
 * Props:
 *   modelValue  {Object}   表单数据（v-model）
 *   initialCosts {Array}   回填的费用列表（编辑时传入，单位：分）
 *   initialPhotos {Array}  回填的图片列表 [{ id, hash, name, ... }]
 *
 * Emits:
 *   update:modelValue
 */
import { reactive, ref, computed, watch, onMounted } from "vue";
import { useI18n } from "vue-i18n";
import tagadder from "@/components/tagadder.vue";
import useDropzone from "@/components/useDropzone.vue";

const props = defineProps({
  modelValue: {
    type: Object,
    required: true,
  },
  /** 回填图片列表 [{ Sha256, Name, ... }] */
  initialPhotos: {
    type: Array,
    default: () => [],
  },
});

const emit = defineEmits(["update:modelValue"]);

const { t } = useI18n();

// ==================== 常量 ====================
const textMaxLen = 256;
const currencyOptions = { 1: "CNY", 2: "MOP", 3: "HKD", 4: "USD" };

const costType = computed(() => ({
  1: t("cost_type.unit_price"),
  2: t("cost_type.freight"),
}));

// ==================== 费用明细 ====================
/** 展示用（元为单位） */
const costEntries = reactive([]);

const newCost = reactive({
  type: 1,
  int: 1,
  cost: 0,
  currencyType: 1,
});

const costError = ref(false);

const newCostTotal = computed(() =>
  parseFloat((newCost.int * newCost.cost).toFixed(2)),
);

function addCostEntry() {
  if (newCost.cost <= 0) {
    costError.value = true;
    return;
  }
  costError.value = false;
  costEntries.push({
    type: newCost.type,
    int: newCost.int,
    cost: newCost.cost,
    costt: newCostTotal.value,
    currencytype: newCost.currencyType,
  });
  newCost.type = 1;
  newCost.int = 1;
  newCost.cost = 0;
  newCost.currencyType = 1;
  syncCosts();
}

function removeCostEntry(index) {
  costEntries.splice(index, 1);
  syncCosts();
}

/** 将当前 costEntries（元）转换为分并同步到父组件 */
function syncCosts() {
  // 直接更新父组件的 form._costs，跳过 emit 链路
  props.modelValue._costs = costEntries.map((h) => ({
    ...h,
    cost: Math.round(h.cost * 100),
    costt: Math.round(h.costt * 100),
  }));
}

watch(
  () => newCost.cost,
  (val) => {
    const fixed = parseFloat(val).toFixed(2);
    if (parseFloat(fixed) !== val) newCost.cost = parseFloat(fixed);
    if (val > 0) costError.value = false;
  },
);

// 回填费用（父组件填充 form._costs 后直接消费）
watch(
  () => props.modelValue._costs,
  (list) => {
    if (!list || list.length === 0) return;
    // 先清空默认空白行，再填入 API 返回的数据
    costEntries.splice(0, costEntries.length);
    list.forEach((c) => {
      costEntries.push({ ...c });
    });
  },
);

// ==================== 图片上传 ====================
const photosRef = ref(null);

/**
 * 供父组件调用，获取当前上传图片的哈希列表
 */
function getPhotoHashes() {
  return photosRef.value?.return_files().map((f) => f.hash) ?? [];
}

defineExpose({ getPhotoHashes, costEntries });

// ==================== 表单字段双向绑定 ====================
function update(field, value) {
  emit("update:modelValue", { ...props.modelValue, [field]: value });
}
</script>

<template>
  <div class="space-y-4 px-6 py-5">
    <!-- 标题（必填） -->
    <div>
      <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
        {{ t("purchase_addorder.part_name") }}
        <span class="text-red-500">*</span>
      </label>
      <input
        :value="modelValue.title"
        @input="update('title', $event.target.value)"
        type="text"
        maxlength="50"
        class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
        :placeholder="t('purchase_addorder.part_name')"
      />
    </div>

    <!-- 备注 -->
    <div>
      <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
        {{ t("purchase_addorder.remarks") }}
        <span class="text-gray-400">{{ modelValue.remark.length }}/{{ textMaxLen }}</span>
      </label>
      <textarea
        :value="modelValue.remark"
        @input="update('remark', $event.target.value)"
        class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
        rows="4"
        :placeholder="t('purchase_addorder.remarks_text')"
      />
    </div>

    <!-- 采购链接 -->
    <div>
      <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
        {{ t("purchase_addorder.link") }}
      </label>
      <textarea
        :value="modelValue.link"
        @input="update('link', $event.target.value)"
        class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
        rows="2"
        placeholder="url"
      />
    </div>

    <!-- 款式标签 -->
    <div>
      <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
        {{ t("purchase_addorder.style_remarks") }}
      </label>
      <tagadder
        :placeholder="t('purchase_addorder.add_style')"
        :modelValue="modelValue.styles"
        @update:modelValue="update('styles', $event)"
      />
    </div>

    <!-- 费用明细 -->
    <div>
      <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
        {{ t("purchase_addorder.cost") }}
      </label>

      <!-- 已有费用列表 -->
      <div v-if="costEntries.length" class="mb-4 overflow-x-auto">
        <table class="w-full text-left text-sm text-gray-900">
          <thead>
            <tr class="border-b border-gray-200 bg-gray-50 text-gray-500 dark:border-dk-muted dark:bg-dk-base">
              <th class="px-3 py-2 font-medium">{{ t("purchase_addorder.type") }}</th>
              <th class="px-3 py-2 font-medium">{{ t("purchase_addorder.quantity") }}</th>
              <th class="px-3 py-2 font-medium">{{ t("purchase_addorder.fee") }}</th>
              <th class="px-3 py-2 font-medium">{{ t("purchase_addorder.total_price") }}</th>
              <th class="px-3 py-2 font-medium">{{ t("purchase_addorder.currency") }}</th>
              <th class="px-3 py-2 font-medium">{{ t("purchase_addorder.operation") }}</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="(item, idx) in costEntries"
              :key="idx"
              class="border-b border-gray-100 dark:border-dk-muted"
            >
              <td class="px-3 py-2 font-medium text-gray-900 dark:text-white">{{ costType[item.type] }}</td>
              <td class="px-3 py-2 text-gray-500">{{ item.int }}</td>
              <td class="px-3 py-2 text-gray-500">{{ item.cost }}</td>
              <td class="px-3 py-2 text-gray-500">{{ item.costt }}</td>
              <td class="px-3 py-2 text-gray-500">{{ currencyOptions[item.currencytype] }}</td>
              <td class="px-3 py-2">
                <button
                  class="rounded px-2 py-1 text-xs font-medium text-red-600 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/20"
                  @click="removeCostEntry(idx)"
                >
                  {{ t("purchase_addorder.remove") }}
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 添加费用表单 -->
      <div class="grid grid-cols-2 gap-4 sm:grid-cols-5">
        <div>
          <label class="mb-1 block text-xs font-medium text-gray-500">{{ t("purchase_addorder.fee_type") }}</label>
          <select
            v-model="newCost.type"
            class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
          >
            <template v-for="(label, key) in costType" :key="key">
              <option :value="key">{{ label }}</option>
            </template>
          </select>
        </div>
        <div>
          <label class="mb-1 block text-xs font-medium text-gray-500">{{ t("purchase_addorder.input_quantity") }}</label>
          <input
            v-model.number="newCost.int"
            type="number"
            class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
            min="1"
          />
        </div>
        <div>
          <label class="mb-1 block text-xs font-medium text-gray-500">{{ t("purchase_addorder.input_fee") }}</label>
          <input
            v-model="newCost.cost"
            type="number"
            class="w-full rounded-lg border bg-white px-3 py-2 text-sm dark:bg-dk-base dark:text-white"
            :class="costError ? 'border-red-500' : 'border-gray-300 dark:border-dk-muted'"
            step="0.01"
            min="0"
          />
        </div>
        <div>
          <label class="mb-1 block text-xs font-medium text-gray-500">{{ t("purchase_addorder.select_currency") }}</label>
          <select
            v-model="newCost.currencyType"
            class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
          >
            <template v-for="(label, key) in currencyOptions" :key="key">
              <option :value="key">{{ label }}</option>
            </template>
          </select>
        </div>
        <div class="flex items-end">
          <button
            class="w-full rounded-lg border border-gray-300 bg-blue-600 px-3 py-2 text-sm font-semibold text-blue-100 transition-colors hover:bg-blue-700 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-blue-600 dark:text-white"
            @click="addCostEntry"
          >
            {{ t("purchase_addorder.add") }}
          </button>
        </div>
      </div>
    </div>

    <!-- 图片上传 -->
    <div>
      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
        {{ t("purchase_addorder.photo_remarks") }}
      </label>
      <useDropzone
        acceptFiles="image/*"
        uploadURL="/api/files/upload/image"
        :maxFiles="10"
        :initialFiles="initialPhotos"
        ref="photosRef"
      />
    </div>
  </div>
</template>
