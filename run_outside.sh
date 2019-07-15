#!/bin/bash
make install
export IPERF_EMAIL_USER=li-changjiang@163.com
export IPERF_EMAIL_PWD=lcj89712
export IPERF_EMAIL_SMTP=smtp.163.com
export IPERF_EMAIL_PORT=465
export IPERF_OPERATOR_RUNENV=out

./build/iperf-operator operator 
