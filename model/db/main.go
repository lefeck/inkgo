package main

import (
	"app/model"
	"app/repository"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

func NewMysql() (*gorm.DB, error) {
	dsn := "root:123456@tcp(192.168.10.143:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})

	//db.AutoMigrate(&model.Likes{}, &model.User{}, &model.Comment{}, &model.Post{}, &model.Category{}, &model.Tag{}, &model.Follow{}, &model.Activity{}, &model.Favorite{}, &model.Share{})
	return db, nil
}

func create_tag() {
	db := repository.NewRepository()
	tag1 := &model.Tag{
		Name: "golang",
	}
	tag2 := &model.Tag{
		Name: "python",
	}

	tag3 := &model.Tag{
		Name: "c",
	}
	tag4 := &model.Tag{
		Name: "c++",
	}

	tags := []*model.Tag{tag1, tag2, tag3, tag4}

	for _, tag := range tags {
		db.Create(&tag)
	}
}

func Get_tag() {
	db, _ := NewMysql()
	var tags []model.Tag
	db.Find(&tags, "name in ?", []string{"python"})
	fmt.Println(tags)
}

func Get_Tags() {
	db, _ := NewMysql()
	tags := make([]model.Tag, 0)
	db.Find(&tags)
	fmt.Println(tags)
}

func createTestSample() {
	db, err := NewMysql()
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
		return
	}

	author := model.User{
		Name:     "johndoe",
		Password: "johndoe123",
		Mobile:   "9876543210",
		Email:    "johndoe@example.com",
		Avatar:   "https://www.example.com/avatar.jpg",
	}
	db.Create(&author)

	tag := model.Tag{Name: "programming"}
	db.Create(&tag)

	category := model.Category{
		Name:        "webdev",
		Description: "web development",
		Image:       "https://www.example.com/category.jpg",
	}
	db.Create(&category)

	post1 := &model.Post{
		Title:    "Go Programming Intro",
		Content:  "Learn the basics of Go programming language...",
		AuthorID: author.ID,
		Tags:     []model.Tag{tag},
		//Categories: []model.Category{category},
		ViewCount: 20,
		LikeCount: 8,
		UserLiked: true,
		State:     0, // 0:正常发布，1:草稿箱
	}
	db.Create(&post1)

	comment1 := &model.Comment{
		Content:  "Great intro to Go!",
		AuthorID: author.ID,
		PostID:   post1.ID,
	}
	db.Create(&comment1)

	reply1 := &model.Comment{
		AuthorID: author.ID,
		Content:  "Thanks! I'm glad you liked it.",
		PostID:   post1.ID,
	}
	db.Create(&reply1)

	// Add reply1 to comment1 as replies
	comment1.Replies = append(comment1.Replies, *reply1)

	// 创建一个新的 like 对象
	like1 := &model.Likes{
		AuthorID: author.ID,
		PostID:   post1.ID,
	}

	db.Create(&like1)
}

func GetTagsByPost(postID uint) {
	db, err := NewMysql()
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
		return
	}

	//post := new(model.Post)
	//if err := db.First(post, postID).Error; err != nil {
	//	log.Fatalf("Failed to get the post: %v", err)
	//	return
	//}
	post := model.Post{ID: postID}
	tags := make([]model.Tag, 0)
	err = db.Model(post).Association(model.TagsAssociation).Find(&tags)
	if err != nil {
		fmt.Errorf("error :%v", err)
	}
	fmt.Println(tags)
}

func GetAllTags() {
	db, err := NewMysql()
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
		return
	}

	tags := make([]model.Tag, 0)
	err = db.Find(&tags).Error
	if err != nil {
		fmt.Errorf("error :%v", err)
	}
	fmt.Println(tags)
}

