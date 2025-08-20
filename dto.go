package submail

import "mime/multipart"

// 基础响应结构
type BaseResponse struct {
	Status string `json:"status" form:"status" xml:"status"`
	Code   int    `json:"code,omitempty" form:"code" xml:"code"`
	Msg    string `json:"msg,omitempty" form:"msg" xml:"msg"`
}

// 短信发送请求
type SMSSendRequest struct {
	To      string `json:"to" form:"to" xml:"to"`                // 收件人手机号码
	Content string `json:"content" form:"content" xml:"content"` // 短信正文（支持@var(key)和@date()变量）
	Tag     string `json:"tag,omitempty" form:"tag" xml:"tag"`   // 自定义标签，最多32个字符
}

// 短信模板发送请求
type SMSXSendRequest struct {
	To           string            `json:"to" form:"to" xml:"to"`                                            // 收件人手机号码
	Project      string            `json:"project" form:"project" xml:"project"`                             // 短信模板ID
	Vars         map[string]string `json:"vars,omitempty" form:"vars" xml:"vars"`                            // 模板变量
	SMSSignature string            `json:"sms_signature,omitempty" form:"sms_signature" xml:"sms_signature"` // 自定义短信签名（v4.002新增）
	Tag          string            `json:"tag,omitempty" form:"tag" xml:"tag"`                               // 自定义标签，最多32个字符
}

// 短信一对多发送请求
type SMSMultiSendRequest struct {
	Content string         `json:"content" form:"content" xml:"content"` // 短信正文（支持@var(key)变量）
	Multi   []SMSMultiItem `json:"multi" form:"multi" xml:"multi"`       // 收件人列表
	Tag     string         `json:"tag,omitempty" form:"tag" xml:"tag"`   // 自定义标签
}

type SMSMultiItem struct {
	To   string            `json:"to" form:"to" xml:"to"`                 // 收件人手机号码
	Vars map[string]string `json:"vars,omitempty" form:"vars" xml:"vars"` // 文本变量
}

// 短信模板一对多发送请求
type SMSMultiXSendRequest struct {
	Multi        []SMSMultiXItem `json:"multi" form:"multi" xml:"multi"`                                   // 收件人列表
	Project      string          `json:"project" form:"project" xml:"project"`                             // 短信模板ID
	SMSSignature string          `json:"sms_signature,omitempty" form:"sms_signature" xml:"sms_signature"` // 自定义短信签名（v4.002新增）
	Tag          string          `json:"tag,omitempty" form:"tag" xml:"tag"`                               // 自定义标签
}

type SMSMultiXItem struct {
	To           string            `json:"to" form:"to" xml:"to"`                                            // 收件人手机号码
	Vars         map[string]string `json:"vars,omitempty" form:"vars" xml:"vars"`                            // 模板变量
	SMSSignature string            `json:"sms_signature,omitempty" form:"sms_signature" xml:"sms_signature"` // 自定义短信签名（单独设置）
}

// 短信批量群发请求
type SMSBatchSendRequest struct {
	Content string `json:"content" form:"content" xml:"content"` // 短信正文（支持@var(key)和@date()变量）
	To      string `json:"to" form:"to" xml:"to"`                // 收件人手机号码，多个号码用逗号分隔
	Tag     string `json:"tag,omitempty" form:"tag" xml:"tag"`   // 自定义标签
}

// 短信批量模板群发请求
type SMSBatchXSendRequest struct {
	Project      string            `json:"project" form:"project" xml:"project"`                             // 短信模板ID
	To           string            `json:"to" form:"to" xml:"to"`                                            // 收件人手机号码，多个号码用逗号分隔
	Vars         map[string]string `json:"vars,omitempty" form:"vars" xml:"vars"`                            // 模板变量
	SMSSignature string            `json:"sms_signature,omitempty" form:"sms_signature" xml:"sms_signature"` // 自定义短信签名（v4.002新增）
	Tag          string            `json:"tag,omitempty" form:"tag" xml:"tag"`                               // 自定义标签
}

