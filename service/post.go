package service

import (
	"inkgo/model"
	"inkgo/repository"
	"strconv"
)

type postService struct {
	postRepository repository.PostRepository
	likeRepository repository.LikeRepository
}

func NewPostService(postRepository repository.PostRepository) PostService {
	return &postService{
		postRepository: postRepository,
	}
}

func (p *postService) GetPostByID(id string) (*model.Post, error) {
	aid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	post, err := p.postRepository.GetPostByID(uint(aid))
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (p *postService) GetPostByName(name string) (*model.Post, error) {
	post, err := p.postRepository.GetPostByName(name)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (p *postService) HasPublished(page, pageSize int) ([]model.Post, int64, error) {
	posts, total, err := p.postRepository.ListHasPublished(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return posts, total, nil

}

func (p *postService) ListDrafts(page, pageSize int) ([]model.Post, int64, error) {
	posts, total, err := p.postRepository.ListDrafts(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

func (p *postService) UpdateStatus(id string, state model.PostState) (*model.Post, error) {
	aid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	post, err := p.postRepository.UpdateStatus(uint(aid), state)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (p *postService) SortByViewCountDesc(page int, pageSize int) ([]model.Post, int64, error) {
	posts, total, err := p.postRepository.SortByViewCountDesc(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

func (p *postService) ListHotPosts(limit int) ([]model.Post, error) {
	posts, err := p.postRepository.ListHotPosts(limit)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (p *postService) ListRecentPosts(limit int) ([]model.Post, error) {
	posts, err := p.postRepository.ListRecentPosts(limit)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (p *postService) Create(user *model.User, Post *model.Post) (*model.Post, error) {
	return p.postRepository.Create(user, Post)
}

func (p *postService) Get(user *model.User, id string) (*model.Post, error) {
	aid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	if err := p.postRepository.IncView(uint(aid)); err != nil {
		return nil, err
	}
	Post, err := p.postRepository.GetPostByID(uint(aid))
	if err != nil {
		return nil, err
	}
	Post.UserLiked, _ = p.likeRepository.IsLiked(uint(aid), user.ID)

	return Post, nil
}

func (p *postService) Update(id string, Post *model.Post) (*model.Post, error) {
	aid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	Post.ID = uint(aid)
	return p.postRepository.Update(Post)
}

func (p *postService) Delete(id string) error {
	aid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return p.postRepository.Delete(uint(aid))
}
