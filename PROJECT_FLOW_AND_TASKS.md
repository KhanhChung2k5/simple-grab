# Grab-like Backend Project (Web + Go)

## 1) Mục tiêu dự án

Xây dựng hệ thống backend giống Grab cho nền tảng web (không làm mobile app), sử dụng Go.

**Mốc phát triển:**
- **MVP (Phase 0–8):** Luồng đặt xe cơ bản + đánh giá + lịch sử
- **V1 (Phase 9–16):** Gần trải nghiệm Grab ride-hailing (real-time, matching, thanh toán, thông báo…)
- **V2 (Phase 17+):** Vận hành, bảo mật nâng cao, tính năng mở rộng

---

## 2) Phạm vi theo giai đoạn

### MVP (Phase 0–8) — Đã lên kế hoạch
- Auth với JWT (rider/driver)
- Quản lý driver online/offline + vị trí
- Tạo và quản lý ride (manual accept)
- Đánh giá tài xế sau chuyến đi
- Thống kê tài xế: số đơn hoàn thành, điểm trung bình
- Lịch sử chuyến của khách kèm thông tin tài xế và chuyến đi
- Tính giá tạm tính đơn giản
- PostgreSQL + REST API cho web frontend

### V1 (Phase 9–16) — Bổ sung so với Grab ride-hailing
- Hủy chuyến + phí hủy
- Real-time tracking + ETA (WebSocket)
- Auto-matching tài xế theo bán kính
- Nhiều loại xe + địa chỉ lưu
- OTP xác minh SĐT
- Thanh toán + hóa đơn
- Thông báo + chat/gọi ẩn số
- Voucher + surge pricing + đặt trước
- KYC tài xế + admin dashboard cơ bản

### V2 (Phase 17+) — Ngoài MVP ride, gần hệ sinh thái Grab
- Đánh giá 2 chiều, tip, chia sẻ chuyến, SOS
- Routing thật (Google Maps / OSRM)
- Driver incentive, fraud detection, support ticket
- Loyalty / referral
- Deploy production + monitoring

### Không nằm trong scope hiện tại (super app)
- GrabFood, GrabMart, GrabExpress
- GrabPay / ví đa dịch vụ đầy đủ
- Microservices toàn hệ thống

---

## 3) So sánh nhanh với app Grab

| Nhóm tính năng | Grab app | Plan hiện tại |
|----------------|----------|---------------|
| Đăng ký/đăng nhập | Có (OTP, social) | Phase 3 (email/password) |
| Đặt xe cơ bản | Có | Phase 4 |
| Driver online + vị trí | Có | Phase 5 |
| Đánh giá tài xế | Có | Phase 8 |
| Lịch sử + thống kê | Có | Phase 8 |
| Real-time map/ETA | Có | Phase 10 |
| Auto-match tài xế | Có | Phase 11 |
| Hủy chuyến + phí | Có | Phase 9 |
| Nhiều loại xe | Có | Phase 12 |
| Địa chỉ yêu thích | Có | Phase 12 |
| OTP SĐT | Có | Phase 12 |
| Thanh toán đa phương thức | Có | Phase 13 |
| Thông báo push/SMS | Có | Phase 14 |
| Chat/gọi trong app | Có | Phase 14 |
| Voucher / surge / đặt trước | Có | Phase 15 |
| KYC tài xế | Có | Phase 16 |
| Admin / ops | Có | Phase 16 |
| Food / Express / Mart | Có | Không scope |

---

## 4) Luồng hoạt động hệ thống

### 4.1 Đăng ký/Đăng nhập
1. Người dùng (rider hoặc driver) đăng ký tài khoản.
2. Hệ thống hash password (bcrypt) và lưu vào DB.
3. Người dùng đăng nhập, hệ thống trả JWT access token.
4. Mọi request cần bảo mật đều gửi Bearer token.
5. *(V1)* Xác minh SĐT bằng OTP trước khi kích hoạt tài khoản.

### 4.2 Rider tạo chuyến
1. Rider chọn loại xe *(V1)* và địa chỉ pickup/dropoff (lat, lng hoặc địa chỉ lưu).
2. Backend validate dữ liệu.
3. Backend tính khoảng cách + giá tạm tính *(V1: áp surge + voucher)*.
4. Tạo bản ghi ride với trạng thái `pending`.
5. *(V1)* Hệ thống auto-match tài xế gần nhất thay vì chờ driver tự browse.
6. Trả về thông tin ride cho rider.

