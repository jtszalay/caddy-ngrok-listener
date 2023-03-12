package ngroklistener

import (
	"fmt"

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
	opts []config.TCPEndpointOption

	// The remote TCP address to request for this edge
	RemoteAddr string `json:"remoteAddr,omitempty"`

	// opaque metadata string for this tunnel.
	Metadata string `json:"metadata,omitempty"`

	// Rejects connections that do not match the given CIDRs
	AllowCIDR []string `json:"allowCidr,omitempty"`

	// Rejects connections that match the given CIDRs and allows all other CIDRs.
	DenyCIDR []string `json:"denyCidr,omitempty"`

	l *zap.Logger
}

// Provision implements caddy.Provisioner
func (t *TCP) Provision(ctx caddy.Context) error {
	t.l = ctx.Logger()

	err := t.DoReplace()
	if err != nil {
		return fmt.Errorf("loading doing replacements: %v", err)
	}

	return nil
}

func (t *TCP) ProvisionOpts() error {
	t.opts = append(t.opts, config.WithRemoteAddr(t.RemoteAddr))

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

func (t *TCP) DoReplace() error {
	repl := caddy.NewReplacer()
	replaceableFields := []*string{
		&t.Metadata,
	}

	for _, field := range replaceableFields {
		actual, err := repl.ReplaceOrErr(*field, false, true)
		if err != nil {
			return fmt.Errorf("error replacing fields: %v", err)
		}

		*field = actual
	}

	replacedAllowCIDR := make([]string, len(t.AllowCIDR))

	for index, cidr := range t.AllowCIDR {
		actual, err := repl.ReplaceOrErr(cidr, false, true)
		if err != nil {
			return fmt.Errorf("error replacing fields: %v", err)
		}

		replacedAllowCIDR[index] = actual
	}

	t.AllowCIDR = replacedAllowCIDR

	replacedDenyCIDR := make([]string, len(t.DenyCIDR))

	for index, cidr := range t.DenyCIDR {
		actual, err := repl.ReplaceOrErr(cidr, false, true)
		if err != nil {
			return fmt.Errorf("error replacing fields: %v", err)
		}

		replacedDenyCIDR[index] = actual
	}

	t.DenyCIDR = replacedDenyCIDR

	return nil
}

// convert to ngrok's Tunnel type
func (t *TCP) NgrokTunnel() config.Tunnel {
	if err := t.ProvisionOpts(); err != nil {
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
			case "metadata":
				if !d.AllArgs(&t.Metadata) {
					return d.ArgErr()
				}
			case "remote_addr":
				if !d.AllArgs(&t.RemoteAddr) {
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
	_ caddy.Module          = (*TCP)(nil)
	_ Tunnel                = (*TCP)(nil)
	_ caddy.Provisioner     = (*TCP)(nil)
	_ caddyfile.Unmarshaler = (*TCP)(nil)
)
