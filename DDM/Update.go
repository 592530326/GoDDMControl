package DDM

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

func getFileMd5(path string) (string, error) {
	// 2. 打开可执行文件
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("打开文件失败:", err)
		return "", err
	}
	defer file.Close()

	// 3. 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("读取文件失败:", err)
		return "", err
	}

	// 4. 计算 MD5 值
	hash := md5.Sum(content)
	md5Value := fmt.Sprintf("%x", hash) // 将 MD5 值转换为十六进制字符串
	return md5Value, nil
}
func 获取下载服务器() (string, int) {
	// 构造 URL
	url := fmt.Sprintf("http://%s:%s/HotDownloadServer", IP, PORT)

	// 发起 HTTP GET 请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		// 如果请求失败，返回默认服务器
		defaultServer := fmt.Sprintf("http://%s:%s", IP, PORT)
		return defaultServer, http.StatusInternalServerError
	}
	defer resp.Body.Close()

	// 读取响应体
	成功访问的服务器 := ""
	if resp.StatusCode == http.StatusOK {
		buf := make([]byte, 1024)
		n, _ := resp.Body.Read(buf)
		成功访问的服务器 = string(buf[:n])
	}

	// 如果响应体为空，使用默认服务器
	if 成功访问的服务器 == "" {
		成功访问的服务器 = fmt.Sprintf("http://%s:%s", IP, PORT)
	} else {
		成功访问的服务器 = fmt.Sprintf("http://%s", 成功访问的服务器)
	}

	// 如果状态码是 404，返回默认服务器
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Sprintf("http://%s:%s", IP, PORT), http.StatusNotFound
	}

	// 返回成功访问的服务器和状态码
	return 成功访问的服务器, http.StatusOK
}

// 文件复制函数
func 复制文件(srcFile, dstFile string) error {
	// 打开源文件
	src, err := os.Open(srcFile)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(dstFile)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dst.Close()

	// 复制文件内容
	bytesCopied, err := io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}

	fmt.Printf("文件复制成功，复制了 %d 字节\n", bytesCopied)
	return nil
}

