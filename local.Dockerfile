FROM golang:1.23 AS builder

WORKDIR /build

COPY . /build

RUN go mod download
RUN CGO_ENABLED=0 go build -a -o adguard-exporter main.go

FROM alpine:3.21.0

RUN apk add ca-certificates curl --no-cache
ARG SREP_VERSION
ENV SREP_VERSION ${SREP_VERSION}
ENV PORT 9618

COPY --from=builder /build/adguard-exporter /adguard-exporter

HEALTHCHECK --interval=2s --timeout=5s --retries=5 \
  CMD curl --fail http://localhost:$PORT/health || exit 1

ENTRYPOINT [ "/adguard-exporter" ]
