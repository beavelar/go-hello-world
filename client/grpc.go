package client

import (
	"context"
	"fmt"
	"io"
	"log"
	proto "simple_project/genproto/server"
	"simple_project/server"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func endpointOneReq(client proto.SimpleServiceClient) {
	response, err := client.EndpointOne(context.Background(), &proto.SimpleRequest{Message: "EndpointOne client"})
	if err != nil {
		log.Printf("grpc client: failed to send request for EndpointOne - %v\n", err)
	}

	log.Printf("grpc client: received EndpointOne response - %v\n", response)
}

func endpointTwoReq(client proto.SimpleServiceClient) {
	response, err := client.EndpointTwo(context.Background(), &proto.SimpleRequest{Message: "EndpointTwo client"})
	if err != nil {
		log.Printf("grpc client: failed to send request for EndpointTwo - %v\n", err)
	}
	log.Printf("grpc client: received EndpointTwo response - %v\n", response)
}

func responseStreamReq(client proto.SimpleServiceClient) {
	stream, err := client.ResponseStream(context.Background(), &proto.SimpleRequest{Message: "ResponseStream client"})
	if err != nil {
		log.Printf("grpc client: field to start stream to ResponseStream - %v\n", err)
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("grpc client: error occurred receiving messages from ResponseStream - %v\n", err)
			return
		}

		log.Printf("grpc client: received message from ResponseStream - %v", response)
	}
}

func requestStreamReq(client proto.SimpleServiceClient) {
	stream, err := client.RequestStream(context.Background())
	if err != nil {
		log.Printf("grpc client: failed to start RequestStream - %v\n", err)
	}

	for {
		err := stream.Send(&proto.SimpleRequest{Message: "RequestStream client"})
		if err != nil {
			log.Printf("grpc client: failed to send message through RequestStream -%v\n", err)
			break
		}

		time.Sleep(time.Second)
	}
}

func bidirectionStreamReq(client proto.SimpleServiceClient) {
	stream, err := client.BidirectionalStream(context.Background())
	if err != nil {
		log.Printf("grpc client: failed to start BidirectionalStream - %v\n", err)
	}

	for {
		response, err := stream.Recv()
		log.Printf("grpc client: received message from BidirectionalStream - %v\n", response)

		time.Sleep(time.Second)
		if err == io.EOF {
			break
		}

		err = stream.Send(&proto.SimpleRequest{Message: "BidirectionalStream client"})
		if err != nil {
			log.Printf("grpc client: failed to send message through BidirectionalStream -%v\n", err)
			break
		}
	}
}

func StartGrpcClient() {
	log.Println("grpc client: setting up client..")

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", server.GrpcServerConfig.Host, server.GrpcServerConfig.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("grpc client: failed to connect to grpc server - %v\n", err)
	}
	defer conn.Close()

	client := proto.NewSimpleServiceClient(conn)
	go responseStreamReq(client)
	go requestStreamReq(client)
	go bidirectionStreamReq(client)

	for {
		go endpointOneReq(client)
		go endpointTwoReq(client)
		time.Sleep(time.Second)
	}
}
