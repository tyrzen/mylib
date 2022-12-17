# mylib
> Let's go
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
├── mig/
│   └── ###_migration.sql
├── docker-compose.yml
├── Dockerfile
├── README.md
└── .env 
```