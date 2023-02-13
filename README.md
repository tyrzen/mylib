# mylib - API based on REST and Clean Architecture principles.

## Technologies

| technology                                                     | purpose                                                                            |
|----------------------------------------------------------------|------------------------------------------------------------------------------------|
| [Go](https://go.dev/)                                          | awesome programming language                                                       |
| [Git](https://git-scm.com/)                                    | by default                                                                         |
| [OpenAPI](https://www.openapis.org/)                           | used for describing the structure, operations, parameters, and responses of an API |
| [Makefile](https://www.gnu.org/software/make/manual/make.html) | helps automate tasks, such as building, testing, and deploying applications.       |
| [Docker](https://www.docker.com/)                              | used to containerize the application and its dependencies                          |
| [Postgres](https://www.postgresql.org/)                        | system used to store and retrieve main data                                        |
| [Redis](https://redis.io/)                                     | An in-memory data structure store user session data                                |
| [OData](https://www.odata.org/)                                | used as protocol for querying and ~~updating~~ data using query parameters         |
| [Letsencrypt](https://letsencrypt.org/)                        | open certificate authority                                                         |

## Awesome Dependencies

| package                                                    | description             |
|------------------------------------------------------------|-------------------------|
| [autocert](golang.org/x/crypto/acme/autocert)              | access to certificates  |
| [bcrypt](golang.org/x/crypto/bcrypt)                       | hash validation         |
| [chi](github.com/go-chi/chi/v5)                            | router                  |
| [errors](github.com/pkg/errors)                            | graceful error-handling |
| [goose](github.com/pressly/goose/v3)                       | SQL migrations          |
| [golangci-lint](https://github.com/golangci/golangci-lint) | Go liners aggregator    |
| [golang-jwt](github.com/golang-jwt/jwt/v4)                 | JSON Web Tokens         |
| [pgx](github.com/jackc/pgx/v5/stdlib)                      | postgres driver         |
| [uuid](github.com/google/uuid)                             | UUID generator          |
| [zap](go.uber.org/zap)                                     | logger                  |

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

## Helpful materials

* [HTTP status codes with kittens](https://httpcats.com/)

