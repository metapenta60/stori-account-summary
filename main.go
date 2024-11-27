package main

import (
	"context"
	"log"
	"os"

	"stori-account-summary/db"
	"stori-account-summary/model"
	"stori-account-summary/pkg"
	"stori-account-summary/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const host = "https://api.sendgrid.com"

func handler(ctx context.Context, s3Event events.S3Event) error {
	client, err := initS3Client(ctx)
	if err != nil {
		log.Printf("Failed to create S3 client, %v", err)

		return nil
	}

	downloader := services.NewDownloader(ctx, s3Event, client)
	rows, err := downloader.DownloadFile(ctx)

	if err != nil {
		log.Printf("Failed to download file, %v", err)
		return err
	}

	db, err := initDB(ctx)
	if err != nil {
		log.Printf("Failed to create DynamoDB client, %v", err)
		return err
	}

	db.AddTransactions(rows)

	reportService := services.NewReportService(rows)
	report := reportService.AnalyseAccount()

	emailSender := initEmailSender()
	err = emailSender.SendMail(getEmailConfiguration(report))
	if err != nil {
		log.Printf("Failed to send email, %v", err)
		return err
	}

	return nil
}

func initS3Client(ctx context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Printf("Failed to load configuration, %v", err)
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}

func initEmailSender() services.EmailSender {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	emailFrom := os.Getenv("EMAIL_FROM")

	client := pkg.NewSendGridClient(host, apiKey)
	emailSender := services.NewEmailSender(client, emailFrom)

	return emailSender
}

func getEmailConfiguration(report model.AccountReport) model.EmailConfig {
	templateID := os.Getenv("SENDGRID_TEMPLATE_ID")
	emailTo := os.Getenv("EMAIL_TO")

	return model.EmailConfig{
		TemplateID: templateID,
		To:         emailTo,
		Report:     report,
	}
}

func initDB(ctx context.Context) (db.Db, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if err != nil {
		log.Printf("Failed to load configuration, %v", err)
		return db.Db{}, err
	}

	client := dynamodb.NewFromConfig(cfg)

	return db.New(ctx, tableName, client), nil
}

func main() {
	lambda.Start(handler)
}
