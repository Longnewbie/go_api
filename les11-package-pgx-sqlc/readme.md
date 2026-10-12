## Cài đặt Golang-Migrate

```curl
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

hoặc cài đặt theo tài liệu

```curl
https://github.com/golang-migrate/migrate/tree/v4.18.3/cmd/migrate
```

kiểm tra version
```curl
migrate --version
```

migrate create -ext sql -dir internal/db/migrations -seq users
migrate -path internal/db/migrations -database "postgresql://root:123456vV%40@localhost:5433/master-golang?sslmode=disable" up

migrate -path internal/db/migrations -database "postgresql://root:123456vV%40@localhost:5433/master-golang?sslmode=disable" down 1