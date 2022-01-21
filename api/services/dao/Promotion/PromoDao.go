package Promotion

import (
	"api/services/Enum"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"time"
)

func InsertPromo(data entity.Promotion) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	_, err := engine.Session.Table(entity.Promotion{}).Insert(&data)
	if err != nil {
		log.Error("Promotion insert error", err.Error())
	}
	return nil
}

func GetPromo(promoId string) (entity.Promotion, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var data entity.Promotion
	_, err := engine.Engine.Table(entity.Promotion{}).Where("id=?", promoId).Get(&data)
	if err != nil {
		log.Error("Get one promo error", err.Error())
		return data, err
	}
	return data, nil
}

func GetPromos(storeId string) ([]entity.Promotion, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var data []entity.Promotion
	err := engine.Engine.Table(entity.Promotion{}).Where("store_id=?", storeId).Desc("created").Find(&data)
	if err != nil {
		log.Error("Get promos error", err.Error())
		return data, err
	}
	return data, err
}

func CountUsedCoupon(engine *database.MysqlSession, promoId string) (int64, error) {

	counts, err := engine.Engine.Table(entity.CouponUsedRecord{}).Where("promotion_id = ? AND record_status =?", promoId, Enum.RecordStatusSuccess).Count()
	if err != nil {
		log.Error("Count promo used code error", err.Error())
		return 0, err
	}
	return counts, err
}

func InsertUsedCoupon(data entity.CouponUsedRecord) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	_, err := engine.Session.Table(entity.CouponUsedRecord{}).Insert(&data)
	if err != nil {
		log.Error("Insert Use Coupon Record", err.Error())
		return err
	}
	return nil
}

func GetPromoCodes(engine *database.MysqlSession, promoId int64) ([]entity.PromotionCode, error) {
	var codes []entity.PromotionCode
	err := engine.Engine.Table(entity.PromotionCode{}).Where("promotion_id =?", promoId).
		Cols("promotion_id, code").
		Find(&codes)
	if err != nil {
		log.Error("Find promo codes", err.Error())
		return codes, err
	}
	return codes, nil
}

func InsertPromoCodes(engine *database.MysqlSession, codes []entity.PromotionCode) error {
	_, err := engine.Session.Table(entity.PromotionCode{}).Insert(&codes)
	if err != nil {
		log.Error("Insert promo code error", err.Error())
		return err
	}
	return nil
}

