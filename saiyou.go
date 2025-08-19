package submail

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// SUBMAIL API常量
const (
	DefaultBaseURL       = "https://api-v4.mysubmail.com"
	SMSEndpoint          = "/sms/send"
	SMSXEndpoint         = "/sms/xsend"
	SMSMultisendEndpoint = "/sms/multisend"
	SMSMultiXEndpoint    = "/sms/multixsend"
	SMSBatchSendEndpoint = "/sms/batchsend"
	SMSBatchXEndpoint    = "/sms/batchxsend"
	SMSUnionSendEndpoint = "/sms/unionsend"
	SMSTemplateEndpoint  = "/sms/template"
	SMSReportsEndpoint   = "/sms/reports"
	SMSLogEndpoint       = "/sms/log"
	SMSMOEndpoint        = "/sms/mo"
	BalanceEndpoint      = "/sms/balance"
	TimestampEndpoint    = "/service/timestamp"
	StatusEndpoint       = "/service/status"
	AddressBookEndpoint  = "/addressbook/sms"

	// 响应格式常量
	FormatJSON = "json"
	FormatXML  = "xml"

	// 默认连接池配置
	DefaultMaxIdleConns        = 100              // 最大空闲连接数
	DefaultMaxIdleConnsPerHost = 10               // 每个主机的最大空闲连接数
	DefaultMaxConnsPerHost     = 0                // 每个主机的最大连接数，0表示无限制
	DefaultIdleConnTimeout     = 90 * time.Second // 空闲连接超时时间
	DefaultTLSHandshakeTimeout = 10 * time.Second // TLS握手超时时间
	DefaultDialTimeout         = 30 * time.Second // 拨号超时时间
	DefaultKeepAlive           = 30 * time.Second // TCP Keep-Alive间隔
)

// ConnectionPoolConfig 连接池配置
type ConnectionPoolConfig struct {
	// 最大空闲连接数
	MaxIdleConns int
	// 每个主机的最大空闲连接数
	MaxIdleConnsPerHost int
	// 每个主机的最大连接数，0表示无限制
	MaxConnsPerHost int
	// 空闲连接超时时间
	IdleConnTimeout time.Duration
	// TLS握手超时时间
	TLSHandshakeTimeout time.Duration
	// 拨号超时时间
	DialTimeout time.Duration
	// TCP Keep-Alive间隔
	KeepAlive time.Duration
	// 客户端请求超时时间
	RequestTimeout time.Duration
}

// DefaultConnectionPoolConfig 返回默认连接池配置
func DefaultConnectionPoolConfig() *ConnectionPoolConfig {
	return &ConnectionPoolConfig{
		MaxIdleConns:        DefaultMaxIdleConns,
		MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
		MaxConnsPerHost:     DefaultMaxConnsPerHost,
		IdleConnTimeout:     DefaultIdleConnTimeout,
		TLSHandshakeTimeout: DefaultTLSHandshakeTimeout,
		DialTimeout:         DefaultDialTimeout,
		KeepAlive:           DefaultKeepAlive,
		RequestTimeout:      30 * time.Second,
	}
}

// createHTTPClientWithPool 根据连接池配置创建HTTP客户端
func createHTTPClientWithPool(config *ConnectionPoolConfig) *http.Client {
	// 创建自定义的传输层配置
	transport := &http.Transport{
		// 连接池配置
		MaxIdleConns:        config.MaxIdleConns,
		MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		MaxConnsPerHost:     config.MaxConnsPerHost,
		IdleConnTimeout:     config.IdleConnTimeout,
		TLSHandshakeTimeout: config.TLSHandshakeTimeout,

		// 拨号配置
		DialContext: (&net.Dialer{
			Timeout:   config.DialTimeout,
			KeepAlive: config.KeepAlive,
		}).DialContext,

		// 启用HTTP/2
		ForceAttemptHTTP2: true,

		// 其他优化配置
		DisableCompression: false, // 启用压缩
		DisableKeepAlives:  false, // 启用Keep-Alive
	}

	return &http.Client{
		Transport: transport,
		Timeout:   config.RequestTimeout,
	}
}

