package product

import (
	"api/services/Enum"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/member"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"strings"
	"time"
)

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

//取出一般商品資料
//func GetProductsByProductId(engine *database.MysqlSession, productId string) ([]entity.ProductsData, error) {
//	var productData []entity.ProductsData
//	var err = engine.Engine.Table(entity.ProductData{}).
//		Select("*").
//		Join("Left", entity.ProductSpecData{}, "product_data.product_id = product_spec_data.product_id").
//		Where("product_data.product_id = ?", productId).
//		And("spec_status = ?", Enum.ProductStatusSuccess).
//		OrderBy("product_spec_data.spec_id ASC").Find(&productData)
//	if err != nil {
//		log.Error("get products Database Error", err)
//		return nil, err
//	}
//	return productData, nil
//}

//取出即時商品資料
func GetRealTimeByProductId(engine *database.MysqlSession, productId string) (entity.ProductData, error) {
	var data entity.ProductData
	_, err := engine.Engine.Table(entity.ProductData{}).Select("*").Where("product_id = ?", productId).
		And("is_realtime = ?", 1).Get(&data)
	if err != nil {
		log.Error("get RealTim Database Error", err)
		return data, err
	}
	return data, nil
}

//取商品資料BY spec
func GetProductByProductSpecId(engine *database.MysqlSession, productSpecId string) (entity.ProductsData, error) {
	var productData entity.ProductsData
	var _, err = engine.Engine.Table(entity.ProductData{}).
		Join("Left", entity.ProductSpecData{}, "product_data.product_id = product_spec_data.product_id").
		Select("*").Where("product_spec_data.spec_id = ?", productSpecId).
		Limit(1, 0).Get(&productData)
	if err != nil {
		log.Error("get products Database Error", err)
		return productData, err
	}
	return productData, nil
}

func GetProductListByStoreId(engine *database.MysqlSession, where []string, bind []interface{}, limit int, start int) ([]entity.ProductsData, error) {
	var productData []entity.ProductsData
	start = (start - 1) * limit
	where = append(where, "product_status != ?")
	bind = append(bind, Enum.ProductStatusDelete)
	if err := engine.Engine.Table(entity.ProductData{}).Select("*").
		Where(strings.Join(where, " AND "), bind...).Limit(limit, start).Desc("update_time").Find(&productData); err != nil {
		log.Error("get products Database Error", err)
		return nil, err
	}
	return productData, nil
}

func CountProductListByStoreId(engine *database.MysqlSession, where []string, bind []interface{}) (int64, error) {
	where = append(where, "product_status != ?")
	bind = append(bind, Enum.ProductStatusDelete)
	count, err := engine.Engine.Table(entity.ProductData{}).Select("count(*)").
		Where(strings.Join(where, " AND "), bind...).Count()
	if err != nil {
		log.Error("get products Database Error", err)
		return count, err
	}
	return count, nil
}

//計算收銀機商品數
func CountProductTotalByStoreId(engine *database.MysqlSession, storeId string) (int64, error) {
	sql := fmt.Sprintf("SELECT count(*) FROM product_data WHERE store_id = ? AND is_realtime = ? AND product_status != ?")
	result, err := engine.Engine.SQL(sql, storeId, 0, Enum.ProductStatusDelete).Count()
	if err != nil {
		log.Error("count products Database Error", err)
		return 0, err
	}
	return result, nil
}

//計算收銀機上架商品數
func CountProductSellByStoreId(engine *database.MysqlSession, storeId string) (int64, error) {
	sql := fmt.Sprintf("SELECT count(*) FROM product_data WHERE store_id = ? AND product_status = ? AND is_realtime = ?")
	result, err := engine.Engine.SQL(sql, storeId, Enum.ProductStatusSuccess, 0).Count()
	if err != nil {
		log.Error("count products Database Error", err)
		return 0, err
	}
	log.Debug("count", result)
	return result, nil
}

//計算收銀機下架商品數
func CountProductDownByStoreId(engine *database.MysqlSession, storeId string) (int64, error) {
	sql := fmt.Sprintf("SELECT count(*) FROM product_data WHERE store_id = ? AND product_status = ? AND is_realtime = ?")
	result, err := engine.Engine.SQL(sql, storeId, Enum.ProductStatusDown, 0).Count()
	if err != nil {
		log.Error("count products Database Error", err)
		return 0, err
	}
	return result, nil
}

