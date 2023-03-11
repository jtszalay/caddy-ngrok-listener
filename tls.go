package ngroklistener

import (
	"context"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"go.uber.org/zap"
	"golang.ngrok.com/ngrok/config"
)

func init() {
	caddy.RegisterModule(new(TLS))
}

// ngrok TLS tunnel
// Note: only available for ngrok Enterprise user
type TLS struct {
	opts []config.TLSEndpointOption

	// the domain for this edge.
	Domain string `json:"domain,omitempty"`

	ctx context.Context
	l   *zap.Logger
}

// CaddyModule implements caddy.Module
func (*TLS) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "caddy.listeners.ngrok.tunnels.tls",
		New: func() caddy.Module {
			return new(TLS)
		},
	}
}

// Provision implements caddy.Provisioner
func (t *TLS) Provision(caddy.Context) error {
	t.ctx = ctx
	t.l = ctx.Logger()
	return nil
}

func (t *TLS) ProvisionOpts() error {
	if t.Domain != "" {
		t.opts = append(t.opts, config.WithDomain(t.Domain))
	}
	return nil
}

// convert to ngrok's Tunnel type
func (t *TLS) NgrokTunnel() config.Tunnel {
	err := t.ProvisionOpts()
	if err != nil {
		panic(err)
	}
	return config.TLSEndpoint(t.opts...)
}

func (t *TLS) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			subdirective := d.Val()
			switch subdirective {
			case "domain":
				if !d.AllArgs(&t.Domain) {
					d.ArgErr()
				}
			default:
				return d.ArgErr()
			}
		}
	}
	return nil
}

var (
	_ caddy.Module          = (*TLS)(nil)
	_ Tunnel                = (*TLS)(nil)
	_ caddy.Provisioner     = (*TLS)(nil)
	_ caddyfile.Unmarshaler = (*TLS)(nil)
)
