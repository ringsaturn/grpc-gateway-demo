package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	helloworldpb "github.com/ringsaturn/grpc-gateway-demo/proto/helloworld"
)

type server struct {
	helloworldpb.UnimplementedGreeterServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) SayHello(ctx context.Context, in *helloworldpb.HelloRequest) (*helloworldpb.HelloReply, error) {
	return &helloworldpb.HelloReply{Message: in.Name + " world"}, nil
}

var customErr = status.Errorf(codes.Code(10010), "custom error")

func (s *server) Bar(ctx context.Context, in *helloworldpb.BarRequest) (*helloworldpb.BarResponse, error) {
	if in.Name == "fake" {
		return nil, status.Errorf(codes.InvalidArgument, "fake not support")
	}
	if in.Name == "error" {
		return nil, errors.New("cause error")
	}
	if in.Name == "cus" {
		return nil, customErr
	}
	return &helloworldpb.BarResponse{
		Name: in.Name,
	}, nil
}

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8070")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	helloworldpb.RegisterGreeterServer(s, &server{})
	// Serve gRPC server
	log.Println("Serving gRPC on localhost:8070")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	// ctx :=
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// if err != nil {
	// 	log.Fatalln("Failed to dial server:", err)
	// }
	conn, err := grpc.DialContext(
		ctx,
		"localhost:8070",
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// Register Greeter
	err = helloworldpb.RegisterGreeterHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
