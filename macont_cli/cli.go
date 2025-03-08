package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/brotherlogic/macont/proto"
)

func main() {
	conn, err := grpc.NewClient(os.Args[1], grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	if status.GetMachineState() == pb.PingResponse_MACHINE_STATE_SHUTDOWN {
		fmt.Printf("Shutting down\n")
		err = exec.Command("shutdown", "now").Run()
		if err != nil {
			log.Printf("Unable to shutdown: %v", err)
		}
	}
}
