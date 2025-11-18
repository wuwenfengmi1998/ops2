<script setup>
import { onMounted, watch, ref } from "vue";
import settingNavigation from "@/components/settingNavigation.vue";
import { useI18n } from "vue-i18n";
import datePicker from "@/components/datePicker.vue";
const { t } = useI18n();

const birthday = ref();
const username = ref();
const userremark = ref();

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
    console.log("用户信息有误，无法保存");
    return;
  }
  console.log("保存用户信息");
  console.log("用户名:", usernameValue);
  console.log("备注:", userremarkValue);
  console.log("生日:", birthdayValue);
}
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
                    <span
                      class="avatar avatar-xl"
                      style="background-image: url(./static/avatars/000m.jpg)"
                    ></span>
                  </div>
                  <div class="col-auto">
                    <button class="btn">
                      {{ t("settings.change_avatar") }}
                    </button>
                  </div>
                </div>
                <h3 class="card-title mt-4">-</h3>
                <div class="row g-3">
                  <div class="col-md">
                    <div class="form-label">{{ t("settings.name") }}</div>
                    <input ref="username" type="text" class="form-control" />
                  </div>
                  <div class="col-md">
                    <div class="form-label">{{ t("settings.remark") }}</div>
                    <input ref="userremark" type="text" class="form-control" />
                  </div>
                  <div class="col-md">
                    <div class="form-label">{{ t("settings.birthday") }}</div>
                    <datePicker ref="birthday" />
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
                      <input type="text" class="form-control w-auto" />
                    </div>
                    <div class="col-auto">
                      <button class="btn">
                        {{ t("settings.change_email") }}
                      </button>
                    </div>
                  </div>
                </div>
                <h3 class="card-title mt-4">{{ t("settings.password") }}</h3>
                <input type="password" class="form-control w-auto mb-1" />
                <input type="password" class="form-control w-auto mb-1" />
                <input type="password" class="form-control w-auto mb-1" />

                <div>
                  <button class="btn">
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
</template>
