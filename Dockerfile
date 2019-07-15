FROM centos:7
WORKDIR .

COPY ./build/iperf-operator /usr/local/bin
CMD [""]
