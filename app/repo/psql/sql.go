package psql

const (
	queryCreateReader = `INSERT INTO readers (id, first_name, last_name, email, password, created_at) 
							VALUES(gen_random_uuid(), $1, $2, $3, $4, NOW())--RETURNING *`
	queryCreateBook = ``
)
