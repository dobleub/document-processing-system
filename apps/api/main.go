package main

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"nx-recipes/dps/lambda/config"
	docs "nx-recipes/dps/lambda/docs"
	"nx-recipes/dps/lambda/interfaces"
	"nx-recipes/dps/lambda/lib/database"
	"nx-recipes/dps/lambda/logger"
	"nx-recipes/dps/lambda/middlewares"
	processDomainHandlers "nx-recipes/dps/lambda/src/processDomain/handlers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

var (
	router     *gin.Engine
	appContext context.Context
	appLogger  *zap.Logger
	env        *config.Config
	appErr     error
)

func init() {
	appLogger = logger.Instance
	appContext = context.Background()

	env, appErr = config.SetUp(appContext)
	if appErr != nil {
		panic(appErr)
	}

	// add mongodb to context
	mongodbClient, err := database.ConnectMongoDB(*env.MongoDBConfig())
	if err != nil {
		appLogger.Fatal("error connecting to MongoDB", zap.Error(err))
	}
	dbContext := interfaces.MongoDBContext{
		Client: mongodbClient,
		DB:     mongodbClient.Database(env.MongoDBConfig().DB),
	}
	appLogger.Info("DB Connection Setted Up")

	// setup router
	router = gin.New()
	var state sync.Map // Initialize the state map, this states will be used to store the status of each process, it will be shared across all handlers
	middlewares.Setup(router, appLogger, env, dbContext, &state)

	// setup routes
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})
	docs.SwaggerInfo.BasePath = "/process"
	processRouter := router.Group("/process")
	{
		processRouter.POST("/start", processDomainHandlers.StartProcessHandler)
		processRouter.POST("/stop/:id", processDomainHandlers.StopProcessHandler)
		processRouter.GET("/status/:id", processDomainHandlers.StatusProcessHandler)
		processRouter.GET("/list", processDomainHandlers.ListProcessHandler)
		processRouter.GET("/results/:id", processDomainHandlers.ResultsProcessHandler)
	}
	// add swagger docs route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// add health check route
	router.GET("/health", func(c *gin.Context) {
		// TODO: add real health check logic
		c.String(http.StatusOK, "OK")
	})
	appLogger.Info("API Handler Setted Up")
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	defer func() {
		_ = logger.Instance.Sync()
	}()

	return httpadapter.New(router).ProxyWithContext(ctx, req)
}

func main() {
	isRunningAtLambda := strings.Contains(env.AWSConfig().ExecutionEnv, "AWS_Lambda_")
	appLogger.Info("Running at Lambda:", zap.Bool("isRunningAtLambda", isRunningAtLambda))

	if isRunningAtLambda {
		lambda.StartWithOptions(handler, lambda.WithContext(appContext))
	} else {
		appLogger.Info("Running locally, starting HTTP server on port 8080")
		if err := router.Run(":8080"); err != nil {
			appLogger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}
}
