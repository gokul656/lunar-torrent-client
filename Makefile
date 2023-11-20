BINARY_NAME=luna

build:
	GOARCH=amd64 GOOS=darwin go build -o bin/${BINARY_NAME}-darwin .
	GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-linux .
	GOARCH=amd64 GOOS=windows go build -o bin/${BINARY_NAME}-windows .