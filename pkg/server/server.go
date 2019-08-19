package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	userData   interface{}
	grpcServer *grpc.Server
	cfg        *ServerConfig
}

type ServerConfig struct {
	ListenAddress string
}

type ServiceServer interface {
	RegisterService(grpcServer *grpc.Server)
}

func NewServer(cfg *ServerConfig) (*Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("invalid configuration: %v", cfg)
	}

	ctxtags_opts := []grpc_ctxtags.Option{
		grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.TagBasedRequestFieldExtractor("log_fields")),
	}

	logger := zap.L()
	grpc_zap.ReplaceGrpcLogger(logger)
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(ctxtags_opts...),
			grpc_zap.StreamServerInterceptor(logger),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(ctxtags_opts...),
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	return &Server{grpcServer: s, cfg: cfg}, nil
}

func (s *Server) Start(ctx context.Context, services []ServiceServer) error {
	lis, err := net.Listen("tcp", s.cfg.ListenAddress)
	if err != nil {
		return err
	}

	for _, service := range services {
		service.RegisterService(s.grpcServer)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)
	go func() {
		defer s.stop()
		select {
		case <-sigs:
		case <-ctx.Done():
		}
	}()

	return s.grpcServer.Serve(lis)
}

func (s *Server) stop() {
	s.grpcServer.GracefulStop()
}
