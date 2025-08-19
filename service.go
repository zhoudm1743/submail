package submail

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 常量定义
const (
	// API基础URL
	DefaultBaseURL = "https://api-v4.mysubmail.com"
	LogBaseURL     = "https://log.mysubmail.com"

	// 响应格式
	FormatJSON = "json"
	FormatXML  = "xml"

	// 签名类型
	SignTypeMD5    = "md5"
	SignTypeSHA1   = "sha1"
	SignTypeNormal = "normal"

	// API端点
	EndpointSMSSend          = "/sms/send"
	EndpointSMSXSend         = "/sms/xsend"
	EndpointSMSMultiSend     = "/sms/multisend"
	EndpointSMSMultiXSend    = "/sms/multixsend"
	EndpointSMSBatchSend     = "/sms/batchsend"
	EndpointSMSBatchXSend    = "/sms/batchxsend"
	EndpointSMSUnionSend     = "/sms/unionsend"
	EndpointSMSTemplate      = "/sms/template"
	EndpointSMSReports       = "/sms/reports"
	EndpointSMSBalance       = "/balance/sms"
	EndpointSMSBalanceLog    = "/message/balancelog"
	EndpointSMSLog           = "/sms/log"
	EndpointSMSMO            = "/sms/mo"
	EndpointSMSAppextend     = "/sms/appextend"
	EndpointServiceTimestamp = "/service/timestamp"
	EndpointServiceStatus    = "/service/status"
)

// Client 赛邮云SDK客户端
type Client struct {
	AppID          string             // App ID (应用ID)
	AppKey         string             // App Key (应用密钥，用于签名计算)
	BaseURL        string             // API基础URL
	client         *http.Client       // HTTP客户端
	format         string             // 响应格式 (json/xml)
	useDigitalSign bool               // 是否使用数字签名模式，false为明文模式
	signType       string             // 签名类型：md5 或 sha1，仅数字签名模式使用
	timeout        time.Duration      // 请求超时时间
	varProcessor   *VariableProcessor // 变量处理器
}

// Config 客户端配置
type Config struct {
	AppID          string        // App ID (必填)
	AppKey         string        // App Key (必填)
	BaseURL        string        // API基础URL (可选，默认为官方API地址)
	Format         string        // 响应格式 (可选，默认json)
	UseDigitalSign bool          // 是否使用数字签名模式 (可选，默认false)
	SignType       string        // 签名类型 (可选，默认md5)
	Timeout        time.Duration // 请求超时时间 (可选，默认30秒)
}

