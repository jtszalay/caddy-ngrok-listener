package ngroklistener

import (
	"context"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"go.uber.org/zap"
	"golang.ngrok.com/ngrok/config"
)

func init() {
	caddy.RegisterModule(new(TCP))
}

// ngrok TCP tunnel
type TCP struct {
	// The remote TCP address to request for this edge
	RemoteAddr string `json:"remote_address,omitempty"`

	opts []config.TCPEndpointOption

	ctx context.Context
	l   *zap.Logger
}

// Provision implements caddy.Provisioner
func (t *TCP) Provision(caddy.Context) error {
	t.ctx = ctx
	t.l = ctx.Logger()

	return nil
}

func (t *TCP) ProvisionOpts() error {
	t.opts = append(t.opts, config.WithRemoteAddr(t.RemoteAddr))
	return nil
}

// convert to ngrok's Tunnel type
func (t *TCP) NgrokTunnel() config.Tunnel {
	err := t.ProvisionOpts()
	if err != nil {
		panic(err)
	}
	return config.TCPEndpoint(t.opts...)
}

// CaddyModule implements caddy.Module
func (*TCP) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "caddy.listeners.ngrok.tunnels.tcp",
		New: func() caddy.Module {
			return new(TCP)
		},
	}
}

func (t *TCP) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			subdirective := d.Val()
			switch subdirective {
			case "remote_address":
				if !d.AllArgs(&t.RemoteAddr) {
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
	_ caddy.Module          = (*TCP)(nil)
	_ Tunnel                = (*TCP)(nil)
	_ caddy.Provisioner     = (*TCP)(nil)
	_ caddyfile.Unmarshaler = (*TCP)(nil)
)
