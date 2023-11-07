package app

import (
	"context"
	"fmt"
	"github.com/patyukin/banking-system/notifier/internal/config/env"
	"github.com/patyukin/banking-system/notifier/internal/queue/kafka"
	"github.com/patyukin/banking-system/notifier/internal/sender/email"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/patyukin/banking-system/notifier/internal/interceptor"
	"github.com/patyukin/banking-system/notifier/pkg/closer"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/patyukin/banking-system/notifier/internal/config"
	descNotifier "github.com/patyukin/banking-system/notifier/pkg/notifier_v1"
	_ "github.com/patyukin/banking-system/notifier/statik"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
	swaggerServer   *http.Server
}

const EnvFilePath = "ENV_FILE_PATH"

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			log.Fatalf("failed to run GRPC server: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			log.Fatalf("failed to run Swagger server: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := a.runConsumer(ctx)
		if err != nil {
			log.Fatalf("failed to run consumer: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
		a.initSender,
		a.initConsumer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	configPath := os.Getenv(EnvFilePath)
	err := config.Load(configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),
	)

	reflection.Register(a.grpcServer)

	descNotifier.RegisterNotifierV1Server(a.grpcServer, a.serviceProvider.NotifierImpl(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := descNotifier.RegisterNotifierV1HandlerFromEndpoint(ctx, mux, a.serviceProvider.GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:    a.serviceProvider.HTTPConfig().Address(),
		Handler: corsMiddleware.Handler(mux),
	}

	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:    a.serviceProvider.SwaggerConfig().Address(),
		Handler: mux,
	}

	return nil
}

func (a *App) initSender(_ context.Context) error {
	emailProviders, err := env.NewEmailProvider()
	if err != nil {
		return fmt.Errorf("failed to init email provider: %w", err)
	}

	if a.serviceProvider.sender == nil {
		a.serviceProvider.sender, err = email.New(emailProviders)
	}

	return nil
}

func (a *App) initConsumer(ctx context.Context) error {
	var err error
	a.serviceProvider.consumer, err = kafka.NewConsumer([]string{"localhost:9092"}, "my-topic")
	if err != nil {
		return fmt.Errorf("failed to init consumer: %w", err)
	}

	closer.Add(a.serviceProvider.consumer.Close)

	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	err := a.serviceProvider.consumer.RunConsume(ctx, a.serviceProvider.sender.Send)
	if err != nil {
		return fmt.Errorf("failed to init consumer: %w", err)
	}

	return nil
}
func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running on %s", a.serviceProvider.GRPCConfig().Address())

	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runHTTPServer() error {
	log.Printf("HTTP server is running on %s", a.serviceProvider.HTTPConfig().Address())

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	log.Printf("Swagger server is running on %s", a.serviceProvider.SwaggerConfig().Address())

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer func(file http.File) {
			err = file.Close()
			if err != nil {
				log.Printf("failed to close http.File, err: %v", err)
			}
		}(file)

		log.Printf("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Served swagger file: %s", path)
	}
}
