<script setup>
import { ref, reactive, computed, onMounted } from "vue";
import { useRoute, useRouter, RouterLink } from "vue-router";
import { useI18n } from "vue-i18n";
import { useToastStore } from "@/stores/toast";
import { usePageTitle } from "@/composables/usePageTitle";
import { purchaseApi } from "@/api/purchase";
import { useUserStore } from "@/stores/user";
import { useUsersStore } from "@/stores/users";
import {
  IconChevronLeft,
  IconExternalLink,
  IconPhoto,
  IconCheck,
  IconLoader2,
  IconX,
  IconUpload,
  IconTrash,
} from "@tabler/icons-vue";

usePageTitle("purchase.order_detail");
const { t, locale } = useI18n();
const route = useRoute();
const router = useRouter();
const toast = useToastStore();
const usersStore = useUsersStore();
const userStore = useUserStore();

const orderId = computed(() => parseInt(route.params.id));

const order = ref(null);
const costs = ref([]);
const photos = ref([]);
const commits = ref([]);
const loading = ref(true);
const notFound = ref(false);
const updatingStatus = ref(false);
const showStatusDialog = ref(false);
const pendingStatus = ref("");
const pendingComment = ref("");

// 状态变更附带的图片
const pendingPhotos = ref([]); // { hash, url, uploading, error }
const photoInputRef = ref(null);

// 状态选项
const statusOptions = [
  { value: "pending", labelKey: "status_pending", color: "yellow" },
  { value: "ordered", labelKey: "status_ordered", color: "blue" },
  { value: "arrived", labelKey: "status_arrived", color: "purple" },
  { value: "received", labelKey: "status_received", color: "green" },
  { value: "lost", labelKey: "status_lost", color: "red" },
  { value: "returned", labelKey: "status_returned", color: "gray" },
];

// 状态颜色映射
const statusColorClass = computed(() => ({
  pending:
    "bg-yellow-100 text-yellow-700 dark:bg-yellow-900/40 dark:text-yellow-400",
  ordered: "bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-400",
  arrived:
    "bg-purple-100 text-purple-700 dark:bg-purple-900/40 dark:text-purple-400",
  received:
    "bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400",
  lost:
    "bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400",
  returned:
    "bg-gray-200 text-gray-600 dark:bg-gray-700 dark:text-gray-300",
}));

// 货币选项
const currencyOptions = { 1: "CNY", 2: "MOP", 3: "HKD", 4: "USD" };

// 费用类型映射
const costTypeMap = computed(() => ({
  1: t("cost_type.unit_price"),
  2: t("cost_type.freight"),
}));

// 合计费用
const costTotalYuan = computed(() => {
  return (
    costs.value.reduce(
      (sum, c) => sum + (c.Price || 0) * (c.Quantity || 0),
      0,
    ) / 100
  );
});

// 按货币分组统计
const costsByCurrency = computed(() => {
  const groups = {};
  costs.value.forEach((c) => {
    const cur = currencyOptions[c.CurrencyType] || "Unknown";
    const amount =
      c.Price && c.Quantity
        ? ((c.Price * c.Quantity) / 100).toFixed(2)
        : "0.00";
    if (!groups[cur]) groups[cur] = 0;
    groups[cur] += parseFloat(amount);
  });
  return Object.entries(groups).map(([currency, total]) => ({
    currency,
    total: total.toFixed(2),
  }));
});

function formatDate(dateStr) {
  if (!dateStr) return "-";
  const d = new Date(dateStr);
  if (isNaN(d.getTime())) return "-";
  return new Intl.DateTimeFormat(locale.value, {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
  }).format(d);
}

function formatPrice(priceInCents) {
  if (!priceInCents) return "0.00";
  return (priceInCents / 100).toFixed(2);
}

function getPhotoUrl(file) {
  return `/api/files/get/${file.Sha256}`;
}

