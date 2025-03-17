FROM ubuntu:22.04 AS builder

RUN apt update && apt install -y wget

RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /build

COPY ./consumers/code_processor/go.mod ./consumers/code_processor/go.sum ./
RUN go mod download

COPY ./consumers/code_processor .
RUN go build -o main ./cmd/code_processor.go

FROM ubuntu:22.04 AS runner

WORKDIR /app

COPY --from=builder /build/main ./main
COPY --from=builder /build/build ./build
COPY --from=builder /build/config.yaml /app/config.yaml

CMD [ "./main" ]
