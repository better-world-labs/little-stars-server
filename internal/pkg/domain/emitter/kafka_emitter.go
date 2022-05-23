package emitter

import (
	"aed-api-server/internal/pkg/async"
	"aed-api-server/internal/pkg/domain/config"
	"aed-api-server/internal/pkg/utils"
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"sync"
	"time"
)

const (
	_eventTypeTag = "event-type"
)

type (
	DeadMessage struct {
		OriginTopic string `json:"originTopic"`
		OriginType  string `json:"originType"`
		ErrorMsg    string `json:"errorMsg"`
		Timestamp   int64  `json:"timestamp"`
		Payload     []byte `json:"payload"`
	}

	HandlerKeeper struct {
		slice   []DomainEventHandler
		decoder Decoder
	}

	KafkaEmitter struct {
		ctx             context.Context
		cancel          context.CancelFunc
		producer        *kafka.Producer
		consumer        *kafka.Consumer
		handlers        *HandlerRegistry
		group           sync.WaitGroup
		closed          bool
		topic           string
		deadLetterTopic string
	}
)

func (receiver *HandlerKeeper) Handlers() []DomainEventHandler {
	return receiver.slice
}

func NewKafkaEmitter(ctx context.Context, c *config.DomainEventConfig) (Emitter, error) {
	withCancel, cancelFunc := context.WithCancel(ctx)
	e := &KafkaEmitter{
		handlers:        NewHandlerRegistry(),
		ctx:             withCancel,
		cancel:          cancelFunc,
		topic:           c.Topic,
		deadLetterTopic: c.DeadLetterTopic,
	}

	err := e.initProducer(c)
	if err != nil {
		return nil, err
	}

	err = e.initConsumer(c)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (e *KafkaEmitter) Start() {
	e.group = sync.WaitGroup{}
	e.group.Add(2)
	go func() {
		e.startReadLoop()
		e.group.Done()
	}()
	go func() {
		e.startReportLoop()
		e.group.Done()
	}()
}

func (e *KafkaEmitter) Close() {
	e.cancel()
	e.group.Wait()
}

func (e *KafkaEmitter) initProducer(config *config.DomainEventConfig) error {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.Server,
		"acks":              "all",
	})

	if err != nil {
		return err
	}

	e.producer = producer
	return nil
}

func (e *KafkaEmitter) initConsumer(config *config.DomainEventConfig) error {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  config.Server,
		"group.id":           config.GroupId,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	if err != nil {
		return err
	}

	err = consumer.Subscribe(config.Topic, nil)
	if err != nil {
		return err
	}

	e.consumer = consumer
	return nil
}

func (e *KafkaEmitter) startReportLoop() {
	for e := range e.producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Printf("[emitter.KafkaEmitter] Delivery failed: %v\n", ev.TopicPartition)
			} else {
				log.Printf("[emitter.KafkaEmitter] Delivered message to %v\n", ev.TopicPartition)
			}
		}
	}
}

func (e *KafkaEmitter) startReadLoop() {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("[emitter.KafkaEmitter] ReedLoop error: %v\n %s", err, utils.PanicTrace(2))
		}

		if e.closed {
			e.producer.Close()
			log.Printf("[emitter.KafkaEmitter] producer closed\n")

			err := e.consumer.Close()
			if err != nil {
				log.Printf("[emitter.KafkaEmitter] consumer close error: %v\n", err)
			}

			log.Printf("[emitter.KafkaEmitter] consumer closed\n")
		} else {
			e.startReadLoop()
		}
	}()

	for {
		select {
		case <-e.ctx.Done():
			e.closed = true
			return
		default:
			message, err := e.consumer.ReadMessage(5 * time.Second)
			if err == nil {
				e.kafkaMessageDeal(message)
			} else {
				switch err.(type) {
				case kafka.Error:
					kafkaError := err.(kafka.Error)
					if kafkaError.Code() == -185 { // 超时不打log
					}
					continue
				}

				log.Printf("[emitter.KafkaEmitter] ReadMessage error: %v\n", err)
			}
		}
	}
}

