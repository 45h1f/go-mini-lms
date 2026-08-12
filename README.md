# Mini LMS (Learning Management System)

Enterprise-level Go backend & Next.js frontend project structure.

## Directory Layout

```
mini-lms/
├── backend/                  # Go Enterprise Backend
│   ├── cmd/
│   │   ├── api/              # Main REST API application entry point
│   │   └── worker/           # Async background worker entry point
│   ├── internal/             # Private application code
│   │   ├── config/           # App configuration
│   │   ├── domain/           # Core domain models and interfaces
│   │   ├── handler/          # Transport layer handlers (HTTP/gRPC)
│   │   ├── middleware/       # HTTP Middlewares (Auth, CORS, Logging)
│   │   ├── pkg/              # Internal utilities (DB, Logger)
│   │   ├── repository/       # Data persistence layer
│   │   └── service/          # Business logic layer
│   ├── migrations/           # Database migration files
│   ├── docs/                 # API documentation / Swagger specs
│   ├── go.mod                # Go module file
│   └── README.md
│
├── frontend/                 # Next.js Frontend Application
│   ├── src/                  # Next.js App Router source code
│   ├── public/               # Static assets
│   ├── package.json
│   └── tsconfig.json
│
└── docker-compose.yml        # Development environment services setup
```
# go-mini-lms
