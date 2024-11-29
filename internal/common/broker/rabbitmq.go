package broker

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"time"
)

const (
	DLX                = "dlx"
	DLQ                = "dlq"
	amqpRetryHeaderKey = "x-retry-count"
)

var (
	maxRetryCount = viper.GetInt64("rabbitmq.max-retry")
)

func Connect(user, password, host, port string) (*amqp.Channel, func() error) {
	address := fmt.Sprintf("amqp://%s:%s@%s:%s", user, password, host, port)
	conn, err := amqp.Dial(address)
	if err != nil {
		logrus.Fatal(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatal(err)
	}
	err = ch.ExchangeDeclare(EventOrderCreated, "direct", true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	err = ch.ExchangeDeclare(EventOrderPaid, "fanout", true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	if err = createDLX(ch); err != nil {
		logrus.Fatal(err)
	}
	return ch, conn.Close
}

func createDLX(ch *amqp.Channel) error {
	q, err := ch.QueueDeclare("share_queue", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = ch.ExchangeDeclare(DLX, "fanout", true, false, false, false, nil)
	err = ch.QueueBind(q.Name, "", DLX, false, nil)
	if err != nil {
		return err
	}
	_, err = ch.QueueDeclare(DLQ, true, false, false, false, nil)
	return err
}

func HandleRetry(ctx context.Context, ch *amqp.Channel, d *amqp.Delivery) error {
	logrus.Info("handle retry_max-retry-count", maxRetryCount)
	if d.Headers == nil {
		d.Headers = amqp.Table{}
	}
	retryCount, ok := d.Headers[amqpRetryHeaderKey].(int64)
	if !ok {
		retryCount = 0
	}
	retryCount++
	d.Headers[amqpRetryHeaderKey] = retryCount

	if retryCount >= maxRetryCount {
		logrus.Infof("moving message %s to dlq", d.MessageId)
		return ch.PublishWithContext(ctx, "", DLQ, false, false, amqp.Publishing{
			Headers:      d.Headers,
			ContentType:  "application/json",
			Body:         d.Body,
			DeliveryMode: amqp.Persistent,
		})
	}
	logrus.Infof("retring message %s, count=%d", d.MessageId, retryCount)
	time.Sleep(time.Second * time.Duration(retryCount))
	return ch.PublishWithContext(ctx, d.Exchange, d.RoutingKey, false, false, amqp.Publishing{
		Headers:      d.Headers,
		ContentType:  "application/json",
		Body:         d.Body,
		DeliveryMode: amqp.Persistent,
	})
}

type RabbitMQHeaderCarrier map[string]interface{}

func (r RabbitMQHeaderCarrier) Get(key string) string {
	value, ok := r[key]
	if !ok {
		return ""
	}
	return value.(string)
}

func (r RabbitMQHeaderCarrier) Set(key string, value string) {
	r[key] = value
}

func (r RabbitMQHeaderCarrier) Keys() []string {
	keys := make([]string, len(r))
	i := 0
	for key := range r {
		keys[i] = key
		i++
	}
	return keys
}

// InjectRabbitMQHeaders
// InjectRabbitMQHeaders函数用于将相关的追踪上下文信息注入到RabbitMQ的消息头中
// 它接受一个上下文ctx作为参数，并返回一个包含注入后的RabbitMQ消息头的映射（map）类型数据结构
func InjectRabbitMQHeaders(ctx context.Context) map[string]interface{} {
	// 创建一个空的RabbitMQHeaderCarrier类型的变量carrier，用于存储即将注入的消息头信息
	// RabbitMQHeaderCarrier应该是一个自定义的类型，用于表示RabbitMQ消息头的容器
	carrier := make(RabbitMQHeaderCarrier)

	// 使用OpenTelemetry的全局文本映射传播器（通过otel.GetTextMapPropagator()获取）
	// 将当前上下文ctx中的追踪相关信息注入到刚才创建的carrier中
	// 这样carrier就会包含与追踪上下文相关的键值对，以便后续添加到RabbitMQ消息头中
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	// 将包含了注入追踪信息的carrier作为函数的最终返回值返回
	// 这个返回值就是一个可以直接用于设置RabbitMQ消息头的映射类型数据结构
	return carrier
}

// ExtractRabbitMQHeaders
// ExtractRabbitMQHeaders函数用于从RabbitMQ的消息头（以映射类型表示）中提取相关的追踪上下文信息
// 并将提取到的信息更新到给定的上下文ctx中，最后返回更新后的上下文ctx
func ExtractRabbitMQHeaders(ctx context.Context, headers map[string]interface{}) context.Context {
	// 使用OpenTelemetry的全局文本映射传播器（通过otel.GetTextMapPropagator()获取）
	// 从给定的RabbitMQ消息头（以RabbitMQHeaderCarrier类型包装后的形式，这里通过将headers转换为RabbitMQHeaderCarrier类型）
	// 提取相关的追踪上下文信息，并将这些信息更新到当前给定的上下文ctx中
	return otel.GetTextMapPropagator().Extract(ctx, RabbitMQHeaderCarrier(headers))
}
