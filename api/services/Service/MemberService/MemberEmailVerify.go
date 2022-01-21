package MemberService

import (
	"api/services/Enum"
	"api/services/Service/Notification"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Email"
	"api/services/dao/Store"
	"api/services/dao/member"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"time"
)
//取出EMAIL驗證資料
func HandleGetEmailVerify(engine *database.MysqlSession, params Request.MemberEmailVerifyRequest) (Response.MemberEmailVerifyResponse, error) {
	var resp Response.MemberEmailVerifyResponse
	data, err := Email.GetEmailVerifyDataByCode(engine, tools.Trim(params.VerifyCode))
	if err != nil {
		return resp, fmt.Errorf("1001001")
	}
	if len(data.UserId) == 0 {
		return resp, fmt.Errorf("1001001")
	}
	if data.VerifyType == Enum.EmailVerifyTypeStore {
		store, err := Store.GetStoreDataByStoreId(engine, data.StoreId)
		if err != nil {
			log.Debug("Get Store Data Error")
			return resp, err
		}
		resp.VerifyStoreName = store.StoreName
	}
	resp.VerifyType = data.VerifyType
	//驗證成功
	if data.VerifyStatus == Enum.EmailVerifySuccess {
		resp.VerifyStatus = Enum.EmailVerifyAlready
		return resp, nil
	}
	//驗證類型
	if data.VerifyStatus == Enum.EmailVerifyWait {
		//WAIT時先檢查到期時間
		if !time.Now().Before(data.ExpiredTime) {
			resp.VerifyStatus = Enum.EmailVerifyExpired
			return resp, nil
		}
	}
	resp.VerifyStatus = Enum.EmailVerifyWait
	return resp, nil
}
//驗證會員MAIL
func HandleMemberEmailVerify(engine *database.MysqlSession, userData entity.MemberData, params Request.MemberEmailVerifyRequest) (Response.MemberEmailVerifyResponse, error) {
	var resp Response.MemberEmailVerifyResponse
	verify, err := Email.GetEmailVerifyDataByCode(engine, params.VerifyCode)
	if err != nil {
		return resp, fmt.Errorf("1001001")
	}
	user, err := member.GetMemberDataByUid(engine, verify.UserId)
	if err != nil {
		return resp, fmt.Errorf("1001001")
	}
	if len(user.Uid) == 0 || userData.Uid != verify.UserId {
		log.Error("Email Verify uid != uid", userData.Uid, verify.UserId)
		return resp, fmt.Errorf("1003004")
	}
	//驗證 類別 各自處理
	if verify.VerifyType == Enum.EmailVerifyTypeUser {
		//驗證使用者EMAIL
		if err := checkMemberEmailData(engine, userData, verify, &resp); err != nil  {
			return resp, err
		}
	} else if verify.VerifyType == Enum.EmailVerifyTypeStore {
		//驗證管理者EMAIL
		if err := checkStoreManagerEmailData(engine, verify, &resp); err != nil {
			return resp, err
		}
	}
	//驗證資料 目前狀態
	resp.VerifyType = verify.VerifyType
	if verify.VerifyStatus != Enum.EmailVerifyWait {
		//非WAIT即是已驗證
		resp.VerifyStatus = Enum.EmailVerifyAlready
		return resp, nil
	}
	//WAIT時先檢是否驗證過期
	if !time.Now().Before(verify.ExpiredTime) {
		resp.VerifyStatus = Enum.EmailVerifyExpired
		return resp, nil
	}
	//驗證成功
	resp.VerifyStatus = Enum.EmailVerifySuccess
	if err := Email.UpdateEmailVerifyDataByCode(engine, Enum.EmailVerifySuccess, verify.Id); err != nil {
		log.Debug("Update Email Verify Data Error")
		return resp, fmt.Errorf("1001001")
	}
	return resp, nil
}
//驗證會員EMAIL資料
func checkMemberEmailData(engine *database.MysqlSession, userData entity.MemberData, verify entity.EmailVerifyData, resp *Response.MemberEmailVerifyResponse) error {
	//驗證發送EMAIL驗證是否與登入的會員相同
	if userData.VerifyEmail != verify.Email {
		log.Error("Email Verify email != email", userData.VerifyEmail, verify.Email)
		return fmt.Errorf("1003004")
	}
	userData.Email = userData.VerifyEmail
	userData.VerifyEmail = ""
	//更新會員資料
	if _, err := member.UpdateMember(engine, &userData); err != nil {
		return fmt.Errorf("1001001")
	}
	resp.VerifyType = verify.VerifyType
	//發送驗證系統訊息
	if err := Notification.SendEmailVerifySuccess(engine, userData.Uid); err != nil {
		return fmt.Errorf("1001001")
	}
	return nil
}
//驗證管理者EMAIL資料
func checkStoreManagerEmailData(engine *database.MysqlSession, verify entity.EmailVerifyData, resp *Response.MemberEmailVerifyResponse) error {
	//取出賣場資料
	store, err := Store.GetStoreDataByStoreId(engine, verify.StoreId)
	if err != nil {
		log.Debug("Get Store Data Error")
		return fmt.Errorf("1001001")
	}
	//管理者MAIL認證
	manager, err := Store.GetStoreManagerByStoreIdAndUserId(engine, verify.StoreId, verify.UserId)
	if err != nil {
		log.Debug("Get Store Manager Data Error")
		return fmt.Errorf("1001001")
	}
	if len(manager.UserId) == 0 || manager.Email != verify.Email {
		log.Error("Email Verify email != email", manager.Email, verify.Email)
		return fmt.Errorf("1003004")
	}
	resp.VerifyStoreName = store.StoreName
	manager.RankStatus = Enum.StoreRankSuccess
	//更新管理者資訊
	if _, err := Store.UpdateStoreRankData(engine, manager); err != nil {
		log.Debug("Update Store Rank Data Error")
		return fmt.Errorf("1001001")
	}
	//發送驗證系統訊息
	if err := Notification.SendManagerVerifyCompleteMessage(engine, verify.StoreId, verify.UserId); err != nil {
		return fmt.Errorf("1001001")
	}
	return nil
}
