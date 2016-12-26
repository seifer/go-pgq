package pgq

import (
	"database/sql"
	"fmt"
)

// Old function that returns just batch_id.
// Parameters
//      queue_name    Name of the queue
//      consumer_name Name of the consumer
// Returns
//      Batch ID or NULL if there are no more events available.
func (h *PGQHandle) NextBatch(queue_name, consumer_name string) (int64, error) {
	var batch_id sql.NullInt64

	err := h.q.QueryRow(
		`SELECT pgq.next_batch($1, $2)`,
		queue_name,
		consumer_name,
	).Scan(&batch_id)

	return batch_id.Int64, err
}

// Get all events in batch.
// Parameters
//      batch_id  ID of active batch.
// Returns
//      List of events.
func (h *PGQHandle) GetBatchEvents(batch_id int64) ([]Event, error) {
	rows, err := h.q.Query(
		`SELECT
            ev_id,
            ev_time,
            ev_txid,
            ev_retry,
            ev_type,
            ev_data,
            ev_extra1,
            ev_extra2,
            ev_extra3,
            ev_extra4
        FROM
            pgq.get_batch_events($1)
        `,
		batch_id,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var batch []Event

	for rows.Next() {
		var event Event

		err := rows.Scan(
			&event.Ev_id,
			&event.Ev_time,
			&event.Ev_txid,
			&event.Ev_retry,
			&event.Ev_type,
			&event.Ev_data,
			&event.Ev_extra1,
			&event.Ev_extra2,
			&event.Ev_extra3,
			&event.Ev_extra4,
		)

		if err != nil {
			return nil, fmt.Errorf("GetBatchEvents error during row processing: %s. %+v", err, rows)
		}

		batch = append(batch, event)
	}

	return batch, nil
}

// Closes a batch.  No more operations can be done with events of this batch.
// Parameters
//      batch_id  id of batch.
// Returns
//      1 if batch was found, 0 otherwise.
// Calls: None Tables directly manipulated: update - pgq.subscription
func (h *PGQHandle) FinishBatch(batch_id int64) (out int, err error) {
	err = h.q.QueryRow(
		`SELECT pgq.finish_batch($1)`,
		batch_id,
	).Scan(&out)

	return out, nil
}

// Put the event into retry queue, to be processed later again.
// Parameters
//      batch_id      ID of active batch.
//      event_id      event id
//      retry_seconds Time when the event should be put back into queue
// Returns
//      1   success
//      0   event already in retry queue
// Calls:
//      pgq.event_retry(3a)
// Tables directly manipulated:
//      None
func (h *PGQHandle) EventRetry3b(batch_id, event_id, retry_seconds int64) (out int, err error) {
	err = h.q.QueryRow(
		`SELECT pgq.event_retry($1, $2, $3::integer)`,
		batch_id,
		event_id,
		retry_seconds,
	).Scan(&out)

	return out, nil
}
