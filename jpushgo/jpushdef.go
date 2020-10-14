package jpushgo

// Notice 推送通知
type Notice struct {
	Alert    string         `json:"alert,omitempty"`
	Android  *AndroidNotice `json:"android,omitempty"`
	IOS      *IOSNotice     `json:"ios,omitempty"`
	WINPhone *WPNotice      `json:"winphone,omitempty"`
}

// NewNotice 创建一个推送通知
func NewNotice(alert string) *Notice {
	return &Notice{Alert: alert}
}

// SetNotice 设置推送
func (n *Notice) SetNotice(nn interface{}) *Notice {
	switch nn.(type) {
	case *AndroidNotice:
		n.Android = nn.(*AndroidNotice)
	case *IOSNotice:
		n.IOS = nn.(*IOSNotice)
	case *WPNotice:
		n.WINPhone = nn.(*WPNotice)
	default:
		panic("notice type not support of jpush")
	}
	return n
}

// AndroidNotice 安卓通知
type AndroidNotice struct {
	Alert      string                 `json:"alert"`                  // 必填	通知内容	这里指定了，则会覆盖上级统一指定的 alert 信息；内容可以为空字符串，则表示不展示到通知栏。
	Title      string                 `json:"title,omitempty"`        // 可选	通知标题	如果指定了，则通知里原来展示 App名称的地方，将展示成这个字段。
	BuilderID  int                    `json:"builder_id,omitempty"`   // 可选	通知栏样式ID	Android SDK 可设置通知栏样式，这里根据样式 ID 来指定该使用哪套样式。
	Priority   int                    `json:"priority,omitempty"`     // 可选	通知栏展示优先级	默认为0，范围为 -2～2 ，其他值将会被忽略而采用默认。
	Category   string                 `json:"category,omitempty"`     // 可选	通知栏条目过滤或排序	完全依赖 rom 厂商对 category 的处理策略
	Style      int                    `json:"style,omitempty"`        // 可选	通知栏样式类型	默认为0，还有1，2，3可选，用来指定选择哪种通知栏样式，其他值无效。有三种可选分别为bigText=1，Inbox=2，bigPicture=3。
	AlertType  int                    `json:"alert_type,omitempty"`   // 可选	通知提醒方式	可选范围为 -1 ～ 7 ，对应 Notification.DEFAULT_ALL = -1 或者 Notification.DEFAULT_SOUND = 1, Notification.DEFAULT_VIBRATE = 2, Notification.DEFAULT_LIGHTS = 4 的任意 “or” 组合。默认按照 -1 处理。
	BigText    string                 `json:"big_text,omitempty"`     // 可选	大文本通知栏样式	当 style = 1 时可用，内容会被通知栏以大文本的形式展示出来。支持 api 16以上的rom。
	BigPicPath string                 `json:"big_pic_path,omitempty"` // 可选	大图片通知栏样式	当 style = 3 时可用，可以是网络图片 url，或本地图片的 path，目前支持.jpg和.png后缀的图片。图片内容会被通知栏以大图片的形式展示出来。如果是 http／https 的url，会自动下载；如果要指定开发者准备的本地图片就填sdcard 的相对路径。支持 api 16以上的rom。
	Inbox      map[string]interface{} `json:"inbox,omitempty"`        // 可选	文本条目通知栏样式	当 style = 2 时可用， json 的每个 key 对应的 value 会被当作文本条目逐条展示。支持 api 16以上的rom。
	Extras     map[string]interface{} `json:"extras,omitempty"`       // 可选	扩展字段	这里自定义 JSON 格式的 Key/Value 信息，以供业务使用。
}

