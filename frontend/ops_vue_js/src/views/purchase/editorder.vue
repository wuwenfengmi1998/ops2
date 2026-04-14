<script setup>
/**
 * 采购订单编辑页面
 *
 * 功能概述：
 * - 通过路由参数 :id 加载已有订单数据
 * - 使用 PurchaseOrderForm 组件展示可编辑表单
 * - 提交时调用 /purchase/updateorder 保存修改
 */

import { reactive, ref, onMounted } from "vue";
import { useI18n } from "vue-i18n";
import { useRoute, useRouter } from "vue-router";
import { useToastStore } from "@/stores/toast";
import { usePageTitle } from "@/composables/usePageTitle";
import { useValidation } from "@/composables";
import { purchaseApi } from "@/api/purchase";
import PurchaseOrderForm from "@/components/PurchaseOrderForm.vue";

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

/** 回填的费用明细（分为单位，传给 PurchaseOrderForm） */
const initialCosts = ref([]);
/** 回填的图片列表 */
const initialPhotos = ref([]);

/** 表单数据 */
const form = reactive({
  title: "",
  remark: "",
  link: "",
  styles: "",
  photos: [],
  costs: [],
  _costs: [], // 由 PurchaseOrderForm 组件同步的分为单位费用数组
});

/** PurchaseOrderForm 组件引用（用于获取图片哈希） */
const formRef = ref(null);

// ==================== 加载订单数据 ====================
onMounted(async () => {
  if (!orderId) {
    pageError.value = t("purchase.order_not_found");
    pageLoading.value = false;
    return;
  }

  try {
    const res = await purchaseApi.getOrder(orderId);
    console.log(res)
    if (res.errCode !== 0 || res.raw?.err_code !== 0) {
      pageError.value = t("purchase.order_not_found");
      pageLoading.value = false;
      return;
    }

    const { order, costs, photos } = res.raw.data;

    // 回填基本信息
    form.title = order.Title ?? "";
    form.remark = order.Remark ?? "";
    form.link = order.Link ?? "";
    form.styles = order.Styles ?? "";

    // 回填费用（传给子组件，由子组件转换为元展示）
    initialCosts.value = costs ?? [];
    // 回填图片
    initialPhotos.value = photos ?? [];
  } catch {
    pageError.value = t("purchase.order_not_found");
  } finally {
    pageLoading.value = false;
  }
});

// ==================== 提交 ====================
async function handleSubmit() {
  clearErrors();
  const ok = validate("title", form.title, t("purchase_addorder.title"));
  if (!ok) return;

  // 获取图片哈希
  form.photos = formRef.value?.getPhotoHashes() ?? [];
  // 使用子组件同步的费用（分为单位）
  form.costs = form._costs ?? [];

  loading.value = true;
  try {
    const res = await purchaseApi.updateOrder(orderId, {
      title: form.title,
      remark: form.remark,
      link: form.link,
      styles: form.styles,
      photos: form.photos,
      costs: form.costs,
    });

    if (res.errCode === 0 && res.raw?.err_code === 0) {
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
    <div v-if="pageLoading" class="flex items-center justify-center py-20 text-gray-400">
      <svg class="mr-2 h-5 w-5 animate-spin" viewBox="0 0 24 24" fill="none">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
      </svg>
      {{ t("message.loading") }}
    </div>

    <!-- 订单不存在 -->
    <div v-else-if="pageError" class="rounded-xl border border-red-200 bg-red-50 px-6 py-10 text-center text-red-500 dark:border-red-800 dark:bg-red-900/20">
      {{ pageError }}
    </div>

    <!-- 主卡片 -->
    <div
      v-else
      class="flex flex-col gap-0 rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card"
    >
      <!-- 顶部标题栏 -->
      <div class="flex items-center justify-between border-b border-gray-200 px-6 py-4 dark:border-dk-muted">
        <h4 class="text-sm font-semibold text-gray-900 dark:text-white">
          {{ t("purchase_addorder.edit_order") }}
        </h4>
        <!-- 返回按钮 -->
        <button
          class="flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-dk-base"
          @click="router.back()"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
          </svg>
          {{ t("purchase.back_to_list") }}
        </button>
      </div>

      <!-- 错误提示（字段验证） -->
      <div v-if="errors.title" class="mx-6 mt-4 rounded-lg bg-red-50 px-4 py-2 text-sm text-red-600 dark:bg-red-900/20 dark:text-red-400">
        {{ errors.title }}
      </div>

      <!-- 表单主体（公共组件） -->
      <PurchaseOrderForm
        v-model="form"
        :initialCosts="initialCosts"
        :initialPhotos="initialPhotos"
        ref="formRef"
      />

      <!-- 底部操作栏 -->
      <div class="flex justify-end border-t border-gray-200 px-6 py-4 dark:border-dk-muted">
        <button
          class="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500/20 disabled:opacity-60"
          :disabled="loading"
          @click="handleSubmit"
        >
          <svg v-if="loading" class="h-4 w-4 animate-spin text-white" viewBox="0 0 24 24" fill="none">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
          {{ t("message.save_ok") }}
        </button>
      </div>
    </div>
  </div>
</template>
