<script setup>
import { onMounted, watch, ref } from "vue";
import settingNavigation from "@/components/settingNavigation.vue";
import { useI18n } from "vue-i18n";
import datePicker from "@/components/datePicker.vue";
import imageCropper from "@/components/imageCropper.vue";
import { useUserStore } from "@/stores/user";
import { my_network_func } from "@/my_network_func";
import MyOffcanvas from "@/components/MyOffcanvas.vue";
import { myfuncs } from "@/myfunc";
import { useRouter } from "vue-router";
const mos = ref();

const { t, locale } = useI18n();

const router = useRouter();

const birthday = ref();
const username = ref();
const userremark = ref();

const emailInput = ref();

const userStore = useUserStore();

const oldPassInput = ref();
const newPassInput = ref();
const cnfPassInput = ref();

const isShowPassword = ref(false);

function togglePasswordVisibility() {
  isShowPassword.value = !isShowPassword.value;
}

function updataInfo() {
  let isDataErr = false;

  let birthdayValue = birthday.value.datepicker.value;
  let usernameValue = username.value.value;
  let userremarkValue = userremark.value.value;

  username.value?.classList.remove("is-invalid");
  userremark.value?.classList.remove("is-invalid");
  birthday.value?.datepicker.classList.remove("is-invalid");

  if (!usernameValue) {
    isDataErr = true;
    username.value?.classList.add("is-invalid");
  }

  if (!userremarkValue) {
    isDataErr = true;
    userremark.value?.classList.add("is-invalid");
  }

  if (!birthdayValue) {
    isDataErr = true;
    birthday.value?.datepicker.classList.add("is-invalid");
  }

  if (isDataErr) {
    //console.log("用户信息有误，无法保存");
    return;
  }
  // console.log("保存用户信息");
  // console.log("用户名:", usernameValue);
  // console.log("备注:", userremarkValue);
  // console.log("生日:", birthdayValue);
  my_network_func.postJson(
    "/users/updateInfo",
    {
      username: usernameValue,
      remark: userremarkValue,
      birthday: birthdayValue,
    },
    (r) => {
      //console.log(r);
      switch (r.statusCode) {
        case 200:
          switch (r.data.err_code) {
            case 0:
              mos.value?.showAlert("success", t("message.save_ok"), 1000);
              // 更新用户信息到store
              userStore.getUserInfoFromCookie();
              break;
            default:
              mos.value?.showAlert("danger", t("message.server_error"), 5000);
              break;
          }
          break;
        default:
          mos.value?.showAlert("danger", t("message.network_err"), 5000);
          break;
      }
    }
  );
}

function changeEmail() {
  if (emailInput.value.value == "") {
    emailInput.value.classList.add("is-invalid");
    mos.value?.showAlert("warning", t("message.please_enter_your_email"), 3000);
    return;
  } else {
    emailInput.value.classList.remove("is-invalid");
  }
  //判断是否是合法邮箱
  if (myfuncs.isValidEmail(emailInput.value.value) == false) {
    emailInput.value.classList.add("is-invalid");
    mos.value?.showAlert("danger", t("message.this_not_email"), 3000);
    return;
  } else {
    emailInput.value.classList.remove("is-invalid");
  }

  my_network_func.postJson(
    "/users/changeEmail",
    {
      newemail: emailInput.value.value,
    },
    (r) => {
      switch (r.statusCode) {
        case 200:
          switch (r.data.err_code) {
            case 0:
              mos.value?.showAlert("success", t("message.change_ok"), 5000);
              // 更新用户信息到store
              userStore.getUserInfoFromCookie();
              break;
            case -43:
              emailInput.value.classList.add("is-invalid");
              mos.value?.showAlert("danger", t("message.this_not_email"), 3000);
              break;

            default:
              mos.value?.showAlert("danger", t("message.server_error"), 5000);
              break;
          }
          break;
        default:
          mos.value?.showAlert("danger", t("message.network_err"), 5000);
          break;
      }
    }
  );
}

function changePassword() {
  let isDataErr = false;

  let oldPass = oldPassInput.value.value;
  let newPass = newPassInput.value.value;
  let cnfPass = cnfPassInput.value.value;

  oldPassInput.value.classList.remove("is-invalid");
  newPassInput.value.classList.remove("is-invalid");
  cnfPassInput.value.classList.remove("is-invalid");

  if (!oldPass) {
    isDataErr = true;
    oldPassInput.value.classList.add("is-invalid");
  }
  if (!newPass) {
    isDataErr = true;
    newPassInput.value.classList.add("is-invalid");
  }

  if (!cnfPass) {
    isDataErr = true;
    cnfPassInput.value.classList.add("is-invalid");
  }

  if (newPass !== cnfPass) {
    isDataErr = true;
    newPassInput.value.classList.add("is-invalid");
    cnfPassInput.value.classList.add("is-invalid");
    mos.value?.showAlert(
      "warning",
      t("message.confirm_password_incorrect"),
      3000
    );
  }
  if (isDataErr) {
    return;
  }

  my_network_func.postJson(
    "/users/changePassword",
    {
      oldpass: oldPass,
      newpass: newPass,
    },
    (r) => {
      switch (r.statusCode) {
        case 200:
          switch (r.data.err_code) {
            case 0:
              // 清空输入框
              oldPassInput.value.value = "";
              newPassInput.value.value = "";
              cnfPassInput.value.value = "";
              mos.value?.showAlert(
                "success",
                t("message.change_ok"),
                2000,
                () => {
                  userStore.logout();
                  router.push("/");
                }
              );

              break;
            case -42:
              oldPassInput.value.classList.add("is-invalid");
              mos.value?.showAlert(
                "danger",
                t("message.old_pass_incorrect"),
                3000
              );
              break;

            default:
              mos.value?.showAlert("danger", t("message.server_error"), 5000);
              break;
          }
          break;
        default:
          mos.value?.showAlert("danger", t("message.network_err"), 5000);
          break;
      }
    }
  );
}

