package routers

import (
	"io"
	"net/http"
	"ops/models"
	"path"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type TabFileInfo_ struct {
	ID     uint      `gorm:"primaryKey;autoIncrement"`
	Name   string    `gorm:"not null;size:256;index"` // 前端报告的文件名
	Path   string    `gorm:"not null;size:300"`       //
	Sha256 string    `gorm:"not null;size:64;index"`  //
	Mime   string    `gorm:"size:64;index"`
	Type   string    `gorm:"size:64;index"`
	Const  uint      `gorm:"default:1;index"`
	Per    uint      `gorm:"default:1"`
	UserID uint      `gorm:"not null;index"`
	Date   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"` // 默认当前时间
}

func file_save() {

}

func ApiFilesInit() {

	models.DB.AutoMigrate(&TabFileInfo_{})
}

func ApiFiles(r *gin.RouterGroup) {

	//getfile := r.Group("/get") //定义上传组
	r.GET("/:mode/:hash", func(ctx *gin.Context) {
		hash := ctx.Param("hash")
		mode := ctx.Param("mode")
		// filename := ctx.Param("filename")
		// fmt.Println(filename)

		download := false
		isPartOK := false

		if mode == "get" {
			isPartOK = true
			download = true
		}
		if mode == "download" {
			isPartOK = true
			download = false
		}
		if isPartOK {
			file_info := TabFileInfo_{
				Sha256: hash,
			}
			if models.DB.Where(&file_info).First(&file_info).Error == nil {
				ReturnFile(ctx, &file_info, download)
			} else {
				//fmt.Println("not fund")
				ReturnJson(ctx, "file_not_found", nil)
			}
		} else {
			ReturnJson(ctx, "file_part_err", nil)
		}

	})

	upload := r.Group("/upload") //定义上传组
	//上传文件的总接口，能上传什么文件应该由后端决定，前端仅做相应限制

	upload.POST("/image", func(ctx *gin.Context) {

		cookie := ctx.PostForm("cookie") //首先需要判断用户是否登录

		//通过cookie获取用户信息
		user, err := AuthenticationAuthorityFromCookie(cookie)
		if err == nil {
			file, err := ctx.FormFile("file")

			if err == nil {
				if file.Filename != "" {
					//限制文件大小
					if file.Size > 512 {
						if file.Size < int64(models.ConfigsFile.MaxSize) {

							//判断文件mime是否合法
							// 打开文件流
							src_mime, _ := file.Open()
							defer src_mime.Close()
							// 读取前512字节用于MIME检测
							buffer := make([]byte, 512)
							io.ReadFull(src_mime, buffer)
							// 检测MIME类型
							mimeType := http.DetectContentType(buffer)
							file_extname := models.ConfigsFile.AllowImageMime[mimeType]
							if file_extname != "" {
								filename := filepath.Base(file.Filename) // 防御性处理路径分隔符
								// 计算哈希值
								hash_str, err := models.SHA256HashFile(file)
								if err == nil {
									//fmt.Println(hash_str)
									//fmt.Println(filename)
									//这是上传的真实路径
									dst := path.Join(models.ConfigsFile.Pahts["image"], hash_str)
									//fmt.Println(dst)
									//判断文件是否存在避免重复保存
									if models.FileExists(dst) {
										//fmt.Println("文件存在")

									} else {
										//fmt.Println("文件no存在")
										ferr := ctx.SaveUploadedFile(file, dst)
										if ferr == nil {
											//文件保存成功

										} else {

											ReturnJson(ctx, "file_save_err", nil)
											ctx.Abort() //end
											return
										}
									}
									//记录到数据库
									//先检查数据库有没有数据
									fund_file_info := TabFileInfo_{
										Name:   filename,
										Sha256: hash_str,
										Mime:   mimeType,
										Type:   "image",
										UserID: user.ID,
									}
									fund_file_info2 := TabFileInfo_{}

									models.DB.Where(&fund_file_info).Find(&fund_file_info2)

									if fund_file_info2.ID != 0 {
										fund_file_info2.Const += 1
										models.DB.Where(&fund_file_info).Updates(&fund_file_info2)
									} else {
										fund_file_info.Path = dst
										models.DB.Create(&fund_file_info) // 传入指针
									}

									//返回后台存储的URL
									download_URL := path.Join("/api/files/download/", hash_str)
									get_URL := path.Join("/api/files/get/", hash_str)
									re := map[string]interface{}{
										"download": download_URL,
										"get":      get_URL,
										"hash":     hash_str,
									}

									ReturnJson(ctx, "apiOK", re)

								} else {
									ReturnJson(ctx, "file_hash_err", nil)
								}

							} else {
								ReturnJson(ctx, "file_mime_err", nil)
							}
						} else {
							ReturnJson(ctx, "file_size_err", nil)
						}

					} else {
						ReturnJson(ctx, "file_size_err", nil)
					}

				} else {
					ReturnJson(ctx, "file_name_err", nil)
				}

			} else {
				ReturnJson(ctx, "file_get_err", nil)
			}

		} else {
			ReturnJson(ctx, "userCookieError", nil)
		}

		//ReturnJson(ctx, "apiErr", nil)
	})

	// r.GET("/upload", func(ctx *gin.Context) {
	// 	ReturnJson(ctx, "apiOK", nil)
	// })

}
