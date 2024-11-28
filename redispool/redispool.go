package redispool

import "time"
import "github.com/gomodule/redigo/redis"
import (
	"github.com/name5566/leaf/log"
)

var Pool PoolMap
var redisHost string

func init() {
}

// 池Map
type PoolMap struct {
	Map map[string]*redis.Pool
}

func (p *PoolMap) Run(id string, host string, passwd string, db string) {
	if p.Map == nil {
		p.Map = make(map[string]*redis.Pool)
	}
	_, e := p.Map[id]
	if e != false {
		log.Debug("Run dupliacate redis db", id, host)
		return
	}

	pool, err := p.dial(host, passwd, db)
	if err != true {
		return
	}
	p.Map[id] = pool
	redisHost = host
}

func (p *PoolMap) Host() string {
	return redisHost
}

func (p *PoolMap) dial(host, passwd, db string) (*redis.Pool, bool) {
	log.Debug("Dial to Redis", host, passwd, db)
	redisPool := &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     10,
		MaxActive:   500,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host)
			if err != nil {
				log.Error(err.Error())
				return nil, err
			}
			if passwd != "" {
				c.Do("AUTH", passwd)
			}
			if db != "" {
				// 选择db
				c.Do("SELECT", db)
			}
			return c, nil
		},
	}
	return redisPool, true
}

func (p *PoolMap) Get(id string) *redis.Pool {
	pool, err := p.Map[id]
	if err != true {
		return nil
	}
	return pool
}

func (p *PoolMap) GetGlobal() *redis.Pool {
	return p.Get("Global")
}
