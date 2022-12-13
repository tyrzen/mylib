goose postgres "host=localhost port=5432 user=postgres password=postgrespw dbname=postgres sslmode=disable"
goose dir ./mig

goose create readers sql
goose create authors sql
goose create books sql

goose fix