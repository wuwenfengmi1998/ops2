<script>
	import { userStore } from './store/user.js'

	export default {
		onLaunch: function() {
			// 应用启动时恢复登录状态（对标 PC 端 App.vue 的 restoreSession）
			userStore.restoreSession()
		},
		onShow: function() {},
		onHide: function() {}
	}
</script>




<style>
	/*每个页面公共css */
	@import url("/static/dist/css/tabler.min.css");
	@import url("/mycss.css");

</style>
