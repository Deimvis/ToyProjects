package main

import (
	"bufio"
	"database/sql"
	"flag"
	"io"
	"os"

	_ "github.com/lib/pq"

	normalize "github.com/Deimvis/toyprojects/phone_number_normalizer/src"
)

func main() {
	filePath := flag.String("file", "numbers.txt", "Path to file with phone numbers (one by line)")
	flag.Parse()
	file, err := os.Open(*filePath)
	check(err)
	defer file.Close()
	phoneNumbers := readPhoneNumbers(file)
	db, err := normalize.ConnectToPostgres()
	check(err)
	defer db.Close()
	err = initDB(db, phoneNumbers)
	check(err)
	err = normalizePhoneNumbers(db)
	check(err)
}

func normalizePhoneNumbers(db *sql.DB) error {
	rows, err := db.Query(`SELECT "id", "value" FROM "phone_number"`)
	if err != nil {
		return err
	}
	defer rows.Close()
	var normalizedVals []interface{}
	for rows.Next() {
		var id int
		var value string
		err = rows.Scan(&id, &value)
		if err != nil {
			return err
		}
		normalizedVals = append(normalizedVals, normalize.NormalizePhoneNumber(value))
	}
	err = normalize.InsertMany(db, "normalized_phone_number", "value", normalizedVals)
	return err
}

func initDB(db *sql.DB, phoneNumbers []interface{}) error {
	query := `
		CREATE TABLE IF NOT EXISTS "phone_number" (
			id SERIAL,
			value VARCHAR(255)
		)`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	query = `
		CREATE TABLE IF NOT EXISTS "normalized_phone_number" (
			id SERIAL,
			value VARCHAR(255)
		)`
	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	_, err = db.Exec(`TRUNCATE ONLY "phone_number", "normalized_phone_number" RESTART IDENTITY`)
	if err != nil {
		return err
	}
	err = normalize.InsertMany(db, "phone_number", "value", phoneNumbers)
	return err
}

func readPhoneNumbers(r io.Reader) []interface{} {
	var phoneNumbers []interface{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		phoneNumbers = append(phoneNumbers, scanner.Text())
	}
	return phoneNumbers
}

func phoneNumbers() []interface{} {
	return []interface{}{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
