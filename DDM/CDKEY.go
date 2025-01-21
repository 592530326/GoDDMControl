package DDM

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func MD5(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}
func F数据处理(s string) string {
	// 匹配所有空格、制表符、换行符和其他不可见字符
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(s, "")
}

func 卡密心跳(卡密, token, 项目密钥, 项目UUID string, 心跳回调函数 func(S卡密心跳回调结构体)) {
	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop()
	for {

		设备UUID, err := F读取设备UUID()
		if err != nil {
			var 心跳结果 S卡密心跳回调结构体
			心跳结果.Code = -8888
			心跳结果.Msg = "设备信息文件被删除,请重新启动程序"
			心跳回调函数(心跳结果)
			return
		}
		var 时间戳 = strconv.FormatInt(time.Now().Unix(), 10)
		var Sign = MD5(卡密 + token + 项目密钥 + 项目UUID + 设备UUID + 时间戳)

		var postData = "Sign=%s&CDKEY=%s&CDKEYDeviceInfo=%s&ProjectUUID=%s&Token=%s&Timestamp=%s"
		postData = fmt.Sprintf(postData, Sign, 卡密, 设备UUID, 项目UUID, token, 时间戳)
		var 请求地址 = "http://" + IP + ":" + PORT + "/cdkey/v2/script/verify/heartbeat"
		请求结果, code := HttpPost(请求地址, postData, 60)
		LOG(请求结果)
		if code == 200 {
			var 心跳结果 S卡密心跳回调结构体
			err := json.Unmarshal([]byte(请求结果), &心跳结果)
			if err != nil {
				LOG("服务器返回数据解析错误,正在重试", err, 请求结果)
			}
			if 心跳结果.Code == 1 {
				var 本地Sign = 卡密 + token + 项目密钥 + 项目UUID + 设备UUID + 心跳结果.Timestamp
				if MD5(本地Sign) != 心跳结果.Sign {
					心跳结果.Code = -6666
					心跳结果.Msg = "有人尝试破解卡密系统"
					心跳回调函数(心跳结果)
					return
				} else {
					心跳回调函数(心跳结果)
				}
			} else {
				心跳回调函数(心跳结果)
				return
			}
		} else {
			LOG("网络错误,卡密登录失败,正在重试", err)
			var 心跳结果 S卡密心跳回调结构体
			心跳结果.Code = -9999
			心跳结果.Msg = "访问服务器失败,心跳错误"
			心跳回调函数(心跳结果)
		}
		<-ticker.C
	}

}
func F卡密_卡密登录(项目UUID, 项目加密密钥, 卡密 string, 心跳回调函数 func(S卡密心跳回调结构体)) (S卡密登录结果结构体, error) {
	var 登录结果 S卡密登录结果结构体
	卡密 = F数据处理(卡密)
	项目UUID = F数据处理(项目UUID)
	项目加密密钥 = F数据处理(项目加密密钥)
	设备UUID, err := F读取设备UUID()
	if err != nil {
		LOG("卡密登录失败,读取设备UUID出现错误", err)
		return 登录结果, errors.New("卡密登录失败,读取设备UUID出现错误," + err.Error())
	}
	时间戳 := strconv.FormatInt(time.Now().Unix(), 10)
	var Sign = MD5(卡密 + 项目加密密钥 + 项目UUID + 设备UUID + 时间戳)

	var postData = "Sign=" + Sign + "&Timestamp=" + 时间戳 + "&CDKEY=" + 卡密 + "&CDKEYDeviceInfo=" + 设备UUID + "&ProjectUUID=" + 项目UUID
	var 请求地址 = "http://" + IP + ":" + PORT + "/cdkey/v2/script/verify/login"
	for {
		请求结果, code := HttpPost(请求地址, postData, 60)
		if code == 200 {
			err := json.Unmarshal([]byte(请求结果), &登录结果)
			if err != nil {
				LOG("服务器返回数据解析错误,正在重试", err, 请求结果)
			}
			LOG(请求结果, 登录结果)
			if 登录结果.Code == 1 {
				var 本地Sign = 卡密 + 登录结果.Token + 项目加密密钥 + 项目UUID + 设备UUID + 登录结果.S时间戳
				if MD5(本地Sign) != 登录结果.Sign {
					return S卡密登录结果结构体{}, errors.New("卡密登录失败,签名校验失败")
				} else {
					go 卡密心跳(卡密, 登录结果.Token, 项目加密密钥, 项目UUID, 心跳回调函数)
					return 登录结果, nil
				}
			}

		} else {
			LOG("网络错误,卡密登录失败,正在重试", err)
		}
		time.Sleep(time.Second)
	}
}
