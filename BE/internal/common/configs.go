package common

import (
	"os"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// 基础配置
type Basis struct {
	Port           int  // 端口
	Debug          bool // 调试模式
	SessionTimeout int  // 会话超时时间
}

// 数据源
type Datasource struct {
	Host     string // 主机名
	Port     int    // 端口
	Database string // 数据库
	Username string // 用户名
	Password string // 密码
	Debug    bool   // 调试模式
}

// 微信
type WeiXin struct {
	APPID  string // ID
	Secret string // 密钥
}

// 全局变量
var (
	basis      Basis
	datasource Datasource
	weixin     WeiXin
)

// 初始化函数
func init() {
	// 获取工作目录
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// 设置配置文件
	viper.SetConfigName("configs")
	viper.SetConfigType("yml")
	viper.AddConfigPath("../configs/")
	viper.AddConfigPath(fmt.Sprintf("%s/configs/", cwd[0:strings.Index(cwd, "BE")+2]))

	// 读取配置文件
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// 基础配置
	basis = Basis{
		Port:           viper.GetInt("BASIS.PORT"),
		Debug:          viper.GetBool("BASIS.DEBUG"),
		SessionTimeout: viper.GetInt("BASIS.SESSION_TIMEOUT"),
	}

	// 数据源
	datasource = Datasource{
		Host:     viper.GetString("DATASOURCE.HOST"),
		Port:     viper.GetInt("DATASOURCE.PORT"),
		Database: viper.GetString("DATASOURCE.DATABASE"),
		Username: viper.GetString("DATASOURCE.USERNAME"),
		Password: viper.GetString("DATASOURCE.PASSWORD"),
		Debug:    viper.GetBool("DATASOURCE.DEBUG"),
	}

	// 微信
	weixin = WeiXin{
		APPID:  viper.GetString("WEIXIN.APPID"),
		Secret: viper.GetString("WEIXIN.SECRET"),
	}
}

// 返回基础配置
func GetBasis() Basis {
	return basis
}

// 返回数据源
func GetDatasource() Datasource {
	return datasource
}

// 返回微信
func GetWeiXin() WeiXin {
	return weixin
}
