package submail

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
)

// SubhookHandler SUBHOOK 处理器
type SubhookHandler struct {
	client *Client
}

// NewSubhookHandler 创建 SUBHOOK 处理器
func NewSubhookHandler(client *Client) *SubhookHandler {
	return &SubhookHandler{
		client: client,
	}
}

// ===== SUBHOOK 数据验证功能 =====

// ValidateSubhookSignature 验证 SUBHOOK 签名
// 参数：
//   - token: 32位随机字符串（从POST数据中获取）
//   - signature: 数字签名（从POST数据中获取）
//   - key: SUBHOOK 密匙（创建SUBHOOK时获得的密匙）
//
// 返回：
//   - bool: 签名是否有效
func ValidateSubhookSignature(token, signature, key string) bool {
	if token == "" || signature == "" || key == "" {
		return false
	}

	// 拼接字符串：token + key
	combinedString := token + key

	// 计算 MD5 签名
	hash := md5.Sum([]byte(combinedString))
	generatedSignature := fmt.Sprintf("%x", hash)

	// 比对签名（不区分大小写）
	return strings.EqualFold(generatedSignature, signature)
}

// ParseSubhookEvent 解析 SUBHOOK 事件数据
// 从 HTTP 请求中解析 SUBHOOK 事件通知
func ParseSubhookEvent(r *http.Request) (*SubhookEventData, error) {
	if r.Method != "POST" {
		return nil, fmt.Errorf("SUBHOOK 事件通知必须使用 POST 方法")
	}

	// 解析表单数据
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("解析表单数据失败: %v", err)
	}

	event := &SubhookEventData{
		Token:     r.FormValue("token"),
		Signature: r.FormValue("signature"),
		Event:     r.FormValue("event"),
		AppID:     r.FormValue("appid"),
	}

	// 解析时间戳
	if timestampStr := r.FormValue("timestamp"); timestampStr != "" {
		var timestamp int64
		if _, err := fmt.Sscanf(timestampStr, "%d", &timestamp); err == nil {
			event.Timestamp = timestamp
		}
	}

	// 解析事件数据（根据事件类型处理不同的数据格式）
	event.Data = make(map[string]interface{})
	for key, values := range r.Form {
		if key != "token" && key != "signature" && key != "event" && key != "appid" && key != "timestamp" {
			if len(values) == 1 {
				event.Data[key] = values[0]
			} else {
				event.Data[key] = values
			}
		}
	}

	return event, nil
}

// ParseSMSSubhookEvent 解析短信相关的 SUBHOOK 事件
func ParseSMSSubhookEvent(eventData *SubhookEventData) (*SMSSubhookEventData, error) {
	if eventData == nil || eventData.Data == nil {
		return nil, fmt.Errorf("事件数据为空")
	}

	smsEvent := &SMSSubhookEventData{}

	// 解析基础字段
	if sendID, ok := eventData.Data["send_id"].(string); ok {
		smsEvent.SendID = sendID
	}
	if to, ok := eventData.Data["to"].(string); ok {
		smsEvent.To = to
	}
	if content, ok := eventData.Data["content"].(string); ok {
		smsEvent.Content = content
	}
	if status, ok := eventData.Data["status"].(string); ok {
		smsEvent.Status = status
	}

	// 解析数值字段
	if feeStr, ok := eventData.Data["fee"].(string); ok {
		var fee int
		if _, err := fmt.Sscanf(feeStr, "%d", &fee); err == nil {
			smsEvent.Fee = fee
		}
	}

	// 解析时间戳字段
	if sendAtStr, ok := eventData.Data["send_at"].(string); ok {
		var sendAt int64
		if _, err := fmt.Sscanf(sendAtStr, "%d", &sendAt); err == nil {
			smsEvent.SendAt = sendAt
		}
	}
	if reportAtStr, ok := eventData.Data["report_at"].(string); ok {
		var reportAt int64
		if _, err := fmt.Sscanf(reportAtStr, "%d", &reportAt); err == nil {
			smsEvent.ReportAt = reportAt
		}
	}

	return smsEvent, nil
}

