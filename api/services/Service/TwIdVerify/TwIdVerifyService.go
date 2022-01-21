package TwIdVerify

import (
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/util/curl"
	"api/services/util/log"
	"api/services/util/tools"
	"crypto/sha1"
	"fmt"
	"github.com/spf13/viper"
)

// 換領補驗證
func SendIdentityVerify(params *Request.TWIDParams, userId, applyTime string) (Response.TWIDVerify, error) {
	var resp Response.TWIDVerify
	tokenHeader, err := GeneratorTokenHeader(params, applyTime, userId)
	if err != nil {
		log.Error("Generator Token Header error", err)
		return resp, fmt.Errorf("系統錯誤")
	}
	log.Debug("Generator Token Info", tokenHeader)
	twIdUrl := viper.GetString("MOI.SERVER_URL")
	resp, err = curl.GetIDCheck(twIdUrl, tokenHeader)
	if err != nil {
		log.Error("curl post xml error", err)
		return resp, fmt.Errorf("系統錯誤")
	}
	return resp, nil
}

/**
 * 產生 Header Token
 */
func GeneratorTokenHeader(params *Request.TWIDParams, applyTime string, userId string) (string, error) {
	h := sha1.New()
	h.Write([]byte(params.IdentityId))
	condition := tools.Condition{
		PersonId: params.IdentityId,
		IdMark: params.IssueType,
		IdMarkDate: tools.StringPadLeft(applyTime, 7),
		IssueAreaCode: params.IssueCounties,
	}
	tokenHeader, err := tools.GeneratorIDCheckJWT(userId, condition)
	if err != nil {
		log.Debug("Generator Token Header error", err)
		return "", err
	}
	tokenHeader = "Bearer " + tokenHeader
	return tokenHeader, nil
}
