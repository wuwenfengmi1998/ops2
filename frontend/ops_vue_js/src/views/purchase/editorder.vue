<script setup>
/**
 * 采购订单编辑页面
 *
 * 功能概述：
 * - 通过路由参数 :id 加载已有订单数据
 * - 费用明细直接在本页管理（与 addorder.vue 相同模式）
 * - 提交时调用 /purchase/updateorder 保存修改
 */

import { reactive, ref, computed, watch, onMounted, nextTick } from "vue";
import { useI18n } from "vue-i18n";
import { useRoute, useRouter } from "vue-router";
import { useToastStore } from "@/stores/toast";
import { usePageTitle } from "@/composables/usePageTitle";
import { useValidation } from "@/composables";
import { purchaseApi } from "@/api/purchase";
import tagadder from "@/components/tagadder.vue";
import useDropzone from "@/components/useDropzone.vue";
import ConfirmDialog from "@/components/ConfirmDialog.vue";

usePageTitle("purchase_addorder.edit_order");

const route = useRoute();
const router = useRouter();
const { t } = useI18n();
const toast = useToastStore();
const { validate, errors, clearErrors } = useValidation();

const orderId = Number(route.params.id);

// ==================== 状态 ====================
const loading = ref(false);
const pageLoading = ref(true);
const pageError = ref("");

// ==================== 表单数据 ====================
const form = reactive({
  title: "",
  remark: "",
  link: "",
  styles: "",
  photos: [],
});

// ==================== 费用明细（与 addorder.vue 完全一致） ====================
const textMaxLen = 256;
const currencyOptions = { 1: "CNY", 2: "MOP", 3: "HKD", 4: "USD" };
const costType = computed(() => ({
  1: t("cost_type.unit_price"),
  2: t("cost_type.freight"),
}));

/** 已添加的费用列表 */
const costEntries = reactive([]);
const costError = ref(false);

const newCost = reactive({
  type: 1,
  int: 1,
  cost: 0,
  currencytype: 1,
});

function addCostEntry() {
  if (!newCost.cost || parseFloat(newCost.cost) <= 0) {
    costError.value = true;
    return;
  }
  const cost = parseFloat(newCost.cost);
  costEntries.push({
    type: newCost.type,
    int: newCost.int,
    cost,
    costt: parseFloat((cost * newCost.int).toFixed(2)),
    currencytype: newCost.currencytype,
  });
  newCost.cost = 0;
  newCost.type = 1;
  newCost.int = 1;
  newCost.currencytype = 1;
  costError.value = false;
}

function removeCostEntry(idx) {
  costEntries.splice(idx, 1);
}

watch(
  () => newCost.cost,
  (val) => {
    const fixed = parseFloat(val).toFixed(2);
    if (parseFloat(fixed) !== val) newCost.cost = parseFloat(fixed);
    if (val > 0) costError.value = false;
  },
);

// ==================== 图片上传 ====================
const dropzoneRef = ref(null);
const showDeleteConfirm = ref(false);

function getPhotoHashes() {
  return dropzoneRef.value?.return_files().map((f) => f.hash) ?? [];
}

// ==================== 加载订单数据 ====================
onMounted(async () => {
  if (!orderId) {
    pageError.value = t("purchase.order_not_found");
    pageLoading.value = false;
    return;
  }

  try {
    const res = await purchaseApi.getOrder(orderId);
    if (res.errCode !== 0 || !res.data) {
      pageError.value = t("purchase.order_not_found");
      pageLoading.value = false;
      return;
    }

    const { order, costs, photos } = res.data;

    // 回填基本信息
    form.title = order.Title ?? "";
    form.remark = order.Remark ?? "";
    form.link = order.Link ?? "";
    form.styles = order.Styles ?? "";

    // 回填费用（分→元，直接写 costEntries）
    if (costs && costs.length > 0) {
      costs.forEach((c) => {
        costEntries.push({
          type: c.CostType,
          int: c.Quantity,
          cost: parseFloat((c.Price / 100).toFixed(2)),
          costt: parseFloat(((c.Price * c.Quantity) / 100).toFixed(2)),
          currencytype: c.CurrencyType,
        });
      });
    }

    // 回填图片
    await nextTick();
    if (photos && photos.length > 0) {
      //dropzoneRef.value?.loadInitialFiles(photos);
      //console.log(photos)
      form.photos=photos
    }
  } catch {
    pageError.value = t("purchase.order_not_found");
  } finally {
    pageLoading.value = false;
  }
});

// ==================== 提交 ====================
// ==================== 删除订单 ====================
async function handleDelete() {
  showDeleteConfirm.value = true;
}

async function doDelete() {

  loading.value = true;
  try {
    const res = await purchaseApi.deleteOrder(orderId);
    if (res.errCode === 0) {
      toast.success(t("message.delete_ok"));
      router.replace("/purchase");
    } else {
      toast.error(t("message.server_error"));
    }
  } catch {
    toast.error(t("message.server_error"));
  } finally {
    loading.value = false;
  }
}

