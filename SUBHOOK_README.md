# SUBMAIL SUBHOOK 功能说明

## 概述

SUBHOOK 是 SUBMAIL 提供的事件推送通知接口，类似于 Webhook。当特定事件发生时（如短信发送成功、失败、用户回复等），SUBMAIL 会向您指定的回调 URL 发送 HTTP POST 请求，通知您相关事件信息。

## 主要功能

### 1. SUBHOOK 管理
- **创建 SUBHOOK**：设置回调 URL 和监听的事件类型
- **查询 SUBHOOK**：查看已创建的 SUBHOOK 列表
- **删除 SUBHOOK**：删除不需要的 SUBHOOK

### 2. 事件类型
支持以下事件类型：
- `request`: 发送请求被接收
- `delivered`: 发送成功
- `dropped`: 发送失败
- `sending`: 正在发送
- `mo`: 短信上行（用户回复）
- `template_accept`: 短信模板审核通过
- `template_reject`: 短信模板审核未通过

### 3. 数据验证
实现了完整的签名验证机制，确保接收到的事件通知数据的安全性和完整性。

## 快速开始

### 1. 创建 SUBHOOK

```go
package main

import (
    "fmt"
    "log"
    "github.com/your-org/submail"
)

func main() {
    // 创建客户端
    client := submail.NewClient(submail.Config{
        AppID:  "your_app_id",
        AppKey: "your_app_key",
    })

    // 创建短信相关事件的 SUBHOOK
    resp, err := client.SubhookCreateForSMS("https://your-domain.com/subhook", "sms_events")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("SUBHOOK 创建成功:\n")
    fmt.Printf("ID: %s\n", resp.Target)
    fmt.Printf("密匙: %s\n", resp.Key) // 请妥善保存此密匙
}
```

### 2. 创建事件处理器

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/your-org/submail"
)

