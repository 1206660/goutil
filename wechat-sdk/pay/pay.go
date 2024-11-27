package pay

import (
	"fmt"
	"strings"
	"time"

	"encoding/xml"
	"wechat-sdk/utils"
)

type (
	// WePay 微信支付配置类
	WePay struct {
		AppID     string // 微信应用APPId或小程序APPId
		MchID     string // 商户号
		PayKey    string // 支付密钥
		NotifyURL string // 回调地址
		TradeType string // 小程序写"JSAPI",客户端写"APP"
	}

	// AppRet 返回的基本内容
	AppRet struct {
		Timestamp string `json:"timestamp,omitempty"` // 时间戳
		NonceStr  string `json:"noncestr,omitempty"`  // 随机字符串
	}

	// AppPayRet 下单返回内容
	AppPayRet struct {
		AppRet

		AppID     string `json:"appid,omitempty"`     // 应用ID
		PartnerID string `json:"partnerid,omitempty"` // 微信支付分配的商户号
		PrepayID  string `json:"prepayid,omitempty"`  // 预支付交易会话ID
		Package   string `json:"package,omitempty"`   // 扩展字段 暂填写固定值Sign=WXPay
		Sign      string `json:"sign,omitempty"`      // 签名
	}

	// WaxRet 返回的基本内容
	WaxRet struct {
		Timestamp string `json:"timeStamp,omitempty"` // 时间戳
		NonceStr  string `json:"nonceStr,omitempty"`  // 随机字符串
	}

	// WaxPayRet 微信小程序下单返回内容
	WaxPayRet struct {
		WaxRet

		AppID    string `json:"appId,omitempty"`    // 应用ID
		Package  string `json:"package,omitempty"`  // 扩展字段 统一下单接口返回的 prepay_id 参数值，提交格式如：prepay_id=*
		SignType string `json:"signType,omitempty"` // 签名算法，暂支持 MD5
		PaySign  string `json:"paySign,omitempty"`  // 签名
	}

	// 微信充值结果返回信息
	WxPayNotifyData struct {
		XMLName        xml.Name `xml:"xml"`                                                    //xml标签
		AppID          string   `xml:"appid" json:"appid"`                                     //Appid
		MchID          string   `xml:"mch_id" json:"mch_id"`                                   //微信支付分配的商户号，必须
		DeviceInfo     string   `xml:"device_info" json:"device_info"`                         //微信支付填"WEB"，必须
		NonceStr       string   `xml:"nonce_str" json:"nonce_str"`                             //随机字符串，必须
		Sign           string   `xml:"sign" json:"sign"`                                       //签名，必须
		ResultCode     string   `xml:"result_code" json:"result_code"`                         //SUCCESS/FAIL
		ErrCode        string   `xml:"err_code" json:"err_code"`                               //SYSTEMERROR
		ErrCodeDes     string   `xml:"err_code_des" json:"err_code_des"`                       //错误返回的信息描述
		TransactionId  string   `xml:"transaction_id" json:"transaction_id"`                   //微信支付订单号
		OpenId         string   `xml:"openid" json:"openid"`                                   //用户在商户appid下的唯一标识
		Body           string   `xml:"body" json:"body"`                                       //商品简单描述，必须
		Detail         string   `xml:"detail,omitempty" json:"detail,omitempty"`               //商品详细列表，使用json格式
		Attach         string   `xml:"attach" json:"attach"`                                   //附加数据，如"贵阳分店"，非必须
		OutTradeNo     string   `xml:"out_trade_no" json:"out_trade_no"`                       //订单号，必须
		FeeType        string   `xml:"fee_type,omitempty" json:"fee_type,omitempty"`           //默认人民币：CNY，非必须
		TotalFee       int      `xml:"total_fee" json:"total_fee"`                             //订单金额，单位分，必须
		CashFee        int      `xml:"cash_fee" json:"cash_fee"`                               //现金支付金额订单现金支付金额，详见支付金额
		CashType       string   `xml:"cash_fee_type,omitempty" json:"cash_fee_type,omitempty"` //默认人民币：CNY，非必须
		SpBillCreateIP string   `xml:"spbill_create_ip" json:"spbill_create_ip"`               //支付提交客户端IP，如“123.123.123.123”，必须
		TimeStart      string   `xml:"time_start,omitempty" json:"time_start,omitempty"`       //订单生成时间，格式为yyyyMMddHHmmss，如20170324094700，非必须
		TimeExpire     string   `xml:"time_expire,omitempty" json:"time_expire,omitempty"`     //订单结束时间，格式同上，非必须
		GoodsTag       string   `xml:"goods_tag,omitempty" json:"goods_tag,omitempty"`         //商品标记，代金券或立减优惠功能的参数，非必须
		NotifyURL      string   `xml:"notify_url" json:"notify_url"`                           //接收微信支付异步通知回调地址，不能携带参数，必须
		TradeType      string   `xml:"trade_type" json:"trade_type"`                           //交易类型，小程序写"JSAPI"，APP 写 APP
		LimitPay       string   `xml:"limit_pay,omitempty" json:"limit_pay,omitempty"`         //限制某种支付方式，非必须
	}

	CDATA struct {
		Text string `xml:",cdata"`
	}
	WxPayReturnNotifyData struct {
		XMLName     xml.Name `xml:"xml"`                            //xml标签
		ReturenCode CDATA    `xml:"return_code" json:"return_code"` //SUCCESS/FAIL
		ReturnMsg   CDATA    `xml:"return_msg" json:"return_msg"`   //return_msg
	}
)

