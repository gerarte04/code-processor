FROM ubuntu:22.04 AS builder

RUN apt update && apt install -y wget

RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /build

COPY http_server/go.mod http_server/go.sum ./
RUN go mod download

COPY http_server .
RUN go build cmd/main.go

FROM ubuntu:22.04 AS runner

COPY --from=builder /build/main /app/main
COPY --from=builder /build/config.yaml /app/config.yaml

CMD [ "/app/main", "--config=/app/config.yaml" ]
