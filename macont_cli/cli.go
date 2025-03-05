package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/brotherlogic/macont/proto"
)

func main() {
	conn, err := grpc.Dial(os.Args[1], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Dial fail: %v", err)
	}

	name, err := os.Hostname()
	if err != nil {
		log.Fatalf("Getting a hostname failed: %v", err)

	}

	mclient := pb.NewMacontServiceClient(conn)
	status, err := mclient.Ping(context.Background(), &pb.PingRequest{
		MachineName: name,
	})

	fmt.Printf("%v [%v]\n", status, err)
}
