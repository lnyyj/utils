package jpushgo

import "testing"

const (
	Key    = ""
	Secret = ""
)

func TestJPush(t *testing.T) {
	// 推送消息
	n := NewNotice("hello hes")
	n.SetNotice(&AndroidNotice{Alert: "Android 安卓 push.....", AlertType: 7, Priority: 1})
	n.SetNotice(&IOSNotice{Alert: "ios 苹果 push.....", Sound: "default", Badge: "+1"})

	jp := NewJPush(Key, Secret)                              // 设置KEY, secret
	jp.SetPlatform().SetAudience(Audiences[3], "1390374904") // 设置推送平台和目标
	jp.SetNotification(n)                                    // 设置推送消息
	jp.SetOptions(&Options{ApnsProduction: false})           // 设置推送可选项
	if err := jp.Push(); err != nil {                        // 推送
		t.Log(err)
	}
}
