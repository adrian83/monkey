

fmt:
	echo "formatting code"
	gofmt -w .

test:
	echo "running tests" 
	go test ./... -cover

lint:
	echo "running linters" 
	golangci-lint run ./...

repl:
	echo "Running REPL"
	go run cmd/repl/main.go

qa: fmt test lint