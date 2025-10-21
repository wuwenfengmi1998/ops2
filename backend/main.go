package main

import (
	"fmt"
	"net/http"
	"ops/models"
	"ops/routers"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-yaml"
)

func main() {

	fmt.Println("OPS Backend Service started")

	// 直接尝试创建所有必要的目录
	err := os.MkdirAll("./data", 0755)
	if err != nil {
		fmt.Printf("创建文件夹失败: %v\n", err)
		panic("创建文件夹失败")

	}

	config_file_path := "./data/config.yaml"
	config_temp_path := "./defConfig/configTemp.yaml"

	//尝试读取配置
	if !models.FileExists(config_file_path) {
		fmt.Println("读取配置失败")

		//复制配置模板
		fmt.Println("复制配置模板")
		input, err := os.ReadFile(config_temp_path)
		if err != nil {
			panic(err)

		}
		err = os.WriteFile(config_file_path, input, 0644)
		if err != nil {
			panic(err)

		}
		fmt.Printf("需要修改此配置:%s\n", config_file_path)

	}

	//读取默认配置
	data, err := os.ReadFile(config_file_path)
	if err != nil {
		panic(err)

	}

	if err := yaml.Unmarshal(data, &models.Configs); err != nil {
		panic(err)

	}

	if models.Configs["configed"] == false {
		fmt.Printf("需要将:%s 内的configed设置为true", config_file_path)
		panic("need config")
	}

	//统一初始化
	models.ConfigAllInit()

	//启动gin服务
	r := gin.Default()

	// 静态文件服务
	fs := http.FileServer(http.Dir("./dist"))
	// 中间件处理路由
	r.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Next() // 继续处理API请求
			return
		}

		// 处理静态文件
		fs.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	})

	// API路由
	routers.ApiRoot(r.Group("/api/"))

	var http_port = models.ConfigsWed.Host + ":" + models.ConfigsWed.Port
	var gin_port = "0.0.0.0" + ":" + models.ConfigsWed.Port
	if models.ConfigsWed.Tls {
		if models.ConfigsWed.CertPublicPath == "" || models.ConfigsWed.CertPrivatePath == "" {
			fmt.Printf("需要配置证书路径")
			return
		} else {
			fmt.Println("https://" + http_port)
			r.RunTLS(gin_port, models.ConfigsWed.CertPublicPath, models.ConfigsWed.CertPrivatePath)
		}
	} else {
		fmt.Println("http://" + http_port)
		r.Run(gin_port)
	}
}
