package psql

const (
	queryCreateReader = `INSERT INTO readers (id, first_name, last_name, email, password, created_at) 
							VALUES(gen_random_uuid(), $1, $2, $3, $4, NOW())--RETURNING *;`

	queryFindReader = `SELECT id, first_name, last_name, email, password, created_at 
								FROM readers 
								WHERE 
								    CASE 
								        WHEN $1!='' AND $2!='' THEN id=$1 AND email=$2
								        WHEN $1!='' THEN id=$1
								       	WHEN $2!='' THEN email=$2
									END;`
)