// SaiyouService SUBMAIL赛邮服务结构体
type SaiyouService struct {
	AppID          string                // App ID (应用ID)
	AppKey         string                // App Key (应用密钥，用于签名计算)
	BaseURL        string                // API基础URL
	client         *http.Client          // HTTP客户端
	format         string                // 响应格式 (json/xml)
	poolConfig     *ConnectionPoolConfig // 连接池配置
	useDigitalSign bool                  // 是否使用数字签名模式，false为明文模式
	signType       string                // 签名类型：md5 或 sha1，仅数字签名模式使用
}

// NewSaiyouService 创建新的赛邮服务实例（使用默认连接池配置，默认数字签名模式）
func NewSaiyouService(appID, appKey string) *SaiyouService {
	poolConfig := DefaultConnectionPoolConfig()
	return &SaiyouService{
		AppID:          appID,
		AppKey:         appKey,
		BaseURL:        DefaultBaseURL,
		client:         createHTTPClientWithPool(poolConfig),
		format:         FormatJSON, // 默认使用JSON格式
		poolConfig:     poolConfig,
		useDigitalSign: true,  // 默认使用数字签名模式
		signType:       "md5", // 默认使用MD5签名
	}
}

// NewSaiyouServiceWithFormat 创建新的赛邮服务实例并指定响应格式（使用默认连接池配置，默认数字签名模式）
func NewSaiyouServiceWithFormat(appID, appKey, format string) *SaiyouService {
	// 验证格式参数
	if format != FormatJSON && format != FormatXML {
		format = FormatJSON // 无效格式时使用默认JSON
	}

	poolConfig := DefaultConnectionPoolConfig()
	return &SaiyouService{
		AppID:          appID,
		AppKey:         appKey,
		BaseURL:        DefaultBaseURL,
		client:         createHTTPClientWithPool(poolConfig),
		format:         format,
		poolConfig:     poolConfig,
		useDigitalSign: true,  // 默认使用数字签名模式
		signType:       "md5", // 默认使用MD5签名
	}
}

// NewSaiyouServiceWithPool 创建新的赛邮服务实例并指定连接池配置（默认数字签名模式）
func NewSaiyouServiceWithPool(appID, appKey string, poolConfig *ConnectionPoolConfig) *SaiyouService {
	if poolConfig == nil {
		poolConfig = DefaultConnectionPoolConfig()
	}

	return &SaiyouService{
		AppID:          appID,
		AppKey:         appKey,
		BaseURL:        DefaultBaseURL,
		client:         createHTTPClientWithPool(poolConfig),
		format:         FormatJSON, // 默认使用JSON格式
		poolConfig:     poolConfig,
		useDigitalSign: true,  // 默认使用数字签名模式
		signType:       "md5", // 默认使用MD5签名
	}
}

// NewSaiyouServiceWithPoolAndFormat 创建新的赛邮服务实例并指定连接池配置和响应格式（默认数字签名模式）
func NewSaiyouServiceWithPoolAndFormat(appID, appKey, format string, poolConfig *ConnectionPoolConfig) *SaiyouService {
	// 验证格式参数
	if format != FormatJSON && format != FormatXML {
		format = FormatJSON // 无效格式时使用默认JSON
	}

	if poolConfig == nil {
		poolConfig = DefaultConnectionPoolConfig()
	}

	return &SaiyouService{
		AppID:          appID,
		AppKey:         appKey,
		BaseURL:        DefaultBaseURL,
		client:         createHTTPClientWithPool(poolConfig),
		format:         format,
		poolConfig:     poolConfig,
		useDigitalSign: true,  // 默认使用数字签名模式
		signType:       "md5", // 默认使用MD5签名
	}
}

// SetBaseURL 设置自定义的基础URL
func (s *SaiyouService) SetBaseURL(baseURL string) {
	s.BaseURL = baseURL
}

// SetFormat 设置响应格式 (json 或 xml)
func (s *SaiyouService) SetFormat(format string) {
	if format == FormatJSON || format == FormatXML {
		s.format = format
	}
}

// NewSaiyouServiceWithPlaintextAuth 创建使用明文验证模式的赛邮服务实例
func NewSaiyouServiceWithPlaintextAuth(appID, appKey string) *SaiyouService {
	poolConfig := DefaultConnectionPoolConfig()
	return &SaiyouService{
		AppID:          appID,
		AppKey:         appKey,
		BaseURL:        DefaultBaseURL,
		client:         createHTTPClientWithPool(poolConfig),
		format:         FormatJSON,
		poolConfig:     poolConfig,
		useDigitalSign: false, // 使用明文验证模式
		signType:       "",    // 明文模式不需要签名类型
	}
}

