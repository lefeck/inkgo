package utils

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"inkgo/test/forms"
	"time"
)

func SignUpParamStructLevelValidation(sl validator.StructLevel) {
	su := sl.Current().Interface().(forms.CreateUserForm)
	if su.Password != su.RePassword {
		// 输出错误提示信息，最后一个参数就是传递的param
		sl.ReportError(su.RePassword, "re_password", "RePassword", "eqfield", "password")
	}
}

func CustemFunc(fl validator.FieldLevel) bool {
	date, err := time.Parse("2006-01-02", fl.Field().String())
	if err != nil {
		return false
	}
	if date.Before(time.Now()) {
		return true
	}
	return false
}

func RegisterTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(trans ut.Translator) error {
		if err := trans.Add(tag, msg, true); err != nil {
			return err
		}
		return nil
	}
}

func Translate(trans ut.Translator, fe validator.FieldError) string {
	t, err := trans.T(fe.Tag(), fe.Field())
	if err != nil {
		panic(fe.(error).Error())
	}
	return t
}
