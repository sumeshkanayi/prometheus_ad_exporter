default:    run
setup:
	git pull
	echo "Starting copying the source ${PWD} to ${GOPATH}"
	cp -r ${PWD} ${GOPATH}/src
	echo "copy completed"
build:	setup
	echo "pre installing packages , you will require internet access.Setup your proxy accoridngly"
  go get ./...
	go build -o ad_exporter main.go
	echo "pre requisite installation competed"
run:	build
	echo "Running exporter with basic configuration"
	./ad_exporter
