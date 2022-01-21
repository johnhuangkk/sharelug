package model

import (
	"api/config/middleware"
	"api/services/Enum"
	"api/services/Service/Notification"
	"api/services/VO/Response"
	"api/services/VO/TokenVo"
	"api/services/dao/device"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/session"
	"api/services/util/tools"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

//更換UUID
func HandleChangeUuid(ctx *gin.Context, params TokenVo.TokenParams, ip string) (TokenVo.ChangeUuidResponse, error) {
	Session := generateUUID(ctx, params, ip)
	log.Info("Response Session Value", Session)
	response := TokenVo.ChangeUuidResponse {
		UUID: Session,
	}
	return response, nil
}

func HandleNotice(userData entity.MemberData, storeData entity.StoreDataResp) (Response.NoticeResponse, error) {
	var resp Response.NoticeResponse
	resp.Message = GetMessage(userData, storeData)
	return resp, nil
}
//檢查Token是否過期需更換
func HandleChangeToken(ctx *gin.Context, params TokenVo.TokenParams, ip string) (TokenVo.ChangeTokenResponse, error) {
	Count := int64(0)
	Session := generateUUID(ctx, params, ip)
	NewToken := checkToken(ctx, Session, ip)
	if len(NewToken) != 0 {
		userData := middleware.GetUserData(ctx)
		storeData := middleware.GetStoreData(ctx)
		Count = GetMessage(userData, storeData)
	}
	log.Info("Response Session Value", Session)
	log.Info("Response Token Value", NewToken)
	response := TokenVo.ChangeTokenResponse {
		UUID: Session,
		Token: NewToken,
		Message: Count,
	}
	return response, nil
}

// 產生UUID
func generateUUID(ctx *gin.Context, params TokenVo.TokenParams, ip string) string {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	Session := middleware.GetSessionValue(ctx)
	if len(Session) == 0 {
		data, err := device.InsertDevice(engine, params, ip)
		if err != nil {
			log.Error("Insert Device Data Error", err)
		}
		log.Info("New Session Value", Session)
		Session = data.DeviceUuid
	}
	log.Info("Get Session Value", Session)
	return Session
}
//驗證TOKEN是否過期
func IsVerifyTokenIsExpire(Session, token, Ip string) (bool, tools.Claims) {
	//解析TOKEN
	claim, _ := tools.ParseToken(token)
	var data tools.Claims
	claims, _ := json.Marshal(claim)
	_ = json.Unmarshal(claims, &data)
	if time.Now().Before(data.Exp) && Session == data.UUid {
		log.Info("Get IP UUID Exp", data.Exp.Sub(time.Now()) < time.Minute * 60, data.Exp.Sub(time.Now()), Ip,  Session)
		return false, data
	} else {
		log.Info("Get IP UUID Exp", data.Exp.Sub(time.Now()) < time.Minute * 60, data.Exp.Sub(time.Now()), Ip,  Session)
	}
	return true, data
}
//檢查到期時間及IP位置 更新TOKEN FIXME
func checkToken(ctx *gin.Context, Session string, Ip string) string {
	token := middleware.GetTokenValue(ctx)
	var NewToken string
	if token != "" {
		expire, data := IsVerifyTokenIsExpire(Session, token, Ip)
		if !expire {
			if data.Exp.Sub(time.Now()) < time.Minute * 30 {
				if err := session.OldUser(Session).Put("token", token); err != nil {
					log.Error("Session Set Old Error", err)
				}
				NewToken = GeneratorToken(Session, Ip, data.UserId, data.StoreId)
			} else {
				NewToken = token
			}
		} else {
			if err := session.NewUser(Session).Destroy(); err != nil{
				log.Error("Session Destroy Error", err)
			}
			NewToken = ""
		}
	} else {
		if err := session.NewUser(Session).Destroy(); err != nil{
			log.Error("Session Destroy Error", err)
		}
	}
	return NewToken
}
//產生新的TOKEN
func GeneratorToken(Session, Ip, UserId, StoreId string) string {
	token, _ := tools.GeneratorJWT(Session, UserId, StoreId, Ip)
	_ = session.NewUser(Session).Put("token", token)
	return token
}

func GetMessage(userData entity.MemberData, storeData entity.StoreDataResp) int64 {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	counts := int64(0)
	if storeData.Rank != Enum.StoreRankMaster {
		//改讀資料庫
		user := []string{storeData.StoreId}
		count, err := Notification.GetNotifyCount(engine, user)
		if err != nil {
			log.Debug("Get Redis Message Count Error", err)
			return 0
		}
		counts = count
	} else {
		user := []string{storeData.StoreId, userData.Uid}
		count, err := Notification.GetNotifyCount(engine, user)
		if err != nil {
			log.Debug("Get Redis Message Count Error", err)
			return 0
		}
		counts = count
	}
	return counts
}
//檢查UUID是否已登入
func HandleCheckToken(params TokenVo.CheckTokenParams) error {
	SessionToken := session.GetSession(params.Uuid, "token")
	if len(SessionToken) == 0 {
		log.Error("check token", SessionToken)
		return fmt.Errorf("1001001")
	}
	if SessionToken != params.Token {
		log.Error("check token", SessionToken, params.Token)
		return fmt.Errorf("1001001")
	}
	return nil
}