package DDM

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"sync"
	"time"
)

var G_云控回调函数 = func(自定义云控命令名称, 调用参数 string) {}

type S云控对象 struct {
	conn               *websocket.Conn
	lock               sync.Mutex
	CancelContext      context.CancelFunc
	DeviceInfo         string
	DDMcontrolSessinId string
}

func (ws *S云控对象) SetConn(conn *websocket.Conn) {
	if ws.conn != nil {
		ws.conn.Close()
	}
	ws.conn = conn
}
func (ws *S云控对象) Send(m map[string]interface{}) error {
	if ws.conn == nil {
		return errors.New("云控暂未连接成功,无法发送数据")
	}
	ws.lock.Lock()
	m["SessinId"] = ws.DDMcontrolSessinId
	err := ws.conn.WriteMessage(websocket.BinaryMessage, JsonEncode(&m))
	ws.lock.Unlock()
	return err
}

var DDMWSConn S云控对象

func F绑定用户后台(用户后台用户名, 云控UUID, 卡密UUID string, 不循环绑定 bool) (bool, error) {
	G_用户后台用户名 = 用户后台用户名
	var postData = "DeviceInfo=%s&UserAdminName=%s&CloudControlUUID=%s&CDKEYUUID=%s"
	设备UUID, err := F读取设备UUID()
	if err != nil {
		LOG("F绑定用户后台,读取设备UUID出现错误", err)
		return false, err
	}
	var 请求地址 = "http://" + IP + ":" + PORT + "/coludControl/v2/script/BindCloudControlAndCDKEYAndDevice"
	postData = fmt.Sprintf(postData, 设备UUID, url.QueryEscape(用户后台用户名), 云控UUID, 卡密UUID)
	for {
		请求结果, code := HttpPost(请求地址, postData, 60)
		if 结果, err := JsonDecode([]byte(请求结果)); code == 200 && err == nil {
			初始化状态CODE, _ := 结果["code"].(float64)
			if 初始化状态CODE == 200 {
				return true, nil
			} else {
				LOG("绑定云控出现错误,解析格式错误", err)
				if 不循环绑定 {
					return false, err
				}
			}
		} else {
			LOG("绑定云控出现错误,解析格式错误", err)
			if 不循环绑定 {
				return false, err
			}
		}
		time.Sleep(time.Second)
	}
}
func F云控_修改设备名字(设备名字 string) error {
	设备UUID, err := F读取设备UUID()
	if err != nil {
		LOG("连接云控,读取设备UUID出现错误", err)
		return err
	}
	var postData = "DeviceInfo=%s&UserUUID=%s&DeviceName=%s"
	postData = fmt.Sprintf(postData, 设备UUID, G_作者UUID, url.QueryEscape(设备名字))
	请求结果, code := HttpPost("http://"+IP+":"+PORT+"/ChangeDeviceName", postData, 60)
	if _, err := JsonDecode([]byte(请求结果)); code == 200 && err == nil {
		return nil
	} else {
		return err
	}

}
func F云控_上传运行日志(data, 字体颜色, 背景颜色 string) bool {
	设备UUID, err := F读取设备UUID()
	if err != nil {
		LOG("连接云控,读取设备UUID出现错误", err)
		return false
	}
	data = url.QueryEscape(data)
	m := make(map[string]interface{})
	m["op"] = "ScriptDIYState"

	m["op"] = 10016
	m["data"] = data
	m["DeviceInfo"] = 设备UUID
	m["FontColor"] = 字体颜色
	m["Background"] = 背景颜色
	err = DDMWSConn.Send(m)
	if err != nil {
		return false
	}
	return true
}
func F云控_上传脚本状态(状态码 string) bool {
	设备UUID, err := F读取设备UUID()
	if err != nil {
		LOG("连接云控,读取设备UUID出现错误", err)
		return false
	}

	m := make(map[string]interface{})
	m["Cmd"] = "ScriptDIYState"
	m["StateCode"] = 状态码
	m["DeviceInfo"] = 设备UUID
	err = DDMWSConn.Send(m)
	if err != nil {
		return false
	}
	return true
}

func F云控_连接云控系统(云控回调函数 func(自定义云控命令名称, 调用参数 string), 基础云控回调函数自定义 func(基础事件名称, 事件值 string)) bool {
	G_云控回调函数 = 云控回调函数
	设备UUID, err := F读取设备UUID()
	if err != nil {
		LOG("连接云控,读取设备UUID出现错误", err)
		return false
	}
	go func() {
		for {
			urlStr := "ws://" + IP + ":" + PORT + "/Control"
			LOG(urlStr)
			云控连接, _, err := websocket.DefaultDialer.Dial(urlStr, nil)
			if err != nil {
				LOG("云控连接失败,3秒后重连...", err.Error())
				time.Sleep(time.Second * 3)
				continue
			}
			DDMWSConn.DeviceInfo = 设备UUID
			DDMWSConn.SetConn(云控连接)

			var m = make(map[string]interface{})
			m["op"] = 2
			m["CloudControlUUID"] = G_云控UUID
			m["UserUUID"] = G_作者UUID
			m["DeviceInfo"] = 设备UUID
			m["PackageName"] = PKG
			DDMWSConn.Send(m)
			if G_用户后台用户名 != "" {
				m = make(map[string]interface{})
				m["Cmd"] = "BindUserAdminNameReq"
				m["UserAdminName"] = G_用户后台用户名
				m["DeviceInfo"] = 设备UUID
				DDMWSConn.Send(m)
			}
			for {
				_, message, err := 云控连接.ReadMessage()
				LOG(string(message))
				if err != nil {
					云控连接.Close()
					LOG("云控与服务器断开连接...", err.Error())
					break
				}
				data, err := JsonDecode(message)
				if err != nil {
					LOG("云控连接数据解析格式不正确,断开连接", err.Error())
					continue
				}
				op, ok := data["op"].(float64)
				if ok {
					switch op {
					case 200005: //云控自定义命令
						OrderName, _ := data["OrderName"].(string)
						OrderValue, _ := data["OrderValue"].(string)

						go G_云控回调函数(OrderName, OrderValue)
						m = make(map[string]interface{})
						m["DeviceInfo"] = DDMWSConn.DeviceInfo
						m["op"] = 200006
						err := DDMWSConn.Send(m)
						if err != nil {
							LOG("发送数据出现错误,断开连接...", err.Error())
							DDMWSConn.conn.Close()
							break
						}

					}
				}
			}

		}

	}()
	return true
}
