package main

import (
	"fmt"
	"log"
	"time"

	"github.com/zhoudm1743/submail"
)

func main() {
	// 创建客户端配置
	config := submail.Config{
		AppID:          "your-app-id",          // 替换为您的App ID
		AppKey:         "your-app-key",         // 替换为您的App Key
		BaseURL:        submail.DefaultBaseURL, // 使用默认API地址
		Format:         submail.FormatJSON,     // 使用JSON格式
		UseDigitalSign: false,                  // 使用明文模式（推荐测试时使用）
		SignType:       submail.SignTypeMD5,    // MD5签名（数字签名模式时使用）
		Timeout:        30 * time.Second,       // 30秒超时
	}

	// 创建客户端
	client := submail.NewClient(config)

	// 诊断网络连接（如果遇到连接问题可以使用）
	if err := client.DiagnoseConnection(); err != nil {
		log.Printf("网络诊断失败，请检查网络连接和配置: %v", err)
		return
	}

	// 示例1: 获取服务器时间戳
	fmt.Println("=== 获取服务器时间戳 ===")
	timestampResp, err := client.ServiceTimestamp()
	if err != nil {
		log.Printf("获取时间戳失败: %v", err)
	} else {
		serverTime := time.Unix(timestampResp.Timestamp, 0)
		localTime := time.Now()
		timeDiff := serverTime.Sub(localTime)

		fmt.Printf("服务器时间戳: %d\n", timestampResp.Timestamp)
		fmt.Printf("服务器时间: %s\n", serverTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("本地时间: %s\n", localTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("时间差: %v\n", timeDiff)

		// 便捷方法示例
		currentTimestamp, err := client.GetCurrentTimestamp()
		if err != nil {
			log.Printf("获取当前时间戳失败: %v", err)
		} else {
			fmt.Printf("当前时间戳（便捷方法）: %d\n", currentTimestamp)
		}
	}

	// 示例1.1: 获取服务器状态
	fmt.Println("\n=== 获取服务器状态 ===")
	statusResp, err := client.ServiceStatus()
	if err != nil {
		log.Printf("获取服务状态失败: %v", err)
	} else {
		fmt.Printf("服务状态: %s\n", statusResp.Status)
		fmt.Printf("响应时间: %.3f 秒\n", statusResp.Runtime)

		// 便捷方法示例
		isRunning, err := client.IsServiceRunning()
		if err != nil {
			log.Printf("检查服务状态失败: %v", err)
		} else {
			fmt.Printf("服务是否正常运行: %t\n", isRunning)
		}

		runtime, err := client.GetServiceRuntime()
		if err != nil {
			log.Printf("获取响应时间失败: %v", err)
		} else {
			fmt.Printf("服务响应时间（便捷方法）: %.3f 秒\n", runtime)

			// 使用数据处理方法
			fmt.Printf("性能等级: %s\n", statusResp.GetPerformanceLevel())
			fmt.Printf("是否健康: %t\n", statusResp.IsHealthy())
			fmt.Printf("状态描述: %s\n", statusResp.GetStatusDescription())

			if runtime > 1.0 {
				fmt.Println("⚠️  服务响应时间较慢，可能存在网络延迟")
			} else if runtime < 0.1 {
				fmt.Println("✅ 服务响应时间良好")
			}
		}
	}

	// 示例2: 查询短信余额
	fmt.Println("\n=== 查询短信余额 ===")
	balanceResp, err := client.SMSBalance()
	if err != nil {
		log.Printf("查询余额失败: %v", err)
	} else {
		fmt.Printf("通用类短信余额: %s\n", balanceResp.Balance)
		fmt.Printf("事务类短信余额: %s\n", balanceResp.TransactionalBalance)
	}

	// 示例2.1: 查询短信余额日志
	fmt.Println("\n=== 查询短信余额日志 ===")
	balanceLogResp, err := client.SMSBalanceLogLast7Days()
	if err != nil {
		log.Printf("查询余额日志失败: %v", err)
	} else {
		fmt.Printf("查询到 %d 条余额变更记录\n", len(balanceLogResp.Data))

		// 统计总变更
		transactionalTotal, marketingTotal := balanceLogResp.GetTotalChanges()
		fmt.Printf("最近7天总变更 - 事务类: %d, 运营类: %d\n", transactionalTotal, marketingTotal)

		// 按类型分组统计
		changesByType := balanceLogResp.GetChangesByType()
		fmt.Printf("变更记录分类 - 事务类: %d条, 运营类: %d条\n",
			len(changesByType["transactional"]), len(changesByType["marketing"]))

		// 显示前几条记录
		for i, entry := range balanceLogResp.Data {
			if i >= 3 { // 只显示前3条
				break
			}

			changeTime, _ := entry.ParseDateTime()
			fmt.Printf("  记录%d - 时间: %s, 说明: %s\n",
				i+1, changeTime.Format("2006-01-02 15:04:05"), entry.Message)

			if entry.IsTransactionalSMSChange() {
				transactional, _ := entry.GetChangeAmount()
				fmt.Printf("           事务类变更: %d (变更前: %s, 变更后: %s)\n",
					transactional, entry.TMessageBeforeCredits, entry.TMessageAfterCredits)
			}

			if entry.IsMarketingSMSChange() {
				_, marketing := entry.GetChangeAmount()
				fmt.Printf("           运营类变更: %d (变更前: %s, 变更后: %s)\n",
					marketing, entry.MessageBeforeCredits, entry.MessageAfterCredits)
			}
		}
	}

	// 示例3: 发送短信（需要替换为真实的手机号码）
	fmt.Println("\n=== 发送短信 ===")
	smsReq := &submail.SMSSendRequest{
		To:      "13800138000", // 替换为真实的手机号码
		Content: "【测试签名】这是一条测试短信，验证码：123456。如非本人操作，请忽略。",
		Tag:     "test",
	}

	sendResp, err := client.SMSSend(smsReq)
	if err != nil {
		log.Printf("发送短信失败: %v", err)
	} else {
		fmt.Printf("短信发送成功 - SendID: %s, 费用: %d, 短信条数: %d\n",
			sendResp.SendID, sendResp.Fee, sendResp.Sms)
	}

	// 示例3.1: 使用变量发送短信
	fmt.Println("\n=== 使用变量发送短信 ===")
	// 演示文本变量和日期变量的使用
	contentWithVars := "【测试签名】尊敬的@var(name)，您的验证码是@var(code)，请在@var(expire)分钟内输入。发送时间：@date(Y)年@date(m)月@date(d)日 @date(h):@date(i):@date(s)"

	// 提取变量名
	varNames := client.ExtractVariableNames(contentWithVars)
	fmt.Printf("提取到的变量名: %v\n", varNames)

	// 验证变量格式
	if errors := client.ValidateVariables(contentWithVars); len(errors) > 0 {
		log.Printf("变量格式错误: %v", errors)
	} else {
		fmt.Println("变量格式验证通过")
	}

	// 设置变量值
	vars := map[string]string{
		"name":   "张三",
		"code":   "888888",
		"expire": "10",
	}

	// 使用带变量的发送方法
	varResp, err := client.SMSSendWithVariables("13800138000", contentWithVars, vars, "var-test")
	if err != nil {
		log.Printf("带变量短信发送失败: %v", err)
	} else {
		fmt.Printf("带变量短信发送成功 - SendID: %s, 费用: %d, 短信条数: %d\n",
			varResp.SendID, varResp.Fee, varResp.Sms)
	}

	// 演示处理后的内容
	processedContent := client.ProcessVariables(contentWithVars, vars)
	fmt.Printf("处理后的短信内容: %s\n", processedContent)

	// 示例4: 使用短信模板发送
	fmt.Println("\n=== 使用短信模板发送 ===")
	templateReq := &submail.SMSXSendRequest{
		To:      "13800138000",      // 替换为真实的手机号码
		Project: "your-template-id", // 替换为您的模板ID
		Vars: map[string]string{
			"code": "123456",
			"time": "5",
		},
		Tag: "template-test",
	}

	templateResp, err := client.SMSXSend(templateReq)
	if err != nil {
		log.Printf("模板短信发送失败: %v", err)
	} else {
		fmt.Printf("模板短信发送成功 - SendID: %s, 费用: %d, 短信条数: %d\n",
			templateResp.SendID, templateResp.Fee, templateResp.Sms)
	}

	// 示例4.1: 使用自定义签名发送模板短信（v4.002新功能）
	fmt.Println("\n=== 使用自定义签名发送模板短信 ===")
	customSignatureResp, err := client.SMSXSendWithSignature(
		"13800138000",      // 收件人
		"your-template-id", // 模板ID
		"自定义签名",            // 自定义签名
		map[string]string{ // 变量
			"code": "888888",
			"time": "10",
		},
		"custom-signature-test", // 标签
	)
	if err != nil {
		log.Printf("自定义签名模板短信发送失败: %v", err)
	} else {
		fmt.Printf("自定义签名模板短信发送成功 - SendID: %s, 费用: %d, 短信条数: %d\n",
			customSignatureResp.SendID, customSignatureResp.Fee, customSignatureResp.Sms)
	}

	// 示例5: 一对多发送短信
	fmt.Println("\n=== 一对多发送短信 ===")
	multiReq := &submail.SMSMultiSendRequest{
		Content: "【测试签名】您好，@var(name)，您的取货码为 @var(code)，请在@var(expire)分钟内使用。",
		Multi: []submail.SMSMultiItem{
			{
				To: "13800138000", // 替换为真实的手机号码
				Vars: map[string]string{
					"name":   "张三",
					"code":   "A001",
					"expire": "30",
				},
			},
			{
				To: "13800138001", // 替换为真实的手机号码
				Vars: map[string]string{
					"name":   "李四",
					"code":   "B002",
					"expire": "60",
				},
			},
			{
				To: "13800138002", // 替换为真实的手机号码
				Vars: map[string]string{
					"name":   "王五",
					"code":   "C003",
					"expire": "45",
				},
			},
		},
		Tag: "multi-test",
	}

	multiResp, err := client.SMSMultiSend(multiReq)
	if err != nil {
		log.Printf("一对多短信发送失败: %v", err)
	} else {
		// 处理多条发送结果
		success, failed, totalFee := multiResp.GetStatistics()
		fmt.Printf("一对多短信发送完成 - 成功: %d条, 失败: %d条, 总费用: %d\n", success, failed, totalFee)

		// 显示成功的结果
		successResults := multiResp.GetSuccessResults()
		for i, result := range successResults {
			if i >= 3 { // 只显示前3条
				break
			}
			fmt.Printf("  成功%d - 收件人: %s, SendID: %s, 费用: %d\n",
				i+1, result.To, result.SendID, result.Fee)
		}

		// 显示失败的结果
		failedResults := multiResp.GetFailedResults()
		for i, result := range failedResults {
			if i >= 3 { // 只显示前3条
				break
			}
			fmt.Printf("  失败%d - 收件人: %s, 错误码: %d, 错误信息: %s\n",
				i+1, result.To, result.Code, result.Msg)
		}
	}

	// 示例5.1: 使用便捷方法发送一对多短信
	fmt.Println("\n=== 使用便捷方法发送一对多短信 ===")
	recipients := []submail.SMSMultiItem{
		{
			To: "13800138003",
			Vars: map[string]string{
				"name":   "赵六",
				"code":   "D004",
				"expire": "20",
			},
		},
		{
			To: "13800138004",
			Vars: map[string]string{
				"name":   "钱七",
				"code":   "E005",
				"expire": "25",
			},
		},
	}

	convenienceResp, err := client.SMSMultiSendWithVariables(
		"【测试签名】亲爱的@var(name)，您的验证码是@var(code)，@var(expire)分钟内有效。",
		recipients,
		"convenience-test",
	)
	if err != nil {
		log.Printf("便捷方法一对多短信发送失败: %v", err)
	} else {
		success, failed, totalFee := convenienceResp.GetStatistics()
		fmt.Printf("便捷方法发送完成 - 成功: %d条, 失败: %d条, 总费用: %d\n", success, failed, totalFee)
	}

	// 示例5.2: 模板一对多发送
	fmt.Println("\n=== 模板一对多发送 ===")
	templateMultiReq := &submail.SMSMultiXSendRequest{
		Project: "your-template-id", // 替换为您的模板ID
		Multi: []submail.SMSMultiXItem{
			{
				To: "13800138005",
				Vars: map[string]string{
					"name": "用户A",
					"code": "123456",
					"time": "5",
				},
			},
			{
				To: "13800138006",
				Vars: map[string]string{
					"name": "用户B",
					"code": "654321",
					"time": "10",
				},
				SMSSignature: "特殊签名", // 单独为这个用户设置签名
			},
			{
				To: "13800138007",
				Vars: map[string]string{
					"name": "用户C",
					"code": "789012",
					"time": "3",
				},
			},
		},
		SMSSignature: "通用签名", // 全局签名，没有单独设置签名的用户会使用这个
		Tag:          "template-multi-test",
	}

	templateMultiResp, err := client.SMSMultiXSend(templateMultiReq)
	if err != nil {
		log.Printf("模板一对多短信发送失败: %v", err)
	} else {
		success, failed, totalFee := templateMultiResp.GetStatistics()
		fmt.Printf("模板一对多发送完成 - 成功: %d条, 失败: %d条, 总费用: %d\n", success, failed, totalFee)

		// 显示详细结果
		for i, result := range *templateMultiResp {
			if i >= 3 { // 只显示前3条
				break
			}
			if result.Status == "success" {
				fmt.Printf("  结果%d - 收件人: %s, SendID: %s, 费用: %d\n",
					i+1, result.To, result.SendID, result.Fee)
			} else {
				fmt.Printf("  结果%d - 收件人: %s, 错误: %s\n",
					i+1, result.To, result.Msg)
			}
		}
	}

	// 示例5.3: 使用便捷方法发送带自定义签名的模板一对多短信
	fmt.Println("\n=== 使用自定义签名的模板一对多发送 ===")
	customMultiRecipients := []submail.SMSMultiXItem{
		{
			To: "13800138008",
			Vars: map[string]string{
				"name": "VIP用户",
				"code": "VIP001",
				"time": "15",
			},
		},
		{
			To: "13800138009",
			Vars: map[string]string{
				"name": "普通用户",
				"code": "REG002",
				"time": "5",
			},
		},
	}

	customMultiResp, err := client.SMSMultiXSendWithSignature(
		"your-template-id", // 模板ID
		"自定义签名",            // 自定义签名
		customMultiRecipients,
		"custom-multi-test",
	)
	if err != nil {
		log.Printf("自定义签名模板一对多发送失败: %v", err)
	} else {
		success, failed, totalFee := customMultiResp.GetStatistics()
		fmt.Printf("自定义签名模板一对多发送完成 - 成功: %d条, 失败: %d条, 总费用: %d\n", success, failed, totalFee)
	}

	// 示例5.4: 批量群发短信
	fmt.Println("\n=== 批量群发短信 ===")
	phones := []string{"13800138010", "13800138011", "13800138012", "13800138013", "13800138014"}
	batchContent := "【测试签名】尊敬的用户，您的账户余额变动通知：当前余额@date(Y)年@date(m)月@date(d)日更新。"

	batchResp, err := client.SMSBatchSendWithPhones(batchContent, phones, "batch-test")
	if err != nil {
		log.Printf("批量群发短信失败: %v", err)
	} else {
		success, failed, totalFee := batchResp.GetStatistics()
		fmt.Printf("批量群发完成 - 任务ID: %s, 成功: %d条, 失败: %d条, 总费用: %d\n",
			batchResp.BatchList, success, failed, totalFee)

		// 显示部分结果
		successResults := batchResp.GetSuccessResults()
		for i, result := range successResults {
			if i >= 3 { // 只显示前3条成功的
				break
			}
			fmt.Printf("  成功%d - 收件人: %s, SendID: %s, 费用: %d\n",
				i+1, result.To, result.SendID, result.Fee)
		}

		failedResults := batchResp.GetFailedResults()
		for i, result := range failedResults {
			if i >= 3 { // 只显示前3条失败的
				break
			}
			fmt.Printf("  失败%d - 收件人: %s, 错误: %s\n",
				i+1, result.To, result.Msg)
		}
	}

	// 示例5.5: 批量模板群发短信
	fmt.Println("\n=== 批量模板群发短信 ===")
	templatePhones := []string{"13800138015", "13800138016", "13800138017"}
	templateVars := map[string]string{
		"name":   "批量用户",
		"code":   "BATCH001",
		"expire": "30",
	}

	batchTemplateResp, err := client.SMSBatchXSendWithPhones(
		"your-template-id", // 模板ID
		templatePhones,
		templateVars,
		"批量签名", // 自定义签名
		"batch-template-test",
	)
	if err != nil {
		log.Printf("批量模板群发失败: %v", err)
	} else {
		success, failed, totalFee := batchTemplateResp.GetStatistics()
		fmt.Printf("批量模板群发完成 - 任务ID: %s, 成功: %d条, 失败: %d条, 总费用: %d\n",
			batchTemplateResp.BatchList, success, failed, totalFee)
	}

	// 示例5.6: 国内外短信联合发送
	fmt.Println("\n=== 国内外短信联合发送 ===")

	// 测试号码判断功能
	testNumbers := []string{"13800138000", "+8613800138000", "+1234567890", "+852987654321"}
	for _, number := range testNumbers {
		isInternational := submail.IsInternationalNumber(number)
		fmt.Printf("号码 %s: %s\n", number, map[bool]string{true: "国际号码", false: "国内号码"}[isInternational])
	}

	// 发送国内号码（会使用国内短信通道）
	domesticResp, err := client.SMSUnionSendWithConfig(
		"13800138000", // 国内号码
		"【测试签名】您的验证码是1234，请在10分钟内输入。", // 国内短信内容
		"your-international-app-id",              // 国际短信AppID
		"your-international-app-key",             // 国际短信AppKey
		"[Test] Your verify code is: @var(code)", // 国际短信内容模板
		"union-domestic-test",                    // 标签
		true,                                     // 启用验证码提取转换
	)
	if err != nil {
		log.Printf("国内联合发送失败: %v", err)
	} else {
		fmt.Printf("国内联合发送成功 - SendID: %s, 费用: %d\n", domesticResp.SendID, domesticResp.Fee)
	}

	// 发送国际号码（会使用国际短信通道）
	internationalResp, err := client.SMSUnionSendWithConfig(
		"+852987654321", // 国际号码（香港）
		"【测试签名】您的验证码是5678，请在10分钟内输入。", // 国内短信内容（作为验证码提取源）
		"your-international-app-id",              // 国际短信AppID
		"your-international-app-key",             // 国际短信AppKey
		"[Test] Your verify code is: @var(code)", // 国际短信内容模板
		"union-international-test",               // 标签
		true,                                     // 启用验证码提取转换
	)
	if err != nil {
		log.Printf("国际联合发送失败: %v", err)
	} else {
		fmt.Printf("国际联合发送成功 - SendID: %s, 费用: %d\n", internationalResp.SendID, internationalResp.Fee)
	}

	// 示例6: 查询短信发送历史
	fmt.Println("\n=== 查询短信发送历史 ===")
	logResp, err := client.SMSLogLast7Days()
	if err != nil {
		log.Printf("查询短信历史失败: %v", err)
	} else {
		fmt.Printf("查询到 %d 条短信记录（总共 %d 条）\n", len(logResp.Data), logResp.Total)

		// 统计信息
		success, failed, pending, totalFee := logResp.GetLogStatistics()
		fmt.Printf("统计信息 - 成功: %d, 失败: %d, 未知: %d, 总费用: %d\n",
			success, failed, pending, totalFee)

		// 按运营商分组统计
		operatorStats := logResp.GetLogsByOperator()
		fmt.Printf("运营商分布: ")
		for operator, logs := range operatorStats {
			fmt.Printf("%s(%d) ", operator, len(logs))
		}
		fmt.Println()

		// 失败原因统计
		failureReasons := logResp.GetFailureReasons()
		if len(failureReasons) > 0 {
			fmt.Printf("失败原因统计: ")
			for reason, count := range failureReasons {
				fmt.Printf("%s(%d) ", reason, count)
			}
			fmt.Println()
		}

		// 显示前几条记录
		for i, logEntry := range logResp.Data {
			if i >= 3 { // 只显示前3条
				break
			}
			sendTime := logEntry.GetSendTime()
			duration := logEntry.GetDeliveryDuration()
			fmt.Printf("  记录%d - 手机号: %s, 状态: %s, 发送时间: %s",
				i+1, logEntry.To, logEntry.Status, sendTime.Format("2006-01-02 15:04:05"))
			if duration > 0 {
				fmt.Printf(", 送达耗时: %v", duration)
			}
			if logEntry.IsDropped() && logEntry.DroppedReason != "" {
				fmt.Printf(", 失败原因: %s", logEntry.DroppedReason)
			}
			fmt.Println()
		}
	}

	// 示例6.1: 根据手机号查询短信历史
	fmt.Println("\n=== 根据手机号查询短信历史 ===")
	phoneLogResp, err := client.SMSLogByPhone("13800138000") // 替换为真实手机号
	if err != nil {
		log.Printf("根据手机号查询失败: %v", err)
	} else {
		fmt.Printf("手机号 13800138000 的短信记录: %d 条\n", len(phoneLogResp.Data))
	}

	// 示例6.2: 查询发送失败的短信
	fmt.Println("\n=== 查询发送失败的短信 ===")
	failedLogResp, err := client.SMSLogByStatus("dropped")
	if err != nil {
		log.Printf("查询失败短信失败: %v", err)
	} else {
		fmt.Printf("发送失败的短信记录: %d 条\n", len(failedLogResp.Data))
		failedLogs := failedLogResp.GetFailedLogs()
		for i, log := range failedLogs {
			if i >= 2 { // 只显示前2条
				break
			}
			fmt.Printf("  失败记录%d - 手机号: %s, 原因: %s, 运营商状态: %s\n",
				i+1, log.To, log.DroppedReason, log.ReportState)
		}
	}

	// 示例7: 获取短信分析报告
	fmt.Println("\n=== 获取短信分析报告 ===")

	// 使用便捷方法获取最近7天的报告
	reportResp, err := client.SMSReportsLast7Days()
	if err != nil {
		log.Printf("获取分析报告失败: %v", err)
	} else {
		fmt.Printf("分析报告时间范围: %s 到 %s\n", reportResp.StartDate, reportResp.EndDate)

		// 显示概览数据
		overview := reportResp.Overview
		fmt.Printf("总体概览: 请求 %d 次, 成功 %d 次, 失败 %d 次, 计费 %d 条\n",
			overview.Request, overview.Deliveryed, overview.Dropped, overview.Fee)
		fmt.Printf("成功率: %.2f%%, 失败率: %.2f%%\n",
			overview.GetSuccessRate(), overview.GetFailureRate())

		// 显示运营商占比
		if overview.Operators.GetTotalOperators() > 0 {
			fmt.Println("运营商分布:")
			percentages := overview.Operators.GetOperatorPercentage()
			fmt.Printf("  移动: %d (%.1f%%), 联通: %d (%.1f%%), 电信: %d (%.1f%%)\n",
				overview.Operators.ChinaMobile, percentages["移动"],
				overview.Operators.ChinaUnicom, percentages["联通"],
				overview.Operators.ChinaTelecom, percentages["电信"])
		}

		// 显示主要省份分布
		if len(overview.Location.Province) > 0 {
			fmt.Println("发送量最多的省份:")
			topProvinces := overview.Location.GetTopProvinces(5)
			for i, province := range topProvinces {
				fmt.Printf("  %d. %s: %d 条\n", i+1, province.Province, province.Count)
			}
		}

		// 显示主要失败原因
		if len(overview.DroppedReasonAnalysis) > 0 {
			fmt.Println("主要失败原因:")
			topReasons := overview.GetTopFailureReasons(3)
			for i, reason := range topReasons {
				fmt.Printf("  %d. %s: %d 次\n", i+1, reason.Reason, reason.Count)
			}
		}

		// 显示时间线数据（前3天）
		if len(reportResp.Timeline) > 0 {
			fmt.Println("时间线数据:")
			for i, timeline := range reportResp.Timeline {
				if i >= 3 { // 只显示前3天
					break
				}
				report := timeline.Report
				fmt.Printf("  %s - 请求: %d, 成功: %d, 失败: %d, 计费: %d\n",
					timeline.Date, report.Request, report.Deliveryed, report.Dropped, report.Fee)
			}
		}
	}

	// 示例8: 日期变量演示
	fmt.Println("\n=== 日期变量演示 ===")
	fmt.Println("支持的日期变量:")
	dateDescriptions := client.GetDateVariableDescription()
	for variable, description := range dateDescriptions {
		fmt.Printf("  %s - %s\n", variable, description)
	}

	// 演示日期变量的实际效果
	dateContent := "【测试签名】当前时间：@date()，年份：@date(Y)，月份：@date(m)，日期：@date(d)，星期：@date(l)"
	processedDateContent := client.ProcessVariables(dateContent, nil)
	fmt.Printf("\n日期变量示例:\n原始内容: %s\n处理后内容: %s\n", dateContent, processedDateContent)

	// 示例9: 短信签名管理
	fmt.Println("\n=== 短信签名管理 ===")

	// 查询现有签名
	fmt.Println("查询现有签名:")
	queryReq := &submail.SMSSignatureQueryRequest{
		// 可以指定查询特定签名
		// SMSSignature: "测试签名",
	}

	queryResp, err := client.SMSSignatureQuery(queryReq)
	if err != nil {
		log.Printf("查询签名失败: %v", err)
	} else {
		fmt.Printf("查询到 %d 个签名:\n", len(queryResp.SMSSignatures))
		for i, sig := range queryResp.SMSSignatures {
			if i >= 3 { // 只显示前3个
				break
			}
			status := submail.GetSignatureStatus(sig.Status)
			fmt.Printf("  签名%d - %s (状态: %s)\n", i+1, sig.SMSSignature, status)
		}
	}

	// 创建新签名（需要提供真实的企业信息和证明材料）
	fmt.Println("\n创建新签名（示例，需要真实信息）:")

	// 注意：实际使用时，您需要从HTTP请求中获取multipart.FileHeader
	// 这里只是展示结构体的定义，实际不能执行
	createReq := &submail.SMSSignatureCreateRequest{
		SMSSignature:       "新测试签名", // 可省略【】符号
		Company:            "测试公司有限公司",
		CompanyLisenceCode: "91000000000000000X",
		LegalName:          "张三",
		// Attachments:        fileHeaders,  // 从HTTP表单获取的[]multipart.FileHeader
		AgentName:  "李四",
		AgentID:    "110000000000000000",
		AgentMob:   "13800138000",
		SourceType: 0, // 0=营业执照、1=商标、2=APP
		Contact:    "13800138000",
	}

	// 实际使用示例（在Web服务器中）:
	// r.ParseMultipartForm(32 << 20) // 32MB
	// form := r.MultipartForm
	// files := form.File["attachments"]
	// createReq.Attachments = files

	// 注意：这里只是演示结构体，实际使用时需要提供真实的企业信息和文件
	fmt.Printf("签名创建请求: %s (公司: %s, 材料类型: %s)\n",
		createReq.SMSSignature,
		createReq.Company,
		submail.GetSourceTypeDescription(createReq.SourceType))

	// createResp, err := client.SMSSignatureCreate(createReq)
	// if err != nil {
	//     log.Printf("创建签名失败: %v", err)
	// } else {
	//     fmt.Printf("签名创建成功，状态: %s\n", createResp.Status)
	// }

	// 示例10: 短信模板管理
	fmt.Println("\n=== 短信模板管理 ===")

	// 查询模板列表
	fmt.Println("查询模板列表:")
	getReq := &submail.SMSTemplateGetRequest{
		// TemplateID: "specific-template-id", // 查询特定模板
		Offset: 0, // 数据偏移量
	}

	getResp, err := client.SMSTemplateGet(getReq)
	if err != nil {
		log.Printf("查询模板失败: %v", err)
	} else {
		fmt.Printf("查询到 %d 个模板 (第%d-%d行):\n", len(getResp.Templates), getResp.StartRow, getResp.EndRow)
		for i, template := range getResp.Templates {
			if i >= 3 { // 只显示前3个
				break
			}
			status := submail.GetTemplateStatus(template.TemplateStatus)
			addTime := submail.GetTemplateAddTime(template.AddDate)
			fmt.Printf("  模板%d - ID: %s, 标题: %s, 状态: %s, 创建时间: %s\n",
				i+1, template.TemplateID, template.SMSTitle, status, addTime.Format("2006-01-02 15:04:05"))
			fmt.Printf("           签名: %s, 内容: %s\n", template.SMSSignature, template.SMSContent)
			if template.TemplateRejectReason != "" {
				fmt.Printf("           驳回原因: %s\n", template.TemplateRejectReason)
			}
		}
	}

	// 创建新模板（示例）
	fmt.Println("\n创建新模板（示例）:")
	createTemplateReq := &submail.SMSTemplateCreateRequest{
		SMSTitle:     "验证码模板",
		SMSSignature: "【测试公司】",
		SMSContent:   "您的验证码是@var(code)，请在@var(expire)分钟内输入。如非本人操作，请忽略此短信。",
	}

	fmt.Printf("模板创建请求: 标题=%s, 签名=%s\n", createTemplateReq.SMSTitle, createTemplateReq.SMSSignature)
	fmt.Printf("模板内容: %s\n", createTemplateReq.SMSContent)

	// 提取模板中的变量
	templateVarNames := client.ExtractVariableNames(createTemplateReq.SMSContent)
	fmt.Printf("模板中包含的变量: %v\n", templateVarNames)

	// createTemplateResp, err := client.SMSTemplateCreate(createTemplateReq)
	// if err != nil {
	//     log.Printf("创建模板失败: %v", err)
	// } else {
	//     fmt.Printf("模板创建成功，模板ID: %s\n", createTemplateResp.TemplateID)
	// }

	// 更新模板（示例）
	fmt.Println("\n更新模板（示例）:")
	updateTemplateReq := &submail.SMSTemplateUpdateRequest{
		TemplateID:   "existing-template-id", // 需要替换为真实的模板ID
		SMSTitle:     "更新的验证码模板",
		SMSSignature: "【测试公司】",
		SMSContent:   "您的验证码是@var(code)，请在@var(expire)分钟内输入。此验证码仅用于身份验证，请勿泄露。",
	}

	fmt.Printf("模板更新请求: ID=%s, 新标题=%s\n", updateTemplateReq.TemplateID, updateTemplateReq.SMSTitle)

	// updateResp, err := client.SMSTemplateUpdate(updateTemplateReq)
	// if err != nil {
	//     log.Printf("更新模板失败: %v", err)
	// } else {
	//     fmt.Printf("模板更新成功，状态: %s\n", updateResp.Status)
	// }

	fmt.Println("\n=== 示例程序执行完成 ===")
	fmt.Println("注意：以上示例中的手机号码和模板ID需要替换为真实有效的值")
	fmt.Println("建议在测试环境中先使用明文模式（UseDigitalSign: false）进行测试")
	fmt.Println("生产环境建议使用数字签名模式（UseDigitalSign: true）以提高安全性")
	fmt.Println("\n变量功能说明:")
	fmt.Println("- 文本变量格式: @var(变量名)，如 @var(name), @var(code)")
	fmt.Println("- 日期变量格式: @date() 或 @date(格式)，如 @date(Y), @date(m), @date(d)")
	fmt.Println("- 变量会在发送前自动处理和替换")
	fmt.Println("- 支持时区设置: client.SetTimezone(\"Asia/Shanghai\")")

	// 示例10: 查询短信上行回复
	fmt.Println("\n=== 查询短信上行回复 ===")
	moResp, err := client.SMSMOLast7Days()
	if err != nil {
		log.Printf("查询短信上行失败: %v", err)
	} else {
		fmt.Printf("查询到 %d 条上行回复（总共 %d 条）\n", len(moResp.MO), moResp.Total)

		// 统计信息
		total, validReplies, unsubscribes := moResp.GetMOStatistics()
		fmt.Printf("统计信息 - 总数: %d, 有效回复: %d, 退订: %d\n",
			total, validReplies, unsubscribes)

		// 按回复内容分组
		contentStats := moResp.GetMOByContent()
		fmt.Printf("回复内容分布: ")
		for content, mos := range contentStats {
			fmt.Printf("'%s'(%d) ", content, len(mos))
		}
		fmt.Println()

		// 显示前几条记录
		for i, mo := range moResp.MO {
			if i >= 3 { // 只显示前3条
				break
			}
			replyTime := mo.GetReplyTime()
			replyType := "有效回复"
			if mo.IsReturnReceipt() {
				replyType = "退订"
			}
			fmt.Printf("  回复%d - 手机号: %s, 内容: '%s', 类型: %s, 时间: %s\n",
				i+1, mo.From, mo.Content, replyType, replyTime.Format("2006-01-02 15:04:05"))
			if mo.SMSContent != "" {
				fmt.Printf("           原短信: %s\n", mo.SMSContent)
			}
		}
	}

	// 示例10.1: 根据手机号查询上行回复
	fmt.Println("\n=== 根据手机号查询上行回复 ===")
	phoneResp, err := client.SMSMOByPhone("13800138000") // 替换为真实手机号
	if err != nil {
		log.Printf("根据手机号查询上行失败: %v", err)
	} else {
		fmt.Printf("手机号 13800138000 的上行回复: %d 条\n", len(phoneResp.MO))
	}

	// 示例10.2: 分析退订情况
	fmt.Println("\n=== 分析退订情况 ===")
	unsubscribes := moResp.GetUnsubscribes()
	if len(unsubscribes) > 0 {
		fmt.Printf("发现 %d 条退订回复\n", len(unsubscribes))
		for i, unsubscribe := range unsubscribes {
			if i >= 2 { // 只显示前2条
				break
			}
			fmt.Printf("  退订%d - 手机号: %s, 内容: '%s', 时间: %s\n",
				i+1, unsubscribe.From, unsubscribe.Content,
				unsubscribe.GetReplyTime().Format("2006-01-02 15:04:05"))
		}
	} else {
		fmt.Println("没有发现退订回复")
	}
}
