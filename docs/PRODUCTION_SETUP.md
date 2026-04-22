# Production Setup Guide

Panduan lengkap untuk setup production environment untuk API (Go) dan Web (Next.js).

## Daftar Isi

1. [Prerequisites](#prerequisites)
2. [Port Configuration](#port-configuration)
3. [Database Setup](#database-setup)
4. [API Production Setup](#api-production-setup)
5. [Web Production Setup](#web-production-setup)
6. [Docker Production Setup](#docker-production-setup)
7. [Environment Variables](#environment-variables)
8. [Security Checklist](#security-checklist)
9. [Monitoring & Maintenance](#monitoring--maintenance)

---

## Port Configuration

**PENTING**: Port configuration sudah di-setup untuk menghindari konflik dengan project lain.

| Service | Container Port | Host Port | Access | Notes |
|---------|---------------|-----------|--------|-------|
| **API** | 8081 | 8081 | http://localhost:8081 | Port 8081 (bukan 8080) untuk menghindari konflik |
| **Web** | 3000 | 3001 | http://localhost:3001 | Port 3001 (bukan 3000) untuk menghindari konflik |
| **Database** | 5432 | 15432 | Internal only | Host port 15432 untuk menghindari konflik, hanya accessible dari localhost |

**Catatan Penting**:
- Semua port hanya accessible dari `127.0.0.1` (localhost) untuk security
- Gunakan **Nginx reverse proxy** untuk public access dengan SSL
- Port database (5432) tidak di-expose ke host, hanya accessible dari container lain di network yang sama
- Jika port 8081 atau 3001 sudah digunakan, ubah di docker-compose files

**Deploy Order**:
1. **Deploy API terlebih dahulu** (akan membuat network dan database)
2. **Deploy Web** (akan menggunakan network yang sama)

---

## Prerequisites

### Server Requirements

- **OS**: Linux (Ubuntu 20.04+ recommended)
- **CPU**: Minimum 2 cores
- **RAM**: Minimum 4GB (8GB recommended)
- **Storage**: Minimum 20GB SSD
- **Network**: Public IP dengan domain name (optional tapi recommended)

### Software Requirements

- **Docker** & **Docker Compose** (untuk containerized deployment)
- **PostgreSQL 16+** (jika tidak menggunakan Docker)
- **Node.js 18+** (untuk Next.js build)
- **Go 1.25+** (jika build API secara manual)
- **Nginx** atau reverse proxy lainnya (untuk production web server)
- **SSL Certificate** (Let's Encrypt recommended)

---

## Database Setup

### Option 1: PostgreSQL dengan Docker (Recommended)

```bash
# Buat file docker-compose untuk production database
cd apps/api
cat > docker-compose.prod.yml << EOF
version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    container_name: crm-healthcare-db-prod
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    ports:
      - "127.0.0.1:5432:5432"  # Bind ke localhost saja untuk security
    volumes:
      - postgres_data_prod:/var/lib/postgresql/data
      - ./backups:/backups
    networks:
      - crm-network-prod
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5

networks:
  crm-network-prod:
    driver: bridge

volumes:
  postgres_data_prod:
EOF

# Start database
docker-compose -f docker-compose.prod.yml up -d postgres
```

### Option 2: PostgreSQL Manual Installation

```bash
# Install PostgreSQL
sudo apt update
sudo apt install postgresql-16 postgresql-contrib

# Create database dan user
sudo -u postgres psql

# Di dalam PostgreSQL shell:
CREATE DATABASE crm_healthcare;
CREATE USER crm_user WITH PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE crm_healthcare TO crm_user;
ALTER USER crm_user CREATEDB;
\q

# Setup SSL (recommended untuk production)
sudo nano /etc/postgresql/16/main/postgresql.conf
# Set: ssl = on

sudo nano /etc/postgresql/16/main/pg_hba.conf
# Update authentication method untuk production
```

### Database Migration & Seeding

```bash
# Setelah database running, jalankan migration dan seeders
cd apps/api

# Pastikan environment variables sudah di-set
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=crm_user
export DB_PASSWORD=your_secure_password
export DB_NAME=crm_healthcare
export DB_SSLMODE=require

# Run seeders (hanya sekali untuk initial setup)
go run seeders/seed_all.go
```

---

## API Production Setup

### Prerequisites

- Docker & Docker Compose sudah terinstall
- Go 1.25+ (untuk build, atau gunakan Docker build)
- PostgreSQL credentials sudah disiapkan

### Step-by-Step Deployment

#### Step 1: Setup Environment Variables

```bash
cd apps/api

# Buat file .env.production
cat > .env.production << EOF
# Server Configuration
PORT=8081
ENV=production

# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=crm_user
DB_PASSWORD=your_secure_password_here
DB_NAME=crm_healthcare
DB_SSLMODE=require

# JWT Configuration - GENERATE SECURE SECRET!
JWT_SECRET=$(openssl rand -base64 32)
JWT_ACCESS_TTL=24
JWT_REFRESH_TTL=7

# Cerebras AI Configuration (optional)
CEREBRAS_BASE_URL=https://api.cerebras.ai
CEREBRAS_API_KEY=your_cerebras_api_key
CEREBRAS_MODEL=llama-3.1-8b
EOF

# Set secure permissions
chmod 600 .env.production

# Load environment variables
export $(cat .env.production | grep -v '^#' | xargs)
```

**PENTING**: 
- **JWT_SECRET**: Generate secret yang kuat (min 32 karakter)
  ```bash
  openssl rand -base64 32
  ```
- **DB_PASSWORD**: Gunakan password yang kuat untuk production
- **DB_SSLMODE**: Wajib `require` untuk production

#### Step 2: Build Docker Image

```bash
# Pastikan berada di directory apps/api
cd apps/api

# Build Docker image
docker-compose -f docker-compose.production.yml build

# Build akan compile Go binary dan create image
# Build time: ~1-3 menit pertama kali
```

**Catatan**: 
- Build pertama kali akan download Go base image
- Go binary akan di-compile di dalam container
- Image size: ~20-30MB (alpine-based)

#### Step 3: Deploy API dan Database

```bash
# Start semua services (API + Database)
docker-compose -f docker-compose.production.yml up -d

# Check status semua containers
docker ps | grep -E "(api-prod|db-prod)"

# Check logs API
docker logs crm-healthcare-api-prod --tail 50

# Follow logs (real-time)
docker logs -f crm-healthcare-api-prod
```

**Expected Output**:
- Database container: `Up` dan `healthy`
- API container: `Up` dan `healthy` (setelah database ready)
- Logs menunjukkan: `Server running on port 8081`

#### Step 4: Run Database Migrations & Seeders

```bash
# Pastikan database sudah ready
docker exec crm-healthcare-db-prod pg_isready -U crm_user

# Run seeders (hanya sekali untuk initial setup)
cd apps/api
export $(cat .env.production | grep -v '^#' | xargs)
go run seeders/seed_all.go

# Atau jika menggunakan Docker:
docker exec -it crm-healthcare-api-prod /bin/sh
# Di dalam container:
# go run seeders/seed_all.go
```

#### Step 5: Verify Deployment

```bash
# Test API health endpoint
curl http://localhost:8081/health

# Expected response:
# {"status":"ok","message":"API is running"}

# Test API version endpoint
curl http://localhost:8081/api/v1/

# Check database connection dari API
docker logs crm-healthcare-api-prod | grep -i "database\|postgres"
```

**Expected Results**:
- Health endpoint: `200 OK` dengan response JSON
- API version endpoint: `200 OK`
- Database connection: Success (no errors in logs)

### Troubleshooting API Deployment

#### API tidak bisa connect ke database

```bash
# Check database status
docker ps | grep db-prod
docker logs crm-healthcare-db-prod

# Test database connection
docker exec -it crm-healthcare-db-prod psql -U crm_user -d crm_healthcare

# Check network
docker network inspect crm-network-prod

# Check API environment variables
docker exec crm-healthcare-api-prod env | grep DB_
```

**Common Solutions**:
1. Database belum ready: Tunggu beberapa detik, database perlu waktu untuk start
2. Wrong credentials: Check `.env.production` file
3. Network issue: Pastikan API dan database di network yang sama

#### Port 8081 sudah digunakan

```bash
# Check process menggunakan port 8081
lsof -i :8081
# atau
sudo netstat -tlnp | grep 8081

# Stop process atau ubah port di docker-compose.production.yml
```

#### API container restart terus

```bash
# Check logs untuk error
docker logs crm-healthcare-api-prod

# Common causes:
# 1. Database connection failed
# 2. Missing environment variables
# 3. Port conflict
```

### Dockerfile dan Docker Compose

File-file berikut sudah tersedia:

- **Dockerfile**: `apps/api/Dockerfile` (multi-stage build, optimized)
- **docker-compose.production.yml**: `apps/api/docker-compose.production.yml`

**Port Configuration**:
- **API Container Port**: 8081 (internal)
- **API Host Port**: 8081 (external, untuk menghindari konflik dengan project lain)
- **Database Port**: 5432 (internal), tidak di-expose ke host (hanya accessible dari container)
- **Binding**: `127.0.0.1:8081:8081` (hanya accessible dari localhost, gunakan Nginx untuk public access)

**Network**:
- Network `crm-network-prod` akan dibuat otomatis
- API dan Database akan terhubung di network yang sama

### Option 2: Manual Binary Deployment

#### 1. Build Binary

```bash
cd apps/api

# Build untuk Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o bin/server ./cmd/server/main.go

# Atau menggunakan Makefile (jika ada)
make build-linux
```

#### 2. Setup Systemd Service

```bash
sudo nano /etc/systemd/system/crm-api.service
```

```ini
[Unit]
Description=CRM Healthcare API
After=network.target postgresql.service
Requires=postgresql.service

[Service]
Type=simple
User=crm
Group=crm
WorkingDirectory=/opt/crm-healthcare/api
ExecStart=/opt/crm-healthcare/api/bin/server
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=crm-api

# Environment variables
Environment="PORT=8081"
Environment="ENV=production"
Environment="DB_HOST=localhost"
Environment="DB_PORT=5432"
Environment="DB_USER=crm_user"
Environment="DB_PASSWORD=your_secure_password"
Environment="DB_NAME=crm_healthcare"
Environment="DB_SSLMODE=require"
Environment="JWT_SECRET=your-super-secret-jwt-key-minimum-32-characters-long"
Environment="JWT_ACCESS_TTL=24"
Environment="JWT_REFRESH_TTL=7"

[Install]
WantedBy=multi-user.target
```

#### 3. Deploy dan Start Service

```bash
# Copy binary ke production server
scp bin/server user@production-server:/opt/crm-healthcare/api/bin/

# Create user (jika belum ada)
sudo useradd -r -s /bin/false crm

# Set permissions
sudo chown -R crm:crm /opt/crm-healthcare/api
sudo chmod +x /opt/crm-healthcare/api/bin/server

# Enable dan start service
sudo systemctl daemon-reload
sudo systemctl enable crm-api
sudo systemctl start crm-api

# Check status
sudo systemctl status crm-api
sudo journalctl -u crm-api -f
```

---

## Web Production Setup

### Prerequisites

- Docker & Docker Compose sudah terinstall
- Network `crm-network-prod` sudah dibuat (akan dibuat otomatis oleh API docker-compose)
- Environment variable `NEXT_PUBLIC_API_URL` sudah disiapkan

### Step-by-Step Deployment

#### Step 1: Setup Environment Variables

```bash
cd apps/web

# Copy environment example
cp env.example .env.production

# Edit environment variables
nano .env.production
```

**`.env.production` content:**
```env
# API Configuration - Ganti dengan URL API production Anda
NEXT_PUBLIC_API_URL=http://localhost:8081
# Atau jika sudah setup Nginx:
# NEXT_PUBLIC_API_URL=https://api.yourdomain.com

# Environment
NODE_ENV=production
```

**PENTING**: 
- `NEXT_PUBLIC_API_URL` harus sesuai dengan URL API yang akan digunakan
- Jika menggunakan Nginx reverse proxy, gunakan URL domain production
- Jika testing lokal, gunakan `http://localhost:8081` (port API production)

#### Step 2: Build Docker Image

```bash
# Pastikan berada di directory apps/web
cd apps/web

# Build Docker image
docker-compose -f docker-compose.production.yml build

# Build akan memakan waktu beberapa menit pertama kali
# Build selanjutnya akan lebih cepat karena cache
```

**Catatan**: 
- Build pertama kali akan download base image dan install dependencies
- Pastikan koneksi internet stabil
- Build time: ~2-5 menit (tergantung koneksi)

#### Step 3: Create Docker Network (Jika Belum Ada)

```bash
# Network akan dibuat otomatis oleh API docker-compose
# Tapi jika ingin membuat manual:
docker network create crm-network-prod

# Atau skip step ini jika akan deploy API terlebih dahulu
```

#### Step 4: Deploy Web Container

```bash
# Start web container
docker-compose -f docker-compose.production.yml up -d

# Check status
docker ps | grep web-prod

# Check logs
docker logs crm-healthcare-web-prod --tail 50

# Follow logs (real-time)
docker logs -f crm-healthcare-web-prod
```

#### Step 5: Verify Deployment

```bash
# Test web server
curl http://localhost:3001

# Atau buka di browser
# http://localhost:3001

# Check health
docker exec crm-healthcare-web-prod wget -q -O- http://localhost:3001
```

**Expected Output**: 
- Container status: `Up` dan `healthy`
- Logs menunjukkan: `✓ Ready in XXXms`
- Web accessible di `http://localhost:3001`

### Troubleshooting Web Deployment

#### Container tidak start

```bash
# Check logs untuk error
docker logs crm-healthcare-web-prod

# Common issues:
# 1. Port 3001 sudah digunakan
#    Solution: Stop process yang menggunakan port 3001
#    lsof -i :3001
#    kill <PID>

# 2. Network tidak ditemukan
#    Solution: Buat network terlebih dahulu
#    docker network create crm-network-prod

# 3. Build error
#    Solution: Rebuild dengan --no-cache
#    docker-compose -f docker-compose.production.yml build --no-cache
```

#### Web tidak bisa connect ke API

```bash
# Check environment variable
docker exec crm-healthcare-web-prod env | grep NEXT_PUBLIC_API_URL

# Test koneksi ke API dari web container
docker exec crm-healthcare-web-prod wget -q -O- http://crm-healthcare-api-prod:8081/health

# Pastikan API container sudah running
docker ps | grep api-prod
```

### Dockerfile dan Docker Compose

File-file berikut sudah tersedia dan tidak perlu dibuat manual:

- **Dockerfile**: `apps/web/Dockerfile` (sudah dioptimalkan dengan Next.js standalone output)
- **docker-compose.production.yml**: `apps/web/docker-compose.production.yml`

**Port Configuration**:
- **Container Port**: 3000 (internal)
- **Host Port**: 3001 (external, untuk menghindari konflik dengan project lain)
- **Binding**: `127.0.0.1:3001:3000` (hanya accessible dari localhost, gunakan Nginx untuk public access)

### 2. Option B: PM2 Deployment

```bash
# Install PM2
npm install -g pm2

# Create ecosystem file
cd apps/web
cat > ecosystem.config.js << EOF
module.exports = {
  apps: [{
    name: 'crm-web',
    script: 'node_modules/.bin/next',
    args: 'start',
    cwd: '/opt/crm-healthcare/web',
    instances: 2,
    exec_mode: 'cluster',
    env: {
      NODE_ENV: 'production',
      PORT: 3001,
      NEXT_PUBLIC_API_URL: 'https://api.yourdomain.com'
    },
    error_file: '/var/log/crm-web/error.log',
    out_file: '/var/log/crm-web/out.log',
    log_date_format: 'YYYY-MM-DD HH:mm:ss Z',
    merge_logs: true,
    autorestart: true,
    max_memory_restart: '1G'
  }]
}
EOF

# Start dengan PM2
pm2 start ecosystem.config.js

# Save PM2 configuration
pm2 save
pm2 startup
```

---

## Nginx Reverse Proxy Setup

### 1. Install Nginx

```bash
sudo apt update
sudo apt install nginx
```

### 2. Setup SSL dengan Let's Encrypt

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx

# Generate SSL certificate
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com
```

### 3. Configure Nginx

```bash
sudo nano /etc/nginx/sites-available/crm-healthcare
```

```nginx
# API Upstream - Port 8081 untuk menghindari konflik
upstream api_backend {
    server 127.0.0.1:8081;
}

# Web Upstream - Port 3001 untuk menghindari konflik
upstream web_backend {
    server 127.0.0.1:3001;
}

# API Server
server {
    listen 80;
    server_name api.yourdomain.com;
    
    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.yourdomain.com/privkey.pem;
    
    # SSL Configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Security Headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # CORS Headers (adjust sesuai kebutuhan)
    add_header Access-Control-Allow-Origin "https://yourdomain.com" always;
    add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS" always;
    add_header Access-Control-Allow-Headers "Authorization, Content-Type" always;

    # Logging
    access_log /var/log/nginx/api-access.log;
    error_log /var/log/nginx/api-error.log;

    # Proxy Settings
    location / {
        proxy_pass http://api_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Health check endpoint (optional, bisa di-restrict)
    location /health {
        proxy_pass http://api_backend;
        access_log off;
    }
}

# Web Server
server {
    listen 80;
    server_name yourdomain.com www.yourdomain.com;
    
    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com www.yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    
    # SSL Configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Security Headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # Logging
    access_log /var/log/nginx/web-access.log;
    error_log /var/log/nginx/web-error.log;

    # Gzip Compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/x-javascript application/xml+rss application/json application/javascript;

    # Static files caching
    location /_next/static {
        proxy_pass http://web_backend;
        proxy_cache_valid 200 60m;
        add_header Cache-Control "public, immutable";
    }

    # Proxy Settings
    location / {
        proxy_pass http://web_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

### 4. Enable Site dan Test

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/crm-healthcare /etc/nginx/sites-enabled/

# Test configuration
sudo nginx -t

# Reload Nginx
sudo systemctl reload nginx
```

---

## Environment Variables

### API Environment Variables

File: `apps/api/.env.production`

```env
# Server Configuration
PORT=8081  # Port 8081 untuk menghindari konflik dengan project lain
ENV=production

# Database Configuration
DB_HOST=postgres  # Nama service di docker-compose (untuk Docker) atau localhost (untuk manual)
DB_PORT=5432      # Container port (internal), host port adalah 15432
DB_USER=crm_user
DB_PASSWORD=your_secure_password_here
DB_NAME=crm_healthcare
DB_SSLMODE=disable  # 'disable' untuk Docker internal network, 'require' untuk external connection

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-minimum-32-characters-long-change-this-in-production
JWT_ACCESS_TTL=24  # hours
JWT_REFRESH_TTL=7  # days

# Cerebras AI Configuration (optional)
CEREBRAS_BASE_URL=https://api.cerebras.ai
CEREBRAS_API_KEY=your_cerebras_api_key
CEREBRAS_MODEL=llama-3.1-8b
```

**Port Configuration**:
- **API Port**: `8081` (bukan 8080 untuk menghindari konflik)
- **Database Port**: `5432` (internal, tidak di-expose ke host)

### Web Environment Variables

File: `apps/web/.env.production`

```env
# API Configuration
# Jika menggunakan Nginx reverse proxy:
NEXT_PUBLIC_API_URL=https://api.yourdomain.com
# Jika testing lokal tanpa Nginx:
# NEXT_PUBLIC_API_URL=http://localhost:8081

# Environment
NODE_ENV=production
```

**Port Configuration**:
- **Web Port**: `3001` (bukan 3000 untuk menghindari konflik)
- **API URL**: Sesuaikan dengan setup Anda (localhost:8081 atau domain production)

**PENTING**: 
- `NEXT_PUBLIC_*` variables akan di-expose ke client-side
- Jangan simpan sensitive data di `NEXT_PUBLIC_*` variables
- Gunakan server-side environment variables untuk sensitive data

---

## Security Checklist

### ✅ Pre-Deployment Security

- [ ] **JWT Secret**: Generate strong random secret (minimum 32 characters)
  ```bash
  openssl rand -base64 32
  ```

- [ ] **Database Password**: Gunakan strong password
  ```bash
  openssl rand -base64 24
  ```

- [ ] **Environment Files**: Set proper permissions
  ```bash
  chmod 600 .env.production
  ```

- [ ] **Database SSL**: Enable SSL connection (`DB_SSLMODE=require`)

- [ ] **Firewall**: Setup firewall rules
  ```bash
  sudo ufw allow 22/tcp    # SSH
  sudo ufw allow 80/tcp     # HTTP
  sudo ufw allow 443/tcp    # HTTPS
  sudo ufw enable
  ```

- [ ] **Docker Security**: 
  - Jangan expose database port ke public
  - Gunakan Docker secrets untuk sensitive data
  - Run containers dengan non-root user

- [ ] **SSL/TLS**: Setup SSL certificate dengan Let's Encrypt

- [ ] **CORS**: Configure CORS properly di API
  - Hanya allow domain production
  - Jangan gunakan wildcard `*` di production

- [ ] **Rate Limiting**: Implement rate limiting di API (recommended)

- [ ] **Backup Strategy**: Setup automated database backups

### ✅ Post-Deployment Security

- [ ] **Health Checks**: Monitor `/health` endpoint
- [ ] **Logs**: Setup log rotation dan monitoring
- [ ] **Updates**: Keep dependencies updated
- [ ] **Monitoring**: Setup monitoring tools (Prometheus, Grafana, etc.)

---

## Monitoring & Maintenance

### Database Backups

```bash
# Create backup script
cat > /opt/crm-healthcare/scripts/backup-db.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/opt/crm-healthcare/backups"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="crm_healthcare"
DB_USER="crm_user"

mkdir -p $BACKUP_DIR

# Backup dengan pg_dump
docker exec crm-healthcare-db-prod pg_dump -U $DB_USER $DB_NAME > $BACKUP_DIR/backup_$DATE.sql

# Compress backup
gzip $BACKUP_DIR/backup_$DATE.sql

# Keep only last 30 days
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +30 -delete

echo "Backup completed: backup_$DATE.sql.gz"
EOF

chmod +x /opt/crm-healthcare/scripts/backup-db.sh

# Setup cron job (daily at 2 AM)
echo "0 2 * * * /opt/crm-healthcare/scripts/backup-db.sh" | sudo crontab -
```

### Log Rotation

```bash
sudo nano /etc/logrotate.d/crm-healthcare
```

```
/var/log/crm-api/*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 crm crm
    sharedscripts
    postrotate
        systemctl reload crm-api > /dev/null 2>&1 || true
    endscript
}
```

### Health Monitoring

```bash
# Create health check script
cat > /opt/crm-healthcare/scripts/health-check.sh << 'EOF'
#!/bin/bash

API_URL="https://api.yourdomain.com/health"
WEB_URL="https://yourdomain.com"

# Check API
if curl -f -s $API_URL > /dev/null; then
    echo "API: OK"
else
    echo "API: FAILED"
    # Send alert (email, Slack, etc.)
fi

# Check Web
if curl -f -s $WEB_URL > /dev/null; then
    echo "Web: OK"
else
    echo "Web: FAILED"
    # Send alert
fi
EOF

chmod +x /opt/crm-healthcare/scripts/health-check.sh

# Setup cron (every 5 minutes)
echo "*/5 * * * * /opt/crm-healthcare/scripts/health-check.sh" | sudo crontab -
```

---

## Troubleshooting

### API tidak bisa connect ke database

```bash
# Check database connection
docker exec -it crm-healthcare-db-prod psql -U crm_user -d crm_healthcare

# Check API logs
docker logs crm-healthcare-api-prod

# Check network
docker network inspect crm-network-prod
```

### Web tidak bisa connect ke API

```bash
# Check environment variable
docker exec crm-healthcare-web-prod env | grep NEXT_PUBLIC_API_URL

# Test API dari web container
docker exec crm-healthcare-web-prod wget -O- http://crm-healthcare-api-prod:8081/health
```

### SSL Certificate Issues

```bash
# Renew certificate
sudo certbot renew

# Test renewal
sudo certbot renew --dry-run
```

---

## Quick Start Commands

### Deploy Order (PENTING!)

**Deploy API terlebih dahulu**, karena:
1. API akan membuat network `crm-network-prod`
2. API akan membuat database container
3. Web membutuhkan network yang sama untuk connect ke API

### Start All Services (Docker)

```bash
# Step 1: Deploy API (termasuk Database)
cd apps/api
docker-compose -f docker-compose.production.yml up -d

# Wait untuk database ready (sekitar 10-15 detik)
sleep 15

# Step 2: Deploy Web
cd ../web
docker-compose -f docker-compose.production.yml up -d

# Step 3: Check status semua services
docker ps | grep -E "(api-prod|web-prod|db-prod)"

# Step 4: Verify services
curl http://localhost:8081/health  # API
curl http://localhost:3001         # Web
```

### Port Summary

| Service | Container Port | Host Port | URL |
|---------|---------------|-----------|-----|
| API | 8081 | 8081 | http://localhost:8081 |
| Web | 3000 | 3001 | http://localhost:3001 |
| Database | 5432 | 15432 | Internal only (accessible dari localhost:15432 untuk maintenance) |

**Catatan**: 
- Port di-host (8081, 3001) dipilih untuk menghindari konflik dengan project lain
- Database port tidak di-expose ke host untuk security
- Semua port hanya accessible dari localhost (127.0.0.1), gunakan Nginx untuk public access

### Stop All Services

```bash
# API
cd apps/api
docker-compose -f docker-compose.production.yml down

# Web
cd apps/web
docker-compose -f docker-compose.production.yml down
```

### View Logs

```bash
# API logs
docker logs -f crm-healthcare-api-prod

# Web logs
docker logs -f crm-healthcare-web-prod

# Database logs
docker logs -f crm-healthcare-db-prod
```

---

## Next Steps

1. **Setup CI/CD**: Configure automated deployment
2. **Monitoring**: Setup monitoring tools (Prometheus, Grafana)
3. **Alerting**: Configure alerts untuk critical issues
4. **Performance**: Optimize database queries dan API responses
5. **Scaling**: Setup load balancing jika diperlukan

---

## Support

Jika ada masalah dengan setup production, silakan:
1. Check logs terlebih dahulu
2. Review dokumentasi ini
3. Check GitHub issues
4. Contact team lead