// 短信联合发送请求
type SMSUnionSendRequest struct {
	To                          string `json:"to" form:"to" xml:"to"`                                                                                               // 收件人手机号码（支持国内11位和国际E164格式）
	Content                     string `json:"content" form:"content" xml:"content"`                                                                                // 国内短信正文
	InterAppID                  string `json:"inter_appid" form:"inter_appid" xml:"inter_appid"`                                                                    // 国际短信AppID
	InterSignature              string `json:"inter_signature" form:"inter_signature" xml:"inter_signature"`                                                        // 国际短信应用密钥
	InterContent                string `json:"inter_content,omitempty" form:"inter_content" xml:"inter_content"`                                                    // 国际短信正文（可选）
	IntersmsVerifyCodeTransform string `json:"intersms_verify_code_transform,omitempty" form:"intersms_verify_code_transform" xml:"intersms_verify_code_transform"` // 是否提取验证码替换@var(code)
	Tag                         string `json:"tag,omitempty" form:"tag" xml:"tag"`                                                                                  // 自定义标签
}

// 短信发送响应
type SMSSendResponse struct {
	BaseResponse
	SendID string `json:"send_id,omitempty" form:"send_id" xml:"send_id"` // 发送ID
	Fee    int    `json:"fee,omitempty" form:"fee" xml:"fee"`             // 扣费金额
	Sms    int    `json:"sms,omitempty" form:"sms" xml:"sms"`             // 短信条数
}

// 短信多条发送响应（用于MultiSend等批量发送API）
type SMSMultiSendResponse []SMSSendResult

// 单条短信发送结果
type SMSSendResult struct {
	Status string `json:"status" form:"status" xml:"status"`              // 请求状态：success/error
	To     string `json:"to,omitempty" form:"to" xml:"to"`                // 收件人手机号码
	SendID string `json:"send_id,omitempty" form:"send_id" xml:"send_id"` // 发送ID
	Fee    int    `json:"fee,omitempty" form:"fee" xml:"fee"`             // 扣费金额
	Code   int    `json:"code,omitempty" form:"code" xml:"code"`          // 错误代码（失败时）
	Msg    string `json:"msg,omitempty" form:"msg" xml:"msg"`             // 错误信息（失败时）
}

// 短信批量发送响应（用于BatchSend和BatchXSend）
type SMSBatchSendResponse struct {
	BaseResponse
	BatchList string          `json:"batchlist,omitempty" form:"batchlist" xml:"batchlist"` // 批量任务ID
	TotalFee  int             `json:"total_fee,omitempty" form:"total_fee" xml:"total_fee"` // 总计费条数
	Responses []SMSSendResult `json:"responses,omitempty" form:"responses" xml:"responses"` // 各号码发送结果
}

// 短信模板管理相关结构

// 短信模板查询请求（GET方法）
type SMSTemplateGetRequest struct {
	TemplateID string `json:"template_id,omitempty" form:"template_id" xml:"template_id"` // 模板ID（获取单个模板时使用）
	Offset     int    `json:"offset,omitempty" form:"offset" xml:"offset"`                // 数据偏移指针，默认0
}

// 短信模板创建请求（POST方法）
type SMSTemplateCreateRequest struct {
	SMSTitle     string `json:"sms_title,omitempty" form:"sms_title" xml:"sms_title"`   // 模板标题
	SMSSignature string `json:"sms_signature" form:"sms_signature" xml:"sms_signature"` // 短信模板签名（必填）
	SMSContent   string `json:"sms_content" form:"sms_content" xml:"sms_content"`       // 短信正文（必填）
}