### 4.3 Driver nhận chuyến
1. Driver bật trạng thái `online` *(V1: chỉ sau khi KYC approved)*.
2. Driver nhận chuyến qua auto-match *(V1)* hoặc xem danh sách `pending` *(MVP)*.
3. Driver accept chuyến.
4. Backend kiểm tra ride chưa có driver.
5. Gán `driver_id`, đổi status thành `accepted`.

### 4.4 Chạy chuyến và hoàn thành
1. Driver cập nhật status `in_progress` khi đón khách.
2. *(V1)* Rider theo dõi vị trí tài xế real-time qua WebSocket.
3. Driver cập nhật status `completed` khi kết thúc.
4. *(V1)* Thanh toán và phát hóa đơn.
5. Rider/Driver xem lịch sử ride theo role.

### 4.5 Hủy chuyến *(V1 — Phase 9)*
1. Rider hoặc driver gửi yêu cầu hủy kèm lý do.
2. Hệ thống kiểm tra trạng thái ride cho phép hủy.
3. Tính phí hủy theo rule (thời điểm hủy, ai hủy).
4. Cập nhật status `cancelled`, ghi audit log.

### 4.6 Đánh giá tài xế và thống kê
1. Sau khi ride `completed`, rider đánh giá tài xế (1–5 sao + comment).
2. Mỗi ride chỉ đánh giá một lần.
3. Hệ thống tính `avg_rating` và `completed_rides` cho tài xế.
4. *(V2)* Tài xế đánh giá lại rider (2 chiều).

### 4.7 Lịch sử đơn
1. Tài xế xem lịch sử chuyến đã nhận/chạy.
2. Khách xem lịch sử chuyến kèm thông tin tài xế và chuyến đi.
3. Phân trang `page`, `limit`.

### 4.8 Thanh toán *(V1 — Phase 13)*
1. Rider chọn phương thức: tiền mặt / ví / thẻ.
2. Sau `completed`, hệ thống tạo payment record.
3. Trừ voucher nếu có, tính số tiền cuối.
4. Phát hóa đơn điện tử.

### 4.9 Thông báo & liên lạc *(V1 — Phase 14)*
1. Gửi thông báo khi: có tài xế, tài xế đến, hoàn thành, hủy.
2. Chat in-app giữa rider và driver trong thời gian chuyến active.
3. Gọi ẩn số (masked phone) qua provider bên thứ 3.

### 4.10 Luồng lỗi chuẩn
- JWT không hợp lệ → 401
- Sai role truy cập → 403
- Ride không tồn tại → 404
- Ride đã có driver / đã hủy → 409
- Input thiếu/không hợp lệ → 400

---

## 5) Kiến trúc đề xuất

```text
cmd/api/main.go
internal/
  config/
  handler/
  service/
  repository/
  model/
  middleware/
  websocket/        # Phase 10
  notification/     # Phase 14
migrations/
docker-compose.yml
README.md
```

Nguyên tắc:
- `handler`: nhận request/response HTTP
- `service`: xử lý nghiệp vụ
- `repository`: truy vấn DB
- `model`: entity và DTO

---

## 6) Các phase công việc (break task)

> **Cách dùng trên Jira:** Mỗi Phase = 1 Epic. Mỗi bullet `- [ ]` = 1 Story/Task. Sub-task tách theo handler/service/repo nếu cần.

---

### Phase 0 — Setup
**Epic:** `PH0 - Setup môi trường`

- [ ] Cài Go 1.22+, Docker Desktop, Postman
- [ ] `go mod init` cho project
- [ ] Tạo cấu trúc thư mục `cmd`, `internal`, `migrations`
- [ ] Viết `docker-compose.yml` (PostgreSQL)
- [ ] Tạo `.env.example` + config loader
- [ ] Verify kết nối DB thành công

**Deliverable:** API skeleton chạy local, DB connect OK

---

### Phase 1 — API Foundation
**Epic:** `PH1 - Nền tảng API`

- [ ] Setup Gin/Echo HTTP server
- [ ] `GET /api/v1/health`
- [ ] Middleware: logger, recover, request timeout
- [ ] Chuẩn hóa response `{ "data": ..., "error": ... }`
- [ ] Global error handler

**Deliverable:** Skeleton API ổn định, cấu trúc rõ ràng

