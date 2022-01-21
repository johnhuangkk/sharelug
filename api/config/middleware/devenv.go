package middleware

import (
	"github.com/gin-gonic/gin"
)

func GetDevSiteAllow(ctx *gin.Context) {
	//log.Debug("Request URL", ctx.Request.RequestURI)
	//ENV := viper.GetString("ENV")
	//if ENV == "local" {
	//	ctx.Next()
	//	return
	//}
	//key ,err := ctx.Cookie("dev_cookie")
	//if err != nil || key != "allow"  {
	//	log.Debug("cookie Error", err)
	//	ctx.Redirect(302,"/v1/devview")
	//	return
	//}
	ctx.Next()
}
