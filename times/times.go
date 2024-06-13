package times

import "time"

var location *time.Location

var (
	UTC0Zone  = time.FixedZone("UTC0", 0*60*60) // UTC+0
	UTC7Zone  = time.FixedZone("UTC7", 7*60*60) // UTC+7 越南时区
	UTC8Zone  = time.FixedZone("UTC8", 8*60*60)
	UTC9Zone  = time.FixedZone("UTC9", 9*60*60)   // UTC+9 韩国时间
	UTCN5Zone = time.FixedZone("UTCN5", -5*60*60) // UTC-5 EST 东部标准时区
	UTCN4Zone = time.FixedZone("UTCN4", -4*60*60) // UTC-4 EST 东部标准时区
)

func SetLocation(loc *time.Location) {
	location = loc
}

func InitTimeZone(language string) {
	switch language {
	case "cn":
		SetLocation(UTC8Zone)
	case "kr":
		SetLocation(UTC9Zone)
	case "en":
		SetLocation(UTC0Zone)
	case "vn":
		SetLocation(UTC7Zone)
	default:
		SetLocation(UTC8Zone)
	}
}

const OneDay = 24 * time.Hour

func Now() time.Time {
	return time.Now().In(location)
}

func Date(year, month, day, hour, min, sec, nsec int) time.Time {
	return time.Date(year, time.Month(month), day, hour, min, sec, nsec, location)
}

func Time2Loc(src time.Time) time.Time {
	return src.In(location)
}

func Unix(sec int64) time.Time {
	return time.Unix(sec, 0).In(location)
}

func NowUnix() int64 {
	return Now().Unix()
}

func Date2int(t time.Time) int {
	year, month, day := t.Date()
	return year*10000 + int(month*100) + day
}

func Int2Date(date int) time.Time {
	year, month, day := date/10000, (date/100)%100, date%100
	return Date(year, month, day, 0, 0, 0, 0)
}

func DateTimeCombine(date int, hour, min, sec int) time.Time {
	year, month, day := date/10000, (date/100)%100, date%100
	return Date(year, month, day, hour, min, sec, 0)
}

func init() {
	location = UTC8Zone
}
