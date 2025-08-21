package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zhoudm1743/submail" // 根据实际路径调整
)

func main() {
	fmt.Println("SUBMAIL SUBHOOK 示例")
	fmt.Println("请根据需要取消注释相应的示例函数")
	// runSubhookExample()
	// exampleHTTPServer()
}

func runSubhookExample() {
	// 创建 SUBMAIL 客户端
	client := submail.NewClient(submail.Config{
		AppID:  "your_app_id",
		AppKey: "your_app_key",
	})

	// ===== 示例1：创建 SUBHOOK =====
	fmt.Println("=== 创建 SUBHOOK 示例 ===")

	// 创建短信相关事件的 SUBHOOK
	createResp, err := client.SubhookCreateForSMS("https://your-domain.com/subhook", "sms_events")
	if err != nil {
		log.Printf("创建 SUBHOOK 失败: %v", err)
	} else {
		fmt.Printf("创建 SUBHOOK 成功:\n")
		fmt.Printf("  SUBHOOK ID: %s\n", createResp.Target)
		fmt.Printf("  密匙: %s\n", createResp.Key)
	}

	// 创建所有事件的 SUBHOOK
	createAllResp, err := client.SubhookCreateForAll("https://your-domain.com/subhook/all", "all_events")
	if err != nil {
		log.Printf("创建全事件 SUBHOOK 失败: %v", err)
	} else {
		fmt.Printf("创建全事件 SUBHOOK 成功:\n")
		fmt.Printf("  SUBHOOK ID: %s\n", createAllResp.Target)
		fmt.Printf("  密匙: %s\n", createAllResp.Key)
	}

	// ===== 示例2：查询 SUBHOOK =====
	fmt.Println("\n=== 查询 SUBHOOK 示例 ===")

	queryResp, err := client.SubhookQueryAll()
	if err != nil {
		log.Printf("查询 SUBHOOK 失败: %v", err)
	} else {
		fmt.Printf("查询到 %d 个 SUBHOOK:\n", len(queryResp.Subhooks))
		for i, subhook := range queryResp.Subhooks {
			fmt.Printf("  %d. ID: %s, URL: %s, 事件: %v\n",
				i+1, subhook.Target, subhook.URL, subhook.Event)
		}
	}

	// ===== 示例3：创建 SUBHOOK HTTP 处理器 =====
	fmt.Println("\n=== 创建 SUBHOOK HTTP 处理器示例 ===")

	// 假设从创建响应中获得的密匙
	subhookKey := "your_subhook_key" // 实际使用时应该从创建响应中获取

	// 创建事件处理器
	eventHandler := &submail.DefaultSubhookEventHandler{
		OnDelivered: func(eventData *submail.SubhookEventData, smsData *submail.SMSSubhookEventData) error {
			fmt.Printf("短信发送成功: SendID=%s, To=%s\n", smsData.SendID, smsData.To)
			return nil
		},
		OnDropped: func(eventData *submail.SubhookEventData, smsData *submail.SMSSubhookEventData) error {
			fmt.Printf("短信发送失败: SendID=%s, To=%s, Status=%s\n",
				smsData.SendID, smsData.To, smsData.Status)
			return nil
		},
		OnMO: func(eventData *submail.SubhookEventData, moData *submail.SMSMOSubhookEventData) error {
			fmt.Printf("收到短信回复: From=%s, Content=%s\n", moData.From, moData.Content)
			return nil
		},
		OnTemplateAccept: func(eventData *submail.SubhookEventData, templateData *submail.TemplateSubhookEventData) error {
			fmt.Printf("模板审核通过: TemplateID=%s\n", templateData.TemplateID)
			return nil
		},
		OnTemplateReject: func(eventData *submail.SubhookEventData, templateData *submail.TemplateSubhookEventData) error {
			fmt.Printf("模板审核未通过: TemplateID=%s, 原因=%s\n",
				templateData.TemplateID, templateData.Reason)
			return nil
		},
	}

	// 创建 HTTP 处理器
	_ = submail.CreateSubhookHTTPHandler(subhookKey, eventHandler) // 示例中未使用，实际使用时移除 _

	// ===== 示例4：启动 HTTP 服务器（可选） =====
	fmt.Println("\n=== HTTP 服务器示例（注释掉以避免阻塞） ===")
	fmt.Println("// 取消注释以下代码来启动 HTTP 服务器:")
	fmt.Println("// http.HandleFunc(\"/subhook\", httpHandler)")
	fmt.Println("// fmt.Println(\"SUBHOOK 服务器启动在 :8080/subhook\")")
	fmt.Println("// log.Fatal(http.ListenAndServe(\":8080\", nil))")

	/*
		// 实际使用时取消注释以下代码
		http.HandleFunc("/subhook", httpHandler)
		fmt.Println("SUBHOOK 服务器启动在 :8080/subhook")
		log.Fatal(http.ListenAndServe(":8080", nil))
	*/

	// ===== 示例5：手动验证签名 =====
	fmt.Println("\n=== 手动验证签名示例 ===")

	// 模拟接收到的 SUBHOOK 数据
	token := "abcd1234567890abcd1234567890abcd"     // 32位随机字符串
	signature := "e10adc3949ba59abbe56e057f20f883e" // 假设的签名
	key := "your_subhook_key"                       // SUBHOOK 密匙

	// 验证签名
	isValid := submail.ValidateSubhookSignature(token, signature, key)
	fmt.Printf("签名验证结果: %v\n", isValid)

	// 演示正确的签名计算
	fmt.Println("\n=== 签名计算演示 ===")
	testToken := "12345678901234567890123456789012"
	testKey := "test_key"
	// 手动计算签名用于测试
	combinedString := testToken + testKey
	fmt.Printf("拼接字符串: %s\n", combinedString)

	// 使用我们的验证函数测试
	testSignature := "d41d8cd98f00b204e9800998ecf8427e" // 这是一个示例签名
	isTestValid := submail.ValidateSubhookSignature(testToken, testSignature, testKey)
	fmt.Printf("测试签名验证: %v\n", isTestValid)

	// ===== 示例6：事件类型验证 =====
	fmt.Println("\n=== 事件类型验证示例 ===")

	validEvents := []string{
		submail.SubhookEventDelivered,
		submail.SubhookEventDropped,
		submail.SubhookEventMO,
	}

	invalidEvents := []string{
		submail.SubhookEventDelivered,
		"invalid_event",
		submail.SubhookEventMO,
	}

	fmt.Printf("验证有效事件类型: %v\n", submail.ValidateEventTypes(validEvents))
	fmt.Printf("验证无效事件类型: %v\n", submail.ValidateEventTypes(invalidEvents))

	// ===== 示例7：获取事件类型描述 =====
	fmt.Println("\n=== 事件类型描述示例 ===")

	eventTypes := []string{
		submail.SubhookEventRequest,
		submail.SubhookEventDelivered,
		submail.SubhookEventDropped,
		submail.SubhookEventSending,
		submail.SubhookEventMO,
		submail.SubhookEventTemplateAccept,
		submail.SubhookEventTemplateReject,
	}

	for _, eventType := range eventTypes {
		fmt.Printf("%s: %s\n", eventType, submail.GetEventTypeDescription(eventType))
	}

	// ===== 示例8：删除 SUBHOOK（可选） =====
	fmt.Println("\n=== 删除 SUBHOOK 示例（注释掉以避免误删） ===")
	fmt.Println("// 取消注释以下代码来删除 SUBHOOK:")
	fmt.Println("// if createResp != nil && createResp.Target != \"\" {")
	fmt.Println("//     deleteResp, err := client.SubhookDeleteByID(createResp.Target)")
	fmt.Println("//     if err != nil {")
	fmt.Println("//         log.Printf(\"删除 SUBHOOK 失败: %v\", err)")
	fmt.Println("//     } else {")
	fmt.Println("//         fmt.Printf(\"删除 SUBHOOK 成功: %s\\n\", deleteResp.Status)")
	fmt.Println("//     }")
	fmt.Println("// }")

	/*
		// 实际使用时根据需要取消注释
		if createResp != nil && createResp.Target != "" {
			deleteResp, err := client.SubhookDeleteByID(createResp.Target)
			if err != nil {
				log.Printf("删除 SUBHOOK 失败: %v", err)
			} else {
				fmt.Printf("删除 SUBHOOK 成功: %s\n", deleteResp.Status)
			}
		}
	*/

	fmt.Println("\n=== SUBHOOK 示例完成 ===")
}

