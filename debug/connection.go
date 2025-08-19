package main

import (
	"fmt"
	"log"
	"time"

	"github.com/zhoudm1743/submail"
)

func main() {
	fmt.Println("SUBMAIL Go SDK 连接诊断工具")
	fmt.Println("================================")

	// 配置客户端
	config := submail.Config{
		AppID:          "112455",                           // 替换为您的App ID
		AppKey:         "a086473503b38c3ddd71ee38c44819f6", // 替换为您的App Key
		BaseURL:        submail.DefaultBaseURL,             // 使用默认API地址
		Format:         submail.FormatJSON,                 // 使用JSON格式
		UseDigitalSign: true,                               // 使用明文模式（推荐测试时使用）
		SignType:       submail.SignTypeMD5,                // MD5签名（数字签名模式时使用）
		Timeout:        30 * time.Second,                   // 30秒超时
	}

	client := submail.NewClient(config)

	// 运行诊断
	fmt.Printf("\n开始诊断...\n")
	if err := client.DiagnoseConnection(); err != nil {
		log.Printf("❌ 诊断失败: %v", err)

		fmt.Printf("\n故障排除建议:\n")
		fmt.Printf("1. 检查网络连接是否正常\n")
		fmt.Printf("2. 确认防火墙没有阻止HTTPS连接\n")
		fmt.Printf("3. 尝试在浏览器访问: https://api-v4.mysubmail.com/service/timestamp.json\n")
		fmt.Printf("4. 如果使用代理，请检查代理设置\n")
		fmt.Printf("5. 尝试增加超时时间（如60秒）\n")

		return
	}

	fmt.Printf("\n✅ 所有测试通过！SDK可以正常工作。\n")
	fmt.Printf("\n现在您可以使用真实的AppID和AppKey进行短信发送了。\n")
}
