dependencies:
	go mod tidy && go mod vendor

run:
	go run main.go

build:
	GOOS=linux go build -o bootstrap main.go

package: build
	zip bootstrap.zip bootstrap

mock:
	mockgen stori-account-summary/services/send-email SendGridClient > mocks/send-grid-client.go