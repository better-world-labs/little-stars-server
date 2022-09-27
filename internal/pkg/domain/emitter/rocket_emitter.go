package emitter

import (
	"aed-api-server/internal/pkg/async"
	"aed-api-server/internal/pkg/crypto"
	"aed-api-server/internal/pkg/domain/config"
	"aed-api-server/internal/pkg/utils"
	"context"
	rocket "github.com/aliyunmq/mq-http-go-sdk"
	"github.com/gogap/errors"
	log "github.com/sirupsen/logrus"
	"sort"
	"strings"
	"sync"
	"time"
)

func NewRocketEmitter(conf interface{}) (Emitter, error) {
	c, ok := conf.(*config.RocketConf)
	if !ok {
		panic("NewKafkaEmitter params err")
	}

	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	client := rocket.NewAliyunMQClient(c.EndPoint, c.AccessKey, c.SecretKey, "")
	rocketEmitter := RocketEmitter{
		MQProducer: client.GetProducer(c.InstanceId, c.Topic),
		conf:       c,
		client:     client,
		cancel:     cancelFunc,
		ctx:        cancelCtx,
		handlers:   NewHandlerRegistry(),
	}
	return &rocketEmitter, nil
}

type RocketEmitter struct {
	rocket.MQProducer
	conf          *config.RocketConf
	client        rocket.MQClient
	ctx           context.Context
	cancel        context.CancelFunc
	group         *sync.WaitGroup
	handlers      *HandlerRegistry
	messageTagStr string
}

func (e *RocketEmitter) Start() {
	e.group = &sync.WaitGroup{}
	e.group.Add(1)

	go e.loop()
}

func (e *RocketEmitter) Close() {
	if e.cancel != nil {
		e.cancel()
	}
	if e.group != nil {
		e.group.Wait()
	}
}
func (e *RocketEmitter) Emit(events ...DomainEvent) error {
	return e.DelayEmit(0, events...)
}

func (e *RocketEmitter) DelayEmit(duration time.Duration, events ...DomainEvent) error {
	for i := range events {
		evt := events[i]
		msgBody, err := evt.Encode()
		if err != nil {
			log.Errorf("encode event:%v err:%v", evt, err)
			return err
		}

		eventType := GetStructType(evt)
		msgTag := e.getStructTypeHash(eventType) //使用消息的类型Hash做messageTag
		msg := rocket.PublishMessageRequest{
			MessageBody: string(msgBody),
			MessageTag:  msgTag,
			Properties: map[string]string{
				_traceIdTag:   utils.GetTraceId(),
				_eventTypeTag: eventType,
			},
		}
		if duration > 0 {
			delayTime := time.Now().Add(duration)
			msg.StartDeliverTime = delayTime.UnixMilli()
		}

		ret, err := e.PublishMessage(msg)
		if err != nil {
			log.Errorf("send event:%v err:%v", evt, err)
			return err
		}
		log.Infof("message|send msg:Id=%s, type=%s, tag=%s, content=%s", ret.MessageId, eventType, msgTag, msgBody)
	}
	return nil
}

func (e *RocketEmitter) On(evt DomainEvent, handlers ...DomainEventHandler) Emitter {
	for _, h := range handlers {
		e.handlers.Register(evt, h)
	}
	e.messageTagStr = e.genEvtListHash(e.handlers.GetEventTypes())
	return e
}

func (e *RocketEmitter) Off(evt DomainEvent, handlers ...DomainEventHandler) Emitter {
	for _, h := range handlers {
		e.handlers.Delete(evt, h)
	}
	e.messageTagStr = e.genEvtListHash(e.handlers.GetEventTypes())
	return e
}

func (e *RocketEmitter) genEvtListHash(eventsName []string) string {
	arr := make([]string, 0, len(eventsName))
	for i := range eventsName {
		arr = append(arr, e.getStructTypeHash(eventsName[i]))
	}

	sort.Sort(sort.StringSlice(arr))

	return strings.Join(arr, "||")
}

