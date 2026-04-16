<template>

	<tabler-header type="mini" ref="heard"></tabler-header>
	<div class="page-wrapper">
		<!-- Page header -->
		<div class="page-header d-print-none">
			<div class="container-xl">
				<div class="row g-2 align-items-center">
					<div class="col">
						<h2 class="page-title">
							设置
						</h2>
					</div>
				</div>
			</div>
		</div>
		<!-- Page body -->
		<div class="page-body">
			<div class="container-xl">
				<div class="card">
					<div class="row g-0">

						<setting-menu></setting-menu>

						<div class="col-12 col-md-9 d-flex flex-column">
							<div class="card-body">

								<h2 class="mb-4">信息设置</h2>
								<!-- <h3 class="card-title">Profile Details</h3> -->
								<div class="row align-items-center">
									<div class="col-auto">

										<avatar size="xl" :url="user_info.AvatarPath"></avatar>
									</div>
									<div class="col-auto"><a class="btn" @click="showPopup">更改头像 </a></div>
									<!-- <div class="col-auto"><a href="#" class="btn btn-ghost-danger">
	                        删除头像
	                      </a></div> -->
								</div>
								<!-- <h3 class="card-title mt-4">Business Profile</h3> -->
								<div class="row g-3">
									<div class="col-md">
										<div class="form-label">名字</div>
										<input type="text" class="my_input_field" v-model="user_info.Username"
											maxlength="30">
									</div>

									<div class="col-md">
										<div class="form-label">备注</div>
										<input type="text" class="my_input_field" v-model="user_info.FirstName"
											maxlength="50">
									</div>
								</div>

								<div class="mb-3">
									<label class="form-label">生日</label>
									<div class="row g-2">
										<div class="col-5">

											<uni-datetime-picker type="date" :clear-icon="false" v-model="birthdate" />
										</div>
									</div>
								</div>


							</div>

							<div class="card-footer bg-transparent mt-auto">
								<div class="btn-list justify-content-end">

									<a href="#" class="btn btn-primary">
										提交
									</a>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>

	</div>

	<tabler-footer></tabler-footer>

	<uni-popup ref="popup" type="center">
		<view class="popup-content">
			<div class="card card-body container">
				<h1>头像裁剪工具</h1>

				<div class="flex-wrapper">
					<!-- 左侧裁剪区 -->
					<div class="crop-section">
						<div ref="image_wrapper">
							<img ref="cropper_image"
								src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII=">
						</div>
						<!-- 上传进度 -->
						<div class="progress-container">
							<div class="progress-bar"></div>
						</div>



						<div class="preview-stats">
							<!-- <p>当前缩放: <span id="zoomValue">100%</span></p> -->
							<p>图片尺寸: <span ref="imageSize">0 x 0</span></p>
						</div>
					</div>

					<!-- 右侧预览区 -->
					<div class="preview-section">
						<h3>实时预览</h3>
						<div class="preview-box">
							<img ref="preview_img"
								src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII=">
						</div>

						<!-- 控制按钮 -->
						<div class="controls">

							<button class="btn btn-primary" @click="selefile">📁 选择图片</button>
							<button class="btn btn-secondary " onclick="rotateImage(-90)">↩️ 左旋</button>
							<button class="btn btn-success " id="uploadBtn">✂️ 裁剪头像</button>

							<button class="btn btn-danger " style="margin-top: 150px;" @click="closePopup">❌ 取消</button>



						</div>
					</div>
				</div>
			</div>
		</view>
	</uni-popup>




</template>

