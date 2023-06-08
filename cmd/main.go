package main

import (
	"errors"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"ozon-test-unzhakov/constants"
	"ozon-test-unzhakov/internal/app"
	"ozon-test-unzhakov/internal/config"
	"ozon-test-unzhakov/internal/service"
	"ozon-test-unzhakov/internal/storage/cache"
	"ozon-test-unzhakov/internal/storage/pg"
	"ozon-test-unzhakov/internal/storage/storage"
	desc "ozon-test-unzhakov/pkg"
	"path/filepath"
)

func main() {

	err := config.InitConfig(filepath.Join("..", "config"), "config", "yaml")
	if err != nil {
		log.Fatal(err)
	}

	var linkStorage storage.LinkStorage
	storageType := viper.GetString("linkStorage")

	switch storageType {
	case constants.CacheStorageType:
		linkStorage, err = cache.NewLinkStorage()
	case constants.PostgresStorageType:
		linkStorage, err = pg.NewLinkStorage()
	default:
		log.Fatal(errors.New("unsupported linkStorage type"))
	}

	linkService, err := service.NewLinkService(linkStorage)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	microservice := app.NewMicroservice(linkService)

	desc.RegisterMicroserviceServer(grpcServer, microservice)

}
