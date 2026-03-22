package sqlite

import (
	"context"
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const (
	databaseTypeSQLLite    = "sqlite"
	defaultMaxOpenConnSize = 3
)

type Store struct {
	db *sql.DB
}

func New(databasePath, schemaFile, seedFile string) (*Store, error) {
	db, err := sql.Open(databaseTypeSQLLite, databasePath)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(defaultMaxOpenConnSize)

	store := &Store{db: db}
	if err := store.initialize(schemaFile, seedFile); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) queryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

func (s *Store) query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

func (s *Store) exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

// Note: 테스트 프로젝트로 initial query 처리를 위해 작성한 함수
func (s *Store) initialize(schemaFile, seedFile string) error {
	if err := s.execSQLFile(schemaFile); err != nil {
		return err
	}
	return s.execSQLFile(seedFile)
}

func (s *Store) execSQLFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(string(content))
	return err
}
