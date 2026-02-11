package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "hi-go/docs" // 引入生成的docs包
)

// SetupDocsRoutes 设置 API 文档路由
func SetupDocsRoutes(r *gin.Engine) {
	// Swagger UI（原始界面，支持在线测试）
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Redoc（更美观友好的文档界面）
	r.GET("/docs", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, `<!DOCTYPE html>
<html>
<head>
    <title>Hi-Go API 文档</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
    <style>
      body {
        margin: 0;
        padding: 0;
      }
    </style>
</head>
<body>
    <redoc spec-url='/docs/swagger.json'></redoc>
    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"> </script>
</body>
</html>`)
	})

	// 提供 swagger.json 文件访问
	r.StaticFile("/docs/swagger.json", "./docs/swagger.json")
}
