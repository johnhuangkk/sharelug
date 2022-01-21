package router

import (
	"api"
	"api/config/middleware"
	"api/controllers"
	"api/services/Service/Erp"
)

func SetupRouter(router *api.Website) *api.Website {

	//gw
	gateway := router.Group("/gw")
	{
		gateway.POST("/transfer/notify", controllers.TransferNotifyAction)
		gateway.POST("/credit/auth", controllers.AuthNotifyAction)
		ship := gateway.Group("/ship")
		{
			ship.POST("/family/notify/send", controllers.FamilySendAction)     //寄件即時通知
			ship.POST("/family/notify/leave", controllers.FamilyLeaveAction)   //寄件離店即時通知
			ship.POST("/family/notify/enter", controllers.FamilyEnterAction)   //進店即時通知
			ship.POST("/family/notify/pickup", controllers.FamilyPickupAction) //取件即時通知
			ship.POST("/family/notify/switch", controllers.FamilySwitchAction) //閉轉店即時通知
		}
		sevenMyship := gateway.Group("/seven")
		{
			sevenMyship.POST("/oltp", controllers.CreateOLTP)
		}
	}

	//api
	v1 := router.Group("/v1")
	{
		v1.PUT("/Kms", controllers.Kms)
		v1.GET("/TriggerOk", controllers.TriggerOk)
		v1.POST("/reCaptcha", controllers.InvoiceReCaptchaAction)
		v1.GET("/", controllers.IndexAction)
		//聯絡我們
		v1.POST("/contact", controllers.PostContactAction)
		//取TOKEN
		v1.PUT("/token", controllers.PostTokenAction)
		v1.PUT("/uuid", controllers.PostUuidAction)
		v1.PUT("/token/check", controllers.TokenCheckAction)
		//登入OTP
		v1.PUT("/login/otp", controllers.SendPayOtpAction)
		//登入驗證
		v1.PUT("/login/verify", controllers.ValidateOtpAction)
		//登出
		v1.PUT("/logout", controllers.LogoutAction)
		v1.GET("/memberData", middleware.Auth(), controllers.GetMemberDataAction)
		//取得商品資料
		v1.GET("/get/product/:productId", controllers.GetProductDataAction)
		//身份認證
		v1.POST("/twid/check", middleware.Auth(), controllers.PostTWIDCheckAction)
		// 取得託運單號
		v1.POST("/getShipNumber", middleware.Auth(), controllers.GetShipNumberAction)
		// 取得託運單
		v1.POST("/getConsignment", middleware.Auth(), controllers.GetConsignmentAction)
		// 查詢貨態資料
		v1.POST("/getShipStatus", middleware.Auth(), controllers.GetShipStatusAction)
		// 取得超商托運資料
		v1.GET("/getCvsShipData", middleware.Auth(), controllers.GetShipDataAction)
		// 閉轉店變更
		v1.PUT("/csvSwitchOrder", middleware.Auth(), controllers.PutCvsShipDataAction)
		//取小叮鈴的訊息數
		v1.GET("/getNotice", middleware.Auth(), controllers.GetNoticeAction)
		//上線通知我
		v1.POST("/message", controllers.PostNotifyMessageAction)
		//Email驗證
		v1.GET("/email/verify", controllers.GetEmailVerifyAction)
		v1.PUT("/email/verify", middleware.Auth(), controllers.MemberEmailVerifyAction)
		v1.POST("/member/company", controllers.MemberCompanyVerify)
		//購物車
		cart := v1.Group("/cart")
		{
			cart.PUT("/add", controllers.AddCartAction)
			cart.PUT("/check", controllers.CheckCartAction)
			cart.POST("/get", controllers.GetCartAction)
			cart.PUT("/delete", controllers.DeleteCartAction)
			cart.PUT("/ship/change", controllers.ChangeShippingAction)
			cart.PUT("/quantity/change", controllers.ChangeQuantityAction)
			cart.PUT("/coupon", controllers.ImportCouponNumberAction)
			cart.PUT("/coupon/delete", controllers.DeleteCouponNumberAction)
			cart.GET("/zipcode", controllers.GetShipZipcodeAction)
			cart.PUT("/otp", controllers.SendPayOtpAction)
			cart.PUT("/otp/verify", controllers.ValidateOtpAction)
			cart.POST("/pay", middleware.Auth(), controllers.PostPayAction)
			//cart.GET("/order/:orderId", middleware.Auth(), controllers.GetOrderDataAction)
			cart.GET("/count", controllers.GetCartsCountAction)
			//取信用卡資訊
			cart.GET("/card", middleware.Auth(), controllers.GetCardAction)
			//最後使用的地址
			cart.GET("/address", middleware.Auth(), controllers.GetDeliveryLastAddressAction)
		}

		pay := v1.Group("/pay")
		{
			pay.GET("/credit/3d/:content", controllers.CreditCheckAction)
			pay.POST("/credit/confirm/:pay", controllers.CreditConfirmAction)
			pay.GET("/credit/check", middleware.Auth(), controllers.GetOrderPaymentCheckAction)
		}

		store := v1.Group("/store")
		{
			//取得賣家商品列表(買家)
			store.GET("/product/list/:storeId", controllers.GetStoreProductsListAction)
			store.POST("/product/new/post", middleware.Auth(), controllers.NewProductPostAction)
			//商品編輯
			store.POST("/product/edit/post", middleware.Auth(), controllers.EditProductPostAction)
			store.GET("/product/edit/:productId", middleware.Auth(), controllers.GetProductAction)
			//賣家商品列表
			store.GET("/products", middleware.Auth(), controllers.GetProductsListAction)
			//取得賣家帳單列表
			store.GET("/realtime", middleware.Auth(), controllers.GetRealTimesListAction)
			//賣家帳單取消
			store.PUT("/realtime/cancel", middleware.Auth(), controllers.SetRealTimesCancelAction)
			//賣家帳單延期
			store.PUT("/realtime/extension", middleware.Auth(), controllers.SetRealTimesExtensionAction)
			//批次修改
			store.PUT("/product/batch/ship", middleware.Auth(), controllers.BatchProductShippingAction)
			store.PUT("/product/batch/payway", middleware.Auth(), controllers.BatchProductPayWayAction)
			store.PUT("/product/batch/status", middleware.Auth(), controllers.BatchProductStatusAction)
			//批次修改免運
			store.PUT("/product/batch/free", middleware.Auth(), controllers.BatchProductFreeShipAction)
			//訂單列表
			store.GET("/order/list", middleware.Auth(), controllers.GetSellerOrderListAction)
			store.GET("/ship/list", middleware.Auth(), controllers.GetSellerShipListAction)
			//訂單已讀處理
			store.PUT("/order/read", middleware.Auth(), controllers.PutOrderReadAction)
			//訂單MEMO
			store.PUT("/order/memo", middleware.Auth(), controllers.SetOrderMemo)
			//訂單退款列表
			store.GET("/order/refund/list", middleware.Auth(), controllers.GetRefundListAction)
			//取可退款
			store.GET("/order/refund", middleware.Auth(), controllers.GetRefundAction)
			//執行退款
			store.POST("/order/refund", middleware.Auth(), controllers.PostRefundAction)
			//訂單退貨列表
			store.GET("/order/return/list", middleware.Auth(), controllers.GetReturnListAction)
			//取可退貨
			store.GET("/order/return", middleware.Auth(), controllers.GetReturnAction)
			//執行退貨
			store.POST("/order/return", middleware.Auth(), controllers.PostReturnAction)
			//退貨完成確認
			store.PUT("/order/return/confirm", middleware.Auth(), controllers.PostReturnConfirmAction)
			//執行 設定運送 方式和運送單號
			store.PUT("/order/set/ship/number", middleware.Auth(), controllers.SetOrderShippingAction)
			//宅配匯出批次出貨單
			store.PUT("/order/export/ship", middleware.Auth(), controllers.ExportOrderShippingAction)
			//宅配匯出批次出貨單PDF
			store.PUT("/order/export/pdf", middleware.Auth(), controllers.ExportOrderShippingPdfAction)
			//宅配匯出批次交寄單PDF
			store.PUT("/order/export/delivery", middleware.Auth(), controllers.ExportOrderShipSendAction)
			//宅配匯入批次出貨單
			store.POST("/order/import/ship", middleware.Auth(), controllers.ImportBatchShippingAction)
			store.GET("/order/import/batch/:batchId", middleware.Auth(), controllers.ProcessBatchShippingAction)
			//面交匯出批次出貨單
			//store.PUT("/order/export/f2f", middleware.Auth(), controllers.ExportOrderF2fAction)
			//store.PUT("/order/export/f2f/pdf", middleware.Auth(), controllers.ExportOrderF2fPdfAction)
			////退貨退款列表
			store.GET("/refund/list", middleware.Auth(), controllers.GetOrderRefundListAction)
			//我的收銀機
			store.GET("/my", middleware.Auth(), controllers.GetMyStoreListAction)
			//切換收銀機
			store.PUT("/my/exchanges", middleware.Auth(), controllers.ExchangeStoreAction)
			//我的收銀機首頁
			store.GET("/my/index", middleware.Auth(), controllers.GetMyStoreAction)
			//銷售統計
			store.GET("/my/report", middleware.Auth(), controllers.GetSalesReportAction)
			store.GET("/my/report/download", middleware.Auth(), controllers.DownLoadSalesReportAction)
			//使用者操作記錄
			store.GET("/my/operate", middleware.Auth(), controllers.UserOperateRecordAction)
			//帳戶明細
			store.GET("/my/balance", middleware.Auth(), controllers.MyBalanceAction)
			//帳戶明細下載
			store.GET("/my/balance/export", controllers.ExportMyBalanceAction)
			//保留款項
			store.GET("/my/retain", middleware.Auth(), controllers.MyRetainAction)
			//帳務管理
			store.GET("/my/account", middleware.Auth(), controllers.MyAccountAction)
		}

		iPostBox := v1.Group("/iPostBox")
		{
			// 郵局呼叫
			iPostBox.POST("/notification", middleware.AllowIP(), controllers.IPostShipStatusUpdateNotification)
		}

		v1.POST("/HiLife", controllers.HiLifeNotification)
		v1.GET("/HiLife", controllers.HiLifeNotification)
		//v1.GET("/SevenEleven", controllers.HiLifeNotification)
		//v1.POST("/SevenEleven", controllers.HiLifeNotification)

		v1.GET("/MartPrintTest", controllers.MartPrintOrder)
		v1.POST("/MartSwitch", controllers.MartSwitchOrder)
		v1.POST("/MartFetching", controllers.MartFetching)

		address := v1.Group("/address")
		{
			address.POST("/add", middleware.Auth(), controllers.AddAddressAction)
			address.PUT("/delete", middleware.Auth(), controllers.DeleteAddressAction)
			address.GET("/receiveAddress", middleware.Auth(), controllers.GetReceiveAddressAction)
			address.GET("/receiveAddress/:ship", middleware.Auth(), controllers.GetReceiveAddressAction)
			address.GET("/sendAddress", middleware.Auth(), controllers.GetSendAddressAction)
			address.GET("/checkSendAddress/:ship", middleware.Auth(), controllers.CheckShipSendAddressExistAction)
		}

		order := v1.Group("/order")
		{
			order.GET("/messageBoardInfo/:orderId", middleware.Auth(), controllers.GetBuyerAndStoreDataAction)
			order.GET("/getMessage/:orderId", middleware.Auth(), controllers.GetOrderMessageBoardAction)
			order.PUT("/addMessage", middleware.Auth(), controllers.AddOrderMessageBoardAction)
			order.GET("/seller/messageBoard/list", middleware.Auth(), controllers.GetSellerOrderMessageBoardListAction)
			order.GET("/buyer/messageBoard/list", middleware.Auth(), controllers.GetBuyerOrderMessageBoardListAction)
		}
		//取得買家訂單列表
		orders := v1.Group("/orders")
		{
			orders.GET("/detail/:orderId", middleware.Auth(), controllers.GetOrderDataAction)
			orders.GET("/list", middleware.Auth(), controllers.GetBuyerOrderListAction)
			set := orders.Group("/set")
			{
				//執行提前付款
				set.PUT("/advance", middleware.Auth(), controllers.SetAdvancePaymentAction)
				//延長撥付
				set.PUT("/extension", middleware.Auth(), controllers.SetExtensionPaymentAction)
				//訂單完成交易 fixme
				set.PUT("/confirm", middleware.Auth(), controllers.SetConfirmPaymentAction)
				//執行訂單取消交易
				set.PUT("/cancel", middleware.Auth(), controllers.SetCancelOrderAction)
			}
		}
		customer := v1.Group("/customer")
		{
			//聯絡客服
			customer.GET("", middleware.Auth(), controllers.GetCustomerAction)
			customer.POST("", middleware.Auth(), controllers.PostCustomerAction)
			customer.GET("/:questionId", middleware.Auth(), controllers.GetCustomerQuestion)
		}
		//設定
		setting := v1.Group("/set")
		{
			//取設定收銀機
			setting.GET("/store", middleware.Auth(), controllers.SettingStoreAction)
			setting.PUT("/store", middleware.Auth(), controllers.PutSettingStoreAction)
			//取出免運設定
			setting.GET("/shipping", middleware.Auth(), controllers.SettingFreeShipAction)
			setting.PUT("/shipping", middleware.Auth(), controllers.PutSettingFreeShipAction)
			//取社群帳號連結
			setting.GET("/media", middleware.Auth(), controllers.SettingSocialMediaAction)
			setting.PUT("/media", middleware.Auth(), controllers.PutSettingSocialMediaAction)
			//取資料設定
			setting.GET("/user", middleware.Auth(), controllers.SettingAccountAction)
			setting.PUT("/user", middleware.Auth(), controllers.PutSettingAccountAction)
			setting.PUT("/user/invite", middleware.Auth(), controllers.InviteAccountAction)
			// 刪除信用卡
			setting.PUT("/credit/delete", middleware.Auth(), controllers.DeleteCreditAction)
			// 變更信用卡預設
			setting.PUT("/credit/default", middleware.Auth(), controllers.ChangeDefaultCreditAction)
			//檢查會員EMAIL是否認證
			setting.GET("/member/email/verify", middleware.Auth(), controllers.CheckMemberEmailVerifyAction)
			//訊息通知列表
			setting.GET("/notify", middleware.Auth(), controllers.NotificationListAction)
			//訊息通知已讀
			setting.PUT("/notify", middleware.Auth(), controllers.NotificationReadAction)
			//收銀機升級
			setting.GET("/upgrade", middleware.Auth(), controllers.UpgradeAction)
			//收銀機升級取結帳
			setting.POST("/upgrade/get", middleware.Auth(), controllers.GetUpgradePayAction)
			//收銀機升級付款
			setting.POST("/upgrade", middleware.Auth(), controllers.UpgradePayAction)
			//升級訂單未付款數
			setting.GET("/upgrade/count", middleware.Auth(), controllers.CountUpgradeOrderAction)
			//升級訂單未付款列表
			setting.GET("/upgrade/list", middleware.Auth(), controllers.GetUpgradeOrderListAction)
			//取升級訂單
			setting.GET("/upgrade/order/:orderId", middleware.Auth(), controllers.GetUpgradeOrderAction)
			//新增收銀機
			setting.POST("/create/store", middleware.Auth(), controllers.CreateStoreAction)
			//新增管理者
			setting.GET("/list/manager", middleware.Auth(), controllers.ManagerListAction)
			setting.POST("/create/manager", middleware.Auth(), controllers.CreateManagerAction)
			//重發邀請信
			setting.PUT("/invite/manager", middleware.Auth(), controllers.InviteManagerAction)
			//刪除管理員
			setting.PUT("/delete/manager", middleware.Auth(), controllers.DeleteManagerAction)
			//取出發票載具資料
			setting.GET("/carrier/get", middleware.Auth(), controllers.GetCarrierAction)
			//設定發票載具資料
			setting.PUT("/carrier/set", middleware.Auth(), controllers.PostCarrierAction)
			setting.GET("/invoice/list", middleware.Auth(), controllers.GetInvoiceListAction)
			setting.GET("/invoice/detail", middleware.Auth(), controllers.GetInvoiceDetailAction)

			setting.PUT("/seller/email/verify", middleware.Auth(), controllers.VerifyEmailAction)
			setting.GET("/seller/email/check", middleware.Auth(), controllers.VerifyCheckEmailAction)
			setting.GET("/seller/industry", middleware.Auth(), controllers.GetIndustryListAction)
			setting.PUT("/seller/industry", middleware.Auth(), controllers.SetIndustryAction)
			//外送
			setting.GET("/store/:storeId/self-delivery", controllers.GetStoreSelfDelivery)
			setting.POST("/store/self-delivery/charege-free", middleware.Auth(), controllers.SetStoreSelfDeliveryChargeFree)
			setting.POST("/store/self-delivery", middleware.Auth(), controllers.SetStoreSelfDelivery)
			setting.GET("/self-delivery/city/taiwan", middleware.Auth(), controllers.GetDeliveryCities)
			setting.GET("/self-delivery/city/taiwna/:cityCode/area", middleware.Auth(), controllers.GetDeliveryArea)

			//優惠活動相關設定
			setting.POST("/promo", middleware.Auth(), controllers.CreatePromo)
			setting.POST("/promo/enable-promo", middleware.Auth(), controllers.EnablePromo)
			setting.GET("/promo", middleware.Auth(), controllers.GetPromotion)
			setting.GET("/promo/:promoId", middleware.Auth(), controllers.GetPromotion)
			setting.PUT("/promo", middleware.Auth(), controllers.TakePromotionCoupon)
			setting.PUT("/promo/:promoId/stop-promo", middleware.Auth(), controllers.StopPromotion)
			setting.PUT("/promo/:promoId/coupon/unuse", middleware.Auth(), controllers.GetUnuseCoupon)
			setting.PUT("/promo/:promoId/coupon/used", middleware.Auth(), controllers.GetUsedCoupon)
			setting.PUT("/promo/:promoId/coupon/is-copy/:couponId", middleware.Auth(), controllers.CopyCoupon)
			setting.GET("/promo/:promoId/coupon/report", middleware.Auth(), controllers.GetUsedCouponExcel)
		}
		//反向帳單
		bill := v1.Group("/bill")
		{
			bill.POST("", middleware.Auth(), controllers.PostBillOrderAction)
			bill.GET("/review/:billId", middleware.Auth(), controllers.ReviewBillOrderAction)
			bill.GET("/detail/:billId", controllers.GetBillOrderAction)
			bill.PUT("/confirm", middleware.Auth(), controllers.BillConfirmAction)
			bill.GET("/list", middleware.Auth(), controllers.GetBillListsAction)
			bill.PUT("/extension", middleware.Auth(), controllers.BillExtensionAction)
			bill.PUT("/cancel", middleware.Auth(), controllers.BillCancelAction)
			bill.GET("/all", middleware.Auth(), controllers.GetAllBillListsAction)
		}
		withdraw := v1.Group("/withdraw")
		{ //提領
			withdraw.POST("", middleware.Auth(), controllers.WithdrawAction)
			withdraw.GET("/bank", middleware.Auth(), controllers.WithdrawBankCodeAction)
			withdraw.PUT("/delete", middleware.Auth(), controllers.WithdrawDeleteAction)
			withdraw.PUT("/default", middleware.Auth(), controllers.WithdrawChangeDefaultAction)
		}
		short := v1.Group("/short")
		{
			short.GET("/:short", controllers.ShortUrlAction)
		}

		erp := v1.Group("/erp")
		{
			//搜尋訂單列表
			erp.PUT("/orders", controllers.SearchOrdersAction)
			//搜尋訂單內容
			erp.GET("/order", controllers.GetOrderDetailAction)
			//搜尋訂單退款內容
			erp.GET("/refund", controllers.GetRefundAndReturnAction)
			//搜尋訂單運送資訊
			erp.GET("/shipping", controllers.GetShippingAction)
			erp.GET("/payment", controllers.GetSearchPaymentAction)
			//搜尋商品列表
			erp.PUT("/products", controllers.GetProductsAction)
			//搜尋會員資料
			erp.PUT("/member", controllers.SearchMemberAction)
			erp.GET("/member/:account/store", controllers.MemberStore)
			//賣場相關訊息

			erp.GET("/ach", controllers.ExporterKgiAchAction)
			erp.GET("/each", controllers.ExporterKgiEachAction)
			erp.POST("/each/upload", controllers.EachResponseAction)
			//審單
			erp.PUT("/audit/release", controllers.CreditAuditAction)
			erp.PUT("/audit/memo", controllers.CreditAuditMemoAction)
			erp.PUT("/audit/list", controllers.CreditAuditListAction)
			//中止升級服務
			erp.PUT("/user/suspend", controllers.SuspendUpgradeAction)

			erp.GET("/credit/capture", controllers.CreditCaptureAction)
			//重新讀取信用卡回應檔
			erp.GET("/credit/read/respond", controllers.CreditAgainReadRespondAction)
			erp.GET("/export/member", controllers.GetMemberReportAction)
			erp.GET("/export/order", controllers.GetOrderReportAction)
			erp.GET("/export/daily/:day", controllers.ExporterDayStatementAction)
			erp.GET("/export/invoice/:day", controllers.ExporterInvoiceReportAction)
			erp.GET("/export/invoices/:id", controllers.ExporterUserInvoiceReportAction)
			erp.GET("/export/balances", controllers.ExporterBalancesReportAction)
			//上傳次特店代碼表
			erp.POST("/import/special", controllers.ImportSpecialStoreAction)
			// 匯出銀行賣家資料
			erp.GET("/export/bank", controllers.ExporterBankReportAction)
			//取出提領列表
			erp.PUT("/withdraw", controllers.GetWithdrawAction)
			erp.PUT("/withdraw/change", controllers.GetWithdrawChangeStatusAction)
			//重送發票
			erp.PUT("/invoice/resend", controllers.ResendInvoiceAction)
			erp.GET("/withdraw", controllers.GetWithdrawAction)
			// 超商寄件核帳
			erp.GET("/cvsSendChecked", controllers.GetCvsSendCheckedAction)
			// 超商取件核帳
			erp.GET("/cvsPickUpChecked", controllers.GetCvsPickUpCheckedAction)
			// erp.POST("/login", Erp.Login)
			// erp.POST("/create", Erp.CreateErpUser)
			erp.POST("/login", Erp.Login)
			erp.POST("/create", Erp.CreateErpUser)
			erp.GET("/customer", controllers.FetchCustomerAction)
			erp.GET("/order/:orderId", controllers.FetchErpOrder)
			erp.POST("/customer/question/:questionId", controllers.ReplyCustomer)
			erp.PUT("/customer/question/:questionId/pending", controllers.PendingCustomer)
			erp.PUT("/customer/question/:questionId/finish", controllers.FinishCustomer)
			erp.POST("/customer/question/:questionId/memo", controllers.CustomerMemo)
			//erp.PUT("/customer/reply", controllers.SendCustomerReply)
			erp.GET("/customer/count", controllers.FetchCountCustomerQuestionsByStatus)
			erp.GET("/history/question/", controllers.FetchCustomerRelatedQuestions)
			//Erp身份認證
			erp.POST("/check/twId", controllers.ErpTWIDCheckAction)
			//更改帳戶為公司戶
			erp.GET("/company/verify-pending", controllers.MemberCompanyVerifyPendingList)
			erp.GET("/company/verify-pending/:memberId", controllers.MemberSpecialStoreVerifyInfo)
			erp.POST("/company/:memberId", controllers.UpdateMemberCompany)
			erp.GET("/company/verify-send", controllers.MemberSpecialStoreVerifyList)
			erp.PUT("/company/verify-send", controllers.FindMemberSpecialStoreVerifyRecord)
			erp.POST("/member/company", controllers.MemberCompanyVerify)
			// erp.GET("/test", controllers.GetMemberSendToBank)
			erp.GET("/ftp/member-kgi", controllers.GetMemberSendToKgiFtp)
			erp.PUT("/ftp/kgi", controllers.GetKgiSpecialStoreExcel)
			//手動跑轉帳被動查詢
			erp.GET("/transfer/check/:date", controllers.SearchTransferAction)
			erp.GET("/credit/check/:OrderId", controllers.SearchCreditAction)
			erp.GET("/credit/cancel/:OrderId", controllers.CancelCreditOrderAction)
			erp.GET("/ship/order/:OrderId", controllers.SetOrderShipAction)

			erp.POST("/platform/msg", controllers.SendPlatformMessageAction)
			erp.GET("/again/daily/:day", controllers.AgainDailyReportAction)
			erp.GET("/refund/platform/:orderId", controllers.RefundPlatformFeeAction)
			erp.POST("/area/taiwan", controllers.CreateTaiwanArea)
			erp.POST("/mail/check", controllers.PostMailAction)
			//上傳發票字軌
			erp.POST("/invoice/track", controllers.ImportInvoiceTrackAction)
			//上傳發票開獎結果檔
			erp.POST("/invoice/awarded", controllers.ImportAwardedAction)
			erp.GET("/invoice/awarded/:orderId", controllers.GetAwardedInvoiceAction)
			//發送EMD 簡訊
			erp.POST("/edm/send/message", controllers.SendEdmMessageAction)
			erp.POST("/edm/generate/short", controllers.GenerateShortAction)

			erp.PUT("/translate/addr", controllers.TranslateAddress)
			//查詢會員餘額
			erp.PUT("/balances", controllers.SearchUserBalanceAction)
			//重新計算會員待撥付餘額
			erp.GET("/balances/retain/recalculate/:uid", controllers.RecalculateBalanceRetainAction)
			erp.GET("/balances/balances/recalculate/:uid", controllers.RecalculateBalanceAction)
		}
		sevenMyship := v1.Group("/seven")
		{
			sevenMyship.GET("/address", controllers.GetCityShopsAddress)
			sevenMyship.GET("/shop", controllers.FetchDailyShops)
			sevenMyship.GET("/test", controllers.Test)
		}

		invoice := v1.Group("/invoice")
		{
			invoice.POST("/carrier/bind", controllers.InvoiceBindCarrierAction)
			invoice.POST("/carrier/verify", middleware.Auth(), controllers.InvoiceVerifyCarrierAction)
		}

	}

	//修改退款 要審退款

	//電子發票綁定
	//router.POST("/invoice/carrier", middleware.Auth(), controllers.InvoiceBindPlatformAction)
	//router.GET("/invoice/carrier", middleware.Auth(), controllers.CarrierAction)
	//router.POST("/setting/invoice", middleware.Auth(), controllers.InvoiceAction)
	//router.GET("/setting/invoice", middleware.Auth(), controllers.InvoiceAction)
	//router.GET("/setting/invoicebind", middleware.Auth(), controllers.InvoiceBindAction)
	////商品頁
	//router.GET("/product/:productId", middleware.Auth(), controllers.ProductAction)
	////結帳頁
	//router.GET("/pay", middleware.Auth(), middleware.GetDevSiteAllow, controllers.GetPayAction)
	//router.GET("/expire", middleware.Auth(), controllers.GetExpireAction)
	//router.GET("/pay/succ/:orderId", middleware.Auth(), middleware.GetDevSiteAllow, controllers.GetSuccess)
	//router.GET("/pay/fail", middleware.Auth(), middleware.GetDevSiteAllow, controllers.GetFail)
	////router.GET("/fcm", controllers.PushAction)
	//

	//router.GET("/login", middleware.Auth(), middleware.GetDevSiteAllow, controllers.LoginGetAction)
	//router.GET("/login/otp/:phone", middleware.Auth(), middleware.GetDevSiteAllow, controllers.LoginOtpGetAction)
	//router.PUT("/otp/send", middleware.Auth(), controllers.SendOtpAction)
	//router.PUT("/otp/validate", middleware.Auth(), controllers.ValidateOtpAction)
	//
	//router.GET("/logout", middleware.Auth(), middleware.GetDevSiteAllow, controllers.LogoutAction)
	////router.GET("/register", controllers.RegisterGetAction)
	////router.POST("/register", controllers.RegisterPostAction)

	//測試頁
	simulator := router.Group("/simulator")
	{
		simulator.Use(middleware.GetDevSiteAllow)
		simulator.GET("/sms", controllers.GetSmsAction)
		simulator.POST("/post/sms", controllers.PostSmsAction)
	}
	return router
}
