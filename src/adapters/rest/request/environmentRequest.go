package request

type EnvironmentCreationRequest struct {
	Name string `binding:"required"`
}
