package ngroklistener

import (
	"context"

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

	// the domain for this edge.
	Domain string `json:"domain,omitempty"`

	// opaque metadata string for this tunnel.
	Metadata string `json:"metadata,omitempty"`

	ctx context.Context
	l   *zap.Logger
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
func (t *HTTP) Provision(caddy.Context) error {
	t.ctx = ctx
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

	return nil
}

// convert to ngrok's Tunnel type
func (t *HTTP) NgrokTunnel() config.Tunnel {
	err := t.ProvisionOpts()
	if err != nil {
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
					d.ArgErr()
				}
			case "metadata":
				if !d.AllArgs(&t.Metadata) {
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
	_ caddy.Module          = (*HTTP)(nil)
	_ Tunnel                = (*HTTP)(nil)
	_ caddy.Provisioner     = (*HTTP)(nil)
	_ caddyfile.Unmarshaler = (*HTTP)(nil)
)
