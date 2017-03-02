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

func (s *Storage) Get(rowID int) (string, error) {
	sqlGet := `
	SELECT content FROM clips
	WHERE rowid = ?
	`

	stmt, _ := s.db.Prepare(sqlGet)
	defer stmt.Close()
	var content string
	rows, _ := stmt.Query(rowID)
	for rows.Next() {
		_ = rows.Scan(&content)
		log.Printf("got it: %v\n", content)
	}

	return content, nil
}

func (s *Storage) List() {
	sqlReadall := `
	SELECT rowid, id, content FROM clips
	ORDER BY datetime(created_at) DESC
	`

	rows, _ := s.db.Query(sqlReadall)
	defer rows.Close()

	for rows.Next() {
		var clipShortcut int
		item := NewClip()
		_ = rows.Scan(&clipShortcut, &item.Hash, &item.Content)
		fmt.Printf("[ %d ]\t%s\n", clipShortcut, item.Content)
	}
}

func (s *Storage) SaveIfNew(c *Clip) error {
	sql_add := `
	INSERT OR REPLACE INTO clips(
		id,
		content,
		created_at
	) VALUES(?, ?, CURRENT_TIMESTAMP)
	`

	stmt, _ := s.db.Prepare(sql_add)
	defer stmt.Close()
	_, err := stmt.Exec(c.Hash, c.Content)
	return err
}
