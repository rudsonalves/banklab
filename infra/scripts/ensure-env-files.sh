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

TRANSACTION_PASSWORD_PEPPER="$(extract_env_value "$API_ENV" "TRANSACTION_PASSWORD_PEPPER")"
if [[ -z "$TRANSACTION_PASSWORD_PEPPER" ]]; then
  TRANSACTION_PASSWORD_PEPPER="$(random_base64_32)"
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

create_if_missing "$API_ENV" "APP_TOKEN=$APP_TOKEN
JWT_SECRET=$JWT_SECRET
TRANSACTION_PASSWORD_PEPPER=$TRANSACTION_PASSWORD_PEPPER"

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
