package models

import "github.com/mitchellh/mapstructure"

var Configs map[string]interface{}

//mime信息转换位拓展名

type ConfigsWeb_ struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Tls             bool   `mapstructure:"tls"`
	CertPrivatePath string `mapstructure:"certPrivatePath"`
	CertPublicPath  string `mapstructure:"certPublicPath"`
}

type ConfigsUser_ struct {
	CookieTimeout int    `mapstructure:"cookieTimeout"`
	PassHashType  string `mapstructure:"passHashType"`
}

type ConfigsFile_ struct {
	MaxSize        uint64            `mapstructure:"maxSize"`
	Pahts          map[string]string `mapstructure:"pahts"`
	AllowImageMime map[string]string `mapstructure:"allowImageMime"`
	AllowVideoMime map[string]string `mapstructure:"allowVideoMime"`
	AllowMusicMime map[string]string `mapstructure:"allowMusicMime"`
	AllowPdfMime   map[string]string `mapstructure:"allowPdfMime"`
}

var ConfigsWed ConfigsWeb_
var ConfigsUser ConfigsUser_
var ConfigsFile ConfigsFile_

func ConfigAllInit() error {

	//初始化数据库
	DatabaseInit()

	//读取web配置
	err := mapstructure.Decode(Configs["web"].(map[string]interface{}), &ConfigsWed)
	if err != nil {
		panic(err)
	}

	//初始化user config
	err = mapstructure.Decode(Configs["user"].(map[string]interface{}), &ConfigsUser)
	if err != nil {
		panic(err)
	}

	//初始化file config
	err = mapstructure.Decode(Configs["file"].(map[string]interface{}), &ConfigsFile)
	if err != nil {
		panic(err)
	}

	return nil
}
