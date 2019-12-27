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

// SaveEntry insert an entry into the database
func (db *Database) SaveEntry(data map[string]interface{}) (string, error) {
	table := fmt.Sprintf("%v", data["table"])
	if table == "" {
		return "", errors.New("table to set")
	}
	db.checkTable(table)
	if data["id"] == nil {
		ID := uuid.New()
		data["id"] = ID.String()
	}

	if data["provider"] == nil {
		data["provider"] = "N/A"
	}

	date := time.Now().Format(time.RFC3339)
	if data["created"] == nil {
		data["created"] = date
	}
	if data["imported"] == nil {
		data["imported"] = date
	}
	data["updated"] = date

	return db.insertOrUpdateEntry(data)
}

func (db *Database) insertOrUpdateEntry(data map[string]interface{}) (string, error) {
	id := fmt.Sprintf("%v", data["id"])
	log.Printf("inserting id %s in table %s\n", data["id"], data["table"])
	tx, err := db.DB.Begin()
	if err != nil {
		return "", err
	}

	sqlStmt := fmt.Sprintf("INSERT INTO %s(id, provider, imported, created, updated, content) VALUES(?, ?, ?, ?, ?, ?) ON CONFLICT(id) DO UPDATE SET provider=?, imported=?, created=?, updated=?, content=? WHERE id = ?;", data["table"])

	stmt, err := tx.Prepare(sqlStmt)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	content, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	_, err = stmt.Exec(
		id,
		data["provider"],
		data["imported"],
		data["created"],
		data["updated"],
		content,
		data["provider"],
		data["imported"],
		data["created"],
		data["updated"],
		content,
		id)
	if err != nil {
		return "", err
	}
	tx.Commit()
	return id, nil

}

func (db *Database) checkTable(table string) error {
	stmt, err := db.DB.Prepare("SELECT name FROM sqlite_master WHERE type='table' AND name=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	var name string
	var sqlStmt string
	err = stmt.QueryRow(table).Scan(&name)
	if err != nil && err == sql.ErrNoRows {
		sqlStmt = fmt.Sprintf("CREATE TABLE %s (id TEXT not null primary key, provider TEXT, imported TEXT, created TEXT, updated TEXT, content TEXT);", table)
		log.Printf("created table %s\n", table)
		_, err = db.DB.Exec(sqlStmt)
		if err != nil {
			return err
		}

		sqlStmt = fmt.Sprintf("CREATE UNIQUE INDEX idx_%s_id ON %s(id);", table, table)
		log.Printf("created index for table %s\n", table)
		_, err = db.DB.Exec(sqlStmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadEntries - queries the database for all entries
func (db *Database) ReadEntries(query map[string]interface{}) (map[string]interface{}, error) {
	output := make(map[string]interface{})
	if query["table"] == nil {
		return output, errors.New("'table' need to be set for now")
	}
	return db.queryTable(fmt.Sprintf("%v", query["table"]), query)
}

func (db *Database) queryTable(table string, query map[string]interface{}) (map[string]interface{}, error) {
	i := 0
	l := len(query) - 1
	output := make(map[string]interface{})
	queryStmt := fmt.Sprintf("SELECT * FROM %s WHERE ", table)
	for k, v := range query {
		if k != "table" {
			i++
			queryStmt += fmt.Sprintf(" %s='%s'", k, v)
			if i < l {
				queryStmt += " AND "
			}
		}
	}
	log.Printf("query: %s\n", queryStmt)
	rows, err := db.DB.Query(queryStmt)
	if err != nil {
		return output, err
	}

	data := make(map[string]interface{})
	var id, provider, imported, created, updated, content string
	for rows.Next() {
		err = rows.Scan(&id, &provider, &imported, &created, &updated, &content)
		if err == nil {
			data["id"] = id
			data["imported"] = imported
			data["created"] = created
			data["updated"] = updated
			data["content"] = content
			output[id] = data
		} else {
			log.Printf("Error %s while loading row.\n", err.Error())
		}
	}
	rows.Close()
	return output, nil
}

// Close - closes the main DB
func (db *Database) Close() {
	db.DB.Close()
}
