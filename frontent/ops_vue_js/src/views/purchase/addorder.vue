<script setup>
import { onMounted, watch, ref, reactive } from "vue";
import { useI18n } from "vue-i18n";

import tagadder from "@/components/tagadder.vue";
import dateTimePicker from "@/components/dateTimePicker.vue";

import useDropzone from "@/components/useDropzone.vue";

import { useUserStore } from "@/stores/user";
const userStore = useUserStore();

import { useRouter } from "vue-router";
const router = useRouter();

import TomSelect from "tom-select";
import "tom-select/dist/css/tom-select.css";

const textarea_maxlen = 256;
const textarea_len = ref(0);
const textarea_val = ref();

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
  4: t("order_status.completed"),
  5: t("order_status.refund_requested"),
  6: t("order_status.returning"),
  7: t("order_status.refunded"),
  8: t("order_status.lost_package"),
});

function update_order_status() {
  order_status["1"] = t("order_status.pending_order");
  order_status["2"] = t("order_status.order_placed");
  order_status["3"] = t("order_status.in_transit");
  order_status["4"] = t("order_status.completed");
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
  currency_type: "1",
});

function del_cost(key) {
  cost_sheet_tab.splice(key, 1);
}
function add_cost() {
  cost_sheet_tab.push(JSON.parse(JSON.stringify(cost_sheet)));
  // var t = {
  //   type: cost_type[cost_sheet.type],
  //   int: cost_sheet.int,
  //   cost: cost_sheet.cost,
  //   currency_type: currency_type[cost_sheet.currency_type],
  // };
  // cost_sheet_tab.push(t);
  //console.log(t);
}

function submit_order() {
  console.log("up");
}

function textarea_change(a) {
  //console.log(textarea_val.value.length)

  textarea_len.value = textarea_val.value.length;

  // if(a.inputType=="insertText"){
  //   textarea_len.value+=1;
  // }
  // if(a.inputType=="deleteContentBackward"){
  //   textarea_len.value-=1;
  // }
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
  }
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
                />
              </div>
              <div class="mb-3">
                <label class="form-label"
                  >{{ t("purchase_addorder.remarks") }}
                  <span class="form-label-description"
                    >{{ textarea_len }}/{{ textarea_maxlen }}</span
                  ></label
                >
                <textarea
                  class="form-control mt-2 mb-2"
                  name="example-textarea-input"
                  rows="6"
                  :placeholder="t('purchase_addorder.remarks_text')"
                  :maxlength="textarea_maxlen"
                  @input="textarea_change"
                  v-model="textarea_val"
                ></textarea>
                <useDropzone
                  acceptedFiles="image/*"
                  uploadURL="/api/files/upload/image"
                  maxFiles="10"
                ></useDropzone>
              </div>
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
                <input
                  name="url"
                  type="url"
                  class="form-control"
                  placeholder="http"
                  value=""
                />
                <div class="mb-3 mt-3">
                  <label class="form-label">{{
                    t("purchase_addorder.part_name")
                  }}</label>
                  <input
                    type="text"
                    class="form-control"
                    name="example-text-input"
                    :placeholder="t('purchase_addorder.part_name')"
                  />
                </div>
                <div class="mt-3">
                  <label class="form-label">{{
                    t("purchase_addorder.style_remarks")
                  }}</label>

                  <tagadder
                    :placeholder="t('purchase_addorder.add_style')"
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
                          {{ value.cost * value.int }}
                        </td>
                        <td class="text-secondary">
                          {{ currency_type[value.currency_type] }}
                        </td>
                        <td>
                          <button
                            class="btn btn-outline-danger"
                            @click="del_cost(key)"
                          >
                            {{ t("purchase_addorder.remove") }}
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
                        autocomplete="off"
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
                        min="0"
                        value="0"
                        v-model="cost_sheet.cost"
                      />
                    </div>
                    <div class="col-xl-2">
                      {{ t("purchase_addorder.select_currency") }}
                      <select
                        ref="select_beast"
                        class="form-control"
                        autocomplete="off"
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
                    <dateTimePicker></dateTimePicker>
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
                      value=""
                    />
                  </div>
                  <div class="col-xl-4">
                    {{ t("purchase_addorder.order_status") }}
                    <select
                      ref="select_beast"
                      class="form-control"
                      autocomplete="off"
                      value="1"
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
</template>

<style></style>
