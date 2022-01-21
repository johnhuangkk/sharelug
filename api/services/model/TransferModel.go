package model

import (
	"api/services/Enum"
	"api/services/Service/Balance"
	"api/services/Service/Excel"
	"api/services/Service/OrderService"
	"api/services/Service/Soap"
	"api/services/Service/TransferService"
	"api/services/dao/Credit"
	"api/services/dao/transfer"
	"api/services/database"
	"api/services/entity"
	"api/services/entity/Response"
	"api/services/util/log"
	"api/services/util/tools"
	"api/services/util/xml"
	"fmt"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

//轉帳入金處理
func HandleTransfer(param entity.TransferParams) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	err := engine.Session.Begin()
	data, err := transfer.GetTransferByAccount(engine, param.Accno)
	if err != nil {
		log.Error("get transfer data error", err)
		return err
	}
	t, _ := time.ParseInLocation("20060102150405", param.Tdate+param.Ttime, time.Local)
	if len(data.OrderId) != 0 && data.TransferStatus == Enum.TransferInit {
		data.RecdBankAccount = param.Raccno
		data.RecdAmount = param.Amt
		data.RecdDate = t
		data.Seqno = param.Seqno
		data.TransferStatus = Enum.TransferSuccess
		if err := transfer.UpdateTransferDate(engine, data.Id, data); err != nil {
			log.Error("update transfer data Error", err)
			engine.Session.Rollback()
			return err
		}
		//更改訂單狀態
		if data.TransType == Enum.OrderTransC2c {
			if err := Balance.OrderCheckout(engine, data.OrderId, Enum.OrderSuccess); err != nil {
				log.Error("update order status Error", err)
				engine.Session.Rollback()
				return err
			}
		} else {
			if err := Balance.B2cOrderCheckout(engine, data.OrderId, Enum.OrderSuccess, ""); err != nil {
				log.Error("update order status Error", err)
				engine.Session.Rollback()
				return err
			}
		}
	} else if param.Resend == "I" {
		data.BankAccount = param.Accno
		data.RecdBankAccount = param.Raccno
		data.RecdAmount = param.Amt
		data.RecdDate = t
		data.Seqno = param.Seqno
		data.TransferStatus = Enum.TransferDuplicate
		if _, err := transfer.InsertTransfer(engine, data); err != nil {
			log.Error("insert transfer data Error", err)
			engine.Session.Rollback()
			return err
		}
	}
	err = engine.Session.Commit()
	if err != nil {
		return err
	}
	return nil
}
//轉帳回傳結果 寫入LOG
func TransferResponse(params entity.TransferParams) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	json, _ := tools.JsonEncode(params)
	if err := transfer.InsertTransferLog(engine, json); err != nil {
		return fmt.Errorf("1001001")
	}
	return nil
}
//轉帳入金處理(上海銀行)
func HandleTransferForSCSBank(param entity.ScsBankAccountData) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	err := engine.Session.Begin()
	data, err := transfer.GetTransferByAccount(engine, param.VAccount)
	if err != nil {
		log.Error("get transfer data error", err)
		return err
	}
	t, _ := time.Parse("20060102150405", param.TxDate+param.TxTime)
	if data.OrderId != "" && data.TransferStatus == Enum.TransferInit {
		data.RecdBankAccount = param.SAccount
		data.RecdAmount = strconv.Itoa(int(param.Amount))
		data.RecdDate = t
		data.Seqno = param.SeqNo
		data.TransferStatus = Enum.TransferSuccess
		if err := transfer.UpdateTransferDate(engine, data.Id, data); err != nil {
			log.Error("update transfer data Error", err)
			engine.Session.Rollback()
			return err
		}
		//更改訂單狀態
		if err := Balance.OrderCheckout(engine, data.OrderId, Enum.OrderSuccess); err != nil {
			log.Error("update order status Error", err)
			engine.Session.Rollback()
			return err
		}
	}
	err = engine.Session.Commit()
	if err != nil {
		return err
	}
	return nil
}
//C2C轉帳入金查詢
func QueryC2cTransDateTransfer(StartDate, EndData, Temp string) (Response.SMX, []Response.DETAIL, error) {
	var smx Response.SMX
	content, err := TransferService.GenerateTransferXml("", StartDate, EndData, Temp, Enum.OrderTransC2c)
	if err != nil {
		log.Error("Generate Transfer Xml error", err)
		return smx, nil, err
	}
	log.Debug("xml", content)
	hostname := viper.GetString("KgiBank.Path")
	body, _ := Soap.Call(hostname, content)
	smx, detail, err := xml.TransferXmlDecoder(body)
	if err != nil {
		log.Debug("Xml Decoder Error", err)
		return smx, detail, err
	}
	return smx, detail, nil
}
//B2C轉帳入金查詢
func QueryB2cTransDateTransfer(StartDate, EndData, Temp string) (Response.SMX, []Response.DETAIL, error) {
	var smx Response.SMX
	content, err := TransferService.GenerateTransferXml("", StartDate, EndData, Temp, Enum.OrderTransB2c)
	if err != nil {
		log.Error("Generate Transfer Xml error", err)
		return smx, nil, err
	}
	log.Debug("xml", content)
	hostname := viper.GetString("KgiBank.Path")
	body, _ := Soap.Call(hostname, content)
	smx, detail, err := xml.TransferXmlDecoder(body)
	if err != nil {
		log.Debug("Xml Decoder Error", err)
		return smx, detail, err
	}
	return smx, detail, nil
}
//查詢帳號
func QueryC2cAccountTransfer(Account, startDate, endData string) (Response.SMX, []Response.DETAIL, error) {
	var smx Response.SMX
	content, err := TransferService.GenerateTransferAccountQueryXml(Account, startDate, endData, Enum.OrderTransC2c)
	log.Error("Generate Transfer Xml", content)
	if err != nil {
		log.Error("Generate Transfer Xml error", err)
		return smx, nil, err
	}
	hostname := viper.GetString("KgiBank.Path")
	body, _ := Soap.Call(hostname, content)
	smx, response, err := xml.TransferXmlDecoder(body)
	if err != nil {
		log.Debug("Xml Decoder Error", err)
		return smx, response, err
	}
	log.Debug("response", response)
	return smx, response, nil
}

