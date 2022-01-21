package model

import (
	"api/services/Enum"
	"api/services/Service/History"
	"api/services/Service/Product"
	"api/services/Service/StoreService"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Store"
	"api/services/dao/member"
	"api/services/dao/product"
	"api/services/database"
	"api/services/entity"
	"api/services/util/images"
	"api/services/util/log"
	"api/services/util/qrcode"
	"api/services/util/tools"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

//建立商品資訊
func HandleCreateProduct(storeData entity.StoreDataResp, params *Request.NewProductParams, filename []string) (entity.ProductData, error) {
	var Entity entity.ProductData
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := engine.Session.Begin(); err != nil {
		log.Error("Set Begin Error", err)
		return Entity, fmt.Errorf("1001001")
	}
	var productId = tools.GeneratorProductId()
	uri := fmt.Sprintf("/product/%s", productId)
	tiny, err := Product.GeneratorShortUrl(engine, uri)
	if err != nil {
		return Entity, fmt.Errorf("1001001")
	}
	if err := insertProduct(engine, storeData, productId, tiny, params); err != nil {
		engine.Session.Rollback()
		return Entity, fmt.Errorf("1001001")
	}

	if len(params.ProductSpecList) == 0 {
		var ProductSpecId = tools.GeneratorProductSpecId(productId, 1)
		price, _ := strconv.Atoi(params.Price)
		if price == 0 {
			return Entity, fmt.Errorf("1007008")
		}
		if err := insertProductSpec(engine, productId, ProductSpecId, "", params.ProductQty, price); err != nil {
			engine.Session.Rollback()
			return Entity, fmt.Errorf("1001001")
		}
		if err := product.UpdateProductStock(engine, productId, params.ProductQty); err != nil {
			engine.Session.Rollback()
			return Entity, fmt.Errorf("1001001")
		}
	} else {
		stock := 0
		for k, v := range params.ProductSpecList {
			var ProductSpecId = tools.GeneratorProductSpecId(productId, k+1)
			price, _ := strconv.Atoi(params.Price)
			if price == 0 {
				return Entity, fmt.Errorf("1007008")
			}
			stock += v.Quantity
			if err := insertProductSpec(engine, productId, ProductSpecId, v.ProductSpec, v.Quantity, price); err != nil {
				engine.Session.Rollback()
				return Entity, fmt.Errorf("1001001")
			}
		}
		if err := product.UpdateProductStock(engine, productId, stock); err != nil {
			engine.Session.Rollback()
			return Entity, fmt.Errorf("1001001")
		}
	}
	if err := CreateProductImagesData(engine, productId, filename); err != nil {
		engine.Session.Rollback()
		return Entity, fmt.Errorf("1001001")
	}
	if err := engine.Session.Commit(); err != nil {
		return Entity, fmt.Errorf("1001001")
	}
	productData, err := product.GetProductDataByProductId(engine, productId)

	if err != nil {
		return Entity, fmt.Errorf("1001001")
	}
	History.GenerateNewProductLog(productData, storeData.UserId)
	return productData, nil
}

//建立圖片資料
func CreateProductImagesData(engine *database.MysqlSession, productId string, filename []string) error {
	for k, value := range filename {
		//URL decode
		img, err := url.QueryUnescape(value)
		if err != nil {
			log.Error("Images URL Decode Error", value)
			return err
		}
		file := value
		if strings.Contains(img, "/") {
			split := strings.Split(img, "/")
			file = split[len(split)-1]
		}
		if err := createProductImages(engine, file, productId, k); err != nil {
			return err
		}
	}
	return nil
}

//變更圖片資料
func CheckProductImagesData(engine *database.MysqlSession, productId string, filename []string) error {
	for k, value := range filename {
		img, err := url.QueryUnescape(value)
		if err != nil {
			log.Error("Images URL Decode Error", value)
			return err
		}
		file := value
		if strings.Contains(img, "/") {
			split := strings.Split(img, "/")
			file = split[len(split)-1]
		}
		image, err := product.GetProductImageByImageAndProductId(engine, file, productId)
		if err != nil {
			return err
		}
		if len(image.Image) == 0 {
			if err := createProductImages(engine, value, productId, k); err != nil {
				return err
			}
		} else {
			image.ImageSeq = k
			image.ImageStatus = Enum.ProductStatusSuccess
			if err := product.UpdateProductImageData(engine, image); err != nil {
				return err
			}
		}
	}
	return nil
}

//商品編輯
func HandleEditProduct(storeData entity.StoreDataResp, params *Request.EditProductParams, filename []string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := engine.Session.Begin(); err != nil {
		log.Error("Set Begin Error", err)
		return fmt.Errorf("1001001")
	}
	data, err := product.GetProductDataByProductId(engine, params.ProductId)
	if err != nil {
		return fmt.Errorf("1001010")
	}
	//為了存取原始資料
	tempOldData := data
	if err := updateProduct(engine, storeData, params, &data); err != nil {
		engine.Session.Rollback()
		return fmt.Errorf("1001001")
	}
	History.GenerateEditProductLog(params, tempOldData, storeData.UserId)
	if err := product.UpdateProductSpecDelete(engine, data.ProductId); err != nil {
		engine.Session.Rollback()
		return fmt.Errorf("1001001")
	}
	//判斷是否有規格
	if len(params.ProductSpecList) == 0 {
		//沒規格
		spec, err := product.GetProductSpecByProductId(engine, data.ProductId)
		if err != nil {
			engine.Session.Rollback()
			return fmt.Errorf("1001001")
		}
		for k, v := range spec {
			if k == 0 {
				price, _ := strconv.Atoi(params.Price)
				if price == 0 {
					engine.Session.Rollback()
					return fmt.Errorf("1007008")
				}
				if params.ProductQty == 0 {
					engine.Session.Rollback()
					return fmt.Errorf("1007007")
				}
				v.SpecName = ""
				v.Quantity = int64(params.ProductQty)
				v.SpecPrice = int64(price)
				v.SpecStatus = Enum.ProductStatusSuccess
				if err := product.UpdateProductSpec(engine, v.SpecId, v); err != nil {
					engine.Session.Rollback()
					return fmt.Errorf("1001001")
				}
			}
		}
		if err := product.UpdateProductStock(engine, data.ProductId, params.ProductQty); err != nil {
			return err
		}
	} else {
		//有規格
		stock := 0
		for k, v := range params.ProductSpecList {
			stock += v.Quantity
			//判斷Params SpecID是否有值
			if len(v.ProductSpecId) == 0 {
				ProductSpecId := tools.GeneratorProductSpecId(data.ProductId, k+1)
				spec, err := product.GetProductSpecByProductSpecId(engine, ProductSpecId)
				if err != nil {
					engine.Session.Rollback()
					return fmt.Errorf("1001001")
				}
				//判斷SpecID是否存在
				if ProductSpecId == spec.SpecId {
					if v.Quantity == 0 {
						engine.Session.Rollback()
						return fmt.Errorf("1007007")
					}
					price, _ := strconv.Atoi(params.Price)
					if price == 0 {
						engine.Session.Rollback()
						return fmt.Errorf("1007008")
					}
					spec.SpecName = v.ProductSpec
					spec.Quantity = int64(v.Quantity)
					spec.SpecPrice = int64(price)
					spec.SpecStatus = Enum.ProductStatusSuccess
					if err := product.UpdateProductSpec(engine, spec.SpecId, spec); err != nil {
						engine.Session.Rollback()
						return fmt.Errorf("1001001")
					}
				} else {
					if v.Quantity == 0 {
						engine.Session.Rollback()
						return fmt.Errorf("1007007")
					}
					price, _ := strconv.Atoi(params.Price)
					if price == 0 {
						engine.Session.Rollback()
						return fmt.Errorf("1007008")
					}
					if err := insertProductSpec(engine, data.ProductId, ProductSpecId, v.ProductSpec, v.Quantity, price); err != nil {
						engine.Session.Rollback()
						return err
					}
				}
			} else {
				spec, err := product.GetProductSpecByProductSpecId(engine, v.ProductSpecId)
				if err != nil {
					engine.Session.Rollback()
					return fmt.Errorf("1001001")
				}
				if v.Quantity == 0 {
					engine.Session.Rollback()
					return fmt.Errorf("1007007")
				}
				price, _ := strconv.Atoi(params.Price)
				spec.SpecPrice = int64(price)
				if price == 0 {
					engine.Session.Rollback()
					return fmt.Errorf("1007008")
				}
				spec.SpecName = v.ProductSpec
				spec.Quantity = int64(v.Quantity)
				spec.SpecStatus = Enum.ProductStatusSuccess
				if err := product.UpdateProductSpec(engine, spec.SpecId, spec); err != nil {
					engine.Session.Rollback()
					return fmt.Errorf("1001001")
				}
			}
		}
		//商品資料寫入總庫存數
		if err := product.UpdateProductStock(engine, data.ProductId, stock); err != nil {
			engine.Session.Rollback()
			return err
		}
	}
	if err := product.UpdateProductImageStatus(engine, data.ProductId); err != nil {
		engine.Session.Rollback()
		return fmt.Errorf("1001001")
	}
	if err := CheckProductImagesData(engine, data.ProductId, filename); err != nil {
		engine.Session.Rollback()
		return fmt.Errorf("1001001")
	}
	if err := engine.Session.Commit(); err != nil {
		return fmt.Errorf("1001001")
	}
	return nil
}

//新增商品
func insertProduct(engine *database.MysqlSession, storeData entity.StoreDataResp, productId string, TinyUrl string, params *Request.NewProductParams) error {
	var data entity.ProductData
	data.ProductId = productId
	data.ProductName = params.ProductName
	price, _ := strconv.Atoi(params.Price)
	data.Price = int64(price)
	shipping, _ := tools.JsonEncode(rearrangeShipMode(params.ShippingList))
	data.Shipping = shipping
	data.ShipMerge = params.ShipMerge
	payWay, _ := tools.JsonEncode(rearrangesPayWayMode(params.PayWayList, params.ShippingList))
	data.PayWay = payWay
	if StoreStatusIsClose(storeData) {
		data.ProductStatus = Enum.ProductStatusPending
	} else {
		data.ProductStatus = Enum.ProductStatusSuccess
	}
	switch params.LimitKey {
	case Enum.ProductLimitLeast:
		data.LimitKey = Enum.ProductLimitLeast
		data.LimitQty = int64(params.LimitQty)
	case Enum.ProductLimitMost:
		data.LimitKey = Enum.ProductLimitMost
		data.LimitQty = int64(params.LimitQty)
	default:
		data.LimitKey = Enum.ProductLimitNone
	}
	data.StoreId = storeData.StoreId
	data.IsSpec = params.IsSpec
	data.TinyUrl = TinyUrl
	data.IsRealtime = params.IsRealtime
	data.FormUrl = params.FormUrl
	data.IsFreeShip = true
	if params.IsRealtime == 1 {
		data.ExpireDate = time.Now().Add(72 * time.Hour)
	}
	if _, err := product.InsertProduct(engine, data); err != nil {
		return err
	}
	return nil
}

//更新商品資料
func updateProduct(engine *database.MysqlSession, storeData entity.StoreDataResp, params *Request.EditProductParams, data *entity.ProductData) error {

	data.ProductName = params.ProductName
	price, _ := strconv.Atoi(params.Price)
	data.Price = int64(price)
	shipping, _ := json.Marshal(rearrangeShipMode(params.ShippingList))
	data.Shipping = string(shipping)
	data.ShipMerge = params.ShipMerge
	payWay, _ := json.Marshal(rearrangesPayWayMode(params.PayWayList, params.ShippingList))
	data.PayWay = string(payWay)
	if params.StatusDown == 1 {
		data.ProductStatus = Enum.ProductStatusDown
	} else if params.StatusDelete == 1 {
		data.ProductStatus = Enum.ProductStatusDelete
	} else {
		if StoreStatusIsClose(storeData) {
			data.ProductStatus = Enum.ProductStatusPending
		} else {
			data.ProductStatus = Enum.ProductStatusSuccess
		}
	}
	switch params.LimitKey {
	case Enum.ProductLimitLeast:
		data.LimitKey = Enum.ProductLimitLeast
		data.LimitQty = int64(params.LimitQty)
	case Enum.ProductLimitMost:
		data.LimitKey = Enum.ProductLimitMost
		data.LimitQty = int64(params.LimitQty)
	default:
		data.LimitKey = Enum.ProductLimitNone
	}
	//data.IsFreeShip = params.IsFreeShip
	data.IsSpec = params.IsSpec
	data.FormUrl = params.FormUrl
	err := product.UpdateProductData(engine, data.ProductId, data)
	if err != nil {
		return err
	}

	return nil
}

//新增規格
func insertProductSpec(engine *database.MysqlSession, productId string, productSpecId string, specName string, qty int, price int) error {
	var Entity entity.ProductSpecData
	Entity.ProductId = productId
	Entity.SpecId = productSpecId
	Entity.SpecName = specName
	Entity.SpecPrice = int64(price)
	Entity.Quantity = int64(qty)
	Entity.SpecStatus = Enum.ProductStatusSuccess
	_, err := product.InsertProductSpec(engine, Entity)
	if err != nil {
		return err
	}
	return nil
}

//重組運送方式
func rearrangeShipMode(shipping []Request.NewShipping) []entity.NewShippingMode {
	var ShippingMode []entity.NewShippingMode
	for _, val := range shipping {
		var data entity.NewShippingMode
		data.Type = val.ShipType
		data.Text = Enum.Shipping[val.ShipType]
		data.Price = val.ShipFee
		data.Remark = val.ShipRemark
		ShippingMode = append(ShippingMode, data)
	}
	return ShippingMode
}

func checkShipIsCvs(shipping []Request.NewShipping) bool {
	for _, v := range shipping {
		if tools.InArray([]string{Enum.CVS_HI_LIFE, Enum.CVS_FAMILY, Enum.CVS_OK_MART, Enum.CVS_7_ELEVEN}, v.ShipType) {
			return true
		}
	}
	return false
}

func checkShippingIsCvs(shipping []Response.ShippingMode) bool {
	for _, v := range shipping {
		if tools.InArray([]string{Enum.CVS_HI_LIFE, Enum.CVS_FAMILY, Enum.CVS_OK_MART, Enum.CVS_7_ELEVEN}, v.Type) {
			return true
		}
	}
	return false
}

//重組付款方式
func rearrangesPayWayMode(payway []string, shipping []Request.NewShipping) []Response.PayWayMode {
	var paywayMode []Response.PayWayMode
	for _, val := range payway {
		if val == Enum.CvsPay {
			if checkShipIsCvs(shipping) {
				mode := Response.PayWayMode{
					Type: val,
					Text: Enum.PayWay[val],
				}
				paywayMode = append(paywayMode, mode)
			}
		} else {
			mode := Response.PayWayMode{
				Type: val,
				Text: Enum.PayWay[val],
			}
			paywayMode = append(paywayMode, mode)
		}
	}
	return paywayMode
}

//重組付款方式
func rearrangesPayWay(payway []string) []Response.PayWayMode {
	var paywayMode []Response.PayWayMode
	for _, val := range payway {

		mode := Response.PayWayMode{
			Type: val,
			Text: Enum.PayWay[val],
		}
		paywayMode = append(paywayMode, mode)
	}
	return paywayMode
}

//重組付款方式
func rearrangesPayment(payway []Response.PayWayMode, shipping []Response.ShippingMode) []Response.PayWayMode {
	var payment []Response.PayWayMode
	for _, val := range payway {
		if val.Type == Enum.CvsPay {
			if checkShippingIsCvs(shipping) {
				payment = append(payment, val)
			}
		} else {
			payment = append(payment, val)
		}
	}
	return payment
}

//建立商品圖
func createProductImages(engine *database.MysqlSession, filename string, productId string, seq int) error {
	var Ent entity.ProductImagesData
	Ent.ProductId = productId
	Ent.Image = filename
	Ent.ImageSeq = seq
	Ent.ImageStatus = Enum.ProductStatusSuccess
	err := product.InsertProductImage(engine, Ent)
	if err != nil {
		return err
	}
	return nil
}

//取一般商品資料
func GetProductDataByTinyUrl(tinyUrl string) (Response.ProductResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.ProductResponse
	productData, err := product.GetProductByTinyUrl(engine, tinyUrl)
	if err != nil {
		return resp, err
	}
	if productData.ProductStatus == Enum.ProductStatusPaid {
		return resp, fmt.Errorf("1006001")
	}
	if productData.ProductStatus == Enum.ProductStatusCancel {
		return resp, fmt.Errorf("1006002")
	}
	if productData.ProductStatus == Enum.ProductStatusDown || productData.ProductStatus == Enum.ProductStatusPending {
		return resp, fmt.Errorf("1006003")
	}
	if productData.ProductStatus == Enum.ProductStatusOverdue {
		return resp, fmt.Errorf("1006004")
	}
	if productData.ProductStatus == Enum.ProductStatusDelete {
		return resp, fmt.Errorf("1006005")
	}
	if len(productData.ProductStatus) == 0 {
		return resp, fmt.Errorf("1006005")
	}
	resp = composeProductData(engine, productData, true, false)
	return resp, nil
}

//取一般商品資料
func GetProductDataByProductId(ProductId string, all bool, edit bool) (Response.ProductResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.ProductResponse
	productData, err := product.GetProductsByProductId(engine, ProductId)
	if err != nil {
		return resp, err
	}
	if productData.ProductStatus == Enum.ProductStatusPaid {
		return resp, fmt.Errorf("1006001")
	}
	if productData.ProductStatus == Enum.ProductStatusCancel {
		return resp, fmt.Errorf("1006002")
	}
	if productData.ProductStatus == Enum.ProductStatusDown || productData.ProductStatus == Enum.ProductStatusPending {
		return resp, fmt.Errorf("1006003")
	}
	if productData.ProductStatus == Enum.ProductStatusOverdue {
		return resp, fmt.Errorf("1006004")
	}
	if productData.ProductStatus == Enum.ProductStatusDelete {
		return resp, fmt.Errorf("1006005")
	}
	if len(productData.ProductStatus) == 0 {
		return resp, fmt.Errorf("1006005")
	}
	resp = composeProductData(engine, productData, all, edit)
	return resp, nil
}

//取得編輯商品資料
func GetEditProductDataByProductId(storeData entity.StoreDataResp, ProductId string, all bool, edit bool) (Response.ProductResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.ProductResponse

	productData, err := product.GetProductsByProductIdByStoreId(engine, ProductId, storeData.StoreId)
	if err != nil {
		return resp, err
	}
	resp = composeProductData(engine, productData, all, edit)
	return resp, nil
}

//判斷收銀機是否關閉
func StoreStatusIsClose(storeData entity.StoreDataResp) bool {
	if storeData.StoreStatus != Enum.StoreStatusSuccess {
		return true
	}
	return false
}

//判斷收銀機是否可編輯商品
func ProductCheckStore(storeData entity.StoreDataResp, ProductId string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if storeData.StoreStatus == Enum.StoreStatusSuspend {
		return fmt.Errorf("1007010")
	}
	productData, err := product.GetProductsByProductId(engine, ProductId)
	if err != nil {
		return err
	}
	if storeData.StoreId != productData.StoreId {
		return fmt.Errorf("1006005")
	}
	return nil
}

//重組商品資料
func composeProductData(engine *database.MysqlSession, productData entity.ProductsData, all bool, edit bool) Response.ProductResponse {
	var data Response.ProductResponse
	var payWayList []Response.PayWayMode

	productImages, err := RearrangeProductImage(engine, productData.ProductId, edit, productData.IsRealtime)
	if err != nil {
		log.Error("Get product spec Error")
	}
	storeData, _ := Store.GetStoreDataByStoreId(engine, productData.StoreId)
	//取出免運設定
	setting, _ := StoreService.GetStoreFreeShipping(engine, productData.StoreId)
	data.ProductId = productData.ProductId
	data.ProductName = productData.ProductName
	data.ProductImageList = productImages
	data.Price = productData.Price
	data.TotalStock = productData.Stock
	data.IsSpec = productData.IsSpec
	data.IsRealTime = productData.IsRealtime
	data.StoreId = productData.StoreId
	data.ShipMerge = productData.ShipMerge
	data.StoreName = storeData.StoreName
	data.StorePicture = storeData.StorePicture
	data.LimitKey = productData.LimitKey
	data.LimitQty = int(productData.LimitQty)
	if productData.IsFreeShip == true && setting.FreeShipKey != Enum.FreeShipNone && data.ShipMerge == 1 {
		data.FreeShipKey = setting.FreeShipKey
		data.FreeShip = setting.FreeShip
	} else {
		data.FreeShipKey = Enum.FreeShipNone
		data.FreeShip = 0
	}

	SsmData, err := Store.GetStoreSocialMediaDataByStoreId(engine, storeData.StoreId)
	if err != nil {
		log.Error("Get store Social Media Data Error", err)
	}
	data.StoreSocialMedia = SsmData.GetStoreSocialMediaInfo()
	data.ExpireDate = ""
	if !productData.ExpireDate.IsZero() {
		data.ExpireDate = productData.ExpireDate.Format("2006-01-02 15:04:05")
	}
	data.QRCode = qrcode.GetQrcodeImageLink(productData.ProductId)
	data.ShortUrl = qrcode.GetTinyUrl(productData.TinyUrl)
	data.Status = productData.ProductStatus
	data.FormUrl = productData.FormUrl
	shipList := Product.ModifyShipList(productData.Shipping, setting)
	data.ShippingList = shipList
	if setting.SelfDelivery && data.ShipMerge == 1 {
		for _, ship := range data.ShippingList {
			if ship.Type == Enum.SELF_DELIVERY {
				data.SelfDeliveryKey = setting.SelfDeliveryKey
				data.SelfDeliveryFree = setting.SelfDeliveryFree
			}
		}
	} else {
		data.SelfDeliveryKey = Enum.FreeShipNone
		data.SelfDeliveryFree = 0
	}
	if err := tools.JsonDecode([]byte(productData.PayWay), &payWayList); err != nil {
		log.Error("Json Decode PayWay Error")
	}
	productPaymentList := rearrangesPayment(payWayList, data.ShippingList)
	data.ProductPayWay = productPaymentList
	productSpec, err := RearrangeProductSpec(engine, productData.ProductId, all)
	if err != nil {
		log.Error("Get product spec Error")
	}
	data.ProductSpecList = productSpec
	data.UpdateTime = productData.UpdateTime.Format("2006/01/02 15:04")
	return data
}

//取出商品圖片
func RearrangeProductImage(engine *database.MysqlSession, ProductId string, edit bool, IsRealtime int) ([]string, error) {
	var ProductImages []string
	productImage, err := product.GetProductImageByProductId(engine, ProductId)
	if err != nil {
		return nil, err
	}
	for _, elem := range productImage {
		var img = images.GetImageUrl(elem.Image)
		ProductImages = append(ProductImages, img)
	}
	if !edit && IsRealtime != 1 {
		ProductImages = append(ProductImages, fmt.Sprintf("/static/images/qrcode/%s.jpg", ProductId))
	}
	return ProductImages, nil
}

//取出商品規格
func RearrangeProductSpec(engine *database.MysqlSession, ProductId string, all bool) ([]Response.ProductSpec, error) {
	var productSpec []Response.ProductSpec
	ProductSpecData, err := product.GetProductSpecByProductId(engine, ProductId)
	if err != nil {
		return productSpec, err
	}
	for _, elem := range ProductSpecData {
		if elem.Quantity > 0 || all {
			spec := Response.ProductSpec{}
			spec.ProductSpecId = elem.SpecId
			spec.Spec = elem.SpecName
			spec.Price = elem.SpecPrice
			spec.Quantity = elem.Quantity
			productSpec = append(productSpec, spec)
		}
	}
	return productSpec, nil
}

//組即時帳單搜尋條件
func setRealTimeParams(sellerId string, params Request.RealTimesRequest) ([]string, []interface{}, string) {
	var sql []string
	var bind []interface{}
	var orderBy string
	sql = append(sql, "store_id = ?")
	bind = append(bind, sellerId)
	sql = append(sql, "is_realtime = ?")
	bind = append(bind, 1)

	switch params.Tab {
	case "ValidBill":
		sql = append(sql, "product_status = ?")
		bind = append(bind, Enum.ProductStatusSuccess)
	case "CancelBill":
		sql = append(sql, "product_status = ?")
		bind = append(bind, Enum.ProductStatusCancel)
	case "OverdueBill":
		sql = append(sql, "product_status = ?")
		bind = append(bind, Enum.ProductStatusOverdue)
	}
	orderBy = "ORDER BY create_time DESC"
	return sql, bind, orderBy
}

//取得即時帳單列表
func GetStoreRealTimeList(storeId string, params Request.RealTimesRequest) (Response.StoreRealTimesResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.StoreRealTimesResponse
	var storeProduct []Response.StoreProduct
	sql, bind, _ := setRealTimeParams(storeId, params)
	count, err := product.GetItemListCountByStoreId(engine, sql, bind)
	productsData, err := product.GetProductListByStoreId(engine, sql, bind, params.Limit, params.Start)
	if err != nil {
		log.Error("Get Product List Error", err)
		return resp, fmt.Errorf("系統錯誤！")
	}
	storeProduct, err = composeProduct(engine, productsData)
	if err != nil {
		log.Error("compose Product Error", err)
		return resp, fmt.Errorf("系統錯誤！")
	}
	resp.ProductCount = count
	resp.Tabs.CancelBill = product.CountRealTimeByStoreId(engine, storeId, Enum.ProductStatusCancel)
	resp.Tabs.OverdueBill = product.CountRealTimeByStoreId(engine, storeId, Enum.ProductStatusOverdue)
	resp.Tabs.ValidBill = product.CountRealTimeByStoreId(engine, storeId, Enum.ProductStatusSuccess)
	resp.ProductList = storeProduct
	return resp, nil
}

//取出收銀機商品列表
func GetStoreProductList(storeId string, request Request.ProductList, isRealtime int) (Response.StoreProductsResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.StoreProductsResponse
	var storeProduct []Response.StoreProduct
	storeData, err := Store.GetStoreDataByStoreId(engine, storeId)
	if err != nil {
		return resp, fmt.Errorf("系統錯誤！")
	}
	if len(storeData.StoreId) == 0 {
		return resp, fmt.Errorf("無此收銀機")
	}
	setting, err := StoreService.GetStoreFreeShipping(engine, storeData.StoreId)
	if err != nil {
		return resp, fmt.Errorf("系統錯誤！")
	}
	where, bind, _ := setParams(storeId, isRealtime, request)
	count, err := product.CountProductListByStoreId(engine, where, bind)
	productsData, err := product.GetProductListByStoreId(engine, where, bind, request.Length, request.Page)
	if err != nil {
		log.Error("Get Product List Error", err)
		return resp, fmt.Errorf("系統錯誤！")
	}
	for _, v := range productsData {
		var storeProductSpec []Response.StoreProductSpec
		var productShipList []Response.ShippingMode
		var productPayWayList []Response.PayWayMode
		productData, err := product.GetProductSpecByProductId(engine, v.ProductId)
		if err != nil {
			return resp, err
		}
		for _, elem := range productData {
			spec := Response.StoreProductSpec{}
			spec.ProductSpecId = elem.SpecId
			spec.Spec = elem.SpecName
			spec.Price = elem.SpecPrice
			spec.Quantity = elem.Quantity
			storeProductSpec = append(storeProductSpec, spec)
		}
		productImage, err := product.GetProductImageByProductId(engine, v.ProductId)
		if err != nil {
			return resp, err
		}
		var imgs []string
		for _, elem := range productImage {
			imgs = append(imgs, images.GetImageUrl(elem.Image))
		}
		if len(productData) != 0 {
			rep := Response.StoreProduct{}
			rep.ProductId = v.ProductId
			rep.ProductName = v.ProductName
			rep.ProductImage = imgs
			rep.ProductIsSpec = v.IsSpec
			rep.ProductSpecList = storeProductSpec
			rep.ProductPrice = v.Price
			rep.ProductStatus = v.ProductStatus
			rep.ProductShipMerge = v.ShipMerge
			rep.TotalStock = v.Stock
			rep.FormUrl = v.FormUrl
			rep.LimitKey = v.LimitKey
			rep.LimitQty = int(v.LimitQty)
			if v.IsFreeShip == true && setting.FreeShipKey != Enum.FreeShipNone && v.ShipMerge == 1 {
				rep.FreeShipKey = setting.FreeShipKey
				rep.FreeShip = setting.FreeShip
			} else {
				rep.FreeShipKey = Enum.FreeShipNone
				rep.FreeShip = 0
			}
			_ = json.Unmarshal([]byte(v.Shipping), &productShipList)
			rep.ProductShipList = productShipList
			if setting.SelfDelivery && v.ShipMerge == 1 {
				for _, ship := range productShipList {
					if ship.Type == Enum.SELF_DELIVERY {
						rep.SelfDeliveryFree = setting.SelfDeliveryFree
						rep.SelfDeliveryKey = setting.SelfDeliveryKey
					}
				}
			} else {
				rep.SelfDeliveryFree = 0
				rep.SelfDeliveryKey = Enum.FreeShipNone
			}

			_ = json.Unmarshal([]byte(v.PayWay), &productPayWayList)
			productPaymentList := rearrangesPayment(productPayWayList, productShipList)
			rep.ProductPayWayList = productPaymentList
			storeProduct = append(storeProduct, rep)
		}
	}
	total, err := product.CountProductTotalByStoreId(engine, storeId)
	sell, err := product.CountProductSellByStoreId(engine, storeId)
	down, err := product.CountProductDownByStoreId(engine, storeId)
	stock, err := product.CountProductStockByStoreId(engine, storeId)
	userData, _ := member.GetMemberDataByUid(engine, storeData.SellerId)
	SsmData, err := Store.GetStoreSocialMediaDataByStoreId(engine, storeData.StoreId)
	if err != nil {
		log.Error("Get store Social Media Data Error", err)
		return resp, err
	}
	resp.StoreSocialMedia = SsmData.GetStoreSocialMediaInfo()
	resp.StoreId = storeId
	resp.ProductCount = count
	resp.ProductTotal = total
	resp.ProductSellCount = sell
	resp.ProductDownCount = down
	resp.ProductStockCount = stock
	resp.StoreName = storeData.StoreName
	resp.StoreImage = storeData.StorePicture
	resp.StoreStatus = storeData.StoreStatus
	resp.VerifyIdentity = userData.VerifyIdentity
	resp.FreeShipKey = setting.FreeShipKey
	resp.FreeShip = setting.FreeShip
	resp.SelfDeliveryFree = setting.SelfDeliveryFree
	resp.SelfDeliveryKey = setting.SelfDeliveryKey
	resp.ProductList = storeProduct
	return resp, nil
}

//組即商品搜尋條件
func setParams(sellerId string, isRealtime int, params Request.ProductList) ([]string, []interface{}, string) {
	var sql []string
	var bind []interface{}
	var or []string
	var orderBy string

	sql = append(sql, "store_id = ?")
	bind = append(bind, sellerId)
	sql = append(sql, "is_realtime = ?")
	bind = append(bind, isRealtime)
	if isRealtime == 1 {
		sql = append(sql, "expire_date > ?")
		bind = append(bind, time.Now())
	}
	switch params.Status {
	case "onSell":
		sql = append(sql, "product_status = ?")
		bind = append(bind, Enum.ProductStatusSuccess)
	case "notSell":
		sql = append(sql, "product_status = ?")
		bind = append(bind, Enum.ProductStatusDown)
	case "stock":
		sql = append(sql, "stock = ?")
		bind = append(bind, 0)
	case "all":
	default:
		//下架
		sql = append(sql, "product_status != ?")
		bind = append(bind, Enum.ProductStatusDown)
		//官方下架
		sql = append(sql, "product_status != ?")
		bind = append(bind, Enum.ProductStatusPending)
	}
	if params.Name != "" {
		sql = append(sql, "product_name LIKE '%"+params.Name+"%'")
	}
	if params.Ship != "" {
		res := strings.Split(params.Ship, ",")
		for _, v := range res {
			or = append(or, "shipping LIKE '%"+v+"%'")
		}
	}
	if params.PayWay != "" {
		res := strings.Split(params.PayWay, ",")
		for _, v := range res {
			pay := ""
			switch v {
			case "cvsPay":
				pay = Enum.CvsPay
			case "atm":
				pay = Enum.Transfer
			case "credit":
				pay = Enum.Credit
			case "balance":
				pay = Enum.Balance
			}
			or = append(or, "pay_way LIKE '%"+pay+"%'")
		}
	}
	var order []string
	if params.Price != "" {
		order = append(order, "price "+params.Price)
	}
	if params.ShipMerge != "" {
		order = append(order, "ship_merge "+params.ShipMerge)
	}
	if or != nil {
		sql = append(sql, "("+strings.Join(or, " OR ")+")")
	}

	if order != nil {
		orderBy = "ORDER BY " + strings.Join(order, ", ")
	} else {
		orderBy = "ORDER BY create_time DESC"
	}
	return sql, bind, orderBy
}

//批次修改商品運送方式
func BatchChangeShipping(params Request.BatchProductShipping) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := engine.Session.Begin(); err != nil {
		log.Error("Set Begin Error", err)
		return fmt.Errorf("系統錯誤")
	}
	for _, v := range params.ProductId {
		productData, err := product.GetProductByProductId(engine, v)
		if err != nil {
			engine.Session.Rollback()
			return err
		}
		tempOldData := productData
		shipping, _ := json.Marshal(rearrangeShipMode(params.ShippingList))
		productData.Shipping = string(shipping)
		productData.ShipMerge = 0
		if params.ShipMerge == 1 {
			productData.ShipMerge = 1
		}
		if err := product.UpdateProductsData(engine, productData.ProductId, productData); err != nil {
			engine.Session.Rollback()
			log.Error("Update Product Error", err)
			return err
		}
		History.GenerateBatchShipWayProductLog(engine, tempOldData, productData, "")
	}
	if err := engine.Session.Commit(); err != nil {
		return fmt.Errorf("系統錯誤")
	}
	return nil
}

