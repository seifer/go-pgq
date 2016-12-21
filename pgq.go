package pgq

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type Event struct {
	Ev_id     int64
	Ev_time   time.Time
	Ev_txid   int64
	Ev_retry  sql.NullInt64
	Ev_type   sql.NullString
	Ev_data   sql.NullString
	Ev_extra1 sql.NullString
	Ev_extra2 sql.NullString
	Ev_extra3 sql.NullString
	Ev_extra4 sql.NullString
}

type PGQHandle struct {
	db *sql.DB
}

func NewPGQHandle(dsn string) *PGQHandle {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	return &PGQHandle{
		db: db,
	}
}
