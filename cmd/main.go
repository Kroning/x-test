package main

import (
	"github.com/Kroning/x-test/internal/config"
	"github.com/Kroning/x-test/internal/database/postgresql"
	"github.com/Kroning/x-test/internal/handlers"
	loglib "github.com/Kroning/x-test/internal/logger"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// logger
	logLevel := new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}))
	slog.SetDefault(logger)

	// config
	var mainConfig config.Config
	err := config.ReadConfig(&mainConfig)
	if err != nil {
		log.Fatalf("Read config error: %s", err.Error())
	}

	// logger set level
	logLevel.Set(loglib.Level[mainConfig.Logger.Level])
	logger.Info("Starting...")
	logger.Debug("Config vars", "Config", mainConfig)

	db, err := postgresql.InitAndMigrate(mainConfig.Postgres)
	if err != nil {
		log.Fatalf("DB initialization error: %s", err.Error())
	}

	app := handlers.App{Db: db, Logger: logger, JWTSecret: []byte(mainConfig.Server.JWTKey)}

	router := gin.New()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})
	// Auth can be a separate service, but I will use func for simplicity
	router.POST("/company", app.CheckAuth, app.CreateCompany)
	router.PATCH("/company/:id", app.CheckAuth, app.PatchCompany)
	router.DELETE("/company/:id", app.CheckAuth, app.DeleteCompany)
	router.GET("/company/:id", app.GetCompany)

	logger.Info("Starting gin server on", "port", mainConfig.Server.Port)
	router.Run(mainConfig.Server.Port)

	logger.Info("Stopping service")
	err = db.DB.Close()
	if err != nil {
		logger.Warn("db.DB.Close()", "error", err.Error())
	}
}
