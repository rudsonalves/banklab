#!/usr/bin/env bash
set -euo pipefail

# Creates independent API and mobile environment files.
# Existing values are preserved and only missing API keys are appended.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

API_LEGACY_ENV="$REPO_ROOT/api/.env"
API_DEV_ENV="$REPO_ROOT/api/dev.env"
API_STAGING_ENV="$REPO_ROOT/api/staging.env"
API_PROD_ENV="$REPO_ROOT/api/prod.env"
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

  date +%s%N | shasum -a 256 | awk '{print $1}'
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

  date +%s%N | shasum -a 256 | awk '{print $1}'
}

extract_env_value() {
  local file="$1"
  local key="$2"

  if [[ ! -f "$file" ]]; then
    return 0
  fi

  awk -F'=' -v key="$key" '
    $1 == key {
      value = substr($0, index($0, "=") + 1)
      gsub(/^"|"$/, "", value)
      gsub(/^\047|\047$/, "", value)
      print value
      exit
    }
  ' "$file"
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

set_env_value() {
  local file="$1"
  local key="$2"
  local value="$3"
  local current
  local tmp

  current="$(extract_env_value "$file" "$key")"
  if [[ "$current" == "$value" ]]; then
    return 0
  fi

  tmp="$(mktemp)"
  awk -v key="$key" -v value="$value" '
    BEGIN { replaced = 0 }
    {
      if (!replaced && index($0, key "=") == 1) {
        print key "=" value
        replaced = 1
      } else {
        print
      }
    }
    END {
      if (!replaced) {
        print key "=" value
      }
    }
  ' "$file" > "$tmp"
  mv "$tmp" "$file"
  echo "Updated: $key in $file"
}

create_api_env() {
  local file="$1"
  local app_env="$2"
  local db_port="$3"
  local db_name="$4"
  local expose_debug_token="$5"
  local seed_file="${6:-}"
  local public_base_url="${7:-http://localhost:8080}"
  local api_published_host="${8:-127.0.0.1}"

  local app_token
  local jwt_secret
  local pepper
  local db_password

  app_token="$(extract_env_value "$file" "APP_TOKEN")"
  jwt_secret="$(extract_env_value "$file" "JWT_SECRET")"
  pepper="$(extract_env_value "$file" "TRANSACTION_PASSWORD_PEPPER")"
  db_password="$(extract_env_value "$file" "DB_PASSWORD")"

  if [[ -z "$app_token" && -n "$seed_file" ]]; then
    app_token="$(extract_env_value "$seed_file" "APP_TOKEN")"
  fi
  if [[ -z "$jwt_secret" && -n "$seed_file" ]]; then
    jwt_secret="$(extract_env_value "$seed_file" "JWT_SECRET")"
  fi
  if [[ -z "$pepper" && -n "$seed_file" ]]; then
    pepper="$(extract_env_value "$seed_file" "TRANSACTION_PASSWORD_PEPPER")"
  fi
  if [[ -z "$db_password" && -n "$seed_file" ]]; then
    db_password="$(extract_env_value "$seed_file" "DB_PASSWORD")"
  fi

  app_token="${app_token:-$(random_hex_64)}"
  jwt_secret="${jwt_secret:-$(random_hex_64)}"
  pepper="${pepper:-$(random_base64_32)}"
  db_password="${db_password:-$(random_hex_64)}"

  if [[ ! -f "$file" ]]; then
    touch "$file"
    echo "Created: $file"
  else
    echo "Exists: $file"
  fi

  append_missing_env_value "$file" "APP_ENV" "$app_env"
  append_missing_env_value "$file" "PUBLIC_BASE_URL" "$public_base_url"
  append_missing_env_value "$file" "EXPOSE_DEBUG_VERIFICATION_TOKEN" "$expose_debug_token"
  append_missing_env_value "$file" "APP_TOKEN" "$app_token"
  append_missing_env_value "$file" "JWT_SECRET" "$jwt_secret"
  append_missing_env_value "$file" "JWT_ACCESS_TOKEN_DURATION" "15m"
  append_missing_env_value "$file" "JWT_REFRESH_TOKEN_DURATION" "168h"
  append_missing_env_value "$file" "TRANSACTION_PASSWORD_PEPPER" "$pepper"
  append_missing_env_value "$file" "SERVER_HOST" "0.0.0.0"
  append_missing_env_value "$file" "SERVER_PORT" "8080"
  append_missing_env_value "$file" "API_PUBLISHED_HOST" "$api_published_host"
  append_missing_env_value "$file" "API_PUBLISHED_PORT" "8080"
  append_missing_env_value "$file" "DB_HOST" "postgres"
  append_missing_env_value "$file" "DB_PORT" "5432"
  append_missing_env_value "$file" "DB_PUBLISHED_PORT" "$db_port"
  append_missing_env_value "$file" "DB_NAME" "$db_name"
  append_missing_env_value "$file" "DB_USER" "postgres"
  append_missing_env_value "$file" "DB_PASSWORD" "$db_password"

  # These values describe the Docker network and must remain consistent.
  set_env_value "$file" "SERVER_HOST" "0.0.0.0"
  set_env_value "$file" "SERVER_PORT" "8080"
  set_env_value "$file" "API_PUBLISHED_HOST" "$api_published_host"
  set_env_value "$file" "API_PUBLISHED_PORT" "8080"
  set_env_value "$file" "DB_HOST" "postgres"
  set_env_value "$file" "DB_PORT" "5432"
  set_env_value "$file" "DB_PUBLISHED_PORT" "$db_port"
}

create_mobile_env() {
  local file="$1"
  local base_url="$2"
  local app_mode="$3"
  local api_env_file="$4"
  local app_token

  app_token="$(extract_env_value "$api_env_file" "APP_TOKEN")"

  if [[ ! -f "$file" ]]; then
    cat > "$file" <<EOF
BASE_URL=$base_url

CONNECT_TIMEOUT=30000
RECEIVE_TIMEOUT=30000

APP_MODE=$app_mode

APP_ACCESS_TOKEN=$app_token
EOF
    echo "Created: $file"
  else
    echo "Exists: $file"
  fi
}

create_api_env "$API_DEV_ENV" "dev" "5432" "bank" "true" "$API_LEGACY_ENV" "http://localhost:8080" "0.0.0.0"
create_api_env "$API_STAGING_ENV" "staging" "5433" "bank_staging" "true" "" "https://api.rralves.dev.br" "127.0.0.1"
create_api_env "$API_PROD_ENV" "production" "5434" "bank_production" "false" "" "https://api.rralves.dev.br" "127.0.0.1"

create_mobile_env "$MOBILE_DEV_ENV" "http://localhost:8080" "dev" "$API_DEV_ENV"
create_mobile_env "$MOBILE_STAGING_ENV" "https://api.rralves.dev.br" "staging" "$API_STAGING_ENV"
create_mobile_env "$MOBILE_PROD_ENV" "https://api.rralves.dev.br" "prod" "$API_PROD_ENV"
