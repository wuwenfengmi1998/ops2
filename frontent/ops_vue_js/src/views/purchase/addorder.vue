<script setup>
import { onMounted, watch, ref } from "vue";
import { useI18n } from "vue-i18n";

const { t, locale } = useI18n();

function functionupdataTitle() {
  document.title = "Operations." + t("purchase.add_part");
}
onMounted(() => {
  functionupdataTitle();
});
// 监听语言变化，更新标题
watch(locale, () => {
  functionupdataTitle();
}); 
</script>

<template>
add_part
</template>