package messaging

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var ErrorMetadataConfigIncomplete = errors.New("service metadata config incomplete")

var _ Broker = &natsBroker{}

type natsBroker struct {
	natsConnection   *nats.Conn
	jetstreamManager jetstream.JetStream
	subscriptions    []*nats.Subscription
	streams          []jetstream.Stream
	cancelFunctions  []context.CancelFunc
	consumer         []jetstream.Consumer
}

func newNatsBroker(url, userName, password, serviceName string) (Broker, error) {
	natsOptions := make([]nats.Option, 0, 6)

	natsOptions = append(natsOptions, nats.Name(serviceName))
	natsOptions = append(natsOptions, nats.UserInfo(userName, password))
	natsOptions = append(natsOptions, nats.MaxReconnects(maxReconnects))
	natsOptions = append(natsOptions, nats.ReconnectWait(reconnectWait))
	natsOptions = append(natsOptions, nats.ReconnectJitter(reconnectJitterNoTLS, reconnectJitterTLS))
	natsOptions = append(natsOptions, nats.ReconnectBufSize(reconnectBuffer))

	conn, err := nats.Connect(url, natsOptions...)
	if err != nil {
		return nil, err
	}

	jsManager, err := jetstream.New(conn)
	if err != nil {
		return nil, err
	}

	return &natsBroker{natsConnection: conn, jetstreamManager: jsManager}, nil
}

// CreatePersistentMessageStore creates a jetstream stream.
func (nb *natsBroker) CreatePersistentMessageStore(name string, subjects []string) (jetstream.Stream, error) {
	ctx, cancelFun := context.WithCancel(context.Background())

	cfg := jetstream.StreamConfig{
		Name:      name,
		Subjects:  subjects,
		Retention: jetstream.LimitsPolicy,
		MaxMsgs:   200,
		Storage:   jetstream.MemoryStorage,
	}

	stream, err := nb.jetstreamManager.CreateOrUpdateStream(ctx, cfg)

	nb.streams = append(nb.streams, stream)
	nb.cancelFunctions = append(nb.cancelFunctions, cancelFun)

	return stream, err
}

// CreatePersistentMessageReceiver creates a jetstream consumer.
func (nb *natsBroker) CreatePersistentMessageReceiver(name, id, address, port, cType string, stream jetstream.Stream) (jetstream.Consumer, error) {
	ctx, cancelFun := context.WithCancel(context.Background())

	messageStartTime := time.Now().Add(-15 * time.Minute).Round(time.Minute)

	cfg := jetstream.ConsumerConfig{
		Name:          name,
		DeliverPolicy: jetstream.DeliverByStartTimePolicy,
		OptStartTime:  &messageStartTime,
		AckPolicy:     jetstream.AckNonePolicy,
		Metadata:      map[string]string{"id": id, "address": address, "port": port, "type": cType},
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, cfg)

	nb.consumer = append(nb.consumer, consumer)
	nb.cancelFunctions = append(nb.cancelFunctions, cancelFun)

	return consumer, err
}

// Close closes the connection to the broker and unsubscribes on all topics.
func (nb *natsBroker) Close() error {
	for _, subscription := range nb.subscriptions {
		if err := subscription.Unsubscribe(); err != nil {
			return err
		}
	}

	for _, cancelFunc := range nb.cancelFunctions {
		cancelFunc()
	}

	nb.natsConnection.Close()

	return nil
}

// Subscribe subscribes to a pubsub topic.
func (nb *natsBroker) Subscribe(topic string, callback nats.MsgHandler) error {
	subscription, err := nb.natsConnection.Subscribe(topic, callback)
	if err != nil {
		return err
	}

	nb.subscriptions = append(nb.subscriptions, subscription)

	return nil
}

// QueueGroupSubscribe implements Broker.
func (nb *natsBroker) QueueGroupSubscribe(topic, group string, callback nats.MsgHandler) error {
	subscription, err := nb.natsConnection.QueueSubscribe(topic, group, callback)
	if err != nil {
		return err
	}

	nb.subscriptions = append(nb.subscriptions, subscription)

	return nil
}

// Publish send a payload to a topic.
func (nb *natsBroker) Publish(topic string, payload []byte) error {
	return nb.natsConnection.Publish(topic, payload)
}

// PublishSync send a payload to a topic and wait for a reply.
func (nb *natsBroker) PublishSync(topic string, payload []byte, timeout time.Duration) (*nats.Msg, error) {
	return nb.natsConnection.Request(topic, payload, timeout)
}

// Unsubscribe unsubscribes from a pubsub topic.
func (nb *natsBroker) Unsubscribe(topic string) error {
	for _, subscription := range nb.subscriptions {
		if subscription.Subject != topic {
			continue
		}
		if err := subscription.Unsubscribe(); err != nil {
			return err
		}
	}

	return nil
}

// RetrieveRegisteredConsumers returns registered consumers.
func (nb *natsBroker) RetrieveRegisteredConsumers(stream jetstream.Stream) ([]ResolvedService, error) {
	serviceList := make([]ResolvedService, 0, 10)

	consumerList := stream.ListConsumers(context.TODO())
	for consumer := range consumerList.Info() {
		cType, ok := consumer.Config.Metadata["type"]
		if !ok || cType != "secondary" {
			continue
		}

		id, ok := consumer.Config.Metadata["id"]
		if !ok {
			return nil, ErrorMetadataConfigIncomplete
		}
		address, ok := consumer.Config.Metadata["address"]
		if !ok {
			return nil, ErrorMetadataConfigIncomplete
		}
		portString, ok := consumer.Config.Metadata["port"]
		if !ok {
			return nil, ErrorMetadataConfigIncomplete
		}

		port, err := strconv.Atoi(portString)
		if err != nil {
			return nil, err
		}

		serviceList = append(serviceList, ResolvedService{ID: id, Address: address, Port: port})
	}

	return serviceList, nil
}

func (nb *natsBroker) Consume(consumer jetstream.Consumer, handler func(msg jetstream.Msg)) (jetstream.ConsumeContext, error) {
	return consumer.Consume(handler)
}

func (nb *natsBroker) PersistedPublish(topic string, payload []byte) (err error) {
	_, err = nb.jetstreamManager.Publish(context.TODO(), topic, payload)

	return
}
