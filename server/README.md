The Heartbleed test server. With caching by Mozilla.

Install the `SAMPLE.aws-config.json` as `~/.aws-config.json` or `/etc/aws-config.json` or in a path specified by the `GODYNAMO_CONF_FILE` env var.

```
Usage:
  HBserver --redir-host=<host> [--listen=<addr:port> --expiry=<duration>]
           [--key=<key> --cert=<cert>]
  HBserver -h | --help
  HBserver --version

Options:
  --redir-host HOST   Redirect requests to "/" to this host.
  --listen ADDR:PORT  Listen and serve requests to this address:port [default: :8082].
  --expiry DURATION   ENABLE CACHING. Expire records after this period.
                      Uses Go's parse syntax
                      e.g. 10m = 10 minutes, 600s = 600 seconds, 1d = 1 day, etc.
  --key KEY           TLS key .pem file -- enable TLS
  --cert CERT         TLS cert .pem file -- enable TLS
  -h --help           Show this screen.
  --version           Show version.
```