// ===== 自定义事件处理器示例 =====

// CustomSubhookHandler 自定义 SUBHOOK 事件处理器
type CustomSubhookHandler struct {
	// 可以添加自定义字段，如数据库连接、日志记录器等
}

// HandleEvent 实现 SubhookEventHandler 接口
func (h *CustomSubhookHandler) HandleEvent(eventType string, eventData *submail.SubhookEventData) error {
	fmt.Printf("收到事件: %s (%s)\n", eventType, submail.GetEventTypeDescription(eventType))

	switch eventType {
	case submail.SubhookEventDelivered:
		return h.handleDelivered(eventData)
	case submail.SubhookEventDropped:
		return h.handleDropped(eventData)
	case submail.SubhookEventMO:
		return h.handleMO(eventData)
	default:
		fmt.Printf("未处理的事件类型: %s\n", eventType)
	}

	return nil
}

func (h *CustomSubhookHandler) handleDelivered(eventData *submail.SubhookEventData) error {
	smsData, err := submail.ParseSMSSubhookEvent(eventData)
	if err != nil {
		return err
	}

	// 自定义处理逻辑
	fmt.Printf("处理发送成功事件: SendID=%s, To=%s, Fee=%d\n",
		smsData.SendID, smsData.To, smsData.Fee)

	// 这里可以添加数据库记录、日志记录等操作

	return nil
}