// NewClient 创建新的赛邮云客户端
func NewClient(config Config) *Client {
	if config.BaseURL == "" {
		config.BaseURL = DefaultBaseURL
	}
	if config.Format == "" {
		config.Format = FormatJSON
	}
	if config.SignType == "" {
		config.SignType = SignTypeMD5
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &Client{
		AppID:          config.AppID,
		AppKey:         config.AppKey,
		BaseURL:        config.BaseURL,
		client:         &http.Client{Timeout: config.Timeout},
		format:         config.Format,
		useDigitalSign: config.UseDigitalSign,
		signType:       config.SignType,
		timeout:        config.Timeout,
		varProcessor:   NewVariableProcessor(),
	}
}

// ===== 变量处理方法 =====

// SetTimezone 设置时区
func (c *Client) SetTimezone(timezone string) error {
	return c.varProcessor.SetTimezone(timezone)
}

// ProcessVariables 处理短信内容中的变量
func (c *Client) ProcessVariables(content string, vars map[string]string) string {
	return c.varProcessor.ProcessVariables(content, vars)
}

// ValidateVariables 验证变量格式
func (c *Client) ValidateVariables(content string) []string {
	return c.varProcessor.ValidateVariables(content)
}

// ExtractVariableNames 提取内容中的自定义变量名
func (c *Client) ExtractVariableNames(content string) []string {
	return c.varProcessor.ExtractVariableNames(content)
}

// GetDateVariableDescription 获取日期变量说明
func (c *Client) GetDateVariableDescription() map[string]string {
	return c.varProcessor.GetDateVariableDescription()
}

// SMSSendWithVariables 发送带变量的短信（自动处理变量）
func (c *Client) SMSSendWithVariables(to, content string, vars map[string]string, tag string) (*SMSSendResponse, error) {
	// 验证变量格式
	if errors := c.ValidateVariables(content); len(errors) > 0 {
		return nil, fmt.Errorf("变量格式错误: %v", errors)
	}

	// 处理变量
	processedContent := c.ProcessVariables(content, vars)

	// 创建请求
	req := &SMSSendRequest{
		To:      to,
		Content: processedContent,
		Tag:     tag,
	}

	return c.SMSSend(req)
}

// SMSXSendWithSignature 使用自定义签名发送模板短信
func (c *Client) SMSXSendWithSignature(to, project, signature string, vars map[string]string, tag string) (*SMSSendResponse, error) {
	req := &SMSXSendRequest{
		To:           to,
		Project:      project,
		Vars:         vars,
		SMSSignature: signature,
		Tag:          tag,
	}

	return c.SMSXSend(req)
}

// SMSMultiSendWithVariables 使用变量发送一对多短信
func (c *Client) SMSMultiSendWithVariables(content string, recipients []SMSMultiItem, tag string) (*SMSMultiSendResponse, error) {
	// 验证变量格式
	if errors := c.ValidateVariables(content); len(errors) > 0 {
		return nil, fmt.Errorf("变量格式错误: %v", errors)
	}

	req := &SMSMultiSendRequest{
		Content: content,
		Multi:   recipients,
		Tag:     tag,
	}

	return c.SMSMultiSend(req)
}

// SMSMultiXSendWithSignature 使用自定义签名发送模板一对多短信
func (c *Client) SMSMultiXSendWithSignature(project, signature string, recipients []SMSMultiXItem, tag string) (*SMSMultiSendResponse, error) {
	req := &SMSMultiXSendRequest{
		Project:      project,
		Multi:        recipients,
		SMSSignature: signature,
		Tag:          tag,
	}

	return c.SMSMultiXSend(req)
}

// ===== 多条发送结果处理方法 =====

// GetSuccessResults 获取成功的发送结果
func (resp *SMSMultiSendResponse) GetSuccessResults() []SMSSendResult {
	var results []SMSSendResult
	for _, result := range *resp {
		if result.Status == "success" {
			results = append(results, result)
		}
	}
	return results
}

// GetFailedResults 获取失败的发送结果
func (resp *SMSMultiSendResponse) GetFailedResults() []SMSSendResult {
	var results []SMSSendResult
	for _, result := range *resp {
		if result.Status == "error" {
			results = append(results, result)
		}
	}
	return results
}

// GetTotalFee 获取总费用
func (resp *SMSMultiSendResponse) GetTotalFee() int {
	total := 0
	for _, result := range *resp {
		if result.Status == "success" {
			total += result.Fee
		}
	}
	return total
}

// GetStatistics 获取发送统计信息
func (resp *SMSMultiSendResponse) GetStatistics() (success, failed, totalFee int) {
	for _, result := range *resp {
		if result.Status == "success" {
			success++
			totalFee += result.Fee
		} else {
			failed++
		}
	}
	return
}

// ===== 批量发送结果处理方法 =====

// GetSuccessResults 获取成功的发送结果
func (resp *SMSBatchSendResponse) GetSuccessResults() []SMSSendResult {
	var results []SMSSendResult
	for _, result := range resp.Responses {
		if result.Status == "success" {
			results = append(results, result)
		}
	}
	return results
}

// GetFailedResults 获取失败的发送结果
func (resp *SMSBatchSendResponse) GetFailedResults() []SMSSendResult {
	var results []SMSSendResult
	for _, result := range resp.Responses {
		if result.Status == "error" {
			results = append(results, result)
		}
	}
	return results
}

// GetStatistics 获取发送统计信息
func (resp *SMSBatchSendResponse) GetStatistics() (success, failed int, totalFee int) {
	totalFee = resp.TotalFee
	for _, result := range resp.Responses {
		if result.Status == "success" {
			success++
		} else {
			failed++
		}
	}
	return
}

// SMSBatchSendWithPhones 批量发送短信（便捷方法）
func (c *Client) SMSBatchSendWithPhones(content string, phones []string, tag string) (*SMSBatchSendResponse, error) {
	// 验证变量格式
	if errors := c.ValidateVariables(content); len(errors) > 0 {
		return nil, fmt.Errorf("变量格式错误: %v", errors)
	}

	// 将手机号码数组转换为逗号分隔的字符串
	phoneStr := ""
	for i, phone := range phones {
		if i > 0 {
			phoneStr += ","
		}
		phoneStr += phone
	}

	req := &SMSBatchSendRequest{
		Content: content,
		To:      phoneStr,
		Tag:     tag,
	}

	return c.SMSBatchSend(req)
}

// SMSBatchXSendWithPhones 批量模板发送短信（便捷方法）
func (c *Client) SMSBatchXSendWithPhones(project string, phones []string, vars map[string]string, signature, tag string) (*SMSBatchSendResponse, error) {
	// 将手机号码数组转换为逗号分隔的字符串
	phoneStr := ""
	for i, phone := range phones {
		if i > 0 {
			phoneStr += ","
		}
		phoneStr += phone
	}

	req := &SMSBatchXSendRequest{
		Project:      project,
		To:           phoneStr,
		Vars:         vars,
		SMSSignature: signature,
		Tag:          tag,
	}

	return c.SMSBatchXSend(req)
}

// getTimestampFromServer 从服务器获取时间戳（内部使用，避免循环依赖）
func (c *Client) getTimestampFromServer() (int64, error) {
	// Service/Timestamp API 不需要授权参数
	params := make(map[string]string)

	body, err := c.doRequest("GET", EndpointServiceTimestamp, params)
	if err != nil {
		return 0, err
	}

	// 直接解析时间戳响应，格式为 {"timestamp": 1414253462}
	var timestampResp struct {
		Timestamp int64 `json:"timestamp"`
	}
	if err := json.Unmarshal(body, &timestampResp); err != nil {
		return 0, fmt.Errorf("解析时间戳响应失败: %v", err)
	}

	return timestampResp.Timestamp, nil
}

// buildSignature 构建签名
func (c *Client) buildSignature(params map[string]string) (string, error) {
	if !c.useDigitalSign {
		// 明文模式直接返回AppKey
		return c.AppKey, nil
	}

	// 数字签名模式
	timestamp, err := c.getTimestampFromServer()
	if err != nil {
		return "", fmt.Errorf("获取时间戳失败: %v", err)
	}
	params["timestamp"] = strconv.FormatInt(timestamp, 10)

	// 排序参数
	var keys []string
	for k := range params {
		// tag和sms_signature参数不参与加密计算
		if k != "signature" && k != "tag" && k != "sms_signature" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 构建签名字符串
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	signStr := strings.Join(parts, "&")

	// 加上APPID和APPKEY
	finalStr := c.AppID + c.AppKey + signStr + c.AppID + c.AppKey

	// 计算签名
	var signature string
	switch c.signType {
	case SignTypeMD5:
		hash := md5.Sum([]byte(finalStr))
		signature = fmt.Sprintf("%x", hash)
	case SignTypeSHA1:
		hash := sha1.Sum([]byte(finalStr))
		signature = fmt.Sprintf("%x", hash)
	default:
		return "", fmt.Errorf("不支持的签名类型: %s", c.signType)
	}

	return signature, nil
}

// buildAuthParams 构建认证参数
func (c *Client) buildAuthParams(params map[string]string) error {
	if params == nil {
		params = make(map[string]string)
	}

	// 添加基础参数
	params["appid"] = c.AppID

	if c.useDigitalSign {
		// 数字签名模式
		params["sign_type"] = c.signType

		signature, err := c.buildSignature(params)
		if err != nil {
			return err
		}
		params["signature"] = signature
	} else {
		// 明文模式
		params["signature"] = c.AppKey
	}

	return nil
}

// doRequest 执行HTTP请求
func (c *Client) doRequest(method, endpoint string, params map[string]string) ([]byte, error) {
	return c.doRequestWithBaseURL(method, endpoint, params, c.BaseURL)
}

// doRequestWithBaseURL 使用指定基础URL执行HTTP请求
func (c *Client) doRequestWithBaseURL(method, endpoint string, params map[string]string, baseURL string) ([]byte, error) {
	// 如果不是获取时间戳的请求，则构建认证参数
	if endpoint != EndpointServiceTimestamp {
		if err := c.buildAuthParams(params); err != nil {
			return nil, fmt.Errorf("构建认证参数失败: %v", err)
		}
	}

	// 构建URL
	requestURL := baseURL + endpoint
	if c.format == FormatXML {
		requestURL += ".xml"
	}

	var req *http.Request
	var err error

	if method == "GET" {
		// GET请求，参数放在URL中
		values := url.Values{}
		for k, v := range params {
			values.Set(k, v)
		}
		if len(values) > 0 {
			requestURL += "?" + values.Encode()
		}
		req, err = http.NewRequest("GET", requestURL, nil)
	} else {
		// POST请求，参数放在body中
		values := url.Values{}
		for k, v := range params {
			values.Set(k, v)
		}
		req, err = http.NewRequest("POST", requestURL, strings.NewReader(values.Encode()))
		if err == nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 执行请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("执行请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP错误: %d - %s", resp.StatusCode, string(body))
	}

	// 检查API错误
	if err := ParseAPIError(body); err != nil {
		return nil, err
	}

	return body, nil
}

// doJSONRequest 执行JSON请求
func (c *Client) doJSONRequest(method, endpoint string, data interface{}) ([]byte, error) {
	return c.doJSONRequestWithBaseURL(method, endpoint, data, c.BaseURL)
}

// doJSONRequestWithBaseURL 使用指定基础URL执行JSON请求
func (c *Client) doJSONRequestWithBaseURL(method, endpoint string, data interface{}, baseURL string) ([]byte, error) {
	// 将结构体转换为map[string]string
	params := make(map[string]string)

	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("序列化请求数据失败: %v", err)
		}

		var dataMap map[string]interface{}
		if err := json.Unmarshal(jsonData, &dataMap); err != nil {
			return nil, fmt.Errorf("解析请求数据失败: %v", err)
		}

		for k, v := range dataMap {
			switch val := v.(type) {
			case string:
				params[k] = val
			case float64:
				params[k] = strconv.FormatFloat(val, 'f', -1, 64)
			case bool:
				params[k] = strconv.FormatBool(val)
			case map[string]interface{}, []interface{}:
				// 对于复杂类型，重新序列化为JSON字符串
				jsonBytes, _ := json.Marshal(val)
				params[k] = string(jsonBytes)
			default:
				params[k] = fmt.Sprintf("%v", val)
			}
		}
	}

	return c.doRequestWithBaseURL(method, endpoint, params, baseURL)
}

