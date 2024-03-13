package model

// User representa a estrutura dos dados do usuário
type User struct {
	ID       string `json:"id" dynamodbav:"ID"`
	Name     string `json:"name" dynamodbav:"Name"`
	Email    string `json:"email" dynamodbav:"Email"`
	Username string `json:"username" dynamodbav:"Username"`
	Password string `json:"password" dynamodbav:"Password"`
}
