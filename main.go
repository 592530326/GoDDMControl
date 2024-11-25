package main

import (
	"GoDDMControl/DDM"
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
func main() {
	DDM.F初始化("127.0.0.1", "9000", "2bcdbcebf5b14d65a659292fa721d507", "d99f33a3-4b17-49f8-891a-bedf77886598", "", 60)
	DDM.F绑定用户后台("592530326", "2bcdbcebf5b14d65a659292fa721d507", "cd8a0ccd-1ead-4a53-825b-1a27d167cac2", true)
	DDM.F云控_连接云控系统(云控连接回调, nil)
	time.Sleep(time.Second)
	DDM.F云控_上传脚本状态("等待任务")
	DDM.F云控_上传运行日志("等待任务", "", "")
	DDM.F云控_修改设备名字("等待任务")
	Token, err := DDM.F用户云数据_登录账号("592530326", "592530326", "592530326")
	if err != nil {
		DDM.LOG(err)
	}
	DDM.LOG(Token)
	var 查询条件 = make(map[string]interface{})
	查询条件["密码"] = DDM.F用户云数据_创建数据(DDM.G字符串, "muZzmR36", "", "")
	一条未读取的数据, err := DDM.F用户云数据_读取一条未读取的数据(Token, "微软邮箱", 查询条件, false)
	if err != nil {
		DDM.LOG(err)
		return
	}
	DDM.LOG(一条未读取的数据)
	select {}
	DDM.F读取设备UUID()

}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
