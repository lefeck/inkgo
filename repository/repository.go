package repository

import (
	"context"
	"gorm.io/gorm"
	"inkgo/database"
)

type repository struct {
	user     UserRepository
	post     PostRepository
	category CategoryRepository
	comment  CommentRepository
	tag      TagRepository
	like     LikeRepository
	activity ActivityRepository
	favorite FavoriteRepository
	follow   FollowRepository
	share    ShareRepository
	token    TokenRepository
	auth     AuthRepository // 假设有一个 AuthRepository 接口
	db       *gorm.DB
	rdb      *database.RedisDB
	migrants []Migrant
}

func NewRepository(db *gorm.DB, rdb *database.RedisDB) Repository {
	r := &repository{
		user:     NewUserRepository(db, rdb),
		post:     NewPostRepository(db),
		tag:      NewTagRepository(db),
		comment:  NewCommentRepository(db),
		like:     NewLikeRepository(db),
		category: NewCategoryRepository(db),
		activity: NewActivityRepository(db),
		follow:   NewFollowRepository(db),
		favorite: NewFavoriteRepository(db),
		share:    NewShareRepository(db),
		token:    NewTokenRepository(rdb),
		auth:     NewAuthRepository(rdb),
		db:       db,
		rdb:      rdb,
	}
	r.migrants = getMigrants(
		r.user,
		r.Post,
		r.like,
		r.category,
		r.comment,
		r.tag,
		r.follow,
		r.activity,
		r.favorite,
		r.share,
		r.auth,
		r.token,
	)
	return r
}

// getMigrants 从传入的对象中提取实现了 Migrant 接口的对象
func getMigrants(objs ...interface{}) []Migrant {
	var migrants []Migrant
	for _, obj := range objs {
		if m, ok := obj.(Migrant); ok {
			migrants = append(migrants, m)
		}
	}
	return migrants
}

// Migrate 执行所有的迁移操作
func (r *repository) Migrate() error {
	for _, m := range r.migrants {
		if err := m.Migrate(); err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) User() UserRepository {
	return r.user
}

func (r *repository) Auth() AuthRepository {
	return r.auth
}

// Token 获取 TokenRepository
func (r *repository) Token() TokenRepository {
	return r.token
}

func (r *repository) Post() PostRepository {
	return r.post
}

func (r *repository) Comment() CommentRepository {
	return r.comment
}

func (r *repository) Category() CategoryRepository {
	return r.category
}

func (r *repository) Like() LikeRepository {
	return r.like
}

func (r *repository) Tag() TagRepository {
	return r.tag
}

func (r *repository) Activity() ActivityRepository {
	return r.activity
}

func (r *repository) Favorite() FavoriteRepository {
	return r.favorite
}

func (r *repository) Follow() FollowRepository {
	return r.follow
}

func (r *repository) Share() ShareRepository {
	return r.share
}

func (r *repository) Close() error {
	db, _ := r.db.DB()
	if db != nil {
		if err := db.Close(); err != nil {
			return err
		}
	}
	if r.rdb != nil {
		if err := r.rdb.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) Ping(ctx context.Context) error {
	db, _ := r.db.DB()
	if db != nil {
		if err := db.PingContext(ctx); err != nil {
			return err
		}
	}
	if r.rdb != nil {
		if _, err := r.rdb.Ping(ctx).Result(); err != nil {
			return err
		}
	}
	return nil
}
