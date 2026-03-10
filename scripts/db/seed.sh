#!/usr/bin/env bash
set -euo pipefail

MYSQL_DSN="${MYSQL_DSN:-root:123456@tcp(127.0.0.1:3306)/gpu_scheduler}"
SEED_FILE="${1:-database/seeds/dev_seed.sql}"

if ! command -v mysql >/dev/null 2>&1; then
  echo "mysql client not found"
  exit 1
fi

echo "Seeding database using: ${SEED_FILE}"
mysql --default-character-set=utf8mb4 "${MYSQL_DSN}" < "${SEED_FILE}"
echo "Seed completed."