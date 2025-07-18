package service

import (
	"inkgo/model"
	"inkgo/repository"
	"strconv"
)

type likeService struct {
	likeRepository repository.LikeRepository
	//postRepository repository.PostRepository
}

func NewLikeService(likeRepository repository.LikeRepository) LikeService {
	return &likeService{
		likeRepository: likeRepository,
	}
}

func (l *likeService) LikePost(user *model.User, id string) error {
	pid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return l.likeRepository.LikePost(uint(pid), user.ID)
}

func (l *likeService) UnLikePost(user *model.User, id string) error {
	pid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return l.likeRepository.UnLikePost(uint(pid), user.ID)
}

func (l *likeService) CountLikes(id string) (int64, error) {
	pid, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	count, err := l.likeRepository.CountLikes(uint(pid))
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (l *likeService) IsLiked(user *model.User, id string) (bool, error) {
	pid, err := strconv.Atoi(id)
	if err != nil {
		return false, err
	}
	return l.likeRepository.IsLiked(user.ID, uint(pid))
}