---

### Phase 2 — Database
**Epic:** `PH2 - Database schema`

- [ ] Migration bảng `users` (id, email, phone, password_hash, role, created_at)
- [ ] Migration bảng `drivers` (user_id, vehicle_type, is_online, lat, lng)
- [ ] Migration bảng `rides` (rider_id, driver_id, pickup/dropoff, status, fare, timestamps)
- [ ] Migration bảng `ride_ratings` (ride_id, rider_id, driver_id, score, comment)
- [ ] Repository layer + connection pool
- [ ] Seed data test (optional)

**Deliverable:** Schema sẵn sàng cho auth + ride + rating

---

### Phase 3 — Auth
**Epic:** `PH3 - Xác thực`

- [ ] `POST /api/v1/auth/register` (role: rider | driver)
- [ ] `POST /api/v1/auth/login` → JWT
- [ ] Hash password bcrypt
- [ ] Middleware `AuthRequired`
- [ ] Middleware `RequireRole("rider" | "driver")`
- [ ] Refresh token (optional)

**Deliverable:** Rider/Driver đăng ký, đăng nhập, gọi API bảo mật

---

### Phase 4 — Ride Core Flow
**Epic:** `PH4 - Luồng đặt chuyến`

- [ ] `POST /api/v1/rides` — rider tạo chuyến
- [ ] Fare estimator: base_fare + distance_km × price_per_km (Haversine)
- [ ] `GET /api/v1/rides/{id}`
- [ ] `GET /api/v1/rides` — danh sách theo role
- [ ] `GET /api/v1/rides/available` — driver xem pending
- [ ] `POST /api/v1/rides/{id}/accept`
- [ ] `PATCH /api/v1/rides/{id}/status` (accepted → in_progress → completed)
- [ ] Chặn conflict khi ride đã có driver

**Deliverable:** Full luồng đặt chuyến end-to-end (manual accept)

---

### Phase 5 — Driver Status & Location
**Epic:** `PH5 - Trạng thái tài xế`

- [ ] `PATCH /api/v1/drivers/me/online` (online/offline)
- [ ] `PATCH /api/v1/drivers/me/location` (lat, lng)
- [ ] Rule: chỉ driver online mới accept được chuyến
- [ ] Lưu `last_location_at` timestamp

**Deliverable:** Driver sẵn sàng nhận chuyến + cập nhật vị trí

---

### Phase 6 — Test & Hardening
**Epic:** `PH6 - Kiểm thử MVP`

- [ ] Unit test fare calculator (Haversine)
- [ ] Test auth middleware
- [ ] Integration test flow: register → ride → complete
- [ ] Postman collection đầy đủ
- [ ] Error codes nhất quán (400/401/403/404/409)

**Deliverable:** MVP ổn định, demo được

---

### Phase 7 — Docs & Runbook
**Epic:** `PH7 - Tài liệu`

- [ ] README: setup, env, chạy local
- [ ] Danh sách endpoint + ví dụ curl
- [ ] Mô tả scope đã làm / chưa làm
- [ ] Diagram luồng ride (optional)

**Deliverable:** Người mới clone và chạy được trong < 15 phút

---

### Phase 8 — Rating, History & Driver Stats
**Epic:** `PH8 - Đánh giá & lịch sử`

- [ ] `POST /api/v1/rides/{id}/rating` — rider đánh giá (1–5 + comment)
- [ ] Rule: chỉ ride `completed`, mỗi ride 1 lần, đúng rider
- [ ] `GET /api/v1/drivers/{id}/rating-summary` (avg_rating, total_ratings)
- [ ] `GET /api/v1/drivers/{id}/ratings` — danh sách đánh giá
- [ ] `GET /api/v1/drivers/me/stats` (completed_rides, avg_rating, total_earnings)
- [ ] `GET /api/v1/drivers/me/rides/history` — phân trang
- [ ] `GET /api/v1/riders/me/rides/history` — kèm info tài xế + chuyến đi

**Deliverable:** Đánh giá, lịch sử, thống kê tài xế đầy đủ

---

### Phase 9 — Hủy chuyến & Phí hủy
**Epic:** `PH9 - Cancel ride`

