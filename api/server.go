package api

import (
	db "firstprj/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.GetAccount)
	router.GET("/accounts", server.GetListAccount)
	router.DELETE("/accounts/:id", server.DeleteAccount)
	router.PUT("/accounts/:id", server.UpdateAccount)

	server.router = router
	return server
}
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func createResponseError(err error) gin.H {
	return gin.H{"error": err.Error()}
}

type responseJson struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func createResponse(ctx *gin.Context, code int, data any, message string) {
	res := responseJson{Data: data, Message: message}
	ctx.JSONP(code, res)
}
