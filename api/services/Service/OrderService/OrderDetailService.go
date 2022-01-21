package OrderService

import (
	"api/services/Enum"
	"api/services/VO/CartsVo"
	"api/services/VO/Response"
	"api/services/dao/Orders"
	"api/services/dao/Store"
	"api/services/dao/product"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"encoding/json"
)

//處理Order Detail
func ProcessOrderDetail(engine *database.MysqlSession, carts CartsVo.Carts, products []entity.ProductsData, orderData entity.OrderData) error {
	for _, value := range carts.Products {
		if err := createOrderDetail(engine, orderData.OrderId, value, products, carts.Shipping); err != nil {
			log.Error("Create Order Detail", err)
			return err
		}
		for _, v := range products {
			if value.ProductSpecId == v.ProductSpecData.SpecId {
				if v.IsRealtime != 1 {
					log.Debug("Minus Product stock", v.ProductId, value.ProductSpecId, value.Quantity, orderData.SellerId)
					if err := product.UpdateMinusProductStockByProductSpecId(engine, v.ProductId, value.ProductSpecId, value.Quantity); err != nil {
						return err
					}
				} else {
					if err := ProcessBillingClose(engine, v.ProductId, Enum.ProductStatusPaid); err != nil {
						log.Error("close billing Error!!")
						return err
					}
				}
			}
		}

	}
	return nil
}

// 建立訂單Detail
func createOrderDetail(engine *database.MysqlSession, orderId string, data CartsVo.CartProduct, products []entity.ProductsData, choose string) error {
	for _, value := range products {
		if value.SpecId == data.ProductSpecId {
			storeData, _ := Store.GetStoreDataByStoreId(engine, value.StoreId)
			var Entity entity.OrderDetail
			ShipFee := GetProductShipFee(value.Shipping, choose)
			Entity.SellerId = storeData.SellerId
			Entity.OrderId = orderId
			Entity.ProductSpecId = data.ProductSpecId
			Entity.ProductSpecName = value.SpecName
			Entity.ProductName = value.ProductName
			Entity.ProductQuantity = int64(data.Quantity)
			Entity.ProductPrice = value.Price
			Entity.Subtotal = value.Price * int64(data.Quantity)
			Entity.ShipMerge = int64(value.ShipMerge)
			Entity.ShipFee = int64(ShipFee)
			Entity.IsFreeShip = value.IsFreeShip
			err := Orders.InsertOrderDetail(engine, Entity)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 取得商品運費
func GetProductShipFee(shipList string, choose string) int {
	var fee int
	var ship []Response.ShippingMode
	_ = json.Unmarshal([]byte(shipList), &ship)
	for _, value := range ship {
		if value.Type == choose {
			fee = value.Price
		}
	}
	return fee
}

func PaymentFailReturnStock(engine *database.MysqlSession, OrderId string) error {
	//訂購失敗退還庫存
	detail, err := Orders.GetOrderDetailByOrderId(engine, OrderId)
		if err != nil {
		log.Error("Get Order Detail Error", err)
		return err
	}
	for _, v := range detail {
		data, err := product.GetProductByProductSpecId(engine, v.ProductSpecId)
			if err != nil {
			log.Error("Get product data Error", err)
			return err
		}
		//判斷是否為即時帳單
		if data.IsRealtime != 1 {
			if err := product.UpdatePlusProductStockByProductSpecId(engine, data.ProductId, data.SpecId, int(v.ProductQuantity)); err != nil {
				return err
			}
		} else {
			if err := ProcessBillingClose(engine, data.ProductId, Enum.ProductStatusSuccess); err != nil {
				log.Error("close billing Error!!")
				return err
			}
		}
	}
	return nil
}

//訂購失敗退還庫存
func ProcessReturnStock(engine *database.MysqlSession, OrderId string) error {
	detail, err := Orders.GetOrderDetailByOrderId(engine, OrderId)
	if err != nil {
		log.Error("Get Order Detail Error", err)
		return err
	}
	for _, v := range detail {
		data, err := product.GetProductByProductSpecId(engine, v.ProductSpecId)
		if err != nil {
			log.Error("Get product data Error", err)
			return err
		}
		if err := product.UpdatePlusProductStockByProductSpecId(engine, data.ProductId, data.SpecId, int(v.ProductQuantity)); err != nil {
			return err
		}
	}
	return nil
}

// 建立訂單Detail
func CreateBillOrderDetail(engine *database.MysqlSession, orderId, sellerId string, bill entity.BillOrderData) error {
	var Entity entity.OrderDetail
	Entity.SellerId = sellerId
	Entity.OrderId = orderId
	Entity.ProductSpecId = orderId
	Entity.ProductSpecName = bill.ProductSpec
	Entity.ProductName = bill.ProductName
	Entity.ProductQuantity = bill.ProductQty
	Entity.ProductPrice = bill.ProductPrice
	Entity.Subtotal = bill.ProductPrice * int64(bill.ProductQty)
	Entity.ShipMerge = 1
	Entity.ShipFee = int64(bill.ShipFee)
	err := Orders.InsertOrderDetail(engine, Entity)
	if err != nil {
				return err
			}

	return nil
}
//取出訂單內的表單連結
func GetDetailProductFormUrl(engine *database.MysqlSession, detail []entity.ProductsData) string {
	for _, v := range detail {
		data, _ := product.GetProductsByProductId(engine, v.ProductId)
		if data.FormUrl != "" {
			return data.FormUrl
		}
	}
	return ""
}


