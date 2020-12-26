//nolint:funlen
package servicediscovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/consul/api"
)

func TestServiceRegistration_generateServiceDiscoveryRegistration(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		sr := ServiceRegistration{
			ServiceDiscoveryURL:      "http://test.local:1234",
			ServiceDiscoveryUsername: "",
			ServiceDiscoveryPassword: "",
			MicroserviceID:           "a1b2-c3d4-e5f6",
			MicroserviceName:         "testservice",
			MicroserviceTags:         []string{"test", "test2"},
			MicroserviceMetadata:     map[string]string{"test": "TEST"},
			MicroserviceURL:          "http://test.local:4321",
		}

		actual, err := sr.generateServiceDiscoveryRegistration()

		require.NoError(t, err)
		require.NotNil(t, actual)

		assert.Equal(t, actual, api.AgentServiceRegistration{
			ID:      "a1b2-c3d4-e5f6",
			Name:    "testservice",
			Tags:    []string{"test", "test2"},
			Port:    4321,
			Address: "test.local",
			Meta:    map[string]string{"test": "TEST"},
		})
	})

	t.Run("fail", func(t *testing.T) {
		t.Run("url_parsing", func(t *testing.T) {
			sr := ServiceRegistration{
				ServiceDiscoveryURL:      "http://test.local:1234",
				ServiceDiscoveryUsername: "",
				ServiceDiscoveryPassword: "",
				MicroserviceID:           "a1b2-c3d4-e5f6",
				MicroserviceName:         "testservice",
				MicroserviceTags:         []string{"test", "test2"},
				MicroserviceMetadata:     map[string]string{"test": "TEST"},
				MicroserviceURL:          "http::://test.local:4321",
			}

			actual, err := sr.generateServiceDiscoveryRegistration()

			require.Error(t, err)
			require.NotNil(t, actual)

			assert.Equal(t, actual, api.AgentServiceRegistration{})
		})
		t.Run("port_parsing", func(t *testing.T) {
			sr := ServiceRegistration{
				ServiceDiscoveryURL:      "http://test.local:1234",
				ServiceDiscoveryUsername: "",
				ServiceDiscoveryPassword: "",
				MicroserviceID:           "a1b2-c3d4-e5f6",
				MicroserviceName:         "testservice",
				MicroserviceTags:         []string{"test", "test2"},
				MicroserviceMetadata:     map[string]string{"test": "TEST"},
				MicroserviceURL:          "http://test.local:43k21",
			}

			actual, err := sr.generateServiceDiscoveryRegistration()

			require.Error(t, err)
			require.NotNil(t, actual)

			assert.Equal(t, actual, api.AgentServiceRegistration{})
		})
	})
}

func TestServiceRegistration_isURLSet(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		sr := ServiceRegistration{
			ServiceDiscoveryURL: "http://test.local",
		}

		actual := sr.isURLSet()

		assert.True(t, actual)
	})

	t.Run("fail", func(t *testing.T) {
		sr := ServiceRegistration{
			ServiceDiscoveryURL: "",
		}

		actual := sr.isURLSet()

		assert.False(t, actual)
	})
}