// ===== 工具类API =====

// ServiceTimestamp 获取服务器时间戳
func (c *Client) ServiceTimestamp() (*ServiceTimestampResponse, error) {
	// Service/Timestamp API 不需要授权参数
	params := make(map[string]string)

	body, err := c.doRequest("GET", EndpointServiceTimestamp, params)
	if err != nil {
		return nil, err
	}

	// 直接解析时间戳响应，格式为 {"timestamp": 1414253462}
	var timestampResp struct {
		Timestamp int64 `json:"timestamp"`
	}
	if err := json.Unmarshal(body, &timestampResp); err != nil {
		return nil, fmt.Errorf("解析时间戳响应失败: %v", err)
	}

	return &ServiceTimestampResponse{
		BaseResponse: BaseResponse{Status: "success"},
		Timestamp:    timestampResp.Timestamp,
	}, nil
}

// GetCurrentTimestamp 获取当前服务器时间戳（便捷方法）
func (c *Client) GetCurrentTimestamp() (int64, error) {
	resp, err := c.ServiceTimestamp()
	if err != nil {
		return 0, err
	}
	return resp.Timestamp, nil
}

// ServiceStatus 获取服务器状态
func (c *Client) ServiceStatus() (*ServiceStatusResponse, error) {
	// Service/Status API 不需要授权参数
	params := make(map[string]string)

	body, err := c.doRequest("GET", EndpointServiceStatus, params)
	if err != nil {
		return nil, err
	}

	var resp ServiceStatusResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析服务状态响应失败: %v", err)
	}

	return &resp, nil
}

// IsServiceRunning 检查服务是否正常运行（便捷方法）
func (c *Client) IsServiceRunning() (bool, error) {
	resp, err := c.ServiceStatus()
	if err != nil {
		return false, err
	}
	return resp.Status == "runing", nil
}

// GetServiceRuntime 获取服务响应时间（便捷方法）
func (c *Client) GetServiceRuntime() (float64, error) {
	resp, err := c.ServiceStatus()
	if err != nil {
		return 0, err
	}
	return resp.Runtime, nil
}

// ===== 短信发送API =====

// SMSSend 短信发送
func (c *Client) SMSSend(req *SMSSendRequest) (*SMSSendResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("请求参数不能为空")
	}

	body, err := c.doJSONRequest("POST", EndpointSMSSend, req)
	if err != nil {
		return nil, err
	}

	var resp SMSSendResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析短信发送响应失败: %v", err)
	}

	return &resp, nil
}

// SMSXSend 短信模板发送
func (c *Client) SMSXSend(req *SMSXSendRequest) (*SMSSendResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("请求参数不能为空")
	}

	body, err := c.doJSONRequest("POST", EndpointSMSXSend, req)
	if err != nil {
		return nil, err
	}

	var resp SMSSendResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析短信模板发送响应失败: %v", err)
	}

	return &resp, nil
}

