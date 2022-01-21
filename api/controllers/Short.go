package controllers

import (
	"api/services/dao/Short"
	"api/services/database"
	"api/services/util/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ShortUrlAction(ctx *gin.Context) {
	short := ctx.Param("short")
	log.Debug("short", short)
	engine := database.GetMysqlEngine()
	defer engine.Close()

	data, err := Short.GetShortUrlDataByShort(engine, short)
	if err != nil {
		log.Debug("short error1", short)
		ctx.Redirect(http.StatusFound, "/error/404")
	}
	log.Debug("short", data, len(data.Url))
	if len(data.Url) == 0 {
		log.Debug("short error2", short)
		ctx.Redirect(http.StatusFound, "/error/404")
	} else {
		ctx.Redirect(http.StatusMovedPermanently, data.Url)
	}
}




