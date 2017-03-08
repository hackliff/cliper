// storage.go
// Copyright (C) 2017 hackliff <xavier.bruhiere@gmail.com>
// Distributed under terms of the MIT license.

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const DB_DRIVER string = "sqlite3"
const CLIPS_TABLE string = "clips"

type Storage struct {
	db *sql.DB
}

func NewStorage(dbPath string, reset bool) (*Storage, error) {
	if reset {
		log.Println("reseting database")
		os.Remove(dbPath)
	}

	db, err := sql.Open(DB_DRIVER, dbPath)

	return &Storage{db}, err
}

func (s *Storage) Init() error {
	sql_table := `
	CREATE TABLE IF NOT EXISTS clips(
		id TEXT NOT NULL PRIMARY KEY,
		content TEXT NOT NULL,
		created_at DATETIME
	);
	`
	_, err := s.db.Exec(sql_table)
	return err
}

func (s *Storage) Get(rowID int) (*Clip, error) {
	sqlGet := `
	SELECT content FROM clips
	WHERE rowid = ?
	`

	stmt, _ := s.db.Prepare(sqlGet)
	defer stmt.Close()

	clip := &Clip{}
	rows, _ := stmt.Query(rowID)
	for rows.Next() {
		_ = rows.Scan(&clip.Content)
	}

	return clip, nil
}

func (s *Storage) List(limit int) (clips []*Clip) {
	sqlReadall := fmt.Sprintf(`
	SELECT rowid, id, content FROM %s
	ORDER BY datetime(created_at) DESC
	LIMIT %d
	`, CLIPS_TABLE, limit)

	log.Printf("scaning for clips (limit=%d)\n", limit)
	rows, _ := s.db.Query(sqlReadall)
	defer rows.Close()

	for rows.Next() {
		item := &Clip{}
		_ = rows.Scan(&item.Shortcut, &item.Hash, &item.Content)
		clips = append(clips, item)
	}

	return clips
}

func (s *Storage) SaveIfNew(c *Clip) error {
	sqlAdd := `
	INSERT OR REPLACE INTO clips(
		id,
		content,
		created_at
	) VALUES(?, ?, CURRENT_TIMESTAMP)
	`

	stmt, _ := s.db.Prepare(sqlAdd)
	defer stmt.Close()

	_, err := stmt.Exec(c.Hash, c.Content)
	return err
}
