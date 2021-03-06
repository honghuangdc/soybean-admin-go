package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/honghuangdc/soybean-admin-go/api/e"
	db "github.com/honghuangdc/soybean-admin-go/db/sqlc"
	"github.com/honghuangdc/soybean-admin-go/token"
	"github.com/honghuangdc/soybean-admin-go/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	if err := InitTrans("zh"); err != nil {
		log.Fatalf("初始化验证翻译器错误: %s", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/users/renew_access", server.renewAccessToken)

	authRouters := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRouters.POST("/testauth", func(ctx *gin.Context) {
		appg := Gin{C: ctx}
		appg.Response(http.StatusOK, e.Success, "认证成功")
	})

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
