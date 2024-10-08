// Package main is the entry point to the server. It reads configuration, sets up logging and error handling,
// handles signals from the OS, and starts and stops the server.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	awss3 "stockinos.com/api/aws_s3"
	"stockinos.com/api/server"
	"stockinos.com/api/storage"
	"stockinos.com/api/utils"
)

// release is set through the linker at build time, generally from a git sha
// User for logging and error reporting
var release string

func init() {
	godotenv.Load()
}

func main() {
	fmt.Println("🤓")
	os.Exit(start())
}

func start() int {
	logEnv := utils.GetDefault("LOG_ENV", "development")
	log, err := createLogger(logEnv)
	if err != nil {
		fmt.Println("Error setting up the logger: ", err)
		return 1
	}

	log = log.With(zap.String("release", release))

	defer func() {
		_ = log.Sync()
	}()

	host := utils.GetDefault("HOST", "0.0.0.0")
	port := utils.GetIntDefault("PORT", 8080)

	database := storage.NewDatabase(storage.NewDatabaseOptions{
		URI:  utils.GetDefault("MONGODB_URI", "mongodb://localhost:27017/stockinos"),
		Name: utils.GetDefault("MONGDB_DBNAME", "stockinos"),
		Log:  log,
	})

	if err := database.Connect(); err != nil {
		log.Fatal("error connecting to database: %w", zap.Error(err))
	}

	s3Client := awss3.NewS3Client()

	s := server.New(server.Options{
		Database: database,
		S3:       s3Client,
		Host:     host,
		Port:     port,
		Log:      log,
	})

	// gs := grpc.New(grpc.Options{
	// 	Database: database,
	// 	Host:     host,
	// 	Port:     port + 10,
	// 	Log:      log,
	// })

	// var eg errgroup.Group
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := s.Start(); err != nil {
			return err
		}
		return nil
	})

	<-ctx.Done()

	eg.Go(func() error {
		<-ctx.Done()
		if err := s.Stop(); err != nil {
			return err
		}

		return nil
	})

	// wg := new(sync.WaitGroup)
	// wg.Add(2)

	// eg.Go(func() error {
	// 	if err := s.Start(); err != nil {
	// 		log.Info("Error starting http server", zap.Error(err))
	// 		// return 1
	// 		// wg.Done()
	// 		return err
	// 	}

	// 	return nil
	// })

	// eg.Go(func() error {
	// 	if err := gs.Start(); err != nil {
	// 		log.Info("Error starting grpc server", zap.Error(err))
	// 		// return 1
	// 		// wg.Done()
	// 		return err
	// 	}

	// 	return nil
	// })
	// go func() {
	// 	if err := gs.Start(); err != nil {
	// 		log.Info("Error starting grpc server", zap.Error(err))
	// 		// return 1
	// 		wg.Done()
	// 	}
	// }()

	// go func() {
	// 	if err := s.Start(); err != nil {
	// 		log.Info("Error starting server", zap.Error(err))
	// 		// return 1
	// 		wg.Done()
	// 	}
	// }()

	if err := eg.Wait(); err != nil {
		return 1
	}

	return 0
}

func createLogger(env string) (*zap.Logger, error) {
	switch env {
	case "production":
		return zap.NewProduction()
	case "development":
		return zap.NewDevelopment()
	default:
		return zap.NewNop(), nil
	}
}
