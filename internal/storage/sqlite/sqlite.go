package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"lil-url/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.sqllite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", fn, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY,
		lil TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);

		CREATE INDEX IF NOT EXISTS idx_lil ON url(lil);
		`)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", fn, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s : %w", fn, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUrl(urlToSave, lil string) (int64, error) {
	const fn = "storage.sqlite.SaveUrl"

	stmt, err := s.db.Prepare("INSERT INTO urls(url, lil) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s : %w", fn, err)
	}

	res, err := stmt.Exec(urlToSave, lil)
	if err != nil {
		sqliteErr, ok := err.(sqlite3.Error)
		if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, storage.ErrUrlExists
		}

		return 0, fmt.Errorf("%s : %w", fn, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s : failed to get last insert id : %w", fn, err)
	}

	return id, nil
}

func (s *Storage) GetUrl(lil string) (string, error) {
	const fn = "storage.sqlite.GetUrl"

	stmt, err := s.db.Prepare("SELECT url FROM urls WHERE lil = ?")
	if err != nil {
		return "", fmt.Errorf("%s : %w", fn, err)
	}

	var resUrl string

	err = stmt.QueryRow(lil).Scan(&resUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrUrlNotFound
		}
		return "", fmt.Errorf("%s : %w", fn, err)
	}

	return resUrl, nil
}

func (s *Storage) DeleteUrl(lil string) error {
	const fn = "storage.sqlite.DeleteUrl"

	stmt, err := s.db.Prepare("DELETE FROM urls WHERE lil = ?")
	if err != nil {
		return fmt.Errorf("%s : %w", fn, err)
	}

	_, err = stmt.Exec(lil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrLilNotFound
		}
		return fmt.Errorf("%s : %w", fn, err)
	}

	return nil
}
