package UserAddressService

import (
	"api/services/Enum"
	"api/services/VO/UserAddress"
	"api/services/dao/Mart"
	"api/services/dao/UserAddressData"
	"api/services/dao/iPost"
	"api/services/dao/member"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"strings"
)

func HandleAddress(userData entity.MemberData, params UserAddress.AddressInfo) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	err := AddAddress(engine, userData, params)
	if err != nil {
		return err
	}
	return nil
}

// 新增地址
func AddAddress(engine *database.MysqlSession, userData entity.MemberData, params UserAddress.AddressInfo) error {
	addressEnt := entity.UserAddressData{}
	addressEnt.Uid = userData.Uid
	addressEnt.Address = params.Address
	addressEnt.RealName = params.Name
	addressEnt.Phone = params.Phone
	addressEnt.Type = strings.ToUpper(params.Type)
	addressEnt.Ship = strings.ToUpper(params.Ship)
	addressEnt.Status = "Y"

	// 檢查為宅配
	if strings.Index(addressEnt.Ship, "DELIVERY") != -1 {
		addressEnt.Ship = "DELIVERY"
		if len(strings.Split(addressEnt.Address, ",")) != 4 {
			log.Error("格式錯誤", addressEnt.Address)
			return fmt.Errorf("格式錯誤")
		}
	}
	// 同一個配送方式只能有一個寄送地址
	if addressEnt.Type == "S" {
		booleanT, _ := UserAddressData.CheckSendAddressUniqOfShip(engine, addressEnt)
		if booleanT {
			log.Error("不可重複新增相同配送方式寄送地址", addressEnt)
			return fmt.Errorf("不可重複新增相同配送方式寄送地址")
		}

		if len(addressEnt.RealName) == 0 {
			log.Error("退貨收件人不可為空值", addressEnt)
			return fmt.Errorf("退貨收件人不可為空值")
		}
	}
	// 檢查地址是否已存在
	boolean, _ := UserAddressData.CheckExistAddress(engine, addressEnt)
	if boolean {
		log.Info("不可重複新增相同地址", addressEnt)
		return fmt.Errorf("不可重複新增相同地址")
	}
	// 新增地址資訊
	_, err := UserAddressData.InsertAddressInfo(engine, addressEnt)
	if err != nil {
		return err
	}

	// 賣家新增地址時 將真實姓名寫至member Data SendName
	if addressEnt.Type == "S" {
		memberData, err := member.GetMemberDataByUid(engine, addressEnt.Uid)
		if err != nil {
			log.Error(`GetMemberDataByUid`, err.Error())
			return err
		}
		memberData.SendName = addressEnt.RealName
		_, _ = member.UpdateMember(engine, &memberData)
	}

	return err
}

// 確認寄送狀態是否有寄送地址
func InputCheckShipSendAddressExist(ship string, uid string) bool {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	return CheckShipSendAddressExist(engine, ship, uid)
}

// 確認寄送狀態是否有寄送地址
func CheckShipSendAddressExist(engine *database.MysqlSession, ship string, uid string) bool {
	addressEnt := entity.UserAddressData{}
	addressEnt.Uid = uid
	addressEnt.Type = "S"
	addressEnt.Ship = strings.ToUpper(ship)
	addressEnt.Status = "Y"
	booleanT, _ := UserAddressData.CheckSendAddressUniqOfShip(engine, addressEnt)

	return booleanT
}

// 軟刪除地址
func DeleteAddress(userData entity.MemberData, params UserAddress.DeleteAddress) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	// 取得地址資訊
	info, err := UserAddressData.QueryAddressByUaId(engine, params.UaId, userData.Uid)
	if err != nil {
		return fmt.Errorf("系統錯誤")
	}
	if len(info.UaId) == 0 {
		return fmt.Errorf("查無資訊無法變更")
	}
	info.Status = "N"
	err = UserAddressData.UpdateSenderInfo(engine, info)
	if err != nil {
		return fmt.Errorf("刪除失敗")
	}
	return nil
}

// 取得寄送方式預設地址
func GetSendDefaultAddressByShip(engine *database.MysqlSession, ship string, uid string) entity.UserAddressData {
	address, _ := UserAddressData.QuerySendAddressByShip(engine, uid, ship, "S")
	return address
}

// 取得地址資訊
func GetAddresses(aType string, ship string, member entity.MemberData) ([]UserAddress.AddressInfoResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var addressArray []UserAddress.AddressInfoResponse

	if ship == Enum.DELIVERY_I_POST_BAG1 {
		ship = "DELIVERY"
	}
	addresses, err := UserAddressData.QueryAddresses(engine, member.Uid, aType, ship)
	log.Info("addresses", addresses)

	if err != nil {
		return addressArray, fmt.Errorf("系統錯誤")
	}
	for _, a := range addresses {

		address := GetShipAddress(engine, a)
		if address.Ship == "DELIVERY" {
			addr := strings.Split(address.Address, ",")
			address.Address = tools.MaskerAddressLater(addr[3])
		}
		addressArray = append(addressArray, address)
	}

	return addressArray, nil
}

