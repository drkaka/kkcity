package db

import (
	"errors"
	"fmt"
	"strings"

	"github.com/drkaka/kkpanic"
	"github.com/jackc/pgx"
)

// ErrCountryID to define the wrong country ID.
var ErrCountryID = errors.New("Country ID must be 2 charactor.")

// ErrCountryExisted to define the country already existed.
var ErrCountryExisted = errors.New("Country is already existed.")

func prepareCountry(tx *pgx.Tx) {
	// create country info table
	// id is uppercased like EN, AR.
	// name is English version.
	s := `CREATE TABLE IF NOT EXISTS country_info (
	id text primary key);`

	_, err := tx.Exec(s)
	kkpanic.P(err)

	// setup the language name and address column
	for i := range languages {
		nameColumn, err := getCountryColumnName(i)
		kkpanic.P(err)

		if existed, err := CheckColumnExisted("country_info", nameColumn); err != nil {
			panic(err)
		} else if !existed {
			_, err := tx.Exec(fmt.Sprintf("ALTER TABLE country_info ADD %s text;", nameColumn))
			kkpanic.P(err)
		}
	}
}

// getCountryColumnName to get the name of country name column.
func getCountryColumnName(tp int) (string, error) {
	llength := len(languages)
	if tp < 0 || tp >= llength {
		return "", ErrLanguageIndex
	}

	return fmt.Sprintf("name_%s", languages[tp]), nil
}

// checkCountryID to check whether country id is valid.
func checkCountryID(id string) error {
	if len(id) != 2 {
		return ErrCountryID
	}
	return nil
}

// AddCountry to add a country.
func AddCountry(id, name string, tp int) error {
	if err := checkCountryID(id); err != nil {
		return err
	}

	nameColumn, err := getCountryColumnName(tp)
	if err != nil {
		return err
	}

	s := fmt.Sprintf("INSERT INTO country_info(id,%s) VALUES($1,$2)", nameColumn)

	upperID := strings.ToUpper(id)
	_, err = dbPool.Exec(s, upperID, name)
	if err != nil {
		if err, ok := err.(pgx.PgError); ok && err.Code == "23505" {
			return ErrCountryExisted
		}
	}
	return err
}

// GetCountryName to get certain country name.
func GetCountryName(id string, tp int) (bool, string, error) {
	if err := checkCountryID(id); err != nil {
		return false, "", err
	}

	nameColumn, err := getCountryColumnName(tp)
	if err != nil {
		return false, "", err
	}

	var countryID, countryName pgx.NullString
	s := fmt.Sprintf("SELECT id,%s FROM country_info WHERE id=$1", nameColumn)

	upperID := strings.ToUpper(id)
	err = dbPool.QueryRow(s, upperID).Scan(&countryID, &countryName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, "", nil
		}
		return false, "", err
	}
	return true, countryName.String, nil
}

// UpdateCountryInfo to update a certain language.
func UpdateCountryInfo(id, name string, tp int) error {
	nameColumn, err := getCountryColumnName(tp)
	if err != nil {
		return err
	}

	s := fmt.Sprintf("UPDATE country_info SET %s=$1 WHERE id=$2", nameColumn)
	_, err = dbPool.Exec(s, name, id)
	return err
}

// GetCountries to get country and their names.
func GetCountries(tp int) ([]string, []string, error) {
	nameColumn, err := getCountryColumnName(tp)
	if err != nil {
		return nil, nil, err
	}

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
