#!/usr/bin/env bash

ganache_port=8545

start_ganache() {
  ./node_modules/.bin/ganache-cli -p $ganache_port -m gravity top burden flip student usage spell purchase hundred improve check genre > /dev/null 2>&1 &
  ganache_pid=$!
  echo "ganache-cli started with pid $ganache_pid and private words: gravity top burden flip student usage spell purchase hundred improve check genre"
  echo $ganache_pid > ganache.pid
}

echo "Starting our own ganache instance at port $ganache_port. The resulting PID is saved to file to then be able to kill it by reading that file"
start_ganache
