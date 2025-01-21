package DDM

type S卡密心跳回调结构体 struct {
	Code          int    `json:"code"`
	Msg           string `json:"msg"`
	Timestamp     string `json:"Timestamp"`
	Sign          string `json:"Sign"`
	RemainingTime int    `json:"RemainingTime"`
	EndTime       string `json:"endTime"`
}

type S卡密登录结果结构体 struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	S剩余时间     int    `json:"RemainingTime"`
	Sign      string `json:"Sign"`
	S时间戳      string `json:"Timestamp"`
	Token     string `json:"Token"`
	S已经使用窗口数量 int    `json:"UseWindow"` //只针对限制设备卡有意义,返回的是已经绑定的设备,非在线设备
	S可使用窗口数量  int    `json:"WindowNumber"`
	S到期时间     string `json:"endTime"`
}
