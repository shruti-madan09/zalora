package bennjerry

import (
	"github.com/gin-gonic/gin"

	"authenticator"
)

func RoutesBenNJerry(group *gin.RouterGroup) {
	// to create and save new ice cream data in DB
	group.POST("/", authenticator.IsAuthorized , CreateData)

	// to read ice cream data for a specific product id
	group.GET("/:product_id/", authenticator.IsAuthorized, ReadData)

	// to update ice cream data for a specific product id
	group.PUT("/:product_id/", authenticator.IsAuthorized, UpdateData)

	// to soft/permanent delete ice cream data for a specific product id
	group.DELETE("/:product_id/", authenticator.IsAuthorized, DeleteData)
}