// SetAuthMode 设置验证模式
func (s *SaiyouService) SetAuthMode(useDigitalSign bool, signType string) {
	s.useDigitalSign = useDigitalSign
	if useDigitalSign {
		if signType == "md5" || signType == "sha1" {
			s.signType = signType
		} else {
			s.signType = "md5" // 默认使用MD5
		}
	} else {
		s.signType = "" // 明文模式不需要签名类型
	}
}

// GetAuthMode 获取当前验证模式
func (s *SaiyouService) GetAuthMode() (useDigitalSign bool, signType string) {
	return s.useDigitalSign, s.signType
}

// GetFormat 获取当前响应格式
func (s *SaiyouService) GetFormat() string {
	return s.format
}

// SetConnectionPoolConfig 设置连接池配置并重新创建HTTP客户端
func (s *SaiyouService) SetConnectionPoolConfig(config *ConnectionPoolConfig) {
	if config == nil {
		config = DefaultConnectionPoolConfig()
	}
	s.poolConfig = config
	s.client = createHTTPClientWithPool(config)
}

// GetConnectionPoolConfig 获取当前连接池配置
func (s *SaiyouService) GetConnectionPoolConfig() *ConnectionPoolConfig {
	return s.poolConfig
}

// UpdateConnectionPoolConfig 更新连接池配置的特定参数
func (s *SaiyouService) UpdateConnectionPoolConfig(updater func(*ConnectionPoolConfig)) {
	if s.poolConfig == nil {
		s.poolConfig = DefaultConnectionPoolConfig()
	}
	updater(s.poolConfig)
	s.client = createHTTPClientWithPool(s.poolConfig)
}

// CloseIdleConnections 关闭空闲连接
func (s *SaiyouService) CloseIdleConnections() {
	if transport, ok := s.client.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}
}

// generateSignature 生成API签名
// SUBMAIL API 签名算法：
// 明文模式：直接返回 appkey
// 数字签名模式：
// 1. 添加 appid 和 timestamp，排除 tag 参数
// 2. 将所有参数按键名字典序排序
// 3. 按 key=value&key=value... 格式拼接
// 4. 拼接为：appid + appkey + signature_string + appid + appkey
// 5. 对整个字符串进行MD5或SHA1哈希
func (s *SaiyouService) generateSignature(params url.Values) string {
	// 明文验证模式：直接返回 appkey
	if !s.useDigitalSign {
		return s.AppKey
	}

	// 数字签名模式
	// 添加必要的公共参数
	params.Set("appid", s.AppID)
	params.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))

	// 获取所有参数键并按字典序排序，排除 signature 和 tag 参数
	keys := make([]string, 0, len(params))
	for key := range params {
		// tag 参数不参与签名计算
		if key != "signature" && key != "tag" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	// 构建签名字符串：key=value&key=value...
	var signParts []string
	for _, key := range keys {
		value := params.Get(key)
		signParts = append(signParts, key+"="+value)
	}

	// 拼接参数字符串
	signatureString := strings.Join(signParts, "&")

	// 按照SUBMAIL官方要求：appid + appkey + signature_string + appid + appkey
	finalString := s.AppID + s.AppKey + signatureString + s.AppID + s.AppKey

	// 根据签名类型计算哈希
	if s.signType == "sha1" {
		hash := sha1.Sum([]byte(finalString))
		return hex.EncodeToString(hash[:])
	} else {
		// 默认使用MD5
		hash := md5.Sum([]byte(finalString))
		return hex.EncodeToString(hash[:])
	}
}

