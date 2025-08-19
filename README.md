# SUBMAIL 赛邮云 Go SDK

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

SUBMAIL 赛邮云官方 Go SDK，支持短信发送、模板管理、统计分析等完整功能。

## 特性

- ✅ 完整的短信发送功能
- ✅ 支持文本变量和日期变量
- ✅ 动态短信签名支持（v4.002）
- ✅ 多种发送模式：单发、一对多、批量群发、国内外联合发送
- ✅ 明文和数字签名两种认证方式
- ✅ 完善的错误处理机制（100+错误码定义）
- ✅ 统计分析和历史查询
- ✅ 余额查询和变更日志
- ✅ 短信上行回复查询
- ✅ 签名和模板管理
- ✅ 服务状态监控
- ✅ 丰富的便捷方法和数据处理功能

## 安装

```bash
go get github.com/zhoudm1743/submail
```

## 快速开始

### 基础配置

```go
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
        AppID:          "your-app-id",     // 您的App ID
        AppKey:         "your-app-key",    // 您的App Key
        BaseURL:        submail.DefaultBaseURL, // 默认API地址
        Format:         submail.FormatJSON,     // JSON格式
        UseDigitalSign: false,                  // 明文模式（测试推荐）
        SignType:       submail.SignTypeMD5,    // 数字签名类型
        Timeout:        30 * time.Second,       // 超时时间
    }

    // 创建客户端
    client := submail.NewClient(config)
}
```

### 基本短信发送

```go
// 1. 简单短信发送
req := &submail.SMSSendRequest{
    To:      "13800138000",
    Content: "【测试签名】您的验证码是123456，请在10分钟内输入。",
    Tag:     "test",
}

resp, err := client.SMSSend(req)
if err != nil {
    log.Printf("发送失败: %v", err)
} else {
    fmt.Printf("发送成功 - SendID: %s, 费用: %d\n", resp.SendID, resp.Fee)
}
```

### 模板短信发送

```go
// 2. 模板短信发送
templateReq := &submail.SMSXSendRequest{
    To:      "13800138000",
    Project: "your-template-id",
    Vars: map[string]string{
        "code": "123456",
        "time": "10",
    },
    Tag: "template-test",
}

templateResp, err := client.SMSXSend(templateReq)
```

### 自定义签名发送（v4.002新功能）

```go
// 3. 使用自定义签名发送模板短信
resp, err := client.SMSXSendWithSignature(
    "13800138000",        // 收件人
    "your-template-id",   // 模板ID
    "自定义签名",          // 自定义签名
    map[string]string{    // 变量
        "code": "888888",
        "time": "10",
    },
    "custom-signature-test", // 标签
)
```

## 变量功能

### 文本变量

支持在短信内容中使用 `@var(变量名)` 格式的文本变量：

```go
// 使用变量发送短信
content := "【测试签名】尊敬的@var(name)，您的验证码是@var(code)，请在@var(expire)分钟内输入。"
vars := map[string]string{
    "name":   "张三",
    "code":   "123456", 
    "expire": "10",
}

resp, err := client.SMSendWithVariables("13800138000", content, vars, "var-test")
```

### 日期变量

支持多种日期时间变量：

```go
content := "【测试签名】当前时间：@date()，年份：@date(Y)，月份：@date(m)，日期：@date(d)"

// 支持的日期变量：
// @date()   - 完整日期时间 (2024-01-15 14:30:25)
// @date(Y)  - 年份 (2024)
// @date(m)  - 月份 (01)
// @date(d)  - 日期 (15)
// @date(h)  - 小时 (14)
// @date(i)  - 分钟 (30)
// @date(s)  - 秒钟 (25)
// 更多格式请参考文档
```

### 变量验证和提取

```go
// 验证变量格式
if errors := client.ValidateVariables(content); len(errors) > 0 {
    log.Printf("变量格式错误: %v", errors)
}

// 提取变量名
varNames := client.ExtractVariableNames(content)
fmt.Printf("变量名: %v\n", varNames)

// 获取日期变量说明
descriptions := client.GetDateVariableDescription()
```

## 发送模式

### 1. 单条发送
- `SMSSend` - 普通短信发送
- `SMSXSend` - 模板短信发送

### 2. 一对多发送
支持对不同用户发送个性化内容：

```go
// 一对多发送（支持每个用户不同的变量）
multiReq := &submail.SMSMultiSendRequest{
    Content: "【测试签名】您好，@var(name)，您的取货码为 @var(code)",
    Multi: []submail.SMSMultiItem{
        {
            To: "13800138000",
            Vars: map[string]string{
                "name": "张三",
                "code": "A001",
            },
        },
        {
            To: "13800138001", 
            Vars: map[string]string{
                "name": "李四",
                "code": "B002",
            },
        },
    },
    Tag: "multi-test",
}

multiResp, err := client.SMSMultiSend(multiReq)

// 处理结果
success, failed, totalFee := multiResp.GetStatistics()
fmt.Printf("成功: %d条, 失败: %d条, 总费用: %d\n", success, failed, totalFee)
```

