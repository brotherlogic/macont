package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	pb "github.com/brotherlogic/macont/proto"
)

var (
	port        = flag.Int("port", 8080, "gRPC Port")
	metricsPort = flag.Int("metrics_port", 8081, "Metrics port")
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

var (
	serverRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "macont_requests",
		Help: "The number of server requests",
	}, []string{"method", "status"})
)

func (s *Server) ServerInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	h, err := handler(ctx, req)
	serverRequests.With(
		prometheus.Labels{
			"status": status.Convert(err).Code().String(),
			"method": info.FullMethod},
	).Inc()
	return h, err
}

func main() {
	flag.Parse()

	s := NewServer()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("macont is unable to listen on the grpc port %v: %v", *port, err)
	}
	gs := grpc.NewServer()
	pb.RegisterMacontServiceServer(gs, s)
	go func() {
		if err := gs.Serve(lis); err != nil {
			log.Fatalf("macont is unable to serve grpc: %v", err)
		}
		log.Fatalf("macont has closed the grpc port for some reason")
	}()

	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(fmt.Sprintf(":%v", *metricsPort), nil)
	if err != nil {
		log.Fatalf("macont is unable to serve metrics: %v", err)
	}
	log.Printf("Exiting after safe shutdown")
}