// generateSignatureWithTimestamp 使用指定时间戳生成API签名
func (s *SaiyouService) generateSignatureWithTimestamp(params url.Values, timestamp int64) string {
	// 明文验证模式：直接返回 appkey
	if !s.useDigitalSign {
		return s.AppKey
	}

	// 数字签名模式
	// 添加必要的公共参数
	params.Set("appid", s.AppID)
	params.Set("timestamp", strconv.FormatInt(timestamp, 10))

	// 获取所有参数键并按字典序排序，排除 signature 和 tag 参数
	keys := make([]string, 0, len(params))
	for key := range params {
		// tag 参数不参与签名计算
		if key != "signature" && key != "tag" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	// 构建签名字符串：key=value&key=value...
	var signParts []string
	for _, key := range keys {
		value := params.Get(key)
		signParts = append(signParts, key+"="+value)
	}

	// 拼接参数字符串
	signatureString := strings.Join(signParts, "&")

	// 按照SUBMAIL官方要求：appid + appkey + signature_string + appid + appkey
	finalString := s.AppID + s.AppKey + signatureString + s.AppID + s.AppKey

	// 根据签名类型计算哈希
	if s.signType == "sha1" {
		hash := sha1.Sum([]byte(finalString))
		return hex.EncodeToString(hash[:])
	} else {
		// 默认使用MD5
		hash := md5.Sum([]byte(finalString))
		return hex.EncodeToString(hash[:])
	}
}

// ValidateSignature 验证签名（用于调试）
// 返回生成的签名字符串和计算的签名，用于排查签名问题
func (s *SaiyouService) ValidateSignature(params url.Values) (signString, signature string) {
	// 复制参数以避免修改原始数据
	paramsCopy := make(url.Values)
	for k, v := range params {
		paramsCopy[k] = v
	}

	// 明文验证模式
	if !s.useDigitalSign {
		return "明文模式，无需签名字符串", s.AppKey
	}

	// 数字签名模式
	// 添加必要的公共参数
	paramsCopy.Set("appid", s.AppID)
	paramsCopy.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))

	// 获取所有参数键并按字典序排序，排除 signature 和 tag 参数
	keys := make([]string, 0, len(paramsCopy))
	for key := range paramsCopy {
		// tag 参数不参与签名计算
		if key != "signature" && key != "tag" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	// 构建签名字符串
	var signParts []string
	for _, key := range keys {
		value := paramsCopy.Get(key)
		signParts = append(signParts, key+"="+value)
	}

	// 拼接参数字符串
	signatureString := strings.Join(signParts, "&")

	// 按照SUBMAIL官方要求：appid + appkey + signature_string + appid + appkey
	signString = s.AppID + s.AppKey + signatureString + s.AppID + s.AppKey

	// 根据签名类型计算哈希
	if s.signType == "sha1" {
		hash := sha1.Sum([]byte(signString))
		signature = hex.EncodeToString(hash[:])
	} else {
		// 默认使用MD5
		hash := md5.Sum([]byte(signString))
		signature = hex.EncodeToString(hash[:])
	}

	return signString, signature
}

// SyncServerTime 手动同步服务器时间
// 当遇到时间相关错误时，可以主动调用此方法同步时间
func (s *SaiyouService) SyncServerTime() (int64, error) {
	timestampResp, err := s.GetTimestamp()
	if err != nil {
		return 0, fmt.Errorf("获取服务器时间戳失败: %w", err)
	}

	serverTimestamp, err := s.extractTimestampFromResponse(timestampResp)
	if err != nil {
		return 0, fmt.Errorf("解析服务器时间戳失败: %w", err)
	}

	return serverTimestamp, nil
}

// GetTimeOffset 获取本地时间与服务器时间的偏移量（秒）
// 正值表示本地时间快于服务器时间，负值表示本地时间慢于服务器时间
func (s *SaiyouService) GetTimeOffset() (int64, error) {
	serverTime, err := s.SyncServerTime()
	if err != nil {
		return 0, err
	}

	localTime := time.Now().Unix()
	return localTime - serverTime, nil
}

// SMSBatchSendRequest 短信批量群发请求结构
type SMSBatchSendRequest struct {
	To      []string `json:"to"`                // 收件人手机号码列表 (必填)
	Text    string   `json:"text"`              // 短信正文 (必填)
	Project string   `json:"project,omitempty"` // 项目标记
	Tag     string   `json:"tag,omitempty"`     // 自定义标签
}

// SMSBatchXSendRequest 短信批量模板群发请求结构
type SMSBatchXSendRequest struct {
	To      []string          `json:"to"`             // 收件人手机号码列表 (必填)
	Project string            `json:"project"`        // 短信模板标记 (必填)
	Vars    map[string]string `json:"vars,omitempty"` // 文本变量
	Tag     string            `json:"tag,omitempty"`  // 自定义标签
}

