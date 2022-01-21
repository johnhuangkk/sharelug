package model

import (
	"api/services/Enum"
	"api/services/Service/Excel"
	"api/services/Service/Mail"
	"api/services/Service/Notification"
	"api/services/Service/OrderService"
	"api/services/Service/StoreService"
	"api/services/Service/SysLog"
	"api/services/Service/Upgrade"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Area"
	"api/services/dao/Balance"
	"api/services/dao/Orders"
	"api/services/dao/Store"
	"api/services/dao/SysLogDao"
	"api/services/dao/Withdraw"
	"api/services/dao/member"
	"api/services/dao/product"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"api/services/util/upload"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetStoreData(engine *database.MysqlSession, uid string) (entity.StoreData, error) {
	var storeData entity.StoreData
	storeData, err := Store.GetStoreDefaultDataByUid(engine, uid)
	if err != nil {
		log.Error("Get store Data Error", err)
		return storeData, err
	}
	return storeData, nil
}

func GetStoreDataByStoreId(engine *database.MysqlSession, StoreId string, userData entity.MemberData) (Response.StoreInfo, error) {
	var storeData Response.StoreInfo
	store, err := Store.GetStoreByUserIdAndStoreId(engine, userData.Uid, StoreId)
	if err != nil {
		log.Error("Get store Data Error", err)
		return storeData, fmt.Errorf("1002004")
	}
	if len(store.StoreId) == 0 {
		log.Error("Get not store Error", err, StoreId)
		return storeData, fmt.Errorf("1002004")
	}
	if store.RankStatus == Enum.StoreRankInit {
		return storeData, fmt.Errorf("1002005")
	}
	if store.RankStatus != Enum.StoreRankSuccess {
		return storeData, fmt.Errorf("1002004")
	}
	storeData.Sid = store.StoreId
	storeData.Name = store.StoreName
	storeData.IsExpire = false
	if store.StoreDefault != 1 {
		if !time.Now().Before(userData.UpgradeExpire) {
			storeData.IsExpire = true
		}
	}
	storeData.Picture = store.StorePicture
	storeData.Rank = store.Rank
	count, _ := Store.CountStoreRankByUid(engine, store.UserId)
	storeData.Count = count

	return storeData, nil
}

//取得收銀機列表
func GetStoreList(userData entity.MemberData) ([]entity.StoreDataResp, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := Store.GetStoreRankListByUid(engine, userData.Uid)
	if err != nil {
		log.Error("Get Store Rank Data Error", err)
		return nil, fmt.Errorf("系統錯誤！")
	}
	for k, v := range data {
		Seller, err := member.GetMemberDataByUid(engine, v.SellerId)
		if err != nil {
			log.Error("Get Member Data Error", err)
			return nil, fmt.Errorf("系統錯誤！")
		}
		if v.StoreDefault != 1 {
			if Seller.UpgradeExpire.IsZero() {
				data[k].StoreData.ExpireTime = time.Now().Add(time.Hour * time.Duration(24) * 7).Format("2006-01-02 15:04:05")
			} else {
				data[k].StoreData.ExpireTime = Seller.UpgradeExpire.Format("2006-01-02 15:04:05")
			}
		} else {
			data[k].StoreData.ExpireTime = ""
		}
	}
	return data, nil
}

func GetMyStore(StoreData entity.StoreDataResp, UserData entity.MemberData) (Response.MyStoreResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.MyStoreResponse
	balance, err := Balance.GetBalanceAccountLastByUserId(engine, StoreData.SellerId)
	if err != nil {
		log.Error("Get Balance Data Error", err)
		return resp, fmt.Errorf("系統錯誤")
	}
	//計算未回覆數
	reply, err := Orders.CountSellerOrderMessageBoardNotReply(engine, StoreData.StoreId)
	if err != nil {
		log.Error("Count Order Message Not reply Error", err)
		return resp, fmt.Errorf("系統錯誤！")
	}
	//計算待出貨數
	shipWait := Orders.CountShipWaitOrderData(engine, StoreData.StoreId)
	//計算待退貨數
	refundWait := Orders.GetRefundCount(engine, StoreData.StoreId, Enum.TypeRefund, Enum.RefundStatusWait)
	//fixme 應付帳單數

	resp.StoreName = StoreData.StoreName
	resp.Rank = StoreData.Rank
	resp.RankText = Enum.StoreRank[StoreData.Rank]
	resp.MemberPhone = tools.MaskerPhoneLater(UserData.Mphone)
	resp.BalanceAmount = strconv.Itoa(int(balance.Balance))
	resp.OrderWarn.CustomerMessage = reply
	resp.OrderWarn.ShipMessage = shipWait
	resp.OrderWarn.RefundMessage = refundWait
	resp.OrderWarn.BillMessage = 0

	now := time.Now()
	day, _ := time.ParseDuration("-24h")
	start := now.Add(day)

	activityList, err := member.GetAccountActivityByUserId(engine, UserData.Uid, start.Format("2006/01/02"), now.Format("2006/01/02"))
	if err != nil {
		log.Error("Get Account Activity Error", err)
		return resp, fmt.Errorf("系統錯誤！")
	}
	resp.AccountActivityList = composeActivity(activityList)
	return resp, nil
}

func composeActivity(data []entity.AccountActivityData) []Response.AccountActivityList {
	var resp []Response.AccountActivityList
	for _, v := range data {
		res := Response.AccountActivityList{}
		res.Message = fmt.Sprintf("%s : %s", Enum.ActivityStatus[v.Action], v.Message)
		res.Time = v.CreateTime.Format("2006/01/02")
		resp = append(resp, res)
	}
	return resp
}

func DownLoadSalesStatisticsReport(storeData entity.StoreDataResp, params *Request.SalesReportRequest) (string, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	start, end := setSalesParams(params)
	order, err := Orders.GetOrderDataByStoreIdAndDay(engine, storeData.StoreId, start, end)
	if err != nil {
		log.Error("Get Order Data Error", err)
		return "", fmt.Errorf("系統錯誤！")
	}
	var OrderReport []Excel.OrderReport
	for k, v := range order {
		detail, err := Orders.GetOrderDetailSingleByOrderId(engine, v.OrderId)
		if err != nil {
			log.Error("Get Order Data Error", err)
			return "", fmt.Errorf("系統錯誤！")
		}
		var report Excel.OrderReport
		report.Id = strconv.Itoa(k + 1)
		report.OrderTime = v.CreateTime.Format("2006/01/02")
		report.OrderId = v.OrderId
		report.ProductName = detail.ProductName
		report.ProductSpec = detail.ProductSpecName
		report.PaymentTime = ""
		if !v.PayWayTime.IsZero() {
			report.PaymentTime = v.PayWayTime.Format("2006/01/02")
		}
		report.PaymentType = Enum.PayWay[v.PayWay]
		report.OrderStatus = OrderService.GetOrderStatusText(v)
		report.ShipTime = ""
		if !v.ShipTime.IsZero() {
			report.ShipTime = v.ShipTime.Format("2006/01/02")
		}
		report.ShipType = Enum.Shipping[v.ShipType]
		report.ShipNumber = v.ShipNumber
		report.OrderAmount = int64(v.TotalAmount)
		report.ProductAmount = detail.Subtotal
		report.ShipFee = int64(v.ShipFee)
		report.PlatformFee = int64(v.PlatformPayFee + v.PlatformInfoFee + v.PlatformShipFee + v.PlatformTransFee)
		report.CreditAmount = int64(v.CaptureAmount)
		OrderReport = append(OrderReport, report)
	}

	filename, err := Excel.New().ToReportFile(OrderReport, start, end)
	if err != nil {
		return "", fmt.Errorf("系統錯誤！")
	}
	log.Debug("filename", filename)
	return filename, nil
}

//銷售統計
func GetSalesStatisticsReport(storeData entity.StoreDataResp, params *Request.SalesReportRequest) (Response.SalesReportResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()

	var resp Response.SalesReportResponse
	start, end := setSalesParams(params)

	resp.StartTime = start
	resp.EndTime = end

	//取訂單總數
	orderCount, err := Orders.CountOrderDataByStoreIdAndDay(engine, storeData.StoreId, start, end)
	if err != nil {
		log.Error("Count Order all Error", err)
		return resp, fmt.Errorf("系統錯誤！")
	}
	resp.SalesCount = orderCount
	//取訂單銷售金額
	orderSum, err := Orders.SumOrderDataByStoreIdAndDay(engine, storeData.StoreId, start, end)
	if err != nil {
		log.Error("sum Order all Error", err)
		return resp, fmt.Errorf("系統錯誤！")
	}
	resp.SalesAmount = orderSum
	//取訂單撥付金額
	ApprnSum, err := Orders.SumOrderAppropriationByStoreIdAndDay(engine, storeData.StoreId, start, end)
	if err != nil {
		log.Error("sum Order Appropriation Error", err)
		return resp, fmt.Errorf("系統錯誤！")
	}
	resp.CaptureAmount = ApprnSum
	//fixme 取應撥付金款
	RecCaptureAmount, err := Orders.SumOrderRecAppropriationByStoreIdAndDay(engine, storeData.StoreId, start, end)
	if err != nil {
		log.Error("sum Order Rec Appropriation Error", err)
		return resp, fmt.Errorf("系統錯誤！")
	}
	resp.RecCaptureAmount = RecCaptureAmount
	return resp, nil
}

func setSalesParams(params *Request.SalesReportRequest) (string, string) {
	var start string
	var end string
	now := time.Now()
	startDay, _ := time.ParseDuration(fmt.Sprintf("-%vh", 24*1))
	e := now.Add(startDay)
	switch params.Tab {
	case "OneDay":
		start = e.Format("2006-01-02")
		end = now.Format("2006-01-02")
	case "SevenDay":
		day, _ := time.ParseDuration(fmt.Sprintf("-%vh", 24*7))
		s := now.Add(day)
		start = s.Format("2006-01-02")
		end = now.Format("2006-01-02")
	case "TenDay":
		day, _ := time.ParseDuration(fmt.Sprintf("-%vh", 24*10))
		s := now.Add(day)
		start = s.Format("2006-01-02")
		end = now.Format("2006-01-02")
	case "Custom":
		start = params.StartTime
		end = params.EndTime
	}
	return start, end
}

func GetMyStoreInfo(storeData entity.StoreDataResp) (Response.SettingStoreResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.SettingStoreResponse
	userData, err := member.GetMemberDataByUid(engine, storeData.SellerId)
	if err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	ManagerMax, ManagerCount, err := Upgrade.ComputeManager(engine, userData, storeData)
	if err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}

	resp.StoreStatus = storeData.StoreStatus
	resp.FreeShip = storeData.FreeShip
	resp.FreeShipKey = storeData.FreeShipKey
	resp.SelfDelivery = storeData.SelfDelivery
	resp.SelfDeliveryFree = storeData.FreeSelfDelivery
	resp.SelfDeliveryKey = storeData.SelfDeliveryKey
	resp.StoreStatusText = Enum.StoreStatus[storeData.StoreStatus]
	resp.StoreName = storeData.StoreName
	resp.StorePicture = storeData.StorePicture
	resp.ManagerCount = ManagerCount
	resp.ManagerMax = ManagerMax

	stores, err := Store.GetStoresByUid(engine, storeData.SellerId)
	if err != nil {
		log.Error("Get Stores Error", err)
		return resp, fmt.Errorf("系統錯誤！")
	}
	for _, v := range stores {
		var res Response.StoreDataResp
		res.StoreId = v.StoreId
		res.SellerId = v.SellerId
		res.StoreName = v.StoreName
		res.StoreTax = v.StoreTax
		res.StorePicture = v.StorePicture
		res.StoreDefault = v.StoreDefault
		res.StoreStatus = v.StoreStatus

		Seller, err := member.GetMemberDataByUid(engine, v.SellerId)
		if err != nil {
			log.Error("Get Member Data Error", err)
			return resp, fmt.Errorf("系統錯誤！")
		}
		if v.StoreDefault != 1 {
			if Seller.UpgradeExpire.IsZero() {
				res.ExpireTime = time.Now().Add(time.Hour * time.Duration(24) * 7).Format("2006-01-02 15:04:05")
			} else {
				res.ExpireTime = Seller.UpgradeExpire.Format("2006-01-02 15:04:05")
			}
		} else {
			res.ExpireTime = ""
		}
		res.RankId = v.RankId
		res.UserId = v.UserId
		res.Rank = v.Rank
		res.RankStatus = v.RankStatus
		res.ManagerCount = Store.CountStoreManager(engine, v.StoreId)
		resp.MyStoreList = append(resp.MyStoreList, res)
	}
	return resp, nil
}

