package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/Nishad4140/order_service/db"
	"github.com/Nishad4140/order_service/initializer"
	"github.com/Nishad4140/order_service/service"
	servicediscoveryconsul "github.com/Nishad4140/order_service/servicediscovery_consul"
	"github.com/Nishad4140/proto_files/pb"
	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	jaegar "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err.Error())
	}
	cartConn, err := grpc.Dial("localhost:3003", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		cartConn.Close()
	}()
	cartRes := pb.NewCartServiceClient(cartConn)
	service.CartClient = cartRes
	addr := os.Getenv("DB_KEY")

	DB, err := db.InitDB(addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	services := initializer.Initialize(DB)
	server := grpc.NewServer()

	pb.RegisterOrderServiceServer(server, services)

	lis, err := net.Listen("tcp", ":3004")
	if err != nil {
		log.Fatal("failed to listen on port 3004")
	}
	log.Println("order server listen on port 3004")

	go func() {
		time.Sleep(5 * time.Second)
		servicediscoveryconsul.RegisterService()
	}()

	healthService := &service.HealthChecker{}
	grpc_health_v1.RegisterHealthServer(server, healthService)
	tracer, closer := InitTracer()
	defer closer.Close()
	service.RetrieveTracer(tracer)
	if err := server.Serve(lis); err != nil {
		log.Fatal("failed to listen on port 3004")
	}
}

func InitTracer() (tracer opentracing.Tracer, closer io.Closer) {
	jaegarEndpoint := "http://localhost:14268/api/tracer"
	cfg := &config.Configuration{
		ServiceName: "order-service",
		Sampler: &config.SamplerConfig{
			Type:  jaegar.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: jaegarEndpoint,
		},
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("updated")
	return
}
