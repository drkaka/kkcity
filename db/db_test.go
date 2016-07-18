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

		columnName := getCityNameColumn(i)
		existed, err = CheckColumnExisted("city_info", columnName)
		suite.True(existed, columnName, " should be existed.")
		suite.Nil(err, "There should be no error while check exist.")

		columnAddress := getCityAddressColumn(i)
		existed, err = CheckColumnExisted("city_info", columnAddress)
		suite.True(existed, columnAddress, " should be existed.")
		suite.Nil(err, "There should be no error while check exist.")
	}
}

func (suite *DBHandleSuite) TestCityInfo() {

}
