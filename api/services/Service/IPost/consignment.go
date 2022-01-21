package IPost

import (
	"api/services/Enum"
	"api/services/dao/iPost"
	"api/services/database"
	"api/services/entity"
	"api/services/model"
	"api/services/util/log"
	"github.com/spf13/viper"
	"strings"
)

type address struct {
	Name, Mobile, Zipcode, Alias, Address string
}

type PostConsignmentData struct {
	MerchantId, MerchantName, ShopName, ShipNumber, Mark string
	Seller, Receiver address

}

// 取得 托運單資料
func GetPostConsignmentData(orderId []string, StoreDataResp entity.StoreDataResp) []PostConsignmentData {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var consignmentDataAry []PostConsignmentData
	var consignmentData = &PostConsignmentData{}

	data, _ := iPost.QueryPostConsignmentData(engine, orderId)
	//
	for _, d := range data {
		consignmentData.Mark = "若收件人未於送達 i 郵箱後 3 日內取件，會直接退回寄件人交寄之 i 郵箱。"
		consignmentData.MerchantId = d.MerchantId
		consignmentData.MerchantName = viper.GetString("IPOST.company")
		consignmentData.ShopName = StoreDataResp.StoreName
		consignmentData.ShipNumber = d.ShipNumber
		consignmentData.Seller.Name = d.SellerName
		consignmentData.Seller.Mobile = d.SellerPhone
		consignmentData.Seller.Zipcode = d.SellerZip
		sAddr := strings.Split(d.SellerAddr, ",")

		if len(sAddr) > 1 {
			consignmentData.Seller.Alias = sAddr[0]
			consignmentData.Seller.Address = sAddr[1]
		} else {
			consignmentData.Seller.Alias = ""
			consignmentData.Seller.Address = d.SellerAddr
		}

		consignmentData.Receiver.Name = d.ReceiverName
		consignmentData.Receiver.Mobile = d.ReceiverPhone

		if d.ShipType == Enum.I_POST {
			rIPostAddr := model.GetPostBoxAddressById(engine, d.ReceiverAddress)
			consignmentData.Receiver.Zipcode = rIPostAddr.Zip
			consignmentData.Receiver.Alias = rIPostAddr.Alias
			consignmentData.Receiver.Address = rIPostAddr.Address
		} else {
			rAddr := strings.Split(d.ReceiverAddress, ",")
			consignmentData.Receiver.Zipcode = rAddr[0]
			consignmentData.Receiver.Alias = ""
			consignmentData.Receiver.Address = strings.Join(rAddr[1:], "")
		}

		consignmentDataAry = append(consignmentDataAry, *consignmentData)
	}

	log.Info("consignmentDataAry", consignmentDataAry)

	return consignmentDataAry
}
