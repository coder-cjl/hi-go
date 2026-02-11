package middleware

import (
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	allowedOriginsOnce sync.Once
	allowedOriginsMap  map[string]struct{}
	allowAnyOrigin     bool
	allowCredentials   bool
)

func loadCORSConfig() {
	allowedOriginsMap = make(map[string]struct{})

	rawOrigins := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if rawOrigins != "" {
		for _, origin := range strings.Split(rawOrigins, ",") {
			origin = strings.TrimSpace(origin)
			if origin == "" {
				continue
			}
			if origin == "*" {
				allowAnyOrigin = true
				continue
			}
			allowedOriginsMap[origin] = struct{}{}
		}
	}

	allowCredentials = strings.EqualFold(strings.TrimSpace(os.Getenv("CORS_ALLOW_CREDENTIALS")), "true")

	// If credentials are allowed, never allow wildcard origin.
	if allowCredentials {
		allowAnyOrigin = false
	}
}

// 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOriginsOnce.Do(loadCORSConfig)

		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")

		if origin != "" {
			_, allowed := allowedOriginsMap[origin]
			if allowAnyOrigin || allowed {
				if allowAnyOrigin {
					c.Header("Access-Control-Allow-Origin", "*")
				} else {
					c.Header("Access-Control-Allow-Origin", origin)
				}
				c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
				c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
				c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
				if allowCredentials {
					c.Header("Access-Control-Allow-Credentials", "true")
				}
			}
		}

		// 放行所有 OPTIONS 方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