//批次修改商品付款方式
func BatchChangePayWay(params Request.BatchProductPayWay) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := engine.Session.Begin(); err != nil {
		log.Error("Set Begin Error", err)
		return fmt.Errorf("系統錯誤")
	}
	for _, v := range params.ProductId {
		productData, err := product.GetProductByProductId(engine, v)
		if err != nil {
			engine.Session.Rollback()
			return err
		}
		tempOldData := productData
		payWay, _ := json.Marshal(rearrangesPayWay(params.PayWayList))
		productData.PayWay = string(payWay)
		if err := product.UpdateProductsData(engine, productData.ProductId, productData); err != nil {
			engine.Session.Rollback()
			log.Error("Update Product Error", err)
			return err
		}
		History.GenerateBatchPayWayProductLog(engine, tempOldData, productData, "")
	}
	if err := engine.Session.Commit(); err != nil {
		return fmt.Errorf("系統錯誤")
	}
	return nil
}

func BatchChangeStatus(storeData entity.StoreDataResp, params Request.BatchProductStatus) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := engine.Session.Begin(); err != nil {
		log.Error("Set Begin Error", err)
		return fmt.Errorf("系統錯誤")
	}
	for _, v := range params.ProductId {
		productData, err := product.GetProductByProductId(engine, v)
		if err != nil {
			engine.Session.Rollback()
			return err
		}
		tempOldData := productData
		if StoreStatusIsClose(storeData) {
			productData.ProductStatus = Enum.ProductStatusPending
		} else {
			productData.ProductStatus = Enum.ProductStatusSuccess
		}
		if params.ProductDown == true {
			productData.ProductStatus = Enum.ProductStatusDown
		}
		if err := product.UpdateProductsData(engine, productData.ProductId, productData); err != nil {
			engine.Session.Rollback()
			log.Error("Update Product Error", err)
			return err
		}
		History.GenerateBatchDownProductLog(engine, tempOldData, productData, storeData.UserId)
	}
	if err := engine.Session.Commit(); err != nil {
		return fmt.Errorf("系統錯誤")
	}
	return nil
}

