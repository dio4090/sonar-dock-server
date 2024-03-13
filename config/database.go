package config

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// NewSession cria e retorna uma nova sessão do AWS, não um cliente DynamoDB.
func NewSession() *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),                // Ou qualquer outra região necessária
		Endpoint: aws.String(os.Getenv("DYNAMODB_HOST")), // Endpoint local para testes ou remova para uso real na AWS
	})
	if err != nil {
		fmt.Printf("Erro ao criar sessão do AWS SDK: %s\n", err)
		os.Exit(1)
	}
	return sess
}

// NewSession cria e retorna uma nova sessão do DynamoDB.
func NewClientDB() *dynamodb.DynamoDB {
	sess, _ := session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String(os.Getenv("DYNAMODB_HOST")),
	})
	return dynamodb.New(sess)
}
