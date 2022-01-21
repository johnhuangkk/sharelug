package Task

import (
	"api/services/Enum"
	"api/services/Service/Notification"
	"api/services/Service/Sms"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"time"
)

func getNoneShip() []string {
	return []string{Enum.F2F, Enum.NONE}
}

// 排程掃描正向到店第四天未領取 發出小鈴鐺通知
func ScanOrderShipStatusIsShopEqual4Days() {
	engine := database.GetMysqlEngine()
	defer engine.Close()

	var data []entity.OrderData

	now := time.Now().Add(-(time.Hour * time.Duration(24) * 3)).Format(`2006-01-02`)
	start := fmt.Sprintf(`%s 00:00:00`, now)
	end := fmt.Sprintf(`%s 23:59:59`, now)

	err := engine.Engine.Table(entity.OrderData{}).Select("*").
		Where(`arrived_time >= ? and arrived_time <= ?`, start, end).
		NotIn(`ship_type`, getNoneShip()).
		And(`ship_status = ?`, Enum.OrderShipShop).Find(&data)

	if err != nil {
		log.Error(`ScanOrderShipStatusIsShop Get Data Error`, err.Error())
		return
	}

	for _, d := range data {
		_ = Notification.SendShipToShopFourDayMessage(engine, d)
	}
}

// 排程掃描取號後兩天未寄出 發出小鈴鐺通知
func ScanOrderShipStatusIsInitEqual2Days() {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var data []entity.OrderData
	var seller []string

	now := time.Now().Add(-(time.Hour * time.Duration(24) * 1)).Format(`2006-01-02`)
	start := fmt.Sprintf(`%s 00:00:00`, now)
	end := fmt.Sprintf(`%s 23:59:59`, now)

	err := engine.Engine.Table(entity.OrderData{}).Distinct(`store_id`).
		Where(`pay_way_time >= ? and pay_way_time <= ?`, start, end).
		NotIn(`ship_type`, getNoneShip()).
		And(`ship_status = ?`, Enum.OrderSuccess).
		Find(&seller)

	if err != nil {
		log.Error(`ScanOrderShipStatusIsShop Get sellers Error`, err.Error())
		return
	}

	for _, s := range seller {
		err := engine.Engine.Table(entity.OrderData{}).Select("*").
			Where(`pay_way_time >= ? and pay_way_time <= ?`, start, end).
			NotIn(`ship_type`, getNoneShip()).
			And(`ship_status = ?`, Enum.OrderSuccess).
			And(`store_id = ?`, s).Find(&data)

		if err != nil {
			log.Error(`ScanOrderShipStatusIsShop Get sellers Error`, err.Error())
		}

		_ = Notification.SendNotShippedMessage(engine, s, data, 2)
	}
}


// OK 順向到貨通知 到店 四天 七天
func OkOrderShipIsShopSmsNotification() {
	OkOrderShipIsShopFewDays(3)
	OkOrderShipIsShopFewDays(6)
}

func OkOrderShipIsShopFewDays(days int64) {
	log.Info(`OkOrderShipIsShopFewDays`, days)
	engine := database.GetMysqlEngine()
	defer engine.Close()

	var data []entity.OrderData
	var storeName string

	now := time.Now().Add(-(time.Hour * time.Duration(24*days))).Format(`2006-01-02`)
	start := fmt.Sprintf(`%s 00:00:00`, now)
	end := fmt.Sprintf(`%s 23:59:59`, now)

	err := engine.Engine.Table(entity.OrderData{}).Select("*").
		Where(`arrived_time >= ? and arrived_time <= ?`, start, end).
		And(`ship_type = ?`, Enum.CVS_OK_MART).
		And(`ship_status = ?`, Enum.OrderShipShop).Find(&data)

	if err != nil {
		log.Error(`ScanOrderShipStatusIsShop Get sellers Error`, err.Error())
	}

	for _, d := range data {
		_, err = engine.Engine.Table(entity.MartOkStoreData{}).Where(`store_id = ?`, d.ReceiverAddress).Get(&storeName)

		template := fmt.Sprintf(
			`%s已送達%s%s，請於%s前攜帶證件前往取貨，若已取件請忽略此訊息。`,
			d.ShipNumber,
			Enum.Shipping[d.ShipType],
			storeName,
			d.ArrivedTime.Add(time.Hour*time.Duration(24*7)).Format(`2006-01-02`),
		)

		_ = Sms.PushMessageSms(d.ReceiverPhone, template)
	}

}
