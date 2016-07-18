package db

import (
	"errors"
	"strings"

	"github.com/drkaka/kkpanic"
	"github.com/jackc/pgx"
)

// ErrCountryID to define the wrong country ID.
var ErrCountryID = errors.New("Country ID must be 2 charactor.")

func prepareCountry(tx *pgx.Tx) {
	// create country info table
	// id is uppercased like EN, AR.
	// name is English version.
	s := `CREATE TABLE IF NOT EXISTS country_info (
	id text primary key,
    name text);`

	_, err := tx.Exec(s)
	kkpanic.P(err)
}

// checkCountryID to check whether country id is valid.
func checkCountryID(id string) error {
	if len(id) != 2 {
		return ErrCountryID
	}
	return nil
}

// AddCountry to add a country.
func AddCountry(id, name string) error {
	if err := checkCountryID(id); err != nil {
		return err
	}

	upperID := strings.ToUpper(id)
	_, err := dbPool.Exec("INSERT INTO country_info(id,name) VALUES($1,$2)", upperID, name)
	if err != nil {
		if err, ok := err.(pgx.PgError); ok && err.Code == "23505" {
			return nil
		}
	}
	return err
}

// CheckCountryExisted to check whether a country is already added.
func CheckCountryExisted(id string) (bool, error) {
	if err := checkCountryID(id); err != nil {
		return false, err
	}

	var countryID pgx.NullString

	err := dbPool.QueryRow("SELECT id FROM country_info WHERE id=$1", id).Scan(&countryID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetCountries to get country and their names.
func GetCountries() ([]string, []string) {
	rows, _ := dbPool.Query("SELECT id,name FROM country_info")

	var countries []string
	var countryNames []string
	for rows.Next() {
		var country pgx.NullString
		var countryName pgx.NullString

		kkpanic.P(rows.Scan(&country, &countryName))

		countries = append(countries, country.String)
		countryNames = append(countryNames, countryName.String)
	}
	return countries, countryNames
}
