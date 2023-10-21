package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	authHttpHandler "github.com/patyukin/banking-system/auth/internal/api/http/auth"
	userHttpHandler "github.com/patyukin/banking-system/auth/internal/api/http/user"
	"github.com/patyukin/banking-system/auth/internal/closer"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
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

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	var errHTTP, errGRPC error
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		errHTTP = a.runHTTPServer()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		errGRPC = a.runGRPCServer()
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
	err := config.Load(".env")
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
	addr := a.serviceProvider.GRPCConfig().Address()
	log.Printf("HTTP server is running on %s", addr)

	authHandler := authHttpHandler.New()
	userHandler := userHttpHandler.New()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Post("/sign-in", authHandler.SignIn)
	r.Post("/sign-in", userHandler.Create)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("Not found"))
		if err != nil {
			log.Println(err)
		}
	})

	a.httpServer = &http.Server{
		Addr:           addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        r,
	}

	return nil
}

func (a *App) runHTTPServer() error {
	err := a.httpServer.ListenAndServe()
	if err != nil {
		// log
		return err
	}

	return nil
}
