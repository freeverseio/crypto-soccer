#!/bin/bash

# Ports used
ganache_port=8545
dapp_port=8080
dappchain_port_1=46657
dappchain_port_2=46658
build_number=416

# Check available platforms
platform='unknown'
unamestr=`uname`
if [[ "$unamestr" == 'Linux' ]]; then
  platform='linux'
elif [[ "$unamestr" == 'Darwin' ]]; then
  platform='osx'
else
  echo "Platform not supported on this script yet"
  exit -1
fi

function check_file_exists {
  if [ -f $1 ]; then
    echo 1
  else
    echo 0
  fi
}

function check_directory_exists {
  if [ -d $1 ]; then
    echo 1
  else
    echo 0
  fi
}

function check_port {
  if (nc -z localhost $1); then
    echo 1
  else
    echo 0
  fi
}

function is_setup_already {
  if [ $(check_directory_exists truffle-ethereum/node_modules) = 1 ] &&
     [ $(check_directory_exists truffle-dappchain/node_modules) = 1 ] &&
     [ $(check_file_exists dappchain/loom) = 1 ] &&
     [ $(check_directory_exists webclient/node_modules) = 1 ] &&
     [ $(check_directory_exists transfer-gateway-scripts/node_modules) = 1 ]; then
    echo 1
  else
    echo 0
  fi
}

# Setup function does the first work of download node_packages and loom binary
function setup {
  cd webclient
  echo "Install Webclient"
  yarn
  cd ../truffle-dappchain
  echo "Install DappChain"
  yarn
  cd ../truffle-ethereum
  echo "Install Truffle Ethereum"
  yarn
  cd ../transfer-gateway-scripts
  echo "Install Transfer Gateway Scripts"
  yarn
  cd ../dappchain
  echo "Install DappChain"
  wget "https://private.delegatecall.com/loom/$platform/build-$build_number/loom"
  chmod +x loom
  ./loom init -f
  sleep 5
  cp genesis.example.json genesis.json
  cd ..
}

function start_ganache {
  if [ $(check_file_exists truffle-ethereum/ganache.pid) = 0 ]; then
    echo "Start and Deploy Truffle Ethereum"
    echo "...we start Ganache via: yarn ganache-cli:dev, which calls: truffle-ethereum/scripts/ganache-cli.sh"
    echo "...and we send all ganache output to truffle-ethereum/ganache.log"
    echo "...Ganache PID was read and saved to file to be able to kill it later."
    cd truffle-ethereum
    yarn ganache-cli:dev > ganache.log
    echo "You should wait for 5 secs here..."
    #sleep 5
    echo "...DONE"
    cd ..
  else
    echo "Mmm... Ganache seems to be already up, please check..."
  fi
}

function deploy_truffle_ethereum {
    cd truffle-ethereum
    echo "...we deploy truffle-ethereum via: yarn deploy, which calls: rm -rf build; truffle migrate --reset --network development"
    echo "...note that the development network has the same parameters as our standard ganache network"
    echo "...we send the output to deploy.log"
    yarn deploy > deploy.log
    echo "...DONE"
    cd ..
}

function stop_ganache {
  if [ $(check_file_exists truffle-ethereum/ganache.pid) = 1 ]; then
    echo "Stop Truffle Ethereum by reading ganached PID saved to file, and doing a kill -9"
    cd truffle-ethereum
    pid=$(cat ganache.pid)
    kill -9 $pid
    rm ganache.pid
    cd ..
  else
    echo "Truffle Ethereum not running"
  fi
}

function start_loomchain {
  if [ $(check_file_exists dappchain/loom.pid) = 0 ]; then
    echo "Start DAppChain by executing cd dappchain; ./loom reset; ./loom run"
    cd dappchain
    ./loom reset; 
    echo "LOOM ENVIRONMENT INFO:"; 
    echo "------ ";
    ./loom env; 
    echo "------";
    ./loom run > /dev/null 2>&1 &
    loom_pid=$!
    echo $loom_pid > loom.pid
    echo "You should wait for 10 secs to complete safely..."
    #sleep 10
    cd ..
  else
    echo "DAppChain is running"
  fi
}

