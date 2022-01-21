package controllers

import (
	"api/config/middleware"
	PromotionService "api/services/Service/Promotion"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/log"
	"api/services/util/response"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreatePromo(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.PromoCreate{}
	tNow := time.Now()
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error("post coupon params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}

	errorCode := PromotionService.ValidateCreatePromotion(params, tNow)
	if errorCode != 0 {
		resp.Fail(errorMessage.GetMessageByCode(errorCode)).Send()
		return
	}
	member := middleware.GetUserData(ctx)
	store := middleware.GetStoreData(ctx)
	flag, err := model.CheckPromoActiveLimit(store.StoreId, tNow)
	if err != nil {
		log.Error(err.Error())
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	if flag {
		resp.Fail(errorMessage.GetMessageByCode(1011010)).Send()
		return
	}
	_, err = model.CreatePromotion(params, member.Uid, store.StoreId, tNow)
	if err != nil {
		log.Error(err.Error())
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	resp.Success("OK").Send()
}
func EnablePromo(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.PromoEnable{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error("post coupon params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	member := middleware.GetUserData(ctx)
	store := middleware.GetStoreData(ctx)
	err := model.HandleStorePromoEnable(params, member.Uid, store.StoreId)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").Send()
}

func GetPromotion(ctx *gin.Context) {
	nowTime := time.Now
	resp := response.New(ctx)
	promoId := ctx.Param("promoId")
	var RespPromo Response.PromoList
	store := middleware.GetStoreData(ctx)
	data, err := model.GetPromotion(store.StoreId, promoId, nowTime())
	RespPromo.PromoEnable = store.EnablePromo
	RespPromo.Promos = data
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	resp.Success("OK").SetData(RespPromo).Send()
}

func TakePromotionCoupon(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.PromoCouponTake{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error("post coupon params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	qty, _ := strconv.Atoi(params.Quantity)
	if qty <= 0 || qty > 1000 {
		resp.Fail(errorMessage.GetMessageByCode(1011009)).Send()
		return
	}
	member := middleware.GetUserData(ctx)
	store := middleware.GetStoreData(ctx)
	err := model.TakeCoupons(qty, params.PromoId, member.Uid, store.StoreId)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").Send()
}

func StopPromotion(ctx *gin.Context) {
	resp := response.New(ctx)
	promoId := ctx.Param("promoId")
	err := model.TerminatePromotion(promoId)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	resp.Success("OK").Send()
}

func GetUsedCoupon(ctx *gin.Context) {
	resp := response.New(ctx)
	promoId := ctx.Param("promoId")
	params := Request.PromoCouponListGet{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error("Get coupon list params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := model.GetUsedCouponList(params, promoId)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

func GetUnuseCoupon(ctx *gin.Context) {
	resp := response.New(ctx)
	promoId := ctx.Param("promoId")
	params := Request.PromoCouponListGet{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error("Get coupon list params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := model.GetUnuseCouponList(params, promoId)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

func CopyCoupon(ctx *gin.Context) {
	resp := response.New(ctx)
	promoId := ctx.Param("promoId")
	couponId := ctx.Param("couponId")
	err := model.SetCouponIsCopy(promoId, couponId)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	resp.Success("OK").Send()
}

func GetUsedCouponExcel(ctx *gin.Context) {
	resp := response.New(ctx)
	promoId := ctx.Param("promoId")
	filename, err := model.GenerateUsedCouponRecordExcel(promoId)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	file, err := os.Open(filename) //Create a file
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	defer file.Close()
	ctx.Writer.Header().Add("Content-Disposition", "attachment; filename=Coupon-"+time.Now().Format("20060102")+".xlsx")
	ctx.Writer.Header().Add("Content-type", "application/octet-stream")
	_, err = io.Copy(ctx.Writer, file)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
}
