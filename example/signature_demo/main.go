// Copyright 2025 zhoudm1743
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net/url"

	"github.com/zhoudm1743/submail"
)

func main() {
	fmt.Println("=== SUBMAIL API 签名算法演示 ===")

	// 使用您的实际AppID和AppKey
	appID := "your_app_id"
	appKey := "your_app_key"

	fmt.Println("\n1. 数字签名模式演示 (MD5)")
	// 创建数字签名模式的服务实例（默认MD5）
	digitalService := submail.NewSaiyouService(appID, appKey)

	// 演示签名生成过程
	params := url.Values{}
	params.Set("to", "13800138000")
	params.Set("text", "您的验证码是：123456")
	params.Set("tag", "test_tag") // tag参数不参与签名计算

	signString, signature := digitalService.ValidateSignature(params)
	fmt.Printf("签名字符串: %s\n", signString)
	fmt.Printf("生成的签名: %s\n", signature)

	fmt.Println("\n2. 数字签名模式演示 (SHA1)")
	// 切换为SHA1签名
	digitalService.SetAuthMode(true, "sha1")
	signString2, signature2 := digitalService.ValidateSignature(params)
	fmt.Printf("SHA1签名字符串: %s\n", signString2)
	fmt.Printf("SHA1生成的签名: %s\n", signature2)

	fmt.Println("\n3. 明文验证模式演示")
	// 创建明文验证模式的服务实例
	plaintextService := submail.NewSaiyouServiceWithPlaintextAuth(appID, appKey)

	signString3, signature3 := plaintextService.ValidateSignature(params)
	fmt.Printf("明文模式签名字符串: %s\n", signString3)
	fmt.Printf("明文模式签名: %s\n", signature3)

	fmt.Println("\n4. 验证模式切换演示")
	// 动态切换验证模式
	service := submail.NewSaiyouService(appID, appKey)

	// 获取当前验证模式
	useDigital, signType := service.GetAuthMode()
	fmt.Printf("当前验证模式: 数字签名=%v, 签名类型=%s\n", useDigital, signType)

	// 切换为明文模式
	service.SetAuthMode(false, "")
	useDigital, signType = service.GetAuthMode()
	fmt.Printf("切换后验证模式: 数字签名=%v, 签名类型=%s\n", useDigital, signType)

	// 切换回SHA1数字签名模式
	service.SetAuthMode(true, "sha1")
	useDigital, signType = service.GetAuthMode()
	fmt.Printf("再次切换后验证模式: 数字签名=%v, 签名类型=%s\n", useDigital, signType)

	fmt.Println("\n=== 签名算法说明 ===")
	fmt.Println("1. 明文验证模式：")
	fmt.Println("   - 直接在signature参数中提交appkey")
	fmt.Println("   - 集成简单，但安全性较低")
	fmt.Println("   - 不需要timestamp和sign_type参数")

	fmt.Println("\n2. 数字签名验证模式：")
	fmt.Println("   - 按照官方算法生成签名：appid + appkey + signature_string + appid + appkey")
	fmt.Println("   - 参数按字典序排列，tag参数不参与签名计算")
	fmt.Println("   - 需要timestamp参数，建议从服务器获取时间戳")
	fmt.Println("   - 需要sign_type参数（md5或sha1）")
	fmt.Println("   - 安全性高，推荐生产环境使用")

	fmt.Println("\n=== 使用建议 ===")
	fmt.Println("- 开发测试阶段：可使用明文模式，简化集成")
	fmt.Println("- 生产环境：强烈建议使用数字签名模式")
	fmt.Println("- 时间敏感应用：使用数字签名模式并从服务器同步时间")
}