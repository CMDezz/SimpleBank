package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/token"
	"github.com/techschool/simplebank/util"
	"github.com/techschool/simplebank/worker"
)

type Server struct {
	store           db.Store
	tokenMaker      token.Maker
	router          *gin.Engine
	config          util.Config
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	//add custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.SetUpRouter()
	return server, nil
}

func (server *Server) SetUpRouter() {
	router := gin.Default()

	//user
	router.POST("/users/login", server.loginUser)
	router.POST("/users/", server.createUser)
	router.POST("/users/refreshAccessToken", server.RefreshAccessToken)

	//auth group
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.GET("/users/", server.getUserByUserName)

	//router section
	authRoutes.POST("/accounts/", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccountById)
	authRoutes.GET("/accounts/", server.listAccounts)
	// router.DELETE("/accounts/:id", server.deleteAccount)
	authRoutes.PUT("/accounts", server.updateAccount)

	authRoutes.POST("/transfers/", server.createTransfer)

	//end
	server.router = router

}

func (server *Server) Start(address string) error {
	return server.router.Run()
}

func errorResponse(err error) gin.H {
	return gin.H{"Error: ": err.Error()}
}
