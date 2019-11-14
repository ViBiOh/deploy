FROM golang:1.13 as builder

WORKDIR /app
COPY . .

RUN make \
 && git diff -- *.go \
 && git diff --quiet -- *.go

ARG CODECOV_TOKEN
RUN curl -q -sSL --max-time 10 https://codecov.io/bash | bash

FROM docker/compose:1.24.1

EXPOSE 1080

RUN apk --update add bash coreutils ca-certificates \
 && rm -rf /var/cache/apk/*

HEALTHCHECK --retries=10 CMD [ "/deploy", "-url", "http://localhost:1080/health" ]
ENTRYPOINT [ "/deploy" ]

ARG APP_VERSION
ENV VERSION=${APP_VERSION}

COPY --from=builder /app/bin/deploy /
COPY deploy /deploy
