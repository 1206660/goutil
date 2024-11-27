package monitor

import (
	"common"
	"common/consulMgr"
	"fmt"
	"github.com/gamecat/cache2go"
	"github.com/gamecat/corpwechat"
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/module"
	"github.com/name5566/leaf/timer"
	"time"
)

const (
	SERVICE_SHUTDOWN_MSG         = "!!! 服务%v宕了，速速查看! (%v)"
	SERVICE_SHUTDOWN_RECOVER_MSG = "!!! 服务%v恢复了，恢复耗时! (%v)"
)

type Monitor struct {
	cache    *cache2go.CacheTable
	skeleton *module.Skeleton
}

func (m *Monitor) OnInit(consulIp string) {
	m.cache = cache2go.Cache("CORP_WECHAT_BROKEN_CACHE")
	consulMgr.Mgr.Run(consulIp)

	skeleton := new(module.Skeleton)
	skeleton.Init()

	cr, err := timer.NewCronExpr("* * * * * *")
	if err != nil {
		log.Debug(err.Error())
		return
	}
	skeleton.CronFunc(cr, func() {
		log.Debug("skeleton cronFunc")
		go m.checkService(common.WebManagerServiceName)
		go m.checkService(common.WebgameServiceName)
	})
}

func (m *Monitor) checkService(serviceName string) {
	ok, serviceAddr := consulMgr.Mgr.GetConsulService(serviceName)
	_ = serviceAddr
	if !ok {
		cacheItem, _ := m.cache.Value(serviceName)
		if cacheItem == nil {
			m.cache.Add(serviceName, 0, time.Now())
			corpwechat.SendMsg(fmt.Sprintf(SERVICE_SHUTDOWN_MSG, serviceName, time.Now().String()))
		} else {
			tick := cacheItem.Data().(int64)
			_ = tick
		}
	} else {
		cacheItem, _ := m.cache.Value(serviceName)
		if cacheItem != nil {
			tm := cacheItem.Data().(time.Time)
			corpwechat.SendMsg(fmt.Sprintf(SERVICE_SHUTDOWN_RECOVER_MSG, serviceName, time.Since(tm)))
			m.cache.Delete(serviceName)
		}
	}
}
