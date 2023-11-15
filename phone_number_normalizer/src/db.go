package normalize

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

func InsertMany(db *sql.DB, tableName string, columnName string, vals []interface{}) error {
	if !IsSafeTableName(tableName) {
		return fmt.Errorf("got unsafe table name: %q", tableName)
	}
	if !IsSafeTableName(columnName) {
		return fmt.Errorf("got unsafe column name: %q", columnName)
	}
	queryBase := fmt.Sprintf(`INSERT INTO "%s" ("%s") VALUES `, tableName, columnName)
	queryVals := make([]string, len(vals))
	for i := range vals {
		queryVals[i] = fmt.Sprintf("($%d)", i+1)
	}
	query := queryBase + strings.Join(queryVals, ",")
	_, err := db.Exec(query, vals...)
	return err
}

func ConnectToPostgres() (*sql.DB, error) {
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

func IsSafeTableName(tableName string) bool {
	return !strings.ContainsAny(tableName, " ;")
}
