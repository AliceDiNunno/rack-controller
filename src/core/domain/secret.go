package domain

type Secret struct {
	Name    string
	Content []Environment
}

type SecretCreationRequest struct {
	Name string `binding:"required"`
}

type SecretUpdateRequest struct {
	Content []Environment `binding:"required"`
}