// AppPay App支付
func (m *WePay) AppPay(totalFee int, productID string, clientIP string, attach string) (results *AppPayRet, outTradeNo string, err error) {

	outTradeNo = utils.GetTradeNO(m.MchID)

	appUnifiedOrder := &AppUnifiedOrder{
		UnifiedOrder: UnifiedOrder{
			AppID:          m.AppID,
			MchID:          m.MchID,
			NotifyURL:      m.NotifyURL,
			TradeType:      m.TradeType,
			SpBillCreateIP: clientIP, // 客户端IP 必填
			OutTradeNo:     outTradeNo,
			TotalFee:       totalFee,
			Body:           productID,
			NonceStr:       utils.RandomNumString(16, 32),
			Attach:         attach,
		},
	}

	t, err := utils.Struct2Map(appUnifiedOrder)

	if err != nil {
		return results, outTradeNo, err
	}

	// 获取签名
	sign, err := utils.GenWeChatPaySign(t, m.PayKey)
	if err != nil {
		return results, outTradeNo, err
	}
	appUnifiedOrder.Sign = strings.ToUpper(sign)

	unifiedOrderResp, err := NewUnifiedOrder(appUnifiedOrder)

	if err != nil {
		return results, outTradeNo, err
	}
	results = &AppPayRet{
		AppRet: AppRet{
			Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
			NonceStr:  unifiedOrderResp.NonceStr,
		},
		AppID:     unifiedOrderResp.AppID,
		PartnerID: unifiedOrderResp.MchID,
		PrepayID:  unifiedOrderResp.PrepayID,
		Package:   "Sign=WXPay",
	}

	r, err := utils.Struct2Map(results)

	if err != nil {
		return results, outTradeNo, err
	}

	sign, err = utils.GenWeChatPaySign(r, m.PayKey)
	if err != nil {
		return results, outTradeNo, err
	}
	results.Sign = strings.ToUpper(sign)

	return
}

