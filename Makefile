default: run
setup:
	echo "Starting copying the source ${PWD} to ${GOPATH}"
	go get github.com/sumeshkanayi/prometheus_ad_exporter
	echo "copy completed"
build:	setup
	echo "pre installing packages , you will require internet access.Setup your proxy accoridngly"
	cd  ${GOPATH}/src/github.com/sumeshkanayi/prometheus_ad_exporter;go get ./...;go build -o ad_exporter main.go
	echo "pre requisite installation competed"
run:	build
	echo "Running exporter with basic configuration"
	./ad_exporter