func createTestSample1() {
	db, err := NewMysql()
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
		return
	}

	author := model.User{
		Name:     "janedoe",
		Password: "janedoe123",
		Mobile:   "1234567890",
		Email:    "janedoe@example.com",
		Avatar:   "https://www.example.com/avatar2.jpg",
	}
	db.Create(&author)

	tag1 := model.Tag{Name: "mobiledev"}
	db.Create(&tag1)

	category1 := model.Category{
		Name:        "ios",
		Description: "iOS development",
		Image:       "https://www.example.com/ios-category.jpg",
	}
	db.Create(&category1)

	post1 := &model.Post{
		Title:    "Swift Programming Basics",
		Content:  "Learn the foundations of Swift programming...",
		AuthorID: author.ID,
		Tags:     []model.Tag{tag1},
		//Categories: []model.Category{category1},
		ViewCount: 15,
		LikeCount: 6,
		UserLiked: true,
		State:     0,
	}
	db.Create(&post1)

	comment1 := &model.Comment{
		Content:  "Nice intro to Swift!",
		AuthorID: author.ID,
		PostID:   post1.ID,
	}
	db.Create(&comment1)

	reply1 := &model.Comment{
		AuthorID: author.ID,
		Content:  "Thank you! I'm glad you found it helpful.",
		PostID:   post1.ID,
	}
	db.Create(&reply1)

	comment1.Replies = append(comment1.Replies, *reply1)

	like1 := &model.Likes{
		AuthorID: author.ID,
		PostID:   post1.ID,
	}
	db.Create(&like1)
}

func createTestSample2() {
	db, err := NewMysql()
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
		return
	}

	author := model.User{
		Name:     "alice",
		Password: "alice123",
		Mobile:   "1112223334",
		Email:    "alice@example.com",
		Avatar:   "https://www.example.com/avatar3.jpg",
	}
	db.Create(&author)

	tag1 := model.Tag{Name: "frontend"}
	tag2 := model.Tag{Name: "react"}
	db.Create(&tag1)
	db.Create(&tag2)

	category1 := model.Category{
		Name:        "javascript",
		Description: "JavaScript development",
		Image:       "https://www.example.com/js-category.jpg",
	}
	db.Create(&category1)

	post1 := &model.Post{
		Title:    "Starting with React",
		Content:  "Learn to build React applications from scratch...",
		AuthorID: author.ID,
		Tags:     []model.Tag{tag1, tag2},
		//Categories: []model.Category{category1},
		ViewCount: 44,
		LikeCount: 10,
		UserLiked: true,
		State:     0,
	}
	db.Create(&post1)

	comment1 := &model.Comment{
		Content:  "Great content for React beginners!",
		AuthorID: author.ID,
		PostID:   post1.ID,
	}
	db.Create(&comment1)

	reply1 := &model.Comment{
		AuthorID: author.ID,
		Content:  "Thanks! I'm happy to help!",
		PostID:   post1.ID,
	}
	db.Create(&reply1)

	comment1.Replies = append(comment1.Replies, *reply1)

	like1 := &model.Likes{
		AuthorID: author.ID,
		PostID:   post1.ID,
	}
	db.Create(&like1)
}

func Get_post() {
	db, _ := NewMysql()
	// 添加文章选择标签
	//多对对查询

	// 查询文章, 显示文章的标签列表
	var post model.Post
	db.Preload("Tags").Find(&post, 1)
	fmt.Println(post)

	// 查询标签, 显示文章的列表
	var tag model.Tag
	db.Preload("Posts").Take(&tag, 3)
	fmt.Println(tag)
}

func countview() {
	db, err := NewMysql()
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
		return
	}

	var posts []model.Post
	// 按阅读量降序排列文章
	result := db.Order("view_count desc").Find(&posts)

	if result.Error != nil {
		log.Fatalf("Failed to get the posts ordered by view_count: %v", result.Error)
	}

	// 打印结果
	fmt.Printf("Posts ordered by ViewCount:\n")
	for i, post := range posts {
		fmt.Printf("%d. %s (ViewCount: %d)\n", i+1, post.Title, post.ViewCount)
	}
}

func creationtime() {
	db, err := NewMysql()
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
		return
	}

	var posts []model.Post
	// 按创建时间降序排列文章
	result := db.Preload("Author").Order("created_at desc").Find(&posts)

	if result.Error != nil {
		log.Fatalf("Failed to get the posts ordered by created_at: %v", result.Error)
	}

	// 打印结果
	fmt.Printf("Posts ordered by CreatedAt:\n")
	for i, post := range posts {
		fmt.Printf("%d. %s (CreatedAt: %s, Author: %s)\n", i+1, post.Title, post.BaseModel.CreatedAt.Format("2006-01-02 15:04:05"), post.Author.Name)
	}
}

