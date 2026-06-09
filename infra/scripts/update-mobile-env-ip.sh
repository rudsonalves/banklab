#!/usr/bin/env bash
set -euo pipefail

# Updates BASE_URL in mobile/dev.env with the current host LAN IP.
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

update_env_file "$MOBILE_DIR/dev.env"
