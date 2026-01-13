<script setup>
import { onMounted, watch, ref } from "vue";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";
const { t, locale } = useI18n();

function functionupdataTitle() {
  document.title = "Operations." + t("errorpage.404_title");
}

// 监听语言变化，更新标题
watch(locale, () => {
  functionupdataTitle();
});

onMounted(() => {
  //console.log("account mounted");
  //username.value.value="Kevin";
  functionupdataTitle();
});
</script>
<template>
  <div>
    <div class="page page-center">
      <div class="container-tight py-4">
        <div class="empty">
          <div class="empty-header">404</div>
          <p class="empty-title">{{ t("errorpage.404_msg_title") }}</p>
          <p class="empty-subtitle text-secondary">
            {{ t("errorpage.404_msg") }}
          </p>
          <div class="empty-action">
            <router-link to="/" class="btn btn-primary">
              <!-- Download SVG icon from http://tabler-icons.io/i/arrow-left -->
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="icon"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                stroke-width="2"
                stroke="currentColor"
                fill="none"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                <path d="M5 12l14 0" />
                <path d="M5 12l6 6" />
                <path d="M5 12l6 -6" />
              </svg>
              {{ t("errorpage.404_back_home") }}
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
