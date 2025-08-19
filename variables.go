package submail

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// VariableProcessor 变量处理器
type VariableProcessor struct {
	timezone *time.Location // 时区设置
}

// NewVariableProcessor 创建变量处理器
func NewVariableProcessor() *VariableProcessor {
	// 默认使用中国时区
	location, _ := time.LoadLocation("Asia/Shanghai")
	return &VariableProcessor{
		timezone: location,
	}
}

// SetTimezone 设置时区
func (vp *VariableProcessor) SetTimezone(timezone string) error {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return fmt.Errorf("无效的时区: %v", err)
	}
	vp.timezone = location
	return nil
}

// ProcessVariables 处理短信内容中的变量
func (vp *VariableProcessor) ProcessVariables(content string, vars map[string]string) string {
	// 处理自定义变量 @var(key_name)
	content = vp.processCustomVariables(content, vars)

	// 处理日期时间变量 @date()
	content = vp.processDateVariables(content)

	return content
}

// processCustomVariables 处理自定义变量
func (vp *VariableProcessor) processCustomVariables(content string, vars map[string]string) string {
	if vars == nil {
		return content
	}

	// 匹配 @var(key_name) 格式的变量
	re := regexp.MustCompile(`@var\(([^)]+)\)`)

	return re.ReplaceAllStringFunc(content, func(match string) string {
		// 提取变量名
		matches := re.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match // 如果匹配失败，返回原字符串
		}

		varName := matches[1]
		if value, exists := vars[varName]; exists {
			return value
		}

		// 如果变量不存在，返回原字符串
		return match
	})
}

// processDateVariables 处理日期时间变量
func (vp *VariableProcessor) processDateVariables(content string) string {
	now := time.Now().In(vp.timezone)

	// 定义日期变量映射
	dateVars := map[string]string{
		"@date()":  now.Format("2006-01-02 15:04:05"),     // Y-m-d H:i:s
		"@date(Y)": now.Format("2006"),                    // 4位年份
		"@date(y)": now.Format("06"),                      // 2位年份
		"@date(m)": now.Format("01"),                      // 月份 01-12
		"@date(M)": now.Format("Jan"),                     // 英文简写月份
		"@date(F)": now.Format("January"),                 // 英文完整月份
		"@date(d)": now.Format("02"),                      // 日期 01-31
		"@date(D)": now.Format("Mon"),                     // 英文简写星期
		"@date(l)": now.Format("Monday"),                  // 英文完整星期
		"@date(h)": now.Format("15"),                      // 小时 00-23
		"@date(i)": now.Format("04"),                      // 分钟 00-59
		"@date(s)": now.Format("05"),                      // 秒钟 00-59
		"@date(N)": fmt.Sprintf("%d", int(now.Weekday())), // 星期几数字 (0=Sunday, 1=Monday, ...)
	}

	// 替换所有日期变量
	for pattern, value := range dateVars {
		content = strings.ReplaceAll(content, pattern, value)
	}

	return content
}

// ValidateVariables 验证变量格式
func (vp *VariableProcessor) ValidateVariables(content string) []string {
	var errors []string

	// 检查自定义变量格式
	varRe := regexp.MustCompile(`@var\([^)]*\)`)
	varMatches := varRe.FindAllString(content, -1)

	for _, match := range varMatches {
		if !regexp.MustCompile(`^@var\([a-zA-Z_][a-zA-Z0-9_]*\)$`).MatchString(match) {
			errors = append(errors, fmt.Sprintf("无效的变量格式: %s", match))
		}
	}

	// 检查日期变量格式
	dateRe := regexp.MustCompile(`@date\([^)]*\)`)
	dateMatches := dateRe.FindAllString(content, -1)

	validDateFormats := map[string]bool{
		"@date()": true, "@date(Y)": true, "@date(y)": true, "@date(m)": true,
		"@date(M)": true, "@date(F)": true, "@date(d)": true, "@date(D)": true,
		"@date(l)": true, "@date(h)": true, "@date(i)": true, "@date(s)": true,
		"@date(N)": true,
	}

	for _, match := range dateMatches {
		if !validDateFormats[match] {
			errors = append(errors, fmt.Sprintf("无效的日期变量格式: %s", match))
		}
	}

	return errors
}

// ExtractVariableNames 提取内容中的自定义变量名
func (vp *VariableProcessor) ExtractVariableNames(content string) []string {
	var varNames []string

	re := regexp.MustCompile(`@var\(([^)]+)\)`)
	matches := re.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 2 {
			varName := match[1]
			// 检查是否已存在，避免重复
			found := false
			for _, existing := range varNames {
				if existing == varName {
					found = true
					break
				}
			}
			if !found {
				varNames = append(varNames, varName)
			}
		}
	}

	return varNames
}

// GetDateVariableDescription 获取日期变量说明
func (vp *VariableProcessor) GetDateVariableDescription() map[string]string {
	return map[string]string{
		"@date()":  "输出 Y-m-d H:i:s 格式日期 (如: 2024-01-15 14:30:25)",
		"@date(Y)": "输出当前年份 (如: 2024)",
		"@date(y)": "输出2位年份 (如: 24)",
		"@date(m)": "输出当前月份 (如: 01)",
		"@date(M)": "输出英文简写月份 (如: Jan)",
		"@date(F)": "输出英文完整月份 (如: January)",
		"@date(d)": "输出当日日期 (如: 15)",
		"@date(D)": "输出英文简写星期 (如: Mon)",
		"@date(l)": "输出英文完整星期 (如: Monday)",
		"@date(h)": "输出当前小时 (如: 14)",
		"@date(i)": "输出当前分钟 (如: 30)",
		"@date(s)": "输出当前秒钟 (如: 25)",
		"@date(N)": "输出星期几数字 (0=Sunday, 1=Monday, ...)",
	}
}
