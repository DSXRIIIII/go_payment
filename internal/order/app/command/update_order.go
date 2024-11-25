package command

import (
	"context"
	"github.com/dsxriiiii/l3x_pay/common/decorator"
	domain "github.com/dsxriiiii/l3x_pay/order/domain/order"
	"github.com/sirupsen/logrus"
)

type UpdateOrder struct {
	Order    *domain.Order
	UpdateFn func(ctx context.Context, order *domain.Order) (*domain.Order, error)
}

type UpdateOrderHandler decorator.CommandHandler[UpdateOrder, interface{}]

type updateOrderHandler struct {
	orderRepo domain.Repository
	// TODO stock rpc
}

func NewUpdateOrderHandler(orderRepo domain.Repository, logger *logrus.Entry, metricClient decorator.MetricsClient) UpdateOrderHandler {
	if orderRepo == nil {
		panic("update order err,nil repo found")
	}
	return decorator.ApplyCommandDecorators[UpdateOrder, interface{}](
		updateOrderHandler{orderRepo: orderRepo},
		logger,
		metricClient,
	)
}

func (u updateOrderHandler) Handle(ctx context.Context, cmd UpdateOrder) (interface{}, error) {
	if cmd.UpdateFn == nil {
		logrus.Warnf("updateOrderHandle got nil UpdateFn,order = %#v", cmd.Order)
		cmd.UpdateFn = func(_ context.Context, order *domain.Order) (*domain.Order, error) {
			return order, nil
		}
	}
	if err := u.orderRepo.Update(ctx, cmd.Order, cmd.UpdateFn); err != nil {
		return nil, err
	}
	return nil, nil
}
