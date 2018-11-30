FROM golang:1.11.2
MAINTAINER sumeshkanayi@gmail.com
ADD Makefile /tmp
WORKDIR /tmp
ENV https_proxy http://172.17.0.1:3128
ENV http_proxy http://172.17.0.1:3128
RUN make build
RUN ln -s $GOPATH/src/github.com/sumeshkanayi/prometheus_ad_exporter/ad_exporter /usr/local/sbin
CMD ad_exporter  -server ${AD_SERVER} --user ${USER_NAME}

