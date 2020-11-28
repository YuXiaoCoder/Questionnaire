package common

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// 全局翻译器
var Trans ut.Translator

// 初始化翻译器
func InitTrans(locale string) (err error) {
	// 自定义Gin框架中的Validator引擎属性
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		// 注册标签函数
		validate.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			} else {
				return name
			}
		})

		enTrans := en.New() // 英文翻译器
		zhTrans := zh.New() // 中文翻译器

		uni := ut.New(enTrans, zhTrans)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// 注册翻译器
		switch locale {
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(validate, Trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(validate, Trans)
		}
		return err
	}
	return nil
}

// 从验证错误消息中移除结构体名称
// https://github.com/go-playground/validator/issues/633#issuecomment-654382345
func RemoveTopStruct(fields map[string]string) map[string]string {
	result := map[string]string{}
	for field, err := range fields {
		result[field[strings.Index(field, ".")+1:]] = err
	}
	return result
}
