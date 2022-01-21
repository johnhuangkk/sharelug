package model

import (
	"api/services/Enum"
	"api/services/Service/Invoice"
	"api/services/VO/InvoiceVo"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/InvoiceDao"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"github.com/spf13/viper"
)

func HandleGetInvoiceList(userData entity.MemberData, params Request.InvoiceListRequest) (Response.InvoiceListResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.InvoiceListResponse
	count, err := InvoiceDao.CountInvoiceListByBuyerId(engine, userData.Uid)
	if err != nil {
		log.Error("Count Invoice List Error", err)
		return resp, fmt.Errorf("系統錯誤")
	}
	data, err := InvoiceDao.GetInvoiceListByBuyerId(engine, userData.Uid, int(params.Limit), int(params.Start))
	if err != nil {
		log.Error("Get Invoice List Error", err)
		return resp, fmt.Errorf("系統錯誤")
	}
	for _, v := range data {
		var res Response.InvoiceList
		res.OrderId = v.OrderId
		res.InvoiceNumber = fmt.Sprintf("%s-%s", v.InvoiceTrack, v.InvoiceNumber)
		res.InvoiceStatus = v.InvoiceStatus
		if v.DonateMark == 1 {
			res.InvoiceStatus = Enum.InvoiceTypeDonate
			res.InvoiceNumber = tools.MaskerInvoice(res.InvoiceNumber)
			res.InvoiceStatusText = Invoice.GetDonateInfo(engine, v.DonateBan)
		}
		log.Debug("Identifier:", v.Identifier != "0000000000")
		if v.Identifier != "0000000000" {
			res.InvoiceStatusText = fmt.Sprintf("已輸入統編(%s)", v.Identifier)
		}

		//如果已列印發票
		if v.PrintMark == "Y" {
			res.CarrierType = Enum.InvoiceCarrierTypePrint
		} else {
			res.CarrierType = v.CarrierType
		}
		if v.CarrierType == Enum.InvoiceCarrierTypeMobile {
			res.CarrierTypeText = fmt.Sprintf("(存手機條碼載具%s)", v.Carrier)
		}
		res.CreateTime = v.CreateTime.Format("2006/01/02 15:04")
		resp.InvoiceList = append(resp.InvoiceList, res)
	}
	resp.InvoiceCount = count
	return resp, nil
}
func HandleGetInvoiceDetail(userData entity.MemberData, params InvoiceVo.InvoiceRequest) (InvoiceVo.InvoiceDetailVo, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp InvoiceVo.InvoiceDetailVo
	data, err := InvoiceDao.GetInvoiceByOrderId(engine, params.OrderId)
	if err != nil {
		log.Error("Get Invoice Data Error", err)
		return resp, fmt.Errorf("系統錯誤")
	}
	if data.BuyerId != userData.Uid {
		return resp, fmt.Errorf("系統錯誤")
	}
	config := viper.GetStringMapString("PLATFORM")
	resp.InvoiceNumber = fmt.Sprintf("%s-%s", data.InvoiceTrack, data.InvoiceNumber)
	if data.Identifier != "0000000000" {
		if data.CarrierType == Enum.InvoiceCarrierTypeMember {
			resp.InvoiceStatusText = "已輸入統編 || 會員載具"
		}
		if data.CarrierType == Enum.InvoiceCarrierTypeMobile {
			resp.InvoiceStatusText = "已輸入統編 || 手機載具"
		}
		if data.CarrierType == Enum.InvoiceCarrierTypeCert {
			resp.InvoiceStatusText = "已輸入統編 || 憑證載具"
		}
		if data.PrintMark == "Y" {
			resp.InvoiceStatusText = "已輸入統編 || 已索取紙本電子發票"
		}
	} else {
		if data.CarrierType == Enum.InvoiceCarrierTypeMember {
			resp.InvoiceStatusText = "會員載具"
		}
		if data.CarrierType == Enum.InvoiceCarrierTypeMobile {
			resp.InvoiceStatusText = "手機載具"
		}
		if data.CarrierType == Enum.InvoiceCarrierTypeCert {
			resp.InvoiceStatusText = "憑證載具"
		}
	}
	if data.DonateMark == 1 {
		resp.InvoiceNumber = tools.MaskerInvoice(resp.InvoiceNumber)
		resp.InvoiceStatus = Enum.InvoiceTypeDonate
	} else {
		resp.InvoiceStatus = data.InvoiceStatus
	}
	if data.InvoiceStatus == Enum.InvoiceStatusWin {
		if data.AwardModel == "Y" {
			resp.InvoiceStatusText = "已索取紙本電子發票"
		}
		if data.AwardModel == "Z" {
			resp.InvoiceStatusText = "已於整合服務平台設定中獎獎金自動匯入帳戶功能"
		}
		if data.AwardModel == "X" {
			resp.InvoiceStatusText = "歸戶至共通性載具可至超商多媒體服務機KIOS以手機條碼或自然人憑證列印"
		}
		if data.AwardModel == "A" {
			resp.InvoiceStatusText = fmt.Sprintf("將於 %v/10 前寄送", tools.StringToInt64(data.Month) + 2)
		}
	}
	resp.InvoiceYear = data.Year
	resp.InvoiceMonth = data.Month
	resp.CreateTime = data.CreateTime.Format("2006/01/02 15:04:05")
	resp.InvoiceRandom = data.RandomNumber
	resp.SellerBan = config["company_ban"]
	if data.Identifier == "0000000000" {
		resp.BuyerBan = ""
		resp.BuyerName = ""
	} else {
		resp.BuyerBan = data.Identifier
		resp.BuyerName = data.Buyer
	}
	var items []InvoiceVo.ItemList
	_ = tools.JsonDecode([]byte(data.Detail), &items)
	resp.ItemList = items
	resp.Sales = data.Sales
	resp.Tax = data.Tax
	resp.Amount = data.Amount

	InvVo := resp.GetQRCodeInvVo()
	code, err := Invoice.QRCodeINV(InvVo)
	if err != nil {
		log.Error("Get QRCodeINV Error", err)
		return resp, fmt.Errorf("系統錯誤")
	}
	resp.QRCode1 = code
	resp.QRCode2 = Invoice.QRCodeProduct(InvVo)
	log.Debug("HandleGetInvoiceDetail", resp)
	return resp, nil
}

