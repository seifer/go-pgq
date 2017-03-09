package pgq

import "fmt"

// Subscribe a subconsumer on a queue.
// Subconsumer will be registered as another consumer on queue, whose name will be i_consumer_name and i_subconsumer_name combined.
// Returns
//      0   if already registered
//      1   if this is a new registration
// Calls
//      pgq.register_consumer(i_queue_name, i_consumer_name);
//      pgq.register_consumer(i_queue_name, _subcon_name);
// Tables directly manipulated
//      update  pgq.subscription
func (h *PGQCOOPHandle) RegisterSubconsumer(queue_name, consumer_name, subconsumer_name string) (out int, err error) {
	err = h.q.QueryRow(
		"SELECT pgq_coop.register_subconsumer($1, $2, $3)",
		queue_name,
		consumer_name,
		subconsumer_name,
	).Scan(&out)

	return
}

// Unregisters subconsumer from the queue.
// If consumer has active batch, then behviour depends on i_batch_handling parameter.
// Values for i_batch_handling
//      0   Fail with an exception.
//      1   Close the batch, ignoring the events.
// Returns
//      0   no consumer found
//      1   consumer found and unregistered
// Tables directly manipulated
//          delete  pgq.subscription
func (h *PGQCOOPHandle) UnregisterSubconsumer(queue_name, consumer_name, subconsumer_name string, batch_handling int) (out int, err error) {
	err = h.q.QueryRow(
		"SELECT pgq_coop.unregister_subconsumer($1, $2, $3, $4)",
		queue_name,
		consumer_name,
		subconsumer_name,
		batch_handling,
	).Scan(&out)

	return
}

// Full unsubscribe subconsumers from the queue. Non standart function
// Parameters
//      queue_name    Name of the queue
//      consumer_name Name of the consumer
// Returns
//      None
// Calls:
//      pgq.unregister_consumer(2)
// Tables directly manipulated:
// 		delete - pgq.consumer
func (h *PGQHandle) FullUnregisterSubconsumers(queue_name, consumer_name string) (err error) {
	_, err = h.q.Exec(fmt.Sprintf(`
		do $$
		DECLARE
		    a text[];
		    r RECORD;
		BEGIN
		    FOR r IN
		        SELECT
		            consumer_name
		        FROM
		            pgq.get_consumer_info('%s')
		        WHERE
		            consumer_name != '%s'
		    LOOP
		        a := regexp_split_to_array(r.consumer_name, '\.');
		        PERFORM pgq_coop.unregister_subconsumer('%s', a[1], a[2], 1);
		    END LOOP;
		END $$;`,
		queue_name,
		consumer_name,
		queue_name,
	))

	return
}
