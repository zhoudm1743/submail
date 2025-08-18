# SUBMAIL 赛邮 Go SDK

这是一个用于 SUBMAIL 赛邮云通信服务的 Go SDK，支持短信发送、模板短信、余额查询等功能。

## 功能特性

- ✅ 短信发送 (SMS/Send)
- ✅ 模板短信发送 (SMS/XSend)
- ✅ 短信一对多发送 (SMS/Multisend)
- ✅ 模板短信一对多发送 (SMS/MultiXSend)
- ✅ 短信模板管理 (SMS/Template)
- ✅ 短信分析报告 (SMS/Reports)
- ✅ 历史明细查询 (SMS/Log)
- ✅ 短信上行查询 (SMS/MO)
- ✅ 账户余额查询 (SMS/Balance)
- ✅ 服务器状态查询 (Service/Status)
- ✅ 时间戳获取 (Service/Timestamp)
- ✅ API 签名验证
- ✅ 错误处理
- ✅ 超时控制
- ✅ 自定义 API 地址
- ✅ 响应格式支持 (JSON/XML)

## 安装

```bash
go get github.com/zhoudm1743/submail
```

## 快速开始

### 1. 获取 API 凭证

在使用之前，您需要在 [SUBMAIL 控制台](https://www.mysubmail.com) 创建应用并获取：

- `AppID`: 应用ID
- `AppKey`: 应用密钥  
- `Signature`: 应用签名

### 2. 初始化服务

```go
package main

import (
    "github.com/zhoudm1743/submail"
)

func main() {
    // 创建服务实例（默认JSON格式）
    service := submail.NewSaiyouService(
        "your_app_id",
        "your_app_key", 
        "your_signature",
    )
    
    // 或者创建指定XML格式的服务实例
    serviceXML := submail.NewSaiyouServiceWithFormat(
        "your_app_id",
        "your_app_key", 
        "your_signature",
        submail.FormatXML,
    )
}
```

### 3. 发送短信

```go
// 发送普通短信
smsReq := &submail.SMSRequest{
    To:   "13800138000",
    Text: "您的验证码是1234，请在5分钟内使用。",
    Project: "test", // 可选
}

resp, err := service.SendSMS(smsReq)
if err != nil {
    log.Printf("发送失败: %v", err)
} else {
    fmt.Printf("发送成功! ID: %s\n", resp.SendID)
}
```

### 4. 发送模板短信

```go
// 发送模板短信
templateReq := &submail.SMSXRequest{
    To:      "13800138000",
    Project: "your_template_id", // 模板ID
    Vars: map[string]string{
        "code": "1234",
        "time": "5",
    },
}

resp, err := service.SendSMSTemplate(templateReq)
if err != nil {
    log.Printf("发送失败: %v", err)
} else {
    fmt.Printf("发送成功! ID: %s\n", resp.SendID)
}
```

### 5. 查询余额

```go
balance, err := service.GetBalance()
if err != nil {
    log.Printf("查询失败: %v", err)
} else {
    fmt.Printf("账户余额: %s\n", balance.Balance)
}
```

### 6. 短信一对多发送

```go
// 短信一对多发送
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

resp, err := service.SendSMSMulti(multiReq)
if err != nil {
    log.Printf("发送失败: %v", err)
} else {
    fmt.Printf("发送成功! 状态: %s\n", resp.Status)
}
```

### 7. 短信模板管理

```go
// 获取模板列表
templates, err := service.GetSMSTemplates()

// 创建模板
resp, err := service.CreateSMSTemplate("您的验证码是{code}", 1)

// 更新模板
resp, err := service.UpdateSMSTemplate("template_id", "新的模板内容{code}", 1)

// 删除模板
resp, err := service.DeleteSMSTemplate("template_id")
```

### 8. 查询分析报告和历史记录

```go
// 查询分析报告
reportsReq := &submail.SMSReportsRequest{
    Project:   "test",
    StartDate: "2024-01-01 00:00:00",
    EndDate:   "2024-12-31 23:59:59",
}
reports, err := service.GetSMSReports(reportsReq)

// 查询历史明细
logReq := &submail.SMSLogRequest{
    Project: "test",
    Limit:   10,
    Offset:  0,
}
logs, err := service.GetSMSLog(logReq)

// 查询短信上行
moReq := &submail.SMSMORequest{
    StartDate: "2024-01-01 00:00:00",
    EndDate:   "2024-12-31 23:59:59",
    Limit:     10,
}
mo, err := service.GetSMSMO(moReq)
```

### 9. 响应格式设置

```go
// 获取当前响应格式
format := service.GetFormat() // 返回 "json" 或 "xml"

// 设置响应格式为XML
service.SetFormat(submail.FormatXML)

// 设置响应格式为JSON（默认）
service.SetFormat(submail.FormatJSON)

// 创建时指定格式
serviceXML := submail.NewSaiyouServiceWithFormat(
    "your_app_id",
    "your_app_key", 
    "your_signature",
    submail.FormatXML,
)
```

### 10. 签名调试

```go
// 用于调试签名问题
params := url.Values{}
params.Set("to", "13800138000")
params.Set("text", "测试短信")

signString, signature := service.ValidateSignature(params)
fmt.Printf("签名字符串: %s\n", signString)
fmt.Printf("计算签名: %s\n", signature)
```

## API 文档

### 数据结构

#### SMSRequest - 短信发送请求

```go
type SMSRequest struct {
    To      string            // 收件人手机号码 (必填)
    Text    string            // 短信正文 (必填)
    Vars    map[string]string // 文本变量 (可选)
    Project string            // 项目标记 (可选)
    Tag     string            // 自定义标签 (可选)
}
```

#### SMSXRequest - 模板短信发送请求

```go
type SMSXRequest struct {
    To      string            // 收件人手机号码 (必填)
    Project string            // 短信模板标记 (必填)
    Vars    map[string]string // 文本变量 (可选)
    Tag     string            // 自定义标签 (可选)
}
```

#### SMSMultisendRequest - 短信一对多发送请求

```go
type SMSMultisendRequest struct {
    Multi   []SMSMultiItem // 联系人列表 (必填)
    Text    string         // 短信正文 (必填)
    Project string         // 项目标记 (可选)
}

type SMSMultiItem struct {
    To   string            // 收件人手机号 (必填)
    Vars map[string]string // 个性化变量 (可选)
}
```

#### SMSMultiXSendRequest - 短信模板一对多发送请求

```go
type SMSMultiXSendRequest struct {
    Multi   []SMSMultiItem // 联系人列表 (必填)
    Project string         // 短信模板标记 (必填)
}
```

#### APIResponse - API响应

```go
type APIResponse struct {
    Status string // 请求状态
    Code   int    // 状态码
    Msg    string // 响应消息
    SendID string // 发送ID
    Fee    int    // 消费费用
    Sms    int    // 短信条数
}
```

### 方法

#### NewSaiyouService

创建新的服务实例（默认JSON格式）。

```go
func NewSaiyouService(appID, appKey, signature string) *SaiyouService
```

#### NewSaiyouServiceWithFormat

创建新的服务实例并指定响应格式。

```go
func NewSaiyouServiceWithFormat(appID, appKey, signature, format string) *SaiyouService
```

#### SetBaseURL

设置自定义的 API 基础地址（用于私有部署）。

```go
func (s *SaiyouService) SetBaseURL(baseURL string)
```

#### SetFormat

设置响应格式 (json 或 xml)。

```go
func (s *SaiyouService) SetFormat(format string)
```

#### GetFormat

获取当前响应格式。

```go
func (s *SaiyouService) GetFormat() string
```

#### SendSMS

发送普通短信。

```go
func (s *SaiyouService) SendSMS(req *SMSRequest) (*APIResponse, error)
```

#### SendSMSTemplate

发送模板短信。

```go
func (s *SaiyouService) SendSMSTemplate(req *SMSXRequest) (*APIResponse, error)
```

#### GetBalance

查询账户余额。

```go
func (s *SaiyouService) GetBalance() (*BalanceResponse, error)
```

#### SendSMSMulti

发送短信一对多。

```go
func (s *SaiyouService) SendSMSMulti(req *SMSMultisendRequest) (*MultiSendResponse, error)
```

#### SendSMSMultiTemplate

发送短信模板一对多。

```go
func (s *SaiyouService) SendSMSMultiTemplate(req *SMSMultiXSendRequest) (*MultiSendResponse, error)
```

#### GetSMSTemplates

获取短信模板列表。

```go
func (s *SaiyouService) GetSMSTemplates() (*APIResponse, error)
```

#### CreateSMSTemplate

创建短信模板。

```go
func (s *SaiyouService) CreateSMSTemplate(sms string, templateType int) (*APIResponse, error)
```

#### UpdateSMSTemplate

更新短信模板。

```go
func (s *SaiyouService) UpdateSMSTemplate(templateID, sms string, templateType int) (*APIResponse, error)
```

#### DeleteSMSTemplate

删除短信模板。

```go
func (s *SaiyouService) DeleteSMSTemplate(templateID string) (*APIResponse, error)
```

#### GetSMSReports

获取短信分析报告。

```go
func (s *SaiyouService) GetSMSReports(req *SMSReportsRequest) (*APIResponse, error)
```

#### GetSMSLog

查询短信历史明细。

```go
func (s *SaiyouService) GetSMSLog(req *SMSLogRequest) (*APIResponse, error)
```

#### GetSMSMO

查询短信上行。

```go
func (s *SaiyouService) GetSMSMO(req *SMSMORequest) (*APIResponse, error)
```

#### GetTimestamp

获取服务器时间戳。

```go
func (s *SaiyouService) GetTimestamp() (*APIResponse, error)
```

#### GetStatus

获取服务器状态。

```go
func (s *SaiyouService) GetStatus() (*APIResponse, error)
```

#### SendSMSBatch

短信批量群发。

```go
func (s *SaiyouService) SendSMSBatch(req *SMSBatchSendRequest) (*APIResponse, error)
```

**请求参数：**
- `req.To`: 收件人手机号码列表
- `req.Text`: 短信正文
- `req.Project`: 项目标记（可选）
- `req.Tag`: 自定义标签（可选）

**示例：**
```go
req := &submail.SMSBatchSendRequest{
    To:   []string{"13800138000", "13800138001", "13800138002"},
    Text: "您的验证码是：123456，请勿泄露给他人。",
    Project: "verification",
    Tag:   "batch_send",
}

resp, err := service.SendSMSBatch(req)
```

#### SendSMSBatchTemplate

短信批量模板群发。

```go
func (s *SaiyouService) SendSMSBatchTemplate(req *SMSBatchXSendRequest) (*APIResponse, error)
```

**请求参数：**
- `req.To`: 收件人手机号码列表
- `req.Project`: 短信模板标记
- `req.Vars`: 文本变量（可选）
- `req.Tag`: 自定义标签（可选）

**示例：**
```go
req := &submail.SMSBatchXSendRequest{
    To:      []string{"13800138000", "13800138001", "13800138002"},
    Project: "verification_template",
    Vars:    map[string]string{"code": "123456"},
    Tag:     "batch_template",
}

resp, err := service.SendSMSBatchTemplate(req)
```

#### SendSMSUnion

国内短信与国际短信联合发送。

```go
func (s *SaiyouService) SendSMSUnion(req *SMSUnionSendRequest) (*APIResponse, error)
```

**请求参数：**
- `req.To`: 收件人手机号码
- `req.Text`: 短信正文
- `req.Project`: 项目标记（可选）
- `req.Tag`: 自定义标签（可选）
- `req.Country`: 国家代码（可选，不传默认为中国）

**示例：**
```go
req := &submail.SMSUnionSendRequest{
    To:      "13800138000",
    Text:    "您的验证码是：123456，请勿泄露给他人。",
    Tag:     "union_send",
    Country: "US", // 美国
}

resp, err := service.SendSMSUnion(req)
```

#### SubscribeSMS

短信订阅。

```go
func (s *SaiyouService) SubscribeSMS(to, project string) (*APIResponse, error)
```

**请求参数：**
- `to`: 手机号码
- `project`: 项目标记（可选）

**示例：**
```go
resp, err := service.UnsubscribeSMS("13800138000", "newsletter")
```

#### UnsubscribeSMS

短信退订。

```go
func (s *SaiyouService) UnsubscribeSMS(to, project string) (*APIResponse, error)
```

**请求参数：**
- `to`: 手机号码
- `project`: 项目标记（可选）

**示例：**
```go
resp, err := service.SubscribeSMS("13800138000", "newsletter")
```

#### ValidateSignature

验证签名（用于调试）。

```go
func (s *SaiyouService) ValidateSignature(params url.Values) (signString, signature string)
```

#### SyncServerTime

手动同步服务器时间。当遇到时间相关错误时，可以主动调用此方法同步时间。

```go
func (s *SaiyouService) SyncServerTime() (int64, error)
```

**返回值：**
- `int64`: 服务器时间戳
- `error`: 错误信息

**示例：**
```go
serverTime, err := service.SyncServerTime()
if err != nil {
    log.Printf("同步服务器时间失败: %v", err)
} else {
    fmt.Printf("服务器时间戳: %d\n", serverTime)
}
```

#### GetTimeOffset

获取本地时间与服务器时间的偏移量。

```go
func (s *SaiyouService) GetTimeOffset() (int64, error)
```

**返回值：**
- `int64`: 时间偏移量（秒）
  - 正值：本地时间快于服务器时间
  - 负值：本地时间慢于服务器时间
- `error`: 错误信息

**示例：**
```go
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
```

## 错误处理

SDK 会返回详细的错误信息，包括：

- 参数验证错误
- 网络请求错误  
- API 响应错误
- JSON 解析错误

### 自动时间同步重试

SDK 内置了智能重试机制，当遇到时间相关的鉴权错误时，会自动：

1. 检测到时间相关错误（如签名过期、时间戳无效等）
2. 调用服务器时间戳接口获取准确时间
3. 使用服务器时间重新生成签名
4. 自动重试原始请求

这大大提高了SDK的容错性，减少了因时间差异导致的鉴权失败。

### 手动时间同步

如果遇到特殊的时间同步需求，也可以手动调用时间同步方法：

```go
// 获取服务器时间戳
serverTime, err := service.SyncServerTime()

// 获取时间偏移量
offset, err := service.GetTimeOffset()
```

```go
resp, err := service.SendSMS(req)
if err != nil {
    // 处理错误
    log.Printf("发送短信失败: %v", err)
    return
}

// 检查 API 响应状态
if resp.Status != "success" {
    log.Printf("API 调用失败: %s", resp.Msg)
    return
}
```

## 配置选项

### 响应格式配置

```go
// 默认JSON格式
service := submail.NewSaiyouService(appID, appKey, signature)

// 指定XML格式
serviceXML := submail.NewSaiyouServiceWithFormat(appID, appKey, signature, submail.FormatXML)

// 动态切换格式
service.SetFormat(submail.FormatXML)  // 切换到XML
service.SetFormat(submail.FormatJSON) // 切换到JSON
```

### 自定义超时时间

```go
service := submail.NewSaiyouService(appID, appKey, signature)
// 默认超时时间为 30 秒，如需自定义可通过反射或重构代码
```

### 自定义 API 地址

```go
service.SetBaseURL("https://your-custom-api.example.com")
```

## 示例代码

完整的示例代码请查看 `example/main.go` 文件。

## 支持的功能

- [x] 短信发送 (SMS/Send)
- [x] 模板短信发送 (SMS/XSend)
- [x] 短信一对多发送 (SMS/Multisend)
- [x] 模板短信一对多发送 (SMS/MultiXSend)
- [x] 短信批量群发 (SMS/BatchSend)
- [x] 短信批量模板群发 (SMS/BatchXSend)
- [x] 国内短信与国际短信联合发送 (SMS/UnionSend)
- [x] 短信模板管理 (SMS/Template)
- [x] 短信分析报告 (SMS/Reports)
- [x] 历史明细查询 (SMS/Log)
- [x] 短信上行查询 (SMS/MO)
- [x] 账户余额查询 (SMS/Balance)
- [x] 短信订阅管理 (AddressBook/SMS/Subscribe)
- [x] 短信退订管理 (AddressBook/SMS/Unsubscribe)
- [x] 服务器状态查询 (Service/Status)
- [x] 时间戳获取 (Service/Timestamp)

## 参考文档

- [SUBMAIL 官方文档](https://www.mysubmail.com/documents/LJ4xa2)
- [API 授权验证机制](https://www.mysubmail.com/documents/VBcbe)

## 许可证

本项目采用 Apache 2.0 许可证，详见 LICENSE 文件。

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目。