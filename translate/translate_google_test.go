package translate

import (
	"github.com/name5566/leaf/log"
	"testing"
)

// 测试翻译100次的效率，因为用了缓存
func Test_GoogleTranslate(t *testing.T) {
	translateText, err := GoogleTranslate("今天天气正好", "zh-ch", "en")
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	log.Debug("%v", translateText)
}
