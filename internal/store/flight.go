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

type FlightStore struct {
	client 		*dynamodb.Client
	tableName 	string
}

func NewFlightStore(client *dynamodb.Client, tName string) *FlightStore {
	return &FlightStore{
		client: client,
		tableName: tName,
	}
}

func (d *FlightStore) All(ctx context.Context) ([]*types.Flight, error) {
	return nil, nil
}

func (d *FlightStore) Get(ctx context.Context, userID, flightKey string) (*types.Flight, error) {
	response, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &d.tableName,
		Key: map[string]ddbtypes.AttributeValue{
			"user_id": &ddbtypes.AttributeValueMemberS{Value: userID},
			"flight_key": &ddbtypes.AttributeValueMemberS{Value: flightKey},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get item from Flight: %w", err)
	}

	if len(response.Item) == 0 {
		return nil, nil
	}

	flight := types.Flight{}
	err = attributevalue.UnmarshalMap(response.Item, &flight)
	if err != nil {
		return nil, fmt.Errorf("error getting item %w", err)
	}

	return &flight, nil
}

func (d *FlightStore) Put(ctx context.Context, flight types.Flight) error {
	if flight.FlightKey == "" {
		if flight.ArrivalTime != "" {
			flight.FlightKey = fmt.Sprintf("%s#%s", flight.Number, flight.ArrivalTime)
		} else {
			return fmt.Errorf("departure_time required to generate flight_key")
		}
	}

	item, err := attributevalue.MarshalMap(&flight)
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
