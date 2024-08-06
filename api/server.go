package api

import (
	db "banking_application/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	//add routes to router

	router.POST("/accounts", server.createAccount)
	router.GET("/getAccount/:account_number", server.getAccount)
	router.GET("/getAccounts/", server.listAccounts)
	router.PUT("/updateAccount/", server.updateAccount)
	router.DELETE("/deleteAccount/", server.deleteAccount)

	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

// Run HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
