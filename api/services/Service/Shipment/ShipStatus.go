package Shipment

import (
	"api/services/Enum"
	"api/services/VO/ShipmentVO"
	"api/services/dao/Cvs"
	"api/services/dao/Orders"
	"api/services/dao/iPost"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"time"
)

type ShipStatus struct {
	Time string `json:"time"`
	Text string `json:"text"`
}

// 處理超商貨態紀錄
func handleCvsRecord(record []entity.CvsShippingLogData) []ShipStatus {
	var shipStatusAry []ShipStatus
	var shipStatus ShipStatus

	if len(record) == 0 {
		return shipStatusAry
	}

	for _, r := range record {
		t, _ := time.Parse(time.RFC3339, r.DateTime)
		shipStatus.Time = t.Format(`2006-01-02 15:04:05`)
		shipStatus.Text = r.Text
		shipStatusAry = append(shipStatusAry, shipStatus)
	}

	return shipStatusAry
}

func getShipStatus(engine *database.MysqlSession, orderData entity.OrderData) ([]ShipStatus, error) {

	var shipStatusAry []ShipStatus
	var shipStatus ShipStatus

	switch orderData.ShipType {
	case Enum.CVS_7_ELEVEN, Enum.CVS_FAMILY, Enum.CVS_HI_LIFE, Enum.CVS_OK_MART:
		data := Cvs.GetCvsShippingLogData(engine, orderData)
		return handleCvsRecord(data), nil
	default: // I_POST
		ip, err := iPost.QueryPostShippingStatus(engine, orderData.ShipNumber, orderData.CreateTime)

		if err != nil {
			log.Error("QueryPostShippingStatus Error: [%]", err)
			return shipStatusAry, fmt.Errorf("貨態查詢異常")
		}

		for _, d := range ip {
			shipStatus.Time = d.HandleTime
			shipStatus.Text = d.Branch + d.ShippingStatus
			shipStatusAry = append(shipStatusAry, shipStatus)
		}
	}

	return shipStatusAry, nil
}

func checkOwner(orderData entity.OrderData, uid, sid string) error {
	if orderData.StoreId != sid && orderData.BuyerId != uid {
		log.Error("使用者[%s]異常查詢", uid)
		return fmt.Errorf("系統錯誤")
	}
	return nil
}

// 查詢貨態資料
func SearchSippingStatus(order ShipmentVO.Order, uid, sid string) ([]ShipStatus, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()

	var shipStatusAry []ShipStatus

	orderData, _ := Orders.GetOrderByOrderId(engine, order.OrderId)

	if err := checkOwner(orderData, uid, sid); err != nil {

		return shipStatusAry, err
	}
//
	shipStatusAry, err := getShipStatus(engine, orderData)

	log.Debug("shipStatusAry", shipStatusAry)
	if err != nil {
		return shipStatusAry, err
	}

	return shipStatusAry, nil

}