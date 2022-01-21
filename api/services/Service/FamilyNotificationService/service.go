package FamilyNotificationService

import (
	"api/services/Enum"
	"api/services/Service/CvsShipping"
	"api/services/VO/FamilyMart"
	"api/services/database"
	"api/services/util/log"
	"api/services/util/tools"
	"encoding/xml"
	"fmt"
	"reflect"
	"regexp"
)

const (
	Sent      = "Sent"
	SentLeave = "SentLeave"
	Enter     = "Enter"
	PickUp    = "PickUp"
	Switch    = "Switch"
)

func checkEmptyFields(v interface{}) error {
	rv := reflect.ValueOf(v)
	typeOfS := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		if rv.Field(i).Interface() == "" {
			if OrderDate, ok := typeOfS.FieldByName(typeOfS.Field(i).Name); ok {
				match, _ := regexp.MatchString(`omitempty`, fmt.Sprintf(`%v`, OrderDate))
				if !match {
					return fmt.Errorf(`%s 欄位不能為空`, typeOfS.Field(i).Name)
				}
			} else {
				return fmt.Errorf(`%s 欄位不能為空`, typeOfS.Field(i).Name)
			}
		}
	}
	return nil
}

func handleNotifyReqParams(i interface{}) (FamilyMart.ErrorInfo, error)  {
	var errInfo = FamilyMart.ErrorInfo{ErrorMessage: `成功`, ErrorCode: `000`}
	var err error
	// 確認欄位是否有空值
	if err = checkEmptyFields(i); err != nil {
		errInfo.ErrorCode = `004`
		errInfo.ErrorMessage = err.Error()
		log.Error("handleNotifyReqParams Error [%s]", err.Error())
		return errInfo, err
	}

	return errInfo, nil
}

// 處理通知訊息
func HandleNotify(paramsXml FamilyMart.FamilyParams, action string) FamilyMart.ResponseXml {

	var rsp FamilyMart.ResponseXml

	log.Info("action : [%s]", action)
	log.Info("paramsXml : [%v]", paramsXml)

	engine := database.GetMysqlEngine()
	defer engine.Close()

	switch action {
	case Sent,SentLeave:
		var _xml FamilyMart.SendDoc
		_ = xml.Unmarshal([]byte(paramsXml.DataUrlEncode()), &_xml)
		rsp = sentNotify(engine, _xml, action)
	case Enter:
		var _xml FamilyMart.EnterDoc
		_ = xml.Unmarshal([]byte(paramsXml.DataUrlEncode()), &_xml)
		rsp = enterNotify(engine, _xml)
	case PickUp:
		var _xml FamilyMart.PickupDoc
		_ = xml.Unmarshal([]byte(paramsXml.DataUrlEncode()), &_xml)
		rsp = pickUpNotify(engine, _xml)
	case Switch:
		var _xml FamilyMart.SwitchDoc
		_ = xml.Unmarshal([]byte(paramsXml.DataUrlEncode()), &_xml)
		rsp = switchNotify(engine, _xml)
	}

	return rsp
}

// 寄件即時通知 寄件離店通知
func sentNotify(engine *database.MysqlSession, _xml FamilyMart.SendDoc, action string) FamilyMart.ResponseXml {
	var rsp FamilyMart.ResponseXml
	var rspBody FamilyMart.ResponseBody
	var errInfo FamilyMart.ErrorInfo
	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	var err error


	for _, i := range _xml.ShipmentNos {
		errInfo, err = handleNotifyReqParams(i)

		UpdateCvsShipping.ShipNo = i.OrderNo
		UpdateCvsShipping.ShipType = Enum.CVS_FAMILY
		UpdateCvsShipping.DateTime = i.GetDateTime()
		UpdateCvsShipping.DetailStatus = ""
		UpdateCvsShipping.Log = tools.XmlToString(i)

		if err == nil {
			switch action {
			case Sent:
				UpdateCvsShipping.Type = `R22Api`
				// 訂單貨運狀態為 未出貨 則修改狀態為 已出貨
				err = UpdateCvsShipping.UpdateCvsShippingShipment(engine)
			case SentLeave:
				UpdateCvsShipping.Type = `R23Api`
				// 訂單貨運狀態為 已出貨 則修改狀態為 配送中
				err = UpdateCvsShipping.UpdateCvsShippingTransit(engine)
			}

			if err != nil {
				errInfo.ErrorCode = `999`
				errInfo.ErrorMessage = err.Error()
			}
		}

		rspBody.SetErrorInfo(errInfo, i.ShipmentNos)
		rsp.ShipmentNos = append(rsp.ShipmentNos, rspBody)
	}

	return rsp
}

