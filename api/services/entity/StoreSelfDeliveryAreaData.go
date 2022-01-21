package entity

import "time"

type StoreSelfDeliveryArea struct {
	StoreId  string    `xorm:"varchar(50) comment('商店ID')"`
	CityCode string    `xorm:"varchar(15) comment('縣市代碼')"`
	AreaZip  string    `xorm:"text comment('區域郵遞區號集')"`
	Created  time.Time `xorm:"timestamp created"`
	Updated  time.Time `xorm:"timestamp updated"`
}
