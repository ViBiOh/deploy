FROM golang:1.12 as builder

ARG CODECOV_TOKEN
ENV APP_NAME deploy

WORKDIR /app
COPY . .

RUN make ${APP_NAME} \
 && bash <(curl -s https://codecov.io/bash)

FROM docker/compose:1.24.0

ENV APP_NAME deploy
EXPOSE 1080

RUN apk --update add bash coreutils \
 && rm -rf /var/cache/apk/*

HEALTHCHECK --retries=10 CMD [ "/deploy", "-url", "http://localhost:1080/health" ]
ENTRYPOINT [ "/deploy" ]

COPY --from=builder /app/bin/${APP_NAME} /
COPY deploy.sh /deploy.sh