function openLink() {
  if (!order.value?.Link) return;
  let url = order.value.Link.trim();
  if (!/^https?:\/\//i.test(url)) url = "https://" + url;
  window.open(url, "_blank");
}

function getStatusLabel(status) {
  if (!status) return "";
  const opt = statusOptions.find((o) => o.value === status);
  return opt ? t("purchase." + opt.labelKey) || status : status;
}

function getStatusColorClass(status) {
  return statusColorClass.value[status] || "bg-gray-100 text-gray-600";
}

function openStatusDialog(newStatus) {
  if (newStatus === order.value?.OrderStatus) return;
  pendingStatus.value = newStatus;
  pendingComment.value = "";
  pendingPhotos.value = [];
  showStatusDialog.value = true;
}

function closeStatusDialog() {
  showStatusDialog.value = false;
  pendingStatus.value = "";
  pendingComment.value = "";
  pendingPhotos.value = [];
}

// 触发文件选择
function openPhotoPicker() {
  photoInputRef.value?.click();
}

// 选择文件后上传
async function handlePhotoChange(event) {
  const files = Array.from(event.target.files || []);
  if (!files.length) return;
  event.target.value = ""; // 清空，允许重复选同一文件

  for (const file of files) {
    if (pendingPhotos.value.length >= 10) break;
    const tempId = Date.now() + Math.random();
    const entry = {
      tempId,
      url: URL.createObjectURL(file),
      uploading: true,
      error: false,
      hash: null,
    };
    pendingPhotos.value.push(entry);

    try {
      const formData = new FormData();
      formData.append("file", file);
      formData.append("cookie", userStore.cookieValue);
      const res = await fetch("/api/files/upload/image", {
        method: "POST",
        body: formData,
      });
      const json = await res.json();
      if (json.errCode === 0 || json.return?.hash) {
        const p = pendingPhotos.value.find((p) => p.tempId === tempId);
        if (p) {
          p.hash = json.return.hash;
          p.uploading = false;
        }
      } else {
        const p = pendingPhotos.value.find((p) => p.tempId === tempId);
        if (p) {
          p.uploading = false;
          p.error = true;
        }
      }
    } catch {
      const p = pendingPhotos.value.find((p) => p.tempId === tempId);
      if (p) {
        p.uploading = false;
        p.error = true;
      }
    }
  }
}

// 移除待上传的图片
function removePendingPhoto(tempId) {
  const idx = pendingPhotos.value.findIndex((p) => p.tempId === tempId);
  if (idx !== -1) {
    const p = pendingPhotos.value[idx];
    if (p.url) URL.revokeObjectURL(p.url);
    pendingPhotos.value.splice(idx, 1);
  }
}

async function confirmStatusChange() {
  // 等所有图片上传完
  await new Promise((resolve) => setTimeout(resolve, 200));
  const stillUploading = pendingPhotos.value.some((p) => p.uploading);
  if (stillUploading) {
    toast.error("图片正在上传中，请稍候");
    return;
  }
  const photoHashes = pendingPhotos.value
    .filter((p) => !p.error)
    .map((p) => p.hash);

  updatingStatus.value = true;
  showStatusDialog.value = false;
  try {
    const { errCode } = await purchaseApi.updateOrderStatus(
      orderId.value,
      pendingStatus.value,
      pendingComment.value,
      photoHashes,
    );
    if (errCode === 0) {
      order.value.OrderStatus = pendingStatus.value;
      toast.success(t("message.save_success"));
      await fetchOrder();
    } else {
      toast.error(t("message.server_error"));
    }
  } catch {
    toast.error(t("message.server_error"));
  } finally {
    updatingStatus.value = false;
    pendingStatus.value = "";
    pendingComment.value = "";
  }
}

async function fetchOrder() {
  loading.value = true;
  try {
    const { errCode, data } = await purchaseApi.getOrder(orderId.value);
    if (errCode === 0 && data) {
      order.value = data.order ?? null;
      costs.value = data.costs ?? [];
      photos.value = data.photos ?? [];
      commits.value = data.commits ?? [];
    } else {
      notFound.value = true;
    }
  } catch {
    notFound.value = true;
  } finally {
    loading.value = false;
  }
}

onMounted(fetchOrder);
</script>

<template>
  
  <div class="mx-auto max-w-6xl px-6 py-6">
    <!-- 返回按钮 -->
    <div class="mb-4">
      <RouterLink
        to="/purchase"
        class="inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-dk-card dark:hover:text-gray-200"
      >
        <IconChevronLeft :size="16" />
        {{ t("purchase.back_to_list") }}
      </RouterLink>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-24">
      <svg
        class="h-8 w-8 animate-spin text-blue-500"
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
      <span class="ml-3 text-gray-500">Loading...</span>
    </div>

    <!-- Not Found -->
    <div
      v-else-if="notFound"
      class="rounded-xl border border-gray-200 bg-white py-16 text-center shadow-lg dark:border-dk-muted dark:bg-dk-card"
    >
      <p class="text-gray-400">{{ t("purchase.order_not_found") }}</p>
    </div>

    <!-- Order Detail -->
    <div v-else class="flex flex-col gap-6">
      <!-- Header -->
      <div
        class="rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card"
      >
        <div
          class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dk-muted"
        >
          <div class="flex items-center gap-3">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ t("purchase.order_detail") }} #{{ orderId }}
            </h2>
            <!-- 当前状态标签 -->
            <span
              v-if="order?.OrderStatus"
              class="inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-semibold"
              :class="getStatusColorClass(order.OrderStatus)"
            >
              <span
                v-if="updatingStatus"
                class="inline-flex items-center gap-1"
              >
                <IconLoader2 :size="10" class="animate-spin" />
              </span>
              {{ getStatusLabel(order.OrderStatus) }}
            </span>
            <!-- 创建者 -->
            <span
              v-if="order?.UserID"
              class="ml-1 flex items-center gap-1.5 rounded-full border border-gray-200 bg-gray-50 px-2 py-0.5 text-xs text-gray-500 dark:border-dk-muted dark:bg-dk-base dark:text-gray-400"
            >
              <img
                :src="usersStore.getAvatarUrlFromUserID(order.UserID)"
                class="rounded-full object-cover"
                style="width:18px;height:18px;"
              />
              {{ usersStore.getUsernameFromUserID(order.UserID) }}
            </span>
          </div>
          <span class="text-sm text-gray-400">{{
            formatDate(order?.CreatedAt)
          }}</span>
        </div>

        <!-- 状态快捷切换按钮 -->
        <div
          class="flex flex-wrap gap-2 border-b border-gray-100 px-6 py-3 dark:border-dk-muted"
        >
          <span class="text-sm text-gray-500 dark:text-gray-400"
            >{{ t("purchase.change_status") }}:</span
          >
          <button
            v-for="opt in statusOptions.slice(0, 4)"
            :key="opt.value"
            class="inline-flex items-center gap-1 rounded-full border px-3 py-1 text-xs font-medium transition-all"
            :class="
              order?.OrderStatus === opt.value
                ? [getStatusColorClass(opt.value), 'border-transparent']
                : 'border-gray-200 text-gray-500 hover:border-gray-300 hover:bg-gray-50 dark:border-dk-muted dark:text-gray-400 dark:hover:bg-dk-base'
            "
            :disabled="updatingStatus"
            @click="openStatusDialog(opt.value)"
          >
            <IconCheck v-if="order?.OrderStatus === opt.value" :size="12" />
            {{ t("purchase." + opt.labelKey) }}
          </button>
          <!-- 异常状态右对齐 -->
          <span class="flex-1" />
          <button
            v-for="opt in statusOptions.slice(4)"
            :key="opt.value"
            class="inline-flex items-center gap-1 rounded-full border px-3 py-1 text-xs font-medium transition-all"
            :class="
              order?.OrderStatus === opt.value
                ? [getStatusColorClass(opt.value), 'border-transparent']
                : 'border-gray-200 text-gray-500 hover:border-gray-300 hover:bg-gray-50 dark:border-dk-muted dark:text-gray-400 dark:hover:bg-dk-base'
            "
            :disabled="updatingStatus"
            @click="openStatusDialog(opt.value)"
          >
            <IconCheck v-if="order?.OrderStatus === opt.value" :size="12" />
            {{ t("purchase." + opt.labelKey) }}
          </button>
        </div>

        <!-- Order Info -->
        <div class="space-y-4 px-6 py-5">
          <h4 class="text-sm font-semibold text-gray-500 dark:text-gray-400">
            {{ t("purchase.order_info") }}
          </h4>

          <div class="grid gap-4 sm:grid-cols-2">
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-400">{{
                t("purchase_addorder.part_name")
              }}</label>
              <p class="font-medium text-gray-900 dark:text-white">
                {{ order?.Title || "-" }}
              </p>
            </div>

            <div>
              <label class="mb-1 block text-xs font-medium text-gray-400">{{
                t("purchase.link")
              }}</label>
              <div v-if="order?.Link" class="flex items-center gap-2">
                <p class="max-w-xs truncate text-blue-600 dark:text-blue-400">
                  {{ order.Link }}
                </p>
                <button
                  class="inline-flex items-center gap-1 rounded px-2 py-0.5 text-xs text-blue-600 hover:bg-blue-50 dark:text-blue-400 dark:hover:bg-blue-900/20"
                  @click="openLink"
                >
                  <IconExternalLink :size="14" />
                  {{ t("purchase.open_link") }}
                </button>
              </div>
              <p v-else class="text-gray-400">-</p>
            </div>

            <div v-if="order?.Styles">
              <label class="mb-1 block text-xs font-medium text-gray-400">{{
                t("purchase_addorder.style_remarks")
              }}</label>
              <p class="text-gray-700 dark:text-gray-200">{{ order.Styles }}</p>
            </div>

            <div v-if="order?.Remark" class="sm:col-span-2">
              <label class="mb-1 block text-xs font-medium text-gray-400">{{
                t("purchase_addorder.remarks")
              }}</label>
              <p class="whitespace-pre-wrap text-gray-700 dark:text-gray-200">
                {{ order.Remark }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- 费用明细 -->
      <div
        class="rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card"
      >
        <div class="border-b border-gray-100 px-6 py-4 dark:border-dk-muted">
          <h4 class="text-sm font-semibold text-gray-500 dark:text-gray-400">
            {{ t("purchase.cost_detail") }}
          </h4>
        </div>

        <div v-if="costs.length" class="overflow-x-auto px-0">
          <table class="w-full text-left text-sm">
            <thead>
              <tr
                class="border-b border-gray-100 bg-gray-50 text-gray-500 dark:border-dk-muted dark:bg-dk-base"
              >
                <th class="px-6 py-3 font-medium">
                  {{ t("purchase_addorder.fee_type") }}
                </th>
                <th class="px-6 py-3 font-medium">
                  {{ t("purchase_addorder.quantity") }}
                </th>
                <th class="px-6 py-3 font-medium">
                  {{ t("purchase.unit_price") }}
                </th>
                <th class="px-6 py-3 font-medium">
                  {{ t("purchase.total_price") }}
                </th>
                <th class="px-6 py-3 font-medium">
                  {{ t("purchase_addorder.currency") }}
                </th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(item, idx) in costs"
                :key="idx"
                class="border-b border-gray-50 dark:border-dk-muted/50"
              >
                <td class="px-6 py-3 font-medium text-gray-800 dark:text-white">
                  {{ costTypeMap[item.CostType] || item.CostType }}
                </td>
                <td class="px-6 py-3 text-gray-600 dark:text-gray-300">
                  {{ item.Quantity }}
                </td>
                <td class="px-6 py-3 text-gray-600 dark:text-gray-300">
                  {{ formatPrice(item.Price) }}
                </td>
                <td class="px-6 py-3 font-medium text-gray-900 dark:text-white">
                  {{ formatPrice(item.Price * item.Quantity) }}
                </td>
                <td class="px-6 py-3 text-gray-600 dark:text-gray-300">
                  {{ currencyOptions[item.CurrencyType] || "-" }}
                </td>
              </tr>
            </tbody>
            <tfoot>
              <tr class="bg-gray-50 dark:bg-dk-base">
                <td
                  colspan="3"
                  class="px-6 py-3 text-sm font-semibold text-gray-600 dark:text-gray-300"
                >
                  {{ t("purchase.cost_total") }}
                </td>
                <td
                  class="px-6 py-3 text-base font-bold text-gray-900 dark:text-white"
                >
                  {{ costTotalYuan.toFixed(2) }}
                </td>
                <td class="px-6 py-3">
                  <span
                    v-for="g in costsByCurrency"
                    :key="g.currency"
                    class="mr-2 inline-block rounded bg-blue-100 px-1.5 py-0.5 text-xs font-medium text-blue-700 dark:bg-blue-900/40 dark:text-blue-300"
                  >
                    {{ g.currency }} {{ g.total }}
                  </span>
                </td>
              </tr>
            </tfoot>
          </table>
        </div>

        <div v-else class="px-6 py-8 text-center text-gray-400">
          {{ t("purchase.no_costs") }}
        </div>
      </div>

      <!-- 图片备注 -->
      <div
        class="rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card"
      >
        <div class="border-b border-gray-100 px-6 py-4 dark:border-dk-muted">
          <h4 class="text-sm font-semibold text-gray-500 dark:text-gray-400">
            {{ t("purchase.photo_remarks") }}
          </h4>
        </div>

        <div v-if="photos.length" class="flex flex-wrap gap-3 px-6 py-5">
          <a
            v-for="photo in photos"
            :key="photo.ID"
            :href="getPhotoUrl(photo)"
            target="_blank"
            class="group relative block overflow-hidden rounded-lg border border-gray-200 dark:border-dk-muted"
            style="width: 120px; height: 120px"
          >
            <img
              :src="getPhotoUrl(photo)"
              :alt="photo.Name || 'photo'"
              class="h-full w-full object-cover transition-transform group-hover:scale-105"
            />
          </a>
        </div>

        <div
          v-else
          class="flex flex-col items-center justify-center gap-2 px-6 py-10 text-gray-400"
        >
          <IconPhoto :size="32" class="opacity-40" />
          <span class="text-sm">{{ t("purchase.no_photos") }}</span>
        </div>
      </div>

      <!-- 状态记录（Commit History） -->
      <div
        class="rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card"
      >
        <div class="border-b border-gray-100 px-6 py-4 dark:border-dk-muted">
          <h4 class="text-sm font-semibold text-gray-500 dark:text-gray-400">
            {{ t("purchase.commit_history") }}
          </h4>
        </div>

        <div
          v-if="commits.length"
          class="divide-y divide-gray-50 px-6 py-2 dark:divide-dk-muted/50"
        >
          <div
            v-for="commit in commits"
            :key="commit.id"
            class="flex items-start gap-3 py-3"
          >
            <!-- 左侧：头像 + 用户名 -->
            <div class="flex w-20 flex-shrink-0 flex-col items-center gap-1">
              <img
                :src="usersStore.getAvatarUrlFromUserID(commit.userId)"
                class="rounded-full border border-gray-200 object-cover dark:border-dk-muted"
                style="width:32px;height:32px;"
              />
              <span
                class="w-full truncate text-center"
                style="font-size:11px;color:var(--text-secondary,#6b7280);"
                :title="usersStore.getUsernameFromUserID(commit.userId)"
              >
                {{ usersStore.getUsernameFromUserID(commit.userId) }}
              </span>
            </div>

            <!-- 中间：时间线点 -->
            <div class="flex flex-shrink-0 flex-col items-center pt-1">
              <div
                class="h-3 w-3 rounded-full border-2 border-white dark:border-dk-base"
                :class="getStatusColorClass(commit.status)"
              />
            </div>

            <!-- 右侧：状态 + 备注 + 图片 -->
            <div class="flex-1 min-w-0">
              <div class="flex flex-wrap items-center gap-2">
                <span
                  class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium"
                  :class="getStatusColorClass(commit.status)"
                >
                  {{ getStatusLabel(commit.status) }}
                </span>
                <span
                  v-if="commit.action === 'create'"
                  class="text-xs text-gray-400"
                >
                  {{ t("purchase.commit_create") }}
                </span>
                <span v-else class="text-xs text-gray-400">
                  {{
                    commit.oldStatus
                      ? getStatusLabel(commit.oldStatus) + " → "
                      : ""
                  }}{{ getStatusLabel(commit.status) }}
                </span>
                <span class="ml-auto text-xs text-gray-400">{{
                  formatDate(commit.createdAt)
                }}</span>
              </div>
              <p
                v-if="commit.comment"
                class="mt-1 text-sm text-gray-600 dark:text-gray-300"
              >
                {{ commit.comment }}
              </p>
              <div
                v-if="commit.photos?.length"
                class="mt-2 flex flex-wrap gap-1.5"
              >
                <a
                  v-for="hash in commit.photos"
                  :key="hash"
                  :href="`/api/files/get/${hash}`"
                  target="_blank"
                  class="block overflow-hidden rounded border border-gray-200 dark:border-dk-muted transition-transform hover:scale-105"
                  style="width:48px;height:48px;"
                >
                  <img
                    :src="`/api/files/get/${hash}`"
                    class="h-full w-full object-cover"
                  />
                </a>
              </div>
            </div>
          </div>
        </div>

        <div v-else class="px-6 py-8 text-center text-sm text-gray-400">
          {{ t("purchase.no_commits") }}
        </div>
      </div>
    </div>
  </div>

  <!-- 状态变更弹窗 -->
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="showStatusDialog"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4"
        @click.self="closeStatusDialog"
      >
        <div
          class="w-full max-w-md rounded-xl border border-gray-200 bg-white shadow-2xl dark:border-dk-muted dark:bg-dk-card"
        >
          <!-- Header -->
          <div
            class="flex items-center justify-between border-b border-gray-100 px-5 py-4 dark:border-dk-muted"
          >
            <h3 class="font-semibold text-gray-900 dark:text-white">
              {{ t("purchase.change_status") }}
            </h3>
            <button
              class="rounded p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-600 dark:hover:bg-dk-base"
              @click="closeStatusDialog"
            >
              <IconX :size="18" />
            </button>
          </div>

          <!-- Body -->
          <div class="space-y-4 px-5 py-5">
            <!-- 新状态 -->
            <div>
              <label
                class="mb-2 block text-sm text-gray-500 dark:text-gray-400"
              >
                {{ t("purchase.status") }}
              </label>
              <div class="flex items-center gap-2">
                <span
                  class="inline-flex items-center gap-1 rounded-full px-3 py-1 text-sm font-semibold"
                  :class="getStatusColorClass(pendingStatus)"
                >
                  <IconLoader2
                    v-if="updatingStatus"
                    :size="12"
                    class="animate-spin"
                  />
                  <IconCheck v-else :size="12" />
                  {{ getStatusLabel(pendingStatus) }}
                </span>
              </div>
            </div>

            <!-- 变更备注 -->
            <div>
              <label
                class="mb-2 block text-sm text-gray-500 dark:text-gray-400"
              >
                {{ t("purchase.change_remark") }}
              </label>
              <textarea
                v-model="pendingComment"
                rows="3"
                class="w-full rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm text-gray-900 placeholder-gray-400 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 disabled:bg-gray-50 dark:border-dk-muted dark:bg-dk-base dark:text-white dark:placeholder-gray-500"
                :placeholder="t('purchase.commit_placeholder')"
                :disabled="updatingStatus"
                @keydown.ctrl.enter="confirmStatusChange"
              />
            </div>

            <!-- 变更图片 -->
            <div>
              <label
                class="mb-2 block text-sm text-gray-500 dark:text-gray-400"
              >
                {{ t("purchase.status_photos") }}
              </label>

              <!-- 隐藏的文件选择框 -->
              <input
                ref="photoInputRef"
                type="file"
                accept="image/*"
                multiple
                class="hidden"
                @change="handlePhotoChange"
              />

              <!-- 已选图片预览 -->
              <div
                v-if="pendingPhotos.length"
                class="mb-2 flex flex-wrap gap-2"
              >
                <div
                  v-for="p in pendingPhotos"
                  :key="p.tempId"
                  class="group relative overflow-hidden rounded-lg border border-gray-200 dark:border-dk-muted"
                  style="width: 60px; height: 60px"
                >
                  <img :src="p.url" class="h-full w-full object-cover" />
                  <!-- 上传中遮罩 -->
                  <div
                    v-if="p.uploading"
                    class="absolute inset-0 flex items-center justify-center bg-black/40"
                  >
                    <IconLoader2 :size="16" class="animate-spin text-white" />
                  </div>
                  <!-- 失败遮罩 -->
                  <div
                    v-else-if="p.error"
                    class="absolute inset-0 flex items-center justify-center bg-red-500/60"
                  >
                    <IconX :size="14" class="text-white" />
                  </div>
                  <!-- 移除按钮 -->
                  <button
                    v-if="!p.uploading"
                    class="absolute -top-1 -right-1 flex h-4 w-4 items-center justify-center rounded-full bg-red-500 text-white opacity-0 transition-opacity group-hover:opacity-100"
                    @click="removePendingPhoto(p.tempId)"
                  >
                    <IconX :size="10" />
                  </button>
                </div>
              </div>

              <!-- 上传按钮 -->
              <button
                class="inline-flex items-center gap-1.5 rounded-lg border border-dashed border-gray-300 px-3 py-2 text-sm text-gray-500 transition-colors hover:border-blue-400 hover:text-blue-500 disabled:opacity-50 dark:border-dk-muted dark:text-gray-400 dark:hover:border-blue-500 dark:hover:text-blue-400"
                :disabled="
                  updatingStatus ||
                  pendingPhotos.filter((p) => !p.uploading).length >= 10
                "
                @click="openPhotoPicker"
              >
                <IconUpload :size="14" />
                {{ t("purchase.upload_photos") }}
              </button>
            </div>
          </div>

          <!-- Footer -->
          <div
            class="flex justify-end gap-3 border-t border-gray-100 px-5 py-4 dark:border-dk-muted"
          >
            <button
              class="rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-50 disabled:opacity-50 dark:border-dk-muted dark:text-gray-300 dark:hover:bg-dk-base"
              :disabled="updatingStatus"
              @click="closeStatusDialog"
            >
              {{ t("settings.cancel") }}
            </button>
            <button
              class="flex items-center gap-2 rounded-lg bg-blue-500 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-600 disabled:opacity-50"
              :disabled="updatingStatus"
              @click="confirmStatusChange"
            >
              <IconLoader2
                v-if="updatingStatus"
                :size="14"
                class="animate-spin"
              />
              {{ t("message.submit") }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