//func composePayWay(payway string) map[string]string {
//	var data = make(map[string]string)
//	out := strings.Split(payway, ",")
//	for _, v := range out {
//		val := strings.Trim(v, "\"")
//		data[val] = Enum.PayWay[val]
//	}
//	return data
//}
//重組商品資料
func composeProduct(engine *database.MysqlSession, productsData []entity.ProductsData) ([]Response.StoreProduct, error) {
	var storeProduct []Response.StoreProduct
	for _, v := range productsData {
		var storeProductSpec []Response.StoreProductSpec
		var productShipList []Response.ShippingMode
		var productPayWayList []Response.PayWayMode
		productData, err := product.GetProductSpecByProductId(engine, v.ProductId)
		if err != nil {
			return storeProduct, err
		}
		for _, elem := range productData {
			spec := Response.StoreProductSpec{}
			spec.ProductSpecId = elem.SpecId
			spec.Spec = elem.SpecName
			spec.Price = elem.SpecPrice
			spec.Quantity = elem.Quantity
			storeProductSpec = append(storeProductSpec, spec)
		}
		productImage, err := product.GetProductImageByProductId(engine, v.ProductId)
		if err != nil {
			return storeProduct, err
		}
		var imgs []string
		for _, elem := range productImage {
			imgs = append(imgs, images.GetImageUrl(elem.Image))
		}
		storeProductData := Response.StoreProduct{}
		storeProductData.ProductId = v.ProductId
		storeProductData.ProductName = v.ProductName
		storeProductData.ProductImage = imgs
		storeProductData.ProductIsSpec = v.IsSpec
		storeProductData.ProductSpecList = storeProductSpec
		storeProductData.ProductPrice = v.Price
		storeProductData.ProductStatus = v.ProductStatus
		storeProductData.ProductStatusText = ""
		storeProductData.ProductExpireTime = ""
		if v.IsRealtime == 1 {
			storeProductData.ProductStatusText = Enum.RealtimeStatus[v.ProductStatus]
			storeProductData.ProductExpireTime = v.ExpireDate.Format("2006/01/02 15:04")
		}
		storeProductData.ProductCancelTime = ""
		if v.ProductStatus == Enum.ProductStatusCancel {
			storeProductData.ProductCancelTime = v.UpdateTime.Format("2006/01/02 15:04")
		}
		storeProductData.ProductShipMerge = v.ShipMerge
		storeProductData.ProductQrcode = fmt.Sprintf("/static/images/qrcode/%s.jpg", productData[0].ProductId)
		_ = json.Unmarshal([]byte(v.Shipping), &productShipList)
		storeProductData.ProductShipList = productShipList
		_ = json.Unmarshal([]byte(v.PayWay), &productPayWayList)
		storeProductData.ProductPayWayList = productPayWayList
		storeProductData.FormUrl = v.FormUrl
		storeProduct = append(storeProduct, storeProductData)
	}
	return storeProduct, nil
}

