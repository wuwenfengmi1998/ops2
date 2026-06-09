<template>
  <view class="page">
    <view class="group">
      <text class="group-title">打印机控制</text>
      <view class="group-body">
        <view class="cell" @tap="initPrinter">
          <text class="cell-label">初始化打印</text>
        </view>
        <view class="cell" @tap="closePrinter">
          <text class="cell-label">关闭打印</text>
        </view>
      </view>
    </view>

    <view class="group">
      <text class="group-title">打印测试</text>
      <view class="group-body">
        <view class="cell" @tap="printerText">
          <text class="cell-label">打印小票</text>
        </view>
        <view class="cell" @tap="printerLabel">
          <text class="cell-label">测试打印（文字+二维码）</text>
        </view>
        <view class="cell" @tap="printLabel">
          <text class="cell-label">打印标签（条形码）</text>
        </view>
		<view class="cell" @tap="myprinttest">
		  <text class="cell-label">我的demo</text>
		</view>
      </view>
    </view>

    <!-- 日志区域 -->
    <view class="log-box" v-if="logs.length">
      <text class="group-title">日志</text>
      <view class="log-content">
        <text class="log-line" v-for="(log, i) in logs" :key="i">{{ log }}</text>
      </view>
      <view class="btn-wrapper">
        <button class="clear-btn" @tap="logs = []">清空日志</button>
      </view>
    </view>
  </view>
</template>

<script setup>
  import { ref, onMounted } from 'vue'
  import { onLoad } from '@dcloudio/uni-app'

  const printer = uni.requireNativePlugin('LcPrinter')
  const modal = uni.requireNativePlugin('modal')
  const globalEvent = uni.requireNativePlugin('globalEvent')

  const logs = ref([])

  function log(msg) {
    logs.value.push(msg)
  }

  onLoad(() => {
    // 添加打印状态监听
    globalEvent.addEventListener('onPrintCallback', (e) => {
      log('state: ' + JSON.stringify(e))
      uni.showToast({
        title: 'state: ' + JSON.stringify(e),
        duration: 2000
      })
      if (e.key == 0) {
        uni.showToast({
          title: '打印成功',
          duration: 2000
        })
      } else if (e.key == 3) {
        uni.showToast({
          title: '缺纸',
          duration: 2000
        })
      }
    })

    // 打印机版本获取回调
    globalEvent.addEventListener('onVersion', (e) => {
      log('version: ' + JSON.stringify(e))
      uni.showToast({
        title: 'version: ' + JSON.stringify(e),
        duration: 2000
      })
    })

    globalEvent.addEventListener('getsupportprint', (e) => {
      log('getsupportprint: ' + JSON.stringify(e))
      uni.showToast({
        title: 'key: ' + JSON.stringify(e),
        duration: 2000
      })
    })
  })

  function initPrinter() {
    console.log('初始化')
    log('初始化打印机...')
    const ret = printer.initPrinter({})
    modal.toast({
      message: ret,
      duration: 1.5
    })
    printer.setConcentration({
      level: 39
    })
    printer.setLineSpacing({
      spacing: 1
    })
    printer.getsupportprint()
    log('初始化完成')
  }

  function closePrinter() {
    console.log('关闭')
    log('关闭打印机')
    printer.closePrinter()
  }

  function printerLabel() {
    console.log('测试打印')
    log('测试打印（文字+二维码）...')
    
    // 标签打印，使用黑标
    printer.printEnableMark({
      enable: true
    })

    printer.setFontSize({
      fontSize: 0
    })
    printer.setTextBold({
      bold: true
    })

    printer.printText({
      content: '垃圾收运小票凭证'
    })
    printer.printLine({
      line_length: 1
    })
    printer.printText({
      content: 'asdads'
    })
    printer.printLine({
      line_length: 1
    })

    console.log('测试打印QR')
    log('打印二维码...')
    printer.printQR({
      text: 'title',
      height: 400,
      offset: 1
    })
    printer.printLine({
      line_length: 2
    })

    printer.start()
    console.log('测试打印QR结束')
    log('打印完成')
  }

  function printerText() {
    console.log('打印小票')
    log('打印小票...')
    
    // 普通打印（小票），不使用黑标
    printer.printEnableMark({
      enable: false
    })

    printer.setFontSize({
      fontSize: 0
    })
    printer.setTextBold({
      bold: true
    })

    printer.printText({
      content: '这是一张测试小票'
    })
    printer.setTextBold({
      bold: false
    })
    printer.printLine({
      line_length: 1
    })

    printer.printText({
      content: 'commodity'
    })
    printer.printLine({
      line_length: 1
    })

    printer.printText({
      content: 'Quantity'
    })
    printer.printLine({
      line_length: 1
    })
    printer.printText({
      content: 'unit price'
    })
    printer.printLine({
      line_length: 1
    })

    printer.printBarcode({
      text: '123456',
      height: 80,
      barcodeType: 73
    })
    printer.printLine({
      line_length: 1
    })
    printer.printQR({
      text: '1234456'
    })
    printer.printLine({
      line_length: 2
    })
    printer.start()
    log('小票打印完成')
  }

  function printLabel() {
    log('打印标签（条形码）...')
    
    // 标签打印，使用黑标
    printer.printEnableMark({
      enable: true
    })
    
    printer.printBarcode({
      text: '1234567890123456789',
      height: 80,
      barcodeType: 73
    })
    printer.printLine({
      line_length: 5
    })
    printer.printGoToNextMark()
    log('标签打印完成')
  }
  
  function myprinttest(){
	  // 标签打印，使用黑标
	  printer.printEnableMark({
	    enable: true
	  })
	  printer.printText({
	    content: 'type:1\n123'
	  })
	  printer.printBarcode({
	       text: 'abcdefg:1234',
	       height: 40,
	       barcodeType: 73,
		   
	     })
	  printer.printGoToNextMark()
  }
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f2f2f7;
  padding-top: 20rpx;
}
.group {
  margin-bottom: 32rpx;
}
.group-title {
  font-size: 24rpx;
  color: #8e8e93;
  padding: 0 32rpx 12rpx;
  display: block;
}
.group-body {
  background-color: #ffffff;
  border-radius: 12rpx;
  margin: 0 24rpx;
  overflow: hidden;
}
.cell {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  padding: 0 32rpx;
  height: 96rpx;
  background-color: #ffffff;
  border-bottom: 1rpx solid #e5e5ea;
}
.cell:last-child {
  border-bottom: none;
}
.cell:active {
  background-color: #f2f2f7;
}
.cell-label {
  font-size: 30rpx;
  color: #1c1c1e;
}
.btn-wrapper {
  padding: 20rpx 32rpx;
}
.clear-btn {
  background-color: #ff3b30;
  color: #fff;
  font-size: 28rpx;
  border-radius: 12rpx;
  padding: 12rpx 0;
}
.log-box {
  margin: 0 24rpx 32rpx;
  background: #1c1c1e;
  border-radius: 12rpx;
  padding: 20rpx;
}
.log-content {
  max-height: 400rpx;
  overflow-y: auto;
}
.log-line {
  display: block;
  font-size: 22rpx;
  color: #32d74b;
  line-height: 1.6;
  word-break: break-all;
}
</style>
