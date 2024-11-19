package command

import (
	"context"
	"github.com/dsxriiiii/l3x_pay/common/decorator"
	"github.com/dsxriiiii/l3x_pay/common/genproto/orderpb"
	"github.com/sirupsen/logrus"
	domain "github/dsxriiiii/l3x_pay/order/domain/order"
)

type CreateOrder struct {
	CustomerID string
	Items      []*orderpb.ItemWithQuantity
}

type CreateOrderResult struct {
	OrderID string
}

type CreateOrderHandle decorator.CommandHandler[CreateOrder, *CreateOrderResult]

type createOrderHandle struct {
	orderRepo domain.Repository
}

func NewCreateOrderHandle(
	orderRepo domain.Repository,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CreateOrderHandle {
	if orderRepo == nil {
		panic("create order err, orderRp is nil")
	}
	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
		createOrderHandle{orderRepo: orderRepo},
		logger,
		metricClient,
	)
}

func (c createOrderHandle) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
	// TODO call stock grpc to get items
	var stockResponse []*orderpb.Item
	for _, item := range cmd.Items {
		stockResponse = append(stockResponse, &orderpb.Item{
			ID:       item.ID,
			Quantity: item.Quantity,
		})
	}
	order, err := c.orderRepo.Create(ctx, &domain.Order{
		CustomerID: cmd.CustomerID,
		Items:      stockResponse,
	})
	if err != nil {
		return nil, err
	}
	return &CreateOrderResult{OrderID: order.ID}, nil
}
