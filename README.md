# mylib

## Structure
```
mylib/
├── app/
│   ├── ent/
│   │   └── user.go
│   ├── repo/
│   │   └── pg/
│   │         └── user.sql
│   ├── trans/
│   │   └── rest/
│   │       ├── contracts.go
│   │       ├── handler.go
│   │       ├── resp.go
│   │       ├── router.go
│   │       ├── server.go
│   │       └── user.go
│   └── logic/
│       ├── contracts.go
│       └── user.goo
├── cfg/
│   ├── config.go
│   └── config.yml
├── cmd/
│   └── rest/
│       └── main.go
├── lib/
│   └── revalid/
│       └── revalid.go
├── mig/
│   └── ###_migration.sql
├── docker-compose.yml
├── Dockerfile
├── README.md
└── .env 
```
# Helpful materials
* [HTTP status codes with kittens](https://httpcats.com/)