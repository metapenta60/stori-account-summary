package main

import (
	"context"
	"log"
	"os"

	"stori-account-summary/model"
	"stori-account-summary/pkg"
	"stori-account-summary/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
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

	reportService := services.NewReportService(rows)
	report := reportService.AnalyseAccount()

	sendGridClient := initSendGridClient()
	emailSender := services.NewEmailSender(sendGridClient)

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

func initSendGridClient() pkg.SendGridClient {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	emailFrom := os.Getenv("EMAIL_FROM")

	return pkg.NewSendGridClient(host, apiKey, emailFrom)
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
func main() {
	lambda.Start(handler)
}
