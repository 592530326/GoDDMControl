package main

import (
	"encoding/json"
	"github.com/592530326/GoDDMControl/DDM"
	"os"
	"time"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func 云控连接回调(自定义云控命令名称, 调用参数 string) {
	DDM.LOG(自定义云控命令名称, 调用参数)
}
func 自定义云控连接回调(自定义云控命令名称, 调用参数 string) {
	DDM.LOG(自定义云控命令名称, 调用参数)
}

var 失败次数 = 0

func 心跳回调函数(心跳结果 DDM.S卡密心跳回调结构体) {
	DDM.LOG(心跳结果)
	switch 心跳结果.Code {
	case 1:
		失败次数 = 0
		DDM.LOG("心跳成功")
	case -5:
		DDM.LOG("请重新登录,一般是卡密被禁用,删除,设备被解绑!")
		os.Exit(0)
	case -8:
		DDM.LOG("卡密到期!")
		os.Exit(0)
	case -9999:
		DDM.LOG("心跳失败,网络错误!")
		失败次数 = 失败次数 + 1
		if 失败次数 >= 30 { //连续心跳失败30次(每次60秒),停止运行
			os.Exit(0)
		}
	case -11:
		DDM.LOG("错误原因:" + 心跳结果.Msg)
		os.Exit(0)
	case -6666:
		DDM.LOG("有人尝试破解卡密系统!")
		os.Exit(0)
	default:
		DDM.LOG("错误原因:" + 心跳结果.Msg)
		os.Exit(0)
	}
}
func mapToJSON(tempMap *map[string]interface{}) string {
	data, err := json.Marshal(tempMap)

	if err != nil {
		panic(err)
	}

	return string(data)
}
func main() {

	DDM.F初始化("192.168.3.14", "9000", "68896ff898d74a91b2b8e6f8a5b850dd", "d99f33a3-4b17-49f8-891a-bedf77886598", "", 60)
	var 热更新回调函数 = func(进度 float64) {
		//DDM.LOG("下载进度:", 进度)
	}
	状态, 信息 := DDM.F热更新_检测热更新并直接更新(-1, "auto.go", 热更新回调函数)
	DDM.LOG(状态, 信息, "1")
	for {
		time.Sleep(time.Second)
	}
	os.Exit(0)
	//DDM.F卡密_卡密登录("6a012cc5-917d-4bcb-b83d-975b3955cd07", "a0ee6b4e-7eb2-4fb0-953b-694aa50a5a0e", "2xvd5h", 心跳回调函数)
	DDM.F绑定用户后台("592530326", "2bcdbcebf5b14d65a659292fa721d507", "cd8a0ccd-1ead-4a53-825b-1a27d167cac2", true)
	DDM.F云控_连接云控系统(云控连接回调, nil)
	//time.Sleep(time.Second)
	//DDM.F云控_上传脚本状态("等待任务")
	//DDM.F云控_上传运行日志("等待任务", "", "")
	//DDM.F云控_修改设备名字("等待任务")

	Token, err := DDM.F用户云数据_登录账号("kuaiwan01", "kuaiwan01", "592530326")
	if err != nil {
		DDM.LOG(err)
	}
	DDM.LOG(Token)
	var 查询条件 = make(map[string]interface{})
	//查询条件["密码"] = DDM.F用户云数据_创建数据(DDM.G字符串, "muZzmR36", "", "")
	//一条未读取的数据, err := DDM.F用户云数据_读取一条未读取的数据(Token, "微软邮箱", 查询条件, false)
	//if err != nil {
	//	DDM.LOG(err)
	//	return
	//}
	//DDM.LOG(一条未读取的数据)

	所有数据, err := DDM.F用户云数据_搜索数据(Token, "快玩账号数据", 查询条件, 1, 2)
	DDM.LOG(所有数据)

}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
