//nolint:scopelint,lll,dupl,testpackage,noctx
package httputils

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsStatusCodeSuccesful(t *testing.T) {
	type args struct {
		statuscode int
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "status 200", args: args{statuscode: 200}, want: true},
		{name: "status 202", args: args{statuscode: 202}, want: true},
		{name: "status 300", args: args{statuscode: 300}, want: false},
		{name: "status 100", args: args{statuscode: 300}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStatusCodeSuccesful(tt.args.statuscode); got != tt.want {
				t.Errorf("IsStatusCodeSuccesful() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAddressFromURL(t *testing.T) {
	type args struct {
		url string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "url with scheme and port and path", args: args{url: "https://test.example.com:1234/testme"}, want: "test.example.com:1234", wantErr: false},
		{name: "url with port and path", args: args{url: "test.example.com:1234/testme"}, want: "", wantErr: true},
		{name: "false scheme typ", args: args{url: "httpq://test.example.com:1234/testme"}, want: "test.example.com:1234", wantErr: false},
		{name: "false scheme format", args: args{url: "http:/test.example.com:1234/testme"}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHostAndPortFromURL(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHostAndPortFromURL() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("GetHostAndPortFromURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHostnameFromURL(t *testing.T) {
	type args struct {
		url string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "url with scheme and port and path", args: args{url: "https://test.example.com:1234/testme"}, want: "test.example.com", wantErr: false},
		{name: "url with port and path", args: args{url: "test.example.com:1234/testme"}, want: "", wantErr: true},
		{name: "false scheme typ", args: args{url: "httpq://test.example.com:1234/testme"}, want: "test.example.com", wantErr: false},
		{name: "false scheme format", args: args{url: "http:/test.example.com:1234/testme"}, want: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHostnameFromURL(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHostnameFromURL() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("GetHostnameFromURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getZoneIDFromRequestPath(t *testing.T) {
	t.Run("found-zone-id", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("GET /{zone_id}", func(writer http.ResponseWriter, request *http.Request) {
			_, _ = fmt.Fprint(writer, getZoneIDFromRequestPath(request))
		})

		s := httptest.NewServer(mux)
		defer s.Close()

		resp, err := http.Get(s.URL + "/example.org")
		require.NoError(t, err)

		defer resp.Body.Close()

		actualZone, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, "example.org", string(actualZone))
	})
}

func TestGetZoneIDFromRequest(t *testing.T) { //nolint:funlen
	t.Run("success-path-zone-id", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("GET /{zone_id}", func(writer http.ResponseWriter, request *http.Request) {
			actualZone, err := GetZoneIDFromRequest(request)
			require.NoError(t, err)

			_, _ = fmt.Fprint(writer, actualZone)
		})

		s := httptest.NewServer(mux)
		defer s.Close()

		resp, err := http.Get(s.URL + "/example.org")
		require.NoError(t, err)

		defer resp.Body.Close()

		actualZone, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, "example.org.", string(actualZone))
	})
	t.Run("success-payload-zone-id", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("POST /", func(writer http.ResponseWriter, request *http.Request) {
			actualZone, err := GetZoneIDFromRequest(request)
			require.NoError(t, err)

			_, _ = fmt.Fprint(writer, actualZone)
		})

		s := httptest.NewServer(mux)
		defer s.Close()

		rBody := strings.NewReader("{\"id\":\"example.org\"}")

		resp, err := http.Post(s.URL+"/", "application/json", rBody)
		require.NoError(t, err)

		defer resp.Body.Close()

		actualZone, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, "example.org.", string(actualZone))
	})
	t.Run("fail-no-zone-id", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("POST /", func(writer http.ResponseWriter, request *http.Request) {
			actualZone, err := GetZoneIDFromRequest(request)
			require.EqualError(t, err, errZoneIDNotFound.Error())

			_, _ = fmt.Fprint(writer, actualZone)
		})

		s := httptest.NewServer(mux)
		defer s.Close()

		rBody := strings.NewReader("{}")

		resp, err := http.Post(s.URL+"/", "application/json", rBody)
		require.NoError(t, err)

		defer resp.Body.Close()

		actualZone, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, "", string(actualZone))
	})
}

func Test_getZoneIDFromRequestBody(t *testing.T) {
	t.Run("success-id", func(t *testing.T) {
		rBody := strings.NewReader("{\"id\":\"example.org\"}")

		r, err := http.NewRequest(http.MethodPost, "", rBody)
		require.NoError(t, err)

		actualZoneID, err := getZoneIDFromRequestBody(r)
		require.NoError(t, err)

		assert.Equal(t, "example.org", actualZoneID)
	})
	t.Run("success-name", func(t *testing.T) {
		rBody := strings.NewReader("{\"name\":\"example.org\"}")

		r, err := http.NewRequest(http.MethodPost, "", rBody)
		require.NoError(t, err)

		actualZoneID, err := getZoneIDFromRequestBody(r)
		require.NoError(t, err)

		assert.Equal(t, "example.org", actualZoneID)
	})
	t.Run("fail-invalid-json", func(t *testing.T) {
		rBody := strings.NewReader("")

		r, err := http.NewRequest(http.MethodPost, "", rBody)
		require.NoError(t, err)

		actualZoneID, err := getZoneIDFromRequestBody(r)

		require.EqualError(t, err, "unexpected end of JSON input")
		assert.Empty(t, actualZoneID)
	})
	t.Run("fail-no-zone-id", func(t *testing.T) {
		rBody := strings.NewReader("{}")

		r, err := http.NewRequest(http.MethodPost, "", rBody)
		require.NoError(t, err)

		actualZoneID, err := getZoneIDFromRequestBody(r)

		require.EqualError(t, err, errZoneIDNotFound.Error())
		assert.Empty(t, actualZoneID)
	})
}

func TestExecutePowerDNSRequest(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			require.Equal(t, "1234", request.Header.Get("X-API-Key"))
			data, err := io.ReadAll(request.Body)
			require.NoError(t, err)
			require.Equal(t, "{}", string(data))

			_, _ = fmt.Fprint(writer, "ok")
		}))
		defer s.Close()

		rBody := strings.NewReader("{}")

		resp, err := ExecutePowerDNSRequest(http.MethodPost, s.URL, "1234", rBody)
		require.NoError(t, err)

		assert.Equal(t, "ok", resp)
	})
}
