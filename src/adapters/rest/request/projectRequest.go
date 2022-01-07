package request

type CreateProjectRequest struct {
	Name string `json:"name" binding:"required"`
}
