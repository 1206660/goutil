package translate

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gamecat/cache2go"
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/util"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	cachesMap util.Map
)

func init() {
}

type Response struct {
	Translation []string
	ErrorCode   string
}

/*
101 	缺少必填的参数，出现这个情况还可能是et的值和实际加密方式不对应
102 	不支持的语言类型
103 	翻译文本过长
104 	不支持的API类型
105 	不支持的签名类型
106 	不支持的响应类型
107 	不支持的传输加密类型
108 	appKey无效
109 	batchLog格式不正确
201 	解密失败，可能为DES,BASE64,URLDecode的错误
202 	签名检验失败
301 	辞典查询失败
302 	小语种查询失败
303 	服务端的其它异常
401 	账户已经欠费
*/

/*
中文	zh-CHS
中文繁体	zh-TW
日文	ja
英文	EN
韩文	ko
法文	fr
俄文	ru
葡萄牙文	pt
西班牙文	es
越南文	vi
阿拉伯语 	ar
德语 	de
印度尼西亚语 	id
泰语 	th
土耳其语 	tr
*/
// http://fanyi.youdao.com/openapi.do?keyfrom=<keyfrom>&key=<key>&type=data&doctype=<doctype>&version=1.1&q=要翻译的文本
func YoudaoTranslate(text string, source string, target string) (string, error) {
	switch strings.ToLower(target) {
	case "zh-chs":
	case "ja":
	case "en":
	case "ko":
	case "fr":
	case "ru":
	case "pt":
	case "es":
	case "vi":
	case "ar":
	case "de":
	case "id":
	case "th":
	case "tr":
	default:
		target = "en" // 不在列表内的都改为英文
	}

	// 先走一遍缓存机制
	ok, translateText := getCacheText(text, source, target)
	if ok {
		return translateText, nil
	}

	appKey := "06e0059c27996b39"
	appSecret := "6aQq0ihHzAFjlBqoOdlndIEEUkzzWX0t"
	v := make(url.Values)
	v.Set("appKey", appKey)
	v.Set("to", target)
	if source != "" {
		v.Set("from", source)
	}
	salt := strconv.Itoa(rand.Intn(9999999))
	v.Set("salt", salt)
	v.Set("q", text)
	m := md5.New()
	m.Write([]byte(fmt.Sprintf("%v%v%v%v", appKey, text, salt, appSecret)))
	sign := m.Sum(nil)
	strSign := hex.EncodeToString(sign)
	v.Set("sign", strSign)

	url := "http://openapi.youdao.com/api?" + v.Encode()
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
	var r Response
	if err := json.Unmarshal(data, &r); err != nil {
		log.Error("%v", err.Error())
		return "", err
	}
	if r.Translation == nil || len(r.Translation) == 0 {
		return "", errors.New("not found")
	}
	addCacheText(text, source, target, r.Translation[0]) // 添加到缓存里
	return r.Translation[0], nil
}

// 获取缓存的table
func getCacheTable(source string, target string) *cache2go.CacheTable {
	cacheMapKName := fmt.Sprintf("%v->%v", source, target)
	cache := cachesMap.Get(cacheMapKName)
	if cache != nil {
		return cache.(*cache2go.CacheTable)
	}
	mapKName := fmt.Sprintf("TRANSLATE_%v", cacheMapKName)
	newCache := cache2go.Cache(mapKName)
	cachesMap.Set(mapKName, newCache)
	return newCache
}

// 添加一个永久的缓存
func addCacheText(text string, source string, target string, translateText string) {
	cache := getCacheTable(source, target)
	cache.Add(text, 0, translateText)
}

// 获取缓存字符串
func getCacheText(text string, source string, target string) (bool, string) {
	cache := getCacheTable(source, target)
	translateText, err := cache.Value(text)
	if err != nil {
		return false, ""
	}
	if translateText == nil {
		return false, ""
	}
	return true, translateText.Data().(string)
}
