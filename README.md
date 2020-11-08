# deploy

[![Build](https://github.com/ViBiOh/deploy/workflows/Build/badge.svg)](https://github.com/ViBiOh/deploy/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/deploy)](https://goreportcard.com/report/github.com/ViBiOh/deploy)
[![codecov](https://codecov.io/gh/ViBiOh/deploy/branch/master/graph/badge.svg)](https://codecov.io/gh/ViBiOh/deploy)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FViBiOh%2Fdeploy.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2FViBiOh%2Fdeploy?ref=badge_shield)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ViBiOh_deploy&metric=alert_status)](https://sonarcloud.io/dashboard?id=ViBiOh_deploy)

docker-compose deploy API

```bash
curl -O https://raw.githubusercontent.com/ViBiOh/deploy/master/deploy
chmod +x deploy

./deploy PROJECT_NAME DOCKER-COMPOSE-FILE
```

```bash
Usage: deploy [PROJECT_NAME] [DOCKER-COMPOSE-FILE]
  where
    - PROJECT_NAME         Name of your compose project
    - DOCKER_COMPOSE_FILE  Path to your compose file (default: docker-compose.yml in current dir)
```

## Golang API

You can execute the `deploy` script through HTTP API.

```bash
curl -X POST http://localhost:1080/[project_name]/ --data-binary @docker-compose.yml
```

We recommend putting an `Authorization` in front of your server (e.g. reverse-proxy, nginx, etc) if you plan to expose it to the internet.

If something goes wrong during the deploy process, the uploaded `docker-compose.yml` is kept in order to manually retry or debug what's going on. Otherwise, the file is deleted.

### CLI of HTTP Server

```bash
Usage of deploy:
  -address string
        [http] Listen address {DEPLOY_ADDRESS}
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
        [http] Certificate file {DEPLOY_CERT}
  -csp string
        [owasp] Content-Security-Policy {DEPLOY_CSP} (default "default-src 'self'; base-uri 'self'")
  -frameOptions string
        [owasp] X-Frame-Options {DEPLOY_FRAME_OPTIONS} (default "deny")
  -graceDuration string
        [http] Grace duration when SIGTERM received {DEPLOY_GRACE_DURATION} (default "30s")
  -hsts
        [owasp] Indicate Strict Transport Security {DEPLOY_HSTS} (default true)
  -idleTimeout string
        [http] Idle Timeout {DEPLOY_IDLE_TIMEOUT} (default "2m")
  -key string
        [http] Key file {DEPLOY_KEY}
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
  -mailerPass string
        [mailer] Pass {DEPLOY_MAILER_PASS}
  -mailerURL string
        [mailer] URL (an instance of github.com/ViBiOh/mailer) {DEPLOY_MAILER_URL}
  -mailerUser string
        [mailer] User {DEPLOY_MAILER_USER}
  -okStatus int
        [http] Healthy HTTP Status code {DEPLOY_OK_STATUS} (default 204)
  -port uint
        [http] Listen port {DEPLOY_PORT} (default 1080)
  -prometheusIgnore string
        [prometheus] Ignored path prefixes for metrics, comma separated {DEPLOY_PROMETHEUS_IGNORE}
  -prometheusPath string
        [prometheus] Path for exposing metrics {DEPLOY_PROMETHEUS_PATH} (default "/metrics")
  -readTimeout string
        [http] Read Timeout {DEPLOY_READ_TIMEOUT} (default "5s")
  -shutdownTimeout string
        [http] Shutdown Timeout {DEPLOY_SHUTDOWN_TIMEOUT} (default "10s")
  -url string
        [alcotest] URL to check {DEPLOY_URL}
  -userAgent string
        [alcotest] User-Agent for check {DEPLOY_USER_AGENT} (default "Alcotest")
  -writeTimeout string
        [http] Write Timeout {DEPLOY_WRITE_TIMEOUT} (default "2m")
```

## License

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FViBiOh%2Fdeploy.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2FViBiOh%2Fdeploy?ref=badge_large)
