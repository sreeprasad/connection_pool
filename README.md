## Simulate connection pool and too many connections error

### Steps

Install postgres
Create new user with privileges. Creating new user is not strictly required. I
added this for creating users specific to this project.

### Too many connection test

```go
go run simulate.go simulate_too_many_conn_error.go
```

### Too simulate connection pool

```go
go run simulate.go connectionPool.go
```
