package config

import (
	"github.com/name5566/leaf/log"
	"strings"
	"testing"
)

func TestStartWatchDog(t *testing.T) {
	filename := "/data/gamecat/rpo/server/conf/excel/.StringLang.txt.GSwFPk"
	__loadMap = map[string]*ConfigMap{
		"String.txt":      &StringInfoMgr,
		"HeroExpInfo.txt": &HeroExpInfoMgr,
		"DropInfo.txt":    &DropInfoMgr,
		"GameInfo.txt":    &GameInfoMgr,
		"StringLang.txt":  &StringInfoMgr,
	}

	for name, mgr := range __loadMap {
		_, _ = name, mgr
		contain := strings.Contains(filename, name)
		log.Debug("%v %v = %v", filename, name, contain)
	}
}
