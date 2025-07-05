// @title Test Go Project
// @version 1.0
// @description A part of auth-service made for an interview

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"net/http"
	"test_go_project/internal/auth"
	"test_go_project/internal/common"
	"test_go_project/internal/controller"
	"test_go_project/internal/db"
	"test_go_project/internal/repository"
	"test_go_project/internal/service"

	_ "test_go_project/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	db := db.Init(os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	jwtUtils := auth.JWTUtils{JwtKey: []byte(os.Getenv("JWT_KEY"))}
	notifier := common.WebhookNotifier{URL: os.Getenv("IP_NOTIFICATION_WEBHOOK"), Client: &http.Client{}}

	tokenRepository := repository.NewTokenRepo(db)
	tokenService := service.NewTokenService(tokenRepository, &jwtUtils, &notifier)
	tokenController := controller.NewTokenController(tokenService)

	r := gin.Default()

	r.GET("/api/tokens/:guid", tokenController.GetNewTokenPair)

	tokenGroup := r.Group("/api", auth.JWTMiddleware(jwtUtils))
	tokenGroup.GET("/tokens/refresh/:refresh_token", tokenController.UpdateTokenPair)
	tokenGroup.DELETE("/tokens", tokenController.ClearTokens)
	tokenGroup.GET("/user", tokenController.GetGUID)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
