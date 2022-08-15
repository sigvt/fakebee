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
	defaultBacklogSize     = 100
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

			err := eventWorker.produce()

			if err != nil {
				fmt.Println(err.Error())
			}

			eventWorker.BacklogSize -= 1
		}
	}()
}

// Creates and encodes an Event for a given topic
// Returns the encoded `event` and `key`
func EncodeEvent(topic string) (payload []byte, key []byte, err error) {
	switch topic {
	case "chats":
		chat := ytl.ChatFactory()
		payload, err = json.Marshal(chat)
		key = []byte(fmt.Sprintf("ct-%s", chat.ID))
	case "superchats":
		superchat := ytl.SuperChatFactory()
		payload, err = json.Marshal(superchat)
		key = []byte(fmt.Sprintf("sc-%s", superchat.ID))
	case "superstickers":
		supersticker := ytl.SuperStickerFactory()
		payload, err = json.Marshal(supersticker)
		key = []byte(fmt.Sprintf("stk-%s", supersticker.ID))
	case "memberships":
		membership := ytl.MembershipFactory()
		payload, err = json.Marshal(membership)
		key = []byte(fmt.Sprintf("mem-%s", membership.ID))
	case "milestones":
		milestone := ytl.MilestoneFactory()
		payload, err = json.Marshal(milestone)
		key = []byte(fmt.Sprintf("mil-%s", milestone.ID))
	case "banactions":
		banaction := ytl.BanFactory()
		payload, err = json.Marshal(banaction)
		key = []byte(fmt.Sprintf("ban-%s", banaction.ID))
	case "deleteactions":
		deleteaction := ytl.DeletionFactory()
		payload, err = json.Marshal(deleteaction)
		key = []byte(fmt.Sprintf("del-%s", deleteaction.ID))
	default:
		payload, key, err = nil, []byte(""), fmt.Errorf("unknown topic: %s", topic)
	}

	return
}

// Create a Kafka message from a Youtube Live event struct
func createKafkaMsg(topic string) (kafka.Message, error) {
	var res, key []byte
	var err error

	res, key, err = EncodeEvent(topic)
	if err != nil {
		return kafka.Message{}, err
	}

	return kafka.Message{
		Key:   key,
		Value: res,
	}, nil
}

func (ew *EventWorker) produce() error {
	switch ew.Backend {
	case "printer":
		data, _, err := EncodeEvent(ew.Topic)
		if err != nil {
			return err
		}
		fmt.Printf("Printing to topic [%s] with message: %+s\n", ew.Topic, data)
	case "kafka":
		// only other backend atm is kafka
		payload, err := createKafkaMsg(ew.Topic)
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
