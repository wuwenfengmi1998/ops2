<script setup>
import { onMounted, ref } from "vue";
import Litepicker from "litepicker";

const datepicker_inline = ref(null);

onMounted(() => {
  // @formatter:off

  new Litepicker({
    element: datepicker_inline.value,
    lang: 'zh-cn',
    firstDay: 0,
    buttonText: {
      previousMonth: `<!-- Download SVG icon from http://tabler-icons.io/i/chevron-left -->
    <svg xmlns="http://www.w3.org/2000/svg" class="icon" width="24" height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M15 6l-6 6l6 6" /></svg>`,
      nextMonth: `<!-- Download SVG icon from http://tabler-icons.io/i/chevron-right -->
    <svg xmlns="http://www.w3.org/2000/svg" class="icon" width="24" height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M9 6l6 6l-6 6" /></svg>`,
    },
    inlineMode: true,
  });

  // @formatter:on
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
      placeholder="Select a date"
      id="datepicker-icon-prepend"
      value="2020-06-20"
    />
  </div>

  <div class="mb-3">
    <label class="form-label">Inline datepicker</label>
    <div class="datepicker-inline" ref="datepicker_inline"></div>
  </div>
</template>
