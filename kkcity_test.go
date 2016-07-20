package kkcity

import (
	"net"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var testLangs = []string{"en", "zh"}

func TestMain(t *testing.T) {
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
	assert.NoError(t, err, "Should be able to create pool.")

	Use(testLangs, "", pool)

	suite.Run(t, new(dbHandleSuite))
	suite.Run(t, new(languageHandleSuite))
}
