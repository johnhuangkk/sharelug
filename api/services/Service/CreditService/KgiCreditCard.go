package CreditService

import (
	"api/services/entity"
	"api/services/util/curl"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"github.com/spf13/viper"
	"net/url"
	"strconv"
)

type SetConnectionConfig struct {
	ResponseLink string
	HostName     string
	Path         string
	MerchantId   string
	TerminalId   string
	Type         string
	OrderId     string
	SellerId    string
	BuyerId     string
	CardId      string
	CardNumber  string
	ExpireDate  string
	Security    string
	Amount      int64
}

func New(params entity.AuthRequest) *SetConnectionConfig {
	auth := &SetConnectionConfig{}
	auth.MerchantId = params.MerchantId
	auth.TerminalId = params.TerminalId
	auth.ResponseLink = fmt.Sprintf("https://www.checkne.com/v1/pay/credit/confirm/%s", params.Type)
	auth.HostName = viper.GetString("KgiCredit.hostname")
	auth.Path = viper.GetString("KgiCredit.AuthPath")
	auth.OrderId = params.OrderId
	auth.SellerId = params.SellerId
	auth.BuyerId = params.BuyerId
	auth.CardId = params.CardId
	auth.CardNumber = params.CardNumber
	auth.ExpireDate = params.ExpireDate
	auth.Security = params.Security
	auth.Amount = params.Amount
	return auth
}

//授權交易
func (auth SetConnectionConfig) DoAuth() (entity.AuthResult, error) {
	var result entity.AuthResult
	var err error
	PostValue := url.Values{}
	PostValue.Add("MerchantID", auth.MerchantId)                    //商店代號
	PostValue.Add("TerminalID", auth.TerminalId)                    //端末機代號
	PostValue.Add("OrderID", auth.OrderId)                        //訂單編號
	PostValue.Add("CardNumber", auth.CardNumber)                  //交易卡號
	PostValue.Add("ExpireDate", auth.ExpireDate)                  //卡片到期日
	PostValue.Add("Cvv", auth.Security)                           //cvv
	PostValue.Add("TransCode", "00")                                //交易代碼
	PostValue.Add("TransMode", "0")                                 //交易類別
	PostValue.Add("Install", "0")                                   //分期期數
	PostValue.Add("TransAmt", strconv.FormatInt(auth.Amount, 10)) //交易金額
	PostValue.Add("NotifyURL", auth.ResponseLink)                      //回應網址
	PostValue.Add("HostName", auth.HostName)
	PostValue.Add("Path", auth.Path)
	log.Info("Do Auth", PostValue)
	body, err := curl.PostValues("http://172.27.0.6:8000/auth", PostValue)
	if err != nil {
		log.Error("curl post auth error", err)
		return result, err
	}
	if err := tools.JsonDecode(body, &result); err != nil {
		log.Error("Json Decode Error", err)
	}
	log.Info("Do Auth Result", result)
	return result, nil
}

//取消交易
func (auth SetConnectionConfig) DoVoid(params entity.CancelRequest) (entity.AuthResult, error) {
	PostValue := url.Values{}
	PostValue.Add("MerchantID", params.MerchantId) //商店代號
	PostValue.Add("OrderID", params.OrderId)     //訂單編號
	PostValue.Add("TransCode", "01")             //交易代碼
	PostValue.Add("HostName", auth.HostName)
	PostValue.Add("Path", auth.Path)
	log.Debug("post data", PostValue.Encode())
	//
	var result entity.AuthResult
	body, err := curl.PostValues("http://172.27.0.6:8000/auth", PostValue)
	if err != nil {
		log.Error("curl post xml error", err)
		return result, err
	}
	log.Debug("result", body)
	if err := tools.JsonDecode(body, &result); err != nil {
		log.Error("Json Decode Error", err)
	}
	return result, nil
}

//查詢交易
func (auth SetConnectionConfig) DoQuery(params entity.QueryRequest) (entity.QueryResponse, error) {
	PostValue := url.Values{}
	PostValue.Add("MerchantID", params.MerchantId) //商店代號
	PostValue.Add("OrderID", params.OrderId)     //訂單編號
	PostValue.Add("HostName", auth.HostName)
	PostValue.Add("Path", auth.Path)
	//
	log.Debug("post data", PostValue.Encode())
	var result entity.QueryResponse
	body, err := curl.PostValues("http://172.27.0.6:8000/authQuery", PostValue)
	if err != nil {
		log.Error("curl post xml error", err)
		return result, err
	}
	log.Debug("result", body)
	if err := tools.JsonDecode(body, &result); err != nil {
		log.Error("Json Decode Error", err)
	}
	return result, nil
}

