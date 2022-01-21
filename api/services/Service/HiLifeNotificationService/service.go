package HiLifeNotificationService

import (
	"api/services/Enum"
	"api/services/Service/CvsShipping"
	"api/services/Service/MartHiLife"
	"api/services/VO/HiLifeMart"
	"api/services/database"
	"api/services/util/tools"
	"encoding/xml"
)

// 處理閉轉資訊
func handleSwitchShipmentNos(engine *database.MysqlSession, body HiLifeMart.SwitchBody) (code, message string) {
	err := MartHiLife.SwitchCheckSum(body)
	if err != nil {
		return `999`, `驗證碼檢核失敗`
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = `Switch`
	UpdateCvsShipping.ShipType = Enum.CVS_HI_LIFE
	UpdateCvsShipping.ShipNo = body.OrderNo
	UpdateCvsShipping.DateTime = tools.Now(`YmdHis`)
	UpdateCvsShipping.FlowType = `N`
	UpdateCvsShipping.Log = tools.XmlToString(body)

	if body.StoreType == `1` {
		UpdateCvsShipping.FlowType = `R`
	}

	// 更新托運閉轉資訊
	if err = UpdateCvsShipping.UpdateCvsShippingSwitch(engine); err != nil {
		return `999`, err.Error()
	}

	return `000`, `成功`
}

// 閉轉通知
func SwitchNotification(params HiLifeMart.HiLifParams) HiLifeMart.Response {

	var xmlRsp HiLifeMart.Response
	var xmlRspBody HiLifeMart.ResponseBody

	var _xml HiLifeMart.Switch

	_ = xml.Unmarshal([]byte(params.Data), &_xml)

	engine := database.GetMysqlEngine()
	defer engine.Close()

	for _, s := range _xml.ShipmentNos {
		code, message := handleSwitchShipmentNos(engine, s)


		xmlRspBody.RspSwitch(s, code, message)
		xmlRsp.ShipmentNos = append(xmlRsp.ShipmentNos, xmlRspBody)
	}

	return xmlRsp
}