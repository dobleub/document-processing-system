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
	pd_interfaces "nx-recipes/dps/lambda/src/processDomain/interfaces"
	summarizerDomainHandlers "nx-recipes/dps/lambda/src/summarizerDomain/handlers"
	websocketDomainHandlers "nx-recipes/dps/lambda/src/websocketDomain/handlers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/gin-gonic/gin"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	dbContext := &interfaces.MongoDBContext{
		DB: mongodbClient.Database(env.MongoDBConfig().DB),
	}
	appLogger.Info("DB Connection Setted Up")
	appContext = context.WithValue(appContext, interfaces.MongodbKey, dbContext)

	// add mcp client to context
	mcpHandlerSetUp := summarizerDomainHandlers.SetUpMCPHandler(appContext, env, appLogger)

	// setup router
	router = gin.New()
	// initialize the state map, this states will be used to store the status of each process, it will be shared across all handlers
	var state sync.Map

	// fill the state with the existing processes from MongoDB, this is useful when the API restarts and we want to keep track of the existing processes
	mongoClient := &interfaces.MongoCollection{}
	mongoClient.SetDBContext(appContext)
	mongoClient.SetCollectionName(pd_interfaces.CollectionName)
	// fetch existing processes from MongoDB and populate the state map
	limit := int64(10)
	cursor, count, err := mongoClient.Find(bson.M{}, &options.FindOptions{Limit: &limit})
	if err != nil {
		appLogger.Fatal("error fetching existing processes from MongoDB", zap.Error(err))
	}
	appLogger.Info("Existing processes fetched from MongoDB", zap.Int32("count", count))
	for cursor.Next(appContext) {
		var process pd_interfaces.OperationResponse
		if err := cursor.Decode(&process); err != nil {
			appLogger.Error("error decoding process from MongoDB", zap.Error(err))
			continue
		}
		state.Store(process.ID, &process)
	}

	// setup middlewares
	middlewares.Setup(router, appLogger, env, dbContext, &state, mcpHandlerSetUp.Client)
	// setup routes
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})
	// setup process routes
	docs.SwaggerInfo.BasePath = "/process"
	processRouter := router.Group("/process")
	{
		processRouter.POST("/start", processDomainHandlers.StartProcessHandler)
		processRouter.POST("/stop/:id", processDomainHandlers.StopProcessHandler)
		processRouter.GET("/status/:id", processDomainHandlers.StatusProcessHandler)
		processRouter.GET("/list", processDomainHandlers.ListProcessHandler)
		processRouter.GET("/results/:id", processDomainHandlers.ResultsProcessHandler)
	}
	// setup summarizer routes
	docs.SwaggerInfo.BasePath = "/summarizer"
	summarizerRouter := router.Group("/summarizer")
	{
		mcpHandler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
			return mcpHandlerSetUp.Server
		}, &mcp.StreamableHTTPOptions{})
		summarizerRouter.POST("/mcp", func(c *gin.Context) {
			mcpHandler.ServeHTTP(c.Writer, c.Request)
		})
		summarizerRouter.POST("/summarize", summarizerDomainHandlers.HttpSummarizerHandler)
	}
	// setup websocket route
	docs.SwaggerInfo.BasePath = "/ws"
	wsRouter := router.Group("/ws")
	{
		wsRouter.GET("/status", websocketDomainHandlers.StatusProcessHandler)
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