func HandleGetAwardedInvoice(orderId string) (InvoiceVo.InvoiceDetailVo, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp InvoiceVo.InvoiceDetailVo
	data, err := InvoiceDao.GetInvoiceByOrderId(engine, orderId)
	if err != nil {
		log.Error("Get Invoice Data Error", err)
		return resp, fmt.Errorf("系統錯誤")
	}
	config := viper.GetStringMapString("PLATFORM")
	resp.InvoiceNumber = fmt.Sprintf("%s-%s", data.InvoiceTrack, data.InvoiceNumber)
	if data.Identifier != "0000000000" {
		if data.CarrierType == Enum.InvoiceCarrierTypeMember {
			resp.InvoiceStatusText = "已輸入統編 || 會員載具"
		}
		if data.CarrierType == Enum.InvoiceCarrierTypeMobile {
			resp.InvoiceStatusText = "已輸入統編 || 手機載具"
		}
		if data.CarrierType == Enum.InvoiceCarrierTypeCert {
			resp.InvoiceStatusText = "已輸入統編 || 憑證載具"
		}
		if data.PrintMark == "Y" {
			resp.InvoiceStatusText = "已輸入統編 || 已索取紙本電子發票"
		}
	} else {
		if data.CarrierType == Enum.InvoiceCarrierTypeMember {
			resp.InvoiceStatusText = "會員載具"
		}
		if data.CarrierType == Enum.InvoiceCarrierTypeMobile {
			resp.InvoiceStatusText = "手機載具"
		}
		if data.CarrierType == Enum.InvoiceCarrierTypeCert {
			resp.InvoiceStatusText = "憑證載具"
		}
	}
	if data.DonateMark == 1 {
		resp.InvoiceNumber = tools.MaskerInvoice(resp.InvoiceNumber)
		resp.InvoiceStatus = Enum.InvoiceTypeDonate
	} else {
		resp.InvoiceStatus = data.InvoiceStatus
	}
	if data.InvoiceStatus == Enum.InvoiceStatusWin {
		if data.AwardModel == "Y" {
			resp.InvoiceStatusText = "已索取紙本電子發票"
		}
		if data.AwardModel == "Z" {
			resp.InvoiceStatusText = "已於整合服務平台設定中獎獎金自動匯入帳戶功能"
		}
		if data.AwardModel == "X" {
			resp.InvoiceStatusText = "歸戶至共通性載具可至超商多媒體服務機KIOS以手機條碼或自然人憑證列印"
		}
		if data.AwardModel == "A" {
			resp.InvoiceStatusText = fmt.Sprintf("將於 %v/10 前寄送", tools.StringToInt64(data.Month) + 2)
		}
	}
	resp.InvoiceYear = data.Year
	resp.InvoiceMonth = data.Month
	resp.CreateTime = data.CreateTime.Format("2006/01/02 15:04:05")
	resp.InvoiceRandom = data.RandomNumber
	resp.SellerBan = config["company_ban"]
	if data.Identifier == "0000000000" {
		resp.BuyerBan = ""
		resp.BuyerName = ""
	} else {
		resp.BuyerBan = data.Identifier
		resp.BuyerName = data.Buyer
	}
	var items []InvoiceVo.ItemList
	_ = tools.JsonDecode([]byte(data.Detail), &items)
	resp.ItemList = items
	resp.Sales = data.Sales
	resp.Tax = data.Tax
	resp.Amount = data.Amount

	InvVo := resp.GetQRCodeInvVo()
	code, err := Invoice.QRCodeINV(InvVo)
	if err != nil {
		log.Error("Get QRCodeINV Error", err)
		return resp, fmt.Errorf("系統錯誤")
	}
	resp.QRCode1 = code
	resp.QRCode2 = Invoice.QRCodeProduct(InvVo)
	log.Debug("HandleGetInvoiceDetail", resp)
	return resp, nil
}
