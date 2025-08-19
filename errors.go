package submail

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// APIError 赛邮云API错误
type APIError struct {
	Code        int    `json:"code"`        // 错误代码
	Msg         string `json:"msg"`         // 错误信息
	Description string `json:"description"` // 错误描述
}

func (e *APIError) Error() string {
	return fmt.Sprintf("SubMail API Error %d: %s", e.Code, e.Msg)
}

// 错误代码常量定义
const (
	// 应用相关错误 (101-120)
	ErrIncorrectAppID             = 101 // 不正确的 APP ID
	ErrAppDisabled                = 102 // 此应用已被禁用
	ErrDeveloperNotAvailable      = 103 // 未启用的开发者
	ErrDeveloperNotVerified       = 104 // 此开发者未通过验证
	ErrAccountExpired             = 105 // 此账户已过期
	ErrAccountDisabled            = 106 // 此账户已被禁用
	ErrInvalidSignType            = 107 // sign_type 必须设置为 MD5 或 SHA1 或 normal
	ErrInvalidSignature           = 108 // signature 参数无效
	ErrInvalidAppKey              = 109 // appkey 无效
	ErrWrongSignType              = 110 // sign_type 错误
	ErrEmptySignatureParam        = 111 // 空的 signature 参数
	ErrSubscriptionDisabled       = 112 // 应用的订阅与退订功能已禁用
	ErrIPNotInWhitelist           = 113 // 您的 IP 不在白名单范围
	ErrPhoneInBlacklist           = 114 // 该手机号码在账户黑名单中
	ErrPhoneFrequencyLimit        = 115 // 该手机号码请求超限
	ErrSignatureUsedByOther       = 116 // 签名错误，该签名已被其他应用使用
	ErrTemplateSignatureInconsist = 117 // 该模板已失效，短信模板签名与固定签名不一致
	ErrTemplateInvalid            = 118 // 该模板已失效
	ErrPermissionDenied           = 119 // 您不具备使用该API的权限
	ErrTemplateExpired            = 120 // 模板已失效
	ErrSignatureNotReported       = 126 // 短信签名还未报备成功
	ErrSignatureAlreadyExists     = 127 // 短信签名已存在，无需创建新签名

	// 时间戳相关错误 (151-154)
	ErrTimestampError       = 151 // 错误的 UNIX 时间戳
	ErrInvalidTimestamp     = 152 // 错误的 UNIX 时间戳，请将请求时间控制在6秒以内
	ErrNoAvailableSignature = 154 // appid 下无可用签名

	// 地址簿相关错误 (201-203)
	ErrUnknownAddressbookModel = 201 // 未知的 addressbook 模式
	ErrIncorrectEmailAddress   = 202 // 错误的收件人地址
	ErrEmptyAddressbook        = 203 // 地址薄不包含任何联系人

	// 短信相关错误 (251-253)
	ErrIncorrectMessageAddress = 251 // 错误的收件人地址（message）
	ErrEmptyMessageAddressbook = 252 // 地址薄不包含任何联系人（message）
	ErrContactUnsubscribed     = 253 // 此联系人已退订你的短信系统

	// 项目相关错误 (305-310)
	ErrEmptyProjectID   = 305 // 没有填写项目标记
	ErrInvalidProjectID = 306 // 无效的项目标记
	ErrIncorrectJSON    = 307 // 错误的 json 格式
	ErrTagTooLong       = 310 // tag参数长度不能超过32个字符

	// 短信内容相关错误 (401-422)
	ErrEmptyMessageSignature      = 401 // 短信签名不能为空
	ErrSignatureTooLong           = 402 // 请将短信签名控制在40个字符以内
	ErrEmptyContent               = 403 // 短信正文不能为空
	ErrContentTooLong             = 404 // 请将短信内容（加上签名）控制在1000个字符以内
	ErrForbiddenWords             = 405 // 短信中包含禁用词或短语
	ErrEmptyProjectIDForContent   = 406 // 项目标记不能为空
	ErrInvalidProjectIDForContent = 407 // 无效的项目标记
	ErrDuplicateMessage           = 408 // 不能发送完全相同的短信
	ErrMessageUnderReview         = 409 // 短信项目正在审核中
	ErrInvalidMultiParam          = 410 // multi 参数无效
	ErrMissingSignatureInTemplate = 411 // 短信模板缺少签名
	ErrSignatureTooLongInTemplate = 412 // 短信签名超过10个字符
	ErrSignatureLengthInvalid     = 413 // 短信签名字数应在2到10个字符之间
	ErrEmptyContentInTemplate     = 414 // 请提交短信正文
	ErrContentTooLongInTemplate   = 415 // 短信正文超过1000个字符
	ErrTitleTooLong               = 416 // 短信标题超过64个字符
	ErrEmptyTemplateID            = 417 // 请提交需要更新的模板ID
	ErrTemplateNotExists          = 418 // 尝试更新的模板不存在
	ErrEmptyContentForUpdate      = 419 // 短信正文不能为空
	ErrNoMatchingTemplate         = 420 // 找不到可匹配的模板
	ErrTemplateTooLong            = 422 // 模板长度超过255个字符

	// 地址簿相关错误 (501)
	ErrInvalidAddressbookSign = 501 // 错误的目标地址簿标识

	// 配额和余额相关错误 (901-905)
	ErrQuotaExhausted               = 901 // 今日发送配额已用尽
	ErrInsufficientCredit           = 903 // 短信发送许可已用尽或余额不足
	ErrInsufficientBalance          = 904 // 账户余额已用尽或余额不足
	ErrInsufficientTransactionalSMS = 905 // 事务性短信余额不足
)

