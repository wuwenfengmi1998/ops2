<template>

	<view class="d-flex flex-column">
		<tabler-header type="mini"></tabler-header>

		<div class="page page-center">

			<div class="container container-normal py-4">
				<div class="row align-items-center g-4">
					<div class="col-lg">
						<div class="container-tight">
							<div class="text-center mb-4">
								<div class="navbar-brand navbar-brand-autodark"><img src="/static/logo.svg" height="36"
										alt=""></div>
							</div>
							<div class="card card-md">
								<div class="card-body">
									<h2 class="h2 text-center mb-4">登录你的账号</h2>


									<div class="mb-3">
										<label class="form-label mb-0 ms-1">用户名</label>

										<input type="text" class="my_input_field"
											:placeholder-class="is_username_err?'my_input_placeholder_error':'my_input_placeholder'"
											:placeholder="ph_username" autocomplete="off" maxlength="25"
											v-model="post_data.username">
									</div>
									<div class="mb-2">
										<label class="form-label mb-0 ms-1">密码</label>

										<input type="password" class="my_input_field"
											:placeholder-class="is_password_err?'my_input_placeholder_error':'my_input_placeholder'"
											password="true" :placeholder="ph_password" autocomplete="off"
											maxlength="100" v-model="post_data.password">

									</div>
<!-- 									<div class="mb-2">
										<checkbox-group @change="chack_box_change">
											<label class="d-flex">
												<checkbox value="keep_login_in" />
												<span class="form-check-label">保持登录</span>
											</label>
										</checkbox-group>
									</div> -->
									<div class="form-footer">
										<button class="btn btn-primary w-100" @click="submit_data">登录</button>
									</div>

								</div>


							</div>
							<div class="text-center text-secondary mt-3">
								还没账号？
								<a href="/sign-up/" tabindex="-1">点击注册</a>
								<span class="form-label-description">
									<a href="./forgot-password.html">忘记密码？</a>
								</span>
							</div>

						</div>
					</div>
					<div class="col-lg d-none d-lg-block">
						<img src="/static/illustrations/undraw_secure_login_pdn4.svg" height="300"
							class="d-block mx-auto" alt="">
					</div>
				</div>
			</div>
		</div>

		<tabler-footer ref="footer"></tabler-footer>
	</view>



</template>

<script>
	import {
		my_network_func
	} from '../my_network_func'
	import {
		myfunc
	} from '../myfunc'


	export default {
		data() {
			return {
				ph_username: "",
				ph_password: "",
				is_username_err: false,
				is_password_err: false,
				post_data: {
					is_keep_login: true,
					username: "",
					password: "",
				}
			}
		},
		methods: {
			initdata(){
				this.ph_username="输入你的用户名"
				this.ph_password="输入你的密码"
				this.is_username_err=false
				this.is_password_err=false
				this.post_data.is_keep_login=true
				this.post_data.username=""
				this.post_data.password=""
			},

			chack_box_change(val) {
				//console.log(val.detail.value[0])
				if (val.detail.value[0] === 'keep_login_in') {
					this.post_data.is_keep_login = true
				} else {
					this.post_data.is_keep_login = false
				}
			},
			submit_data() {
				//提交登录数据，
				//先验证数据合法性


				if (this.post_data.username == "") {
					this.is_username_err = true
					this.ph_username = "用户名不能为空"

				} else {
					this.is_username_err = false
				}

				if (this.post_data.password == "") {
					this.is_password_err = true
					this.ph_password = "密码不能为空"

				} else {
					this.is_password_err = false
				}

				if (this.is_username_err === false && this.is_password_err === false) {
					
					my_network_func.post_json("/user/login", this.post_data, (c) => {
						
						if (c.statusCode == 200) {
							if (c.data.err_code == 0) {
								this.$refs.footer.alert('success', "登录成功")
								myfunc.save_json("cookie", c.data.return.cookie)
								myfunc.save_json("user_info", c.data.return.user_info)
								this.initdata()
								setTimeout(() => {
									uni.navigateTo({
										url: '/'
									});
								}, 1000);

							} else {
								this.$refs.footer.alert('warning', "账号或密码不正确")
							}
						} else {
							this.$refs.footer.alert('danger', "网络连接错误：" + c.statusCode)
						}


					})
				} else {

				}

			}
		},
		mounted() {
			this.initdata()
		}
	}
</script>

<style>
	.my123 {
		/* 设置边框：宽度、样式、颜色 */
		border: 1px solid #e3e3e3;
		/* 2像素宽的红色实线边框 */

		/* 设置圆角，值越大，角越圆 */
		border-radius: 3px;
		/* 也可以使用百分比，例如50%会变成圆形 */

		font-size: 24px;

	}
</style>