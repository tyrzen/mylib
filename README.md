# mylib - API based on REST and Clean Architecture principles.

## Structure

```
mylib/
├── app/
│   ├── exceptions/
│   │   └── errors.go
│   ├── models/
│   │   ├──abstract.go
│   │   ├──author.go
│   │   ├──book.go
│   │   ├──credentials.go
│   │   ├──filter.go
│   │   ├──reader.go
│   │   └──token.go
│   ├── presenters/
│   │   └── rest/
│   │       ├── abstract.go
│   │       ├── book_handler.go
│   │       ├── const.go
│   │       ├── cookie.go
│   │       ├── errors.go
│   │       ├── middleware.go
│   │       ├── reader_handler.go
│   │       ├── responder.go
│   │       ├── router.go
│   │       ├── server.go
│   │       └── token.go
│   ├─── repository/
│   │    ├── psql/
│   │    │   ├── author.go 
│   │    │   ├── book.go
│   │    │   ├── conn.go
│   │    │   ├── filter.go
│   │    │   └── reader.go
│   │    └── rds/
│   │        ├── client.go
│   │        └── token.yml
│   └── usecases/
│       ├── abstract.go 
│       ├── author.go   
│       ├── book.go   
│       ├── reader.go   
│       └── token.go   
├── cmd/
│   └── main.go 
├── doc/
│   └── openapi.yml 
├── lib/
│   ├── banderlog/
│   │    └── logger.go
│   ├── env/
│   │    └── load.go
│   ├── hash/
│   │    └── hash.go
│   ├── revalid/
│   │    └── validator.go
│   └── tokay/
│       └── jwt.go
├── mig/
│   ├── ######_*.sql
│   └── mig.go 
├── .env 
├── .gitignore 
├── .golangci.yml 
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── LICENSE
├── log.json
├── Makefile
└── README.md
```

# Helpful materials

* [HTTP status codes with kittens](https://httpcats.com/)
* [Letsencrypt](https://letsencrypt.org/) gives production ready certificates for free.