// ErrorMessages 错误信息映射
var ErrorMessages = map[int]string{
	ErrIncorrectAppID:               "不正确的 APP ID",
	ErrAppDisabled:                  "此应用已被禁用，请至 submail > 应用集成 > 应用 页面开启此应用",
	ErrDeveloperNotAvailable:        "未启用的开发者，此应用的开发者身份未验证，请更新您的开发者资料",
	ErrDeveloperNotVerified:         "此开发者未通过验证或此开发者资料发生更改。请至应用集成页面更新你的开发者资料",
	ErrAccountExpired:               "此账户已过期",
	ErrAccountDisabled:              "此账户已被禁用",
	ErrInvalidSignType:              "sign_type （验证模式）必须设置为 MD5（MD5签名模式）或 SHA1（SHA1签名模式）或 normal (密匙模式)",
	ErrInvalidSignature:             "signature 参数无效",
	ErrInvalidAppKey:                "appkey 无效",
	ErrWrongSignType:                "sign_type 错误",
	ErrEmptySignatureParam:          "空的 signature 参数",
	ErrSubscriptionDisabled:         "应用的订阅与退订功能已禁用",
	ErrIPNotInWhitelist:             "请求的 APPID 已设置 IP 白名单，您的 IP 不在此白名单范围",
	ErrPhoneInBlacklist:             "该手机号码在账户黑名单中，已被屏蔽",
	ErrPhoneFrequencyLimit:          "该手机号码请求超限",
	ErrSignatureUsedByOther:         "签名错误，该签名已被其他应用使用并已申请固定签名",
	ErrTemplateSignatureInconsist:   "该模板已失效，短信模板签名与固定签名不一致或你的账户已取消固签，请联系 SUBMAIL 管理员",
	ErrTemplateInvalid:              "该模板已失效，请联系SUBMAIL管理员",
	ErrPermissionDenied:             "您不具备使用该API的权限，请联系SUBMAIL管理员",
	ErrTemplateExpired:              "模板已失效",
	ErrSignatureNotReported:         "短信签名还未报备成功",
	ErrSignatureAlreadyExists:       "短信签名已存在，无需创建新签名",
	ErrTimestampError:               "错误的 UNIX 时间戳",
	ErrInvalidTimestamp:             "错误的 UNIX 时间戳，请将请求 UNIX 时间戳 至 发送 API 的过程控制在6秒以内",
	ErrNoAvailableSignature:         "appid 下无可用签名",
	ErrUnknownAddressbookModel:      "未知的 addressbook 模式",
	ErrIncorrectEmailAddress:        "错误的收件人地址",
	ErrEmptyAddressbook:             "错误的收件人地址。如果你正在使用 adressbook , 你所标记的地址薄不包含任何联系人",
	ErrIncorrectMessageAddress:      "错误的收件人地址（message）",
	ErrEmptyMessageAddressbook:      "错误的收件人地址（message）如果你正在使用 adressbook 模式，你所标记的地址薄不包含任何联系人",
	ErrContactUnsubscribed:          "此联系人已退订你的短信系统",
	ErrEmptyProjectID:               "没有填写项目标记",
	ErrInvalidProjectID:             "无效的项目标记",
	ErrIncorrectJSON:                "错误的 json 格式。 请检查 vars 和 links 参数",
	ErrTagTooLong:                   "tag参数长度不能超过32个字符",
	ErrEmptyMessageSignature:        "短信签名不能为空",
	ErrSignatureTooLong:             "请将短信签名控制在40个字符以内",
	ErrEmptyContent:                 "短信正文不能为空",
	ErrContentTooLong:               "请将短信内容（加上签名）控制在1000个字符以内",
	ErrForbiddenWords:               "依据当地法律法规，以下词或短语不能出现在短信中",
	ErrEmptyProjectIDForContent:     "项目标记不能为空",
	ErrInvalidProjectIDForContent:   "无效的项目标记",
	ErrDuplicateMessage:             "你不能向此联系人或此地址簿中包含的联系人发送完全相同的短信",
	ErrMessageUnderReview:           "尝试发送的短信项目正在审核中，请稍候再试",
	ErrInvalidMultiParam:            "multi 参数无效",
	ErrMissingSignatureInTemplate:   "您必须为每条短信模板提交一个短信签名，且该签名必须使用全角大括号【和】包括起来，请将短信签名的字数控制在2至10字符以内（括号不计算字符数）",
	ErrSignatureTooLongInTemplate:   "请将短信签名的字数控制在10字符以内（括号不计算字符数）",
	ErrSignatureLengthInvalid:       "请将短信签名的字数控制在2到10个字符之间（括号不计算字符数）",
	ErrEmptyContentInTemplate:       "请提交短信正文",
	ErrContentTooLongInTemplate:     "请将短信正文的字数控制在1000个字符以内",
	ErrTitleTooLong:                 "请将短信标题的字数控制在64个字符以内",
	ErrEmptyTemplateID:              "请提交需要更新的模板ID",
	ErrTemplateNotExists:            "尝试更新的模板不存在",
	ErrEmptyContentForUpdate:        "短信正文不能为空",
	ErrNoMatchingTemplate:           "找不到可匹配的模板",
	ErrTemplateTooLong:              "请控制您的模板长度在255个字符内",
	ErrInvalidAddressbookSign:       "错误的目标地址簿标识",
	ErrQuotaExhausted:               "你今日的发送配额已用尽。如需提高发送配额，请至 submail > 应用集成 >应用 页面开启更多发送配额",
	ErrInsufficientCredit:           "您的短信发送许可已用尽或您的余额不支持本次的请求数量。如需继续发送，请至 submail.cn > 商店 页面购买更多发送许可后重试",
	ErrInsufficientBalance:          "您的账户余额已用尽或您的余额不支持本次的请求数量。如需继续充值，请至 submail.cn > 商店 页面购买更多发送许可后重试",
	ErrInsufficientTransactionalSMS: "您的账户余额已用尽或您的余额不支持本次的请求数量。如需继续充值，请至 submail.cn > 商店 页面购买更多发送许可后重试",
}

// NewAPIError 创建API错误
func NewAPIError(code int, msg string) *APIError {
	description := ErrorMessages[code]
	if description == "" {
		description = "未知错误"
	}

	return &APIError{
		Code:        code,
		Msg:         msg,
		Description: description,
	}
}

// ParseAPIError 从响应中解析API错误
func ParseAPIError(data []byte) error {
	var errorResp struct {
		Status string      `json:"status"`
		Code   interface{} `json:"code"`
		Msg    string      `json:"msg"`
	}

	if err := json.Unmarshal(data, &errorResp); err != nil {
		// 如果无法解析为错误格式，说明可能是正常响应
		return nil
	}

	// 只有明确标记为error状态的响应才认为是错误
	if errorResp.Status != "error" {
		return nil
	}

	var code int
	switch v := errorResp.Code.(type) {
	case float64:
		code = int(v)
	case string:
		if c, err := strconv.Atoi(v); err == nil {
			code = c
		}
	case int:
		code = v
	}

	return NewAPIError(code, errorResp.Msg)
}