//計算收銀機無庫存商品數
func CountProductStockByStoreId(engine *database.MysqlSession, storeId string) (int64, error) {
	sql := fmt.Sprintf("SELECT count(*) FROM product_data WHERE store_id = ? AND stock = ? AND is_realtime = ? AND product_status != ?")
	result, err := engine.Engine.SQL(sql, storeId, 0, 0, Enum.ProductStatusDelete).Count()
	if err != nil {
		log.Error("count products Database Error", err)
		return 0, err
	}
	return result, nil
}

func GetItemListCountByStoreId(engine *database.MysqlSession, where []string, bind []interface{}) (int64, error) {
	where = append(where, "product_status != ?")
	bind = append(bind, Enum.ProductStatusDelete)
	sql := fmt.Sprintf("SELECT count(*) as count FROM product_data"+
		" WHERE %s ", strings.Join(where, " AND "))
	result, err := engine.Engine.SQL(sql, bind...).Count()
	if err != nil {
		log.Error("count products Database Error", err)
		return 0, err
	}
	return result, nil
}

func GetProductsByStoreId(engine *database.MysqlSession, storeId string, isRealtime int, limit int, start int) ([]entity.ProductsData, error) {
	var productData []entity.ProductsData
	var err = engine.Engine.Table(Response.ProductSpec{}).
		Join("Left", entity.ProductSpecData{}, "product_data.product_id = product_spec_data.product_id").
		Select("product_id").Where("store_id = ?", storeId).
		And("product_status = ?", "SUCC").And("is_realtime = ?", isRealtime).
		GroupBy("product_id").OrderBy("product_id DESC").Limit(limit, start).Find(&productData)
	if err != nil {
		log.Error("get products Database Error", err)
		return nil, err
	}
	return productData, nil
}

func UpdateMinusProductStockByProductSpecId(engine *database.MysqlSession, productId, ProductSpecId string, qty int) error {
	sql1 := "UPDATE product_data SET stock = stock - ? WHERE product_id = ?"
	if _, err := engine.Session.Exec(sql1, qty, productId); err != nil {
		log.Error("update product stock Database Error", err)
		return err
	}
	sql2 := "UPDATE product_spec_data SET Quantity = Quantity - ? WHERE spec_id = ?"
	if _, err := engine.Session.Exec(sql2, qty, ProductSpecId); err != nil {
		log.Error("update product stock Database Error", err)
		return err
	}
	return nil
}

func UpdatePlusProductStockByProductSpecId(engine *database.MysqlSession, ProductId, ProductSpecId string, qty int) error {
	log.Debug("return product stock", ProductId, qty)
	sql1 := "UPDATE product_data SET stock = stock + ? WHERE product_id = ?"
	if _, err := engine.Session.Exec(sql1, qty, ProductId); err != nil {
		log.Error("update product stock Database Error", err)
		return err
	}

	sql2 := "UPDATE product_spec_data SET Quantity = Quantity + ? WHERE spec_id = ?"
	if _, err := engine.Session.Exec(sql2, qty, ProductSpecId); err != nil {
		log.Error("update product stock Database Error", err)
		return err
	}
	return nil
}

func GetProductListBySellerId(engine *database.MysqlSession, sellerId string, limit int, start int) ([]entity.ProductsData, error) {
	var productData []entity.ProductsData
	var err = engine.Engine.Table("products_data").
		Select("*").Where("seller_id = ?", sellerId).
		And("is_realtime = ?", "0").
		GroupBy("product_id").OrderBy("product_id DESC").Limit(limit, start).Find(&productData)
	if err != nil {
		log.Error("get products Database Error", err)
		return nil, err
	}
	return productData, nil
}

//更新商品資訊
func UpdateProductsData(engine *database.MysqlSession, Id string, Data entity.ProductData) error {
	Data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.ProductData{}).ID(Id).AllCols().Update(Data)
	if err != nil {
		return err
	}
	return nil
}

//更新商品資訊
func UpdateProductData(engine *database.MysqlSession, Id string, Data *entity.ProductData) error {
	Data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.ProductData{}).ID(Id).AllCols().Update(Data)
	if err != nil {
		return err
	}
	return nil
}

func CountRealTimeByStoreId(engine *database.MysqlSession, storeId, status string) int64 {
	result, err := engine.Engine.Table(entity.ProductData{}).And("store_id = ?", storeId).
		And("is_realtime = ?", 1).
		And("product_status = ?", status).Count()
	if err != nil {
		log.Error("count RealTime Database Error", err)
		return 0
	}
	return result
}

