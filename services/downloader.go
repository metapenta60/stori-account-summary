package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"io"
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
	output, err := d.Client.GetObject(ctx, d.createInput())
	if err != nil {
		return model.Rows{}, err
	}
	defer output.Body.Close()

	return parseFile(output.Body)
}

func (d *Downloader) createInput() *s3.GetObjectInput {
	record := d.S3Event.Records[0]
	bucket := record.S3.Bucket.Name
	key := record.S3.Object.Key

	return &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
}

func parseFile(ioReader io.Reader) (model.Rows, error) {
	rows := model.Rows{}

	content, err := getStringContent(ioReader)
	if err != nil {
		return model.Rows{}, err
	}

	records, err := getRows(content)
	if err != nil {
		return model.Rows{}, err
	}

	for i := 1; i < len(records); i++ {
		row := createRow(records[i])
		rows = append(rows, row)
	}

	return rows, nil
}

func getStringContent(ioReader io.Reader) (string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(ioReader)

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func getRows(file string) ([][]string, error) {
	csvReader := csv.NewReader(strings.NewReader(file))
	records, err := csvReader.ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}

func getDate(date string) model.Date {
	dateArray := strings.Split(date, "/")
	return model.Date{
		Day:   dateArray[1],
		Month: dateArray[0],
	}
}

func createRow(record []string) model.Row {
	date := getDate(record[1])
	value, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return model.Row{}
	}

	return model.Row{
		Id:          record[0],
		Transaction: value,
		Date:        date,
	}
}
