// server.go

package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartServer(db *sql.DB) {
	router := gin.Default()
	// Setup routes
	SetupRoutes(router, db)

	// Start the server
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
