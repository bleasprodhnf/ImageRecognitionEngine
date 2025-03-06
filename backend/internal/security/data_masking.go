package security

import (
	"regexp"
	"strings"
)

// DataMaskingService 提供数据脱敏功能
type DataMaskingService struct {
	encryption *EncryptionService
}

// NewDataMaskingService 创建一个新的数据脱敏服务实例
func NewDataMaskingService(encryption *EncryptionService) *DataMaskingService {
	return &DataMaskingService{
		encryption: encryption,
	}
}

// MaskPhoneNumber 对手机号码进行脱敏处理
// 例如：13812345678 -> 138****5678
func (s *DataMaskingService) MaskPhoneNumber(phone string) string {
	if len(phone) < 7 {
		return phone
	}

	// 保留前3位和后4位，中间用*替代
	return phone[:3] + strings.Repeat("*", len(phone)-7) + phone[len(phone)-4:]
}

// MaskEmail 对电子邮件进行脱敏处理
// 例如：user@example.com -> u***@example.com
func (s *DataMaskingService) MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	username := parts[0]
	domain := parts[1]

	if len(username) <= 1 {
		return email
	}

	// 保留用户名第一个字符，其余用*替代
	maskedUsername := username[:1] + strings.Repeat("*", len(username)-1)
	return maskedUsername + "@" + domain
}

// MaskIDCard 对身份证号进行脱敏处理
// 例如：110101199001011234 -> 1101**********1234
func (s *DataMaskingService) MaskIDCard(idCard string) string {
	if len(idCard) < 10 {
		return idCard
	}

	// 保留前4位和后4位，中间用*替代
	return idCard[:4] + strings.Repeat("*", len(idCard)-8) + idCard[len(idCard)-4:]
}

// MaskBankCard 对银行卡号进行脱敏处理
// 例如：6222021234567890123 -> 622202******0123
func (s *DataMaskingService) MaskBankCard(bankCard string) string {
	if len(bankCard) < 10 {
		return bankCard
	}

	// 保留前6位和后4位，中间用*替代
	return bankCard[:6] + strings.Repeat("*", len(bankCard)-10) + bankCard[len(bankCard)-4:]
}

// MaskJSON 对JSON数据中的敏感字段进行脱敏处理
func (s *DataMaskingService) MaskJSON(jsonStr string) string {
	// 手机号码脱敏
	phoneRegex := regexp.MustCompile(`"(phone|mobile|tel)"\s*:\s*"(\d{11})"`)
	jsonStr = phoneRegex.ReplaceAllStringFunc(jsonStr, func(match string) string {
		parts := phoneRegex.FindStringSubmatch(match)
		if len(parts) == 3 {
			maskedPhone := s.MaskPhoneNumber(parts[2])
			return `"` + parts[1] + `":"` + maskedPhone + `"`
		}
		return match
	})

	// 邮箱脱敏
	emailRegex := regexp.MustCompile(`"(email|mail)"\s*:\s*"([^"]+@[^"]+)"`)
	jsonStr = emailRegex.ReplaceAllStringFunc(jsonStr, func(match string) string {
		parts := emailRegex.FindStringSubmatch(match)
		if len(parts) == 3 {
			maskedEmail := s.MaskEmail(parts[2])
			return `"` + parts[1] + `":"` + maskedEmail + `"`
		}
		return match
	})

	// 身份证号脱敏
	idCardRegex := regexp.MustCompile(`"(idCard|idNumber)"\s*:\s*"(\d{15,18})"`)
	jsonStr = idCardRegex.ReplaceAllStringFunc(jsonStr, func(match string) string {
		parts := idCardRegex.FindStringSubmatch(match)
		if len(parts) == 3 {
			maskedIDCard := s.MaskIDCard(parts[2])
			return `"` + parts[1] + `":"` + maskedIDCard + `"`
		}
		return match
	})

	// 银行卡号脱敏
	bankCardRegex := regexp.MustCompile(`"(bankCard|cardNumber)"\s*:\s*"(\d{16,19})"`)
	jsonStr = bankCardRegex.ReplaceAllStringFunc(jsonStr, func(match string) string {
		parts := bankCardRegex.FindStringSubmatch(match)
		if len(parts) == 3 {
			maskedBankCard := s.MaskBankCard(parts[2])
			return `"` + parts[1] + `":"` + maskedBankCard + `"`
		}
		return match
	})

	return jsonStr
}

// MaskSensitiveData 对敏感数据进行脱敏处理
func (s *DataMaskingService) MaskSensitiveData(data map[string]string) map[string]string {
	result := make(map[string]string)

	for key, value := range data {
		switch {
		case strings.Contains(strings.ToLower(key), "phone") || strings.Contains(strings.ToLower(key), "mobile") || strings.Contains(strings.ToLower(key), "tel"):
			result[key] = s.MaskPhoneNumber(value)
		case strings.Contains(strings.ToLower(key), "email") || strings.Contains(strings.ToLower(key), "mail"):
			result[key] = s.MaskEmail(value)
		case strings.Contains(strings.ToLower(key), "idcard") || strings.Contains(strings.ToLower(key), "idnumber"):
			result[key] = s.MaskIDCard(value)
		case strings.Contains(strings.ToLower(key), "bankcard") || strings.Contains(strings.ToLower(key), "cardnumber"):
			result[key] = s.MaskBankCard(value)
		default:
			result[key] = value
		}
	}

	return result
}