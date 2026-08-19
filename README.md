# Simple Grab — Backend (Go)

Backend API cho hệ thống gọi xe kiểu Grab (web), viết bằng Go + Gin + GORM + PostgreSQL.

## Yêu cầu

- Go 1.22+
- Docker (chạy PostgreSQL qua `docker-compose`)

## Setup & chạy local

```bash
# 1. Khởi động PostgreSQL
docker compose up -d

# 2. Copy file env
copy .env.example .env   # Windows
# cp .env.example .env   # macOS/Linux

# 3. Áp migration
psql "postgres://grab:grab123@localhost:5433/grab_db?sslmode=disable" -f migrations/001_init.up.sql
psql "postgres://grab:grab123@localhost:5433/grab_db?sslmode=disable" -f migrations/002_active_ride_constraints.up.sql

# 4. Chạy server
go run ./cmd/api
```

Server mặc định chạy ở `http://localhost:8080`.

## Biến môi trường (`.env`)

| Biến | Mô tả | Mặc định |
|---|---|---|
| `PORT` | Cổng HTTP server | `8080` |
| `DATABASE_URL` | Connection string PostgreSQL | `postgres://grab:grab123@localhost:5433/grab_db?sslmode=disable` |
| `JWT_SECRET` | Secret ký JWT | — (đổi khi deploy) |
| `JWT_EXPIRY_HOURS` | Thời gian sống của access token (giờ) | `24` |
| `GIN_MODE` | `debug` hoặc `release` | `debug` |

## Kiểm thử

```bash
go test ./...
go test -race ./...
go vet ./...
```

Integration test cần một PostgreSQL riêng (tên database phải chứa `test`):

```bash
export TEST_DATABASE_URL="postgres://grab:grab123@localhost:5432/grab_test?sslmode=disable"
go test ./internal/server -run TestIntegration -v
```

Xem chi tiết ma trận HTTP status và cách chạy Postman collection trong [PHASE6_TESTING.md](PHASE6_TESTING.md).

## Danh sách endpoint

Base URL: `http://localhost:8080/api/v1`. Các route có 🔒 cần header `Authorization: Bearer <token>` (lấy từ response của `/auth/login`).

### Health

| Method | Path | Mô tả |
|---|---|---|
| GET | `/health` | Kiểm tra server sống |

```bash
curl http://localhost:8080/api/v1/health
```

### Auth

| Method | Path | Mô tả |
|---|---|---|
| POST | `/auth/register` | Đăng ký tài khoản (`role`: `rider` hoặc `driver`) |
| POST | `/auth/login` | Đăng nhập, trả về `access_token` |
| GET 🔒 | `/auth/me` | Lấy thông tin user hiện tại |

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"rider1@example.com","password":"secret123","role":"rider"}'

curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"rider1@example.com","password":"secret123"}'

curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer <access_token>"
```

### Rides 🔒 (yêu cầu đăng nhập)

| Method | Path | Role | Mô tả |
|---|---|---|---|
| POST | `/rides` | rider | Tạo chuyến đi mới |
| GET | `/rides` | rider/driver | Danh sách chuyến của chính mình |
| GET | `/rides/available` | driver | Danh sách chuyến đang chờ (pending, chưa có tài xế) |
| GET | `/rides/:id` | rider/driver | Xem chi tiết 1 chuyến (chỉ chủ sở hữu) |
| POST | `/rides/:id/accept` | driver | Nhận chuyến |
| PATCH | `/rides/:id/status` | driver | Cập nhật trạng thái (`in_progress`, `completed`) |
| POST | `/rides/:id/cancel` | rider/driver | Hủy chuyến |

```bash
# Tạo chuyến (rider)
curl -X POST http://localhost:8080/api/v1/rides \
  -H "Authorization: Bearer <rider_token>" \
  -H "Content-Type: application/json" \
  -d '{"pickup_lat":10.7769,"pickup_lng":106.7009,"dropoff_lat":10.7629,"dropoff_lng":106.6602}'

