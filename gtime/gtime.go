package gtime

import (
	"fmt"
	"github.com/astaxie/beego/toolbox"
	"github.com/name5566/leaf/log"
	"path/filepath"
	"runtime"
	"time"
)

const DAY_SECONDS = 86400

// 比较从t2到t1过了多少天
func SubDays(t1, t2 time.Time) int {
	// 如果t1小于t2 则交换一下
	if t1.Unix() < t2.Unix() {
		tmp := t2
		t2, t1 = t1, tmp
	}

	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)
	return int(t1.Sub(t2).Hours() / 24)
}

// 比较是否为同一天
func IsToday(t1 int) bool {
	if t1 == 0 {
		return false
	}
	return SubDays(time.Now(), time.Unix(int64(t1), 0)) == 0
}

func SubWeeks(t1, t2 time.Time) int {
	sdays := SubDays(t1, t2)
	return sdays / 7
}

func TodayLastSecondTick() int64 {
	// "2006-01-02 15:04:05"是golang诞生的时间，必须要晚过这个时间才行。
	timeStr := time.Now().Format("2006-01-02")
	//使用Parse 默认获取为UTC时区 需要获取本地时区 所以使用ParseInLocation
	tm, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 23:59:59", time.Local)
	return tm.Unix()
}

// 是否统计性能分析
var __TIME_COST bool

func StartTimeCost() {
	__TIME_COST = true
}
func StopTimeCost() {
	__TIME_COST = false
}

func TimeCost(start time.Time) {
	if __TIME_COST == true {
		pc, file, line, ok := runtime.Caller(2)
		_, _ = pc, ok
		terminal := time.Since(start)
		millisecond := terminal.Seconds() * 1e3
		if millisecond >= 10 {
			log.Error("pc %v Func(%v:%v) Cost %v millisecond(%v)", pc, file, line, terminal, millisecond)
		}
		// 增加性能分析
		fileWithSuffix := filepath.Base(file)
		go toolbox.StatisticsMap.AddStatistics("", fmt.Sprintf("%v:%v", fileWithSuffix, line), "", terminal)
	}
}
