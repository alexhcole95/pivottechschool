default: test

test:
	cd packages/calculator && go test -v ./...

build:
	cd cmd/writing-tests && go build -o calculator