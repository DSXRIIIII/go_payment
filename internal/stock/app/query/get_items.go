package query

import (
	"context"
	"github.com/dsxriiiii/l3x_pay/common/decorator"
	"github.com/dsxriiiii/l3x_pay/common/genproto/orderpb"
	"github.com/sirupsen/logrus"
	domain "github/dsxriiiii/l3x_pay/stock/domain/stock"
)

type GetItems struct {
	ItemIDs []string
}

type GetItemsHandler decorator.QueryHandler[GetItems, []*orderpb.Item]

type getItemsHandler struct {
	stockRepo domain.Repository
}

func NewGetItemsHandler(stockRepo domain.Repository, logger *logrus.Entry, metricClient decorator.MetricsClient) GetItemsHandler {
	if stockRepo == nil {
		panic("nil stockRepo")
	}
	return decorator.ApplyQueryDecorators[GetItems, []*orderpb.Item](
		getItemsHandler{stockRepo: stockRepo},
		logger,
		metricClient,
	)
}

func (g getItemsHandler) Handle(ctx context.Context, query GetItems) ([]*orderpb.Item, error) {
	items, err := g.stockRepo.GetItems(ctx, query.ItemIDs)
	if err != nil {
		return nil, err
	}
	return items, nil
}
