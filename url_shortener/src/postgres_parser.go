package urlshortener

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type URLMappingDB struct {
	Path string `db:"path"`
	URL  string `db:"url"`
}

func parsePostgres(tableName string) (map[string]string, error) {
	db, err := connectToPostgres()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf("SELECT path, url FROM \"%s\"", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	pathToUrls := make(map[string]string)
	for rows.Next() {
		mapping := URLMappingDB{}
		err = rows.Scan(&mapping.Path, &mapping.URL)
		if err != nil {
			return nil, err
		}
		pathToUrls[mapping.Path] = mapping.URL
	}
	return pathToUrls, nil
}

func connectToPostgres() (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func isSafeTableName(tableName string) bool {
	return !strings.ContainsAny(tableName, " ;")
}