// SMSUnionSendRequest 国内短信与国际短信联合发送请求结构
type SMSUnionSendRequest struct {
	To      string `json:"to"`                // 收件人手机号码 (必填)
	Text    string `json:"text"`              // 短信正文 (必填)
	Project string `json:"project,omitempty"` // 项目标记
	Tag     string `json:"tag,omitempty"`     // 自定义标签
	Country string `json:"country,omitempty"` // 国家代码，不传默认为中国
}

// SMSAddressBookRequest 短信地址簿管理请求结构
type SMSAddressBookRequest struct {
	Action  string `json:"action"`            // 操作类型: subscribe/unsubscribe
	To      string `json:"to"`                // 手机号码
	Project string `json:"project,omitempty"` // 项目标记
}

// doRequest 执行HTTP请求
func (s *SaiyouService) doRequest(endpoint string, params url.Values) (*APIResponse, error) {
	return s.doRequestWithRetry(endpoint, params, false)
}

// doRequestWithRetry 执行HTTP请求，支持时间同步重试
func (s *SaiyouService) doRequestWithRetry(endpoint string, params url.Values, isRetry bool) (*APIResponse, error) {
	// 数字签名模式需要添加 sign_type 参数
	if s.useDigitalSign {
		params.Set("sign_type", s.signType)
	}

	// 生成签名
	signature := s.generateSignature(params)
	params.Set("signature", signature)

	// 构建请求URL，根据format添加适当的扩展名
	var requestURL string
	if s.format == FormatXML {
		requestURL = s.BaseURL + endpoint + ".xml"
	} else {
		requestURL = s.BaseURL + endpoint + ".json"
	}

	// 创建POST请求
	req, err := http.NewRequest("POST", requestURL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "SUBMAIL-GO-SDK/1.0")

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析JSON响应
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应JSON失败: %w", err)
	}

	// 统一错误处理：若非 success，返回带中文描述的错误
	if apiErr := FromAPIResponse(&apiResp); apiErr != nil {
		// 对于时间相关错误，且不是重试请求，触发时间同步重试
		if !isRetry && IsTimeRelatedError(apiErr.Code) {
			return s.retryWithServerTimestamp(endpoint, params)
		}
		return nil, apiErr
	}

	return &apiResp, nil
}

// retryWithServerTimestamp 使用服务器时间戳重试请求
func (s *SaiyouService) retryWithServerTimestamp(endpoint string, params url.Values) (*APIResponse, error) {
	// 获取服务器时间戳
	timestampResp, err := s.GetTimestamp()
	if err != nil {
		return nil, fmt.Errorf("获取服务器时间戳失败: %w", err)
	}

	// 从响应中提取时间戳
	serverTimestamp, err := s.extractTimestampFromResponse(timestampResp)
	if err != nil {
		return nil, fmt.Errorf("解析服务器时间戳失败: %w", err)
	}

	// 使用服务器时间戳重新生成签名
	// 注意：这里需要复制参数，避免修改原始参数
	paramsCopy := make(url.Values)
	for k, v := range params {
		paramsCopy[k] = v
	}
	// 移除之前的签名
	paramsCopy.Del("signature")

	// 使用服务器时间戳生成新签名
	newSignature := s.generateSignatureWithTimestamp(paramsCopy, serverTimestamp)
	paramsCopy.Set("signature", newSignature)

	// 重新发送请求
	return s.doRequestWithRetry(endpoint, paramsCopy, true)
}

// extractTimestampFromResponse 从响应中提取时间戳
func (s *SaiyouService) extractTimestampFromResponse(resp *APIResponse) (int64, error) {
	// 尝试从不同的字段中提取时间戳
	if resp.Timestamp > 0 {
		return resp.Timestamp, nil
	}

	// 如果响应中有data字段，尝试从中提取
	if resp.Data != nil {
		if timestampStr, ok := resp.Data["timestamp"].(string); ok {
			return strconv.ParseInt(timestampStr, 10, 64)
		}
		if timestampFloat, ok := resp.Data["timestamp"].(float64); ok {
			return int64(timestampFloat), nil
		}
	}

	// 如果都没有，返回当前时间
	return time.Now().Unix(), nil
}

