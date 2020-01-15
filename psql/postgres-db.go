package psql

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

type Db struct {
	//Reader is readonly connection to database
	reader *pgxpool.Pool

	//Writer is read/write connection to database
	writer *pgxpool.Pool
}

//Conn provide reader or writer connection as per readonly state
func (db *Db) Conn(readonly bool) *pgxpool.Pool {
	if readonly {
		return db.reader
	}
	return db.writer
}

//PDbConfig is config for disk persistent database (PostgreSQL)
type PDbConfig struct {
	PDbHost string
	PDbPort uint16
	PDbName string
	PDbUser string
	PDbPwd  string
	//DbTimeout is Connection timeout in seconds
	//if could not connect to server in given time then giveup and raise error
	PDbTimeout int
	//DbSSLMode flag to enable disable SSL for database connection
	PDbSSLMode string
}

//InitDbPool Initialize database connection ppol for PostgreSQL database
func InitDbPool(reader, writer PDbConfig, l zerolog.Logger) *Db {
	db := Db{}
	db.reader = initdb(reader, "pgx-reader", l)
	db.writer = initdb(writer, "pgx-writer", l)

	return &db
}

func initdb(config PDbConfig, connName string, l zerolog.Logger) *pgxpool.Pool {
	var err error

	s := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?connect_timeout=%d&sslmode=%s",
		config.PDbUser, config.PDbPwd, config.PDbHost, config.PDbPort, config.PDbName, config.PDbTimeout, config.PDbSSLMode)

	cfg, err := pgxpool.ParseConfig(s)
	if err != nil {
		log.Fatalf("Psql Connection parse error %v\n", err)
	}

	//cfg.LogLevel = pgx.LogLevelWarn
	//cfg.Logger = newLogger(l, connName)

	// pgxConnPoolConfig := pgxpool.Config{
	// 	ConnConfig:   cfg,
	// 	MaxConns:     8,
	// 	AfterConnect: nil,
	// }

	p, err := pgxpool.ConnectConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalf("Psql Connection error %v\n", err)
	}
	return p
}
