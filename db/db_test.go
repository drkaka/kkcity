package db

import (
	"net"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/suite"
)

var testLangs = []string{"en", "cn"}

func TestDB(t *testing.T) {
	suite.Run(t, new(DBHandleSuite))
}

type DBHandleSuite struct {
	suite.Suite
}

func (suite *DBHandleSuite) SetupTest() {
	DBName := os.Getenv("dbname")
	DBHost := os.Getenv("dbhost")
	DBUser := os.Getenv("dbuser")
	DBPassword := os.Getenv("dbpassword")

	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     DBHost,
			User:     DBUser,
			Password: DBPassword,
			Database: DBName,
			Dial:     (&net.Dialer{KeepAlive: 5 * time.Minute, Timeout: 5 * time.Second}).Dial,
		},
		MaxConnections: 10,
	}

	var err error
	var pool *pgx.ConnPool

	pool, err = pgx.NewConnPool(connPoolConfig)
	suite.NoError(err, "Should be able to create pool.")

	err = Use(testLangs, pool)
	suite.NoError(err, "Should be able to use DB.")
}

func (suite *DBHandleSuite) TearDownSuite() {
	var err error
	_, err = dbPool.Exec("DROP TABLE country_info;")
	suite.NoError(err, "country_info should be able to be dropped.")

	_, err = dbPool.Exec("DROP TABLE city_info;")
	suite.NoError(err, "city_info should be able to be dropped.")
}

func (suite *DBHandleSuite) TestColumnsExist() {
	for i := range testLangs {
		var existed bool
		var err error
		var columnName, columnAddress string

		columnName, columnAddress, err = getCityColumnNames(i)
		suite.NoError(err, "Should be able to get.")

		existed, err = CheckColumnExisted("city_info", columnName)
		suite.True(existed, columnName, " should be existed.")
		suite.Nil(err, "There should be no error while check exist.")

		existed, err = CheckColumnExisted("city_info", columnAddress)
		suite.True(existed, columnAddress, " should be existed.")
		suite.Nil(err, "There should be no error while check exist.")
	}
}

func (suite *DBHandleSuite) TestCityInfo() {
	_, _, err := getCityColumnNames(2)
	suite.Error(err, "Should have err")

	pid1 := "placeid1"
	err = AddCityInfo(pid1, "CN")
	suite.NoError(err, "Should be able to add city info.")
	err = AddCityInfo(pid1, "CN")
	suite.NoError(err, "Should be able to add the same city info without any error.")

	cityName := "Xiamen"
	cityAddress := "Xiamen, Fujian, China"
	tp := 0

	UpdateCityInfo(pid1, cityName, cityAddress, tp)

	existed, resultName, resultAddress, err := GetCityInfo(pid1, tp)
	suite.True(existed, "The result should be existed.")
	suite.NoError(err, "Should be able to get.")
	suite.EqualValues(cityName, resultName, "The name should be equal")
	suite.EqualValues(cityAddress, resultAddress, "The address should be equal")

	noTp := 1
	existed, resultName, resultAddress, err = GetCityInfo(pid1, noTp)
	suite.True(existed, "The result should be existed.")
	suite.NoError(err, "Should be able to get.")
	suite.EqualValues("", resultName, "The name should be empty")
	suite.EqualValues("", resultAddress, "The address should be empty")

	noPlace := "placeid2"
	existed, _, _, err = GetCityInfo(noPlace, tp)
	suite.False(existed, "The place should be not existed.")
	suite.NoError(err, "Should be able to get.")

	tpOutRange := 2
	_, _, _, err = GetCityInfo(noPlace, tpOutRange)
	suite.Error(err, "Should have error while tp is out of range.")
}

func (suite *DBHandleSuite) TestCountryInfo() {
	id1 := "EN"
	id2 := "CN"
	name1 := "English"
	name2 := "Chinese"
	id2Lower := "cn"

	badID := "123"
	err := AddCountry(badID, "bad")
	suite.Equal(ErrCountryID, err, "country id format is wrong.")

	err = AddCountry(id1, name1)
	suite.NoError(err, "Should have no error.")

	err = AddCountry(id2Lower, name2)
	suite.NoError(err, "Shoule have no error.")

	ids, names := GetCountries()
	suite.EqualValues(2, len(ids), "Shoule have 2 result.")
	suite.EqualValues(2, len(names), "Shoule have 2 result.")
	for _, one := range ids {
		if one != id1 && one != id2 {
			suite.Fail("id is wrong.")
		}
	}

	for _, one := range names {
		if one != name1 && one != name2 {
			suite.Fail("name is wrong.")
		}
	}

	var existed bool
	existed, err = CheckCountryExisted(badID)
	suite.False(existed, "Country should not existed.")
	suite.EqualValues(ErrCountryID, err, "Should have bad country id error.")

	existed, err = CheckCountryExisted(id1)
	suite.True(existed, "Country should existed.")
	suite.NoError(err, "Should be able to get country.")
}