// IOSNotice 苹果通知
type IOSNotice struct {
	Alert            interface{}            `json:"alert"`                       //	必填	通知内容	这里指定内容将会覆盖上级统一指定的 alert 信息；内容为空则不展示到通知栏。支持字符串形式也支持官方定义的alert payload 结构
	Sound            string                 `json:"sound,omitempty"`             //	可选	通知提示声音	如果无此字段，则此消息无声音提示；有此字段，如果找到了指定的声音就播放该声音，否则播放默认声音,如果此字段为空字符串，iOS 7 为默认声音，iOS 8及以上系统为无声音。(消息) 说明：JPush 官方 API Library (SDK) 会默认填充声音字段。提供另外的方法关闭声音。
	Badge            string                 `json:"badge,omitempty"`             //	可选	应用角标	如果不填，表示不改变角标数字；否则把角标数字改为指定的数字；为 0 表示清除。JPush 官方 API Library(SDK) 会默认填充badge值为"+1",详情参考：badge +1
	ContentAvailable bool                   `json:"content-available,omitempty"` //	可选	推送唤醒	推送的时候携带"content-available":true 说明是 Background Remote Notification，如果不携带此字段则是普通的Remote Notification。详情参考：Background Remote Notification
	MutableContent   bool                   `json:"mutable-content,omitempty"`   //	可选	通知扩展	推送的时候携带”mutable-content":true 说明是支持iOS10的UNNotificationServiceExtension，如果不携带此字段则是普通的Remote Notification。详情参考：UNNotificationServiceExtension
	Category         string                 `json:"category,omitempty"`          //	可选	IOS8才支持。设置APNs payload中的"category"字段值
	Extras           map[string]interface{} `json:"extras,omitempty"`            //	可选	附加字段	这里自定义 Key/value 信息，以供业务使用。
}

// IOSAlert 苹果通知alert
type IOSAlert struct {
	Title        string   `json:"title"`          // A short string describing the purpose of the notification. Apple Watch displays this string as part of the notification interface. This string is displayed only briefly and should be crafted so that it can be understood quickly. This key was added in iOS 8.2.
	Body         string   `json:"body"`           // The text of the alert message.
	TitleLocKey  string   `json:"title-loc-key"`  // String or null The key to a title string in the Localizable.strings file for the current localization. The key string can be formatted with %@ and %n$@ specifiers to take the variables specified in the title-loc-args array. See Localizing the Content of Your Remote Notifications for more information. This key was added in iOS 8.2.
	TitleLocArgs []string `json:"title-loc-args"` // Array of strings or null Variable string values to appear in place of the format specifiers in title-loc-key. See Localizing the Content of Your Remote Notifications for more information. This key was added in iOS 8.2.
	ActionLocKey string   `json:"action-loc-key"` // String or null If a string is specified, the system displays an alert that includes the Close and View buttons. The string is used as a key to get a localized string in the current localization to use for the right button’s title instead of “View”. See Localizing the Content of Your Remote Notifications for more information.
	LocKey       string   `json:"loc-key"`        // String A key to an alert-message string in a Localizable.strings file for the current localization (which is set by the user’s language preference). The key string can be formatted with %@ and %n$@ specifiers to take the variables specified in the loc-args array. See Localizing the Content of Your Remote Notifications for more information.
	LocArgs      []string `json:"loc-args"`       // Array of strings Variable string values to appear in place of the format specifiers in loc-key. See Localizing the Content of Your Remote Notifications for more information.
	LaunchImage  string   `json:"launch-image"`   // String The filename of an image file in the app bundle, with or without the filename extension. The image is used as the launch image when users tap the action button or move the action slider. If this property is not specified, the system either uses the previous snapshot, uses the image identified by the UILaunchImageFile key in the app’s Info.plist file, or falls back to Default.png.
}

