package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dsxriiiii/l3x_pay/common/broker"
	"github.com/dsxriiiii/l3x_pay/common/decorator"
	"github.com/dsxriiiii/l3x_pay/order/app/query"
	"github.com/dsxriiiii/l3x_pay/order/convertor"
	domain "github.com/dsxriiiii/l3x_pay/order/domain/order"
	"github.com/dsxriiiii/l3x_pay/order/entity"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type CreateOrder struct {
	CustomerID string
	Items      []*entity.ItemWithQuantity
}

type CreateOrderResult struct {
	OrderID string
}

type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]

type createOrderHandler struct {
	orderRepo domain.Repository
	stockGRPC query.StockService
	channel   *amqp.Channel
}

func NewCreateOrderHandler(
	orderRepo domain.Repository,
	stockGRPC query.StockService,
	channel *amqp.Channel,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CreateOrderHandler {
	if orderRepo == nil {
		panic("create order err, orderRp is nil")
	}
	if stockGRPC == nil {
		panic("nil stockGRPC")
	}
	if channel == nil {
		panic("nil channel ")
	}
	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
		createOrderHandler{orderRepo: orderRepo, stockGRPC: stockGRPC, channel: channel},
		logger,
		metricClient,
	)
}

func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
	q, err := c.channel.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	t := otel.Tracer("rabbitmq")
	ctx, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.publish", q.Name))
	defer span.End()
	validItems, err := c.validate(ctx, cmd.Items)
	if err != nil {
		return nil, err
	}
	o, err := c.orderRepo.Create(ctx, &domain.Order{
		CustomerID: cmd.CustomerID,
		Items:      validItems,
	})
	if err != nil {
		return nil, err
	}

	// change created order message to json format
	marshalledOrder, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	header := broker.InjectRabbitMQHeaders(ctx)
	err = c.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		// report json message
		Body:    marshalledOrder,
		Headers: header,
	})
	return &CreateOrderResult{OrderID: o.ID}, nil
}

func (c createOrderHandler) validate(ctx context.Context, items []*entity.ItemWithQuantity) ([]*entity.Item, error) {
	if len(items) == 0 {
		return nil, errors.New("must have at least one item")
	}
	items = packItems(items)
	resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, convertor.NewItemWithQuantityConvertor().EntitiesToProtos(items))
	if err != nil {
		return nil, err
	}
	return convertor.NewItemConvertor().ProtosToEntities(resp.Items), nil
}

func packItems(items []*entity.ItemWithQuantity) []*entity.ItemWithQuantity {
	merged := make(map[string]int32)
	for _, item := range items {
		merged[item.ID] += item.Quantity
	}
	var res []*entity.ItemWithQuantity
	for id, quantity := range merged {
		res = append(res, &entity.ItemWithQuantity{
			ID:       id,
			Quantity: quantity,
		})
	}
	return res
}
