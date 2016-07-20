package kkcity

import (
	"errors"
	"fmt"
	"strings"

	"github.com/drkaka/kkpanic"
	"github.com/jackc/pgx"
)

var (
	// ErrCountryID to define the wrong country ID.
	ErrCountryID = errors.New("Country ID must be 2 charactor.")

	// ErrCountryExisted to define the country already existed.
	ErrCountryExisted = errors.New("Country is already existed.")
)

// dbPool the pgx database pool.
var dbPool *pgx.ConnPool

// use the pool to do further operations.
func use(pool *pgx.ConnPool, langs []string) {
	dbPool = pool

	tx, err := dbPool.Begin()
	kkpanic.P(err)

	prepareCountry(tx, langs)
	prepareCity(tx, langs)

	kkpanic.P(tx.Commit())
}

// checkDBColumnExisted to check whether the column is existed in table.
func checkDBColumnExisted(table, column string) (bool, error) {
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

func prepareCity(tx *pgx.Tx, langs []string) {
	var err error
	// create city info table
	s := `CREATE TABLE IF NOT EXISTS city_info (
	placeid text primary key,
	country_id text);`

	_, err = tx.Exec(s)
	kkpanic.P(err)

	_, err = tx.Exec("CREATE INDEX IF NOT EXISTS index_city_info_country_id ON city_info (country_id);")
	kkpanic.P(err)

	// setup the language name and address column
	for _, one := range langs {
		nameColumn, addressColumn := getCityColumnNames(one)

		if existed, err := checkDBColumnExisted("city_info", nameColumn); err != nil {
			panic(err)
		} else if !existed {
			_, err := tx.Exec(fmt.Sprintf("ALTER TABLE city_info ADD %s text;", nameColumn))
			kkpanic.P(err)

			_, err = tx.Exec(fmt.Sprintf("ALTER TABLE city_info ADD %s text;", addressColumn))
			kkpanic.P(err)
		}
	}
}

func prepareCountry(tx *pgx.Tx, langs []string) {
	// create country info table
	// id is uppercased like EN, AR.
	// name is English version.
	s := `CREATE TABLE IF NOT EXISTS country_info (
	id text primary key);`

	_, err := tx.Exec(s)
	kkpanic.P(err)

	// setup the language name and address column
	for _, one := range langs {
		nameColumn := getCountryColumnName(one)

		if existed, err := checkDBColumnExisted("country_info", nameColumn); err != nil {
			panic(err)
		} else if !existed {
			_, err := tx.Exec(fmt.Sprintf("ALTER TABLE country_info ADD %s text;", nameColumn))
			kkpanic.P(err)
		}
	}
}

// getCityColumnNames to get the name of city name and address column.
func getCityColumnNames(lang string) (string, string) {
	return fmt.Sprintf("name_%s", lang), fmt.Sprintf("address_%s", lang)
}

// addCityInfo to add a city.
func addCityInfo(placeid, country string) error {
	_, err := dbPool.Exec("INSERT INTO city_info(placeid,country_id) VALUES($1,$2)", placeid, country)
	if err != nil {
		if err, ok := err.(pgx.PgError); ok && err.Code == "23505" {
			return nil
		}
	}
	return err
}

// getCityInfo to get city information of a certain language.
// Return place existed, name, address, error.
func getCityInfo(placeid, lang string) (bool, string, string, error) {
	nameColumn, addressColumn := getCityColumnNames(lang)

	s := fmt.Sprintf("SELECT %s,%s FROM city_info WHERE placeid=$1", nameColumn, addressColumn)

	var name, address pgx.NullString
	if err := dbPool.QueryRow(s, placeid).Scan(&name, &address); err != nil {
		if err == pgx.ErrNoRows {
			return false, "", "", nil
		}
		return false, "", "", err
	}
	return true, name.String, address.String, nil
}

// updateCityInfo to update a certain language.
func updateCityInfo(placeid, name, address, lang string) error {
	nameColumn, addressColumn := getCityColumnNames(lang)

	s := fmt.Sprintf("UPDATE city_info SET %s=$1,%s=$2 WHERE placeid=$3", nameColumn, addressColumn)

	_, err := dbPool.Exec(s, name, address, placeid)
	return err
}

// getCountryColumnName to get the name of country name column.
func getCountryColumnName(lang string) string {
	return fmt.Sprintf("name_%s", lang)
}

// checkCountryID to check whether country id is valid.
func checkCountryID(id string) error {
	if len(id) != 2 {
		return ErrCountryID
	}
	return nil
}

// addCountry to add a country.
func addCountry(id, name, lang string) error {
	if err := checkCountryID(id); err != nil {
		return err
	}

	nameColumn := getCountryColumnName(lang)

	s := fmt.Sprintf("INSERT INTO country_info(id,%s) VALUES($1,$2)", nameColumn)

	upperID := strings.ToUpper(id)
	_, err := dbPool.Exec(s, upperID, name)
	if err != nil {
		if err, ok := err.(pgx.PgError); ok && err.Code == "23505" {
			return ErrCountryExisted
		}
	}
	return err
}

// getCountryName to get certain country name.
func getCountryName(id, lang string) (bool, string, error) {
	if err := checkCountryID(id); err != nil {
		return false, "", err
	}

	nameColumn := getCountryColumnName(lang)

	var countryID, countryName pgx.NullString
	s := fmt.Sprintf("SELECT id,%s FROM country_info WHERE id=$1", nameColumn)

	upperID := strings.ToUpper(id)
	err := dbPool.QueryRow(s, upperID).Scan(&countryID, &countryName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, "", nil
		}
		return false, "", err
	}
	return true, countryName.String, nil
}

// updateCountryInfo to update a certain language.
func updateCountryInfo(id, name, lang string) error {
	nameColumn := getCountryColumnName(lang)

	s := fmt.Sprintf("UPDATE country_info SET %s=$1 WHERE id=$2", nameColumn)

	_, err := dbPool.Exec(s, name, id)
	return err
}

// getCountries to get country and their names.
func getCountries(lang string) ([]string, []string, error) {
	nameColumn := getCountryColumnName(lang)

	s := fmt.Sprintf("SELECT id,%s FROM country_info", nameColumn)
	rows, _ := dbPool.Query(s)

	var countries []string
	var countryNames []string
	for rows.Next() {
		var country pgx.NullString
		var countryName pgx.NullString

		if err := rows.Scan(&country, &countryName); err != nil {
			return countries, countryNames, err
		}

		countries = append(countries, country.String)
		countryNames = append(countryNames, countryName.String)
	}
	return countries, countryNames, nil
}