function changeAvatar() {
  mos.value?.showAlert(
    "info",
    t("message.functionality_not_yet_developed"),
    5000
  );
}

function functionupdataTitle() {
  document.title = "Operations." + t("settings.account_settings");
}

// 监听语言变化，更新标题
watch(locale, () => {
  functionupdataTitle();
});

onMounted(() => {
  //console.log("account mounted");
  //username.value.value="Kevin";
  functionupdataTitle();

  if (!userStore.isLoggedIn) {
    router.push("/login");
  }
});
</script>

<template>
  <div class="page-wrapper">
    <div class="page-header d-print-none">
      <div class="container-xl">
        <h2 class="page-title">{{ t("settings.account_settings") }}</h2>
      </div>
    </div>

    <div class="page-body">
      <div class="container-xl">
        <div class="card">
          <div class="row g-0">
            <settingNavigation />
            <div class="col-12 col-md-9 d-flex flex-column">
              <div class="card-body">
                <h2 class="mb-4">{{ t("settings.my_account") }}</h2>
                <h3 class="card-title">
                  {{ t("settings.profile_information") }}
                </h3>
                <div class="row align-items-center">
                  <div class="col-auto">
                    <img
                      :src="userStore.getUserAvatarPath()"
                      alt=""
                      class="avatar avatar-xl"
                    />
                  </div>
                  <!-- <imageCropper /> -->
                  <div class="col-auto">
                    <button class="btn" @click="changeAvatar">
                      {{ t("settings.change_avatar") }}
                    </button>
                  </div>
                </div>
                <h3 class="card-title mt-4">-</h3>
                <div class="row g-3">
                  <div class="col-md">
                    <div class="form-label">{{ t("settings.name") }}</div>
                    <input
                      ref="username"
                      type="text"
                      class="form-control"
                      :value="
                        userStore.userInfo ? userStore.userInfo.Username : ''
                      "
                    />
                  </div>
                  <div class="col-md">
                    <div class="form-label">{{ t("settings.remark") }}</div>
                    <input
                      ref="userremark"
                      type="text"
                      class="form-control"
                      :value="
                        userStore.userInfo ? userStore.userInfo.FirstName : ''
                      "
                    />
                  </div>
                  <div class="col-md">
                    <div class="form-label">{{ t("settings.birthday") }}</div>
                    <datePicker
                      ref="birthday"
                      :setdef="userStore.getUserBirthday()"
                    />
                  </div>
                  <div>
                    <button class="btn" @click="updataInfo">
                      {{ t("settings.save_changes") }}
                    </button>
                  </div>
                </div>
                <h3 class="card-title mt-4">{{ t("settings.email") }}</h3>
                <div>
                  <div class="row g-2">
                    <div class="col-auto">
                      <input
                        ref="emailInput"
                        type="text"
                        class="form-control w-auto"
                        :value="userStore.user.Email"
                        :placeholder="t('message.your_email_address')"
                      />
                    </div>
                    <div class="col-auto">
                      <button class="btn" @click="changeEmail">
                        {{ t("settings.change_email") }}
                      </button>
                    </div>
                  </div>
                </div>
                <h3 class="card-title mt-4">
                  {{ t("settings.password") }}
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
                </h3>

                <input
                  ref="oldPassInput"
                  :type="isShowPassword ? 'text' : 'password'"
                  class="form-control w-auto mb-1"
                  :placeholder="t('message.type_old_pass')"
                />

                <input
                  ref="newPassInput"
                  :type="isShowPassword ? 'text' : 'password'"
                  class="form-control w-auto mb-1"
                  :placeholder="t('message.type_new_pass')"
                />
                <input
                  ref="cnfPassInput"
                  :type="isShowPassword ? 'text' : 'password'"
                  class="form-control w-auto mb-1"
                  :placeholder="t('message.type_cof_pass')"
                />

                <div>
                  <button class="btn" @click="changePassword">
                    {{ t("settings.set_new_password") }}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <MyOffcanvas ref="mos" />
</template>