<script>
	import "/static/dist/libs/cropper/cropper.min.js"

	export default {
		data() {
			return {
				user_info: {
					AvatarPath: "",
					Username: "",
					FirstName: "",
					Birthdate: "",
				},
				birthdate: "",
			}
		},
		methods: {
			showPopup() {
				this.$refs.popup.open();
			},
			closePopup() {
				this.$refs.popup.close();
			},
			updateImageInfo() {
				if (this.cropper) {
					const data = this.cropper.getData();
					// document.getElementById('zoomValue').textContent = 
					//     `${Math.round(currentScale * 100)}%`;
					// document.getElementById('imageSize').textContent =
					//   `${Math.round(data.width)} x ${Math.round(data.height)}`;
				}
			},
			updatePreview() {
				const canvas = this.cropper.getCroppedCanvas({
					width: 250,
					height: 250,
					imageSmoothingQuality: 'high'
				});

				if (canvas) {
					var that = this
					canvas.toBlob(blob => {
						const previewUrl = URL.createObjectURL(blob);
						that.$refs.preview_img.src = previewUrl;
						// 清理旧URL
						setTimeout(() => URL.revokeObjectURL(previewUrl), 1000);
					}, 'image/jpeg', 0.85);
					this.updateImageInfo();
				}
			},
			initCropper(imageSrc) {
				var that = this
				// 初始化Cropper
				if (this.cropper) {
					this.cropper.destroy();
				}
				const image = this.$refs.cropper_image
				console.log(image)
				image.src = imageSrc;
				this.cropper = new Cropper(image, {
					aspectRatio: 1,
					viewMode: 2,
					autoCropArea: 0.8,
					zoomable: true,
					zoomOnWheel: true,
					zoomOnTouch: true,
					wheelZoomRatio: 0.1,
					//minCanvasWidth: 400,
					//minCanvasHeight: 400,
					crop: that.updatePreview,
					ready() {
						that.updateImageInfo();
						//document.querySelector('.progress-container').style.display = 'none';
					}
				});
			},
			selefile() {
				var that=this
				uni.chooseImage({
					count: 1, // 最多9张
					sourceType: ['album', 'camera'], // 来源相册或相机
					success: (res) => {

						console.log(res);
						// 如果需要得到原始的File对象（在H5中），则使用res.tempFiles
						const file = res.tempFiles[0];
						if (!file) return;

						if (!file.type.startsWith('image/')) {
							that.$refs.footer.alert('warning', "⚠️ 请选择有效的图片文件")
							return;
						}

						const reader = new FileReader();
						reader.onload = () => {
							that.initCropper(reader.result);
							//currentScale = 1;
						};
						reader.readAsDataURL(file);
					}
				});
				// this.initCropper(
				// 	'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII='
				// );

			},

		},

		mounted() {

			if (this.$refs.heard.is_login) {
				this.user_info = this.$refs.heard.user_info
				this.birthdate = this.user_info.Birthdate.substring(0, 10)
			} else {

			}

			// 初始化默认图片



		}
	}
</script>

<style>
	@import url("/static/dist/libs/cropper/cropper.min.css");
	/* 头像裁剪器样式*/

	.container {
		width: 95%;
		/* 改为百分比宽度 */
		margin: 20px auto;
		/* 增加上下边距 */
		max-width: 1200px;
		/* 保留最大宽度 */
		background: white;
		padding: 30px;
		border-radius: 12px;

	}

	.flex-wrapper {
		display: flex;
		gap: 30px;
		margin-top: 20px;
		flex-wrap: wrap;
		/* 添加换行支持 */
	}

	/* 裁剪区域 */
	.crop-section {
		background-color: aqua;
		flex: 1 1 60%;
		/* 弹性布局基础宽度 */
		min-width: 300px;
		/* 降低最小宽度 */
		height: auto;
		/* 移除固定高度 */
		min-height: 400px;
		/* 设置最小高度 */
	}

	#image-wrapper {
		width: 100%;
		height: 60vh;
		/* 改用视窗单位 */
		max-height: 600px;
		/* 设置最大高度 */
		background: #f8f9fa;
		border: 2px dashed #ddd;
		border-radius: 8px;
		overflow: hidden;
	}

	/* 预览区域自适应 */
	.preview-section {
		flex: 1 1 35%;
		/* 弹性布局基础宽度 */
		min-width: 250px;
		/* 设置合理最小宽度 */
	}

	/* 移动端适配 */
	@media (max-width: 768px) {
		.container {
			padding: 15px;
			/* 减少内边距 */
		}

		.flex-wrapper {
			flex-direction: column;
			/* 垂直排列 */
		}

		.crop-section,
		.preview-section {
			width: 100% !important;
			/* 强制全宽 */
			min-width: unset;
			/* 移除最小宽度 */
		}

		#image-wrapper {
			height: 50vh;
			/* 调整移动端高度 */
		}

		.preview-box {
			width: 120px;
			/* 缩小预览区域 */
			height: 120px;
		}

		.controls {
			flex-direction: column;
			/* 垂直排列按钮 */
			margin-top: 10px;
		}
	}


	#cropper-image {
		max-width: none !important;
		max-height: none !important;
	}

	/* 控制区域 */
	.controls {
		margin-top: 20px;
		display: flex;
		flex-direction: column;
		grid-template-columns: repeat(3, 1fr);
		gap: 10px;
	}

	/* 预览区域 */
	.preview-section {
		flex: 1;
		min-width: 50px;
	}

	.preview-box {
		width: 150px;
		height: 150px;
		/* border-radius: 50%; */
		border: 3px solid var(--primary-color);
		overflow: hidden;
		margin: 0 auto 20px;
	}

	#preview-img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}



	/* 上传进度 */
	.progress-container {
		height: 8px;
		background: #eee;
		border-radius: 4px;
		margin-top: 20px;
		overflow: hidden;
		display: none;
	}

	.progress-bar {
		width: 0%;
		height: 100%;
		background: var(--primary-color);
		transition: width 0.3s ease;
	}
</style>