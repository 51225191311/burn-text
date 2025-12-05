package main

import (
	"burn-text/internal/config"
	"burn-text/internal/handler"
	"burn-text/storage"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	//初始化配置
	config.InitConfig()

	//初始化Redis
	if err := storage.InitRedis(); err != nil {
		fmt.Printf("连接Redis失败: %v\n", err)
		fmt.Println("请检查Docker是否启动成功，以及是否执行了docker run ...")
		return
	}
	fmt.Println("Redis连接成功")

	//设置Gin模式
	gin.SetMode(config.GlobalConfig.Server.Mode)
	r := gin.Default()

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
		api.POST("/burn", handler.CreateSecret)
		api.GET("/view/:id", handler.GetSecret)
	}

	port := config.GlobalConfig.Server.Port
	fmt.Printf("服务启动在 http://localhost:%s\n", port)
	r.Run(":" + port)
}
