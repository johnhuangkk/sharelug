package model

import (
	"api/services/Enum"
	"api/services/Service/Excel"
	"api/services/VO/ExcelVo"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Promotion"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"strconv"
	"time"
)

const coupnMaxLen = 8

func CreatePromotion(params Request.PromoCreate, userId string, storeId string, tNow time.Time) (Response.PromoCreate, error) {
	var resp Response.PromoCreate
	var data entity.Promotion
	data.Name = params.Name
	data.StartTime = tNow
	data.Remaining = params.Quantity
	data.Amount = params.Quantity
	data.Type = Enum.PromoTypeAmount
	data.Creator = userId
	data.Value = float64(params.Amount)
	data.Status = Enum.PromoStatusActive
	data.StoreId = storeId
	closeTime, _ := time.ParseInLocation("2006-01-02 15:04", params.EndTime, time.Local)
	data.CloseTime = closeTime
	err := Promotion.InsertPromo(data)
	if err != nil {
		log.Error("Create Promotion Error", err)
		return resp, err
	}
	resp.Name = data.Name
	resp.StartTime = data.StartTime.Format("2006-01-02 15:04")
	resp.EndTime = data.CloseTime.Format("2006-01-02 15:04")
	resp.Quantity = data.Amount
	resp.Amount = int64(data.Value)
	resp.Status = data.Status
	return resp, nil
}

func GetPromotion(storeId string, promoId string, now time.Time) ([]Response.PromoDetail, error) {
	var resp []Response.PromoDetail
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if len(promoId) == 0 {
		promos, err := Promotion.GetPromos(storeId)
		if err != nil {
			return resp, nil
		}
		for _, promo := range promos {
			promoResp := promo.GetPromotionDetailResponse(now)
			promoResp.Used, _ = Promotion.CountUsedCoupon(engine, string(rune(promo.Id)))
			promoResp.UnUsed = promoResp.Picked - promoResp.Used
			resp = append(resp, promoResp)
		}

	} else {
		promo, err := Promotion.GetPromo(promoId)
		if err != nil {
			return resp, nil
		}
		promoResp := promo.GetPromotionDetailResponse(now)
		promoResp.Used, _ = Promotion.CountUsedCoupon(engine, promoId)
		promoResp.UnUsed = promoResp.Picked - promoResp.Used
		resp = append(resp, promoResp)
	}
	if len(resp) > 0 {
		resp = sortPromoDetails(resp)
	}
	return resp, nil
}

func TakeCoupons(qty int, promoId string, userId string, storeId string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	//generate coupons
	intPromoId, _ := strconv.Atoi(promoId)
	operateId := tools.GenerateCouponActionId()
	promo, err := Promotion.GetPromo(promoId)
	if err != nil {
		return fmt.Errorf("1001001")
	}
	if qty > int(promo.Remaining) {
		return fmt.Errorf("1011008")
	}
	if promo.Remaining > 0 {
		existCodeMap, err := getPromoExistCodeMap(engine, promoId, storeId)
		if err != nil {
			return fmt.Errorf("1001001")
		}
		newCodes := generateNewCodes(qty, operateId, intPromoId, existCodeMap)
		err = Promotion.InsertPromoCodes(engine, newCodes)
		if err != nil {
			return fmt.Errorf("1001001")
		}
		nowRemaining := promo.Remaining - int64(qty)
		err = Promotion.UpdatePromotionRemaining(engine, promoId, nowRemaining)
		if err != nil {
			return fmt.Errorf("1001001")
		}
		err = Promotion.InsertPromoOperateRecord(engine, int64(intPromoId), qty, operateId, userId)
		if err != nil {
			return fmt.Errorf("1001001")
		}
	}
	return nil
}
func TerminatePromotion(promoId string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()

	err := Promotion.UpdatePromotionStatus(engine, promoId, Enum.PromoStatusStop)
	if err != nil {
		return fmt.Errorf("1001001")
	}
	return nil
}

func GetUnuseCouponList(params Request.PromoCouponListGet, promoId string) (Response.CouponUnuseList, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.CouponUnuseList
	var couponList []Response.CouponUnuse
	counts, err := Promotion.CountUnuseCoupon(engine, promoId)
	if err != nil {
		return resp, fmt.Errorf("1001001")
	}
	codes, err := Promotion.GetUnuseCodeWithPromoIdAndPageLimit(engine, promoId, params.Page, params.Limit)
	if err != nil {
		return resp, fmt.Errorf("1001001")
	}

	for _, code := range codes {
		var data Response.CouponUnuse
		data.Code = code.Code
		data.Id = code.Id
		data.IsCopy = code.IsSellerCopy
		data.PromotionId = code.PromotionId
		data.StatusText = Enum.PromoCodeUnuseText
		couponList = append(couponList, data)
	}
	resp.Coupons = couponList
	resp.Counts = counts
	return resp, nil
}
func GetUsedCouponList(params Request.PromoCouponListGet, promoId string) (Response.CouponUsedList, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.CouponUsedList
	var couponList []Response.CouponUsed
	counts, err := Promotion.CountUsedCoupon(engine, promoId)
	if err != nil {
		return resp, fmt.Errorf("1001001")
	}
	records, err := Promotion.GetUsedCodeWithPromoIdAndPageLimit(engine, promoId, params.Page, params.Limit)
	if err != nil {
		return resp, fmt.Errorf("1001001")
	}
	for _, record := range records {
		var data Response.CouponUsed
		data.Amount = record.Amount
		data.DiscountAmount = record.DiscountAmount
		data.StatusText = Enum.PromoCodeUsedText
		data.OrderId = record.OrderId
		data.BuyerPhone = tools.MaskerPhone(record.BuyerPhone)
		data.Code = record.Code
		data.UseDate = record.TransTime.Format("2006-01-02 15:04")
		couponList = append(couponList, data)
	}
	resp.Counts = counts
	resp.Coupons = couponList
	return resp, nil
}
func SetCouponIsCopy(promoId, couponId string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	err := Promotion.UpdateCouponIsCopy(engine, promoId, couponId)
	if err != nil {
		return fmt.Errorf("1001001")
	}
	return nil
}

