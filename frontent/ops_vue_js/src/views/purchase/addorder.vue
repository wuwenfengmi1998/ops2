<script setup>
import { onMounted, watch, ref ,reactive} from "vue";
import { useI18n } from "vue-i18n";

import tagadder from "@/components/tagadder.vue";
import datePicker from "@/components/datePicker.vue";

import TomSelect from "tom-select";
import "tom-select/dist/css/tom-select.css";

const { t, locale } = useI18n();

function functionupdataTitle() {
  document.title = "Operations." + t("purchase.add_part");
}

const cost_sheet_tab = reactive([]);
// 表单对象
const cost_sheet = reactive({
  type: '1',
  int: 1,
  cost: 0.0,
  currency_type:'1'
})

function add_cost(){
  var t={
     type: cost_sheet.type,
      int: cost_sheet.int,
      cost: cost_sheet.cost,
      currency_type:cost_sheet.currency_type
  }
cost_sheet_tab.push(t)
console.log(cost_sheet_tab)

}

// const select_beast = ref()
// const select_type = ref()
// function sele_init(){
// new TomSelect(select_beast.value,{
// 	create: true,
// 	sortField: {
// 		//field: "text",
// 		//direction: "asc"
// 	}
// });

// new TomSelect(select_type.value,{
// 	create: true,
// 	sortField: {
// 		//field: "text",
// 		//direction: "asc"
// 	}
// });

// }

onMounted(() => {
  functionupdataTitle();
  //sele_init();
});
// 监听语言变化，更新标题
watch(locale, () => {
  functionupdataTitle();
});
</script>

<template>
  <div class="page-header d-print-none">
    <div class="container-xl">
      <div class="row g-2 align-items-center">
        <div class="col">
          <h2 class="page-title">添加订单</h2>
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
              <h4 class="card-title">订单信息</h4>
            </div>
            <div class="card-body">
              <div class="mb-3">
                <label class="form-label required">标题</label>
                <input
                  type="text"
                  class="form-control"
                  name="example-text-input"
                  placeholder="输入订单标题"
                />
              </div>
              <div class="mb-3">
                <label class="form-label"
                  >用途 <span class="form-label-description">0/100</span></label
                >
                <textarea
                  class="form-control"
                  name="example-textarea-input"
                  rows="6"
                  placeholder="Content.."
                ></textarea>
              </div>
            </div>

            <div class="card-header">
              <h4 class="card-title">采购途径</h4>
            </div>
            <div class="card-body">
              <div class="mb-3">
                <label class="form-label">URL</label>
                <input
                  name="url"
                  type="url"
                  class="form-control"
                  placeholder="http"
                  value=""
                />
                <div class="mt-3">
                  <label class="form-label">样式备注</label>

                  <tagadder placeholder="添加样式"></tagadder>
                </div>

                <div class="mt-3">
                  <label class="form-label">成本</label>

                  <table class="table table-vcenter card-table table-striped">
                    <thead>
                      <tr>
                        <th>类型</th>
                        <th>数量</th>
                        <th>费用</th>
                        <th>总价</th>
                        <th>货币</th>
                        <th class="w-1">操作</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="value in cost_sheet_tab">
                        <td>{{value.type}}</td>
                        <td class="text-secondary">{{value.int}}</td>
                        <td class="text-secondary">{{value.cost}}</td>
                        <td class="text-secondary">{{value.cost*value.int}}</td>
                        <td class="text-secondary">{{value.currency_type}}</td>
                        <td>
                          <button class="btn btn-outline-danger">Del</button>
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
                      类型
                      <select
                        ref="select_type"
                        class="form-control"
                        placeholder="选择费用类型"
                        autocomplete="off"
                        value="1"
                        v-model="cost_sheet.type"
                      >
                        <option value="1">单价</option>
                        <option value="2">运输</option>
                      </select>
                    </div>
                    <div class="col-xl-3">
                      数量
                      <input
                        type="number"
                        class="form-control"
                        min="1"
                        value="1"
                        v-model="cost_sheet.int"
                      />
                    </div>
                    <div class="col-xl-3">
                      费用
                      <input
                        type="number"
                        class="form-control"
                        step="0.01"
                        min="0"
                        value="0.0"
                        v-model="cost_sheet.cost"
                      />
                    </div>
                    <div class="col-xl-2">
                      货币
                      <select
                        ref="select_beast"
                        class="form-control"
                        placeholder="选择货币类型"
                        autocomplete="off"
                        value="1"
                        v-model="cost_sheet.currency_type"
                      >
                        <option value="1">RMB</option>
                        <option value="2">MOP</option>
                        <option value="3">HKD</option>
                        <option value="4">USD</option>
                      </select>
                    </div>
                    <div class="col-xl-2">
                      操作
                      <button class="form-control btn btn-outline-primary" @click="add_cost">
                        添加
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div class="card-header">
              <h4 class="card-title">其他状态</h4>
            </div>
            <div class="card-body">
              <div class="mb-3">
                <div class="row g-5">
                  <div class="col-xl-4">
                    <label class="form-label required">更新时间</label>
                    <datePicker></datePicker>
                  </div>
                  <div class="col-xl-4">
                    <label class="form-label">快递单号</label>
                    <input
                      type="text"
                      class="form-control"
                      placeholder="输入快递单号"
                      value=""
                    />
                  </div>
                  <div class="col-xl-4">
                    订单状态
                    <select
                      ref="select_beast"
                      class="form-control"
                      placeholder="选择订单状态"
                      autocomplete="off"
                      value="1"
                    >
                      <option value="1">待下单</option>
                      <option value="2">已下单</option>
                      <option value="3">运输中</option>
                      <option value="4">已完成</option>
                      <option value="5">申请退款</option>
                      <option value="6">退回中</option>
                      <option value="7">已退款</option>
                    </select>
                  </div>
                </div>
              </div>
            </div>
            <div class="card-footer text-end">
              <div class="d-flex">
                <button class="btn btn-link">Cancel</button>
                <button type="submit" class="btn btn-primary ms-auto">
                  Send data
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
