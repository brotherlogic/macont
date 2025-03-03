package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/habridge/proto"

	auth_client "github.com/brotherlogic/auth/client"
)

var (
	port        = flag.Int("port", 8080, "gRPC Port")
	metricsPort = flag.Int("metrics_port", 8081, "Metrics port")
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	authModule, err := auth_client.NewAuthInterceptor(ctx)
	if err != nil {
		log.Fatalf("Unable to register with auth: %v", err)
	}
	cancel()

	s := NewServer()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("gramophile is unable to listen on the grpc port %v: %v", *port, err)
	}
	gs := grpc.NewServer(grpc.UnaryInterceptor(authModule.AuthIntercept))
	pb.RegisterHabridgeServiceServer(gs, s)
	go func() {
		if err := gs.Serve(lis); err != nil {
			log.Fatalf("gramophile is unable to serve grpc: %v", err)
		}
		log.Fatalf("gramophile has closed the grpc port for some reason")
	}()

	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(fmt.Sprintf(":%v", *metricsPort), nil)
	if err != nil {
		log.Fatalf("gramophile is unable to serve metrics: %v", err)
	}
	log.Printf("Exiting after safe shutdown")
}
