package main

import (
	"api"
	"api/config/middleware"
	"api/config/router"
	"api/cron"
	"api/services/util/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	defer func() {
		if err := recover(); err != nil {
			log.Error("recover error", err)
		}
	}()

	httpsRouter := api.NewDevelopment()
	router.SetupRouter(httpsRouter)

	httpRouter := api.NewDevelopment()
	httpRouter.GET("/", middleware.GetDevSiteAllow, func(ctx *gin.Context) {
		log.Debug("aaa => ", ctx.Request.Host)
		ctx.Redirect(http.StatusMovedPermanently, "https://"+ctx.Request.Host+"/"+ctx.Param("variable"))
	})

	go cron.Run()

	err := httpsRouter.Run(":8001")
	if err != nil {
		log.Error("https router run Error", err)
	}
}