// WPNotice winphone通知
type WPNotice struct {
	Alert    string                 `json:"alert"`                // 必填	通知内容	会填充到 toast 类型 text2 字段上。这里指定了，将会覆盖上级统一指定的 alert 信息；内容为空则不展示到通知栏。
	Title    string                 `json:"title,omitempty"`      // 可选	通知标题	会填充到 toast 类型 text1 字段上。
	OpenPage string                 `json:"_open_page,omitempty"` // 可选	点击打开的页面名称	点击打开的页面。会填充到推送信息的 param 字段上，表示由哪个 App 页面打开该通知。可不填，则由默认的首页打开。
	Extras   map[string]interface{} `json:"extras,omitempty"`     // 可选	扩展字段	作为参数附加到上述打开页面的后边。
}

// Message 推送自定义消息
type Message struct {
	Content     string                 `json:"msg_content"`            // 必填 消息内容本身
	Title       string                 `json:"title,omitempty"`        // 可选 消息标题
	ContentType string                 `json:"content_type,omitempty"` // 可选 消息内容类型
	Extras      map[string]interface{} `json:"extras,omitempty"`       // 可选 扩展字段
}

// SMSMessage 推送短信通知
type SMSMessage struct {
	Content   string `json:"content"`    //	必填	不能超过480个字符。"你好,JPush"为8个字符。70个字符记一条短信费，如果超过70个字符则按照每条67个字符拆分，逐条计费。单个汉字、标点、英文都算一个字。
	DelayTime int    `json:"delay_time"` //	 必填	单位为秒，不能超过24小时。设置为0，表示立即发送短信。该参数仅对android平台有效，iOS 和 Winphone平台则会立即发送短信
}

// Options 推送选项
type Options struct {
	Sendno          int    `json:"sendno,omitempty"`            // 	可选	推送序号	纯粹用来作为 API 调用标识，API 返回时被原样返回，以方便 API 调用方匹配请求与返回。值为 0 表示该 messageid 无 sendno，所以字段取值范围为非 0 的 int.
	TimeToLive      int    `json:"time_to_live,omitempty"`      // 	可选	离线消息保留时长(秒)	推送当前用户不在线时，为该用户保留多长时间的离线消息，以便其上线时再次推送。默认 86400 （1 天），最长 10 天。设置为 0 表示不保留离线消息，只有推送当前在线的用户可以收到。
	OverrideMsgID   int64  `json:"override_msg_id,omitempty"`   // 	可选	要覆盖的消息ID	如果当前的推送要覆盖之前的一条推送，这里填写前一条推送的 msg_id 就会产生覆盖效果，即：1）该 msg_id 离线收到的消息是覆盖后的内容；2）即使该 msg_id Android 端用户已经收到，如果通知栏还未清除，则新的消息内容会覆盖之前这条通知；覆盖功能起作用的时限是：1 天。如果在覆盖指定时限内该 msg_id 不存在，则返回 1003 错误，提示不是一次有效的消息覆盖操作，当前的消息不会被推送。
	ApnsProduction  bool   `json:"apns_production"`             // 	可选	APNs是否生产环境	True 表示推送生产环境，False 表示要推送开发环境；如果不指定则为推送生产环境。但注意，JPush 服务端 SDK 默认设置为推送 “开发环境”。
	ApnsCollapseID  string `json:"apns_collapse_id,omitempty"`  // 	可选	更新 iOS 通知的标识符	APNs 新通知如果匹配到当前通知中心有相同 apns-collapse-id 字段的通知，则会用新通知内容来更新它，并使其置于通知中心首位。collapse id 长度不可超过 64 bytes。
	BigPushDuration int    `json:"big_push_duration,omitempty"` // 	可选	定速推送时长(分钟)	又名缓慢推送，把原本尽可能快的推送速度，降低下来，给定的n分钟内，均匀地向这次推送的目标用户推送。最大值为1400.未设置则不是定速推送。
}

// JPushRet 返回结构
type JPushRet struct {
	SendNO string      `json:"sendno,omitempty"`
	MsgID  interface{} `json:"msg_id,omitempty"`
	Err    struct {
		Msg  string `json:"message"`
		Code int    `json:"code"`
	} `json:"error,omitempty"`
}
