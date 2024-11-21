package decorator

import (
	"context"
	"github.com/sirupsen/logrus"
)

type QueryHandler[C any, R any] interface {
	Handle(ctx context.Context, query C) (R, error)
}

func ApplyQueryDecorators[H, R any](handler QueryHandler[H, R], logger *logrus.Entry, metricsClient MetricsClient) QueryHandler[H, R] {
	return queryLoggingDecorator[H, R]{
		logger: logger,
		base: queryMetricsDecorator[H, R]{
			base:   handler,
			client: metricsClient,
		},
	}
}