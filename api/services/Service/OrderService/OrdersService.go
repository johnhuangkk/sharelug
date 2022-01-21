package OrderService

import (
	"api/services/Enum"
	"api/services/Service/History"
	"api/services/Service/Product"
	"api/services/Service/StoreService"
	"api/services/VO/CartsVo"
	"api/services/VO/Response"
	"api/services/dao/Orders"
	"api/services/dao/product"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"time"
)
//檢查購物車內所有商品庫存數
func CheckProductStock(engine *database.MysqlSession, carts CartsVo.Carts) ([]entity.ProductsData, error) {
	products, err := GetProducts(engine, carts)
	if err != nil {
		return nil, err
	}
	if len(products) == 0 {
		log.Error("購物車無商品", len(products))
		return nil, fmt.Errorf("1005014")
	}
	//檢查購物車內所有商品庫存數
	for k, v := range products {
		if carts.Products[k].ProductSpecId == v.SpecId {
			if carts.Products[k].Quantity == 0 {
				log.Error("商品目前庫存不足", carts.Products[k].Quantity, v.Quantity)
				return nil, fmt.Errorf("1005013")
			}
			if int64(carts.Products[k].Quantity) > v.Quantity {
				log.Error("商品目前庫存不足", carts.Products[k].Quantity, v.Quantity)
				return nil, fmt.Errorf("1005009")
			}
		}
	}
	return products, nil
}
//取多個商品資料
func GetProducts(engine *database.MysqlSession, data CartsVo.Carts) ([]entity.ProductsData, error) {
	var products []entity.ProductsData
	for _, value := range data.Products {
		productData, err := Product.GetProductDataByProductSpecId(engine, value.ProductSpecId)
		if err != nil {
			return products, err
		}
		products = append(products, productData)
	}
	return products, nil
}
//檢查購買人所購買的商品是否為本人的
func CheckNonOwnedProduct(engine *database.MysqlSession, buyerId string, storeId string) error {
	store, err := StoreService.GetUserAllStore(engine, buyerId)
	if err != nil {
		return err
	}
	if tools.InArray(store, storeId) {
		return fmt.Errorf("1005001")
	}
	return nil
}
//變更訂單狀態
func ChangeOrderStatus(engine *database.MysqlSession, orderData entity.OrderData, status string) error {
	OldStatus := Enum.OrderInit
	if orderData.OrderStatus != "" {
		OldStatus = orderData.OrderStatus
	}
	orderData.OrderStatus = status
	//checkOrderSuccessAndNotDelivery(&orderData)
	err := UpdateOrderDataByOrderStatus(engine, OldStatus, orderData)
	if err != nil {
		log.Error("Update Order Data Error", err)
		return err
	}
	return nil
}

//變更訂單撥付狀態
func ChangeOrderCaptureStatus(engine *database.MysqlSession, orderData entity.OrderData, status string) error {
	OldStatus := Enum.OrderCaptureInit
	if orderData.CaptureStatus != "" {
		OldStatus = orderData.CaptureStatus
	}
	orderData.CaptureStatus = status
	orderData.UpdateTime = time.Now()
	if _, err := Orders.UpdateOrderData(engine, orderData.OrderId, orderData); err != nil {
		log.Error("Update Order Data Error", err)
		return err
	}
	if err := SetOrderStatusLog("CaptureStatus",  orderData.OrderId, OldStatus, orderData.CaptureStatus, ""); err != nil {
		log.Error("Update Status History Log Error", err)
		return err
	}
	return nil
}


//待付款
func OrderWaitPayment(engine *database.MysqlSession, OrderData entity.OrderData, Status string) error {
	if err := ChangeOrderStatus(engine, OrderData, Status); err != nil {
		return err
	}
	return nil
}


//檢查無需配送
//func checkOrderSuccessAndNotDelivery(OrderData *entity.OrderData) {
//	if OrderData.OrderStatus == Enum.OrderSuccess  && tools.InArray([]string{Enum.F2F, Enum.NONE}, OrderData.ShipType) {
//		OrderData.ShipStatus = Enum.OrderShipNone
//		OrderCaptureRelease(OrderData, time.Time{})
//	}
//}

