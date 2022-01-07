package clusterDomain

type ConfigMap struct {
	Name    string
	Content []Environment
}

type ConfigMapCreationRequest struct {
	Name string `binding:"required"`
}

type ConfigMapUpdateRequest struct {
	Content []Environment `binding:"required"`
}
