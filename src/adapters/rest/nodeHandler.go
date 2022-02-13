package rest

import "github.com/gin-gonic/gin"

func (rH RoutesHandler) getNodesHandler(context *gin.Context) {
	nodes, err := rH.usecases.GetNodes()

	if err != nil {
		rH.handleError(context, err)
		return
	}

	context.JSON(200, success(nodes))
}
