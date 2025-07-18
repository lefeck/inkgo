package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"inkgo/model"
)

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{
		db: db,
	}
}

// 通过id获取文章
func (a *postRepository) GetPostByID(id uint) (*model.Post, error) {
	Post := new(model.Post)
	//文章关联了,creator, tag, category, comments, comment.user(评论的用户)
	// gorm.Preload() 会自动填充结构体的关联字段，gin 会自动将结构体转为 JSON，
	if err := a.db.Preload(model.AuthorAssociation).Preload(model.TagsAssociation).Preload(model.CategoryAssociation).
		Preload("Comments.Author").Preload(model.CommentsAssociation).
		Find(Post).Error; err != nil {
		return nil, err
	}
	// 获取文章的点赞数
	like, err := a.CountLike(id)
	if err != nil {
		return nil, err
	}
	Post.LikeCount = uint(like)
	// 文章阅读量自增
	if err := a.IncView(id); err != nil {
		return nil, err
	}

	return Post, nil
}

// 按创建时间降序排列文章
func (p *postRepository) SortByTimeOfCreation() ([]model.Post, error) {
	var posts []model.Post
	// Preload 预加载关联数据，Order 排序
	//"created_at desc" 是 SQL 中的排序语句，表示按照创建时间降序排列
	// select * from posts order by created_at desc
	if err := p.db.Preload(model.AuthorAssociation).Order("created_at desc").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// 按阅读量降序排列文章
func (p *postRepository) SortByViewCount() ([]model.Post, error) {
	var posts []model.Post
	if err := p.db.Order("view_count desc").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// 按点赞降序排列文章
func (p *postRepository) SortByViewLike() ([]model.Post, error) {
	var posts []model.Post
	if err := p.db.Order("like_count desc").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// 统计文章的点赞数
func (p *postRepository) CountLike(id uint) (int64, error) {
	var count int64
	// select count(*) from likes where post_id = ?
	if err := p.db.Model(&model.Likes{}).Where("post_id = ?", id).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// 通过名字获取文章,这是个搜索功能
func (a *postRepository) GetPostByName(title string) (*model.Post, error) {
	post := new(model.Post)
	if err := a.db.Where("title = ?", title).First(post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

// 获取所有的已发布文章
func (p *postRepository) ListHasPublished(page, pageSize int) ([]model.Post, int64, error) {
	posts := make([]model.Post, 0)

	var total int64
	//db.Omit 在create，update和query 忽略哪些字段
	/*
		order 从数据库检索记录时指定顺序
		db.Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true})
		// 表示按照创建时间顺序来获取输出
	*/

	db := p.db.Model(&model.Post{})
	// 计算总数
	if err := db.Where(" state = ?", model.PostPublished).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := p.db.Omit("content").Preload(model.AuthorAssociation).Preload(model.TagsAssociation).Preload(model.CategoryAssociation).
		Where(" state = ?", model.PostPublished).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&posts).Error; err != nil {
		return nil, 0, err
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
	if err := p.db.Model(&model.Likes{}).Select("post_id as id, count(likes.post_id) as likes").Where("post_id in ?", ids).Group("post_id").Scan(&results).Error; err != nil {
		return nil, 0, err
	}

	resMap := make(map[uint]uint, len(results))
	for _, r := range results {
		resMap[r.ID] = r.Likes
	}

	for i := range posts {
		posts[i].LikeCount = resMap[posts[i].ID]
	}
	return posts, total, nil
}

// 获取所有的未发布文章
func (p *postRepository) ListDrafts(page, pageSize int) ([]model.Post, int64, error) {
	posts := make([]model.Post, 0)
	var total int64

	db := p.db.Model(&model.Post{})
	// 计算总数
	if err := db.Where(" state = ?", model.PostDraft).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	//db.Omit 在create，update和query 忽略哪些字段
	if err := p.db.Omit("content").Preload(model.AuthorAssociation).Preload(model.TagsAssociation).Preload(model.CategoryAssociation).
		Where(" state = ?", model.PostDraft).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

// 创建文章
func (a *postRepository) Create(user *model.User, Post *model.Post) (*model.Post, error) {
	// 检查作者是否存在
	if result := a.db.First(&user, Post.AuthorID); result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound // 作者不存在
	}

	err := a.db.Create(&Post).Error
	return Post, err
}

// 更新文章
func (a *postRepository) Update(post *model.Post) (*model.Post, error) {
	var existing model.Post
	// 检查文章是否存在
	if result := a.db.First(&existing, post.ID); result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound // 文章不存在
	}
	// 我这omit 添了哪些字段, 但是我update中传了这些字段,还是被修改了为啥?
	err := a.db.Model(&existing).Omit("view_count", "author_id", "view_count", "like_count", "user_liked", "state").Updates(&post).Error
	if err != nil {
		return nil, err
	}
	return post, err
}

// 更新文章的状态, 例如从草稿到已发布
func (p *postRepository) UpdateStatus(id uint, state model.PostState) (*model.Post, error) {
	Post := &model.Post{Model: gorm.Model{ID: id}}
	// 检查文章是否存在
	if result := p.db.First(Post); result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound // 文章不存在
	}

	// 更新文章状态
	if err := p.db.Model(Post).Update("state", state).Error; err != nil {
		return nil, err
	}
	return Post, nil
}

// 删除文章
func (a *postRepository) Delete(id uint) error {
	post := &model.Post{}
	// 检查文章是否存在
	if result := a.db.First(post, id); result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // 文章不存在
	}

	return a.db.Delete(post, id).Error
}

// 文章阅读数自增
func (p *postRepository) IncView(id uint) error {
	Post := &model.Post{Model: gorm.Model{ID: id}}
	// select * from posts where id = ?
	return p.db.Model(Post).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// 统计已经发布文章个数
func (p *postRepository) CountNumberOfPost() (int64, error) {
	var count int64
	// 根据发布状态统计数量
	post := &model.Post{}
	if err := p.db.Model(post).Where("state = ?", model.PostPublished).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// 统计未发布文章个数
func (p *postRepository) CountNumberOfDraft() (int64, error) {
	var count int64
	// 根据未发布状态统计数量
	post := &model.Post{}
	if err := p.db.Model(post).Where("state = ?", model.PostDraft).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// 根据文章阅读量从高到低排序输出
func (p *postRepository) SortByViewCountDesc(page int, pageSize int) ([]model.Post, int64, error) {
	posts := make([]model.Post, 0)
	var total int64
	// 计算总数
	if err := p.db.Model(&model.Post{}).Where("state = ?", model.PostPublished).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 按照阅读量降序排列
	if err := p.db.Order("view_count desc").
		Preload(model.AuthorAssociation).Preload(model.TagsAssociation).Preload(model.CategoryAssociation).
		Where("state = ?", model.PostPublished).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

/*
✔ 适合的场景：
首页展示“最近发布的前 N 篇文章”
推荐模块、侧边栏模块（如“最新更新”）
固定数量展示
*/
// 列出热门文章
func (p *postRepository) ListHotPosts(limit int) ([]model.Post, error) {
	posts := make([]model.Post, 0)
	err := p.db.Where("state = ?", model.PostPublished).Order("view_count desc").Limit(limit).Find(&posts).Error
	return posts, err
}

// 列出最近的文章
func (p *postRepository) ListRecentPosts(limit int) ([]model.Post, error) {
	posts := make([]model.Post, 0)
	err := p.db.Where("state = ?", model.PostPublished).Order("created_at desc").Limit(limit).Find(&posts).Error
	return posts, err
}

// 自动创建表结构到db
func (a *postRepository) Migrate() error {
	return a.db.AutoMigrate(&model.Post{})
}