func UpdatePromotionRemaining(engine *database.MysqlSession, promoId string, remaining int64) error {
	_, err := engine.Session.Table(entity.Promotion{}).
		Cols("id,remaining").
		Where("id =?", promoId).Update(entity.Promotion{Remaining: remaining})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func InsertPromoOperateRecord(engine *database.MysqlSession, promoId int64, qty int, operateId, userId string) error {
	var data entity.PromoCodeOperateRecord
	data.Amount = int64(qty)
	data.OperateId = operateId
	data.PromotionId = promoId
	data.Type = Enum.PromoOperateTake
	data.UserId = userId
	_, err := engine.Session.Table(entity.PromoCodeOperateRecord{}).Insert(&data)
	if err != nil {
		log.Error("Insert promo code error", err.Error())
		return err
	}
	return nil
}

func GetUnuseCodeWithPromoIdAndPageLimit(engine *database.MysqlSession, promoId string, page, limit int) ([]entity.PromotionCode, error) {
	var data []entity.PromotionCode
	err := engine.Engine.Table(entity.PromotionCode{}).Where("promotion_id = ? and is_used =?", promoId, false).Limit(limit, (page-1)*limit).Asc("is_seller_copy").Find(&data)
	if err != nil {
		log.Error("Insert promo code error", err.Error())
		return data, err
	}
	return data, nil
}

func GetUsedCodeWithPromoIdAndPageLimit(engine *database.MysqlSession, promoId string, page, limit int) ([]entity.CouponUsedRecord, error) {
	var data []entity.CouponUsedRecord
	err := engine.Engine.Table(entity.CouponUsedRecord{}).Where("promotion_id = ? AND record_status =?", promoId, Enum.RecordStatusSuccess).Limit(limit, (page-1)*limit).Find(&data)
	if err != nil {
		log.Error("Insert promo code error", err.Error())
		return data, err
	}
	return data, nil
}

func UpdatePromotionStatus(engine *database.MysqlSession, promoId, status string) error {
	now := time.Now()
	_, err := engine.Session.Table(entity.Promotion{}).Cols("id,status,stop_time").
		Where("id =?", promoId).
		Update(entity.Promotion{Status: status, StopTime: now})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

//取出指定賣場優惠卷未使用過的
func GetPromotionCodeByCodeAndStoreId(engine *database.MysqlSession, storeId, code string) (entity.PromotionAndPromotionCode, error) {
	var data entity.PromotionAndPromotionCode
	if _, err := engine.Engine.Table(entity.Promotion{}).
		Join("LEFT", entity.PromotionCode{}, "promotion.id = promotion_code.promotion_id").
		Where("promotion.store_id = ? AND promotion_code.code = ?", storeId, code).Get(&data); err != nil {
		log.Error("get promo codes", err)
		return data, err
	}
	return data, nil
}

//新增一筆使用過的優惠卷
func InsertUsedCouponRecord(engine *database.MysqlSession, data entity.CouponUsedRecord) error {
	if _, err := engine.Session.Table(entity.CouponUsedRecord{}).Insert(&data); err != nil {
		log.Error("Insert Use Coupon Record Error", err)
		return err
	}
	return nil
}

//取出優惠卷使用記錄
func GetUsedCouponRecordByOrderId(engine *database.MysqlSession, orderId string) (entity.CouponUsedRecord, error) {
	var data entity.CouponUsedRecord
	if _, err := engine.Engine.Table(entity.CouponUsedRecord{}).Select("*").Where("order_id = ?", orderId).Get(&data); err != nil {
		log.Error("Count promo used code error", err.Error())
		return data, err
	}
	return data, nil
}

//更新優惠卷使用記錄
func UpdateUsedCouponRecord(engine *database.MysqlSession, data entity.CouponUsedRecord) error {
	if _, err := engine.Session.Table(entity.CouponUsedRecord{}).Where("order_id = ?", data.OrderId).Update(&data); err != nil {
		log.Error("Insert Use Coupon Record", err.Error())
		return err
	}
	return nil
}

func UpdatePromotionCode(engine *database.MysqlSession, data entity.PromotionCode) error {
	data.Updated = time.Now()
	if _, err := engine.Session.Table(entity.PromotionCode{}).AllCols().Where("id = ?", data.Id).Update(&data); err != nil {
		log.Error("Update PromotionCode Error", err)
		return err
	}
	return nil
}

func CountUnuseCoupon(engine *database.MysqlSession, promoId string) (int64, error) {

	counts, err := engine.Engine.Table(entity.PromotionCode{}).Where("promotion_id=? AND is_used =?", promoId, 0).Count()
	if err != nil {
		log.Error("Count promo used code error", err.Error())
		return 0, err
	}
	return counts, err
}

func UpdateCouponIsCopy(engine *database.MysqlSession, promoId string, couponId string) error {
	_, err := engine.Session.Table(entity.PromotionCode{}).Cols("is_seller_copy").Where("promotion_id=? AND id =?", promoId, couponId).Update(entity.PromotionCode{
		IsSellerCopy: true,
	})
	if err != nil {
		log.Error("Update coupon is copy error", err.Error())
		return err
	}
	return nil

}

func GetUsedCouponRecordWithPromoId(engine *database.MysqlSession, promoId string) ([]entity.CouponUsedRecord, error) {
	var data []entity.CouponUsedRecord
	err := engine.Engine.Table(entity.CouponUsedRecord{}).Where("promotion_id =?", promoId).Find(&data)
	if err != nil {
		log.Error("Get Coupon used record")
		return data, err
	}
	return data, nil
}

func GetActivePromoCountWithStore(engine *database.MysqlSession, storeId string, now time.Time) (int, error) {
	var promos []entity.Promotion
	err := engine.Engine.Table(entity.Promotion{}).Where("store_id=?", storeId).Find(&promos)
	if err != nil {
		log.Error("Get Coupon used record")
		return 0, err
	}
	i := 0
	for _, promo := range promos {
		if !promo.StopTime.IsZero() {
			continue
		}
		if promo.StopTime.IsZero() && now.Sub(promo.CloseTime) > 0 {
			continue
		}
		i++
	}
	return i, nil
}
