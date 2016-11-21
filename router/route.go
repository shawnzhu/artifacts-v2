package router

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/travis-ci/artifacts-v2/server"
	"github.com/travis-ci/artifacts-v2/store"
)

// Routes load middlewares
func Routes(middleware ...gin.HandlerFunc) http.Handler {
	var (
		router    = gin.Default()
		jwtSecret = os.Getenv("JWT_SECRET")
	)

	// TODO add other middlewares

	router.GET("/status", server.HealthCheck)

	router.Use(JWT(jwtSecret))

	router.Use(store.Store())

	router.POST("/upload/:build_id", server.UploadArtifact)

	artifacts := router.Group("/builds/:build_id")
	{
		artifacts.GET("", server.ListArtifacts)
		artifacts.GET("/artifacts/:artifact_id", server.GetArtifact)
	}

	return router
}
