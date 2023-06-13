package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ervinismu/purplestore/internal/app/controller"
	"github.com/ervinismu/purplestore/internal/app/repository"
	"github.com/ervinismu/purplestore/internal/app/service"
	"github.com/ervinismu/purplestore/internal/app/util"
	"github.com/ervinismu/purplestore/internal/pkg/config"
	"github.com/ervinismu/purplestore/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

var cfg config.Config
var dbConn *sqlx.DB
var err error

func init() {
	// load configuration based on app.env
	cfg, err = config.LoadConfig()
	if err != nil {
		panic("failed to load config")
	}

	// Create database connection
	dbConn, err = sqlx.Open(cfg.DatabaseDriver, cfg.DatabaseURL)
	if err != nil {
		errMsg := fmt.Errorf("err database connect: %w", err)
		panic(errMsg)
	}

	err = dbConn.Ping()
	if err != nil {
		errMsg := fmt.Errorf("err database ping: %w", err)
		panic(errMsg)
	}

	// setup logrus logging
	logLevel, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		logLevel = log.InfoLevel
	}

	log.SetLevel(logLevel)
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	r := gin.New()
	r.Use(middleware.LogMiddleware())

	// init repo
	categoryRepository := repository.NewCategoryRepository(dbConn)
	userRepo := repository.NewUserRepository(dbConn)

	// init service
	categoryService := service.NewCategorySerivce(categoryRepository)
	userService := service.NewUserService(userRepo)

	// init controller
	categoryCotroller := controller.NewCategoryController(categoryService)
	userController := controller.NewUserController(&userService)

	// categories routes
	v1Routes := r.Group("api/v1")
	{

		//auth route
		v1Routes.POST("/auth/register", userController.Register)
		v1Routes.POST("/auth/login", userController.Login)

		v1Routes.GET("/secured/example",

			//middleware
			func(c *gin.Context) {
				reqToken := c.Request.Header.Get("Authorization")
				splitToken := strings.Split(reqToken, "Bearer ")
				accessToken := splitToken[1]

				err = util.VerifyJWT(accessToken)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{
						"data": err.Error(),
					})
					c.Abort()
					return
				}

				c.Next()
			},

			//main function
			func(c *gin.Context) {
				fmt.Println("this is the main function")
				c.JSON(http.StatusOK, gin.H{
					"data": "Success",
				})
			})

		v1Routes.GET("/categories", categoryCotroller.GetList)
		v1Routes.POST("/categories", categoryCotroller.Create)
		v1Routes.GET("/categories/:id", categoryCotroller.Detail)

	}

	// run server
	appPort := fmt.Sprintf(":%s", cfg.AppPort)
	err := r.Run(appPort)
	if err != nil {
		log.Panic("cannot start app")
	}
}
