<script setup>
import { onMounted, watch, ref, reactive } from "vue";
import { useI18n } from "vue-i18n";

import MyOffcanvas from "@/components/MyOffcanvas.vue";

import tagadder from "@/components/tagadder.vue";
import dateTimePicker from "@/components/dateTimePicker.vue";

import useDropzone from "@/components/useDropzone.vue";

import { useUserStore } from "@/stores/user";
const userStore = useUserStore();

import { useRouter } from "vue-router";
const router = useRouter();

import "tom-select/dist/css/tom-select.css";
import { my_network_func } from "@/my_network_func";

const textarea_maxlen = 256;

const title_input_dom = ref();

const photos_hash = ref();

const mos = ref();

const { t, locale } = useI18n();

//货币类型
const currency_type = reactive({
  1: "RMB",
  2: "MOP",
  3: "HKD",
  4: "USD",
});
//成本类型
const cost_type = reactive({
  1: t("cost_type.unit_price"),
  2: t("cost_type.freight"),
});
function update_cost_type() {
  cost_type["1"] = t("cost_type.unit_price");
  cost_type["2"] = t("cost_type.freight");
}

//订单状态
const order_status = reactive({
  1: t("order_status.pending_order"),
  2: t("order_status.order_placed"),
  3: t("order_status.in_transit"),
  4: t("order_status.compvared"),
  5: t("order_status.refund_requested"),
  6: t("order_status.returning"),
  7: t("order_status.refunded"),
  8: t("order_status.lost_package"),
});

function update_order_status() {
  order_status["1"] = t("order_status.pending_order");
  order_status["2"] = t("order_status.order_placed");
  order_status["3"] = t("order_status.in_transit");
  order_status["4"] = t("order_status.compvared");
  order_status["5"] = t("order_status.refund_requested");
  order_status["6"] = t("order_status.returning");
  order_status["7"] = t("order_status.refunded");
  order_status["8"] = t("order_status.lost_package");
}

const cost_sheet_tab = reactive([]);
// 表单对象
const cost_sheet = reactive({
  type: "1",
  int: 1,
  cost: 0.0,
  cost_t: 0.0,
  currency_type: "1",
});

function del_cost(key) {
  cost_sheet.type = cost_sheet_tab[key].type;
  cost_sheet.int = cost_sheet_tab[key].int;
  cost_sheet.cost = cost_sheet_tab[key].cost;
  cost_sheet.cost_t = cost_sheet_tab[key].cost_t;
  cost_sheet.currency_type = cost_sheet_tab[key].currency_type;

  cost_sheet_tab.splice(key, 1);
}
function add_cost() {
  if (cost_sheet.cost <= 0) {
  } else {
    // 四舍五入到2位小数
    var t = parseFloat((cost_sheet.int * cost_sheet.cost).toFixed(2));
    cost_sheet.cost_t = t;

    cost_sheet_tab.push(JSON.parse(JSON.stringify(cost_sheet)));

    cost_sheet.type = "1";
    cost_sheet.int = 1;
    cost_sheet.cost = 0.0;
    cost_sheet.cost_t = 0.0;
    cost_sheet.currency_type = "1";
  }
}

const submit_sheet = reactive({
  title: "",
  remark: "",
  photos: [],
  link: "",
  partname: "",
  styles: "",
  costs: [],
  updatetime: "",
  trackingnumber: "",
  orderstatus: "1",
});

function submit_order() {
  if (submit_sheet.title == "") {
    title_input_dom.value.classList.add("is-invalid");
    title_input_dom.value.addEventListener("input", function () {
      if (this.value.trim() !== "") {
        this.classList.remove("is-invalid");
        //this.removeEventListener('input');
      }
    });

    mos.value?.showAlert("danger", t("purchase_addorder.title"), 1000);
    return;
  }
  //载入图片哈希列表
  submit_sheet.photos = [];
  var photos = photos_hash.value.return_files();
  for (var i = 0; i < photos.length; i++) {
    submit_sheet.photos.push(photos[i].hash);
  }

  //载入价格表
  submit_sheet.costs = [];
  for (var i = 0; i < cost_sheet_tab.length; i++) {
    //var t=cost_sheet_tab[i]
    submit_sheet.costs.push(JSON.parse(JSON.stringify(cost_sheet_tab[i])));
  }

  //修改价格表里的小数，将所有价值*100去掉小数
  for (var i = 0; i < submit_sheet.costs.length; i++) {
    submit_sheet.costs[i].cost *= 100;
    submit_sheet.costs[i].cost_t *= 100;
  }

  console.log(submit_sheet);
  my_network_func.postJson("/purchase/addorder", submit_sheet, (r) => {
    console.log(r);
  });
}

