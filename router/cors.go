package router

import (
	"time"

	"github.com/rs/cors"
)

// CORS returns a new negroni middleware to inject headers for CORS usage
func CORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowCredentials: true,
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Accept", "Authorization"},
		ExposedHeaders:   []string{"Content-Length", "Location"},
		MaxAge:           int(24 * time.Hour / time.Second),
	})
}
