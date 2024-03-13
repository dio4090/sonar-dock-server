package service

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/dio4090/sonar-dock-server/model"
)

// Interface para operações de usuário
type UserDao interface {
	CreateUser(user model.User) error
	CreateAdminUser(user model.User) error
	GetUser(id string) (*model.User, error)
	GetAllUsers() ([]model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	UpdateUser(user model.User) error
	DeleteUser(id string) (*model.User, error)
}

func NewUserDao(db *dynamodb.DynamoDB) UserDao {
	return &dao{db: db}
}

type dao struct {
	db *dynamodb.DynamoDB
}

// CreateAdminUser implements UserDao.
func (d *dao) CreateAdminUser(user model.User) error {
	panic("unimplemented")
}

// CreateUser implements UserDao.
func (d *dao) CreateUser(user model.User) error {
	panic("unimplemented")
}

// DeleteUser implements UserDao.
func (d *dao) DeleteUser(id string) (*model.User, error) {
	panic("unimplemented")
}

// GetAllUsers implements UserDao.
func (d *dao) GetAllUsers() ([]model.User, error) {
	panic("unimplemented")
}

// GetUser implements UserDao.
func (d *dao) GetUser(id string) (*model.User, error) {
	panic("unimplemented")
}

// GetUserByEmail busca um usuário pelo email usando um GSI no DynamoDB.
func (d *dao) GetUserByEmail(email string) (*model.User, error) {
	// Construa a expressão de chave para consulta
	keyCond := expression.KeyEqual(expression.Key("Email"), expression.Value(email))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, err
	}

	// Defina os inputs da consulta
	input := &dynamodb.QueryInput{
		TableName:                 aws.String("User"),
		IndexName:                 aws.String("EmailIndex"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		Limit:                     aws.Int64(1), // Opcional: limita o resultado para melhorar a eficiência
	}

	// Execute a consulta
	result, err := d.db.Query(input)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, nil // Usuário não encontrado
	}

	var user model.User
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser implements UserDao.
func (d *dao) UpdateUser(user model.User) error {
	panic("unimplemented")
}