// SMSMultiSend 短信一对多发送
func (c *Client) SMSMultiSend(req *SMSMultiSendRequest) (*SMSMultiSendResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("请求参数不能为空")
	}

	body, err := c.doJSONRequest("POST", EndpointSMSMultiSend, req)
	if err != nil {
		return nil, err
	}

	var resp SMSMultiSendResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析短信一对多发送响应失败: %v", err)
	}

	return &resp, nil
}

// SMSMultiXSend 短信模板一对多发送
func (c *Client) SMSMultiXSend(req *SMSMultiXSendRequest) (*SMSMultiSendResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("请求参数不能为空")
	}

	body, err := c.doJSONRequest("POST", EndpointSMSMultiXSend, req)
	if err != nil {
		return nil, err
	}

	var resp SMSMultiSendResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析短信模板一对多发送响应失败: %v", err)
	}

	return &resp, nil
}

// SMSBatchSend 短信批量群发
func (c *Client) SMSBatchSend(req *SMSBatchSendRequest) (*SMSBatchSendResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("请求参数不能为空")
	}

	body, err := c.doJSONRequest("POST", EndpointSMSBatchSend, req)
	if err != nil {
		return nil, err
	}

	var resp SMSBatchSendResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析短信批量群发响应失败: %v", err)
	}

	return &resp, nil
}

// SMSBatchXSend 短信批量模板群发
func (c *Client) SMSBatchXSend(req *SMSBatchXSendRequest) (*SMSBatchSendResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("请求参数不能为空")
	}

	body, err := c.doJSONRequest("POST", EndpointSMSBatchXSend, req)
	if err != nil {
		return nil, err
	}

	var resp SMSBatchSendResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析短信批量模板群发响应失败: %v", err)
	}

	return &resp, nil
}

// SMSUnionSend 国内短信与国际短信联合发送
func (c *Client) SMSUnionSend(req *SMSUnionSendRequest) (*SMSSendResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("请求参数不能为空")
	}

	body, err := c.doJSONRequest("POST", EndpointSMSUnionSend, req)
	if err != nil {
		return nil, err
	}

	var resp SMSSendResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析短信联合发送响应失败: %v", err)
	}

	return &resp, nil
}

// SMSUnionSendWithConfig 国内外短信联合发送（便捷方法）
func (c *Client) SMSUnionSendWithConfig(to, content, interAppID, interSignature string, interContent, tag string, enableCodeTransform bool) (*SMSSendResponse, error) {
	codeTransform := "false"
	if enableCodeTransform {
		codeTransform = "true"
	}

	req := &SMSUnionSendRequest{
		To:                          to,
		Content:                     content,
		InterAppID:                  interAppID,
		InterSignature:              interSignature,
		InterContent:                interContent,
		IntersmsVerifyCodeTransform: codeTransform,
		Tag:                         tag,
	}

	return c.SMSUnionSend(req)
}

// IsInternationalNumber 判断是否为国际号码
func IsInternationalNumber(phoneNumber string) bool {
	// 国际号码以+开头且不是+86
	if len(phoneNumber) > 3 && phoneNumber[0] == '+' {
		return phoneNumber[:3] != "+86"
	}
	// 11位数字为国内号码
	if len(phoneNumber) == 11 {
		for _, r := range phoneNumber {
			if r < '0' || r > '9' {
				return true // 包含非数字字符，可能是国际号码
			}
		}
		return false // 纯11位数字，国内号码
	}
	// 其他情况视为国际号码
	return true
}

// ===== 短信签名管理API =====

// SMSSignatureQuery 查询短信签名
func (c *Client) SMSSignatureQuery(req *SMSSignatureQueryRequest) (*SMSSignatureQueryResponse, error) {
	body, err := c.doJSONRequest("GET", EndpointSMSAppextend, req)
	if err != nil {
		return nil, err
	}

	var resp SMSSignatureQueryResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析短信签名查询响应失败: %v", err)
	}

	return &resp, nil
}

// SMSSignatureCreate 创建短信签名
func (c *Client) SMSSignatureCreate(req *SMSSignatureCreateRequest) (*SMSSignatureOperationResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("请求参数不能为空")
	}

	body, err := c.doJSONRequest("POST", EndpointSMSAppextend, req)
	if err != nil {
		return nil, err
	}

	var resp SMSSignatureOperationResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析短信签名创建响应失败: %v", err)
	}

	return &resp, nil
}

// SMSSignatureUpdate 更新短信签名
func (c *Client) SMSSignatureUpdate(req *SMSSignatureUpdateRequest) (*SMSSignatureOperationResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("请求参数不能为空")
	}

	body, err := c.doJSONRequest("PUT", EndpointSMSAppextend, req)
	if err != nil {
		return nil, err
	}

	var resp SMSSignatureOperationResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析短信签名更新响应失败: %v", err)
	}

	return &resp, nil
}

// SMSSignatureDelete 删除短信签名
func (c *Client) SMSSignatureDelete(req *SMSSignatureDeleteRequest) (*SMSSignatureOperationResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("请求参数不能为空")
	}

	body, err := c.doJSONRequest("DELETE", EndpointSMSAppextend, req)
	if err != nil {
		return nil, err
	}

	var resp SMSSignatureOperationResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析短信签名删除响应失败: %v", err)
	}

	return &resp, nil
}

// GetSignatureStatus 获取签名状态描述
func GetSignatureStatus(status int) string {
	switch status {
	case 0:
		return "审核中"
	case 1:
		return "审核通过"
	default:
		return "审核不通过"
	}
}

// GetSourceTypeDescription 获取材料类型描述
func GetSourceTypeDescription(sourceType int) string {
	switch sourceType {
	case 0:
		return "营业执照"
	case 1:
		return "商标"
	case 2:
		return "APP"
	default:
		return "未知类型"
	}
}

