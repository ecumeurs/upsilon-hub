#!/bin/bash
# scripts/setup_prod.sh - Initialize production environment secrets

set -e

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="$ROOT_DIR/.env"
TEMPLATE_FILE="$ROOT_DIR/env.example"

echo "---------------------------------------"
echo "Initializing Production Environment..."
echo "---------------------------------------"

FORCE=false

# Simple argument parsing
while [[ "$#" -gt 0 ]]; do
    case $1 in
        -f|--force) FORCE=true ;;
        *) echo "[!] Unknown parameter passed: $1"; exit 1 ;;
    esac
    shift
done

if [ ! -f "$TEMPLATE_FILE" ]; then
    echo "[!] Error: env.example not found at root."
    exit 1
fi

if [ -f "$ENV_FILE" ] && [ "$FORCE" = false ]; then
    echo "[!] .env already exists. Skipping recreation to avoid overriding secrets."
    echo "[TIP] Use --force or -f to overwrite the existing .env file."
    exit 0
fi

if [ "$FORCE" = true ]; then
    echo "[!] Overwriting existing .env as requested..."
fi

echo "[+] Copying template to .env..."
cp "$TEMPLATE_FILE" "$ENV_FILE"

# Function to generate a secure random string (alphanumeric)
generate_secret() {
    base64 < /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c 32
}

# Function to generate a Laravel APP_KEY (32 random bytes, base64 encoded)
generate_app_key() {
    echo "base64:$(head -c 32 /dev/urandom | base64)"
}

echo "[+] Generating secure secrets..."

# Generate App Key
APP_KEY=$(generate_app_key)
# Use ^ to anchor the replacement to the start of the key name to avoid substring collision
sed -i "s|^APP_KEY=GENERATED_SECRET|APP_KEY=$APP_KEY|g" "$ENV_FILE"

# Generate Reverb IDs
REVERB_APP_ID=$(generate_secret)
REVERB_APP_KEY=$(generate_secret)
REVERB_APP_SECRET=$(generate_secret)

sed -i "s|^REVERB_APP_ID=GENERATED_SECRET|REVERB_APP_ID=$REVERB_APP_ID|g" "$ENV_FILE"
sed -i "s|^REVERB_APP_KEY=GENERATED_SECRET|REVERB_APP_KEY=$REVERB_APP_KEY|g" "$ENV_FILE"
sed -i "s|^REVERB_APP_SECRET=GENERATED_SECRET|REVERB_APP_SECRET=$REVERB_APP_SECRET|g" "$ENV_FILE"

echo "---------------------------------------"
echo "[OK] Production environment initialized."
echo "[!] Shared secrets generated and propagated to .env"
echo "---------------------------------------"
