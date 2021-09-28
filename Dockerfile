FROM golang:1.16 AS builder

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apt-get update
RUN apt-get install -y build-essential git unzip curl
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" 

RUN unzip awscliv2.zip
RUN ./aws/install
RUN aws --version

WORKDIR /var/app

COPY . .