// ===== 短信管理API =====

// SMSTemplateGet 获取短信模板列表或单个模板
func (c *Client) SMSTemplateGet(req *SMSTemplateGetRequest) (*SMSTemplateGetResponse, error) {
	body, err := c.doJSONRequest("GET", EndpointSMSTemplate, req)
	if err != nil {
		return nil, err
	}

	var resp SMSTemplateGetResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析短信模板查询响应失败: %v", err)
	}

	return &resp, nil
}

// SMSTemplateCreate 创建短信模板
func (c *Client) SMSTemplateCreate(req *SMSTemplateCreateRequest) (*SMSTemplateCreateResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("请求参数不能为空")
	}

	body, err := c.doJSONRequest("POST", EndpointSMSTemplate, req)
	if err != nil {
		return nil, err
	}

	var resp SMSTemplateCreateResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析短信模板创建响应失败: %v", err)
	}

	return &resp, nil
}

// SMSTemplateUpdate 更新短信模板
func (c *Client) SMSTemplateUpdate(req *SMSTemplateUpdateRequest) (*SMSTemplateOperationResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("请求参数不能为空")
	}

	body, err := c.doJSONRequest("PUT", EndpointSMSTemplate, req)
	if err != nil {
		return nil, err
	}

	var resp SMSTemplateOperationResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析短信模板更新响应失败: %v", err)
	}

	return &resp, nil
}

// SMSTemplateDelete 删除短信模板
func (c *Client) SMSTemplateDelete(req *SMSTemplateDeleteRequest) (*SMSTemplateOperationResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("请求参数不能为空")
	}

	body, err := c.doJSONRequest("DELETE", EndpointSMSTemplate, req)
	if err != nil {
		return nil, err
	}

	var resp SMSTemplateOperationResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析短信模板删除响应失败: %v", err)
	}

	return &resp, nil
}

// GetTemplateStatus 获取模板状态描述
func GetTemplateStatus(status string) string {
	switch status {
	case "0":
		return "未提交审核"
	case "1":
		return "正在审核"
	case "2":
		return "审核通过"
	case "3":
		return "未通过审核"
	default:
		return "未知状态"
	}
}

// GetTemplateAddTime 将UNIX时间戳转换为时间
func GetTemplateAddTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

// SMSReports 短信分析报告
func (c *Client) SMSReports(req *SMSReportsRequest) (*SMSReportsResponse, error) {
	body, err := c.doJSONRequest("POST", EndpointSMSReports, req)
	if err != nil {
		return nil, err
	}

	var resp SMSReportsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析短信分析报告响应失败: %v", err)
	}

	return &resp, nil
}

// SMSReportsWithDateRange 使用日期范围查询短信分析报告（便捷方法）
func (c *Client) SMSReportsWithDateRange(startDate, endDate time.Time) (*SMSReportsResponse, error) {
	req := &SMSReportsRequest{
		StartDate: startDate.Unix(),
		EndDate:   endDate.Unix(),
	}

	return c.SMSReports(req)
}

// SMSReportsLast7Days 获取最近7天的短信分析报告（便捷方法）
func (c *Client) SMSReportsLast7Days() (*SMSReportsResponse, error) {
	now := time.Now()
	startDate := now.AddDate(0, 0, -7) // 7天前

	return c.SMSReportsWithDateRange(startDate, now)
}

// SMSReportsLastMonth 获取上个月的短信分析报告（便捷方法）
func (c *Client) SMSReportsLastMonth() (*SMSReportsResponse, error) {
	now := time.Now()
	startDate := now.AddDate(0, -1, 0) // 1个月前

	return c.SMSReportsWithDateRange(startDate, now)
}

// SMSBalance 短信余额查询
func (c *Client) SMSBalance() (*SMSBalanceResponse, error) {
	body, err := c.doJSONRequest("POST", EndpointSMSBalance, nil)
	if err != nil {
		return nil, err
	}

	var resp SMSBalanceResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析短信余额响应失败: %v", err)
	}

	return &resp, nil
}

// SMSBalanceLog 短信余额日志查询
func (c *Client) SMSBalanceLog(req *SMSBalanceLogRequest) (*SMSBalanceLogResponse, error) {
	body, err := c.doJSONRequest("POST", EndpointSMSBalanceLog, req)
	if err != nil {
		return nil, err
	}

	var resp SMSBalanceLogResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析短信余额日志响应失败: %v", err)
	}

	return &resp, nil
}

// SMSBalanceLogWithDateRange 使用日期范围查询短信余额日志（便捷方法）
func (c *Client) SMSBalanceLogWithDateRange(startDate, endDate time.Time) (*SMSBalanceLogResponse, error) {
	req := &SMSBalanceLogRequest{
		StartDate: startDate.Unix(),
		EndDate:   endDate.Unix(),
	}

	return c.SMSBalanceLog(req)
}

// SMSBalanceLogLast7Days 获取最近7天的短信余额日志（便捷方法）
func (c *Client) SMSBalanceLogLast7Days() (*SMSBalanceLogResponse, error) {
	now := time.Now()
	startDate := now.AddDate(0, 0, -7) // 7天前

	return c.SMSBalanceLogWithDateRange(startDate, now)
}

// SMSBalanceLogLastMonth 获取上个月的短信余额日志（便捷方法）
func (c *Client) SMSBalanceLogLastMonth() (*SMSBalanceLogResponse, error) {
	now := time.Now()
	startDate := now.AddDate(0, -1, 0) // 1个月前

	return c.SMSBalanceLogWithDateRange(startDate, now)
}

// SMSLog 短信历史明细查询
func (c *Client) SMSLog(req *SMSLogRequest) (*SMSLogResponse, error) {
	body, err := c.doJSONRequestWithBaseURL("POST", EndpointSMSLog, req, LogBaseURL)
	if err != nil {
		return nil, err
	}

	var resp SMSLogResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析短信历史明细响应失败: %v", err)
	}

	return &resp, nil
}

