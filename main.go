package main

import (
	"GoDDMControl/DDM"
	"context"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func 云控连接回调(ctx context.Context, 自定义云控命令名称, 调用参数 string) {
	DDM.LOG(自定义云控命令名称, 调用参数)
}
func 自定义云控连接回调(自定义云控命令名称, 调用参数 string) {
	DDM.LOG(自定义云控命令名称, 调用参数)
}
func main() {
	DDM.F初始化("127.0.0.1", "9000", "2bcdbcebf5b14d65a659292fa721d507", "d99f33a3-4b17-49f8-891a-bedf77886598", "", 60)
	DDM.F绑定用户后台("123456", "2bcdbcebf5b14d65a659292fa721d507", "6a012cc5-917d-4bcb-b83d-975b3955cd07", true)
	DDM.F云控_连接云控系统(云控连接回调, nil)
	DDM.F读取设备UUID()
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
