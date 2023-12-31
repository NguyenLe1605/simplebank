package api

import (
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)

// Server serve HTTP request for the client.
type Server struct {
	store db.Store
	router *gin.Engine
}

// NewServer create a new server and set up routing.
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// Assign the route handler
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	server.router = router
	return server
}

// Start runs the HTTP server at the specified address.
func (server *Server) Start(address  string) error {
	return server.router.Run(address)
}

// Create an approriate error response
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}