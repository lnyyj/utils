# 激光推送服务端封装

## 基本使用
```
	// 设置推送消息
	n := NewNotice("hello jpush")
	n.SetNotice(&AndroidNotice{Alert: "Android test push.....", AlertType: 7, Priority: 2})
	n.SetNotice(&IOSNotice{Alert: "ios test push.....", Sound: "default", Badge: +1})

	jp := NewJPush(Key, Secret)                               // 设置KEY, secret
	jp.SetPlatform().SetAudience(Audiences[3], "1", "2", "3") // 设置推送平台和目标
	jp.SetNotification(n)                                     // 设置推送消息
	jp.SetOptions(&Options{ApnsProduction: false})            // 设置推送可选项
	jp.Push() // 推送
```