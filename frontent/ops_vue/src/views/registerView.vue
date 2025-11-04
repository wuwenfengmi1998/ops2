<script setup lang="ts">
import { onMounted, watch, ref } from 'vue'
import MyOffcanvas from '@/components/MyOffcanvas.vue'
import { myfuncs } from '@/myfunc.ts'
import { useI18n } from 'vue-i18n'
// 使用 vue-i18n 的 Composition API
const { t, locale } = useI18n()
const mos = ref<InstanceType<typeof MyOffcanvas> | null>(null)
const isShowPassword = ref(false)
const username = ref<HTMLInputElement | null>(null)
const useremail = ref<HTMLInputElement | null>(null)
const userpassword = ref<HTMLInputElement | null>(null)


function functionupdataTitle() {
  document.title = 'Operations.' + t('appname.register')
}
function togglePasswordVisibility() {
  isShowPassword.value = !isShowPassword.value
}
function createAccount() {
  // 在这里处理创建新账户的逻辑
  const user = username.value?.value
  const email = useremail.value?.value
  const pass = userpassword.value?.value

  username.value?.classList.remove('is-invalid');
  useremail.value?.classList.remove('is-invalid');
  userpassword.value?.classList.remove('is-invalid');

  if (
    !user ||
    !email ||
    !pass
  ) {

    if(!user){
      username.value?.classList.add('is-invalid');
    }
    if(!email){
      useremail.value?.classList.add('is-invalid');
    }
    if(!pass){
      userpassword.value?.classList.add('is-invalid');
    }

    mos.value?.showAlert('info', t('message.please_enter_username_and_password'), 5000)
    return
  }
  if (!myfuncs.isValidEmail(email)) {
    useremail.value?.classList.add('is-invalid');
    mos.value?.showAlert('warning', t('message.this_not_email'), 5000)
    return
  }
console.log('创建新账户信息:', {
    user: username.value?.value,
    email: useremail.value?.value,
    pass: userpassword.value?.value,
  })
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
  <div class="page page-center">
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
          <h2 class="card-title text-center mb-4">{{ t('message.create_new_account') }}</h2>
          <div class="mb-3">
            <label class="form-label">{{ t('message.user_name') }}</label>
            <input
              ref="username"
              type="text"
              class="form-control"
              :placeholder="t('message.your_user_name')"
            />
          </div>
          <div class="mb-3">
            <label class="form-label">{{ t('message.email_address') }}</label>
            <input
              ref="useremail"
              type="email"
              class="form-control"
              :placeholder="t('message.your_email_address')"
            />
          </div>
          <div class="mb-3">
            <label class="form-label">{{ t('message.password') }}</label>
            <div class="input-group input-group-flat">
              <input
                ref="userpassword"
                :type="isShowPassword ? 'text' : 'password'"
                class="form-control"
                :placeholder="t('message.your_password')"
                autocomplete="off"
              />
              <span class="input-group-text">
                <div class="link-secondary" title="Show password" data-bs-toggle="tooltip">
                  <!-- Download SVG icon from http://tabler-icons.io/i/eye -->
                  <svg
                    v-if="!isShowPassword"
                    @click="togglePasswordVisibility"
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
                    <path d="M10 12a2 2 0 1 0 4 0a2 2 0 0 0 -4 0" />
                    <path
                      d="M21 12c-2.4 4 -5.4 6 -9 6c-3.6 0 -6.6 -2 -9 -6c2.4 -4 5.4 -6 9 -6c3.6 0 6.6 2 9 6"
                    />
                  </svg>
                  <svg
                    v-if="isShowPassword"
                    @click="togglePasswordVisibility"
                    xmlns="http://www.w3.org/2000/svg"
                    class="icon"
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                    <path d="M10.585 10.587a2 2 0 0 0 2.829 2.828" />
                    <path
                      d="M16.681 16.673a8.717 8.717 0 0 1 -4.681 1.327c-3.6 0 -6.6 -2 -9 -6c1.272 -2.12 2.712 -3.678 4.32 -4.674m2.86 -1.146a9.055 9.055 0 0 1 1.82 -.18c3.6 0 6.6 2 9 6c-.666 1.11 -1.379 2.067 -2.138 2.87"
                    />
                    <path d="M3 3l18 18" />
                  </svg>
                </div>
              </span>
            </div>
          </div>
          <!-- <div class="mb-3">
              <label class="form-check">
                <input type="checkbox" class="form-check-input"/>
                <span class="form-check-label">Agree the <a href="./terms-of-service.html" tabindex="-1">terms and policy</a>.</span>
              </label>
            </div> -->
          <div class="form-footer">
            <button @click="createAccount" class="btn btn-primary w-100">
              {{ t('message.create_new_account') }}
            </button>
          </div>
        </div>
      </div>
      <div class="text-center text-secondary mt-3">
        {{ t('message.already_have_an_account') }}
        <router-link to="/login">{{ t('message.back_to_login') }}</router-link>
      </div>
    </div>
  </div>
  <MyOffcanvas ref="mos" />
</template>