//更新訂單狀態
func UpdateOrderDataByOrderStatus(engine *database.MysqlSession, OldStatus string, OrderData entity.OrderData) error {
	OrderData.UpdateTime = time.Now()
	_, err := Orders.UpdateOrderData(engine, OrderData.OrderId, OrderData)
	if err != nil {
		log.Error("Update Order Data Error", err)
		return err
	}

	err = SetOrderStatusLog("OrderStatus",  OrderData.OrderId, OldStatus, OrderData.OrderStatus, "")
	if err != nil {
		log.Error("Update Status History Log Error", err)
		return err
	}
	return nil
}

// 記錄狀態變更
func SetOrderStatusLog(Field string, OrderId string, OldStatus string, NewStatus string, UserName string) error {
	err := History.GenerateStatusLog("OrderData", Field, OrderId, OldStatus, NewStatus, UserName)
	if err != nil {
		log.Error("Update Status History Log Error", err)
		return err
	}
	return nil
}

//關閉帳單
func ProcessBillingClose(engine *database.MysqlSession, ProductId, status string) error {
	productData, err := product.GetProductDataByProductId(engine, ProductId)
	if err != nil {
		log.Error("Get product data Error", err)
		return err
	}
	if err := product.UpdateProductStatusByProductId(engine, productData.ProductId, status); err != nil {
		log.Error("Update product data Error", err)
		return err
	}
	return nil
}

//設定撥款時間
func OrderCaptureRelease(data *entity.OrderData, date time.Time)  {
	if data.CaptureStatus != Enum.OrderCaptureSuccess {
		if date.IsZero() {
			//判斷撥款狀態是否為待撥付
			if data.CaptureStatus == Enum.OrderCaptureProgress {
				//撥款狀態 變更為延後撥付
				data.CaptureStatus = Enum.OrderCapturePostpone
				data.CaptureApply = Enum.OrderCapturePostpone
				//取撥款時間再加10天
				data.CaptureTime = tools.GenerateExpireTime(10)
			} else {
				//撥款狀態 變更為待撥付
				data.CaptureStatus = Enum.OrderCaptureProgress
				data.CaptureApply = Enum.OrderCaptureProgress
				data.CaptureTime = tools.GenerateExpireTime(10)
			}
		} else {
			//撥款狀態 變更為提前撥付
			data.CaptureTime = date
			data.CaptureStatus = Enum.OrderCaptureAdvance
			data.CaptureApply = Enum.OrderCaptureAdvance
		}
	}
}

//取退貨退款
func GetReturnAndRefundByOrderId(engine *database.MysqlSession, OrderData *Response.OrderResponse) error {
	countRefund, err := Orders.CountOrderRefundByOrderId(engine, OrderData.OrderId)
	if err != nil {
		return err
	}
	OrderData.IsRefund = false
	if countRefund > 0 {
		OrderData.IsRefund = true
	}
	countReturn, err := Orders.CountOrderReturnByOrderId(engine, OrderData.OrderId)
	if err != nil {
		return err
	}
	OrderData.IsReturn = false
	if countReturn > 0 {
		OrderData.IsReturn = true
	}
	return nil	
}

func OrderOpenInvoice(engine *database.MysqlSession, data *entity.OrderData) error {
	data.AskInvoice = true
	data.InvoiceStatus = Enum.InvoiceOpenStatusNot
	if err := Orders.UpdateOrdersData(engine, data); err != nil {
		log.Error("update order Error", err)
		return err
	}
	return nil
}

//判斷目前訂單狀態 fixme
func GetOrderStatusText(orderData entity.OrderData) string {
	if orderData.OrderStatus == Enum.OrderSuccess {
		if orderData.ShipType != Enum.F2F {
			return Enum.OrderShipStatus[orderData.ShipStatus]
		} else if orderData.ShipType == Enum.F2F {
			return Enum.OrderF2fStatus[orderData.ShipStatus]
		}
	}
	return Enum.OrderStatus[orderData.OrderStatus]
}

func ChangeReturnStatus(engine *database.MysqlSession, orderId string) error {
	data, err :=  Orders.GetOrderRefundsByOrderIdAndStatus(engine, orderId)
	if err != nil {
		log.Error("Get Refund Database Error", err)
		return err
	}
	data.Status = Enum.ReturnStatusSuccess
	if err := Orders.UpdateOrderRefundData(engine, data); err != nil {
		return err
	}
	return nil
}
//取新的訂單編號
func GetNewOrderId(Realtime int, PayType string) string {
	if Realtime == 1 {
		return tools.GeneratorRealtimeOrderId()
	} else if Realtime == 0 && PayType == Enum.PayTypeMarket {
		return tools.GeneratorMarketingOrderId()
	}
	return tools.GeneratorOrderId()
}