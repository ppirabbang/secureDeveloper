package service

import (
	"context"

	"gosecureskeleton/pkg/dtos"
	"gosecureskeleton/pkg/errors"
	"gosecureskeleton/pkg/ext/db/sqlite"
	"gosecureskeleton/pkg/session"
	"gosecureskeleton/pkg/util"

	"github.com/sirupsen/logrus"
)

type AuthService struct {
	store    *sqlite.Store
	sessions *session.Store
}

func NewAuthService(store *sqlite.Store, sessions *session.Store) *AuthService {
	return &AuthService{store: store, sessions: sessions}
}

func (s *AuthService) Register(ctx context.Context, req dtos.RegisterRequest) (dtos.UserResponse, error) {
	_, exists, err := s.store.FindUserByUsername(ctx, req.Username)
	if err != nil {
		return dtos.UserResponse{}, err
	}
	if exists {
		return dtos.UserResponse{}, errors.ErrUserAlreadyExists
	}

	user := dtos.User{
		Username: req.Username,
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	}
	if err = s.store.CreateUser(ctx, user); err != nil {
		return dtos.UserResponse{}, err
	}

	created, _, err := s.store.FindUserByUsername(ctx, req.Username)
	if err != nil {
		return dtos.UserResponse{}, err
	}

	util.LogInfo(ctx, "회원가입 완료", logrus.Fields{"username": req.Username})
	return dtos.MakeUserResponse(created), nil
}

func (s *AuthService) Login(ctx context.Context, req dtos.LoginRequest) (dtos.User, string, error) {
	user, ok, err := s.store.FindUserByUsername(ctx, req.Username)
	if err != nil {
		return dtos.User{}, "", err
	}
	if !ok || user.Password != req.Password {
		return dtos.User{}, "", errors.ErrInvalidCredentials
	}

	token, err := s.sessions.Create(user.ID)
	if err != nil {
		return dtos.User{}, "", err
	}

	util.LogInfo(ctx, "로그인 완료", logrus.Fields{"username": req.Username})
	return user, token, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	userID, ok := s.sessions.Lookup(token)
	if !ok {
		return errors.ErrInvalidToken
	}
	s.sessions.Delete(token)
	util.LogInfo(ctx, "로그아웃 완료", logrus.Fields{"user_id": userID})
	return nil
}

func (s *AuthService) WithdrawAccount(ctx context.Context, token string, password string) (dtos.UserResponse, error) {
	userID, ok := s.sessions.Lookup(token)
	if !ok {
		return dtos.UserResponse{}, errors.ErrInvalidToken
	}

	user, found, err := s.store.FindUserByID(ctx, userID)
	if err != nil {
		return dtos.UserResponse{}, err
	}
	if !found {
		return dtos.UserResponse{}, errors.ErrUserNotFound
	}
	if user.Password != password {
		return dtos.UserResponse{}, errors.ErrPasswordMismatch
	}

	if err = s.store.DeleteUser(ctx, user.ID); err != nil {
		return dtos.UserResponse{}, err
	}

	s.sessions.Delete(token)
	util.LogInfo(ctx, "회원탈퇴 완료", logrus.Fields{"username": user.Username})
	return dtos.MakeUserResponse(user), nil
}
