FROM ubuntu:22.04 AS builder

RUN apt update && apt install -y wget

RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /build
ARG service_path="consumers/code_processor"

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ${service_path}/cmd/main.go

FROM ubuntu:22.04 AS runner

WORKDIR /app
ARG service_path="consumers/code_processor"

COPY --from=builder /build/main ./main
COPY --from=builder /build/${service_path}/build ./build
COPY --from=builder /build/${service_path}/config.yaml /app/config.yaml

CMD [ "./main" ]
