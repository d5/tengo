vet:
	go vet ./...

lint:
	golint -set_exit_status ./...

test: vet lint
	go test -race -cover ./...

fmt:
	go fmt ./...

build:
	go build -o ./bin/tengo github.com/d5/tengo
