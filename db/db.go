package db

import (
	"errors"
	"strings"

	"github.com/drkaka/kkpanic"
	"github.com/jackc/pgx"
)

// languages used to generate db column
var languages []string

// dbPool the pgx database pool.
var dbPool *pgx.ConnPool

// ErrLanguageIndex the language index is wrong.
var ErrLanguageIndex = errors.New("Language index wrong.")

// Use the pool to do further operations.
// langs must follow ISO-639-1 (https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes)
func Use(langs []string, pool *pgx.ConnPool) {
	for _, one := range langs {
		if len(one) != 2 {
			panic("Language length is not 2.")
		}
		languages = append(languages, strings.ToLower(one))
	}

	dbPool = pool

	prepareDB()
}

// prepareDB to prepare the database.
func prepareDB() {
	tx, err := dbPool.Begin()
	kkpanic.P(err)

	prepareCountry(tx)
	prepareCity(tx)

	kkpanic.P(tx.Commit())
}

// CheckColumnExisted to check whether the column is existed in table.
func CheckColumnExisted(table, column string) (bool, error) {
	var columnName pgx.NullString

	err := dbPool.QueryRow("SELECT column_name FROM information_schema.columns WHERE table_name=$1 AND column_name=$2", table, column).Scan(&columnName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	if !columnName.Valid {
		return false, nil
	} else if columnName.String == column {
		return true, nil
	}
	return false, nil
}
