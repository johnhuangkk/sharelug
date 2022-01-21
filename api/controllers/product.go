package controllers

import (
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/response"
	"fmt"
	"github.com/gin-gonic/gin"
)

func GetProductDataAction(ctx *gin.Context) {
	resp := response.New(ctx)
	productId := ctx.Param("productId")
	var productData Response.ProductResponse
	var err error
	if len(productId) > 10 {
		productData, err = model.GetProductDataByProductId(productId, false, false)
	} else {
		productData, err = model.GetProductDataByTinyUrl(productId)
	}
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("完成").SetData(productData).Send()
}

//商品列表 買家
func GetStoreProductsListAction(ctx *gin.Context) {
	resp := response.New(ctx)
	storeId := ctx.Param("storeId")

	var productList Request.ProductList
	err := ctx.BindQuery(&productList)

	data, err := model.GetStoreProductList(storeId, productList, 0)
	//data, err := model.GetProductList(storeId,false,0,20,0, ctx.Request.Host)
	if err != nil {
		resp.Fail(200, fmt.Sprintf("%v", err)).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}