function stop_loomchain {
  if [ $(check_file_exists dappchain/loom.pid) = 1 ]; then
    echo "Stop DAppChain"
    cd dappchain
    pid=$(cat loom.pid)
    kill -9 $pid
    rm loom.pid
    cd ..
  else
    echo "DAppChain not running"
  fi
}

function deploy_truffle_dappchain {
  if [ $(check_file_exists webclient/webclient.pid) = 0 ]; then
    echo "Deploy Truffle DAppChain"
    cd truffle-dappchain
    yarn deploy > /dev/null 2>&1 &
    cd ..
    echo "You should wait for 20 secs to complete safely"
    #sleep 20
  else
    echo "Truffle DAppChain is deployed"
  fi
}

# Mapping is necessary to "mirroring" the token on mainnet and dappchain
function run_mapping {
  echo "Running mapping"
  cd transfer-gateway-scripts
  node mapping_crypto_cards.js > /dev/null 2>&1
  node mapping_game_token.js > /dev/null 2>&1
  cd ..
}

function start_webapp {
  if [ $(check_file_exists webclient/webclient.pid) = 0 ]; then
    echo "Running DApp"
    cd webclient
    yarn serve > /dev/null 2>&1 &
    echo "You should wait for 5 secs..."
    #sleep 5
    cd ..
  else
    echo "Dapp is running"
  fi
}

function stop_webdapp {
  if [ $(check_file_exists webclient/webclient.pid) = 1 ]; then
    echo "Stop DApp"
    cd webclient
    pid=$(cat webclient.pid)
    kill -9 $pid
    rm webclient.pid
    cd ..
  else
    echo "DApp not running"
  fi
}

case "$1" in
################################ SETUP ##############################
setup)
  echo "------------------------------------------------------------------------------------------"
  echo "Installing necessary packages, this can take up to 3 minutes (depending on internet speed)"
  echo "------------------------------------------------------------------------------------------"
  echo

  if [ $(is_setup_already) = 1 ]; then
    echo "Setup already ran"
    exit -1
  fi

  setup

  echo
  echo "-------------------------------------"
  echo "Done, packages installed with success"
  echo "-------------------------------------"

  ;;
################################ STATUS ##############################
status)
  echo "-----------------"
  echo "Services statuses"
  echo "-----------------"
  echo

  [[ $(check_file_exists truffle-ethereum/ganache.pid) = 1 ]] && echo "Ganache running" || echo "Ganache stopped"
  [[ $(check_file_exists dappchain/loom.pid) = 1 ]] && echo "Loomchain running" || echo "Loomchain stopped"
  [[ $(check_file_exists webclient/webclient.pid) = 1 ]] && echo "Webclient running" || echo "Webclient stopped"

  echo

  ;;
################################ START EVERYTHING ##############################
start)
  echo "-------------------------------------------------------------------"
  echo "Initializing background services, it can take (40 seconds) ..."
  echo "-------------------------------------------------------------------"
  echo

  if [ $(is_setup_already) = 0 ]; then
    echo "Please use the setup command first: ./transfer_gateway setup"
    echo
    exit -1
  fi

  if [ $(check_port $ganache_port) != 0 ]; then
    echo "Ganache port $ganache_port is already in use"
    echo
    exit -1
  fi

  if [ $(check_port $dapp_port) != 0 ]; then
    echo "Dapp port $dapp_port is already in use"
    echo
    exit -1
  fi

  if [ $(check_port $dappchain_port_1) != 0 ] || [ $(check_port $dappchain_port_2) != 0 ]; then
    echo "Some port from DAppChain already in use [$dappchain_port_1 or $dappchain_port_2]"
    echo
    exit -1
  fi

  start_ganache
  sleep 5;
  deploy_truffle_ethereum
  start_loomchain
  sleep 10;
  deploy_truffle_dappchain
  sleep 20;
  start_webapp
  sleep 5;
  run_mapping

  echo
  echo "-----------------------------------------------------------"
  echo "Services initialized and ready, check http://localhost:8080"
  echo "-----------------------------------------------------------"

  ;;
