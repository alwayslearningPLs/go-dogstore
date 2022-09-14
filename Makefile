run:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o . ./...
	./wso2-rest-api-example

