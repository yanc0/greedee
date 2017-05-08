default:
	make setup
	make build
	make test

setup:
	go get -u github.com/golang/dep
	go install github.com/golang/dep/cmd/dep

build:
	dep ensure
	go build -i

clean:
	go clean

fmt:
	go fmt ./...

test:
	go test

beforecommit:
	make fmt
	make test
	make clean

