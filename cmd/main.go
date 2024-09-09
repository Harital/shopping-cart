package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	httpHandlers "github.com/Harital/shopping-cart/internal/adapters/handlers/http"
	"github.com/Harital/shopping-cart/internal/adapters/repositories/mysql"
	"github.com/Harital/shopping-cart/internal/core/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	localPort       = "8080"
	defaultLogLevel = zerolog.InfoLevel
	shutdownTimeout = 5 * time.Second
)

func setupRouter() *gin.Engine {
	// Using gin, as it is a very useful (and easy to use) http server engine
	gin.SetMode(gin.DebugMode) // Debug mode for testing purposes. Should be changed to release in pruduction
	router := gin.Default()    // Getting a default router. It will  do the job

	return router
}

func main() {

	// Using zerolog as it allows to write the logs in json format. It´s useful for parsing them
	// Default log level should be set in an variable that could be adjustable, ideally in run time to debug potential problems.
	zerolog.SetGlobalLevel(defaultLogLevel)

	router := setupRouter()

	// default context that handles the signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	// We should setup an authenticated endpoint in production. In order to do that, middlewares should be added
	routerGroupWithoutAuth := router.Group("/shopping-cart/v1/")

	// Create all handlers, services and repos. In production code a dependency injection tool may be advisable,
	// as this part could get potentially big
	db, dbErr := mysql.InitMySqlDB()
	if dbErr != nil {
		log.Error().Msg("Cannot connect to the database " + dbErr.Error())
	}
	repo := mysql.NewCartItemsRepository(db)

	// These host values should be retrieved from env variables depending on the stac
	svc := services.NewCartItemsService(repo, "http://www.reservationhost.com", 60*time.Second)
	h := httpHandlers.NewCartItemsHandler(routerGroupWithoutAuth, svc)
	h.Register()

	// Start gin service
	log.Debug().Msg("Running")
	ginSrv := &http.Server{
		Addr:    ":" + localPort,
		Handler: router,
	}

	// Run gin service in a goroutine in order to catch signals in the main thread
	go func() {
		runErr := ginSrv.ListenAndServe()
		if runErr != nil && !errors.Is(runErr, http.ErrServerClosed) {
			log.Fatal().Msg("could not start http server: " + runErr.Error())
		}
	}()

	<-ctx.Done()
	stop()

	// Shut down with timeout in order to allow any running operation to finish
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if shutdownErr := ginSrv.Shutdown(ctx); shutdownErr != nil {
		log.Error().Msg("Forcing shutdown: " + shutdownErr.Error())
	}

	log.Debug().Msg("Stopped")
}
