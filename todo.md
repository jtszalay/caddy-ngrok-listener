* additional Replacers should be used where appropriate. 

* The following config options are not yet surfaced.

```golang

// https://github.com/ngrok/ngrok-go/blob/main/session.go
// WithProxyURL configures the session to connect to ngrok through an outbound
// HTTP or SOCKS5 proxy. This parameter is ignored if you override the dialer
// with [WithDialer].
//
// See the [proxy url paramter in the ngrok docs] for additional details.
//
// [proxy url paramter in the ngrok docs]: https://ngrok.com/docs/ngrok-agent/config#proxy_url
func WithProxyURL(url *url.URL) ConnectOption


// WithCA configures the CAs used to validate the TLS certificate returned by
// the ngrok service while establishing the session. Use this option only if
// you are connecting through a man-in-the-middle or deep packet inspection
// proxy.
//
// See the [root_cas parameter in the ngrok docs] for additional details.
//
// [root_cas parameter in the ngrok docs]: https://ngrok.com/docs/ngrok-agent/config#root_cas
func WithCA(pool *x509.CertPool) ConnectOption

// https://github.com/ngrok/ngrok-go/blob/main/config/tls_termination.go
// WithTermination sets the key and certificate in PEM format for TLS termination at the ngrok
// edge.
func WithTermination(certPEM, keyPEM []byte) TLSEndpointOption 


// https://github.com/ngrok/ngrok-go/blob/main/config/mutual_tls.go

type mutualTLSEndpointOption []*x509.Certificate

// WithMutualTLSCA adds a list of [x509.Certificate]'s to use for mutual TLS
// authentication.
// These will be used to authenticate client certificates for requests at the
// ngrok edge.
func WithMutualTLSCA(certs ...*x509.Certificate) interface {
	HTTPEndpointOption
	TLSEndpointOption
}


https://github.com/ngrok/ngrok-go/blob/main/config/http_headers.go
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


// https://github.com/ngrok/ngrok-go/blob/main/config/oauth.go
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


// https://github.com/ngrok/ngrok-go/blob/main/config/oidc.go
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


// https://github.com/ngrok/ngrok-go/blob/main/config/webhook_verification.go
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




Future adds
```golang

// WithConnectHandler configures a function which is called each time the ngrok
// [Session] successfully connects to the ngrok service. Use this option to
// receive events when ngrok successfully reconnects a [Session] that was
// disconnected because of a network failure.
func WithConnectHandler(handler SessionConnectHandler) ConnectOption {
	return func(cfg *connectConfig) {
		cfg.ConnectHandler = handler
	}
}

// WithDisconnectHandler configures a function which is called each time the
// ngrok [Session] disconnects from the ngrok service. Use this option to detect
// when the ngrok session has gone temporarily offline.
func WithDisconnectHandler(handler SessionDisconnectHandler) ConnectOption {
	return func(cfg *connectConfig) {
		cfg.DisconnectHandler = handler
	}
}

// WithHeartbeatHandler configures a function which is called each time the
// [Session] successfully heartbeats the ngrok service. The callback receives
// the latency of the round trip time from initiating the heartbeat to
// receiving an acknowledgement back from the ngrok service.
func WithHeartbeatHandler(handler SessionHeartbeatHandler) ConnectOption {
	return func(cfg *connectConfig) {
		cfg.HeartbeatHandler = handler
	}
}

// WithStopHandler configures a function which is called when the ngrok service
// requests that this [Session] stops. Your application may choose to interpret
// this callback as a request to terminate the [Session] or the entire process.
//
// Errors returned by this function will be visible to the ngrok dashboard or
// API as the response to the Stop operation.
//
// Do not block inside this callback. It will cause the Dashboard or API stop
// operation to hang. Do not call [Session].Close or [os.Exit] inside this
// callback, it will also cause the operation to hang.
//
// Instead, either return an error or if you intend to Stop, spawn a goroutine
// to asynchronously call [Session].Close or [os.Exit].
func WithStopHandler(handler ServerCommandHandler) ConnectOption {
	return func(cfg *connectConfig) {
		cfg.StopHandler = handler
	}
}
```


* Test plan:
	For each tunnel type test the following:
	* parse
	* provision
	* validate

	For ngrok.go test:
	* parse of toplevel options
	* basic parse of each tunnel type
	* provision
	* validate