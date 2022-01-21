package model

import (
	"api/services/Enum"
	"api/services/VO/Response"
	"api/services/database"
	"api/services/entity"
)

//出貨明細
func GetErpOrderShippingInfo(engine *database.MysqlSession, OrderData entity.OrderData) (Response.ErpOrderShipInfoResponse, error) {
	var resp Response.ErpOrderShipInfoResponse
	resp.ShipMethod = Enum.Shipping[OrderData.ShipType]
	switch OrderData.ShipType {
	case Enum.DELIVERY_POST_BAG1, Enum.DELIVERY_POST_BAG2, Enum.DELIVERY_POST_BAG3, Enum.DELIVERY_I_POST_BAG1, Enum.DELIVERY_I_POST_BAG2, Enum.DELIVERY_I_POST_BAG3, Enum.I_POST:
		resp.ShipTrader = "中華郵政"
	case Enum.CVS_7_ELEVEN:
		resp.ShipTrader = "統一超商"
	case Enum.CVS_FAMILY:
		resp.ShipTrader = "全家超商"
	case Enum.CVS_OK_MART:
		resp.ShipTrader = "OK超商"
	case Enum.CVS_HI_LIFE:
		resp.ShipTrader = "萊爾富超商"
	case Enum.DELIVERY_T_CAT:
		resp.ShipTrader = "黑貓宅急便"
	case Enum.DELIVERY_E_CAN:
		resp.ShipTrader = "宅配通"
	case Enum.DELIVERY_OTHER:
		resp.ShipTrader = OrderData.ShipText
	}
	resp.RecipientName = OrderData.ReceiverName
	resp.RecipientPhone = OrderData.ReceiverPhone
	//運送地止
	ship := setShipInfo(engine, OrderData)
	resp.ReceiptAddress = ship.Receiver.ReceiverAddress
	//fixme 退貨地址
	resp.RefundAddress = "---"
	return resp, nil
}

