package service

import (
	"inkgo/model"
	"inkgo/repository"
	"strconv"
)

type commentService struct {
	commentRepository repository.CommentRepository
}

func NewCommentService(commentRepository repository.CommentRepository) CommentService {
	return &commentService{
		commentRepository: commentRepository,
	}
}

func (c *commentService) Add(comment *model.Comment, id string, user *model.User) (*model.Comment, error) {
	cid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	comment.ID = uint(cid)
	comment.AuthorID = user.ID

	return c.commentRepository.Add(comment)
}

func (c *commentService) Delete(id string) error {
	return c.commentRepository.Delete(id)
}

func (c *commentService) List(aid string) ([]model.Comment, error) {
	return c.commentRepository.List(aid)
}
