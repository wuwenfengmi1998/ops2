<script setup>
import { onMounted, watch, ref ,defineProps} from "vue";
import { useI18n } from "vue-i18n";

import TomSelect from "tom-select";
import "tom-select/dist/css/tom-select.css";

const { t, locale } = useI18n();

const disable_backspace=ref()

function sele_tag_init() {
  new TomSelect(disable_backspace.value, {
    plugins: [ "remove_button"],
    
    persist: false,
	createOnBlur: true,
	create: true,

    // 自定义提示文本
        render: {
            no_results: function(data, escape) {
                return '<div class="no-results">'+t("tagadder.not_fund_item")+'</div>';
            },
            loading: function(data, escape) {
                return '<div class="loading">'+t("tagadder.loding")+'</div>';
            },
            option_create: function(data, escape) {
                return '<div class="create">'+t("tagadder.add")+'<strong>' + escape(data.input) + '</strong></div>';
            }
        }
  });
}

defineProps({
  placeholder: {
    type: String,
    default: "",
  },
})

onMounted(() => {
  sele_tag_init();
});
</script>

<template>
  
    <div ref="example_wrapper">
      <input type="text" ref="disable_backspace"  value="" autocomplete="off" :placeholder="placeholder"/>
    </div>
 
</template>
