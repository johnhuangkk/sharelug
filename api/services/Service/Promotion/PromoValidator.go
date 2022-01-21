package PromotionService

import (
	"api/services/VO/Request"
	"api/services/util/log"
	"strings"
	"time"
)

const (
	LimitMonth        = 6
	MaxDiscountAmount = 100000
)

func ValidateCreatePromotion(params Request.PromoCreate, tNow time.Time) int64 {
	if len([]rune(strings.TrimSpace(params.Name))) > 20 || len([]rune(strings.TrimSpace(params.Name))) == 0 {
		return 1011002
	}
	if params.Quantity > 1000 || params.Quantity == 0 {
		return 1011003
	}
	endTime, err := time.ParseInLocation("2006-01-02 15:04", params.EndTime, time.Local)
	if err != nil {
		log.Error("Promo Set End Time", params.EndTime, err.Error())
		return 1011006
	}
	if tNow.Sub(endTime) >= 0 {
		return 1011011
	}

	endTimeDate := time.Date(tNow.Year(), tNow.Month(), tNow.Day(), tNow.Hour(), tNow.Minute(), 0, 0, time.Local).AddDate(0, LimitMonth, 0)
	if endTimeDate.Sub(endTime).Hours()/24/30 < 0 {
		return 1011004
	}

	if params.Quantity*params.Amount > MaxDiscountAmount {
		return 1011005
	}
	return 0
}