// SMSLogWithDateRange 使用日期范围查询短信历史明细（便捷方法）
func (c *Client) SMSLogWithDateRange(startDate, endDate time.Time) (*SMSLogResponse, error) {
	req := &SMSLogRequest{
		StartDate: startDate.Unix(),
		EndDate:   endDate.Unix(),
		Rows:      50, // 默认返回50条
	}

	return c.SMSLog(req)
}

// SMSLogLast7Days 获取最近7天的短信历史明细（便捷方法）
func (c *Client) SMSLogLast7Days() (*SMSLogResponse, error) {
	now := time.Now()
	startDate := now.AddDate(0, 0, -7) // 7天前

	return c.SMSLogWithDateRange(startDate, now)
}

// SMSLogByPhone 根据手机号查询短信历史明细（便捷方法）
func (c *Client) SMSLogByPhone(phone string) (*SMSLogResponse, error) {
	req := &SMSLogRequest{
		To:   phone,
		Rows: 50,
	}

	return c.SMSLog(req)
}

// SMSLogBySendID 根据Send ID查询短信历史明细（便捷方法）
func (c *Client) SMSLogBySendID(sendID string) (*SMSLogResponse, error) {
	req := &SMSLogRequest{
		SendID: sendID,
		Rows:   50,
	}

	return c.SMSLog(req)
}

// SMSLogByStatus 根据状态查询短信历史明细（便捷方法）
func (c *Client) SMSLogByStatus(status string) (*SMSLogResponse, error) {
	req := &SMSLogRequest{
		Status: status, // "delivered" 或 "dropped"
		Rows:   50,
	}

	return c.SMSLog(req)
}

// SMSMO 短信上行查询
func (c *Client) SMSMO(req *SMSMORequest) (*SMSMOResponse, error) {
	body, err := c.doJSONRequest("POST", EndpointSMSMO, req)
	if err != nil {
		return nil, err
	}

	var resp SMSMOResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析短信上行查询响应失败: %v", err)
	}

	return &resp, nil
}

// SMSMOWithDateRange 使用日期范围查询短信上行（便捷方法）
func (c *Client) SMSMOWithDateRange(startDate, endDate time.Time) (*SMSMOResponse, error) {
	req := &SMSMORequest{
		StartDate: startDate.Unix(),
		EndDate:   endDate.Unix(),
		Rows:      50, // 默认返回50条
	}

	return c.SMSMO(req)
}

// SMSMOLast7Days 获取最近7天的短信上行（便捷方法）
func (c *Client) SMSMOLast7Days() (*SMSMOResponse, error) {
	now := time.Now()
	startDate := now.AddDate(0, 0, -7) // 7天前

	return c.SMSMOWithDateRange(startDate, now)
}

// SMSMOByPhone 根据手机号查询短信上行（便捷方法）
func (c *Client) SMSMOByPhone(phone string) (*SMSMOResponse, error) {
	req := &SMSMORequest{
		From: phone,
		Rows: 50,
	}

	return c.SMSMO(req)
}

// ===== 分析报告数据处理方法 =====

// GetSuccessRate 获取成功率
func (overview *SMSReportOverview) GetSuccessRate() float64 {
	if overview.Request == 0 {
		return 0
	}
	return float64(overview.Deliveryed) / float64(overview.Request) * 100
}

// GetFailureRate 获取失败率
func (overview *SMSReportOverview) GetFailureRate() float64 {
	if overview.Request == 0 {
		return 0
	}
	return float64(overview.Dropped) / float64(overview.Request) * 100
}

// GetTotalOperators 获取运营商总数
func (operators *SMSReportOperators) GetTotalOperators() int {
	return operators.ChinaMobile + operators.ChinaUnicom + operators.ChinaTelecom
}

// GetOperatorPercentage 获取运营商占比
func (operators *SMSReportOperators) GetOperatorPercentage() map[string]float64 {
	total := operators.GetTotalOperators()
	if total == 0 {
		return map[string]float64{
			"移动": 0,
			"联通": 0,
			"电信": 0,
		}
	}

	return map[string]float64{
		"移动": float64(operators.ChinaMobile) / float64(total) * 100,
		"联通": float64(operators.ChinaUnicom) / float64(total) * 100,
		"电信": float64(operators.ChinaTelecom) / float64(total) * 100,
	}
}

// GetTopProvinces 获取发送量最多的省份（前N个）
func (location *SMSReportLocation) GetTopProvinces(topN int) []ProvinceCount {
	var provinces []ProvinceCount
	for province, count := range location.Province {
		if province != "UNKOWN" { // 排除未知省份
			provinces = append(provinces, ProvinceCount{Province: province, Count: count})
		}
	}

	// 按数量排序
	for i := 0; i < len(provinces)-1; i++ {
		for j := i + 1; j < len(provinces); j++ {
			if provinces[i].Count < provinces[j].Count {
				provinces[i], provinces[j] = provinces[j], provinces[i]
			}
		}
	}

	if len(provinces) > topN {
		provinces = provinces[:topN]
	}

	return provinces
}

// GetTopFailureReasons 获取主要失败原因（前N个）
func (overview *SMSReportOverview) GetTopFailureReasons(topN int) []ReasonCount {
	var reasons []ReasonCount
	for reason, count := range overview.DroppedReasonAnalysis {
		reasons = append(reasons, ReasonCount{Reason: reason, Count: count})
	}

	// 按数量排序
	for i := 0; i < len(reasons)-1; i++ {
		for j := i + 1; j < len(reasons); j++ {
			if reasons[i].Count < reasons[j].Count {
				reasons[i], reasons[j] = reasons[j], reasons[i]
			}
		}
	}

	if len(reasons) > topN {
		reasons = reasons[:topN]
	}

	return reasons
}

// ProvinceCount 省份统计
type ProvinceCount struct {
	Province string
	Count    int
}

