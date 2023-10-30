package auth

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/patyukin/banking-system/auth/internal/api/auth"
	appInternal "github.com/patyukin/banking-system/auth/internal/app"
	desc "github.com/patyukin/banking-system/auth/pkg/auth_v1"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCTestSuite struct {
	suite.Suite

	server        *auth.Implementation
	clientConn    *grpc.ClientConn
	serviceClient desc.AuthV1Client
}

func (s *GRPCTestSuite) SetupSuite() {
	err := os.Setenv(appInternal.EnvFilePath, "./../../.env.local")
	if err != nil {
		s.FailNow("Failed to init app: " + err.Error())
	}

	ctx := context.Background()
	app, err := appInternal.NewApp(ctx)
	if err != nil {
		s.FailNow("Failed to init app: " + err.Error())
	}

	go app.Run()

	conn, err := grpc.Dial("0.0.0.0:11110", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.FailNow("Failed to init app: " + err.Error())
	}

	s.clientConn = conn
	s.serviceClient = desc.NewAuthV1Client(conn)
}

func (s *GRPCTestSuite) TearDownSuite() {
	s.clientConn.Close()
}

func (s *GRPCTestSuite) TestSignIn() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
	defer cancel()

	res, err := s.serviceClient.SignIn(ctx, &desc.AuthRequest{
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, false, false, false, false, 10),
	})

	s.Require().NoError(err)
	s.Require().NotNil(res)
}

func TestGRPCTestSuite(t *testing.T) {
	suite.Run(t, new(GRPCTestSuite))
}
