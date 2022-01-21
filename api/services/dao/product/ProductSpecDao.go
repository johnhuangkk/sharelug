package product

import (
	"api/services/Enum"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"time"
)

func InsertProductSpec(engine *database.MysqlSession, data entity.ProductSpecData) (entity.ProductSpecData, error) {
	data.CreateTime = time.Now()
	_, err := engine.Session.Table(entity.ProductSpecData{}).Insert(&data)
	if err != nil {
		log.Error("insert product spec database error", err)
		return data, err
	}
	return data, nil
}

func GetProductSpecByProductSpecId(engine *database.MysqlSession, ProductSpecId string) (entity.ProductSpecData, error) {
	var data entity.ProductSpecData
	_, err := engine.Engine.Table(entity.ProductSpecData{}).Select("*").
		Where("spec_id = ?", ProductSpecId).Get(&data)
	if err != nil {
		log.Error("Get ProductSpec Database Error", err)
		return data, err
	}
	return data, nil
}

func GetProductSpecByProductId(engine *database.MysqlSession, productId string) ([]entity.ProductSpecData, error) {
	var data []entity.ProductSpecData
	var err = engine.Engine.Table(entity.ProductSpecData{}).Select("*").
		Where("product_id = ?", productId).And("spec_status = ?", Enum.ProductStatusSuccess).Asc("spec_id").Find(&data)
	if err != nil {
		log.Error("Database Error", err)
		return data, err
	}
	return data, nil
}

func UpdateProductSpec(engine *database.MysqlSession, SpecId string, data entity.ProductSpecData) error {
	_, err := engine.Session.Table(entity.ProductSpecData{}).ID(SpecId).AllCols().Update(data)
	if err != nil {
		return err
	}
	return nil
}

func UpdateProductSpecDelete(engine *database.MysqlSession, ProductId string) error {
	sql := fmt.Sprintf("UPDATE product_spec_data SET spec_status = ? WHERE product_id = ?")
	_, err := engine.Session.Exec(sql, Enum.ProductStatusDelete, ProductId)
	if err != nil {
		return err
	}
	return nil
}