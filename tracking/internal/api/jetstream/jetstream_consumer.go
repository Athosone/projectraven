package api

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type MessageHandler func(ctx context.Context, payload []byte, messageId string) error

type JSServerConfig struct {
	Broker   string
	ClientID string
}

type JSServer struct {
	jsctx  nats.JetStreamContext
	rwLock sync.RWMutex
	logger *zap.SugaredLogger

	handlers map[string][]MessageHandler
}

func NewJSServer(logger *zap.SugaredLogger, jsctx nats.JetStreamContext) (*JSServer, error) {
	listener := &JSServer{jsctx: jsctx, handlers: make(map[string][]MessageHandler), logger: logger.With("component", "mqtt")}
	return listener, nil
}

func (l *JSServer) Stop(ctx context.Context) {
	l.rwLock.Lock()
	defer l.rwLock.Unlock()
	// for topic, handlers := range l.handlers {
	// 	for _, handler := range handlers {
	// 		handler(ctx, nil, "")
	// 	}
	// 	delete(l.handlers, topic)
	// }
	// l.jsctx.Close()
}

func (l *JSServer) Subscribe(ctx context.Context, topic string, name string, handler MessageHandler) error {
	l.rwLock.Lock()
	defer l.rwLock.Unlock()
	l.handlers[topic] = append(l.handlers[topic], handler)
	jsctx := l.jsctx
	fetcher, err := jsctx.PullSubscribe(topic, name,
		nats.AckExplicit(),
		nats.ManualAck())
	if err != nil {
		return fmt.Errorf("error subscribing to topic %s: %w", topic, err)
	}

	go func() {

		if err != nil {
			panic(err)
		}
		// pull messages
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			msg, err := fetcher.Fetch(1, nats.MaxWait(1*time.Second))
			if err != nil {
				if err == context.Canceled {
					return
				}
				// if err is deadline exceeded, then no messages were received
				if err == context.DeadlineExceeded {
					continue
				}
				panic(err)
			}
			for _, m := range msg {
				if err := handler(ctx, m.Data, fmt.Sprint(m.Header.Get("message_id"))); err != nil {
					l.logger.Errorf("Error handling message: %v", err)
					m.Nak()
					continue
				}
				m.AckSync()
			}
		}
	}()
	return nil
}
