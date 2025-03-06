package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	hapb "github.com/brotherlogic/habridge/proto"
	pb "github.com/brotherlogic/macont/proto"
	mdbpb "github.com/brotherlogic/mdb/proto"
)

func (s *Server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	conn, err := grpc.NewClient("mdb.mdb:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Bad mdb %w", err)
	}
	defer conn.Close()

	mdbclient := mdbpb.NewMDBServiceClient(conn)
	entry, err := mdbclient.GetMachine(ctx, &mdbpb.GetMachineRequest{Hostname: req.GetMachineName()})
	if err != nil {
		return nil, fmt.Errorf("unable to get machine: %w", err)
	}

	hconn, err := grpc.NewClient("habridge.habridge:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Bad ha: %w", err)
	}
	defer hconn.Close()

	haclient := hapb.NewHabridgeServiceClient(hconn)
	state, err := haclient.GetState(ctx, &hapb.GetStateRequest{
		ButtonId: "pixel_7.location",
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to get ha state: %w", err)
	}

	if entry.GetDetails().GetStability() == mdbpb.MachineStability_MACHINE_STABILITY_SHUTDOWN_ON_LEAVE &&
		state.GetUserState() == hapb.GetStateResponse_USER_STATE_AWAY {
		return &pb.PingResponse{MachineState: pb.PingResponse_MACHINE_STATE_SHUTDOWN}, nil
	}

	return &pb.PingResponse{MachineState: pb.PingResponse_MACHINE_STATE_DO_NOTHING}, nil
}