// 获取对象存储下载地址
func 获取对象存储下载地址(作者UUID, PkgName string) string {
	// 构造 URL
	url := fmt.Sprintf("http://%s:%s/ObjectStorageUrl?UserUUID=%s&PkgName=%s", IP, PORT, 作者UUID, PkgName)

	// 显示提示信息（模拟 HUD）
	fmt.Println("正在获取热更新文件对象下载地址....")

	// 发起 HTTP GET 请求
	client := &http.Client{Timeout: 60}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("请求失败:", err)
		return ""
	}
	defer resp.Body.Close()

	// 读取响应体
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("解析 JSON 失败:", err)
		return ""
	}

	// 处理响应
	if code, ok := result["code"].(float64); ok && code == 1 {
		if downloadUrl, ok := result["DownLoadUrl"].(string); ok {
			return downloadUrl
		}
	} else {
		if msg, ok := result["msg"].(string); ok {
			fmt.Printf("下载热更新出现错误: %s\n", msg)
		}
	}

	// 模拟 sleep
	time.Sleep(1 * time.Second)
	return ""
}
func F热更新_检测热更新并直接更新(版本号 int, 包名 string, 下载回调 interface{}) (int, string) {
	for {
		exePath, err := os.Executable()
		if err != nil {
			fmt.Println("获取当前进程路径失败:", err)
			return 500, err.Error()
		}
		// 获取可执行文件的绝对路径
		absPath, err := filepath.Abs(exePath)
		if err != nil {
			fmt.Println("获取绝对路径失败:", err)
			return 500, err.Error()
		}
		fmt.Println("当前进程路径:", absPath)
		fileMd5, err := getFileMd5(absPath)
		if err != nil {
			return 500, err.Error()
		}
		LOG(fileMd5)
		url := "http://" + IP + ":" + PORT + "/CheckUpdate?UserUUID=" + G_作者UUID + "&PkgName=" + 包名
		LOG(url)
		请求结果, code := HttpGet(url, 60)
		LOG(请求结果, code)
		if 初始化结果, err := JsonDecode([]byte(请求结果)); code == 200 && err == nil {
			LOG(初始化结果)
			LrFileMD5 := 初始化结果["UpdateInfo"].(map[string]interface{})["LrFileMD5"].(string)
			if LrFileMD5 == "" {
				return -1, "无需更新,服务器无更新文件"
			}
			if 版本号 != -1 {
				Version := 初始化结果["UpdateInfo"].(map[string]interface{})["Version"].(float64)
				if float64(版本号) > Version {
					return -1, "无需更新,本地版本比服务器新"
				} else if float64(版本号) == Version {
					return -1, "无需更新,本地版本和服务器一样"
				}
			} else {
				if LrFileMD5 == fileMd5 {
					return -1, "无需更新,本地已经是最新版"
				}
			}
			LOG(LrFileMD5)
			downloadUrl := ""
			ObjectStorage, ObjectStorageOk := 初始化结果["UpdateInfo"].(map[string]interface{})["ObjectStorage"].(bool)
			if ObjectStorageOk && ObjectStorage {
				downloadUrl = 获取对象存储下载地址(G_作者UUID, 包名)
			} else {
				下载服务器, _ := 获取下载服务器()
				downloadUrl = 下载服务器 + "/HotDownloadV2?UserUUID=" + G_作者UUID + "&PkgName=" + 包名
			}

			// 创建 HTTP 请求
			resp, err := http.Get(downloadUrl)
			if err != nil {
				LOG("请求失败:", err)
				continue
			}
			defer resp.Body.Close()
			// 检查响应状态码
			if resp.StatusCode != http.StatusOK {
				LOG("服务器返回错误状态码:", resp.Status)
				continue
			}
			// 获取文件总大小
			fileSize, err := strconv.Atoi(resp.Header.Get("Content-Length"))
			if err != nil {
				fileSize = 0
			}
			// 创建本地文件
			outputFile, err := os.Create(absPath + "_bak")
			if err != nil {
				LOG("创建文件失败:", err)
				continue
			}
			defer outputFile.Close()
			var f func(float64)
			var 回调转换结果 bool
			if 下载回调 != nil {
				f, 回调转换结果 = 下载回调.(func(float64))
			}
			// 分块下载文件并显示进度
			buffer := make([]byte, 1024) // 每次读取 1KB
			var downloadedBytes int
			for {
				// 读取数据
				n, err := resp.Body.Read(buffer)
				if err != nil && err != io.EOF {
					LOG("读取数据失败:", err)
					break
				}
				if n != 0 {
					// 写入本地文件
					_, err = outputFile.Write(buffer[:n])
					if err != nil {
						LOG("写入文件失败:", err)
						break
					}
				}

				// 更新已下载字节数
				downloadedBytes += n
				// 计算并显示进度
				progress := float64(downloadedBytes) / float64(fileSize) * 100
				//LOG("下载进度: ", progress)
				if 回调转换结果 {
					f(progress)
				}
				// 如果下载完成，退出循环
				if err == io.EOF {
					// 新进程的脚本内容
					// 新进程的脚本内容
					script := `
#!/system/bin/sh
# 结束当前进程
kill -9 ` + fmt.Sprint(os.Getpid()) + ` 2>/data/local/tmp/replace_and_restart_error.log
# 等待进程结束
sleep 1

# 替换文件（直接覆盖目标文件）
if command -v cp >/dev/null 2>&1; then
    cp -f /data/local/tmp/app_bak /data/local/tmp/app 2>>/data/local/tmp/replace_and_restart_error.log
elif command -v cat >/dev/null 2>&1; then
    cat /data/local/tmp/app_bak > /data/local/tmp/app 2>>/data/local/tmp/replace_and_restart_error.log
elif command -v dd >/dev/null 2>&1; then
    dd if=/data/local/tmp/app_bak of=/data/local/tmp/app 2>>/data/local/tmp/replace_and_restart_error.log
elif command -v mv >/dev/null 2>&1; then
    mv -f /data/local/tmp/app_bak /data/local/tmp/app 2>>/data/local/tmp/replace_and_restart_error.log
else
    echo "没有可用的文件复制或移动命令" >> /data/local/tmp/replace_and_restart_error.log
fi

# 重新启动程序
/data/local/tmp/app & 2>>/data/local/tmp/replace_and_restart_error.log
`
					//确保 必须有cp命令,否则热更失败
					// 将脚本写入临时文件
					tmpScript := "/data/local/tmp/replace_and_restart.sh"
					err := os.WriteFile(tmpScript, []byte(script), 0755)
					if err != nil {
						fmt.Println("写入脚本文件失败:", err)
						break
					}

					// 启动新进程执行脚本
					cmd := exec.Command("sh", tmpScript)
					err = cmd.Start()
					if err != nil {
						fmt.Println("启动新进程失败:", err)
						os.Exit(0)
						//break
					}

					//err := 复制文件(absPath+"_bak", absPath)
					//LOG(err)
					//if err != nil {
					//	LOG(err)
					//	break
					//}
					return 1, "热更完成"
				}

			}

			LOG("重新下载")
		}
	}

}
