package bee

import (
	"fmt"
	"sync"
	"time"

	"github.com/holodata/fakebee/ytl"
	"github.com/pioz/faker"
)

const (
	defaultBacklogSize     = 10
	defaultIntervalSeconds = 1
)

type EventWorker struct {
	Topic, OriginChannelId, OriginVideoId string
	IntervalSeconds, BacklogSize          int
}

func NewEventWorker(originChannelId, originVideoId, topic string) *EventWorker {
	return &EventWorker{
		topic, originChannelId, originVideoId, defaultIntervalSeconds, defaultBacklogSize}
}

// Starts the control loop implemented as a goroutine
func (eventWorker *EventWorker) Run(wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Duration(eventWorker.IntervalSeconds) * time.Second)

		for eventWorker.BacklogSize > 0 {
			<-ticker.C

			chat := ytl.Chat{}
			err := faker.Build(&chat)
			if err != nil {
				panic(err)
			}

			fmt.Printf("Producing to topic [%s] with message: %v+\n", eventWorker.Topic, chat)
			eventWorker.BacklogSize -= 1
		}
		ticker.Stop()
	}()
}
