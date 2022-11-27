default: test

test:
	cd packages/calculator && go test -v ./...

build:
	cd cmd/calculator && go build -o calculator