// 短信模板更新请求（PUT方法）
type SMSTemplateUpdateRequest struct {
	TemplateID   string `json:"template_id" form:"template_id" xml:"template_id"`       // 需要更新的模板ID（必填）
	SMSTitle     string `json:"sms_title,omitempty" form:"sms_title" xml:"sms_title"`   // 模板标题
	SMSSignature string `json:"sms_signature" form:"sms_signature" xml:"sms_signature"` // 短信模板签名（必填）
	SMSContent   string `json:"sms_content" form:"sms_content" xml:"sms_content"`       // 短信正文（必填）
}

// 短信模板删除请求（DELETE方法）
type SMSTemplateDeleteRequest struct {
	TemplateID string `json:"template_id" form:"template_id" xml:"template_id"` // 需要删除的模板ID（必填）
}

// 短信模板信息
type SMSTemplate struct {
	TemplateID                string `json:"template_id" form:"template_id" xml:"template_id"`                                                 // 模板ID
	SMSTitle                  string `json:"sms_title" form:"sms_title" xml:"sms_title"`                                                       // 短信标题
	SMSSignature              string `json:"sms_signature" form:"sms_signature" xml:"sms_signature"`                                           // 短信签名
	SMSContent                string `json:"sms_content" form:"sms_content" xml:"sms_content"`                                                 // 短信内容
	AddDate                   int64  `json:"add_date" form:"add_date" xml:"add_date"`                                                          // 创建时间（UNIX时间戳）
	EditDate                  int64  `json:"edit_date" form:"edit_date" xml:"edit_date"`                                                       // 编辑时间（UNIX时间戳）
	TemplateStatus            string `json:"template_status" form:"template_status" xml:"template_status"`                                     // 模板状态：0=未提交、1=审核中、2=通过、3=未通过
	TemplateStatusDescription string `json:"template_status_description" form:"template_status_description" xml:"template_status_description"` // 模板状态描述
	TemplateRejectReason      string `json:"template_reject_reson,omitempty" form:"template_reject_reson" xml:"template_reject_reson"`         // 驳回原因（注意：API返回的字段名是reson不是reason）
}

// 短信模板查询响应
type SMSTemplateGetResponse struct {
	BaseResponse
	StartRow  int           `json:"start_row,omitempty" form:"start_row" xml:"start_row"` // 起始行
	EndRow    int           `json:"end_row,omitempty" form:"end_row" xml:"end_row"`       // 结束行
	Templates []SMSTemplate `json:"templates,omitempty" form:"templates" xml:"templates"` // 模板列表
}

// 短信模板创建响应
type SMSTemplateCreateResponse struct {
	BaseResponse
	TemplateID string `json:"template_id,omitempty" form:"template_id" xml:"template_id"` // 创建成功返回的模板ID
}

// 短信模板操作响应（更新、删除）
type SMSTemplateOperationResponse struct {
	BaseResponse
}

// 短信分析报告请求
type SMSReportsRequest struct {
	StartDate int64 `json:"start_date,omitempty" form:"start_date" xml:"start_date"` // 报告开始时间（UNIX时间戳）
	EndDate   int64 `json:"end_date,omitempty" form:"end_date" xml:"end_date"`       // 报告结束时间（UNIX时间戳）
}

// 短信分析报告响应
type SMSReportsResponse struct {
	BaseResponse
	StartDate string              `json:"start_date,omitempty" form:"start_date" xml:"start_date"` // 开始日期
	EndDate   string              `json:"end_date,omitempty" form:"end_date" xml:"end_date"`       // 结束日期
	Overview  SMSReportOverview   `json:"overview,omitempty" form:"overview" xml:"overview"`       // 概览数据
	Timeline  []SMSReportTimeline `json:"timeline,omitempty" form:"timeline" xml:"timeline"`       // 时间线数据
}

