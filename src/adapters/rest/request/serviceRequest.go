package request

type ServiceCreationRequest struct {
	Name      string `binding:"required"`
	ImageName string `binding:"required"`
}
