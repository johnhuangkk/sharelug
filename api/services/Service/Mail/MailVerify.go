package Mail

import (
	"api/services/Enum"
	"api/services/Service/SendMail"
	"api/services/Service/Upgrade"
	"api/services/dao/Email"
	"api/services/dao/member"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

func EmailVerify(engine *database.MysqlSession, userId string, VerifyEmail string) error {
	userData, err := member.GetMemberDataByUid(engine, userId)
	if err != nil {
		log.Error("Get Member data Error [%s]", err)
		return err
	}
	userData.VerifyEmail = VerifyEmail

	Url, err := GeneratorVerifyEmail(engine, VerifyEmail, userData.Uid, "", Enum.EmailVerifyTypeUser)
	if err != nil {
		log.Error("generator Verify Mail Error", err)
		return err
	}
	err = SendVerifyMail("", userData.VerifyEmail, Url)
	if err != nil {
		log.Error("Send Mail Verify Error", err)
		return err
	}
	_, err = member.UpdateMember(engine, &userData)
	if err != nil {
		log.Error("Update member data Error", err)
		return err
	}
	return nil
}

func GeneratorVerifyEmail(engine *database.MysqlSession, email, userId, storeId, verifyType string) (string, error) {
	code := tools.GeneratorValidationCode(email)
	link := fmt.Sprintf("https://%s/verify/mail?code=%s&type=%s", viper.GetString("WEBHOST"), code, verifyType)
	var data entity.EmailVerifyData
	data.Email = email
	data.UserId = userId
	data.StoreId = storeId
	data.VerifyCode = code
	data.SendTime = time.Now()
	data.ExpiredTime = time.Now().Add(24 * time.Hour)
	data.VerifyStatus = Enum.EmailVerifyWait
	data.VerifyType = verifyType
	err := Email.InsertEmailVerifyData(engine, data)
	if err != nil {
		log.Error("get user Info Error", err)
		return link, fmt.Errorf("系統錯誤！")
	}
	return link, nil
}

func SendVerifyMail(username, email, link string) error {
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 電子郵件驗證信",
		Username: username,
		Title:    "請驗證你的電子郵件",
		Content:  "<p class='text'>你已申請電子郵件驗證，請依照下方進行：</p>",
		Subcontent: fmt.Sprintf("<p align='center'><a href='%s' target='_blank' style='color: #00b896'>%s</a></p><p align='center' style='font-size:13px; font-weight: bold'>%s</p>",
			link, "請點此連結驗證電子郵件，此連結有效時間為48小時，請盡早完成。", "若你沒有註冊此帳號，請勿點此連結，謝謝。"),
	}
	err := mail.ToSendMail(email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendStoreVerifyMail(UserData entity.MemberData, StoreData entity.StoreDataResp, email, link string) error {
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 電子郵件驗證信",
		Username: UserData.Username,
		Title:    "請驗證你的電子郵件",
		Content:  fmt.Sprintf("<p class='text'>%s 指派你為收銀機管理員：共同經營收銀機、上架商品、出貨及訂單客服等收銀機服務功能。<br><br><br>請點選以下電子郵件驗證連結完成驗證後，正式啟用收銀機管理功能。<br>。</p>", StoreData.StoreName),
		Subcontent: fmt.Sprintf(
			"<p align='center'><a href='%s' target='_blank' style='color: #00b896'>%s</a></p><p align='center' style='font-size:13px; font-weight: bold'>%s</p>",
			link, "請點此連結認證你的 Mail 帳號，<br>此連結在 24 小時內有效，請儘早完成認證。", "若你沒有申請Check'Ne.com收銀機服務，請勿點選此連結，謝謝。"),
	}
	err := mail.ToSendMail(email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendUpgradeApplyMail(UserData entity.MemberData, OrderData entity.B2cOrderData) error {
	text := "<p class='text'>感謝你在 Check'Ne 網站上完成帳號升級服務。</p><p class='text'><b> %s </b>，方案費用為 NT：%v 元，將於%s生效。</p><p class='text'>下一筆帳單將於每月 %s 日前發出付款通知。</p>"
	now := time.Now().Format("2006/01/02")
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 帳號升級-付款完成通知",
		Username: UserData.Username,
		Title:    "帳號升級-付款完成通知",
		Content:  fmt.Sprintf(text, OrderData.ProductName, OrderData.Amount, now),
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendUpgradeStopMail(UserData entity.MemberData) error {
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 收銀機服務方案終止通知",
		Username: UserData.Username,
		Title:    "收銀機服務方案終止通知",
		Content:  fmt.Sprintf("<p class='text'>已收到你申請提前終止「%s」。<br> %s 將於 %s 到期，在此之前你可以繼續使用本收銀機服務。<br><br>如有需要購買其他方案，請至「收銀機管理」>「我的收銀機」中選擇「升級收銀機」中挑選合適的收銀機方案。</p>",
			Upgrade.GetChangeProduct(UserData.UpgradeLevel), Upgrade.GetChangeProduct(UserData.UpgradeLevel), UserData.UpgradeExpire.Format("2006/01/02")),
	}
	if err := mail.ToSendMail(UserData.Email); err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendAddManagerMail(UserData entity.MemberData, StoreData entity.StoreData) error {
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 收銀機新增管理權限通知",
		Username: UserData.Username,
		Title:    "收銀機新增管理權限通知",
		Content:  fmt.Sprintf("<p class='text'>已收到你申請[%s]：新增一名管理員權限，如你未曾申請新增管理員權限，請儘速與我們聯繫。</p>", StoreData.StoreName),
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendDeleteManagerMail(UserData entity.MemberData, StoreData entity.StoreData) error {
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 收銀機刪除管理權限通知",
		Username: UserData.Username,
		Title:    "收銀機刪除管理權限通知",
		Content:  fmt.Sprintf("<p class='text'>已收到你申請[%s]：刪除[%s]的收銀機管理權限，如你未曾申請刪除管理員權限，請儘速與我們聯繫。</p>", StoreData.StoreName, UserData.Username),
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendUpgradePaymentMail(UserData entity.MemberData, StoreData entity.StoreDataResp, OrderData entity.B2cBillingData) error {
	OrderData.Expiration.AddDate(0, 1, 0).Format("2006/01/02")
	url := "https://www.checkne.com/store/setting/upgrade-cart/"
	text := "<p style='margin: 0; line-height: 1.5; color: #545454;'><b>%s收銀機</b>新增一筆收銀機帳單<b>%s</b>。<br>方案有效時間：<b>%s-%s。</b><br>每月方案費用：<b>NT＄%v</b>。<br>如有任何收銀機管理相關問題，請與我們聯繫。</p>"
	text2 := "<a href='%s' style='display: inline-block; padding: 12px 36px; text-align: center; border-radius: 5px; background: #00b896; text-decoration: none; color: #ffffff;'>立即付款</a>"

	link := "<div style='border-bottom: 24px solid #ffffff;'><table width='100%' style='text-align: center; background-color: #f7f7f7; border-top: 24px solid #f7f7f7; border-bottom: 24px solid #f7f7f7;'><tr><td>" + fmt.Sprintf(text2, url) + "</td></tr></table></div>"
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 收銀機服務繳費通知",
		Username: UserData.Username,
		Title:    "收銀機服務繳費通知",
		Content: fmt.Sprintf(text,
			StoreData.StoreName, OrderData.ProductName, OrderData.Expiration.Format("2006/01/02"), OrderData.Expiration.AddDate(0, 1, 0).Format("2006/01/02"), OrderData.Amount),
		Subcontent: link,
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendApplyWithdrawMail(UserData entity.MemberData, Amount int64) error {
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 申請提領通知",
		Username: UserData.Username,
		Title:    "申請提領通知",
		Content:  fmt.Sprintf("<p class='text'>新增一筆提領款項，提領金額NT$ %v 元。<br>如你未曾申請款項提領，請儘速與我們聯繫。</p>", Amount),
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendSuspendStoreMail(UserData entity.MemberData, StoreData entity.StoreData) error {
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 申請暫停收銀機通知",
		Username: UserData.Username,
		Title:    "申請暫停收銀機通知",
		Content:  fmt.Sprintf("<p class='text'>[%s]已申請暫停對外營業，商品下架、結帳連結無法使用。<br>如要恢復收銀機經營，請至「我的收銀機」>開啟「收銀機狀態」；並更新商品列表中商品狀態：上架商品。<br><br>如你未曾申請「暫停收銀機」功能，請儘速與我們聯繫。</p>", StoreData.StoreName),
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendPlatformSuspendStoreMail(UserData entity.MemberData, StoreData entity.StoreData) error {
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 收銀機暫停通知",
		Username: UserData.Username,
		Title:    "收銀機暫停通知",
		Content:  fmt.Sprintf("<p class='text'>[%s]因違反平台管理規則，目前強制暫停中，無法進行商品上架及結帳成交交易。<br>如需重新開啟收銀機管理功能，請儘速與我們聯繫。</p>", StoreData.StoreName),
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendSuccessfullySellerOrderedMail(UserData entity.MemberData, StoreData entity.StoreData, OrderData entity.OrderData) error {
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 成功訂購通知信",
		Username: UserData.Username,
		Title:    "成功訂購通知信",
		Content: fmt.Sprintf("<p class='text'>[%s]有一筆新訂單：<br>訂單編號：%s <br>訂購時間：%s <br>請登入CheckNe.com安排出貨。<br><br>如有任何收銀機管理相關問題，請與我們聯絡。</p>",
			StoreData.StoreName, OrderData.OrderId, OrderData.CreateTime.Format("2006/01/02 15:04")),
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendWaitShipMail(UserData entity.MemberData, StoreData entity.StoreData, OrderData entity.OrderData) error {
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 訂單待出貨通知信",
		Username: UserData.Username,
		Title:    "訂單待出貨通知信",
		Content: fmt.Sprintf("<p class='text'>目前有一筆訂單狀態為：待出貨，請儘速安排出貨。<br><br>Check'Ne[%s] 訂單編號：%s <br>訂購時間：%s <br><br>如有任何收銀機管理相關問題，請與我們聯繫。</p>",
			StoreData.StoreName, OrderData.OrderId, OrderData.CreateTime.Format("2006/01/02 15:04")),
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendNotShippedMail(UserData entity.MemberData, count int) error {
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 訂單待出貨通知信",
		Username: UserData.Username,
		Title:    "訂單待出貨通知信",
		Content:  fmt.Sprintf("<p class='text'>目前有 %v 筆訂單狀態為：待出貨，請儘速安排出貨，<br>如有任何收銀機管理相關問題，請與我們聯繫。</p>", count),
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendOrderShipOverdueMail(UserData entity.MemberData, StoreData entity.StoreData, OrderData entity.OrderData) error {
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 訂單逾期未交寄通知信",
		Username: UserData.Username,
		Title:    "訂單逾期未交寄通知信",
		Content: fmt.Sprintf("<p class='text'>Check'Ne[%s] 訂單編號：%s<br>訂購時間：%s<br>這筆訂單的出貨單號，已超過交寄期限，無法再使用該出貨單號交寄，若需退款給買方，請儘早辦理。<br><br>如有任何退款相關問題，請與我們聯繫。</p>",
			StoreData.StoreName, OrderData.OrderId, OrderData.CreateTime.Format("2006/01/02 15:04")),
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendSuccessfullyBuyerOrderedMail(UserData entity.MemberData, StoreData entity.StoreData, OrderData entity.OrderData) error {
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 成功訂購通知信",
		Username: UserData.Username,
		Title:    "成功訂購通知信",
		Content: fmt.Sprintf("<p class='text'>你已在[%s]完成訂購：<br>訂單編號：%s<br>訂購時間：%s<br>你可以登入CheckNe.com追蹤出貨進度。<br>Check'Ne不會主動通知你辦理付款、取消付款、或解除分期等相關作業，也不會請你到ATM解除錯誤設定或進行任何操作。<br>如有任何訂單相關問題，請與我們聯絡。</p>",
			StoreData.StoreName, OrderData.OrderId, OrderData.CreateTime.Format("2006/01/02 15:04")),
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendAtmSuccessfullyBuyerOrderedMail(UserData entity.MemberData, StoreData entity.StoreData, OrderData entity.OrderData, data entity.TransferData) error {
	date := data.ExpireDate.Format("2006/01/02") + " 23:59"
	bank := fmt.Sprintf("%s%s", data.BankCode, data.BankName)
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 成功訂購通知信",
		Username: UserData.Username,
		Title:    "成功訂購通知信",
		Content: fmt.Sprintf("<p class='text'>你已在[%s]完成訂購：<br>訂單編號：%s<br>訂購時間：%s<br>請於%s前繳款，金額NT$%d元，<br>轉帳銀行：%s，轉帳帳號%s<br>你可以登入CheckNe.com追蹤出貨進度。<br>Check'Ne不會主動通知你辦理付款、取消付款、或解除分期等相關作業，也不會請你到ATM解除錯誤設定或進行任何操作。<br>如有任何訂單相關問題，請與我們聯絡。</p>",
			StoreData.StoreName, OrderData.OrderId, OrderData.CreateTime.Format("2006/01/02 15:04"),
			date, data.Amount, bank, data.BankAccount),
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendForSellerShipClosedShopMail(UserData entity.MemberData, StoreData entity.StoreData, OrderData entity.OrderData) error {
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 重新選擇取件超商通知信",
		Username: UserData.Username,
		Title:    "重新選擇取件超商通知信",

		Content: fmt.Sprintf("<p class='text'>Check'Ne[%s] 訂單編號：%s<br>訂購時間：%s <br>此筆訂單選擇配送之超商門市，目前已無法取件，請儘速聯繫買家在2日內重新選擇配送門市位置，以免逾期無法完成配送。</p>",
			StoreData.StoreName, OrderData.OrderId, OrderData.CreateTime.Format("2006/01/02 15:04")),
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendForBuyerShipClosedShopMail(UserData entity.MemberData, StoreData entity.StoreData, OrderData entity.OrderData) error {
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 重新選擇取件超商通知信",
		Username: UserData.Username,
		Title:    "重新選擇取件超商通知信",
		Content: fmt.Sprintf("<p class='text'>Check'Ne[%s] 訂單編號：%s <br>訂購時間：%s <br>此筆訂單選擇配送之超商門市，目前已無法取件，請在2日內重新選擇配送門市位置，以免逾期無法完成配送。</p>",
			StoreData.StoreName, OrderData.OrderId, OrderData.CreateTime.Format("2006/01/02 15:04")),
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendSellerAtmBillOrderMail(UserData entity.MemberData, StoreData entity.StoreData, OrderData entity.OrderData) error {
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 成功接受訂購通知信",
		Username: UserData.Username,
		Title:    "Check'Ne 成功接受訂購通知信",
		Content: fmt.Sprintf("<p class='text'>[%s]有一筆新接受的訂單：<br>訂單編號：%s <br>訂購時間：%s <br>此為ATM繳款帳單，買家繳款成功後，請儘速安排出貨流程。</p>",
			StoreData.StoreName,
			OrderData.OrderId,
			OrderData.CreateTime.Format("2006/01/02 15:04"),
		),
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendBuyerAtmBillOrderMail(UserData entity.MemberData, StoreData entity.StoreData, OrderData entity.OrderData, data entity.TransferData) error {
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 成功接受訂購通知信",
		Username: UserData.Username,
		Title:    "Check'Ne 成功接受訂購通知信",
		Content: fmt.Sprintf("<p class='text'>你已在[%s]完成訂購：<br>訂單編號：%s<br>訂購時間：%s<br>請於 %s 前繳款，金額NT$%v元，<br>轉帳銀行：%s，轉帳帳號%s <br>你可以登入CheckNe.com追蹤出貨進度。<br><br>Check'Ne不會主動通知你辦理付款、取消付款、或解除分期等相關作業，也不會請你到ATM解除錯誤設定或進行任何操作。</p>",
			StoreData.StoreName,
			OrderData.OrderId,
			OrderData.CreateTime.Format("2006/01/02 15:04"),
			data.ExpireDate.Format("2006/01/02")+" 23:59",
			OrderData.TotalAmount,
			fmt.Sprintf("%s%s", data.BankCode, data.BankName),
			data.BankAccount,
		),
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendSellerBillOrderSuccessMail(UserData entity.MemberData, StoreData entity.StoreData, OrderData entity.OrderData) error {
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 成功接受訂購通知信",
		Username: UserData.Username,
		Title:    "Check'Ne 成功接受訂購通知信",
		Content: fmt.Sprintf("<p class='text'>[%s]有一筆新接受的訂單：<br>訂單編號：%s<br>訂購時間：%s<br>請登入CheckNe.com安排出貨。</p>",
			StoreData.StoreName,
			OrderData.OrderId,
			OrderData.CreateTime.Format("2006/01/02 15:04"),
		),
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendBuyerBillOrderSuccessMail(UserData entity.MemberData, StoreData entity.StoreData, OrderData entity.OrderData) error {
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 成功接受訂購通知信",
		Username: UserData.Username,
		Title:    "Check'Ne 成功接受訂購通知信",
		Content: fmt.Sprintf("<p class='text'>[%s]已接受你的訂購：<br>訂單編號：%s<br>訂購時間：%s<br>你可以登入CheckNe.com追蹤出貨進度。<br><br>Check'Ne不會主動通知你辦理付款、取消付款、或解除分期等相關作業，也不會請你到ATM解除錯誤設定或進行任何操作。</p>",
			StoreData.StoreName,
			OrderData.OrderId,
			OrderData.CreateTime.Format("2006/01/02 15:04"),
		),
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendBuyerOrderCustomerMail(UserData entity.MemberData, data entity.OrderMessageBoardData) error {
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 訂單客服問題",
		Username: UserData.Username,
		Title:    "Check'Ne 訂單客服問題",
		Content: fmt.Sprintf("<p class='text'>訂單編號：%s，提出一筆客服問題，請儘速安排回覆。<br>訂單留言內容：%s</p>",
			data.OrderId,
			data.Message,
		),
	}
	if err := mail.ToSendMail(UserData.Email); err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}


func SendOtpMail(UserData entity.MemberData, code string) error {
	text := "<p class='text'>你的 Check'Ne 登入網站驗證碼，請於 15 分鐘內輸入。</p>"

	codeText := "<p style='margin: 0 0 1em; font-size: .75em; line-height: 1; color: #999999;'>驗證碼</p><p style='margin: 0; font-size: 1.5em; line-height: 1; color: #000000;'><b>%s</b></p>"
	mail := SendMail.SetMassage{
		Subject:    "Check'Ne 電子郵件驗證碼通知信",
		Username:   UserData.Username,
		Title:      "電子郵件驗證碼通知信",
		Content:    text,
		Subcontent: "<div style='border-bottom: 24px solid #ffffff;'><table width='100%' style='text-align: center; background-color: #f7f7f7; border-top: 24px solid #f7f7f7; border-bottom: 24px solid #f7f7f7;'><tr><td>" + fmt.Sprintf(codeText, code) + "</td></tr></table></div>",
	}
	err := mail.ToMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendHostCustmerRequestEmail(email string, phone string, name string, data entity.CustomerData) error {
	mail := SendMail.SetMassage{
		Subject:    data.Question,
		Username:   "系統轉寄",
		Title:      "客服詢問",
		Subcontent: data.Contents,
	}
	if data.OrderId != "" {
		mail.Content = fmt.Sprintf("訂單編號 : %s <br>會員 : %s <br>電話 : %s <br>email : %s", data.OrderId, name, phone, email)
	} else {
		mail.Content = fmt.Sprintf("會員 : %s <br>電話 : %s <br>email : %s", name, phone, email)
	}
	err := mail.ToSendMail("service@checkne.com")
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendShippedMessageEmail(UserData entity.MemberData, OrderData entity.OrderData) error {
	link := fmt.Sprintf("https://www.checkne.com/order/detail/%s", OrderData.OrderId)
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 訂單已出貨通知信",
		Username: UserData.Username,
		Title:    "訂單已出貨通知信",
		Body:     "",
		Content:  fmt.Sprintf("訂單編號：%s 已出貨，查看<a href='%s' target='_blank' style='color: #00b896'>配送進度</a>。", OrderData.ShipNumber, link),
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendOrderInvoiceEmail(UserData entity.MemberData, data entity.InvoiceData) error {
	number := fmt.Sprintf("%s%s", data.InvoiceTrack, data.InvoiceNumber)
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 電子發票開立通知",
		Username: UserData.Username,
		Title:    "電子發票開立通知",
		Body:     "",
		Content:  fmt.Sprintf("感謝你使用Check'Ne服務，平台服務費電子發票已經開立，電子發票資訊如下：<br>發票號碼：%s <br>隨機碼：%s <br>發票金額：NT$ %v元<br>開立日期：%s <br><br>" +
			"你可以登入Check'Ne，在會員中心「發票資訊」查詢發票明細資料，也可以在收到電子發票開立通知48小時後於財政部電子發票整合服務平台查詢發票明細資料。<br>" +
			"財政部電子發票整合服務平台之連結：<br>https://www.einvoice.nat.gov.tw/APCONSUMER/BTC601W/", number, data.RandomNumber, data.Amount, data.CreateTime.Format("2006/01/02 15:04:05")),
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}


func SendServiceInvoiceEmail(UserData entity.MemberData, data entity.InvoiceData) error {
	number := fmt.Sprintf("%s%s", data.InvoiceTrack, data.InvoiceNumber)
	mail := SendMail.SetMassage{
		Subject:  "Check'Ne 加值服務費電子發票開立通知",
		Username: UserData.Username,
		Title:    "加值服務費電子發票開立通知",
		Body:     "",
		Content:  fmt.Sprintf("感謝你使用Check'Ne服務，加值服務費電子發票已經開立，電子發票資訊如下：<br>發票號碼：%s <br>隨機碼：%s <br>發票金額：NT$ %v元<br>開立日期：%s <br><br>" +
			"你可以登入Check'Ne，在會員中心「發票資訊」查詢發票明細資料，也可以在收到電子發票開立通知48小時後於財政部電子發票整合服務平台查詢發票明細資料。<br>" +
			"財政部電子發票整合服務平台之連結：<br>https://www.einvoice.nat.gov.tw/APCONSUMER/BTC601W/", number, data.RandomNumber, data.Amount, data.CreateTime.Format("2006/01/02 15:04:05")),
	}
	err := mail.ToSendMail(UserData.Email)
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendHostCustmerContactEmail(company string, email string, phone string, name string, content string) error {
	mail := SendMail.SetMassage{
		Subject:    "訪客詢問",
		Username:   "系統轉寄",
		Title:      "客服詢問",
		Subcontent: content,
	}

	mail.Content = fmt.Sprintf("稱呼 : %s <br>公司名稱 : %s<br>電話 : %s <br>email : %s", name, company, phone, email)

	err := mail.ToSendMail("promo@checkne.com")
	if err != nil {
		log.Error("send mail Error", err)
		return err
	}
	return nil
}

func SendInvoiceSystemMail(last int64) error {
	mail := SendMail.SetMassage{
		Subject:  "發票系統通知",
		Username: "",
		Title:    "發票系統通知",
		Content:  "",
		Subcontent: fmt.Sprintf("發票系統通知，目前發票剩餘張數%v，已少餘300張", last),
	}
	log.Debug("mail", mail)
	email := []string {"john@sharelug.com", "john@sharelug.com"}
	for _, v := range email {
		err := mail.ToSendMail(v)
		if err != nil {
			log.Error("send mail Error", err)
			return err
		}
	}
	return nil
}

func SendCreditSystemMail(OrderId string) error {
	mail := SendMail.SetMassage{
		Subject:  "信用卡請款系統通知",
		Username: "",
		Title:    "信用卡請款系統通知",
		Content:  "",
		Subcontent: fmt.Sprintf("信用卡請款系統通知，訂單編號%s目前出現交易時間有誤，代處理。", OrderId),
	}
	log.Debug("mail", mail)
	email := []string {"john@sharelug.com", "eric@sharelug.com"}
	for _, v := range email {
		err := mail.ToSendMail(v)
		if err != nil {
			log.Error("send mail Error", err)
			return err
		}
	}
	return nil
}
