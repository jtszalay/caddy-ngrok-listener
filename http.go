package ngroklistener

import (
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"go.uber.org/zap"
	"golang.ngrok.com/ngrok/config"
)

func init() {
	caddy.RegisterModule(new(HTTP))
}

// ngrok HTTP tunnel
type HTTP struct {
	opts []config.HTTPEndpointOption

	// Rejects connections that do not match the given CIDRs
	AllowCIDR []string `json:"allowCidr,omitempty"`

	// Rejects connections that match the given CIDRs and allows all other CIDRs.
	DenyCIDR []string `json:"denyCidr,omitempty"`

	// the domain for this edge.
	Domain string `json:"domain,omitempty"`

	// opaque metadata string for this tunnel.
	Metadata string `json:"metadata,omitempty"`

	// sets the scheme for this edge.
	Scheme string `json:"scheme,omitempty"`

	// the 5XX response ratio at which the ngrok edge will stop sending requests to this tunnel.
	CircuitBreaker *float64 `json:"circuitBreaker,omitempty"`

	// enables gzip compression.
	EnableCompression bool `json:"enableCompression,omitempty"`

	// enables the websocket-to-tcp converter.
	EnableWebsocketTCPConversion bool `json:"enableWebsocketTcpConversion,omitempty"`

	// A map of basicauth, username and password value pairs for this tunnel.
	BasicAuth map[string]string `json:"basicauth,omitempty"`

	l *zap.Logger
}

// CaddyModule implements caddy.Module
func (*HTTP) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "caddy.listeners.ngrok.tunnels.http",
		New: func() caddy.Module {
			return new(HTTP)
		},
	}
}

// Provision implements caddy.Provisioner
func (t *HTTP) Provision(ctx caddy.Context) error {
	t.l = ctx.Logger()

	return nil
}

func (t *HTTP) ProvisionOpts() error {
	if t.Domain != "" {
		t.opts = append(t.opts, config.WithDomain(t.Domain))
	}

	if t.Metadata != "" {
		t.opts = append(t.opts, config.WithMetadata(t.Metadata))
	}

	if t.AllowCIDR != nil {
		t.opts = append(t.opts, config.WithAllowCIDRString(t.AllowCIDR...))
	}

	if t.DenyCIDR != nil {
		t.opts = append(t.opts, config.WithDenyCIDRString(t.DenyCIDR...))
	}

	if t.CircuitBreaker != nil {
		t.opts = append(t.opts, config.WithCircuitBreaker(*t.CircuitBreaker))
	}

	if t.EnableCompression {
		t.opts = append(t.opts, config.WithCompression())
	}

	if t.Scheme != "" {
		if t.Scheme == "http" {
			t.opts = append(t.opts, config.WithScheme(config.SchemeHTTP))
		} else if t.Scheme == "https" {
			t.opts = append(t.opts, config.WithScheme(config.SchemeHTTPS))
		}
	}

	if t.EnableWebsocketTCPConversion {
		t.opts = append(t.opts, config.WithWebsocketTCPConversion())
	}

	for username, password := range t.BasicAuth {
		t.opts = append(t.opts, config.WithBasicAuth(username, password))
	}

	return nil
}

// convert to ngrok's Tunnel type
func (t *HTTP) NgrokTunnel() config.Tunnel {
	if err := t.ProvisionOpts(); err != nil {
		panic(err)
	}

	return config.HTTPEndpoint(t.opts...)
}

func (t *HTTP) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}

		for nesting := d.Nesting(); d.NextBlock(nesting); {
			subdirective := d.Val()
			switch subdirective {
			case "domain":
				if !d.AllArgs(&t.Domain) {
					return d.ArgErr()
				}
			case "metadata":
				if !d.AllArgs(&t.Metadata) {
					return d.ArgErr()
				}
			case "allow":
				if err := t.unmarshalAllowCidr(d); err != nil {
					return err
				}
			case "deny":
				if err := t.unmarshalDenyCidr(d); err != nil {
					return err
				}
			case "circuit_breaker":
				if err := t.unmarshalCircuitBreaker(d); err != nil {
					return err
				}
			case "enable_compression":
				if err := t.unmarshalEnableCompression(d); err != nil {
					return err
				}
			case "scheme":
				if !d.AllArgs(&t.Scheme) {
					return d.ArgErr()
				}
			case "enable_websocket_tcp_conversion":
				if err := t.unmarshalEnableWebsocketTCPConversion(d); err != nil {
					return err
				}
			case "basicauth":
				if err := t.unmarshalBasicAuth(d); err != nil {
					return err
				}
			default:
				return d.ArgErr()
			}
		}
	}

	return nil
}

func (t *HTTP) unmarshalEnableWebsocketTCPConversion(d *caddyfile.Dispenser) error {
	var value string
	if !d.Args(&value) { // no arg default is true
		t.EnableWebsocketTCPConversion = true
	} else { // arg was given check it
		var err error
		t.EnableWebsocketTCPConversion, err = strconv.ParseBool(value)
		if err != nil {
			return d.Errf(`parsing enable_websocket_tcp_conversion value %+v: %w`, value, err)
		}
	}

	return nil
}

func (t *HTTP) unmarshalEnableCompression(d *caddyfile.Dispenser) error {
	var value string
	if !d.Args(&value) { // no arg default is true
		t.EnableCompression = true
	} else { // arg was given check it
		var err error
		t.EnableCompression, err = strconv.ParseBool(value)
		if err != nil {
			return d.Errf(`parsing enable_compression value %+v: %w`, value, err)
		}
	}

	return nil
}

func (t *HTTP) unmarshalAllowCidr(d *caddyfile.Dispenser) error {
	if d.CountRemainingArgs() == 0 {
		return d.ArgErr()
	}

	t.AllowCIDR = append(t.AllowCIDR, d.RemainingArgs()...)

	return nil
}

func (t *HTTP) unmarshalDenyCidr(d *caddyfile.Dispenser) error {
	if d.CountRemainingArgs() == 0 {
		return d.ArgErr()
	}

	t.DenyCIDR = append(t.DenyCIDR, d.RemainingArgs()...)

	return nil
}

func (t *HTTP) unmarshalBasicAuth(d *caddyfile.Dispenser) error {
	for nesting := d.Nesting(); d.NextBlock(nesting); {
		username := d.Val()
		if username == "}" || username == "{" {
			continue
		}

		var password string

		if !d.AllArgs(&password) {
			return d.ArgErr()
		}

		if username == "" || password == "" {
			return d.Err("username and password cannot be empty or missing")
		}

		minLenPassword := 8
		if len(password) < minLenPassword {
			return d.Err("password must be at least eight characters.")
		}

		if t.BasicAuth == nil {
			t.BasicAuth = map[string]string{}
		}

		t.BasicAuth[username] = password
	}

	return nil
}

func (t *HTTP) unmarshalCircuitBreaker(d *caddyfile.Dispenser) error {
	var ratio string
	if !d.AllArgs(&ratio) {
		return d.ArgErr()
	}

	circuitBreaker, err := strconv.ParseFloat(ratio, 64)
	if err != nil {
		return d.ArgErr()
	}

	t.CircuitBreaker = &circuitBreaker

	return nil
}

var (
	_ caddy.Module          = (*HTTP)(nil)
	_ Tunnel                = (*HTTP)(nil)
	_ caddy.Provisioner     = (*HTTP)(nil)
	_ caddyfile.Unmarshaler = (*HTTP)(nil)
)
