docker compose --profile test run --rm test sh -c "
    echo '=== GO ENVIRONMENT ===' && \
    go env | grep -E 'VERSION|MODULE|PROXY|CACHE' && \
    echo '' && \
    echo '=== TEST PACKAGES ===' && \
    go list -m && \
    go list ./... | head -10 && \
    echo '' && \
    echo '=== STARTING TESTS ===' && \
    echo 'Command: go test -v -cover -race -count=1 -timeout=30s -failfast ./...' && \
    echo '' && \
    go test -v -cover -race -count=1 -timeout=30s -failfast -json ./... 2>&1 | tee /tmp/test-output.json && \
    echo '' && \
    echo '=== TEST SUMMARY ===' && \
    go test -v -cover ./... 2>&1 | tail -30
" 2>&1 | tee test-output.log