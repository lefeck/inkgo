package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

// UserInfo 用户信息
type UserInfo struct {
	ID     uint `json:"id" gorm:"autoIncrement;primaryKey"`
	Name   string
	Gender string
	Hobby  string
}

/*
创建数据库的时候：设置默认编码为utf8：
CREATE DATABASE `test` CHARACTER SET 'utf8' COLLATE 'utf8_general_ci';
*/

func NewMysqls() (*gorm.DB, error) {
	dsn := "root:123456@tcp(192.168.10.168:3306)/dbtest?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})

	//自动创建表结构到数据库
	//db.AutoMigrate(&UserInfo{})
	return db, nil
}

func CreateUserInfo() {
	db, err := NewMysqls()
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
		return
	}

	//u1 := UserInfo{Name: "tom", Gender: "男", Hobby: "篮球"}
	//u2 := UserInfo{Name: "luce", Gender: "女", Hobby: "足球"}
	//u3 := UserInfo{Name: "jack", Gender: "女", Hobby: "跳舞"}
	//u4 := UserInfo{Name: "john", Gender: "女", Hobby: "乒乓球"}
	//u5 := UserInfo{Name: "dany", Gender: "女", Hobby: "台球"}
	//u6 := UserInfo{Name: "trump", Gender: "女", Hobby: "橄榄球"}
	//u7 := UserInfo{Name: "lee", Gender: "女", Hobby: "足球"}
	//u8 := UserInfo{Name: "nmary", Gender: "女", Hobby: "羽毛球"}
	//
	//// 创建测试数据
	//users := []UserInfo{u1, u2, u3, u4, u5, u6, u7, u8}
	//for _, user := range users {
	//	db.Create(&user)
	//}

	//批量插入数据
	userInfo := []UserInfo{
		{Name: "like", Gender: "女", Hobby: "保龄球"},
		{Name: "julic", Gender: "女", Hobby: "排球"},
	}

	db.Create(userInfo)

}

func GetUserInfo() {
	db, err := NewMysqls()
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
		return
	}
	// 数据库查询所有的记录
	var users []UserInfo
	db.Find(&users)
	for _, user := range users {
		fmt.Printf("%v\n", user)
	}
}

func GetSingleUserInfo() {
	db, err := NewMysqls()
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
		return
	}
	// 获取第一条记录（主键升序）
	var user UserInfo
	db.First(&user)

	// 获取一条记录，没有指定排序字段
	//db.Take(&user)
	//SELECT * FROM users LIMIT 1;

	// 获取最后一条记录（主键降序）
	//db.Last(&user)
	// SELECT * FROM users ORDER BY id DESC LIMIT 1;
	fmt.Println(user)
}

func GetUserInfoByWhere() {
	db, err := NewMysqls()
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
		return
	}
	var users []UserInfo
	db.Where("hobby = ?", "足球").Find(&users)
	//第二种方式:
	//db.Find(&users, "hobby = ?", "足球")
	for _, user := range users {
		fmt.Printf("%v\n", user)
	}
}

func DeleteUserInfo() {
	db, err := NewMysqls()
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
		return
	}
	// 删除
	var user UserInfo
	db.Where("id = ?", 1).Delete(&user)
}

func UpdateUserInfo() {
	db, err := NewMysqls()
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
		return
	}
	var user UserInfo
	// 更新
	db.Model(&user).Where("name = ?", "luce").Update("hobby", "双色球")

	db.First(&user, "name = ?", "luce")
	fmt.Println(user)
}

func main() {
	//GetUserInfo()
	GetSingleUserInfo()
	//GetUserInfoByWhere()
	//UpdateUserInfo()
	//DeleteUserInfo()
	//CreateUserInfo()
	//GetUserInfo()
}
