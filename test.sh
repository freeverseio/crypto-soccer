#!/bin/sh

MY_DIR=`dirname "$0"`
MY_DIR=`cd "$MY_DIR" ; pwd`

function exit_on_error
{
    if [ $? -ne 0 ]; then
        echo "An error occured. Exiting"
        exit -1
    fi
}

function help
{
    echo "Option list:"
    echo "  --clean      remove build directories, compile and test"
}

clean=0
# -----------------------------------------------------------------------
# parse command line arguments
# -----------------------------------------------------------------------
until [ $# -eq 0 ]
do
    arg=$1
    if [ $arg == '--clean' ]; then
        clean=1
        shift 1
    elif [ $arg == '-h' ]; then
         help
         exit 0
    elif [ $arg == '--h' ]; then
         help
         exit 0
    elif [ $arg == '--help' ]; then
         help
         exit 0
    elif [ $arg == '-help' ]; then
         help
         exit 0
    else
         echo "Unknown parameter: $1"
         help
         exit 1
    fi
 done

# -----------------------------------------------------------------------
# truffle
# -----------------------------------------------------------------------
function install_truffle
{
  cd ${MY_DIR}/truffle-core && npm install
  exit_on_error
}

function clean_truffle
{
  rm -r ${MY_DIR}/truffle-core/build
  exit_on_error
}
function build_truffle
{
  if [ ! -d "${MY_DIR}/truffle-core/build" ]; then
    cd ${MY_DIR}/truffle-core && node_modules/.bin/truffle compile
  fi
  exit_on_error
}
function test_truffle
{
  cd ${MY_DIR}/truffle-core && node_modules/.bin/truffle test
  exit_on_error
}

# -----------------------------------------------------------------------
# graphql
# -----------------------------------------------------------------------

function install_graph_ql
{
  cd ${MY_DIR}/nodejs-horizon && npm install
  exit_on_error
}
function test_graph_ql
{
  cd ${MY_DIR}/nodejs-horizon && npm test
  exit_on_error
}

# -----------------------------------------------------------------------
# go
# -----------------------------------------------------------------------
function clean_go
{
  cd ${MY_DIR}/go-soccer && go clean
  exit_on_error
}
function build_go
{
  cd ${MY_DIR}/go-soccer && go build
  exit_on_error
}
function test_go
{
  #TODO
  exit_on_error
}


# -----------------------------------------------------------------------

if [ $clean -eq 1 ]; then
  clean_truffle
  clean_go
fi

install_truffle
build_truffle
test_truffle

install_graph_ql
test_graph_ql

build_go
test_go
