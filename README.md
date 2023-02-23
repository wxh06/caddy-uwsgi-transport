# caddy-uwsgi-transport

This module adds [uwsgi](https://uwsgi-docs.readthedocs.io/en/latest/Protocol.html) reverse proxying support to Caddy.

ID: `http.reverse_proxy.transport.uwsgi`

## Installation

By using [`xcaddy`](https://caddyserver.com/docs/build#xcaddy)

```sh
xcaddy build \
    --with github.com/wxh06/caddy-uwsgi-transport
```

## Usage

### Caddyfile

```caddyfile
reverse_proxy [<matcher>] [<upstreams...>] {
	transport uwsgi
}
```

### JSON

```json
{
  "apps": {
    "http": {
      "servers": {
        "": {
          "routes": [
            {
              "handle": [
                {
                  "handler": "reverse_proxy",
                  "transport": { "protocol": "uwsgi" },
                  "upstreams": [{ "dial": "<upstream>" }]
                }
              ]
            }
          ]
        }
      }
    }
  }
}
```
