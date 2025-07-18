package config

type SMTPConfig struct {
	Host string `yaml:"host" json:"host"` // SMTP服务器地址
	Port int    `yaml:"port" json:"port"` // SMTP服务器端口
}

var smtpConfig = map[string]SMTPConfig{
	"qq.com": {
		Host: "smtp.qq.com",
		Port: 587,
	},
	"163.com": {
		Host: "smtp.163.com",
		Port: 25,
	},
	"126.com": {
		Host: "smtp.126.com",
		Port: 25,
	},
	"gmail.com": {
		Host: "smtp.gmail.com",
		Port: 587,
	},
	"aliyun.com": {
		Host: "smtp.aliyun.com",
		Port: 465,
	},
	"sina.com": {
		Host: "smtp.sina.com",
		Port: 25,
	},
}

// LookupSMTPConfig 根据邮箱域名获取SMTP配置
func LookupSMTPConfig(domain string) (SMTPConfig, bool) {
	config, exists := smtpConfig[domain]
	if !exists {
		return SMTPConfig{}, false
	}
	return config, true
}
