<script setup>
import { RouterLink } from "vue-router";
import { useI18n } from "vue-i18n";
const { t } = useI18n();

</script>

<template>
  <div class="col-12 col-md-3 border-end">
    <div class="card-body">
      <div class="list-group list-group-transparent">
        <router-link
          to="/settings/account"
          class="list-group-item list-group-item-action d-flex align-items-center active"
        >
          {{t('settings.account_settings')}}</router-link
        >
      
        
      </div>
      <h4 class="subheader mt-4">{{t('settings.admin')}}</h4>
      <div class="list-group list-group-transparent">
        <router-link
          to="/settings/account"
          class="list-group-item list-group-item-action d-flex align-items-center active"
        >
          {{t('settings.website_settings')}}</router-link
        >
      </div>
    </div>
  </div>
</template>
