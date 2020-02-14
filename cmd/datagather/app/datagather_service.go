package app

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/ijsong/farseer/internal/service"
	"github.com/ijsong/farseer/pkg/queue"
	"google.golang.org/grpc"
)

const (
	datagatherTopic = "_topic_datagather"
)

type DatagatherService struct {
	numberOfProducers int
	producers         []*queue.EmbeddedQueueProducer
}

func NewDatagatherService(producers []*queue.EmbeddedQueueProducer) *DatagatherService {
	return &DatagatherService{
		numberOfProducers: len(producers),
		producers:         producers,
	}
}

func (ds *DatagatherService) RegisterService(grpcServer *grpc.Server) {
	service.RegisterDatagatherServiceServer(grpcServer, ds)
}

func (ds *DatagatherService) CreateEvent(ctx context.Context, req *service.CreateEventRequest) (*types.Empty, error) {
	// TODO: Check message
	if req.Event == nil {
		return nil, service.NewNotInitiatedMessageError("CreateEvent")
		//return nil, fmt.Errorf("nil event")
	}
	if req.Event.Timestamp.IsZero() {
		req.Event.Timestamp = time.Now()
	}
	if err := ds.publishDatagatherRequest(req); err != nil {
		return nil, err
	}
	return &types.Empty{}, nil
}

func (ds *DatagatherService) CreateItem(ctx context.Context, req *service.CreateItemRequest) (*types.Empty, error) {
	if req.Item == nil {
		return nil, fmt.Errorf("nil event")
	}
	if err := ds.publishDatagatherRequest(req); err != nil {
		return nil, err
	}
	return &types.Empty{}, nil
}

func (ds *DatagatherService) DeleteItem(ctx context.Context, req *service.DeleteItemRequest) (*types.Empty, error) {
	if err := ds.publishDatagatherRequest(req); err != nil {
		return nil, err
	}
	return &types.Empty{}, nil
}

func (ds *DatagatherService) UpdateItem(ctx context.Context, req *service.UpdateItemRequest) (*types.Empty, error) {
	if err := ds.publishDatagatherRequest(req); err != nil {
		return nil, err
	}
	return &types.Empty{}, nil
}

func (ds *DatagatherService) CreateUser(ctx context.Context, req *service.CreateUserRequest) (*types.Empty, error) {
	if err := ds.publishDatagatherRequest(req); err != nil {
		return nil, err
	}
	return &types.Empty{}, nil
}

func (ds *DatagatherService) DeleteUser(ctx context.Context, req *service.DeleteUserRequest) (*types.Empty, error) {
	if err := ds.publishDatagatherRequest(req); err != nil {
		return nil, err
	}
	return &types.Empty{}, nil
}

func (ds *DatagatherService) UpdateUser(ctx context.Context, req *service.UpdateUserRequest) (*types.Empty, error) {
	if err := ds.publishDatagatherRequest(req); err != nil {
		return nil, err
	}
	return &types.Empty{}, nil
}

func (ds *DatagatherService) publishDatagatherRequest(req interface{}) error {
	dataReq := &service.DatagatherRequest{}
	if !dataReq.SetValue(req) {
		return fmt.Errorf("invalid argument: %v", req)
	}
	bytes, err := proto.Marshal(dataReq)
	if err != nil {
		return err
	}
	producer := ds.producers[rand.Intn(ds.numberOfProducers)]
	if err := producer.Publish(bytes); err != nil {
		return err
	}
	return nil
}
