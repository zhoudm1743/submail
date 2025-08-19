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
	"log"
	"net/url"
	"time"

	"github.com/zhoudm1743/submail"
)

func main() {
	// 创建SUBMAIL服务实例（默认JSON格式，使用默认连接池配置，数字签名模式）
	// 请替换为您的实际AppID和AppKey
	service := submail.NewSaiyouService(
		"your_app_id",
		"your_app_key",
	)

	// 或者创建使用明文验证模式的服务实例（更简单，但安全性较低）
	// servicePlaintext := submail.NewSaiyouServiceWithPlaintextAuth(
	//     "your_app_id",
	//     "your_app_key",
	// )

	// 或者创建指定XML格式的服务实例
	// serviceXML := submail.NewSaiyouServiceWithFormat(
	//     "your_app_id",
	//     "your_app_key",
	//     submail.FormatXML,
	// )

	// 动态切换验证模式示例
	// service.SetAuthMode(false, "") // 切换为明文模式
	// service.SetAuthMode(true, "sha1") // 切换为SHA1数字签名模式

	// 示例: 使用自定义连接池配置创建服务实例
	// customPoolConfig := &submail.ConnectionPoolConfig{
	//     MaxIdleConns:        200,              // 增加最大空闲连接数
	//     MaxIdleConnsPerHost: 20,               // 增加每个主机的最大空闲连接数
	//     MaxConnsPerHost:     100,              // 限制每个主机的最大连接数
	//     IdleConnTimeout:     120 * time.Second, // 延长空闲连接超时时间
	//     TLSHandshakeTimeout: 15 * time.Second,  // TLS握手超时时间
	//     DialTimeout:         20 * time.Second,  // 拨号超时时间
	//     KeepAlive:           60 * time.Second,  // TCP Keep-Alive间隔
	//     RequestTimeout:      45 * time.Second,  // 请求超时时间
	// }

	// 使用自定义连接池配置创建服务（注释掉，避免重复创建）
	// serviceWithPool := submail.NewSaiyouServiceWithPool(
	//     "your_app_id",
	//     "your_app_key",
	//     customPoolConfig,
	// )

	// 示例1: 发送普通短信
	fmt.Println("=== 发送普通短信 ===")
	smsReq := &submail.SMSRequest{
		To:   "13800138000", // 收件人手机号
		Text: "您的验证码是1234，请在5分钟内使用。",
		Vars: map[string]string{
			"code": "1234",
			"time": "5",
		},
		Project: "test",    // 项目标记，可选
		Tag:     "example", // 自定义标签，可选
	}

	smsResp, err := service.SendSMS(smsReq)
	if err != nil {
		log.Printf("发送短信失败: %v", err)
	} else {
		fmt.Printf("短信发送成功! 状态: %s, 发送ID: %s, 费用: %d\n",
			smsResp.Status, smsResp.SendID, smsResp.Fee)
	}

	// 示例2: 发送模板短信
	fmt.Println("\n=== 发送模板短信 ===")
	templateReq := &submail.SMSXRequest{
		To:      "13800138000",      // 收件人手机号
		Project: "your_template_id", // 模板ID
		Vars: map[string]string{
			"code": "5678",
			"time": "10",
		},
		Tag: "template_example", // 自定义标签，可选
	}

	templateResp, err := service.SendSMSTemplate(templateReq)
	if err != nil {
		log.Printf("发送模板短信失败: %v", err)
	} else {
		fmt.Printf("模板短信发送成功! 状态: %s, 发送ID: %s, 费用: %d\n",
			templateResp.Status, templateResp.SendID, templateResp.Fee)
	}

	// 示例3: 查询账户余额
	fmt.Println("\n=== 查询账户余额 ===")
	balance, err := service.GetBalance()
	if err != nil {
		log.Printf("查询余额失败: %v", err)
	} else {
		fmt.Printf("账户余额查询成功! 状态: %s, 余额: %s\n",
			balance.Status, balance.Balance)
	}

	// 示例4: 短信一对多发送
	fmt.Println("\n=== 短信一对多发送 ===")
	multiReq := &submail.SMSMultisendRequest{
		Text: "亲爱的{name}，您的验证码是{code}",
		Multi: []submail.SMSMultiItem{
			{
				To: "13800138000",
				Vars: map[string]string{
					"name": "张三",
					"code": "1234",
				},
			},
			{
				To: "13800138001",
				Vars: map[string]string{
					"name": "李四",
					"code": "5678",
				},
			},
		},
		Project: "test_multi",
	}

	multiResp, err := service.SendSMSMulti(multiReq)
	if err != nil {
		log.Printf("一对多发送失败: %v", err)
	} else {
		fmt.Printf("一对多发送成功! 状态: %s\n", multiResp.Status)
		for _, send := range multiResp.Sends {
			fmt.Printf("  收件人: %s, 状态: %s, ID: %s\n",
				send.To, send.Status, send.SendID)
		}
	}

	// 示例5: 模板管理
	fmt.Println("\n=== 短信模板管理 ===")

	// 获取模板列表
	templates, err := service.GetSMSTemplates()
	if err != nil {
		log.Printf("获取模板列表失败: %v", err)
	} else {
		fmt.Printf("模板列表获取成功! 状态: %s\n", templates.Status)
	}

	// 创建模板
	createResp, err := service.CreateSMSTemplate("您的验证码是{code}，请在{time}分钟内使用。", 1)
	if err != nil {
		log.Printf("创建模板失败: %v", err)
	} else {
		fmt.Printf("创建模板成功! 状态: %s\n", createResp.Status)
	}

	// 示例6: 查询报告和日志
	fmt.Println("\n=== 查询分析报告 ===")
	reportsReq := &submail.SMSReportsRequest{
		Project:   "test",
		StartDate: "2024-01-01 00:00:00",
		EndDate:   "2024-12-31 23:59:59",
	}

	reports, err := service.GetSMSReports(reportsReq)
	if err != nil {
		log.Printf("获取分析报告失败: %v", err)
	} else {
		fmt.Printf("分析报告获取成功! 状态: %s\n", reports.Status)
	}

	fmt.Println("\n=== 查询历史明细 ===")
	logReq := &submail.SMSLogRequest{
		Project: "test",
		Limit:   10,
		Offset:  0,
	}

	logs, err := service.GetSMSLog(logReq)
	if err != nil {
		log.Printf("获取历史明细失败: %v", err)
	} else {
		fmt.Printf("历史明细获取成功! 状态: %s\n", logs.Status)
	}

	// 示例7: 获取服务器状态和时间戳
	fmt.Println("\n=== 服务器信息 ===")

	timestamp, err := service.GetTimestamp()
	if err != nil {
		log.Printf("获取时间戳失败: %v", err)
	} else {
		fmt.Printf("服务器时间戳: %s\n", timestamp.Status)
	}

	status, err := service.GetStatus()
	if err != nil {
		log.Printf("获取服务器状态失败: %v", err)
	} else {
		fmt.Printf("服务器状态: %s\n", status.Status)
	}

	// 示例8: 响应格式设置
	fmt.Println("\n=== 响应格式设置 ===")
	fmt.Printf("当前响应格式: %s\n", service.GetFormat())

	// 切换到XML格式
	service.SetFormat(submail.FormatXML)
	fmt.Printf("切换后的响应格式: %s\n", service.GetFormat())

	// 切换回JSON格式
	service.SetFormat(submail.FormatJSON)
	fmt.Printf("切换回JSON格式: %s\n", service.GetFormat())

	// 示例9: 签名调试（开发调试用）
	fmt.Println("\n=== 签名调试 ===")
	debugParams := url.Values{}
	debugParams.Set("to", "13800138000")
	debugParams.Set("text", "测试签名")

	signString, signature := service.ValidateSignature(debugParams)
	fmt.Printf("签名字符串: %s\n", signString)
	fmt.Printf("计算签名: %s\n", signature)

	// 示例10: 时间同步
	fmt.Println("\n=== 时间同步 ===")
	serverTime, err := service.SyncServerTime()
	if err != nil {
		log.Printf("同步服务器时间失败: %v", err)
	} else {
		fmt.Printf("服务器时间戳: %d\n", serverTime)
	}

	offset, err := service.GetTimeOffset()
	if err != nil {
		log.Printf("获取时间偏移量失败: %v", err)
	} else {
		if offset > 0 {
			fmt.Printf("本地时间快于服务器时间 %d 秒\n", offset)
		} else if offset < 0 {
			fmt.Printf("本地时间慢于服务器时间 %d 秒\n", -offset)
		} else {
			fmt.Println("本地时间与服务器时间同步")
		}
	}

	// 示例10: 连接池管理
	fmt.Println("\n=== 连接池管理示例 ===")

	// 获取当前连接池配置
	currentConfig := service.GetConnectionPoolConfig()
	fmt.Printf("当前最大空闲连接数: %d\n", currentConfig.MaxIdleConns)
	fmt.Printf("当前每主机最大空闲连接数: %d\n", currentConfig.MaxIdleConnsPerHost)
	fmt.Printf("当前请求超时时间: %v\n", currentConfig.RequestTimeout)

	// 动态更新连接池配置
	service.UpdateConnectionPoolConfig(func(config *submail.ConnectionPoolConfig) {
		config.MaxIdleConns = 150                // 调整最大空闲连接数
		config.RequestTimeout = 60 * time.Second // 调整请求超时时间
	})
	fmt.Println("连接池配置已更新")

	// 关闭空闲连接（在应用关闭前或需要释放资源时调用）
	service.CloseIdleConnections()
	fmt.Println("空闲连接已关闭")

	// 示例11: 自定义API地址（如果使用私有部署）
	fmt.Println("\n=== 使用自定义API地址 ===")
	service.SetBaseURL("https://your-custom-api.example.com")
	fmt.Println("API地址已设置为自定义地址")

	// 示例12: 短信批量群发
	fmt.Println("\n=== 短信批量群发 ===")
	batchReq := &submail.SMSBatchSendRequest{
		To:      []string{"13800138000", "13800138001", "13800138002"},
		Text:    "您的验证码是：123456，请勿泄露给他人。",
		Project: "verification",
		Tag:     "batch_send",
	}
	batchResp, err := service.SendSMSBatch(batchReq)
	if err != nil {
		log.Printf("短信批量群发失败: %v", err)
	} else {
		fmt.Printf("批量群发结果: %+v\n", batchResp)
	}

	// 示例13: 短信批量模板群发
	fmt.Println("\n=== 短信批量模板群发 ===")
	batchTemplateReq := &submail.SMSBatchXSendRequest{
		To:      []string{"13800138000", "13800138001", "13800138002"},
		Project: "verification_template",
		Vars:    map[string]string{"code": "123456"},
		Tag:     "batch_template",
	}
	batchTemplateResp, err := service.SendSMSBatchTemplate(batchTemplateReq)
	if err != nil {
		log.Printf("短信批量模板群发失败: %v", err)
	} else {
		fmt.Printf("批量模板群发结果: %+v\n", batchTemplateResp)
	}

	// 示例14: 国内短信与国际短信联合发送
	fmt.Println("\n=== 国内短信与国际短信联合发送 ===")
	unionReq := &submail.SMSUnionSendRequest{
		To:      "13800138000",
		Text:    "您的验证码是：123456，请勿泄露给他人。",
		Project: "verification",
		Tag:     "union_send",
		Country: "US", // 美国
	}
	unionResp, err := service.SendSMSUnion(unionReq)
	if err != nil {
		log.Printf("联合发送失败: %v", err)
	} else {
		fmt.Printf("联合发送结果: %+v\n", unionResp)
	}

	// 示例15: 短信订阅管理
	fmt.Println("\n=== 短信订阅管理 ===")
	subscribeResp, err := service.SubscribeSMS("13800138000", "newsletter")
	if err != nil {
		log.Printf("短信订阅失败: %v", err)
	} else {
		fmt.Printf("订阅结果: %+v\n", subscribeResp)
	}

	unsubscribeResp, err := service.UnsubscribeSMS("13800138000", "newsletter")
	if err != nil {
		log.Printf("短信退订失败: %v", err)
	} else {
		fmt.Printf("退订结果: %+v\n", unsubscribeResp)
	}
}
