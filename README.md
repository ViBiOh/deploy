# deploy

Zero-downtime docker-compose deploy

```bash
curl -O https://raw.githubusercontent.com/ViBiOh/docker-compose-deploy/master/deploy.sh
chmod +x deploy.sh

./deploy.sh PROJECT_NAME SHA1 DOCKER-COMPOSE-FILE
```

```bash
Usage: deploy [PROJECT_NAME] [SHA1] [DOCKER-COMPOSE-FILE]
  where
    - PROJECT_NAME         Name of your compose project
    - SHA1                 Unique identifier of your project (default: git sha1 of commit)
    - DOCKER_COMPOSE_FILE  Path to your compose file (default: docker-compose.yml in current dir)
```

## Golang API

You can execute the `deploy.sh` script through HTTP API.

```bash
curl -X POST http://localhost:1080/[project_name]/[sha1_version] --data-binary @docker-compose.yml
```

We recommend putting an `Authorization` in front of your server (e.g. reverse-proxy, nginx, etc) if you plan to expose it to the internet.

If something goes wrong during the deploy process, the uploaded `docker-compose.yml` is kept in order to manually retry ou debug what's going on. Otherwise, the file is deleted.

### CLI og HTTP Server

```bash
Usage of deploy:
Usage of deploy:
  -apiNotification string
        [api] Email notificiation when deploy ends (possibles values ares 'never', 'onError', 'all') (default "onError")
  -apiNotificationEmail string
        [api] Email address to notify
  -apiTempFolder string
        [api] Temp folder for uploading files (default "/tmp")
  -csp string
        [owasp] Content-Security-Policy (default "default-src 'self'; base-uri 'self'")
  -frameOptions string
        [owasp] X-Frame-Options (default "deny")
  -hsts
        [owasp] Indicate Strict Transport Security (default true)
  -mailerPass string
        [mailer] Mailer Pass
  -mailerURL string
        [mailer] Mailer URL (default "https://mailer.vibioh.fr")
  -mailerUser string
        [mailer] Mailer User
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

## Containers

We provide a `docker-compose.yml` which contains the HTTP API and a [Portainer](https://www.portainer.io) container for having a Docker GUI.
