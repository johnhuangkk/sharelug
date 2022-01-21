package middleware

import (
	"api/services/util/log"
	"api/services/util/response"
	"api/services/util/tools"
	"github.com/gin-gonic/gin"
)

// 允許連線的IP
func AllowIP() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp := response.New(ctx)
		allowIP := []string {
			"172.27.0.1", // local
			"210.200.25.105", // i 郵局 正式
			"61.220.55.12",
		}

		// 郵局測試
		ipostDevIP := []string {
			"210.200.25.252",
			"210.200.26.252",
			"210.200.27.252",
			"124.219.114.252",
			"124.219.115.252",
		}

		allowIP = append(allowIP, ipostDevIP...)

		log.Debug("connect server ip", allowIP)
		log.Debug("connect server ip", ctx.ClientIP())
		if !tools.InArray(allowIP, ctx.ClientIP()) {
			resp.Conflict("不允許連線").Send()
			return
		}
		ctx.Next()
	}
}

// 允許連線的IP
func ErpAllowIP() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp := response.New(ctx)
		allowIP := []string {
			"172.27.0.1", // local
			"210.200.25.105", // i 郵局 正式
		}
		if !tools.InArray(allowIP, ctx.ClientIP()) {
			resp.Conflict("不允許連線").Send()
			return
		}
		ctx.Next()
	}
}