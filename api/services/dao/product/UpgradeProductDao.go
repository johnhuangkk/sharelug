package product

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"time"
)

func InsertUpgradeProductData(engine *database.MysqlSession, data entity.UpgradeProductData) error {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.UpgradeProductData{}).Insert(&data)
	if err != nil {
		log.Error("insert product spec database error", err)
		return err
	}
	return nil
}

func GetUpgradeProductData(engine *database.MysqlSession) ([]entity.UpgradeProductData, error) {
	var data []entity.UpgradeProductData
	err := engine.Engine.Table(entity.UpgradeProductData{}).
		Select("*").Asc("product_id").Find(&data)
	if err != nil {
		log.Error("Get Upgrade Product Database Error", err)
		return data, err
	}
	return data, nil
}

func GetUpgradeProductDataByLevel(engine *database.MysqlSession, Level int64) (entity.UpgradeProductData, error) {
	var data entity.UpgradeProductData
	_, err := engine.Engine.Table(entity.UpgradeProductData{}).
		Select("*").Where("upgrade_level = ?", Level).Get(&data)
	if err != nil {
		log.Error("Get Upgrade Product Database Error", err)
		return data, err
	}
	return data, nil
}

func GetUpgradeProductByProductId(engine *database.MysqlSession, ProductId string) (entity.UpgradeProductData, error) {
	var data entity.UpgradeProductData
	_, err := engine.Engine.Table(entity.UpgradeProductData{}).
		Select("*").Where("product_id = ?", ProductId).Get(&data)
	if err != nil {
		log.Error("Get Upgrade Product Database Error", err)
		return data, err
	}
	return data, nil
}

func CreateUpgradeProductData() error {
	var resp []entity.UpgradeProductData
	res := entity.UpgradeProductData {
		ProductId:   "C0000011",
		ProductName: "NT$99方案",
		Description: "1個賣場,1個管理員帳號",
		Note: "",
		Amount: 99,
		UpgradeLevel: 1,
		Store: 0,
		Manager: 1,
	}
	resp = append(resp, res)
	res = entity.UpgradeProductData {
		ProductId:   "C0000015",
		ProductName: "NT$199方案",
		Description: "1個賣場,5個管理員帳號",
		Note: "",
		Amount: 199,
		UpgradeLevel: 2,
		Store: 0,
		Manager: 5,
	}
	resp = append(resp, res)
	res = entity.UpgradeProductData {
		ProductId:   "C0000051",
		ProductName: "NT$299方案",
		Description: "5個賣場,1個管理員帳號",
		Note: "(每個賣場1個管理帳號)",
		Amount: 299,
		UpgradeLevel: 3,
		Store: 5,
		Manager: 1,
	}
	resp = append(resp, res)
	res = entity.UpgradeProductData {
		ProductId:   "C0000055",
		ProductName: "NT$399方案",
		Description: "5個賣場,25個管理員帳號",
		Note: "(每個賣場5個管理帳號)",
		Amount: 399,
		UpgradeLevel: 4,
		Store: 5,
		Manager: 5,
	}
	resp = append(resp, res)
	engine := database.GetMysqlEngine()
	defer engine.Close()

	for _, v := range resp {
		err := InsertUpgradeProductData(engine, v)
		if err != nil {
			log.Error("Insert Upgrade Product Data", err)
			return err
		}
	}
	return nil
}

