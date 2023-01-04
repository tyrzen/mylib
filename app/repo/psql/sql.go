package psql

const (
	queryCreateReader = `INSERT INTO readers (id, first_name, last_name, email, password, created_at) 
							VALUES($1, $2, $3, $4, $5, NOW())--RETURNING *`
	queryCreateBook = ``
)
