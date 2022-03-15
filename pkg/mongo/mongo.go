package mongo

import (
	"context"

	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"github.com/imkuqin-zw/courier/pkg/config"
	"github.com/qiniu/qmgo"
	options2 "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Base         qmgo.Config
	EnableTracer bool
}

func New() *qmgo.Database {
	cfg := &Config{}
	if err := config.Get("mongo").Scan(&cfg); err != nil {
		logger.Errorf("fault to new mongo client, error: %s", err.Error())
		logger.Fatalf("fault to scan mongo config, error: %s", err.Error())
	}
	ops := make([]options2.ClientOptions, 0)
	var monitor *event.CommandMonitor
	if cfg.EnableTracer {
		monitor = WithTracerMonitor(monitor)

	}
	if monitor != nil {
		ops = append(ops, options2.ClientOptions{
			ClientOptions: options.Client().SetMonitor(monitor),
		})
	}
	client, err := qmgo.NewClient(context.Background(), &cfg.Base, ops...)
	if err != nil {
		logger.Errorf("fault to new mongo client, error: %s", err.Error())
		panic(err)
	}
	if err := client.Ping(3); err != nil {
		logger.Errorf("fault to ping mongo, error: %s", err.Error())
		panic(err)
	}
	return client.Database(cfg.Base.Database)
}

func WithTracerMonitor(baseMonitor *event.CommandMonitor) *event.CommandMonitor {
	monitor := newTracerCommandMonitor()
	if baseMonitor != nil {
		return &event.CommandMonitor{
			Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
				monitor.Started(ctx, startedEvent)
				baseMonitor.Started(ctx, startedEvent)
			},
			Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {
				monitor.Succeeded(ctx, succeededEvent)
				baseMonitor.Succeeded(ctx, succeededEvent)
			},
			Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {
				monitor.Failed(ctx, failedEvent)
				baseMonitor.Failed(ctx, failedEvent)
			},
		}
	}
	return monitor
}

func newTracerCommandMonitor() *event.CommandMonitor {
	return &event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {

		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {

		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {

		},
	}
}
