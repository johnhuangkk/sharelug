package middleware

import (
	"api/services/util/log"
	"api/services/util/session"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		setCookie(ctx)
		getCookie(ctx)
	}
}

// set cookie
func setCookie(ctx *gin.Context) {
	var cookieName = viper.GetString("cookie.name")
	cookie, _ := ctx.Request.Cookie(cookieName)
	if cookie == nil {
		uuidStr, err := uuid.NewUUID()
		if err != nil {
			log.Error("Generate Cookie error", err)
		}
		ctx.SetCookie(cookieName, uuidStr.String(), 3600, "/", ".checkne.com", false, true)
	}
}

// 抓UUID
func getCookie(ctx *gin.Context) {
	var SessionValue string
	var Token string

	header := ctx.Request.Header
	s := strings.SplitAfter(ctx.Request.RequestURI, "/")
	if s[1] == "v1/" || s[1] == "gw/" {
		if header.Values("Sharelug-Id") != nil {
			SessionValue = header.Values("Sharelug-Id")[0]
			SharelugToken := header.Values("Sharelug-Token")[0]
			log.Info("Header UUID And Token ", ctx.Request.Method, ctx.Request.RequestURI, SessionValue, SharelugToken, GetClientIP())
			SessionToken := session.GetSession(SessionValue, "token")
			//檢查 TOKEN 是否正確
			if SharelugToken != SessionToken {
				oldToken := session.GetOldSession(SessionValue, "token")
				if SharelugToken != oldToken {
					Token = ""
				} else {
					Token = SharelugToken
				}
			} else {
				Token = SharelugToken
			}
			log.Info("header Get UUID And Token", SessionValue, Token)
		} else {
			log.Error("header Get Error", header.Values("Sharelug-Id"), header.Values("Sharelug-Token"))
			SessionValue = ""
			Token = ""
		}
	} else {
		//log.Info("URL", ctx.Request.Method, ctx.Request.RequestURI, GetClientIP())
	}
}
