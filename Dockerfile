FROM --platform=$BUILDPLATFORM golang:1.19 AS builder

ARG SKAFFOLD_GO_GCFLAGS
ARG TARGETOS
ARG TARGETARCH

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
ENV GO_BIN=/go/bin/app
ENV GRPC_HEALTH_PROBE_VERSION=v0.3.6

RUN apt-get update
RUN apt-get install -y build-essential git unzip curl wget file default-jre python3-pip python3

RUN \
  if [ "${TARGETARCH}" = "amd64" ]; \ 
  then curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" ; \
  else curl "https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip" -o "awscliv2.zip" ; \ 
  fi

RUN unzip awscliv2.zip
RUN ./aws/install

RUN wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-${GOOS}-${GOARCH} && \
  chmod +x /bin/grpc_health_probe

RUN curl https://downloads.apache.org/kafka/3.3.2/kafka_2.13-3.3.2.tgz -o "/tmp/kafka.tgz"
RUN tar -xzf /tmp/kafka.tgz -C /tmp

WORKDIR /var/app

COPY . .
