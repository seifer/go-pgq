package pgq

import "database/sql"

// Makes next block of events active
// Result NULL means nothing to work with, for a moment
// Parameters
//      queue_name        Name of the queue
//      consumer_name     Name of the consumer
//      subconsumer_name  Name of the subconsumer
// Calls
//      pgq.register_consumer(i_queue_name, i_consumer_name);
//      pgq.register_consumer(i_queue_name, _subcon_name);
// Tables directly manipulated
//      update  pgq.subscription
func (h *PGQCOOPHandle) NextBatch(queue_name, consumer_name, subconsumer_name string) (int64, error) {
	var batch_id sql.NullInt64

	err := h.q.QueryRow(
		`SELECT pgq_coop.next_batch($1, $2, $3)`,
		queue_name,
		consumer_name,
		subconsumer_name,
	).Scan(&batch_id)

	return batch_id.Int64, err
}

// Closes a batch.
// Parameters
//      batch_id  id of the batch to be closed
// Returns
//      1 if success (batch was found), 0 otherwise
// Calls:
//      None
// Tables directly manipulated:
//      update - pgq.subscription
func (h *PGQCOOPHandle) FinishBatch(batch_id int64) (out int, err error) {
	err = h.q.QueryRow(
		`SELECT pgq_coop.finish_batch($1)`,
		batch_id,
	).Scan(&out)

	return out, nil
}