// 短信报告概览
type SMSReportOverview struct {
	Request               int                `json:"request" form:"request" xml:"request"`                                                           // API请求总数
	Deliveryed            int                `json:"deliveryed" form:"deliveryed" xml:"deliveryed"`                                                  // 成功数
	Dropped               int                `json:"dropped" form:"dropped" xml:"dropped"`                                                           // 失败数
	Fee                   int                `json:"fee" form:"fee" xml:"fee"`                                                                       // 计费数
	Operators             SMSReportOperators `json:"operators,omitempty" form:"operators" xml:"operators"`                                           // 运营商占比
	Location              SMSReportLocation  `json:"location,omitempty" form:"location" xml:"location"`                                              // 地区分类
	DroppedReasonAnalysis map[string]int     `json:"dropped_reason_analysis,omitempty" form:"dropped_reason_analysis" xml:"dropped_reason_analysis"` // 失败原因分析
}

// 运营商占比
type SMSReportOperators struct {
	ChinaMobile  int `json:"china_mobile" form:"china_mobile" xml:"china_mobile"`    // 移动
	ChinaUnicom  int `json:"china_unicom" form:"china_unicom" xml:"china_unicom"`    // 联通
	ChinaTelecom int `json:"china_telecom" form:"china_telecom" xml:"china_telecom"` // 电信
}

// 地区分类
type SMSReportLocation struct {
	Province map[string]int `json:"province,omitempty" form:"province" xml:"province"` // 省份统计
	Cities   map[string]int `json:"cities,omitempty" form:"cities" xml:"cities"`       // 城市统计
}

// 时间线数据
type SMSReportTimeline struct {
	Date   string                  `json:"date" form:"date" xml:"date"`       // 日期
	Report SMSReportTimelineDetail `json:"report" form:"report" xml:"report"` // 报告详情
}

// 时间线报告详情
type SMSReportTimelineDetail struct {
	Request    int `json:"request" form:"request" xml:"request"`          // API请求
	Deliveryed int `json:"deliveryed" form:"deliveryed" xml:"deliveryed"` // 成功数
	Dropped    int `json:"dropped" form:"dropped" xml:"dropped"`          // 失败数
	Fee        int `json:"fee" form:"fee" xml:"fee"`                      // 计费数
}

// 短信余额查询响应
type SMSBalanceResponse struct {
	BaseResponse
	Balance              string `json:"balance,omitempty" form:"balance" xml:"balance"`                                           // 通用类短信余额
	TransactionalBalance string `json:"transactional_balance,omitempty" form:"transactional_balance" xml:"transactional_balance"` // 事务类短信余额
}

// 短信余额日志查询请求
type SMSBalanceLogRequest struct {
	StartDate int64 `json:"start_date,omitempty" form:"start_date" xml:"start_date"` // 开始时间（UNIX时间戳）
	EndDate   int64 `json:"end_date,omitempty" form:"end_date" xml:"end_date"`       // 结束时间（UNIX时间戳）
}

// 短信余额日志响应
type SMSBalanceLogResponse struct {
	BaseResponse
	Data []SMSBalanceLogEntry `json:"data,omitempty" form:"data" xml:"data"` // 余额变更记录
}

// 短信余额日志条目
type SMSBalanceLogEntry struct {
	Datetime              string `json:"datetime" form:"datetime" xml:"datetime"`                                                        // 变更时间
	Message               string `json:"message" form:"message" xml:"message"`                                                           // 订单编号或变更说明
	TMessageAddCredits    string `json:"tmessage_add_credits,omitempty" form:"tmessage_add_credits" xml:"tmessage_add_credits"`          // 事务类短信增加余额
	TMessageAfterCredits  string `json:"tmessage_after_credits,omitempty" form:"tmessage_after_credits" xml:"tmessage_after_credits"`    // 事务类短信增加后账户余额
	TMessageBeforeCredits string `json:"tmessage_before_credits,omitempty" form:"tmessage_before_credits" xml:"tmessage_before_credits"` // 事务类短信增加前账户余额
	MessageAddCredits     string `json:"message_add_credits,omitempty" form:"message_add_credits" xml:"message_add_credits"`             // 运营类短信增加余额
	MessageAfterCredits   string `json:"message_after_credits,omitempty" form:"message_after_credits" xml:"message_after_credits"`       // 运营类短信增加后账户余额
	MessageBeforeCredits  string `json:"message_before_credits,omitempty" form:"message_before_credits" xml:"message_before_credits"`    // 运营类短信增加前账户余额
}

