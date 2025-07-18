package controller

//var trans ut.Translator
//
//type TransController struct {
//}
//
//func NewTransController() *TransController {
//	return &TransController{}
//}
//
//func (tc *TransController) Trans(locale string) (err error) {
//	// 修改gin框架中的Validator引擎属性，实现自定制
//	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
//		// 注册一个获取json tag的自定义方法
//		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
//			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
//			if name == "-" {
//				return ""
//			}
//			return name
//		})
//		//为SignUpParam注册自定义校验方法
//
//		//v.RegisterStructValidation(nil, request.Login{})
//
//		v.RegisterStructValidation(utils.SignUpParamStructLevelValidation, forms.CreateUserForm{})
//
//		v.RegisterStructValidation(nil, forms.UpdateUserForm{})
//
//		if err := v.RegisterValidation("checkDate", utils.CustemFunc); err != nil {
//			return err
//		}
//
//		//校验手机号码
//		//if err := v.RegisterValidation("mobile", ValidateMobile); err != nil {
//		//	return err
//		//}
//
//		zhT := zh.New() // 中文翻译器
//		enT := en.New() // 英文翻译器
//
//		// 第一个参数是备用（fallback）的语言环境
//		// 后面的参数是应该支持的语言环境（支持多个）
//		// uni := ut.New(zhT, zhT) 也是可以的
//		uni := ut.New(enT, zhT, enT)
//
//		// locale 通常取决于 http 请求头的 'Accept-Language'
//		var ok bool
//		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
//		trans, ok = uni.GetTranslator(locale)
//		if !ok {
//			return fmt.Errorf(
//
//				"uni.GetTranslator(%s) failed", locale)
//		}
//
//		// 注册翻译器
//		switch locale {
//		case "en":
//			err = enTranslations.RegisterDefaultTranslations(v, trans)
//		case "zh":
//			err = zhTranslations.RegisterDefaultTranslations(v, trans)
//		default:
//			err = enTranslations.RegisterDefaultTranslations(v, trans)
//		}
//		//注册date翻译
//		if err := v.RegisterTranslation("checkDate", trans, utils.RegisterTranslator("checkDate", "{0}必须要晚于当前日期"), utils.Translate); err != nil {
//			return err
//		}
//		//注册手机号码翻译
//		if err := v.RegisterTranslation("mobile", trans, utils.RegisterTranslator("mobile", "{0}非法的手机号码"), utils.Translate); err != nil {
//			return err
//		}
//		return
//	}
//	return
//}
