package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"io"
	"log"
	"stori-account-summary/model"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type Downloader struct {
	S3Event events.S3Event
	Client  S3Client
}

func NewDownloader(ctx context.Context, s3Event events.S3Event, client S3Client) Downloader {
	return Downloader{
		S3Event: s3Event,
		Client:  client,
	}
}

func (d *Downloader) DownloadFile(ctx context.Context) (model.Rows, error) {
	record := d.S3Event.Records[0]
	bucket := record.S3.Bucket.Name
	key := record.S3.Object.Key

	input := s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	output, err := d.Client.GetObject(ctx, &input)
	if err != nil {
		return model.Rows{}, err
	}
	defer output.Body.Close()

	return parseFile(output.Body)
}

func parseFile(ioReader io.Reader) (model.Rows, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(ioReader)

	if err != nil {
		return model.Rows{}, err
	}

	content := buf.String()
	csvReader := csv.NewReader(strings.NewReader(content))
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Printf("Unable to parse file as CSV, %v", err)
		return model.Rows{}, err
	}

	rows := model.Rows{}
	for _, record := range records {
		dateArray := strings.Split(record[1], "/")
		value, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			return model.Rows{}, err
		}

		row := model.Row{
			Transaction: value,
			Date: model.Date{
				Day:   dateArray[1],
				Month: dateArray[0],
			},
		}

		rows = append(rows, row)
	}

	return rows, nil
}
