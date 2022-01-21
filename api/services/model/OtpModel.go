package model

import (
	"api/config/middleware"
	"api/services/Enum"
	"api/services/Service/Mail"
	"api/services/Service/Sms"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Store"
	"api/services/dao/member"
	"api/services/dao/otp"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/session"
	"api/services/util/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"github.com/spf13/viper"
	"time"
)

//type HOTP struct {
//	Secret         string // The secret used to generate the token
//	Length         uint8  // The token size, with a maximum determined by MaxLength
//	Counter        uint64 // The counter used as moving factor
//	IsBase32Secret bool   // If true, the secret will be used as a Base32 encoded string
//}

//產生 OTP Code
func GenerateOtpCode() string {
	code, err := totp.GenerateCode("asdfghj", time.Now())
	if err != nil {
		log.Error("Generate Otp Code Error", err)
	}
	return code
}

//發送OTP SMS
func HandlePushOtpSms(engine *database.MysqlSession, userData entity.MemberData, phone string) (entity.OtpData, error) {
	var resp entity.OtpData
	//檢查登入失敗次數
	if err := validateLoginErrorFrequency(userData); err != nil {
		log.Error("validate otp sms error", err)
		return resp, fmt.Errorf("1002006")
	}
	resp, err := sendOtpSms(engine, phone, userData)
	if err != nil {
		log.Error("Send Otp SMS error", err)
		return resp, fmt.Errorf("1001001")
	}
	return resp, nil
}

func sendOtpSms(engine *database.MysqlSession, phone string, userData entity.MemberData) (entity.OtpData, error) {
	var resp entity.OtpData
	//取出之前的OTP
	data, err := otp.GetOtpByPhone(engine, phone)
	if err != nil {
		log.Error("get otp db error", err)
		return resp, fmt.Errorf("系統錯誤")
	}
	now := time.Now()
	if !now.Before(data.ExpireTime) {
		//產生OTP
		code := GenerateOtpCode()
		//寫入資料庫
		email := ""
		if len(userData.Email) != 0 {
			email = userData.Email
		}
		resp, err = otp.InsertOtp(engine, phone, code, userData.Uid, email)
		if err != nil {
			log.Error("Insert otp sms error", err)
			return data, err
		}
		data = resp
	} else {
		resp = data
		data.SendFreq += 1
		if err := otp.UpdateOtpData(engine, data); err != nil {
			return data, err
		}
	}
	//判斷會員是否有EMAIL
	if len(userData.Email) != 0 {
		ENV := viper.GetString("ENV")
		if ENV != "prod" {
			return data, nil
		}
		if err := Mail.SendOtpMail(userData, resp.OtpNumber); err != nil {
			content := fmt.Sprintf("你的Check'Ne驗證碼為%s，請於15分鐘內輸入。", resp.OtpNumber)
			if err := Sms.PushMessageSms(phone, content); err != nil {
				log.Error("Push SMS Error", err)
				return data, err
			}
			data.Email = ""
		}
	} else {
		if resp.SendFreq > 2 {
			return data, nil
		}
		log.Debug("發送簡訊!!")
		//發送簡訊 fix
		content := fmt.Sprintf("你的Check'Ne驗證碼為%s，請於15分鐘內輸入。", resp.OtpNumber)
		err = Sms.PushMessageSms(phone, content)
		if err != nil {
			log.Error("Push SMS Error", err)
			return data, err
		}
	}
	return data, nil
}

//處理OTP驗證
func HandleValidateOtp(SessionValue, code, phone, ip string) (*Response.ValidateOtpResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp *Response.ValidateOtpResponse
	userData, err := checkOtp(engine, code, phone)
	if err != nil {
		log.Info("會員登入", phone, middleware.GetClientIP(), Enum.SysLogFail)
		return resp, fmt.Errorf("1002003")
	}
	log.Info("會員登入", phone, userData.Uid, middleware.GetClientIP(), Enum.SyslogSuccess)
	//之後需改為使用者選擇
	storeData, _ := GetStoreData(engine, userData.Uid)
	token, err := SetSession(SessionValue, userData.Uid, storeData.StoreId, ip)
	if err != nil {
		log.Error("Set Session Error", err)
		return resp, fmt.Errorf("1002001")
	}
	info, err := GetMemberInfo(engine, userData, storeData.StoreId)
	if err != nil {
		log.Error("Get Member Error", err)
		return resp, fmt.Errorf("1002001")
	}
	data := &Response.ValidateOtpResponse{
		Uuid: SessionValue,
		Token: token,
		Member: info.Member,
		Store: info.Store,
	}
	return data, nil
}

