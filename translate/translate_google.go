package translate

import (
	"encoding/base64"
	"github.com/name5566/leaf/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//http://translate.klince.top/translate?to=%s&text=%s
func GoogleTranslate(text string, source string, target string) (string, error) {
	// 先走一遍缓存机制
	ok, translateText := getCacheText(text, source, target)
	if ok {
		return translateText, nil
	}
	inputText := strings.Replace(text, "+", "-", -1)
	inputText = strings.Replace(inputText, "/", "_", -1)
	inputText = strings.Replace(inputText, "=", ".", -1)
	inputText = base64.URLEncoding.EncodeToString([]byte(inputText))

	v := make(url.Values)
	v.Set("text", inputText)
	v.Set("to", target)

	url := "http://translate.klince.top/translate?" + v.Encode()
	res, err := http.Get(url)
	if err != nil {
		log.Error("%v", err.Error())
		return "", err
	}
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Error("%v", err.Error())
		return "", err
	}
	outputStr := string(data)
	outputStr = strings.Replace(outputStr, "-", "+", -1)
	outputStr = strings.Replace(outputStr, "_", "/", -1)
	outputStr = strings.Replace(outputStr, ".", "=", -1)
	dec, err := base64.StdEncoding.DecodeString(outputStr)
	if err != nil {
		log.Error("decode failed %v error:%v", string(data), err.Error())
		return "", err
	}
	//log.Debug("%v  ->   %v", text, string(dec))
	outputStr = string(dec)
	outputStr = strings.Replace(outputStr, "＃。", "#.", -1)
	outputStr = strings.Replace(outputStr, "[SIZE.", "[SIZE=", -1)
	outputStr = strings.Replace(outputStr, "[ - ]", "[-]", -1)
	return outputStr, nil
}
