#!/usr/bin/env bash
# Mechanical layering migrate for one service: httpapi→app, enrich api, FORCE_REGEN, gen invoke logics.
# Usage: ./scripts/migrate-service-layering.sh <service-name>
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SVC="${1:?service name}"
DIR="$ROOT/services/$SVC"
export PATH="$(go env GOPATH)/bin:$PATH"

echo "==> migrate $SVC"

# 1) httpapi → app
if [[ -d "$DIR/internal/httpapi" ]]; then
  rm -rf "$DIR/internal/app"
  cp -a "$DIR/internal/httpapi" "$DIR/internal/app"
  find "$DIR/internal/app" -name '*.go' -print0 | xargs -0 sed -i '' 's|internal/httpapi/|internal/app/|g'
  find "$DIR/internal/app" -name '*.go' -print0 | xargs -0 perl -i -pe 's/biz\.New(\w+)\(context\.Background\(\),\s*/biz.New$1(/g'
  find "$DIR/internal/app" -name '*.go' -print0 | xargs -0 perl -i -pe 's/biz\.New(\w+)\(r\.Context\(\),\s*/biz.New$1(/g'
  rm -rf "$DIR/internal/httpapi"
  echo "httpapi -> app"
fi

# catalog nested httpapi
for nest in product content shopops notify; do
  if [[ -d "$DIR/internal/$nest/httpapi" ]]; then
    rm -rf "$DIR/internal/$nest/app"
    cp -a "$DIR/internal/$nest/httpapi" "$DIR/internal/$nest/app"
    find "$DIR/internal/$nest/app" -name '*.go' -print0 | xargs -0 sed -i '' "s|internal/$nest/httpapi/|internal/$nest/app/|g"
    find "$DIR/internal/$nest/app" -name '*.go' -print0 | xargs -0 perl -i -pe 's/\.New(\w+)\(context\.Background\(\),\s*/.New$1(/g'
    find "$DIR/internal/$nest/app" -name '*.go' -print0 | xargs -0 perl -i -pe 's/\.New(\w+)\(r\.Context\(\),\s*/.New$1(/g'
    rm -rf "$DIR/internal/$nest/httpapi"
    echo "$nest/httpapi -> app"
  fi
done

# 2) enrich api if needed
API="$(find "$DIR/api" -name '*.api' | head -1)"
python3 "$ROOT/scripts/enrich-api-types-generic.py" "$API"

# 3) FORCE_REGEN
FORCE_REGEN=1 "$ROOT/scripts/gen-api.sh" "$SVC"

# 4) generate invoke logics from thin-wrapper map discovered in app + routes
python3 "$ROOT/scripts/gen-invoke-logics.py" "$SVC"

echo "==> build $SVC"
go build -o /dev/null "./services/$SVC/..."
"$ROOT/scripts/check-api-routes.sh" "$SVC"
echo "OK $SVC"