//更新商品狀態
func UpdateProductStatus(engine *database.MysqlSession, storeId string, status, oldStatus string) error {
	sql := fmt.Sprintf("UPDATE product_data SET product_status = ? WHERE store_id = ? AND product_status = ?")
	_, err := engine.Session.Exec(sql, status, storeId, oldStatus)
	if err != nil {
		log.Error("update product status Error", err)
		return err
	}
	return nil
}

func UpdateProductStatusByProductId(engine *database.MysqlSession, ProductId string, status string) error {
	sql := fmt.Sprintf("UPDATE product_data SET product_status = ? WHERE product_id = ?")
	_, err := engine.Session.Exec(sql, status, ProductId)
	if err != nil {
		log.Error("update product status Error", err)
		return err
	}
	return nil
}

func UpdateProductStock(engine *database.MysqlSession, ProductId string, stock int) error {
	log.Debug("plus product stock", ProductId, stock)
	sql := fmt.Sprintf("UPDATE product_data SET stock = ? WHERE product_id = ?")
	_, err := engine.Session.Exec(sql, stock, ProductId)
	if err != nil {
		log.Error("update product status Error", err)
		return err
	}
	return nil
}

//取出一般商品資料
func GetProductDataByProductId(engine *database.MysqlSession, productId string) (entity.ProductData, error) {
	var productData entity.ProductData
	_, err := engine.Engine.Table(entity.ProductData{}).
		Select("*").Where("product_id = ?", productId).Get(&productData)
	if err != nil {
		log.Error("get products Database Error", err)
		return productData, err
	}
	return productData, nil
}

//---- 圖片 ----

func InsertProductImage(engine *database.MysqlSession, data entity.ProductImagesData) error {
	data.CreateTime = time.Now()
	_, err := engine.Session.Table("product_images_data").Insert(&data)
	if err != nil {
		log.Error("Database Error", err)
		return err
	}
	return nil
}

func UpdateProductImageData(engine *database.MysqlSession, data entity.ProductImagesData) error {
	sql := fmt.Sprintf(" UPDATE product_data SET image = ?, image_seq=?, image_status = ?, image_seq = ? WHERE id = ?")
	_, err := engine.Session.SQL(sql).ID(data.Id).Update(&data)
	if err != nil {
		return err
	}
	return nil
}

func UpdateProductImageStatus(engine *database.MysqlSession, ProductId string) error {
	sql := fmt.Sprintf("UPDATE product_images_data SET image_status = ? WHERE product_id = ?")
	_, err := engine.Session.Exec(sql, Enum.ProductStatusDelete, ProductId)
	if err != nil {
		return err
	}
	return nil
}

//取出商品圖片
func GetProductImageByProductId(engine *database.MysqlSession, productId string) ([]entity.ProductImagesData, error) {
	var productImage []entity.ProductImagesData
	var err = engine.Engine.Table("product_images_data").
		Select("*").Where("product_id = ?", productId).
		And("image_status = ?", Enum.ProductStatusSuccess).
		Asc("image_seq").Find(&productImage)
	if err != nil {
		log.Error("Database Error", err)
		return nil, err
	}
	return productImage, nil
}
//取出商品圖片
func GetProductImageByImageAndProductId(engine *database.MysqlSession, image string, productId string) (entity.ProductImagesData, error) {
	var data entity.ProductImagesData
	_, err := engine.Engine.Table("product_images_data").
		Select("*").Where("product_id = ?", productId).
		And("image_status = ?", Enum.ProductStatusSuccess).
		And("image = ?", image).
		Get(&data)
	if err != nil {
		log.Error("Database Error", err)
		return data, err
	}
	return data, nil
}
//取出商品第一張圖片
func GetProductFirstImageByProductId(engine *database.MysqlSession, productId string) (entity.ProductImagesData, error) {
	var data entity.ProductImagesData
	_, err := engine.Engine.Table("product_images_data").
		Select("*").Where("product_id = ?", productId).
		And("image_status = ?", Enum.ProductStatusSuccess).
		Asc("image_seq").
		Get(&data)
	if err != nil {
		log.Error("Database Error", err)
		return data, err
	}
	return data, nil
}
//取出商品資料 BY ProductId
func GetProductsByProductIdByStoreId(engine *database.MysqlSession, productId, storeId string) (entity.ProductsData, error) {
	var productData entity.ProductsData
	_, err := engine.Engine.Table(entity.ProductData{}).Select("*").Where("product_id = ?", productId).And("store_id = ?", storeId).Get(&productData)
	if err != nil {
		log.Error("get products Database Error", err)
		return productData, err
	}
	return productData, nil
}

