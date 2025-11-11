<script setup>
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>
<template>
  <footer class="footer footer-transparent d-print-none">
    <div class="container-xl">
      <div class="row text-center align-items-center flex-row-reverse">
        <div class="col-lg-auto ms-lg-auto">
          <ul class="list-inline list-inline-dots mb-0">
            <li class="list-inline-item">
              <a
                href="https://github.com/wuwenfengmi1998/ops2/blob/main/readme.md"
                target="_blank"
                class="link-secondary"
                rel="noopener"
                >{{t('footer.doc')}}</a
              >
            </li>
            <li class="list-inline-item">
              <a
                href="https://github.com/wuwenfengmi1998/ops2/blob/main/LICENSE"
                target="_blank"
                class="link-secondary"
                >{{t('footer.license')}}</a
              >
            </li>
            <li class="list-inline-item">
              <a
                href="https://github.com/wuwenfengmi1998/ops2"
                target="_blank"
                class="link-secondary"
                rel="noopener"
                >{{t('footer.source_code')}}</a
              >
            </li>
            <li class="list-inline-item">
              <a
                href="https://github.com/wuwenfengmi1998"
                target="_blank"
                class="link-secondary"
                rel="noopener"
              >

                kevin
              </a>
            </li>
          </ul>
        </div>
        <div class="col-12 col-lg-auto mt-3 mt-lg-0">
          <ul class="list-inline list-inline-dots mb-0">
            <li class="list-inline-item">
              {{ t('footer.copy')}}
            </li>
            <li class="list-inline-item">
              <a
                href="https://git.lmve.net/kevin/ops2/-/commits/main"
                target="_blank"
                class="link-secondary"
                rel="noopener"
              >
                v0.0.1
              </a>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </footer>
</template>
