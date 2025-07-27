package datasource

import (
	"context"
	"fmt"
	"time"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/port"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/database"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type customerDynamoDataSource struct {
	db *database.DynamoDatabase
}

type CustomerDynamoModel struct {
	ID    int    `dynamodbav:"id"`
	CPF   string `dynamodbav:"cpf"`
	Name  string `dynamodbav:"name"`
	Email string `dynamodbav:"email"`
}

func NewCustomerDynamoDataSource(db *database.DynamoDatabase) port.CustomerDataSource {
	return &customerDynamoDataSource{
		db: db,
	}
}

func (ds *customerDynamoDataSource) FindByID(ctx context.Context, id int) (*entity.Customer, error) {
	startTime := time.Now()

	input := &dynamodb.GetItemInput{
		TableName: aws.String(ds.db.TableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", id)},
		},
	}

	result, err := ds.db.Client.GetItem(ctx, input)

	duration := time.Since(startTime)
	ds.db.LogOperation(ctx, "FindByID", ds.db.TableName, duration, err)

	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	var customerModel CustomerDynamoModel
	err = attributevalue.UnmarshalMap(result.Item, &customerModel)
	if err != nil {
		return nil, err
	}

	return &entity.Customer{
		ID:    customerModel.ID,
		CPF:   customerModel.CPF,
		Name:  customerModel.Name,
		Email: customerModel.Email,
	}, nil
}

func (ds *customerDynamoDataSource) FindByCPF(ctx context.Context, cpf string) (*entity.Customer, error) {
	startTime := time.Now()

	// Use GSI for CPF lookup
	input := &dynamodb.QueryInput{
		TableName:              aws.String(ds.db.TableName),
		IndexName:              aws.String("cpf-index"),
		KeyConditionExpression: aws.String("cpf = :cpf"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":cpf": &types.AttributeValueMemberS{Value: cpf},
		},
		Limit: aws.Int32(1),
	}

	result, err := ds.db.Client.Query(ctx, input)

	duration := time.Since(startTime)
	ds.db.LogOperation(ctx, "FindByCPF", ds.db.TableName, duration, err)

	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var customerModel CustomerDynamoModel
	err = attributevalue.UnmarshalMap(result.Items[0], &customerModel)
	if err != nil {
		return nil, err
	}

	return &entity.Customer{
		ID:    customerModel.ID,
		CPF:   customerModel.CPF,
		Name:  customerModel.Name,
		Email: customerModel.Email,
	}, nil
}

func (ds *customerDynamoDataSource) FindAll(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entity.Customer, int64, error) {
	startTime := time.Now()

	// For simplicity, we'll scan the table (not optimal for large datasets)
	// In production, consider implementing pagination tokens
	input := &dynamodb.ScanInput{
		TableName: aws.String(ds.db.TableName),
		Limit:     aws.Int32(int32(limit)),
	}

	// Add filters if provided
	if len(filters) > 0 {
		filterExpression := ""
		expressionAttributeValues := make(map[string]types.AttributeValue)
		expressionAttributeNames := make(map[string]string)

		for key, value := range filters {
			if filterExpression != "" {
				filterExpression += " AND "
			}

			// Handle reserved keywords by using expression attribute names
			attributeName := fmt.Sprintf("#%s", key)
			expressionAttributeNames[attributeName] = key
			filterExpression += fmt.Sprintf("%s = :%s", attributeName, key)
			expressionAttributeValues[":"+key] = &types.AttributeValueMemberS{Value: fmt.Sprint(value)}
		}

		if filterExpression != "" {
			input.FilterExpression = aws.String(filterExpression)
			input.ExpressionAttributeValues = expressionAttributeValues
			input.ExpressionAttributeNames = expressionAttributeNames
		}
	}

	result, err := ds.db.Client.Scan(ctx, input)

	duration := time.Since(startTime)
	ds.db.LogOperation(ctx, "FindAll", ds.db.TableName, duration, err)

	if err != nil {
		return nil, 0, err
	}

	customers := make([]*entity.Customer, len(result.Items))
	for i, item := range result.Items {
		var customerModel CustomerDynamoModel
		err = attributevalue.UnmarshalMap(item, &customerModel)
		if err != nil {
			return nil, 0, err
		}

		customers[i] = &entity.Customer{
			ID:    customerModel.ID,
			CPF:   customerModel.CPF,
			Name:  customerModel.Name,
			Email: customerModel.Email,
		}
	}

	// For total count, we need another scan operation (simplified approach)
	countInput := &dynamodb.ScanInput{
		TableName: aws.String(ds.db.TableName),
		Select:    types.SelectCount,
	}

	if input.FilterExpression != nil {
		countInput.FilterExpression = input.FilterExpression
		countInput.ExpressionAttributeValues = input.ExpressionAttributeValues
		countInput.ExpressionAttributeNames = input.ExpressionAttributeNames
	}

	countResult, err := ds.db.Client.Scan(ctx, countInput)
	if err != nil {
		return customers, int64(len(customers)), nil // Return partial result
	}

	return customers, int64(countResult.Count), nil
}

func (ds *customerDynamoDataSource) getNextID(ctx context.Context) (int, error) {
	// Scan table to find the highest ID
	input := &dynamodb.ScanInput{
		TableName:            aws.String(ds.db.TableName),
		ProjectionExpression: aws.String("id"),
	}

	result, err := ds.db.Client.Scan(ctx, input)
	if err != nil {
		return 0, err
	}

	maxID := 0
	for _, item := range result.Items {
		var customerModel CustomerDynamoModel
		err = attributevalue.UnmarshalMap(item, &customerModel)
		if err != nil {
			continue
		}
		if customerModel.ID > maxID {
			maxID = customerModel.ID
		}
	}

	return maxID + 1, nil
}

func (ds *customerDynamoDataSource) Create(ctx context.Context, customer *entity.Customer) error {
	startTime := time.Now()

	if customer.ID == 0 {
		nextID, err := ds.getNextID(ctx)
		if err != nil {
			return err
		}
		customer.ID = nextID
	}

	customerModel := CustomerDynamoModel{
		ID:    customer.ID,
		CPF:   customer.CPF,
		Name:  customer.Name,
		Email: customer.Email,
	}

	item, err := attributevalue.MarshalMap(customerModel)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName:           aws.String(ds.db.TableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(id)"), // Prevent overwrite
	}

	_, err = ds.db.Client.PutItem(ctx, input)

	duration := time.Since(startTime)
	ds.db.LogOperation(ctx, "Create", ds.db.TableName, duration, err)

	return err
}

func (ds *customerDynamoDataSource) Update(ctx context.Context, customer *entity.Customer) error {
	startTime := time.Now()

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(ds.db.TableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", customer.ID)},
		},
		UpdateExpression: aws.String("SET #name = :name, email = :email, cpf = :cpf"),
		ExpressionAttributeNames: map[string]string{
			"#name": "name", // 'name' is a reserved keyword in DynamoDB
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name":  &types.AttributeValueMemberS{Value: customer.Name},
			":email": &types.AttributeValueMemberS{Value: customer.Email},
			":cpf":   &types.AttributeValueMemberS{Value: customer.CPF},
		},
		ConditionExpression: aws.String("attribute_exists(id)"), // Ensure item exists
	}

	_, err := ds.db.Client.UpdateItem(ctx, input)

	duration := time.Since(startTime)
	ds.db.LogOperation(ctx, "Update", ds.db.TableName, duration, err)

	return err
}

func (ds *customerDynamoDataSource) Delete(ctx context.Context, id int) error {
	startTime := time.Now()

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(ds.db.TableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", id)},
		},
		ConditionExpression: aws.String("attribute_exists(id)"), // Ensure item exists
	}

	_, err := ds.db.Client.DeleteItem(ctx, input)

	duration := time.Since(startTime)
	ds.db.LogOperation(ctx, "Delete", ds.db.TableName, duration, err)

	return err
}
