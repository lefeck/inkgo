package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"inkgo/common/request"
	"inkgo/utils"

	"golang.org/x/crypto/bcrypt"
	"inkgo/model"
	"inkgo/repository"
	"strconv"
)

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

// Deactivate 用于禁用用户账号( 例如：注销账号)
func (u *userService) Deactivate(user *model.User) error {
	if user.ID == 0 {
		return errors.New("user ID is empty")
	}
	if err := u.userRepository.Deactivate(user.ID); err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}
	return nil
}

// FindByUserName 根据用户名查找用户
func (u *userService) FindByUserName(username string) (*model.User, error) {
	if username == "" {
		return nil, errors.New("username is empty")
	}
	user, err := u.userRepository.FindByUserName(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userService) CreateOAuthUser(user *model.User) (*model.User, error) {
	if user.OAuthAccounts == nil {
		return nil, errors.New("user is nil")
	}
	oauthAccount := user.OAuthAccounts[0]
	if oauthAccount.Provider == "" || oauthAccount.ProviderID == "" {
		return nil, errors.New("oauth provider or provider ID is empty")
	}
	existingUser, err := u.userRepository.FindByOauth(oauthAccount.Provider, oauthAccount.ProviderID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to find user by oauth: %w", err)
	}
	if existingUser != nil {
		return existingUser, nil
	}
	return u.userRepository.Create(user)

}

// FindByMobile 根据手机号查找用户
func (u *userService) FindByMobile(mobile string) (*model.User, error) {
	if mobile == "" {
		return nil, errors.New("mobile is empty")
	}
	user, err := u.userRepository.FindByMobile(mobile)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindByEmail 根据邮箱查找用户
func (u *userService) FindByEmail(email string) (*model.User, error) {
	if email == "" {
		return nil, errors.New("email is empty")
	}
	user, err := u.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 登录方式识别
/*
 密码登录
UserName,Mobile或Email 这3类登录, 都是密码登录blog
当用户输入UserName登录时, 只需要输入密码就可以登录, 如果密码忘记了, 可以通过email 获取 验证吗, 修改登录密码
当用户输入Mobile登录时, 也需要输入密码登录
当用户输入Email登录时, 也需要输入密码登录
*/
func (u *userService) LoginByPassword(identify, password string, user *model.User) (*model.User, error) {
	if identify == "" || password == "" {
		return nil, errors.New("identifier or password is empty")
	}
	var err error
	switch request.Detect(identify) {
	case "username":
		user, err = u.userRepository.FindByUserName(identify)
	case "mobile":
		user, err = u.userRepository.FindByMobile(identify)
	case "email":
		user, err = u.userRepository.FindByEmail(identify)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %v", err)
	}
	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}
	return user, nil
}

// 手机号/邮箱登录登录

func (u *userService) LoginByCode(identify string) (*model.User, error) {
	return nil, errors.New("login with password reset is not implemented yet")
}

// 验证码登录
func (u *userService) LoginWithCode(identify, code string) (*model.User, error) {
	return nil, errors.New("login with code is not implemented yet")
}

func (u *userService) UpdatePassword(userID string, password string) error {
	if userID == "" || password == "" {
		return errors.New("user ID or password is empty")
	}

	uid, err := strconv.Atoi(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	if len(password) > 0 {
		hashed, err := utils.HashPassword(password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		if err := u.userRepository.UpdatePassword(uint(uid), hashed); err != nil {
			return fmt.Errorf("failed to update password: %w", err)
		}
	}

	return nil
}

// UpdateRole 更新用户角色
//func (u *userService) UpdateRole(userID string, role string) error {
//	if userID == "" || role == "" {
//		return errors.New("user ID or role is empty")
//	}
//
//	uid, err := strconv.Atoi(userID)
//	if err != nil {
//		return fmt.Errorf("invalid user ID: %w", err)
//	}
//
//	if err := u.userRepository.UpdateRole(uint(uid), role); err != nil {
//		return fmt.Errorf("failed to update role: %w", err)
//	}
//	return nil
//}

// FindByOauth 根据 OAuth 提供商和提供商 ID 查找用户
func (u *userService) FindByOauth(provider string, providerID string) (*model.User, error) {
	if provider == "" || providerID == "" {
		return nil, errors.New("provider or providerID is empty")
	}

	user, err := u.userRepository.FindByOauth(provider, providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by oauth: %w", err)
	}
	return user, nil
}

// IsUsernameTaken 检查用户名是否已被使用
func (u *userService) IsUsernameTaken(username string) (bool, error) {
	if username == "" {
		return false, errors.New("username is empty")
	}

	taken, err := u.userRepository.IsUsernameTaken(username)
	if err != nil {
		return false, fmt.Errorf("failed to check if username is taken: %w", err)
	}
	return taken, nil
}

// IsEmailTaken 检查邮箱是否已被使用
func (u *userService) IsEmailTaken(email string) (bool, error) {
	if email == "" {
		return false, errors.New("email is empty")
	}

	taken, err := u.userRepository.IsEmailTaken(email)
	if err != nil {
		return false, fmt.Errorf("failed to check if email is taken: %w", err)
	}
	return taken, nil
}

// List 获取用户列表
func (u *userService) List(pageSize int, page int) ([]model.User, int64, error) {
	user, total, err := u.userRepository.List(pageSize, page)
	if err != nil {
		return nil, 0, err
	}
	return user, total, nil
}

// GetUserByName 根据用户名获取用户
func (u *userService) GetUserByName(name string) (*model.User, error) {
	if name == "" {
		return nil, errors.New("user name is empty")
	}
	user, err := u.userRepository.FindByUserName(name)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Register 创建用户（注册）
func (u *userService) Create(user *model.User) (*model.User, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(password)
	return u.userRepository.Create(user)
}

// GetUserByID 根据用户ID获取用户
func (u *userService) GetUserByID(id string) (*model.User, error) {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return u.userRepository.GetUserByID(uint(uid))
}

func (u *userService) Update(user *model.User) (*model.User, error) {
	if len(user.Password) > 0 {
		hashed, err := utils.HashPassword(user.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashed
	}
	return u.userRepository.Update(user)
}

func (u *userService) Delete(id string) error {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return u.userRepository.Delete(uint(uid))
}

//func (u *userService) Login(user *model.User) (*model.User, error) {
//	var login request.LoginRequest
//	user, err := u.userRepository.FindByUserName(user.UserName)
//	if err != nil {
//		return nil, err
//	}
//
//	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
//		return nil, err
//	}
//	return user, nil
//}
