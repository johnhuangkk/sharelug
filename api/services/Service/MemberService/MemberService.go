package MemberService

import (
	"api/services/Enum"
	"api/services/VO/Request"
	"api/services/dao/Store"
	"api/services/dao/member"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

//建立新會員資料
func BuildNewMember(engine *database.MysqlSession, phone string) (entity.MemberData, error) {
	var data entity.MemberData
	if err := engine.Session.Begin(); err != nil {
		log.Error("Set Begin Error", err)
		return data, fmt.Errorf("系統錯誤")
	}
	data.Uid = fmt.Sprintf("U%s%s", time.Now().Format("20060102150405"), tools.RangeNumber(99999, 5))
	data.Mphone = phone
	data.InvoiceCarrier = fmt.Sprintf("IC%s%s", time.Now().Format("060102150405"), tools.RangeNumber(99999, 5))
	data.Picture = fmt.Sprintf("/static/img/default-%s.jpg", tools.RangeNumber(10, 2))
	data.Username = fmt.Sprintf("Check'Ne會員")
	data.MemberStatus = Enum.MemberStatusSuccess
	data.Category = Enum.CategoryMember
	data.TerminalId = fmt.Sprintf("CK%s%s", time.Now().Format("060102"), tools.RangeNumber(9999, 4))
	data.Error = 0
	data.LastTime = time.Now()
	data.RegisterTime = time.Now()
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	data.UpgradeType = Enum.UpgradeTypeRenew
	if err := member.InsertMember(engine, data); err != nil {
		engine.Session.Rollback()
		return data, fmt.Errorf("系統錯誤")
	}
	picture := fmt.Sprintf("/static/img/store-%s.jpg", tools.RangeNumber(5, 2))
	if _, err := createStoreData(engine, data.Uid, "checkne收銀機", picture, ""); err != nil {
		engine.Session.Rollback()
		return data, fmt.Errorf("系統錯誤")
	}
	params := Request.CarrierRequest{
		InvoiceType: Enum.InvoiceTypePersonal,
		CarrierType: Enum.InvoiceCarrierTypeMember,
	}
	if err := ChangePostCarrier(engine, params, data); err != nil {
		engine.Session.Rollback()
		return data, fmt.Errorf("系統錯誤")
	}
	data.Unsubscribe = false
	if err := engine.Session.Commit(); err != nil {
		return data, fmt.Errorf("系統錯誤")
	}
	return data, nil
}

func ResetMemberTerminalId() {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := member.GetAllMemberData(engine)
	if err != nil {
		log.Error("Get All Member Error", err)
	}
	for _, v := range data {
		v.TerminalId = getTerminalId()
		log.Debug("sss", v.TerminalId)
		if _, err := member.UpdateMember(engine, &v); err != nil {
			log.Error("Update Member Error", err)
		}
	}
}

func getTerminalId() string {
	return fmt.Sprintf("CK%s%s", time.Now().Format("060102"), tools.RangeNumber(9999, 4))
}

func createStoreData(engine *database.MysqlSession, uid, storeName, picture, expire string) (entity.StoreData, error) {
	storeData, err := Store.InsertStoreData(engine, uid, storeName, picture, expire)
	if err != nil {
		return storeData, err
	}
	//建立管理帳號
	if _, err := Store.InsertStoreRankData(engine, uid, storeData.StoreId, Enum.StoreRankMaster, Enum.StoreRankSuccess, ""); err != nil {
		return storeData, err
	}
	return storeData, nil
}

func ChangePostCarrier(engine *database.MysqlSession, params Request.CarrierRequest, userData entity.MemberData) error {
	data, err := member.GetMemberCarrierByMemberId(engine, userData.Uid)
	if err != nil {
		log.Error("Get Member Carrier Data Error", err)
		return err
	}
	switch params.InvoiceType {
	case Enum.InvoiceTypePersonal:
		data.InvoiceType = Enum.InvoiceTypePersonal
	case Enum.InvoiceTypeCompany:
		data.InvoiceType = Enum.InvoiceTypeCompany
		data.CompanyName = params.CompanyName
		data.CompanyBan = params.CompanyBan
	case Enum.InvoiceTypeDonate:
		data.InvoiceType = Enum.InvoiceTypeDonate
		//驗證捐贈碼
		if err := VerifyDonateCode(params.DonateBan); err != nil {
			log.Error("Verify Donate Code Error", err)
			return err
		}
		data.DonateBan = params.DonateBan
	}
	switch params.CarrierType {
	case Enum.InvoiceCarrierTypeMember:
		data.CarrierType = Enum.InvoiceCarrierTypeMember
		data.CarrierId = userData.InvoiceCarrier
	case Enum.InvoiceCarrierTypeMobile:
		//手機條碼驗證
		data.CarrierType = Enum.InvoiceCarrierTypeMobile
		if err := VerifyMobileCode(params.CarrierId); err != nil {
			log.Error("Verify Mobile Code Error", err)
			return err
		}
		data.CarrierId = params.CarrierId
	case Enum.InvoiceCarrierTypeCert:
		data.CarrierType = Enum.InvoiceCarrierTypeCert
		data.CarrierId = params.CarrierId
	}
	if len(data.MemberId) == 0 {
		data.MemberId = userData.Uid
		if err := member.InsertMemberCarrierData(engine, data); err != nil {
			log.Error("Insert Member Carrier Data Error", err)
			return err
		}
	} else {
		if err := member.UpdateMemberCarrierData(engine, data); err != nil {
			log.Error("Update Member Carrier Data Error", err)
			return err
		}
	}
	return nil
}

func GetMemberCarrierByMemberId(engine *database.MysqlSession, memberId string) (entity.MemberCarrierData, error) {
	data, err := member.GetMemberCarrierByMemberId(engine, memberId)
	if err != nil {
		log.Error("Get Member Carrier Data Error", err)
		return data, err
	}
	if len(data.MemberId) == 0 {
		user, err := member.GetMemberDataByUid(engine, memberId)
		if err != nil {
			log.Error("Get Member Data Error", err)
			return data, err
		}
		data.MemberId = memberId
		data.CarrierType = Enum.InvoiceCarrierTypeMember
		data.CarrierId = user.InvoiceCarrier
		data.InvoiceType = Enum.InvoiceTypePersonal
		if err := member.InsertMemberCarrierData(engine, data); err != nil {
			log.Error("Insert Member Carrier Data Error", err)
			return data, err
		}
	}
	return data, nil
}

func AlterMemberIdentity(engine *database.MysqlSession, userId, IdentityId, realName string) error {
	Key := viper.GetString("EncryptKey")
	identityId := tools.AesEncrypt(IdentityId, Key)
	if err := member.UpdateMemberVerifyIdentity(engine, userId, identityId, realName, Enum.VerifyIdentitySuccess); err != nil {
		log.Error("Update Member data Error", err)
		return fmt.Errorf("1001001")
	}
	if err := Store.UpdateStoreVerifyIdentity(engine, userId, Enum.VerifyIdentitySuccess); err != nil {
		log.Error("Update Store data Error", err)
		return fmt.Errorf("1001001")
	}
	return nil
}

func TakeMemberIdentity(engine *database.MysqlSession, userId string) string {
	data, err := member.GetMemberDataByUid(engine, userId)
	if err != nil {
		log.Error("Get Member data Error", err)
	}
	Key := viper.GetString("EncryptKey")
	return tools.AesDecrypt(data.Identity, Key)
}
