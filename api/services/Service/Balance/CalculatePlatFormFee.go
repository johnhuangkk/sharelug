package Balance

import (
	"api/services/Enum"
	"api/services/VO/OrderVo"
	"api/services/entity"
	"api/services/util/tools"
)

/**
 * 計算平台費用
 */
func CalculatePlatFormFee(vo OrderVo.CalculatePlatFormFeeVo) OrderVo.CalculatePlatFormFeeResponse {
	var resp OrderVo.CalculatePlatFormFeeResponse
	resp.PlatformTransFee = calculateTransFee()
	resp.PlatformShipFee = calculateShipFee(vo.ShipType)
	resp.PlatformPayFee = calculatePayFee(vo.Amount, vo.PayWayType)
	resp.PlatformInfoFee = calculateInfoFee()
	resp.CaptureAmount = vo.Amount - (resp.PlatformTransFee + resp.PlatformShipFee + resp.PlatformPayFee + resp.PlatformInfoFee)
	return resp
}
//計算運費費用
func calculateShipFee(ShipType string) float64 {
	fee := 0
	switch  ShipType {
	case Enum.CVS_FAMILY:
		fee = 60
	case Enum.CVS_OK_MART:
		fee = 60
	case Enum.CVS_HI_LIFE:
		fee = 60
	case Enum.DELIVERY_I_POST_BAG1, Enum.DELIVERY_I_POST_BAG2, Enum.DELIVERY_I_POST_BAG3:
		fee = 70
	case Enum.DELIVERY_POST_BAG1:
		fee = 36
	case Enum.DELIVERY_POST_BAG2:
		fee = 48
	case Enum.DELIVERY_POST_BAG3:
		fee = 60
	case Enum.I_POST:
		fee = 70
	}
	return float64(fee)
}

//計算金流處理費
func calculatePayFee(Amount float64, PayWay string) float64 {
	//       信用卡 虛擬帳號 餘額 台灣Pay 超商取貨付款
	//系統公式 X% X%+Y X%+Y X%+Y X%+Y
	//    收費 2% 5 0 2% 5
	TransferFee := 0
	ratio := 0
	switch PayWay {
	case Enum.Credit:
		ratio = 2
	case Enum.Transfer:
		ratio = 0
		TransferFee = 0
	case Enum.Balance:
		ratio = 0
		TransferFee = 0
	case Enum.CvsPay:
		ratio = 0
		TransferFee = 0
	case Enum.TaiwanPay:
		ratio = 2
		TransferFee = 0
	}
	return tools.Round(Amount * float64(ratio) / 100) + float64(TransferFee)
}

//計算資訊處理費
func calculateInfoFee() float64 {
	fee := 10
	return float64(fee)
}

// 計算成交手續費
func calculateTransFee() float64 {
	TransFee := 0
	return float64(TransFee)
}
/**
 * 計算應收金額
 */
func CalculateIncome(order entity.OrderData) int64 {
	income := order.TotalAmount - (order.PlatformTransFee + order.PlatformShipFee + order.PlatformPayFee + order.PlatformInfoFee)
	return int64(income)
}
