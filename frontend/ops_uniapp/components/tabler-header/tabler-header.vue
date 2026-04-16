<template>

	<header class="navbar navbar-expand-md d-print-none">
		<div class="container-xl">
			<h1 class="navbar-brand navbar-brand-autodark d-none-navbar-horizontal pe-0 pe-md-3">
				<a href="/">
					<img src="/static/logo.svg" width="110" height="32" alt="Tabler" class="navbar-brand-image">
				</a>
			</h1>
			<div class="navbar-nav flex-row order-md-last">
				<!-- 深色模式按钮 -->
				<div class="d-none d-md-flex">
					<a @click="switchTheme('dark')" class="nav-link px-0 hide-theme-dark" title="Enable dark mode"
						data-bs-toggle="tooltip" data-bs-placement="bottom">
						<!-- Download SVG icon from http://tabler-icons.io/i/moon -->
						<svg xmlns="http://www.w3.org/2000/svg" class="icon" width="24" height="24" viewBox="0 0 24 24"
							stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round"
							stroke-linejoin="round">
							<path stroke="none" d="M0 0h24v24H0z" fill="none" />
							<path d="M12 3c.132 0 .263 0 .393 0a7.5 7.5 0 0 0 7.92 12.446a9 9 0 1 1 -8.313 -12.454z" />
						</svg>
					</a>
					<a @click="switchTheme('light')" class="nav-link px-0 hide-theme-light" title="Enable light mode"
						data-bs-toggle="tooltip" data-bs-placement="bottom">
						<!-- Download SVG icon from http://tabler-icons.io/i/sun -->
						<svg xmlns="http://www.w3.org/2000/svg" class="icon" width="24" height="24" viewBox="0 0 24 24"
							stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round"
							stroke-linejoin="round">
							<path stroke="none" d="M0 0h24v24H0z" fill="none" />
							<path d="M12 12m-4 0a4 4 0 1 0 8 0a4 4 0 1 0 -8 0" />
							<path
								d="M3 12h1m8 -9v1m8 8h1m-9 8v1m-6.4 -15.4l.7 .7m12.1 -.7l-.7 .7m0 11.4l.7 .7m-12.1 -.7l-.7 .7" />
						</svg>
					</a>
				</div>
				<!-- 深色模式按钮 -->

				<div v-if="is_login" class="nav-item dropdown">
					<a class="nav-link d-flex lh-1 text-reset p-0" data-bs-toggle="dropdown"
						aria-label="Open user menu">
						
						<avatar size="sm" :url="user_info.AvatarPath"></avatar>
						<div class="d-none d-xl-block ps-2">
							<div>{{user_info.Username}}</div>
							<div class="mt-1 small text-secondary">{{user_info.FirstName}}</div>
						</div>
					</a>
					<div class="dropdown-menu dropdown-menu-end dropdown-menu-arrow ">
						<!-- <a href="#" class="dropdown-item">Status</a> -->
						<a href="./profile" class="dropdown-item">个人资料</a>
						<!-- <a href="#" class="dropdown-item">Feedback</a> -->
						<a href="/#/pages/setting/my_info" class="dropdown-item">设置</a>
						<div class="dropdown-divider"></div>
						<a href="#" @click="logout(true)" class="dropdown-item">登出</a>
					</div>
				</div>





				<div v-if="type!='mini' && !is_login" class="nav-item btn-list">
					<a href="/#/pages/signin" class="btn" rel="noreferrer">
						<!-- Download SVG icon from http://tabler-icons.io/i/heart -->
						<svg xmlns="http://www.w3.org/2000/svg" class="icon text-blue" width="24" height="24"
							viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
							stroke-linecap="round" stroke-linejoin="round">
							<path stroke="none" d="M0 0h24v24H0z" fill="none" />
							<path d="M12 17v4" />
							<path d="M10 20l4 -2" />
							<path d="M10 18l4 2" />
							<path d="M5 17v4" />
							<path d="M3 20l4 -2" />
							<path d="M3 18l4 2" />
							<path d="M19 17v4" />
							<path d="M17 20l4 -2" />
							<path d="M17 18l4 2" />
							<path d="M9 6a3 3 0 1 0 6 0a3 3 0 0 0 -6 0" />
							<path d="M7 14a2 2 0 0 1 2 -2h6a2 2 0 0 1 2 2" />
						</svg>
						登录/注册
					</a>
				</div>




			</div>
		</div>
	</header>
	<header v-if="type!='mini'" class="navbar-expand-md">
		<div class="collapse navbar-collapse" id="navbar-menu">
			<div class="navbar">
				<div class="container-xl">
					<ul class="navbar-nav">
						<li class="nav-item">
							<a class="nav-link" href="/">
								<span
									class="nav-link-icon d-md-none d-lg-inline-block"><!-- Download SVG icon from http://tabler-icons.io/i/home -->
									<svg xmlns="http://www.w3.org/2000/svg" class="icon" width="24" height="24"
										viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none"
										stroke-linecap="round" stroke-linejoin="round">
										<path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
										<path d="M5 12l-2 0l9 -9l9 9l-2 0"></path>
										<path d="M5 12v7a2 2 0 0 0 2 2h10a2 2 0 0 0 2 -2v-7"></path>
										<path d="M9 21v-6a2 2 0 0 1 2 -2h2a2 2 0 0 1 2 2v6"></path>
									</svg>
								</span>
								<span class="nav-link-title">
									主页
								</span>
							</a>
						</li>
						<li class="nav-item">
							<a class="nav-link" href="/warehouses">
								<span
									class="nav-link-icon d-md-none d-lg-inline-block"><!-- Download SVG icon from http://tabler-icons.io/i/home -->
									<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
										fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
										stroke-linejoin="round"
										class="icon icon-tabler icons-tabler-outline icon-tabler-building-warehouse">
										<path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
										<path d="M3 21v-13l9 -4l9 4v13"></path>
										<path d="M13 13h4v8h-10v-6h6"></path>
										<path d="M13 21v-9a1 1 0 0 0 -1 -1h-2a1 1 0 0 0 -1 1v3"></path>
									</svg>
								</span>
								<span class="nav-link-title">
									仓库
								</span>
							</a>
						</li>
						<li class="nav-item">
							<a class="nav-link" href="/workorders">
								<span
									class="nav-link-icon d-md-none d-lg-inline-block"><!-- Download SVG icon from http://tabler-icons.io/i/home -->
									<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
										fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
										stroke-linejoin="round"
										class="icon icon-tabler icons-tabler-outline icon-tabler-file-description">
										<path stroke="none" d="M0 0h24v24H0z" fill="none" />
										<path d="M14 3v4a1 1 0 0 0 1 1h4" />
										<path
											d="M17 21h-10a2 2 0 0 1 -2 -2v-14a2 2 0 0 1 2 -2h7l5 5v11a2 2 0 0 1 -2 2z" />
										<path d="M9 17h6" />
										<path d="M9 13h6" />
									</svg>
								</span>
								<span class="nav-link-title">
									工单
								</span>
							</a>
						</li>

					</ul>

				</div>
			</div>
		</div>
	</header>
