FROM alpine:3.21.3 AS builder

RUN apk update && apk --no-cache add bash build-base clang python3 py3-pip

ARG file
ARG translator
ENV FILE_NAME=${file}
ENV TRANSLATOR=${translator}

WORKDIR /build
COPY . .

CMD [ "/bin/bash", "-c", "./run.sh" ]
