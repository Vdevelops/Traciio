#!/bin/sh
# Entrypoint script untuk CRM Healthcare API
# Otomatis setup database indexes saat container start

set -e

echo "🚀 CRM Healthcare API - Starting up..."

# Tunggu PostgreSQL ready
echo "⏳ Waiting for PostgreSQL..."
until nc -z -v -w30 $DB_HOST ${DB_PORT:-5432} 2>/dev/null; do
  echo "   Waiting for database connection..."
  sleep 2
done
echo "✅ PostgreSQL is ready"

# Tunggu Redis ready jika enabled
if [ "$REDIS_ENABLED" = "true" ]; then
  echo "⏳ Waiting for Redis..."
  until nc -z -v -w30 $(echo $REDIS_SENTINEL_ADDRS | cut -d',' -f1 | cut -d':' -f1) 26379 2>/dev/null; do
    echo "   Waiting for Redis connection..."
    sleep 2
  done
  echo "✅ Redis is ready"
fi

# Jalankan database migrations terlebih dahulu
echo "🔄 Running database migrations..."
export PGPASSWORD=$DB_PASSWORD

# Tunggu dan verifikasi koneksi database bisa execute query
echo "   Verifying database connection..."
until psql -h $DB_HOST -p ${DB_PORT:-5432} -U $DB_USER -d $DB_NAME -c "SELECT 1;" > /dev/null 2>&1; do
  echo "   Waiting for database to accept queries..."
  sleep 2
done
echo "   ✅ Database connection verified"

# Setup database indexes (setelah migrations)
echo "📊 Setting up database indexes..."
INDEX_FILE="/app/database_indexes.sql"

if [ -f "$INDEX_FILE" ]; then
    echo "   Checking if tables exist..."
    # Tunggu sampai minimal 1 table ada (indikasi migrations sudah berjalan)
    TABLE_COUNT=$(psql -h $DB_HOST -p ${DB_PORT:-5432} -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE';" 2>/dev/null | xargs)
    
    if [ "$TABLE_COUNT" -gt 0 ] 2>/dev/null; then
        echo "   Found $TABLE_COUNT tables. Creating performance indexes..."
        psql -h $DB_HOST -p ${DB_PORT:-5432} -U $DB_USER -d $DB_NAME -f "$INDEX_FILE" 2>&1 | grep -E "(CREATE|already exists|ERROR)" || true
        echo "✅ Database indexes creation attempted"
    else
        echo "   ⚠️  No tables found. Indexes will be created by application."
        echo "   📋 Note: Run indexes manually later if needed:"
        echo "      psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f $INDEX_FILE"
    fi
else
    echo "   ⚠️  database_indexes.sql not found at $INDEX_FILE"
fi

# Jalankan aplikasi utama
echo "🎯 Starting API Server..."
exec ./server
