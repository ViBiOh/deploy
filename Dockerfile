FROM golang:1.12 as builder

ENV APP_NAME deploy
ENV WORKDIR ${GOPATH}/src/github.com/ViBiOh/deploy

WORKDIR ${WORKDIR}
COPY ./ ${WORKDIR}/

RUN make ${APP_NAME} \
 && mkdir -p /app \
 && cp bin/${APP_NAME} /app/

FROM docker/compose:1.24.0

ENV APP_NAME deploy
EXPOSE 1080

HEALTHCHECK --retries=10 CMD [ "/deploy", "-url", "https://localhost:1080/health" ]
ENTRYPOINT [ "/deploy" ]

COPY --from=builder /app/${APP_NAME} /
COPY deploy.sh /deploy.sh
