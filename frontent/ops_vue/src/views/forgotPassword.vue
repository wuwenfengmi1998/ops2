<script setup lang="ts">
import { onMounted, watch, ref } from 'vue'
import MyOffcanvas from '@/components/MyOffcanvas.vue'
import { myfuncs } from '@/myfunc.ts'
import { useI18n } from 'vue-i18n'
// 使用 vue-i18n 的 Composition API
const { t, locale } = useI18n()
const email = ref<HTMLInputElement | null>(null)
const mos = ref<InstanceType<typeof MyOffcanvas> | null>(null)

function resetPassword() {
  // 在这里处理重置密码逻辑
  const emailValue = email.value?.value
  if (emailValue === undefined || emailValue.trim() === '') {
    mos.value?.showAlert('info', t('message.please_enter_your_email'), 5000)
    return
  }

  if (!myfuncs.isValidEmail(emailValue)) {
    mos.value?.showAlert('warning', t('message.this_not_email'), 5000)
    return
  }

  console.log('sending password reset to:', emailValue)
}

function functionupdataTitle() {
  document.title = 'Operations.' + t('appname.forgot_password')
}
onMounted(() => {
  functionupdataTitle()
})
// 监听语言变化，更新标题
watch(locale, () => {
  functionupdataTitle()
})
</script>
<template>
  <div class="container container-tight py-4">
    <div class="text-center mb-4">
      <a href="." class="navbar-brand navbar-brand-autodark">
        <img
          src="/static/logo.svg"
          width="110"
          height="32"
          alt="Tabler"
          class="navbar-brand-image"
        />
      </a>
    </div>
    <div class="card card-md">
      <div class="card-body">
        <h2 class="card-title text-center mb-4">{{ t('message.forgot_password') }}</h2>
        <p class="text-secondary mb-4">
          {{ t('message.enter_your_email_to_reset_password') }}
        </p>
        <div class="mb-3">
          <label class="form-label">{{ t('message.email_address') }}</label>
          <input
            ref="email"
            type="email"
            class="form-control"
            :placeholder="t('message.your_email_address')"
          />
        </div>
        <div class="form-footer">
          <button @click="resetPassword" class="btn btn-primary w-100">
            <!-- Download SVG icon from http://tabler-icons.io/i/mail -->
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
                d="M3 7a2 2 0 0 1 2 -2h14a2 2 0 0 1 2 2v10a2 2 0 0 1 -2 2h-14a2 2 0 0 1 -2 -2v-10z"
              />
              <path d="M3 7l9 6l9 -6" />
            </svg>
            {{ t('button.send_me_new_password') }}
          </button>
        </div>
      </div>
    </div>
    <div class="text-center text-secondary mt-3">
      <router-link to="/login">{{ t('message.back_to_login') }}</router-link>
    </div>
  </div>
  <MyOffcanvas ref="mos" />
</template>
