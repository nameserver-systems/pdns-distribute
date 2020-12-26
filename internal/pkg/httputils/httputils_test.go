//nolint:scopelint,lll,dupl,testpackage
package httputils

import (
	"net/http"
	"testing"
)

func Test_ensureTrailingDot(t *testing.T) {
	type args struct {
		zoneid string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "without trailing dot", args: args{zoneid: "test.de"}, want: "test.de."},
		{name: "with trailing dot", args: args{zoneid: "test.de."}, want: "test.de."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ensureTrailingDot(tt.args.zoneid); got != tt.want {
				t.Errorf("ensureTrailingDot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasNotFoundZoneID(t *testing.T) {
	type args struct {
		zoneid string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "with trailing dot", args: args{zoneid: "test.de."}, want: false},
		{name: "without trailing dot", args: args{zoneid: "test.de"}, want: false},
		{name: "has no zoneid", args: args{zoneid: ""}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasNotFoundZoneID(tt.args.zoneid); got != tt.want {
				t.Errorf("hasNotFoundZoneID() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestCloseResponseBody(t *testing.T) {
	type args struct {
		response *http.Response
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "ok", args: args{response: &http.Response{}}},
		{name: "nil response", args: args{response: nil}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
