FROM golang:1.12 as builder

WORKDIR /app
COPY . .

RUN make

ARG CODECOV_TOKEN
RUN curl -s https://codecov.io/bash | bash

FROM docker/compose:1.25.0-rc1

EXPOSE 1080

RUN apk --update add bash coreutils \
 && rm -rf /var/cache/apk/*

HEALTHCHECK --retries=10 CMD [ "/deploy", "-url", "http://localhost:1080/health" ]
ENTRYPOINT [ "/deploy" ]

ARG APP_VERSION
ENV VERSION=${APP_VERSION}

COPY --from=builder /app/bin/deploy /
COPY deploy.sh /deploy.sh