// SendSMS 发送短信
func (s *SaiyouService) SendSMS(req *SMSRequest) (*APIResponse, error) {
	if req.To == "" {
		return nil, fmt.Errorf("收件人手机号不能为空")
	}
	if req.Text == "" {
		return nil, fmt.Errorf("短信内容不能为空")
	}

	// 构建请求参数
	params := url.Values{}
	params.Set("to", req.To)
	params.Set("text", req.Text)

	if req.Project != "" {
		params.Set("project", req.Project)
	}

	if req.Tag != "" {
		params.Set("tag", req.Tag)
	}

	// 处理变量
	if len(req.Vars) > 0 {
		varsJSON, err := json.Marshal(req.Vars)
		if err != nil {
			return nil, fmt.Errorf("序列化变量失败: %w", err)
		}
		params.Set("vars", string(varsJSON))
	}

	return s.doRequest(SMSEndpoint, params)
}

// SendSMSTemplate 发送模板短信
func (s *SaiyouService) SendSMSTemplate(req *SMSXRequest) (*APIResponse, error) {
	if req.To == "" {
		return nil, fmt.Errorf("收件人手机号不能为空")
	}
	if req.Project == "" {
		return nil, fmt.Errorf("短信模板标记不能为空")
	}

	// 构建请求参数
	params := url.Values{}
	params.Set("to", req.To)
	params.Set("project", req.Project)

	if req.Tag != "" {
		params.Set("tag", req.Tag)
	}

	// 处理变量
	if len(req.Vars) > 0 {
		varsJSON, err := json.Marshal(req.Vars)
		if err != nil {
			return nil, fmt.Errorf("序列化变量失败: %w", err)
		}
		params.Set("vars", string(varsJSON))
	}

	return s.doRequest(SMSXEndpoint, params)
}

// GetBalance 查询账户余额
func (s *SaiyouService) GetBalance() (*BalanceResponse, error) {
	params := url.Values{}

	// 数字签名模式需要添加 sign_type 参数
	if s.useDigitalSign {
		params.Set("sign_type", s.signType)
	}

	// 生成签名
	signature := s.generateSignature(params)
	params.Set("signature", signature)

	// 构建请求URL，根据format添加适当的扩展名
	var requestURL string
	if s.format == FormatXML {
		requestURL = s.BaseURL + BalanceEndpoint + ".xml"
	} else {
		requestURL = s.BaseURL + BalanceEndpoint + ".json"
	}

	// 创建POST请求
	req, err := http.NewRequest("POST", requestURL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "SUBMAIL-GO-SDK/1.0")

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析JSON响应
	var balanceResp BalanceResponse
	if err := json.Unmarshal(body, &balanceResp); err != nil {
		return nil, fmt.Errorf("解析响应JSON失败: %w", err)
	}

	return &balanceResp, nil
}

// SendSMSMulti 发送短信一对多
func (s *SaiyouService) SendSMSMulti(req *SMSMultisendRequest) (*MultiSendResponse, error) {
	if len(req.Multi) == 0 {
		return nil, fmt.Errorf("联系人列表不能为空")
	}
	if req.Text == "" {
		return nil, fmt.Errorf("短信内容不能为空")
	}

	// 构建请求参数
	params := url.Values{}
	params.Set("text", req.Text)

	if req.Project != "" {
		params.Set("project", req.Project)
	}

	// 序列化联系人列表
	multiJSON, err := json.Marshal(req.Multi)
	if err != nil {
		return nil, fmt.Errorf("序列化联系人列表失败: %w", err)
	}
	params.Set("multi", string(multiJSON))

	// 发送请求
	return s.doMultiRequest(SMSMultisendEndpoint, params)
}

// SendSMSMultiTemplate 发送短信模板一对多
func (s *SaiyouService) SendSMSMultiTemplate(req *SMSMultiXSendRequest) (*MultiSendResponse, error) {
	if len(req.Multi) == 0 {
		return nil, fmt.Errorf("联系人列表不能为空")
	}
	if req.Project == "" {
		return nil, fmt.Errorf("短信模板标记不能为空")
	}

	// 构建请求参数
	params := url.Values{}
	params.Set("project", req.Project)

	// 序列化联系人列表
	multiJSON, err := json.Marshal(req.Multi)
	if err != nil {
		return nil, fmt.Errorf("序列化联系人列表失败: %w", err)
	}
	params.Set("multi", string(multiJSON))

	// 发送请求
	return s.doMultiRequest(SMSMultiXEndpoint, params)
}

