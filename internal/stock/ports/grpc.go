package ports

import (
	"context"
	"github.com/dsxriiiii/l3x_pay/common/genproto/stockpb"
)

type GRPCServer struct {
}

func NewGrpcServer() *GRPCServer {
	return &GRPCServer{}
}

func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	//TODO implement me
	panic("implement me")
}