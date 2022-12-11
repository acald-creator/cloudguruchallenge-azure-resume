BINARY=handler

build:
	GOOS=linux GOARCH=amd64 go build -o ${BINARY} handler.go