<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Offcanvas } from 'bootstrap'

const offcanvasTop = ref(null)
let ov: Offcanvas
const alertType=ref() // 可选值：'success', 'warning', 'danger', 'info'
const alertText=ref()
let autoCloseTimeout:number
onMounted(() => {
  // 确保在组件挂载后初始化
  if (offcanvasTop.value) {
    ov = new Offcanvas(offcanvasTop.value,{
      backdrop: false,
    })
    //ov.show();
    //console.log('Offcanvas initialized:', ov)
  }
})


function showAlert(type: string, text: string) {
  alertText.value = text;

  alertType.value = type;
  //console.log(ov);
  if (ov) {
    ov.hide();
    ov.show();
    if(autoCloseTimeout)
    {
      clearTimeout(autoCloseTimeout);
    }
    autoCloseTimeout=setTimeout(() => {
        //console.log("timeout");
        ov.hide();
      }, 5000);


  }
}

defineExpose({
  showAlert
});


</script>

<style scoped>
.my_offcanvas_top {
  position: fixed;

  height: 45px;

  top: 0;
  right: 0;
  left: 0;

  margin-left: 20px;
  margin-right: 20px;
  margin-top: 20px;

  max-height: 100%;
  transform: translateY(-100%);
}
</style>

<template>
  <div
    class="offcanvas alert alert-important alert-dismissible my_offcanvas_top"
    :class="{
      'alert-success': alertType === 'success',
      'alert-warning': alertType === 'warning',
      'alert-danger': alertType === 'danger',
      'alert-info': alertType === 'info',
    }"
    role="alert"
    tabindex="-1"
    ref="offcanvasTop"
  >
    <div class="d-flex">
      <div>
        <!-- Download SVG icon from http://tabler-icons.io/i/check -->
        <svg
          v-if="alertType === 'success'"
          xmlns="http://www.w3.org/2000/svg"
          class="icon alert-icon"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          stroke-width="2"
          stroke="currentColor"
          fill="none"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
          <path d="M5 12l5 5l10 -10"></path>
        </svg>
        <svg
          v-if="alertType === 'warning'"
          xmlns="http://www.w3.org/2000/svg"
          class="icon alert-icon"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          stroke-width="2"
          stroke="currentColor"
          fill="none"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
          <path
            d="M10.24 3.957l-8.422 14.06a1.989 1.989 0 0 0 1.7 2.983h16.845a1.989 1.989 0 0 0 1.7 -2.983l-8.423 -14.06a1.989 1.989 0 0 0 -3.4 0z"
          ></path>
          <path d="M12 9v4"></path>
          <path d="M12 17h.01"></path>
        </svg>
        <svg
          v-if="alertType === 'danger'"
          xmlns="http://www.w3.org/2000/svg"
          class="icon alert-icon"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          stroke-width="2"
          stroke="currentColor"
          fill="none"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
          <path d="M3 12a9 9 0 1 0 18 0a9 9 0 0 0 -18 0"></path>
          <path d="M12 8v4"></path>
          <path d="M12 16h.01"></path>
        </svg>
        <svg
          v-if="alertType === 'info'"
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="icon icon-tabler icons-tabler-outline icon-tabler-brand-hipchat"
        >
          <path stroke="none" d="M0 0h24v24H0z" fill="none" />
          <path
            d="M17.802 17.292s.077 -.055 .2 -.149c1.843 -1.425 3 -3.49 3 -5.789c0 -4.286 -4.03 -7.764 -9 -7.764c-4.97 0 -9 3.478 -9 7.764c0 4.288 4.03 7.646 9 7.646c.424 0 1.12 -.028 2.088 -.084c1.262 .82 3.104 1.493 4.716 1.493c.499 0 .734 -.41 .414 -.828c-.486 -.596 -1.156 -1.551 -1.416 -2.29z"
          />
          <path d="M7.5 13.5c2.5 2.5 6.5 2.5 9 0" />
        </svg>
      </div>
      <div>
        {{ alertText }}
      </div>
    </div>
    <a class="btn-close" data-bs-dismiss="offcanvas" aria-label="close"></a>
  </div>

  <!-- <div>
    <button @click="showAlert('success','success')">success</button>
    <button @click="showAlert('warning','warning')">warning</button>
    <button @click="showAlert('danger','danger')">danger</button>
    <button @click="showAlert('info','info')">info</button>
  </div> -->
</template>
