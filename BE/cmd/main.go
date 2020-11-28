package main

import (
	"BE/internal/common"
	"BE/internal/routers"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func main() {
	// 获取配置
	basis := common.GetBasis()

	// 初始化翻译器
	err := common.InitTrans("zh")
	if err != nil {
		fmt.Printf("init trans failed, err: [%v]\n", err.Error())
		return
	}

	// 运行模式
	if basis.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 默认实例
	engine := gin.Default()

	// 中间件
	// 创建基于内存的存储引擎, 加密密钥: secret
	store := memstore.NewStore([]byte("secret"))
	// 设置过期时间
	store.Options(sessions.Options{
		MaxAge: basis.SessionTimeout * 60, // 超时时间
		Path:   "/",                       // 路径
	})

	// 会话名称: session
	engine.Use(sessions.Sessions("session", store))

	// 中间件: 允许所有来源
	engine.Use(cors.Default())

	// 注册路由
	routers.Register(engine)

	// 运行服务
	engine.Run(fmt.Sprintf(":%d", basis.Port))
}
