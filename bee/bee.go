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

// Create a new EventWorker with some options
// `WithTopic` and `WithOrigin` are required
func NewEventWorker(options ...func(*EventWorker)) *EventWorker {
	ew := &EventWorker{}

	ew.IntervalSeconds = defaultIntervalSeconds
	ew.BacklogSize = defaultBacklogSize

	for _, o := range options {
		o(ew)
	}

	return ew
}

func WithOrigin(origin ytl.Origin) func(*EventWorker) {
	return func(ew *EventWorker) {
		ew.OriginVideoId = origin.VideoId
		ew.OriginChannelId = origin.ChannelId
	}
}

func WithTopic(topic string) func(*EventWorker) {
	return func(ew *EventWorker) {
		ew.Topic = topic
	}
}

func WithBacklogSize(size int) func(*EventWorker) {
	return func(ew *EventWorker) {
		ew.BacklogSize = size
	}
}

func WithInterval(interval time.Duration) func(*EventWorker) {
	return func(ew *EventWorker) {
		ew.IntervalSeconds = int(interval.Seconds())
	}
}

// Starts the control loop implemented as a goroutine
func (eventWorker *EventWorker) Run(wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Duration(eventWorker.IntervalSeconds) * time.Second)

		for eventWorker.BacklogSize > 0 {
			<-ticker.C

			event := ytl.Chat{} // TODO: Replace with switch by topic type
			err := faker.Build(&event)
			if err != nil {
				panic(err)
			}

			produce(eventWorker.Topic, &event)

			eventWorker.BacklogSize -= 1
		}
		ticker.Stop()
	}()
}

func produce[T ytl.Event](topic string, msg *T) {
	fmt.Printf("Producing to topic [%s] with message: %+v\n", topic, msg)
}
