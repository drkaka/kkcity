package db

import (
	"fmt"

	"github.com/drkaka/kkpanic"
	"github.com/jackc/pgx"
)

func prepareCity(tx *pgx.Tx) {
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
	for i := range languages {
		nameColumn, addressColumn, err := getCityColumnNames(i)
		kkpanic.P(err)

		if existed, err := CheckColumnExisted("city_info", nameColumn); err != nil {
			panic(err)
		} else if !existed {
			_, err := tx.Exec(fmt.Sprintf("ALTER TABLE city_info ADD %s text;", nameColumn))
			kkpanic.P(err)

			_, err = tx.Exec(fmt.Sprintf("ALTER TABLE city_info ADD %s text;", addressColumn))
			kkpanic.P(err)
		}
	}
}

// getCityColumnNames to get the name of city name and address column.
func getCityColumnNames(tp int) (string, string, error) {
	llength := len(languages)
	if tp < 0 || tp >= llength {
		return "", "", ErrLanguageIndex
	}

	return fmt.Sprintf("name_%s", languages[tp]), fmt.Sprintf("address_%s", languages[tp]), nil
}

// AddCityInfo to add a city.
func AddCityInfo(placeid, country string) error {
	_, err := dbPool.Exec("INSERT INTO city_info(placeid,country_id) VALUES($1,$2)", placeid, country)
	if err != nil {
		if err, ok := err.(pgx.PgError); ok && err.Code == "23505" {
			return nil
		}
	}
	return err
}

// GetCityInfo to get city information of a certain language.
// Return place existed, name, address.
func GetCityInfo(placeid string, tp int) (bool, string, string, error) {
	nameColumn, addressColumn, err := getCityColumnNames(tp)
	if err != nil {
		return false, "", "", err
	}

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

// UpdateCityInfo to update a certain language.
func UpdateCityInfo(placeid, name, address string, tp int) error {
	nameColumn, addressColumn, err := getCityColumnNames(tp)
	kkpanic.P(err)

	s := fmt.Sprintf("UPDATE city_info SET %s=$1,%s=$2 WHERE placeid=$3", nameColumn, addressColumn)
	_, err = dbPool.Exec(s, name, address, placeid)
	return err
}
