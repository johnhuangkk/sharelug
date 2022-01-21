package member

import (
	"api/services/Enum"
	"api/services/VO/Request"
	"api/services/dao/Store"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// 取出 GetMemberByMemID ...
func GetMemberDataByUid(engine *database.MysqlSession, Uid string) (entity.MemberData, error) {
	var member entity.MemberData
	_, err := engine.Engine.Table(entity.MemberData{}).
		Select("*").Where("uid = ?", Uid).Get(&member)
	if err != nil {
		return member, err
	}
	return member, nil
}

// 用Phone取得會員資料
func GetMemberDataByPhone(engine *database.MysqlSession, Phone string) (entity.MemberData, error) {
	var member entity.MemberData
	_, err := engine.Engine.Table(entity.MemberData{}).
		Select("*").Where("mphone = ?", Phone).Get(&member)
	if err != nil {
		return member, err
	}
	return member, nil
}

func GetMemberDataByTerminalId(engine *database.MysqlSession, terminalId string) (entity.MemberData, error) {
	var member entity.MemberData
	_, err := engine.Engine.Table(entity.MemberData{}).
		Select("*").Where("terminal_id = ?", terminalId).Get(&member)
	if err != nil {
		return member, err
	}
	return member, nil
}

// 判斷PHONE是否存在
func IsExistsMemberDataByPhone(engine *database.MysqlSession, Phone string) bool {
	count, err := engine.Engine.Table(entity.MemberData{}).
		Where("mphone = ?", Phone).Count()
	if err != nil {
		panic(err)
	}
	return count != 0
}

// 新增 Insert Member
func InsertMember(engine *database.MysqlSession, data entity.MemberData) error {
	_, err := engine.Session.Table(entity.MemberData{}).Insert(&data)
	if err != nil {
		log.Error("Database Error", err)
		return err
	}
	return nil
}

// 用Phone取得會員資料
func GetMemberAndStoreByAccount(engine *database.MysqlSession, Phone string) ([]entity.QueryUserStore, error) {
	var data []entity.QueryUserStore
	if err := engine.Engine.Table(entity.MemberData{}).Select("*").
		Join("LEFT", entity.StoreData{}, "member_data.uid = store_data.seller_id").
		Where("member_data.mphone = ?", Phone).Find(&data); err != nil {
		return data, err
	}
	return data, nil
}

func GetMemberAndStoreByTerminalId(engine *database.MysqlSession, TerminalId string) ([]entity.QueryUserStore, error) {
	var data []entity.QueryUserStore
	if err := engine.Engine.Table(entity.MemberData{}).Select("*").
		Join("LEFT", entity.StoreData{}, "member_data.uid = store_data.seller_id").
		Where("member_data.Terminal_id = ?", TerminalId).Find(&data); err != nil {
		return data, err
	}
	return data, nil
}

func GetMemberAndStoreByStoreName(engine *database.MysqlSession, StoreName string) ([]entity.QueryUserStore, error) {
	var data []entity.QueryUserStore
	if err := engine.Engine.Table(entity.MemberData{}).Select("*").
		Join("LEFT", entity.StoreData{}, "member_data.uid = store_data.seller_id").
		Where("store_data.store_name = ?", StoreName).Find(&data); err != nil {
		return data, err
	}
	return data, nil
}

// 更新Member
func UpdateMember(engine *database.MysqlSession, member *entity.MemberData) (int64, error) {
	affected, err := engine.Session.Table(entity.MemberData{}).ID(member.Uid).AllCols().Update(member)
	if err != nil {
		log.Error("Update member Error", err)
		return affected, err
	}
	return affected, nil
}

//更新會員等級
func UpdateMemberLevel(engine *database.MysqlSession, uid string, expire time.Time, level int64) error {
	var where []string
	var bind []interface{}
	where = append(where, "upgrade_level = ?")

	if !expire.IsZero() {
		where = append(where, "upgrade_expire = ?")
	}
	where = append(where, "update_time = ?")

	sql := fmt.Sprintf("UPDATE member_data SET %s WHERE uid = ?", strings.Join(where, ", "))
	bind = append(bind, sql)
	bind = append(bind, level)
	if !expire.IsZero() {
		bind = append(bind, expire)
	}
	bind = append(bind, time.Now())
	bind = append(bind, uid)
	_, err := engine.Session.Exec(bind...)
	if err != nil {
		log.Error("UpdateOrderMessageBoardData Error", err)
		return err
	}
	return nil
}


//更新會員身份證認證
func UpdateMemberVerifyIdentity(engine *database.MysqlSession, uid, id, realName string, is int64) error {
	sql := fmt.Sprintf("UPDATE member_data SET verify_identity = ?, identity = ?, identity_name = ?, update_time = ? WHERE uid = ?")
	_, err := engine.Session.Exec(sql, is, id, realName, time.Now(), uid)
	if err != nil {
		log.Error("UpdateOrderMessageBoardData Error", err)
		return err
	}
	return nil
}

func GetAllMemberData(engine *database.MysqlSession) ([]entity.MemberData, error) {
	var data []entity.MemberData
	if err := engine.Engine.Table(entity.MemberData{}).Select("*").
		Asc("create_time").Find(&data); err != nil {
		return data, err
	}
	return data, nil
}

func GetMemberDataBySubscribe(engine *database.MysqlSession) ([]entity.MemberData, error) {
	var data []entity.MemberData
	if err := engine.Engine.Table(entity.MemberData{}).Select("*").
		Where("unsubscribe = ?", false).Asc("create_time").Find(&data); err != nil {
		return data, err
	}
	return data, nil
}


func GetAllMemberByVerifyIdentity(engine *database.MysqlSession) ([]entity.MemberData, error) {
	var data []entity.MemberData
	if err := engine.Engine.Table(entity.MemberData{}).Select("*").
		Where("verify_identity = ?", Enum.VerifyIdentitySuccess).And("report_bank != ?", true).
		Asc("create_time").Find(&data); err != nil {
		return data, err
	}
	return data, nil
}

func GetMemberUpgradeExpire(engine *database.MysqlSession, expire time.Time) ([]entity.MemberData, error) {
	var member []entity.MemberData
	if err := engine.Engine.Table(entity.MemberData{}).Select("*").
		Where("upgrade_level > ? AND upgrade_expire < ?", 0, expire).Find(&member); err != nil {
		return member, err
	}
	return member, nil
}

func CountSearchMemberData(engine *database.MysqlSession, params Request.SearchMemberRequest) (int64, error) {
	where, bind := ComposeSearchMemberParams(engine, params.Search)
	count, err := engine.Engine.Table(entity.MemberData{}).Select("count(*)").
		Where(strings.Join(where, " AND "), bind...).Count()
	if err != nil {
		log.Error("Count Order Database Error", err)
		return count, err
	}
	return count, nil
}

func TakeSearchMemberData(engine *database.MysqlSession, params Request.SearchMemberRequest) ([]entity.MemberData, error) {
	where, bind := ComposeSearchMemberParams(engine, params.Search)
	var data []entity.MemberData
	if err := engine.Engine.Table(entity.MemberData{}).Select("*").
		Where(strings.Join(where, " AND "), bind...).Desc("uid").Find(&data); err != nil {
		return data, err
	}
	return data, nil
}

func ComposeSearchMemberParams(engine *database.MysqlSession, params Request.MemberRequest) ([]string, []interface{}) {
	var where []string
	var bind []interface{}
	if len(params.TerminalId) != 0 {
		where = append(where, "terminal_id = ?")
		bind = append(bind, params.TerminalId)
	}
	if len(params.Account) != 0 {
		where = append(where, "mphone = ?")
		bind = append(bind, params.Account)
	}
	if len(params.Nickname) != 0 {
		where = append(where, "username = ?")
		bind = append(bind, params.Nickname)
	}
	if len(params.StoreName) != 0 {
		var orWhere []string
		var orBind []interface{}
		data, _ := Store.GetStoresByStoreName(engine, params.StoreName)
		for _, v := range data {
			orWhere = append(where, "uid = ?")
			orBind = append(bind, v.SellerId)
		}
		where = append(where, strings.Join(orWhere, " OR "))
		bind = append(bind, orBind...)
	}
	return where, bind
}

func GetMemberStores(engine *database.MysqlSession, account string) (entity.MemberData, []entity.StoreData, map[string][]entity.StoreRankData, error) {
	var memberAccount entity.MemberData

	var data []entity.StoreData
	_, err := engine.Engine.Table(entity.MemberData{}).Select("*").Where("mphone = ?", account).Desc("uid").Get(&memberAccount)

	storeManagers := make(map[string][]entity.StoreRankData)
	if err != nil {
		return memberAccount, data, storeManagers, err
	}
	if err := engine.Engine.Table(entity.StoreData{}).Select("*").Where("seller_id=?", memberAccount.Uid).Find(&data); err != nil {
		return memberAccount, data, storeManagers, err
	}

	for _, row := range data {

		managers, err := Store.GetStoreManagersByStoreId(engine, row.StoreId)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		storeManagers[row.StoreId] = managers
	}
	return memberAccount, data, storeManagers, nil
}

func GetTaiwanCity(engine *database.MysqlSession, city string) (entity.TaiwanCity, error) {
	var data entity.TaiwanCity
	if _, err := engine.Engine.Table(entity.TaiwanCity{}).Select("*").
		Where("name = ?", city).Get(&data); err != nil {
		log.Error("Get Taiwan City Database Error", err)
		return data, err
	}
	return data, nil
}
func GetMemberCompanyVerifyAndCount(engine *database.MysqlSession) (int64, []entity.MemberData, error) {
	var datas []entity.MemberData
	counts, err := engine.Engine.Table(entity.MemberData{}).
		Where("member_data.verify_business =? AND member_data.verify_identity=? AND member_data.report_bank=?", true, true, false).
		FindAndCount(&datas)
	if err != nil {
		log.Error("Count Member Database Error", err)
		return counts, datas, err
	}
	// 先暫時不需要提領成功
	// for i, data := range datas {
	// 	counts := Withdraw.CountWithdrawSuccessByUserId(engine, data.Uid)
	// 	if counts == 0 {
	// 		datas = removeMember(datas, i)
	// 	}
	// }
	return int64(len(datas)), datas, nil
}
func GetMemberCompanyInfo(engine *database.MysqlSession, memberId string) (entity.MemberWithSpecialStore, error) {
	var data entity.MemberWithSpecialStore
	_, err := engine.Engine.Table(entity.MemberData{}).Where("uid =?", memberId).
		Join("LEFT", entity.KgiSpecialStore{}, "kgi_special_store.user_id = member_data.uid").Get(&data)
	if err != nil {
		log.Error("Get Member Database Error", err)
		return data, err
	}
	return data, nil
}

func UpdateCompanyInfo(engine *database.MysqlSession, memberId string, params Request.MemberCompanyRequest) error {
	Key := viper.GetString("EncryptKey")
	_, err := engine.Session.Table(entity.MemberData{}).Where("uid = ? AND verify_business = ? AND verify_identity = ?", memberId, 1, 1).
		Cols("representative,representativeId,identity,represent_last,represent_first,capital,zip_code,company_addr,company_address_en,company_name,establish,contact,contact_phone").
		Update(entity.MemberData{
			Representative:   params.Representative,
			RepresentativeId: tools.AesEncrypt(params.Identity, Key),
			Identity:         tools.AesEncrypt(params.RepresentativeId, Key),
			RepresentLast:    params.RepresentLast,
			RepresentFirst:   params.RepresentFirst,
			ZipCode:          params.ZipCode,
			CompanyAddr:      params.CompanyAddr,
			CompanyAddressEn: params.CompanyAddrEn,
			CompanyName:      params.CompanyName,
			Capital:          params.Capital,
			Establish:        params.Establish,
			Contact:          params.Contact,
			ContactPhone:     params.ContactPhone,
		})
	if err != nil {
		log.Error("Update Member Database Error", err)
		return err
	}
	return nil
}
func UpdateMemberInfo(engine *database.MysqlSession, memberId string, params Request.MemberPersonalRequest) error {
	_, err := engine.Session.Table(entity.MemberData{}).Where("uid = ? AND verify_business = ? AND verify_identity = ?", memberId, 1, 1).
		Cols("represent_last,represent_first,company_address_en,report_bank").
		Update(entity.MemberData{
			RepresentLast:    params.RepresentLast,
			RepresentFirst:   params.RepresentFirst,
			CompanyAddressEn: params.AddrEn,
			ReportBank:       true,
		})
	if err != nil {
		log.Error("Update Member Database Error", err)
		return err
	}
	return nil
}
func InsertMemberKgiRecord(engine *database.MysqlSession, memberId, merchantId string) error {
	_, err := engine.Session.Table(entity.MemberSendKgiBank{}).Insert(&entity.MemberSendKgiBank{Uid: memberId, MerchantId: merchantId})
	if err != nil {
		log.Error("Update Member Database Error", err)
		return err
	}
	return nil
}
func GetMemberKgiRecord(engine *database.MysqlSession) ([]entity.MemberWithSendToSpecial, error) {
	var datas []entity.MemberWithSendToSpecial
	err := engine.Engine.Table(entity.MemberSendKgiBank{}).Limit(10, 0).Desc("created").
		Join("LEFT", entity.MemberData{}, "member_data.uid = member_send_kgi_bank.uid").
		Join("LEFT", entity.KgiSpecialStore{}, "kgi_special_store.user_id = member_send_kgi_bank.uid").
		Find(&datas)
	if err != nil {
		log.Error("GetMemberKgiRecord Error", err)
		return datas, err
	}
	return datas, nil

}

func GetMemberKgiRecordWithMemberMerchant(engine *database.MysqlSession, uid, merchantId string) ([]entity.MemberWithSendToSpecial, error) {
	var datas []entity.MemberWithSendToSpecial
	sql := engine.Engine.Table(entity.MemberSendKgiBank{}).Limit(10, 0).Desc("created").
		Join("LEFT", entity.MemberData{}, "member_data.uid = member_send_kgi_bank.uid").
		Join("LEFT", entity.KgiSpecialStore{}, "kgi_special_store.user_id = member_send_kgi_bank.uid")
	if len(uid) != 0 {
		sql = sql.Where("member_send_kgi_bank.uid = ?", uid)
	}
	if len(merchantId) != 0 {
		sql = sql.Where("member_send_kgi_bank.merchant_id = ?", merchantId)

	}
	err := sql.Find(&datas)
	if err != nil {
		log.Error("GetMemberKgiRecord Error", err)
		return datas, err
	}
	return datas, nil

}
func GetMemberKgiRecordExcel(engine *database.MysqlSession) ([]entity.MemberWithSendToSpecial, error) {
	tNow := time.Now().AddDate(0, 0, -1)

	start := time.Date(tNow.Year(), tNow.Month(), tNow.Day(), 0, 0, 0, 0, tNow.Location())

	end := time.Date(tNow.Year(), tNow.Month(), tNow.Day(), 23, 59, 59, 999999999, tNow.Location())
	var datas []entity.MemberWithSendToSpecial
	sql := engine.Engine.Table(entity.MemberSendKgiBank{}).Desc("created").
		Join("LEFT", entity.MemberData{}, "member_data.uid = member_send_kgi_bank.uid").
		Join("LEFT", entity.KgiSpecialStore{}, "kgi_special_store.user_id = member_send_kgi_bank.uid").
		Where("member_send_kgi_bank.created >=? AND member_send_kgi_bank.created <=?", start, end)
	err := sql.Find(&datas)
	if err != nil {
		log.Error("GetMemberKgiRecord Error", err)
		return datas, err
	}
	return datas, nil

}
func UpdateMemberKgiRecordWithIds(engine *database.MysqlSession, ids []int64, filename string) error {

	_, err := engine.Session.Table(entity.MemberSendKgiBank{}).In("id", ids).Cols("id,is_send,file_name").Update(&entity.MemberSendKgiBank{IsSend: true, FileName: filename})

	if err != nil {
		log.Error("GetMemberKgiRecord Error", err)
		return err
	}
	return nil

}
func removeMember(slice []entity.MemberData, s int) []entity.MemberData {
	return append(slice[:s], slice[s+1:]...)
}
