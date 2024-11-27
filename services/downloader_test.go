package services

import (
	"context"
	"io"
	"os"
	"stori-account-summary/mocks"
	"stori-account-summary/model"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/mock/gomock"
)

func TestDownloadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event := events.S3Event{
		Records: []events.S3EventRecord{
			{
				S3: events.S3Entity{
					Bucket: events.S3Bucket{
						Name: "bucket",
					},
					Object: events.S3Object{
						Key: "key",
					},
				},
			},
		},
	}

	ctx := context.Background()
	s3Client := mocks.NewMockS3Client(ctrl)

	output := &s3.GetObjectOutput{
		Body: readExampleTransactions(t),
	}
	defer output.Body.Close()

	s3Client.EXPECT().GetObject(ctx, gomock.Any()).Return(output, nil)

	downloader := NewDownloader(ctx, event, s3Client)
	rows, err := downloader.DownloadFile(ctx)

	expectedRows := model.Rows{
		model.Row{
			Id: "0",
			Date: model.Date{
				Day:   "15",
				Month: "7",
			},
			Transaction: 60.5,
		},
		model.Row{
			Id: "1",
			Date: model.Date{
				Day:   "28",
				Month: "7",
			},
			Transaction: -10.3,
		},
		model.Row{
			Id: "2",
			Date: model.Date{
				Day:   "2",
				Month: "8",
			},
			Transaction: -20.46,
		},
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedRows, rows)
}

func readExampleTransactions(t *testing.T) io.ReadCloser {
	file, err := os.Open("../data/transactions.csv")
	if err != nil {
		t.Fatal("Error al abrir el archivo:", err)
	}

	var ioReadCloser io.ReadCloser = file

	return ioReadCloser
}