func ProcessTransfer(engine *database.MysqlSession, resp Response.DETAIL) error {
	if err := QueryTransferResponse(engine, resp); err != nil {
		log.Error("insert transfer log data error", err)
	}
	data, err := transfer.GetTransferByAccount(engine, resp.REOMAIL)
	if err != nil {
		log.Error("get transfer data error", err)
		return err
	}
	//log.Debug("order not success", data)
	if len(data.OrderId) != 0 && data.TransferStatus != Enum.TransferSuccess {
		log.Debug("order not success", data)
		t, _ := time.ParseInLocation("20060102150405", resp.TXNDATE + resp.TXNTIME, time.Local)
		data.RecdBankAccount = fmt.Sprintf("%s%s", resp.REOSND[:3], resp.REONAME)
		data.RecdAmount = resp.AMT
		data.RecdDate = t
		data.Seqno = resp.NUM
		data.TransferStatus = Enum.TransferSuccess
		log.Debug("write order success", data)
		if err := transfer.UpdateTransferDate(engine, data.Id, data); err != nil {
			log.Error("update transfer data Error", err)
			return err
		}
		//更改訂單狀態
		if data.TransType == Enum.OrderTransC2c {
			if err := Balance.TransferC2CCheckout(engine, data.OrderId, Enum.OrderSuccess, t); err != nil {
				log.Error("update order status Error", err)
				return err
			}
		} else {
			if err := Balance.B2cOrderCheckout(engine, data.OrderId, Enum.OrderSuccess, ""); err != nil {
				log.Error("update order status Error", err)
				return err
			}
		}
	}
	return nil
}

func ProcessTransferExpire(engine *database.MysqlSession) error {
	now := time.Now().Format("2006-01-02 15:04")
	log.Debug("now time", now)
	//取出過期未付的轉帳單
	data, err := transfer.GetTransferExpire(engine, now)
	if err != nil {
		log.Error("Get Transfer ")
		return err
	}
	log.Debug("Expire time", data)
	for _, v := range data {
		//過期修改狀態
		if v.TransType == Enum.OrderTransC2c {
			if err := Balance.OrderCheckout(engine, v.OrderId, Enum.OrderExpire); err != nil {
				log.Error("update order status Error", err)
				return err
			}
			if err := OrderService.ProcessReturnStock(engine, v.OrderId); err != nil {
				log.Error("Process Return Stock Error", err)
			}
		} else {
			if err := Balance.B2cOrderCheckout(engine, v.OrderId, Enum.OrderExpire, ""); err != nil {
				log.Error("update order status Error", err)
				return err
			}
		}
		v.TransferStatus = Enum.OrderExpire
		if err := transfer.UpdateTransferDate(engine, v.Id,  v); err != nil {
			log.Error("update Transfer status Error", err)
		}
	}
	return nil
}

//轉帳回傳結果 寫入LOG
func QueryTransferResponse(engine *database.MysqlSession, params Response.DETAIL) error {
	json, _ := tools.JsonEncode(params)
	if err := transfer.InsertTransferLog(engine, json); err != nil {
		return fmt.Errorf("1001001")
	}
	return nil
}


func ProcessTransferQuery(engine *database.MysqlSession) error {
	//now := time.Now().AddDate(0, 0, -1).Format("20060102")
	//log.Debug("date", "20210707")
	smx, res, err := QueryC2cAccountTransfer("8000100057411894", "20210705", "20210711")
	if err != nil {
		log.Error(" Query Transfer ")
		return err
	}
	log.Debug("res", smx, res)
	//for _, v := range res.SvcRs.DETAIL {
	//	log.Debug("Account", v.REOMAIL)
	//	//取出未付的轉帳單
	//	data, err := transfer.GetTransferByAccount(engine, v.REOMAIL)
	//	if err != nil {
	//		log.Error("Get Transfer Data Error", err)
	//		return err
	//	}
	//	if len(data.OrderId) != 0 {
	//		log.Debug("OrderId", data)
	//	}
	//}
	return nil
}

func HandleSpecialStoreCode(filename string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	excel, err := Excel.ReadSpecialStoreExcel(filename)
	if err != nil {
		log.Debug("Read Order Ship Excel Error", err)
		return fmt.Errorf("1001001")
	}
	for _, v := range excel {
		var data entity.KgiSpecialStore
		data.MerchantId = v.MerchantId
		data.Terminal3dId = v.Terminal3DId
		data.Terminaln3dId = v.TerminalN3DId
		data.ChStoreName = v.ChStoreName
		data.EnStoreName = v.EnStoreName
		if err := Credit.InsertSellerMerchantId(engine, data); err != nil {
			log.Error("Insert Seller Merchant Id Data Error", err)
			continue
		}
	}
	return nil
}