//執行帳單取消
func HandleRealTimesCancel(params Request.SetRealTimesRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := engine.Session.Begin(); err != nil {
		log.Error("Set Begin Error", err)
		return fmt.Errorf("系統錯誤")
	}
	data, err := product.GetRealTimeByProductId(engine, params.ProductId)
	if err != nil {
		engine.Session.Rollback()
		return fmt.Errorf("無帳單")
	}
	data.ProductStatus = Enum.ProductStatusCancel
	data.UpdateTime = time.Now()
	if err := product.UpdateProductData(engine, data.ProductId, &data); err != nil {
		engine.Session.Rollback()
		return fmt.Errorf("系統錯誤")
	}
	if err := engine.Session.Commit(); err != nil {
		return fmt.Errorf("系統錯誤")
	}
	return nil
}

//執行帳單延長
func HandleRealTimesExtension(params Request.SetRealTimesRequest) (Response.RealTimesExtensionResponse, error) {
	var resp Response.RealTimesExtensionResponse
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := engine.Session.Begin(); err != nil {
		log.Error("Set Begin Error", err)
		return resp, fmt.Errorf("系統錯誤")
	}
	data, err := product.GetRealTimeByProductId(engine, params.ProductId)
	if err != nil {
		engine.Session.Rollback()
		return resp, fmt.Errorf("無帳單")
	}
	create, _ := time.Parse("2006/01/02", data.CreateTime.Format("2006/01/02"))
	expire, _ := time.Parse("2006/01/02", data.ExpireDate.Format("2006/01/02"))
	day := create.Sub(expire).Hours() / 24
	if day > 3 {
		engine.Session.Rollback()
		return resp, fmt.Errorf("已延期過不能再延期")
	}
	data.ExpireDate = time.Now().Add(72 * time.Hour)
	data.UpdateTime = time.Now()
	if err := product.UpdateProductData(engine, data.ProductId, &data); err != nil {
		engine.Session.Rollback()
		return resp, fmt.Errorf("系統錯誤")
	}
	resp.Time = data.ExpireDate.Format("2006/01/02 15:04")
	if err := engine.Session.Commit(); err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	return resp, nil
}

