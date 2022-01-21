package model

import (
	"api/services/Enum"
	"api/services/Service/Balance"
	"api/services/Service/Excel"
	"api/services/VO/ExcelVo"
	"api/services/VO/Response"
	"api/services/dao/Orders"
	"api/services/dao/Store"
	"api/services/dao/TwId"
	"api/services/dao/Withdraw"
	"api/services/dao/member"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"github.com/spf13/viper"
)

//訂單報表 是否有扣除服務費
func HandleGetOrderReportExporter() (string, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var data []entity.OrderData
	data1, err := Orders.GetOrderByOrderStatus(engine, Enum.OrderInit)
	if err != nil {
		log.Debug("Read Ach Response Excel Error", err)
		return "", err
	}
	data = append(data, data1...)
	where := []string{"B210416019453", "B210416019560", "B210416019873"}
	data2, _ := Orders.GetOrderById(engine, where)
	data = append(data, data2...)

	var PlatformFee = map[int64]string{
		0: "未扣平台服務費",
		1: "已扣平台服務費",
	}
	var report []ExcelVo.OrderReportVo
	for k, v := range data {
		count, _ := Balance.CountBalanceByUserIdAndOrderId(engine, v.SellerId, v.OrderId)

		seller, err := member.GetMemberDataByUid(engine, v.SellerId)
		if err != nil {
			log.Debug("Get Member data Error", err)
		}
		buyer, err := member.GetMemberDataByUid(engine, v.BuyerId)
		if err != nil {
			log.Debug("Get Member data Error", err)
		}
		shipTime := ""
		if tools.InArray([]string{Enum.OrderShipment, Enum.OrderShipTransit, Enum.OrderShipShop, Enum.OrderShipSuccess, Enum.OrderShipNone}, v.ShipStatus) {
			shipTime = v.ShipTime.Format("2006/01/02 15:04")
		}
		payWayTime := ""
		if !v.PayWayTime.IsZero() {
			payWayTime = v.PayWayTime.Format("2006/01/02 15:04")
		}
		res := ExcelVo.OrderReportVo{
			Id:          int64(k + 1),
			OrderTime:   v.CreateTime.Format("2006/01/02 15:04"),
			PayWay:      Enum.PayWay[v.PayWay],
			PayWayTime:  payWayTime,
			ShipType:    Enum.Shipping[v.ShipType],
			ShipTime:    shipTime,
			OrderId:     v.OrderId,
			SellerId:    seller.TerminalId,
			BuyerId:     buyer.TerminalId,
			Amount:      int64(v.TotalAmount),
			PlatformFee: int64(v.PlatformTransFee + v.PlatformInfoFee + v.PlatformPayFee + v.PlatformShipFee),
			IsFee:       PlatformFee[count],
		}
		report = append(report, res)
	}
	filename, err := Excel.OrderNew().ToOrderReportFile(report)
	if err != nil {
		return "", fmt.Errorf("系統錯誤！")
	}
	log.Debug("filename", filename)
	return filename, nil
}
//報送銀行賣家資料
func HandleGetBankReportExporter() (string, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := member.GetAllMemberByVerifyIdentity(engine)
	if err != nil {
		log.Debug("Get Member Error", err)
		return "", err
	}
	Key := viper.GetString("EncryptKey")
	var report []ExcelVo.BankReportVo
	for _, v := range data {
		if Withdraw.CountWithdrawByUserId(engine, v.Uid) == 0 {
			continue
		}
		var res ExcelVo.BankReportVo
		res.Account = v.Mphone
		res.TerminalId = v.TerminalId
		res.Category = v.Category
		res.Name = v.IdentityName
		if len(v.Identity) != 0 {
			res.Identity = tools.AesDecrypt(v.Identity, Key)
		}
		if v.Category == Enum.CategoryCompany {
			res.Head = v.Representative
			res.HeadIdentity = tools.AesDecrypt(v.RepresentativeId, Key)
			res.CompanyAddr = v.CompanyAddr
		}
		storeData, err := getStoreName(engine, v.Uid)
		if err != nil {
			log.Debug("Get Store Error", err)
			return "", err
		}
		res.Store1 = storeData[0]
		res.Store2 = storeData[1]
		res.Store3 = storeData[2]
		res.Store4 = storeData[3]
		res.Store5 = storeData[4]
		report = append(report, res)
		v.ReportBank = true
		if _, err := member.UpdateMember(engine, &v); err != nil {
			log.Error("Update Member Error", err)
			return "", err
		}
	}
	filename, err := Excel.BankNew().ToBankReportFile(report)
	if err != nil {
		return "", fmt.Errorf("系統錯誤！")
	}
	log.Debug("filename", filename)
	return filename, nil
}
//取出賣家餘額，待撥付餘額，保留餘額，扣留餘額
func HandleGetBalancesReportExporter() (string, error) {
	res := make(chan Response.SellerBalanceResponse)
	engine := database.GetMysqlEngine()
	defer engine.Close()
	user, err := member.GetAllMemberData(engine)
	if err != nil {
		log.Debug("Get All Members Error", err)
		return "", fmt.Errorf("1001001")
	}
	var  report []Response.SellerBalanceResponse
	//log.Debug("user = ", user)
	for _, v := range user {
		go GetUser(engine, v, res)
		s := <-res
		report = append(report, s)
		//close(res)
		log.Debug("chan = ", s)
	}
	log.Debug("set excel")
	filename, err := Excel.SellerBalanceNew().ToSellerBalanceReportFile(report)
	if err != nil {
		return "", fmt.Errorf("1001001")
	}
	log.Debug("filename", filename)
	return filename, nil
}


func GetUser(engine *database.MysqlSession, user entity.MemberData, resp chan Response.SellerBalanceResponse) {
	log.Debug("get user", user.Uid)
	var res Response.SellerBalanceResponse
	if Orders.CountOrderDataBySeller(engine, user.Uid) > 0 {
		res.Account = user.Mphone
		res.SellerId = user.TerminalId
		//取出餘額
		res.Balance = int64(Balance.GetBalanceByUid(engine, user.Uid))
		res.RetainBalance = int64(Balance.GetBalanceRetainsByUid(engine, user.Uid))
		res.DetainBalance = 0
		res.WithholdBalance = 0
	}
	resp <- res
}


func getStoreName(engine *database.MysqlSession, userId string) ([]string, error) {
	name := []string{"", "", "", "", ""}
	storeData, err := Store.GetStoresByUid(engine, userId)
	if err != nil {
		log.Debug("Get Store Error", err)
		return nil, err
	}
	for k, v := range storeData {
		name[k] = v.StoreName
	}
	return name, nil
}

func ModifyMemberIdentity() error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := member.GetAllMemberByVerifyIdentity(engine)
	if err != nil {
		log.Debug("Get Member Error", err)
		return err
	}
	for _, v := range data {
		if v.Category == Enum.CategoryMember {
			identity, err := TwId.GetTwIdLogDataByUserId(engine, v.Uid)
			if err != nil {
				log.Debug("Get Identity Error", err)
				return err
			}
			if len(identity.IdentityName) != 0 {
				if len(v.IdentityName) == 0 || len(v.Identity) == 0 {
					v.IdentityName = identity.IdentityName
					if len(v.Identity) == 0 {
						Key := viper.GetString("EncryptKey")
						v.Identity = tools.AesEncrypt(identity.IdentityId, Key)
					}
					if _, err := member.UpdateMember(engine, &v); err != nil {
						log.Error("Update Member Error", err)
						return err
					}
				}
			}
		}
	}
	return nil
}