- [ ] Migration: `cancellation_reasons`, cột `cancelled_by`, `cancel_fee` trên `rides`
- [ ] `POST /api/v1/rides/{id}/cancel` — rider hoặc driver hủy
- [ ] Rule hủy theo status (pending/accepted/in_progress)
- [ ] Tính phí hủy theo thời điểm và vai trò người hủy
- [ ] `GET /api/v1/rides/{id}/cancellation` — chi tiết hủy
- [ ] Audit log mỗi lần đổi trạng thái ride
- [ ] Thông báo cho bên còn lại khi bị hủy *(liên kết Phase 14)*

**Deliverable:** Hủy chuyến có lý do, phí, audit

---

### Phase 10 — Real-time Tracking & ETA
**Epic:** `PH10 - Theo dõi real-time`

- [ ] Setup WebSocket hub (Gorilla WebSocket)
- [ ] Channel theo `ride_id`: rider subscribe vị trí driver
- [ ] Driver push location mỗi N giây qua WS hoặc REST
- [ ] `GET /api/v1/rides/{id}/eta` — ước tính thời gian đến
- [ ] Redis lưu vị trí driver tạm thời (TTL)
- [ ] Fallback polling nếu WS disconnect

**Deliverable:** Rider thấy vị trí tài xế cập nhật gần real-time

---

### Phase 11 — Auto-matching Driver
**Epic:** `PH11 - Gán tài xế tự động`

- [ ] PostGIS hoặc query Haversine: tìm driver online trong bán kính R km
- [ ] Job/worker gán ride `pending` → push tới driver gần nhất
- [ ] Timeout: không accept trong T giây → chuyển driver tiếp theo
- [ ] `POST /api/v1/rides/{id}/decline` — driver từ chối
- [ ] Ưu tiên: khoảng cách, rating, số chuyến hoàn thành
- [ ] Cấu hình `MATCH_RADIUS_KM`, `MATCH_TIMEOUT_SEC` trong env

**Deliverable:** Không cần driver tự browse; hệ thống tự match

---

### Phase 12 — Loại xe, Địa chỉ lưu & OTP
**Epic:** `PH12 - Đa dịch vụ xe & xác minh`

**12a — Nhiều loại xe**
- [ ] Migration `vehicle_types` (bike, car, car_plus…)
- [ ] Bảng `pricing_rules` theo loại xe
- [ ] Rider chọn `vehicle_type` khi tạo ride
- [ ] Fare tính theo loại xe

**12b — Địa chỉ lưu**
- [ ] Migration `saved_addresses` (user_id, label, address, lat, lng)
- [ ] `POST/GET/PUT/DELETE /api/v1/riders/me/addresses`
- [ ] `GET /api/v1/riders/me/addresses/recent` — điểm đến gần đây

**12c — OTP SĐT**
- [ ] `POST /api/v1/auth/send-otp` (phone)
- [ ] `POST /api/v1/auth/verify-otp`
- [ ] Tích hợp SMS provider (Twilio / mock dev)
- [ ] Chặn tạo ride nếu phone chưa verify

**Deliverable:** Đặt xe đa loại, địa chỉ nhanh, xác minh SĐT

---

### Phase 13 — Thanh toán & Hóa đơn
**Epic:** `PH13 - Payment`

- [ ] Migration `payments`, `invoices`
- [ ] Phương thức: `cash`, `wallet`, `card` (mock gateway trước)
- [ ] `POST /api/v1/rides/{id}/pay` — thanh toán sau completed
- [ ] `GET /api/v1/riders/me/payments` — lịch sử thanh toán
- [ ] `GET /api/v1/rides/{id}/invoice` — hóa đơn PDF/JSON
- [ ] Tích hợp payment gateway thật (Stripe/MoMo/VNPay) — optional
- [ ] Idempotency key chống thanh toán trùng

**Deliverable:** Thanh toán và hóa đơn sau mỗi chuyến

---

### Phase 14 — Thông báo & Chat
**Epic:** `PH14 - Notification & Communication`

**14a — Thông báo**
- [ ] Migration `notifications` (user_id, type, payload, read_at)
- [ ] `GET /api/v1/notifications`, `PATCH .../read`
- [ ] Push qua FCM/Web Push (web) hoặc mock queue
- [ ] Event: matched, driver_arrived, completed, cancelled

**14b — Chat**
- [ ] Migration `ride_messages` (ride_id, sender_id, content, created_at)
- [ ] `GET/POST /api/v1/rides/{id}/messages`
- [ ] Chỉ chat khi ride active (accepted/in_progress)
- [ ] WebSocket channel chat theo ride_id

