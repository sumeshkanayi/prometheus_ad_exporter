FROM golang:1.11.2
MAINTAINER SUMESH
ADD Makefile /tmp
RUN make build
