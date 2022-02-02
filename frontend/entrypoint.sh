#!/bin/sh

echo "Check that we have NEXT_PUBLIC_API_URL var"
test -n "$NEXT_PUBLIC_API_URL"

echo "Check that we have NEXT_PUBLIC_ALGOD_PROTOCOL var"
test -n "$NEXT_PUBLIC_ALGOD_PROTOCOL"

echo "Check that we have NEXT_PUBLIC_ALGOD_ADDR var"
test -n "$NEXT_PUBLIC_ALGOD_ADDR"

echo "Check that we have NEXT_PUBLIC_ALGOD_TOKEN var"
test -n "$NEXT_PUBLIC_ALGOD_TOKEN"

find /app/.next \( -type d -name .git -prune \) -o -type f -print0 | xargs -0 sed -i "s#APP_NEXT_PUBLIC_API_URL#$NEXT_PUBLIC_API_URL#g"
find /app/.next \( -type d -name .git -prune \) -o -type f -print0 | xargs -0 sed -i "s#APP_NEXT_PUBLIC_ALGOD_PROTOCOL#$NEXT_PUBLIC_ALGOD_PROTOCOL#g"
find /app/.next \( -type d -name .git -prune \) -o -type f -print0 | xargs -0 sed -i "s#APP_NEXT_PUBLIC_ALGOD_ADDR#$NEXT_PUBLIC_ALGOD_ADDR#g"
find /app/.next \( -type d -name .git -prune \) -o -type f -print0 | xargs -0 sed -i "s#APP_NEXT_PUBLIC_ALGOD_TOKEN#$NEXT_PUBLIC_ALGOD_TOKEN#g"

echo "Starting Nextjs"
exec "$@"
