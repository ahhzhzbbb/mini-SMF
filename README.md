# mini-SMF

# Project Structure
```bash
.
├── database
├── docker-compose.yml
├── gateway
│   ├── cmd
│   │   └── main.go
│   ├── Dockerfile
│   └── internal
│       ├── config
│       │   └── config.go
│       ├── handler
│       │   └── handler.go
│       ├── middleware
│       │   ├── auth.go
│       │   └── logging.go
│       ├── proxy
│       │   ├── proxy.go
│       │   └── routes.go
│       ├── registry
│       │   └── registry.go
│       └── router
│           └── round_robin.go
├── go.mod
├── go.sum
├── pdu-session
│   ├── cmd
│   │   └── main.go
│   ├── Dockerfile
│   ├── internal
│   │   ├── config
│   │   │   └── config.go
│   │   ├── handler
│   │   │   └── handler.go
│   │   ├── middleware
│   │   │   └── logging.go
│   │   └── server
│   │       ├── routes.go
│   │       └── server.go
│   └── pdu-session
├── pkg
│   └── logger
│       └── logger.go
└── README.md

20 directories, 23 files
```