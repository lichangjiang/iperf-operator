VERSION := v0.0.12

.PHONY: all

all: docker.push

install:
	go clean
	go build
	mkdir -p ./build
	mv iperf-operator ./build


build.img: install
	sudo docker build -t riverlcj/iperf-operator:$(VERSION) ./

docker.push: build.img
	sudo docker tag riverlcj/iperf-operator:$(VERSION) 172.17.8.101:30002/iperf/iperf-operator:$(VERSION)
	sudo docker push  172.17.8.101:30002/iperf/iperf-operator:$(VERSION)

clean:
	rm -rf ./build
	sudo docker rmi 192.168.38.13/iperf/iperf-operator:$(VERSION) 
	sudo docker rmi riverlcj/iperf-operator:$(VERSION) 