func GenerateUsedCouponRecordExcel(promoId string) (string, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var datas []ExcelVo.CouponUsedReportVo
	promo, err := Promotion.GetPromo(promoId)
	if err != nil {
		return "", fmt.Errorf("1001001")
	}
	records, err := Promotion.GetUsedCouponRecordWithPromoId(engine, promoId)
	if err != nil {
		return "", fmt.Errorf("1001001")
	}
	for i, record := range records {
		data := ExcelVo.CouponUsedReportVo{
			Id:             int64(i + 1),
			Status:         Enum.PromoCodeUsedText,
			Code:           record.Code,
			OrderId:        record.OrderId,
			TransTime:      record.TransTime.Format("2006-01-02 15:04"),
			BuyerPhone:     tools.MaskerPhone(record.BuyerPhone),
			Amount:         int64(record.Amount),
			DiscountAmount: int64(record.DiscountAmount),
		}
		datas = append(datas, data)
	}
	filename, err := Excel.CouponRecordNew().ToCouponUsedRecordFile(datas, promo.Name)
	if err != nil {
		log.Debug("Generator Report File Error", err)
		return filename, err
	}
	return filename, nil

}
func getPromoExistCodeMap(engine *database.MysqlSession, promoId string, storeId string) (map[string]bool, error) {

	tempCodesList := make(map[string]bool)
	promos, err := Promotion.GetPromos(storeId)
	if err != nil {
		return tempCodesList, fmt.Errorf("1001001")
	}
	var totalExistCodes []entity.PromotionCode
	for _, promo := range promos {
		existCodes, err := Promotion.GetPromoCodes(engine, promo.Id)
		if err != nil {
			return tempCodesList, fmt.Errorf("1001001")
		}
		totalExistCodes = append(totalExistCodes, existCodes...)
	}

	if len(totalExistCodes) == 0 {
		return tempCodesList, nil
	}
	for _, existCode := range totalExistCodes {
		tempCodesList[existCode.Code] = true
	}
	return tempCodesList, nil
}
func generateNewCodes(qty int, operateId string, promoId int, codeMap map[string]bool) []entity.PromotionCode {
	var newCodes []entity.PromotionCode
	i := 0
	for i < qty {
		code := tools.RandLowerString(coupnMaxLen)
		if !codeMap[code] {
			codeMap[code] = true
			var codeData entity.PromotionCode
			codeData.Code = code
			codeData.OperateId = operateId
			codeData.PromotionId = int64(promoId)
			newCodes = append(newCodes, codeData)
			i++
		}
	}
	return newCodes
}

func HandleCouponUsedRecord(engine *database.MysqlSession, order entity.OrderData) error {
	if len(order.CouponNumber) == 0 {
		return nil
	}
	data, err := Promotion.GetPromotionCodeByCodeAndStoreId(engine, order.StoreId, order.CouponNumber)
	if err != nil {
		log.Debug("Get Promotion Code Error", err)
		return err
	}
	data.PromotionCode.IsUsed = true
	if err := Promotion.UpdatePromotionCode(engine, data.PromotionCode); err != nil {
		log.Debug("Update Promotion Code Error", err)
		return err
	}
	record := order.GeneratorCouponUsedRecordData(data.Promotion.Id)
	if err := Promotion.InsertUsedCouponRecord(engine, record); err != nil {
		log.Error("Insert Used Coupon Record Error", err)
		return err
	}
	return nil
}

func sortPromoDetails(details []Response.PromoDetail) []Response.PromoDetail {
	var active []Response.PromoDetail
	var disable []Response.PromoDetail
	for _, detail := range details {
		if detail.Status == Enum.PromoStatusActive {
			active = append(active, detail)
		} else {
			disable = append(disable, detail)
		}
	}
	active = append(active, disable...)
	return active
}

func CheckPromoActiveLimit(storeId string, now time.Time) (bool, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	counts, err := Promotion.GetActivePromoCountWithStore(engine, storeId, now)
	if err != nil {
		return false, fmt.Errorf("1001001")
	}
	if counts >= 5 {
		return true, nil
	}
	return false, nil
}
