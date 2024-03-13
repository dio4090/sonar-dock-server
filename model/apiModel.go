package model

// ErrorResponse representa uma resposta de erro padrão da API.
type ErrorResponse struct {
	Error string `json:"error"`
}

// NotFoundResponse representa uma resposta para recursos não encontrados.
type NotFoundResponse struct {
	Message string `json:"message"`
}

// ForbiddenResponse representa uma falha na autenticação.
type ForbiddenResponse struct {
	Message string `json:"Invalid username or password"`
}