func (e *RocketEmitter) getNewConsumer() rocket.MQConsumer {
	return e.client.GetConsumer(e.conf.InstanceId, e.conf.Topic, e.conf.GroupId, e.messageTagStr)
}

func (e *RocketEmitter) loop() {
	for {
		mqConsumer := e.getNewConsumer()

		endChan := make(chan int)
		respChan := make(chan rocket.ConsumeMessageResponse)
		errChan := make(chan error)

		go func() {
			select {
			case resp := <-respChan:
				e.dealMsgRes(mqConsumer, resp)
			case err := <-errChan:
				if !strings.Contains(err.(errors.ErrCode).Error(), "MessageNotExist") {
					log.Errorf("consume msg err:%v", err)
					time.Sleep(time.Duration(3) * time.Second)
				}
			case <-time.After(35 * time.Second):
				log.Warn("Timeout of consumer message")
			}
			endChan <- 1
		}()

		// 长轮询消费消息，网络超时时间默认为35s。
		// 长轮询表示如果Topic没有消息，则客户端请求会在服务端挂起3s，3s内如果有消息可以消费则立即返回响应。
		mqConsumer.ConsumeMessage(respChan, errChan,
			3, // 一次最多消费3条（最多可设置为16条）。
			3, // 长轮询时间3s（最多可设置为30s）。
		)

		select {
		case <-e.ctx.Done():
			log.Warn("rocket emitter closing")
			e.group.Done()
			return
		case <-endChan:
		}
	}
}

func (e *RocketEmitter) dealMsgRes(mqConsumer rocket.MQConsumer, msgRes rocket.ConsumeMessageResponse) {
	var handles []string
	for i := range msgRes.Messages {
		msg := msgRes.Messages[i]

		err := e.consumeMsg(msg)

		if err == nil {
			handles = append(handles, msg.ReceiptHandle)
		} else {
			log.Errorf("dealMsgRes: consumeMsg error: %v", err)
		}
	}

	if len(handles) == 0 {
		return
	}

	ackErr := mqConsumer.AckMessage(handles)
	if ackErr != nil {
		log.Errorf("ack err:%v", ackErr)
		if errAckItems, ok := ackErr.(errors.ErrCode).Context()["Detail"].([]rocket.ErrAckItem); ok {
			for _, errAckItem := range errAckItems {
				log.Errorf("ErrorHandle:%s, ErrorCode:%s, ErrorMsg:%s",
					errAckItem.ErrorHandle, errAckItem.ErrorCode, errAckItem.ErrorMsg)
			}
		}
		time.Sleep(time.Duration(3) * time.Second)
	}
}

func (e *RocketEmitter) consumeMsg(msg rocket.ConsumeMessageEntry) (err error) {
	traceId := msg.Properties[_traceIdTag]
	utils.SetTraceId(traceId, func() {
		t := msg.Properties[_eventTypeTag]
		log.Infof("message|revice msg:Id=%s, type=%s, tag=%s, content=%s", msg.MessageId, t, msg.MessageTag, msg.MessageBody)
		err = e.routeMsg(msg.MessageId, t, msg.MessageBody)
	})
	return err
}

func (e *RocketEmitter) routeMsg(msgId string, msgType string, msgContent string) error {
	keeper, exists := e.handlers.Get(msgType)
	if !exists {
		log.Warnf("messageId=%s, do not found handler", msgId)
		return nil
	}

	if keeper.decoder == nil {
		log.Warnf("messageId=%s, do not found decoder", msgId)
		return nil
	}

	domainEvent, err := keeper.decoder.Decode([]byte(msgContent))
	if err != nil {
		log.Warnf("messageId=%s, decode err:%v", msgId, err)
		return nil
	}

	var futures []*async.Future[any]
	for _, fn := range keeper.slice {
		future := async.RunTask[any](dealEventProxy(fn, domainEvent))
		futures = append(futures, future)
	}
	return async.CompositeFutureAll(futures)
}

func (e *RocketEmitter) getStructTypeHash(t string) string {
	return crypto.Md5(e.conf.DebugTag + t)[:5]
}
