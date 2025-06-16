package store

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/erobx/bobBot/internal/types"
)

type PollStore struct {
	client 		*dynamodb.Client
	tableName 	string
}

func NewPollStore(client *dynamodb.Client, tName string) *PollStore {
	return &PollStore{
		client: client,
		tableName: tName,
	}
}

func (d *PollStore) All(ctx context.Context) {

}

func (d *PollStore) Get(ctx context.Context, question string) (*types.Poll, error) {
	response, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &d.tableName,
		Key: map[string]ddbtypes.AttributeValue{
			"question": &ddbtypes.AttributeValueMemberS{Value: question},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get item from Flight: %w", err)
	}

	if len(response.Item) == 0 {
		return nil, nil
	}

	poll := types.Poll{}
	err = attributevalue.UnmarshalMap(response.Item, &poll)
	if err != nil {
		return nil, fmt.Errorf("error getting item %w", err)
	}

	return &poll, nil
}

func (d *PollStore) Put(ctx context.Context, poll types.Poll) error {
	item, err := attributevalue.MarshalMap(&poll)
	if err != nil {
		return fmt.Errorf("unable to marshal flight: %v", err)
	}

	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item: item,
	})

	if err != nil {
		return fmt.Errorf("cannot put item: %w", err)
	}

	return nil
}
