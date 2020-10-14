package jpushgo

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/lnyyj/links/httplink"
)

var (
	// ALL 全部
	ALL = "all"

	// Platforms 推送平台
	Platforms = []string{"ios", "android", "winphone"}

	// Audiences 推送目标
	Audiences = []string{"tag", "tag_and", "tag_not", "alias", "registration_id", "segment", "abtest"}

	
	// JPushURL 推送地址
	JPushURL = "https://api.jpush.cn/v3/push"

	// base64对象
	base64Coder = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
)

// JPush 一个推送对象
type JPush struct {
	Platform     interface{}     `json:"platform"`               // 必填	推送平台设置
	Audience     interface{}     `json:"audience"`               // 必填	推送设备指定
	Notification *Notice         `json:"notification,omitempty"` // 可选	通知内容体。是被推送到客户端的内容。与 message 一起二者必须有其一，可以二者并存
	Message      *Message        `json:"message,omitempty"`      // 可选	消息内容体。是被推送到客户端的内容。与 notification 一起二者必须有其一，可以二者并存
	SmsMessage   *SMSMessage     `json:"sms_message,omitempty"`  // 可选	短信渠道补充送达内容体
	Options      *Options        `json:"options,omitempty"`      // 可选	推送参数
	Cid          string          `json:"cid,omitempty"`          // 可选	用于防止 api 调用端重试造成服务端的重复推送而定义的一个标识符。
	MasterSecret string          `json:"-"`
	AppKey       string          `json:"-"`
	AuthCode     string          `json:"-"`
	HC           *httplink.HTTPC `json:"-"`
}

// NewJPush 创建一个推送对象
func NewJPush(appKey, secret string) *JPush {
	auth := "Basic " + base64Coder.EncodeToString([]byte(appKey+":"+secret))
	return &JPush{MasterSecret: secret, AppKey: appKey, AuthCode: auth, HC: httplink.NewHTTPC()}
}

// Clone 创建一个推送对象
func (jp *JPush) Clone() *JPush {
	jp = &JPush{MasterSecret: jp.MasterSecret, AppKey: jp.AppKey, AuthCode: jp.AuthCode, HC: jp.HC}
	return jp
}

// SetPlatform 设置推送平台
func (jp *JPush) SetPlatform(p ...string) *JPush {
	if len(p) == 0 || p[0] == ALL {
		jp.Platform = ALL
		return jp
	}
	var platform []string
	for k := range p {
		platform = append(platform, p[k])
	}
	jp.Platform = platform
	return jp
}

// SetAudience 设置推送目标
// 第一个参数可指定all,或者tag等元素
func (jp *JPush) SetAudience(a ...string) *JPush {
	if len(a) == 0 || a[0] == ALL {
		jp.Audience = ALL
		return jp
	}
	if jp.Audience != nil && jp.Audience.(string) == ALL {
		return jp
	}
	if jp.Audience == nil {
		jp.Audience = make(map[string][]string)
	}
	jp.Audience.(map[string][]string)[a[0]] = a[1:]
	return jp
}

// SetNotification 设置推送通知内容
func (jp *JPush) SetNotification(n *Notice) *JPush {
	jp.Notification = n
	return jp
}

// SetMessage 设置推送自定义消息
func (jp *JPush) SetMessage(m *Message) *JPush {
	jp.Message = m
	return jp
}

// SetSMSMessage 设置推送短信消息
func (jp *JPush) SetSMSMessage(m *SMSMessage) *JPush {
	jp.SmsMessage = m
	return jp
}

// SetOptions 设置推送选项
func (jp *JPush) SetOptions(o *Options) *JPush {
	jp.Options = o
	return jp
}

// Push 开始推送
func (jp *JPush) Push() error {
	req, err := json.Marshal(jp)
	if err != nil {
		return err
	}
	h := jp.HC.NewHeader("Authorization", jp.AuthCode, "Charset", "UTF-8", "Content-Type", "application/json")
	resp, err := jp.HC.SendUpstream(h, string(req), "POST", JPushURL)
	if err != nil {
		return err
	}
	ret := &JPushRet{}
	if err = json.Unmarshal(resp, ret); err != nil {
		return err
	}
	if ret.Err.Code != 0 {
		return fmt.Errorf("jpush error %d:%s", ret.Err.Code, ret.Err.Msg)
	}
	return nil
}
