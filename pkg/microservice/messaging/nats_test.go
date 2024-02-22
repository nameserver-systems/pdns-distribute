package messaging

import (
	"context"
	"testing"
	"time"

	"github.com/nameserver-systems/pdns-distribute/internal/pkg/mocks"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen --destination internal/pkg/mocks/jetstream.go --package mocks github.com/nats-io/nats.go/jetstream JetStream
//go:generate mockgen --destination internal/pkg/mocks/jetstream_stream.go --package mocks github.com/nats-io/nats.go/jetstream Stream
//go:generate mockgen --destination internal/pkg/mocks/jetstream_consumer_info_listener.go --package mocks github.com/nats-io/nats.go/jetstream ConsumerInfoLister
func Test_natsBroker_CreatePersistentMessageStore(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedConfig := jetstream.StreamConfig{
			Name:      "test-store",
			Subjects:  []string{"sub01", "sub02"},
			Retention: jetstream.LimitsPolicy,
			MaxMsgs:   200,
			Storage:   jetstream.MemoryStorage,
		}

		mockJetStreamManager := mocks.NewMockJetStream(ctrl)
		mockJetStreamManager.EXPECT().CreateOrUpdateStream(gomock.Any(), expectedConfig).Return(nil, nil)

		nb := natsBroker{
			jetstreamManager: mockJetStreamManager,
		}

		stream, err := nb.CreatePersistentMessageStore("test-store", []string{"sub01", "sub02"})

		require.NoError(t, err)
		assert.Nil(t, stream)

		assert.Len(t, nb.streams, 1)
		assert.Len(t, nb.cancelFunctions, 1)
	})
}

func Test_natsBroker_CreatePersistentMessageReceiver(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		messageStartTime := time.Now().Add(-15 * time.Minute).Round(time.Minute)

		expectedConfig := jetstream.ConsumerConfig{
			Name:          "test-client",
			DeliverPolicy: jetstream.DeliverByStartTimePolicy,
			OptStartTime:  &messageStartTime,
			AckPolicy:     jetstream.AckNonePolicy,
			Metadata:      map[string]string{"id": "0000", "address": "127.0.0.1", "port": "53", "type": "secondary"},
		}

		mockStream := mocks.NewMockStream(ctrl)
		mockStream.EXPECT().CreateOrUpdateConsumer(gomock.Any(), expectedConfig).Return(nil, nil)

		nb := natsBroker{}

		consumer, err := nb.CreatePersistentMessageReceiver("test-client", "0000", "127.0.0.1", "53", "secondary", mockStream)

		require.NoError(t, err)
		assert.Nil(t, consumer)

		assert.Len(t, nb.consumer, 1)
		assert.Len(t, nb.cancelFunctions, 1)
	})
}

func Test_natsBroker_PersistedPublish(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStream := mocks.NewMockJetStream(ctrl)
		mockStream.EXPECT().Publish(context.TODO(), "top1", []byte("content"))

		nb := natsBroker{
			jetstreamManager: mockStream,
		}

		err := nb.PersistedPublish("top1", []byte("content"))

		require.NoError(t, err)
	})
}

func Test_natsBroker_RetrieveRegisteredConsumers(t *testing.T) {
	t.Run("success-secondary", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockInfoListener := mocks.NewMockConsumerInfoLister(ctrl)
		mockInfoListener.EXPECT().Info().DoAndReturn(func() <-chan *jetstream.ConsumerInfo {
			cf := make(chan *jetstream.ConsumerInfo, 2)
			defer close(cf)

			cf <- &jetstream.ConsumerInfo{
				Config: jetstream.ConsumerConfig{
					Metadata: map[string]string{"type": "secondary", "id": "1", "address": "10.0.0.2", "port": "1234"},
				},
			}

			return cf
		})

		mockStream := mocks.NewMockStream(ctrl)
		mockStream.EXPECT().ListConsumers(context.TODO()).Return(mockInfoListener)

		nb := natsBroker{}
		list, err := nb.RetrieveRegisteredConsumers(mockStream)

		require.NoError(t, err)
		assert.EqualValues(t, []ResolvedService{{
			ID:      "1",
			Address: "10.0.0.2",
			Port:    1234,
		}}, list)
	})
	t.Run("success-primary", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockInfoListener := mocks.NewMockConsumerInfoLister(ctrl)
		mockInfoListener.EXPECT().Info().DoAndReturn(func() <-chan *jetstream.ConsumerInfo {
			cf := make(chan *jetstream.ConsumerInfo, 2)
			defer close(cf)

			cf <- &jetstream.ConsumerInfo{
				Config: jetstream.ConsumerConfig{
					Metadata: map[string]string{"type": "primary", "id": "1", "address": "10.0.0.2", "port": "1234"},
				},
			}

			return cf
		})

		mockStream := mocks.NewMockStream(ctrl)
		mockStream.EXPECT().ListConsumers(context.TODO()).Return(mockInfoListener)

		nb := natsBroker{}
		list, err := nb.RetrieveRegisteredConsumers(mockStream)

		require.NoError(t, err)
		assert.EqualValues(t, []ResolvedService{}, list)
	})
}
