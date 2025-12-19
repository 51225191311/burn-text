package main

import (
	"burn-text/internal/config"
	"burn-text/internal/handler"
	"burn-text/internal/logger"
	"burn-text/internal/middleware"
	"burn-text/storage"
	"context"
	"os"
	"os/signal"
	"syscall"
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

	//增加优雅关机
	port := config.GlobalConfig.Server.Port
	logger.Log.Info("服务器启动", zap.String("port", port))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	//在独立的Goroutine中启动服务，让主线程不被阻塞继续监听
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("服务启动失败", zap.Error(err))
		}
	}()

	//创建信号通道
	quit := make(chan os.Signal, 1)

	//监听Ctrl+C和Docker stop
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit //阻塞直到收到信号
	logger.Log.Info("正在关闭服务...")

	//创建5秒的超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal("服务被强制关闭", zap.Error(err))
	}
	logger.Log.Info("服务已退出")
}
