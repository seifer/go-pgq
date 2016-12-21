package pgq

type PGQCOOPHandle struct {
	PGQHandle
}

func NewPGQCOOPHandle(q Querier) *PGQCOOPHandle {
	return &PGQCOOPHandle{PGQHandle{q: q}}
}