// ==================== 提交 ====================
async function handleSubmit() {
  clearErrors();
  const ok = validate("title", form.title, t("purchase_addorder.title"));
  if (!ok) return;

  form.photos = getPhotoHashes();
  // 费用（转为分）
  const rawCosts = costEntries.map((h) => ({
    type: h.type,
    int: h.int,
    cost: Math.round(h.cost * 100),
    costt: Math.round(h.costt * 100),
    currencytype: h.currencytype,
  }));

  loading.value = true;
  console.log(form.photos)
  try {
    const res = await purchaseApi.updateOrder(orderId, {
      title: form.title,
      remark: form.remark,
      link: form.link,
      styles: form.styles,
      photos: form.photos,
      costs: rawCosts,
    });

    if (res.errCode === 0) {
      toast.success(t("message.save_ok"));
      setTimeout(() => {
        router.replace(`/purchase/showorder/${orderId}`);
      }, 800);
    } else {
      toast.error(t("message.server_error"));
    }
  } catch {
    toast.error(t("message.server_error"));
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="mx-auto max-w-6xl px-6 py-6">
    <!-- 加载中 -->
    <div
      v-if="pageLoading"
      class="flex items-center justify-center py-20 text-gray-400"
    >
      <svg class="mr-2 h-5 w-5 animate-spin" viewBox="0 0 24 24" fill="none">
        <circle
          class="opacity-25"
          cx="12"
          cy="12"
          r="10"
          stroke="currentColor"
          stroke-width="4"
        />
        <path
          class="opacity-75"
          fill="currentColor"
          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
        />
      </svg>
      {{ t("message.loading") }}
    </div>

    <!-- 订单不存在 -->
    <div
      v-else-if="pageError"
      class="rounded-xl border border-red-200 bg-red-50 px-6 py-10 text-center text-red-500 dark:border-red-800 dark:bg-red-900/20"
    >
      {{ pageError }}
    </div>

    <!-- 主卡片 -->
    <div
      v-else
      class="flex flex-col gap-0 rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card"
    >
      <!-- 顶部标题栏 -->
      <div
        class="flex items-center justify-between border-b border-gray-200 px-6 py-4 dark:border-dk-muted"
      >
        <h4 class="text-sm font-semibold text-gray-900 dark:text-white">
          {{ t("purchase_addorder.edit_order") }}
        </h4>
        <!-- 操作按钮组 -->
        <div class="flex items-center gap-2">
          <!-- 删除按钮 -->
          <button
            class="flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20"
            :disabled="loading"
            @click="handleDelete"
          >
            <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            {{ t("purchase.delete_order") }}
          </button>
          <!-- 返回按钮 -->
          <button
            class="flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-dk-base"
            @click="router.back()"
          >
          <svg
            class="h-4 w-4"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M15 19l-7-7 7-7"
            />
          </svg>
          {{ t("purchase.back") }}
        </button>
        </div>
      </div>

      <!-- 错误提示（字段验证） -->
      <div
        v-if="errors.title"
        class="mx-6 mt-4 rounded-lg bg-red-50 px-4 py-2 text-sm text-red-600 dark:bg-red-900/20 dark:text-red-400"
      >
        {{ errors.title }}
      </div>

      <!-- ==================== 订单信息区块 ==================== -->
      <div class="border-b border-gray-200 px-6 py-4 dark:border-dk-muted">
        <h4 class="text-sm font-semibold text-gray-900 dark:text-white">
          {{ t("purchase_addorder.order_info") }}
        </h4>
      </div>

      <div class="space-y-4 px-6 py-5">
        <!-- 标题字段（必填） -->
        <div>
          <label
            class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300"
          >
            {{ t("purchase_addorder.part_name") }}
            <span class="text-red-500">*</span>
          </label>
          <input
            v-model="form.title"
            type="text"
            maxlength="50"
            class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
            :class="errors.title ? 'border-red-500' : 'border-gray-300'"
            :placeholder="t('purchase_addorder.part_name')"
          />
          <span v-if="errors.title" class="mt-1 block text-xs text-red-500">{{
            errors.title
          }}</span>
        </div>

        <!-- 备注字段 -->
        <div>
          <label
            class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300"
          >
            {{ t("purchase_addorder.remarks") }}
            <span class="text-gray-400"
              >{{ form.remark.length }}/{{ textMaxLen }}</span
            >
          </label>
          <textarea
            v-model="form.remark"
            class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
            rows="4"
            :placeholder="t('purchase_addorder.remarks_text')"
          />
        </div>

        <!-- 采购链接 -->
        <div>
          <label
            class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300"
            >{{ t("purchase_addorder.link") }}</label
          >
          <textarea
            v-model="form.link"
            class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
            rows="2"
          />
        </div>

        <!-- 款式标签 -->
        <div>
          <label
            class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300"
            >{{ t("purchase_addorder.style_remarks") }}</label
          >
          <tagadder
            :placeholder="t('purchase_addorder.add_style')"
            v-model="form.styles"
          />
        </div>

        <!-- ==================== 费用明细表格 ==================== -->
        <div>
          <label
            class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300"
            >{{ t("purchase_addorder.cost") }}</label
          >

          <!-- 已添加的费用列表 -->
          <div v-if="costEntries.length" class="mb-4 overflow-x-auto">
            <table class="w-full text-left text-sm text-gray-900">
              <thead>
                <tr
                  class="border-b border-gray-200 bg-gray-50 text-gray-500 dark:border-dk-muted dark:bg-dk-base"
                >
                  <th class="px-3 py-2 font-medium">
                    {{ t("purchase_addorder.type") }}
                  </th>
                  <th class="px-3 py-2 font-medium">
                    {{ t("purchase_addorder.quantity") }}
                  </th>
                  <th class="px-3 py-2 font-medium">
                    {{ t("purchase_addorder.fee") }}
                  </th>
                  <th class="px-3 py-2 font-medium">
                    {{ t("purchase_addorder.total_price") }}
                  </th>
                  <th class="px-3 py-2 font-medium">
                    {{ t("purchase_addorder.currency") }}
                  </th>
                  <th class="px-3 py-2 font-medium">
                    {{ t("purchase_addorder.operation") }}
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="(item, idx) in costEntries"
                  :key="idx"
                  class="border-b border-gray-100 dark:border-dk-muted"
                >
                  <td
                    class="px-3 py-2 font-medium text-gray-900 dark:text-white"
                  >
                    {{ costType[item.type] }}
                  </td>
                  <td class="px-3 py-2 text-gray-500">{{ item.int }}</td>
                  <td class="px-3 py-2 text-gray-500">{{ item.cost }}</td>
                  <td class="px-3 py-2 text-gray-500">{{ item.costt }}</td>
                  <td class="px-3 py-2 text-gray-500">
                    {{ currencyOptions[item.currencytype] }}
                  </td>
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
              <label class="mb-1 block text-xs font-medium text-gray-500">{{
                t("purchase_addorder.fee_type")
              }}</label>
              <select
                v-model="newCost.type"
                class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
              >
                <option
                  v-for="(label, key) in costType"
                  :key="key"
                  :value="Number(key)"
                >
                  {{ label }}
                </option>
              </select>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-500">{{
                t("purchase_addorder.input_quantity")
              }}</label>
              <input
                v-model.number="newCost.int"
                type="number"
                class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
                min="1"
              />
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-500">{{
                t("purchase_addorder.input_fee")
              }}</label>
              <input
                v-model="newCost.cost"
                type="number"
                class="w-full rounded-lg border bg-white px-3 py-2 text-sm dark:bg-dk-base dark:text-white"
                :class="
                  costError
                    ? 'border-red-500'
                    : 'border-gray-300 dark:border-dk-muted'
                "
                step="0.01"
                min="0"
              />
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-500">{{
                t("purchase_addorder.select_currency")
              }}</label>
              <select
                v-model="newCost.currencytype"
                class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
              >
                <option
                  v-for="(label, key) in currencyOptions"
                  :key="key"
                  :value="Number(key)"
                >
                  {{ label }}
                </option>
              </select>
            </div>
            <div class="flex items-end">
              <button
                class="w-full rounded-lg border border-gray-300 bg-gray-50 px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-100 dark:border-dk-muted dark:bg-dk-base dark:text-gray-200 dark:hover:bg-dk-muted"
                @click="addCostEntry"
              >
                {{ t("purchase_addorder.add_cost") }}
              </button>
            </div>
          </div>
        </div>

        <!-- ==================== 图片上传 ==================== -->
        <div>
          <label
            class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300"
            >{{ t("purchase_addorder.upload_photos") }}</label
          >
          <useDropzone 
            ref="dropzoneRef"
            acceptFiles="image/*"
            uploadURL="/api/files/upload/image" 
            :initialFiles="form.photos" 
            :maxFiles="10" 
          />
        </div>
      </div>

      <!-- 底部操作栏 -->
      <div
        class="flex justify-end border-t border-gray-200 px-6 py-4 dark:border-dk-muted"
      >
        <button
          class="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500/20 disabled:opacity-60"
          :disabled="loading"
          @click="handleSubmit"
        >
          <svg
            v-if="loading"
            class="h-4 w-4 animate-spin text-white"
            viewBox="0 0 24 24"
            fill="none"
          >
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            />
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
            />
          </svg>
          {{ t("purchase.submit_changes") }}
        </button>
      </div>
    </div>
  </div>

  <!-- 通用确认弹窗 -->
  <ConfirmDialog
    v-model="showDeleteConfirm"
    :title="t('purchase.confirm_delete')"
    danger
    @confirm="doDelete"
  />
</template>
