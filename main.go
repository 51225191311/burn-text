package main

import (
	"burn-text/internal/config"
	"burn-text/internal/handler"
	"burn-text/internal/logger"
	"burn-text/internal/middleware"
	"burn-text/storage"
	"time"

	ginzap "github.com/gin-contrib/zap"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	//初始化配置
	config.InitConfig()

	//初始化日志
	logger.InitLogger(config.GlobalConfig.Server.Mode)
	defer logger.Log.Sync() //退出时刷新缓存

	//初始化Redis
	if err := storage.InitRedis(); err != nil {
		logger.Log.Fatal("Redis连接失败", zap.Error(err))
		return
	}
	logger.Log.Info("Redis连接成功", zap.String("addr", config.GlobalConfig.Redis.Addr))

	//设置Gin模式
	gin.SetMode(config.GlobalConfig.Server.Mode)
	r := gin.New()

	r.Use(ginzap.Ginzap(logger.Log, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger.Log, true))

	//加载HTML模板，并且HTML文件都在templates目录下
	r.LoadHTMLGlob("templates/*")

	//首页
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	//当用户访问http://localhost:8080/xxx时，返回view.html页面
	r.GET("/view/:id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "view.html", nil)
	})

	//创建接口
	api := r.Group("/api")
	{
		api.POST("/burn", middleware.IPwRateLimiter(), handler.CreateSecret)
		api.GET("/view/:id", handler.GetSecret)
	}

	port := config.GlobalConfig.Server.Port
	logger.Log.Info("服务器启动", zap.String("port", port))
	r.Run(":" + port)
}
