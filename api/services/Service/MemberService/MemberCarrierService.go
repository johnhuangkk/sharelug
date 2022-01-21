package MemberService

import (
	"api/services/VO/Response"
	"api/services/dao/member"
	"api/services/database"
	"api/services/util/curl"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"github.com/spf13/viper"
	"net/url"
	"sort"
	"strings"
	"time"
)

//取出捐贈單位列表
func GetDonateUnitList(engine *database.MysqlSession) ([]Response.DonateListResponse, error) {
	var resp []Response.DonateListResponse
	data, err := member.GetDonateData(engine)
	if err != nil {
		return resp, err
	}
	for _, v := range data {
		res := Response.DonateListResponse {
			DonateName: v.DonateShort,
			DonateCode: v.DonateCode,
		}
		resp = append(resp, res)
	}
	res := Response.DonateListResponse {
		DonateName: "其它",
		DonateCode: "other",
	}
	resp = append(resp, res)
	return resp, nil
}
//驗證捐贈單位
func VerifyDonateCode(code string) error {
	config := viper.GetStringMapString("INVOICE.VC")
	PostValue := url.Values{}
	PostValue.Add("version", "1.0")
	PostValue.Add("action", "preserveCodeCheck")
	PostValue.Add("pCode", code)
	PostValue.Add("appId", config["appid"])
	PostValue.Add("TxID", time.Now().Format("20060102150405"))
	if err := postVerify(config["url"], PostValue); err != nil {
		return fmt.Errorf("捐贈單位無效")
	}
	return nil
}
//驗證手機載具號碼
func VerifyMobileCode(code string) error {
	config := viper.GetStringMapString("INVOICE.VC")
	PostValue := url.Values{}
	PostValue.Add("version", "1.0")
	PostValue.Add("action", "bcv")
	PostValue.Add("barCode", code)
	PostValue.Add("appId", config["appid"])
	PostValue.Add("TxID", time.Now().Format("20060102150405"))
	log.Debug("PostValue", PostValue)
	if err := postVerify(config["url"], PostValue); err != nil {
		return fmt.Errorf("手機條碼載具無效")
	}
	return nil
}
//驗證是否為營業人
func VerifyCompanyBan(code string) error {
	config := viper.GetStringMapString("INVOICE.VC")
	log.Debug("config", config, time.Now().Unix() + 10, time.Now().Unix())
	PostValue := url.Values{}
	PostValue.Add("version", "1.0")
	PostValue.Add("serial", "0000000004")
	PostValue.Add("action", "qryRecvRout")
	PostValue.Add("ban", code)
	PostValue.Add("timeStamp", fmt.Sprintf("%d",time.Now().Unix() + 180))
	PostValue.Add("appId", config["appid"])
	signature := BindSignature(PostValue, config["apikey"])
	PostValue.Add("signature", signature)
	log.Debug("PostValue", PostValue)
	result, err := curl.Post(config["url"], PostValue.Encode())
	if err != nil {
		return err
	}
	res := Response.VerifyCompanyBan{}
	if err := tools.JsonDecode(result, &res); err != nil{
		log.Error("Json Decode Error", err)
		return err
	}
	if res.BanUnitTpStatus != "Y" {
		log.Error("Verify Company Ban Fail", res)
		return fmt.Errorf("驗證失敗")
	}
	return nil
}
//SEND API CURL
func postVerify(url string, data url.Values) error {
	result, err := curl.Post(url, data.Encode())
	if err != nil {
		return err
	}
	res := Response.VerifyDonateCodeResponse{}
	if err := tools.JsonDecode(result, &res); err != nil{
		log.Error("Json Decode Error", err)
		return err
	}
	if res.IsExist != "Y" {
		log.Error("Verify Donate Code Fail", res)
		return fmt.Errorf("驗證失敗")
	}
	return nil
}
//產生簽名
func BindSignature(value url.Values, apikey string) string {
	var data []string
	for k, v := range value {
		data = append(data, fmt.Sprintf("%s=%s", k, v[0]))
	}
	sort.Strings(data)
	log.Debug("ssss", strings.Join(data, "&"))
	return tools.SHA256ByKey(apikey, strings.Join(data, "&"))
}