// WaxPay 小程序支付
func (m *WePay) WaxPay(totalFee int, openID string, productId string) (results *WaxPayRet, outTradeNo string, err error) {

	outTradeNo = utils.GetTradeNO(m.MchID)

	wxaUnifiedOrder := &WxaUnifiedOrder{
		UnifiedOrder: UnifiedOrder{
			AppID:          m.AppID,
			MchID:          m.MchID,
			NotifyURL:      m.NotifyURL,
			TradeType:      m.TradeType,
			SpBillCreateIP: "123.123.123.123", // Ip
			OutTradeNo:     outTradeNo,
			TotalFee:       totalFee,
			Body:           productId,
			NonceStr:       utils.RandomNumString(16, 32),
		},
		OpenID: openID,
	}

	t, err := utils.Struct2Map(wxaUnifiedOrder)

	if err != nil {
		return results, outTradeNo, err
	}

	// 获取签名
	sign, err := utils.GenWeChatPaySign(t, m.PayKey)
	if err != nil {
		return results, outTradeNo, err
	}
	wxaUnifiedOrder.Sign = strings.ToUpper(sign)

	unifiedOrderResp, err := NewUnifiedOrder(wxaUnifiedOrder)

	if err != nil {
		return results, outTradeNo, err
	}
	results = &WaxPayRet{
		WaxRet: WaxRet{
			Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
			NonceStr:  unifiedOrderResp.NonceStr,
		},
		AppID:    m.AppID,
		Package:  "prepay_id=" + unifiedOrderResp.PrepayID,
		SignType: "MD5",
	}

	r, err := utils.Struct2Map(results)

	if err != nil {
		return results, outTradeNo, err
	}

	sign, err = utils.GenWeChatPaySign(r, m.PayKey)
	if err != nil {
		return results, outTradeNo, err
	}
	results.PaySign = strings.ToUpper(sign)

	return
}

//// 公众号支付
//func (m *WePay) H5Pay(totalFee int, openId string) (results *WaxPayRet, outTradeNo string, err error) {
//
//	return m.WaxPay(totalFee, openId)
//}
//
//// 网页支付
//func (m *WePay) WebPay(totalFee int, openId string) (results *WaxPayRet, outTradeNo string, err error) {
//	return m.WaxPay(totalFee, openId)
//}

func (m *WePay) GetWxPayTestResult() []byte {
	xml_test_data := []byte(`
	<xml>
  		<appid><![CDATA[wx2421b1c4370ec43b]]></appid>
  		<attach><![CDATA[支付测试]]></attach>
  		<bank_type><![CDATA[CFT]]></bank_type>
  		<fee_type><![CDATA[CNY]]></fee_type>
  		<is_subscribe><![CDATA[Y]]></is_subscribe>
  		<mch_id><![CDATA[10000100]]></mch_id>
  		<nonce_str><![CDATA[5d2b6c2a8db53831f7eda20af46e531c]]></nonce_str>
  		<openid><![CDATA[oUpF8uMEb4qRXf22hE3X68TekukE]]></openid>
		<out_trade_no><![CDATA[1409811653]]></out_trade_no>
  		<result_code><![CDATA[SUCCESS]]></result_code>
  		<return_code><![CDATA[SUCCESS]]></return_code>
  		<sign><![CDATA[B552ED6B279343CB493C5DD0D78AB241]]></sign>
  		<sub_mch_id><![CDATA[10000100]]></sub_mch_id>
  		<time_end><![CDATA[20140903131540]]></time_end>
  		<total_fee>1</total_fee><coupon_fee><![CDATA[10]]></coupon_fee>
		<coupon_count><![CDATA[1]]></coupon_count>
		<coupon_type><![CDATA[CASH]]></coupon_type>
		<coupon_id><![CDATA[10000]]></coupon_id>
		<coupon_fee><![CDATA[100]]></coupon_fee>
  		<trade_type><![CDATA[JSAPI]]></trade_type>
  		<transaction_id><![CDATA[1004400740201409030005092168]]></transaction_id>
	</xml> `)
	return xml_test_data
}

func (m *WePay) BuildWxPaySuccessResponse() string {
	retData := WxPayReturnNotifyData{
		ReturenCode: CDATA{"SUCCESS"},
		ReturnMsg:   CDATA{"OK"},
	}
	d, _ := xml.Marshal(retData)
	return string(d)
}

func (m *WePay) BuildWxPayFailResponse(return_msg string) string {
	retData := WxPayReturnNotifyData{
		ReturenCode: CDATA{"FAIL"},
		ReturnMsg:   CDATA{return_msg},
	}
	d, _ := xml.Marshal(retData)
	return string(d)
}
