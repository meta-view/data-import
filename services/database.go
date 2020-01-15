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
	err := db.checkTable(table)
	if err != nil {
		return "", err
	}

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

	sqlStmt := fmt.Sprintf("INSERT INTO %s(id, provider, imported, created, updated, content) VALUES (?, ?, ?, ?, ?, json(?)) ON CONFLICT(id) DO UPDATE SET provider=?, imported=?, created=?, updated=?, content=json(?) WHERE id = ?;", data["table"])

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
		_, err = db.DB.Exec(sqlStmt)
		if err != nil {
			return err
		}
		log.Printf("created table %s\n", table)

		sqlStmt = fmt.Sprintf("CREATE UNIQUE INDEX idx_%s_id ON %s(id);", table, table)
		_, err = db.DB.Exec(sqlStmt)
		if err != nil {
			return err
		}
		log.Printf("created index for table %s\n", table)

		//sqlStmt = fmt.Sprintf("CREATE INDEX idx_%s_json_name ON %s ( json_value(json_text, 'name') COLLATE NOCASE );", table, table)
		//_, err = db.DB.Exec(sqlStmt)
		//if err != nil {
		//	return err
		//}
		//log.Printf("created JSON index for table %s\n", table)

	}
	return nil
}

// ReadEntries - queries the database for all entries
func (db *Database) ReadEntries(query map[string]interface{}) (map[string]interface{}, error) {
	output := make(map[string]interface{})
	if query["table"] == nil {
		return output, errors.New("'table' need to be set for now")
	}
	if query["limit"] == nil {
		query["limit"] = "100"
	}
	if query["content"] == nil {
		query["content"] = make(map[string]interface{})
	}
	return db.queryTable(fmt.Sprintf("%v", query["table"]), query)
}

func (db *Database) queryTable(table string, query map[string]interface{}) (map[string]interface{}, error) {
	i := 0
	l := len(query) - 3
	hasWhere := false
	contentSelect := ""
	contentQuery := query["content"].(map[string]interface{})
	cl := len(contentQuery)
	output := make(map[string]interface{})
	queryStmt := fmt.Sprintf("SELECT id, provider, imported, created, updated, content %s FROM %s ", contentSelect, table)
	if len(query) > 3 {
		queryStmt += " WHERE "
		hasWhere = true
	}
	for k, v := range query {
		if k != "table" && k != "limit" && k != "content" {
			i++
			queryStmt += fmt.Sprintf(" %s='%s'", k, v)
			if i < l {
				queryStmt += " AND "
			}
		}
	}
	i = 0
	if cl > 0 {
		if !hasWhere {
			queryStmt += " WHERE "
		}
		for k, v := range contentQuery {
			i++
			queryStmt += fmt.Sprintf(" json_extract(content, '$.%s')='%s'", k, v)
			if i < cl {
				queryStmt += " AND "
			}
		}
	}
	queryStmt += fmt.Sprintf(" LIMIT %s;", query["limit"])
	log.Printf("query: %s\n", queryStmt)
	rows, err := db.DB.Query(queryStmt)
	if err != nil {
		return output, err
	}

	log.Printf("mapping results of %s to elements", query)
	for rows.Next() {
		data := make(map[string]interface{})
		var id, provider, imported, created, updated, content string
		err = rows.Scan(&id, &provider, &imported, &created, &updated, &content)
		if err == nil {
			data["id"] = id
			data["table"] = table
			data["imported"] = imported
			data["created"] = created
			data["updated"] = updated
			data["provider"] = provider
			data["content"] = content
			log.Printf("Reading entry %s for %s\n", data["id"], data["provider"])
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
