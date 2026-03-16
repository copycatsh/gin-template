# gin-template

REST API example — Go 1.22, Gin, GORM, MySQL 8, Docker

## Quick Start
```bash
make up
```

**API:** http://localhost:8054  
**DB:** localhost:3307 · `gouser / gopassword` · `godb`

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /users | List all users |
| GET | /users/:id | Get user by ID |
| POST | /users | Create user |
| PUT | /users/:id | Update user |
| DELETE | /users/:id | Delete user |

## Development
```bash
make test   # run tests
make lint   # golangci-lint
make fmt    # gofmt
```
