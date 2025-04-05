package app

import (
	"context"
	"errors"
	"fmt"
	"gym-core/internal/config"
	"gym-core/internal/http/handlers"
	"gym-core/internal/http/middlewares"
	"gym-core/internal/http/route"
	"gym-core/internal/repositories"
	"gym-core/internal/services"
	"gym-core/pkg/logger"
	"gym-core/pkg/storage"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	config     *config.Config
	pgConn     *pgxpool.Pool
	httpServer *http.Server

	closers []func() error
}

func NewApp() (*App, error) {

	app := &App{}
	// - config
	app.config = config.MustLoadConfig()
	// - logger
	app.initLogger()
	// - db
	app.initDbConnection()

	app.initRest()

	return app, nil
}

func (a *App) initLogger() {
	switch a.config.Env {
	case config.EnvLocal:
		logger.SetupPrettySlog()
	case config.EnvDev, config.EnvProd:

		logFileCloser := logger.GetSlogFileConsoleJsonHandler()
		a.closers = append(a.closers, logFileCloser)
	}
	slog.Debug("config", slog.Any("data", a.config))
}

func (a *App) initDbConnection() {

	db, err := storage.PostgresConnect(a.config.PostgresConfig)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	slog.Info("Successfully connected to the database")

	a.pgConn = db
	a.closers = append(a.closers, func() error {
		db.Close()
		return nil
	})
}

func (a *App) initRest() {

	if a.config.HTTPServer.GinModeRelease {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middlewares.LoggingMiddleware)

	a.registerRoute(router)

	server := &http.Server{
		Addr:         a.config.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  a.config.HTTPServer.Timeout,
		WriteTimeout: a.config.HTTPServer.Timeout,
		IdleTimeout:  a.config.HTTPServer.IdleTimeout,
	}

	a.httpServer = server
}

func (a *App) registerRoute(router *gin.Engine) {

	// init handlers, services and repositories

	// User repository and service
	userRepository := repositories.NewUserRepository(a.pgConn)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	// Exercise repository and service
	exerciseRepository := repositories.NewExerciseRepository(a.pgConn)
	exerciseService := services.NewExerciseService(exerciseRepository)
	exerciseHandler := handlers.NewExerciseHandler(exerciseService)

	// Program repository and service
	programRepository := repositories.NewProgramRepository(a.pgConn)
	programService := services.NewProgramService(programRepository)
	programHandler := handlers.NewProgramHandler(programService)

	route.RegisterRoute(router, userHandler, exerciseHandler, programHandler)
}

func (a *App) Run() {

	go func() {
		slog.Info(fmt.Sprintf("Starting server on %s", a.config.HTTPServer.Address))
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start http server: \n %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), a.config.HTTPServer.ShutdownTimeout)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown:", slog.String("error", err.Error()))
		return
	}

	for i := len(a.closers) - 1; i >= 0; i-- {

		select {
		case <-ctx.Done():
			slog.Error("Closers time failed:", slog.String("error", ctx.Err().Error()))
			return
		default:
			err := a.closers[i]()
			if err != nil {
				slog.Error("failed to close resource", "i", i, "error", err)
			}
		}

	}

	slog.Info("Server exiting")

}
