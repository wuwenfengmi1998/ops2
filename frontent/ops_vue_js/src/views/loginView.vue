<script setup>
import { onMounted, watch, ref } from "vue";
import MyOffcanvas from "@/components/MyOffcanvas.vue";
import { useI18n } from "vue-i18n";
// 使用 vue-i18n 的 Composition API
const { t, locale } = useI18n();

const mos = ref();

const username = ref();
const password = ref();
const isRemember = ref();

const isShowPassword = ref(false);
function togglePasswordVisibility() {
  isShowPassword.value = !isShowPassword.value;
}

function login() {
  // 在这里处理登录逻辑

  const user = username.value?.value;
  const pass = password.value?.value;
  const remember = isRemember.value?.checked;

  username.value?.classList.remove("is-invalid");
  password.value?.classList.remove("is-invalid");

  if (!user || !pass) {
    if (!user) {
      username.value?.classList.add("is-invalid");
    }
    if (!pass) {
      password.value?.classList.add("is-invalid");
    }

    mos.value?.showAlert(
      "info",
      t("message.please_enter_username_and_password"),
      5000
    );
    return;
  }

  console.log("登录信息:", { user, pass, remember });
}

function functionupdataTitle() {
  document.title = "Operations." + t("appname.login");
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
  <div class="page page-center">
    <div class="container container-normal py-6">
      <div class="row align-items-center g-4">
        <div class="col-lg">
          <div class="container-tight">
            <div class="card card-md">
              <div class="card-body">
                <h2 class="h2 text-center mb-4">
                  {{ t("message.login_to_your_account") }}
                </h2>

                <div class="mb-3">
                  <label class="form-label">{{ t("message.user_name") }}</label>
                  <input
                    ref="username"
                    type="text"
                    class="form-control"
                    :placeholder="t('message.your_user_name')"
                    autocomplete="off"
                  />
                </div>
                <div class="mb-2">
                  <label class="form-label">
                    {{ t("message.password") }}
                    <span class="form-label-description">
                      <router-link to="/forgot_password">{{
                        t("message.i_forgot_password")
                      }}</router-link>
                    </span>
                  </label>
                  <div class="input-group input-group-flat">
                    <input
                      ref="password"
                      :type="isShowPassword ? 'text' : 'password'"
                      class="form-control"
                      :placeholder="t('message.your_password')"
                      autocomplete="off"
                    />
                    <span class="input-group-text">
                      <div
                        class="link-secondary"
                        :title="
                          isShowPassword
                            ? t('message.hidden_Password')
                            : t('message.show_password')
                        "
                        data-bs-toggle="tooltip"
                      >
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
                <div class="mb-2">
                  <label class="form-check">
                    <input
                      ref="isRemember"
                      type="checkbox"
                      class="form-check-input"
                    />
                    <span class="form-check-label">{{
                      t("message.remember_me_on_this_device")
                    }}</span>
                  </label>
                </div>
                <div class="form-footer">
                  <button @click="login" class="btn btn-primary w-100">
                    {{ t("button.sign_in") }}
                  </button>
                </div>
              </div>
            </div>
            <div class="text-center text-secondary mt-3">
              {{ t("message.dont_have_account_yet") }}
              <router-link to="/register">{{
                t("message.register_now")
              }}</router-link>
            </div>
          </div>
        </div>
        <div class="col-lg d-none d-lg-block">
          <img
            src="/static/illustrations/undraw_secure_login_pdn4.svg"
            height="300"
            class="d-block mx-auto"
            alt=""
          />
        </div>
      </div>
    </div>
  </div>

  <MyOffcanvas ref="mos" />
</template>
