FROM docker/compose:1.27.4

EXPOSE 1080

RUN apk --update add bash coreutils ca-certificates \
 && rm -rf /var/cache/apk/*

HEALTHCHECK --retries=10 CMD [ "/deploy", "-url", "http://localhost:1080/health" ]
ENTRYPOINT [ "/deploy" ]

ARG APP_VERSION
ENV VERSION=${APP_VERSION}

COPY --from=builder /app/bin/deploy /
COPY deploy /deploy
COPY clean /clean
