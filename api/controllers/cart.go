package controllers

import (
	"api/config/middleware"
	"api/services/Service/Carts"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/log"
	"api/services/util/response"
	"github.com/gin-gonic/gin"
	"github.com/smokezl/govalidators"
)

func CheckCartAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.AddCartParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	cookie := middleware.GetSessionValue(ctx)
	if len(cookie) == 0 {
		return
	}
	if err := model.CheckProductAndCartSameStore(cookie, params); err != nil {
		resp.Success("OK").SetData(false).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}
/**
 * 購物車加入商品
 */
func AddCartAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.AddCartParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	validator := govalidators.New()
	if err := validator.LazyValidate(params); err != nil {
		resp.Fail(errorMessage.GetMessageFormatCode(1001002, err)).Send()
		return
	}
	cookie := middleware.GetSessionValue(ctx)
	if len(cookie) == 0 {
		resp.Success("OK").SetData(true).Send()
		return
	}
	if err := model.AddProductToCart(cookie, params); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

/**
 * 取出購物車資訊
 */
func GetCartAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.AddCartParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	cookie := middleware.GetSessionValue(ctx)
	if len(cookie) == 0 {
		log.Error("Cookie Error", cookie)
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	data, err := model.GetCartsData(cookie, params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

/**
 * 購物車刪除商品
 */
func DeleteCartAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.DeleteCartParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	cookie := middleware.GetSessionValue(ctx)
	if len(cookie) == 0 {
		log.Error("Cookie Error", cookie)
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	err := model.DeleteCartsProduct(cookie, params.ProductSpecId)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	carts, err := Carts.GetCarts(cookie)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	data, err := model.GetUserCartData(cookie, carts)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

/**
 * 變更運送方式
 */
func ChangeShippingAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.ChangeShippingParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	cookie := middleware.GetSessionValue(ctx)
	if len(cookie) == 0 {
		log.Error("Cookie Error", cookie)
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	if err := model.ChangeUserCartsShipping(cookie, params.ShipType); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	carts, err := Carts.GetCarts(cookie)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	data, err := model.GetUserCartData(cookie, carts)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

/**
 * 變更商品數量
 */
func ChangeQuantityAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.ChangeQuantityParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	cookie := middleware.GetSessionValue(ctx)
	if len(cookie) == 0 {
		log.Error("Cookie Error", cookie)
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	if err := model.ChangeUserCartProductQuantity(cookie, params.ProductSpecId, params.Type); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	carts, err := Carts.GetCarts(cookie)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	data, err := model.GetUserCartData(cookie, carts)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}
//輸入優惠卷
func ImportCouponNumberAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.CheckCouponNumberParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	cookie := middleware.GetSessionValue(ctx)
	if len(cookie) == 0 {
		log.Error("Cookie Error", cookie)
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	if err := model.CheckCartCoupon(cookie, params.CouponNumber);err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	carts, err := Carts.GetCarts(cookie)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	data, err := model.GetUserCartData(cookie, carts)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}
//刪除優惠卷
func DeleteCouponNumberAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.CheckCouponNumberParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	cookie := middleware.GetSessionValue(ctx)
	if len(cookie) == 0 {
		log.Error("Cookie Error", cookie)
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	if err := model.DeleteCartCoupon(cookie);err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	carts, err := Carts.GetCarts(cookie)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	data, err := model.GetUserCartData(cookie, carts)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}


/**
 * 取購物車目前數量
 */
func GetCartsCountAction(ctx *gin.Context) {
	resp := response.New(ctx)
	cookie := middleware.GetSessionValue(ctx)
	count, err := model.GetCartsCount(cookie)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	cartsCountResponse := Response.CartsCountResponse{
		Count: count,
	}
	resp.Success("OK").SetData(cartsCountResponse).Send()
}

/**
 * 取 Zipcode
 */
func GetShipZipcodeAction(ctx *gin.Context) {
	resp := response.New(ctx)
	data, err := model.GetZipcode()
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}
