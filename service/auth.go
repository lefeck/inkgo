package service

import (
	"context"
	"fmt"
	"inkgo/repository"
	"time"

	"gopkg.in/gomail.v2"
	"inkgo/config"
	"inkgo/utils"
)

type AuthService interface {
	SendEmailVerificationCode(toEmail, fromEmail, smtpAuthCode string) error
	VerifyEmailCode(email, code string) (resetToken string, err error)
	GetEmailByResetToken(token string) (email string, err error)
}

type authService struct {
	authRepository repository.AuthRepository
}

func NewAuthService(authRepository repository.AuthRepository) AuthService {
	return &authService{
		authRepository: authRepository,
	}
}

// 发送验证码到用户邮箱
func (auth *authService) SendEmailVerificationCode(from, to, authCode string) error {
	// 生成验证码
	code := utils.GenerateCode()

	// 提取发送邮箱的域名
	_, domain := utils.ParseEmail(from)
	// 获取对应 SMTP 配置
	smtpConfig, exist := config.LookupSMTPConfig(domain)
	if !exist {
		return fmt.Errorf("不支持该邮箱服务商: %s", domain)
	}

	// 构建邮件内容
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "找回密码 - 验证码")
	m.SetBody("text/html", fmt.Sprintf(`
		<h2>您的验证码</h2>
		<p>请使用以下验证码完成操作：</p>
		<p style="font-size: 24px; color: blue;"><strong>%s</strong></p>
		<p>验证码将在 5 分钟后失效。</p>
	`, code))

	// 发送邮件
	d := gomail.NewDialer(smtpConfig.Host, smtpConfig.Port, from, authCode)

	// 启用 SSL（根据服务商自动判断，也可以手动配置）
	d.SSL = smtpConfig.Port == 465

	expiration := 5 * time.Minute
	// 将验证码存入 Redis，设置过期时间为5分钟
	if err := auth.authRepository.SaveEmailCode(context.Background(), to, code, expiration); err != nil {
		return fmt.Errorf("存储验证码失败: %v", err)
	}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	return nil
}

// VerifyCode 用于验证验证码
func (auth *authService) VerifyEmailCode(email string, code string) (string, error) {
	if email == "" || code == "" {
		return "", fmt.Errorf("email or code is empty")
	}
	// 获取存储在 Redis 中的验证码
	value, err := auth.authRepository.FindEmailCode(context.Background(), email)
	if err != nil {
		return "", fmt.Errorf("获取验证码失败: %v", err)
	}
	if value != code {
		return "", fmt.Errorf("验证码不匹配")
	}

	// 验证成功后生成reset_token
	token := "reset_" + utils.GenerateCode()
	err = auth.authRepository.SaveResetToken(context.Background(), token, email, 15*time.Minute)
	if err != nil {
		return "", fmt.Errorf("生成 token 失败: %v", err)
	}

	return token, nil
}

// 从缓存或数据库中获取用户的电子邮件地址
func (auth *authService) GetEmailByResetToken(token string) (string, error) {
	value, err := auth.authRepository.FindEmailByResetToken(context.Background(), token)
	if err != nil {
		return "", fmt.Errorf("获取电子邮件失败: %v", err)
	}
	return value, nil
}
