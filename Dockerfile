FROM alpine:3.21.0
ENV PORT 9618

RUN apk add ca-certificates curl --no-cache

WORKDIR /

COPY adguard-exporter /adguard-exporter
USER 65532:65532
HEALTHCHECK --interval=2s --timeout=5s --retries=5 \
  CMD curl --fail http://localhost:$PORT/health || exit 1

ENTRYPOINT ["/adguard-exporter"]
