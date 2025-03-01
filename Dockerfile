FROM ubuntu:22.04

RUN apt update && apt install -y wget

RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

COPY go.mod go.sum main.go /

RUN go mod download
RUN go build main.go

CMD [ "./main" ]
