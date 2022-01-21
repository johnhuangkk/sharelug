package model

import (
	"api/services/Enum"
	"api/services/Service/MemberService"
	"api/services/VO/InvoiceVo"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/database"
	"api/services/entity"
	"api/services/util/curl"
	"api/services/util/log"
	"api/services/util/session"
	"api/services/util/tools"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

//取出發票載具資料
func HandleGetCarrier(userData entity.MemberData) (Response.CarrierResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.CarrierResponse
	unit, err := MemberService.GetDonateUnitList(engine)
	if err != nil {
		log.Error("Get Donate Unit Data Error", err)
		return resp, err
	}
	resp.DonateUnit = unit
	if userData.Category != Enum.CategoryCompany {
		data, err := MemberService.GetMemberCarrierByMemberId(engine, userData.Uid)
		if err != nil {
			log.Error("Get Member Carrier Data Error", err)
			return resp, err
		}
		if data.InvoiceType == Enum.InvoiceTypeDonate {
			for _, v := range unit {
				if v.DonateCode == data.DonateBan {
					resp.Choose = data.DonateBan
				}
			}
		}
		if  len(resp.Choose) == 0 && len(data.DonateBan) != 0 {
			resp.Choose = "other"
		}
		resp.InvoiceType = data.InvoiceType
		resp.CompanyBan = data.CompanyBan
		resp.CompanyName = data.CompanyName
		resp.DonateBan = data.DonateBan
		resp.CarrierType = data.CarrierType
		resp.CarrierId = data.CarrierId
	} else {
		Key := viper.GetString("EncryptKey")
		resp.InvoiceType = Enum.InvoiceTypeCompany
		resp.CompanyBan = tools.AesDecrypt(userData.Identity, Key)
		resp.CompanyName = userData.CompanyName
		resp.DonateBan = ""
		resp.CarrierType = Enum.InvoiceCarrierTypeMember
		resp.CarrierId = userData.InvoiceCarrier
	}
	return resp, nil
}

//寫入發票載具資料
func HandlePostCarrier(userData entity.MemberData, params Request.CarrierRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if userData.Category == Enum.CategoryCompany {
		return fmt.Errorf("公司戶不得變更")
	}
	if err := MemberService.ChangePostCarrier(engine, params, userData); err != nil {
		return err
	}
	return nil
}

//發票大平台綁定傳送綁定
func HandleBindPlatform(params Request.BindPlatformRequest) (string, error) {
	data := InvoiceVo.BindCarrierVo{
		Token: params.Token,
		Ban:   params.Ban,
	}
	uud := uuid.New()
	str, _ := tools.JsonEncode(data)
	cookie := uud.String()
	log.Debug("token ", str)
	if err := session.Bind(cookie).Put("bind", str); err != nil {
		log.Error("Session Set Old Error", err)
		return cookie, fmt.Errorf("系統錯誤")
	}
	return cookie, nil
}

//發票大平台綁定驗證綁定資料
func HandleVerifyBindCarrier(userData entity.MemberData, params Request.BindCarrierRequest) (Response.BindCarrierResponse, error) {
	config := viper.GetStringMapString("INVOICE.API")
	log.Debug("config", config)
	var resp Response.BindCarrierResponse
	//判斷是否登入
	if len(userData.Uid) == 0 {
		return resp, fmt.Errorf("尚未登入")
	}
	//取出Redis的資料
	redis := session.Bind(params.Token).Get("bind")
	var token InvoiceVo.BindCarrierVo
	if err := tools.JsonDecode([]byte(redis.(string)), &token); err != nil {
		return resp, fmt.Errorf("驗證失敗！")
	}
	//驗證資料
	if len(token.Token) == 0 || len(token.Ban) == 0 {
		return resp, fmt.Errorf("驗證失敗！")
	}
	//驗證是否為大平台請求
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())
	vo := InvoiceVo.VerifyBindCarrierRequest{
		Token: token.Token,
		Nonce: nonce,
	}
	result, err := curl.PostJson(config["apifrom"], vo)
	if err != nil {
		return resp, fmt.Errorf("驗證失敗！")
	}
	var body InvoiceVo.VerifyBindCarrierResponse
	if err := tools.JsonDecode(result, &body); err != nil {
		return resp, fmt.Errorf("驗證失敗！")
	}
	//回傳結果是否為Y
	if body.TokenFlag != "Y" || body.Nonce != nonce {
		return resp, fmt.Errorf("驗證失敗！")
	}

	signature := bindDataSignature(config["apikey"], config["apicardban"], userData.InvoiceCarrier, userData.InvoiceCarrier, config["apicardtype"], token.Token)
	log.Debug("bindDataSignature", signature)
	resp.CardBan = config["apicardban"]
	resp.CardNo1 = tools.Base64EncodeByString(userData.InvoiceCarrier)
	resp.CardNo2 = tools.Base64EncodeByString(userData.InvoiceCarrier)
	resp.CardType = tools.Base64EncodeByString(config["apicardtype"])
	resp.Token = token.Token
	resp.Signature = signature
	resp.Action = config["apifrom"]
	return resp, nil
}

//產生簽名
func bindDataSignature(key, cardBan, cardNo1, cardNo2, cardType, token string) string {
	PostValue := url.Values{}
	PostValue.Add("card_ban", cardBan)
	PostValue.Add("card_no1", cardNo1)
	PostValue.Add("card_no2", cardNo2)
	PostValue.Add("card_type", cardType)
	PostValue.Add("token", token)
	return MemberService.BindSignature(PostValue, key)
}
