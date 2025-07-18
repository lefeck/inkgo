package main

import (
	"gorm.io/gorm"
)

var DB *gorm.DB

//func init() {
//	dns := "root:123456@tcp(192.168.10.168:3306)/testuser?charset=utf8mb4&parseTime=True&loc=Local"
//	newLogger := logger.New(
//		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
//		logger.Config{
//			SlowThreshold: time.Second,   // 慢 SQL 阈值
//			LogLevel:      logger.Silent, // Log level
//			Colorful:      true,          // 禁用彩色打印
//		},
//	)
//
//	var err error
//	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
//		NamingStrategy: schema.NamingStrategy{
//			SingularTable: true,
//		},
//		Logger: newLogger,
//	})
//	if err != nil {
//		panic(err)
//	}
//	DB = db
//}

//func main() {
//	r := gin.Default()
//
//	r.GET("/v1/user/list", GetList)
//	r.Run()
//_ = db.AutoMigrate(&model.User{})

//user1 := model.User{gorm.Model{1}, "tom", "12345234", 1}
//user2 := model.User{gorm.Model{2}, "jack", "12345324", 1}
//db.Create(&user1)
//db.Create(&user2)
//
//// 查询
//var u = new(model.User)
//db.First(u)
//fmt.Printf("%#v\n", u)
//}
//
//func GetUserList(pageSize int, pageNum int) (int, []interface{}) {
//	var users []model.User
//	userList := make([]interface{}, 0, len(users))
//
//	if err := DB.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Error; err != nil {
//		return 0, nil
//	}
//	total := len(users)
//	for _, user := range users {
//		userItemMap := map[string]interface{}{
//			"id":       user.ID,
//			"name":     user.Name,
//			"password": user.Password,
//			"email   ": user.Email,
//			"authType": user.AuthType,
//			"authId  ": user.AuthId,
//			"avatar  ": user.Avatar,
//		}
//		userList = append(userList, userItemMap)
//	}
//	return total, userList
//}
//
//func GetList(c *gin.Context) {
//	userListForm := forms.UserListForm{}
//	if err := c.ShouldBind(&userListForm); err != nil {
//		c.JSON(http.StatusBadRequest, err)
//		return
//	}
//	fmt.Println(userListForm.PageSize, userListForm.PageNum)
//	total, userlist := GetUserList(userListForm.PageSize, userListForm.PageNum)
//	if (total + len(userlist)) == 0 {
//		common.Err(c, 400, 400, "未获取到到数据", gin.H{
//			"total":    total,
//			"userlist": userlist,
//		})
//		return
//	}
//	common.Success(c, 200, "获取用户列表成功", gin.H{
//		"total":    total,
//		"userlist": userlist,
//	})
//}