// 短信历史明细请求
type SMSLogRequest struct {
	App       string `json:"app,omitempty" form:"app" xml:"app"`                      // 指定appid
	StartDate int64  `json:"start_date,omitempty" form:"start_date" xml:"start_date"` // 开始时间（UNIX时间戳）
	EndDate   int64  `json:"end_date,omitempty" form:"end_date" xml:"end_date"`       // 结束时间（UNIX时间戳）
	To        string `json:"to,omitempty" form:"to" xml:"to"`                         // 查询特定手机号码
	SendID    string `json:"send_id,omitempty" form:"send_id" xml:"send_id"`          // 查询特定Send ID
	SendList  string `json:"sendlist,omitempty" form:"sendlist" xml:"sendlist"`       // 查询特定发送任务
	Status    string `json:"status,omitempty" form:"status" xml:"status"`             // delivered或dropped
	Rows      int    `json:"rows,omitempty" form:"rows" xml:"rows"`                   // 返回数据行数
	Offset    int    `json:"offset,omitempty" form:"offset" xml:"offset"`             // 数据偏移值
}

// 短信历史明细响应
type SMSLogResponse struct {
	BaseResponse
	StartDate int64    `json:"start_date" form:"start_date" xml:"start_date"` // 查询开始日期
	EndDate   int64    `json:"end_date" form:"end_date" xml:"end_date"`       // 查询结束日期
	Total     int      `json:"total" form:"total" xml:"total"`                // 记录数
	Offset    int      `json:"offset" form:"offset" xml:"offset"`             // 数据偏移值
	Results   int      `json:"results" form:"results" xml:"results"`          // 每页行数
	Data      []SMSLog `json:"data" form:"data" xml:"data"`                   // 数据
}

type SMSLog struct {
	SendID        string `json:"sendID" form:"sendID" xml:"sendID"`                                   // Send ID
	To            string `json:"to" form:"to" xml:"to"`                                               // 手机号码
	AppID         string `json:"appid" form:"appid" xml:"appid"`                                      // AppID
	TemplateID    string `json:"template_id" form:"template_id" xml:"template_id"`                    // 模板ID
	SMSSignature  string `json:"sms_signature" form:"sms_signature" xml:"sms_signature"`              // 短信签名
	SMSContent    string `json:"sms_content" form:"sms_content" xml:"sms_content"`                    // 短信正文
	Fee           int    `json:"fee" form:"fee" xml:"fee"`                                            // 计费条数
	Status        string `json:"status" form:"status" xml:"status"`                                   // 发送状态：delivered=成功，dropped=失败，pending=未知
	ReportState   string `json:"report_state" form:"report_state" xml:"report_state"`                 // 运营商返回的实际状态
	DroppedReason string `json:"dropped_reason,omitempty" form:"dropped_reason" xml:"dropped_reason"` // 失败原因
	Location      string `json:"location" form:"location" xml:"location"`                             // 手机号归属地
	MobileType    string `json:"mobile_type" form:"mobile_type" xml:"mobile_type"`                    // 手机运营商
	IPAddress     string `json:"ip_address" form:"ip_address" xml:"ip_address"`                       // 发送IP
	SendAt        int64  `json:"send_at" form:"send_at" xml:"send_at"`                                // 请求时间
	SentAt        int64  `json:"sent_at" form:"sent_at" xml:"sent_at"`                                // 平台发送时间
	ReportAt      int64  `json:"report_at" form:"report_at" xml:"report_at"`                          // 运营商状态汇报时间
}

