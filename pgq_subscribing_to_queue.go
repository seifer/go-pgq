// TODO Implememt full API

package pgq

import "fmt"

// Subscribe consumer on a queue.
// From this moment forward, consumer will see all events in the queue.
// Parameters
//      queue_name     Name of queue
//      consumer_name  Name of consumer
// Returns
//      0   if already registered
//      1   if new registration
// Calls:
//      pgq.register_consumer_at(3)
// Tables directly manipulated:
//       None
func (h *PGQHandle) RegisterConsumer(queue_name, consumer_name string) (out int, err error) {
	err = h.q.QueryRow(
		"SELECT pgq.register_consumer($1, $2)",
		queue_name,
		consumer_name,
	).Scan(&out)

	return
}

// Extended registration, allows to specify tick_id.
// Note
//      For usage in special situations.
// Parameters
//      queue_name     Name of a queue
//      consumer_name  Name of consumer
//      tick_pos       Tick ID
// Returns
//      0/1 whether consumer has already registered.
// Calls:
//      None
// Tables directly manipulated:
//      update/insert - pgq.subscription
func (h *PGQHandle) RegisterConsumerAt(queue_name, consumer_name string, tick_pos int) (out int, err error) {
	err = h.q.QueryRow(
		"SELECT pgq.register_consumer_at($1, $2, $3)",
		queue_name,
		consumer_name,
		tick_pos,
	).Scan(&out)

	return
}

// Unsubscribe consumer from the queue.  Also consumer’s retry events are deleted.
// Parameters
//      queue_name    Name of the queue
//      consumer_name Name of the consumer
// Returns
//      number of (sub)consumers unregistered
// Calls:
//      None
// Tables directly manipulated:
//      delete - pgq.retry_queue
//      delete - pgq.subscription
func (h *PGQHandle) UnregisterConsumer(queue_name, consumer_name string) (out int, err error) {
	err = h.q.QueryRow(
		"SELECT pgq.unregister_consumer($1, $2)",
		queue_name,
		consumer_name,
	).Scan(&out)

	return
}

// Full unsubscribe consumer from the queue.  Also consumer’s retry events are deleted. Non standart function
// Parameters
//      queue_name    Name of the queue
//      consumer_name Name of the consumer
// Returns
//      number of (sub)consumers unregistered
// Calls:
//      pgq.unregister_consumer(2)
// Tables directly manipulated:
//      delete - pgq.retry_queue
//      delete - pgq.subscription
// 		delete - pgq.consumer
func (h *PGQHandle) FullUnregisterConsumer(queue_name, consumer_name string) (out int, err error) {
	if out, err = h.UnregisterConsumer(queue_name, consumer_name); err != nil {
		return
	}

	_, err = h.q.Exec(fmt.Sprintf("DELETE FROM consumer WHERE co_name LIKE '%s.%%'", consumer_name))

	return
}
