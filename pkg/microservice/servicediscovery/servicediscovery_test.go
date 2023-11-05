//nolint:scopelint,lll,testpackage
package servicediscovery

import (
	"testing"

	consul "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceDiscovery_newClient(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		sd := ServiceDiscovery{
			serviceID:              "1234",
			consulURL:              "http://test.local:1234",
			consulBasicAuth:        &consul.HttpBasicAuth{},
			ConsulClient:           nil,
			serviceListenAddresses: map[string]consul.ServiceAddress{},
			serviceHealthCheck:     &consul.AgentServiceCheck{},
		}

		err := sd.newClient()

		require.NoError(t, err)
	})

	t.Run("fail_client_exists", func(t *testing.T) {
		sd := ServiceDiscovery{
			serviceID:              "1234",
			consulURL:              "http://test.local:1234",
			consulBasicAuth:        &consul.HttpBasicAuth{},
			ConsulClient:           &consul.Client{},
			serviceListenAddresses: map[string]consul.ServiceAddress{},
			serviceHealthCheck:     &consul.AgentServiceCheck{},
		}

		err := sd.newClient()

		require.NoError(t, err)
	})
}

func TestServiceDiscovery_isClientInitiated(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		sd := ServiceDiscovery{
			ConsulClient: &consul.Client{},
		}

		actual := sd.isClientInitiated()

		assert.True(t, actual)
	})

	t.Run("fail", func(t *testing.T) {
		sd := ServiceDiscovery{
			ConsulClient: nil,
		}

		actual := sd.isClientInitiated()

		assert.False(t, actual)
	})
}

func TestServiceDiscovery_setServiceDiscoveryBasicAuthCredentials(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		sd := ServiceDiscovery{
			consulBasicAuth: &consul.HttpBasicAuth{},
		}

		sd.setServiceDiscoveryBasicAuthCredentials("TEST-USER", "1234")

		assert.EqualValues(t, "TEST-USER", sd.consulBasicAuth.Username)
		assert.EqualValues(t, "1234", sd.consulBasicAuth.Password)
	})
}

func TestServiceDiscovery_isHealthCheckSet(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		sd := ServiceDiscovery{serviceHealthCheck: &consul.AgentServiceCheck{}}

		actual := sd.isHealthCheckSet()

		assert.True(t, actual)
	})

	t.Run("fail", func(t *testing.T) {
		sd := ServiceDiscovery{}

		actual := sd.isHealthCheckSet()

		assert.False(t, actual)
	})
}

func TestServiceDiscovery_isListenAddressesSet(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		sd := ServiceDiscovery{serviceListenAddresses: map[string]consul.ServiceAddress{"test": {}}}

		actual := sd.isListenAddressesSet()

		assert.True(t, actual)
	})

	t.Run("fail", func(t *testing.T) {
		sd := ServiceDiscovery{}

		actual := sd.isListenAddressesSet()

		assert.False(t, actual)
	})
}

func TestServiceDiscovery_insertListenAddressesInRegistration(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		sr := &consul.AgentServiceRegistration{TaggedAddresses: map[string]consul.ServiceAddress{}}
		sd := ServiceDiscovery{}

		sd.insertListenAddressesInRegistration(sr)

		assert.IsType(t, map[string]consul.ServiceAddress{}, sd.serviceListenAddresses)
	})

	t.Run("fail", func(t *testing.T) {
		sr := &consul.AgentServiceRegistration{}
		sd := ServiceDiscovery{}

		sd.insertListenAddressesInRegistration(sr)

		assert.Empty(t, sd.serviceListenAddresses)
	})
}

func TestServiceDiscovery_insertHealthCheckInRegistration(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		sr := &consul.AgentServiceRegistration{Check: &consul.AgentServiceCheck{
			Name: "test",
		}}
		sd := ServiceDiscovery{}

		sd.insertHealthCheckInRegistration(sr)

		assert.IsType(t, &consul.AgentServiceCheck{}, sd.serviceHealthCheck)
	})

	t.Run("fail", func(t *testing.T) {
		sr := &consul.AgentServiceRegistration{}
		sd := ServiceDiscovery{}

		sd.insertHealthCheckInRegistration(sr)

		assert.Empty(t, sd.serviceListenAddresses)
	})
}
