package model

import (
	"api/config/middleware"
	"api/services/Enum"
	"api/services/Service/Excel"
	"api/services/Service/MemberService"
	"api/services/Service/Upgrade"
	"api/services/VO/ExcelVo"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Credit"
	"api/services/dao/member"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type LoginParams struct {
	Email  string `form:"email" validate:"required"`
	Passwd string `form:"passwd" validate:"required"`
}

type RegisterParams struct {
	Username        string `form:"username" validate:"required"`
	Email           string `form:"email" validate:"required"`
	Password        string `form:"password" validate:"required"`
	ConfirmPassword string `form:"confirmpassword" validate:"required"`
}

//更新最後登入時間 AND 登入失敗次數歸零
func UpdateMemberDataLastTimeAndErrorZero(engine *database.MysqlSession, userData *entity.MemberData) error {
	userData.Error = 0
	userData.LastTime = time.Now()
	userData.UpdateTime = time.Now()
	log.Debug("user data =>", userData)
	_, err := member.UpdateMember(engine, userData)
	if err != nil {
		log.Error("update last time error", err)
		return fmt.Errorf("系統錯誤")
	}
	return nil
}

//更新登入失敗次數
func UpdateMemberDataErrorFrequency(engine *database.MysqlSession, userData *entity.MemberData) error {
	userData.Error += 1
	userData.ErrorTime = time.Now()
	userData.UpdateTime = time.Now()
	log.Debug("user error data =>", userData)
	_, err := member.UpdateMember(engine, userData)
	if err != nil {
		log.Error("update error frequency err", err)
		return err
	}
	return nil
}

//驗證密碼和確認密碼是否一致
//func ValidatePassword(params *RegisterParams) error {
//	if params.Password != params.ConfirmPassword {
//		return fmt.Errorf("密碼及密碼確認不一致")
//	}
//	if err := ValidatePasswordRule(params); err != nil {
//		return err
//	}
//	return nil
//}

//驗證密碼規則
//func ValidatePasswordRule(params *RegisterParams) error {
//	match, err := regexp.MatchString("[a-zA-Z]{3,}\\.^[0-9a-z]$\\.^[0-9A-Z]$", params.Password)
//	if err != nil {
//		return err
//	}
//	if !match {
//		return fmt.Errorf("密碼不否符合規則")
//	}
//	return nil
//}

//驗證登入失敗 (失敗3次 一小時之內不能再登入)
func validateLoginErrorFrequency(userData entity.MemberData) error {
	now := time.Now()
	hour, _ := time.ParseDuration("1h")
	hours := now.Add(hour)
	if userData.Error%3 == 0 && hours.Before(userData.ErrorTime) {
		return fmt.Errorf("登入錯誤已達3次 %v", userData.Error%3)
	}
	return nil
}

//登入處理
func HandleLogin(phone string) (Response.PaySendOtpResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.PaySendOtpResponse
	//取得此手機號碼會員
	userData, err := getMemberData(engine, phone)
	if err != nil {
		log.Error("get Member Data Error", err)
		return resp, fmt.Errorf("1001001")
	}
	//無此會員
	if len(userData.Uid) == 0 {
		var day time.Time
		if viper.GetString("ENV") == "prod" {
			day, _ = time.ParseInLocation("20060102 150405", "20210930 170000", time.Local)
		} else {
			day, _ = time.ParseInLocation("20060102 150405", "20210923 000000", time.Local)
		}
		now := time.Now()
		if !now.Before(day) {
			return resp, fmt.Errorf("1002010")
		} else {
			//建立會員
			userData, err = MemberService.BuildNewMember(engine, phone)
			if err != nil {
			return resp, fmt.Errorf("1001001")
		}
		}
	}
	//發送OTP
	log.Info("會員登入 送出OTP", userData.Uid, middleware.GetClientIP())
	data, err := HandlePushOtpSms(engine, userData, phone)
	if err != nil {
		log.Error("Send Sms Otp Error", err)
		return resp, err
	}
	resp.Otp = 1
	resp.Email = ""
	if len(data.Email) != 0 {
		resp.Email = tools.MaskerEMail(userData.Email)
	}
	resp.OtpExpire = fmt.Sprintf("%v", data.ExpireTime.Unix())
	return resp, nil
}

//取出會員資料
func getMemberData(engine *database.MysqlSession, phone string) (entity.MemberData, error) {
	userData, err := member.GetMemberDataByPhone(engine, phone)
	if err != nil {
		return userData, err
	}
	//判斷帳號是否存在
	if userData.Uid == "" {
		log.Debug("user data =>", userData)
		return userData, nil
	}
	return userData, nil
}

func HandleGetEmailVerify(params Request.MemberEmailVerifyRequest) (Response.MemberEmailVerifyResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.MemberEmailVerifyResponse
	resp, err := MemberService.HandleGetEmailVerify(engine, params)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

//驗證會員MAIL
func HandleMemberEmailVerify(userData entity.MemberData, params Request.MemberEmailVerifyRequest) (Response.MemberEmailVerifyResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.MemberEmailVerifyResponse
	resp, err := MemberService.HandleMemberEmailVerify(engine, userData, params)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func HandleGetMemberInfo(userData entity.MemberData, storeData entity.StoreDataResp) (Response.LoginInfoResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.LoginInfoResponse
	resp, err := GetMemberInfo(engine, userData, storeData.StoreId)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

//取會員資訊
func GetMemberInfo(engine *database.MysqlSession, UserData entity.MemberData, storeId string) (Response.LoginInfoResponse, error) {
	var resp Response.LoginInfoResponse
	resp.Member = UserData.GetMemberLoginInfo()
	StoreMax, count, err := Upgrade.ComputeStore(engine, UserData)
	if err != nil {
		return resp, fmt.Errorf("1001001")
	}
	resp.Member.StoreMax = StoreMax
	resp.Member.StoreCount = count
	log.Debug("ssss", UserData, storeId)
	data, err := GetStoreDataByStoreId(engine, storeId, UserData)
	if err != nil {
		log.Debug("Update Email Verify Data Error")
		return resp, err
	}
	resp.Store = data
	return resp, nil
}

//變更會員 訂購人姓名
func ChangeMemberRealName(engine *database.MysqlSession, userId, name string) error {
	data, err := member.GetMemberDataByUid(engine, userId)
	if err != nil {
		return err
	}
	if data.RealName != name {
		data.RealName = name
		if _, err := member.UpdateMember(engine, &data); err != nil {
			return err
		}
	}
	return nil
}

func HandleMemberCompany(params Request.MemberCompanyVerifyRequest) (Response.MemberCompanyVerifyResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.MemberCompanyVerifyResponse
	Key := viper.GetString("EncryptKey")
	data, err := member.GetMemberDataByPhone(engine, params.Mphone)
	if err != nil {
		return resp, fmt.Errorf("1001004")
	}
	if len(data.Uid) == 0 {
		return resp, fmt.Errorf("1001004")
	}
	data.Category = Enum.CategoryCompany
	data.Representative = params.Representative
	data.CompanyAddr = params.CompanyAddr
	data.IdentityName = params.CompanyName
	data.CompanyName = params.CompanyName
	data.RepresentativeId = tools.AesEncrypt(params.RepresentativeId, Key)
	data.Identity = tools.AesEncrypt(params.Identity, Key)
	data.VerifyIdentity = 1
	if _, err := member.UpdateMember(engine, &data); err != nil {
		return resp, fmt.Errorf("1001001")
	}
	resp.Category = Enum.CategoryCompany
	resp.Representative = params.Representative
	resp.CompanyAddr = params.CompanyAddr
	resp.IdentityName = params.CompanyName
	resp.CompanyName = params.CompanyName
	resp.VerifyStatus = "OK"
	resp.Mphone = params.Mphone
	return resp, nil
}

func HandleCompanyPendingList() (Response.CompanyVerifyPendingList, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	Key := viper.GetString("EncryptKey")
	var resp Response.CompanyVerifyPendingList
	counts, datas, err := member.GetMemberCompanyVerifyAndCount(engine)
	if err != nil {
		return resp, err
	}
	for _, data := range datas {
		company := Response.CompanyVerifyInfo{
			UserId:           data.Uid,
			UserPhone:        data.Mphone,
			CompanyName:      data.CompanyName,
			Representative:   data.Representative,
			RepresentativeId: tools.AesDecrypt(data.Identity, Key),
		}
		resp.Companies = append(resp.Companies, company)
	}
	resp.Counts = counts
	return resp, nil
}

//處理會員帳戶中的公司資料
func HandleMemberCompanyInfo(memberId string) (Response.MemberSpecialStore, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.MemberSpecialStore
	data, err := member.GetMemberCompanyInfo(engine, memberId)
	if err != nil {
		return resp, err
	}
	Key := viper.GetString("EncryptKey")
	resp.SpecialStoreName = data.KgiSpecialStore.ChStoreName
	resp.SpecialStoreNameEn = data.KgiSpecialStore.EnStoreName
	resp.MerchantId = data.KgiSpecialStore.MerchantId
	resp.MccCode = data.MemberData.MccCode
	resp.JobCode = data.MemberData.JobCode
	resp.CityCode = data.MemberData.CityCode
	resp.Account = data.MemberData.Uid
	resp.Addr = data.MemberData.CompanyAddr
	resp.AddrEn = data.MemberData.CompanyAddressEn
	resp.Capital = data.MemberData.Capital
	resp.CompanyName = data.MemberData.CompanyName
	resp.ZipCode = data.MemberData.ZipCode
	resp.MemberPhone = data.MemberData.Mphone
	resp.Representative = data.MemberData.Representative
	resp.RepresentFirst = data.MemberData.RepresentFirst
	resp.RepresentLast = data.MemberData.RepresentLast
	resp.MemberName = data.MemberData.SendName
	resp.RepresentativeId = tools.AesDecrypt(data.MemberData.Identity, Key)
	resp.IdentityId = tools.AesDecrypt(data.MemberData.RepresentativeId, Key)
	return resp, nil
}

//更新會員帳戶的公司資料
func UpdateMemberCompanyInfo(memberId string, params Request.MemberPersonalRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := member.UpdateMemberInfo(engine, memberId, params); err != nil {
		return err
	}
	specialStore, err := Credit.GetSellerMerchantIdBySellerId(engine, memberId)
	if err != nil {
		return err
	}
	fmt.Println(specialStore)
	if err := member.InsertMemberKgiRecord(engine, memberId, specialStore.MerchantId); err != nil {
		return err
	}
	return nil
}

func GetMemberSpecialStoreVerifyList() ([]Response.MemberSpecialStore, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	Key := viper.GetString("EncryptKey")
	var resp []Response.MemberSpecialStore
	records, err := member.GetMemberKgiRecord(engine)
	if err != nil {
		return resp, err
	}

	for _, record := range records {
		data := Response.MemberSpecialStore{
			MemberPhone:      record.Member.MemberData.Mphone,
			RepresentLast:    record.Member.MemberData.RepresentLast,
			RepresentFirst:   record.Member.MemberData.RepresentFirst,
			MccCode:          record.Member.MemberData.MccCode,
			JobCode:          record.Member.MemberData.JobCode,
			CityCode:         record.Member.MemberData.CityCode,
			Terminal3D:       record.Member.KgiSpecialStore.Terminal3dId,
			TerminalId:       record.Member.KgiSpecialStore.Terminaln3dId,
			MerchantId:       record.Member.KgiSpecialStore.MerchantId,
			Representative:   record.Member.MemberData.Representative,
			RepresentativeId: tools.AesDecrypt(record.Member.MemberData.RepresentativeId, Key),
			Addr:             record.Member.MemberData.CompanyAddr,
			AddrEn:           record.Member.MemberData.CompanyAddressEn,
			Account:          record.Member.MemberData.Uid,
			Created:          record.MemberSendKgiBank.Created.Format("2006/01/02 15:04"),
		}
		resp = append(resp, data)
	}
	return resp, err
}
func GetMemberSpecialStoreRecords(params Request.MemberPersonalSpecialStoreRequest) ([]Response.MemberSpecialStore, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp []Response.MemberSpecialStore
	Key := viper.GetString("EncryptKey")
	var memberUid string
	if len(params.MemberPhone) != 0 {
		userData, err := member.GetMemberDataByPhone(engine, params.MemberPhone)
		if err != nil {
			return resp, err
		}
		memberUid = userData.Uid
	}
	records, err := member.GetMemberKgiRecordWithMemberMerchant(engine, memberUid, params.MerchantId)
	if err != nil {
		return resp, err
	}

	for _, record := range records {
		data := Response.MemberSpecialStore{
			RepresentLast:    record.Member.MemberData.RepresentLast,
			RepresentFirst:   record.Member.MemberData.RepresentFirst,
			MccCode:          record.Member.MemberData.MccCode,
			JobCode:          record.Member.MemberData.JobCode,
			CityCode:         record.Member.MemberData.CityCode,
			Terminal3D:       record.Member.KgiSpecialStore.Terminal3dId,
			TerminalId:       record.Member.KgiSpecialStore.Terminaln3dId,
			MerchantId:       record.Member.KgiSpecialStore.MerchantId,
			Representative:   record.Member.MemberData.Representative,
			RepresentativeId: tools.AesDecrypt(record.Member.MemberData.RepresentativeId, Key),
			Addr:             record.Member.MemberData.CompanyAddr,
			AddrEn:           record.Member.MemberData.CompanyAddressEn,
			Account:          record.Member.MemberData.Uid,
			Created:          record.MemberSendKgiBank.Created.Format("2006/01/02 15:04"),
		}
		resp = append(resp, data)
	}
	return resp, err
}
func GetMemberSpecialStoreRecordExcelWithIds() (string, string, []int64, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var datas []ExcelVo.SpecialStoreRecordVo
	var recordIds []int64
	records, err := member.GetMemberKgiRecordExcel(engine)
	if err != nil {
		return "", "", recordIds, fmt.Errorf("1001001")
	}
	Key := viper.GetString("EncryptKey")
	for i, record := range records {
		data := ExcelVo.SpecialStoreRecordVo{
			Id:             int64(i + 1),
			RepresentLast:  record.Member.MemberData.RepresentLast,
			RepresentFirst: record.Member.MemberData.RepresentFirst,
			MccCode:        record.Member.MemberData.MccCode,
			JobCode:        record.Member.MemberData.JobCode,
			Terminal3DId:   record.Member.KgiSpecialStore.Terminal3dId,
			TerminalId:     record.Member.KgiSpecialStore.Terminaln3dId,
			MerchantId:     record.Member.KgiSpecialStore.MerchantId,
			CityCode:       record.Member.MemberData.CityCode,
			CityName:       record.Member.MemberData.CityNameEn,
			StoreName:      record.Member.KgiSpecialStore.ChStoreName,
			StoreNameEn:    record.Member.KgiSpecialStore.EnStoreName,
			Represent:      record.Member.MemberData.Representative,
			RepresentId:    tools.AesDecrypt(record.Member.MemberData.RepresentativeId, Key),
			Addr:           record.Member.MemberData.CompanyAddr,
			AddrEn:         record.Member.MemberData.CompanyAddressEn,
		}
		datas = append(datas, data)
		recordIds = append(recordIds, record.MemberSendKgiBank.Id)
	}
	log.Info(fmt.Sprintf("Today total %v kgi special store records upload", len(datas)))
	filename, path, err := Excel.SpecialStoreRecordNew().ToSpecialStoreRecordFile(datas)
	if err != nil {
		return "", "", recordIds, fmt.Errorf("1001001")
	}
	return filename, path, recordIds, nil
}

func UpdateMemberSpecialStoreRecordIsSendWithFilename(ids []int64, filename string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	err := member.UpdateMemberKgiRecordWithIds(engine, ids, filename)
	if err != nil {
		log.Error("UpdateMemberKgiRecord Error", err)
		return err
	}
	return nil

}