################################ STOP EVERYTHING ##############################
  stop)
  echo "-----------------"
  echo "Stopping services"
  echo "-----------------"
  echo

  stop_ganache
  stop_webdapp
  stop_loomchain
  stop_ganache

  echo
  echo "----------------"
  echo "Services stopped"
  echo "----------------"

  ;;
################################ START ONLY WEB-APP ##############################
start-webapp)
  echo "-------------------------------------------------------------------"
  echo "Starting webapp..."
  echo "-------------------------------------------------------------------"
  echo

  if [ $(is_setup_already) = 0 ]; then
    echo "Please use the setup command first: ./transfer_gateway setup"
    echo
    exit -1
  fi
  start_webapp

  echo
  echo "-----------------------------------------------------------"
  echo "Web app started, check http://localhost:8080"
  echo "-----------------------------------------------------------"

  ;;
################################ STOP ONLY WEBAPP ##############################
stop-webapp)
  echo "-------------------------------------------------------------------"
  echo "Stopping webapp..."
  echo "-------------------------------------------------------------------"
  echo
  stop_webapp

  echo
  echo "-----------------------------------------------------------"
  echo "Web app stopped"
  echo "-----------------------------------------------------------"

  ;;
################################ START LOOM CHAIN ##############################
start-loomchain)
  echo "-------------------------------------------------------------------"
  echo "Initializing Loomchain ..."
  echo "-------------------------------------------------------------------"
  echo

  if [ $(is_setup_already) = 0 ]; then
    echo "Please use the setup command first: ./transfer_gateway setup"
    echo
    exit -1
  fi

  if [ $(check_port $dapp_port) != 0 ]; then
    echo "Dapp port $dapp_port is already in use"
    echo
    exit -1
  fi

  if [ $(check_port $dappchain_port_1) != 0 ] || [ $(check_port $dappchain_port_2) != 0 ]; then
    echo "Some port from DAppChain already in use [$dappchain_port_1 or $dappchain_port_2]"
    echo
    exit -1
  fi

  start_loomchain

  echo
  echo "-----------------------------------------------------------"
  echo "Services initialized and ready, check http://localhost:8080"
  echo "-----------------------------------------------------------"

  ;;
################################ STOP ONLY LOOM CHAIN ##############################
start-loomchain)
  echo "-------------------------------------------------------------------"
  echo "Stopping Loomchain ..."
  echo "-------------------------------------------------------------------"
  echo

  if [ $(is_setup_already) = 0 ]; then
    echo "Please use the setup command first: ./transfer_gateway setup"
    echo
    exit -1
  fi

  stop_loomchain

  echo
  echo "-----------------------------------------------------------"
  echo "Services initialized and ready, check http://localhost:8080"
  echo "-----------------------------------------------------------"

  ;;
################################ DEPLOY DAPP CONTRACT ##############################
deploy-dapp)
  echo "-------------------------------------------------------------------"
  echo "Deploying DappChain ..."
  echo "-------------------------------------------------------------------"
  echo

  if [ $(is_setup_already) = 0 ]; then
    echo "Please use the setup command first: ./transfer_gateway setup"
    echo
    exit -1
  fi

  if [ $(check_port $dapp_port) != 0 ]; then
    echo "Dapp port $dapp_port is already in use"
    echo
    exit -1
  fi

  if [ $(check_port $dappchain_port_1) != 0 ] || [ $(check_port $dappchain_port_2) != 0 ]; then
    echo "Some port from DAppChain already in use [$dappchain_port_1 or $dappchain_port_2]"
    echo
    exit -1
  fi

  deploy_truffle_dappchain

  echo
  echo "-----------------------------------------------------------"
  echo "Services initialized and ready, check http://localhost:8080"
  echo "-----------------------------------------------------------"

  ;;
