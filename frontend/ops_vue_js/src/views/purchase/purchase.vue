<script setup>
import { onMounted, watch, ref, reactive } from "vue";
import { useI18n } from "vue-i18n";
import MyOffcanvas from "@/components/MyOffcanvas.vue";
const mos = ref();
import { my_network_func } from "@/my_network_func";
import { myfuncs } from "@/myfunc";

import { useRouter } from "vue-router";
const router = useRouter();

const { t, locale } = useI18n();

const all_items = ref(0);
const all_pages = ref(0);
const page_items = ref(10);
const now_page = ref(1);

const all_orders = ref({});

const page_start = ref(0);
const page_end = ref(0);

const page_input = ref();
const page_items_items = ref("10");

function jump_to_order(order_id) {
  //console.log(order_id);

  var order_str=order_id.toString()

  const resolved = router.resolve({
    path: "/purchase/showorder/" + order_str,
  });
  window.open(resolved.href, "_blank");
}

//获取订单列表
function get_orders() {
  my_network_func.postJson(
    "/purchase/getorders",
    {
      search: "",
      entries: page_items.value,
      page: now_page.value,
    },
    (r) => {
      //console.log(r);
      switch (r.statusCode) {
        case 200:
          switch (r.data.err_code) {
            case 0:
              all_orders.value = r.data.return.all_orders;
              all_items.value = r.data.return.all_count;
              all_pages.value = Math.ceil(all_items.value / page_items.value);
              if (now_page.value < 3) {
                page_start.value = 1;
              } else {
                if (now_page.value > all_pages.value - 3) {
                  page_start.value = all_pages.value - 4;
                  if (page_start.value <= 0) {
                    page_start.value = 1;
                  }
                } else {
                  page_start.value = now_page.value - 2;
                }
              }
              if (now_page.value > all_pages.value - 3) {
                page_end.value = all_pages.value;
              } else {
                if (now_page.value < 3) {
                  page_end.value = 5;
                } else {
                  page_end.value = now_page.value + 2;
                }
              }
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
    },
  );
}

function change_page(page) {
  now_page.value = page;
  get_orders();
}

function functionupdataTitle() {
  document.title = "Operations." + t("appname.purchase");
}

function range(start, end) {
  return Array.from({ length: end - start + 1 }, (_, i) => start + i);
}

function page_input_change(c) {
  //console.log(page_input.value);

  var t = parseInt(page_input.value);
  if (t > 0) {
    if (t <= all_pages.value) {
      page_input.value = "";
      change_page(t);
    }
  }
}

function page_input_input(c) {
  page_input.value = page_input.value.replace(/[^\d]/g, "");

  //console.log(c)
}

function page_items_input_change(c) {
  var t = parseInt(page_items_items.value);
  page_items.value = t;
  now_page.value = 1;
  get_orders();
  //console.log(t)
}

function page_items_input_input(c) {
  page_items_items.value = page_items_items.value.replace(/[^\d]/g, "");

  var t = parseInt(page_items_items.value);
  if (t > 300) {
    page_items_items.value = "300";
  }

  //console.log(c)
}

onMounted(() => {
  functionupdataTitle();

  get_orders();
});
// 监听语言变化，更新标题
watch(locale, () => {
  functionupdataTitle();
});
</script>

<template>
  <div class="page-body">
    <div class="container-xl">
      <div class="card">
        <div class="card-header">
          <h3 class="card-title">{{ t("purchase.purchase_list") }}</h3>
        </div>
        <div class="card-body border-bottom py-3">
          <div class="d-flex">
            <div class="text-secondary">
              <router-link to="/purchase/addorder" class="btn btn-info m-1">
                {{ t("purchase.add_part") }}
              </router-link>
              <button class="btn m-1">
                {{ t("purchase.exp_report") }}
              </button>
            </div>

            <!-- //搜索dom -->
            <div class="ms-auto text-secondary">
              {{ t("purchase.search") }}
              <div class="ms-2 d-inline-block mr-2">
                <input
                  type="text"
                  class="form-control form-control-sm"
                  aria-label="Search invoice"
                />
              </div>
            </div>

            <div class="ms-auto text-secondary"></div>
          </div>
        </div>
        <div class="table-responsive">
          <table class="table card-table table-vcenter text-nowrap datatable">
            <thead>
              <tr>
                <th class="w-1">
                  <input
                    class="form-check-input m-0 align-middle"
                    type="checkbox"
                    aria-label="Select all invoices"
                  />
                </th>
                <th class="col-1">
                  No.
                  <!-- Download SVG icon from http://tabler-icons.io/i/chevron-up -->
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="icon icon-sm icon-thick"
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
                    <path d="M6 15l6 -6l6 6" />
                  </svg>
                </th>
                <th class="col-3">{{ t("purchase.item_name") }}</th>
                <th class="col-3">{{ t("purchase.purpose") }}</th>
                <th class="w-1">{{ t("purchase.quantity") }}</th>
                <th class="w-1">{{ t("purchase.created_at") }}</th>
                <th class="w-1">{{ t("purchase.updated_at") }}</th>
                <th class="w-1">{{ t("purchase.status") }}</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="value in all_orders"
                class="element"
                @click="jump_to_order(value.ID)"
              >
                <td>
                  <input
                    class="form-check-input m-0 align-middle"
                    type="checkbox"
                    aria-label="Select invoice"
                  />
                </td>
                <td>
                  <span class="text-muted">{{ value.ID }}</span>
                </td>
                <td>
                  {{ value.Title }}
                </td>
                <td>{{ value.Remark }}</td>
                <td>1</td>
                <td>
                  {{ myfuncs.formatLocalizedDate(value.CreatedAt, locale) }}
                </td>
                <td>{{ myfuncs.formatLocalizedDate(value.UpdatedAt) }}</td>
                <td>1</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="card-footer d-flex align-items-center">
          <p class="m-0 text-secondary">
            {{ t("purchase.show") }}
          </p>
          <div class="mx-2 d-inline-block">
            <input
              type="text"
              class="form-control form-control-sm w-6"
              v-model="page_items_items"
              aria-label="Invoices count"
              @change="page_items_input_change"
              @input="page_items_input_input"
            />
          </div>
          <p class="m-0 text-secondary">
            {{ t("purchase.entries") }}
            {{ t("purchase.There_are_a_total_of") }} {{ all_items }}
            {{ t("purchase.entries") }}
          </p>

          <ul class="pagination m-0 ms-auto">
            <li class="page-item" :class="now_page == 1 ? 'disabled' : ''">
              <div
                class="page-link"
                :tabindex="now_page == 1 ? '-1' : ''"
                :aria-disabled="now_page == 1 ? 'true' : ''"
                @click="change_page(1)"
              >
                <!-- Download SVG icon from http://tabler-icons.io/i/chevron-left -->
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  class="icon icon-tabler icons-tabler-outline icon-tabler-arrow-bar-to-left"
                >
                  <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                  <path d="M10 12l10 0" />
                  <path d="M10 12l4 4" />
                  <path d="M10 12l4 -4" />
                  <path d="M4 4l0 16" />
                </svg>
              </div>
            </li>
            <li class="page-item" :class="now_page == 1 ? 'disabled' : ''">
              <div
                class="page-link"
                :tabindex="now_page == 1 ? '-1' : ''"
                :aria-disabled="now_page == 1 ? 'true' : ''"
                @click="change_page(now_page - 1)"
              >
                <!-- Download SVG icon from http://tabler-icons.io/i/chevron-left -->
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
                  <path d="M15 6l-6 6l6 6" />
                </svg>
                <!-- prev -->
              </div>
            </li>

            <li
              v-for="value in range(page_start, page_end)"
              class="page-item"
              :class="value == now_page ? 'active' : ''"
            >
              <div class="page-link" @click="change_page(value)">
                {{ value }}
              </div>
            </li>

            <li
              class="page-item"
              :class="now_page == all_pages ? 'disabled' : ''"
            >
              <div
                class="page-link"
                :tabindex="now_page == all_pages ? '-1' : ''"
                :aria-disabled="now_page == all_pages ? 'true' : ''"
                @click="change_page(now_page + 1)"
              >
                <!-- next -->
                <!-- Download SVG icon from http://tabler-icons.io/i/chevron-right -->
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
                  <path d="M9 6l6 6l-6 6" />
                </svg>
              </div>
            </li>

            <li
              class="page-item"
              :class="now_page == all_pages ? 'disabled' : ''"
            >
              <div
                class="page-link"
                :tabindex="now_page == all_pages ? '-1' : ''"
                :aria-disabled="now_page == all_pages ? 'true' : ''"
                @click="change_page(all_pages)"
              >
                <!-- Download SVG icon from http://tabler-icons.io/i/chevron-right -->
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  class="icon icon-tabler icons-tabler-outline icon-tabler-arrow-bar-to-right"
                >
                  <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                  <path d="M14 12l-10 0" />
                  <path d="M14 12l-4 4" />
                  <path d="M14 12l-4 -4" />
                  <path d="M20 4l0 16" />
                </svg>
              </div>
            </li>
            <li>
              <input
                type="text"
                class="form-control form-control-sm w-6"
                @change="page_input_change"
                @input="page_input_input"
                v-model="page_input"
              />
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>

  <MyOffcanvas ref="mos" />
</template>

<style lang="scss" scoped>
.element:hover {
  background-color: #4299e11c;
}
</style>