### 3. 批量群发
适合向大量用户发送相同内容（最多10000个号码）：

```go
// 批量群发
phones := []string{"13800138000", "13800138001", "13800138002"}
content := "【测试签名】系统维护通知，预计@date(h)点完成。"

batchResp, err := client.SMSBatchSendWithPhones(content, phones, "batch-test")

// 查看结果
success, failed, totalFee := batchResp.GetStatistics()
fmt.Printf("批量发送 - 任务ID: %s, 成功: %d条, 总费用: %d\n", 
    batchResp.BatchList, success, totalFee)
```

### 4. 国内外联合发送
支持在一个接口中同时发送国内和国际短信：

```go
// 国内外联合发送配置
config := &submail.SMSUnionConfig{
    DomesticAppID:    "domestic-app-id",
    DomesticAppKey:   "domestic-app-key",
    InternationalAppID:    "international-app-id", 
    InternationalAppKey:   "international-app-key",
    DomesticContent:  "【测试签名】您的验证码是123456，请在10分钟内输入。",
    InternationalContent: "[Test] Your verification code is 123456, valid for 10 minutes.",
    VerifyCodeTransform: true, // 自动提取验证码
}

// 发送到国内号码
resp1, err := client.SMSUnionSendWithConfig("+86138000138000", config, "union-test")

// 发送到国际号码  
resp2, err := client.SMSUnionSendWithConfig("+1234567890", config, "union-test")
```

## API 列表

### 短信发送
- `SMSSend` - 短信发送
- `SMSXSend` - 短信模板发送  
- `SMSMultiSend` - 短信一对多发送
- `SMSMultiXSend` - 短信模板一对多发送
- `SMSBatchSend` - 短信批量群发
- `SMSBatchXSend` - 短信批量模板群发
- `SMSUnionSend` - 国内外短信联合发送

### 便捷方法
- `SMSSendWithVariables` - 带变量的短信发送
- `SMSXSendWithSignature` - 自定义签名模板发送
- `SMSMultiSendWithVariables` - 带变量的一对多发送
- `SMSMultiXSendWithSignature` - 自定义签名模板一对多发送
- `SMSBatchSendWithPhones` - 批量发送便捷方法
- `SMSBatchXSendWithPhones` - 批量模板发送便捷方法
- `SMSUnionSendWithConfig` - 国内外联合发送便捷方法

### 查询和分析
- `SMSBalance` - 余额查询
- `SMSBalanceLog` - 余额变更日志查询
- `SMSReports` - 分析报告
- `SMSLog` - 历史明细查询
- `SMSMO` - 上行回复查询

### 模板和签名管理
- `SMSTemplateGet` - 获取短信模板
- `SMSTemplateCreate` - 创建短信模板
- `SMSTemplateUpdate` - 更新短信模板
- `SMSTemplateDelete` - 删除短信模板
- `SMSSignatureQuery` - 查询短信签名
- `SMSSignatureCreate` - 创建短信签名
- `SMSSignatureUpdate` - 更新短信签名
- `SMSSignatureDelete` - 删除短信签名

### 工具功能
- `ServiceTimestamp` - 获取服务器时间戳
- `ServiceStatus` - 获取服务器状态
- `GetCurrentTimestamp` - 获取当前时间戳（便捷方法）
- `IsServiceRunning` - 检查服务运行状态（便捷方法）

## 认证方式

### 明文模式（推荐测试使用）
```go
config := submail.Config{
    AppID:          "your-app-id",
    AppKey:         "your-app-key", 
    UseDigitalSign: false,  // 明文模式
}
```

### 数字签名模式（推荐生产使用）
```go
config := submail.Config{
    AppID:          "your-app-id",
    AppKey:         "your-app-key",
    UseDigitalSign: true,           // 数字签名模式
    SignType:       submail.SignTypeMD5, // MD5或SHA1
}
```

## 错误处理

SDK 提供完善的错误处理机制：

```go
resp, err := client.SMSSend(req)
if err != nil {
    // 检查是否为API错误
    if apiErr, ok := err.(*submail.APIError); ok {
        fmt.Printf("API错误 - 代码: %d, 消息: %s, 描述: %s\n", 
            apiErr.Code, apiErr.Msg, apiErr.Description)
    } else {
        fmt.Printf("其他错误: %v\n", err)
    }
}
```

