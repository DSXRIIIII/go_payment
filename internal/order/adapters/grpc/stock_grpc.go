package grpc

import (
	"context"
	"github.com/dsxriiiii/l3x_pay/common/genproto/orderpb"
	"github.com/dsxriiiii/l3x_pay/common/genproto/stockpb"
	"github.com/sirupsen/logrus"
)

type StockGRPC struct {
	client stockpb.StockServiceClient
}

func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
	return &StockGRPC{
		client: client,
	}
}

func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, itemIDs []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error) {
	resp, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{
		Items: itemIDs,
	})
	logrus.Info("stock_grpc response", resp)
	return resp, err

}

func (s StockGRPC) GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error) {
	resp, err := s.client.GetItems(ctx, &stockpb.GetItemsRequest{
		ItemIDs: itemIDs,
	})
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}
