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

func (rH RoutesHandler) getSpecificNodeHandler(context *gin.Context) {
	nodeID := context.Param("node_id")

	node, err := rH.usecases.GetSpecificNode(nodeID)

	if err != nil {
		rH.handleError(context, err)
		return
	}

	context.JSON(200, success(node))
}

func (rH RoutesHandler) getSpecificNodeInstancesHandler(context *gin.Context) {
	nodeID := context.Param("node_id")

	instances, err := rH.usecases.GetSpecificNodeInstances(nodeID)

	if err != nil {
		rH.handleError(context, err)
		return
	}

	context.JSON(200, success(instances))
}
