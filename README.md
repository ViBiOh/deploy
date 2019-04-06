# deploy

Zero-downtime docker-compose deploy

```bash
curl -O https://raw.githubusercontent.com/ViBiOh/docker-compose-deploy/master/deploy.sh
chmod +x deploy.sh

./deploy.sh awesome_project sha1_default_to_git_sha path_to_your_compose_default_pwd
```

## Golang API

```bash
Usage of deploy:
  -apiTempFolder string
        [api] Temp folder for uploading files
  -csp string
        [owasp] Content-Security-Policy (default "default-src 'self'; base-uri 'self'")
  -frameOptions string
        [owasp] X-Frame-Options (default "deny")
  -hsts
        [owasp] Indicate Strict Transport Security (default true)
  -port int
        Listen port (default 1080)
  -prometheusPath string
        [prometheus] Path for exposing metrics (default "/metrics")
  -tls
        Serve TLS content (default true)
  -tlsCert string
        [tls] PEM Certificate file
  -tlsHosts string
        [tls] Self-signed certificate hosts, comma separated (default "localhost")
  -tlsKey string
        [tls] PEM Key file
  -tlsOrganization string
        [tls] Self-signed certificate organization (default "ViBiOh")
  -tracingAgent string
        [opentracing] Jaeger Agent (e.g. host:port) (default "jaeger:6831")
  -tracingName string
        [opentracing] Service name
  -url string
        [health] URL to check
  -userAgent string
        [health] User-Agent for check (default "Golang alcotest")
```