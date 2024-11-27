# Stori Account Summary

This repository contains the code to generate account summaries from transaction data stored in an S3 bucket and process them using AWS Lambda.

## Table of Contents

- [Dependencies](#dependencies)
- [Configuration](#configuration)
- [Environment Variables](#environment-variables)
- [Makefile Instructions](#makefile-instructions)
- [Usage](#usage)

## Project structure:
This project follows a layer arquitecture with the following folders:
- `data`: this foder contain an example of transactions
- `data`: The storage layer, this folder contains all the logic to connect to a DynamoDB table and to upload the transactions to it with its respective unit tests
- `mocks`: Mocks for unit testing. You can find more information about how to generate it in the section [Makefile Instructions](#makefile-instructions)
- `model`: This folder contains all the required go'structures to manage the report and the request in the files `report.go` and `send-email.go`
- `pkg`:  This folder contains clients abtractiosn for SendGrid and Report. This extra abstraction to the library allows us to break the direct connection with the libraries and facilitates testing.
- `services`: contain all the services of the project `account-summary.go` is in charge of receive the rows of the transaction files and generate the report.  `downloader.go` connects with `s3` bucket to donwload the file and `send-email.go` connects with SendGrid to send the email

## Dependencies
- Go 1.18+: This projects relies en Golang as a prgramming language. You can download Golang: https://go.dev/doc/install
- MockGen for testing: You can download it following the instructions: https://github.com/uber-go/mock
- A Sendgrid account. We rely in Sendgrid as a email provider. The SendGrid free layer allow you 100 messages/day

## Configuration

### AWS Lambda

### S3 Bucket

1. Create an S3 bucket in the AWS Management Console.
2. Upload your `transactions.csv` file to the S3 bucket. This file should be a file that contains all the transactions. You can see an example in ./data folder.

### DynamoDB
1. The free trier of AWS allow you a free usage up to 25GB. In `DynamoDB service -> Create new table -> Choose a name for your table`. We leave the `PARTITION_KEY` as ID and type name and the `SORT_KEY` as the date with type string.


### AWS Lambda
1. Create a new Lambda function in the AWS Management Console.
2. Set the runtime to Amazon Linux 2023
3. Choose the correct arquitecture of your computer of the computer where the buil will be created.
4. Configure the triggers: The proposal was to use s3 callbacks to know when a file was uploaded and process it. To configure this trigger you have to go to `Configuration -> Triggers`, go to Add Trigger and select S3 with the PUT events
5. Environment variables: This project require five environment variables `DYNAMODB_TABLE_NAME` `EMAIL_FROM` `EMAIL_TO` `SENDGRID_API_KEY` `SENDGRID_TEMPLATE_ID`. To add this variables to your lambda go to `Configuration -> Environment variables -> Edit`. You now can add all your variables. the detail of each varaible can be fourn in [Environment Variables](#environment-variables)
6. You now can test your function using `Test window -> Create a new event -> Select a template -> S3 PUT`. Change the file name for transaactions.csv or the file name you upload with the transactions.
7. Add `IAM` permission to your lambda to access to `Dynamo` and `S3`. In `IAM` service go to roles, AWS automatically will create a role for lambda, the template 
`lambda-name-id`. Add a two new policies to your lambda: `AmazonDynamoDBFullAccess` and `AmazonS3ReadOnlyAccess`


## Environment Variables

The following environment variables are required for the Lambda function:

- `DYNAMODB_TABLE_NAME`: The AWS region where your resources are located.
- `EMAIL_FROM`: The name of the S3 bucket containing the `transactions.csv` file.
- `EMAIL_TO`: The name of the DynamoDB table to store transaction data.
- `SENDGRID_API_KEY`: You need to authenticate to sendgrid to be able to send the email
- `SENDGRID_TEMPLATE_ID`: This is the template that you will use for your email.

## Makefile Instructions

The `Makefile` contains several useful commands for managing the project.

Install dependencies: This command will tidy up and vendor the Go modules.

```sh
make dependencies
```
Test: This command will run all the tests in the project.

Run Tests
```sh
make tests
```

Generate mocks
```sh
make generate-mocks
```

This command will automatically re-create all the mocks using `mockgen`

Generate the package 
```sh
make package
```
This command will build and compress you project, leave it ready to be upload to lambda

## Usage

With all required configuration, you can upload the zip to your lambda and upload the file to your s3 bucket and see the file arriving to your mail!

I created this demo video to leave you see the final result! 
[Video](https://drive.google.com/file/d/1fU88FAXrAuM7xzZNAAfuHX5z2xzJg1d4/view?usp=sharing)
