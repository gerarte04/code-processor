FROM ubuntu:22.04 as builder

RUN apt update && apt install -y wget

RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build main.go

FROM ubuntu:22.04 as runner

COPY --from=builder /build/main /app/main

CMD [ "/app/main" ]
