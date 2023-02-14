package grpc

import (
	"fmt"
	"net"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	pb "stockinos.com/api/grpc/protos"
	"stockinos.com/api/grpc/server"
	"stockinos.com/api/storage"
)

type GrpcServer struct {
	// listen *net.Listener
	address  string
	database *storage.Database
	log      *zap.Logger
}

type Options struct {
	Database *storage.Database
	Host     string
	Log      *zap.Logger
	Port     int
}

func New(opts Options) *GrpcServer {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}

	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))

	return &GrpcServer{
		address:  address,
		database: opts.Database,
		log:      opts.Log,
	}
}

func (gs *GrpcServer) Start() error {
	lis, err := net.Listen("tcp", gs.address)
	if err != nil {
		return fmt.Errorf("error starting grpc server: %w", err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterWhatsappServer(
		grpcServer,
		server.NewGrpcWhatsappServer(server.GrpcWhatsappOptions{
			Database: gs.database,
			Log:      gs.log,
		}),
	)

	gs.log.Info("Starting the gRPC server on ", zap.String("address", gs.address))
	grpcServer.Serve(lis)

	return nil
}

// func (gs *GrpcServer) Stop() error {

// 	if err != gs.
// }
