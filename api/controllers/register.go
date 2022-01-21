package controllers

//註冊頁
//func RegisterGetAction(ctx *gin.Context) {
//	ctx.HTML(http.StatusOK, "index/register.html", gin.H{
//		"title": "會員註冊",
//	})
//}

//註冊
//func RegisterPostAction(ctx *gin.Context) {
//	resp := response.New(ctx)
//	params := &model.RegisterParams{}
//
//	if err := ctx.ShouldBind(&params); err != nil {
//		log.Error("Register post params error", err)
//		resp.Fail(200, "欄位未填寫完整").Send()
//		return
//	}
//	validator := govalidators.New()
//	if err := validator.LazyValidate(params); err != nil {
//		resp.Fail(200, fmt.Sprintf("%v", err)).Send()
//		return
//	}
//
//	if err := model.ValidatePassword(params); err != nil {
//		resp.Fail(200, fmt.Sprintf("%v", err)).Send()
//		return
//	}
//
//	//if err := model.ValidateEmailIsExists(params); err != nil {
//	//	resp.Fail(200, fmt.Sprintf("%v", err)).Send()
//	//	return
//	//}
//
//	//if err := model.CreateMember(params); err != nil {
//	//	resp.Fail(200, "系統錯誤").Send()
//	//	return
//	//}
//
//	resp.Success("完成註冊。").Send()
//}
