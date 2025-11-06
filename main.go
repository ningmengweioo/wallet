package main

import (
	"fmt"
	"log"

	"wallet/config"
	"wallet/router"
)

func main() {
	// 初始化配置
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// 初始化数据库
	_, err := config.InitDB(config.GetConf())
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 设置路由
	r := router.SetupRouter()

	// 获取配置的端口
	port := config.GetConf().Http.Port

	// 启动服务器
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
