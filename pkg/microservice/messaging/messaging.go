package messaging

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	maxReconnects        = 5000
	reconnectWait        = 10 * time.Second
	reconnectJitterNoTLS = 500 * time.Millisecond
	reconnectJitterTLS   = 2 * time.Second
	reconnectBuffer      = 1024 * 1024 * 10
)

type ResolvedService struct {
	ID      string
	Address string
	Port    int
}

type Broker interface {
	CreatePersistentMessageStore(name string, subjects []string) (jetstream.Stream, error)
	CreatePersistentMessageReceiver(name, id, address, port, ctype string, stream jetstream.Stream) (jetstream.Consumer, error)
	RetrieveRegisteredConsumers(stream jetstream.Stream) ([]ResolvedService, error)
	PersistedPublish(topic string, payload []byte) error
	Consume(consumer jetstream.Consumer, handler func(msg jetstream.Msg)) (jetstream.ConsumeContext, error)
	Subscribe(topic string, callback nats.MsgHandler) error
	QueueGroupSubscribe(topic, queue string, callback nats.MsgHandler) error
	Publish(topic string, payload []byte) error
	PublishSync(topic string, payload []byte, timeout time.Duration) (*nats.Msg, error)
	Unsubscribe(topic string) error
	Close() error
}

type MessageBroker struct {
	url      string
	userName string
	password string

	broker   Broker
	stream   jetstream.Stream
	consumer jetstream.Consumer
}

func NewMessageBroker(url, userName, password, serviceName string) (*MessageBroker, error) {
	broker, err := newNatsBroker(url, userName, password, serviceName)
	if err != nil {
		return nil, err
	}

	return &MessageBroker{url: url, userName: userName, password: password, broker: broker}, nil
}

func (mb *MessageBroker) CloseConnection() error {
	return mb.broker.Close()
}

func (mb *MessageBroker) SubscribeAsync(topic string, callback nats.MsgHandler) error {
	return mb.broker.Subscribe(topic, callback)
}

func (mb *MessageBroker) SubscribeQueueAsync(topic string, queue string, callback nats.MsgHandler) error {
	return mb.broker.QueueGroupSubscribe(topic, queue, callback)
}

func (mb *MessageBroker) DesubscribeAsync(topic string) error {
	return mb.broker.Unsubscribe(topic)
}

func (mb *MessageBroker) Publish(topic string, payload []byte) error {
	return mb.broker.Publish(topic, payload)
}

func (mb *MessageBroker) PublishRequestAndWait(topic string, payload []byte, timeout time.Duration) (*nats.Msg, error) {
	return mb.broker.PublishSync(topic, payload, timeout)
}

func (mb *MessageBroker) RetrieveRegisteredConsumers(stream jetstream.Stream) ([]ResolvedService, error) {
	return mb.broker.RetrieveRegisteredConsumers(stream)
}

func (mb *MessageBroker) CreatePersistentMessageStore(name string, subjects []string) (jetstream.Stream, error) {
	return mb.broker.CreatePersistentMessageStore(name, subjects)
}

func (mb *MessageBroker) CreatePersistentMessageReceiver(name, id, address, port, cType string, stream jetstream.Stream) (jetstream.Consumer, error) {
	return mb.broker.CreatePersistentMessageReceiver(name, id, address, port, cType, stream)
}

func (mb *MessageBroker) SetStream(stream jetstream.Stream) {
	mb.stream = stream
}

func (mb *MessageBroker) GetStream() jetstream.Stream {
	return mb.stream
}

func (mb *MessageBroker) SetConsumer(consumer jetstream.Consumer) {
	mb.consumer = consumer
}

func (mb *MessageBroker) GetConsumer() jetstream.Consumer {
	return mb.consumer
}

func (mb *MessageBroker) Consume(consumer jetstream.Consumer, handler func(msg jetstream.Msg)) (jetstream.ConsumeContext, error) {
	return mb.broker.Consume(consumer, handler)
}

func (mb *MessageBroker) PersistedPublish(topic string, payload []byte) error {
	return mb.broker.PersistedPublish(topic, payload)
}