## 响应处理

### 单条发送响应
```go
type SMSSendResponse struct {
    Status string `json:"status"`   // 状态
    SendID string `json:"send_id"`  // 发送ID
    Fee    int    `json:"fee"`      // 费用
    Sms    int    `json:"sms"`      // 短信条数
}
```

### 多条发送响应
```go
// 一对多发送响应处理
multiResp, err := client.SMSMultiSend(req)
if err == nil {
    // 获取统计信息
    success, failed, totalFee := multiResp.GetStatistics()
    
    // 获取成功的结果
    successResults := multiResp.GetSuccessResults()
    
    // 获取失败的结果  
    failedResults := multiResp.GetFailedResults()
}

// 批量发送响应处理
batchResp, err := client.SMSBatchSend(req)
if err == nil {
    fmt.Printf("任务ID: %s\n", batchResp.BatchList)
    success, failed, totalFee := batchResp.GetStatistics()
}
```

### 查询和分析功能

```go
// 1. 余额查询
balanceResp, err := client.SMSBalance()
if err == nil {
    fmt.Printf("通用类余额: %s, 事务类余额: %s\n", 
        balanceResp.Balance, balanceResp.TransactionalBalance)
}

// 2. 余额变更日志
balanceLogResp, err := client.SMSBalanceLogLast7Days()
if err == nil {
    transactionalTotal, marketingTotal := balanceLogResp.GetTotalChanges()
    fmt.Printf("最近7天变更 - 事务类: %d, 运营类: %d\n", 
        transactionalTotal, marketingTotal)
}

// 3. 历史明细查询
logResp, err := client.SMSLogLast7Days()
if err == nil {
    success, failed, pending, totalFee := logResp.GetLogStatistics()
    fmt.Printf("最近7天统计 - 成功: %d, 失败: %d, 未知: %d, 费用: %d\n",
        success, failed, pending, totalFee)
    
    // 按运营商统计
    operatorStats := logResp.GetLogsByOperator()
    // 失败原因分析
    failureReasons := logResp.GetFailureReasons()
}

// 4. 上行回复查询
moResp, err := client.SMSMOLast7Days()
if err == nil {
    total, validReplies, unsubscribes := moResp.GetMOStatistics()
    fmt.Printf("上行统计 - 总数: %d, 有效回复: %d, 退订: %d\n",
        total, validReplies, unsubscribes)
    
    // 获取退订回复
    unsubscribeList := moResp.GetUnsubscribes()
}

// 5. 分析报告
reportsResp, err := client.SMSReportsLast7Days()
if err == nil {
    successRate := reportsResp.Overview.GetSuccessRate()
    failureRate := reportsResp.Overview.GetFailureRate()
    fmt.Printf("成功率: %.2f%%, 失败率: %.2f%%\n", successRate, failureRate)
    
    // 运营商分布
    totalOperators := reportsResp.GetTotalOperators()
    // 地区分布
    topProvinces := reportsResp.GetTopProvinces(5)
}
```

### 模板和签名管理

```go
// 1. 模板管理
templates, err := client.SMSTemplateGet(&submail.SMSTemplateGetRequest{})
if err == nil {
    for _, template := range templates.Templates {
        status := submail.GetTemplateStatus(template.TemplateStatus)
        fmt.Printf("模板 %s: %s (%s)\n", 
            template.TemplateID, template.SMSTitle, status)
    }
}

// 2. 签名管理
signatures, err := client.SMSSignatureQuery(&submail.SMSSignatureQueryRequest{})
if err == nil {
    for _, sig := range signatures.SMSSignature {
        status := submail.GetSignatureStatus(sig.Status)
        fmt.Printf("签名 %s: %s\n", sig.SMSSignature, status)
    }
}
```

### 服务状态监控

```go
// 1. 服务器时间戳
timestampResp, err := client.ServiceTimestamp()
if err == nil {
    serverTime := time.Unix(timestampResp.Timestamp, 0)
    fmt.Printf("服务器时间: %s\n", serverTime.Format("2006-01-02 15:04:05"))
}

// 2. 服务器状态
statusResp, err := client.ServiceStatus()
if err == nil {
    fmt.Printf("服务状态: %s\n", statusResp.GetStatusDescription())
    fmt.Printf("性能等级: %s\n", statusResp.GetPerformanceLevel())
    fmt.Printf("是否健康: %t\n", statusResp.IsHealthy())
}
```

## 最佳实践

### 1. 生产环境配置
```go
config := submail.Config{
    AppID:          "your-app-id",
    AppKey:         "your-app-key",
    UseDigitalSign: true,                    // 使用数字签名
    SignType:       submail.SignTypeMD5,     // MD5签名
    Timeout:        30 * time.Second,        // 合理的超时时间
}
```

