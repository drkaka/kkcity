package kkcity

import (
	"sync"

	"github.com/jackc/pgx"
)

// Use the pool to do further operations.
// langs must follow ISO-639-1 (https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes)
func Use(langs []string, gKey string, pool *pgx.ConnPool) {
	setupLanguage(langs)
	googleKey = gKey

	use(pool, getAll())
}

// handleCityInfo to deal with city information with placeid.
// Return city name, address, error
func handleCityInfo(placeid, lang string) (string, string, error) {
	var err error
	var cityExist bool
	var cityName, cityAddress string

	cityExist, cityName, cityAddress, err = getCityInfo(placeid, lang)
	if err != nil {
		return "", "", err
	}

	cityLangExist := len(cityName) > 0
	if !cityExist || !cityLangExist {
		var country, countryName string
		country, countryName, cityName, cityAddress, err = requestPlaceInfo(placeid, lang)
		if err != nil {
			return "", "", err
		}

		var countryExist bool
		var recordedCountryName string
		countryExist, recordedCountryName, err = getCountryName(country, lang)
		if err != nil {
			return "", "", err
		}

		if !countryExist {
			if err = addCountry(country, countryName, lang); err != nil {
				return "", "", err
			}
		} else if len(recordedCountryName) == 0 {
			if err = updateCountryInfo(country, countryName, lang); err != nil {
				return "", "", err
			}
		}

		if !cityExist {
			if err = addCityInfo(placeid, country, cityName, cityAddress, lang); err != nil {
				return "", "", err
			}
		} else {
			if err = updateCityInfo(placeid, cityName, cityAddress, lang); err != nil {
				return "", "", err
			}
		}
	}
	return cityName, cityAddress, nil
}

// GetCountries to get all the countries.
// Return country ids, names, error
func GetCountries(langIndex int) ([]string, []string, error) {
	lang, err := getLanguage(langIndex)
	if err != nil {
		return nil, nil, err
	}
	return getCountries(lang)
}

// GetCityWithLatLng to get city information with lat and lng.
// Return placeid, name, address, error
func GetCityWithLatLng(lat, lng float32, langIndex int) (string, string, string, error) {
	lang, err := getLanguage(langIndex)
	if err != nil {
		return "", "", "", err
	}

	var placeid string
	placeid, err = requestLocationWithLatLng(lat, lng)
	if err != nil {
		return "", "", "", err
	}

	var name, address string
	name, address, err = handleCityInfo(placeid, lang)
	return placeid, name, address, err
}

// GetCitiesWithInput to get cities with input.
// Return place ids, city names, addresses, error
func GetCitiesWithInput(input string, langIndex int) ([]string, []string, []string, error) {
	var placeIDs, cityNames, cityAddresses []string
	lang, err := getLanguage(langIndex)
	if err != nil {
		return placeIDs, cityNames, cityAddresses, err
	}

	placeIDs, cityAddresses, err = requestAutoComplete(input, lang)
	if err != nil {
		return placeIDs, cityNames, cityAddresses, err
	}

	cityNames = make([]string, len(placeIDs))

	var wg sync.WaitGroup
	for i, id := range placeIDs {
		wg.Add(1)

		go func(thisID string, index int) {
			defer wg.Done()

			var cityName string
			cityName, _, err = handleCityInfo(thisID, lang)
			cityNames[index] = cityName
		}(id, i)
	}
	wg.Wait()

	return placeIDs, cityNames, cityAddresses, err
}

// GetCountryCities to get all the cities in one country.
// Return city ids, names, addresses, error
func GetCountryCities(countryID string, langIndex int) ([]string, []string, []string, error) {
	lang, err := getLanguage(langIndex)
	if err != nil {
		return nil, nil, nil, err
	}
	return getCountryCities(countryID, lang)
}
