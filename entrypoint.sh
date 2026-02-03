#!/bin/sh
set -e

cd /app

if [ ! -f ./private.pem ] || [ ! -f ./public.pem ]; then
  echo "RSA keys not found, generating..."
  openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
  openssl rsa -pubout -in private.pem -out public.pem
fi

echo "PostgreSQL is up - executing migrations"
/app/migrate -mode=up

exec /app/bookmark_service


