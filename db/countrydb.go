package db

import (
	"errors"
	"strings"

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

	if _, err := tx.Exec(s); err != nil {
		panic(err)
	}
}

// AddCountry to add a country.
func AddCountry(id, name string) (bool, error) {
	if len(id) != 2 {
		return false, ErrCountryID
	}

	upperID := strings.ToUpper(id)
	_, err := dbPool.Exec("INSERT INTO country_info(id,name) VALUES($1,$2)", upperID, name)
	if err != nil {
		if err, ok := err.(pgx.PgError); ok && err.Code == "23505" {
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
		if err := rows.Scan(&country, &countryName); err != nil {
			panic(err)
		}

		countries = append(countries, country.String)
		countryNames = append(countryNames, countryName.String)
	}
	return countries, countryNames
}
