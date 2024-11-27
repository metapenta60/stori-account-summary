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

func handler(ctx context.Context, s3Event events.S3Event) error {
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Printf("Failed to load configuration, %v", err)
		return err
	}

	client := s3.NewFromConfig(cfg)

	downloader := services.NewDownloader(ctx, s3Event, client)
	rows, err := downloader.DownloadFile(ctx)
	if err != nil {
		log.Printf("Failed to download file, %v", err)
		return err
	}

	reportService := services.NewReportService(rows)
	report := reportService.AnalyseAccount()

	const host = "https://api.sendgrid.com"
	apiKey := os.Getenv("SENDGRID_API_KEY")

	sendGridClient := pkg.NewSendGridClient(host, apiKey, "v25a07@gmail.com")

	emailSender := services.NewEmailSender(sendGridClient)
	emailConfig := model.EmailConfig{
		To:     "v25a07@gmail.com",
		Report: report,
	}
	err = emailSender.SendMail(emailConfig)
	if err != nil {
		log.Printf("Failed to send email, %v", err)
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
