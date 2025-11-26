<script setup>
import { onMounted, ref, watch ,defineProps} from "vue";
import Litepicker from "litepicker";
import { useI18n } from "vue-i18n";
const { t, locale } = useI18n();

const datepicker = ref(null);
var picker = null

watch(locale, () => {
  picker?.setOptions({ lang: locale.value });
});

defineProps({
  setdef: {
    type: String,
    default: "",
  },
})
onMounted(() => {
  // @formatter:off

  picker = new Litepicker({
    element: datepicker.value,
    lang: locale.value,
    firstDay: 0,
    format: "YYYY-MM-DD", // 日期格式

    dropdowns: {
      minYear: 1900, // 最小可选年份
      maxYear: new Date().getFullYear() + 1, // 最大为当前年份
      months: true, // 显示月份下拉
      years: true, // 显示年份下拉
    },
    //inlineMode: true,
  });
});

defineExpose({
  datepicker,
  
});
</script>

<template>
  
  <div class="input-icon">
    <span class="input-icon-addon"
      ><!-- Download SVG icon from http://tabler-icons.io/i/calendar -->
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
        <path
          d="M4 7a2 2 0 0 1 2 -2h12a2 2 0 0 1 2 2v12a2 2 0 0 1 -2 2h-12a2 2 0 0 1 -2 -2v-12z"
        />
        <path d="M16 3v4" />
        <path d="M8 3v4" />
        <path d="M4 11h16" />
        <path d="M11 15h1" />
        <path d="M12 15v3" />
      </svg>
    </span>
    <input
      class="form-control"
      :placeholder="t('message.select_date')"
      ref="datepicker"
      :value="setdef"
    />
  </div>
</template>