### 2. 错误重试机制
```go
func sendSMSWithRetry(client *submail.Client, req *submail.SMSSendRequest, maxRetries int) (*submail.SMSSendResponse, error) {
    var lastErr error
    for i := 0; i < maxRetries; i++ {
        resp, err := client.SMSSend(req)
        if err == nil {
            return resp, nil
        }
        lastErr = err
        time.Sleep(time.Second * time.Duration(i+1)) // 递增延迟
    }
    return nil, lastErr
}
```

### 3. 批量发送优化
```go
// 对于大量号码，建议分批处理
func sendBatchInChunks(client *submail.Client, content string, phones []string, chunkSize int) {
    for i := 0; i < len(phones); i += chunkSize {
        end := i + chunkSize
        if end > len(phones) {
            end = len(phones)
        }
        
        chunk := phones[i:end]
        resp, err := client.SMSBatchSendWithPhones(content, chunk, "batch")
        if err != nil {
            log.Printf("批次 %d 发送失败: %v", i/chunkSize+1, err)
        } else {
            success, failed, fee := resp.GetStatistics()
            log.Printf("批次 %d 完成 - 成功: %d, 失败: %d, 费用: %d", 
                i/chunkSize+1, success, failed, fee)
        }
    }
}
```

## 发送模式对比

| 模式 | API | 适用场景 | 最大数量 | 个性化 | 特殊功能 |
|------|-----|----------|----------|--------|----------|
| 单条发送 | SMSSend/SMSXSend | 单个用户 | 1 | ✅ | 支持变量 |
| 一对多发送 | SMSMultiSend/SMSMultiXSend | 少量个性化 | 50-200 | ✅ | 每用户不同变量 |
| 批量群发 | SMSBatchSend/SMSBatchXSend | 大量相同内容 | 10000 | ❌ | 高效群发 |
| 联合发送 | SMSUnionSend | 国内外混发 | 1 | ✅ | 自动识别国内外 |

## 常见问题

### Q: 如何选择发送模式？
- **单条发送**：适合实时验证码、通知等
- **一对多发送**：适合需要个性化内容的场景，如取货通知
- **批量群发**：适合营销短信、系统通知等大量相同内容
- **联合发送**：适合需要同时向国内外用户发送的场景

### Q: 明文模式和数字签名模式的区别？
- **明文模式**：简单快速，适合测试环境
- **数字签名模式**：更安全，适合生产环境，支持MD5和SHA1

### Q: 如何处理发送失败？
- 检查错误码和错误信息（支持100+错误码）
- 实现重试机制
- 记录失败日志便于分析
- 使用历史明细查询API分析失败原因

### Q: 如何监控短信发送效果？
- 使用 `SMSReports` 查看分析报告（成功率、运营商分布等）
- 使用 `SMSLog` 查询历史明细
- 使用 `SMSMO` 查询用户回复和退订情况
- 使用 `SMSBalance` 监控余额变化

### Q: 如何管理短信模板和签名？
- 使用模板管理API：`SMSTemplateGet`、`SMSTemplateCreate`、`SMSTemplateUpdate`、`SMSTemplateDelete`
- 使用签名管理API：`SMSSignatureQuery`、`SMSSignatureCreate`、`SMSSignatureUpdate`、`SMSSignatureDelete`
- 支持动态签名功能（v4.002新增）

### Q: 变量功能如何使用？
- 文本变量：`@var(变量名)`，支持个性化内容
- 日期变量：`@date()` 或 `@date(格式)`，自动填充时间
- 支持变量验证和提取功能
- 可设置时区：`client.SetTimezone("Asia/Shanghai")`

## 版本更新

### v1.0.2 (2025-08-19)
- ✅ 完整的短信发送功能（单发、一对多、批量、联合发送）
- ✅ 动态短信签名支持（v4.002）
- ✅ 完善的变量处理功能（文本变量 + 日期变量）
- ✅ 短信模板和签名管理
- ✅ 余额查询和变更日志
- ✅ 历史明细查询和分析
- ✅ 短信上行回复查询
- ✅ 统计分析报告
- ✅ 服务状态监控
- ✅ 完善的错误处理（100+错误码定义）
- ✅ 丰富的便捷方法和数据处理功能
- ✅ 支持明文和数字签名两种认证方式

## 许可证

MIT License

## 支持

- 官方文档：https://www.mysubmail.com/documents/
- 技术支持：4008-753-365
- 邮箱支持：service@submail.cn

---

**注意**：使用前请确保已在 [SUBMAIL 控制台](https://www.mysubmail.com) 创建应用并获取 App ID 和 App Key。