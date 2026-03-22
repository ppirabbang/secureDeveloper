package service

import (
	"context"

	"gosecureskeleton/pkg/dtos"
	"gosecureskeleton/pkg/errors"
	"gosecureskeleton/pkg/ext/db/sqlite"
	"gosecureskeleton/pkg/util"

	"github.com/sirupsen/logrus"
)

type PostService struct {
	store *sqlite.Store
}

func NewPostService(store *sqlite.Store) *PostService {
	return &PostService{store: store}
}

func (s *PostService) List(ctx context.Context) ([]dtos.PostView, error) {
	return s.store.ListPosts(ctx)
}

func (s *PostService) Get(ctx context.Context, id uint) (dtos.PostView, error) {
	post, ok, err := s.store.GetPost(ctx, id)
	if err != nil {
		return dtos.PostView{}, err
	}
	if !ok {
		return dtos.PostView{}, errors.ErrPostNotFound
	}
	return post, nil
}

func (s *PostService) Create(ctx context.Context, ownerID uint, title, content string) (dtos.PostView, error) {
	post, err := s.store.CreatePost(ctx, ownerID, title, content)
	if err != nil {
		return dtos.PostView{}, err
	}

	util.LogInfo(ctx, "게시글 생성 완료", logrus.Fields{"owner_id": ownerID, "post_id": post.ID})
	return post, nil
}

func (s *PostService) Update(ctx context.Context, postID, ownerID uint, title, content string) (dtos.PostView, error) {
	post, ok, err := s.store.GetPost(ctx, postID)
	if err != nil {
		return dtos.PostView{}, err
	}
	if !ok {
		return dtos.PostView{}, errors.ErrPostNotFound
	}
	if post.OwnerID != ownerID {
		return dtos.PostView{}, errors.ErrUnauthorizedAction
	}

	updated, err := s.store.UpdatePost(ctx, postID, title, content)
	if err != nil {
		return dtos.PostView{}, err
	}

	util.LogInfo(ctx, "게시글 수정 완료", logrus.Fields{"post_id": postID, "owner_id": ownerID})
	return updated, nil
}

func (s *PostService) Delete(ctx context.Context, postID, ownerID uint) error {
	post, ok, err := s.store.GetPost(ctx, postID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.ErrPostNotFound
	}
	if post.OwnerID != ownerID {
		return errors.ErrUnauthorizedAction
	}

	if err = s.store.DeletePost(ctx, postID); err != nil {
		return err
	}

	util.LogInfo(ctx, "게시글 삭제 완료", logrus.Fields{"post_id": postID, "owner_id": ownerID})
	return nil
}
