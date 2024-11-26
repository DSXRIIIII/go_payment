package command

import (
	"context"
	"github.com/dsxriiiii/l3x_pay/common/genproto/orderpb"
)

type OrderService interface {
	UpdateOrder(ctx context.Context, order *orderpb.Order) error
}
