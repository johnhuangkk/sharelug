package Edm

import (
	"api/services/Service/Product"
	"api/services/Service/Sms"
	"api/services/VO/Request"
	"api/services/dao/Orders"
	"api/services/dao/member"
	"api/services/database"
	"api/services/util/log"
	"api/services/util/qrcode"
	"fmt"
	"strings"
)

func SendEdmAllMember(rule, message string, test bool) (int, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var list []string
	//取出所有會員
	if test {
		users := []string{"0955595590", "0939876808", "0986927039", "0930099161", "0973188223", "0960764133", "0916853214", "0972209578", "0916989876"}
		for _, v := range users {
			phone := strings.Replace(v, "0", "886", 1)
			list = append(list, phone)
		}
	} else {
		users, err := member.GetMemberDataBySubscribe(engine)
		if err != nil {
			return 0, fmt.Errorf("1001001")
		}
		if len(message) == 0 {
			return 0, fmt.Errorf("1012003")
		}
		//SellerId := "U2021052717241740041"
		//排除特定購買過指定的賣家
		for _, v := range users {
			if len(rule) != 0 {
				if Orders.IsExistOrderDataByBuyerAndSeller(engine, v.Uid, rule) == 0 {
					phone := strings.Replace(v.Mphone, "0", "886", 1)
					list = append(list, phone)
				}
			} else {
				phone := strings.Replace(v.Mphone, "0", "886", 1)
				list = append(list, phone)
			}
		}
	}
	content := []byte(message)
	limit := 500
	start := 0
	for start < len(list) {
		end := start + limit
		if end > len(list) {
			end = len(list)
		}
		log.Debug("Send users", list[start:end], content)
		response, err := Sms.FetNetSendMultiSms(list[start:end], content)
		if err != nil {
			log.Error("Send SMS Error", err)
			//return err
		}
		log.Debug("sms response ResultCode", response)
		start += limit
	}
	return len(list), nil
}
//產生短網址
func HandleGenerateShortLink(params Request.GenerateShortRequest) (string, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if len(params.Uri) == 0 {
		return "", fmt.Errorf("1012004")
	}
	tiny, err := Product.GeneratorShortUrl(engine, params.Uri)
	if err != nil {
		return "", fmt.Errorf("1001001")
	}
	return qrcode.GetTinyUrl(tiny), nil
}