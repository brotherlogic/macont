package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

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
	pstat, err := mclient.Ping(context.Background(), &pb.PingRequest{
		MachineName: name,
	})

	if err != nil {
		if status.Code(err) == codes.Unavailable {
			// Fail closed if we can't reach the macont server
			pstat = &pb.PingResponse{MachineState: pb.PingResponse_MACHINE_STATE_SHUTDOWN}
		} else {
			log.Fatalf("Unable to handle response: %v", err)
		}
	}

	if pstat.GetMachineState() == pb.PingResponse_MACHINE_STATE_SHUTDOWN {
		fmt.Printf("Shutting down\n")
		err = exec.Command("shutdown", "now").Run()
		if err != nil {
			log.Printf("Unable to shutdown: %v", err)
		}
	} else {
		fmt.Printf("Not shutting down: %v\n", pstat)
	}
}
