package pgq

import (
	"database/sql"
	"time"
)

type Querier interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

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
	q Querier
}

func NewPGQHandle(q Querier) *PGQHandle {
	return &PGQHandle{q: q}
}
