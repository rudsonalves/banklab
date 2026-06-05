#!/usr/bin/env bash
set -euo pipefail

# Creates API and Mobile env files only when they do not exist.
# Existing files are never overwritten.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

API_ENV="$REPO_ROOT/api/.env"
MOBILE_DEV_ENV="$REPO_ROOT/mobile/dev.env"
MOBILE_STAGING_ENV="$REPO_ROOT/mobile/staging.env"
MOBILE_PROD_ENV="$REPO_ROOT/mobile/prod.env"

random_hex_64() {
  if command -v openssl >/dev/null 2>&1; then
    openssl rand -hex 32
    return 0
  fi

  if command -v xxd >/dev/null 2>&1; then
    xxd -l 32 -p /dev/urandom
    return 0
  fi

  # Last-resort fallback, deterministic but still unique enough for local dev.
  date +%s | shasum -a 256 | awk '{print $1}'
}

random_base64_32() {
  if command -v openssl >/dev/null 2>&1; then
    openssl rand -base64 32
    return 0
  fi

  if command -v xxd >/dev/null 2>&1; then
    xxd -l 32 -p /dev/urandom | xxd -r -p | base64
    return 0
  fi

  # Last-resort fallback for local dev only.
  date +%s | shasum -a 256 | awk '{print $1}'
}

extract_env_value() {
  local file="$1"
  local key="$2"

  if [[ ! -f "$file" ]]; then
    return 0
  fi

  awk -F'=' -v key="$key" '
    $1 == key {
      value = $2
      gsub(/^"|"$/, "", value)
      gsub(/^\047|\047$/, "", value)
      print value
      exit
    }
  ' "$file"
}

APP_TOKEN="$(extract_env_value "$API_ENV" "APP_TOKEN")"
if [[ -z "$APP_TOKEN" ]]; then
  APP_TOKEN="$(random_hex_64)"
fi

JWT_SECRET="$(extract_env_value "$API_ENV" "JWT_SECRET")"
if [[ -z "$JWT_SECRET" ]]; then
  JWT_SECRET="$(random_hex_64)"
fi

JWT_ACCESS_TOKEN_DURATION="$(extract_env_value "$API_ENV" "JWT_ACCESS_TOKEN_DURATION")"
if [[ -z "$JWT_ACCESS_TOKEN_DURATION" ]]; then
  JWT_ACCESS_TOKEN_DURATION="15m"
fi

JWT_REFRESH_TOKEN_DURATION="$(extract_env_value "$API_ENV" "JWT_REFRESH_TOKEN_DURATION")"
if [[ -z "$JWT_REFRESH_TOKEN_DURATION" ]]; then
  JWT_REFRESH_TOKEN_DURATION="168h"
fi

TRANSACTION_PASSWORD_PEPPER="$(extract_env_value "$API_ENV" "TRANSACTION_PASSWORD_PEPPER")"
if [[ -z "$TRANSACTION_PASSWORD_PEPPER" ]]; then
  TRANSACTION_PASSWORD_PEPPER="$(random_base64_32)"
fi

DB_HOST="$(extract_env_value "$API_ENV" "DB_HOST")"
if [[ -z "$DB_HOST" ]]; then
  DB_HOST="localhost"
fi

DB_PORT="$(extract_env_value "$API_ENV" "DB_PORT")"
if [[ -z "$DB_PORT" ]]; then
  DB_PORT="5432"
fi

DB_NAME="$(extract_env_value "$API_ENV" "DB_NAME")"
if [[ -z "$DB_NAME" ]]; then
  DB_NAME="bank"
fi

DB_USER="$(extract_env_value "$API_ENV" "DB_USER")"
if [[ -z "$DB_USER" ]]; then
  DB_USER="postgres"
fi

DB_PASSWORD="$(extract_env_value "$API_ENV" "DB_PASSWORD")"
if [[ -z "$DB_PASSWORD" ]]; then
  DB_PASSWORD="postgres"
fi

create_if_missing() {
  local file="$1"
  local content="$2"

  if [[ -f "$file" ]]; then
    echo "Exists: $file"
    return 0
  fi

  cat > "$file" <<EOF
$content
EOF
  echo "Created: $file"
}

append_missing_env_value() {
  local file="$1"
  local key="$2"
  local value="$3"

  if grep -q "^$key=" "$file"; then
    return 0
  fi

  printf '%s=%s\n' "$key" "$value" >> "$file"
  echo "Added: $key to $file"
}

create_if_missing "$API_ENV" "APP_TOKEN=$APP_TOKEN
JWT_SECRET=$JWT_SECRET
JWT_ACCESS_TOKEN_DURATION=$JWT_ACCESS_TOKEN_DURATION
JWT_REFRESH_TOKEN_DURATION=$JWT_REFRESH_TOKEN_DURATION
TRANSACTION_PASSWORD_PEPPER=$TRANSACTION_PASSWORD_PEPPER
SERVER_PORT=8080
DB_HOST=$DB_HOST
DB_PORT=$DB_PORT
DB_NAME=$DB_NAME
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD"

append_missing_env_value "$API_ENV" "APP_TOKEN" "$APP_TOKEN"
append_missing_env_value "$API_ENV" "JWT_SECRET" "$JWT_SECRET"
append_missing_env_value "$API_ENV" "JWT_ACCESS_TOKEN_DURATION" "$JWT_ACCESS_TOKEN_DURATION"
append_missing_env_value "$API_ENV" "JWT_REFRESH_TOKEN_DURATION" "$JWT_REFRESH_TOKEN_DURATION"
append_missing_env_value "$API_ENV" "TRANSACTION_PASSWORD_PEPPER" "$TRANSACTION_PASSWORD_PEPPER"
append_missing_env_value "$API_ENV" "SERVER_PORT" "8080"
append_missing_env_value "$API_ENV" "DB_HOST" "$DB_HOST"
append_missing_env_value "$API_ENV" "DB_PORT" "$DB_PORT"
append_missing_env_value "$API_ENV" "DB_NAME" "$DB_NAME"
append_missing_env_value "$API_ENV" "DB_USER" "$DB_USER"
append_missing_env_value "$API_ENV" "DB_PASSWORD" "$DB_PASSWORD"

create_if_missing "$MOBILE_DEV_ENV" "BASE_URL=http://localhost:8080

CONNECT_TIMEOUT=30000
RECEIVE_TIMEOUT=30000

APP_MODE=dev

APP_ACCESS_TOKEN=$APP_TOKEN"

create_if_missing "$MOBILE_STAGING_ENV" "BASE_URL=http://localhost:8080

CONNECT_TIMEOUT=30000
RECEIVE_TIMEOUT=30000

APP_MODE=staging

APP_ACCESS_TOKEN=$APP_TOKEN"

create_if_missing "$MOBILE_PROD_ENV" "BASE_URL=http://localhost:8080

CONNECT_TIMEOUT=30000
RECEIVE_TIMEOUT=30000

APP_MODE=prod

APP_ACCESS_TOKEN=$APP_TOKEN"
