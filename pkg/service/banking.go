package service

import (
	"context"

	"gosecureskeleton/pkg/dtos"
	"gosecureskeleton/pkg/errors"
	"gosecureskeleton/pkg/ext/db/sqlite"
	"gosecureskeleton/pkg/util"

	"github.com/sirupsen/logrus"
)

type BankingService struct {
	store *sqlite.Store
}

func NewBankingService(store *sqlite.Store) *BankingService {
	return &BankingService{store: store}
}

func (s *BankingService) Deposit(ctx context.Context, user dtos.User, amount int64) (dtos.UserResponse, error) {
	if err := s.store.UpdateUserBalance(ctx, user.ID, amount); err != nil {
		return dtos.UserResponse{}, err
	}

	updated, _, err := s.store.FindUserByUsername(ctx, user.Username)
	if err != nil {
		return dtos.UserResponse{}, err
	}

	util.LogInfo(ctx, "입금 처리 완료", logrus.Fields{"username": user.Username, "amount": amount})
	return dtos.MakeUserResponse(updated), nil
}

func (s *BankingService) Withdraw(ctx context.Context, user dtos.User, amount int64) (dtos.UserResponse, error) {
	if err := s.store.UpdateUserBalance(ctx, user.ID, -amount); err != nil {
		return dtos.UserResponse{}, err
	}

	updated, _, err := s.store.FindUserByUsername(ctx, user.Username)
	if err != nil {
		return dtos.UserResponse{}, err
	}

	util.LogInfo(ctx, "출금 처리 완료", logrus.Fields{"username": user.Username, "amount": amount})
	return dtos.MakeUserResponse(updated), nil
}

func (s *BankingService) Transfer(ctx context.Context, sender dtos.User, toUsername string, amount int64) (dtos.UserResponse, error) {
	receiver, ok, err := s.store.FindUserByUsername(ctx, toUsername)
	if err != nil {
		return dtos.UserResponse{}, err
	}
	if !ok {
		return dtos.UserResponse{}, errors.ErrUserNotFound
	}

	if err = s.store.UpdateUserBalance(ctx, sender.ID, -amount); err != nil {
		return dtos.UserResponse{}, err
	}
	if err = s.store.UpdateUserBalance(ctx, receiver.ID, amount); err != nil {
		return dtos.UserResponse{}, err
	}

	updated, _, err := s.store.FindUserByUsername(ctx, sender.Username)
	if err != nil {
		return dtos.UserResponse{}, err
	}

	util.LogInfo(ctx, "이체 처리 완료", logrus.Fields{
		"sender":   sender.Username,
		"receiver": toUsername,
		"amount":   amount,
	})
	return dtos.MakeUserResponse(updated), nil
}
