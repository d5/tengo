vet:
	go vet ./...

lint:
	golint -set_exit_status ./...

test: vet lint
	go test -race -cover ./...

fmt:
	go fmt ./...
