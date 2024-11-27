.PHONY: dependencies run build package generate-mocks

dependencies:
	go mod tidy && go mod vendor

tests:
	go test -v ./...

build:
	GOOS=linux go build -o bootstrap main.go

package: build
	zip bootstrap.zip bootstrap

generate-mocks:
	mockgen -package=mocks stori-account-summary/services SendGridClient > mocks/sendGridClient.go && \
	mockgen -package=mocks stori-account-summary/services S3Client > mocks/s3Client.go && \
	mockgen -package=mocks stori-account-summary/db DynamoClient > mocks/dynamoClient.go