// doMultiRequest 执行一对多发送HTTP请求
func (s *SaiyouService) doMultiRequest(endpoint string, params url.Values) (*MultiSendResponse, error) {
	// 数字签名模式需要添加 sign_type 参数
	if s.useDigitalSign {
		params.Set("sign_type", s.signType)
	}

	// 生成签名
	signature := s.generateSignature(params)
	params.Set("signature", signature)

	// 构建请求URL，根据format添加适当的扩展名
	var requestURL string
	if s.format == FormatXML {
		requestURL = s.BaseURL + endpoint + ".xml"
	} else {
		requestURL = s.BaseURL + endpoint + ".json"
	}

	// 创建POST请求
	req, err := http.NewRequest("POST", requestURL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "SUBMAIL-GO-SDK/1.0")

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析JSON响应
	var multiResp MultiSendResponse
	if err := json.Unmarshal(body, &multiResp); err != nil {
		return nil, fmt.Errorf("解析响应JSON失败: %w", err)
	}

	return &multiResp, nil
}

// GetSMSTemplates 获取短信模板列表
func (s *SaiyouService) GetSMSTemplates() (*APIResponse, error) {
	params := url.Values{}
	params.Set("action", "get")

	return s.doRequest(SMSTemplateEndpoint, params)
}

// CreateSMSTemplate 创建短信模板
func (s *SaiyouService) CreateSMSTemplate(sms string, templateType int) (*APIResponse, error) {
	if sms == "" {
		return nil, fmt.Errorf("短信内容不能为空")
	}

	params := url.Values{}
	params.Set("action", "post")
	params.Set("sms", sms)
	params.Set("type", strconv.Itoa(templateType))

	return s.doRequest(SMSTemplateEndpoint, params)
}

// UpdateSMSTemplate 更新短信模板
func (s *SaiyouService) UpdateSMSTemplate(templateID, sms string, templateType int) (*APIResponse, error) {
	if templateID == "" {
		return nil, fmt.Errorf("模板ID不能为空")
	}
	if sms == "" {
		return nil, fmt.Errorf("短信内容不能为空")
	}

	params := url.Values{}
	params.Set("action", "put")
	params.Set("template_id", templateID)
	params.Set("sms", sms)
	params.Set("type", strconv.Itoa(templateType))

	return s.doRequest(SMSTemplateEndpoint, params)
}

// DeleteSMSTemplate 删除短信模板
func (s *SaiyouService) DeleteSMSTemplate(templateID string) (*APIResponse, error) {
	if templateID == "" {
		return nil, fmt.Errorf("模板ID不能为空")
	}

	params := url.Values{}
	params.Set("action", "delete")
	params.Set("template_id", templateID)

	return s.doRequest(SMSTemplateEndpoint, params)
}

// GetSMSReports 获取短信分析报告
func (s *SaiyouService) GetSMSReports(req *SMSReportsRequest) (*APIResponse, error) {
	params := url.Values{}

	if req.Project != "" {
		params.Set("project", req.Project)
	}
	if req.StartDate != "" {
		params.Set("start_date", req.StartDate)
	}
	if req.EndDate != "" {
		params.Set("end_date", req.EndDate)
	}

	return s.doRequest(SMSReportsEndpoint, params)
}

// GetSMSLog 查询短信历史明细
func (s *SaiyouService) GetSMSLog(req *SMSLogRequest) (*APIResponse, error) {
	params := url.Values{}

	if req.Project != "" {
		params.Set("project", req.Project)
	}
	if req.StartDate != "" {
		params.Set("start_date", req.StartDate)
	}
	if req.EndDate != "" {
		params.Set("end_date", req.EndDate)
	}
	if req.Offset > 0 {
		params.Set("offset", strconv.Itoa(req.Offset))
	}
	if req.Limit > 0 {
		params.Set("limit", strconv.Itoa(req.Limit))
	}

	return s.doRequest(SMSLogEndpoint, params)
}

// GetSMSMO 查询短信上行
func (s *SaiyouService) GetSMSMO(req *SMSMORequest) (*APIResponse, error) {
	params := url.Values{}

	if req.StartDate != "" {
		params.Set("start_date", req.StartDate)
	}
	if req.EndDate != "" {
		params.Set("end_date", req.EndDate)
	}
	if req.Offset > 0 {
		params.Set("offset", strconv.Itoa(req.Offset))
	}
	if req.Limit > 0 {
		params.Set("limit", strconv.Itoa(req.Limit))
	}

	return s.doRequest(SMSMOEndpoint, params)
}

