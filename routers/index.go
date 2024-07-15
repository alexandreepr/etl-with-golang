package routers

import (
	"net/http"
	"etl-with-golang/handlers"
	
	"github.com/gin-gonic/gin"
)

//RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Route Not Found"})
	})
	route.GET("/health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"status": "All good, captain!"}) })

	basePath := "/api/v1"

	// V1 Routes
	v1 := route.Group(basePath)
	{
		v1.POST("/file-import", handlers.ImportFile)
		v1.GET("/import-report", handlers.GetImportReport)
	}
}