// ParseSMSMOSubhookEvent 解析短信上行 SUBHOOK 事件
func ParseSMSMOSubhookEvent(eventData *SubhookEventData) (*SMSMOSubhookEventData, error) {
	if eventData == nil || eventData.Data == nil {
		return nil, fmt.Errorf("事件数据为空")
	}

	moEvent := &SMSMOSubhookEventData{}

	// 解析基础字段
	if from, ok := eventData.Data["from"].(string); ok {
		moEvent.From = from
	}
	if content, ok := eventData.Data["content"].(string); ok {
		moEvent.Content = content
	}
	if smsContent, ok := eventData.Data["sms_content"].(string); ok {
		moEvent.SMSContent = smsContent
	}

	// 解析时间戳字段
	if replyAtStr, ok := eventData.Data["reply_at"].(string); ok {
		var replyAt int64
		if _, err := fmt.Sscanf(replyAtStr, "%d", &replyAt); err == nil {
			moEvent.ReplyAt = replyAt
		}
	}

	return moEvent, nil
}

// ParseTemplateSubhookEvent 解析模板审核 SUBHOOK 事件
func ParseTemplateSubhookEvent(eventData *SubhookEventData) (*TemplateSubhookEventData, error) {
	if eventData == nil || eventData.Data == nil {
		return nil, fmt.Errorf("事件数据为空")
	}

	templateEvent := &TemplateSubhookEventData{}

	// 解析基础字段
	if templateID, ok := eventData.Data["template_id"].(string); ok {
		templateEvent.TemplateID = templateID
	}
	if status, ok := eventData.Data["status"].(string); ok {
		templateEvent.Status = status
	}
	if reason, ok := eventData.Data["reason"].(string); ok {
		templateEvent.Reason = reason
	}

	return templateEvent, nil
}

// ===== SUBHOOK 事件处理器接口 =====

// SubhookEventHandler SUBHOOK 事件处理器接口
type SubhookEventHandler interface {
	HandleEvent(eventType string, eventData *SubhookEventData) error
}

// DefaultSubhookEventHandler 默认的 SUBHOOK 事件处理器
type DefaultSubhookEventHandler struct {
	OnRequest        func(*SubhookEventData, *SMSSubhookEventData) error      // 发送请求被接收
	OnDelivered      func(*SubhookEventData, *SMSSubhookEventData) error      // 发送成功
	OnDropped        func(*SubhookEventData, *SMSSubhookEventData) error      // 发送失败
	OnSending        func(*SubhookEventData, *SMSSubhookEventData) error      // 正在发送
	OnMO             func(*SubhookEventData, *SMSMOSubhookEventData) error    // 短信上行
	OnTemplateAccept func(*SubhookEventData, *TemplateSubhookEventData) error // 模板审核通过
	OnTemplateReject func(*SubhookEventData, *TemplateSubhookEventData) error // 模板审核未通过
}

