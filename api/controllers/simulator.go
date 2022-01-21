package controllers

import (
	"api/services/Service/SendMail"
	"api/services/Service/Sms"
	"api/services/errorMessage"
	"api/services/util/curl"
	"api/services/util/log"
	"api/services/util/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
)

func GetSmsAction(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "simulator/sms.html", gin.H{
		"title": "簡訊測試",
	})
}

type SmsParams struct {
	Telecom string `form:"telecom"`
	Phone   string `form:"phone"`
	Content string `form:"content"`
}

type MailParams struct {
	Type       string 	`json:"Type"`
	Username   string	`json:"Username"`
	Subject    string	`json:"Subject"`
	Title      string	`json:"Title"`
	Content    string	`json:"Content"`
	To 		   string   `json:"To"`	
}

type TransferParams struct {
	TransferAcc string `json:"TransferAcc"`
}

func PostSmsAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &SmsParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Debug("new => ", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}

	if params.Phone == "" && params.Content == "" {
		resp.Fail(1001001, "請輸入內容").Send()
		return
	}
	params.Phone = strings.Replace(params.Phone, "0", "886", 1)
	log.Debug("Post Sms ", params.Phone, params.Content)
	var message interface{}
	switch params.Telecom {
	case "fetNet":
		data, err := Sms.FetNetSendSms(params.Phone, []byte(params.Content))
		if err != nil {
			resp.Fail(1001001, err.Error()).Send()
			return
		}
		message = data
	case "miTake":
		data, err := Sms.MiTakeSmsSend(params.Phone, []byte(params.Content))
		if err != nil {
			resp.Fail(1001001, err.Error()).Send()
			return
		}
		message = data
	}
	resp.Success("發送完成").SetData(message).Send()
}

func PostMailAction(ctx *gin.Context)  {
	resp := response.New(ctx)
	params := &MailParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	mail := SendMail.SetMassage{
		Username:   "Test Mail",
		Subject:    params.Subject,
		Title:      params.Title,
		Content:    params.Content,
	}
	if params.Type == "AWS" {
		if err := mail.SendAwsMail(params.To); err != nil {
			log.Error("send mail Error", err)
			resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
			return
		}
	} else {
		if err := mail.SendGmail(params.To); err != nil {
			log.Error("send mail Error", err)
			resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
			return
		}
	}
	resp.Success("OK").SetData(true).Send()
}

func PostTransferAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := TransferParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}

	PostValue := url.Values{}
	//PostValue.Add("RESEND", "I")
	//PostValue.Add("HEAD", "11")
	//PostValue.Add("TDATE", "20210728")
	//PostValue.Add("TTIME", "194629")
	//PostValue.Add("DATE","20210728")
	//PostValue.Add("ACCNO", "8000100121112104")
	//PostValue.Add("DEPTYPE", "2")
	//PostValue.Add("CURRENCY", "TWD")
	//PostValue.Add("SIGN", "+")
	//PostValue.Add("AMT", "1260")
	//PostValue.Add("TYPE","T")
	//PostValue.Add("RACCNO", "0120000610168131456")
	//PostValue.Add("SwiftCode","")
	//PostValue.Add("NOTE", "")
	//PostValue.Add("ATYPE", "")
	//PostValue.Add("IDNO", "")
	//PostValue.Add("BACCNO", "0060070100001513")
	//PostValue.Add("SEQNO", "000939617")
	//PostValue.Add("EOR", "0a")
	res, err := curl.PostJson("http://local.api.sharelug.com/gw/transfer/notify", PostValue.Encode())
	if err != nil {
		log.Error("curl post Transfer Error", err)
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	log.Debug("ssss", res)
	resp.Success("OK").SetData(true).Send()
}

