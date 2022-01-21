package model

import (
	"api/services/VO/IPOSTVO"
	"api/services/dao/iPost"
	"api/services/database"
	"api/services/entity"
)


/**
	寫入賣家貨運寄送地址資訊
 */
func InsertPostConsignmentData(engine *database.MysqlSession, upEcMailInfo IPOSTVO.UP_ECMAILINFO) error {
	var ecSellerInfo = entity.PostConsignmentData{}
	var payLoad = upEcMailInfo.PayLoad

	ecSellerInfo.OrderId = payLoad.EcOrderNo
	ecSellerInfo.MerchantId = payLoad.APLBR_VipNO
	ecSellerInfo.ShipNumber = payLoad.MailNo
	ecSellerInfo.SellerName = payLoad.SenderName
	ecSellerInfo.SellerPhone = payLoad.SenderPhone
	ecSellerInfo.SellerZip = payLoad.SenderZipCode
	ecSellerInfo.SellerAddr = payLoad.SenderAddr

	_, err := iPost.InsertPostConsignmentData(engine, ecSellerInfo)
	if err != nil {
		return err
	}

	return nil
}
