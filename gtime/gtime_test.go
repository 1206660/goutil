package gtime

import (
	"github.com/name5566/leaf/log"
	"testing"
	"time"
)

func TestGetCurPvpType(t *testing.T) {
}

func TestSubDays(t *testing.T) {
	t1 := time.Date(2018, 07, 21, 12, 0, 0, 0, time.Local)
	t2 := time.Date(2018, 07, 20, 13, 0, 1, 0, time.Local)
	log.Debug("Subday = %v", SubDays(t2, t1))
	log.Debug("Subday = %v", SubDays(t1, t2))

	t3 := time.Unix(1536127029, 0)
	t4 := time.Date(2018, 9, 6, 0, 0, 0, 0, time.Local)
	log.Debug("%v  %v", t3.Day(), t4.Day())
	log.Debug("Subday = %v", SubDays(t4, t3))
	log.Debug("Subday = %v", SubDays(t3, t4))

	for i := -20; i < 20; i++ {
		t5 := time.Unix(int64(1536176262+i*3600*2), 0)
		log.Debug("%v IsToday:%v", t5.String(), IsToday(int(t5.Unix())))
	}
}
