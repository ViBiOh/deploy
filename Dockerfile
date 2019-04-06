FROM golang:1.12 as builder

ENV APP_NAME deploy
ENV WORKDIR ${GOPATH}/src/github.com/ViBiOh/deploy

WORKDIR ${WORKDIR}
COPY ./ ${WORKDIR}/

RUN make ${APP_NAME} \
 && mkdir -p /app \
 && curl -s -o /app/cacert.pem https://curl.haxx.se/ca/cacert.pem \
 && curl -s -o /app/zoneinfo.zip https://raw.githubusercontent.com/golang/go/master/lib/time/zoneinfo.zip \
 && cp bin/${APP_NAME} /app/

FROM docker/compose

ENV APP_NAME deploy
EXPOSE 1080

HEALTHCHECK --retries=10 CMD [ "/deploy", "-url", "https://localhost:1080/health" ]
ENTRYPOINT [ "/deploy" ]

COPY --from=builder /app/cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/zoneinfo.zip /app/${APP_NAME} /
COPY deploy.sh /deploy.sh
