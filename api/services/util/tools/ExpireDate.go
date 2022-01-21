package tools

import (
	"api/services/util/log"
	"time"
)

//到期日 00:00:00
func GenerateExpireTime(n int64) time.Time {
	now := time.Now()
	day := ParseInLocation(now)
	expire := day.Add(time.Hour * time.Duration(24 * n))
	return expire
}
//到期日 23:59:59
func GenerateTransferExpireTime(n int64) time.Time {
	now := time.Now()
	day := ParseInLocation(now)
	to := day.Add(time.Hour * time.Duration(24 * n))
	expire := to.Add(- time.Second * 1)
	return expire
}

func GenerateBeforeTime(n int64) time.Time {
	now := time.Now()
	day := ParseInLocation(now)
	before := day.Add(- time.Hour * time.Duration(24 * n))
	log.Debug("end", before.Format("2006/01/02 15:04:05"))
	return before
}

func ParseInLocation(now time.Time) time.Time {
	day, _ := time.ParseInLocation("2006-01-02 15:04:05", now.Format("2006-01-02") + " 00:00:00", time.Local)
	return day
}

//下個月的終止日
func NextMonth(now time.Time) time.Time {
	day, _ := time.ParseDuration("24h")
	t := ParseInLocation(now)
	nextMonth := t.AddDate(0, 1, 0)
	return nextMonth.Add(time.Hour * time.Duration(24*1)).Truncate(day).Add(time.Hour * -8)
}

//上個月的起始日
func LastMonth(now time.Time) time.Time {
	day, _ := time.ParseDuration("24h")
	t := ParseInLocation(now)
	lastMonth := t.AddDate(0, -1, 0)
	return lastMonth.Add(time.Hour * time.Duration(24*1)).Truncate(day).Add(time.Hour * -8)
}
