package rest

import "github.com/gin-gonic/gin"

func (rH RoutesHandler) getNodesHandler(c *gin.Context) {
	nodes, err := rH.usecases.GetNodes()

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, nodes)
}

func (rH RoutesHandler) getSpecificNodeHandler(c *gin.Context) {
	nodeID := c.Param("node_id")

	node, err := rH.usecases.GetSpecificNode(nodeID)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, node)
}

func (rH RoutesHandler) getSpecificNodeInstancesHandler(c *gin.Context) {
	nodeID := c.Param("node_id")

	instances, err := rH.usecases.GetSpecificNodeInstances(nodeID)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, instances)
}