################################ START ONLY GANACHE ##############################
  start-ganache)
  echo "-------------------------------------------------------------------"
  echo "Waking up ganache..."
  echo "-------------------------------------------------------------------"
  echo

  if [ $(is_setup_already) = 0 ]; then
    echo "Please use the setup command first: ./transfer_gateway setup"
    echo
    exit -1
  fi

  if [ $(check_port $ganache_port) != 0 ]; then
    echo "Ganache port $ganache_port is already in use"
    echo
    exit -1
  fi

  start_ganache

  echo
  echo "-----------------------------------------------------------"
  echo "Ganache should be up and running"
  echo "-----------------------------------------------------------"

  ;;
################################ STOP ONLY GANACHE ##############################
  stop-ganache)
  echo "-----------------"
  echo "Stopping Ganache"
  echo "-----------------"
  echo

  stop_ganache

  echo
  echo "----------------"
  echo "Services stopped"
  echo "----------------"

  ;;
################################ DEPLOY ETHEREUM CONTRACTS ##############################
  deploy-ethereum)
  echo "-------------------------------------------------------------------"
  echo "Deploying Ethereum contracts..."
  echo "-------------------------------------------------------------------"
  echo

  if [ $(is_setup_already) = 0 ]; then
    echo "Please use the setup command first: ./transfer_gateway setup"
    echo
    exit -1
  fi

  deploy_truffle_ethereum

  echo
  echo "-----------------------------------------------------------"
  echo "Ethreum contracts deployed"
  echo "-----------------------------------------------------------"

  ;;
################################ RUN MAPPING ONLY ##############################
  mapping)
  echo "-------------------------------------------------------------------"
  echo "Running mapping..."
  echo "-------------------------------------------------------------------"
  echo

  if [ $(is_setup_already) = 0 ]; then
    echo "Please use the setup command first: ./transfer_gateway setup"
    echo
    exit -1
  fi

  run_mapping

  echo
  echo "-----------------------------------------------------------"
  echo "Running mapping...DONE"
  echo "-----------------------------------------------------------"

  ;;
################################ RESTART EVERYTHING ##############################
restart)
  $0 stop
  sleep 1
  $0 start

  ;;
################################ CLEAN IT ALL ##############################
cleanup)
  echo "-----------------------------------------"
  echo "Cleaning packages and binaries downloaded"
  echo "-----------------------------------------"
  echo

  echo "Cleaning DAppChain"
  rm -rf dappchain/loom
  rm -rf dappchain/genesis.json
  rm -rf dappchain/app.db
  rm -rf dappchain/chaindata
  rm -rf dappchain/loom.pid

  echo "Cleaning Transfer Gateway Scripts"
  rm -rf transfer-gateway-scripts/node_modules

  echo "Cleaning Truffle DAppChain"
  rm -rf truffle-dappchain/node_modules
  rm -rf truffle-dappchain/build

  echo "Cleaning Truffle Ethereum"
  rm -rf truffle-ethereum/node_modules
  rm -rf truffle-ethereum/build

  echo "Cleaning DApp"
  rm -rf webclient/node_modules

  echo
  echo "-----"
  echo "Clear"
  echo "-----"

  ;;
*)
   echo ""
   echo "Usage for all-in-one: $0 {setup|start|status|stop|restart|cleanup}"
   echo ""
   echo "Usage for individual: $0 {start-ganache|stop-ganache|start-loomchain|stop-loomchain|start-webapp|stop-webapp|deploy-ethereum|deploy-dapp|mapping}"
   echo "   ..logical order: start-ganache => deploy-ethereum => start-loomchain => deploy-dapp => start-webapp => mapping"
   echo ""
   echo "   ..typical use : "
   echo "     ..first time: $0 start-ganache; $0 start-ganache"
   echo "     ..many times: $0 deploy-ethereum or $0 deploy-dapp"
esac

exit 0