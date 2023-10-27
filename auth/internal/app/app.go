package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/patyukin/banking-system/auth/internal/closer"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/patyukin/banking-system/auth/internal/config"
	descAuth "github.com/patyukin/banking-system/auth/pkg/auth_v1"
	descUser "github.com/patyukin/banking-system/auth/pkg/user_v1"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
}

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

	var errHTTP, errGRPC error
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		errGRPC = a.runGRPCServer()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		errHTTP = a.runHTTPServer(ctx)
	}()

	wg.Wait()
	if errHTTP != nil {
		return fmt.Errorf("HTTP server failed: %w", errHTTP)
	}

	if errGRPC != nil {
		return fmt.Errorf("gRPC server failed: %w", errGRPC)
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
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
	var configPath string
	flag.StringVar(&configPath, "config-path", "", "path to config file")
	flag.Parse()

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
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)

	descUser.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserImpl(ctx))
	descAuth.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthImpl(ctx))

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

func (a *App) initHTTPServer(_ context.Context) error {
	mux := http.NewServeMux()
	a.httpServer = &http.Server{
		Addr:           a.serviceProvider.HTTPConfig().Address(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        mux,
	}

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	gateWayConn, err := grpc.DialContext(
		ctx,
		a.serviceProvider.GRPCConfig().Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)

	if err != nil {
		return err
	}

	grpcGwMux := runtime.NewServeMux()
	//err := descUser.RegisterUserV1HandlerFromEndpoint(ctx, grpcGwMux, a.serviceProvider.GRPCConfig().Address(), []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	err = descUser.RegisterUserV1Handler(ctx, grpcGwMux, gateWayConn)
	if err != nil {
		return err
	}

	//err = descAuth.RegisterAuthV1HandlerFromEndpoint(ctx, grpcGwMux, a.serviceProvider.GRPCConfig().Address(), []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	err = descAuth.RegisterAuthV1Handler(ctx, grpcGwMux, gateWayConn)
	if err != nil {
		return err
	}

	log.Printf("HTTP server is running on %s", a.serviceProvider.HTTPConfig().Address())
	err = a.httpServer.ListenAndServe()
	if err != nil {
		// TODO log
		return err
	}

	return nil
}