//修改收銀機資料
func SetStoreInfo(storeId string, params Request.SettingStoreRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()

	storeData, err := Store.GetStoreDataByStoreId(engine, storeId)
	if err != nil {
		log.Error("Get Store Data Error", err)
		return fmt.Errorf("系統錯誤")
	}
	if len(params.StorePicture) != 0 {
		filename, err := upload.StorePicture(params.StorePicture)
		if err != nil {
			return fmt.Errorf("系統錯誤")
		}
		storeData.StorePicture = "/static/images/store/" + filename
	}
	if len(params.StoreName) != 0 {
		storeData.StoreName = params.StoreName
	}
	if storeData.StoreStatus != Enum.StoreStatusSuspend {
		storeData.StoreStatus = Enum.StoreStatusSuccess
		//關閉收銀機
		if params.StoreStatus == 0 {
			storeData.StoreStatus = Enum.StoreStatusClose
			if err := SysLog.ShutdownStoreSystemLog(storeData.SellerId, storeData.StoreName); err != nil {
				log.Error("System Log Error", err)
			}
		}
		if err := Store.UpdateStoreData(engine, storeData.StoreId, storeData); err != nil {
			log.Error("Update Store Info Error", err)
			return fmt.Errorf("系統錯誤！")
		}
		//處理下架
		if storeData.StoreStatus == Enum.StoreStatusClose {
			err := product.UpdateProductStatus(engine, storeId, Enum.ProductStatusPending, Enum.StoreRankSuccess)
			if err != nil {
				log.Error("Update Product Status Error", err)
				return fmt.Errorf("系統錯誤！")
			}
		} else {
			err := product.UpdateProductStatus(engine, storeId, Enum.StoreRankSuccess, Enum.ProductStatusPending)
			if err != nil {
				log.Error("Update Product Status Error", err)
				return fmt.Errorf("系統錯誤！")
			}
		}
		return nil
	} else {
		return fmt.Errorf("此收銀機不得開啟")
	}

}

