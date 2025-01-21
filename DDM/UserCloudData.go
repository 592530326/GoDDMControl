package DDM

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"time"
)

func F用户云数据_登录账号(账号, 密码, 签名密钥 string) (string, error) {

	var md5hash = md5.New()
	var 时间戳 = time.Now().Format("2006-01-02 15:04:05")
	md5hash.Write([]byte(签名密钥 + 时间戳 + 账号 + 密码))
	Sign := hex.EncodeToString(md5hash.Sum(nil))

	var postData = "UserName=%s&PassWord=%s&Sign=%s&T=%s"
	postData = fmt.Sprintf(postData, 账号, 密码, Sign, 时间戳)
	请求结果, code := HttpPost("http://"+IP+":"+PORT+"/coludControl/v2/script/UserAdminLogin", postData, 60)
	if 结果, err := JsonDecode([]byte(请求结果)); code == 200 && err == nil {
		return 结果["token"].(string), nil
	} else {
		return "", errors.New(请求结果)
	}
}
func F用户云数据_读取一条未读取的数据(Token string, 项目名称 string, table条件 map[string]interface{}, 当无数据时设置数据为未读取 bool) (map[string]interface{}, error) {
	设备UUID, err := F读取设备UUID()
	if err != nil {
		LOG("连接云控,读取设备UUID出现错误", err)
		return nil, err
	}
	Where := JsonEncode(table条件)
	var postData = "SetFalse=%t&DeviceInfo=%s&Where=%s&Token=%s&ProjectName=%s"
	postData = fmt.Sprintf(postData, 当无数据时设置数据为未读取, 设备UUID, Where, Token, url.QueryEscape(项目名称))
	请求结果, code := HttpPost("http://"+IP+":"+PORT+"/userCloudData/v2/script/readOneDataSetRead", postData, 60)
	if 结果, err := JsonDecode([]byte(请求结果)); code == 200 && err == nil {
		return 结果, nil
	} else {
		return nil, errors.New(请求结果)
	}
}

func F用户云数据_搜索数据(Token string, 项目名称 string, table条件 map[string]interface{}, 页编号, 每页数量 int) (map[string]interface{}, error) {
	设备UUID, err := F读取设备UUID()
	if err != nil {
		LOG("连接云控,读取设备UUID出现错误", err)
		return nil, err
	}
	Where := JsonEncode(table条件)

	var postData = "DeviceInfo=%s&Where=%s&Page=%d&PageSize=%d&Token=%s&ProjectName=%s"

	postData = fmt.Sprintf(postData, 设备UUID, Where, 页编号, 每页数量, Token, url.QueryEscape(项目名称))
	LOG(postData)
	请求结果, code := HttpPost("http://"+IP+":"+PORT+"/userCloudData/v2/script/readManyData", postData, 60)
	LOG(请求结果)
	if 结果, err := JsonDecode([]byte(请求结果)); code == 200 && err == nil {
		return 结果, nil
	} else {
		return nil, errors.New(请求结果)
	}
}

type 数据类型 string

const (
	G字符串  数据类型 = "字符串"
	G数字   数据类型 = "数字"
	G图片   数据类型 = "图片"
	G数组   数据类型 = "数组"
	G超链接  数据类型 = "超链接"
	G日期时间 数据类型 = "日期时间"
	G逻辑   数据类型 = "逻辑"
)

func F用户云数据_创建数据(类型 数据类型, 数据, 背景颜色, 字体颜色 string) map[string]interface{} {

	var m = make(map[string]interface{})
	m["value"] = 数据
	if 类型 == G数组 {
		m["valueType"] = 类型
		return m
	}
	m["bgcolor"] = 背景颜色
	m["textcolor"] = 字体颜色
	m["valueType"] = 类型
	return m
}
