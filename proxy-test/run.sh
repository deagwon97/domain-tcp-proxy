#!/bin/bash

go build .

server_range=(1)
client_range=(1)
repeat_range=(1)
data_size_range=(20)
mid_server_port=9980


function kill_mid_server(){
  kill -9 `netstat -nlp | grep $mid_server_port   | awk '{ print $7}' | awk -F"/" '{print $1}'`
}

cleanup(){
  kill_mid_server
}

function reverse_proxy_test(){
  echo "NUM_OF_SERVER, NUM_OF_CLIENT, NUM_OF_REPEAT, DATA_SIZE, TIME"
  for DATA_SIZE in "${data_size_range[@]}"; do
    for NUM_OF_SERVER in "${server_range[@]}"; do
      for NUM_OF_CLIENT in "${client_range[@]}"; do
        for NUM_OF_REPEAT in "${repeat_range[@]}"; do
          ./proxy-test $NUM_OF_SERVER $NUM_OF_CLIENT $NUM_OF_REPEAT $DATA_SIZE 2>/dev/null
          sleep 1
        done
      done
    done
  done
}


cd ../proxy-go
go build 
./proxy-go &
cd ../proxy-test
echo test-proxy-go
reverse_proxy_test
kill_mid_server


cd ../proxy-nodejs
yarn build
yarn serve & 
sleep 2

cd ../proxy-test
echo test-proxy-nodejs
reverse_proxy_test
kill_mid_server