// ReasonCount 失败原因统计
type ReasonCount struct {
	Reason string
	Count  int
}

// ===== 余额日志数据处理方法 =====

// IsTransactionalSMSChange 判断是否为事务类短信余额变更
func (entry *SMSBalanceLogEntry) IsTransactionalSMSChange() bool {
	return entry.TMessageAddCredits != ""
}

// IsMarketingSMSChange 判断是否为运营类短信余额变更
func (entry *SMSBalanceLogEntry) IsMarketingSMSChange() bool {
	return entry.MessageAddCredits != ""
}

// GetChangeAmount 获取余额变更金额（正数为增加，负数为减少）
func (entry *SMSBalanceLogEntry) GetChangeAmount() (transactional, marketing int) {
	if entry.TMessageAddCredits != "" {
		if amount, err := strconv.Atoi(entry.TMessageAddCredits); err == nil {
			transactional = amount
		}
	}
	if entry.MessageAddCredits != "" {
		if amount, err := strconv.Atoi(entry.MessageAddCredits); err == nil {
			marketing = amount
		}
	}
	return
}

// GetBalanceChange 获取余额变更详情
func (entry *SMSBalanceLogEntry) GetBalanceChange() map[string]map[string]int {
	result := make(map[string]map[string]int)

	// 事务类短信余额变更
	if entry.IsTransactionalSMSChange() {
		transactional := make(map[string]int)
		if before, err := strconv.Atoi(entry.TMessageBeforeCredits); err == nil {
			transactional["before"] = before
		}
		if after, err := strconv.Atoi(entry.TMessageAfterCredits); err == nil {
			transactional["after"] = after
		}
		if change, err := strconv.Atoi(entry.TMessageAddCredits); err == nil {
			transactional["change"] = change
		}
		result["transactional"] = transactional
	}

	// 运营类短信余额变更
	if entry.IsMarketingSMSChange() {
		marketing := make(map[string]int)
		if before, err := strconv.Atoi(entry.MessageBeforeCredits); err == nil {
			marketing["before"] = before
		}
		if after, err := strconv.Atoi(entry.MessageAfterCredits); err == nil {
			marketing["after"] = after
		}
		if change, err := strconv.Atoi(entry.MessageAddCredits); err == nil {
			marketing["change"] = change
		}
		result["marketing"] = marketing
	}

	return result
}

// ParseDateTime 解析变更时间
func (entry *SMSBalanceLogEntry) ParseDateTime() (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", entry.Datetime)
}

// GetTotalChanges 获取余额日志的总变更统计
func (resp *SMSBalanceLogResponse) GetTotalChanges() (transactionalTotal, marketingTotal int) {
	for _, entry := range resp.Data {
		transactional, marketing := entry.GetChangeAmount()
		transactionalTotal += transactional
		marketingTotal += marketing
	}
	return
}

// GetChangesByType 按变更类型分组统计
func (resp *SMSBalanceLogResponse) GetChangesByType() map[string][]SMSBalanceLogEntry {
	result := make(map[string][]SMSBalanceLogEntry)

	for _, entry := range resp.Data {
		if entry.IsTransactionalSMSChange() {
			result["transactional"] = append(result["transactional"], entry)
		}
		if entry.IsMarketingSMSChange() {
			result["marketing"] = append(result["marketing"], entry)
		}
	}

	return result
}

// ===== SMS Log数据处理方法 =====

// IsDelivered 判断短信是否成功发送
func (log *SMSLog) IsDelivered() bool {
	return log.Status == "delivered"
}

// IsDropped 判断短信是否发送失败
func (log *SMSLog) IsDropped() bool {
	return log.Status == "dropped"
}

// IsPending 判断短信状态是否未知
func (log *SMSLog) IsPending() bool {
	return log.Status == "pending"
}

// GetSendTime 获取请求时间
func (log *SMSLog) GetSendTime() time.Time {
	return time.Unix(log.SendAt, 0)
}

// GetSentTime 获取平台发送时间
func (log *SMSLog) GetSentTime() time.Time {
	return time.Unix(log.SentAt, 0)
}

// GetReportTime 获取运营商状态汇报时间
func (log *SMSLog) GetReportTime() time.Time {
	return time.Unix(log.ReportAt, 0)
}

// GetDeliveryDuration 获取发送到汇报的时长
func (log *SMSLog) GetDeliveryDuration() time.Duration {
	if log.ReportAt > 0 && log.SentAt > 0 {
		return time.Unix(log.ReportAt, 0).Sub(time.Unix(log.SentAt, 0))
	}
	return 0
}

// GetSuccessLogs 获取成功发送的日志
func (resp *SMSLogResponse) GetSuccessLogs() []SMSLog {
	var results []SMSLog
	for _, log := range resp.Data {
		if log.IsDelivered() {
			results = append(results, log)
		}
	}
	return results
}

// GetFailedLogs 获取发送失败的日志
func (resp *SMSLogResponse) GetFailedLogs() []SMSLog {
	var results []SMSLog
	for _, log := range resp.Data {
		if log.IsDropped() {
			results = append(results, log)
		}
	}
	return results
}

// GetPendingLogs 获取状态未知的日志
func (resp *SMSLogResponse) GetPendingLogs() []SMSLog {
	var results []SMSLog
	for _, log := range resp.Data {
		if log.IsPending() {
			results = append(results, log)
		}
	}
	return results
}

// GetLogStatistics 获取日志统计信息
func (resp *SMSLogResponse) GetLogStatistics() (success, failed, pending, totalFee int) {
	for _, log := range resp.Data {
		switch log.Status {
		case "delivered":
			success++
			totalFee += log.Fee
		case "dropped":
			failed++
		case "pending":
			pending++
		}
	}
	return
}

// GetLogsByOperator 按运营商分组日志
func (resp *SMSLogResponse) GetLogsByOperator() map[string][]SMSLog {
	result := make(map[string][]SMSLog)
	for _, log := range resp.Data {
		operator := log.MobileType
		if operator == "" {
			operator = "未知"
		}
		result[operator] = append(result[operator], log)
	}
	return result
}