//token set Session
func SetSession(SessionValue, UserId, StoreId, Ip string) (string, error) {
	//todo 登入成功 APP 要產生 TOKEN 存入資料庫 SESSION 記錄
	token, err := tools.GeneratorJWT(SessionValue, UserId, StoreId, Ip)
	if err != nil {
		log.Debug("Generator Jwt Token Error", err)
		return token, err
	}
	sess := session.Session{
		Name:"user" + SessionValue,
		TTL: 0,
	}
	_ = sess.Put("uid", UserId)
	_ = sess.Put("uuid", SessionValue)
	_ = sess.Put("sid", StoreId)
	_ = sess.Put("token", token)
	return token, nil
}

//切換收銀機
func HandleExchangeStore(ctx *gin.Context, params *Request.ExchangeStoreParams) (Response.ValidateOtpResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()

	var resp Response.ValidateOtpResponse
	SessionValue := middleware.GetSessionValue(ctx)
	userData := middleware.GetUserData(ctx)
	storeData, err := GetStoreDataByStoreId(engine, params.StoreId, userData)
	if err != nil {
		return resp, err
	}

	count, _ := Store.CountStoreRankByUid(engine, userData.Uid)
	storeData.Count = count
	token, err := SetSession(SessionValue, userData.Uid, storeData.Sid, middleware.GetClientIP())
	if err != nil {
		return resp, fmt.Errorf("1001001")
	}
	resp.Uuid = SessionValue
	resp.Token = token
	resp.Member = userData.GetMemberLoginInfo()
	resp.Store = storeData

	return resp, nil
}

//驗證OTP
func checkOtp(engine *database.MysqlSession, code string, phone string) (entity.MemberData, error) {
	var resp entity.MemberData
	//取得OTP資料
	data, err := otp.GetOtpByPhone(engine, phone)
	if err != nil {
		log.Error("get otp db error", err)
		return resp, fmt.Errorf("1001001")
	}
	now := time.Now()
	if !now.Before(data.ExpireTime) {
		return resp, fmt.Errorf("1002007")
	}
	//取出會員資料
	userData, err := member.GetMemberDataByPhone(engine, data.Phone)
	if err != nil {
		log.Debug("Get Member Data Error", err)
		return resp, fmt.Errorf("1001001")
	}
	if userData.Uid == "" {
		log.Debug("Get Member Uid Error")
		return resp, fmt.Errorf("1001001")
	}
	if err := verifyCode(engine, code, data, userData); err != nil {
		return resp, fmt.Errorf("1002007")
	}
	return userData, nil
}

//驗證 OTP code
func verifyCode(engine *database.MysqlSession, code string, data entity.OtpData, userData entity.MemberData) error {
	arr := []string {data.OtpNumber}
	//測試環境 快速密碼
	if viper.GetString("ENV") != "prod" {
		arr = append(arr, "999999")
	}
	if tools.InArray(arr, code) {
		//更新最後登入時間
		if err := UpdateMemberDataLastTimeAndErrorZero(engine, &userData); err != nil {
			log.Error("update member data Error", err)
			return fmt.Errorf("系統錯誤")
		}
		data.OtpUse = 1
		if err := otp.UpdateOtpData(engine, data); err != nil {
			return fmt.Errorf("系統錯誤")
		}
		return nil
	} else {
		//寫入會員資料庫 登入失敗次數
		if err := UpdateMemberDataErrorFrequency(engine, &userData); err != nil {
			log.Error("update member data Frequency Error", err)
			return fmt.Errorf("系統錯誤")
		}
		return fmt.Errorf("OTP不正確")
	}
}
