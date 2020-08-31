package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joao-fontenele/go-url-shortener/pkg/common"
	"go.uber.org/zap"
)

var conn *pgxpool.Pool

// Connect creates a connection to postgres db
func Connect() (func(), error) {
	var err error
	logger := common.GetLogger()
	dbConf := common.GetConf().Database
	dbURL := fmt.Sprintf(
		// "postgresql://%s:%s@%s:%s/%s",
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		dbConf.User,
		dbConf.Pass,
		dbConf.Host,
		dbConf.Port,
		dbConf.Name,
		dbConf.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		logger.Fatal("failed to parse dburl", zap.String("dbUrl", dbURL), zap.Error(err))
	}
	poolConfig.ConnConfig.Logger = zapadapter.NewLogger(logger)

	conn, err = pgxpool.ConnectConfig(context.Background(), poolConfig)

	if err != nil {
		return nil, err
	}

	return conn.Close, nil
}

// GetConnection returns a previously created connection pool
func GetConnection() *pgxpool.Pool {
	return conn
}
