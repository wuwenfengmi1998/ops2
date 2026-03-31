<script setup>
import { onMounted, watch, ref } from "vue";
import settingNavigation from "@/components/settingNavigation.vue";
import { useI18n } from "vue-i18n";
import { useUserStore } from "@/stores/user";
import { my_network_func } from "@/my_network_func";
import MyOffcanvas from "@/components/MyOffcanvas.vue";
import { myfuncs } from "@/myfunc";
import { useRouter } from "vue-router";
const mos = ref();

const { t, locale } = useI18n();

const router = useRouter();

const emailInput = ref();

const userStore = useUserStore();


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


function functionupdataTitle() {
  document.title = "Operations." + t("settings.contact_information");
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
        <h2 class="page-title">{{ t("settings.my_account") }}</h2>
      </div>
    </div>

    <div class="page-body">
      <div class="container-xl">
        <div class="card">
          <div class="row g-0">
            <settingNavigation />
            <div class="col-12 col-md-9 d-flex flex-column">
              <div class="card-body">

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
                
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <MyOffcanvas ref="mos" />
</template>
