package main

import (
	"github.com/gin-gonic/gin"

	"bennjerry"
	"constants"
	"logger"
	"mysqlc"
)

func main() {
	// connecting to mysql
	mysqlc.Init()
	logger.Init()

	// Creating group route for bennjerry
	mainRouter := gin.Default()
	benNJerryGroup := mainRouter.Group("/bennjerry")
	bennjerry.RoutesBenNJerry(benNJerryGroup)

	// starting the server
	mainRouter.Run(constants.ServerHost + ":" + constants.ServerPort)
}
