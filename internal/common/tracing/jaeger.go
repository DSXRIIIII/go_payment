package tracing

import (
	"context"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("default_tracer")

func InitJaegerProvider(jaegerURL, serviceName string) (func(ctx context.Context) error, error) {
	if jaegerURL == "" {
		panic("empty jaeger url")
	}
	// 根据传入的服务名称重新设置全局追踪器tracer
	tracer = otel.Tracer(serviceName)
	// 创建一个新的Jaeger exporter，指定了收集器端点的URL
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerURL)))
	if err != nil {
		return nil, err
	}
	// 创建一个新的追踪提供程序TracerProvider
	// 使用了批处理模式（WithBatcher）将追踪数据发送到之前创建的Jaeger exporter
	// 并设置了资源信息，包括服务名称等，这里使用了无模式（NewSchemaless）的资源创建方式
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	// 将创建好的追踪提供程序设置为全局的追踪提供程序
	otel.SetTracerProvider(tp)

	// 创建一个新的b3传播器，指定了注入编码为b3.B3MultipleHeader，用于在请求中传播追踪上下文
	b3Propagator := b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader))

	// 创建一个复合文本映射传播器，包含了TraceContext、Baggage和b3传播器
	// 用于在不同的服务间传播追踪相关的上下文信息
	p := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}, b3Propagator,
	)

	// 将创建好的复合文本映射传播器设置为全局的传播器
	otel.SetTextMapPropagator(p)
	return tp.Shutdown, nil
}

func Start(ctx context.Context, name string) (context.Context, trace.Span) {
	return tracer.Start(ctx, name)
}

func TraceID(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	return spanCtx.TraceID().String()
}
