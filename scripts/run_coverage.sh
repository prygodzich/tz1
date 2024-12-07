set -e
set -a

mkdir -p coverage
go test -cover ./...  -coverprofile=coverage/coverage.txt 
go tool cover -html=coverage/coverage.txt -o coverage/coverage.html


totalcover=$(go tool cover -func=coverage/coverage.txt | grep total: | awk '{print $3}')

echo "Total coverage: $totalcover"