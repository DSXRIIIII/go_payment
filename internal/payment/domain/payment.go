package payment

import (
	"context"
	"github.com/dsxriiiii/l3x_pay/common/genproto/orderpb"
)

type Processor interface {
	CreatePaymentLink(context.Context, *orderpb.Order) (string, error)
}