// 到店即時通知
func enterNotify(engine *database.MysqlSession, _xml FamilyMart.EnterDoc) FamilyMart.ResponseXml {
	var rsp FamilyMart.ResponseXml
	var rspBody FamilyMart.ResponseBody
	var errInfo FamilyMart.ErrorInfo
	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	var err error

	for _, i := range _xml.ShipmentNos {
		errInfo, err = handleNotifyReqParams(i)

		UpdateCvsShipping.ShipNo = i.OrderNo
		UpdateCvsShipping.ShipType = Enum.CVS_FAMILY
		UpdateCvsShipping.DateTime = i.GetDateTime()
		UpdateCvsShipping.DetailStatus = ""
		UpdateCvsShipping.Type = `R28Api`
		UpdateCvsShipping.FlowType = i.FlowType
		UpdateCvsShipping.Log = tools.XmlToString(i)

		// 訂單貨運狀態為 配送中 則修改狀態為 到店
		err = UpdateCvsShipping.UpdateCvsShippingShop(engine)

		if err != nil {
			errInfo.ErrorCode = `999`
			errInfo.ErrorMessage = err.Error()
		}

		rspBody.SetErrorInfo(errInfo, i.ShipmentNos)
		rsp.ShipmentNos = append(rsp.ShipmentNos, rspBody)
	}

	return rsp
}

// 取貨即時通知
func pickUpNotify(engine *database.MysqlSession, _xml FamilyMart.PickupDoc) FamilyMart.ResponseXml {
	var rsp FamilyMart.ResponseXml
	var rspBody FamilyMart.ResponseBody
	var errInfo FamilyMart.ErrorInfo
	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	var err error

	for _, i := range _xml.ShipmentNos {
		errInfo, err = handleNotifyReqParams(i)

		if err == nil {

			UpdateCvsShipping.ShipNo = i.OrderNo
			UpdateCvsShipping.ShipType = Enum.CVS_FAMILY
			UpdateCvsShipping.DateTime = i.GetDateTime()
			UpdateCvsShipping.DetailStatus = ""
			UpdateCvsShipping.Type = `R29Api`
			UpdateCvsShipping.FlowType = i.FlowType
			UpdateCvsShipping.Log = tools.XmlToString(i)

			// 訂單貨運狀態為 到店 則修改狀態為 取貨
			err = UpdateCvsShipping.UpdateCvsShippingSuccess(engine)

			if err != nil {
				errInfo.ErrorCode = `999`
				errInfo.ErrorMessage = err.Error()
			}
		}

		rspBody.SetErrorInfo(errInfo, i.ShipmentNos)
		rsp.ShipmentNos = append(rsp.ShipmentNos, rspBody)
	}

	return rsp
}

// 閉轉即時通知
func switchNotify(engine *database.MysqlSession, _xml FamilyMart.SwitchDoc) FamilyMart.ResponseXml {
	var rsp FamilyMart.ResponseXml
	var rspBody FamilyMart.ResponseBody
	var errInfo FamilyMart.ErrorInfo
	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	var err error

	for _, i := range _xml.ShipmentNos {
		errInfo, err = handleNotifyReqParams(i)

		if err == nil {

			UpdateCvsShipping.Type = `Switch`
			UpdateCvsShipping.ShipType = Enum.CVS_FAMILY
			UpdateCvsShipping.ShipNo = i.OrderNo
			UpdateCvsShipping.OrderNo = i.EcOrderNo
			UpdateCvsShipping.DateTime = tools.Now(`YmdHis`)
			UpdateCvsShipping.FlowType = `N`
			UpdateCvsShipping.Log = tools.XmlToString(i)

			if i.StoreType == `1` {
				UpdateCvsShipping.FlowType = `R`
			}

			// 訂單貨運狀態為 到店 則修改狀態為 取貨
			if err = UpdateCvsShipping.UpdateCvsShippingSwitch(engine); err != nil {
				errInfo.ErrorCode = `999`
				errInfo.ErrorMessage = err.Error()
			}
		}

		rspBody.SetSwitchErrorInfo(errInfo, i)
		rsp.ShipmentNos = append(rsp.ShipmentNos, rspBody)
	}

	return rsp
}
