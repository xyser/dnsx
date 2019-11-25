package customtime

import (
	"time"

	"dnsx/internal/consts"

	translator "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

// ValidateTime 验证是否是时间
func ValidateTime(fl validator.FieldLevel) bool {
	_, err := time.ParseInLocation(consts.DateTimeFormat, fl.Field().String(), time.Local)
	if err != nil {
		return false
	}
	return true
}

// ValidateTimeTranslator 翻译
func ValidateTimeTranslator(ut translator.Translator) (err error) {
	return ut.Add("time", "{0} 时间格式错误", true)
}

// ValidateTimeHM 验证是否是时间 eg. 13:09
func ValidateTimeHM(fl validator.FieldLevel) bool {
	_, err := time.ParseInLocation(consts.DateTimeHMFormat, fl.Field().String(), time.Local)
	if err != nil {
		return false
	}
	return true
}

// ValidateTimeHMTranslator 翻译
func ValidateTimeHMTranslator(ut translator.Translator) (err error) {
	return ut.Add("hour", "{0} 时间格式错误", true)
}