func GetProductsByProductId(engine *database.MysqlSession, productId string) (entity.ProductsData, error) {
	var productData entity.ProductsData
	_, err := engine.Engine.Table(entity.ProductData{}).Select("*").Where("product_id = ?", productId).Get(&productData)
	if err != nil {
		log.Error("get products Database Error", err)
		return productData, err
	}
	return productData, nil
}

//取出商品資料 BY ProductId
func GetProductByProductId(engine *database.MysqlSession, productId string) (entity.ProductData, error) {
	var productData entity.ProductData
	_, err := engine.Engine.Table(entity.ProductData{}).Select("*").Where("product_id = ?", productId).Get(&productData)
	if err != nil {
		log.Error("get products Database Error", err)
		return productData, err
	}
	return productData, nil
}

//取商品資料 BY TinyUrl
func GetProductByTinyUrl(engine *database.MysqlSession, tiny string) (entity.ProductsData, error) {
	var productData entity.ProductsData
	_, err := engine.Engine.Table(entity.ProductData{}).Select("*").Where("tiny_url = ?", tiny).Get(&productData)
	if err != nil {
		log.Error("get products Database Error", err)
		return productData, err
	}
	return productData, nil
}

//取出過期訂單
func GetRealtimeProductExpire(engine *database.MysqlSession, date string) ([]entity.ProductsData, error) {
	var data []entity.ProductsData
	sql := fmt.Sprintf("SELECT * FROM product_data WHERE expire_date <= ? AND product_status = ? AND is_realtime = ? ORDER BY create_time ASC")
	err := engine.Engine.SQL(sql, date, Enum.ProductStatusSuccess, 1).Find(&data)
	if err != nil {
		log.Error("Select Appropriation Order Database Error", err)
		return data, err
	}
	return data, nil
}

//取商品資料
func GetProductAll(engine *database.MysqlSession) ([]entity.ProductsData, error) {
	var productData []entity.ProductsData
	if err := engine.Engine.Table(entity.ProductData{}).Select("*").Find(&productData); err != nil {
		log.Error("get products Database Error", err)
		return productData, err
	}
	return productData, nil
}

//取商品資料
func SearchProducts(engine *database.MysqlSession, params Request.ErpSearchProductRequest) ([]entity.ProductData, error) {
	log.Debug("sss", params)
	where, bind := ComposeSearchProductsParams(engine, params)
	var data []entity.ProductData
	if err := engine.Engine.Table(entity.ProductData{}).Select("*").Where(strings.Join(where, " AND "), bind...).Find(&data); err != nil {
		log.Error("get products Database Error", err)
		return data, err
	}
	return data, nil
}

