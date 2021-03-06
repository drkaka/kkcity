package kkcity

import "github.com/stretchr/testify/suite"

type dbHandleSuite struct {
	suite.Suite
}

func (suite *dbHandleSuite) TearDownSuite() {
	var err error
	_, err = dbPool.Exec("DROP TABLE country_info;")
	suite.NoError(err, "country_info should be able to be dropped.")

	_, err = dbPool.Exec("DROP TABLE city_info;")
	suite.NoError(err, "city_info should be able to be dropped.")

	dbPool.Close()
}

func (suite *dbHandleSuite) TestColumnsExist() {
	for _, one := range testLangs {
		var existed bool
		var err error
		var columnName, columnAddress string

		columnName, columnAddress = getCityColumnNames(one)

		existed, err = checkDBColumnExisted("city_info", columnName)
		suite.True(existed, columnName, " should be existed.")
		suite.Nil(err, "There should be no error while check exist.")

		existed, err = checkDBColumnExisted("city_info", columnAddress)
		suite.True(existed, columnAddress, " should be existed.")
		suite.Nil(err, "There should be no error while check exist.")

		columnName = getCountryColumnName(one)
		existed, err = checkDBColumnExisted("country_info", columnName)
		suite.True(existed, columnName, " should be existed.")
		suite.Nil(err, "There should be no error while check exist.")
	}
}

func (suite *dbHandleSuite) TestCityInfo() {
	pid1 := "placeid1"
	countryID := "CN"
	cityName := "Xiamen"
	cityAddress := "Xiamen, Fujian, China"

	lang, err := getLanguage(0)
	suite.NoError(err, "Shoule be able to get language.")

	// add city information
	err = addCityInfo(pid1, countryID, cityName, cityAddress, lang)
	suite.NoError(err, "Should be able to add city info.")

	// add duplicated city information
	err = addCityInfo(pid1, countryID, cityName, cityAddress, lang)
	suite.Equal(ErrCityExisted, err, "City already existed.")

	// update city information
	err = updateCityInfo(pid1, cityName, cityAddress, lang)
	suite.NoError(err, "Should be able to update city info.")

	// get city information
	existed, resultName, resultAddress, err := getCityInfo(pid1, lang)
	suite.True(existed, "The result should be existed.")
	suite.NoError(err, "Should be able to get.")
	suite.EqualValues(cityName, resultName, "The name should be equal")
	suite.EqualValues(cityAddress, resultAddress, "The address should be equal")

	var noLang string
	noLang, err = getLanguage(1)
	suite.NoError(err, "Shoule be able to get language.")

	// get not set language information
	existed, resultName, resultAddress, err = getCityInfo(pid1, noLang)
	suite.True(existed, "The result should be existed.")
	suite.NoError(err, "Should be able to get.")
	suite.EqualValues("", resultName, "The name should be empty")
	suite.EqualValues("", resultAddress, "The address should be empty")

	// check not existed city
	noPlace := "placeid2"
	existed, _, _, err = getCityInfo(noPlace, lang)
	suite.False(existed, "The place should be not existed.")
	suite.NoError(err, "Should be able to get.")

	// add another city information
	err = addCityInfo(noPlace, countryID, "", "", lang)
	suite.NoError(err, "Should be able to add city info.")

	// get all the cities in one country.
	var pids, names, addresses []string
	pids, names, addresses, err = getCountryCities(countryID, lang)
	suite.NoError(err, "Shoule be able to get cities.")
	suite.EqualValues(2, len(pids), "Should have 2 result.")
	suite.EqualValues(2, len(names), "Should have 2 result.")
	suite.EqualValues(2, len(addresses), "Should have 2 result.")

	for i, one := range pids {
		if one != pid1 && one != noPlace {
			suite.Fail("id is wrong.")
		}

		if one == noPlace {
			suite.EqualValues("", names[i], "Name should be empty.")
			suite.EqualValues("", addresses[i], "Address should be empty.")
		}

		if one == pid1 {
			suite.EqualValues(cityName, names[i], "Name is wrong.")
			suite.EqualValues(cityAddress, addresses[i], "Address is wrong.")
		}
	}
}

func (suite *dbHandleSuite) TestCountryInfo() {
	id1 := "EN"
	id2 := "CN"
	name1 := "English"
	name2 := "Chinese"
	name2CN := "中国"
	id2Lower := "cn"

	badID := "123"

	var lang0, lang1 string
	var err error

	lang0, err = getLanguage(0)
	suite.NoError(err, "Shoule be able to get language.")

	lang1, err = getLanguage(1)
	suite.NoError(err, "Shoule be able to get language.")

	err = addCountry(badID, "bad", lang0)
	suite.Equal(ErrCountryID, err, "country id format is wrong.")

	err = addCountry(id1, name1, lang0)
	suite.NoError(err, "Should have no error.")

	err = addCountry(id1, name1, lang1)
	suite.Equal(ErrCountryExisted, err, "Should have error that country is existed.")

	err = addCountry(id2Lower, name2, lang0)
	suite.NoError(err, "Shoule have no error.")

	err = updateCountryInfo(id2, name2CN, lang1)
	suite.NoError(err, "Shoule have no error.")

	var ids, names []string
	ids, names, err = getCountries(lang1)
	suite.NoError(err, "Shoule have no error.")

	suite.EqualValues(2, len(ids), "Shoule have 2 result.")
	suite.EqualValues(2, len(names), "Shoule have 2 result.")
	for i, one := range ids {
		if one != id1 && one != id2 {
			suite.Fail("id is wrong.")
		}

		if one == id1 {
			suite.Equal("", names[i], "Name should be empty.")
		}

		if one == id2 {
			suite.Equal(name2CN, names[i], "Name is wrong.")
		}
	}

	var existed bool
	var name string

	existed, name, err = getCountryName(badID, lang0)
	suite.False(existed, "Country should not existed.")
	suite.Equal("", name, "Name should be empty.")
	suite.EqualValues(ErrCountryID, err, "Should have bad country id error.")

	existed, name, err = getCountryName(id2, lang1)
	suite.True(existed, "Country should existed.")
	suite.Equal(name2CN, name, "Name is wrong.")
	suite.NoError(err, "Should be able to get country.")
}
