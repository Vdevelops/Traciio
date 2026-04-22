#!/bin/bash
# Script untuk create database indexes setelah Docker containers running
# Usage: ./run_indexes.sh

echo "========================================"
echo "Creating Database Indexes"
echo "========================================"
echo ""

cd "$(dirname "$0")"

echo "Waiting for API to be ready..."
until curl -s http://localhost:8080/health >/dev/null 2>&1; do
    echo "  Waiting..."
    sleep 2
done
echo "API is ready!"
echo ""

echo "Creating indexes..."
docker-compose exec -T postgres psql -U postgres -d crm_healthcare -f /docker-entrypoint-initdb.d/01_indexes.sql

echo ""
echo "========================================"
echo "Done! Indexes created."
echo "========================================"