// 短信上行查询请求
type SMSMORequest struct {
	StartDate int64  `json:"start_date,omitempty" form:"start_date" xml:"start_date"` // 开始时间（UNIX时间戳）
	EndDate   int64  `json:"end_date,omitempty" form:"end_date" xml:"end_date"`       // 结束时间（UNIX时间戳）
	From      string `json:"from,omitempty" form:"from" xml:"from"`                   // 查询特定手机号码
	Rows      int    `json:"rows,omitempty" form:"rows" xml:"rows"`                   // 返回数据行数
	Offset    int    `json:"offset,omitempty" form:"offset" xml:"offset"`             // 数据偏移值
}

// 短信上行查询响应
type SMSMOResponse struct {
	BaseResponse
	StartDate int64   `json:"start_date" form:"start_date" xml:"start_date"` // 开始日期
	EndDate   int64   `json:"end_date" form:"end_date" xml:"end_date"`       // 结束日期
	Total     int     `json:"total" form:"total" xml:"total"`                // 查询总数
	Offset    int     `json:"offset" form:"offset" xml:"offset"`             // 数据偏移值
	Results   int     `json:"results" form:"results" xml:"results"`          // 返回结果数
	MO        []SMSMO `json:"mo" form:"mo" xml:"mo"`                         // 上行数据
}

type SMSMO struct {
	AppID      string `json:"appid" form:"appid" xml:"appid"`                   // AppID
	From       string `json:"from" form:"from" xml:"from"`                      // 回复手机号
	Content    string `json:"content" form:"content" xml:"content"`             // 回复正文
	ReplyAt    int64  `json:"reply_at" form:"reply_at" xml:"reply_at"`          // 回复时间
	SMSContent string `json:"sms_content" form:"sms_content" xml:"sms_content"` // 下行短信内容
	SendList   string `json:"sendlist" form:"sendlist" xml:"sendlist"`          // 批次号
}

// 服务时间戳响应
type ServiceTimestampResponse struct {
	BaseResponse
	Timestamp int64 `json:"timestamp,omitempty" form:"timestamp" xml:"timestamp"` // UNIX时间戳
}

// 服务状态响应
type ServiceStatusResponse struct {
	Status  string  `json:"status" form:"status" xml:"status"`    // 服务状态，如 "runing"
	Runtime float64 `json:"runtime" form:"runtime" xml:"runtime"` // 响应时间（秒）
}

// 地址簿订阅请求
type AddressBookSubscribeRequest struct {
	Address string            `json:"address" form:"address" xml:"address"`  // 联系人地址（手机号码）
	Tag     string            `json:"tag,omitempty" form:"tag" xml:"tag"`    // 标签
	Vars    map[string]string `json:"vars,omitempty" form:"vars" xml:"vars"` // 自定义变量
}

// 地址簿退订请求
type AddressBookUnsubscribeRequest struct {
	Address string `json:"address" form:"address" xml:"address"` // 联系人地址（手机号码）
}

// 地址簿响应
type AddressBookResponse struct {
	BaseResponse
	Address string `json:"address,omitempty" form:"address" xml:"address"` // 联系人地址
}

// 短信签名管理相关结构

// 短信签名查询请求
type SMSSignatureQueryRequest struct {
	TargetAppID  string `json:"target_appid,omitempty" form:"target_appid" xml:"target_appid"`    // 目标AppID
	SMSSignature string `json:"sms_signature,omitempty" form:"sms_signature" xml:"sms_signature"` // 要查询的短信签名
}

