package controllers

import (
	"api/config/middleware"
	"api/services/VO/Request"
	"api/services/entity"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/log"
	"api/services/util/response"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smokezl/govalidators"
)

type RequestSearchVO struct {
	Search  string `json:"search"`  //條件
	Start   string `json:"start"`   //啟始
	Length  string `json:"length"`  //幾筆
	OrderBy string `json:"orderBy"` //排序
	Sort    string `json:"sort"`    //DESC ASC
}

//新增商品
func NewProductPostAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.NewProductParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Debug("new => ", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
	//關閉上架
	//validator := govalidators.New()
	//if err := validator.LazyValidate(params); err != nil {
	//	resp.Fail(200, fmt.Sprintf("%v", err)).Send()
	//	return
	//}
	//StoreData := middleware.GetStoreData(ctx)
	//if StoreData.StoreStatus != Enum.StoreStatusSuccess {
	//	resp.Fail(errorMessage.GetMessageByCode(1001006)).Send()
	//	return
	//}
	//var filename []string
	//if len(params.ProductImage) != 0 {
	//	filename, _ = model.HandleProductImages(params.ProductImage)
	//}
	//data, err := model.HandleCreateProduct(StoreData, params, filename)
	//if err != nil {
	//	log.Debug("new product", err)
	//	resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
	//	return
	//}
	//var url = fmt.Sprintf("https://%s/product/%s", ctx.Request.Host, data.ProductId)
	//var imgPath = fmt.Sprintf("./www%s", StoreData.StorePicture)
	//img, err := images.GetImageFromFilePath(imgPath)
	//if err != nil {
	//	log.Error("Get image path Error", err, imgPath)
	//}
	//if err := qrcode.GeneratorQrCode(url, data.ProductId, img); err != nil {
	//	log.Error("generator qr code err ", err)
	//}
	//var product entity.NewProductResponse
	//product.ProductId = data.ProductId
	//
	//History.GenerateNewProductLog(data, StoreData.UserId)
	//resp.Success("OK").SetData(product).Send()
}

//商品編輯
func EditProductPostAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.EditProductParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Debug("new => ", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	day, _ := time.ParseInLocation("20060102 150405", "20210930 170000", time.Local)
	now := time.Now()
	if !now.Before(day) {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
	}
	validator := govalidators.New()
	if err := validator.LazyValidate(params); err != nil {
		log.Error("params Error", err.Error())
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	StoreData := middleware.GetStoreData(ctx)
	if err := model.ProductCheckStore(StoreData, params.ProductId); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	filename, err := model.HandleProductImages(params.ProductImage)
	if err != nil {
		resp.Fail(200, fmt.Sprintf("%v", err)).Send()
		return
	}
	if err := model.HandleEditProduct(StoreData, params, filename); err != nil {
		log.Debug("new product", err)
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	var product entity.NewProductResponse
	product.ProductId = params.ProductId

	resp.Success("OK").SetData(product).Send()
}

//取得商品列表
func GetProductsListAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var productList Request.ProductList
	err := ctx.BindQuery(&productList)
	store := middleware.GetStoreData(ctx)
	data, err := model.GetStoreProductList(store.StoreId, productList, 0)
	if err != nil {
		resp.Fail(200, fmt.Sprintf("%v", err)).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//取得帳單列表
func GetRealTimesListAction(ctx *gin.Context) {
	resp := response.New(ctx)
	StoreData := middleware.GetStoreData(ctx)
	var params Request.RealTimesRequest
	err := ctx.BindQuery(&params)

	data, err := model.GetStoreRealTimeList(StoreData.StoreId, params)
	if err != nil {
		resp.Fail(200, fmt.Sprintf("%v", err)).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//取得商品資訊
func GetProductAction(ctx *gin.Context) {
	resp := response.New(ctx)
	productId := ctx.Param("productId")
	storeData := middleware.GetStoreData(ctx)
	if err := model.ProductCheckStore(storeData, productId); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	productData, err := model.GetEditProductDataByProductId(storeData, productId, true, true)
	if err != nil {
		resp.Fail(200, fmt.Sprintf("%v", err)).Send()
		return
	}
	resp.Success("完成").SetData(productData).Send()
}

//批次設定商品運送方式
func BatchProductShippingAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.BatchProductShipping{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Debug("new => ", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	validator := govalidators.New()
	if err := validator.LazyValidate(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	err := model.BatchChangeShipping(params)
	if err != nil {
		resp.Fail(200, fmt.Sprintf("%v", err)).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

//批次設定商品運送方式
func BatchProductPayWayAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.BatchProductPayWay{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Debug("new => ", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	validator := govalidators.New()
	if err := validator.LazyValidate(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	err := model.BatchChangePayWay(params)
	if err != nil {
		resp.Fail(200, fmt.Sprintf("%v", err)).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

//批次設定商品運送方式
func BatchProductStatusAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.BatchProductStatus{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	validator := govalidators.New()
	if err := validator.LazyValidate(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	storeData := middleware.GetStoreData(ctx)
	if err := model.BatchChangeStatus(storeData, params); err != nil {
		resp.Fail(200, fmt.Sprintf("%v", err)).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

//設定帳單取消
func SetRealTimesCancelAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.SetRealTimesRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	err := model.HandleRealTimesCancel(params)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

//設定帳單延期
func SetRealTimesExtensionAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.SetRealTimesRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := model.HandleRealTimesExtension(params)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//新增收銀機
func CreateStoreAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.CreationStoreRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
	//關閉上架
	//userData := middleware.GetUserData(ctx)
	//if err := model.CreationStore(userData, params); err != nil {
	//	resp.Fail(1001001, err.Error()).Send()
	//	return
	//}
	//resp.Success("OK").SetData(true).Send()
}

//管理員列表
func ManagerListAction(ctx *gin.Context) {
	resp := response.New(ctx)
	userData := middleware.GetUserData(ctx)
	StoreData := middleware.GetStoreData(ctx)
	data, err := model.GetManagerList(userData, StoreData)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//新增管理者
func CreateManagerAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.CreationManagerRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
	//關閉上架
	//storeData := middleware.GetStoreData(ctx)
	//userData := middleware.GetUserData(ctx)
	//if err := model.CreationManager(userData, storeData, params); err != nil {
	//	resp.Fail(1001001, err.Error()).Send()
	//	return
	//}
	//resp.Success("OK").SetData(true).Send()
}

//重發邀請管理員
func InviteManagerAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.InviteManagerRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	StoreData := middleware.GetStoreData(ctx)
	if err := model.PutInviteManager(StoreData, params); err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

//刪除管理員
func DeleteManagerAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.DeleteManagerRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	storeData := middleware.GetStoreData(ctx)
	if err := model.DeleteManager(storeData, params); err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

//批次修改免運設定
func BatchProductFreeShipAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.BatchProductFreeShip{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	validator := govalidators.New()
	if err := validator.LazyValidate(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	err := model.BatchChangeFreeShip(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}
