package pgq

// Insert a event into queue.
// Parameters
//      queue_name  Name of the queue
//      ev_type     User-specified type for the event
//      ev_data     User data for the event
// Returns
//      Event ID Calls: pgq.insert_event(7)
func (h *PGQHandle) InsertEvent3(queue_name, ev_type, ev_data string) (out int, err error) {
	err = h.q.QueryRow(
		"SELECT pgq.insert_event($1, $2, $3)",
		queue_name,
		ev_type,
		ev_data,
	).Scan(&out)

	return
}

// Insert a event into queue with all the extra fields.
// Parameters
//      queue_name  Name of the queue
//      ev_type     User-specified type for the event
//      ev_data     User data for the event
//      ev_extra1   Extra data field for the event
//      ev_extra2   Extra data field for the event
//      ev_extra3   Extra data field for the event
//      ev_extra4   Extra data field for the event
// Returns
//      Event ID Calls: pgq.insert_event_raw(11)
// Tables directly manipulated:
//      insert - pgq.insert_event_raw(11), a C function, inserts into current event_N_M table
func (h *PGQHandle) InsertEvent7(queue_name, ev_type, ev_data, ev_extra1, ev_extra2, ev_extra3, ev_extra4 string) (out int, err error) {
	err = h.q.QueryRow(
		"SELECT pgq.insert_event($1, $2, $3, $4, $5, $6, $7)",
		queue_name,
		ev_type,
		ev_data,
		ev_extra1,
		ev_extra2,
		ev_extra3,
		ev_extra4,
	).Scan(&out)

	return
}

// Return active event table for particular queue.  Event can be added to it without going via functions, e.g. by COPY.
// If queue is disabled and GUC session_replication_role <> ‘replica’ then raises exception.
// or expressed in a different way an even table of a disabled queue is returned only on replica
// Note
//      The result is valid only during current transaction.
// Permissions
//      Actual insertion requires superuser access.
// Parameters
//      x_queue_name    Queue name.
func (h *PGQHandle) CurrentEventTable(queue_name string) (out string, err error) {
	err = h.q.QueryRow(
		"SELECT pgq.current_event_table($1)",
		queue_name,
	).Scan(&out)

	return
}
