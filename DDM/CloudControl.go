package DDM

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"sync"
	"time"
)

var G_云控回调函数 = func(ctx context.Context, 自定义云控命令名称, 调用参数 string) {}

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
	ws.lock.Lock()
	m["SessinId"] = ws.DDMcontrolSessinId
	err := ws.conn.WriteMessage(websocket.BinaryMessage, JsonEncode(&m))
	ws.lock.Unlock()
	return err
}

var DDMWSConn S云控对象

func F绑定用户后台(用户后台用户名, 云控UUID, 卡密UUID string, 不循环绑定 bool) (bool, error) {
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
func F云控_连接云控系统(云控回调函数 func(ctx context.Context, 自定义云控命令名称, 调用参数 string), 基础云控回调函数自定义 func(基础事件名称, 事件值 string)) bool {
	G_云控回调函数 = 云控回调函数
	设备UUID, err := F读取设备UUID()
	if err != nil {
		LOG("连接云控,读取设备UUID出现错误", err)
		return false
	}
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
		for {
			_, message, err := 云控连接.ReadMessage()
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
					if "云控按钮命令" == OrderName && "停止运行" == OrderValue {
						if DDMWSConn.CancelContext != nil {
							DDMWSConn.CancelContext()
						}
					} else {
						cxt, cancel := context.WithCancel(context.Background())
						DDMWSConn.CancelContext = cancel
						go G_云控回调函数(cxt, OrderName, OrderValue)
					}
					var m = make(map[string]interface{})
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

}
