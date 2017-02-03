package router

import (
	"net/http"
	"time"

	"gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/travis-ci/artifacts-v2/server"
	"github.com/travis-ci/artifacts-v2/store"
)

// Routes load middlewares
func Routes(middleware ...gin.HandlerFunc) http.Handler {
	var router = gin.Default()

	// TODO add other middlewares

	router.GET("/status", server.HealthCheck)

	router.Use(store.Store())

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Location"},
		MaxAge:           24 * time.Hour,
	}))

	router.Use(Auth())

	router.POST("/upload/:build_id", server.UploadArtifact)

	artifacts := router.Group("/builds/:build_id")
	{
		artifacts.GET("", server.ListArtifacts)
		artifacts.GET("/artifacts/:artifact_id", server.GetArtifact)
	}

	return router
}