# Xem chuyến đang chờ (driver)
curl http://localhost:8080/api/v1/rides/available \
  -H "Authorization: Bearer <driver_token>"

# Nhận chuyến (driver)
curl -X POST http://localhost:8080/api/v1/rides/<ride_id>/accept \
  -H "Authorization: Bearer <driver_token>"

# Bắt đầu / hoàn thành chuyến (driver)
curl -X PATCH http://localhost:8080/api/v1/rides/<ride_id>/status \
  -H "Authorization: Bearer <driver_token>" \
  -H "Content-Type: application/json" \
  -d '{"status":"in_progress"}'

# Hủy chuyến (rider hoặc driver)
curl -X POST http://localhost:8080/api/v1/rides/<ride_id>/cancel \
  -H "Authorization: Bearer <token>"
```

### Drivers 🔒 (role `driver`)

| Method | Path | Mô tả |
|---|---|---|
| PATCH | `/drivers/me/online` | Bật/tắt trạng thái online |
| PATCH | `/drivers/me/location` | Cập nhật vị trí hiện tại |

```bash
curl -X PATCH http://localhost:8080/api/v1/drivers/me/online \
  -H "Authorization: Bearer <driver_token>" \
  -H "Content-Type: application/json" \
  -d '{"is_online":true}'

curl -X PATCH http://localhost:8080/api/v1/drivers/me/location \
  -H "Authorization: Bearer <driver_token>" \
  -H "Content-Type: application/json" \
  -d '{"latitude":10.7769,"longitude":106.7009}'
```

## Luồng chuyến đi (ride flow)

```
pending --(driver accept)--> accepted --(driver update status)--> in_progress --(driver update status)--> completed
   |                            |
   +-------- cancel ------------+
```

- Rider tạo chuyến → status `pending`.
- Driver (phải đang `online`) accept → status `accepted`, ride bị khóa với driver đó.
- Driver cập nhật `in_progress` rồi `completed` theo thứ tự; fare được tính bằng công thức Haversine khi hoàn thành.
- Rider hoặc driver có thể `cancel` ở các trạng thái chưa `completed`/`cancelled`.

## Mã lỗi HTTP dùng chung

| Status | Khi nào |
|---|---|
| 400 | JSON/UUID/tọa độ hoặc dữ liệu đầu vào không hợp lệ |
| 401 | Thiếu/sai Bearer token hoặc sai thông tin đăng nhập |
| 403 | Đúng token nhưng sai role hoặc không sở hữu tài nguyên |
| 404 | Ride/driver không tồn tại |
| 409 | Email trùng, ride đã được nhận, đang có active ride, sai thứ tự trạng thái |
| 500 | Lỗi không dự kiến |
| 504 | Request vượt quá deadline (30s) |

## Phạm vi đã làm / chưa làm

**Đã làm (Phase 1–6):**
- Đăng ký/đăng nhập, JWT auth, phân quyền theo role (`rider`/`driver`)
- Tạo, xem, liệt kê, nhận, cập nhật trạng thái, hủy chuyến đi
- Tính cước theo khoảng cách (Haversine)
- Driver online/offline, cập nhật vị trí, chỉ driver online mới nhận được chuyến
- Unit test (fare, transition trạng thái, middleware), integration test full flow, Postman collection
- Middleware chung: CORS, timeout, recovery, error handler nhất quán

**Chưa làm (kế hoạch Phase 8+):**
- Đánh giá (rating) chuyến đi, lịch sử chuyến, thống kê tài xế (Phase 8)
- Lý do hủy chi tiết + phí hủy + audit log (Phase 9)
- Theo dõi vị trí real-time qua WebSocket, ETA (Phase 10)
- Các phase sau đó theo `PROJECT_FLOW_AND_TASKS.md` (thanh toán, thông báo, v.v. — nếu có)

Xem toàn bộ roadmap chi tiết trong [PROJECT_FLOW_AND_TASKS.md](PROJECT_FLOW_AND_TASKS.md).
