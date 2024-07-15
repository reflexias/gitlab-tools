FROM golang:latest as builder

ADD . /app
WORKDIR /app
RUN 

FROM golang:latest