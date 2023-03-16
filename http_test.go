package ngroklistener

import (
	"testing"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/stretchr/testify/require"
)

func TestParseHTTP(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
		expected  HTTP
	}{
		{
			name: "default",
			input: `http {
			}`,
			shouldErr: false,
			expected:  HTTP{AllowCIDR: []string{}, DenyCIDR: []string{}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := caddyfile.NewTestDispenser(test.input)
			tun := HTTP{}
			err := tun.UnmarshalCaddyfile(d)
			tun.Provision(caddy.Context{})

			if test.shouldErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestHTTPBasicAuth(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
		expected  HTTP
	}{
		{
			name: "empty",
			input: `http {
			}`,
			shouldErr: false,
			expected:  HTTP{BasicAuth: map[string]string{}},
		},
		{
			name: "single-inline",
			input: `http {
				basic_auth foo barbarbar
			}`,
			shouldErr: false,
			expected:  HTTP{BasicAuth: map[string]string{"foo": "barbarbar"}},
		},
		{
			name: "single-block",
			input: `http {
				basic_auth {
					foo barbarbar
				}
			}`,
			shouldErr: false,
			expected:  HTTP{BasicAuth: map[string]string{"foo": "barbarbar"}},
		},
		{
			name: "multiple",
			input: `http {
				basic_auth foo barbarbar
				basic_auth spam eggsandcheese
				basic_auth {
					bar bazbazbaz
					bam bambinos
				}
			}`,
			shouldErr: false,
			expected:  HTTP{BasicAuth: map[string]string{"bam": "bambinos", "bar": "bazbazbaz", "foo": "barbarbar", "spam": "eggsandcheese"}},
		},
		{
			name: "password-too-short",
			input: `http {
				basic_auth foo bar
			}`,
			shouldErr: true,
			expected:  HTTP{BasicAuth: map[string]string{}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := caddyfile.NewTestDispenser(test.input)
			tun := HTTP{}
			err := tun.UnmarshalCaddyfile(d)
			tun.Provision(caddy.Context{})

			if test.shouldErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
				require.ElementsMatch(t, test.expected.BasicAuth, tun.BasicAuth)
			}
		})
	}
}

func TestHTTPCircuitBreaker(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
		expected  HTTP
	}{
		{
			name: "absent",
			input: `http {
			}`,
			shouldErr: false,
			expected:  HTTP{CircuitBreaker: 0},
		},
		{
			name: "breakered",
			input: `http {
				circuit_breaker 0.5
			}`,
			shouldErr: false,
			expected:  HTTP{CircuitBreaker: 0.5},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := caddyfile.NewTestDispenser(test.input)
			tun := HTTP{}
			err := tun.UnmarshalCaddyfile(d)
			tun.Provision(caddy.Context{})

			if test.shouldErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, test.expected.CircuitBreaker, tun.CircuitBreaker)
			}
		})
	}
}

func TestHTTPCompression(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
		expected  HTTP
	}{
		{
			name: "absent",
			input: `http {
			}`,
			shouldErr: false,
			expected:  HTTP{Compression: false},
		},
		{
			name: "compressed-off",
			input: `http {
				compression off
			}`,
			shouldErr: false,
			expected:  HTTP{Compression: false},
		},
		{
			name: "compressed-false",
			input: `http {
				compression false
			}`,
			shouldErr: false,
			expected:  HTTP{Compression: false},
		},
		{
			name: "compressed-true",
			input: `http {
				compression true
			}`,
			shouldErr: false,
			expected:  HTTP{Compression: true},
		},
		{
			name: "compressed-no-arg",
			input: `http {
				compression
			}`,
			shouldErr: false,
			expected:  HTTP{Compression: true},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := caddyfile.NewTestDispenser(test.input)
			tun := HTTP{}
			err := tun.UnmarshalCaddyfile(d)
			tun.Provision(caddy.Context{})

			if test.shouldErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, test.expected.Compression, tun.Compression)
			}
		})
	}
}

func TestHTTPWebsocketTCPConversion(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
		expected  HTTP
	}{
		{
			name: "absent",
			input: `http {
			}`,
			shouldErr: false,
			expected:  HTTP{WebsocketTCPConverter: false},
		},
		{
			name: "converted-off",
			input: `http {
				websocket_tcp_converter off
			}`,
			shouldErr: false,
			expected:  HTTP{WebsocketTCPConverter: false},
		},
		{
			name: "converted-false",
			input: `http {
				websocket_tcp_converter false
			}`,
			shouldErr: false,
			expected:  HTTP{WebsocketTCPConverter: false},
		},
		{
			name: "converted-true",
			input: `http {
				websocket_tcp_converter true
			}`,
			shouldErr: false,
			expected:  HTTP{WebsocketTCPConverter: true},
		},
		{
			name: "converted-no-arg",
			input: `http {
				websocket_tcp_converter true
			}`,
			shouldErr: false,
			expected:  HTTP{WebsocketTCPConverter: true},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := caddyfile.NewTestDispenser(test.input)
			tun := HTTP{}
			err := tun.UnmarshalCaddyfile(d)
			tun.Provision(caddy.Context{})

			if test.shouldErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, test.expected.WebsocketTCPConverter, tun.WebsocketTCPConverter)
			}
		})
	}
}

