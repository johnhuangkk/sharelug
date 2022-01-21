package middleware

import (
	"api/services/dao/Store"
	"api/services/dao/member"
	"api/services/database"
	"api/services/entity"
	"api/services/errorMessage"
	"api/services/util/response"
	"api/services/util/session"
	"github.com/gin-gonic/gin"
)

//var uid string
//var sid string

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp := response.New(ctx)
		// 抓UUID
		header := ctx.Request.Header
		if header.Values("Sharelug-Token") == nil || len(header.Values("Sharelug-Token")[0]) == 0 {
			resp.Fail(errorMessage.GetMessageByCode(1001000)).Send()
			ctx.Abort()
			return
		}
		// 設定 uid & sid
		if len(GetUserData(ctx).Uid) == 0 {
			resp.Fail(errorMessage.GetMessageByCode(1001000)).Send()
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

//// 設定 uid & sid
//func setUidAndSid(ctx *gin.Context) error {
//	uid, sid := GetSession(ctx)
//	if SessionValue != "" && Token != "" {
//		uid = session.GetSession(SessionValue, "uid")
//		sid = session.GetSession(SessionValue, "sid")
//		log.Debug("auth get uid =>", uid, sid)
//	} else {
//		uid = ""
//		sid = ""
//		log.Debug("尚未登入", SessionValue, Token, GetClientIP())
//		return fmt.Errorf("尚未登入")
//	}
//	return nil
//}

func GetSession(ctx *gin.Context) (string, string) {
	header := ctx.Request.Header
	if header.Values("Sharelug-Id") != nil {
		SessionValue := header.Values("Sharelug-Id")[0]
		uid := session.GetSession(SessionValue, "uid")
		sid := session.GetSession(SessionValue, "sid")
		return uid, sid
	} else {
		return "", ""
	}
}

func GetTokenValue(ctx *gin.Context) string {
	header := ctx.Request.Header
	if header.Values("Sharelug-Id") != nil {
		Token := header.Values("Sharelug-Token")[0]
		return Token
	} else {
		return ""
	}
}

func GetSessionValue(ctx *gin.Context) string {
	header := ctx.Request.Header
	if header.Values("Sharelug-Id") != nil {
		SessionValue := header.Values("Sharelug-Id")[0]
		return SessionValue
	} else {
		return ""
	}
}

func GetUserData(ctx *gin.Context) entity.MemberData {
	uid, _ := GetSession(ctx)
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var memberData entity.MemberData
	if uid != "" {
		memberData, _ = member.GetMemberDataByUid(engine, uid)
	}
	return memberData
}

func GetStoreData(ctx *gin.Context) entity.StoreDataResp {
	uid, sid := GetSession(ctx)
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var storeData entity.StoreDataResp
	if sid != "" {
		storeData, _ = Store.GetStoreDataByUserIdAndStoreId(engine, uid, sid)
	}
	return storeData
}

//暫時使用的商店資料
//func GetStoreName(uid string) entity.StoreData {
//	engine := database.GetMysqlEngine()
//	defer engine.Close()
//	var storeData entity.StoreData
//	storeData, err := Store.GetStoreDataByUid(engine, uid)
//	if err != nil {
//		log.Error("Get store Data Error", err)
//	}
//	return storeData
//}
