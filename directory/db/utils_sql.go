package db

const (
	sqlInsertGenericEvent = `
		INSERT INTO generic_events (event_type, txid, height, attributes)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
)
