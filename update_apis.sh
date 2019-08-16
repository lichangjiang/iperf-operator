#!/bin/bash

export GO111MODULE=off
cd $GOPATH/src/k8s.io/code-generator 

./generate-groups.sh all \
    "github.com/lichangjiang/iperf-operator/pkg/client" \
    "github.com/lichangjiang/iperf-operator/pkg/apis" \
    iperf.test.svc:alpha1
