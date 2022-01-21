package entity

import (
	"api/services/Enum"
	"api/services/VO/Response"
	"time"
)

type Promotion struct {
	Id        int64     `xorm:"pk int(10) autoincr comment('活動代號')"`
	StoreId   string    `xorm:"varchar(50) NOT NULL comment('店舖代號')"`
	Name      string    `xorm:"varchar(20) NOT NULL comment('活動名稱')"`
	Status    string    `xorm:"varchar(20) NOT NULL comment('活動狀態')"`
	Type      string    `xorm:"varchar(15) NOT NULL comment('類型')"`
	Value     float64   `xorm:"decimal(10,2) NOT NULL comment('金額比例')"`
	Amount    int64     `xorm:"int(10) NOT NULL comment('數量')"`
	Remaining int64     `xorm:"int(10) comment('剩餘數量')"`
	StartTime time.Time `xorm:"timestamp comment('開始時間')"`
	StopTime  time.Time `xorm:"timestamp comment('自行中止活動時間')"`
	CloseTime time.Time `xorm:"datetime NOT NULL comment('結束時間')"`
	Creator   string    `xorm:"varchar(50) comment('建立者')"`
	Updater   string    `xorm:"varchar(50) comment('更新者')"`
	Aborter   string    `xorm:"varchar(50) comment('中止者')"`
	Created   time.Time `xorm:"timestamp created "`
	Updated   time.Time `xorm:"timestamp updated "`
}

type PromotionCode struct {
	Id           int64     `xorm:"pk int(10) autoincr"`
	PromotionId  int64     `xorm:"int(10) NOT NULL comment('活動代號')"`
	Code         string    `xorm:"varchar(10) NOT NULL comment('折扣碼')"`
	IsUsed       bool      `xorm:"tinyint(1) default 0 comment('是否已使用')"`
	UserId       string    `xorm:"varchar(50) comment('限定使用者')"`
	IsSellerCopy bool      `xorm:"tinyint(1) default 0 comment('賣家複製紀錄')"`
	OperateId    string    `xorm:"varchar(20) NOT NULL comment('取號序號')"`
	Created      time.Time `xorm:"timestamp created  "`
	Updated      time.Time `xorm:"timestamp updated  "`
}

type CouponUsedRecord struct {
	StoreId        string    `xorm:"varchar(50) NOT NULL comment('賣場編號')"`
	OrderId        string    `xorm:"varchar(50) NOT NULL comment('訂單編號')"`
	PromotionId    int64     `xorm:"int(10) NOT NULL comment('活動編號')"`
	BuyerId        string    `xorm:"varchar(50) NOT NULL comment('訂購者編號')"`
	BuyerPhone     string    `xorm:"varchar(20) NOT NULL comment('訂購者電話')"`
	Code           string    `xorm:"varchar(10) NOT NULL comment('優惠代碼')"`
	Amount         float64   `xorm:"decimal(10,2) NOT NULL comment('原始金額')"`
	DiscountAmount float64   `xorm:"decimal(10,2) NOT NULL comment('折扣後金額')"`
	TransTime      time.Time `xorm:"timestamp NOT NULL comment('訂購日期')"`
	RecordStatus   string    `xorm:"varchar(10) NOT NULL comment('狀態')"`
	Created        time.Time `xorm:"timestamp created "`
}

type PromoCodeOperateRecord struct {
	OperateId   string    `xorm:"varchar(20) NOT NULL comment('取號序號')"`
	PromotionId int64     `xorm:"int(10) NOT NULL comment('活動代號')"`
	Type        string    `xorm:"varchar(10) NOT NULL comment('操作類型')"`
	Amount      int64     `xorm:"int(10) NOT NULL comment('操作數量')"`
	UserId      string    `xorm:"varchar(50) NOT NULL comment('操作者編號')"`
	Created     time.Time `xorm:"timestamp created  "`
}

type PromotionAndPromotionCode struct {
	Promotion     Promotion     `xorm:"extends"`
	PromotionCode PromotionCode `xorm:"extends"`
}

func (p *Promotion) GetPromotionDetailResponse(now time.Time) Response.PromoDetail {
	var resp Response.PromoDetail
	resp.Name = p.Name
	resp.StartTime = p.StartTime.Format("2006-01-02 15:04")
	resp.Quantity = p.Amount
	resp.Amount = int64(p.Value)
	resp.Remain = p.Remaining
	resp.Picked = p.Amount - p.Remaining
	resp.StopTime = ""
	resp.Id = p.Id
	if !p.StopTime.IsZero() {
		resp.StopTime = p.StopTime.Format("2006-01-02 15:04")
	}
	resp.EndTime = p.CloseTime.Format("2006-01-02 15:04")
	resp.Status = setPromoStatus(resp.EndTime, resp.StopTime, now)
	resp.StatusText = setPromoStatusText(resp.Status)
	return resp
}
func setPromoStatus(end string, stop string, now time.Time) string {
	if len(stop) > 0 {
		return Enum.PromoStatusStop
	}
	endTime, _ := time.ParseInLocation("2006-01-02 15:04", end, time.Local)
	if now.Sub(endTime) > 0 {
		return Enum.PromoStatusEnd
	}
	return Enum.PromoStatusActive
}
func setPromoStatusText(status string) string {
	switch status {
	case Enum.PromoStatusStop:
		return Enum.PromoStatusStopText
	case Enum.PromoStatusActive:
		return Enum.PromoStatusActiveText
	case Enum.PromoStatusEnd:
		return Enum.PromoStatusEndText
	}
	return Enum.PromoStatusEndText
}