// GetLogsByLocation 按地区分组日志
func (resp *SMSLogResponse) GetLogsByLocation() map[string][]SMSLog {
	result := make(map[string][]SMSLog)
	for _, log := range resp.Data {
		location := log.Location
		if location == "" {
			location = "未知"
		}
		result[location] = append(result[location], log)
	}
	return result
}

// GetFailureReasons 获取失败原因统计
func (resp *SMSLogResponse) GetFailureReasons() map[string]int {
	result := make(map[string]int)
	for _, log := range resp.Data {
		if log.IsDropped() && log.DroppedReason != "" {
			result[log.DroppedReason]++
		}
	}
	return result
}

// GetLogsByTemplate 按模板ID分组日志
func (resp *SMSLogResponse) GetLogsByTemplate() map[string][]SMSLog {
	result := make(map[string][]SMSLog)
	for _, log := range resp.Data {
		templateID := log.TemplateID
		if templateID == "" {
			templateID = "未知"
		}
		result[templateID] = append(result[templateID], log)
	}
	return result
}

// ===== 服务状态数据处理方法 =====

// IsRunning 判断服务是否正常运行
func (status *ServiceStatusResponse) IsRunning() bool {
	return status.Status == "runing"
}

// IsHealthy 判断服务是否健康（运行正常且响应时间合理）
func (status *ServiceStatusResponse) IsHealthy() bool {
	return status.IsRunning() && status.Runtime < 2.0 // 响应时间小于2秒认为是健康的
}

// GetPerformanceLevel 获取服务性能等级
func (status *ServiceStatusResponse) GetPerformanceLevel() string {
	if !status.IsRunning() {
		return "服务异常"
	}

	if status.Runtime < 0.1 {
		return "优秀"
	} else if status.Runtime < 0.5 {
		return "良好"
	} else if status.Runtime < 1.0 {
		return "一般"
	} else if status.Runtime < 2.0 {
		return "较慢"
	} else {
		return "很慢"
	}
}

// GetStatusDescription 获取状态描述
func (status *ServiceStatusResponse) GetStatusDescription() string {
	if status.IsRunning() {
		return fmt.Sprintf("服务正常运行，响应时间: %.3f秒 (%s)",
			status.Runtime, status.GetPerformanceLevel())
	}
	return fmt.Sprintf("服务状态: %s", status.Status)
}

// ===== SMS MO数据处理方法 =====

// GetReplyTime 获取回复时间
func (mo *SMSMO) GetReplyTime() time.Time {
	return time.Unix(mo.ReplyAt, 0)
}

// IsReturnReceipt 判断是否为回执（退订）
func (mo *SMSMO) IsReturnReceipt() bool {
	content := strings.ToLower(strings.TrimSpace(mo.Content))
	// 常见的退订关键词
	unsubscribeKeywords := []string{"td", "退订", "t", "0000", "00000", "n", "unsubscribe", "stop"}
	for _, keyword := range unsubscribeKeywords {
		if content == keyword {
			return true
		}
	}
	return false
}

// IsValidReply 判断是否为有效回复（非退订）
func (mo *SMSMO) IsValidReply() bool {
	return !mo.IsReturnReceipt() && strings.TrimSpace(mo.Content) != ""
}

// GetMOStatistics 获取上行统计信息
func (resp *SMSMOResponse) GetMOStatistics() (total, validReplies, unsubscribes int) {
	total = len(resp.MO)
	for _, mo := range resp.MO {
		if mo.IsReturnReceipt() {
			unsubscribes++
		} else if mo.IsValidReply() {
			validReplies++
		}
	}
	return
}

// GetValidReplies 获取有效回复
func (resp *SMSMOResponse) GetValidReplies() []SMSMO {
	var results []SMSMO
	for _, mo := range resp.MO {
		if mo.IsValidReply() {
			results = append(results, mo)
		}
	}
	return results
}

// GetUnsubscribes 获取退订回复
func (resp *SMSMOResponse) GetUnsubscribes() []SMSMO {
	var results []SMSMO
	for _, mo := range resp.MO {
		if mo.IsReturnReceipt() {
			results = append(results, mo)
		}
	}
	return results
}

// GetMOByPhone 根据手机号分组上行
func (resp *SMSMOResponse) GetMOByPhone() map[string][]SMSMO {
	result := make(map[string][]SMSMO)
	for _, mo := range resp.MO {
		result[mo.From] = append(result[mo.From], mo)
	}
	return result
}

// GetMOByContent 根据回复内容分组
func (resp *SMSMOResponse) GetMOByContent() map[string][]SMSMO {
	result := make(map[string][]SMSMO)
	for _, mo := range resp.MO {
		content := strings.TrimSpace(mo.Content)
		if content == "" {
			content = "空内容"
		}
		result[content] = append(result[content], mo)
	}
	return result
}

// GetMOByBatch 根据批次号分组上行
func (resp *SMSMOResponse) GetMOByBatch() map[string][]SMSMO {
	result := make(map[string][]SMSMO)
	for _, mo := range resp.MO {
		batch := mo.SendList
		if batch == "" {
			batch = "未知批次"
		}
		result[batch] = append(result[batch], mo)
	}
	return result
}

// GetReplyRate 获取回复率（基于下行短信内容）
func (resp *SMSMOResponse) GetReplyRate() map[string]float64 {
	result := make(map[string]float64)
	contentStats := make(map[string]int)

	// 统计每个下行内容的回复数
	for _, mo := range resp.MO {
		content := mo.SMSContent
		if content == "" {
			content = "未知内容"
		}
		contentStats[content]++
	}

	// 这里只能统计回复数，实际回复率需要结合发送总数计算
	for content, count := range contentStats {
		result[content] = float64(count) // 实际使用时需要除以该内容的发送总数
	}

	return result
}