// 短信签名创建请求
type SMSSignatureCreateRequest struct {
	SMSSignature       string                 `form:"sms_signature" json:"sms_signature" xml:"sms_signature"`                         // 短信签名（可省略【】符号）
	Company            string                 `form:"company" json:"company" xml:"company"`                                           // 公司名称
	CompanyLisenceCode string                 `form:"company_lisence_code" json:"company_lisence_code" xml:"company_lisence_code"`    // 公司税号
	LegalName          string                 `form:"legal_name" json:"legal_name" xml:"legal_name"`                                  // 公司法人姓名
	Attachments        []multipart.FileHeader `form:"attachments" json:"attachments" xml:"attachments"`                               // 证明材料文件（可多个）
	AgentName          string                 `form:"agent_name" json:"agent_name" xml:"agent_name"`                                  // 责任人姓名
	AgentID            string                 `form:"agent_id" json:"agent_id" xml:"agent_id"`                                        // 责任人身份证号
	AgentMob           string                 `form:"agent_mob" json:"agent_mob" xml:"agent_mob"`                                     // 责任人手机号
	SourceType         int                    `form:"source_type,omitempty" json:"source_type,omitempty" xml:"source_type,omitempty"` // 材料类型：0=营业执照、1=商标、2=APP
	Contact            string                 `form:"contact,omitempty" json:"contact,omitempty" xml:"contact,omitempty"`             // 联系人手机号码
}

// 短信签名更新请求
type SMSSignatureUpdateRequest struct {
	SMSSignature       string                 `form:"sms_signature" json:"sms_signature" xml:"sms_signature"`                                                    // 短信签名（必填）
	Company            string                 `form:"company,omitempty" json:"company,omitempty" xml:"company,omitempty"`                                        // 公司名称
	CompanyLisenceCode string                 `form:"company_lisence_code,omitempty" json:"company_lisence_code,omitempty" xml:"company_lisence_code,omitempty"` // 公司税号
	LegalName          string                 `form:"legal_name,omitempty" json:"legal_name,omitempty" xml:"legal_name,omitempty"`                               // 公司法人姓名
	Attachments        []multipart.FileHeader `form:"attachments,omitempty" json:"attachments,omitempty" xml:"attachments,omitempty"`                            // 证明材料文件（可选，可多个）
	AgentName          string                 `form:"agent_name,omitempty" json:"agent_name,omitempty" xml:"agent_name,omitempty"`                               // 责任人姓名
	AgentID            string                 `form:"agent_id,omitempty" json:"agent_id,omitempty" xml:"agent_id,omitempty"`                                     // 责任人身份证号
	AgentMob           string                 `form:"agent_mob,omitempty" json:"agent_mob,omitempty" xml:"agent_mob,omitempty"`                                  // 责任人手机号
	SourceType         int                    `form:"source_type,omitempty" json:"source_type,omitempty" xml:"source_type,omitempty"`                            // 材料类型：0=营业执照、1=商标、2=APP
	Contact            string                 `form:"contact,omitempty" json:"contact,omitempty" xml:"contact,omitempty"`                                        // 联系人手机号码
}

// 短信签名删除请求
type SMSSignatureDeleteRequest struct {
	TargetAppID  string `json:"target_appid,omitempty" form:"target_appid" xml:"target_appid"` // 目标AppID
	SMSSignature string `json:"sms_signature" form:"sms_signature" xml:"sms_signature"`        // 要删除的短信签名
}

// 短信签名信息
type SMSSignatureInfo struct {
	AppID        string `json:"appid" form:"appid" xml:"appid"`                      // AppID
	SMSSignature string `json:"smsSignature" form:"smsSignature" xml:"smsSignature"` // 短信签名
	Status       int    `json:"status" form:"status" xml:"status"`                   // 状态：0=审核中、1=审核通过、其他=审核不通过
}

// 短信签名查询响应
type SMSSignatureQueryResponse struct {
	BaseResponse
	SMSSignatures []SMSSignatureInfo `json:"smsSignature,omitempty" form:"smsSignature" xml:"smsSignature"` // 签名列表
}

// 短信签名操作响应（创建、更新、删除）
type SMSSignatureOperationResponse struct {
	BaseResponse
}