func (e *KafkaEmitter) Emit(events ...DomainEvent) error {
	EmitOne := func(evt DomainEvent) error {
		msgValue, err := evt.Encode()
		if err != nil {
			return err
		}
		eventType := GetStructType(evt)
		kafkaMsg := &kafka.Message{
			Headers: []kafka.Header{
				{
					Key:   _eventTypeTag,
					Value: []byte(eventType),
				},
			},
			Value:          msgValue,
			TopicPartition: kafka.TopicPartition{Topic: &e.topic, Partition: kafka.PartitionAny},
		}

		log.Printf("EmitEvent: type=%s, message=%s\n", eventType, msgValue)
		return e.kafkaPublishMessage(kafkaMsg)
	}

	for i := range events {
		err := EmitOne(events[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *KafkaEmitter) kafkaPublishMessage(msg *kafka.Message) error {
	return e.producer.Produce(msg, nil)
}

func (e *KafkaEmitter) kafkaMessageDeal(msg *kafka.Message) {
	messageType, exists := getHeader(msg, _eventTypeTag)
	if !exists {
		log.Printf("[emitter.KafkaEmitter] eventType not found\n")
		e.commitMessage(msg)
		return
	}

	keeper, exists := e.handlers.Get(messageType)
	if !exists {
		e.commitMessage(msg)
		return
	}

	if keeper.decoder == nil {
		e.commitMessage(msg)
		log.Printf("[emitter.KafkaEmitter] decoder for type %s not found\n", messageType)
		return
	}

	var futures []*async.Future
	domainEvent, err := keeper.decoder.Decode(msg.Value)
	if err != nil {
		log.Printf("[emitter.KafkaEmitter] message decode error: %v\n", err) // Decode error
		e.commitMessage(msg)
		return
	}

	for _, k := range keeper.slice {
		k := k //创建副本，让每个 goroutine 之间不共享
		future := async.RunTask(func() (interface{}, error) {
			err := k(domainEvent)
			return nil, err
		})

		if err != nil {
			e.doDeadLetterQueue(msg, messageType, domainEvent, err)
			return
		}

		futures = append(futures, future)
	}

	err = async.CompositeFutureAll(futures)
	if err != nil {
		e.doDeadLetterQueue(msg, messageType, domainEvent, err)
		return
	}

	e.commitMessage(msg)
}

func (e *KafkaEmitter) doDeadLetterQueue(msg *kafka.Message, messageType string, evt DomainEvent, err error) {
	err = e.handleKafkaMessageFailed(msg, messageType, evt, err)
	if err != nil {
		log.Printf("[emitter.KafkaEmitter] handleKafkaMessageFailed: topic=%s, mesageType=%s, value=%s\n", *msg.TopicPartition.Topic, messageType, msg.Value)
		return
	}

	e.commitMessage(msg)
}

func (e *KafkaEmitter) commitMessage(message *kafka.Message) {
	_, err := e.consumer.CommitMessage(message)
	if err != nil {
		log.Printf("[emitter.KafkaEmitter] message value = %s commit failed: %v\n", message.Value, err)
	}
}

func getHeader(msg *kafka.Message, key string) (string, bool) {
	if key == "" {
		return "", false
	}

	for i := range msg.Headers {
		if msg.Headers[i].Key == key {
			return string(msg.Headers[i].Value), true
		}
	}

	return "", false
}

func (e *KafkaEmitter) getKeeper(evtType string) (*HandlerKeeper, bool) {
	keeper, exists := e.handlers.Get(evtType)
	return keeper, exists
}

func (e *KafkaEmitter) On(evt DomainEvent, handlers ...DomainEventHandler) {
	for _, h := range handlers {
		e.handlers.Register(evt, h)
	}
}

func (e *KafkaEmitter) Off(evt DomainEvent, handlers ...DomainEventHandler) {
	for _, h := range handlers {
		e.handlers.Delete(evt, h)
	}
}

func (e *KafkaEmitter) handleKafkaMessageFailed(msg *kafka.Message, messageType string, evt DomainEvent, evtError error) error {
	encode, err := evt.Encode()
	if err != nil {
		return err
	}

	var dead = DeadMessage{
		OriginTopic: *msg.TopicPartition.Topic,
		OriginType:  messageType,
		ErrorMsg:    evtError.Error(),
		Timestamp:   time.Now().UnixMilli(),
		Payload:     encode,
	}

	bytes, err := json.Marshal(&dead)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Value:          bytes,
		TopicPartition: kafka.TopicPartition{Topic: &e.deadLetterTopic, Partition: kafka.PartitionAny},
	}

	err = e.producer.Produce(&message, nil)
	if err != nil {
		return err
	}

	return nil
}
