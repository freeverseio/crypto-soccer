package relay

import (
	"github.com/gin-gonic/gin"
)

// Server - ...
type Server struct {
}

// Start - starts the server
func (s *Server) Start(port string) {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// GET
	router.GET("/relay/v1/:useraddr/nonce", NonceGET)

	// POST
	router.POST("/relay/v1/:useraddr/action", ActionPOST)
	router.POST("/relay/createuser", CreateUserPOST)

	// DEBUG
	router.GET("/relay/db", dbGET)

	// Listen and Server in localhost:8080
	router.Run(port)
}
