#!/bin/bash

go build .

server_range=(250 500 750 1000 1250 1500 1750 2000 2250 2500 2750 3000)
client_range=(1)
repeat_range=(1)
packet_size_range=(1024)

echo "NUM_OF_SERVER, NUM_OF_CLIENT, NUM_REPEAT, PACKET_SIZE, TIME"
for PACKET_SIZE in "${packet_size_range[@]}"; do
  for NUM_OF_SERVER in "${server_range[@]}"; do
    for NUM_OF_CLIENT in "${client_range[@]}"; do
      for NUM_REPEAT in "${repeat_range[@]}"; do
        ./proxy-test $NUM_OF_SERVER $NUM_OF_CLIENT $NUM_REPEAT $PACKET_SIZE
      done
    done
  done
done