//取使用者資訊
func GetMyUserInfo(userId string) (Response.SettingUserResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.SettingUserResponse
	data, err := member.GetMemberDataByUid(engine, userId)
	if err != nil {
		log.Error("get user Info Error", err)
		return resp, fmt.Errorf("系統錯誤！")
	}
	resp.UserPhone = tools.MaskerPhone(data.Mphone)
	resp.Email = data.Email
	resp.Nickname = data.Username
	resp.Picture = data.Picture
	resp.VerifyEmail = data.VerifyEmail
	//todo 取銀行帳號
	withdraw, err := member.GetMemberWithdrawListByUserId(engine, userId)
	if err != nil {
		log.Error("get Withdraw Info Error", err)
		return resp, fmt.Errorf("系統錯誤！")
	}
	for _, v := range withdraw {
		var res Response.BankAccount
		res.AccountId = fmt.Sprintf("%v", v.Id)
		res.AccountNumber = fmt.Sprintf("%s / %s", v.BankName, tools.MaskerPhoneLater(v.Account))
		res.IsDefault = false
		if v.ActDefault == 1 {
			res.IsDefault = true
		}
		resp.BankAccount = append(resp.BankAccount, res)
	}
	//todo 取信用卡帳號
	card, err := member.GetMemberCreditDataByUserId(engine, userId)
	if err != nil {
		log.Error("get card Info Error", err)
		return resp, fmt.Errorf("系統錯誤！")
	}
	for _, v := range card {
		var res Response.CreditCard
		res.CardId = v.CardId
		res.CardNumber = fmt.Sprintf("**** **** **** %s", v.Last4Digits)
		res.IsDefault = false
		if v.DefaultCard == "1" {
			res.IsDefault = true
		}
		resp.CreditCard = append(resp.CreditCard, res)
	}
	return resp, nil
}