// GetTimestamp 获取服务器时间戳
func (s *SaiyouService) GetTimestamp() (*APIResponse, error) {
	params := url.Values{}
	return s.doRequest(TimestampEndpoint, params)
}

// GetStatus 获取服务器状态
func (s *SaiyouService) GetStatus() (*APIResponse, error) {
	params := url.Values{}
	return s.doRequest(StatusEndpoint, params)
}

// SendSMSBatch 短信批量群发
func (s *SaiyouService) SendSMSBatch(req *SMSBatchSendRequest) (*APIResponse, error) {
	if len(req.To) == 0 {
		return nil, fmt.Errorf("收件人手机号列表不能为空")
	}
	if req.Text == "" {
		return nil, fmt.Errorf("短信内容不能为空")
	}

	// 构建请求参数
	params := url.Values{}

	// 将手机号列表转换为JSON字符串
	toJSON, err := json.Marshal(req.To)
	if err != nil {
		return nil, fmt.Errorf("序列化手机号列表失败: %w", err)
	}
	params.Set("to", string(toJSON))

	params.Set("text", req.Text)

	if req.Project != "" {
		params.Set("project", req.Project)
	}

	if req.Tag != "" {
		params.Set("tag", req.Tag)
	}

	return s.doRequest(SMSBatchSendEndpoint, params)
}

// SendSMSBatchTemplate 短信批量模板群发
func (s *SaiyouService) SendSMSBatchTemplate(req *SMSBatchXSendRequest) (*APIResponse, error) {
	if len(req.To) == 0 {
		return nil, fmt.Errorf("收件人手机号列表不能为空")
	}
	if req.Project == "" {
		return nil, fmt.Errorf("短信模板标记不能为空")
	}

	// 构建请求参数
	params := url.Values{}

	// 将手机号列表转换为JSON字符串
	toJSON, err := json.Marshal(req.To)
	if err != nil {
		return nil, fmt.Errorf("序列化手机号列表失败: %w", err)
	}
	params.Set("to", string(toJSON))

	params.Set("project", req.Project)

	if req.Tag != "" {
		params.Set("tag", req.Tag)
	}

	// 处理变量
	if len(req.Vars) > 0 {
		varsJSON, err := json.Marshal(req.Vars)
		if err != nil {
			return nil, fmt.Errorf("序列化变量失败: %w", err)
		}
		params.Set("vars", string(varsJSON))
	}

	return s.doRequest(SMSBatchXEndpoint, params)
}

// SendSMSUnion 国内短信与国际短信联合发送
func (s *SaiyouService) SendSMSUnion(req *SMSUnionSendRequest) (*APIResponse, error) {
	if req.To == "" {
		return nil, fmt.Errorf("收件人手机号不能为空")
	}
	if req.Text == "" {
		return nil, fmt.Errorf("短信内容不能为空")
	}

	// 构建请求参数
	params := url.Values{}
	params.Set("to", req.To)
	params.Set("text", req.Text)

	if req.Project != "" {
		params.Set("project", req.Project)
	}

	if req.Tag != "" {
		params.Set("tag", req.Tag)
	}

	if req.Country != "" {
		params.Set("country", req.Country)
	}

	return s.doRequest(SMSUnionSendEndpoint, params)
}

// SubscribeSMS 短信订阅
func (s *SaiyouService) SubscribeSMS(to, project string) (*APIResponse, error) {
	if to == "" {
		return nil, fmt.Errorf("手机号不能为空")
	}

	params := url.Values{}
	params.Set("action", "subscribe")
	params.Set("to", to)

	if project != "" {
		params.Set("project", project)
	}

	return s.doRequest(AddressBookEndpoint, params)
}

// UnsubscribeSMS 短信退订
func (s *SaiyouService) UnsubscribeSMS(to, project string) (*APIResponse, error) {
	if to == "" {
		return nil, fmt.Errorf("手机号不能为空")
	}

	params := url.Values{}
	params.Set("action", "unsubscribe")
	params.Set("to", to)

	if project != "" {
		params.Set("project", project)
	}

	return s.doRequest(AddressBookEndpoint, params)
}
