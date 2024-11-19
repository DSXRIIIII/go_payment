package command

import (
	"context"
	"github.com/dsxriiiii/l3x_pay/common/decorator"
	"github.com/sirupsen/logrus"
	domain "github/dsxriiiii/l3x_pay/order/domain/order"
)

type UpdateOrder struct {
	Order    *domain.Order
	UpdateFn func(ctx context.Context, order *domain.Order) (*domain.Order, error)
}

type UpdateOrderHandle decorator.CommandHandler[UpdateOrder, interface{}]

type updateOrderHandle struct {
	orderRepo domain.Repository
	// TODO stock rpc
}

func NewUpdateOrderHandle(orderRepo domain.Repository, logger *logrus.Entry, metricClient decorator.MetricsClient) UpdateOrderHandle {
	if orderRepo == nil {
		panic("update order err,nil repo found")
	}
	return decorator.ApplyCommandDecorators[UpdateOrder, interface{}](
		updateOrderHandle{orderRepo: orderRepo},
		logger,
		metricClient,
	)
}

func (u updateOrderHandle) Handle(ctx context.Context, cmd UpdateOrder) (interface{}, error) {
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
