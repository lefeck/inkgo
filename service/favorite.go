package service

import (
	"inkgo/model"
	"inkgo/repository"
	"strconv"
)

type FavoriteService interface {
	CreateFavorite(uid string, favorite *model.Favorite) (*model.Favorite, error)
	DeleteFavorite(uid string, id string) error
	UpdateFavorite(uid string, favorite *model.Favorite) (*model.Favorite, error)
	GetFavoriteByID(uid, fid string) (*model.Favorite, error)
	// GetUserFavorites 获取某个用户创建的所有收藏夹列表
	GetUserFavorites(userID string, page, pageSize int) ([]model.Favorite, int64, error)

	AddPostToFavorite(uid string, pid string, fid string) error
	RemovePostFromFavorite(uid string, pid string, fid string) error
	IsPostInFavorite(uid string, fid, pid string) (bool, error)

	ListPostsInFavorite(uid string, fid string) ([]model.Post, error) // 列出收藏夹中的文章
	CountPostsInFavorite(uid string, fid string) (int64, error)       // 获取收藏夹中的文章数量
}

type favoriteService struct {
	favoriteRepository repository.FavoriteRepository
}

func NewFavoriteService(favoriteRepository repository.FavoriteRepository) FavoriteService {
	return &favoriteService{
		favoriteRepository: favoriteRepository,
	}
}

func (f *favoriteService) CreateFavorite(uid string, favorite *model.Favorite) (*model.Favorite, error) {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return nil, err
	}
	favorite.UserID = uint(uidInt)
	return f.favoriteRepository.CreateFavorite(uint(uidInt), favorite)
}

func (f *favoriteService) DeleteFavorite(uid string, id string) error {
	fidInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return err
	}
	return f.favoriteRepository.DeleteFavorite(uint(uidInt), uint(fidInt))
}

func (f *favoriteService) UpdateFavorite(uid string, favorite *model.Favorite) (*model.Favorite, error) {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return nil, err
	}
	return f.favoriteRepository.UpdateFavorite(uint(uidInt), favorite)
}

func (f *favoriteService) GetFavoriteByID(uid, fid string) (*model.Favorite, error) {
	fidInt, err := strconv.Atoi(fid)
	if err != nil {
		return nil, err
	}
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return nil, err
	}
	return f.favoriteRepository.GetFavoriteByID(uint(uidInt), uint(fidInt))
}

// GetUserFavorites retrieves favorites for a specific user.
func (f *favoriteService) GetUserFavorites(uid string, page int, pageSize int) ([]model.Favorite, int64, error) {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return nil, 0, err
	}
	favorites, total, err := f.favoriteRepository.GetUserFavorites(uint(uidInt), page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return favorites, total, nil
}

func (f *favoriteService) AddPostToFavorite(uid, pid string, fid string) error {
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		return err
	}
	fidInt, err := strconv.Atoi(fid)
	if err != nil {
		return err
	}
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return err
	}
	return f.favoriteRepository.AddPostToFavorite(uint(uidInt), uint(pidInt), uint(fidInt))
}

func (f *favoriteService) RemovePostFromFavorite(uid, pid string, fid string) error {
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		return err
	}
	fidInt, err := strconv.Atoi(fid)
	if err != nil {
		return err
	}
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return err
	}
	return f.favoriteRepository.RemovePostFromFavorite(uint(uidInt), uint(pidInt), uint(fidInt))
}

func (f *favoriteService) IsPostInFavorite(uid, fid, pid string) (bool, error) {
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		return false, err
	}
	fidInt, err := strconv.Atoi(fid)
	if err != nil {
		return false, err
	}
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false, err
	}
	return f.favoriteRepository.IsPostInFavorite(uint(uidInt), uint(fidInt), uint(pidInt))
}

func (f *favoriteService) ListPostsInFavorite(uid, fid string) ([]model.Post, error) {
	fidInt, err := strconv.Atoi(fid)
	if err != nil {
		return nil, err
	}
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return nil, err
	}
	posts, err := f.favoriteRepository.ListPostsInFavorite(uint(uidInt), uint(fidInt))
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (f *favoriteService) CountPostsInFavorite(uid, fid string) (int64, error) {
	fidInt, err := strconv.Atoi(fid)
	if err != nil {
		return 0, err
	}
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return 0, err
	}
	count, err := f.favoriteRepository.CountPostsInFavorite(uint(uidInt), uint(fidInt))
	if err != nil {
		return 0, err
	}
	return count, nil
}
