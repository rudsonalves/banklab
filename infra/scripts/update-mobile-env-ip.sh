#!/usr/bin/env bash
set -euo pipefail

# Updates BASE_URL in mobile .env files with the current host LAN IP.
# Usage:
#   ./infra/scripts/update-mobile-env-ip.sh
#   ./infra/scripts/update-mobile-env-ip.sh --ip 192.168.0.19 --port 8080

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
MOBILE_DIR="$REPO_ROOT/mobile"

PORT="8080"
CUSTOM_IP=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --ip)
      CUSTOM_IP="${2:-}"
      shift 2
      ;;
    --port)
      PORT="${2:-}"
      shift 2
      ;;
    *)
      echo "Unknown argument: $1" >&2
      exit 1
      ;;
  esac
done

get_default_interface() {
  route get default 2>/dev/null | awk '/interface:/{print $2; exit}'
}

get_lan_ip() {
  local iface
  iface="$(get_default_interface)"

  if [[ -n "$iface" ]]; then
    local ip
    ip="$(ipconfig getifaddr "$iface" 2>/dev/null || true)"
    if [[ -n "$ip" ]]; then
      echo "$ip"
      return 0
    fi
  fi

  # Fallback for uncommon network setups.
  ifconfig | awk '/inet / && $2 != "127.0.0.1" {print $2; exit}'
}

HOST_IP="$CUSTOM_IP"
if [[ -z "$HOST_IP" ]]; then
  HOST_IP="$(get_lan_ip)"
fi

if [[ -z "$HOST_IP" ]]; then
  echo "Could not detect host LAN IP. Use --ip <address>." >&2
  exit 1
fi

BASE_URL_VALUE="http://$HOST_IP:$PORT"

update_env_file() {
  local file="$1"

  if [[ ! -f "$file" ]]; then
    echo "Skipping missing file: $file"
    return 0
  fi

  local tmp
  tmp="$(mktemp)"

  if grep -q '^BASE_URL=' "$file"; then
    awk -v base_url="$BASE_URL_VALUE" '
      BEGIN {replaced = 0}
      {
        if (!replaced && $0 ~ /^BASE_URL=/) {
          print "BASE_URL=" base_url
          replaced = 1
        } else {
          print $0
        }
      }
    ' "$file" > "$tmp"
  else
    {
      echo "BASE_URL=$BASE_URL_VALUE"
      cat "$file"
    } > "$tmp"
  fi

  mv "$tmp" "$file"
  echo "Updated $file -> BASE_URL=$BASE_URL_VALUE"
}

update_postman_environment() {
  local postman_env_file="$REPO_ROOT/tools/postman/Environment.postman_environment.json"
  local app_token="$1"

  if [[ ! -f "$postman_env_file" ]]; then
    echo "Skipping missing Postman environment file: $postman_env_file"
    return 0
  fi

  if ! command -v jq &> /dev/null; then
    echo "Warning: jq is not installed. Skipping Postman environment update." >&2
    return 0
  fi

  local tmp
  tmp="$(mktemp)"

  jq --arg base_url "$BASE_URL_VALUE" \
     --arg app_token "$app_token" \
     '.values |= map(
       if .key == "base_url" then .value = $base_url
       elif .key == "app_token" then .value = $app_token
       else .
       end
     )' "$postman_env_file" > "$tmp"

  mv "$tmp" "$postman_env_file"
  echo "Updated $postman_env_file -> base_url=$BASE_URL_VALUE, app_token=$app_token"
}

update_env_file "$MOBILE_DIR/dev.env"

# Extract APP_ACCESS_TOKEN from dev.env
APP_ACCESS_TOKEN=$(grep '^APP_ACCESS_TOKEN=' "$MOBILE_DIR/dev.env" | cut -d'=' -f2)
if [[ -n "$APP_ACCESS_TOKEN" ]]; then
  update_postman_environment "$APP_ACCESS_TOKEN"
fi
update_env_file "$MOBILE_DIR/staging.env"
update_env_file "$MOBILE_DIR/prod.env"
