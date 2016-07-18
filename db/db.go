package db

import (
	"errors"
	"strings"

	"github.com/jackc/pgx"
)

// dbPool the pgx database pool.
var dbPool *pgx.ConnPool

// Use the pool to do further operations.
func Use(langs []string, pool *pgx.ConnPool) error {
	for _, one := range langs {
		if len(one) != 2 {
			return errors.New("Language length is not 2.")
		}
		languages = append(languages, strings.ToLower(one))
	}

	dbPool = pool
	return prepareDB()
}

// prepareDB to prepare the database.
func prepareDB() error {
	tx, err := dbPool.Begin()
	if err != nil {
		panic(err)
	}

	prepareCountry(tx)
	prepareCity(tx)

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
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
