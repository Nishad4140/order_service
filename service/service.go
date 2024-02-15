package service

import (
	"context"
	"fmt"
	"io"

	"github.com/Nishad4140/order_service/adapter"
	helperstruct "github.com/Nishad4140/order_service/helper_struct"
	"github.com/Nishad4140/proto_files/pb"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

var (
	Tracer     opentracing.Tracer
	CartClient pb.CartServiceClient
)

func RetrieveTracer(tr opentracing.Tracer) {
	Tracer = tr
}

type OrderService struct {
	Adapter adapter.AdapterInterface
	pb.UnimplementedOrderServiceServer
}

func NewOrderService(adapter adapter.AdapterInterface) *OrderService {
	return &OrderService{
		Adapter: adapter,
	}
}

func (order *OrderService) OrderAll(ctx context.Context, req *pb.UserId) (*pb.OrderId, error) {
	span := Tracer.StartSpan("orderall grpc")
	defer span.Finish()

	cartItems, err := CartClient.GetAllCart(context.TODO(), &pb.CartCreate{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get items from cart")
	}
	var cart []helperstruct.OrderAll
	for {
		items, err := cartItems.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		item := helperstruct.OrderAll{
			ProductId: uint(items.ProductId),
			Quantity:  float64(items.Quantity),
			Total:     uint(items.Total),
		}
		cart = append(cart, item)
	}
	if len(cart) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}
	if _, err := CartClient.TruncateCart(context.TODO(), &pb.CartCreate{
		UserId: req.UserId,
	}); err == nil {
		return nil, err
	}
	orderId, err := order.Adapter.OrderAll(cart, uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.OrderId{OrderId: uint32(orderId)}, nil
}

func (order *OrderService) CancelOrder(ctx context.Context, req *pb.OrderId) (*pb.OrderId, error) {
	err := order.Adapter.CancelOrder(uint(req.OrderId))
	if err != nil {
		return nil, err
	}
	return &pb.OrderId{OrderId: req.OrderId}, nil
}

type HealthChecker struct {
	grpc_health_v1.UnimplementedHealthServer
}

func (s *HealthChecker) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	fmt.Println("check called")
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *HealthChecker) Watch(in *grpc_health_v1.HealthCheckRequest, srv grpc_health_v1.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "watching is not supported")
}
