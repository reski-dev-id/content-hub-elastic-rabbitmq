package mysql

import (
	"content-hub/internal/logger"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB(dsn string) (*sqlx.DB, error) {

	db, err := sqlx.Connect("mysql", dsn)

	if err != nil {

		logger.Log.Error().
			Err(err).
			Msg("failed connect mysql")

		return nil, err
	}

	logger.Log.Info().
		Msg("mysql connected")

	return db, nil
}
