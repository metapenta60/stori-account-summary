package db

import (
	"context"
	"fmt"
	"log"
	"stori-account-summary/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoClient interface {
	TransactWriteItems(ctx context.Context, params *dynamodb.TransactWriteItemsInput, optFns ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error)
}

type Db struct {
	DynamoDbClient DynamoClient
	TableName      string
}

func (d Db) AddTransactions(rows model.Rows) error {
	ctx := context.Background()
	var transactItems []types.TransactWriteItem

	for _, row := range rows {
		dateString := fmt.Sprintf("%s/%s", row.Date.Month, row.Date.Day)
		transaction := fmt.Sprintf("%.2f", row.Transaction)
		item := map[string]types.AttributeValue{
			"Id":          &types.AttributeValueMemberN{Value: row.Id},
			"Date":        &types.AttributeValueMemberS{Value: dateString},
			"Transaction": &types.AttributeValueMemberN{Value: transaction},
		}

		transactItems = append(transactItems, types.TransactWriteItem{
			Put: &types.Put{
				TableName: aws.String(d.TableName),
				Item:      item,
			},
		})
	}

	_, err := d.DynamoDbClient.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})

	if err != nil {
		log.Printf("Couldn't add items to table in a single transaction. Here's why: %v\n", err)
		return err
	}

	return nil
}

func New(ctx context.Context, tableName string, client *dynamodb.Client) Db {
	return Db{
		DynamoDbClient: client,
		TableName:      tableName,
	}
}
