package DDM

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var IP = ""
var PORT = ""
var G_云控UUID = ""
var G_作者UUID = ""
var G_云UIUUID = ""
var G_设备UUID存储路径 = "DeviceInfo"
var DEBUG = true
var G_设备UUID string
var PKG = "window.pc.golang"
var G_用户后台用户名 = ""

func LOG(内容 ...any) {
	if DEBUG == false {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	Logtime := time.Now().Format("2006-01-02 15:04:05")
	var LogAgrs []any
	LogAgrs = append(LogAgrs, Logtime)
	LogAgrs = append(LogAgrs, file+":"+strconv.Itoa(line))
	LogAgrs = append(LogAgrs, 内容...)
	fmt.Println(LogAgrs...)
}
func F写入设备UUID(数据 string) {
	err := os.WriteFile(G_设备UUID存储路径, []byte(数据), 0777)
	if err != nil {
		LOG("写入设备UUID出错")
	}
}
func F读取设备UUID() (string, error) {
	data, err := os.ReadFile(G_设备UUID存储路径)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		} else {
			LOG(err)
			return "", err
		}
	}
	return string(data), err
}
func HttpPost(请求地址, 数据 string, 超时时间 int) (string, int) {
	http.DefaultClient.Timeout = time.Second * time.Duration(超时时间)
	resp, err := http.Post(请求地址, "application/x-www-form-urlencoded; charset=UTF-8", strings.NewReader(数据))
	if err != nil {
		return "", -1
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", -1
	}
	return string(body), resp.StatusCode
}
func JsonDecode(数据 []byte) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal(数据, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
func JsonEncode(数据 any) []byte {
	data, err := json.Marshal(数据)
	if err != nil {
		return []byte("")
	}
	return data
}
func F初始化(IP地址 string, PORT端口 string, 云控UUID, 作者UUID, 云UIUUID string, 超时时间 int) bool {
	IP = IP地址
	PORT = PORT端口
	G_云控UUID = 云控UUID
	G_作者UUID = 作者UUID
	G_云UIUUID = 云UIUUID
	设备UUID, err := F读取设备UUID()
	if err != nil {
		LOG(err)
		return false
	}
	var 请求地址 = "http://" + IP + ":" + PORT + "/DeviceInit"

	var postDta = "DeviceInfo=%s&UIProjectUUID=%s&UserUUID=%s&PackageName=%s&CloudControlUUID=%s"
	postDta = fmt.Sprintf(postDta, 设备UUID, 云UIUUID, 作者UUID, PKG, 云控UUID)
	请求结果, code := HttpPost(请求地址, postDta, 超时时间)

	if 初始化结果, err := JsonDecode([]byte(请求结果)); code == 200 && err == nil {
		初始化状态CODE, _ := 初始化结果["code"].(float64)
		if 初始化状态CODE == 1 {
			return true
		}
		if 初始化状态CODE == 2 {
			DeviceInfo, _ := 初始化结果["DeviceInfo"].(string)
			F写入设备UUID(DeviceInfo)
			return true
		}
	}
	return false
}

func F用户后台登录() {

}
