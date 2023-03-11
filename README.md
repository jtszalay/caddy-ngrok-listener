ngrok/Caddy Listener Wrapper
=============================

On March 9, 2023, ngrok announced the release of [ngrok-go](https://blog.ngrok.com/posts/ngrok-go)[^1], a Go package for embedding ngrok into a Go application. The package returns a `net.Listener`. This means it fits right into Caddy's [`listener_wrapper`](https://caddyserver.com/docs/json/apps/http/servers/listener_wrappers/)[^2]. Using this module, when Caddy asks for a listener, it will ask ngrok for the listener, for which ngrok return an ngrok ingress address that is publicly accessible. The public address is printed in logs and avaible on ngrok dashboard.

Currently, the module does not support the extended ngrok options, e.g. allow/deny CIDR. PRs are welcome.

[^1]: [Alan Shreve's tweet](https://twitter.com/inconshreveable/status/1633837669053792260)

[^2]: [`listener_wrappers` Caddyfile docs](https://caddyserver.com/docs/caddyfile/options#listener-wrappers)

## Example

### Caddyfile

```
{
	servers :80 {
		listener_wrappers {
			ngrok {
				auth_token $NGROK_AUTH_TOKEN
				tunnel http {

				}
			}
		}
	}
	servers :8080 {
		listener_wrappers {
			ngrok {
				auth_token $NGROK_AUTH_TOKEN
				tunnel labeled {
					labels {
						edge edghts_LABEL
					}
				}
			}
		}
	}
	servers :8081 {
		listener_wrappers {
			ngrok {
				auth_token $NGROK_AUTH_TOKEN
				tunnel http {
					allow <cidr> ...
					deny <cidr> ...
					domain mydomain.ngrok.io
					metadata bad8c1c0-8fce-11e4-b4a9-0800200c9a66
					scheme http
					circuit_breaker 0.7
					enable_compression
					enable_websocket_tcp_conversion
					basicauth {
						test passw0rd
						test2 password
						user hellohello
					}
				}
			}
		}
	}
}
:80 {
	root * /path/to/site/root
	file_server
}
:8080 {
	reverse_proxy 127.0.0.1:5000
}
:8081 {
	reverse_proxy 127.0.0.1:9090
}
```