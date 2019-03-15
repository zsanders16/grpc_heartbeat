package main

import (
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/zsanders16/grpc_heartbeat/pb"
	"google.golang.org/grpc"
)

const port = ":9000"

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Panic(err)
	}

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)
	pb.RegisterPingServer(server, &pingService{})

	log.Println("Starting server on port ", port)
	server.Serve(lis)
}

type pingService struct{}

func (p *pingService) Connected(stream pb.Ping_ConnectedServer) error {

	inChan := make(chan string)
	go func() {
		for {
			select {
			case msg := <-inChan:
				stream.Send(&pb.PingResponse{Text: msg})
				break
			}
		}
	}()

	for {
		rec, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fmt.Println(rec.Text)
		if rec.Text == "ping" {
			inChan <- "pong"
		}
	}

	return nil
}
