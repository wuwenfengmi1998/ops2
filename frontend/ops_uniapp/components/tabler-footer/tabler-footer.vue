<template>
	<footer class="footer footer-transparent d-print-none my_footer">
		<div class="container-xl">
			<div class="row text-center align-items-center flex-row-reverse">
				<div class="col-lg-auto ms-lg-auto">
					<ul class="list-inline list-inline-dots mb-0">
						<li class="list-inline-item"><a
								href="https://git.lmve.net/kevin/gin_saas/-/blob/main/readme.md?ref_type=heads"
								target="_blank" class="link-secondary" rel="noopener">文档</a></li>
						<li class="list-inline-item"><a
								href="https://git.lmve.net/kevin/gin_saas/-/blob/main/LICENSE?ref_type=heads"
								target="_blank" class="link-secondary">开源协议</a>
						</li>
						<li class="list-inline-item"><a href="https://git.lmve.net/kevin/gin_saas" target="_blank"
								class="link-secondary" rel="noopener">
								<svg xmlns="http://www.w3.org/2000/svg" class="icon text-orange" width="24" height="24"
									viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
									stroke-linecap="round" stroke-linejoin="round">
									<path stroke="none" d="M0 0h24v24H0z" fill="none" />
									<path d="M21 14l-9 7l-9 -7l3 -11l3 7h6l3 -7z" />
								</svg>
								源码</a></li>
						<li class="list-inline-item">
							<a href="https://wnfed.com" target="_blank" class="link-secondary" rel="noopener">
								<!-- Download SVG icon from http://tabler-icons.io/i/heart -->
								<svg xmlns="http://www.w3.org/2000/svg" class="icon text-pink icon-filled icon-inline"
									width="24" height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"
									fill="none" stroke-linecap="round" stroke-linejoin="round">
									<path stroke="none" d="M0 0h24v24H0z" fill="none" />
									<path
										d="M19.5 12.572l-7.5 7.428l-7.5 -7.428a5 5 0 1 1 7.5 -6.566a5 5 0 1 1 7.5 6.572" />
								</svg>
								博客
							</a>
						</li>
					</ul>
				</div>
				<div class="col-12 col-lg-auto mt-3 mt-lg-0">
					<ul class="list-inline list-inline-dots mb-0">
						<li class="list-inline-item">
							Copyright &copy; 2025
							<a href="https://lmve.net" target="_blank" class="link-secondary">Lmve.NET</a>.
							All rights reserved.
						</li>
						<li class="list-inline-item">
							<a href="https://git.lmve.net/kevin/gin_saas/-/commits/main" target="_blank"
								class="link-secondary" rel="noopener">
								v0.0.1
							</a>
						</li>
					</ul>
				</div>
			</div>
		</div>
	</footer>

	<div class="offcanvas my_offcanvas_top alert alert-important alert-dismissible" :class="{
	    'alert-success': alertType === 'success',
	    'alert-warning': alertType === 'warning',
	    'alert-danger': alertType === 'danger',
	    'alert-info': alertType === 'info'
	  }" role="alert" tabindex="-1" ref="offcanvas">
		<div class="d-flex">
			<div>
				<!-- Download SVG icon from http://tabler-icons.io/i/check -->
				<svg v-if="alertType === 'success'" xmlns="http://www.w3.org/2000/svg" class="icon alert-icon" width="24" height="24"
					viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round"
					stroke-linejoin="round">
					<path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
					<path d="M5 12l5 5l10 -10"></path>
				</svg>
				<svg v-if="alertType === 'warning'" xmlns="http://www.w3.org/2000/svg" class="icon alert-icon" width="24" height="24"
					viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round"
					stroke-linejoin="round">
					<path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
					<path
						d="M10.24 3.957l-8.422 14.06a1.989 1.989 0 0 0 1.7 2.983h16.845a1.989 1.989 0 0 0 1.7 -2.983l-8.423 -14.06a1.989 1.989 0 0 0 -3.4 0z">
					</path>
					<path d="M12 9v4"></path>
					<path d="M12 17h.01"></path>
				</svg>
				<svg v-if="alertType === 'danger'" xmlns="http://www.w3.org/2000/svg" class="icon alert-icon" width="24" height="24"
					viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round"
					stroke-linejoin="round">
					<path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
					<path d="M3 12a9 9 0 1 0 18 0a9 9 0 0 0 -18 0"></path>
					<path d="M12 8v4"></path>
					<path d="M12 16h.01"></path>
				</svg>
				<svg v-if="alertType === 'info'" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
					stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
					class="icon icon-tabler icons-tabler-outline icon-tabler-brand-hipchat">
					<path stroke="none" d="M0 0h24v24H0z" fill="none" />
					<path
						d="M17.802 17.292s.077 -.055 .2 -.149c1.843 -1.425 3 -3.49 3 -5.789c0 -4.286 -4.03 -7.764 -9 -7.764c-4.97 0 -9 3.478 -9 7.764c0 4.288 4.03 7.646 9 7.646c.424 0 1.12 -.028 2.088 -.084c1.262 .82 3.104 1.493 4.716 1.493c.499 0 .734 -.41 .414 -.828c-.486 -.596 -1.156 -1.551 -1.416 -2.29z" />
					<path d="M7.5 13.5c2.5 2.5 6.5 2.5 9 0" />
				</svg>



			</div>
			<div>
				{{alertText}}
			</div>
		</div>
		<a class="btn-close" data-bs-dismiss="offcanvas" aria-label="close"></a>

	</div>
</template>

<script>
	import {
		Offcanvas
	} from 'bootstrap'

	export default {
		name: "tabler-footer",
		data() {
			return {
				alertType: 'success',
				alertText: "juest do!",
			};
		},
		methods: {
			alert(type,text) {
				this.alertType=type
				this.alertText=text
				if (!this.offcanvasInstance) {
					const el = this.$refs.offcanvas
					this.offcanvasInstance = new Offcanvas(el, {
						backdrop: false,
					})
				}
				this.offcanvasInstance.show()
				setTimeout(() => {

					this.offcanvasInstance?.hide()
				}, 5000);

			},
			close() {
				this.offcanvasInstance?.hide()
			},
		},
		mounted() {
			
		}
	}
</script>

<style>

</style>