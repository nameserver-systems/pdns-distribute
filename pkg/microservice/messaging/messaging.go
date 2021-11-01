package messaging

import (
	"time"

	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
	"github.com/nats-io/nats.go"
)

type MessageBroker struct {
	URL      string
	Username string
	Password string

	NatsConnection *nats.Conn
	subscriptions  []*nats.Subscription
}

func (mb *MessageBroker) StartMessageBrokerConnection(servicename string) {
	const (
		maxreconnects        = 5000
		reconnectwait        = 10 * time.Second
		reconnectjitternotls = 500 * time.Millisecond
		reconnectjittertls   = 2 * time.Second
		reconnectbuffer      = 1024 * 1024 * 10
	)

	var err error

	var clientname nats.Option

	var auth nats.Option

	clientname = nats.Name(servicename)

	if mb.Username != "" && mb.Password != "" {
		auth = nats.UserInfo(mb.Username, mb.Password)
	}

	mb.NatsConnection, err = nats.Connect(mb.URL, auth, clientname, nats.MaxReconnects(maxreconnects), nats.ReconnectWait(reconnectwait), nats.ReconnectJitter(reconnectjitternotls, reconnectjittertls), nats.ReconnectBufSize(reconnectbuffer))
	if err != nil {
		logger.FatalErrLog(err)
	}
}

func (mb *MessageBroker) CloseConnection() {
	for _, subscription := range mb.subscriptions {
		err := subscription.Unsubscribe()
		if err != nil {
			logger.FatalErrLog(err)
		}
	}

	mb.NatsConnection.Close()
}

func (mb *MessageBroker) SubscribeAsync(topic string, callback nats.MsgHandler) {
	subscription, err := mb.NatsConnection.Subscribe(topic, callback)
	if err != nil {
		logger.FatalErrLog(err)
	}

	mb.subscriptions = append(mb.subscriptions, subscription)
}

func (mb *MessageBroker) SubscribeQueueAsync(topic string, queue string, callback nats.MsgHandler) {
	subscription, err := mb.NatsConnection.QueueSubscribe(topic, queue, callback)
	if err != nil {
		logger.FatalErrLog(err)
	}

	mb.subscriptions = append(mb.subscriptions, subscription)
}

func (mb *MessageBroker) DeSubscribeAsync(topic string) {
	subscriptions := mb.subscriptions
	for _, subscription := range subscriptions {
		if subscription.Subject == topic {
			err := subscription.Unsubscribe()
			if err != nil {
				logger.FatalErrLog(err)
			}
		}
	}
}

func (mb *MessageBroker) Publish(topic string, payload []byte) {
	err := mb.NatsConnection.Publish(topic, payload)
	if err != nil {
		logger.FatalErrLog(err)
	}
}

func (mb *MessageBroker) PublishRequestAndWait(topic string, payload []byte, timeout time.Duration) (*nats.Msg, error) {
	msg, err := mb.NatsConnection.Request(topic, payload, timeout)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
