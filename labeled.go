package ngroklistener

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"go.uber.org/zap"
	"golang.ngrok.com/ngrok/config"
)

func init() {
	caddy.RegisterModule(new(Labeled))
}

// ngrok Labeled Tunnel
type Labeled struct {
	opts []config.LabeledTunnelOption

	// A map of label, value pairs for this tunnel.
	Labels map[string]string `json:"labels,omitempty"`

	// opaque metadata string for this tunnel.
	Metadata string `json:"metadata,omitempty"`

	l *zap.Logger
}

// CaddyModule implements caddy.Module
func (*Labeled) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "caddy.listeners.ngrok.tunnels.labeled",
		New: func() caddy.Module {
			return new(Labeled)
		},
	}
}

// Provision implements caddy.Provisioner
func (t *Labeled) Provision(ctx caddy.Context) error {
	t.l = ctx.Logger()

	return nil
}

func (t *Labeled) ProvisionOpts() error {
	for label, value := range t.Labels {
		t.opts = append(t.opts, config.WithLabel(label, value))
		t.l.Info("applying label", zap.String("label", label), zap.String("value", value))
	}

	if t.Metadata != "" {
		t.opts = append(t.opts, config.WithMetadata(t.Metadata))
	}

	return nil
}

// convert to ngrok's Tunnel type
func (t *Labeled) NgrokTunnel() config.Tunnel {
	if err := t.ProvisionOpts(); err != nil {
		panic(err)
	}

	return config.LabeledTunnel(t.opts...)
}

func (t *Labeled) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
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
			case "labels":
				for nesting := d.Nesting(); d.NextBlock(nesting); {
					label := d.Val()
					if label == "}" || label == "{" {
						continue
					}

					var labelValue string

					if !d.AllArgs(&labelValue) {
						return d.ArgErr()
					}

					if t.Labels == nil {
						t.Labels = map[string]string{}
					}

					t.Labels[label] = labelValue
				}
			default:
				return d.ArgErr()
			}
		}
	}

	return nil
}

var (
	_ caddy.Module          = (*Labeled)(nil)
	_ Tunnel                = (*Labeled)(nil)
	_ caddy.Provisioner     = (*Labeled)(nil)
	_ caddyfile.Unmarshaler = (*Labeled)(nil)
)
