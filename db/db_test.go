package db_test

import (
	"testing"

	"stori-account-summary/db"
	"stori-account-summary/mocks"
	"stori-account-summary/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/assert"
)

func TestAddTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	defer ctrl.Finish()

	mockDynamoClient := mocks.NewMockDynamoClient(ctrl)
	rows := model.Rows{
		{
			Id:          "1",
			Date:        model.Date{Month: "01", Day: "15"},
			Transaction: 100.50,
		},
		{
			Id:          "2",
			Date:        model.Date{Month: "02", Day: "20"},
			Transaction: -50.00,
		},
	}

	expectedItems := []types.TransactWriteItem{
		{
			Put: &types.Put{
				TableName: aws.String("test-table"),
				Item: map[string]types.AttributeValue{
					"Id":          &types.AttributeValueMemberN{Value: "1"},
					"Date":        &types.AttributeValueMemberS{Value: "01/15"},
					"Transaction": &types.AttributeValueMemberN{Value: "100.50"},
				},
			},
		},
		{
			Put: &types.Put{
				TableName: aws.String("test-table"),
				Item: map[string]types.AttributeValue{
					"Id":          &types.AttributeValueMemberN{Value: "2"},
					"Date":        &types.AttributeValueMemberS{Value: "02/20"},
					"Transaction": &types.AttributeValueMemberN{Value: "-50.00"},
				},
			},
		},
	}

	mockDynamoClient.EXPECT().
		TransactWriteItems(gomock.Any(), &dynamodb.TransactWriteItemsInput{
			TransactItems: expectedItems,
		}).
		Return(&dynamodb.TransactWriteItemsOutput{}, nil)

	dbInstance := db.Db{
		DynamoDbClient: mockDynamoClient,
		TableName:      "test-table",
	}

	err := dbInstance.AddTransactions(rows)

	assert.NoError(t, err)
}