//func GeneratorProductShortUrl() error {
//	engine := database.GetMysqlEngine()
//	defer engine.Close()
//	data, err := product.GetProductAll(engine)
//	if err != nil {
//		return err
//	}
//	for _, v := range data {
//		var data entity.ShortUrlData
//		data.Short = v.TinyUrl
//		data.Url = fmt.Sprintf("/product/%s", v.ProductId)
//		if err := Short.InsertShortUrlData(engine, data); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//批次修改免運設定
func BatchChangeFreeShip(params Request.BatchProductFreeShip) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := engine.Session.Begin(); err != nil {
		return fmt.Errorf("1001001")
	}
	for _, v := range params.ProductId {
		productData, err := product.GetProductByProductId(engine, v)
		if err != nil {
			engine.Session.Rollback()
			return fmt.Errorf("1001001")
		}
		if len(productData.ProductId) == 0 {
			return fmt.Errorf("1001010")
		}
		if params.IsFreeShip == true {
			productData.IsFreeShip = true
		} else {
			productData.IsFreeShip = false
		}
		if err := product.UpdateProductsData(engine, productData.ProductId, productData); err != nil {
			engine.Session.Rollback()
			return fmt.Errorf("1001001")
		}
	}
	if err := engine.Session.Commit(); err != nil {
		return fmt.Errorf("1001001")
	}
	return nil
}
