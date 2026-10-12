# Go API Learning Repository

Repository này là tập hợp các bài thực hành Go API theo từng chủ đề từ cơ bản đến nâng cao. Mỗi thư mục `les*` là một mini project riêng, giúp bạn làm quen với các khái niệm và kỹ thuật phổ biến trong phát triển API bằng Go.

## Mục tiêu

- Học cách xây dựng API bằng Go
- Tìm hiểu Gin framework
- Làm việc với routing, validation, middleware
- Kết nối PostgreSQL và Redis
- Sử dụng GORM, pgx, sqlc, migration
- Thiết kế kiến trúc theo tầng (`handler` -> `service` -> `repository` -> `model`)
- Xây dựng một project thương mại điện tử thực tế ở module cuối

## Cấu trúc repository

- `les1-package-http` – Giới thiệu package HTTP cơ bản
- `les2-gin-starter` – Khởi tạo project với Gin
- `les3-route-group` – Route grouping và tổ chức endpoint
- `les4-route-validation` – Validation, middleware và xử lý request
- `les5-user-manager` – Project quản lý người dùng theo kiến trúc rõ ràng
- `les6-docker-postgresql` – Docker + PostgreSQL
- `les7-golang-migrate` – Sử dụng `golang-migrate` để quản lý migration
- `les8-prepare-connection` – Chuẩn bị kết nối database
- `les9-package-db-sql` – Kết nối database bằng `database/sql`
- `les10-package-gorm` – Sử dụng GORM
- `les11-package-pgx-sqlc` – Kết nối PostgreSQL với pgx và sqlc
- `les12-project-shopping-cart` – Project cuối: API thương mại điện tử đầy đủ

## Project chính

Project quan trọng nhất trong repo là:

- `les12-project-shopping-cart`

Đây là phiên bản API hoàn chỉnh với các kỹ thuật thường dùng trong production như:

- Gin framework
- JWT authentication
- PostgreSQL
- Redis
- RabbitMQ / background worker
- Logging
- Validation
- Environment config
- Database migration
- Structured project architecture

## Yêu cầu môi trường

- Go 1.22+
- Docker + Docker Compose
- PostgreSQL
- Redis
- Make (tùy chọn, hỗ trợ chạy lệnh nhanh)

## Chạy project chính

Di chuyển vào thư mục project cuối:

```bash
cd les12-project-shopping-cart
```

Cài đặt dependencies:

```bash
go mod download
```

Nếu bạn dùng Docker toàn bộ môi trường:

```bash
make dev
```

Hoặc chạy ứng dụng trực tiếp:

```bash
make server
```

Hoặc build binary:

```bash
make build
```

## Migration database

Tạo migration mới:

```bash
make migrate-create NAME=profiles
```

Chạy migration:

```bash
make migrate-up
```

Rollback migration cuối cùng:

```bash
make migrate-down
```

## Cấu trúc thư mục mẫu

```text
les12-project-shopping-cart/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── app/
│   ├── config/
│   ├── db/
│   ├── dto/
│   ├── handler/
│   ├── logger/
│   ├── middleware/
│   ├── models/
│   ├── repository/
│   ├── routes/
│   ├── service/
│   └── utils/
├── .env
├── docker-compose.dev.yml
├── docker-compose.prod.yml
├── Makefile
├── go.mod
└── readme.md
```

## Gợi ý học tập

Nên học theo thứ tự sau:

1. `les1-package-http`
2. `les2-gin-starter`
3. `les3-route-group`
4. `les4-route-validation`
5. `les5-user-manager`
6. `les8-prepare-connection`
7. `les9-package-db-sql`
8. `les10-package-gorm`
9. `les11-package-pgx-sqlc`
10. `les12-project-shopping-cart`

Mỗi module sau đều xây dựng trên kiến thức của các module trước đó.

## Lưu ý

- Một số thư mục `les*` là ví dụ học tập, không phải lúc nào cũng hoàn chỉnh như production project.
- Để chạy chính xác, hãy kiểm tra file `.env` và các cấu hình database trong từng project.
- Nếu cần, bạn có thể mở `Makefile` trong từng module để xem các lệnh hỗ trợ nhanh.

## Kết luận

Repo này là một roadmap học tập thực tế cho việc xây dựng API bằng Go, từ cơ bản cho đến ứng dụng có cấu trúc, tích hợp database và service thực tế. Bạn có thể sử dụng nó như tài liệu tham khảo, bài tập thực hành hoặc nền tảng để phát triển project riêng.
