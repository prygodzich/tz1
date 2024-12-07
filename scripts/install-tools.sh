mkdir -p $TOOLS_BIN

go install github.com/pressly/goose/v3/cmd/goose@v3.10.0
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
go install mvdan.cc/gofumpt@v0.6.0