func ComposeSearchProductsParams(engine *database.MysqlSession, params Request.ErpSearchProductRequest) ([]string, []interface{}) {
	var where []string
	var bind []interface{}

	if len(params.UserAccount) != 0 {
		data, _ := member.GetMemberAndStoreByAccount(engine, params.UserAccount)
		if len(data) != 0 {
			var orWhere []string
			var orBind []interface{}
			for _, v := range data {
				orWhere = append(orWhere, "store_id = ?")
				orBind = append(orBind, v.StoreData.StoreId)
			}
			where = append(where, strings.Join(orWhere, " OR "))
			bind = append(bind, orBind...)
		}
	}
	if len(params.UserId) != 0 {
		data, _ := member.GetMemberAndStoreByTerminalId(engine, params.UserId)
		if len(data) != 0 {
			var orWhere []string
			var orBind []interface{}
			for _, v := range data {
				orWhere = append(orWhere, "store_id = ?")
				orBind = append(orBind, v.StoreData.StoreId)
			}
			where = append(where, strings.Join(orWhere, " OR "))
			bind = append(bind, orBind...)
		}
	}
	if len(params.StoreName) != 0 {
		data, _ := member.GetMemberAndStoreByStoreName(engine, params.StoreName)
		if len(data) != 0 {
			var orWhere []string
			var orBind []interface{}
			for _, v := range data {
				orWhere = append(orWhere, "store_id = ?")
				orBind = append(orBind, v.StoreData.StoreId)
			}
			where = append(where, strings.Join(orWhere, " OR "))
			bind = append(bind, orBind...)
		}
	}
	if len(params.Tab) != 0 {
		if params.Tab == "Product" {
			where = append(where, "is_realtime = ?")
			bind = append(bind, 0)
		} else {
			where = append(where, "is_realtime = ?")
			bind = append(bind, 1)
		}
	}

	if len(params.ProductStatus) != 0 {
		where = append(where, "product_status = ?")
		bind = append(bind, params.ProductStatus)
	}
	if len(params.ProductName) != 0 {
		where = append(where, "product_name = ?")
		bind = append(bind, params.ProductName)
	}
	if len(params.ProductId) != 0 {
		where = append(where, "product_id = ?")
		bind = append(bind, params.ProductId)
	}
	if len(params.ShipMode) != 0 {
		where = append(where, "shipping LIKE ?")
		bind = append(bind, "%"+params.ShipMode+"%")
	}
	if len(params.PaymentMode) != 0 {
		where = append(where, "pay_way LIKE ?")
		bind = append(bind, "%"+params.PaymentMode+"%")
	}
	if len(params.CreateDate) != 0 {
		date := strings.Split(params.CreateDate, "-")
		where = append(where, "create_time BETWEEN ? AND ?")
		bind = append(bind, date[0])
		bind = append(bind, date[1])
	}
	if len(params.UpdateDate) != 0 {
		date := strings.Split(params.UpdateDate, "-")
		where = append(where, "update_time BETWEEN ? AND ?")
		bind = append(bind, date[0])
		bind = append(bind, date[1])
	}
	if len(params.Amount) != 0 {
		price := strings.Split(params.Amount, "-")
		where = append(where, "price BETWEEN ? AND ?")
		bind = append(bind, price[0])
		bind = append(bind, price[1])
	}

	return where, bind
}

func GetProductById(engine *database.MysqlSession, productId string) (entity.ProductData, error) {
	var data entity.ProductData
	if _, err := engine.Engine.Table(entity.ProductData{}).Select("*").
		Where("product_id = ?", productId).Get(&data); err != nil {
		log.Error("get RealTim Database Error", err)
		return data, err
	}
	return data, nil
}

type StockData struct {
	ProductId string
	Stock     int64
}

func SumProductsStock(engine *database.MysqlSession) error {
	var data []StockData
	sql := "SELECT product_id, SUM(quantity) as stock FROM product_spec_data GROUP BY product_id"
	err := engine.Engine.SQL(sql).Find(&data)
	if err != nil {
		log.Error("update product status Error", err)
		return err
	}
	for _, v := range data {
		product, _ := GetProductById(engine, v.ProductId)
		log.Debug("product stock", product, v.Stock, product.Stock)
		product.Stock = v.Stock
		err := UpdateProductsData(engine, product.ProductId, product)
		if err != nil {
			log.Debug("Update Product Error", err)
		}
	}
	return nil
}

func CountProductAndInstatOrderByStore(engine *database.MysqlSession, StoreId string) (int64, int64, error) {
	OnlineProducts, err := engine.Engine.Table(entity.ProductData{}).Where("store_id = ? AND product_status = ?", StoreId, Enum.ProductStatusSuccess).Count()
	if err != nil {
		return 0, 0, err
	}
	InstanceProducts, err := engine.Engine.Table(entity.ProductData{}).Where("store_id = ? AND product_status = ? AND is_realtime = ?", StoreId, Enum.ProductStatusSuccess, true).Count()
	if err != nil {
		return OnlineProducts, 0, err
	}
	return OnlineProducts, InstanceProducts, nil
}


func CountProductByStoreId(engine *database.MysqlSession, StoreId string) (int64, error) {
	count, err := engine.Engine.Table(entity.ProductData{}).
		Where("store_id = ? AND product_status = ?", StoreId, Enum.ProductStatusSuccess).Count()
	if err != nil {
		return count, err
	}
	return count, nil
}

func CloseAllProduct(engine *database.MysqlSession) error {
	sql := fmt.Sprintf("UPDATE product_data SET product_status = ?")
	_, err := engine.Session.Exec(sql, Enum.ProductStatusDelete)
	if err != nil {
		log.Error("UpdateOrderMessageBoardData Error", err)
		return err
	}
	return nil
}