**14c — Gọi ẩn số**
- [ ] `POST /api/v1/rides/{id}/call-token` — masked number (mock/provider)
- [ ] Chỉ trong thời gian chuyến active

**Deliverable:** Rider/driver nhận thông báo và liên lạc trong app

---

### Phase 15 — Voucher, Surge & Đặt trước
**Epic:** `PH15 - Pricing nâng cao`

**15a — Voucher**
- [ ] Migration `promo_codes`, `promo_usages`
- [ ] `POST /api/v1/rides/estimate` — báo giá trước khi đặt
- [ ] `POST /api/v1/promos/apply` — áp mã giảm giá
- [ ] Rule: hạn dùng, số lần, giảm % hoặc số tiền cố định

**15b — Surge pricing**
- [ ] Migration `surge_zones` hoặc config theo giờ
- [ ] Hệ số nhân giá theo khu vực/giờ cao điểm
- [ ] Hiển thị `surge_multiplier` trong estimate

**15c — Đặt trước (scheduled ride)**
- [ ] Cột `scheduled_at` trên `rides`
- [ ] `POST /api/v1/rides/scheduled`
- [ ] Cron/worker kích hoạt match trước `scheduled_at` X phút
- [ ] Hủy miễn phí trước deadline (rule)

**Deliverable:** Giá linh hoạt, khuyến mãi, đặt trước

---

### Phase 16 — KYC Tài xế & Admin
**Epic:** `PH16 - Onboarding & vận hành`

**16a — KYC tài xế**
- [ ] Migration `driver_documents` (license, vehicle_reg, status)
- [ ] `POST /api/v1/drivers/me/documents` — upload metadata/URL
- [ ] Trạng thái: `pending_review`, `approved`, `rejected`
- [ ] Chặn `online` nếu chưa `approved`

**16b — Admin dashboard API**
- [ ] Role `admin` + middleware
- [ ] `GET /api/v1/admin/rides` — filter status, date
- [ ] `GET /api/v1/admin/drivers` — filter KYC status
- [ ] `PATCH /api/v1/admin/drivers/{id}/kyc` — duyệt/từ chối
- [ ] `GET /api/v1/admin/stats` — tổng ride, doanh thu, driver active
- [ ] `GET /api/v1/admin/rides/{id}/audit` — lịch sử trạng thái

**Deliverable:** Tài xế được duyệt hồ sơ; admin quản lý hệ thống

---

### Phase 17 — Tính năng bổ sung (V2)
**Epic:** `PH17 - Nâng cao trải nghiệm`

- [ ] Đánh giá 2 chiều: tài xế đánh giá rider
- [ ] Tip sau chuyến: `POST /api/v1/rides/{id}/tip`
- [ ] Chia sẻ chuyến: `POST /api/v1/rides/{id}/share-link`
- [ ] SOS / khẩn cấp: `POST /api/v1/rides/{id}/sos`
- [ ] Routing thật: tích hợp Google Maps Directions / OSRM
- [ ] Đa điểm dừng: `stops[]` trên ride
- [ ] Chi tiết cước: base + km + phút + phụ phí + toll

**Deliverable:** Gần parity ride-hailing Grab (trừ super app)

---

### Phase 18 — Bảo mật, Scale & Production
**Epic:** `PH18 - Production ready`

- [ ] Rate limiting (Redis)
- [ ] Fraud rules: hủy ảo, spam đặt xe
- [ ] Driver incentive API (thưởng theo số chuyến)
- [ ] Support ticket: `POST /api/v1/support/tickets`
- [ ] Structured logging + correlation ID
- [ ] Health check + metrics (Prometheus)
- [ ] Deploy Docker + CI/CD
- [ ] Backup DB + migration strategy

**Deliverable:** Chạy production an toàn, có giám sát

---

## 7) Danh sách API tổng hợp

### Auth
| Method | Endpoint | Phase |
|--------|----------|-------|
| POST | `/api/v1/auth/register` | 3 |
| POST | `/api/v1/auth/login` | 3 |
| POST | `/api/v1/auth/send-otp` | 12 |
| POST | `/api/v1/auth/verify-otp` | 12 |

### Rider
| Method | Endpoint | Phase |
|--------|----------|-------|
| GET | `/api/v1/riders/me/rides/history` | 8 |
| GET/POST/PUT/DELETE | `/api/v1/riders/me/addresses` | 12 |
| GET | `/api/v1/riders/me/payments` | 13 |

