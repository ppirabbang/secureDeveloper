package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"gosecureskeleton/pkg/dtos"
)

const queryAllPosts = `
	SELECT p.id, p.title, p.content, p.owner_id, u.name, u.email, p.created_at, p.updated_at
	FROM posts p
	JOIN users u ON p.owner_id = u.id
	ORDER BY p.created_at DESC
`

func (s *Store) ListPosts(ctx context.Context) ([]dtos.PostView, error) {
	rows, err := s.query(ctx, queryAllPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []dtos.PostView
	for rows.Next() {
		var p dtos.PostView
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.OwnerID, &p.Author, &p.AuthorEmail, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

const queryPostByID = `
	SELECT p.id, p.title, p.content, p.owner_id, u.name, u.email, p.created_at, p.updated_at
	FROM posts p
	JOIN users u ON p.owner_id = u.id
	WHERE p.id = ?
`

func (s *Store) GetPost(ctx context.Context, id uint) (dtos.PostView, bool, error) {
	row := s.queryRow(ctx, queryPostByID, id)

	var p dtos.PostView
	if err := row.Scan(&p.ID, &p.Title, &p.Content, &p.OwnerID, &p.Author, &p.AuthorEmail, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dtos.PostView{}, false, nil
		}
		return dtos.PostView{}, false, err
	}
	return p, true, nil
}

const dmlInsertPost = `
	INSERT INTO posts (title, content, owner_id)
	VALUES (?, ?, ?)
`

func (s *Store) CreatePost(ctx context.Context, ownerID uint, title, content string) (dtos.PostView, error) {
	result, err := s.exec(ctx, dmlInsertPost, title, content, ownerID)
	if err != nil {
		return dtos.PostView{}, err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return dtos.PostView{}, err
	}

	post, _, err := s.GetPost(ctx, uint(insertedID))
	return post, err
}

const dmlUpdatePost = `
	UPDATE posts 
	SET title = ?, content = ?, updated_at = datetime('now')
	WHERE id = ?
`

func (s *Store) UpdatePost(ctx context.Context, id uint, title, content string) (dtos.PostView, error) {
	_, err := s.exec(ctx, dmlUpdatePost, title, content, id)
	if err != nil {
		return dtos.PostView{}, err
	}

	post, _, err := s.GetPost(ctx, id)
	return post, err
}

const dmlDeletePostByID = `
	DELETE FROM posts 
	WHERE id = ?
`

func (s *Store) DeletePost(ctx context.Context, id uint) error {
	_, err := s.exec(ctx, dmlDeletePostByID, id)
	return err
}
