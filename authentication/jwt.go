package authentication

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"inkgo/config"
	"inkgo/model"
	"inkgo/service"
	"time"
)

type CustomClaims struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

type JWT struct {
	secret         []byte
	issuer         string
	expireDuration int64
	tokenService   service.TokenService
}

func NewJWT(cfg *config.JWTConfig, tokenService service.TokenService) *JWT {
	return &JWT{
		secret:         []byte(cfg.Secret),
		issuer:         cfg.Issuer,
		expireDuration: cfg.Expire,
		tokenService:   tokenService,
	}
}

// GenerateToken 生成JWT令牌
func (j *JWT) GenerateToken(user *model.User) (string, error) {
	if user == nil {
		return "", errors.New("empty user type")
	}
	// 创建JWT并签名
	claims := CustomClaims{
		Name: user.UserName,
		ID:   user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.expireDuration) * time.Second)), // 设置过期时间为24小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                                    // 设置签发时间为当前时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                                    // 设置生效时间为当前时间
		},
	}
	// 使用 HMAC SHA256 签名方法创建一个新的 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// ParseToken 解析JWT令牌
func (j *JWT) ParseToken(tokenString string) (*model.User, error) {
	// 解析JWT
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return j.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	user := &model.User{
		Model: gorm.Model{
			ID: claims.ID,
		},
		UserName: claims.Name,
	}
	return user, nil
}

// RevokeToken 注销JWT令牌
func (j *JWT) RevokeToken(tokenString string) error {
	// 这里可以实现将token存储到缓存中，标记为已注销
	// 例如，可以将token存储在一个黑名单中，或者在数据库中记录该token的状态
	// 由于JWT是无状态的，所以通常不需要实际删除token，只需在应用逻辑中忽略它即可
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(j.secret), nil
	})

	if err != nil || !token.Valid {
		return errors.New("invalid or expired token")
	}

	// 获取token中的claims
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return errors.New("invalid token claims") // 在这一步报错啦
	}

	expiration := time.Until(time.Unix(claims.ExpiresAt.Unix(), 0))
	if expiration <= 0 {
		return errors.New("token has already expired")
	}

	key := j.RevokeKey(tokenString)
	// 将token存储到黑名单中，设置过期时间为token的过期时间
	if err := j.tokenService.SetToken(key, "blacklisted", expiration); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	return nil
}

func (j *JWT) RevokeKey(token string) string {
	return fmt.Sprintf("jwt:blacklist:" + token)
}

// IsTokenRevoked 检查token是否被注销
func (j *JWT) IsTokenRevoked(token string) (bool, error) {
	// 检查token是否在黑名单中
	key := j.RevokeKey(token)
	val, err := j.tokenService.GetToken(key)
	if err != nil {
		return false, fmt.Errorf("failed to check token: %w", err)
	}
	return val == "blacklisted", nil
}

// 更新token
func (j *JWT) RefreshToken(oldToken string) (string, error) {
	//解析旧的token
	token, err := jwt.ParseWithClaims(oldToken, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return j.secret, nil
	})
	// 检查解析是否出错或token是否有效
	if err != nil || !token.Valid {
		return "", errors.New("invalid or expired token")
	}

	// 获取旧token中的claims
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	// 检查过期时间，如果还有超过5分钟的有效期，则不需要刷新
	if claims.ExpiresAt != nil && time.Until(claims.ExpiresAt.Time) > 5*time.Minute {
		return "", errors.New("token is still valid, no need to refresh")
	}

	// 创建新的claims
	newClaims := CustomClaims{
		ID:   claims.ID,
		Name: claims.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.expireDuration) * time.Second)), // 设置新的过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                                    // 设置新的签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                                    // 设置新的生效时间
		},
	}
	// 使用 HMAC SHA256 签名方法创建一个新的 JWT
	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	return tokenString.SignedString(j.secret)
}
