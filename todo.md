The following config options are not yet surfaced.

```golang
// WithProxyURL configures the session to connect to ngrok through an outbound
// HTTP or SOCKS5 proxy. This parameter is ignored if you override the dialer
// with [WithDialer].
//
// See the [proxy url paramter in the ngrok docs] for additional details.
//
// [proxy url paramter in the ngrok docs]: https://ngrok.com/docs/ngrok-agent/config#proxy_url
func WithProxyURL(url *url.URL) ConnectOption

// WithServer configures the network address to dial to connect to the ngrok
// service. Use this option only if you are connecting to a custom agent
// ingress.
//
// See the [server_addr parameter in the ngrok docs] for additional details.
//
// [server_addr parameter in the ngrok docs]: https://ngrok.com/docs/ngrok-agent/config#server_addr
func WithServer(addr string) ConnectOption
// WithCA configures the CAs used to validate the TLS certificate returned by
// the ngrok service while establishing the session. Use this option only if
// you are connecting through a man-in-the-middle or deep packet inspection
// proxy.
//
// See the [root_cas parameter in the ngrok docs] for additional details.
//
// [root_cas parameter in the ngrok docs]: https://ngrok.com/docs/ngrok-agent/config#root_cas
func WithCA(pool *x509.CertPool) ConnectOption
// WithHeartbeatTolerance configures the duration to wait for a response to a heartbeat
// before assuming the session connection is dead and attempting to reconnect.
//
// See the [heartbeat_tolerance parameter in the ngrok docs] for additional details.
//
// [heartbeat_tolerance parameter in the ngrok docs]: https://ngrok.com/docs/ngrok-agent/config#heartbeat_tolerance
func WithHeartbeatTolerance(tolerance time.Duration) ConnectOption
// WithHeartbeatInterval configures how often the session will send heartbeat
// messages to the ngrok service to check session liveness.
//
// See the [heartbeat_interval parameter in the ngrok docs] for additional details.
//
// [heartbeat_interval parameter in the ngrok docs]: https://ngrok.com/docs/ngrok-agent/config#heartbeat_interval
func WithHeartbeatInterval(interval time.Duration) ConnectOption 

// HTTP Headers to modify at the ngrok edge.
type headers struct {
	// Headers to add to requests or responses at the ngrok edge.
	Added map[string]string
	// Header names to remove from requests or responses at the ngrok edge.
	Removed []string
}
// WithRequestHeader adds a header to all requests to this edge.
func WithRequestHeader(name, value string) HTTPEndpointOption
// WithRequestHeader adds a header to all responses coming from this edge.
func WithResponseHeader(name, value string) HTTPEndpointOption 
// WithRemoveRequestHeader removes a header from requests to this edge.
func WithRemoveRequestHeader(name string) HTTPEndpointOption
// WithRemoveResponseHeader removes a header from responses from this edge.
func WithRemoveResponseHeader(name string) HTTPEndpointOption



type mutualTLSEndpointOption []*x509.Certificate

// WithMutualTLSCA adds a list of [x509.Certificate]'s to use for mutual TLS
// authentication.
// These will be used to authenticate client certificates for requests at the
// ngrok edge.
func WithMutualTLSCA(certs ...*x509.Certificate) interface {
	HTTPEndpointOption
	TLSEndpointOption
}


// oauthOptions configuration
type oauthOptions struct {
	// The OAuth provider to use
	Provider string
	// Email addresses of users to authorize.
	AllowEmails []string
	// Email domains of users to authorize.
	AllowDomains []string
	// OAuth scopes to request from the provider.
	Scopes []string
}
// Append email addresses to the list of allowed emails.
func WithAllowOAuthEmail(addr ...string) OAuthOption
// Append email domains to the list of allowed domains.
func WithAllowOAuthDomain(domain ...string) OAuthOption
// Append scopes to the list of scopes to request.
func WithOAuthScope(scope ...string) OAuthOption
// WithOAuth configures this edge with the the given OAuth provider.
// Overwrites any previously-set OAuth configuration.
func WithOAuth(provider string, opts ...OAuthOption) HTTPEndpointOption


type oidcOptions struct {
	IssuerURL    string
	ClientID     string
	ClientSecret string
	AllowEmails  []string
	AllowDomains []string
	Scopes       []string
}
// Append scopes to the list of scopes to request.
func WithOIDCScope(scope ...string) OIDCOption
// Append email domains to the list of allowed domains.
func WithAllowOIDCDomain(domain ...string) OIDCOption

// Append email addresses to the list of allowed emails.
func WithAllowOIDCEmail(addr ...string) OIDCOption 
// WithOIDC configures this edge with the the given OIDC provider.
// Overwrites any previously-set OIDC configuration.
func WithOIDC(issuerURL string, clientID string, clientSecret string, opts ...OIDCOption) HTTPEndpointOption 






// WithTermination sets the key and certificate in PEM format for TLS termination at the ngrok
// edge.
func WithTermination(certPEM, keyPEM []byte) TLSEndpointOption 





// Configuration for webhook verification.
type webhookVerification struct {
	// The webhook provider
	Provider string
	// The secret for verifying webhooks from this provider.
	Secret string
}
// WithWebhookVerification configures webhook vericiation for this edge.
func WithWebhookVerification(provider string, secret string) HTTPEndpointOption
```