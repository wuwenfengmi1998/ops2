<template>
	<span class="avatar" :class="'avatar-'+size" :style="'background-image: url('+url+')'"></span>
</template>

<script>
	export default {
		name:"avatar",
		props: {
			url: {
				type: String, // 指定类型
			},
			size: {
				type: String, // 指定类型
				default:"sm"
			}
		},
		data() {
			return {
				
			};
		}
	}
</script>

<style>

</style>