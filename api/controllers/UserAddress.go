package controllers

import (
	"api/config/middleware"
	"api/services/Service/Soap"
	"api/services/Service/UserAddressService"
	"api/services/Task"
	"api/services/VO/Request"
	"api/services/VO/UserAddress"
	"api/services/errorMessage"
	"api/services/util/log"
	"api/services/util/response"
	"api/services/util/tools"
	"encoding/xml"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/smokezl/govalidators"
)

type soapRQ struct {
	XMLName   xml.Name `xml:"soapenv:Envelope"`
	XMLNsSoap string   `xml:"xmlns:soapenv,attr"`
	// XMLNsXSI  string   `xml:"xmlns:xsi,attr"`
	// XMLNsXSD  string   `xml:"xmlns:xsd,attr"`
	Body soap12Body
}

// <SayHello xmlns="http://learnwebservices.com/services/hello">
//          <HelloRequest>
//             <Name>John Doe</Name>
//          </HelloRequest>
//       </SayHello>
type payload struct {
	XMLName      xml.Name `xml:"SayHello"`
	XMLNs        string   `xml:"xmlns,attr"`
	HelloRequest string   `xml:"HelloRequest>Name"`
}

type soap12Body struct {
	XMLName xml.Name `xml:"soapenv:Body"`
	Payload payload
}

// 新增地址
func AddAddressAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := UserAddress.AddressInfo{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Debug("AddAddress => ", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	validator := govalidators.New()
	if err := validator.LazyValidate(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1002104)).Send()
		return
	}
	UserData := middleware.GetUserData(ctx)
	if err := UserAddressService.HandleAddress(UserData, params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1002101)).Send()
		return
	}

	resp.Success("新增成功").SetData(true).Send()
}

// 刪除地址
func DeleteAddressAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := UserAddress.DeleteAddress{}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Debug("DeleteAddress => ", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	userData := middleware.GetUserData(ctx)
	if err := UserAddressService.DeleteAddress(userData, params); err != nil {
		log.Debug("DeleteAddress", err)
		resp.Fail(errorMessage.GetMessageByCode(1002102)).Send()
		return
	}
	resp.Success("刪除成功").SetData(true).Send()
}

// 確認寄送狀態是否有寄送地址
func CheckShipSendAddressExistAction(ctx *gin.Context) {
	resp := response.New(ctx)
	ship := ctx.Param("ship")
	if len(ship) == 0 {
		resp.Fail(errorMessage.GetMessageByCode(1002102)).Send()
		return
	}
	UserData := middleware.GetUserData(ctx)
	boolean := UserAddressService.InputCheckShipSendAddressExist(ship, UserData.Uid)
	if boolean {
		resp.Success("成功").SetData(true).Send()
	} else {
		resp.Fail(errorMessage.GetMessageByCode(1002106)).Send()
	}
}

// 取得收件地址
func GetReceiveAddressAction(ctx *gin.Context) {
	resp := response.New(ctx)
	ship := ctx.Param("ship")
	UserData := middleware.GetUserData(ctx)
	if len(UserData.Uid) == 0 {
		resp.Fail(1001001, "尚未登入").Send()
		return
	}
	//類型 寄件 S 收件 R
	addresses, err := UserAddressService.GetAddresses("R", ship, UserData)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1002103)).Send()
		return
	}
	resp.Success("成功").SetData(addresses).Send()
}

// 取得寄送地址
func GetSendAddressAction(ctx *gin.Context) {
	resp := response.New(ctx)
	//類型 寄件 S 收件 R
	UserData := middleware.GetUserData(ctx)
	addresses, err := UserAddressService.GetAddresses("S", "", UserData)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1002103)).Send()
		return
	}
	resp.Success("成功").SetData(addresses).Send()
}

func Kms(ctx *gin.Context) {
	S, err := tools.AwsKMSEncrypt("1001-2000-3000-4000")
	log.Info("AwsKMSEncrypt [%s]", S)
	if err != nil {
		fmt.Println("Got error encrypting data: ", err)
		os.Exit(1)
	}
	SS, err := tools.AwsKMSDecrypt(S)
	log.Info("AwsKMSDecrypt [%s]", SS)
}

func TriggerOk(ctx *gin.Context) {
	Task.CvsAccountingTask()
}

func TranslateAddress(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.TranslateAddrEn{}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error("Translate addr params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	addr, err := Soap.TranslateAddress(params.Addr)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	resp.Success("OK").SetData(addr).Send()
}
