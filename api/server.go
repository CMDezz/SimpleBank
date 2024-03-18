package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/techschool/simplebank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	//add custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	//router section
	router.POST("/accounts/", server.createAccount)
	router.GET("/accounts/:id", server.getAccountById)
	router.GET("/accounts/", server.listAccounts)
	// router.DELETE("/accounts/:id", server.deleteAccount)
	router.PUT("/accounts", server.updateAccount)

	router.POST("/transfers/", server.createTransfer)

	//user
	router.POST("/users/", server.createUser)
	router.GET("/users/", server.getUserByUserName)
	//end

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run()
}

func errorResponse(err error) gin.H {
	return gin.H{"Error: ": err.Error()}
}
