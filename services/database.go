package services

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"log"
	"os"
	"path"

	"github.com/google/uuid"
	// required for sqlite
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/square/go-jose.v2/json"
)

var dbDataDirectory = path.Join("data", "database")

// Database - basic struct for the database
type Database struct {
	Path string
	DB   *sql.DB
}

func init() {
	folders := []string{dbDataDirectory}
	for _, folder := range folders {
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			os.MkdirAll(folder, 0700)
			log.Printf("created directory %s", folder)
		}
	}

}

// CreateDatabase - create a new database struct or opens an existing one.
func CreateDatabase() (*Database, error) {
	db, err := sql.Open("sqlite3", path.Join(dbDataDirectory, "meta-view.db"))
	if err != nil {
		return nil, err
	}
	return &Database{
		Path: dbDataDirectory,
		DB:   db,
	}, nil
}

// InsertEntry insert an entry into the database
func (db *Database) InsertEntry(data map[string]interface{}) error {
	var id string
	table := fmt.Sprintf("%v", data["table"])
	if table == "" {
		return errors.New("table to set")
	}
	db.checkTable(table)
	if data["id"] == nil {
		ID := uuid.New()
		id = ID.String()
		data["id"] = id
	} else {
		id = fmt.Sprintf("%v", data["id"])
	}

	log.Printf("saving id %s in table %s\n", id, table)

	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(fmt.Sprintf("insert into %s(id, created, updated, content) values(?, ?, ?, ?)", table))
	if err != nil {
		return err
	}
	defer stmt.Close()
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}

	date := time.Now().Format(time.RFC3339)
	_, err = stmt.Exec(id, date, date, content)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (db *Database) checkTable(table string) error {
	stmt, err := db.DB.Prepare("SELECT name FROM sqlite_master WHERE type='table' AND name=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRow(table).Scan(&name)
	if err != nil && err == sql.ErrNoRows {
		sqlStmt := fmt.Sprintf("CREATE TABLE %s (id TEXT not null primary key, created TEXT, updated TEXT, content TEXT);", table)
		log.Printf("created table %s\n", table)
		_, err = db.DB.Exec(sqlStmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// Close - closes the main DB
func (db *Database) Close() {
	db.DB.Close()
}
