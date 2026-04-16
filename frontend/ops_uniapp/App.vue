<script>
	export default {
		onLaunch: function() {
			console.log('App Launch')
		},
		onShow: function() {
			console.log('App Show')
		},
		onHide: function() {
			console.log('App Hide')
		}
	}
	import "/static/dist/js/tabler.min.js"
	import "/static/dist/libs/bootstrap/dist/js/bootstrap.min.js"
</script>




<style>
	/*每个页面公共css */
	@import url("/static/dist/css/tabler.min.css");
	@import url("/mycss.css");

</style>
