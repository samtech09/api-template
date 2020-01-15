package psql

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"
)

type PgLogger struct {
	logger zerolog.Logger
}

// NewLogger accepts a zerolog.Logger as input and returns a new custom pgx
// logging fascade as output.
func newLogger(logger zerolog.Logger, connName string) *PgLogger {
	return &PgLogger{
		logger: logger.With().Str("module", connName).Logger(),
	}
}

func (pl *PgLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	var zlevel zerolog.Level
	switch level {
	case pgx.LogLevelNone:
		zlevel = zerolog.NoLevel
	case pgx.LogLevelError:
		zlevel = zerolog.ErrorLevel
	case pgx.LogLevelWarn:
		zlevel = zerolog.WarnLevel
	case pgx.LogLevelInfo:
		zlevel = zerolog.InfoLevel
	case pgx.LogLevelDebug:
		zlevel = zerolog.DebugLevel
	default:
		zlevel = zerolog.DebugLevel
	}

	pgxlog := pl.logger.With().Fields(data).Logger()
	pgxlog.WithLevel(zlevel).Msg(msg)
}
