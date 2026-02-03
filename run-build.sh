#!/usr/bin/env bash
set -euo pipefail

RESULTS_FILE=".build-results.json"

# Run the build
OUTPUT=$(go build -o gitissuesum . 2>&1) && SUCCESS=true || SUCCESS=false

# Count warnings (Go compiler warnings are rare, but check stderr output)
WARNINGS=$(echo "$OUTPUT" | grep -ci "warning" || true)
ERRORS=0
if [ "$SUCCESS" = false ]; then
  ERRORS=1
fi

# Collect messages
MESSAGES="[]"
if [ -n "$OUTPUT" ]; then
  # Escape output for JSON
  ESCAPED=$(echo "$OUTPUT" | head -20 | python3 -c 'import sys,json; print(json.dumps(sys.stdin.read().splitlines()))')
  MESSAGES="$ESCAPED"
fi

# Write results
cat > "$RESULTS_FILE" <<EOF
{
  "success": $SUCCESS,
  "errors": $ERRORS,
  "warnings": $WARNINGS,
  "messages": $MESSAGES
}
EOF

if [ "$SUCCESS" = true ]; then
  echo "Build succeeded."
else
  echo "Build failed."
  echo "$OUTPUT"
  exit 1
fi