func main() {
    // SUBHOOK 密匙（从创建 SUBHOOK 时获得）
    subhookKey := "your_subhook_key"

    // 创建事件处理器
    handler := &submail.DefaultSubhookEventHandler{
        OnDelivered: func(eventData *submail.SubhookEventData, smsData *submail.SMSSubhookEventData) error {
            fmt.Printf("短信发送成功: SendID=%s, To=%s\n", smsData.SendID, smsData.To)
            return nil
        },
        OnDropped: func(eventData *submail.SubhookEventData, smsData *submail.SMSSubhookEventData) error {
            fmt.Printf("短信发送失败: SendID=%s, To=%s\n", smsData.SendID, smsData.To)
            return nil
        },
        OnMO: func(eventData *submail.SubhookEventData, moData *submail.SMSMOSubhookEventData) error {
            fmt.Printf("收到短信回复: From=%s, Content=%s\n", moData.From, moData.Content)
            return nil
        },
    }

    // 创建 HTTP 处理器
    httpHandler := submail.CreateSubhookHTTPHandler(subhookKey, handler)

    // 设置路由
    http.HandleFunc("/subhook", httpHandler)

    // 启动服务器
    fmt.Println("SUBHOOK 服务器启动在 :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### 3. 自定义事件处理器

```go
type CustomHandler struct {
    // 添加自定义字段
}

// 实现 SubhookEventHandler 接口
func (h *CustomHandler) HandleEvent(eventType string, eventData *submail.SubhookEventData) error {
    switch eventType {
    case submail.SubhookEventDelivered:
        smsData, _ := submail.ParseSMSSubhookEvent(eventData)
        // 处理发送成功事件
        fmt.Printf("短信发送成功: %s -> %s\n", smsData.SendID, smsData.To)
        
    case submail.SubhookEventMO:
        moData, _ := submail.ParseSMSMOSubhookEvent(eventData)
        // 处理短信回复事件
        fmt.Printf("收到回复: %s 说: %s\n", moData.From, moData.Content)
    }
    
    return nil
}
```

## 数据验证

### 签名验证原理

当 SUBHOOK 推送事件通知时，POST 数据中会包含：
- `token`: 32位随机字符串
- `signature`: 数字签名

验证步骤：
1. 将 `token` 和您的 SUBHOOK 密匙拼接：`token + key`
2. 对拼接字符串进行 MD5 哈希
3. 将生成的签名与 `signature` 参数比较

### 手动验证示例

```go
func verifySignature(token, signature, key string) bool {
    return submail.ValidateSubhookSignature(token, signature, key)
}

// 使用示例
isValid := verifySignature("received_token", "received_signature", "your_subhook_key")
if !isValid {
    // 签名验证失败，拒绝处理请求
    return
}
```

## API 参考

### 创建 SUBHOOK

```go
// 创建指定事件的 SUBHOOK
resp, err := client.SubhookCreateWithEvents(url, events, tag)

// 便捷方法
resp, err := client.SubhookCreateForSMS(url, tag)        // 短信相关事件
resp, err := client.SubhookCreateForMO(url, tag)         // 短信上行事件
resp, err := client.SubhookCreateForTemplate(url, tag)   // 模板审核事件
resp, err := client.SubhookCreateForAll(url, tag)        // 所有事件
```

### 查询 SUBHOOK

```go
// 查询所有 SUBHOOK
resp, err := client.SubhookQueryAll()

// 查询指定 ID 的 SUBHOOK
resp, err := client.SubhookQueryByID(target)
```

### 删除 SUBHOOK

```go
// 删除指定 ID 的 SUBHOOK
resp, err := client.SubhookDeleteByID(target)
```

## 事件数据结构

### 短信发送事件数据

```go
type SMSSubhookEventData struct {
    SendID   string // 发送ID
    To       string // 收件人
    Content  string // 短信内容
    Status   string // 状态
    Fee      int    // 费用
    SendAt   int64  // 发送时间
    ReportAt int64  // 汇报时间
}
```

### 短信上行事件数据

```go
type SMSMOSubhookEventData struct {
    From       string // 发送方手机号
    Content    string // 上行内容
    ReplyAt    int64  // 回复时间
    SMSContent string // 对应的下行短信内容
}
```

### 模板审核事件数据

```go
type TemplateSubhookEventData struct {
    TemplateID string // 模板ID
    Status     string // 审核状态
    Reason     string // 审核原因（拒绝时）
}
```

## 最佳实践

1. **安全性**
   - 始终验证签名
   - 使用 HTTPS 端点
   - 妥善保管 SUBHOOK 密匙

2. **可靠性**
   - 实现幂等性处理
   - 添加重试机制
   - 记录事件日志

3. **性能**
   - 异步处理事件
   - 避免长时间阻塞
   - 及时响应 HTTP 200

4. **监控**
   - 监控事件接收情况
   - 记录处理失败的事件
   - 设置告警机制

## 故障排除

### 常见问题

1. **签名验证失败**
   - 检查 SUBHOOK 密匙是否正确
   - 确认 token 和 signature 参数获取正确
   - 验证 MD5 计算逻辑

2. **事件接收不到**
   - 检查回调 URL 是否可访问
   - 确认服务器返回 HTTP 200
   - 查看 SUBMAIL 控制台的推送日志

3. **事件重复接收**
   - 实现幂等性处理
   - 使用事件 ID 去重
   - 检查服务器响应是否及时

### 调试建议

```go
// 开启调试日志
func debugHandler(eventType string, eventData *submail.SubhookEventData) {
    fmt.Printf("收到事件: %s\n", eventType)
    fmt.Printf("Token: %s\n", eventData.Token)
    fmt.Printf("Signature: %s\n", eventData.Signature)
    fmt.Printf("AppID: %s\n", eventData.AppID)
    fmt.Printf("Data: %+v\n", eventData.Data)
}
```

## 完整示例

请参考 `example/subhook_example.go` 文件，其中包含了完整的使用示例和最佳实践。

## 支持

如有问题，请查阅：
- [SUBMAIL 官方文档](https://www.mysubmail.com/documents/Hl8hW)
- 项目 Issue 页面
- 技术支持邮箱
