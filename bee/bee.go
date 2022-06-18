package bee

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/holodata/fakebee/ytl"
	kafka "github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

const (
	defaultBacklogSize     = 10
	defaultIntervalSeconds = 1
)

type EventWorker struct {
	Topic, OriginChannelId, OriginVideoId, Backend string
	IntervalSeconds, BacklogSize                   int
	KafkaWriter                                    *kafka.Writer
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

	if ew.Backend != "printer" {
		broker := viper.GetString("broker")
		ew.KafkaWriter = &kafka.Writer{
			Addr:                   kafka.TCP(broker),
			AllowAutoTopicCreation: true,
			Topic:                  ew.Topic,
			Balancer:               &kafka.Hash{},
			WriteTimeout:           10 * time.Second,
			Logger:                 kafka.LoggerFunc(logf),
			ErrorLogger:            kafka.LoggerFunc(logf),
		}

		// fmt.Printf("Kafka writer config: %+v", ew.KafkaWriter)
	}

	return ew
}

func logf(msg string, a ...interface{}) {
	fmt.Printf(msg, a...)
	fmt.Println()
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

func WithBackend(backend string) func(*EventWorker) {
	return func(ew *EventWorker) {
		ew.Backend = backend
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
		defer ticker.Stop()

		if eventWorker.Backend != "printer" {
			defer eventWorker.KafkaWriter.Close()
		}

		for eventWorker.BacklogSize > 0 {
			<-ticker.C

			event := ytl.NewEventForTopic(eventWorker.Topic)
			err := eventWorker.produce(event)

			if err != nil {
				fmt.Println(err.Error())
			}

			eventWorker.BacklogSize -= 1
		}
	}()
}

// Create a Kafka message from a Youtube Live event struct
func newKafkaMessage(data interface{}, topic string) (kafka.Message, error) {
	var res []byte
	var key []byte
	var err error

	switch topic {
	case "chats":
		msg := data.(ytl.Chat)
		res, err = json.Marshal(msg)
		if err != nil {
			return kafka.Message{}, err
		}
		key = []byte(fmt.Sprintf("ct-%s", msg.ID))

	case "superchats":
		msg := data.(ytl.SuperChat)
		res, err = json.Marshal(msg)
		if err != nil {
			return kafka.Message{}, err
		}
		key = []byte(fmt.Sprintf("sc-%s", msg.ID))

	case "superstickers":
		msg := data.(ytl.SuperSticker)
		res, err = json.Marshal(msg)
		if err != nil {
			return kafka.Message{}, err
		}
		key = []byte(fmt.Sprintf("stk-%s", msg.ID))

	case "memberships":
		msg := data.(ytl.Membership)
		res, err = json.Marshal(msg)
		if err != nil {
			return kafka.Message{}, err
		}
		key = []byte(fmt.Sprintf("mem-%s", msg.ID))

	case "milestones":
		msg := data.(ytl.Milestone)
		res, err = json.Marshal(msg)
		if err != nil {
			return kafka.Message{}, err
		}
		key = []byte(fmt.Sprintf("mil-%s", msg.ID))

	case "banactions":
		msg := data.(ytl.Ban)
		res, err = json.Marshal(msg)
		if err != nil {
			return kafka.Message{}, err
		}
		key = []byte(fmt.Sprintf("ban-%s", msg.ID))

	case "deleactions":
		msg := data.(ytl.Deletion)
		res, err = json.Marshal(msg)
		if err != nil {
			return kafka.Message{}, err
		}
		key = []byte(fmt.Sprintf("del-%s", msg.ID))
	}

	return kafka.Message{
		Key:   key,
		Value: res,
	}, nil
}

func (ew *EventWorker) produce(msg interface{}) error {
	switch ew.Backend {
	case "printer":
		data, _ := json.Marshal(msg)
		fmt.Printf("Printing to topic [%s] with message: %+s\n", ew.Topic, data)
	case "kafka":
		// only other backend atm is kafka
		payload, err := newKafkaMessage(msg, ew.Topic)
		if err != nil {
			return err
		}

		err = ew.KafkaWriter.WriteMessages(context.Background(), payload)
		if err != nil {
			fmt.Println(err)
		}
	default:
		return fmt.Errorf("unkown backend: %s", ew.Backend)
	}

	return nil
}
