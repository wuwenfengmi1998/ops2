<script setup>
/**
 * ConfirmDialog —— 通用确认弹窗
 *
 * 使用方式：
 *   const ok = await confirmDialog({
 *     title: "提示标题",
 *     message: "确认要删除吗？",
 *     confirmText: "删除",
 *     cancelText: "取消",
 *     danger: true,  // 红色确认按钮
 *   })
 *   if (ok) { ... }
 *
 * 或者作为组件使用 v-model：
 *   <ConfirmDialog v-model="show" @confirm="..." @cancel="..." />
 */
import { ref, watch } from "vue";
import { useI18n } from "vue-i18n";

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
    default: "",
  },
  message: {
    type: String,
    default: "",
  },
  confirmText: {
    type: String,
    default: "",
  },
  cancelText: {
    type: String,
    default: "",
  },
  danger: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["update:modelValue", "confirm", "cancel"]);

const { t } = useI18n();

function close() {
  emit("update:modelValue", false);
  emit("cancel");
}

function confirm() {
  emit("update:modelValue", false);
  emit("confirm");
}

watch(
  () => props.modelValue,
  (val) => {
    if (val) {
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "";
    }
  },
);
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="modelValue"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/40"
        @click.self="close"
      >
        <div
          class="min-w-[320px] max-w-sm rounded-xl border border-gray-200 bg-white shadow-2xl dark:border-dk-muted dark:bg-dk-card"
        >
          <!-- 标题 -->
          <div class="flex items-center justify-between border-b border-gray-100 px-5 py-4 dark:border-dk-muted">
            <h3 class="text-base font-semibold text-gray-900 dark:text-white">
              {{ title || t("message.confirm") }}
            </h3>
            <button
              class="ml-4 flex-shrink-0 rounded p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-600 dark:hover:bg-dk-muted dark:hover:text-gray-200"
              @click="close"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- 内容 -->
          <div class="px-5 py-4">
            <p class="text-sm leading-relaxed text-gray-600 dark:text-gray-300">
              {{ message }}
            </p>
          </div>

          <!-- 操作栏 -->
          <div class="flex justify-end gap-3 border-t border-gray-100 px-5 py-4 dark:border-dk-muted">
            <button
              class="rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-base dark:text-gray-200 dark:hover:bg-dk-muted"
              @click="close"
            >
              {{ cancelText || t("message.cancel") }}
            </button>
            <button
              class="rounded-lg px-4 py-2 text-sm font-semibold text-white transition-colors"
              :class="
                danger
                  ? 'bg-red-500 hover:bg-red-600'
                  : 'bg-blue-600 hover:bg-blue-700'
              "
              @click="confirm"
            >
              {{ confirmText || t("message.confirm") }}
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