// HandleEvent 处理 SUBHOOK 事件
func (h *DefaultSubhookEventHandler) HandleEvent(eventType string, eventData *SubhookEventData) error {
	switch eventType {
	case SubhookEventRequest:
		if h.OnRequest != nil {
			smsData, err := ParseSMSSubhookEvent(eventData)
			if err != nil {
				return fmt.Errorf("解析短信事件数据失败: %v", err)
			}
			return h.OnRequest(eventData, smsData)
		}

	case SubhookEventDelivered:
		if h.OnDelivered != nil {
			smsData, err := ParseSMSSubhookEvent(eventData)
			if err != nil {
				return fmt.Errorf("解析短信事件数据失败: %v", err)
			}
			return h.OnDelivered(eventData, smsData)
		}

	case SubhookEventDropped:
		if h.OnDropped != nil {
			smsData, err := ParseSMSSubhookEvent(eventData)
			if err != nil {
				return fmt.Errorf("解析短信事件数据失败: %v", err)
			}
			return h.OnDropped(eventData, smsData)
		}

	case SubhookEventSending:
		if h.OnSending != nil {
			smsData, err := ParseSMSSubhookEvent(eventData)
			if err != nil {
				return fmt.Errorf("解析短信事件数据失败: %v", err)
			}
			return h.OnSending(eventData, smsData)
		}

	case SubhookEventMO:
		if h.OnMO != nil {
			moData, err := ParseSMSMOSubhookEvent(eventData)
			if err != nil {
				return fmt.Errorf("解析短信上行事件数据失败: %v", err)
			}
			return h.OnMO(eventData, moData)
		}

	case SubhookEventTemplateAccept:
		if h.OnTemplateAccept != nil {
			templateData, err := ParseTemplateSubhookEvent(eventData)
			if err != nil {
				return fmt.Errorf("解析模板事件数据失败: %v", err)
			}
			return h.OnTemplateAccept(eventData, templateData)
		}

	case SubhookEventTemplateReject:
		if h.OnTemplateReject != nil {
			templateData, err := ParseTemplateSubhookEvent(eventData)
			if err != nil {
				return fmt.Errorf("解析模板事件数据失败: %v", err)
			}
			return h.OnTemplateReject(eventData, templateData)
		}

	default:
		return fmt.Errorf("未知的事件类型: %s", eventType)
	}

	return nil
}

// ===== HTTP 处理器助手函数 =====

// CreateSubhookHTTPHandler 创建 SUBHOOK HTTP 处理器
// 参数：
//   - subhookKey: SUBHOOK 密匙
//   - handler: 事件处理器
//
// 返回：
//   - http.HandlerFunc: HTTP 处理函数
func CreateSubhookHTTPHandler(subhookKey string, handler SubhookEventHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析事件数据
		eventData, err := ParseSubhookEvent(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("解析事件数据失败: %v", err), http.StatusBadRequest)
			return
		}

		// 验证签名
		if !ValidateSubhookSignature(eventData.Token, eventData.Signature, subhookKey) {
			http.Error(w, "签名验证失败", http.StatusForbidden)
			return
		}

		// 处理事件
		if err := handler.HandleEvent(eventData.Event, eventData); err != nil {
			http.Error(w, fmt.Sprintf("处理事件失败: %v", err), http.StatusInternalServerError)
			return
		}

		// 返回成功响应
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

// ===== 便捷方法 =====

// GetEventTypeDescription 获取事件类型描述
func GetEventTypeDescription(eventType string) string {
	descriptions := map[string]string{
		SubhookEventRequest:        "发送请求被接收",
		SubhookEventDelivered:      "发送成功",
		SubhookEventDropped:        "发送失败",
		SubhookEventSending:        "正在发送",
		SubhookEventMO:             "短信上行（用户回复）",
		SubhookEventTemplateAccept: "短信模板审核通过",
		SubhookEventTemplateReject: "短信模板审核未通过",
	}

	if desc, exists := descriptions[eventType]; exists {
		return desc
	}
	return "未知事件类型"
}

// IsValidEventType 检查事件类型是否有效
func IsValidEventType(eventType string) bool {
	validTypes := []string{
		SubhookEventRequest,
		SubhookEventDelivered,
		SubhookEventDropped,
		SubhookEventSending,
		SubhookEventMO,
		SubhookEventTemplateAccept,
		SubhookEventTemplateReject,
	}

	for _, validType := range validTypes {
		if eventType == validType {
			return true
		}
	}
	return false
}

// ValidateEventTypes 验证事件类型数组
func ValidateEventTypes(eventTypes []string) []string {
	var errors []string
	for _, eventType := range eventTypes {
		if !IsValidEventType(eventType) {
			errors = append(errors, fmt.Sprintf("无效的事件类型: %s", eventType))
		}
	}
	return errors
}
