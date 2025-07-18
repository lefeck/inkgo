package model

// repost 转发文章
type Repost struct {
	ID           uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	Url          string `json:"url" gorm:"type:varchar(200);not null"` // 转发的文章链接
	IsAuthorized bool   `json:"is_authorized" gorm:"default:false"`    // 是否授权转载，默认 false
	RepostID     uint   `json:"repost_id"`                             // 被转发的文章ID
	Repost       Post   `json:"-" gorm:"foreignKey:RepostID"`          // 关联被转发的 Post 实体
}

func (*Repost) TableName() string {
	return "repost"
}
