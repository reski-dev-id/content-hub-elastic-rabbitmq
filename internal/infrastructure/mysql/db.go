package mysql

import (
	"time"

	"content-hub/internal/logger"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB(
	dsn string,
) (*sqlx.DB, error) {

	start := time.Now()

	db, err := sqlx.Connect(
		"mysql",
		dsn,
	)

	duration := time.Since(start).Milliseconds()

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "connection_failed").
			Int64("duration_ms", duration).
			Msg("failed to connect mysql")

		return nil, err
	}

	logger.Info().
		Str("service", "mysql").
		Str("event", "connected").
		Int64("duration_ms", duration).
		Msg("mysql connected")

	return db, nil
}
