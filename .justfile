lint_and_test:
    go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix .
    go fmt .
    go vet .
    golangci-lint run .
    go test .
