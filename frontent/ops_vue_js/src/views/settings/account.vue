<script setup>
import { onMounted, watch, ref } from 'vue'
import settingNavigation from "@/components/settingNavigation.vue";
import { useI18n } from "vue-i18n";
import datePicker from "@/components/datePicker.vue";
const { t } = useI18n();

const datapicker = ref();

function updataInfo() {
  console.log("保存用户信息");
  console.log("生日:", datapicker.value.datepicker.value);
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
                    <input type="text" class="form-control" />
                  </div>
                  <div class="col-md">
                    <div class="form-label">{{ t("settings.remark") }}</div>
                    <input type="text" class="form-control" />
                  </div>
                  <div class="col-md">
                    <div class="form-label">{{ t("settings.birthday") }}</div>
                    <datePicker ref="datapicker"/>
                    
                  </div>
                  <div>
                    <button class="btn" @click="updataInfo">{{ t("settings.save_changes") }}</button>
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
                  <button class="btn">{{ t("settings.set_new_password") }}</button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
