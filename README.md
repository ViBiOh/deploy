# deploy

[![Build](https://github.com/ViBiOh/deploy/workflows/Build/badge.svg)](https://github.com/ViBiOh/deploy/actions)
[![codecov](https://codecov.io/gh/ViBiOh/deploy/branch/main/graph/badge.svg)](https://codecov.io/gh/ViBiOh/deploy)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ViBiOh_deploy&metric=alert_status)](https://sonarcloud.io/dashboard?id=ViBiOh_deploy)

docker-compose deploy API

```bash
curl -O https://raw.githubusercontent.com/ViBiOh/deploy/main/pkg/api/scripts/deploy-compose
chmod +x deploy

./deploy PROJECT_NAME DOCKER-COMPOSE-FILE
```

```bash
Usage: deploy [PROJECT_NAME] [DOCKER-COMPOSE-FILE]
  where
    - PROJECT_NAME         Name of your compose project
    - DOCKER_COMPOSE_FILE  Path to your compose file (default: docker-compose.yaml in current dir)
```

## Golang API

You can execute the `deploy` script through HTTP API.

```bash
curl -X POST http://localhost:1080/[project_name]/ --data-binary @docker-compose.yaml
```

We recommend putting an `Authorization` in front of your server (e.g. reverse-proxy, nginx, etc) if you plan to expose it to the internet.

If something goes wrong during the deploy process, the uploaded `docker-compose.yaml` is kept in order to manually retry or debug what's going on. Otherwise, the file is deleted.

### CLI of HTTP Server

```bash
Usage of deploy:
  -address string
        [server] Listen address {DEPLOY_ADDRESS}
  -annotationPass string
        [annotation] Pass {DEPLOY_ANNOTATION_PASS}
  -annotationURL string
        [annotation] URL of Annotation server (e.g. my.grafana.com/api/annotations) {DEPLOY_ANNOTATION_URL}
  -annotationUser string
        [annotation] User {DEPLOY_ANNOTATION_USER}
  -apiNotification string
        [api] Email notificiation when deploy ends (possibles values ares 'never', 'onError', 'all') {DEPLOY_API_NOTIFICATION} (default "onError")
  -apiNotificationEmail string
        [api] Email address to notify {DEPLOY_API_NOTIFICATION_EMAIL}
  -apiTempFolder string
        [api] Temp folder for uploading files {DEPLOY_API_TEMP_FOLDER} (default "/tmp")
  -cert string
        [server] Certificate file {DEPLOY_CERT}
  -csp string
        [owasp] Content-Security-Policy {DEPLOY_CSP} (default "default-src 'self'; base-uri 'self'")
  -frameOptions string
        [owasp] X-Frame-Options {DEPLOY_FRAME_OPTIONS} (default "deny")
  -graceDuration string
        [http] Grace duration when SIGTERM received {DEPLOY_GRACE_DURATION} (default "30s")
  -hsts
        [owasp] Indicate Strict Transport Security {DEPLOY_HSTS} (default true)
  -idleTimeout string
        [server] Idle Timeout {DEPLOY_IDLE_TIMEOUT} (default "2m")
  -key string
        [server] Key file {DEPLOY_KEY}
  -loggerJson
        [logger] Log format as JSON {DEPLOY_LOGGER_JSON}
  -loggerLevel string
        [logger] Logger level {DEPLOY_LOGGER_LEVEL} (default "INFO")
  -loggerLevelKey string
        [logger] Key for level in JSON {DEPLOY_LOGGER_LEVEL_KEY} (default "level")
  -loggerMessageKey string
        [logger] Key for message in JSON {DEPLOY_LOGGER_MESSAGE_KEY} (default "message")
  -loggerTimeKey string
        [logger] Key for timestamp in JSON {DEPLOY_LOGGER_TIME_KEY} (default "time")
  -mailerName string
        [mailer] HTTP Username or AMQP Exchange name {DEPLOY_MAILER_NAME} (default "mailer")
  -mailerPassword string
        [mailer] HTTP Pass {DEPLOY_MAILER_PASSWORD}
  -mailerURL string
        [mailer] URL (https?:// or amqps?://) {DEPLOY_MAILER_URL}
  -okStatus int
        [http] Healthy HTTP Status code {DEPLOY_OK_STATUS} (default 204)
  -port uint
        [server] Listen port (0 to disable) {DEPLOY_PORT} (default 1080)
  -prometheusAddress string
        [prometheus] Listen address {DEPLOY_PROMETHEUS_ADDRESS}
  -prometheusCert string
        [prometheus] Certificate file {DEPLOY_PROMETHEUS_CERT}
  -prometheusGzip
        [prometheus] Enable gzip compression of metrics output {DEPLOY_PROMETHEUS_GZIP}
  -prometheusIdleTimeout string
        [prometheus] Idle Timeout {DEPLOY_PROMETHEUS_IDLE_TIMEOUT} (default "10s")
  -prometheusIgnore string
        [prometheus] Ignored path prefixes for metrics, comma separated {DEPLOY_PROMETHEUS_IGNORE}
  -prometheusKey string
        [prometheus] Key file {DEPLOY_PROMETHEUS_KEY}
  -prometheusPort uint
        [prometheus] Listen port (0 to disable) {DEPLOY_PROMETHEUS_PORT} (default 9090)
  -prometheusReadTimeout string
        [prometheus] Read Timeout {DEPLOY_PROMETHEUS_READ_TIMEOUT} (default "5s")
  -prometheusShutdownTimeout string
        [prometheus] Shutdown Timeout {DEPLOY_PROMETHEUS_SHUTDOWN_TIMEOUT} (default "5s")
  -prometheusWriteTimeout string
        [prometheus] Write Timeout {DEPLOY_PROMETHEUS_WRITE_TIMEOUT} (default "10s")
  -readTimeout string
        [server] Read Timeout {DEPLOY_READ_TIMEOUT} (default "5s")
  -shutdownTimeout string
        [server] Shutdown Timeout {DEPLOY_SHUTDOWN_TIMEOUT} (default "10s")
  -url string
        [alcotest] URL to check {DEPLOY_URL}
  -userAgent string
        [alcotest] User-Agent for check {DEPLOY_USER_AGENT} (default "Alcotest")
  -writeTimeout string
        [server] Write Timeout {DEPLOY_WRITE_TIMEOUT} (default "2m")
```
