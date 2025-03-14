FROM ubuntu:22.04 AS builder

RUN apt update && apt install -y wget

RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /build

COPY ./consumers/code_processor/. .
# RUN go install github.com/rabbitmq/amqp091-go
RUN go build -o main code_processor.go

FROM ubuntu:22.04 AS runner

COPY --from=builder /build/main /app/main

CMD [ "/app/main" ]
