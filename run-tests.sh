#!/usr/bin/env bash
set -uo pipefail

RESULTS_FILE=".test-results.json"

# Run tests with JSON output
OUTPUT=$(go test ./... -json 2>&1) && TEST_EXIT=0 || TEST_EXIT=$?

# Parse test results
PASSED=0
FAILED=0
FAILURES="[]"

if echo "$OUTPUT" | head -1 | python3 -c 'import sys,json; json.loads(sys.stdin.readline())' 2>/dev/null; then
  # JSON output available — parse it
  PASSED=$(echo "$OUTPUT" | python3 -c '
import sys, json
passed = 0
for line in sys.stdin:
    line = line.strip()
    if not line:
        continue
    try:
        obj = json.loads(line)
        if obj.get("Action") == "pass" and "Test" in obj:
            passed += 1
    except json.JSONDecodeError:
        pass
print(passed)
')
  FAILED_LIST=$(echo "$OUTPUT" | python3 -c '
import sys, json
failed = []
for line in sys.stdin:
    line = line.strip()
    if not line:
        continue
    try:
        obj = json.loads(line)
        if obj.get("Action") == "fail" and "Test" in obj:
            failed.append(obj["Test"])
    except json.JSONDecodeError:
        pass
print(json.dumps(failed))
')
  FAILED=$(echo "$FAILED_LIST" | python3 -c 'import sys,json; print(len(json.loads(sys.stdin.read())))')
  FAILURES="$FAILED_LIST"
else
  # No JSON output (maybe no test files yet)
  if [ $TEST_EXIT -eq 0 ]; then
    PASSED=0
    FAILED=0
  else
    FAILED=1
  fi
fi

TOTAL=$((PASSED + FAILED))

cat > "$RESULTS_FILE" <<EOF
{
  "passed": $PASSED,
  "failed": $FAILED,
  "total": $TOTAL,
  "failures": $FAILURES
}
EOF

if [ $FAILED -eq 0 ]; then
  echo "All $PASSED tests passed."
else
  echo "$FAILED of $TOTAL tests failed."
  exit 1
fi
