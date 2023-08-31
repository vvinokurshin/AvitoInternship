package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	r "github.com/vvinokurshin/AvitoInternship/cmd/router"
	"github.com/vvinokurshin/AvitoInternship/internal/config"
	historyDelivery "github.com/vvinokurshin/AvitoInternship/internal/history/delivery"
	historyRepository "github.com/vvinokurshin/AvitoInternship/internal/history/repository/postgres"
	historyUseCase "github.com/vvinokurshin/AvitoInternship/internal/history/usecase"
	segmentDelivery "github.com/vvinokurshin/AvitoInternship/internal/segment/delivery"
	segmentRepository "github.com/vvinokurshin/AvitoInternship/internal/segment/repository/postgres"
	segmentUseCase "github.com/vvinokurshin/AvitoInternship/internal/segment/usecase"
	userDelivery "github.com/vvinokurshin/AvitoInternship/internal/user/delivery"
	userRepository "github.com/vvinokurshin/AvitoInternship/internal/user/repository/postgres"
	userUseCase "github.com/vvinokurshin/AvitoInternship/internal/user/usecase"
	"github.com/vvinokurshin/AvitoInternship/pkg"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// @title Wishlist Swagger API
// @version 1.0
// @host localhost:8001
// @BasePath	/api/v1
func main() {
	var configFile string

	flag.StringVar(&configFile, "config", "cmd/config/config.yml", "-config=./cmd/config/config.yml")
	flag.Parse()

	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}

	globalLogger := pkg.LoggerInit(log.InfoLevel, *cfg.Logger.LogsUseStdOut, cfg.Logger.LogsFileName, cfg.Logger.LogsTimeFormat,
		cfg.Project.ProjectBaseDir, cfg.Logger.LogsDir)

	var prodCfgPg = postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s port=%s", cfg.DB.DBHost, cfg.DB.DBUser, cfg.DB.DBPassword,
			cfg.DB.DBPort),
	}

	db, err := gorm.Open(postgres.New(prodCfgPg), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	userRepo := userRepository.New(cfg, db)
	segmentRepo, err := segmentRepository.New(cfg, db)
	if err != nil {
		log.Fatal(err)
	}
	historyRepo := historyRepository.New(cfg, db)
	userUC := userUseCase.New(cfg, userRepo)
	segmentUC := segmentUseCase.New(cfg, segmentRepo, userRepo)
	historyUC := historyUseCase.New(cfg, historyRepo)
	userDel := userDelivery.New(cfg, userUC)
	segmentDel := segmentDelivery.New(cfg, segmentUC)
	historyDel := historyDelivery.New(cfg, historyUC)

	router := mux.NewRouter()
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	r.AddRoutes(router, cfg, userDel, segmentDel, historyDel)

	server := http.Server{
		Addr:         ":" + cfg.Project.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		globalLogger.Info("server started")
		if err := server.ListenAndServe(); err != nil {
			globalLogger.Fatalf("server stopped %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Kill, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	if err = server.Shutdown(ctx); err != nil {
		globalLogger.Errorf("failed to gracefully shutdown server")
	}
}
