package main

import (
	"context"
	"net/http"
	"strings"

	"nx-recipes/dps/lambda/config"
	"nx-recipes/dps/lambda/interfaces"
	"nx-recipes/dps/lambda/lib/database"
	"nx-recipes/dps/lambda/logger"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var (
	router     *mux.Router
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

	// add env to context
	appContext = context.WithValue(appContext, interfaces.ConfigKey, env)
	appLogger.Info("Environment Setted Up")
	// add log to context
	appContext = context.WithValue(appContext, interfaces.LoggerKey, appLogger)
	appLogger.Info("Logger Setted Up")
	// add mongodb to context
	mongodbClient, err := database.ConnectMongoDB(*env.MongoDBConfig())
	if err != nil {
		appLogger.Fatal("error connecting to MongoDB", zap.Error(err))
	}
	dbContext := interfaces.MongoDBContext{
		Client: mongodbClient,
		DB:     mongodbClient.Database(env.MongoDBConfig().DB),
	}
	appContext = context.WithValue(appContext, interfaces.MongodbKey, dbContext)
	appLogger.Info("DB Connection Setted Up")

	// setup router
	router = mux.NewRouter()
	// setup routes
	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}))
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
		headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Access-Control-Allow-Origin"})
		credentials := handlers.AllowCredentials()
		methods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})
		ttl := handlers.MaxAge(3600)
		origins := handlers.AllowedOrigins([]string{"*"})
		http.ListenAndServe(":8080", handlers.CORS(headers, credentials, methods, ttl, origins)(router))
	}
}
