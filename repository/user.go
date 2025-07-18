package repository

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"inkgo/database"
	"inkgo/model"
	"strconv"
)

type userRepository struct {
	db  *gorm.DB
	rdb *database.RedisDB
}

// 实例化
func NewUserRepository(db *gorm.DB, rdb *database.RedisDB) UserRepository {
	return &userRepository{
		db:  db,
		rdb: rdb,
	}
}

// UpdatePassword 用于更新用户密码
func (u *userRepository) UpdatePassword(userID uint, hashedPassword string) error {
	user := &model.User{
		Model: gorm.Model{
			ID: userID,
		},
		Password: hashedPassword,
	}
	if err := u.db.Model(&model.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error; err != nil {
		return err
	}
	// 更新缓存
	if err := u.SetCache(user); err != nil {
		logrus.Errorf("failed to set user cache: %v", err)
	}
	return nil
}

// UpdateRole 用于更新用户角色
//func (u *userRepository) UpdateRole(userID uint, role string) error {
//	user := &model.User{
//		Model: gorm.Model{
//			ID: userID,
//		},
//	}
//	if err := u.db.Model(&model.User{}).Where("id = ?", userID).Update("role", role).Error; err != nil {
//		return err
//	}
//	// 更新缓存
//	if err := u.SetCache(user); err != nil {
//		logrus.Errorf("failed to set user cache: %v", err)
//	}
//	return nil
//}

// IsUsernameTaken 检查用户名是否已被使用
func (u *userRepository) IsUsernameTaken(username string) (bool, error) {
	var count int64
	if err := u.db.Model(&model.User{}).Where("user_name = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsEmailTaken 检查电子邮件是否已被使用
func (u *userRepository) IsEmailTaken(email string) (bool, error) {
	var count int64
	if err := u.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsMobileTaken 检查手机号是否已被使用
func (u *userRepository) IsMobileTaken(mobile string) (bool, error) {
	var count int64
	if err := u.db.Model(&model.User{}).Where("mobile = ?", mobile).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetUserByID 用于通过用户ID获取用户信息
func (u *userRepository) GetUserByID(id uint) (*model.User, error) {
	var user *model.User
	// 先尝试从缓存中获取用户信息
	if err := u.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	// 如果缓存中没有用户信息，则从数据库中获取
	if err := u.SetCache(user); err != nil {
		logrus.Errorf("failed to set user: %v", err)
	}

	return user, nil
}

// 登录（支持账号密码或 OAuth）
func (u *userRepository) FindByOauth(provider, provider_id string) (*model.User, error) {
	oauthAccount := new(model.OAuthAccount)
	if err := u.db.Where("provider = ? and provider_id = ?", provider, provider_id).First(&oauthAccount).Error; err != nil {
		return nil, err
	}
	return u.GetUserByID(oauthAccount.UserID)
}

// GetUserByName 用于通过用户名获取用户信息
func (u *userRepository) FindByUserName(name string) (*model.User, error) {
	var user model.User
	if err := u.db.Where("user_name = ?", name).First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return &user, nil
}

// FindByMobile 用于通过手机号查找用户
func (u *userRepository) FindByMobile(mobile string) (*model.User, error) {
	var user model.User
	// 检查手机号是否存在
	if err := u.db.Where("mobile = ?", mobile).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("手机号不存在")
		}
		return nil, err
	}

	return &user, nil
}

// FindByEmail  用于通过电子邮件查找用户
func (u *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		// 如果没有找到记录，返回一个特定的错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email不存在")
		}
		// 其他错误直接返回
		return nil, err
	}
	return &user, nil
}

// List 用于分页获取用户列表
func (u *userRepository) List(pageSize int, page int) ([]model.User, int64, error) {
	var users []model.User
	//userList := make([]interface{}, 0, len(users))

	var total int64
	if err := u.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := u.db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// Create 用于创建新用户
func (u *userRepository) Create(user *model.User) (*model.User, error) {
	userCreateField := []string{"user_name", "email", "password", "mobile", "avatar"}
	if user.Role != model.RoleUser && user.Role != model.RoleAdmin {
		return nil, errors.New("无效的角色值")
	}
	if err := u.db.Unscoped().Select(userCreateField).Create(user).Error; err != nil {
		return nil, err
	}
	u.SetCache(user)
	return user, nil
}

// Update 用于更新用户信息
func (u *userRepository) Update(user *model.User) (*model.User, error) {
	// 检查用户是否被注册
	exist, err := u.IsUsernameTaken(user.UserName)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errors.New("用户已经存在")
	}

	// 检查电子邮件是否已被使用
	emailTaken, err := u.IsEmailTaken(user.Email)
	if err != nil {
		return nil, err
	}
	if emailTaken {
		return nil, errors.New("电子邮件已被使用")
	}

	// 检查手机号是否已被使用
	mobileTaken, err := u.IsMobileTaken(user.Mobile)
	if err != nil {
		return nil, err
	}
	if mobileTaken {
		return nil, errors.New("手机号已被使用")
	}
	if err := u.db.Model(&model.User{}).Where("id = ?", user.ID).Updates(&user).Error; err != nil {
		return nil, err
	}
	u.rdb.HDel(user.CacheKey(), strconv.Itoa(int(user.ID)))
	return user, nil
}

func (u *userRepository) Exists(userID uint) (bool, error) {
	var count int64
	err := u.db.Model(&model.User{}).
		Where("id = ?", userID).
		Count(&count).Error
	return count > 0, err
}

// Delete 用于删除用户
func (u *userRepository) Delete(id uint) error {
	var user model.User
	// 首先检查用户是否存在
	exists, err := u.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("用户不存在")
	}

	if err := u.db.Delete(&user, id).Error; err != nil {
		return err
	}

	u.rdb.HDel(user.CacheKey(), strconv.Itoa(int(user.ID)))
	return nil
}

// Deactvate 用于禁用用户
func (u *userRepository) Deactivate(id uint) error {
	var user model.User
	// 首先检查用户是否存在
	exists, err := u.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("用户不存在")
	}
	if err := u.db.Model(&user).Update("status", "deleted").Error; err != nil {
		return err
	}
	return nil
}

func (u *userRepository) SetCache(user *model.User) error {
	if user == nil {
		return nil
	}
	return u.rdb.HSet(user.CacheKey(), strconv.Itoa(int(user.ID)), user)
}

func (u *userRepository) Migrate() error {
	return u.db.AutoMigrate(&model.User{}, &model.OAuthAccount{})
}