func SetMyUserInfo(userData entity.MemberData, params Request.SettingUserRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := member.GetMemberDataByUid(engine, userData.Uid)
	if err != nil {
		log.Error("get user Info Error", err)
		return fmt.Errorf("1001001")
	}
	//處理圖檔
	if len(params.UserPicture) != 0 {
		filename, err := upload.UserPicture(params.UserPicture)
		if err != nil {
			log.Error("Update user Picture Error", err)
			return fmt.Errorf("1001001")
		}
		data.Picture = "/static/images/user/" + filename
	}
	//更新暱稱
	if len(params.UserName) != 0 {
		data.Username = params.UserName
	}
	if len(params.UserEmail) != 0 && data.Email != params.UserEmail && data.VerifyEmail != params.UserEmail {
		if err := HandleSendVerifyEmail(engine, params.UserEmail, data.Uid, data.Username); err != nil {
			log.Error("Send Verify Email Error", err)
		}
		data.VerifyEmail = params.UserEmail
	}
	if _, err := member.UpdateMember(engine, &data); err != nil {
		log.Error("Update member data Error", err)
		return fmt.Errorf("1001001")
	}
	return nil
}

//發送EMAIL驗證
func HandleSendVerifyEmail(engine *database.MysqlSession, email, userId, userName string) error {
	//產生驗證網址
	link, err := Mail.GeneratorVerifyEmail(engine, email, userId, "", Enum.EmailVerifyTypeUser)
	if err != nil {
		log.Error("generator Verify Mail Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	//發信
	if err := Mail.SendVerifyMail(userName, email, link); err != nil {
		log.Error("Send Mail Verify Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	if err := SysLog.ChangeEmailSystemLog(userId, email); err != nil {
		log.Error("System Log Error", err)
	}
	return nil
}

//建立新收銀機
func CreationStore(userData entity.MemberData, params Request.CreationStoreRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	//取出升級訂單
	StoreLimit, StoreCurrent, err := Upgrade.ComputeStore(engine, userData)
	if StoreCurrent >= StoreLimit {
		return fmt.Errorf("收銀機數已達上限")
	}
	if len(params.StoreName) == 0 {
		return fmt.Errorf("請輸入收銀機名稱")
	}
	storePicture := ""
	if len(params.StorePicture) != 0 {
		filename, err := upload.StorePicture(params.StorePicture)
		if err != nil {
			return fmt.Errorf("系統錯誤")
		}
		storePicture = "/static/images/store/" + filename
	} else {
		storePicture = fmt.Sprintf("/static/img/store-%s.jpg", tools.RangeNumber(5, 2))
	}
	if err := SysLog.NewStoreSystemLog(userData.Uid, params.StoreName); err != nil {
		log.Error("System Log Error", err)
	}
	now := time.Now()
	if _, err = StoreService.CreateStoreData(engine, userData.Uid, params.StoreName, storePicture, now.Format("2006/01/02")); err != nil {
		return err
	}
	return nil
}

//取管理員列表
func GetManagerList(userData entity.MemberData, storeData entity.StoreDataResp) (Response.ManagerListResponse, error) {
	var resp Response.ManagerListResponse
	engine := database.GetMysqlEngine()
	defer engine.Close()

	data, err := Store.GetStoreManagerListByStoreId(engine, storeData.StoreId)
	if err != nil {
		return resp, err
	}
	for _, v := range data {
		var rep Response.ManagerList
		rep.ManagerId = v.StoreRank.RankId
		rep.ManagerEmail = v.StoreRank.Email
		rep.ManagerName = v.Member.Username
		rep.ManagerStatus = v.StoreRank.RankStatus
		rep.ManagerPicture = v.Member.Picture
		rep.ManagerStartTime = ""
		if v.StoreRank.RankStatus == Enum.StoreRankSuccess {
			rep.ManagerStartTime = v.StoreRank.UpdateTime.Format("2006/01/02")
		}
		rep.ManagerStatusText = Enum.StoreRankStatus[v.StoreRank.RankStatus]
		resp.ManagerList = append(resp.ManagerList, rep)
	}
	ManagerLimit, ManagerCurrent, err := Upgrade.ComputeManager(engine, userData, storeData)
	if err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	resp.StoreLimit = ManagerLimit
	resp.StoreCurrent = ManagerCurrent
	return resp, nil
}

//建立管理員
func CreationManager(userData entity.MemberData, storeData entity.StoreDataResp, params Request.CreationManagerRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if storeData.StoreStatus != Enum.StoreStatusSuccess {
		return fmt.Errorf("目前無法管理此收銀機！")
	}
	ManagerLimit, ManagerCurrent, err := Upgrade.ComputeManager(engine, userData, storeData)
	if err != nil {
		return fmt.Errorf("系統錯誤")
	}
	if ManagerCurrent >= ManagerLimit {
		return fmt.Errorf("管理員數已達上限")
	}
	//建立管理者
	userData, err = StoreService.CreateStoreManageData(engine, storeData.StoreId, params.ManagerPhone, params.ManagerEmail)
	if err != nil {
		return err
	}
	//FIXME 產生MAIL驗證URL
	Url, err := Mail.GeneratorVerifyEmail(engine, params.ManagerEmail, userData.Uid, storeData.StoreId, Enum.EmailVerifyTypeStore)
	if err != nil {
		log.Error("generator Verify Mail Error", err)
		//return fmt.Errorf("系統錯誤！")
	}
	//買家發送 認證MAIL
	if err := Mail.SendStoreVerifyMail(userData, storeData, params.ManagerEmail, Url); err != nil {
		log.Error("generator Verify Mail Error", err)
		//return err
	}
	if err := Notification.SendAddManagerMessage(engine, storeData.StoreId, storeData.SellerId); err != nil {
		return err
	}
	if err := SysLog.AssignManagerSystemLog(storeData.SellerId, storeData.StoreName, params.ManagerEmail); err != nil {
		log.Error("System Log Error", err)
	}
	return nil
}

//移除管理員
func DeleteManager(storeData entity.StoreDataResp, params Request.DeleteManagerRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if storeData.StoreStatus != Enum.StoreStatusSuccess {
		return fmt.Errorf("目前無法管理此收銀機！")
	}
	data, err := Store.GetStoreManagerAndMemberByManagerId(engine, params.ManagerId)
	if err != nil {
		log.Debug("Get Store Manager Error")
		return err
	}
	var manager entity.StoreRankData
	manager = data.StoreRank
	manager.RankStatus = Enum.StoreRankDelete
	if err = Store.UpdateStoreManagerData(engine, manager); err != nil {
		log.Debug("Update Store Manager Error")
		return err
	}
	if err := SysLog.CancelManagerSystemLog(storeData.SellerId, storeData.StoreName, data.StoreRank.Email); err != nil {
		log.Error("System Log Error", err)
	}
	//發賣家EMIL 簡訊 系統通知 買家 系統通知
	if err := Notification.SendDeleteManagerMessage(engine, storeData.StoreId, data.Member.Uid); err != nil {
		log.Error("Send Message Error", err)
	}
	return nil
}

//重發邀請管理員
func PutInviteManager(StoreData entity.StoreDataResp, params Request.InviteManagerRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if StoreData.StoreStatus != Enum.StoreStatusSuccess {
		return fmt.Errorf("目前無法管理此收銀機！")
	}
	data, err := Store.GetStoreManagerAndMemberByManagerId(engine, params.ManagerId)
	if err != nil {
		log.Debug("Get Store Manager Error")
		return err
	}
	storeData, err := Store.GetStoreDataByUserIdAndStoreId(engine, data.Member.Uid, data.StoreRank.StoreId)
	if err != nil {
		log.Debug("Get Store Data Error")
		return err
	}
	Url, err := Mail.GeneratorVerifyEmail(engine, data.StoreRank.Email, data.StoreRank.UserId, data.StoreRank.StoreId, Enum.EmailVerifyTypeStore)
	if err != nil {
		log.Error("generator Verify Mail Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	//發送MAIL
	if err := Mail.SendStoreVerifyMail(data.Member, storeData, data.StoreRank.Email, Url); err != nil {
		return err
	}

	if err := Notification.SendAddManagerMessage(engine, storeData.StoreId, storeData.SellerId); err != nil {
		return err
	}
	return nil
}

//帳戶重發認證信
func HandleInviteAccount(userData entity.MemberData) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	//產生驗證
	Url, err := Mail.GeneratorVerifyEmail(engine, userData.VerifyEmail, userData.Uid, "", Enum.EmailVerifyTypeUser)
	if err != nil {
		log.Error("generator Verify Mail Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	//發信
	if err := Mail.SendVerifyMail(userData.Username, userData.VerifyEmail, Url); err != nil {
		log.Error("Send Mail Verify Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	return nil
}

func GetSocialMediaInfo(storeData entity.StoreDataResp) (Response.StoreSocialMediaResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.StoreSocialMediaResponse
	data, err := Store.GetStoreSocialMediaDataByStoreId(engine, storeData.StoreId)
	if err != nil {
		log.Error("Get store Social Media Data Error", err)
		return resp, fmt.Errorf("1001001")
	}
	resp = data.GetStoreSocialMediaInfo()
	return resp, nil
}

func SetSocialMediaInfo(storeData entity.StoreDataResp, params Request.StoreSocialMediaRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := Store.GetStoreSocialMediaDataByStoreId(engine, storeData.StoreId)
	if err != nil {
		log.Error("Get store Social Media Data Error", err)
		return fmt.Errorf("1001001")
	}
	if len(params.Link) > 500 {
		return fmt.Errorf("1007009")
	}
	switch params.Type {
	case Enum.StoreSocialMediaTypeFacebook:
		data.FacebookLink = params.Link
		data.FacebookShow = params.Show
	case Enum.StoreSocialMediaTypeInstagram:
		data.InstagramLink = params.Link
		data.InstagramShow = params.Show
	case Enum.StoreSocialMediaTypeLine:
		data.LineLink = params.Link
		data.LineShow = params.Show
	case Enum.StoreSocialMediaTypeTelegram:
		data.TelegramLink = params.Link
		data.TelegramShow = params.Show
	}
	if len(data.StoreId) == 0 {
		data.StoreId = storeData.StoreId
		if err := Store.InsertStoreSocialMediaData(engine, data); err != nil {
			log.Error("insert store Social Media Data Error", err)
			return fmt.Errorf("1001001")
		}
	} else {
		if err := Store.UpdateStoreSocialMediaData(engine, data); err != nil {
			log.Error("insert store Social Media Data Error", err)
			return fmt.Errorf("1001001")
		}
	}
	return nil
}

//取出使用者操作記錄
func GetUserOperateRecord(userData entity.MemberData) ([]Response.UserOperateRecordResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	start := tools.GenerateBeforeTime(9)
	end := time.Now()
	data, err := SysLogDao.GetSystemLogByUserId(engine, userData.Uid, start.Format("2006-01-02 15:04"), end.Format("2006-01-02 15:04"))
	if err != nil {
		log.Error("Get System Log")
		return nil, fmt.Errorf("1001001")
	}
	var resp []Response.UserOperateRecordResponse
	for _, v := range data {
		var res Response.UserOperateRecordResponse
		res.RecordTime = v.CreateTime.Format("2006/01/02")
		if v.Action == Enum.ActivityCancelEmail {
			s := strings.Split(v.Content, "：")
			res.RecordContent = fmt.Sprintf("%s：%s", s[0], tools.MaskerEMail(s[1]))
		} else {
			res.RecordContent = v.Content
		}
		resp = append(resp, res)
	}
	return resp, nil
}

//取出免運設定
func HandleGetStoreFreeShip(storeData entity.StoreDataResp) (Response.StoreFreeShipResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := StoreService.GetStoreFreeShipping(engine, storeData.StoreId)
	if err != nil {
		return data, err
	}
	return data, nil
}

//寫入免運設定
func HandleSettingFreeShip(storeData entity.StoreDataResp, params Request.SettingFreeShipRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	//取出賣場資料
	store, err := Store.GetStoreDataByStoreId(engine, storeData.StoreId)
	if err != nil {
		log.Error("Get Store Data Error", err)
		return fmt.Errorf("1001001")
	}
	if params.FreeShipKey != Enum.FreeShipNone && params.FreeShip == 0 {
		return fmt.Errorf("1007011")
	}
	switch params.FreeShipKey {
	case Enum.FreeShipQuantity:
		store.FreeShipKey = Enum.FreeShipQuantity
	case Enum.FreeShipAmount:
		store.FreeShipKey = Enum.FreeShipAmount
	default:
		store.FreeShipKey = Enum.FreeShipNone
	}
	store.FreeShip = params.FreeShip
	//更新賣場資料
	if err := Store.UpdateStoreData(engine, store.StoreId, store); err != nil {
		log.Error("Update Store Info Error", err)
		return fmt.Errorf("1001001")
	}
	return nil
}

func HandleSelfDeliveryArea(request Request.SelfDeliveryAreaRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := Store.UpdateStoreSelfDeliveryArea(engine, request.StoreId, request.Enable); err != nil {
		log.Error("Update Store Info Error", err)
		return fmt.Errorf("1001001")
	}
	if !request.Enable {
		err := Store.DeleteStoreSelfDeliveryArea(engine, request.StoreId)
		if err != nil {
			log.Error("Update Store Info Error", err)
			return fmt.Errorf("1001001")
		}
		return nil
	}
	if len(request.Section) > 0 {
		for _, area := range request.Section {
			err := Store.InsertOrUpdateStoreSelfDeliveryArea(engine, request.StoreId, area.CityCode, area.AreaList)
			if err != nil {
				log.Error("Update Store Info Error", err)
			}
			continue
		}
	}
	return nil
}

func HandleSelfDeliveryChargeFree(request Request.SelfDeliveryFeeRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := Store.UpdateStoreSelfDeliveryChargeFree(engine, request.StoreId, request.SelfDeliveryFree, request.SelfDeliveryFreeShipKey); err != nil {
		log.Error("Update Store Info Error", err)
		return fmt.Errorf("1001001")
	}
	return nil
}

func HandleGetStoreSelfDeliveryArea(storeId string) ([]Response.CityWithArea, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	storeAreas, err := Store.GetStoreSelfDeliveryArea(engine, storeId)
	var resp []Response.CityWithArea
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}

	for _, area := range storeAreas {
		var cityArea Response.CityWithArea
		cityArea.CityCode = area.CityCode
		storeAreaList := strings.Split(area.AreaZip, ",")
		allArea, err := Area.GetTaiwanArea(area.CityCode)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		for _, cArea := range allArea {
			var rs Response.AreaFlagWithZipCode
			rs.Name = cArea.AreaName
			rs.ZipCode = cArea.ZipCode
			rs.Enable = false
			if tools.StringInSlice(cArea.ZipCode, storeAreaList) {
				rs.Enable = true
			}
			cityArea.Area = append(cityArea.Area, rs)
		}
		cityArea.CityName = allArea[0].CityName
		resp = append(resp, cityArea)
	}
	return resp, nil
}

func HandleStorePromoEnable(params Request.PromoEnable, userId string, storeId string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	//確認manager
	isManager, err := Store.CheckStoreManager(engine, userId, storeId)
	if err != nil {
		log.Error("Check Manager Error", err.Error())
		return fmt.Errorf("1001001")
	}

	if isManager {
		err := Store.UpdateStorePromoEnable(engine, storeId, params.Enable)
		if err != nil {
			log.Error("Update Store Promo Enable", err.Error())
			return fmt.Errorf("1001001")
		}
		return nil
	}
	return fmt.Errorf("1011007")
}

//EMAIL驗證
func HandleVerifyEmail(params Request.VerifyEmailParams, userData entity.MemberData) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := HandleSendVerifyEmail(engine, params.Email, userData.Uid, userData.Username); err != nil {
		return fmt.Errorf("1001001")
	}
	userData.VerifyEmail = params.Email
	if _, err := member.UpdateMember(engine, &userData); err != nil {
		log.Error("Update member data Error", err)
		return fmt.Errorf("1001001")
	}
	return nil
}

//
func HandleCheckVerifyEmail(userData entity.MemberData) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := member.GetMemberDataByUid(engine, userData.Uid)
	if err != nil {
		return fmt.Errorf("1001001")
	}
	if len(data.VerifyEmail) == 0 && len(data.Email) != 0 {
		return nil
	} else {
		return fmt.Errorf("1003008")
	}
}

//
func HandleGetIndustryList() ([]Response.IndustryCategory, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := Withdraw.GetIndustryData(engine)
	if err != nil {
		log.Error("Get Store Promo Enable", err.Error())
		return nil, fmt.Errorf("1001001")
	}
	var category []Response.IndustryCategory
	var cat Response.IndustryCategory
	var queue Response.IndustryVo

	for _, v := range data {
		if cat.Category != v.Category {
			if len(cat.Category) != 0 {
				category = append(category, cat)
			}
			cat.Category = v.Category
			cat.Industry = nil
		}
		queue.Industry = v.Industry
		queue.Mcc = v.IndustryId
		cat.Industry = append(cat.Industry, queue)
	}
	category = append(category, cat)
	return category, nil
}

func HandleSetStoreIndustry(userData entity.MemberData, params Request.StoreInfoRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	//判斷是否完成身份驗證
	if len(userData.Email) == 0 {
		return fmt.Errorf("1003011")
	}
	if userData.VerifyIdentity == 0 {
		return fmt.Errorf("1003012")
	}
	data, err := Withdraw.GetIndustryCode(engine, params.Industry)
	if err != nil {
		return fmt.Errorf("1003009")
	}
	if len(data.Industry) == 0 {
		return fmt.Errorf("1003009")
	}
	if len(strings.Split(params.Address, ",")) != 4 {
		return fmt.Errorf("1003010")
	}
	addr := strings.Split(params.Address, ",")
	tw, err := member.GetTaiwanCity(engine, tools.ChangeCityName(addr[1]))
	if err != nil {
		return fmt.Errorf("1003009")
	}
	userData.Representative = userData.IdentityName
	userData.RepresentativeId = userData.Identity
	userData.CompanyAddr = fmt.Sprintf("%s%s%s", addr[1], addr[2], addr[3])
	userData.CityCode = tw.Code
	userData.CityNameEn = tw.NameEn
	userData.JobCode = data.IndustryId
	userData.MccCode = data.Mcc
	userData.VerifyBusiness = 1
	userData.UpdateTime = time.Now()
	if _, err := member.UpdateMember(engine, &userData); err != nil {
		return fmt.Errorf("1001001")
	}
	return nil
}
