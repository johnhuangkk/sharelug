package Sms

import (
	"api/services/Enum"
	"api/services/dao/Sms"
	"api/services/database"
	"api/services/entity"
	"api/services/util/curl"
	"api/services/util/log"
	"api/services/util/xml"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/spf13/viper"
	"net/url"
	"strings"
)

// 取得 到店貨態 後發送簡訊  目前只有OK超商
func CvsPushShipStatusShopMessageSms(engine *database.MysqlSession, orderData entity.OrderData)  {
	var template string
	var err error
	var storeName string

	switch orderData.ShipType {
	case Enum.CVS_OK_MART:
		_, err = engine.Engine.Table(entity.MartOkStoreData{}).Select("store_address").Where("store_id = ?", orderData.ReceiverAddress).Get(&storeName)
	default:
		return
	}

	if err != nil {
		log.Error(`getOKConsignmentInfo Error`, err.Error())
		return
	}

	template = fmt.Sprintf(
		`商品到店通知：%s 已送達%s%s，請攜帶證件前往取貨，謝謝。`,
		orderData.ShipNumber,
		Enum.Shipping[orderData.ShipType],
		storeName)


	log.Info(`CvsPushShipStatusShopMessageSms Info [%v]`, template)

	_ = PushMessageSms(orderData.ReceiverPhone, template)
}


func PushMessageSms(phone string, message string) error {
	content := []byte(message)
	phone = strings.Replace(phone, "0", "886", 1)
	ENV := viper.GetString("ENV")
	var response entity.SmsResult
	if ENV == "prod" {
		if true {
			response, _ = FetNetSendSms(phone, content)
			if response.ResultCode != "00000" {
				response, _ = MiTakeSmsSend(phone, content)
			}
			log.Debug("sms response ResultCode", response.ResultCode, response)
		}
	} else {
		_, err := Sms.InsertSmsData(phone, string(content), "FetNet")
		if err != nil {
			log.Error("Sms insert Log Error", err)
			return err
		}
	}
	return nil
}

//發送 遠傳 簡訊
func FetNetSendSms(Number string, Body []byte) (entity.SmsResult, error){
	result := entity.SmsResult{}
	data, err := Sms.InsertSmsData(Number, string(Body), "FetNet")
	if err != nil {
		log.Error("Sms insert Log Error", err)
		return result, err
	}
	//轉Base64 String
	content := base64.StdEncoding.EncodeToString(Body)
	v := entity.SmsSubmitReq{
		SysId:         "",
		SrcAddress:    "",
		DestAddress:   Number,
		SmsBody:       content,
		DrFlag:        true,
		FirstFailFlag: false,
	}

	v.SysId = viper.GetString("sms.sysid")
	v.SrcAddress = viper.GetString("sms.source")
	hostname := viper.GetString("sms.hostname")

	byt, err := xml.SmsXmlEncoder(v)
	if err != nil {
		log.Error("Xml Encoder Error", err)
		return result, err
	}
	contents := bytes.NewBuffer(byt)
	body, err := curl.PostXml(hostname, contents.String())
	if err != nil {
		log.Error("curl post xml error", err)
		return result, err
	}

	response, err := xml.SmsXmlDecoder(string(body))
	if err != nil {
		log.Debug("Xml Decoder Error", err)
		return result, err
	}
	if err := Sms.UpdateSmsData(data, response.ResultCode, response.ResultText); err != nil {
		log.Error("Sms Update Log Error", err)
		return result, err
	}
	return response, nil
}

//發送 三竹 簡訊
func MiTakeSmsSend(Number string, Body []byte) (entity.SmsResult, error) {

	result := entity.SmsResult{}
	data, err := Sms.InsertSmsData(Number, string(Body),"MiTake")
	if err != nil {
		log.Error("Sms insert Log Error", err)
		return result, err
	}
	username := viper.GetString("mitake.username")
	password := viper.GetString("mitake.password")
	hostname := viper.GetString("mitake.hostname")

	PostValue := url.Values{}
	PostValue.Add("username", username)
	PostValue.Add("password", password)
	PostValue.Add("dstaddr", Number)
	PostValue.Add("smbody", string(Body))

	body, err := curl.Post(hostname, PostValue.Encode())
	if err != nil {
		log.Error("curl post xml error", err)
		return result, err
	}
	result.ResultText = string(body)
	if err := Sms.UpdateSmsData(data, "", result.ResultText); err != nil {
		log.Error("Sms Update Log Error", err)
		return result, err
	}
	return result, nil
}

//發送多個手機簡訊
func FetNetSendMultiSms(Number []string, Body []byte) (entity.SmsResult, error){
	result := entity.SmsResult{}
	var smsData []entity.SmsLogData
	for _, v := range Number {
		data, err := Sms.InsertSmsData(v, string(Body), "FetNet")
		if err != nil {
			log.Error("Sms insert Log Error", err)
			return result, err
		}
		smsData = append(smsData, data)
	}
	//轉Base64 String
	content := base64.StdEncoding.EncodeToString(Body)
	var message entity.SmsMultiSubmitReq
	message.SysId = viper.GetString("sms.sysid")
	message.SrcAddress = viper.GetString("sms.source")
	message.DestAddress = Number
	message.SmsBody = content
	message.DrFlag = true
	message.FirstFailFlag = false
	hostname := viper.GetString("sms.hostname")

	byt, err := xml.SmsXmlEncoder(message)
	if err != nil {
		log.Error("Xml Encoder Error", err)
		return result, err
	}
	contents := bytes.NewBuffer(byt)
	log.Debug("contents", contents, byt)
	body, err := curl.PostXml(hostname, contents.String())
	if err != nil {
		log.Error("curl post xml error", err)
		return result, err
	}
	response, err := xml.SmsXmlDecoder(string(body))
	if err != nil {
		log.Debug("Xml Decoder Error", err)
		return result, err
	}
	for _, v := range smsData {
		if err := Sms.UpdateSmsData(v, response.ResultCode, response.ResultText); err != nil {
			log.Error("Sms Update Log Error", err)
			return result, err
		}
	}
	return response, nil
}
