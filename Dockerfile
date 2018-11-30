FROM golang:1.11.2
MAINTAINER sumeshkanayi@gmail.com
ADD Makefile /tmp
CWD /tmp
RUN make docker
