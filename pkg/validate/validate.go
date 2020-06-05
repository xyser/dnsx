package validate

import (
	"errors"
	"reflect"
	"strings"
	"sync"

	customtime "github.com/dingdayu/dnsx/pkg/validate/custom/time"

	"github.com/gin-gonic/gin/binding"
	local "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	translations "gopkg.in/go-playground/validator.v9/translations/zh"
)

var ginValidator *Validator
var validate *validator.Validate
var uni *ut.UniversalTranslator
var trans ut.Translator
var _ binding.StructValidator = &Validator{}

type Validator struct {
	once     sync.Once
	validate *validator.Validate
	trans    *ut.Translator
}

func init() {
	validate = validator.New()
	localZH := local.New()
	uni = ut.New(localZH, localZH)
	trans, _ = uni.GetTranslator("zh")
	_ = translations.RegisterDefaultTranslations(validate, trans)
	customValidator()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	ginValidator = &Validator{
		validate: validate,
	}
}

// GinValidator 初始化验证器
func GinValidator() *Validator {
	return ginValidator
}

// Default 默认验证器
func Default() *validator.Validate {
	return validate
}

func (v *Validator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyinit()
		if err := v.validate.Struct(obj); err != nil {
			errs := err.(validator.ValidationErrors)
			messages := make([]string, len(errs), len(errs))
			for i, e := range errs {
				messages[i] = e.Translate(*v.trans)
			}
			return errors.New(strings.Join(messages, ", "))
		}
	}
	return nil
}

func (v *Validator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *Validator) lazyinit() {
	v.once.Do(func() {
		v.validate = validate
		v.trans = &trans
		// v.validate.SetTagName("binding")
	})
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(fe.Tag(), fe.Field())
	return t
}

// customValidator 统一注册自定义验证器
func customValidator() {
	_ = validate.RegisterValidation("time", customtime.ValidateTime)
	_ = validate.RegisterTranslation("time", trans, customtime.ValidateTimeTranslator, translateFunc)

	_ = validate.RegisterValidation("hour", customtime.ValidateTimeHM)
	_ = validate.RegisterTranslation("hour", trans, customtime.ValidateTimeHMTranslator, translateFunc)
}
