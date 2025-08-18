package submail

// APIResponse SUBMAIL API统一响应结构
type APIResponse struct {
	Status            string                 `json:"status"`                        // 请求状态 (success/error)
	Code              int                    `json:"code,omitempty"`                // 状态码
	Msg               string                 `json:"msg,omitempty"`                 // 响应消息
	SendID            string                 `json:"send_id,omitempty"`             // 发送ID
	Fee               int                    `json:"fee,omitempty"`                 // 消费费用
	Sms               int                    `json:"sms,omitempty"`                 // 短信条数
	To                string                 `json:"to,omitempty"`                  // 收件人
	Credits           string                 `json:"credits,omitempty"`             // 剩余积分
	Timestamp         int64                  `json:"timestamp,omitempty"`           // 时间戳
	JSONDecodingError int                    `json:"json_decoding_error,omitempty"` // JSON 解码错误子码（当 Code 为 307/309 时）
	Data              map[string]interface{} `json:"data,omitempty"`                // 其他数据
}

// SMSRequest 短信发送请求结构
type SMSRequest struct {
	To      string            `json:"to"`                // 收件人手机号码 (必填)
	Text    string            `json:"text"`              // 短信正文 (必填)
	Vars    map[string]string `json:"vars,omitempty"`    // 文本变量，用于模板短信
	Project string            `json:"project,omitempty"` // 项目标记，用于数据分析
	Tag     string            `json:"tag,omitempty"`     // 自定义标签
}

// SMSXRequest 短信模板发送请求结构
type SMSXRequest struct {
	To      string            `json:"to"`             // 收件人手机号码 (必填)
	Project string            `json:"project"`        // 短信模板标记 (必填)
	Vars    map[string]string `json:"vars,omitempty"` // 文本变量
	Tag     string            `json:"tag,omitempty"`  // 自定义标签
}

// BalanceResponse 余额查询响应结构
type BalanceResponse struct {
	Status  string `json:"status"`
	Balance string `json:"balance"`
}

// SMSMultisendRequest 短信一对多发送请求结构
type SMSMultisendRequest struct {
	Multi   []SMSMultiItem `json:"multi"`             // 联系人列表
	Text    string         `json:"text"`              // 短信正文
	Project string         `json:"project,omitempty"` // 项目标记
}

// SMSMultiItem 一对多发送中的联系人信息
type SMSMultiItem struct {
	To   string            `json:"to"`             // 收件人手机号
	Vars map[string]string `json:"vars,omitempty"` // 个性化变量
}

// SMSMultiXSendRequest 短信模板一对多发送请求结构
type SMSMultiXSendRequest struct {
	Multi   []SMSMultiItem `json:"multi"`   // 联系人列表
	Project string         `json:"project"` // 短信模板标记
}

// SMSTemplateRequest 短信模板管理请求结构
type SMSTemplateRequest struct {
	Action     string `json:"action"`                // 操作类型: get/post/put/delete
	TemplateID string `json:"template_id,omitempty"` // 模板ID
	SMS        string `json:"sms,omitempty"`         // 短信内容
	Type       int    `json:"type,omitempty"`        // 模板类型
}

// SMSReportsRequest 短信分析报告请求结构
type SMSReportsRequest struct {
	Project   string `json:"project,omitempty"`    // 项目标记
	StartDate string `json:"start_date,omitempty"` // 开始日期
	EndDate   string `json:"end_date,omitempty"`   // 结束日期
}

// SMSLogRequest 短信历史明细查询请求结构
type SMSLogRequest struct {
	Project   string `json:"project,omitempty"`    // 项目标记
	StartDate string `json:"start_date,omitempty"` // 开始日期
	EndDate   string `json:"end_date,omitempty"`   // 结束日期
	Offset    int    `json:"offset,omitempty"`     // 偏移量
	Limit     int    `json:"limit,omitempty"`      // 获取数据量
}

// SMSMORequest 短信上行查询请求结构
type SMSMORequest struct {
	StartDate string `json:"start_date,omitempty"` // 开始日期
	EndDate   string `json:"end_date,omitempty"`   // 结束日期
	Offset    int    `json:"offset,omitempty"`     // 偏移量
	Limit     int    `json:"limit,omitempty"`      // 获取数据量
}

// MultiSendResponse 一对多发送响应结构
type MultiSendResponse struct {
	Status string                `json:"status"`          // 请求状态
	Code   int                   `json:"code,omitempty"`  // 状态码
	Msg    string                `json:"msg,omitempty"`   // 响应消息
	Sends  []MultiSendResultItem `json:"sends,omitempty"` // 发送结果列表
}

// MultiSendResultItem 一对多发送结果项
type MultiSendResultItem struct {
	To     string `json:"to"`                // 收件人
	Status string `json:"status"`            // 发送状态
	SendID string `json:"send_id,omitempty"` // 发送ID
	Fee    int    `json:"fee,omitempty"`     // 消费费用
	Sms    int    `json:"sms,omitempty"`     // 短信条数
	Code   int    `json:"code,omitempty"`    // 错误码
	Msg    string `json:"msg,omitempty"`     // 错误信息
}