func setCVSAddress(ent interface{}, City, District, StoreAddress, StoreId, StoreName, MdcEndDate, StoreCloseDate string) UserAddress.AddressInfoResponse {
	now := tools.Now("Ymd")

	address := UserAddress.AddressInfoResponse{}
	address.Country = City
	address.City = District
	address.SpecId = StoreId
	address.Alias = StoreId + " " + StoreName

	switch fmt.Sprintf("%T", ent) {
	case "*entity.SevenMyshipShopData":
		address.Address = fmt.Sprintf(`%s%s%s [%s]`, City, District, StoreAddress, StoreName)
	case "*entity.MartOkStoreData", "*entity.MartHiLifeStoreData", "*entity.MartFamilyStoreData":
		address.Address = fmt.Sprintf(`%s [%s]`, StoreAddress, StoreName)
	}

	address.Status = "N"

	if MdcEndDate != now && StoreCloseDate != now {
		address.Status = "Y"
	}

	return address
}

func HandleCVSAddress(engine *database.MysqlSession, ship, storeId string) UserAddress.AddressInfoResponse {
	address := UserAddress.AddressInfoResponse{}
	switch ship {
	case Enum.CVS_7_ELEVEN:
		seven := entity.SevenMyshipShopData{}
		_ = Mart.QueryCVSData(engine, storeId, &seven)
		address = setCVSAddress(&seven, seven.Country, seven.District, seven.Address, seven.StoreID, seven.StoreName, ``, ``)
	case Enum.CVS_FAMILY:
		family := entity.MartFamilyStoreData{}
		_ = Mart.QueryCVSData(engine, storeId, &family)
		address = setCVSAddress(&family, family.City, family.District, family.StoreAddress, family.StoreId, family.StoreName, family.MdcEndDate, family.StoreCloseDate)
	case Enum.CVS_OK_MART:
		ok := entity.MartOkStoreData{}
		_ = Mart.QueryCVSData(engine, storeId, &ok)
		address = setCVSAddress(&ok, ok.City, ok.District, ok.StoreAddress, ok.StoreId, ok.StoreName, ok.MdcEndDate, ok.StoreCloseDate)
	case Enum.CVS_HI_LIFE:
		hiLife := entity.MartHiLifeStoreData{}
		_ = Mart.QueryCVSData(engine, storeId, &hiLife)
		address = setCVSAddress(&hiLife, hiLife.City, hiLife.District, hiLife.StoreAddress, hiLife.StoreId, hiLife.StoreName, hiLife.MdcEndDate, hiLife.StoreCloseDate)
	}

	return address
}

// 取得買賣家選取配送地址
func GetShipAddress(engine *database.MysqlSession, addressData entity.UserAddressData) UserAddress.AddressInfoResponse {
	address := UserAddress.AddressInfoResponse{}

	switch addressData.Ship {
	case Enum.CVS_FAMILY, Enum.CVS_OK_MART, Enum.CVS_HI_LIFE, Enum.CVS_7_ELEVEN:
		address = HandleCVSAddress(engine, addressData.Ship, addressData.Address)
	case Enum.I_POST:
		ip, _ := iPost.SelectPostBoxByAdmId(engine, addressData.Address)
		address.Country = ip.Country
		address.City = ip.City + "#" + ip.Zip
		address.SpecId = ip.AdmId
		address.Alias = ip.AdmId + " " + ip.AdmAlias
		address.Address = ip.Country + ip.City + ip.Address + ip.AdmLocation
		address.Status = ip.BoxStatus
		break
	case "DELIVERY":
		str := strings.Split(addressData.Address, ",")
		address.Country = str[1]
		address.City = str[2] + "#" + str[0]
		address.Address = addressData.Address
		address.Status = "Y"
		break
	}
	address.Name = addressData.RealName
	address.Phone = addressData.Phone
	address.Ship = addressData.Ship
	address.Id = addressData.UaId
	return address
}

func GetReceiverAddress(engine *database.MysqlSession, id string) string {
	data, err := UserAddressData.GetUserAddresses(engine, id)
	if err != nil {
		log.Error("Get Address Error", err)
	}
	return data.Address
}

func NewShipReceiverAddress(engine *database.MysqlSession, userData entity.MemberData, Shipping, BuyerName, BuyerPhone, ReceiverName, ReceiverPhone, Address string) error {
	name := BuyerName
	phone := BuyerPhone
	if len(ReceiverName) != 0 {
		name = ReceiverName
		phone = ReceiverPhone
	}
	addressInfo := UserAddress.AddressInfo{Ship: Shipping, Address: Address, Type: "R", Name: name, Phone: phone}
	if err := AddAddress(engine, userData, addressInfo); err != nil {
		log.Error("Add Address Error", err)
	}
	return nil
}