</template>

<script>
	import {
		my_network_func
	} from '../../my_network_func';
	import {
		myfunc
	} from '../../myfunc';



	myfunc.setTheme(myfunc.getThemefromStorge(), false);



	export default {
		name: "tabler-header",
		props: {
			type: {
				type: String, // 指定类型
				default: 'normal' // 可选的默认值
			}
		},
		data() {
			return {
				is_login: false,
				user_info:{
					AvatarPath:"",
					Username:"",
					FirstName:"",
					Birthdate:"",
				}
			};
		},
		methods: {
			switchTheme(theme) {
				myfunc.setTheme(theme, true);
			},
			logout(logout_server) {
				if (logout_server) {
					my_network_func.post_json("/user/logout", null, (c) => {
						if (c.statusCode == 200) {
							console.log(c)

						}
					})
				}

				this.is_login = false
				myfunc.dele("user_info")
				myfunc.dele("cookie")

			},
			updata_user_info_to_heard() {
				//h5先判断有没有cookie
				var cookie = myfunc.load_json("cookie")
				if (cookie) {
					//判断cookie有没有过期
					const isFuture = new Date(cookie.ExpiresAt) > new Date();
					if (isFuture) {
						//没过期
						this.is_login = true
						//获取用户信息
						this.user_info = myfunc.load_json("user_info")

					} else {
						//过期了
						this.logout(false)
						//如何是h5平台 在请求的时候本地cookie应该就被服务器删了
					}

				} else {
					this.logout(false)
				}
			},
		},
		mounted() {
			this.updata_user_info_to_heard()
		}
	}
</script>

<style>

</style>