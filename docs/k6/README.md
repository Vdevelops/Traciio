# K6 Load Testing Suite — CRM Healthcare

Suite lengkap untuk performance testing API CRM Healthcare menggunakan [K6](https://k6.io/).

## 📁 Struktur

```
docs/k6/
├── run_all_tests.bat              # Batch runner (Windows)
├── README.md                      # Dokumentasi ini
├── lib/
│   ├── config.js                  # Shared config (URL, auth, helpers, endpoints)
│   ├── auth.js                    # Login helper
│   └── reporter.js                # HTML/JSON report generator
├── scripts/
│   ├── 01_smoke_test.js           # Quick validation (1 VU, 1 iter, 30 endpoints)
│   ├── 02_health_check.js         # /health & /ping throughput (up to 20 VUs)
│   ├── 03_auth_flow.js            # Auth lifecycle (login→profile→refresh→logout)
│   ├── 04_load_test.js            # Mixed workload (30 VUs, ~16 min, 18 groups)
│   ├── 05_stress_test.js          # Breaking point (100 VUs, ~20 min)
│   ├── 06_spike_test.js           # Sudden spike (5→80 in 10s)
│   ├── 07_soak_test.js            # Memory leak detection (15 VUs, 34 min)
│   └── cleanup_k6_data.sql        # Cleanup test data dari database
├── reports/                       # Output reports (auto-generated)
└── data/                          # Test data files (optional)
```

## ⚡ Quick Start

### 1. Install K6

```bash
# Windows (winget)
winget install grafana.k6

# Windows (chocolatey)
choco install k6

# macOS
brew install k6

# Linux
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg \
  --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D68
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | \
  sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update && sudo apt-get install k6
```

### 2. Jalankan Test

```bash
# Dari folder docs/k6/:

# Run semua test (total ~80 menit)
run_all_tests.bat all

# Run cepat (smoke + health + auth, ~6 menit)
run_all_tests.bat quick

# Run satu test spesifik
run_all_tests.bat smoke
run_all_tests.bat health
run_all_tests.bat auth
run_all_tests.bat load
run_all_tests.bat stress
run_all_tests.bat spike
run_all_tests.bat soak

# Run manual dengan K6 langsung
k6 run -e BASE_URL=https://api.gilabs.id scripts/01_smoke_test.js
```

### 3. Custom Target URL

```bash
# Default: https://api.gilabs.id (Go backend API)
# Override via environment variable:
k6 run -e BASE_URL=http://localhost:8080 scripts/01_smoke_test.js
```

## 📊 Test Descriptions

| #      | Test                | VUs           | Duration    | Purpose                                       |
| ------ | ------------------- | ------------- | ----------- | --------------------------------------------- |
| 01     | Smoke               | 1             | ~30s        | Validasi 30+ endpoint accessible              |
| 02     | Health Check        | 20            | ~3 min      | Throughput baseline /health & /ping           |
| 03     | Auth Flow           | 10            | ~2.5 min    | Login → profile → refresh → logout            |
| 04     | Load Test           | 30            | ~16 min     | Mixed workload (18 endpoint groups)           |
| 05     | Stress Test         | 100           | ~20 min     | Cari breaking point server                    |
| 06     | Spike Test          | 80            | ~5 min      | Sudden traffic 5→80 dalam 10 detik            |
| 07     | Soak Test           | 15            | ~34 min     | Deteksi memory leak / connection leak         |
| **08** | **High Load Test**  | **1000-5000** | **~35 min** | **Test untuk 1000-5000 concurrent users**     |
| **09** | **Monitoring Test** | **10**        | **~5 min**  | **Verify circuit breakers & runtime metrics** |

## 🔗 Endpoint Coverage

Semua endpoint berikut sudah di-cover oleh test scripts, diverifikasi terhadap route files di `apps/api/internal/api/routes/`:

### Core Business

| Endpoint                  | Route File               | Tests                            |
| ------------------------- | ------------------------ | -------------------------------- |
| `/health`, `/ping`        | `main.go`                | Smoke, Health, Spike             |
| `/api/v1/auth/*`          | `auth_routes.go`         | Smoke, Auth                      |
| `/api/v1/users/*`         | `user_routes.go`         | Smoke, Load                      |
| `/api/v1/leads/*`         | `lead_routes.go`         | Smoke, Load, Stress, Soak        |
| `/api/v1/accounts/*`      | `account_routes.go`      | Smoke, Load, Stress, Spike, Soak |
| `/api/v1/contacts/*`      | `contact_routes.go`      | Smoke, Load, Stress, Spike, Soak |
| `/api/v1/deals/*`         | `pipeline_routes.go`     | Smoke, Load, Stress, Spike, Soak |
| `/api/v1/pipelines/*`     | `pipeline_routes.go`     | Smoke, Load, Soak                |
| `/api/v1/visit-reports/*` | `visit_report_routes.go` | Smoke, Load, Stress, Spike, Soak |
| `/api/v1/activities/*`    | `activity_routes.go`     | Smoke, Load, Soak                |
| `/api/v1/tasks/*`         | `task_routes.go`         | Smoke, Load, Stress, Spike, Soak |
| `/api/v1/schedules/*`     | `schedule_routes.go`     | Smoke, Load, Soak                |
| `/api/v1/products/*`      | `product_routes.go`      | Smoke, Load                      |

### Dashboard & Reports

| Endpoint                            | Route File            | Tests                            |
| ----------------------------------- | --------------------- | -------------------------------- |
| `/api/v1/dashboard/overview`        | `dashboard_routes.go` | Smoke, Load, Stress, Spike, Soak |
| `/api/v1/dashboard/pipeline`        | `dashboard_routes.go` | Smoke                            |
| `/api/v1/reports/sales-performance` | `report_routes.go`    | Smoke, Load                      |
| `/api/v1/reports/pipeline`          | `report_routes.go`    | Smoke, Load                      |
| `/api/v1/reports/visit-reports`     | `report_routes.go`    | Smoke, Load                      |
| `/api/v1/reports/account-activity`  | `report_routes.go`    | Smoke, Load                      |

### Master Data & Others

| Endpoint                   | Route File                 | Tests             |
| -------------------------- | -------------------------- | ----------------- |
| `/api/v1/notifications/*`  | `notification_routes.go`   | Smoke, Load, Soak |
| `/api/v1/roles`            | `role_routes.go`           | Smoke             |
| `/api/v1/lead-sources`     | `lead_source_routes.go`    | Smoke             |
| `/api/v1/lead-statuses`    | `lead_status_routes.go`    | Smoke             |
| `/api/v1/sales-overview/*` | `sales_overview_routes.go` | Smoke             |

## 🔐 Test Credentials

Kredensial default dari `seeders/auth_seeder.go`:

| Role          | Email                      | Password   |
| ------------- | -------------------------- | ---------- |
| Admin         | `admin@example.com`        | `admin123` |
| Sales Manager | `salesmanager@example.com` | `admin123` |
| Sales         | `sales@example.com`        | `admin123` |

## 📈 Thresholds (Pass/Fail Criteria)

### Smoke Test

- Error rate: < 20%
- p95 < 5000ms

### Health Check

- p95 < 500ms
- p99 < 1000ms
- Error rate < 5%

### Auth Flow

- p95 < 1000ms
- Login failures < 10%

### Load Test

- p90 < 500ms
- p95 < 800ms
- p99 < 1500ms
- Error rate < 5%

### Stress Test

- p95 < 3000ms
- Error rate < 20% (lenient, karena tujuannya cari breaking point)

### Spike Test

- p95 < 5000ms
- Error rate < 40% (lenient)

### Soak Test

- p95 < 800ms
- p99 < 1500ms
- Error rate < 5%

## 📋 Reports

Reports otomatis di-generate ke folder `reports/`:

- **JSON**: `reports/<test_name>_<timestamp>.json` — raw data untuk analysis
- **Console**: Summary text langsung di terminal

Contoh membaca report:

```bash
# Lihat report terakhir
type reports\load_test_*.json | python -m json.tool

# Atau buka JSON di browser/editor
```

## 🔧 Konfigurasi

### Rate Limits (Server Production)

Server memiliki rate limiting:
| Type | Limit |
|------|-------|
| General | 100 req / 60s per IP |
| Public | 200 req / 60s per IP |
| High Volume | 3000 req / 300s per IP |
| Mutation | 300 req / 300s per IP |
| Login | 5 req / 900s per IP |
| Refresh | 10 req / hour per IP |

> **Catatan**: Untuk testing beban tinggi (stress/spike), rate limiting mungkin mempengaruhi hasil.
> Koordinasi dengan tim DevOps untuk sementara menaikkan limit atau whitelist IP tester.

### Connection Limits

- PostgreSQL: max_connections = 100
- PgBouncer: MAX_CLIENT_CONN = 1000, DEFAULT_POOL_SIZE = 100
- Go DB Pool: MaxOpenConns = 200, MaxIdleConns = 50
- Redis Pool: Size = 200

### Environment Variables

| Variable     | Default                 | Description             |
| ------------ | ----------------------- | ----------------------- |
| `BASE_URL`   | `https://api.gilabs.id` | Target API server       |
| `MAX_VUS`    | varies per test         | Maximum virtual users   |
| `SPIKE_VUS`  | `80`                    | Spike test peak VUs     |
| `NORMAL_VUS` | `5`                     | Spike test normal VUs   |
| `SOAK_VUS`   | `15`                    | Soak test sustained VUs |

## 🧹 Cleanup

Setelah testing, bersihkan data test dari database:

```sql
-- Jalankan query dari file:
-- scripts/cleanup_k6_data.sql

-- Atau manual:
DELETE FROM leads WHERE first_name = 'K6';
```

Semua K6 test scripts menggunakan `first_name = 'K6'` sebagai marker, sehingga cleanup mudah dilakukan.

## ⚠️ Notes

1. **Jangan run stress/spike test di production tanpa koordinasi** — bisa trigger rate limiting atau down
2. **Login rate limit**: 5 req per 15 menit per IP. Semua script menggunakan 1x login di `setup()` dan share token ke semua VUs
3. **Soak test** membutuhkan ~34 menit — pastikan koneksi internet stabil
4. **Reports** folder akan bertambah besar setelah banyak test runs — bersihkan secara berkala
5. **Token expiry**: Jika test berjalan sangat lama (>1 jam), token mungkin expired. Script tidak auto-refresh di dalam VU loop
6. **Dashboard endpoints** menggunakan `high_volume` rate limiting — throughput testing aman
7. **Lead creation** requires `email` (unique), `first_name`, dan `lead_source` — semua script sudah menghandle ini
