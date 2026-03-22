package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"gosecureskeleton/pkg/dtos"
)

const queryUserByUserName = `
	SELECT id, username, name, email, phone, password, balance, is_admin
	FROM users
	WHERE username = ?
`

func (s *Store) FindUserByUsername(ctx context.Context, username string) (dtos.User, bool, error) {
	row := s.queryRow(ctx, queryUserByUserName, strings.TrimSpace(username))

	var user dtos.User
	var isAdmin int64
	if err := row.Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Phone, &user.Password, &user.Balance, &isAdmin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dtos.User{}, false, nil
		}
		return dtos.User{}, false, err
	}
	user.IsAdmin = isAdmin == 1

	return user, true, nil
}

const queryUserByID = `
	SELECT id, username, name, email, phone, password, balance, is_admin
	FROM users
	WHERE id = ?
`

func (s *Store) FindUserByID(ctx context.Context, id uint) (dtos.User, bool, error) {
	row := s.queryRow(ctx, queryUserByID, id)

	var user dtos.User
	var isAdmin int64
	if err := row.Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Phone, &user.Password, &user.Balance, &isAdmin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dtos.User{}, false, nil
		}
		return dtos.User{}, false, err
	}
	user.IsAdmin = isAdmin == 1

	return user, true, nil
}

const dmlInsertUsers = `
	INSERT INTO users (username, name, email, phone, password, balance, is_admin)
	VALUES (?, ?, ?, ?, ?, ?, ?)
`

func (s *Store) CreateUser(ctx context.Context, user dtos.User) error {
	_, err := s.exec(ctx, dmlInsertUsers, user.Username, user.Name, user.Email, user.Phone, user.Password, user.Balance, boolToInt(user.IsAdmin))
	return err
}

const dmlDeleteUserByID = `
	DELETE FROM users 
	WHERE id = ?
`

func (s *Store) DeleteUser(ctx context.Context, userID uint) error {
	_, err := s.exec(ctx, dmlDeleteUserByID, userID)
	return err
}

const dmlUpdateUser = `
	UPDATE users 
	SET balance = balance + ? 
	WHERE id = ?
`

func (s *Store) UpdateUserBalance(ctx context.Context, userID uint, amount int64) error {
	_, err := s.exec(ctx, dmlUpdateUser, amount, userID)
	return err
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