function functionupdataTitle() {
  document.title = "Operations." + t("purchase.add_part");
}

onMounted(() => {
  functionupdataTitle();
  //sele_init();
  if (!userStore.isLoggedIn) {
    router.push("/login");
  }
});
// 监听语言变化，更新标题
watch(locale, () => {
  functionupdataTitle();
  update_cost_type();
  update_order_status();
});

// 监听 cost 变化，自动限制小数位
watch(
  () => cost_sheet.cost,
  (newVal) => {
    if (newVal !== null && newVal !== undefined) {
      // 四舍五入到2位小数
      const fixed = parseFloat(newVal).toFixed(2);
      if (parseFloat(fixed) !== newVal) {
        cost_sheet.cost = parseFloat(fixed);
      }
    }
  },
);
</script>

<template>
  <div class="page-header d-print-none">
    <div class="container-xl">
      <div class="row g-2 align-items-center">
        <div class="col">
          <h2 class="page-title">{{ t("purchase_addorder.add_order") }}</h2>
        </div>
      </div>
    </div>
  </div>
  <div class="page-body">
    <div class="container-xl">
      <div class="row row-cards">
        <div class="col-12">
          <div class="card">
            <div class="card-header">
              <h4 class="card-title">
                {{ t("purchase_addorder.order_info") }}
              </h4>
            </div>
            <div class="card-body">
              <div class="mb-3">
                <label class="form-label required">{{
                  t("purchase_addorder.title")
                }}</label>
                <input
                  type="text"
                  class="form-control"
                  name="example-text-input"
                  :placeholder="t('purchase_addorder.title')"
                  v-model="submit_sheet.title"
                  ref="title_input_dom"
                />
              </div>
              <div class="mb-3">
                <label class="form-label"
                  >{{ t("purchase_addorder.remarks") }}
                  <span class="form-label-description"
                    >{{ submit_sheet.remark.length }}/{{
                      textarea_maxlen
                    }}</span
                  ></label
                >
                <textarea
                  class="form-control mt-2 mb-2"
                  name="example-textarea-input"
                  rows="6"
                  :placeholder="t('purchase_addorder.remarks_text')"
                  :maxlength="textarea_maxlen"
                  v-model="submit_sheet.remark"
                ></textarea>
              </div>

              <label class="form-label mb-0">{{
                t("purchase_addorder.photo_remarks")
              }}</label>
              <useDropzone
                acceptedFiles="image/*"
                uploadURL="/api/files/upload/image"
                maxFiles="10"
                ref="photos_hash"
              ></useDropzone>
            </div>

            <div class="card-header">
              <h4 class="card-title">
                {{ t("purchase_addorder.purchase_channel") }}
              </h4>
            </div>
            <div class="card-body">
              <div class="mb-3">
                <label class="form-label">{{
                  t("purchase_addorder.link")
                }}</label>
                <textarea
                  name="url"
                  type="url"
                  class="form-control"
                  placeholder="https"
                  v-model="submit_sheet.link"
                ></textarea>
                <div class="mb-3 mt-3">
                  <label class="form-label">{{
                    t("purchase_addorder.part_name")
                  }}</label>
                  <input
                    type="text"
                    class="form-control"
                    name="example-text-input"
                    :placeholder="t('purchase_addorder.part_name')"
                    v-model="submit_sheet.partname"
                  />
                </div>
                <div class="mt-3">
                  <label class="form-label">{{
                    t("purchase_addorder.style_remarks")
                  }}</label>

                  <tagadder
                    :placeholder="t('purchase_addorder.add_style')"
                    v-model="submit_sheet.styles"
                  ></tagadder>
                </div>

                <div class="mt-3">
                  <label class="form-label">{{
                    t("purchase_addorder.cost")
                  }}</label>

                  <table
                    v-show="cost_sheet_tab.length"
                    class="table table-vcenter card-table table-striped"
                  >
                    <thead>
                      <tr>
                        <th>{{ t("purchase_addorder.type") }}</th>
                        <th>{{ t("purchase_addorder.quantity") }}</th>
                        <th>{{ t("purchase_addorder.fee") }}</th>
                        <th>{{ t("purchase_addorder.total_price") }}</th>
                        <th>{{ t("purchase_addorder.currency") }}</th>
                        <th class="w-1">
                          {{ t("purchase_addorder.operation") }}
                        </th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="(value, key) in cost_sheet_tab">
                        <td>{{ cost_type[value.type] }}</td>
                        <td class="text-secondary">{{ value.int }}</td>
                        <td class="text-secondary">{{ value.cost }}</td>
                        <td class="text-secondary">
                          {{ value.cost_t }}
                        </td>
                        <td class="text-secondary">
                          {{ currency_type[value.currency_type] }}
                        </td>
                        <td>
                          <button
                            class="btn btn-outline-danger"
                            @click="del_cost(key)"
                          >
                            {{ t("purchase_addorder.change") }}
                          </button>
                        </td>
                      </tr>
                      <!-- <tr>
                        <td>运输</td>
                        <td class="text-secondary">1</td>
                        <td class="text-secondary">5</td>
                        <td class="text-secondary">MOP</td>
                        <td>
                          <button class="btn btn-outline-danger">Del</button>
                        </td>
                      </tr> -->
                    </tbody>
                  </table>

                  <div class="row g-5">
                    <div class="col-xl-2">
                      {{ t("purchase_addorder.fee_type") }}
                      <select
                        ref="select_type"
                        class="form-control"
                        autocompvare="off"
                        value="1"
                        v-model="cost_sheet.type"
                      >
                        <option v-for="(value, key) in cost_type" :value="key">
                          {{ value }}
                        </option>
                      </select>
                    </div>
                    <div class="col-xl-3">
                      {{ t("purchase_addorder.input_quantity") }}
                      <input
                        type="number"
                        class="form-control"
                        min="1"
                        value="1"
                        v-model="cost_sheet.int"
                      />
                    </div>
                    <div class="col-xl-3">
                      {{ t("purchase_addorder.input_fee") }}
                      <input
                        type="number"
                        class="form-control"
                        step="0.01"
                        min="0.0"
                        value="0.0"
                        v-model="cost_sheet.cost"
                      />
                    </div>
                    <div class="col-xl-2">
                      {{ t("purchase_addorder.select_currency") }}
                      <select
                        ref="select_beast"
                        class="form-control"
                        autocompvare="off"
                        value="1"
                        v-model="cost_sheet.currency_type"
                      >
                        <option
                          v-for="(value, key) in currency_type"
                          :value="key"
                        >
                          {{ value }}
                        </option>
                      </select>
                    </div>
                    <div class="col-xl-2">
                      {{ t("purchase_addorder.operation") }}
                      <button
                        class="form-control btn btn-outline-primary"
                        @click="add_cost"
                      >
                        {{ t("purchase_addorder.add") }}
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div class="card-header">
              <h4 class="card-title">
                {{ t("purchase_addorder.other_status") }}
              </h4>
            </div>
            <div class="card-body">
              <div class="mb-3">
                <div class="row g-5">
                  <div class="col-xl-4">
                    <label class="form-label required">{{
                      t("purchase_addorder.update_time")
                    }}</label>
                    <dateTimePicker
                      v-model="submit_sheet.updatetime"
                    ></dateTimePicker>
                  </div>
                  <div class="col-xl-4">
                    <label class="form-label">{{
                      t("purchase_addorder.tracking_number")
                    }}</label>
                    <input
                      type="text"
                      class="form-control"
                      :placeholder="
                        t('purchase_addorder.input_tracking_number')
                      "
                      v-model="submit_sheet.trackingnumber"
                    />
                  </div>
                  <div class="col-xl-4">
                    {{ t("purchase_addorder.order_status") }}
                    <select
                      ref="select_beast"
                      class="form-control"
                      autocompvare="off"
                      v-model="submit_sheet.orderstatus"
                    >
                      <option v-for="(value, key) in order_status" :value="key">
                        {{ value }}
                      </option>
                    </select>
                  </div>
                </div>
              </div>
            </div>
            <div class="card-footer text-end">
              <div class="d-flex">
                <button
                  type="submit"
                  class="btn btn-primary ms-auto"
                  @click="submit_order"
                >
                  {{ t("purchase_addorder.submit") }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <MyOffcanvas ref="mos" />
</template>

<style></style>
