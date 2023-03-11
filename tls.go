package ngroklistener

import (
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

	// opaque metadata string for this tunnel.
	Metadata string `json:"metadata,omitempty"`

	// Rejects connections that do not match the given CIDRs
	AllowCIDR []string `json:"allowCidr,omitempty"`

	// Rejects connections that match the given CIDRs and allows all other CIDRs.
	DenyCIDR []string `json:"denyCidr,omitempty"`

	l *zap.Logger
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
func (t *TLS) Provision(ctx caddy.Context) error {
	t.l = ctx.Logger()

	return nil
}

func (t *TLS) ProvisionOpts() error {
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

	return nil
}

// convert to ngrok's Tunnel type
func (t *TLS) NgrokTunnel() config.Tunnel {
	if err := t.ProvisionOpts(); err != nil {
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
					return d.ArgErr()
				}
			case "metadata":
				if !d.AllArgs(&t.Metadata) {
					return d.ArgErr()
				}
			case "allow":
				if d.CountRemainingArgs() == 0 {
					return d.ArgErr()
				}

				t.AllowCIDR = append(t.AllowCIDR, d.RemainingArgs()...)
			case "deny":
				if d.CountRemainingArgs() == 0 {
					return d.ArgErr()
				}

				t.DenyCIDR = append(t.DenyCIDR, d.RemainingArgs()...)
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
