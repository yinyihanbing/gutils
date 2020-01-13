package goutils

import (
	"time"
	"strings"
	"errors"
)

// 获取昨天0点的时间
func GetYesterdayZeroTime() time.Time {
	nt := time.Now()
	year, month, day := nt.AddDate(0, 0, -1).Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

// 获取今天0点的时间
func GetTodayZeroTime() time.Time {
	nt := time.Now()
	year, month, day := nt.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

// 获取明天0点的时间
func GetTomorrowZeroTime() time.Time {
	nt := time.Now()
	year, month, day := nt.AddDate(0, 0, 1).Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

// 获取到明天0点的秒数
func GetTomorrowCountDown() int64 {
	nt := time.Now()
	year, month, day := nt.AddDate(0, 0, 1).Date()
	tomorrow := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return tomorrow.Unix() - nt.Unix()
}

// 获取今天0点的秒数
func GetNowDayZeroTs() int64 {
	nt := time.Now()
	year, month, day := nt.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local).Unix()
}

// 检测函数执行时间, limitMillisecond=限制时间, f=执行的函数, limitFunc=超过限制时间执行的函数
func CheckFuncUseTime(limitMillisecond time.Duration, f func(), limitFunc func(elapsedMillisecond time.Duration)) {
	t1 := time.Now()
	f()
	// 如果函数执行超过50毫秒,将打印错误LOG
	elapsed := time.Since(t1)
	if elapsed >= limitMillisecond {
		limitFunc(elapsed)
	}
}

// 字符串时间转时间类型
func ParseTime(value string) (t time.Time, err error) {
	formatTime := "15:04:05"
	formatDate := "2006-01-02"
	formatDateTime := "2006-01-02 15:04:05"
	formatDateTimeT := "2006-01-02T15:04:05"

	if len(value) >= 25 {
		value = value[:25]
		t, err = time.ParseInLocation(time.RFC3339, value, time.Local)
	} else if len(value) >= 19 {
		if strings.Contains(value, "T") {
			value = value[:19]
			t, err = time.ParseInLocation(formatDateTimeT, value, time.Local)
		} else {
			value = value[:19]
			t, err = time.ParseInLocation(formatDateTime, value, time.Local)
		}
	} else if len(value) >= 10 {
		if len(value) > 10 {
			value = value[:10]
		}
		t, err = time.ParseInLocation(formatDate, value, time.Local)
	} else if len(value) >= 8 {
		if len(value) > 8 {
			value = value[:8]
		}
		t, err = time.ParseInLocation(formatTime, value, time.Local)
	}
	return
}

// 将时间戳转换成时间类型
func ParseTimeByTs(ts int64) time.Time {
	return time.Unix(ts, 0)
}

// 获取时间戳的日期部分时间
func GetTimeDateByTs(ts int64) time.Time {
	year, month, day := time.Unix(ts, 0).Date()
	date := time.Date(year, month, day, 0, 0, 0, 0, time.Local)

	return date
}

// 计算2时间相差天数, 正值表示t2比t1大, 负值表示t1比t2大
func TimeSubDays(t1, t2 time.Time) (int, error) {
	// 验证2个时间的时区是否一致
	if t1.Location().String() != t2.Location().String() {
		return 0, errors.New("time zone inconsistency")
	}
	// 2个时间相差的小时数
	hours := t2.Sub(t1).Hours()

	// 转换天数
	return int(hours) / 24, nil
}
