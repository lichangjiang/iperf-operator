#VERSION := v0.0.13
#haboraddr := 172.17.8.101:30002
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
	sudo docker tag riverlcj/iperf-operator:$(VERSION) $(HABORADDR)/$(HABORREPO):$(VERSION)
	sudo docker push  $(HABORADDR)/$(HABORREPO):$(VERSION)

clean:
	rm -rf ./build
	sudo docker rmi  $(HABORADDR)/$(HABORREPO):$(VERSION)
	sudo docker rmi riverlcj/iperf-operator:$(VERSION) 