func TestHTTPDomain(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
		expected  HTTP
	}{
		{
			name: "absent",
			input: `http {
			}`,
			shouldErr: false,
			expected:  HTTP{Domain: ""},
		},
		{
			name: "with domain",
			input: `http {
				domain foo.ngrok.io
			}`,
			shouldErr: false,
			expected:  HTTP{Domain: "foo.ngrok.io"},
		},
		{
			name: "domain-no-args",
			input: `http {
				domain
			}`,
			shouldErr: true,
			expected:  HTTP{Domain: ""},
		},
		{
			name: "domain-too-many-args",
			input: `http {
				domain foo.ngrok.io foo.ngrok.io
			}`,
			shouldErr: true,
			expected:  HTTP{Domain: ""},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := caddyfile.NewTestDispenser(test.input)
			tun := HTTP{}
			err := tun.UnmarshalCaddyfile(d)
			tun.Provision(caddy.Context{})

			if test.shouldErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, test.expected.Domain, tun.Domain)
			}
		})
	}
}

func TestHTTPMetadata(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
		expected  HTTP
	}{
		{
			name: "absent",
			input: `http {
			}`,
			shouldErr: false,
			expected:  HTTP{Metadata: ""},
		},
		{
			name: "with metadata",
			input: `http {
				metadata test
			}`,
			shouldErr: false,
			expected:  HTTP{Metadata: "test"},
		},
		{
			name: "metadata-single-arg-quotes",
			input: `http {
				metadata "Hello, World!"
			}`,
			shouldErr: false,
			expected:  HTTP{Metadata: "Hello, World!"},
		},
		{
			name: "metadata-no-args",
			input: `http {
				metadata
			}`,
			shouldErr: true,
			expected:  HTTP{Metadata: ""},
		},
		{
			name: "metadata-too-many-args",
			input: `http {
				metadata test test2
			}`,
			shouldErr: true,
			expected:  HTTP{Metadata: ""},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := caddyfile.NewTestDispenser(test.input)
			tun := HTTP{}
			err := tun.UnmarshalCaddyfile(d)
			tun.Provision(caddy.Context{})

			if test.shouldErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, test.expected.Metadata, tun.Metadata)
			}
		})
	}
}

func TestHTTPScheme(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
		expected  HTTP
	}{
		{
			name: "default",
			input: `http {
			}`,
			shouldErr: false,
			expected:  HTTP{Scheme: ""},
		},
		{
			name: "set https",
			input: `http {
				scheme https
			}`,
			shouldErr: false,
			expected:  HTTP{Scheme: "https"},
		},
		{
			name: "set http",
			input: `http {
				scheme http
			}`,
			shouldErr: false,
			expected:  HTTP{Scheme: "http"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := caddyfile.NewTestDispenser(test.input)
			tun := HTTP{}
			err := tun.UnmarshalCaddyfile(d)
			tun.Provision(caddy.Context{})

			if test.shouldErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, test.expected.Scheme, tun.Scheme)
			}
		})
	}
}

func TestHTTPCIDRRestrictions(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
		expected  HTTP
	}{
		{
			name: "absent",
			input: `http {
			}`,
			shouldErr: false,
			expected:  HTTP{AllowCIDR: []string{}, DenyCIDR: []string{}},
		},
		{
			name: "allow",
			input: `http {
				allow 127.0.0.0/8
			}`,
			shouldErr: false,
			expected:  HTTP{AllowCIDR: []string{"127.0.0.0/8"}, DenyCIDR: []string{}},
		},
		{
			name: "deny",
			input: `http {
				deny 127.0.0.0/8
			}`,
			shouldErr: false,
			expected:  HTTP{AllowCIDR: []string{}, DenyCIDR: []string{"127.0.0.0/8"}},
		},
		{
			name: "allow multi",
			input: `http {
				allow 127.0.0.0/8
				allow 10.0.0.0/8
			}`,
			shouldErr: false,
			expected:  HTTP{AllowCIDR: []string{"127.0.0.0/8", "10.0.0.0/8"}, DenyCIDR: []string{}},
		},
		{
			name: "allow multi inline",
			input: `http {
				allow 127.0.0.0/8 10.0.0.0/8
			}`,
			shouldErr: false,
			expected:  HTTP{AllowCIDR: []string{"127.0.0.0/8", "10.0.0.0/8"}, DenyCIDR: []string{}},
		},
		{
			name: "deny multi",
			input: `http {
				deny 127.0.0.0/8
				deny 10.0.0.0/8
			}`,
			shouldErr: false,
			expected:  HTTP{AllowCIDR: []string{}, DenyCIDR: []string{"127.0.0.0/8", "10.0.0.0/8"}},
		},
		{
			name: "deny multi inline",
			input: `http {
				deny 127.0.0.0/8 10.0.0.0/8
			}`,
			shouldErr: false,
			expected:  HTTP{AllowCIDR: []string{}, DenyCIDR: []string{"127.0.0.0/8", "10.0.0.0/8"}},
		},
		{
			name: "allow and deny multi",
			input: `http {
				allow 127.0.0.0/8
				allow 10.0.0.0/8
				deny 192.0.0.0/8
				deny 172.0.0.0/8
			}`,
			shouldErr: false,
			expected:  HTTP{AllowCIDR: []string{"127.0.0.0/8", "10.0.0.0/8"}, DenyCIDR: []string{"192.0.0.0/8", "172.0.0.0/8"}},
		},
		{
			name: "allow-no-args",
			input: `http {
				allow
			}`,
			shouldErr: true,
			expected:  HTTP{AllowCIDR: []string{}, DenyCIDR: []string{}},
		},
		{
			name: "deny-no-args",
			input: `http {
				deny
			}`,
			shouldErr: true,
			expected:  HTTP{AllowCIDR: []string{}, DenyCIDR: []string{}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := caddyfile.NewTestDispenser(test.input)
			tun := HTTP{}
			err := tun.UnmarshalCaddyfile(d)
			tun.Provision(caddy.Context{})

			if test.shouldErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
				require.ElementsMatch(t, test.expected.AllowCIDR, tun.AllowCIDR)
				require.ElementsMatch(t, test.expected.DenyCIDR, tun.DenyCIDR)
			}
		})
	}
}
