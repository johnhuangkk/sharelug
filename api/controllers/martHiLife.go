package controllers

import (
	"api/services/Service/FamilyApi"
	"api/services/Service/FamilyXml"
	"api/services/Service/HiLifeApi"
	"api/services/Service/HiLifeNotificationService"
	"api/services/Service/HiLifeXml"
	"api/services/Service/OKXml"
	"api/services/VO/HiLifeMart"
	"api/services/model"
	"api/services/util/log"
	"api/services/util/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

//Data=
//<?xml version='1.0' encoding='UTF-8'?>
//<Doc>
//	<ShipmentNos>
//		<ParentId>124</ParentId>
//		<EshopId>901</EshopId>
//		<OrderNo>20RKB40416476</OrderNo>
//		<EcOrderNo>TSB30001968</EcOrderNo>
//		<OriginStoreId>3750</OriginStoreId>
//		<StoreType>2</StoreType>
//		<ChkMac>0C65AA5CC37A1D415AACD3E4C5D6DB32</ChkMac>
//	</ShipmentNos>
//</Doc>

// 關轉店通知 HiLife -> Self
func HiLifeNotification(ctx *gin.Context) {

	var params HiLifeMart.HiLifParams

	rsp := response.New(ctx)

	if err := ctx.ShouldBind(&params); err != nil {
		log.Error("HiLifeNotification ShouldBind Error", err.Error())
	}

	log.Info("params [%v]", params)

	rsp.XML(HiLifeNotificationService.SwitchNotification(params))

}

func MartPrintOrder(ctx *gin.Context) {
	orderType, _ := ctx.GetQuery("t")
	shipNo, _ := ctx.GetQuery("s")

	log.Debug("MartPrintOrder:", orderType, shipNo)
	var err error
	var data []byte

	if orderType == "OK" {
		data, err = model.MartOkPrintShippingOrder(strings.Split(shipNo, ","))
	}

	if err != nil {
		ctx.String(200, fmt.Sprintf("Error:%v", err))
		return
	}

	ctx.Data(200, "application/pdf; charset=utf-8", data)
}

func MartSwitchOrder(ctx *gin.Context) {
	orderType := ctx.PostForm("ShipType")
	shipNo := ctx.PostForm("ShipNo")
	ecOrderNo := ctx.PostForm("EcOrderNo")
	newStoreId := ctx.PostForm("NewStoreId")
	isReceiveStore := ctx.PostForm("IsReceiveStore") == "true"

	log.Debug("MartSwitchOrder:", orderType, shipNo)
	var err error
	if orderType == "HiLife" {
		if err = HiLifeApi.MartHiLifeSwitchStore(shipNo, ecOrderNo, newStoreId, isReceiveStore); err != nil {
			log.Error("MartSwitchOrder:", err.Error())
		}
	}

	if orderType == "Family" {
		if err = FamilyApi.MartFamilySwitchStore(shipNo, ecOrderNo, newStoreId, isReceiveStore); err != nil {
			log.Error("MartSwitchOrder:", err.Error())
		}
	}

	if orderType == "OK" {
		//if err = OKApi.SwitchStore(shipNo,ecOrderNo,newStoreId,isReceiveStore);err != nil {
		//	log.Error("MartSwitchOrder:",err.Error())
		//}
	}

	if err != nil {
		ctx.String(200, fmt.Sprintf("Error:%v", err))
		return
	}
	ctx.Status(200)
}

func MartFetching(ctx *gin.Context) {
	shipType := ctx.PostForm("ShipType")
	subType := ctx.PostForm("SubType")

	if shipType == "HiLife" {
		if subType == "All" {
			go HiLifeXml.MartHiLifeFetchShipping()
		}
		if subType == "R27" {
			//go HiLifeXml.MartHiLifeFetchR27()
		}
		if subType == "R00" {
			go HiLifeXml.MartHiLifeFetchStoreList()
		}
	}

	if shipType == "Family" {
		if subType == "I00" {
			go FamilyXml.MartFamilyFetchStoreList()
		}
		if subType == "All" {
			go FamilyXml.MartFamilyFetchShipping()
		}
	}

	if shipType == "OK" {
		if subType == "F01" {
			go OKXml.MartOKFetchStoreList()
		}
		if subType == "All" {
			go OKXml.MartOkFetchShipping()
		}
	}

	ctx.String(200, "讚哦")
}
