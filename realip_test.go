package extmiddleware

import (
	"net/http"
	"testing"
)

func TestIsPrivateAddr(t *testing.T) {
	testData := map[string]bool{
		"127.0.0.0":   true,
		"10.0.0.0":    true,
		"169.254.0.0": true,
		"192.168.0.0": true,
		"::1":         true,
		"fc00::":      true,

		"172.15.0.0": false,
		"172.16.0.0": true,
		"172.31.0.0": true,
		"172.32.0.0": false,

		"147.12.56.11": false,
	}

	for addr, isLocal := range testData {
		isPrivate, err := isPrivateAddress(addr)
		if err != nil {
			t.Errorf("fail processing %s: %v", addr, err)
		}

		if isPrivate != isLocal {
			format := "%s should "
			if !isLocal {
				format += "not "
			}
			format += "be local address"

			t.Errorf(format, addr)
		}
	}
}

func Test_realIP(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test remote addr",
			args: args{
				r: &http.Request{
					RemoteAddr: "98.62.45.11",
				},
			},
			want: "98.62.45.11",
		},
		{
			name: "Test X-Forwarded-For single",
			args: args{
				r: &http.Request{
					RemoteAddr: "98.62.45.11",
					Header:     http.Header{http.CanonicalHeaderKey("X-Forwarded-For"): {"100.100.100.100"}},
				},
			},
			want: "100.100.100.100",
		},
		{
			name: "Test X-Forwarded-For multiple",
			args: args{
				r: &http.Request{
					RemoteAddr: "98.62.45.11",
					Header:     http.Header{http.CanonicalHeaderKey("X-Forwarded-For"): {"100.100.100.100,127.0.0.1"}},
				},
			},
			want: "100.100.100.100",
		},
		{
			name: "Test X-Real-IP",
			args: args{
				r: &http.Request{
					RemoteAddr: "98.62.45.11",
					Header:     http.Header{http.CanonicalHeaderKey("X-Real-IP"): {"100.100.100.100"}},
				},
			},
			want: "100.100.100.100",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := realIP(tt.args.r); got != tt.want {
				t.Errorf("realIP() = %v, want %v", got, tt.want)
			}
		})
	}
}
