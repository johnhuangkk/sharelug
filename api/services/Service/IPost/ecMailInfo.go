package IPost

import (
	"api/services/Enum"
	"api/services/VO/IPOSTVO"
	"api/services/VO/ShipmentVO"
	"api/services/dao/Store"
	"api/services/dao/sequence"
	"api/services/database"
	"api/services/entity"
	"api/services/model"
	"api/services/util/curl"
	"api/services/util/log"
	"api/services/util/tools"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"time"
)

// 取得 sequence iPost value
func getMailNo() string {
	value, _ := sequence.GetIPostMailNoSeq()
	return tools.StringPadLeft(value, 6) + "10127770"
}

// 設定參數
func setPayLoad(engine *database.MysqlSession, orderData entity.OrderData, sellerData entity.MemberData, SellerAddress *ShipmentVO.SellerSenderAddress) IPOSTVO.PayLoad {
	var receiverAddrType, vipNO string
	var config = viper.GetStringMapString("IPOST")
	payLoad := IPOSTVO.PayLoad{}

	if orderData.ShipType == "I_POST" {
		receiverAddrType = "002"
		vipNO = config["ipost2ipost"]
	} else {
		receiverAddrType = "001"
		vipNO = config["ipost2home"]
	}

	storeData, err := Store.GetStoreDataByStoreId(engine, orderData.StoreId)
	if err != nil {
		log.Error("Get Store Data Error", err)
	}

	payLoad.APLBR_VipNO = vipNO
	payLoad.APLBR_VipName = config["company"]
	payLoad.MailNo = getMailNo()
	payLoad.PrnAuthCode = orderData.ReceiverPhone[6:]
	payLoad.EcOrderNo = orderData.OrderId
	payLoad.EcOrderDate = time.Now().Format("20060102150405")
	payLoad.PrdName = fmt.Sprintf("Check`Ne[%s]-商品", storeData.StoreName)

	payLoad.SenderName = sellerData.SendName
	payLoad.SenderPhone = sellerData.Mphone

	if len(SellerAddress.Id) != 0 {
		payLoad.SenderAddrType = "002"
		payLoad.SenderAddrIBoxId = SellerAddress.Id
		payLoad.SenderZipCode = SellerAddress.Zip
		payLoad.SenderAddr = SellerAddress.Address
	} else {
		payLoad.SenderAddrType = "001"
		payLoad.SenderAddrIBoxId = ""
		payLoad.SenderZipCode = SellerAddress.Zip
		payLoad.SenderAddr = SellerAddress.Address
	}

	payLoad.ReceiverName = orderData.ReceiverName
	payLoad.ReceiverPhone = orderData.ReceiverPhone
	payLoad.ReceiverAddrType = receiverAddrType

	if orderData.ShipType == Enum.I_POST {
		// 配送方式 郵箱到郵箱
		receiverIPost := model.GetPostBoxAddressById(engine, orderData.ReceiverAddress)
		payLoad.ReceiverZipCode = receiverIPost.Zip
		payLoad.ReceiverAddr = receiverIPost.Address
		payLoad.ReceiverAddrIBoxID = receiverIPost.Id
		//payLoad.Remark = "1.收件人不得申請改投改寄，3日未取於招領期過後逕退寄件人ｉ郵箱。<br />2.倘收件ｉ郵箱已遷移，請務必刷讀未妥投原因為＂無此地址＂後退回寄件人ｉ郵箱。"
	} else {
		// 配送方式 郵箱到宅配
		// 300,新竹市,東區,中華東路1段366號[地下一樓美食街靠近化妝室]
		rAddressSplit := strings.Split(orderData.ReceiverAddress, ",")
		log.Info("rAddressSplit", rAddressSplit)
		payLoad.ReceiverZipCode = rAddressSplit[0]
		payLoad.ReceiverAddr = strings.Join(rAddressSplit[1:], "")
		//payLoad.Remark = "若收件人未於送達 i 郵箱後 3 日內取件，會直接退回 i 郵箱給寄件人。"
	}

	payLoad.ReturnType = 1
	payLoad.ValidDate = orderData.ShipExpire.Format("20060102")
	payLoad.ValidDate_PrintAndSend = payLoad.ValidDate

	return payLoad
}

/**
產生token
*/
func getToken(timeStamp string, payLoad IPOSTVO.PayLoad) string {
	str := payLoad.APLBR_VipNO + payLoad.EcOrderNo + payLoad.EcOrderDate
	str += payLoad.SenderPhone + payLoad.ReceiverPhone + payLoad.ValidDate

	tsMd5 := timeStamp + "," + strings.ToUpper(tools.MD5(str))
	SHA384Key := viper.GetString("IPOST.SHA384Key")
	sha384 := tools.SHA384Mac(tsMd5, SHA384Key)
	binary, _ := hex.DecodeString(sha384)
	base64String := tools.Base64EncodeByString(string(binary))
	big5 := tools.Utf8ToBig5(base64String)
	return big5
}

// 取得交寄取號參數
func getUpEcMailInfo(engine *database.MysqlSession, orderData entity.OrderData, sellerData entity.MemberData, SellerAddress *ShipmentVO.SellerSenderAddress) IPOSTVO.UP_ECMAILINFO {
	upEcMailInfo := IPOSTVO.UP_ECMAILINFO{}
	upEcMailInfo.TimeStamp = time.Now().Format("20060102150405")
	upEcMailInfo.PayLoad = setPayLoad(engine, orderData, sellerData, SellerAddress)
	upEcMailInfo.Token = getToken(upEcMailInfo.TimeStamp, upEcMailInfo.PayLoad)

	return upEcMailInfo
}

// i郵箱取號
func AddShippingOrder(engine *database.MysqlSession, orderData entity.OrderData, sellerData entity.MemberData, sellerAddress *ShipmentVO.SellerSenderAddress) (string, error) {
	upEcMailInfo := getUpEcMailInfo(engine, orderData, sellerData, sellerAddress)
	jsonData, _ := json.Marshal(upEcMailInfo)
	log.Debug(string(jsonData))

	var rspEcMailInfo = &IPOSTVO.RspEcMailInfo{}
	var uri = viper.GetString("IPOST.iPostEcMailInfoPath")
	log.Debug("iPOST URI [%s]", uri)

	rsp, _ := curl.PostJson(uri, upEcMailInfo)
	_ = json.Unmarshal(rsp, rspEcMailInfo)
	jsonData, _ = json.Marshal(rspEcMailInfo)
	log.Info("EcMailChange callPostApi rsp.", string(jsonData))

	// 無錯誤訊息則寫入 郵寄編號
	if len(rspEcMailInfo.Failures) != 0 || ( len(rspEcMailInfo.MailNo) == 0 && len(rspEcMailInfo.Failures) == 0 ) {
		log.Error("%s Fail: %s", orderData.OrderId, rspEcMailInfo.Failures)
		return rspEcMailInfo.MailNo, fmt.Errorf("%s", "系統忙碌中，請稍後再進行出貨。")
	}

	// 成功取號寫入i郵箱託運單
	_ = model.InsertPostConsignmentData(engine, upEcMailInfo)

	return rspEcMailInfo.MailNo, nil
}


