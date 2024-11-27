package corpwechat

import (
	"encoding/json"
	"fmt"
	"github.com/gamecat/cache2go"
	"github.com/gamecat/wechat-sdk/utils"
	"github.com/name5566/leaf/log"
	"time"
)

var debug_mode bool

func SetDebugMode(mode bool) {
	debug_mode = mode
}

type accessToken struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type cachedAccessToken struct {
	AccessToken string
	ExpiredTime int64
}

/*
	发送文本消息
*/
type enterpriseWeChatTextMsg struct {
	ToUser  string            `json:"touser"`
	ToParty string            `json:"toparty"`
	ToTag   string            `json:"totag"`
	MsgType string            `json:"msgtype"`
	AgentId int               `json:"agentid"`
	Text    map[string]string `json:"text"`
	Safe    int               `json:"safe"`
}

func getToken() (string, error) {
	// 先获取缓存Token
	cache := cache2go.Cache("CORP_WECHAT_CACHE")
	kName := "access_token"
	cacheItem, _ := cache.Value(kName)
	if cacheItem != nil {
		cachedAccessToken := cacheItem.Data().(*cachedAccessToken)
		if cachedAccessToken.ExpiredTime > time.Now().Unix() {
			return cachedAccessToken.AccessToken, nil
		}
	}

	serviceUrl := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=wwf8614bc64a060b04&corpsecret=7-tknqW2dnTNqB6_mMMkjhYJ9YOzvqnEzB-M59ChAIE"
	body, err := utils.NewRequest("POST", serviceUrl, nil)
	if err != nil {
		log.Error("CorpWechat query new token err %v", err.Error())
		return "", err
	}
	access_token := new(accessToken)
	err = json.Unmarshal(body, &access_token)
	if err != nil {
		log.Debug("%v", string(body))
		return "", err
	}

	cachedData := new(cachedAccessToken)
	cachedData.AccessToken = access_token.AccessToken
	cachedData.ExpiredTime = time.Now().Unix() + int64(access_token.ExpiresIn-60)
	cache.Add(kName, 0, cachedData)
	log.Debug("CorpWechat New Access_token is %v", access_token.AccessToken)

	return cachedData.AccessToken, nil
}

func SendMsg(text string) error {
	if debug_mode == true {
		log.Error("corpwechat SendMsg : %v", text)
		return nil
	}
	agentID := 1000002
	access_token, err := getToken()
	if err != nil {
		log.Error("CorpWechat get token failed %v", err.Error())
		return err
	}

	msg := new(enterpriseWeChatTextMsg)
	msg.ToUser = "@all"
	msg.MsgType = "text"
	msg.AgentId = agentID
	msg.Text = make(map[string]string)
	msg.Text["content"] = text

	data, err := json.Marshal(msg)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	serviceUrl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%v", access_token)
	body, err := utils.NewRequest("POST", serviceUrl, data)
	if err != nil {
		log.Error("CorpWechat sendMsg %v err is %v", err.Error(), string(body))
		return err
	}

	log.Debug("CorpWechat sendMsg Success : %v", text)
	return nil
}
