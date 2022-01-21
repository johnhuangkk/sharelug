package sevenmyshipdao

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"time"
)

func InsertSevenShipMap(data entity.SevenShipMapData) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	_, err := engine.Session.Table("seven_ship_map_data").Insert(&data)
	if err != nil {
		log.Error("seven_ship_map_data Database insert Error", err)
		return
	}
}
func InsertAddressByDailyFile(datas []entity.SevenMyshipShopData) {
	session := database.Mysql().NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		log.Error("Seven Shop Adderss Database Insert Error", err)
		return
	}

	for _, data := range datas {
		has, err := session.Exist(&entity.SevenMyshipShopData{
			StoreID: data.StoreID,
		})

		if err != nil {
			log.Error("Seven Shop Adderss Database Exist Error", err.Error())
			return
		}
		if has {
			_, err = session.Where("store_id = ?", data.StoreID).Cols("opened", "address", "district", "city", "store_name").Update(&data)
			if err != nil {
				log.Error("Seven Shop Adderss Database update Error", err.Error())
				return
			}
		} else {
			_, err = session.Insert(&data)
			if err != nil {
				log.Error("Seven Shop Adderss Database Insert Error", err.Error())
				return
			}
		}

	}
	err = session.Commit()
	if err != nil {
		log.Error("Seven Shop Adderss Database commit Error", err.Error())
		return
	}

}
func UpdateClosedShop() {
	session := database.Mysql().NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		log.Error("Seven Shop Adderss Database Insert Error", err)
		return
	}
	tNow := time.Now()

	start := time.Date(tNow.Year(), tNow.Month(), tNow.Day(), 0, 0, 0, 0, tNow.Location())

	end := time.Date(tNow.Year(), tNow.Month(), tNow.Day(), 23, 59, 59, 999999999, tNow.Location())

	f, err := session.Table("seven_myship_shop_data").Where(`updated >= ? AND updated <= ?`, start.Format("2006-01-02 15:04"), end.Format("2006-01-02 15:04")).Count()
	if err != nil {
		log.Error("Seven Shop Adderss Database count Error", err.Error())
		return
	}
	if f >= 1 {
		_, err := session.Table("seven_myship_shop_data").Cols("opened", "updated").Where(`updated < ? AND opened = ?`, start.Format("2006-01-02 15:04"), true).Update(&entity.SevenMyshipShopData{
			Opened: false,
		})
		if err != nil {
			log.Error("Seven Shop Adderss Database update Error", err.Error())
			return
		}
		err = session.Commit()
		if err != nil {
			log.Error("Seven Shop Adderss Database commit Error", err.Error())
			return
		}
	}

}
func InsertProduct(engine *database.MysqlSession, data entity.ProductData) (entity.ProductData, error) {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()

	_, err := engine.Session.Table("product_data").Insert(&data)
	if err != nil {
		log.Error("Database Error", err)
		return data, err
	}
	return data, nil
}

func FindShopsAddress(country string, city string) ([]entity.SevenMyshipShopData, error) {
	var shops []entity.SevenMyshipShopData
	engine := database.GetMysqlEngine()
	defer engine.Close()
	err := engine.Engine.Table("seven_myship_shop_data").Where("country = ? AND district = ? AND opened =?", country, city, true).Find(&shops)
	if err != nil {
		return shops, nil
	}
	return shops, nil
}

func InsertChargeOrderRecords(data []entity.SevenChargeOrderData) ([]entity.SevenChargeOrderData, error) {

	engine := database.GetMysqlEngine()
	session := engine.Session
	defer session.Close()
	defer engine.Close()

	_, err := session.Insert(&data)

	if err != nil {
		log.Error("Database Charge Record Error", err)
		return data, err
	}
	return data, nil
}
