package bee

import (
	"fmt"
	"sync"
	"time"
)

const (
	defaultBacklogSize     = 10
	defaultIntervalSeconds = 1
)

type EventWorker struct {
	Topic, OriginChannelId, OriginVideoId string
	IntervalSeconds, BacklogSize          int
}

func New(topic, originChannelId, originVideoId string) *EventWorker {
	return &EventWorker{
		topic, originChannelId, originVideoId, defaultIntervalSeconds, defaultBacklogSize}
}

// Starts the control loop implemented as a goroutine
func (eventWorker *EventWorker) Run(wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Duration(eventWorker.IntervalSeconds) * time.Second)

		for eventWorker.BacklogSize > 0 {
			t := <-ticker.C

			fmt.Println("Doing work for topic", eventWorker.Topic, "got tic:", t)
			eventWorker.BacklogSize -= 1
		}
		ticker.Stop()
	}()
}