// 统计文章的点赞数
func countlike(postID uint) {
	db, err := NewMysql()
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
		return
	}

	Post := new(model.Post)
	//文章关联了,creator, tag, category, comments, comment.user(评论的用户)
	if err := db.Preload(model.CreatorAssociation).Preload(model.TagsAssociation).Preload(model.CategoryAssociation).Preload("Comments.Author").Preload(model.CommentsAssociation).Preload(model.LikesAssociation).First(Post, postID).Error; err != nil {
		log.Fatalf("Failed to get the post: %v", err)
	}

	// 计算并更新 Post.LikeCount
	Post.LikeCount = uint(Post.LikeCount)

	fmt.Printf("count like %d\n", Post.LikeCount)
}

// 获取所有评论数
func countComments() {
	db, err := NewMysql()
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
		return
	}

	var count int64
	if err := db.Model(&model.Comment{}).Where("deleted_at IS NULL").Count(&count).Error; err != nil {
		log.Fatalf("Failed to count comments: %v", err)
	}

	fmt.Printf("Total comments: %d\n", count)
}

// 获取所有的posts
func listposts() {
	posts := make([]model.Post, 0)
	//db.Omit 在create，update和query 忽略哪些字段
	/*
		order 从数据库检索记录时指定顺序
		db.Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true})
		// 表示按照创建时间顺序来获取输出
	*/
	db, err := NewMysql()
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
		return
	}
	if err := db.Omit("content").Preload(model.CreatorAssociation).Preload(model.TagsAssociation).Preload(model.CategoryAssociation).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).Where(" state = ?", 0).Find(&posts).Error; err != nil {
		log.Fatalf("Failed to create connection: %v", err)
	}

	ids := make([]uint, len(posts))
	for i, Post := range posts {
		ids[i] = Post.ID
	}

	type result struct {
		ID    uint
		Likes uint
	}

	results := []result{}
	if err := db.Model(&model.Likes{}).Select("post_id as id, count(like.post_id) as likes").Where("post_id in ?", ids).Group("post_id").Scan(&results).Error; err != nil {
		log.Fatalf("Failed to create connection: %v", err)
	}

	resMap := make(map[uint]uint, len(results))
	for _, r := range results {
		resMap[r.ID] = r.Likes
	}

	for i := range posts {
		posts[i].LikeCount = resMap[posts[i].ID]
	}
	fmt.Println(posts)
}

func ListDrafts() {
	posts := make([]model.Post, 0)
	//db.Omit 在create，update和query 忽略哪些字段
	/*
		order 从数据库检索记录时指定顺序
		db.Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true})
		// 表示按照创建时间顺序来获取输出
	*/
	db, err := NewMysql()
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
		return
	}
	if err := db.Omit("content").Preload(model.CreatorAssociation).Preload(model.TagsAssociation).Preload(model.CategoryAssociation).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).Where(" state = ?", 1).Find(&posts).Error; err != nil {
		log.Fatalf("Failed to create connection: %v", err)
	}
	fmt.Println(posts)
}

func getcomments() {
	db, err := NewMysql()
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
		return
	}
	comments := make([]model.Comment, 0)
	err = db.Where("post_id = ?", 3).Find(&comments).Error
	if err != nil {
		log.Fatalf("Failed to count comments: %v", err)
	}

	b, _ := json.Marshal(comments)

	os.Stdout.Write(b)

	//fmt.Println(comments)

}

func main() {

	NewMysql()

	//create_tag()
	Get_tag()
	//Get_Tags()

	//create_post()
	//Get_post()

	//createTestSample()
	//fmt.Println("Test Sample 0 created.")
	//createTestSample1()
	//fmt.Println("Test Sample 1 created.")
	//createTestSample2()
	//fmt.Println("Test Sample 2 created.")
	//
	//countview()
	//
	//creationtime()

	//likecount()

	//countlike(3)
	//countComments()

	//listposts()

	//GetTagsByPost(3)
	//
	//GetTags()

	//getcomments()

	//ListDrafts()
}
