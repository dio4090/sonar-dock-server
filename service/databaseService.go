package service

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	database "github.com/dio4090/sonar-dock-server/config"
	"github.com/dio4090/sonar-dock-server/model"
	"github.com/dio4090/sonar-dock-server/util"
)

func InitDB(tableNames []string) error {
	sess := database.NewSession()
	db := dynamodb.New(sess) // Cria o cliente DynamoDB com a sessão

	// Supondo a existência de um construtor adequado para UserDao que aceita *dynamodb.DynamoDB
	userDao := NewUserDao(db)

	for _, tableName := range tableNames {
		// Verifica se a tabela existe
		tableExists, err := checkTableExists(db, tableName) // Presume-se que checkTableExists agora aceita *dynamodb.DynamoDB
		if err != nil {
			fmt.Printf("Erro ao verificar a existência da tabela %s: %v\n", tableName, err)
			continue
		}

		if !tableExists {
			fmt.Printf("Criando tabela %s...\n", tableName)
			err = CreateTable(db, tableName) // Presume-se que CreateTable agora aceita *dynamodb.DynamoDB
			if err != nil {
				fmt.Printf("Erro ao criar a tabela %s: %v\n", tableName, err)
				continue
			}
		}
	}

	// Cadastra o usuário admin padrão
	adminUser := model.User{
		ID:       util.GenerateUUID(),
		Name:     "Admin",
		Email:    os.Getenv("ADMIN_EMAIL"),
		Username: os.Getenv("ADMIN_USER"),
		Password: os.Getenv("ADMIN_PASSWORD"),
	}

	existingAdminUser, err := userDao.GetUserByEmail(os.Getenv("ADMIN_EMAIL"))
	if err != nil {
		fmt.Printf("Erro ao verificar a existência do usuário administrador: %v\n", err)
		return err
	}

	if existingAdminUser == nil {
		// O usuário não existe
		err := userDao.CreateAdminUser(adminUser)
		if err != nil {
			fmt.Printf("Erro ao criar usuário administrador: %v\n", err)
			return err
		}
	}

	return nil
}

func CreateTable(svc *dynamodb.DynamoDB, tableName string) error {
	attributeDefinitions := []*dynamodb.AttributeDefinition{
		{
			AttributeName: aws.String("ID"),
			AttributeType: aws.String("S"),
		},
	}

	keySchema := []*dynamodb.KeySchemaElement{
		{
			AttributeName: aws.String("ID"),
			KeyType:       aws.String("HASH"),
		},
	}

	// Inicialmente, define a lista de GSIs como vazia
	var globalSecondaryIndexes []*dynamodb.GlobalSecondaryIndex

	// Se o nome da tabela for "User", adiciona o GSI para "Email"
	if tableName == "User" {
		attributeDefinitions = append(attributeDefinitions, &dynamodb.AttributeDefinition{
			AttributeName: aws.String("Email"),
			AttributeType: aws.String("S"),
		})

		globalSecondaryIndexes = append(globalSecondaryIndexes, &dynamodb.GlobalSecondaryIndex{
			IndexName: aws.String("EmailIndex"),
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("Email"),
					KeyType:       aws.String("HASH"),
				},
			},
			Projection: &dynamodb.Projection{
				ProjectionType: aws.String("ALL"),
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
		})
	}

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: attributeDefinitions,
		KeySchema:            keySchema,
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(tableName),
	}

	if len(globalSecondaryIndexes) > 0 {
		input.GlobalSecondaryIndexes = globalSecondaryIndexes
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		return err
	}

	return nil
}

func checkTableExists(sess *dynamodb.DynamoDB, tableName string) (bool, error) {
	input := &dynamodb.ListTablesInput{}
	result, err := sess.ListTables(input)
	if err != nil {
		return false, err
	}

	for _, name := range result.TableNames {
		if *name == tableName {
			return true, nil
		}
	}

	return false, nil
}
