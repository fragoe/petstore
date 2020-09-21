package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	ps "github.com/f-go/petstore/pkg/service"
)

// ...
func RunServer(ctx context.Context, service ps.PetstoreServiceServer, ms ps.MaintenanceServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":" + port)
	if err != nil {
		return err
	}

	// register service
	server := grpc.NewServer()
	ps.RegisterPetstoreServiceServer(server, service)
	ps.RegisterMaintenanceServiceServer(server, ms)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println("starting gRPC server...")
	return server.Serve(listen)
}