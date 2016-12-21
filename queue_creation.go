package pgq

// Creates new queue with given name.
// Returns
//      0   queue already exists
//      1   queue created
// Calls:
// 		pgq.grant_perms(i_queue_name);
// 		pgq.ticker(i_queue_name);
// 		pgq.tune_storage(i_queue_name);
// Tables directly manipulated:
//      insert  pgq.queue
//      create  pgq.event_N () inherits (pgq.event_template)
//      create  pgq.event_N_0 .. pgq.event_N_M () inherits (pgq.event_N)
func (h *PGQHandle) CreateQueue(queue_name string) (out int, err error) {
	err = h.db.QueryRow(
		"SELECT pgq.create_queue($1)",
		queue_name,
	).Scan(&out)

	return
}

// Drop queue and all associated tables.
// Parameters
//      queue_name  queue name
//      force       ignore (drop) existing consumers Returns:
// Returns
//      1   success
// Calls:
// 		pgq.unregister_consumer(queue_name, consumer_name);
//		pgq.ticker(i_queue_name);
// 		pgq.tune_storage(i_queue_name);
// Tables directly manipulated:
//      delete  pgq.queue
//      drop    pgq.event_N (), pgq.event_N_0 .. pgq.event_N_M
func (h *PGQHandle) DropQueue2(queue_name string, force bool) (out int, err error) {
	err = h.db.QueryRow(
		"SELECT pgq.drop_queue($1, $2)",
		queue_name,
		force,
	).Scan(&out)

	return
}

// Drop queue and all associated tables.  No consumers must be listening on the queue.
func (h *PGQHandle) DropQueue(queue_name string) (out int, err error) {
	err = h.db.QueryRow(
		"SELECT pgq.drop_queue($1)",
		queue_name,
	).Scan(&out)

	return
}

// Set configuration for specified queue.
// Parameters
//      queue_name   Name of the queue to configure.
//      param_name   Configuration parameter name.
//      param_value  Configuration parameter value.
// Returns
//      0 if event was already in queue, 1 otherwise.
// Calls:
// 		None
// Tables directly manipulated:
// 		update - pgq.queue
func (h *PGQHandle) SetQueueConfig(queue_name, param_name, param_value string) (out int, err error) {
	err = h.db.QueryRow(
		"SELECT pgq.set_queue_config($1, $2, $3)",
		queue_name, param_name, param_value,
	).Scan(&out)

	return
}
