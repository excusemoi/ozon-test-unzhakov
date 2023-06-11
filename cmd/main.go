package main

import (
	"context"
	"errors"
	orm "github.com/go-pg/pg/v10"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"ozon-test-unzhakov/constants"
	"ozon-test-unzhakov/internal/app"
	"ozon-test-unzhakov/internal/config"
	"ozon-test-unzhakov/internal/service"
	"ozon-test-unzhakov/internal/storage/cache"
	"ozon-test-unzhakov/internal/storage/pg"
	"ozon-test-unzhakov/internal/storage/storage"
	"syscall"

	desc "ozon-test-unzhakov/pkg"
	"path/filepath"
)

func main() {

	err := godotenv.Load(filepath.Join(".env"))
	if err != nil {
		log.Fatal(err)
	}

	err = config.InitConfig(filepath.Join("config"), os.Getenv("CONFIG_NAME"), "yaml")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	grpcAddress := viper.GetString("grpcAddress")
	httpAddress := viper.GetString("httpAddress")

	listen, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatal(err)
	}

	var linkStorage storage.LinkStorage
	storageType := viper.GetString("storageType")

	switch storageType {
	case constants.CacheStorageType:
		linkStorage, err = cache.NewLinkStorage()
	case constants.PostgresStorageType:
		connectOptions, err := orm.ParseURL("postgres://" +
			os.Getenv("POSTGRES_USER") + ":" +
			os.Getenv("POSTGRES_PASSWORD") + "@" +
			viper.GetString("postgres.host") + ":" +
			viper.GetString("postgres.port") + "/" +
			viper.GetString("postgres.name") + "?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}
		s := orm.Connect(connectOptions)
		err = s.Ping(ctx)
		if err != nil {
			log.Fatal(err)
		}
		linkStorage, err = pg.NewLinkStorage(s)
	default:
		log.Fatal(errors.New("unsupported linkStorage type"))
	}
	if err != nil {
		log.Fatal(err)
	}

	linkService, err := service.NewLinkService(linkStorage, viper.GetInt("linksLength"))
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	microservice := app.NewMicroservice(linkService)

	desc.RegisterMicroserviceServer(grpcServer, microservice)

	mux := runtime.NewServeMux(runtime.WithForwardResponseOption(responseHeaderMatcher))
	err = desc.RegisterMicroserviceHandlerServer(ctx, mux, microservice)

	httpServer := http.Server{Addr: httpAddress, Handler: mux}
	if err != nil {
		log.Fatal(err)
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-exit
		grpcServer.Stop()
		httpServer.Shutdown(ctx)
	}()

	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		err := grpcServer.Serve(listen)
		return err
	})
	g.Go(func() error {
		err := httpServer.ListenAndServe()
		if err.Error() == "http: Server closed" {
			return nil
		}
		return err
	})

	log.Printf("service is ready to accept connections on port %s", httpAddress)

	err = g.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

func responseHeaderMatcher(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	headers := w.Header()
	if location, ok := headers["Grpc-Metadata-Location"]; ok {
		w.Header().Set("Location", location[0])
		w.WriteHeader(http.StatusFound)
	}
	return nil
}
