package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"simple_project/genproto/server"
	"time"

	"google.golang.org/grpc"
)

type simpleGrpcServer struct {
	server.UnimplementedSimpleServiceServer
}

func (s *simpleGrpcServer) EndpointOne(ctx context.Context, req *server.SimpleRequest) (*server.SimpleResponse, error) {
	log.Println("grpc server: received EndpointOne request")
	return &server.SimpleResponse{Message: "EndpointOne server"}, nil
}

func (s *simpleGrpcServer) EndpointTwo(ctx context.Context, req *server.SimpleRequest) (*server.SimpleResponse, error) {
	log.Println("grpc server: received EndpointTwo request")
	return &server.SimpleResponse{Message: "EndpointTwo server"}, nil
}

func (s *simpleGrpcServer) ResponseStream(req *server.SimpleRequest, stream server.SimpleService_ResponseStreamServer) error {
	log.Println("grpc server: received ResponseStream request")
	for {
		time.Sleep(time.Second)
		if err := stream.Send(&server.SimpleResponse{Message: "ResponseStream server"}); err != nil {
			return err
		}
	}
}

func (s *simpleGrpcServer) RequestStream(stream server.SimpleService_RequestStreamServer) error {
	log.Println("grpc server: received RequestStream request")
	for {
		response, err := stream.Recv()
		log.Printf("grpc server: received message from RequestStream - %v\n", response)

		if err == io.EOF {
			return stream.SendAndClose(&server.SimpleResponse{Message: "RequestStream closing"})
		}

		if err != nil {
			return err
		}
	}
}

func (s *simpleGrpcServer) BidirectionalStream(stream server.SimpleService_BidirectionalStreamServer) error {
	log.Println("grpc server: received BidirectionalStream request")
	stream.Send(&server.SimpleResponse{Message: "bidirectional stream start"})

	for {
		response, err := stream.Recv()
		log.Printf("grpc server: received message from BidirectionalStream - %v\n", response)

		time.Sleep(time.Second)
		if err != nil {
			return err
		}

		err = stream.Send(&server.SimpleResponse{Message: "BidirectionalStream server"})
		if err != nil {
			log.Printf("grpc server: failed to send message through BidirectionalStream - %v\n", err)
			return err
		}
	}
}

func StartGrpcServer() {
	log.Println("grpc server: setting up server..")
	if GrpcServerConfig == nil {
		GrpcServerConfig = &Config{Host: "localhost", Port: 55570}
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", GrpcServerConfig.Host, GrpcServerConfig.Port))
	if err != nil {
		log.Fatalf("grpc server: failed to listen - %v", err)
	}

	grpcServer := grpc.NewServer()
	server.RegisterSimpleServiceServer(grpcServer, &simpleGrpcServer{})
	grpcServer.Serve(lis)
}