### Driver
| Method | Endpoint | Phase |
|--------|----------|-------|
| PATCH | `/api/v1/drivers/me/online` | 5 |
| PATCH | `/api/v1/drivers/me/location` | 5 |
| GET | `/api/v1/drivers/me/stats` | 8 |
| GET | `/api/v1/drivers/me/rides/history` | 8 |
| GET | `/api/v1/drivers/{id}/rating-summary` | 8 |
| GET | `/api/v1/drivers/{id}/ratings` | 8 |
| POST | `/api/v1/drivers/me/documents` | 16 |

### Ride
| Method | Endpoint | Phase |
|--------|----------|-------|
| POST | `/api/v1/rides` | 4 |
| POST | `/api/v1/rides/scheduled` | 15 |
| POST | `/api/v1/rides/estimate` | 15 |
| GET | `/api/v1/rides/{id}` | 4 |
| GET | `/api/v1/rides` | 4 |
| GET | `/api/v1/rides/available` | 4 |
| POST | `/api/v1/rides/{id}/accept` | 4 |
| POST | `/api/v1/rides/{id}/decline` | 11 |
| PATCH | `/api/v1/rides/{id}/status` | 4 |
| POST | `/api/v1/rides/{id}/cancel` | 9 |
| POST | `/api/v1/rides/{id}/rating` | 8 |
| POST | `/api/v1/rides/{id}/pay` | 13 |
| GET | `/api/v1/rides/{id}/invoice` | 13 |
| GET | `/api/v1/rides/{id}/eta` | 10 |
| GET/POST | `/api/v1/rides/{id}/messages` | 14 |
| POST | `/api/v1/rides/{id}/tip` | 17 |

### Promo / Notification / Admin
| Method | Endpoint | Phase |
|--------|----------|-------|
| POST | `/api/v1/promos/apply` | 15 |
| GET/PATCH | `/api/v1/notifications` | 14 |
| GET | `/api/v1/admin/rides`, `/drivers`, `/stats` | 16 |
| PATCH | `/api/v1/admin/drivers/{id}/kyc` | 16 |

### System
| Method | Endpoint | Phase |
|--------|----------|-------|
| GET | `/api/v1/health` | 1 |
| WS | `/api/v1/ws/rides/{id}` | 10 |

---

## 8) Lộ trình sprint gợi ý (Jira)

| Sprint | Phase | Mục tiêu |
|--------|-------|----------|
| Sprint 1 | 0, 1, 2 | Setup + API + DB |
| Sprint 2 | 3, 4, 5 | Auth + Ride flow + Driver |
| Sprint 3 | 6, 7, 8 | Test + Docs + Rating/History |
| Sprint 4 | 9, 10 | Hủy chuyến + Real-time |
| Sprint 5 | 11, 12 | Auto-match + Loại xe/OTP |
| Sprint 6 | 13, 14 | Payment + Notification/Chat |
| Sprint 7 | 15, 16 | Voucher/Surge + KYC/Admin |
| Sprint 8 | 17, 18 | Nâng cao + Production |

**Label Jira đề xuất:** `phase-0` … `phase-18`, `backend`, `go`, `mvp`, `v1`, `v2`

---

## 9) Tiêu chí hoàn thành

### MVP (Phase 0–8)
- [ ] Chạy local bằng 1–2 lệnh
- [ ] Đăng ký/đăng nhập rider và driver
- [ ] Tạo ride → accept → complete
- [ ] Đánh giá tài xế + xem lịch sử + stats
- [ ] README đầy đủ

### V1 (Phase 9–16)
- [ ] Hủy chuyến có phí
- [ ] Real-time tracking + auto-match
- [ ] Thanh toán + thông báo
- [ ] Voucher/surge + KYC + admin cơ bản

### V2 (Phase 17–18)
- [ ] Tip, SOS, routing thật
- [ ] Rate limit, monitoring, deploy production

---

## 10) Ghi chú scope

- **Super app** (Food, Mart, Express, GrabPay đầy đủ): không nằm trong plan hiện tại.
- **Web frontend**: gọi REST/WS từ backend; UI là project riêng.
- Ưu tiên build **theo thứ tự phase** — không nhảy phase khi phase trước chưa xong deliverable.
