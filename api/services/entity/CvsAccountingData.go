package entity

import "C"
import (
	"api/services/Enum"
	"time"

	"github.com/spf13/viper"
)

type CvsAccountingData struct {
	CvsType     string    `xorm:"varchar(20) comment('超商名稱')"`
	Type        string    `xorm:"varchar(10) pk comment('寫入類型 核帳檔案代號 (Pick) 寄件資料對帳檔 (Send)  賠償檔 (Compensation) ')"`
	ParentId    string    `xorm:"varchar(5) comment('母特店代號')"`
	DataId      string    `xorm:"varchar(20) pk comment('OK & Family 訂單編號 | HiLife 寄件編號')"`
	Amount      float64   `xorm:"decimal(10,2) default 0.00 comment('代收金額')"`
	ServiceType bool      `xorm:"tinyint(1) NOT NULL comment('是否取貨付款： 1 取貨付款 / 0 取貨不付款')"`
	Status      string    `xorm:"char(1) NOT NULL comment('取件核帳檔 (R99) 1:買家 2:賣家 / 其餘為１')"`
	FileDate    time.Time `xorm:"datetime comment('核帳檔上傳日期')"`
	FileName    string    `xorm:"varchar(50) comment('檔案名稱')"`
	CreateTime  time.Time `xorm:"timestamp not null comment('建立時間')"`
	UpdateTime  time.Time `xorm:"timestamp default CURRENT_TIMESTAMP comment('更新時間')"`
	Log         string    `xorm:"text comment('原始資訊')"`
	Checked     bool      `xorm:"tinyint(1) default 0 NOT NULL comment('是否已對帳： 1 對帳 / 0 未對帳')"`
}

/**
Family 商品核帳檔 (R99) 寄件核帳檔 (R98) 商品賠償檔 (R89)
HiLife 取件核帳檔 (R99) 寄件運費檔 (R98) 遺失賠償檔 (R89)
OK 取件資料對帳檔 (PICK) 寄件資料對帳檔 (SEND)
Seven 寄件資料檔 (OL) 取件資料檔 (CESP) 判賠檔(ACTR)
*/
func (c *CvsAccountingData) SetType(t string, cvsType string) {
	switch cvsType {
	case Enum.CVS_FAMILY:
		c.CvsType = Enum.CVS_FAMILY
		c.ParentId = viper.GetString(`MartFamily901.ParentId`)
		switch t {
		case `R99`:
			c.Type = `P`
		case `R98`:
			c.Type = `S`
		case `R89`:
			c.Type = `C`
		}
	case Enum.CVS_HI_LIFE:
		c.CvsType = Enum.CVS_HI_LIFE
		c.ParentId = viper.GetString(`MartHiLife.ParentId`)
		switch t {
		case `R99`:
			c.Type = `P`
		case `R98`:
			c.Type = `S`
		case `R89`:
			c.Type = `C`
		}
	case Enum.CVS_OK_MART:
		c.CvsType = Enum.CVS_OK_MART
		c.ParentId = `OK`
		switch t {
		case `PICK`:
			c.Type = `P`
		case `SEND`:
			c.Type = `S`
		}
	case Enum.CVS_7_ELEVEN:
		c.CvsType = Enum.CVS_7_ELEVEN
		c.ParentId = `seven`
		switch t {
		case `OL`:
			c.Type = `S`
		case `CESP`:
			c.Type = `P`
		case `ACTR`:
			c.Type = `C`

		}
	}

}
