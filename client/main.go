package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/zsanders16/grpc_heartbeat/pb"

	"google.golang.org/grpc"
)

const port = ":9000"

func main() {

	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial("localhost"+port, opts...)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	client := pb.NewPingClient(conn)

	stream, err := client.Connected(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	sendChan := make(chan string)
	go func() {
		for {
			time.Sleep(3 * time.Second)
			sendChan <- "ping"
		}
	}()

	doneChan := make(chan bool)
	go func() {
		for {
			select {
			case msg := <-sendChan:
				stream.Send(&pb.PingRequest{Text: msg})
				break
			}
		}
	}()

	for {
		rec, err := stream.Recv()
		if err == io.EOF {
			doneChan <- true
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(rec.Text)
	}

	<-doneChan
	stream.CloseSend()

}