func (h *CustomSubhookHandler) handleDropped(eventData *submail.SubhookEventData) error {
	smsData, err := submail.ParseSMSSubhookEvent(eventData)
	if err != nil {
		return err
	}

	// 自定义处理逻辑
	fmt.Printf("处理发送失败事件: SendID=%s, To=%s, Status=%s\n",
		smsData.SendID, smsData.To, smsData.Status)

	// 这里可以添加重试逻辑、告警通知等操作

	return nil
}

func (h *CustomSubhookHandler) handleMO(eventData *submail.SubhookEventData) error {
	moData, err := submail.ParseSMSMOSubhookEvent(eventData)
	if err != nil {
		return err
	}

	// 自定义处理逻辑
	fmt.Printf("处理短信回复事件: From=%s, Content=%s\n",
		moData.From, moData.Content)

	// 这里可以添加自动回复、客服系统集成等操作

	return nil
}

// ===== HTTP 服务器完整示例 =====

func exampleHTTPServer() {
	// 创建客户端
	client := submail.NewClient(submail.Config{
		AppID:  "your_app_id",
		AppKey: "your_app_key",
	})

	// 创建 SUBHOOK（这通常只需要执行一次）
	createResp, err := client.SubhookCreateForAll("https://your-domain.com/subhook", "webhook")
	if err != nil {
		log.Fatalf("创建 SUBHOOK 失败: %v", err)
	}

	fmt.Printf("SUBHOOK 创建成功，密匙: %s\n", createResp.Key)

	// 创建事件处理器
	handler := &CustomSubhookHandler{}

	// 创建 HTTP 处理器
	httpHandler := submail.CreateSubhookHTTPHandler(createResp.Key, handler)

	// 设置路由
	http.HandleFunc("/subhook", httpHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 启动服务器
	fmt.Println("SUBHOOK 服务器启动在 :8080")
	fmt.Println("健康检查: http://localhost:8080/health")
	fmt.Println("SUBHOOK 端点: http://localhost:8080/subhook")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
