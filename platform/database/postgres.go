package database

import (
	"fmt"
	"github.com/Xacnio/aptms-backend/pkg/configs"
	"github.com/Xacnio/aptms-backend/pkg/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"strconv"
	"time"
)

func PostgreSQLConnection() (*sqlx.DB, error) {
	maxConn, _ := strconv.Atoi(configs.Get("DB_MAX_CONNECTIONS"))
	maxIdleConn, _ := strconv.Atoi(configs.Get("DB_MAX_IDLE_CONNECTIONS"))
	maxLifetimeConn, _ := strconv.Atoi(configs.Get("DB_MAX_LIFETIME_CONNECTIONS"))

	postgresConnURL, err := utils.ConnectionURLBuilder("postgres")
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Connect("postgres", postgresConnURL)
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	db.SetMaxOpenConns(maxConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(time.Duration(maxLifetimeConn))

	if err := db.Ping(); err != nil {
		defer db.Close()
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	return db, nil
}
