package query

import (
	"context"
	"github.com/dsxriiiii/l3x_pay/common/decorator"
	"github.com/dsxriiiii/l3x_pay/common/handle/redis"
	domain "github.com/dsxriiiii/l3x_pay/stock/domain/stock"
	"github.com/dsxriiiii/l3x_pay/stock/entity"
	integration "github.com/dsxriiiii/l3x_pay/stock/infrastructure/integration"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

const redisLockPrefix = "check_stock_"

type CheckIfItemsInStock struct {
	Items []*entity.ItemWithQuantity
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*entity.Item]

type checkIfItemsInStockHandler struct {
	stockRepo domain.Repository
	stripeAPI *integration.StripeAPI
}

func NewCheckIfItemsInStockHandler(
	stockRepo domain.Repository,
	stripeAPI *integration.StripeAPI,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CheckIfItemsInStockHandler {
	if stockRepo == nil {
		panic("nil stockRepo")
	}
	if stripeAPI == nil {
		panic("nil stripeAPI")
	}
	return decorator.ApplyQueryDecorators[CheckIfItemsInStock, []*entity.Item](
		checkIfItemsInStockHandler{stockRepo: stockRepo, stripeAPI: stripeAPI},
		logger,
		metricClient,
	)
}

//var stub = map[string]string{
//	"1": "price_1QPb4qE3H882BCCxFAXzD7WZ",
//}

func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*entity.Item, error) {
	if err := lock(ctx, getLockKey(query)); err != nil {
		return nil, err
	}
	defer func() {
		if err := unlock(ctx, getLockKey(query)); err != nil {
			logrus.Warnf("redis unlock fail,err = %v", err)
		}
	}()
	var res []*entity.Item
	for _, i := range query.Items {
		priceID, err := h.stripeAPI.GetPriceByProductID(ctx, i.ID)
		if err != nil || priceID == "" {
			return nil, err
		}
		res = append(res, &entity.Item{
			ID:       i.ID,
			Quantity: i.Quantity,
			PriceID:  priceID,
		})
	}
	// TODO 扣库存
	if err := h.checkStock(ctx, query.Items); err != nil {
		return nil, err
	}
	return res, nil
}

//func getStubPriceID(id string) string {
//	priceID, ok := stub[id]
//	if !ok {
//		priceID = stub["1"]
//	}
//	return priceID
//}

func (h checkIfItemsInStockHandler) checkStock(ctx context.Context, query []*entity.ItemWithQuantity) error {
	var ids []string
	for _, i := range query {
		ids = append(ids, i.ID)
	}
	records, err := h.stockRepo.GetStock(ctx, ids)
	if err != nil {
		return err
	}
	idQuantityMap := make(map[string]int32)
	for _, r := range records {
		idQuantityMap[r.ID] += r.Quantity
	}
	var (
		ok       = true
		failedOn []struct {
			ID   string
			Want int32
			Have int32
		}
	)
	for _, item := range query {
		if item.Quantity > idQuantityMap[item.ID] {
			ok = false
			failedOn = append(failedOn, struct {
				ID   string
				Want int32
				Have int32
			}{ID: item.ID, Want: item.Quantity, Have: idQuantityMap[item.ID]})
		}
	}
	if ok {
		return h.stockRepo.UpdateStock(ctx, query, func(
			ctx context.Context,
			existing []*entity.ItemWithQuantity,
			query []*entity.ItemWithQuantity,
		) ([]*entity.ItemWithQuantity, error) {
			var newItems []*entity.ItemWithQuantity
			for _, e := range existing {
				for _, q := range query {
					if e.ID == q.ID {
						newItems = append(newItems, &entity.ItemWithQuantity{
							ID:       e.ID,
							Quantity: e.Quantity - q.Quantity,
						})
					}
				}
			}
			return newItems, nil
		})
	}
	return domain.ExceedStockError{FailedOn: failedOn}
}

func getLockKey(query CheckIfItemsInStock) string {
	var ids []string
	for _, i := range query.Items {
		ids = append(ids, i.ID)
	}
	return redisLockPrefix + strings.Join(ids, "_")
}

func unlock(ctx context.Context, key string) error {
	return redis.Del(ctx, redis.LocalClient(), key)
}

func lock(ctx context.Context, key string) error {
	return redis.SetNX(ctx, redis.LocalClient(), key, "lock", 5*time.Minute)
}
