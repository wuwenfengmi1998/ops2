<script setup>
/**
 * 采购订单添加页面
 *
 * 功能概述：
 * - 创建新的采购订单
 * - 支持上传订单相关图片
 * - 支持添加多种费用类型（单价、运费）
 * - 支持选择订单状态、填写快递单号等
 *
 * 表单字段：
 * - title: 订单标题（必填）
 * - remark: 备注说明
 * - photos: 订单图片列表
 * - link: 采购链接
 * - partname: 配件名称
 * - styles: 款式标签
 * - costs: 费用明细列表
 * - tracking_number: 快递单号
 * - updatetime: 更新时间
 * - order_status: 订单状态
 */

// ==================== 依赖导入 ====================
import { reactive, ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { useValidation } from '@/composables'
import { purchaseApi } from '@/api/purchase'

// 组件导入
import tagadder from '@/components/tagadder.vue'      // 标签添加组件
import datePicker from '@/components/datePicker.vue'  // 日期选择组件
import useDropzone from '@/components/useDropzone.vue' // 文件上传组件（图片）

// ==================== 页面初始化 ====================
// 设置页面标题，使用 i18n key
usePageTitle('purchase.add_part')

// 获取国际化实例，用于多语言文本
const { t, locale } = useI18n()

// Toast 消息提示
const toast = useToastStore()

// 表单验证工具
const { validate, errors, clearErrors } = useValidation()

// ==================== 常量定义 ====================
/**
 * 备注文字最大长度限制
 * 超过此长度将显示字符计数警告
 */
const textMaxLen = 256

/**
 * 文件上传组件的引用
 * 用于获取已上传的图片列表
 */
const photosRef = ref(null)

/**
 * 货币类型选项
 * key: 数据库中存储的值
 * value: 显示的货币符号
 */
const currencyOptions = { 1: 'CNY', 2: 'MOP', 3: 'HKD', 4: 'USD' }

/**
 * 货币类型对应的国旗图标路径
 */
const currencyFlags = { 1: '/static/flags/cn.svg', 2: '/static/flags/mo.svg', 3: '/static/flags/hk.svg', 4: '/static/flags/us.svg' }

// ==================== 费用类型映射 ====================
// 费用类型：单价 / 运费
const costType = computed(() => ({
  1: t('cost_type.unit_price'),  // 单价
  2: t('cost_type.freight'),     // 运费
}))

// ==================== 订单状态映射 ====================
// 订单的各种状态选项
const orderStatus = computed(() => ({
  1: t('order_status.pending_order'),      // 待下单
  2: t('order_status.order_placed'),       // 已下单
  3: t('order_status.in_transit'),         // 运输中
  4: t('order_status.completed'),           // 已完成
  5: t('order_status.refund_requested'),   // 退款中
  6: t('order_status.returning'),           // 退货中
  7: t('order_status.refunded'),            // 已退款
  8: t('order_status.lost_package'),        // 丢件
}))

// ==================== 费用明细管理 ====================
/**
 * 已添加的费用明细列表
 * 每项包含：类型、数量、单价、总价、货币类型
 */
const costEntries = reactive([])

/**
 * 新费用条目的临时数据
 * 用户填写完表单后点击"添加"按钮加入 costEntries
 */
const newCost = reactive({
  type: '1',       // 费用类型：默认"单价"
  int: 1,          // 数量：默认1
  cost: 0,         // 单价：默认0
  currencyType: '1', // 货币类型：默认人民币
})

// ==================== 表单数据 ====================
/**
 * 订单表单主数据对象
 * 包含所有需要提交的字段
 */
const form = reactive({
  title: '',          // 订单标题（必填）
  remark: '',         // 备注说明
  photos: [],         // 图片列表（上传后由 dropzone 填充）
  link: '',           // 采购链接
  partname: '',       // 配件名称
  styles: '',         // 款式标签
  costs: [],          // 费用明细（提交前由 costEntries 转换）
  tracking_number: '', // 快递单号
  updatetime: '',     // 更新时间
  order_status: '1',  // 订单状态（默认待下单）
})

/**
 * 计算新费用条目的总价
 * 总价 = 数量 × 单价，保留2位小数
 */
const newCostTotal = computed(() =>
  parseFloat((newCost.int * newCost.cost).toFixed(2))
)

/**
 * 添加一条费用明细到列表
 * 条件：单价必须大于0
 */
function addCostEntry() {
  if (newCost.cost <= 0) return
  costEntries.push({
    type: newCost.type,
    int: newCost.int,
    cost: newCost.cost,
    costt: newCostTotal.value,
    currencytype: newCost.currencyType,
  })
  // 添加后重置表单，以便继续添加下一条
  newCost.type = '1'
  newCost.int = 1
  newCost.cost = 0
  newCost.currencyType = '1'
}

/**
 * 从费用明细列表中移除指定条目
 * @param {number} index - 要删除的条目索引
 */
function removeCostEntry(index) {
  costEntries.splice(index, 1)
}

/**
 * 监听单价输入，自动保留2位小数
 * 防止用户输入如 10.999 这样的值
 */
watch(() => newCost.cost, (val) => {
  const fixed = parseFloat(val).toFixed(2)
  if (parseFloat(fixed) !== val) newCost.cost = parseFloat(fixed)
})



/**
 * 提交按钮 loading 状态
 * 防止重复提交
 */
const loading = ref(false)

// ==================== 表单提交 ====================
/**
 * 处理表单提交
 * 1. 验证必填字段
 * 2. 收集图片列表
 * 3. 转换费用数据（金额从元转为分）
 * 4. 调用 API 提交数据
 */
async function handleSubmit() {
  // 清空之前的验证错误
  clearErrors()

  // 验证标题是否填写
  const err = validate('title', form.title, t('purchase_addorder.title'))
  if (!err) return

  // 从 dropzone 组件获取已上传的图片文件名
  form.photos = []
  if (photosRef.value?.has_some_files) {
    const result = photosRef.value.get_some_files()
    form.photos = result.map(f => f.name)
  }

  // 将费用明细转换为提交格式
  // 注意：金额需要从"元"转为"分"（乘以100）存储
  form.costs = costEntries.map(h => ({
    ...h,
    cost: Math.round(h.cost * 100),
    costt: Math.round(h.costt * 100),
  }))

  // 开始 loading
  loading.value = true
  try {
    // 调用采购 API 添加订单
    const { errCode } = await purchaseApi.addOrder(form)
    if (errCode === 0) {
      // 保存成功，显示成功提示
      toast.success(t('message.save_ok'))
    } else {
      // 服务器错误，显示错误提示
      toast.error(t('message.server_error'))
    }
  } catch {
    // 错误已被 HTTP 拦截器处理，此处无需额外处理
  } finally {
    // 无论成功失败，都要关闭 loading
    loading.value = false
  }
}
</script>

<template>
  <!-- 页面容器：居中，最大宽度 6xl，左右内边距，上下内边距 -->
  <div class="mx-auto max-w-6xl px-6 py-6">
    <!-- 主卡片：白色背景，圆角，阴影，支持暗色模式 -->
    <div class="flex flex-col gap-6 rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card">

      <!-- ==================== 订单信息区块 ==================== -->
      <div class="border-b border-gray-200 px-6 py-4 dark:border-dk-muted">
        <h4 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('purchase_addorder.order_info') }}</h4>
      </div>

      <!-- 订单信息表单区域 -->
      <div class="space-y-4 px-6 py-5">
        <!-- 标题字段（必填） -->
        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('purchase_addorder.part_name') }} <span class="text-red-500">*</span>
          </label>
          <input
            v-model="form.title"
            type="text"
            class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
            :class="errors.title ? 'border-red-500' : 'border-gray-300'"
            :placeholder="t('purchase_addorder.part_name')"
          />
          <!-- 验证错误提示 -->
          <span v-if="errors.title" class="mt-1 block text-xs text-red-500">{{ errors.title }}</span>
        </div>

        <!-- 备注字段 -->
        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('purchase_addorder.remarks') }}
            <!-- 字符计数 -->
            <span class="text-gray-400">{{ form.remark.length }}/{{ textMaxLen }}</span>
          </label>
          <textarea
            v-model="form.remark"
            class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
            rows="4"
            :placeholder="t('purchase_addorder.remarks_text')"
          />
        </div>

        <!-- 图片上传区域 -->
        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('purchase_addorder.photo_remarks') }}</label>
          <!--
            useDropzone 组件属性：
            - acceptFiles="image/*": 只接受图片文件
            - uploadURL: 上传接口地址
            - maxFiles="10": 最多10张图片
            - ref="photosRef": 组件引用，用于获取上传的文件列表
          -->
          <useDropzone acceptFiles="image/*" uploadURL="/api/files/upload/image" maxFiles="10" ref="photosRef" />
        </div>
      </div>

      <!-- ==================== 采购渠道区块 ==================== -->
      <div class="border-t border-gray-200 px-6 py-4 dark:border-dk-muted">
        <h4 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('purchase_addorder.purchase_channel') }}</h4>
      </div>

      <div class="space-y-4 px-6 py-5">
        <!-- 采购链接 -->
        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('purchase_addorder.link') }}</label>
          <textarea
            v-model="form.link"
            class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
            rows="2"
            placeholder="url"
          />
        </div>

        <!-- 款式标签 -->
        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('purchase_addorder.style_remarks') }}</label>
          <!-- tagadder: 支持添加多个标签的组件 -->
          <tagadder :placeholder="t('purchase_addorder.add_style')" v-model="form.styles" />
        </div>

        <!-- ==================== 费用明细表格 ==================== -->
        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('purchase_addorder.cost') }}</label>

          <!-- 已添加的费用列表表格（有时才显示） -->
          <div v-if="costEntries.length" class="mb-4 overflow-x-auto">
            <table class="w-full text-left text-sm text-gray-900">
              <thead>
                <tr class="border-b border-gray-200 bg-gray-50 text-gray-500 dark:border-dk-muted dark:bg-dk-base">
                  <th class="px-3 py-2 font-medium">{{ t('purchase_addorder.type') }}</th>
                  <th class="px-3 py-2 font-medium">{{ t('purchase_addorder.quantity') }}</th>
                  <th class="px-3 py-2 font-medium">{{ t('purchase_addorder.fee') }}</th>
                  <th class="px-3 py-2 font-medium">{{ t('purchase_addorder.total_price') }}</th>
                  <th class="px-3 py-2 font-medium">{{ t('purchase_addorder.currency') }}</th>
                  <th class="px-3 py-2 font-medium">{{ t('purchase_addorder.operation') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(item, idx) in costEntries" :key="idx" class="border-b border-gray-100 dark:border-dk-muted">
                  <td class="px-3 py-2 font-medium text-gray-900 dark:text-white">{{ costType[item.type] }}</td>
                  <td class="px-3 py-2 text-gray-500">{{ item.int }}</td>
                  <td class="px-3 py-2 text-gray-500">{{ item.cost }}</td>
                  <td class="px-3 py-2 text-gray-500">{{ item.costt }}</td>
                  <td class="px-3 py-2 text-gray-500">
                    <img :src="currencyFlags[item.currencytype]" class="inline-block h-4 w-6 align-middle" :alt="currencyOptions[item.currencytype]" />
                    {{ currencyOptions[item.currencytype] }}
                  </td>
                  <td class="px-3 py-2">
                    <!-- 删除按钮 -->
                    <button class="rounded px-2 py-1 text-xs font-medium text-red-600 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/20" @click="removeCostEntry(idx)">{{ t('purchase_addorder.remove') }}</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- 添加费用表单：响应式网格布局 -->
          <div class="grid grid-cols-2 gap-4 sm:grid-cols-5">
            <!-- 费用类型选择 -->
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-500">{{ t('purchase_addorder.fee_type') }}</label>
              <select v-model="newCost.type" class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white">
                <template v-for="(label, key) in costType" :key="key">
                  <option :value="key">{{ label }}</option>
                </template>
              </select>
            </div>

            <!-- 数量输入 -->
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-500">{{ t('purchase_addorder.input_quantity') }}</label>
              <input v-model.number="newCost.int" type="number" class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white" min="1" />
            </div>

            <!-- 单价输入 -->
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-500">{{ t('purchase_addorder.input_fee') }}</label>
              <input v-model="newCost.cost" type="number" class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white" step="0.01" min="0" />
            </div>

            <!-- 货币类型选择 -->
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-500">{{ t('purchase_addorder.select_currency') }}</label>
              <select v-model="newCost.currencyType" class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white">
                <template v-for="(label, key) in currencyOptions" :key="key">
                  <option :value="key">{{ label }}</option>
                </template>
              </select>
            </div>

            <!-- 添加按钮（与输入框底部对齐） -->
            <div class="flex items-end">
              <button class="w-full rounded-lg border border-gray-300 bg-blue-600 px-3 py-2 text-sm font-semibold text-blue-100 transition-colors hover:bg-blue-700 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:text-white dark:bg-blue-600" @click="addCostEntry">{{ t('purchase_addorder.add') }}</button>
            </div>
          </div>
        </div>
      </div>

      <!-- ==================== 订单状态区块 ==================== -->
      <div class="border-t border-gray-200 px-6 py-4 dark:border-dk-muted">
        <h4 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('purchase_addorder.order_status') }}</h4>
      </div>

      <div class="px-6 py-5">
        <!-- 三列响应式网格 -->
        <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <!-- 更新时间（必填） -->
          <div>
            <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
              {{ t('purchase_addorder.update_time') }} <span class="text-red-500">*</span>
            </label>
            <datePicker v-model="form.updatetime" />
          </div>

          <!-- 快递单号 -->
          <div>
            <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('purchase_addorder.tracking_number') }}</label>
            <input
              v-model="form.tracking_number"
              type="text"
              class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
              :placeholder="t('purchase_addorder.input_tracking_number')"
            />
          </div>

          <!-- 订单状态下拉框 -->
          <div>
            <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('purchase_addorder.order_status') }}</label>
            <select v-model="form.order_status" class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white">
              <template v-for="(label, key) in orderStatus" :key="key">
                <option :value="key">{{ label }}</option>
              </template>
            </select>
          </div>
        </div>
      </div>

      <!-- ==================== 底部操作栏 ==================== -->
      <div class="flex justify-end border-t border-gray-200 px-6 py-4 dark:border-dk-muted">
        <!-- 提交按钮：提交时显示 loading 动画 -->
        <button
          class="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-blue-700 focus:ring-2 focus:ring-blue-500/20 focus:outline-none disabled:active:scale-100"
          :disabled="loading"
          @click="handleSubmit"
        >
          <!-- loading 状态显示旋转图标 -->
          <svg v-if="loading" class="h-4 w-4 animate-spin text-white" viewBox="0 0 24 24" fill="none">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
          {{ t('purchase_addorder.submit') }}
        </button>
      </div>

    </div>
  </div>
</template>
