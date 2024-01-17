#!/bin/bash -x
pid=$(ps -xu|grep devman |grep -v grep|awk '{print $2}')
test -z ${pid} || kill -9 ${pid}