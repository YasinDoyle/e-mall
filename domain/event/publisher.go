package event

import (
	"context"
	"sync"

	"github.com/YasinDoyle/e-mall/utils/log"
)

type Publisher interface {
	Publish(ctx context.Context, event any) error
}

type InProcessPublisher struct{}

var (
	defaultPublisher Publisher = &InProcessPublisher{}
	publisherMu      sync.RWMutex
)

func SetDefaultPublisher(publisher Publisher) {
	publisherMu.Lock()
	defer publisherMu.Unlock()
	if publisher == nil {
		defaultPublisher = &InProcessPublisher{}
		return
	}
	defaultPublisher = publisher
}

func Publish(ctx context.Context, event any) {
	publisherMu.RLock()
	publisher := defaultPublisher
	publisherMu.RUnlock()
	if err := publisher.Publish(ctx, event); err != nil {
		log.LogrusObj.Errorf("publish domain event failed: %v", err)
	}
}

func (p *InProcessPublisher) Publish(ctx context.Context, event any) error {
	if err := handleNotificationEvent(ctx, event); err != nil {
		log.LogrusObj.Errorf("handle notification event failed: %v", err)
	}
	if err := handleProductIndexEvent(ctx, event); err != nil {
		log.LogrusObj.Errorf("handle product index event failed: %v", err)
	}
	return nil
}
